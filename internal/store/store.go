package store

import (
	"database/sql"
	"fmt"
	"sync"

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

func (s *Store) ListDevices() ([]Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, key, name, token FROM device")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
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
		devices = make([]Device, 0)
	}

	return devices, nil
}

func (s *Store) GetDevice(key string) (*Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var device Device
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

func (s *Store) GetDeviceByName(name string) (*Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var device Device
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

func (s *Store) CreateDevice(key string, name string, token string) (*Device, error) {
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

	return &Device{id, key, name, token}, nil
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

func (s *Store) RecordNotification(notification Notification) error {
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
