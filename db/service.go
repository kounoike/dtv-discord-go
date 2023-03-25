package db

import (
	"context"
	"encoding/json"
)

func (s *Service) UnmarshalJSON(b []byte) error {
	type service Service
	var ss service

	err := json.Unmarshal(b, &ss)
	if err != nil {
		return err
	}
	*s = (Service)(ss)

	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	chmap := m["channel"].(map[string]interface{})

	s.Channel = chmap["channel"].(string)
	s.ChannelType = chmap["type"].(string)
	return nil
}

func (q *Queries) CreateOrUpdateService(ctx context.Context, service Service) error {
	params := createOrUpdateServiceParams{
		ID:                 service.ID,
		ServiceID:          service.ServiceID,
		NetworkID:          service.NetworkID,
		Type:               service.Type,
		LogoID:             service.LogoID,
		RemoteControlKeyID: service.RemoteControlKeyID,
		Name:               service.Name,
		Channel:            service.Channel,
		ChannelType:        service.ChannelType,
		HasLogoData:        service.HasLogoData,
	}
	return q.createOrUpdateService(ctx, params)
}
