package dtv

import (
	"context"
	"database/sql"
	"strings"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) scheduleRecord(ctx context.Context, programMessage db.ProgramMessage) error {
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

	if err := dtv.queries.InsertProgramRecording(ctx, db.InsertProgramRecordingParams{ProgramID: programMessage.ProgramID, ContentPath: contentPath}); err != nil {
		return err
	}

	if err = dtv.mirakc.AddRecordingSchedule(programMessage.ProgramID, contentPath); err != nil {
		return err
	}

	dtv.logger.Debug("録画予約 OK", zap.Int64("ProgramID", programMessage.ProgramID), zap.String("contentPath", contentPath))
	if err := dtv.discord.MessageReactionAdd(programMessage.ChannelID, programMessage.MessageID, discord.OkReactionEmoji); err != nil {
		return err
	}

	return nil
}

func (dtv *DTVUsecase) cancelRecordSchedule(ctx context.Context, programMessage db.ProgramMessage) error {
	if err := dtv.queries.DeleteProgramRecordingByProgramId(ctx, programMessage.ProgramID); err != nil {
		return err
	}

	if err := dtv.mirakc.DeleteRecordingSchedule(programMessage.ProgramID); err != nil {
		return err
	}

	dtv.logger.Debug("録画予約キャンセル OK", zap.Int64("ProgramID", programMessage.ProgramID))
	if err := dtv.discord.MessageReactionRemove(programMessage.ChannelID, programMessage.MessageID, discord.OkReactionEmoji); err != nil {
		return err
	}

	return nil
}

func (dtv *DTVUsecase) checkRecordScheduleForMessage(ctx context.Context, channelID string, messageID string) error {
	programMessage, err := dtv.queries.GetProgramMessageByMessageID(ctx, messageID)
	if errors.Cause(err) == sql.ErrNoRows {
		// NOTE: 番組情報以外の発言の場合は無視する
		return nil
	}
	if err != nil {
		return err
	}

	ngEmojis, err := dtv.discord.GetMessageReactions(channelID, messageID, discord.NgReactionEmoji)
	if err != nil {
		return err
	}
	if len(ngEmojis) > 0 {
		return dtv.cancelRecordSchedule(ctx, programMessage)
	}

	recordingEmojis, err := dtv.discord.GetMessageReactions(channelID, messageID, discord.RecordingReactionEmoji)
	if err != nil {
		return err
	}
	if len(recordingEmojis) > 0 {
		for _, recordingEmoji := range recordingEmojis {
			if recordingEmoji.ID == dtv.discord.Session().State.User.ID {
				// NOTE: 自分のリアクションなので無視
				continue
			}
			return dtv.scheduleRecord(ctx, programMessage)
		}
	}

	asEmojis, err := dtv.discord.GetMessageReactions(channelID, messageID, discord.AutoSearchReactionEmoji)
	if err != nil {
		return err
	}
	if len(asEmojis) > 0 {
		for _, asEmoji := range asEmojis {
			if asEmoji.ID == dtv.discord.Session().State.User.ID {
				return dtv.scheduleRecord(ctx, programMessage)
			}
		}
	}

	return dtv.cancelRecordSchedule(ctx, programMessage)
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
