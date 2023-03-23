package dtv

import (
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord/discord_client"
	"github.com/kounoike/dtv-discord-go/mirakc/mirakc_client"
)

type DTVUsecase struct {
	discord *discord_client.DiscordClient
	mirakc  *mirakc_client.MirakcClient
	queries *db.Queries
}

func NewDTVUsecase(discordClient *discord_client.DiscordClient, mirakcClient *mirakc_client.MirakcClient, queries *db.Queries) *DTVUsecase {
	return &DTVUsecase{
		discord: discordClient,
		mirakc:  mirakcClient,
		queries: queries,
	}
}
