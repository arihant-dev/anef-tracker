package domain

import (
	"time"
)

type Severity string

const (
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// ApplicationStatus holds status code, human-readable label, and administrative description.
type ApplicationStatus struct {
	Code        string   `json:"code"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
}

// Application represents the current state of a residence permit application.
type Application struct {
	ID                  string                 `json:"id"`
	NumeroDemande       string                 `json:"numero_demande"`
	ForeignerID         string                 `json:"foreigner_id"`
	LegalCategory       string                 `json:"legal_category"`
	RegulatoryReference string                 `json:"regulatory_reference"`
	Status              ApplicationStatus      `json:"status"`
	ProcessingSite      string                 `json:"processing_site"`
	CollectionSite      string                 `json:"collection_site"`
	Version             int64                  `json:"version"`
	SubmittedAt         time.Time              `json:"submitted_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Documents           []Document             `json:"documents"`
	RawPayload          map[string]interface{} `json:"raw_payload"`
	RawJSON             []byte                 `json:"raw_json"`
}

// Document represents an attestation or generated document.
type Document struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// Event represents a state transition or attribute change event in the Event Store.
type Event struct {
	ID            int64     `json:"id"`
	ApplicationID string    `json:"application_id"`
	Type          string    `json:"type"`
	Severity      Severity  `json:"severity"`
	Confidence    float64   `json:"confidence"`
	FieldPath     string    `json:"field_path"`
	OldVal        string    `json:"old_val"`
	NewVal        string    `json:"new_val"`
	Timestamp     time.Time `json:"timestamp"`
}

// EndpointObservation represents observed HTTP endpoint metadata and statistics.
type EndpointObservation struct {
	ID             int64     `json:"id"`
	Method         string    `json:"method"`
	URL            string    `json:"url"`
	LastStatusCode int       `json:"last_status_code"`
	LastLatencyMs  int64     `json:"last_latency_ms"`
	FirstSeen      time.Time `json:"first_seen"`
	LastSeen       time.Time `json:"last_seen"`
	Occurrences    int       `json:"occurrences"`
}

// FieldObservation represents a registered API field, data type, and confidence score.
type FieldObservation struct {
	ID          int64     `json:"id"`
	Endpoint    string    `json:"endpoint"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Occurrences int       `json:"occurrences"`
	Confidence  float64   `json:"confidence"`
	Examples    []string  `json:"examples,omitempty"`
}
