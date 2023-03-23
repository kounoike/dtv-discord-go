package discord_handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/dtv"
)

type DiscordHandler struct {
	dtv     *dtv.DTVUsecase
	session *discordgo.Session
}

func NewDiscordHandler(dtv *dtv.DTVUsecase, session *discordgo.Session) *DiscordHandler {
	return &DiscordHandler{
		dtv:     dtv,
		session: session,
	}
}

func (h *DiscordHandler) reactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	fmt.Println("add", reaction.Emoji.Name, reaction.UserID, reaction.ChannelID)

	if reaction.Emoji.Name == discord.RecordingReactionEmoji {
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiAdd(ctx, reaction)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (h *DiscordHandler) reactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	fmt.Println("remove", reaction.Emoji.Name, reaction.UserID, reaction.ChannelID)

	if reaction.Emoji.Name == discord.RecordingReactionEmoji {
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiRemove(ctx, reaction)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (h *DiscordHandler) AddReactionAddHandler() {
	h.session.AddHandler(h.reactionAdd)
}

func (h *DiscordHandler) AddReactionRemoveHandler() {
	h.session.AddHandler(h.reactionRemove)
}
