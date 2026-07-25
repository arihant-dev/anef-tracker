package tui

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/analytics"
	"github.com/arihant-dev/anef-tracker/pkg/audit"
	"github.com/arihant-dev/anef-tracker/pkg/crawler"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/diff"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"github.com/arihant-dev/anef-tracker/pkg/profile"
	"github.com/arihant-dev/anef-tracker/pkg/schema"
	"github.com/arihant-dev/anef-tracker/pkg/security"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
	"github.com/arihant-dev/anef-tracker/pkg/timeline"
	"github.com/arihant-dev/anef-tracker/pkg/workflow"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func renderOverviewLines(app *domain.Application) []string {
	if app == nil {
		return []string{"No active application data loaded. Run 'anef fetch' or 'anef login' first."}
	}

	badge := BadgeLow.Render("LOW")
	switch app.Status.Severity {
	case domain.SeverityMedium:
		badge = BadgeMedium.Render("MEDIUM")
	case domain.SeverityHigh:
		badge = BadgeHigh.Render("HIGH")
	case domain.SeverityCritical:
		badge = BadgeCritical.Render("CRITICAL")
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Application Number: %s", app.NumeroDemande))
	lines = append(lines, fmt.Sprintf("Foreigner ID:       %s", app.ForeignerID))
	lines = append(lines, fmt.Sprintf("Legal Category:     %s", app.LegalCategory))
	lines = append(lines, fmt.Sprintf("Current Status:     %s  %s", app.Status.Label, badge))
	lines = append(lines, fmt.Sprintf("Status Code:        %s", app.Status.Code))
	lines = append(lines, fmt.Sprintf("Description:        %s", app.Status.Description))
	lines = append(lines, fmt.Sprintf("Processing Site:    %s", app.ProcessingSite))
	lines = append(lines, fmt.Sprintf("Collection Site:    %s", app.CollectionSite))
	lines = append(lines, fmt.Sprintf("Version Counter:    v%d", app.Version))
	lines = append(lines, fmt.Sprintf("Last Updated:       %s", app.UpdatedAt.Format(time.RFC1123)))

	if len(app.Documents) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Available Documents Inventory:")
		for _, doc := range app.Documents {
			lines = append(lines, fmt.Sprintf("  ✓ [%s] %s (%s)", doc.Type, doc.Name, doc.URL))
		}
	}

	return lines
}

func renderTimelineLines(app *domain.Application) []string {
	var lines []string
	lines = append(lines, "=== APPLICATION LIFECYCLE TIMELINE ===")
	lines = append(lines, "")

	milestones := []struct {
		code  string
		label string
	}{
		{"DEMANDE_SOUMISE", "1. Application Submitted"},
		{"DOSSIER_DEPOSE", "2. Dossier Deposited"},
		{"VALIDATION_FORMELLE", "3. Formal Validation"},
		{"INSTRUCTION_EN_COURS", "4. Instruction Started"},
		{"DECISION_VALIDEE", "5. Decision Validated"},
		{"TITRE_A_FABRIQUER", "6. Residence Permit in Production"},
		{"TITRE_DISPONIBLE", "7. Ready for Collection"},
	}

	currentCode := "DEMANDE_SOUMISE"
	snapDate := time.Now()
	if app != nil {
		currentCode = app.Status.Code
		snapDate = app.UpdatedAt
	}

	foundActive := false
	for i, m := range milestones {
		prefix := "[ ]"
		style := lipgloss.NewStyle().Foreground(ColorMuted)

		if m.code == currentCode {
			prefix = "[★]"
			style = lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary)
			foundActive = true
		} else if !foundActive {
			prefix = "[✓]"
			style = lipgloss.NewStyle().Foreground(ColorSuccess)
		}

		lines = append(lines, style.Render(fmt.Sprintf("%s %s", prefix, m.label)))
		if i < len(milestones)-1 {
			lines = append(lines, lipgloss.NewStyle().Foreground(ColorMuted).Render("    │"))
		}
	}

	elapsed := time.Since(snapDate)
	days := int(elapsed.Hours() / 24)
	hours := int(elapsed.Hours()) % 24
	mins := int(elapsed.Minutes()) % 60

	statusLabel := currentCode
	if app != nil && app.Status.Label != "" {
		statusLabel = app.Status.Label
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Current State: %s [%s]", statusLabel, currentCode))
	lines = append(lines, fmt.Sprintf("Elapsed in state: %d days %02d hours %02d minutes", days, hours, mins))

	return lines
}

