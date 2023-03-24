package dtv

import "golang.org/x/exp/slog"

func (dtv *DTVUsecase) CreateChannels() error {
	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}
	slog.Debug("ListServices OK", "len(services)", len(services))

	err = dtv.discord.UpdateChannelsCache()
	if err != nil {
		return err
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
