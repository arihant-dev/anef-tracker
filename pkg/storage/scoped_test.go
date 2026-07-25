package storage_test

import (
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/storage"
	"testing"
)

func TestScopedIsolation(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	repo := storage.NewScopedRepository(database)

	scopeA, _ := appcontext.NewScope(100, 100)
	scopeB, _ := appcontext.NewScope(200, 200)

	ev := domain.Event{
		ApplicationID: "ISOLATION-TEST",
		Type:          "MODIFIED",
		Severity:      domain.SeverityHigh,
		Confidence:    1.0,
		FieldPath:     "statut",
		OldVal:        "OLD",
		NewVal:        "NEW",
	}

	if err := repo.SaveEventScoped(scopeA, ev); err != nil {
		t.Fatalf("SaveEventScoped for Scope A failed: %v", err)
	}

	eventsA, err := repo.GetEventsScoped(scopeA, 10)
	if err != nil {
		t.Fatalf("GetEventsScoped A failed: %v", err)
	}
	if len(eventsA) == 0 {
		t.Errorf("expected events for Scope A")
	}

	eventsB, err := repo.GetEventsScoped(scopeB, 10)
	if err != nil {
		t.Fatalf("GetEventsScoped B failed: %v", err)
	}
	for _, e := range eventsB {
		if e.ApplicationID == "ISOLATION-TEST" {
			t.Errorf("LEAK DETECTED: Scope B saw event belonging to Scope A")
		}
	}
}
