package main

import (
	"database/sql"
	"fmt"

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

	usecase := dtv.NewDTVUsecase(discordClient, mirakcClient, queries)

	err = discordClient.Open()
	if err != nil {
		slog.Error("can't open discord session", err)
		return
	}
	slog.Info("Connected!")
	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session())

	err = usecase.CreateChannels()
	if err != nil {
		slog.Error("can't create program infomation channel", err)
		return
	}

	discordHandler.AddReactionAddHandler()
	discordHandler.AddReactionRemoveHandler()

	sseHandler := mirakc_handler.NewSSEHandler(*usecase, config.Mirakc.Host, config.Mirakc.Port)
	sseHandler.Subscribe()
	slog.Info("Subscribed!")
}
