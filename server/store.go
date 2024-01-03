package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore() (*Store, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS device (id INTEGER NOT NULL PRIMARY KEY, name TEXT, token TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (s *Store) AddDevice(device Device) error {
	_, err := s.db.Exec("INSERT INTO device (name, token) VALUES (?, ?)", device.Name, device.Token)
	return err
}

func (s *Store) DeleteDevice(id string) {
}

func (s *Store) ListDevices() ([]Device, error) {
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

	return devices, nil
}
