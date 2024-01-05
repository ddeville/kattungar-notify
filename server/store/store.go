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
	CREATE TABLE IF NOT EXISTS device (id INTEGER NOT NULL PRIMARY KEY, name TEXT UNIQUE, token TEXT);
	CREATE TABLE IF NOT EXISTS notification (id INTEGER NOT NULL PRIMARY KEY, device_name TEXT, title TEXT, subtitle TEXT, body TEXT);
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

	rows, err := s.db.Query("SELECT id, name, token FROM device")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err = rows.Scan(&device.Id, &device.Name, &device.Token)
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

func (s *Store) GetDevice(name string) (*Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var device Device
	row := s.db.QueryRow("SELECT id, name, token FROM device WHERE name = ?", name)
	err := row.Scan(&device.Id, &device.Name, &device.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (s *Store) CreateDevice(device Device) (*Device, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("INSERT INTO device (name, token) VALUES (?, ?)", device.Name, device.Token)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Device{id, device.Name, device.Token}, nil
}

func (s *Store) UpdateDevice(device Device) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("UPDATE device SET name = ?, token = ? WHERE id = ?", device.Name, device.Token, device.Id)
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

func (s *Store) DeleteDevice(device Device) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec("DELETE FROM device WHERE id = ?", device.Id)
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
		"INSERT INTO notification (device_name, title, subtitle, body) VALUES (?, ?, ?, ?)",
		notification.DeviceName,
		notification.Title,
		notification.Subtitle,
		notification.Body,
	)
	return err
}
