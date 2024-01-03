package main

import "log"

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}

	s := NewServer(store, 3000)
	s.Serve()
}
