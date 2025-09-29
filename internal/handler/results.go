package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sidharth-chauhan/task-scheduler/internal/db"
	"github.com/sidharth-chauhan/task-scheduler/internal/models"
)

// GET /tasks/{id}/results → results of one task
func GetTaskResults(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var results []models.TaskResult
	db.DB.Where("task_id = ?", id).Order("run_at DESC").Find(&results)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode results", http.StatusInternalServerError)
	}
}

// GET /results → all results (for all tasks)
func ListResults(w http.ResponseWriter, r *http.Request) {
	var results []models.TaskResult
	db.DB.Order("run_at DESC").Find(&results)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode results", http.StatusInternalServerError)
	}
}
