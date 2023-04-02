package cmd

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/subcommands"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/configor"
	"github.com/kounoike/dtv-discord-go/config"
	sqlcdb "github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/discord/discord_handler"
	"github.com/kounoike/dtv-discord-go/dtv"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_handler"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_model"
	"github.com/lestrrat-go/backoff/v2"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BotCommand struct {
	version string
}

func NewBotCommand(version string) *BotCommand {
	return &BotCommand{version: version}
}

func (c *BotCommand) Name() string { return "bot" }

func (c *BotCommand) Synopsis() string { return "bot subcommand" }

func (c *BotCommand) Usage() string { return "bot" }

func (c *BotCommand) SetFlags(f *flag.FlagSet) {
}

func (c *BotCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// programs, err := getPrograms()
	var config config.Config
	configor.Load(&config, "config.yml")

	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	level := zap.NewAtomicLevel()
	switch config.Log.Level {
	case "DEBUG":
		level.SetLevel(zapcore.DebugLevel)
	case "INFO":
		level.SetLevel(zapcore.InfoLevel)
	case "WARN":
		level.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		level.SetLevel(zapcore.ErrorLevel)
	default:
		fmt.Println("unknown log level:", config.Log.Level)
		level.SetLevel(zapcore.WarnLevel)
	}
	logCfg.Level = level
	logger, err := logCfg.Build()
	if err != nil {
		fmt.Println("can't build logger")
	}
	defer logger.Sync()

	logger.Info("Starting dtv-discord-go bot", zap.String("version", c.version))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name))
	if err != nil {
		logger.Error("can't connect to db server", zap.Error(err))
		return subcommands.ExitFailure
	}
	queries := sqlcdb.New(db)
	migrations := migrate.FileMigrationSource{Dir: "db/migrations"}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	p1 := backoff.Exponential(
		backoff.WithMinInterval(time.Second),
		backoff.WithMaxInterval(time.Minute),
		backoff.WithJitterFactor(0.05),
	)
	retryMigrationFunc := func(db *sql.DB, migrations migrate.FileMigrationSource) (int, error) {
		b := p1.Start(ctx)
		for backoff.Continue(b) {
			n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
			if err == nil {
				return n, nil
			}
		}
		return 0, errors.New("failed to migration")
	}

	n, err := retryMigrationFunc(db, migrations)
	if err != nil {
		logger.Error("db migration error", zap.Error(err))
		return subcommands.ExitFailure
	}
	logger.Info("Applied migrations", zap.Int("count", n))

	var asynqClient *asynq.Client
	var asynqInspector *asynq.Inspector

	if config.Encoding.Enabled {
		redisAddr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
		asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
		defer asynqClient.Close()
		asynqInspector = asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
		defer asynqInspector.Close()
	} else {
		asynqClient = nil
		asynqInspector = nil
	}

	mirakcClient := mirakc_client.NewMirakcClient(config.Mirakc.Host, config.Mirakc.Port, logger)

	p2 := backoff.Constant(
		backoff.WithInterval(5*time.Second),
		backoff.WithJitterFactor(0.05),
		backoff.WithMaxRetries((6*60*60)/5),
	)
	retryMirakcServiceFunc := func() error {
		b := p2.Start(ctx)
		cnt := 0
		for backoff.Continue(b) {
			services, err := mirakcClient.ListServices()
			if err != nil {
				logger.Warn("mirakc ListServices error", zap.Error(err))
			}
			if err == nil && len(services) > 0 {
				return nil
			}
			if cnt%12 == 0 {
				logger.Info("waiting to mirakc service discovery...")
			}
			cnt += 1
		}
		return errors.New("failed to get services")
	}
	err = retryMirakcServiceFunc()
	if err != nil {
		logger.Error("service retrieve error", zap.Error(err))
		return subcommands.ExitFailure
	}
	logger.Info("Get list of services: OK")

	retryMirakcVersionFunc := func() (*mirakc_model.Version, error) {
		b := p2.Start(ctx)
		for backoff.Continue(b) {
			version, err := mirakcClient.GetVersion()
			if err == nil {
				return version, nil
			}
		}
		return nil, errors.New("failed to get services")
	}
	mirakcVersion, err := retryMirakcVersionFunc()
	if err != nil {
		logger.Error("service retrieve error", zap.Error(err))
		return subcommands.ExitFailure
	}

	discordClient, err := discord_client.NewDiscordClient(config, queries, logger)
	if err != nil {
		logger.Error("can't connect to discord", zap.Error(err))
		return subcommands.ExitFailure
	}

	// NOTE: 日本国内のみをターゲットにする
	scheduler := gocron.NewScheduler(time.FixedZone("JST", 9*60*60))
	// scheduler.SetMaxConcurrentJobs(10, gocron.RescheduleMode)

	usecase, err := dtv.NewDTVUsecase(config, asynqClient, asynqInspector, discordClient, mirakcClient, scheduler, queries, logger, config.Match.KanaMatch, config.Match.FuzzyMatch)
	if err != nil {
		logger.Error("can't create DTVUsecase", zap.Error(err))
	}

	err = discordClient.Open()
	if err != nil {
		logger.Error("can't open discord session", zap.Error(err))
		return subcommands.ExitFailure
	}
	logger.Info("Connected!")

	discordClient.UpdateChannelsCache()
	logger.Info("Running!", zap.String("dtv-discord-go version", c.version), zap.String("mirakc version", mirakcVersion.Current))
	logMessage := fmt.Sprintf("起動しました。\ndtv-discord-go version:%s\nmirakc version:%s\n", "v"+c.version, mirakcVersion.Current)
	discordClient.SendMessage(discord.InformationCategory, discord.LogChannel, logMessage)
	if mirakcVersion.Current != mirakcVersion.Latest {
		discordClient.SendMessage(discord.InformationCategory, discord.LogChannel, fmt.Sprintf("mirakcの新しいバージョン(%s)があります", mirakcVersion.Latest))
	}

	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session(), logger)

	err = usecase.InitializeServiceChannels(ctx)
	if err != nil {
		logger.Error("can't create program infomation channel", zap.Error(err))
		return subcommands.ExitFailure
	}
	logger.Info("CreateChannels OK")

	// エンコード結果取得タスク
	if config.Encoding.Enabled {
		scheduler.Every("1m").Do(func() {
			err := usecase.CheckCompletedTask(ctx)
			if err != nil {
				logger.Error("CheckCompletedTask error", zap.Error(err))
			}
			err = usecase.CheckFailedTask(ctx)
			if err != nil {
				logger.Error("CheckFailedTask error", zap.Error(err))
			}
		})
	}

	// バージョンチェックするタスク
	// 適当に12:30に動かしてみる
	scheduler.Every(1).Day().At("12:30").Do(func() {
		err := usecase.CheckUpdateTask(ctx, c.version)
		if err != nil {
			logger.Error("CheckUpdateTask error", zap.Error(err))
		}
	})

	scheduler.StartAsync()

	discordHandler.AddReactionAddHandler()
	discordHandler.AddReactionRemoveHandler()

	logger.Info("AddDiscordHandle done. start subscribe to SSE events.")

	sseHandler := mirakc_handler.NewSSEHandler(*usecase, config.Mirakc.Host, config.Mirakc.Port, logger)
	sseHandler.Subscribe()
	logger.Info("Subscribed!")

	return subcommands.ExitSuccess
}
