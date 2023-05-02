package dtv

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (dtv *DTVUsecase) OnRecordingEmojiRemove(ctx context.Context, reaction *discordgo.MessageReactionRemove) error {
	return dtv.checkRecordScheduleForMessage(ctx, reaction.ChannelID, reaction.MessageID)
}
