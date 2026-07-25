package privacy

import (
	"fmt"
	"strings"
)

type PrivacyPolicy struct {
	RedactExports          bool     `json:"redact_exports" yaml:"redact_exports"`
	SensitiveFieldPatterns []string `json:"sensitive_field_patterns" yaml:"sensitive_field_patterns"`
	RetentionDays          int      `json:"retention_days" yaml:"retention_days"`
}

func DefaultPolicy() PrivacyPolicy {
	return PrivacyPolicy{
		RedactExports: true,
		SensitiveFieldPatterns: []string{
			"telephone", "nom", "prenom", "adresse", "numEtranger", "email", "dateNaissance",
		},
		RetentionDays: 365,
	}
}

func FormatPolicySummary(p PrivacyPolicy) string {
	var sb strings.Builder
	sb.WriteString("=== PRIVACY & DATA PROTECTION POLICY ===\n\n")
	sb.WriteString(fmt.Sprintf("Export Redaction Active:  %t\n", p.RedactExports))
	sb.WriteString(fmt.Sprintf("HTTP Traffic Retention:  %d days\n", p.RetentionDays))
	sb.WriteString(fmt.Sprintf("Monitored PII Patterns:  %s\n", strings.Join(p.SensitiveFieldPatterns, ", ")))
	return sb.String()
}
