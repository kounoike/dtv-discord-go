package dtv

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
)

func (dtv *DTVUsecase) OnRecordingEmojiAdd(ctx context.Context, reaction *discordgo.MessageReactionAdd) error {
	msg, err := dtv.discord.GetChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return err
	}
	for _, r := range msg.Reactions {
		if r.Emoji.Name == discord.RecordingReactionEmoji {
			if r.Count == 1 {
				// 録画しよう！
				programMessage, err := dtv.queries.GetProgramMessageByMessageID(ctx, reaction.MessageID)
				if err != nil {
					return err
				}
				err = dtv.mirakc.AddRecordingSchedule(programMessage.ProgramID)
				if err != nil {
					return err
				}
				err = dtv.discord.MessageReactionAdd(reaction.ChannelID, reaction.MessageID, discord.OkReactionEmoji)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}
