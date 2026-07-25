package privacy_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/privacy"
	"testing"
)

func TestSensitiveFieldDetection(t *testing.T) {
	payload := map[string]interface{}{
		"telephone":       "0745564400",
		"nom":             "JAIN",
		"adresse":         "12 Rue de Paris",
		"numeroEtranger":  "9929006580",
		"statut_sans_pii": "INSTRUCTION_EN_COURS",
	}

	fields := privacy.ScanPayload(payload)
	if len(fields) < 4 {
		t.Errorf("expected at least 4 PII fields, found %d", len(fields))
	}
}

func TestRedactedExport(t *testing.T) {
	payload := map[string]interface{}{
		"telephone": "0745564400",
		"nom":       "JAIN",
	}

	redacted := privacy.RedactPayload(payload)
	if redacted["telephone"] == "0745564400" {
		t.Errorf("expected telephone to be redacted")
	}
	if payload["telephone"] != "0745564400" {
		t.Errorf("original payload mutated by redaction!")
	}
}

func TestPrivacyResultContainsNoPII(t *testing.T) {
	obs := privacy.NewDefaultObserver()
	res := obs.Inspect(map[string]interface{}{"telephone": "0745564400"})
	if !res.PIIDetected {
		t.Errorf("expected PIIDetected = true")
	}
}
