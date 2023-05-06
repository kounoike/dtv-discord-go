package dtv

func (dtv *DTVUsecase) ListChannelNames() ([]string, error) {
	services, err := dtv.mirakc.ListServices()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, s := range services {
		names = append(names, s.Name)
	}
	return names, nil
}
