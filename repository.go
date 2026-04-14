package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Task struct {
	ID          int
	Title       string
	Description string
}

type Storage struct {
	DB *sql.DB
}

func (s *Storage) Close() {
	s.DB.Close()
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type EditTaskRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	log.Println("Connected to PostgreSQL via lib/pq")
	return &Storage{DB: db}, nil
}

type WriteTaskID struct {
	ID int `json:"ID"`
}
