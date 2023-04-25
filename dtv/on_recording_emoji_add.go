package dtv

import (
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

	// 録画しよう！
	program, err := dtv.queries.GetProgram(ctx, programMessage.ProgramID)
	if err != nil {
		return err
	}
	service, err := dtv.queries.GetServiceByProgramID(ctx, programMessage.ProgramID)
	if err != nil {
		return err
	}

	contentPath, err := dtv.getContentPath(ctx, program, service)
	if err != nil {
		return err
	}

	err = dtv.queries.InsertProgramRecording(ctx, db.InsertProgramRecordingParams{ProgramID: programMessage.ProgramID, ContentPath: contentPath})
	if err != nil {
		return err
	}

	err = dtv.mirakc.AddRecordingSchedule(programMessage.ProgramID, contentPath)
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
	if reaction.UserID == dtv.discord.Session().State.User.ID {
		// NOTE: 自分のリアクションなので無視
		return nil
	}
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
