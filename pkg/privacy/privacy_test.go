package privacy_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/privacy"
	"testing"
)

func TestSensitiveFieldDetection(t *testing.T) {
	payload := map[string]interface{}{
		"telephone":       "0600000000",
		"nom":             "DUPONT",
		"adresse":         "12 Rue de Paris",
		"numeroEtranger":  "9999999999",
		"statut_sans_pii": "INSTRUCTION_EN_COURS",
	}

	fields := privacy.ScanPayload(payload)
	if len(fields) < 4 {
		t.Errorf("expected at least 4 PII fields, found %d", len(fields))
	}
}

func TestRedactedExport(t *testing.T) {
	payload := map[string]interface{}{
		"telephone": "0600000000",
		"nom":       "DUPONT",
	}

	redacted := privacy.RedactPayload(payload)
	if redacted["telephone"] == "0600000000" {
		t.Errorf("expected telephone to be redacted")
	}
	if payload["telephone"] != "0600000000" {
		t.Errorf("original payload mutated by redaction!")
	}
}

func TestPrivacyResultContainsNoPII(t *testing.T) {
	obs := privacy.NewDefaultObserver()
	res := obs.Inspect(map[string]interface{}{"telephone": "0600000000"})
	if !res.PIIDetected {
		t.Errorf("expected PIIDetected = true")
	}
}
