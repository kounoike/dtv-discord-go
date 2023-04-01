package dtv

import (
	"context"

	"github.com/kounoike/dtv-discord-go/discord"
	"github.com/kounoike/dtv-discord-go/template"
)

func (dtv *DTVUsecase) OnRecordingFailed(ctx context.Context, programId int64, reason string) error {
	program, err := dtv.queries.GetProgram(ctx, programId)
	if err != nil {
		return err
	}
	service, err := dtv.queries.GetServiceByProgramID(ctx, programId)
	if err != nil {
		return err
	}
	recording, err := dtv.queries.GetProgramRecordingByProgramId(ctx, programId)
	if err != nil {
		return err
	}
	msg, err := template.GetRecordingFailedMessage(program, service, recording.ContentPath, reason)
	if err != nil {
		return err
	}

	dtv.discord.SendMessage(discord.InformationCategory, discord.RecordingFailedChannel, msg)

	return nil
}
