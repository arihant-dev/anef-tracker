package raw

import (
	"encoding/json"
	"fmt"
)

// RawPayload holds raw API JSON bytes as provider-agnostic forensic evidence.
type RawPayload struct {
	RawJSON json.RawMessage        `json:"raw_json"`
	Map     map[string]interface{} `json:"map"`
}

// UnmarshalRaw creates a RawPayload forensic wrapper preserving the exact response bytes.
func UnmarshalRaw(data []byte) (*RawPayload, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty raw JSON payload")
	}

	var parsedMap map[string]interface{}
	if err := json.Unmarshal(data, &parsedMap); err != nil {
		return nil, fmt.Errorf("failed parsing raw JSON payload map: %w", err)
	}

	return &RawPayload{
		RawJSON: json.RawMessage(data),
		Map:     parsedMap,
	}, nil
}
