package mirakc_model

type Channel struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

type Service struct {
	ID                 uint    `json:"id"`
	ServiceID          uint    `json:"serviceId"`
	NetworkID          uint    `json:"networkId"`
	Type               uint    `json:"type"`
	LogoID             int     `json:"logoId"`
	RemoteControlKeyID uint    `json:"remoteControlKeyId"`
	Name               string  `json:"name"`
	Channel            Channel `json:"channel"`
	HasLogoData        bool    `json:"hasLogoData"`
}
