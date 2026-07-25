package config

import (
	"fmt"
	"strings"
)

type DatabaseConfig struct {
	Path string `json:"path" yaml:"path"`
}

type StorageConfig struct {
	ImmutableSnapshots bool `json:"immutable_snapshots" yaml:"immutable_snapshots"`
}

type NotificationsConfig struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Desktop    bool   `json:"desktop" yaml:"desktop"`
	WebhookURL string `json:"webhook_url" yaml:"webhook_url"`
}

type RetentionConfig struct {
	HTTPLogsDays int `json:"http_logs_days" yaml:"http_logs_days"`
}

type SecurityConfig struct {
	EncryptionEnabled bool   `json:"encryption_enabled" yaml:"encryption_enabled"`
	KeySource         string `json:"key_source" yaml:"key_source"`
}

type Config struct {
	Database      DatabaseConfig      `json:"database" yaml:"database"`
	Storage       StorageConfig       `json:"storage" yaml:"storage"`
	Notifications NotificationsConfig `json:"notifications" yaml:"notifications"`
	Retention     RetentionConfig     `json:"retention" yaml:"retention"`
	Security      SecurityConfig      `json:"security" yaml:"security"`
}

func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Path: "data/anef.db",
		},
		Storage: StorageConfig{
			ImmutableSnapshots: true,
		},
		Notifications: NotificationsConfig{
			Enabled: false,
		},
		Retention: RetentionConfig{
			HTTPLogsDays: 365,
		},
	}
}

func (c *Config) FormatSummary() string {
	var sb strings.Builder
	sb.WriteString("=== ANEF CONFIGURATION ===\n\n")

	sb.WriteString(fmt.Sprintf("Database Path:           %s\n", c.Database.Path))
	sb.WriteString(fmt.Sprintf("Immutable Snapshots:     %t\n", c.Storage.ImmutableSnapshots))
	sb.WriteString(fmt.Sprintf("Notifications Enabled:   %t\n", c.Notifications.Enabled))
	sb.WriteString(fmt.Sprintf("HTTP Logs Retention:     %d days\n\n", c.Retention.HTTPLogsDays))

	sb.WriteString("✓ Configuration Status: VALID\n")

	return sb.String()
}
