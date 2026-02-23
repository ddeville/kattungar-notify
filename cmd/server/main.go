package main

import (
	"context"
	"log"
	"os"

	"github.com/ddeville/kattungar-notify/internal/apns"
	"github.com/ddeville/kattungar-notify/internal/gcal"
	"github.com/ddeville/kattungar-notify/internal/server"
	"github.com/ddeville/kattungar-notify/internal/store"
)

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Starting server...")

	store, err := store.NewStore(C.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	apns, err := apns.NewApnsClient(apns.ApnsConfig{
		TeamID:  C.AppleTeamID,
		AppID:   C.AppleAppID,
		KeyID:   C.ApnsKeyID,
		KeyPath: C.ApnsKeyPath,
	})
	if err != nil {
		log.Fatal(err)
	}

	gcal, err := gcal.NewClient(gcal.CalendarConfig{
		GoogleCredentialsPath: C.GoogleCredsPath,
		GoogleRefreshToken:    C.GoogleRefreshToken,
		CalendarID:            C.GoogleCalendarID,
		ApnsClient:            apns,
		Store:                 store,
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go gcal.Run(ctx)

	s, err := server.NewServer(server.ServerConfig{
		Port:        C.Port,
		APIKeysPath: C.ServerAPIKeysPath,
		Store:       store,
		ApnsClient:  apns,
	})
	if err != nil {
		log.Fatal(err)
	}

	s.Serve()
	cancel()
}
