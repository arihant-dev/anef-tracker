package storage

import (
	"fmt"
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"time"
)

type ScopedRepository struct {
	DB *db.DB
}

func NewScopedRepository(database *db.DB) *ScopedRepository {
	return &ScopedRepository{DB: database}
}

func (r *ScopedRepository) SaveSnapshotRefScoped(scope appcontext.Scope, appID, directory string) (int64, error) {
	if r.DB == nil {
		return 0, nil
	}
	res, err := r.DB.Conn.Exec(
		"INSERT INTO snapshots (application_id, snapshot_dir, profile_id, tracked_application_id, created_at) VALUES (?, ?, ?, ?, ?)",
		appID, directory, scope.ProfileID, scope.ApplicationID, time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed saving scoped snapshot: %w", err)
	}
	return res.LastInsertId()
}

func (r *ScopedRepository) GetEventsScoped(scope appcontext.Scope, limit int) ([]domain.Event, error) {
	if r.DB == nil {
		return []domain.Event{}, nil
	}
	rows, err := r.DB.Conn.Query(
		"SELECT id, application_id, event_type, severity, confidence, field_path, old_val, new_val, created_at FROM events WHERE profile_id = ? ORDER BY id DESC LIMIT ?",
		scope.ProfileID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed querying scoped events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		var sevStr string
		if err := rows.Scan(&e.ID, &e.ApplicationID, &e.Type, &sevStr, &e.Confidence, &e.FieldPath, &e.OldVal, &e.NewVal, &e.Timestamp); err != nil {
			return nil, err
		}
		e.Severity = domain.Severity(sevStr)
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *ScopedRepository) SaveEventScoped(scope appcontext.Scope, event domain.Event) error {
	if r.DB == nil {
		return nil
	}
	_, err := r.DB.Conn.Exec(
		"INSERT INTO events (application_id, event_type, severity, confidence, field_path, old_val, new_val, profile_id, tracked_application_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		event.ApplicationID, event.Type, string(event.Severity), event.Confidence, event.FieldPath, event.OldVal, event.NewVal, scope.ProfileID, scope.ApplicationID, time.Now(),
	)
	return err
}
