package dtv

func (dtv *DTVUsecase) CreateChannels() error {
	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		// 無ければ作ってくれるし、キャッシュにも入る
		_, err := dtv.discord.GetChannelID("録画-番組情報", service.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
