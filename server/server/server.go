package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ddeville/kattungar-notify/apns"
	"github.com/ddeville/kattungar-notify/store"
)

type ServerConfig struct {
	Port        int
	ApiKeysPath string
	Store       *store.Store
	ApnsClient  *apns.ApnsClient
}

type Server struct {
	port   int
	router *chi.Mux
	store  *store.Store
	apns   *apns.ApnsClient
}

func NewServer(cfg ServerConfig) (*Server, error) {
	apiKeysData, err := os.Open(cfg.ApiKeysPath)
	if err != nil {
		return nil, err
	}
	defer apiKeysData.Close()

	var apiKeys []string
	err = json.NewDecoder(apiKeysData).Decode(&apiKeys)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	s := Server{cfg.Port, r, cfg.Store, cfg.ApnsClient}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(ApiKeyAuth(apiKeys))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to Kattungar Notify!"))
	})

	r.Route("/devices", func(r chi.Router) {
		r.Get("/", s.listDevices)
		r.Post("/", s.createDevice)
		r.Put("/", s.updateDevice)
		r.Delete("/", s.deleteDevice)
	})

	r.Route("/notify", func(r chi.Router) {
		r.Post("/", s.notify)
	})

	return &s, nil
}

func (s *Server) Serve() {
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}

func (s *Server) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.store.ListDevices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(devices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) createDevice(w http.ResponseWriter, r *http.Request) {
	var device store.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Creating device %v", device)

	if device.Id != 0 {
		http.Error(w, "cannot pass device ID", http.StatusBadRequest)
		return
	}

	d, err := s.store.CreateDevice(device)
	if err != nil {
		if store.IsExistingDeviceError(err) {
			http.Error(w, "A device with this name already exists", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) updateDevice(w http.ResponseWriter, r *http.Request) {
	var device store.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updating device %v", device)

	if device.Id == 0 {
		http.Error(w, "missing device ID", http.StatusBadRequest)
		return
	}

	found, err := s.store.UpdateDevice(device)
	if err != nil {
		if store.IsExistingDeviceError(err) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if !found {
		http.Error(w, "cannot find device", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteDevice(w http.ResponseWriter, r *http.Request) {
	var device store.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleting device %v", device)

	found, err := s.store.DeleteDevice(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "cannot find device", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) notify(w http.ResponseWriter, r *http.Request) {
	var notification store.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(notification.DeviceName) == 0 {
		http.Error(w, "missing device name", http.StatusBadRequest)
		return
	}

	if len(notification.Title) == 0 && len(notification.Subtitle) == 0 {
		http.Error(w, "missing title or subtitle", http.StatusBadRequest)
		return
	}

	device, err := s.store.GetDevice(notification.DeviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if device == nil {
		http.Error(w, "unknown device name", http.StatusBadRequest)
		return
	}

	err = s.store.RecordNotification(notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.apns.Notify(device, notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		http.Error(w, res.Reason, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
