package dtv

import (
	"bytes"
	"context"

	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/tasks"
	"github.com/kounoike/dtv-discord-go/template"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) OnRecordingStopped(ctx context.Context, programId int64) error {
	program, err := dtv.queries.GetProgram(ctx, programId)
	if err != nil {
		return err
	}
	service, err := dtv.queries.GetServiceByProgramID(ctx, programId)
	if err != nil {
		return err
	}
	contentPath, err := dtv.mirakc.GetRecordingScheduleContentPath(programId)
	if err != nil {
		return err
	}
	msg, err := template.GetRecordingStoppedMessage(program, service, contentPath)
	if err != nil {
		return err
	}

	dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingChannel, msg)

	if dtv.asynq != nil {
		var b bytes.Buffer
		data := PathTemplateData{
			Program: PathProgram{
				Name:      program.Name,
				StartTime: program.StartTime(),
			},
			Service: PathService{
				Name: service.Name,
			},
		}
		err = dtv.outputPathTmpl.Execute(&b, data)
		if err != nil {
			return err
		}
		outputPath := b.String()

		task, err := tasks.NewProgramEncodeTask(programId, contentPath, outputPath)
		if err != nil {
			// NOTE: 多分JSONMarshalの失敗なので無視する
			dtv.logger.Error("NewProgramEncodeTask failed", zap.Error(err))
			return nil
		}

		info, err := dtv.asynq.Enqueue(task)
		if err != nil {
			// NOTE: エンキュー失敗は無視する
			dtv.logger.Error("task enqueue failed", zap.Error(err), zap.Int64("programId", programId), zap.String("contentPath", contentPath))
			return nil
		}
		dtv.logger.Debug("task enqueue success", zap.String("Type", info.Type))
	}

	return nil
}
