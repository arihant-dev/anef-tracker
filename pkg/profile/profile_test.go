package profile_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/profile"
	"testing"
)

func TestProfileCreationAndSwitch(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	p, err := profile.CreateProfile(database, "Work Permit Vault")
	if err != nil {
		t.Fatalf("CreateProfile failed: %v", err)
	}

	if err := profile.SwitchProfile(database, p.ID); err != nil {
		t.Fatalf("SwitchProfile failed: %v", err)
	}

	active, err := profile.GetActiveProfile(database)
	if err != nil {
		t.Fatalf("GetActiveProfile failed: %v", err)
	}

	if active.ID != p.ID {
		t.Errorf("expected active profile ID %d, got %d", p.ID, active.ID)
	}
}

func TestLegacyDatabaseMigrationCreatesDefaultProfile(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	profiles, err := profile.ListProfiles(database)
	if err != nil {
		t.Fatalf("ListProfiles failed: %v", err)
	}

	if len(profiles) == 0 {
		t.Errorf("expected at least Default Profile after migrations")
	}
}

func TestDeleteProfileRequiresConfirmation(t *testing.T) {
	err := profile.DeleteProfile(nil, 5, 4)
	if err == nil {
		t.Errorf("expected error when confirmation ID mismatch")
	}
}
