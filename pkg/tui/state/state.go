package state

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TabPreference struct {
	Filter string `yaml:"filter"`
	Search string `yaml:"search"`
	Scroll int    `yaml:"scroll"`
}

type Preferences struct {
	PageSize       int  `yaml:"page_size"`
	ShowTimestamps bool `yaml:"show_timestamps"`
	CompactMode    bool `yaml:"compact_mode"`
}

type TUIState struct {
	LastTab     string                   `yaml:"last_tab"`
	Tabs        map[string]TabPreference `yaml:"tabs"`
	Preferences Preferences              `yaml:"preferences"`
}

func DefaultTUIState() *TUIState {
	return &TUIState{
		LastTab: "OVERVIEW",
		Tabs:    make(map[string]TabPreference),
		Preferences: Preferences{
			PageSize:       25,
			ShowTimestamps: true,
			CompactMode:    false,
		},
	}
}

func GetTUIConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".config", "anef")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "tui.yaml"), nil
}

func LoadTUIState() (*TUIState, error) {
	configPath, err := GetTUIConfigPath()
	if err != nil {
		return DefaultTUIState(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultTUIState(), nil
	}

	st := DefaultTUIState()
	if err := yaml.Unmarshal(data, st); err != nil {
		return DefaultTUIState(), nil
	}
	if st.Tabs == nil {
		st.Tabs = make(map[string]TabPreference)
	}

	return st, nil
}

func SaveTUIState(st *TUIState) error {
	if st == nil {
		return nil
	}
	configPath, err := GetTUIConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(st)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
