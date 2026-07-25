package notify

import (
	"time"
)

type EventType string

const (
	EventStatusChange       EventType = "STATUS_CHANGE"
	EventDocumentDiscovered EventType = "DOCUMENT_DISCOVERED"
	EventAuthWarning        EventType = "AUTH_WARNING"
	EventBackupCompleted    EventType = "BACKUP_COMPLETED"
	EventIntegrityFailure   EventType = "INTEGRITY_FAILURE"
	EventSecurityWarning    EventType = "SECURITY_WARNING"
	EventWatchFailure       EventType = "WATCH_FAILURE"
)

type NotificationEvent struct {
	ID         int64     `json:"id"`
	Type       EventType `json:"type"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	SnapshotID string    `json:"snapshot_id,omitempty"`
	EventID    int64     `json:"event_id,omitempty"`
	HTTPLogID  int64     `json:"http_log_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Delivered  bool      `json:"delivered"`
}
