package dtv

import (
	"context"

	"go.uber.org/zap"
)

func (dtv *DTVUsecase) Reindex(ctx context.Context) error {
	if err := dtv.meili.DeleteProgramIndex(); err != nil {
		return err
	}
	programRows, err := dtv.queries.ListProgramWithMessageAndServiceName(ctx)
	if err != nil {
		return err
	}
	dtv.logger.Info("番組数", zap.Int("count", len(programRows)))
	if err := dtv.meili.UpdatePrograms(programRows, dtv.discord.Session().State.Guilds[0].ID); err != nil {
		return err
	}

	if err := dtv.meili.DeleteRecordedFileIndex(); err != nil {
		return err
	}
	rows, err := dtv.queries.ListRecordedFiles(ctx)
	if err != nil {
		return err
	}
	if err := dtv.meili.UpdateRecordedFiles(rows); err != nil {
		return err
	}
	return nil
}
