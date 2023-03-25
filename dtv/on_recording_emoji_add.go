package dtv

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
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

func (dtv *DTVUsecase) checkRecordingForMessage(ctx context.Context, channelID string, messageID string) error {
	users, err := dtv.discord.GetMessageReactions(channelID, messageID, discord.RecordingReactionEmoji)
	if err != nil {
		return err
	}
	if len(users) == 1 {
		// 録画しよう！
		programMessage, err := dtv.queries.GetProgramMessageByMessageID(ctx, messageID)
		if errors.Cause(err) == sql.ErrNoRows {
			// NOTE: 番組情報以外の発言の場合は無視する
			return nil
		}
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
		slog.Debug("録画予約 OK", "ProgramID", programMessage.ProgramID, "contentPath", contentPath)
		err = dtv.discord.MessageReactionAdd(channelID, messageID, discord.OkReactionEmoji)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dtv *DTVUsecase) OnRecordingEmojiAdd(ctx context.Context, reaction *discordgo.MessageReactionAdd) error {
	return dtv.checkRecordingForMessage(ctx, reaction.ChannelID, reaction.MessageID)
}

func toSafePath(s string) string {
	safePathReplacer := strings.NewReplacer(
		"/", "／",
		":", "：",
		"*", "＊",
		"\\", "￥",
		"?", "？",
		"\"", "‟",
		"<", "＜",
		">", "＞",
		"|", "｜",
	)
	return safePathReplacer.Replace(s)
}
