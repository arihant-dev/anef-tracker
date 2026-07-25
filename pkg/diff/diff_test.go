package diff_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/diff"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"testing"
)

func TestCompareSnapshots(t *testing.T) {
	oldJSON := []byte(`{
		"numero_demande": "9929006580",
		"statut": "INSTRUCTION_EN_COURS",
		"_version": 56
	}`)

	newJSON := []byte(`{
		"numero_demande": "9929006580",
		"statut": "TITRE_A_FABRIQUER",
		"_version": 57,
		"new_flag": true
	}`)

	res, err := diff.CompareSnapshots(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("CompareSnapshots failed: %v", err)
	}

	if !res.HasChanges {
		t.Fatalf("expected HasChanges = true")
	}

	if res.Severity != domain.SeverityHigh {
		t.Errorf("expected Severity HIGH for status change, got %s", res.Severity)
	}

	if len(res.Changes) < 3 {
		t.Errorf("expected at least 3 changes (statut, _version, new_flag), got %d", len(res.Changes))
	}
}
