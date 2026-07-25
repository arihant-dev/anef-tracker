package privacy_test

import (
	"encoding/json"
	"github.com/arihant-dev/anef-tracker/pkg/privacy"
	"testing"
)

func FuzzScanPayload(f *testing.F) {
	f.Add([]byte(`{"telephone":"0745564400","nom":"JAIN"}`))
	f.Add([]byte(`{"adresse":"12 Rue Paris","numEtranger":"9929006580"}`))
	f.Add([]byte(`{}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var payload map[string]interface{}
		if err := json.Unmarshal(data, &payload); err != nil {
			return
		}

		// Verify scanner does not panic on arbitrary fuzz input
		fields := privacy.ScanPayload(payload)
		_ = fields

		// Verify redactor does not panic on arbitrary fuzz input
		redacted := privacy.RedactPayload(payload)
		_ = redacted
	})
}
