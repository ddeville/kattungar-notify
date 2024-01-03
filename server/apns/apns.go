package apns

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"

	"github.com/ddeville/kattungar_notify/store"
)

const teamId = "Q8B696Y8U4"
const keyId = "SZQY3SP3XB"
const appId = "com.ddeville.kattungar_notify"

type ApnsClient struct {
	inner *apns2.Client
}

func NewApnsClient(keyPath string) (*ApnsClient, error) {
	authKey, err := token.AuthKeyFromFile(keyPath)
	if err != nil {
		return nil, err
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   keyId,
		TeamID:  teamId,
	}

	client := apns2.NewTokenClient(token)
	return &ApnsClient{client}, nil
}

func (c *ApnsClient) Push(device store.Device, payload *payload.Payload) (*apns2.Response, error) {
	notification := &apns2.Notification{}
	notification.Topic = appId
	notification.DeviceToken = device.Token
	notification.Payload = payload
	return c.inner.Push(notification)
}
