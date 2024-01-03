package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	store  *Store
	port   int
}

func NewServer(store *Store, port int) Server {
	s := Server{chi.NewRouter(), store, port}
	r := s.router

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to Kattungar Notify!"))
	})
	r.Get("/list_devices", s.listDevices)
	r.Post("/add_device", s.addDevice)

	return s
}

func (s *Server) Serve() {
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
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