func renderDiffLines() []string {
	latest, previous, err := snapshot.GetLatestTwoSnapshots()
	if err != nil || latest == nil {
		return []string{"No snapshot history available to compute diffs. Run 'anef fetch' twice."}
	}

	if previous == nil {
		return []string{fmt.Sprintf("Only 1 snapshot recorded (%s). A second snapshot is required for diff comparison.", latest.Directory)}
	}

	res, err := diff.CompareSnapshots(previous.RawBytes, latest.RawBytes)
	if err != nil {
		return []string{fmt.Sprintf("Diff error: %v", err)}
	}

	return strings.Split(res.Summary, "\n")
}

func renderEventsLines(database *db.DB, eventFilter string) []string {
	if database == nil {
		return []string{"Database not connected."}
	}

	events, err := database.GetEvents(100)
	if err != nil || len(events) == 0 {
		return []string{"No historical events recorded yet."}
	}

	var lines []string
	title := "=== RECORDED STATE EVENTS (Filter: ALL) ==="
	switch eventFilter {
	case "FIELD_DISCOVERED":
		title = "=== RECORDED STATE EVENTS (Filter: FIELD DISCOVERIES) ==="
	case "STATUS_CHANGE":
		title = "=== RECORDED STATE EVENTS (Filter: STATUS CHANGES) ==="
	}
	lines = append(lines, title)
	lines = append(lines, "Filters: [a] All  [f] Field Discoveries  [s] Status Changes")
	lines = append(lines, "")

	filteredCount := 0
	for _, ev := range events {
		if eventFilter == "FIELD_DISCOVERED" && ev.Type != "FIELD_DISCOVERED" {
			continue
		}
		if eventFilter == "STATUS_CHANGE" && ev.Type == "FIELD_DISCOVERED" {
			continue
		}

		filteredCount++
		sevBadge := BadgeLow.Render(string(ev.Severity))
		switch ev.Severity {
		case domain.SeverityHigh:
			sevBadge = BadgeHigh.Render(string(ev.Severity))
		case domain.SeverityCritical:
			sevBadge = BadgeCritical.Render(string(ev.Severity))
		}

		lines = append(lines, fmt.Sprintf("[%s] %s | %s (Confidence: %.2f)", ev.Timestamp.Format("15:04:05"), ev.Type, sevBadge, ev.Confidence))
		lines = append(lines, fmt.Sprintf("  Field: %s", ev.FieldPath))
		lines = append(lines, fmt.Sprintf("  %s → %s", ev.OldVal, ev.NewVal))
		lines = append(lines, "")
	}

	if filteredCount == 0 {
		lines = append(lines, "No events matching active filter.")
	}

	return lines
}

func renderDocumentsLines(app *domain.Application) []string {
	if app == nil || len(app.Documents) == 0 {
		return []string{"No official documents or attestations currently attached."}
	}

	var lines []string
	lines = append(lines, "=== OFFICIAL DOCUMENT INVENTORY ===")
	lines = append(lines, "")

	for i, doc := range app.Documents {
		lines = append(lines, fmt.Sprintf("%d. ✓ [%s] %s", i+1, doc.Type, doc.Name))
		lines = append(lines, fmt.Sprintf("   URL: %s", doc.URL))
		lines = append(lines, fmt.Sprintf("   Generated: %s", doc.CreatedAt.Format("02 Jan 2006 15:04")))
		lines = append(lines, "")
	}

	return lines
}

func renderAPIExplorerLines(database *db.DB) []string {
	if database == nil {
		return []string{"Database not connected."}
	}

	crl := crawler.NewCrawler(database)
	obs, err := crl.DiscoverEndpoints()
	if err != nil || len(obs) == 0 {
		return []string{"No endpoints discovered yet. Run 'anef fetch' first."}
	}

	summaryStr := crawler.FormatEndpointSummary(obs)
	return strings.Split(summaryStr, "\n")
}

