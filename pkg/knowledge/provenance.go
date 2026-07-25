package knowledge

import (
	"time"
)

type SourceType string

const (
	SourceSnapshot SourceType = "SNAPSHOT"
	SourceEvent    SourceType = "EVENT"
	SourceHTTP     SourceType = "HTTP_LOG"
	SourceSchema   SourceType = "SCHEMA"
)

type Provenance struct {
	SourceType SourceType `json:"source_type" yaml:"source_type"`
	SnapshotID string     `json:"snapshot_id,omitempty" yaml:"snapshot_id,omitempty"`
	EventID    int64      `json:"event_id,omitempty" yaml:"event_id,omitempty"`
	HTTPLogID  int64      `json:"http_log_id,omitempty" yaml:"http_log_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at" yaml:"created_at"`
}
