package dtv

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) OnRecordingEmojiRemove(ctx context.Context, reaction *discordgo.MessageReactionRemove) error {
	users, err := dtv.discord.GetMessageReactions(reaction.ChannelID, reaction.MessageID, discord.RecordingReactionEmoji)
	if err != nil {
		return err
	}
	if len(users) == 1 && users[0].ID == dtv.discord.Session().State.User.ID {
		programMessage, err := dtv.queries.GetProgramMessageByMessageID(ctx, reaction.MessageID)
		if err != nil {
			return err
		}
		err = dtv.mirakc.DeleteRecordingSchedule(programMessage.ProgramID)
		if err != nil {
			return err
		}
		err = dtv.discord.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, discord.OkReactionEmoji)
		if err != nil {
			return err
		}
		err = dtv.queries.DeleteProgramRecordingByProgramId(ctx, programMessage.ProgramID)
		if err != nil {
			return err
		}
		dtv.logger.Debug("録画取り消し完了", zap.String("MessageID", reaction.MessageID), zap.Int64("ProgramID", programMessage.ProgramID))
	} else {
		dtv.logger.Debug("他のユーザの予約が残っています", zap.String("MessageID", reaction.MessageID))
	}

	return nil
}
