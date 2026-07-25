package privacy

import (
	"strings"
)

type PIIType string

const (
	PIIName    PIIType = "PII_NAME"
	PIIPhone   PIIType = "PII_PHONE"
	PIIAddress PIIType = "PII_ADDRESS"
	PIIID      PIIType = "PII_ID"
	PIIEmail   PIIType = "PII_EMAIL"
	PIIDOB     PIIType = "PII_DOB"
)

type SensitiveField struct {
	FieldPath string  `json:"field_path"`
	Type      PIIType `json:"type"`
}

func ScanPayload(payload map[string]interface{}) []SensitiveField {
	var fields []SensitiveField
	walkPayload("", payload, &fields)
	return fields
}

func walkPayload(prefix string, data interface{}, results *[]SensitiveField) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			fullPath := k
			if prefix != "" {
				fullPath = prefix + "." + k
			}
			classifyField(fullPath, k, results)
			walkPayload(fullPath, val, results)
		}
	case []interface{}:
		for _, item := range v {
			walkPayload(prefix, item, results)
		}
	}
}

func classifyField(fullPath, key string, results *[]SensitiveField) {
	lowerKey := strings.ToLower(key)
	var piiType PIIType

	switch {
	case strings.Contains(lowerKey, "telephone") || strings.Contains(lowerKey, "phone") || strings.Contains(lowerKey, "mobile"):
		piiType = PIIPhone
	case strings.Contains(lowerKey, "nom") || strings.Contains(lowerKey, "prenom") || strings.Contains(lowerKey, "lastname") || strings.Contains(lowerKey, "firstname"):
		piiType = PIIName
	case strings.Contains(lowerKey, "adresse") || strings.Contains(lowerKey, "street") || strings.Contains(lowerKey, "postal"):
		piiType = PIIAddress
	case strings.Contains(lowerKey, "numetranger") || strings.Contains(lowerKey, "foreignerid") || strings.Contains(lowerKey, "numeroetranger"):
		piiType = PIIID
	case strings.Contains(lowerKey, "email") || strings.Contains(lowerKey, "courriel"):
		piiType = PIIEmail
	case strings.Contains(lowerKey, "datenaissance") || strings.Contains(lowerKey, "dob") || strings.Contains(lowerKey, "birthdate"):
		piiType = PIIDOB
	default:
		return
	}

	*results = append(*results, SensitiveField{
		FieldPath: fullPath,
		Type:      piiType,
	})
}
