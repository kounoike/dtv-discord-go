package dtv

import (
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/hibiken/asynq"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/gpt"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

type DTVUsecase struct {
	asynq                       *asynq.Client
	inspector                   *asynq.Inspector
	discord                     *discord_client.DiscordClient
	mirakc                      *mirakc_client.MirakcClient
	scheduler                   *gocron.Scheduler
	queries                     *db.Queries
	logger                      *zap.Logger
	contentPathTmpl             *template.Template
	encodingOutputPathTmpl      *template.Template
	transcriptionOutputPathTmpl *template.Template
	autoSearchChannel           *discordgo.Channel
	gpt                         *gpt.GPTClient
	kanaMatch                   bool
	fuzzyMatch                  bool
	encodingEnabled             bool
	transcriptionEnabled        bool
}

func fold(str string) string {
	return width.Fold.String(str)
}

func NewDTVUsecase(
	cfg config.Config,
	asynqClient *asynq.Client,
	inspector *asynq.Inspector,
	discordClient *discord_client.DiscordClient,
	mirakcClient *mirakc_client.MirakcClient,
	scheduler *gocron.Scheduler,
	queries *db.Queries,
	logger *zap.Logger,
	kanaMatch bool,
	fuzzyMatch bool,
	gpt *gpt.GPTClient,
) (*DTVUsecase, error) {
	funcMap := map[string]interface{}{
		"fold": fold,
	}
	contentTmpl, err := template.New("content-path").Funcs(funcMap).Parse(cfg.Recording.ContentPathTemplate)
	if err != nil {
		return nil, err
	}
	encodingOutputTmpl, err := template.New("encoding-output-path").Funcs(funcMap).Parse(cfg.Encoding.OutputPathTemplate)
	if err != nil {
		return nil, err
	}
	transcriptionOutputTmpl, err := template.New("transcription-output-path").Funcs(funcMap).Parse(cfg.Transcription.OutputPathTemplate)
	if err != nil {
		return nil, err
	}
	return &DTVUsecase{
		asynq:                       asynqClient,
		inspector:                   inspector,
		discord:                     discordClient,
		mirakc:                      mirakcClient,
		scheduler:                   scheduler,
		queries:                     queries,
		logger:                      logger,
		gpt:                         gpt,
		contentPathTmpl:             contentTmpl,
		encodingOutputPathTmpl:      encodingOutputTmpl,
		transcriptionOutputPathTmpl: transcriptionOutputTmpl,
		kanaMatch:                   kanaMatch,
		fuzzyMatch:                  fuzzyMatch,
		encodingEnabled:             cfg.Encoding.Enabled,
		transcriptionEnabled:        cfg.Transcription.Enabled,
	}, nil
}
