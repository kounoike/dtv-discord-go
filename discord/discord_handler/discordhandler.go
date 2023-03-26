package discord_handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/dtv"
	"go.uber.org/zap"
)

type DiscordHandler struct {
	dtv     *dtv.DTVUsecase
	session *discordgo.Session
	logger  *zap.Logger
}

func NewDiscordHandler(dtv *dtv.DTVUsecase, session *discordgo.Session, logger *zap.Logger) *DiscordHandler {
	return &DiscordHandler{
		dtv:     dtv,
		session: session,
		logger:  logger,
	}
}

func (h *DiscordHandler) reactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	h.logger.Debug("add reaction emoji", zap.String("emoji", reaction.Emoji.Name), zap.String("UserID", reaction.UserID), zap.String("ChannelID", reaction.ChannelID), zap.String("MessageID", reaction.MessageID))

	switch reaction.Emoji.Name {
	case discord.RecordingReactionEmoji:
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiAdd(ctx, reaction)
		if err != nil {
			h.logger.Error("onrecording emoji add error", zap.Error(err), zap.String("UserID", reaction.UserID), zap.String("ChannelID", reaction.ChannelID), zap.String("MessageID", reaction.MessageID))
		}
	case discord.OkReactionEmoji:
		ctx := context.Background()
		err := h.dtv.OnOkEmojiAdd(ctx, reaction)
		if err != nil {
			h.logger.Error("OnOkEmojiAdd error", zap.Error(err), zap.String("UserID", reaction.UserID), zap.String("ChannelID", reaction.ChannelID), zap.String("MessageID", reaction.MessageID))
		}
	default:
		h.logger.Debug("no intent for this Emoji", zap.String("emojiName", reaction.Emoji.Name))
	}
}

func (h *DiscordHandler) reactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	h.logger.Debug("remove reaction emoji", zap.String("emoji", reaction.Emoji.Name), zap.String("UserID", reaction.UserID), zap.String("ChannelID", reaction.ChannelID), zap.String("MessageID", reaction.MessageID))

	if reaction.Emoji.Name == discord.RecordingReactionEmoji {
		ctx := context.Background()
		err := h.dtv.OnRecordingEmojiRemove(ctx, reaction)
		if err != nil {
			h.logger.Error("onrecoding emoji remove error", zap.Error(err), zap.String("UserID", reaction.UserID), zap.String("ChannelID", reaction.ChannelID), zap.String("MessageID", reaction.MessageID))
		}
	}
}

func (h *DiscordHandler) AddReactionAddHandler() {
	h.session.AddHandler(h.reactionAdd)
}

func (h *DiscordHandler) AddReactionRemoveHandler() {
	h.session.AddHandler(h.reactionRemove)
}
