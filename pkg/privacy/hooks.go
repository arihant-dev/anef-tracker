package privacy

import (
	"time"
)

type PrivacyObserver interface {
	Inspect(payload map[string]interface{}) PrivacyResult
}

type PrivacyResult struct {
	PIIDetected bool      `json:"pii_detected"`
	PIITypes    []PIIType `json:"pii_types"`
	FieldCount  int       `json:"field_count"`
	ScannedAt   time.Time `json:"scanned_at"`
}

type DefaultObserver struct{}

func NewDefaultObserver() *DefaultObserver {
	return &DefaultObserver{}
}

func (o *DefaultObserver) Inspect(payload map[string]interface{}) PrivacyResult {
	fields := ScanPayload(payload)
	typeSet := make(map[PIIType]bool)
	var pTypes []PIIType

	for _, f := range fields {
		if !typeSet[f.Type] {
			typeSet[f.Type] = true
			pTypes = append(pTypes, f.Type)
		}
	}

	return PrivacyResult{
		PIIDetected: len(fields) > 0,
		PIITypes:    pTypes,
		FieldCount:  len(fields),
		ScannedAt:   time.Now(),
	}
}
