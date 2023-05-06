package dtv

import (
	"context"

	"github.com/kounoike/dtv-discord-go/discord"
)

func (dtv *DTVUsecase) DeleteAutoSearch(ctx context.Context, threadID string) error {
	autoSearch, err := dtv.queries.GetAutoSearchByThreadID(ctx, threadID)
	if err != nil {
		return err
	}

	if autoSearch.Record {
		founds, err := dtv.queries.ListAutoSearchFoundMessages(ctx, threadID)
		if err != nil {
			return err
		}
		for _, found := range founds {
			programID := found.ProgramID
			cnt, err := dtv.queries.CountAutoSearchFoundMessagesByProgramID(ctx, programID)
			if err != nil {
				return err
			}
			if cnt > 1 {
				// 他の自動予約がある
				continue
			}
			programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, programID)
			if err != nil {
				return err
			}
			users, err := dtv.discord.GetMessageReactions(programMessage.ChannelID, programMessage.MessageID, discord.RecordingReactionEmoji)
			if err != nil {
				return err
			}
			foundOther := false
			for _, u := range users {
				if u.ID != dtv.discord.Session().State.User.ID {
					foundOther = true
					break
				}
			}
			if foundOther {
				// 誰かが個別録画予約しているので、自動検索絵文字だけ外しておく
				if err := dtv.discord.MessageReactionRemove(programMessage.ChannelID, programMessage.MessageID, discord.AutoSearchReactionEmoji); err != nil {
					return err
				}
				continue
			}

			// 他に同じ番組が録画予約されていない場合
			if err := dtv.mirakc.DeleteRecordingSchedule(programID); err != nil {
				return err
			}
			if err := dtv.discord.MessageReactionRemove(programMessage.ChannelID, programMessage.MessageID, discord.AutoSearchReactionEmoji); err != nil {
				return err
			}
			if err := dtv.discord.MessageReactionRemove(programMessage.ChannelID, programMessage.MessageID, discord.OkReactionEmoji); err != nil {
				return err
			}

		}
	}

	if err := dtv.queries.DeleteAutoSearchFoundMessagesByThreadID(ctx, threadID); err != nil {
		return err
	}

	if err := dtv.queries.DeleteAutoSearch(ctx, autoSearch.ID); err != nil {
		return err
	}

	return dtv.discord.DeleteThread(threadID)
}
