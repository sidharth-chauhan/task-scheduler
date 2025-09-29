package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sidharth-chauhan/task-scheduler/internal/db"
	"github.com/sidharth-chauhan/task-scheduler/internal/models"
)

// POST /tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&task).Error; err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}

// GET /tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.DB.Order("created_at DESC").Find(&tasks)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

// GET /tasks/{id}
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var task models.Task
	if err := db.DB.First(&task, id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
	}
}

// PUT /tasks/{id}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var task models.Task
	if err := db.DB.First(&task, id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var in models.Task
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if in.Name != "" {
		task.Name = in.Name
	}
	if in.Status != "" {
		task.Status = in.Status
	}
	if in.Method != "" {
		task.Method = in.Method
	}
	if in.URL != "" {
		task.URL = in.URL
	}
	if in.Headers != "" {
		task.Headers = in.Headers
	}
	if !in.NextRun.IsZero() {
		task.NextRun = in.NextRun
	}
	if in.Cron != "" {
		task.Cron = in.Cron
	}
	if !in.UTCDatetime.IsZero() {
		task.UTCDatetime = in.UTCDatetime
	}
	if in.Type != "" {
		task.Type = in.Type
	}

	if err := db.DB.Save(&task).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(task)
}

// DELETE /tasks/{id} (cancel)
func CancelTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var task models.Task
	if err := db.DB.First(&task, id).Error; err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	task.Status = "cancelled"
	if err := db.DB.Save(&task).Error; err != nil {
		http.Error(w, "cancel failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "task cancelled"})
}
