package config_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/config"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Database.Path == "" {
		t.Errorf("expected non-empty default database path")
	}

	if !cfg.Storage.ImmutableSnapshots {
		t.Errorf("expected ImmutableSnapshots to be true by default")
	}

	summary := cfg.FormatSummary()
	if len(summary) == 0 {
		t.Errorf("expected non-empty config summary")
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatalf("expected non-nil config")
	}
}
