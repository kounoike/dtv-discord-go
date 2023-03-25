package main

import (
	"context"
	"database/sql"
	"fmt"

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
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		logger.Error("db migration error", zap.Error(err))
		return
	}
	logger.Info("Applied migrations", zap.Int("count", n))

	mirakcClient := mirakc_client.NewMirakcClient(config.Mirakc.Host, config.Mirakc.Port, logger)
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
	discordClient.SendMessage(discord.InformationCategory, discord.LogChannel, fmt.Sprintf("起動しました。version:%s", version))

	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session(), logger)

	ctx := context.Background()
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
