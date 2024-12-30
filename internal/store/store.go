package store

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/ddeville/kattungar-notify/internal/types"
	"github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath))
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS device (id INTEGER NOT NULL PRIMARY KEY, key TEXT UNIQUE, name TEXT, token TEXT);
	CREATE TABLE IF NOT EXISTS notification (id INTEGER NOT NULL PRIMARY KEY, device_key TEXT, device_name TEXT, title TEXT, subtitle TEXT, body TEXT);
	CREATE TABLE IF NOT EXISTS calendar_event (id INTEGER NOT NULL PRIMARY KEY, event_id TEXT, notified INTEGER);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) ListDevices() ([]types.Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, key, name, token FROM device")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []types.Device
	for rows.Next() {
		var device types.Device
		err = rows.Scan(&device.Id, &device.Key, &device.Name, &device.Token)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if devices == nil {
		devices = make([]types.Device, 0)
	}

	return devices, nil
}

func (s *Store) GetDevice(key string) (*types.Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var device types.Device
	row := s.db.QueryRow("SELECT id, key, name, token FROM device WHERE key = ?", key)
	err := row.Scan(&device.Id, &device.Key, &device.Name, &device.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (s *Store) GetDeviceByName(name string) (*types.Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var device types.Device
	row := s.db.QueryRow("SELECT id, key, name, token FROM device WHERE name = ?", name)
	err := row.Scan(&device.Id, &device.Key, &device.Name, &device.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (s *Store) CreateDevice(key string, name string, token string) (*types.Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == "" {
		return nil, fmt.Errorf("missing device key")
	}

	res, err := s.db.Exec("INSERT INTO device (key, name, token) VALUES (?, ?, ?)", key, name, token)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.Device{Id: id, Key: key, Name: name, Token: token}, nil
}

func (s *Store) UpdateDeviceName(key string, name string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("UPDATE device SET name = ? WHERE key = ?", name, key)
	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	}

	return true, err
}

func (s *Store) UpdateDeviceToken(key string, token string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("UPDATE device SET token = ? WHERE key = ?", token, key)
	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	}

	return true, err
}

func (s *Store) DeleteDevice(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("DELETE FROM device WHERE key = ?", key)
	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	}

	return true, err
}

func IsExistingDeviceError(err error) bool {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return true
		}
	}
	return false
}

func (s *Store) RecordNotification(notification types.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		"INSERT INTO notification (device_key, device_name, title, subtitle, body) VALUES (?, ?, ?, ?, ?)",
		notification.DeviceKey,
		notification.DeviceName,
		notification.Title,
		notification.Subtitle,
		notification.Body,
	)
	return err
}

func (s *Store) ListNotifications(device *types.Device) ([]types.Notification, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, device_key, device_name, title, subtitle, body FROM notification WHERE device_key = ?", device.Key)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []types.Notification
	for rows.Next() {
		var notif types.Notification
		err = rows.Scan(&notif.Id, &notif.DeviceKey, &notif.DeviceName, &notif.Title, &notif.Subtitle, &notif.Body)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if notifications == nil {
		notifications = make([]types.Notification, 0)
	}

	return notifications, nil
}

func (s *Store) AddCalendarEvent(id string, notified bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("INSERT INTO calendar_event (event_id, notified) VALUES (?, ?)", id, notified)
	return err
}

func (s *Store) HasNotifiedCalendarEvent(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	var notified int
	row := s.db.QueryRow("SELECT notified FROM calendar_event WHERE event_id = ?", id)
	err := row.Scan(&notified)
	if err != nil {
		return false
	}
	return notified != 0
}
