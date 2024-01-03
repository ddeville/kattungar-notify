package main

import (
	"log"

	"github.com/ddeville/kattungar_notify/server"
	"github.com/ddeville/kattungar_notify/store"
)

func main() {
	store, err := store.NewStore("/home/damien/Downloads/store.db")
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(store, 3000)
	s.Serve()
}
