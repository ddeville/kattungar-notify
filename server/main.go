package main

import "log"

func main() {
	store, err := NewStore("/home/damien/Downloads/store.db")
	if err != nil {
		log.Fatal(err)
	}

	s := NewServer(store, 3000)
	s.Serve()
}
