package knowledge

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func GetGeneratedKnowledgeDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cwd, "knowledge", "generated")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// ExportDiscoveredKnowledge writes discovered endpoints and schema fields into knowledge/generated/*.yaml
func ExportDiscoveredKnowledge(endpoints []domain.EndpointObservation, fields []domain.FieldObservation) (string, error) {
	genDir, err := GetGeneratedKnowledgeDir()
	if err != nil {
		return "", err
	}

	// 1. Export endpoints.yaml
	var epList []map[string]interface{}
	for _, ep := range endpoints {
		epList = append(epList, map[string]interface{}{
			"method":      ep.Method,
			"url":         ep.URL,
			"status_code": ep.LastStatusCode,
			"occurrences": ep.Occurrences,
			"first_seen":  ep.FirstSeen.Format(time.RFC3339),
			"last_seen":   ep.LastSeen.Format(time.RFC3339),
		})
	}
	epData, err := yaml.Marshal(map[string]interface{}{"endpoints": epList})
	if err == nil {
		_ = os.WriteFile(filepath.Join(genDir, "endpoints.yaml"), epData, 0644)
	}

	// 2. Export fields.yaml
	var fieldList []map[string]interface{}
	for _, f := range fields {
		fieldList = append(fieldList, map[string]interface{}{
			"endpoint":    f.Endpoint,
			"path":        f.Path,
			"type":        f.Type,
			"occurrences": f.Occurrences,
			"confidence":  f.Confidence,
			"first_seen":  f.FirstSeen.Format(time.RFC3339),
		})
	}
	fieldData, err := yaml.Marshal(map[string]interface{}{"fields": fieldList})
	if err == nil {
		_ = os.WriteFile(filepath.Join(genDir, "fields.yaml"), fieldData, 0644)
	}

	return genDir, nil
}
