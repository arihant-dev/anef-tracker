package evidence_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"testing"
)

func TestEvidenceHashGeneration(t *testing.T) {
	payload := []byte(`{"statut": "TITRE_A_FABRIQUER"}`)
	hash := evidence.CalculateHash(payload)

	if len(hash) != 64 {
		t.Errorf("expected 64-char SHA256 hex string, got %d chars", len(hash))
	}

	if !evidence.VerifyPayloadHash(payload, hash) {
		t.Errorf("expected SHA256 verification to pass for payload")
	}
}

func TestInvalidHashDetection(t *testing.T) {
	payload := []byte(`{"statut": "TITRE_A_FABRIQUER"}`)
	bogusHash := "0000000000000000000000000000000000000000000000000000000000000000"

	if evidence.VerifyPayloadHash(payload, bogusHash) {
		t.Errorf("expected SHA256 verification to fail for bogus hash")
	}
}

func TestEvidenceIntegrityValidation(t *testing.T) {
	rep, err := evidence.VerifyIntegrity(nil)
	if err != nil {
		t.Fatalf("VerifyIntegrity failed: %v", err)
	}

	if !rep.HashesValid {
		t.Errorf("expected snapshot hashes to be valid")
	}
}
