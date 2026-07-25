package export

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

func ToJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func ToYAML(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
