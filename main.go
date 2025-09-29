package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/sidharth-chauhan/task-scheduler/internal/db"
	"github.com/sidharth-chauhan/task-scheduler/internal/handler"
	"github.com/sidharth-chauhan/task-scheduler/internal/handler/runner"
	"github.com/sidharth-chauhan/task-scheduler/internal/models"
)

func main() {

	_ = godotenv.Load()

	db.ConnectDB()

	if err := db.DB.AutoMigrate(&models.Task{}, &models.TaskResult{}); err != nil {
		fmt.Println("migration failed:", err)
		return
	}

	go func() {
		for {
			runner.Tick()
			time.Sleep(2 * time.Second)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", handler.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handler.CancelTask).Methods("DELETE")
	r.HandleFunc("/tasks/{id}/results", handler.GetTaskResults).Methods("GET")
	r.HandleFunc("/results", handler.ListResults).Methods("GET")

	fmt.Println("server on :8080")
	_ = http.ListenAndServe(":8080", r)
}
