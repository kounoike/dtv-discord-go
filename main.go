package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/subcommands"
	"github.com/jinzhu/configor"
	"github.com/kounoike/dtv-discord-go/cmd"
	"github.com/kounoike/dtv-discord-go/config"
)

var (
	version = "0.0.0-develop"
)

func main() {
	var config config.Config
	if err := configor.Load(&config, "config.yml"); err != nil {
		fmt.Fprintf(os.Stderr, "config load Error: %v\n", err)
		return
	}

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.NewBotCommand(version), "")
	subcommands.Register(cmd.NewWorkerCommand(version, config.TaskQueue.DefaultQueueName), "worker")
	if config.TaskQueue.DefaultQueueName != config.TaskQueue.EncodeQueueName {
		subcommands.Register(cmd.NewWorkerCommand(version, config.TaskQueue.EncodeQueueName), "worker")
	}
	if config.TaskQueue.DefaultQueueName != config.TaskQueue.TranscribeQueueName && config.TaskQueue.EncodeQueueName != config.TaskQueue.TranscribeQueueName {
		subcommands.Register(cmd.NewWorkerCommand(version, config.TaskQueue.TranscribeQueueName), "worker")
	}

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
