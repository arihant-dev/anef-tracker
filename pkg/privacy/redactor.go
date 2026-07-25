package privacy

import (
	"encoding/json"
	"strings"
)

func RedactString(piiType PIIType, val string) string {
	if len(val) == 0 {
		return val
	}
	switch piiType {
	case PIIPhone:
		if len(val) >= 4 {
			return val[:2] + "*****" + val[len(val)-2:]
		}
		return "**********"
	case PIIName:
		return string(val[0]) + "***"
	case PIIAddress:
		return "[REDACTED]"
	case PIIID:
		if len(val) >= 4 {
			return val[:2] + "**********" + val[len(val)-2:]
		}
		return "**********"
	case PIIEmail:
		parts := strings.Split(val, "@")
		if len(parts) == 2 {
			return string(parts[0][0]) + "***@" + parts[1]
		}
		return "[REDACTED]"
	default:
		return "[REDACTED]"
	}
}

func RedactPayload(payload map[string]interface{}) map[string]interface{} {
	// Deep copy via JSON
	data, _ := json.Marshal(payload)
	var copyMap map[string]interface{}
	_ = json.Unmarshal(data, &copyMap)

	sensitiveFields := ScanPayload(copyMap)
	for _, sf := range sensitiveFields {
		applyRedactionToPath(copyMap, sf.FieldPath, sf.Type)
	}

	return copyMap
}

func applyRedactionToPath(m map[string]interface{}, path string, piiType PIIType) {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		if valStr, ok := m[parts[0]].(string); ok {
			m[parts[0]] = RedactString(piiType, valStr)
		}
		return
	}

	if subMap, ok := m[parts[0]].(map[string]interface{}); ok {
		applyRedactionToPath(subMap, strings.Join(parts[1:], "."), piiType)
	}
}
