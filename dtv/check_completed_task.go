package dtv

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kounoike/dtv-discord-go/db"
	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/tasks"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) CheckCompletedTask(ctx context.Context) error {
	dtv.logger.Debug("Start CheckCompletedTask")
	if dtv.inspector == nil {
		return nil
	}
	taskInfoList, err := dtv.inspector.ListCompletedTasks("default")
	if err != nil {
		return err
	}
	for _, taskInfo := range taskInfoList {
		_, err := dtv.queries.GetEncodeTaskByTaskID(ctx, taskInfo.ID)
		if errors.Cause(err) != sql.ErrNoRows {
			continue
		}
		var payload tasks.ProgramEncodePayload
		err = json.Unmarshal(taskInfo.Payload, &payload)
		if err != nil {
			dtv.logger.Warn("task payload json.Unmarshal error", zap.Error(err))
			continue
		}
		err = dtv.queries.InsertEncodeTask(ctx, db.InsertEncodeTaskParams{TaskID: taskInfo.ID, Status: "success"})
		if err != nil {
			dtv.logger.Warn("failed to InsertEncodeTask", zap.Error(err))
			continue
		}
		_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingChannel, fmt.Sprintf("**エンコード完了** `%s`のエンコードが完了しました", payload.OutputPath))
		if err != nil {
			dtv.logger.Warn("failed to SendMessage", zap.Error(err))
			continue
		}
		programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, payload.ProgramId)
		if err != nil {
			dtv.logger.Warn("failed to GetProgramMessageByProgramID", zap.Error(err))
		}

		err = dtv.discord.MessageReactionAdd(programMessage.ChannelID, programMessage.MessageID, discord.EncodedReactionEmoji)
		if err != nil {
			dtv.logger.Warn("failed to MessageReactionAdd", zap.Error(err))
		}
	}
	return nil
}
