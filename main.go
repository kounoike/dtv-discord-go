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
)

func main() {
	// programs, err := getPrograms()
	var config config.Config
	configor.Load(&config, "config.yml")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name))
	if err != nil {
		fmt.Println(err)
		return
	}
	queries := sqlcdb.New(db)
	migrations := migrate.FileMigrationSource{Dir: "db/migrations"}
	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Applied %d migrations\n", n)

	mirakcClient := mirakc_client.NewMirakcClient(config.Mirakc.Host, config.Mirakc.Port)
	discordClient, err := discord_client.NewDiscordClient(config, queries)
	if err != nil {
		fmt.Println(err)
		return
	}

	usecase := dtv.NewDTVUsecase(discordClient, mirakcClient, queries)

	err = discordClient.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected!")
	discordHandler := discord_handler.NewDiscordHandler(usecase, discordClient.Session())

	err = usecase.CreateChannels()
	if err != nil {
		fmt.Println(err)
		return
	}

	discordHandler.AddReactionAddHandler()
	discordHandler.AddReactionRemoveHandler()

	sseHandler := mirakc_handler.NewSSEHandler(*usecase, config.Mirakc.Host, config.Mirakc.Port)
	sseHandler.Subscribe()
	fmt.Println("Subscribed!")
}
