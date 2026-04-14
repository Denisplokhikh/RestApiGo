package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv" // To register the driver.
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "host=db port=5432 user=postgres password=postgres dbname=tasks sslmode=disable"
	}
	storage, err := NewStorage(connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer storage.Close()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {

		tasks, err := storage.getAllTasks()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})
	http.HandleFunc("/tasks/delete/{id}", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = storage.deleteTask(id)
		if err != nil {
			http.Error(w, "Task not found or error deleting", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted"))
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req EditTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		storage.updateTask(req.ID, req.Title, req.Description)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte("edited"))
	})
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		storage.createTask(req.Title, req.Description)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte("created"))
	})
	http.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. Получаем ID из URL path
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// 2. Получаем задачу
		task, err := storage.GetOneTask(id)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if task == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// 3. Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK для GET

		json.NewEncoder(w).Encode(task)
	})
	log.Println("Server started on :5001")

	err2 := http.ListenAndServe(":5001", nil)
	if err2 != nil {
		log.Fatal("Server error:", err2)
	}
}
