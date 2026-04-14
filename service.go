package main

import (
	"database/sql"
	"errors"
	"log"
)

func (m *Storage) getAllTasks() ([]Task, error) {
	query := "SELECT id, title, description FROM tasks"

	var tasks []Task
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m *Storage) createTask(title string, description string) {
	task := Task{
		ID:          0, // ID будет сгенерирован базой данных
		Title:       title,
		Description: description,
	}
	query := "INSERT INTO tasks (title, description) VALUES ($1, $2)"
	_, err := m.DB.Exec(query, task.Title, task.Description)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
	}

}

func (s *Storage) deleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := s.DB.Exec(query, id)
	return err
}
func (s *Storage) updateTask(id int, title string, description string) error {
	query := "UPDATE tasks SET title = $1, description = $2 WHERE id = $3"
	_, err := s.DB.Exec(query, title, description, id)
	return err

}
func (s *Storage) GetOneTask(id int) (*Task, error) {
	query := "SELECT id, title, description FROM tasks WHERE id = $1"

	var task Task
	// QueryRow возвращает одну строку
	err := s.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description)

	if err != nil {
		// Если запись не найдена, sql возвращает ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Задача не найдена
		}
		return nil, err // Ошибка базы данных
	}

	return &task, nil
}
