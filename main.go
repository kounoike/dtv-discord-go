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
	"github.com/kounoike/dtv-discord-go/sse_handler"
	"github.com/kounoike/dtv-discord-go/tv"
	"github.com/r3labs/sse/v2"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	// programs, err := getPrograms()
	var config config.Config
	configor.Load(&config, "config.yml")

	ctx := context.Background()
	sseClient := sse.NewClient(fmt.Sprintf("http://%s:%d/events", config.Mirakc.Host, config.Mirakc.Port))
	mirakcClient := tv.NewMirakcClient(config.Mirakc.Host, config.Mirakc.Port)
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

	discordClient, err := discord.NewDiscordClient(ctx, config, queries)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = discordClient.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected!")

	sseHandler := sse_handler.NewSSEHandler(ctx, *mirakcClient, *discordClient, *sseClient, *queries)
	sseHandler.Subscribe()
	fmt.Println("Subscribed!")
}
