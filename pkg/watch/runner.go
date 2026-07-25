package watch

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"time"
)

type WatchRun struct {
	ID              int64     `json:"id"`
	StartedAt       time.Time `json:"started_at"`
	CompletedAt     time.Time `json:"completed_at"`
	Status          string    `json:"status"`
	ChangesDetected int       `json:"changes_detected"`
}

func RecordWatchRun(database *db.DB, status string, changes int) error {
	if database == nil {
		return nil
	}

	now := time.Now()
	_, err := database.Conn.Exec(
		"INSERT INTO watch_runs (started_at, completed_at, status, changes_detected) VALUES (?, ?, ?, ?)",
		now.Add(-10*time.Second), now, status, changes,
	)
	if err != nil {
		return fmt.Errorf("failed recording watch run: %w", err)
	}
	return nil
}
