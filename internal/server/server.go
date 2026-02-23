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

	"github.com/ddeville/kattungar-notify/internal/apns"
	"github.com/ddeville/kattungar-notify/internal/store"
	"github.com/ddeville/kattungar-notify/internal/types"
)

type ServerConfig struct {
	Port        int
	APIKeysPath string
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
	apiKeysData, err := os.Open(cfg.APIKeysPath)
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

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to Kattungar Notify!"))
	})

	// These are admin endpoints to create/list/delete devices and are behind api key auth
	r.Route("/admin/device", func(r chi.Router) {
		r.Use(APIKeyAuth(apiKeys))
		r.Get("/", s.listDevices)
		r.Post("/", s.createDevice)
		r.Delete("/", s.deleteDevice)
	})

	// These are endpoints only gated on the device key itself
	r.Route("/device", func(r chi.Router) {
		r.Use(DeviceAuth(cfg.Store))
		r.Get("/", s.getDevice)
		r.Get("/list_notifications", s.listNotifications)
		r.Put("/name", s.updateDeviceName)
		r.Put("/token", s.updateDeviceToken)
	})

	// Sending a notification doesn't require authentication
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) createDevice(w http.ResponseWriter, r *http.Request) {
	var device types.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Creating device %v", device)

	if device.Key == "" {
		http.Error(w, "missing device key", http.StatusBadRequest)
		return
	}

	d, err := s.store.CreateDevice(device.Key, device.Name, device.Token)
	if err != nil {
		if store.IsExistingDeviceError(err) {
			log.Printf("Failed to create device because it already exists %v", err)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (s *Server) deleteDevice(w http.ResponseWriter, r *http.Request) {
	var device types.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleting device %v", device)

	found, err := s.store.DeleteDevice(device.Key)
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

func (s *Server) getDevice(w http.ResponseWriter, r *http.Request) {
	device := r.Context().Value(DeviceAuthContextKey).(*types.Device)

	jsonData, err := json.Marshal(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	device := r.Context().Value(DeviceAuthContextKey).(*types.Device)

	log.Printf("Listing notifications for device %v", device)

	notifications, err := s.store.ListNotifications(device)
	if err != nil {
		log.Printf("Cannot query notification from store %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(notifications)
	if err != nil {
		log.Printf("Cannot marshal notification into json %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *Server) updateDeviceName(w http.ResponseWriter, r *http.Request) {
	device := r.Context().Value(DeviceAuthContextKey).(*types.Device)

	var update types.Device
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if update.Name == "" {
		http.Error(w, "missing device name in body", http.StatusBadRequest)
		return
	}

	log.Printf("Updating name in device %v", device)

	_, err = s.store.UpdateDeviceName(device.Key, update.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateDeviceToken(w http.ResponseWriter, r *http.Request) {
	device := r.Context().Value(DeviceAuthContextKey).(*types.Device)

	var update types.Device
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if update.Token == "" {
		http.Error(w, "missing device token in body", http.StatusBadRequest)
		return
	}

	log.Printf("Updating token in device %v", device)

	_, err = s.store.UpdateDeviceToken(device.Key, update.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) notify(w http.ResponseWriter, r *http.Request) {
	var notification types.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Printf("Cannot decode notification body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(notification.DeviceKey) == 0 && len(notification.DeviceName) == 0 {
		log.Printf("Missing device_key or device_name")
		http.Error(w, "missing device_key or device_name", http.StatusBadRequest)
		return
	}

	if len(notification.Title) == 0 && len(notification.Subtitle) == 0 && len(notification.Body) == 0 {
		log.Printf("Missing title, subtitle, or body")
		http.Error(w, "missing title, subtitle, or body", http.StatusBadRequest)
		return
	}

	var device *types.Device
	if len(notification.DeviceKey) > 0 {
		device, err = s.store.GetDevice(notification.DeviceKey)
	} else {
		device, err = s.store.GetDeviceByName(notification.DeviceName)
	}
	if err != nil {
		log.Printf("Cannot retrieve device: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if device == nil {
		log.Printf("Unknown device key/name")
		http.Error(w, "unknown device key/name", http.StatusBadRequest)
		return
	}

	notification.DeviceKey = device.Key
	notification.DeviceName = device.Name

	err = s.store.RecordNotification(notification)
	if err != nil {
		log.Printf("Cannot record notification: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.apns.Notify(device, notification)
	if err != nil {
		log.Printf("Cannot connect to APNs: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Request to APNs failed: %v (status code %v)", res.Reason, res.StatusCode)
		http.Error(w, res.Reason, res.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}
