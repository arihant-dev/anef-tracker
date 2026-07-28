package auth_test

import (
	"testing"

	"github.com/arihant-dev/anef-tracker/pkg/session"
)

func TestOpenBrowser(t *testing.T) {
	// OpenBrowser shouldn't panic when given a valid URL format
	_ = session.OpenBrowser("https://example.com")
}
