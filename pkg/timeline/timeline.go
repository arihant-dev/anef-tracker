package timeline

import (
	"time"
)

type TimelineMilestone struct {
	Date          time.Time `json:"date" yaml:"date"`
	StatusCode    string    `json:"status_code" yaml:"status_code"`
	StatusLabel   string    `json:"status_label" yaml:"status_label"`
	Completed     bool      `json:"completed" yaml:"completed"`
	IsCurrent     bool      `json:"is_current" yaml:"is_current"`
	DaysInState   int       `json:"days_in_state" yaml:"days_in_state"`
	SnapshotID    string    `json:"snapshot_id,omitempty" yaml:"snapshot_id,omitempty"`
	EventID       int64     `json:"event_id,omitempty" yaml:"event_id,omitempty"`
	EvidenceCount int       `json:"evidence_count" yaml:"evidence_count"`
}

type HumanTimeline struct {
	ApplicationID string              `json:"application_id"`
	StartedAt     time.Time           `json:"started_at"`
	CurrentStatus string              `json:"current_status"`
	ElapsedDays   int                 `json:"elapsed_days"`
	Milestones    []TimelineMilestone `json:"milestones"`
}
