package storage

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
)

type SnapshotStore interface {
	SaveSnapshotRef(appID, snapshotID, directory string) error
	GetLatestTwoSnapshots() (*snapshot.SnapshotRef, *snapshot.SnapshotRef, error)
}

type EventStore interface {
	SaveEvent(ev domain.Event) error
	GetEvents(limit int) ([]domain.Event, error)
}

type ApplicationStore interface {
	SaveApplication(app *domain.Application) error
}

type Store interface {
	SnapshotStore
	EventStore
	ApplicationStore
}

type SQLiteStore struct {
	DB *db.DB
}

func NewSQLiteStore(database *db.DB) *SQLiteStore {
	return &SQLiteStore{DB: database}
}

func (s *SQLiteStore) SaveSnapshotRef(appID, snapshotID, directory string) error {
	_, err := s.DB.SaveSnapshotRef(appID, directory)
	return err
}

func (s *SQLiteStore) GetLatestTwoSnapshots() (*snapshot.SnapshotRef, *snapshot.SnapshotRef, error) {
	return snapshot.GetLatestTwoSnapshots()
}

func (s *SQLiteStore) SaveEvent(ev domain.Event) error {
	return s.DB.SaveEvent(ev)
}

func (s *SQLiteStore) GetEvents(limit int) ([]domain.Event, error) {
	return s.DB.GetEvents(limit)
}

func (s *SQLiteStore) SaveApplication(app *domain.Application) error {
	return s.DB.SaveApplication(app)
}
