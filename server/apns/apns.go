package apns

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"

	"github.com/ddeville/kattungar-notify/store"
)

type ApnsConfig struct {
	TeamId  string
	AppId   string
	KeyId   string
	KeyPath string
}

type ApnsClient struct {
	inner *apns2.Client
	cfg   ApnsConfig
}

func NewApnsClient(cfg ApnsConfig) (*ApnsClient, error) {
	authKey, err := token.AuthKeyFromFile(cfg.KeyPath)
	if err != nil {
		return nil, err
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   cfg.KeyId,
		TeamID:  cfg.TeamId,
	}

	client := apns2.NewTokenClient(token)
	return &ApnsClient{client, cfg}, nil
}

func (c *ApnsClient) Push(device store.Device, payload *payload.Payload) (*apns2.Response, error) {
	notification := &apns2.Notification{
		Topic:       c.cfg.AppId,
		DeviceToken: device.Token,
		Payload:     payload,
	}
	return c.inner.Push(notification)
}
