package dtv

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"golang.org/x/exp/slog"
)

type ContentPathProgram struct {
	Name      string
	StartTime time.Time
}
type ContentPathService struct {
	Name string
}

type ContentPathTemplateData struct {
	Program ContentPathProgram
	Service ContentPathService
}

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
				program, err := dtv.queries.GetProgram(ctx, programMessage.ProgramID)
				if err != nil {
					return err
				}
				service, err := dtv.queries.GetServiceByProgramID(ctx, programMessage.ProgramID)
				if err != nil {
					return err
				}
				data := ContentPathTemplateData{
					Program: ContentPathProgram{
						Name:      program.Name,
						StartTime: program.StartTime(),
					},
					Service: ContentPathService{
						Name: service.Name,
					},
				}
				var buffer bytes.Buffer
				err = dtv.contentPathTmpl.Execute(&buffer, data)
				if err != nil {
					return err
				}
				contentPath := toSafePath(buffer.String())

				err = dtv.mirakc.AddRecordingSchedule(programMessage.ProgramID, contentPath)
				if err != nil {
					return err
				}
				slog.Debug("AddRecordingSchedule OK", "ProgramID", programMessage.ProgramID, "contentPath", contentPath)
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

func toSafePath(s string) string {
	ret := strings.ReplaceAll(s, "/", "／")
	ret = strings.ReplaceAll(ret, ":", "：")
	ret = strings.ReplaceAll(ret, "*", "＊")
	ret = strings.ReplaceAll(ret, "\\", "￥")
	ret = strings.ReplaceAll(ret, "?", "？")
	ret = strings.ReplaceAll(ret, "\"", "‟")
	ret = strings.ReplaceAll(ret, "<", "＜")
	ret = strings.ReplaceAll(ret, ">", "＞")
	ret = strings.ReplaceAll(ret, "|", "｜")
	return ret
}
