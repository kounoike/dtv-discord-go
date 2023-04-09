package dtv

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/tasks"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) onProgramEncoded(ctx context.Context, taskInfo *asynq.TaskInfo) error {
	_, err := dtv.queries.GetEncodeTaskByTaskID(ctx, taskInfo.ID)
	if errors.Cause(err) != sql.ErrNoRows {
		return err
	}
	var payload tasks.ProgramEncodePayload
	err = json.Unmarshal(taskInfo.Payload, &payload)
	if err != nil {
		dtv.logger.Warn("task payload json.Unmarshal error", zap.Error(err))
		return err
	}
	err = dtv.queries.InsertEncodeTask(ctx, db.InsertEncodeTaskParams{TaskID: taskInfo.ID, Status: "success"})
	if err != nil {
		dtv.logger.Warn("failed to InsertEncodeTask", zap.Error(err))
		return err
	}
	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingChannel, fmt.Sprintf("**エンコード完了** `%s`のエンコードが完了しました", payload.OutputPath))
	if err != nil {
		dtv.logger.Warn("failed to SendMessage", zap.Error(err))
		return err
	}
	programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, payload.ProgramId)
	if errors.Cause(err) == sql.ErrNoRows {
		dtv.logger.Warn("failed to GetProgramMessageByProgramID", zap.Error(err))
		return err
	}
	if err != nil {
		dtv.logger.Warn("failed to GetProgramMessageByProgramID", zap.Error(err))
		return err
	}

	err = dtv.discord.MessageReactionAdd(programMessage.ChannelID, programMessage.MessageID, discord.EncodedReactionEmoji)
	if err != nil {
		dtv.logger.Warn("failed to MessageReactionAdd", zap.Error(err))
		return err
	}
	return nil
}

func (dtv *DTVUsecase) onProgramTranscribed(ctx context.Context, taskInfo *asynq.TaskInfo) error {
	_, err := dtv.queries.GetTranscribeTaskByTaskID(ctx, taskInfo.ID)
	if errors.Cause(err) != sql.ErrNoRows {
		return err
	}
	var payload tasks.ProgramTranscriptionPayload
	err = json.Unmarshal(taskInfo.Payload, &payload)
	if err != nil {
		dtv.logger.Warn("task payload json.Unmarshal error", zap.Error(err))
		return err
	}
	err = dtv.queries.InsertTranscribeTask(ctx, db.InsertTranscribeTaskParams{TaskID: taskInfo.ID, Status: "success"})
	if err != nil {
		dtv.logger.Warn("failed to InsertEncodeTask", zap.Error(err))
		return err
	}
	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingChannel, fmt.Sprintf("**文字起こし完了** `%s`の文字起こしが完了しました", payload.OutputPath))
	if err != nil {
		dtv.logger.Warn("failed to SendMessage", zap.Error(err))
		return err
	}
	programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, payload.ProgramId)
	if errors.Cause(err) == sql.ErrNoRows {
		dtv.logger.Warn("failed to GetProgramMessageByProgramID", zap.Error(err))
		return err
	}
	if err != nil {
		dtv.logger.Warn("failed to GetProgramMessageByProgramID", zap.Error(err))
		return err
	}

	err = dtv.discord.MessageReactionAdd(programMessage.ChannelID, programMessage.MessageID, discord.TranscriptionReactionEmoji)
	if err != nil {
		dtv.logger.Warn("failed to MessageReactionAdd", zap.Error(err))
		return err
	}
	return nil
}

func (dtv *DTVUsecase) CheckCompletedTask(ctx context.Context) error {
	if dtv.inspector == nil {
		return nil
	}
	taskInfoList, err := dtv.inspector.ListCompletedTasks("default")
	if err != nil {
		return err
	}
	for _, taskInfo := range taskInfoList {
		switch taskInfo.Type {
		case tasks.TypeHello:
			continue
		case tasks.TypeProgramEncode:
			_ = dtv.onProgramEncoded(ctx, taskInfo)
		case tasks.TypeProgramTranscription:
			_ = dtv.onProgramTranscribed(ctx, taskInfo)
		}
	}
	return nil
}
