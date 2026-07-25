package storage_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/storage"
	"testing"
)

func TestSQLiteStore(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	store := storage.NewSQLiteStore(database)

	app := &domain.Application{
		ID:            "APP-TEST-123",
		NumeroDemande: "9929006580",
		ForeignerID:   "9929006580",
		LegalCategory: "ETUDIANT",
		Status: domain.ApplicationStatus{
			Code:        "TITRE_A_FABRIQUER",
			Label:       "Residence Permit in Production",
			Description: "Test status description",
		},
		Version: 1,
	}

	if err := store.SaveApplication(app); err != nil {
		t.Fatalf("SaveApplication failed: %v", err)
	}

	ev := domain.Event{
		ApplicationID: "APP-TEST-123",
		Type:          "MODIFIED",
		Severity:      domain.SeverityHigh,
		Confidence:    1.0,
		FieldPath:     "statut",
		OldVal:        "INSTRUCTION_EN_COURS",
		NewVal:        "TITRE_A_FABRIQUER",
	}

	if err := store.SaveEvent(ev); err != nil {
		t.Fatalf("SaveEvent failed: %v", err)
	}

	events, err := store.GetEvents(5)
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("expected at least 1 event in store")
	}
}