func renderSchemaLines(database *db.DB, searchQuery string) []string {
	if database == nil {
		return []string{"Database not connected."}
	}

	reg := schema.NewRegistry(database)
	fields, err := reg.ListFields("")
	if err != nil || len(fields) == 0 {
		return []string{"No schema fields registered yet. Run 'anef fetch' first."}
	}

	var lines []string
	lines = append(lines, "=== DISCOVERED SCHEMA FIELDS ===")
	if searchQuery != "" {
		lines = append(lines, fmt.Sprintf("Search Query: '%s' (Press '/' to change, Esc to clear)", searchQuery))
	} else {
		lines = append(lines, "Press [/] to search fields")
	}
	lines = append(lines, "")

	matchCount := 0
	for _, f := range fields {
		if searchQuery != "" && !strings.Contains(strings.ToLower(f.Path), strings.ToLower(searchQuery)) {
			continue
		}
		matchCount++
		lines = append(lines, fmt.Sprintf("%d. %s [%s]", matchCount, f.Path, f.Type))
		lines = append(lines, fmt.Sprintf("   Occurrences: %d | Confidence: %.2f | First Seen: %s", f.Occurrences, f.Confidence, f.FirstSeen.Format("2006-01-02 15:04")))
		lines = append(lines, "")
	}

	if matchCount == 0 {
		lines = append(lines, fmt.Sprintf("No schema fields matched query '%s'.", searchQuery))
	}

	return lines
}

func renderReplayLines(database *db.DB) []string {
	if database == nil {
		return []string{"Database not connected."}
	}

	rows, err := database.Conn.Query("SELECT id, original_request_id, timestamp, status_code, matched FROM http_replays ORDER BY id DESC LIMIT 50")
	if err != nil {
		return []string{"No HTTP replays recorded yet. Use 'anef replay <id>' CLI command to trigger a request replay."}
	}
	defer rows.Close()

	var lines []string
	lines = append(lines, "=== HTTP REPLAY ENGINE ===")
	lines = append(lines, "")

	count := 0
	for rows.Next() {
		count++
		var id, origID int64
		var statusCode int
		var matched bool
		var ts time.Time
		_ = rows.Scan(&id, &origID, &ts, &statusCode, &matched)

		matchStr := "MATCH [YES]"
		if !matched {
			matchStr = "MISMATCH [NO]"
		}

		lines = append(lines, fmt.Sprintf("Replay #%d | Orig Request #%d | Status: HTTP %d | Result: %s | Time: %s",
			id, origID, statusCode, matchStr, ts.Format("15:04:05")))
	}
	if err := rows.Err(); err != nil {
		return []string{fmt.Sprintf("Replay rows iteration error: %v", err)}
	}

	if count == 0 {
		lines = append(lines, "No replays executed yet. Run 'anef replay 1' to test response matching.")
	}

	return lines
}

func renderHTTPLogsLines(database *db.DB) []string {
	if database == nil {
		return []string{"Database not connected."}
	}

	rows, err := database.Conn.Query("SELECT id, method, url, status_code, latency_ms, created_at FROM http_logs ORDER BY id DESC LIMIT 50")
	if err != nil {
		return []string{fmt.Sprintf("Failed querying HTTP logs: %v", err)}
	}
	defer rows.Close()

	var lines []string
	lines = append(lines, "=== RECENT HTTP TRAFFIC LOGS ===")
	lines = append(lines, "")

	count := 0
	for rows.Next() {
		count++
		var id int64
		var method, rawURL string
		var statusCode, latencyMs int64
		var createdAt time.Time
		_ = rows.Scan(&id, &method, &rawURL, &statusCode, &latencyMs, &createdAt)

		shortURL := rawURL
		if idx := strings.Index(rawURL, "/api/"); idx != -1 {
			shortURL = rawURL[idx:]
		}

		statusStr := fmt.Sprintf("HTTP %d", statusCode)
		statusStyle := lipgloss.NewStyle().Foreground(ColorSuccess)
		if statusCode >= 400 {
			statusStyle = lipgloss.NewStyle().Foreground(ColorDanger)
		}

		lines = append(lines, fmt.Sprintf("#%d [%s] %s %s | %s (%dms)",
			id, createdAt.Format("15:04:05"), method, shortURL, statusStyle.Render(statusStr), latencyMs))
	}
	if err := rows.Err(); err != nil {
		return []string{fmt.Sprintf("HTTP logs iteration error: %v", err)}
	}

	if count == 0 {
		lines = append(lines, "No HTTP traffic recorded yet.")
	}

	return lines
}

