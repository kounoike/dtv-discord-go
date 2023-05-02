package dtv

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (dtv *DTVUsecase) OnRecordingEmojiAdd(ctx context.Context, reaction *discordgo.MessageReactionAdd) error {
	if reaction.UserID == dtv.discord.Session().State.User.ID {
		// NOTE: 自分のリアクションなので無視
		return nil
	}
	return dtv.checkRecordScheduleForMessage(ctx, reaction.ChannelID, reaction.MessageID)
}
