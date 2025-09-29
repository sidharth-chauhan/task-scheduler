package models

import "time"

type Task struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "one-off" or "cron"
	UTCDatetime time.Time `json:"utc_datetime"`
	Cron        string    `json:"cron"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Headers     string    `json:"headers"` // optional JSON string like "{}"
	Status      string    `json:"status"`  // scheduled | cancelled | completed
	NextRun     time.Time `json:"next_run"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
