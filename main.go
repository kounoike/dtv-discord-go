package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/kounoike/dtv-discord-go/config"
	sqlcdb "github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/discord/discord_handler"
	"github.com/kounoike/dtv-discord-go/dtv"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_handler"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/exp/slog"
)

func main() {
	// programs, err := getPrograms()
	var config config.Config
	configor.Load(&config, "config.yml")

	var logLevel slog.Level
	switch config.Log.Level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
		slog.Error("unknown log level", "log.level", config.Log.Level)
	}

	slog.SetDefault(slog.New(slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}.NewTextHandler(os.Stderr)))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name))
	if err != nil {
		slog.Error("can't connect to db server", err)
		return
	}
	queries := sqlcdb.New(db)
	migrations := migrate.FileMigrationSource{Dir: "db/migrations"}
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		slog.Error("db migration error", err)
		return
	}
	slog.Info("Applied migrations", "count", n)

	mirakcClient := mirakc_client.NewMirakcClient(config.Mirakc.Host, config.Mirakc.Port)
	discordClient, err := discord_client.NewDiscordClient(config, queries)
	if err != nil {
		slog.Error("can't connect to discord", err)
		return
	}

	usecase, err := dtv.NewDTVUsecase(config, discordClient, mirakcClient, queries)
	if err != nil {
		slog.Error("can't create DTVUsecase", err)
	}

	err = discordClient.Open()
	if err != nil {
		slog.Error("can't open discord session", err)
		return
	}
	slog.Info("Connected!")
	slog.Debug("Debug!")
	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session())

	ctx := context.Background()
	err = usecase.InitializeServiceChannels(ctx)
	if err != nil {
		slog.Error("can't create program infomation channel", err)
		return
	}
	slog.Info("CreateChannels OK")

	discordHandler.AddReactionAddHandler()
	discordHandler.AddReactionRemoveHandler()

	sseHandler := mirakc_handler.NewSSEHandler(*usecase, config.Mirakc.Host, config.Mirakc.Port)
	sseHandler.Subscribe()
	slog.Info("Subscribed!")
}
