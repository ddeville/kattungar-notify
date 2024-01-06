package store

type Device struct {
	Id    int64
	Key   string `json:"key"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Notification struct {
	Id        int64
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Body      string `json:"body"`
}

type CalendarEvent struct {
	Id       int64
	EventId  string `json:"event_id"`
	Notified bool   `json:"notified"`
}
