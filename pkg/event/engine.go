package event

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/diff"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/notify"
	"log"
)

type EventEngine struct {
	DB       *db.DB
	Notifier notify.Notifier
}

func NewEventEngine(database *db.DB, n notify.Notifier) *EventEngine {
	return &EventEngine{
		DB:       database,
		Notifier: n,
	}
}

// ProcessSnapshotDiff converts structural JSON diff changes into stored domain events and triggers notifications.
func (e *EventEngine) ProcessSnapshotDiff(appID string, diffResult *diff.DiffResult) ([]domain.Event, error) {
	if diffResult == nil || !diffResult.HasChanges {
		return nil, nil
	}

	var emittedEvents []domain.Event

	for _, ch := range diffResult.Changes {
		ev := domain.Event{
			ApplicationID: appID,
			Type:          ch.Change,
			Severity:      ch.Severity,
			Confidence:    1.0, // Concrete diff confidence
			FieldPath:     ch.Path,
			OldVal:        fmt.Sprintf("%v", ch.OldVal),
			NewVal:        fmt.Sprintf("%v", ch.NewVal),
		}

		if e.DB != nil {
			if err := e.DB.SaveEvent(ev); err != nil {
				log.Printf("[WARN] Failed saving event to SQLite: %v", err)
			}
		}

		// Trigger notification for HIGH or CRITICAL severity events
		if ev.Severity == domain.SeverityHigh || ev.Severity == domain.SeverityCritical {
			if e.Notifier != nil {
				if err := e.Notifier.Notify(ev); err != nil {
					log.Printf("[WARN] Failed dispatching notification: %v", err)
				}
			}
		}

		emittedEvents = append(emittedEvents, ev)
	}

	return emittedEvents, nil
}
