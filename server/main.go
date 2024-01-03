package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	store *Store
}

func (s *Server) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.store.ListDevices()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	fmt.Fprintf(w, "%v", devices)
}

func (s *Server) addDevice(w http.ResponseWriter, r *http.Request) {
	var device Device

	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.AddDevice(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	store, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}

	s := Server{store}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to Kattungar Notify!"))
	})
	r.Get("/list_devices", s.listDevices)
	r.Post("/add_device", s.addDevice)

	http.ListenAndServe(":3000", r)
}
