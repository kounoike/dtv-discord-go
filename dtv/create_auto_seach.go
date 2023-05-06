package dtv

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) CreateAutoSearch(userID string, name string, title string, channel string, genre string, kanaSearch bool, fuzzySearch bool, regexSearch bool, record bool) error {
	seachMethod := "部分一致検索"
	if kanaSearch {
		if fuzzySearch {
			seachMethod = "かなあいまい検索"
		} else {
			seachMethod = "かな検索"
		}
	} else {
		if fuzzySearch {
			seachMethod = "あいまい検索"
		}
	}
	if regexSearch {
		seachMethod = "正規表現検索"
	}
	recordStr := "しない"
	if record {
		recordStr = "する"
	}
	content := fmt.Sprintf(
		"タイトル=%s\nチャンネル=%s\nジャンル=%s\n検索方法=%s\n録画=%s\nby <@%s>\n",
		title,
		channel,
		genre,
		seachMethod,
		recordStr,
		userID,
	)

	thID, err := dtv.discord.CreateAutoSearchThread(name, content)
	if err != nil {
		return err
	}

	ctx := context.Background()

	result, err := dtv.queries.InsertAutoSearch(ctx, db.InsertAutoSearchParams{
		Name:        name,
		Title:       normalizeString(title, kanaSearch),
		Channel:     normalizeString(channel, kanaSearch),
		Genre:       normalizeString(genre, kanaSearch),
		KanaSearch:  kanaSearch,
		FuzzySearch: fuzzySearch,
		RegexSearch: regexSearch,
		Record:      record,
		ThreadID:    thID,
	})
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	as, err := dtv.queries.GetAutoSearch(ctx, int32(id))
	if err != nil {
		return err
	}
	dtv.logger.Debug("as", zap.Any("as", as))

	autoSearch := dtv.getAutoSearchFromDB(&as)

	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}

	nowTime := time.Now().Unix() * 1000

	for _, service := range services {
		if autoSearch.IsMatchService(service.Name) {
			programs, err := dtv.mirakc.ListPrograms(uint(service.ID))
			if err != nil {
				dtv.logger.Warn("ListPrograms error", zap.Error(err))
				continue
			}
			for _, program := range programs {
				if program.StartAt+int64(program.Duration) < nowTime {
					// NOTE: 終了済みの番組は無視
					continue
				}
				asp := NewAutoSearchProgram(program)
				if autoSearch.IsMatchProgram(asp) {
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
					if record {
						err = dtv.discord.MessageReactionAdd(msg.ChannelID, msg.ID, discord.AutoSearchReactionEmoji)
						if err != nil {
							dtv.logger.Warn("MessageReactionAdd error", zap.Error(err), zap.String("channelID", msg.ChannelID), zap.String("messageID", msg.ID), zap.String("emoji", discord.AutoSearchReactionEmoji))
							continue
						}
						contentPath, err := dtv.getContentPath(ctx, program, service)
						if err != nil {
							dtv.logger.Warn("getContentPath error", zap.Error(err))
							continue
						}
						if err := dtv.mirakc.AddRecordingSchedule(program.ID, contentPath); err != nil {
							dtv.logger.Warn("AddRecordingSchedule error", zap.Error(err))
							// 重複登録の場合もあるので、エラーは無視して継続
						}
					}
					if err := dtv.queries.InsertAutoSearchFoundMessage(ctx, db.InsertAutoSearchFoundMessageParams{MessageID: msg.ID, ThreadID: thID, ProgramID: program.ID}); err != nil {
						dtv.logger.Warn("InsertAutoSearchFoundMessage error", zap.Error(err))
						continue
					}
				}
			}
		}
	}

	return nil
}
