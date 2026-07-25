package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(home, ".config", "anef", "config.json")
		if data, err := os.ReadFile(userConfigPath); err == nil {
			_ = json.Unmarshal(data, cfg)
			return cfg, nil
		}
	}

	cwd, _ := os.Getwd()
	localConfigPath := filepath.Join(cwd, "config.json")
	if data, err := os.ReadFile(localConfigPath); err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	return cfg, nil
}
