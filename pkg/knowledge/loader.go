package knowledge

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type KnowledgeCatalog struct {
	Endpoints []map[string]interface{} `yaml:"endpoints"`
	Fields    []map[string]interface{} `yaml:"fields"`
}

func LoadCatalog(path string) (*KnowledgeCatalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading knowledge catalog: %w", err)
	}

	var cat KnowledgeCatalog
	if err := yaml.Unmarshal(data, &cat); err != nil {
		return nil, fmt.Errorf("failed parsing YAML catalog: %w", err)
	}

	return &cat, nil
}
