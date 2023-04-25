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
	"github.com/kounoike/dtv-discord-go/meili"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

type DTVUsecase struct {
	asynq                *asynq.Client
	inspector            *asynq.Inspector
	discord              *discord_client.DiscordClient
	mirakc               *mirakc_client.MirakcClient
	scheduler            *gocron.Scheduler
	queries              *db.Queries
	logger               *zap.Logger
	recordingBasePath    string
	contentPathTmpl      *template.Template
	encodedExt           string
	transcribedExt       string
	autoSearchChannel    *discordgo.Channel
	gpt                  *gpt.GPTClient
	meili                *meili.MeiliSearchClient
	kanaMatch            bool
	fuzzyMatch           bool
	encodingEnabled      bool
	transcriptionEnabled bool
	transcriptionType    string
	deleteOriginalFile   bool
	defaultQueueName     string
	encodeQueueName      string
	transcribeQueueName  string
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
	meili *meili.MeiliSearchClient,
) (*DTVUsecase, error) {
	funcMap := map[string]interface{}{
		"fold": fold,
	}
	contentTmpl, err := template.New("content-path").Funcs(funcMap).Parse(cfg.Recording.ContentPathTemplate)
	if err != nil {
		return nil, err
	}
	return &DTVUsecase{
		asynq:                asynqClient,
		inspector:            inspector,
		discord:              discordClient,
		mirakc:               mirakcClient,
		scheduler:            scheduler,
		queries:              queries,
		logger:               logger,
		gpt:                  gpt,
		meili:                meili,
		contentPathTmpl:      contentTmpl,
		encodedExt:           cfg.Encoding.Ext,
		transcribedExt:       cfg.Transcription.Ext,
		kanaMatch:            kanaMatch,
		fuzzyMatch:           fuzzyMatch,
		encodingEnabled:      cfg.Encoding.Enabled,
		transcriptionEnabled: cfg.Transcription.Enabled,
		transcriptionType:    cfg.Transcription.Type,
		deleteOriginalFile:   cfg.Encoding.DeleteOriginalFile,
		defaultQueueName:     cfg.TaskQueue.DefaultQueueName,
		encodeQueueName:      cfg.TaskQueue.EncodeQueueName,
		transcribeQueueName:  cfg.TaskQueue.TranscribeQueueName,
		recordingBasePath:    cfg.Recording.BasePath,
	}, nil
}
