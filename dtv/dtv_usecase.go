package dtv

import (
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"go.uber.org/zap"
	"golang.org/x/text/width"
)

type DTVUsecase struct {
	discord         *discord_client.DiscordClient
	mirakc          *mirakc_client.MirakcClient
	queries         *db.Queries
	logger          *zap.Logger
	contentPathTmpl *template.Template
	autoSearchForum *discordgo.Channel
}

func fold(str string) string {
	return width.Fold.String(str)
}

func NewDTVUsecase(cfg config.Config, discordClient *discord_client.DiscordClient, mirakcClient *mirakc_client.MirakcClient, queries *db.Queries, logger *zap.Logger) (*DTVUsecase, error) {
	funcMap := map[string]interface{}{
		"fold": fold,
	}
	tmpl, err := template.New("content-path").Funcs(funcMap).Parse(cfg.Recording.ContentPathTemplate)
	if err != nil {
		return nil, err
	}
	return &DTVUsecase{
		discord:         discordClient,
		mirakc:          mirakcClient,
		queries:         queries,
		logger:          logger,
		contentPathTmpl: tmpl,
	}, nil
}
