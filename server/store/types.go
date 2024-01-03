package store

type Device struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Notification struct {
	DeviceName string `json:"device_name"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Body       string `json:"body"`
}
