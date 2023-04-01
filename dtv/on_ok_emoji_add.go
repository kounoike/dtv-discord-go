package dtv

import (
	"context"
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) OnOkEmojiAdd(ctx context.Context, reaction *discordgo.MessageReactionAdd) error {
	users, err := dtv.discord.GetMessageReactions(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name)
	if err != nil {
		return err
	}
	if len(users) != 1 {
		// NOTE: 最初のOKリアクション以外は無視
		return nil
	}
	th, err := dtv.discord.GetChannel(reaction.ChannelID)
	if err != nil {
		return err
	}
	if th.Type != discordgo.ChannelTypeGuildPublicThread {
		// NOTE: 自動検索チャンネルじゃないのでリターン
		return nil
	}
	asCh, err := dtv.discord.GetCachedChannel(discord.NotifyAndScheduleCategory, discord.AutoActionChannelName)
	if err != nil {
		return err
	}
	if th.ParentID != asCh.ID {
		// NOTE: チャンネルが違うのでリターン
		return nil
	}

	// 自動検索を登録済みのEPGに対して実行する
	threadMsg, err := dtv.discord.GetChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return err
	}
	autoSearch, err := dtv.getAutoSeachFromMessage(threadMsg)
	if err != nil {
		return err
	}

	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		if autoSearch.IsMatchService(service.Name, dtv.kanaMatch, dtv.fuzzyMatch) {
			programs, err := dtv.mirakc.ListPrograms(uint(service.ID))
			if err != nil {
				dtv.logger.Warn("ListPrograms error", zap.Error(err))
				continue
			}
			for _, program := range programs {
				asp := NewAutoSearchProgram(program, dtv.kanaMatch)
				if autoSearch.IsMatchProgram(asp, dtv.fuzzyMatch) {
					// NOTE: DBに入ってるか確認する
					_, err := dtv.queries.GetProgram(ctx, program.ID)
					if errors.Cause(err) == sql.ErrNoRows {
						// NOTE: DBに入ってないプログラムは後で検索が走るはずなのでそっちで通知
						continue
					}
					programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, program.ID)
					if err != nil {
						dtv.logger.Warn("GetProgramMessageByProgramID error", zap.Error(err))
						continue
					}
					ch, err := dtv.discord.GetCachedChannel(discord.ProgramInformationCategory, service.Name)
					if err != nil {
						dtv.logger.Warn("GetCachedChannel error", zap.Error(err))
						continue
					}
					msg, err := dtv.discord.GetChannelMessage(ch.ID, programMessage.MessageID)
					if err != nil {
						dtv.logger.Warn("GetChannelMessage error", zap.Error(err))
						continue
					}
					err = dtv.sendAutoSearchMatchMessage(ctx, msg, program, &service, autoSearch)
					if err != nil {
						dtv.logger.Warn("sendAutoSearchMatchMessage error", zap.Error(err))
						continue
					}
				}
			}
		}
	}

	return nil
}
