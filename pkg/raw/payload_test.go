package raw_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/raw"
	"testing"
)

func TestUnmarshalRaw(t *testing.T) {
	jsonPayload := []byte(`{
		"numero_demande": "9929006580",
		"statut": "TITRE_A_FABRIQUER",
		"_version": 56,
		"unknown_field_x": "custom_value"
	}`)

	rawPayload, err := raw.UnmarshalRaw(jsonPayload)
	if err != nil {
		t.Fatalf("UnmarshalRaw failed: %v", err)
	}

	if rawPayload.Map["numero_demande"] != "9929006580" {
		t.Errorf("expected numero_demande 9929006580, got %v", rawPayload.Map["numero_demande"])
	}

	if rawPayload.Map["statut"] != "TITRE_A_FABRIQUER" {
		t.Errorf("expected statut TITRE_A_FABRIQUER, got %v", rawPayload.Map["statut"])
	}

	if val, ok := rawPayload.Map["unknown_field_x"]; !ok || val != "custom_value" {
		t.Errorf("expected extra field unknown_field_x = custom_value, got %v", val)
	}
}
