package dtv

import (
	"text/template"

	"github.com/kounoike/dtv-discord-go/config"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
	"golang.org/x/exp/slog"
)

type DTVUsecase struct {
	discord         *discord_client.DiscordClient
	mirakc          *mirakc_client.MirakcClient
	queries         *db.Queries
	contentPathTmpl *template.Template
}

func NewDTVUsecase(cfg config.Config, discordClient *discord_client.DiscordClient, mirakcClient *mirakc_client.MirakcClient, queries *db.Queries) (*DTVUsecase, error) {
	slog.Debug("NewDTVUsecase", "cfg", cfg)
	tmpl, err := template.New("content-path").Parse(cfg.Recording.ContentPathTemplate)
	if err != nil {
		return nil, err
	}
	return &DTVUsecase{
		discord:         discordClient,
		mirakc:          mirakcClient,
		queries:         queries,
		contentPathTmpl: tmpl,
	}, nil
}
