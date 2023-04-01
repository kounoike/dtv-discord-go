package cmd

import (
	"context"
	"flag"
	"fmt"
	"text/template"

	"github.com/google/subcommands"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/configor"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/tasks"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type WorkerCommand struct {
	version string
}

func NewWorkerCommand(version string) *WorkerCommand {
	return &WorkerCommand{version: version}
}

func (c *WorkerCommand) Name() string { return "worker" }

func (c *WorkerCommand) Synopsis() string { return "worker subcommand" }

func (c *WorkerCommand) Usage() string { return "worker" }

func (c *WorkerCommand) SetFlags(f *flag.FlagSet) {
}

func (c *WorkerCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
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

	logger.Info("Starting dtv-discord-go worker", zap.String("version", c.version))

	redisAddr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"default": 5,
			},
		},
	)

	mux := asynq.NewServeMux()
	tmpl := template.Must(template.New("output-name-tmpl").Parse(config.Encoding.EncodeCommandTemplate))
	mux.Handle(tasks.TypeProgramEncode, tasks.NewProgramEncoder(logger, tmpl, config.Recording.BasePath, config.Encoding.BasePath))

	logger.Debug("Starting worker server")
	if err := srv.Run(mux); err != nil {
		logger.Fatal("could not run server", zap.Error(err))
	}

	return subcommands.ExitSuccess
}
