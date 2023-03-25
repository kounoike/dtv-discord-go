package dtv

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"golang.org/x/exp/slog"
)

func (dtv *DTVUsecase) OnRecordingEmojiRemove(ctx context.Context, reaction *discordgo.MessageReactionRemove) error {
	users, err := dtv.discord.GetMessageReactions(reaction.ChannelID, reaction.MessageID, discord.RecordingReactionEmoji)
	if err != nil {
		return err
	}
	if len(users) == 0 {
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
		slog.Debug("録画取り消し完了", "MessageID", reaction.MessageID, "ProgramID", programMessage.ProgramID)
	} else {
		slog.Debug("他のユーザの予約が残っています", "MessageID", reaction.MessageID)
	}

	return nil
}
