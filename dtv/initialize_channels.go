package dtv

import (
	"context"

	"github.com/kounoike/dtv-discord-go/discord"
	"go.uber.org/zap"
)

func (dtv *DTVUsecase) InitializeServiceChannels(ctx context.Context) error {
	forum, err := dtv.discord.CreateNotifyAndScheduleForum()
	if err != nil {
		dtv.logger.Error("can't create notify and schedule forum", zap.Error(err))
	}
	dtv.autoSearchForum = forum

	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}
	dtv.logger.Debug("ListServices OK", zap.Int("len(services)", len(services)))
	for _, service := range services {
		dtv.queries.CreateOrUpdateService(ctx, service)
	}

	for _, service := range services {
		// 無ければ作ってくれるし、キャッシュにも入る
		ch, err := dtv.discord.GetCachedChannel(discord.ProgramInformationCategory, service.Name)
		if err != nil {
			return err
		}
		dtv.logger.Debug("GetCachedChannel", zap.String("ch.ID", ch.ID), zap.String("ch.Name", ch.Name))
	}
	return nil
}
