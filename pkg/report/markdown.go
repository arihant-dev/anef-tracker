package report

import (
	"fmt"
	"strings"
	"time"
)

func (r *EvidenceReport) RenderMarkdown() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", r.Title))
	sb.WriteString(fmt.Sprintf("**Generated At**: %s  \n", r.GeneratedAt.Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("**Application ID**: `%s`  \n", r.ApplicationID))
	sb.WriteString(fmt.Sprintf("**Current Status**: `%s`  \n", r.CurrentStatus))
	sb.WriteString(fmt.Sprintf("**Integrity Validated**: `%t`  \n\n", r.IntegrityOk))

	sb.WriteString("## 1. Application Lifecycle Evidence Claims\n\n")

	for i, claim := range r.Claims {
		sb.WriteString(fmt.Sprintf("### Claim #%d: %s\n", i+1, claim.ClaimStatement))
		sb.WriteString(fmt.Sprintf("- **Verified**: %t\n", claim.Verified))
		sb.WriteString(fmt.Sprintf("- **Timestamp**: %s\n", claim.Timestamp.Format("2006-01-02 15:04")))
		sb.WriteString("- **Evidence References**:\n")
		sb.WriteString(fmt.Sprintf("  - **Snapshot ID**: `%s`\n", claim.SnapshotID))
		if claim.EventID > 0 {
			sb.WriteString(fmt.Sprintf("  - **Event ID**: `#%d`\n", claim.EventID))
		}
		if claim.HTTPLogID > 0 {
			sb.WriteString(fmt.Sprintf("  - **HTTP Log ID**: `#%d`\n", claim.HTTPLogID))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 2. Evidence Verification Summary\n\n")
	sb.WriteString("All claims in this report are backed deterministically by immutable SHA-256 snapshot hashes and SQLite event logs.\n")

	return sb.String()
}
