package context_test

import (
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"testing"
)

func TestScopeValidation(t *testing.T) {
	s := appcontext.DefaultScope()
	if s.ProfileID != 1 || s.ApplicationID != 1 {
		t.Errorf("expected default scope (1, 1), got (%d, %d)", s.ProfileID, s.ApplicationID)
	}

	_, err := appcontext.NewScope(0, 1)
	if err == nil {
		t.Errorf("expected error for invalid profileID 0")
	}

	_, err = appcontext.NewScope(1, 0)
	if err == nil {
		t.Errorf("expected error for invalid applicationID 0")
	}

	valid, err := appcontext.NewScope(2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid.String() != "Profile #2 | App #3" {
		t.Errorf("unexpected string output: %s", valid.String())
	}
}
