package main

import (
	"context"
	"log"

	"github.com/ddeville/kattungar-notify/apns"
	"github.com/ddeville/kattungar-notify/gcal"
	"github.com/ddeville/kattungar-notify/server"
	"github.com/ddeville/kattungar-notify/store"
)

const serverPort = 5000

const teamId = "Q8B696Y8U4"
const appId = "com.ddeville.kattungar-notify"

func main() {
	log.Println("Starting server...")

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	store, err := store.NewStore(cfg.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	apns, err := apns.NewApnsClient(apns.ApnsConfig{
		TeamId:  teamId,
		AppId:   appId,
		KeyId:   cfg.ApnsKeyId,
		KeyPath: cfg.ApnsKeyPath,
	})
	if err != nil {
		log.Fatal(err)
	}

	gcal, err := gcal.NewClient(gcal.CalendarConfig{
		GoogleCredentialsPath: cfg.GoogleCredsPath,
		GoogleRefreshToken:    cfg.GoogleRefreshToken,
		CalendarId:            cfg.GoogleCalendarId,
		ApnsClient:            apns,
		Store:                 store,
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go gcal.Run(ctx)

	s := server.NewServer(serverPort, store, apns)
	s.Serve()

	cancel()
}
