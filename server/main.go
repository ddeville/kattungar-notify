package main

import (
	"log"

	"github.com/ddeville/kattungar_notify/apns"
	"github.com/ddeville/kattungar_notify/server"
	"github.com/ddeville/kattungar_notify/store"
)

func main() {
	store, err := store.NewStore("/home/damien/Downloads/store.db")
	if err != nil {
		log.Fatal(err)
	}

	apns, err := apns.NewApnsClient("/home/damien/Downloads/AuthKey_SZQY3SP3XB.p8")
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(3000, store, apns)
	s.Serve()
}
