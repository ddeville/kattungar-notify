package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath))
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

func (s *Store) CreateDevice(device Device) (*Device, error) {
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

func (s *Store) DeleteDevice(device Device) error {
	_, err := s.db.Exec("DELETE FROM device WHERE id = ?", device.Id)
	return err
}

func (s *Store) UpdateDevice(device Device) error {
	_, err := s.db.Exec("UPDATE device SET name = ?, token = ? WHERE id = ?", device.Name, device.Token, device.Id)
	return err
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
