package runner

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sidharth-chauhan/task-scheduler/internal/db"
	"github.com/sidharth-chauhan/task-scheduler/internal/models"
)

func Tick() {
	now := time.Now().UTC()

	var oneOff []models.Task
	db.DB.Where("status = ? AND type = ? AND next_run <= ?", "scheduled", "one-off", now).Find(&oneOff)

	var cronTasks []models.Task
	db.DB.Where("status = ? AND type = ? AND next_run <= ?", "scheduled", "cron", now).Find(&cronTasks)

	for _, t := range oneOff {
		runAndRecord(&t)
		t.Status = "completed"
		db.DB.Save(&t)
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	for _, t := range cronTasks {
		runAndRecord(&t)
		if spec, err := parser.Parse(strings.TrimSpace(t.Cron)); err == nil {
			t.NextRun = spec.Next(now)
			db.DB.Save(&t)
		}
	}
}

func runAndRecord(t *models.Task) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// no payload/body at all
	req, err := http.NewRequestWithContext(ctx, t.Method, t.URL, nil)
	if err != nil {
		saveError(t.ID, err.Error(), start)
		return
	}

	// optional headers
	if t.Headers != "" && t.Headers != "{}" {
		var hs map[string]string
		_ = json.Unmarshal([]byte(t.Headers), &hs)
		for k, v := range hs {
			req.Header.Set(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		saveError(t.ID, err.Error(), start)
		return
	}
	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	// response headers → JSON string
	hmap := map[string]string{}
	for k, v := range resp.Header {
		if len(v) > 0 {
			hmap[k] = v[0]
		}
	}
	hjson, _ := json.Marshal(hmap)

	db.DB.Create(&models.TaskResult{
		TaskID:          t.ID,
		RunAt:           time.Now().UTC(),
		StatusCode:      resp.StatusCode,
		Success:         resp.StatusCode >= 200 && resp.StatusCode < 300,
		ResponseHeaders: string(hjson),
		ResponseBody:    string(b),
		DurationMS:      int(time.Since(start) / time.Millisecond),
		CreatedAt:       time.Now().UTC(),
	})
}

func saveError(taskID int, msg string, start time.Time) {
	db.DB.Create(&models.TaskResult{
		TaskID:       taskID,
		RunAt:        time.Now().UTC(),
		Success:      false,
		ErrorMessage: msg,
		DurationMS:   int(time.Since(start) / time.Millisecond),
		CreatedAt:    time.Now().UTC(),
	})
}
