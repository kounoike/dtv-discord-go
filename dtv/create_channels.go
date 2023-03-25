package dtv

import (
	"context"

	"golang.org/x/exp/slog"
)

func (dtv *DTVUsecase) InitializeServiceChannels(ctx context.Context) error {
	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}
	slog.Debug("ListServices OK", "len(services)", len(services))
	for _, service := range services {
		dtv.queries.CreateOrUpdateService(ctx, service)
	}

	for _, service := range services {
		// 無ければ作ってくれるし、キャッシュにも入る
		ch, err := dtv.discord.GetCachedChannel("録画-番組情報", service.Name)
		if err != nil {
			return err
		}
		slog.Debug("GetCachedChannel", "ch.ID", ch.ID, "ch.Name", ch.Name)
	}
	return nil
}
