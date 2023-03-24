package dtv

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (dtv *DTVUsecase) OnRecordingEmojiRemove(ctx context.Context, reaction *discordgo.MessageReactionRemove) error {
	msg, err := dtv.discord.GetChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return err
	}
	for _, r := range msg.Reactions {
		fmt.Println(r.Emoji.Name, r.Count)
		// if r.Emoji.Name == recordingReactionEmoji {
		// 	break
		// }
	}
	return nil
}
