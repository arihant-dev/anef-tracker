package domain_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestMapTitreAFabriquerFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "titre_a_fabriquer.json"))
	if err != nil {
		t.Fatalf("failed loading fixture: %v", err)
	}

	app, err := domain.MapJSONToApplication(data, "9999999999")
	if err != nil {
		t.Fatalf("MapJSONToApplication failed: %v", err)
	}

	if app.Status.Code != "TITRE_A_FABRIQUER" {
		t.Errorf("expected status code TITRE_A_FABRIQUER, got %s", app.Status.Code)
	}

	if app.Status.Severity != domain.SeverityHigh {
		t.Errorf("expected severity HIGH, got %s", app.Status.Severity)
	}

	if len(app.Documents) != 1 {
		t.Errorf("expected 1 document in attestation_depot, got %d", len(app.Documents))
	}
}

func TestMapDecisionValideeFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "decision_validee.json"))
	if err != nil {
		t.Fatalf("failed loading fixture: %v", err)
	}

	app, err := domain.MapJSONToApplication(data, "9999999999")
	if err != nil {
		t.Fatalf("MapJSONToApplication failed: %v", err)
	}

	if app.Status.Code != "DECISION_VALIDEE" {
		t.Errorf("expected status code DECISION_VALIDEE, got %s", app.Status.Code)
	}
}

func TestMapDemandeIncompleteFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "demande_incomplete.json"))
	if err != nil {
		t.Fatalf("failed loading fixture: %v", err)
	}

	app, err := domain.MapJSONToApplication(data, "9999999999")
	if err != nil {
		t.Fatalf("MapJSONToApplication failed: %v", err)
	}

	if app.Status.Code != "DEMANDE_INCOMPLETE" {
		t.Errorf("expected status code DEMANDE_INCOMPLETE, got %s", app.Status.Code)
	}
}
