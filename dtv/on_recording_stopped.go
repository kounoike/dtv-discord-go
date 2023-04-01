package dtv

import (
	"context"

	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/template"
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

	return nil
}