func renderWorkflowLines(database *db.DB, app *domain.Application) []string {
	sm := workflow.NewStateMachine(database)
	curCode := "TITRE_A_FABRIQUER"
	if app != nil && app.Status.Code != "" {
		curCode = app.Status.Code
	}

	asciiStr := sm.RenderASCII(curCode)
	return strings.Split(asciiStr, "\n")
}

func renderAnalyticsLines(database *db.DB) []string {
	ae := analytics.NewAnalyticsEngine(database)
	summaryStr := ae.FormatStatisticsSummary()
	return strings.Split(summaryStr, "\n")
}

func renderKnowledgeGraphLines(database *db.DB) []string {
	gb := knowledge.NewGraphBuilder(database)
	g, err := gb.BuildFromDB()
	if err != nil {
		return []string{fmt.Sprintf("Failed building knowledge graph: %v", err)}
	}

	return strings.Split(g.FormatASCII(), "\n")
}

func renderHumanTimelineLines(database *db.DB) []string {
	ht, err := timeline.BuildTimeline(database)
	if err != nil {
		return []string{fmt.Sprintf("Failed building human timeline: %v", err)}
	}
	return strings.Split(ht.FormatASCII(), "\n")
}

func renderNotificationsLines(database *db.DB) []string {
	var lines []string
	lines = append(lines, "=== NOTIFICATION LOG ===")
	lines = append(lines, "")

	if database == nil {
		lines = append(lines, "No database available. Notification history unavailable.")
		return lines
	}

	rows, err := database.Conn.Query(
		"SELECT id, type, title, message, delivered, created_at FROM notifications ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		lines = append(lines, "No notifications recorded yet.")
		lines = append(lines, "")
		lines = append(lines, "Configure notifications:")
		lines = append(lines, "  anef notify configure --webhook <url>")
		lines = append(lines, "  anef notify configure --telegram <bot_token> <chat_id>")
		lines = append(lines, "  anef notify test")
		return lines
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var nType, title, message string
		var delivered bool
		var createdAt time.Time
		if err := rows.Scan(&id, &nType, &title, &message, &delivered, &createdAt); err != nil {
			continue
		}
		count++

		statusIcon := "○"
		if delivered {
			statusIcon = "●"
		}

		lines = append(lines, fmt.Sprintf("%s #%d [%s] %s — %s",
			statusIcon, id, nType, title, createdAt.Format("2006-01-02 15:04")))
		lines = append(lines, fmt.Sprintf("    %s", message))
		lines = append(lines, "")
	}

	if err := rows.Err(); err != nil {
		return []string{fmt.Sprintf("Notifications iteration error: %v", err)}
	}

	if count == 0 {
		lines = append(lines, "No notifications recorded yet.")
		lines = append(lines, "")
		lines = append(lines, "Configure notifications:")
		lines = append(lines, "  anef notify configure --webhook <url>")
		lines = append(lines, "  anef notify configure --telegram <bot_token> <chat_id>")
		lines = append(lines, "  anef notify test")
	}

	return lines
}

func renderProfilesLines(database *db.DB) []string {
	profiles, err := profile.ListProfiles(database)
	if err != nil {
		return []string{fmt.Sprintf("Failed listing profiles: %v", err)}
	}
	return strings.Split(profile.FormatProfileList(profiles), "\n")
}

func renderSecurityStatusLines(_ *db.DB) []string {
	rep := security.AuditSecurity()
	return strings.Split(security.FormatSecurityReport(rep), "\n")
}

func renderAuditLogLines(database *db.DB) []string {
	entries, err := audit.ListAuditLog(database, 50)
	if err != nil {
		return []string{fmt.Sprintf("Failed listing audit log: %v", err)}
	}
	return strings.Split(audit.FormatAuditLog(entries), "\n")
}
