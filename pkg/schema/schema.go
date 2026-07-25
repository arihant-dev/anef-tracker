package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type FieldSchema struct {
	Type     string   `json:"type"`
	Examples []string `json:"examples,omitempty"`
}

type SchemaDocument struct {
	Endpoint string                 `json:"endpoint"`
	Fields   map[string]FieldSchema `json:"fields"`
}

// GenerateJSONSchema constructs a rich JSON schema from a raw JSON payload map.
func GenerateJSONSchema(endpoint string, payload map[string]interface{}) *SchemaDocument {
	doc := &SchemaDocument{
		Endpoint: endpoint,
		Fields:   make(map[string]FieldSchema),
	}

	traversePayload("", payload, doc.Fields)
	return doc
}

func traversePayload(prefix string, m map[string]interface{}, fields map[string]FieldSchema) {
	for k, v := range m {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		if v == nil {
			fields[path] = FieldSchema{Type: "null"}
		} else if subMap, ok := v.(map[string]interface{}); ok {
			traversePayload(path, subMap, fields)
		} else {
			typeName := reflect.TypeOf(v).String()
			ex := fmt.Sprintf("%v", v)
			if len(ex) > 40 {
				ex = ex[:37] + "..."
			}

			existing, exists := fields[path]
			if !exists {
				fields[path] = FieldSchema{
					Type:     typeName,
					Examples: []string{ex},
				}
			} else {
				if len(existing.Examples) < 3 && !containsString(existing.Examples, ex) {
					existing.Examples = append(existing.Examples, ex)
					fields[path] = existing
				}
			}
		}
	}
}

func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func (s *SchemaDocument) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
