package dtv

import (
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
	content, err := template.GetRecordingStoppedMessage(program, service, contentPath)
	if err != nil {
		return err
	}

	_, err = dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingChannel, content)
	if err != nil {
		return err
	}
	programMessage, err := dtv.queries.GetProgramMessageByProgramID(ctx, program.ID)
	if err != nil {
		return err
	}
	err = dtv.discord.MessageReactionAdd(programMessage.ChannelID, programMessage.MessageID, discord.RecordedReactionEmoji)
	if err != nil {
		return err
	}

	if dtv.asynq != nil {
		// NOTE: encoding.enabled = true or transcription.enabled = trueのとき
		pathData := template.PathTemplateData{}

		_ = dtv.gpt.ParseTitle(ctx, program.Name, &pathData)

		pathData.Program = template.PathProgram{
			Name:      program.Name,
			StartTime: program.StartTime(),
		}
		pathData.Service = template.PathService{
			Name: service.Name,
		}
		if dtv.encodingEnabled {
			outputPath, err := dtv.getEncodingOutputPath(ctx, program, service, &pathData)
			if err != nil {
				return err
			}
			task, err := tasks.NewProgramEncodeTask(programId, contentPath, outputPath)
			if err != nil {
				// NOTE: 多分JSONMarshalの失敗なので無視する
				dtv.logger.Warn("NewProgramEncodeTask failed", zap.Error(err))
				return nil
			}

			info, err := dtv.asynq.Enqueue(task)
			if err != nil {
				// NOTE: エンキュー失敗は無視する
				dtv.logger.Warn("task enqueue failed", zap.Error(err), zap.Int64("programId", programId), zap.String("contentPath", contentPath))
				return nil
			}
			dtv.logger.Debug("task enqueue success", zap.String("Type", info.Type))
		}
		if dtv.transcriptionEnabled {
			outputPath, err := dtv.getTranscriptionOutputPath(ctx, program, service, &pathData)
			if err != nil {
				return err
			}
			task, err := tasks.NewProgramTranscriptionTask(programId, contentPath, outputPath)
			if err != nil {
				// NOTE: 多分JSONMarshalの失敗なので無視する
				dtv.logger.Warn("NewProgramTranscriptionTask failed", zap.Error(err))
				return nil
			}

			info, err := dtv.asynq.Enqueue(task)
			if err != nil {
				// NOTE: エンキュー失敗は無視する
				dtv.logger.Warn("task enqueue failed", zap.Error(err), zap.Int64("programId", programId), zap.String("contentPath", contentPath))
				return nil
			}
			dtv.logger.Debug("task enqueue success", zap.String("Type", info.Type))
		}
	}

	return nil
}
