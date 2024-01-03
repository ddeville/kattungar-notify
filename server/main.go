package main

import (
	"log"

	"github.com/ddeville/kattungar_notify/apns"
	"github.com/ddeville/kattungar_notify/server"
	"github.com/ddeville/kattungar_notify/store"
)

const teamId = "Q8B696Y8U4"
const appId = "com.ddeville.kattungar_notify"

func main() {
	store, err := store.NewStore("/home/damien/Downloads/store.db")
	if err != nil {
		log.Fatal(err)
	}

	apns, err := apns.NewApnsClient(apns.ApnsConfig{
		TeamId:  teamId,
		AppId:   appId,
		KeyId:   "SZQY3SP3XB",
		KeyPath: "/home/damien/Downloads/AuthKey_SZQY3SP3XB.p8",
	})
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(3000, store, apns)
	s.Serve()
}
