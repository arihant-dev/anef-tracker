package security_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/security"
	"testing"
)

func TestTokenRedaction(t *testing.T) {
	redactedAuth := security.RedactHeaderValue("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
	if redactedAuth != "[REDACTED]" {
		t.Errorf("expected Authorization header to be [REDACTED], got '%s'", redactedAuth)
	}

	redactedCookie := security.RedactHeaderValue("Cookie", "session_token=abcdef123456")
	if redactedCookie != "[REDACTED]" {
		t.Errorf("expected Cookie header to be [REDACTED], got '%s'", redactedCookie)
	}

	normalHeader := security.RedactHeaderValue("Accept", "application/json")
	if normalHeader != "application/json" {
		t.Errorf("expected Accept header to remain unchanged, got '%s'", normalHeader)
	}
}

func TestPermissionAudit(t *testing.T) {
	rep := security.AuditSecurity()
	if rep == nil {
		t.Fatalf("expected non-nil SecurityReport")
	}

	if !rep.TokensRedacted {
		t.Errorf("expected TokensRedacted to be true")
	}
}
