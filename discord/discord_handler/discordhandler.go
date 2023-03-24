package discord_handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/dtv"
	"golang.org/x/exp/slog"
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
	slog.Debug("add reaction emoji", "emoji", reaction.Emoji.Name, "UserID", reaction.UserID, "ChannelID", reaction.ChannelID, "MessageID", reaction.MessageID)

	if reaction.Emoji.Name == discord.RecordingReactionEmoji {
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiAdd(ctx, reaction)
		if err != nil {
			slog.Error("onrecording emoji add error", err, "UserID", reaction.UserID, "ChannelID", reaction.ChannelID, "MessageID", reaction.MessageID)
		}
	}
}

func (h *DiscordHandler) reactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	slog.Debug("remove reaction emoji", "emoji", reaction.Emoji.Name, "UserID", reaction.UserID, "ChannelID", reaction.ChannelID, "MessageID", reaction.MessageID)

	if reaction.Emoji.Name == discord.RecordingReactionEmoji {
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiRemove(ctx, reaction)
		if err != nil {
			slog.Error("onrecoding emoji remove error", err, "UserID", reaction.UserID, "ChannelID", reaction.ChannelID, "MessageID", reaction.MessageID)
		}
	}
}

func (h *DiscordHandler) AddReactionAddHandler() {
	h.session.AddHandler(h.reactionAdd)
}

func (h *DiscordHandler) AddReactionRemoveHandler() {
	h.session.AddHandler(h.reactionRemove)
}
