package explorer

import (
	"encoding/json"
	"strings"
)

func FormatResponseInspection(respBody string) string {
	var sb strings.Builder
	sb.WriteString("=== RESPONSE PAYLOAD INSPECTION ===\n")

	var obj interface{}
	if err := json.Unmarshal([]byte(respBody), &obj); err == nil {
		indented, err := json.MarshalIndent(obj, "", "  ")
		if err == nil {
			sb.WriteString(string(indented))
			sb.WriteString("\n")
			return sb.String()
		}
	}

	sb.WriteString(respBody)
	sb.WriteString("\n")
	return sb.String()
}
