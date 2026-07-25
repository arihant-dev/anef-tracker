package report

import (
	"fmt"
	"strings"
	"time"
)

func (r *EvidenceReport) RenderHTML() string {
	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title>ANEF Evidence Report</title>\n")
	sb.WriteString("<style>body{font-family:sans-serif;margin:40px;} h1{color:#1a365d;} .claim{background:#f7fafc;border-left:4px solid #3182ce;padding:15px;margin-bottom:15px;}</style>\n")
	sb.WriteString("</head>\n<body>\n")
	sb.WriteString(fmt.Sprintf("<h1>%s</h1>\n", r.Title))
	sb.WriteString(fmt.Sprintf("<p><strong>Application ID:</strong> %s | <strong>Status:</strong> %s</p>\n", r.ApplicationID, r.CurrentStatus))
	sb.WriteString(fmt.Sprintf("<p><strong>Generated:</strong> %s</p>\n<hr>\n", r.GeneratedAt.Format(time.RFC1123)))

	for i, claim := range r.Claims {
		sb.WriteString("<div class=\"claim\">\n")
		sb.WriteString(fmt.Sprintf("<h3>Claim #%d: %s</h3>\n", i+1, claim.ClaimStatement))
		sb.WriteString(fmt.Sprintf("<p>Snapshot: <code>%s</code> | Verified: %t</p>\n", claim.SnapshotID, claim.Verified))
		sb.WriteString("</div>\n")
	}

	sb.WriteString("</body>\n</html>")
	return sb.String()
}
