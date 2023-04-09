package main

import (
	"context"
	"flag"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/subcommands"
	"github.com/kounoike/dtv-discord-go/cmd"
)

var (
	version = "0.0.0-develop"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.NewBotCommand(version), "")
	subcommands.Register(cmd.NewWorkerCommand(version), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
