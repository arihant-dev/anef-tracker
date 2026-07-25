package backup

import (
	"time"
)

type Manifest struct {
	Version       string    `json:"version" yaml:"version"`
	CreatedAt     time.Time `json:"created_at" yaml:"created_at"`
	DatabaseHash  string    `json:"database_hash" yaml:"database_hash"`
	SnapshotCount int       `json:"snapshot_count" yaml:"snapshot_count"`
	EventCount    int       `json:"event_count" yaml:"event_count"`
	HTTPLogCount  int       `json:"http_log_count" yaml:"http_log_count"`
}
