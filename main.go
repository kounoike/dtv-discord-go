package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

var (
	version string
)

func main() {
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

	logger.Info("Starting dtv-discord-go", zap.String("version", version))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name))
	if err != nil {
		logger.Error("can't connect to db server", zap.Error(err))
		return
	}
	queries := sqlcdb.New(db)
	migrations := migrate.FileMigrationSource{Dir: "db/migrations"}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	ctx := context.Background()

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
		return
	}
	logger.Info("Applied migrations", zap.Int("count", n))

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
		return
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
		return
	}

	discordClient, err := discord_client.NewDiscordClient(config, queries, logger)
	if err != nil {
		logger.Error("can't connect to discord", zap.Error(err))
		return
	}

	usecase, err := dtv.NewDTVUsecase(config, discordClient, mirakcClient, queries, logger)
	if err != nil {
		logger.Error("can't create DTVUsecase", zap.Error(err))
	}

	err = discordClient.Open()
	if err != nil {
		logger.Error("can't open discord session", zap.Error(err))
		return
	}
	logger.Info("Connected!")

	discordClient.UpdateChannelsCache()
	logger.Info("Running!", zap.String("dtv-discord-go version", version), zap.String("mirakc version", mirakcVersion.Current))
	logMessage := fmt.Sprintf("起動しました。\ndtv-discord-go version:%s\nmirakc version:%s\n", version, mirakcVersion.Current)
	discordClient.SendMessage(discord.InformationCategory, discord.LogChannel, logMessage)
	if mirakcVersion.Current != mirakcVersion.Latest {
		discordClient.SendMessage(discord.InformationCategory, discord.LogChannel, fmt.Sprintf("mirakcの新しいバージョン(%s)があります", mirakcVersion.Latest))
	}

	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session(), logger)

	err = usecase.InitializeServiceChannels(ctx)
	if err != nil {
		logger.Error("can't create program infomation channel", zap.Error(err))
		return
	}
	logger.Info("CreateChannels OK")

	discordHandler.AddReactionAddHandler()
	discordHandler.AddReactionRemoveHandler()
	// TODO: 自動検索フォーラムに新規スレッドがあったときのハンドラ

	sseHandler := mirakc_handler.NewSSEHandler(*usecase, config.Mirakc.Host, config.Mirakc.Port, logger)
	sseHandler.Subscribe()
	logger.Info("Subscribed!")
}
