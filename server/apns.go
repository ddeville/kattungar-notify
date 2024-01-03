package main

import (
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

const teamId = "Q8B696Y8U4"
const keyId = "SZQY3SP3XB"

type Client struct {
	inner *apns2.Client
}

func NewClient(keyPath string) (*Client, error) {
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
	return &Client{client}, nil
}
