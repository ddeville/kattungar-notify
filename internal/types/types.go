// Package types contains various types used throughout the project.
package types

type Device struct {
	ID    int64
	Key   string `json:"key"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Notification struct {
	ID         int64
	DeviceKey  string `json:"device_key"`
	DeviceName string `json:"device_name"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Body       string `json:"body"`
}

type CalendarEvent struct {
	ID       int64
	EventID  string `json:"event_id"`
	Notified bool   `json:"notified"`
}
