package dtv

import (
	"bytes"
	"context"
	"database/sql"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) scheduleRecordForMessage(ctx context.Context, channelID string, messageID string) error {
	programMessage, err := dtv.queries.GetProgramMessageByMessageID(ctx, messageID)
	if errors.Cause(err) == sql.ErrNoRows {
		// NOTE: 番組情報以外の発言の場合は無視する
		return nil
	}
	if err != nil {
		return err
	}

	_, err = dtv.queries.GetProgramRecordingByProgramId(ctx, programMessage.ProgramID)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return err
	}
	if err == nil {
		// すでに録画がスケジュールされているので何もしない
		return nil
	}

	users, err := dtv.discord.GetMessageReactions(channelID, messageID, discord.RecordingReactionEmoji)
	if err != nil {
		return err
	}
	count := len(users)
	for _, u := range users {
		if u.ID == dtv.discord.Session().State.User.ID {
			count -= 1
		}
	}
	if count == 0 {
		// Bot以外🔴リアクションしていないので無視
		return nil
	}

	// 録画しよう！
	program, err := dtv.queries.GetProgram(ctx, programMessage.ProgramID)
	if err != nil {
		return err
	}
	service, err := dtv.queries.GetServiceByProgramID(ctx, programMessage.ProgramID)
	if err != nil {
		return err
	}
	data := PathTemplateData{
		Program: PathProgram{
			Name:      program.Name,
			StartTime: program.StartTime(),
		},
		Service: PathService{
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

	err = dtv.queries.InsertProgramRecording(ctx, db.InsertProgramRecordingParams{ProgramID: programMessage.ProgramID, ContentPath: contentPath})
	if err != nil {
		return err
	}

	dtv.logger.Debug("録画予約 OK", zap.Int64("ProgramID", programMessage.ProgramID), zap.String("contentPath", contentPath))
	err = dtv.discord.MessageReactionAdd(channelID, messageID, discord.OkReactionEmoji)
	if err != nil {
		return err
	}

	return nil
}

func (dtv *DTVUsecase) OnRecordingEmojiAdd(ctx context.Context, reaction *discordgo.MessageReactionAdd) error {
	return dtv.scheduleRecordForMessage(ctx, reaction.ChannelID, reaction.MessageID)
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
