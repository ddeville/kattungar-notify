package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ddeville/kattungar-notify/internal/types"
)

const alertmanagerDeviceNameKey = "kattungar_notify_device"

type alertmanagerPayload struct {
	Status            string              `json:"status"`
	Receiver          string              `json:"receiver"`
	GroupLabels       map[string]string   `json:"groupLabels"`
	CommonLabels      map[string]string   `json:"commonLabels"`
	CommonAnnotations map[string]string   `json:"commonAnnotations"`
	ExternalURL       string              `json:"externalURL"`
	Alerts            []alertmanagerAlert `json:"alerts"`
}

type alertmanagerAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	GeneratorURL string            `json:"generatorURL"`
}

func (s *Server) alertmanagerWebhook(w http.ResponseWriter, r *http.Request) {
	var payload alertmanagerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deviceName := payload.CommonLabels[alertmanagerDeviceNameKey]
	if deviceName == "" {
		for _, alert := range payload.Alerts {
			deviceName = alert.Labels[alertmanagerDeviceNameKey]
			if deviceName != "" {
				break
			}
		}
	}
	if deviceName == "" {
		http.Error(w, fmt.Sprintf("missing %s label or annotation", alertmanagerDeviceNameKey), http.StatusBadRequest)
		return
	}

	alertName := firstNonEmpty(payload.CommonLabels["alertname"], "Homelab alert")
	status := strings.ToUpper(firstNonEmpty(payload.Status, "firing"))
	severity := firstNonEmpty(payload.CommonLabels["severity"], "unknown")

	body := []string{
		firstNonEmpty(
			payload.CommonAnnotations["summary"],
			payload.CommonAnnotations["description"],
			alertName,
		),
	}

	for _, label := range []string{"namespace", "pod", "persistentvolumeclaim", "instance", "job"} {
		if value := payload.CommonLabels[label]; value != "" {
			body = append(body, fmt.Sprintf("%s: %s", label, value))
		}
	}

	if payload.ExternalURL != "" {
		body = append(body, payload.ExternalURL)
	}

	notif := types.Notification{
		DeviceName: deviceName,
		Title:      fmt.Sprintf("[%s] %s", status, alertName),
		Subtitle:   fmt.Sprintf("%s - %d alert(s)", severity, len(payload.Alerts)),
		Body:       strings.Join(body, "\n"),
	}

	if err := s.deliverNotification(notif); err != nil {
		http.Error(w, err.message, err.status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
