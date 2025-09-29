package models

import "time"

type TaskResult struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID          int       `json:"task_id"`
	RunAt           time.Time `json:"run_at"`
	StatusCode      int       `json:"status_code"`
	Success         bool      `json:"success"`
	ResponseHeaders string    `json:"response_headers"`
	ResponseBody    string    `json:"response_body"`
	ErrorMessage    string    `json:"error_message"`
	DurationMS      int       `json:"duration_ms"`
	CreatedAt       time.Time `json:"created_at"`
}
