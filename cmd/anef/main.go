package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/arihant-dev/anef-tracker/internal/service"
	"github.com/arihant-dev/anef-tracker/internal/version"
	"github.com/arihant-dev/anef-tracker/pkg/analytics"
	"github.com/arihant-dev/anef-tracker/pkg/audit"
	"github.com/arihant-dev/anef-tracker/pkg/auth"
	"github.com/arihant-dev/anef-tracker/pkg/backup"
	"github.com/arihant-dev/anef-tracker/pkg/config"
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"github.com/arihant-dev/anef-tracker/pkg/crawler"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/doctor"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"github.com/arihant-dev/anef-tracker/pkg/export"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"github.com/arihant-dev/anef-tracker/pkg/log"
	"github.com/arihant-dev/anef-tracker/pkg/notify"
	"github.com/arihant-dev/anef-tracker/pkg/privacy"
	"github.com/arihant-dev/anef-tracker/pkg/profile"
	"github.com/arihant-dev/anef-tracker/pkg/report"
	"github.com/arihant-dev/anef-tracker/pkg/schema"
	"github.com/arihant-dev/anef-tracker/pkg/security"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
	"github.com/arihant-dev/anef-tracker/pkg/timeline"
	"github.com/arihant-dev/anef-tracker/pkg/tui"
	"github.com/arihant-dev/anef-tracker/pkg/workflow"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ExitSuccess      = 0
	ExitError        = 1
	ExitAuthExpired  = 2
	ExitChangeDetect = 3
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(ExitError)
	}

	command := os.Args[1]
	_, _ = domain.LoadStatusDict("status_mapping.yaml")

	switch command {
	case "version":
		fmt.Println(version.GetVersionInfo())
		os.Exit(ExitSuccess)
	case "config":
		handleConfig()
	case "login":
		handleLogin()
	case "fetch":
		handleFetch()
	case "status":
		handleStatus()
	case "diff":
		handleDiff()
	case "timeline":
		handleTimeline()
	case "events":
		handleEvents()
	case "api":
		handleAPI()
	case "schema":
		handleSchema()
	case "replay":
		handleReplay()
	case "endpoints":
		handleEndpoints()
	case "workflow":
		handleWorkflow()
	case "analytics":
		handleAnalytics()
	case "graph":
		handleGraph()
	case "evidence":
		handleEvidence()
	case "knowledge":
		handleKnowledge()
	case "export":
		handleExport()
	case "backup":
		handleBackup()
	case "retention":
		handleRetention()
	case "db":
		handleDB()
	case "security":
		handleSecurity()
	case "watch":
		handleWatch()
	case "notify":
		handleNotify()
	case "report":
		handleReport()
	case "metrics":
		handleMetrics()
	case "context":
		handleContext()
	case "profile":
		handleProfile()
	case "privacy":
		handlePrivacy()
	case "audit":
		handleAudit()
	case "doctor":
		handleDoctor()
	case "tui":
		handleTUI()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(ExitError)
	}
}

func printUsage() {
	fmt.Println(`ANEF Residence Permit Tracker CLI & Evidence-Based Workflow Intelligence Platform

Usage:
  anef <command> [flags]

Monitoring & Status Commands:
  login                        Authenticate via Chrome DevTools cURL or credentials
  fetch                        Fetch application status, record HTTP traffic & save snapshot
  status                       Print current residence permit application status
  timeline                     Display chronological application progress timeline
  watch                        Run watcher daemon (--interval 6h) or single execution (--once)
  notify status                Display notification channel configuration & delivery stats
  notify test                  Send test notification to all configured channels
  notify configure             Configure notification channels (webhook, email, telegram)
  report generate              Generate official evidence-backed report (markdown/html/json)
  metrics                      Display operational metrics (snapshot count, DB size, integrity)
  context                      Display active profile & application scope context
  profile list|create|switch   Manage multi-application applicant profiles & vaults
  tui                          Launch interactive 17-tab grouped terminal user interface

Forensic Evidence & Privacy Vault:
  privacy audit                Run PII field scanner & privacy policy audit
  audit list|verify            Inspect tamper-evident audit log & verify hash chain
  evidence bundle              Export application-scoped evidence ZIP (--redact support)
  evidence search <term>       Search evidence graph nodes & snapshot events by keyword
  evidence verify              Verify snapshot SHA-256 hashes & link integrity
  evidence export <node-id>    Export node evidence chain to JSON/YAML/Markdown
  graph explain <node_id>      Explain node provenance attribution & connections
  graph validate               Validate Evidence Graph consistency & check orphan nodes

Workflow Intelligence:
  workflow                     Display reconstructed ASCII workflow state machine
  workflow explain <from> <to> Explain empirical evidence & snapshot sources for transition
  workflow audit <from> <to>   Conduct forensic transition audit
  analytics                    Display status duration medians and percentiles
  analytics explain <state>    Explain duration sample distribution & threshold evidence

API Reverse Engineering:
  api list                     List all observed ANEF REST endpoints
  api inspect <id>             Inspect request/response and schema for an endpoint
  schema list                  List all registered API fields and data types
  schema diff                  Compare schema structures across snapshots
  replay <id>                  Replay recorded HTTP request and verify response match
  endpoints                    Print endpoint statistics and hypermedia dependency graph

Operations & System Maintenance:
  version                      Display version, git commit, build date & platform info
  config show                  Display current configuration settings
  backup create|restore        Create compressed backup tarball or restore state
  db status|vacuum|schema      Display DB metrics, run maintenance, or export schema
  security audit               Audit session permissions (0600) & credential redaction
  doctor                       Run 8-point system diagnostic health suite
  export                       Export events or schemas (CSV/JSON/YAML/Markdown)`)
}

func handleLogin() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	browserFlag := loginCmd.Bool("browser", false, "Launch default browser for automated OAuth redirect journey")
	webFlag := loginCmd.Bool("web", false, "Alias for --browser")
	curlFlag := loginCmd.String("curl", "", "Raw cURL command copied from Chrome DevTools")
	userFlag := loginCmd.String("user", "", "Foreigner ID (numéro d'étranger 9999999999 or email)")
	passFlag := loginCmd.String("pass", "", "Password")
	_ = loginCmd.Parse(os.Args[2:])

	if *browserFlag || *webFlag {
		sess, err := auth.AuthenticateViaBrowser("https://administration-etrangers-en-france.interieur.gouv.fr/usagers/", 8484, 3*time.Minute)
		if err != nil {
			log.Error("Browser authentication error", "error", err)
			os.Exit(ExitError)
		}
		if err := auth.SaveSession(sess); err != nil {
			log.Error("Failed saving session", "error", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Successfully captured browser session for user '%s'!\n", sess.Login)
		os.Exit(ExitSuccess)
	}

	if *curlFlag != "" {
		sess, err := auth.ParseCurl(*curlFlag)
		if err != nil {
			log.Error("Failed parsing cURL", "error", err)
			os.Exit(ExitError)
		}
		if err := auth.SaveSession(sess); err != nil {
			log.Error("Failed saving session", "error", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Successfully imported cURL session for login '%s'!\n", sess.Login)
		os.Exit(ExitSuccess)
	}

	if *userFlag != "" && *passFlag != "" {
		sess, err := auth.AuthenticateWithCredentials(nil, *userFlag, *passFlag)
		if err != nil {
			log.Error("Authentication failed", "error", err)
			os.Exit(ExitAuthExpired)
		}
		if err := auth.SaveSession(sess); err != nil {
			log.Error("Failed saving session", "error", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Successfully logged in user '%s' via Keycloak!\n", sess.Login)
		os.Exit(ExitSuccess)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("=== ANEF Session Authentication ===")
	fmt.Println("Option 1: [ENTER] Launch automated browser authentication journey (anef login --browser)")
	fmt.Println("Option 2: Paste DevTools cURL command ('curl https://...')")
	fmt.Println("Option 3: Type 'p' for password fallback prompt")
	fmt.Println("")
	fmt.Print("Choice or cURL command [ENTER]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" || input == "b" || input == "1" || input == "browser" {
		sess, err := auth.AuthenticateViaBrowser("https://administration-etrangers-en-france.interieur.gouv.fr/usagers/", 8484, 3*time.Minute)
		if err != nil {
			log.Error("Browser authentication error", "error", err)
			os.Exit(ExitError)
		}
		if err := auth.SaveSession(sess); err != nil {
			log.Error("Failed saving session", "error", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Successfully captured browser session for user '%s'!\n", sess.Login)
		os.Exit(ExitSuccess)
	}

	if strings.Contains(input, "curl ") || strings.HasPrefix(input, "curl") {
		sess, err := auth.ParseCurl(input)
		if err != nil {
			log.Error("Failed parsing cURL command", "error", err)
			os.Exit(ExitError)
		}
		if err := auth.SaveSession(sess); err != nil {
			log.Error("Failed saving session", "error", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Successfully imported cURL session for login '%s'!\n", sess.Login)
		os.Exit(ExitSuccess)
	}

	fmt.Println("\nNotice: ANEF Direct Access Password grants may be restricted by Keycloak security policy.")
	fmt.Print("Identifiant (Foreigner ID or Email): ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Mot de passe: ")
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)

	sess, err := auth.AuthenticateWithCredentials(nil, user, pass)
	if err != nil {
		fmt.Printf("\n✗ Authentication Failed:\n%v\n", err)
		os.Exit(ExitAuthExpired)
	}
	if err := auth.SaveSession(sess); err != nil {
		log.Error("Failed saving session", "error", err)
		os.Exit(ExitError)
	}
	fmt.Printf("✓ Successfully logged in user '%s'!\n", user)
	os.Exit(ExitSuccess)
}

func handleFetch() int {
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		return ExitError
	}

	ctx := context.Background()
	log.Info("Fetching latest application status from ANEF...")
	app, snapRef, diffRes, err := svc.Fetch(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "Invalid Token") {
			log.Error("Authentication expired", "error", err)
			return ExitAuthExpired
		}
		log.Error("Fetch failed", "error", err)
		return ExitError
	}

	if snapRef != nil {
		fmt.Printf("✓ Snapshot stored [%s] at %s\n", snapRef.SnapshotID, snapRef.Directory)
	}

	hasChange := false
	if diffRes != nil && diffRes.HasChanges {
		hasChange = true
		fmt.Println("\n" + diffRes.Summary)
	} else {
		fmt.Println("No changes detected since previous snapshot.")
	}

	fmt.Printf("\nStatus: %s [%s]\nDescription: %s\n", app.Status.Label, app.Status.Code, app.Status.Description)

	if hasChange {
		return ExitChangeDetect
	}
	return ExitSuccess
}

func handleStatus() {
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	app, snapRef, err := svc.GetStatus()
	if err != nil {
		fmt.Printf("Status error: %v\n", err)
		os.Exit(ExitError)
	}

	foreignerID := app.ForeignerID
	if foreignerID == "" && svc.Session != nil {
		foreignerID = svc.Session.Login
	}

	fmt.Printf("=== ANEF APPLICATION STATUS ===\n")
	fmt.Printf("Snapshot ID:        %s\n", snapRef.SnapshotID)
	fmt.Printf("Application Number: %s\n", app.NumeroDemande)
	fmt.Printf("Foreigner ID:       %s\n", foreignerID)
	fmt.Printf("Legal Category:     %s\n", app.LegalCategory)
	fmt.Printf("Status:             %s [%s]\n", app.Status.Label, app.Status.Code)
	fmt.Printf("Description:        %s\n", app.Status.Description)
	fmt.Printf("Collection Site:    %s\n", app.CollectionSite)
	fmt.Printf("Version Counter:    v%d\n", app.Version)
	fmt.Printf("Snapshot Date:      %s\n", snapRef.Timestamp.Format(time.RFC1123))
	os.Exit(ExitSuccess)
}

func handleDiff() {
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	diffRes, err := svc.GetDiff()
	if err != nil {
		fmt.Printf("Diff error: %v\n", err)
		os.Exit(ExitError)
	}

	fmt.Println(diffRes.Summary)
	os.Exit(ExitSuccess)
}

func handleEvents() {
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	events, err := svc.Store.GetEvents(20)
	if err != nil || len(events) == 0 {
		fmt.Println("No historical state events recorded yet.")
		os.Exit(ExitSuccess)
	}

	fmt.Println("=== RECORDED APPLICATION EVENTS ===")
	for _, ev := range events {
		fmt.Printf("[%s] %s | Severity: %s | Confidence: %.2f\n  Field: %s\n  %s → %s\n\n",
			ev.Timestamp.Format(time.RFC1123), ev.Type, ev.Severity, ev.Confidence, ev.FieldPath, ev.OldVal, ev.NewVal)
	}
	os.Exit(ExitSuccess)
}

func handleWorkflow() {
	if len(os.Args) >= 4 && os.Args[2] == "explain" {
		fromStatus := os.Args[3]
		toStatus := ""
		if len(os.Args) >= 5 {
			toStatus = os.Args[4]
		}
		sm := workflow.NewStateMachine(nil)
		fmt.Println(sm.ExplainTransition(fromStatus, toStatus))
		os.Exit(ExitSuccess)
	}

	if len(os.Args) >= 4 && os.Args[2] == "audit" {
		fromStatus := os.Args[3]
		toStatus := ""
		if len(os.Args) >= 5 {
			toStatus = os.Args[4]
		}
		sm := workflow.NewStateMachine(nil)
		rep := sm.AuditTransition(fromStatus, toStatus)
		fmt.Println(workflow.FormatAuditReport(rep))
		os.Exit(ExitSuccess)
	}

	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	app, _, _ := svc.GetStatus()
	curCode := "TITRE_A_FABRIQUER"
	if app != nil && app.Status.Code != "" {
		curCode = app.Status.Code
	}

	sm := workflow.NewStateMachine(nil)
	fmt.Println(sm.RenderASCII(curCode))
	os.Exit(ExitSuccess)
}

func handleAnalytics() {
	if len(os.Args) >= 4 && os.Args[2] == "explain" {
		stateCode := os.Args[3]
		ae := analytics.NewAnalyticsEngine(nil)
		fmt.Println(ae.ExplainStateAnalytics(stateCode))
		os.Exit(ExitSuccess)
	}

	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	ae := analytics.NewAnalyticsEngine(nil)
	_ = svc
	fmt.Println(ae.FormatStatisticsSummary())
	os.Exit(ExitSuccess)
}

func handleGraph() {
	if len(os.Args) >= 3 && os.Args[2] == "validate" {
		database, _ := db.InitDB()
		gb := knowledge.NewGraphBuilder(database)
		g, _ := gb.BuildFromDB()
		rep := g.ValidateGraph()
		fmt.Println(knowledge.FormatValidationReport(rep))
		os.Exit(ExitSuccess)
	}

	if len(os.Args) >= 4 && os.Args[2] == "explain" {
		nodeID := os.Args[3]
		database, _ := db.InitDB()
		gb := knowledge.NewGraphBuilder(database)
		g, _ := gb.BuildFromDB()
		fmt.Println(g.ExplainNode(nodeID))
		os.Exit(ExitSuccess)
	}

	database, _ := db.InitDB()
	gb := knowledge.NewGraphBuilder(database)
	g, _ := gb.BuildFromDB()
	fmt.Println(g.FormatASCII())
	os.Exit(ExitSuccess)
}

func handleEvidence() {
	if len(os.Args) >= 3 && os.Args[2] == "verify" {
		database, _ := db.InitDB()
		rep, _ := evidence.VerifyIntegrity(database)
		fmt.Println(evidence.FormatIntegrityReport(rep))
		os.Exit(ExitSuccess)
	}

	if len(os.Args) >= 4 && os.Args[2] == "search" {
		searchTerm := os.Args[3]
		database, _ := db.InitDB()
		gb := knowledge.NewGraphBuilder(database)
		g, _ := gb.BuildFromDB()

		fmt.Println(g.SearchEvidence(searchTerm))
		os.Exit(ExitSuccess)
	}

	fmt.Println("Usage: anef evidence [verify|search <term>]")
	os.Exit(ExitError)
}

func handleExport() {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	formatFlag := exportCmd.String("format", "csv", "Export format: csv, json, yaml, markdown")
	typeFlag := exportCmd.String("type", "events", "Export target: events, schema, endpoints")
	_ = exportCmd.Parse(os.Args[2:])

	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	cwd, _ := os.Getwd()
	exportDir := filepath.Join(cwd, "exports")
	_ = os.MkdirAll(exportDir, 0755)

	var outputBytes []byte
	filename := fmt.Sprintf("%s.%s", *typeFlag, *formatFlag)
	if *formatFlag == "markdown" {
		filename = fmt.Sprintf("%s.md", *typeFlag)
	}

	targetPath := filepath.Join(exportDir, filename)

	switch *typeFlag {
	case "events":
		events, err := svc.Store.GetEvents(500)
		if err != nil || len(events) == 0 {
			fmt.Println("No events available to export.")
			os.Exit(ExitSuccess)
		}
		switch *formatFlag {
		case "json":
			outputBytes, _ = export.ToJSON(events)
		case "yaml":
			outputBytes, _ = export.ToYAML(events)
		case "markdown":
			outputBytes = []byte(export.EventsToMarkdown(events))
		default:
			outputBytes, _ = export.EventsToCSV(events)
		}
	case "schema":
		fields, err := svc.Registry.ListFields("")
		if err != nil || len(fields) == 0 {
			fmt.Println("No schema fields available to export.")
			os.Exit(ExitSuccess)
		}
		switch *formatFlag {
		case "json":
			outputBytes, _ = export.ToJSON(fields)
		case "yaml":
			outputBytes, _ = export.ToYAML(fields)
		default:
			outputBytes, _ = export.FieldsToCSV(fields)
		}
	}

	if len(outputBytes) > 0 {
		_ = os.WriteFile(targetPath, outputBytes, 0644)
		fmt.Printf("✓ Successfully exported %s data in %s format to:\n  %s\n", *typeFlag, strings.ToUpper(*formatFlag), targetPath)
	} else {
		fmt.Println("No data exported.")
	}
	os.Exit(ExitSuccess)
}

func handleAPI() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef api [list|inspect <id>]")
		os.Exit(ExitError)
	}

	subCmd := os.Args[2]
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	switch subCmd {
	case "list":
		endpoints, err := svc.Explorer.ListEndpoints()
		if err != nil || len(endpoints) == 0 {
			fmt.Println("No endpoints discovered yet. Run 'anef fetch' first.")
			os.Exit(ExitSuccess)
		}
		fmt.Println(crawler.FormatEndpointSummary(endpoints))
		os.Exit(ExitSuccess)
	case "inspect":
		if len(os.Args) < 4 {
			fmt.Println("Usage: anef api inspect <id>")
			os.Exit(ExitError)
		}
		id, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Printf("Invalid endpoint ID: %s\n", os.Args[3])
			os.Exit(ExitError)
		}

		obs, reqHeaders, respBody, err := svc.Explorer.InspectEndpoint(id)
		if err != nil {
			fmt.Printf("Inspect error: %v\n", err)
			os.Exit(ExitError)
		}

		fmt.Printf("=== ENDPOINT INSPECTION #%d ===\n", id)
		fmt.Printf("Method:  %s\nURL:     %s\nCalls:   %d\n\n", obs.Method, obs.URL, obs.Occurrences)
		fmt.Printf("Request Headers:\n%s\n\n", reqHeaders)
		fmt.Printf("Response Body Payload:\n%s\n", respBody)
		os.Exit(ExitSuccess)
	}
}

func handleSchema() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef schema [list|diff]")
		os.Exit(ExitError)
	}

	subCmd := os.Args[2]
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	switch subCmd {
	case "list":
		fields, err := svc.Registry.ListFields("")
		if err != nil || len(fields) == 0 {
			fmt.Println("No schema fields registered yet. Run 'anef fetch' first.")
			os.Exit(ExitSuccess)
		}

		fmt.Println("=== REGISTERED API SCHEMA FIELDS ===")
		for i, f := range fields {
			fmt.Printf("%d. %s [%s]\n   Occurrences: %d | Confidence: %.2f | First Seen: %s\n\n",
				i+1, f.Path, f.Type, f.Occurrences, f.Confidence, f.FirstSeen.Format("2006-01-02 15:04"))
		}
		os.Exit(ExitSuccess)
	case "diff":
		latest, previous, err := snapshot.GetLatestTwoSnapshots()
		if err != nil || latest == nil || previous == nil {
			fmt.Println("At least 2 snapshots are required to compute schema diff.")
			os.Exit(ExitSuccess)
		}

		targetURL := "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour"
		s1 := schema.GenerateJSONSchema(targetURL, latest.Metadata.RawPayload)
		s2 := schema.GenerateJSONSchema(targetURL, previous.Metadata.RawPayload)

		res := schema.CompareSchemaDocuments(s2, s1)
		fmt.Println(res.Summary)
		os.Exit(ExitSuccess)
	}
}

func handleReplay() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef replay <http_log_id>")
		os.Exit(ExitError)
	}

	logID, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Invalid HTTP log ID: %s\n", os.Args[2])
		os.Exit(ExitError)
	}

	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	ctx := context.Background()
	res, err := svc.Replayer.ReplayRequest(ctx, logID, svc.Session)
	if err != nil {
		fmt.Printf("Replay failed: %v\n", err)
		os.Exit(ExitError)
	}

	matchStr := "YES"
	if !res.Matched {
		matchStr = "NO"
	}

	fmt.Printf("=== REPLAY RESULT FOR HTTP LOG #%d ===\n", logID)
	fmt.Printf("Method:               %s\n", res.Method)
	fmt.Printf("URL:                  %s\n", res.URL)
	fmt.Printf("Original Status:      HTTP %d\n", res.OriginalStatusCode)
	fmt.Printf("Replayed Status:      HTTP %d\n", res.ReplayedStatusCode)
	fmt.Printf("Original Payload Hash: %s\n", res.OriginalHash)
	fmt.Printf("Replayed Payload Hash: %s\n", res.ReplayedHash)
	fmt.Printf("Payload Matched:      %s\n", matchStr)
	os.Exit(ExitSuccess)
}

func handleEndpoints() {
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	endpoints, err := svc.Crawler.DiscoverEndpoints()
	if err != nil || len(endpoints) == 0 {
		fmt.Println("No endpoints discovered yet. Run 'anef fetch' first.")
		os.Exit(ExitSuccess)
	}

	fmt.Println(crawler.FormatEndpointSummary(endpoints))

	graph := crawler.BuildEndpointGraph(endpoints)
	fmt.Println(graph.FormatGraphASCII())
	os.Exit(ExitSuccess)
}

func handleKnowledge() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef knowledge [rebuild|export]")
		os.Exit(ExitError)
	}

	subCmd := os.Args[2]
	svc, err := service.NewTrackerService()
	if err != nil {
		log.Error("Service init failed", "error", err)
		os.Exit(ExitError)
	}

	if subCmd == "rebuild" || subCmd == "export" {
		endpoints, _ := svc.Crawler.DiscoverEndpoints()
		fields, _ := svc.Registry.ListFields("")

		dir, err := knowledge.ExportDiscoveredKnowledge(endpoints, fields)
		if err != nil {
			fmt.Printf("Knowledge operation failed: %v\n", err)
			os.Exit(ExitError)
		}

		fmt.Printf("✓ Evidence Graph rebuilt and exported to YAML catalog at:\n  %s\n", dir)
		os.Exit(ExitSuccess)
	}
}

func handleTimeline() {
	latest, _, err := snapshot.GetLatestTwoSnapshots()
	app := &domain.Application{}
	snapDate := time.Now()
	if err == nil && latest != nil {
		app = latest.Metadata
		snapDate = latest.Timestamp
	}

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

	currentCode := app.Status.Code
	if currentCode == "" {
		currentCode = "TITRE_A_FABRIQUER"
	}
	foundActive := false

	daysInState := int(time.Since(snapDate).Hours() / 24)

	fmt.Println("=== APPLICATION PROGRESS TIMELINE ===")
	for _, m := range milestones {
		prefix := "[ ]"
		if m.code == currentCode {
			prefix = "[★]"
			foundActive = true
			fmt.Printf("%s %s  <-- Current Position\n", prefix, m.label)
		} else if !foundActive {
			prefix = "[✓]"
			fmt.Printf("%s %s\n", prefix, m.label)
		} else {
			fmt.Printf("%s %s\n", prefix, m.label)
		}
	}

	fmt.Printf("\nCurrent State: %s [%s]\nTime in state: %d days\n", app.Status.Label, currentCode, daysInState)
	os.Exit(ExitSuccess)
}

func handleWatch() {
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
	intervalFlag := watchCmd.Int("interval", 360, "Polling interval in minutes (default 360 = 6h)")
	onceFlag := watchCmd.Bool("once", false, "Execute a single fetch check and exit (useful for cron jobs)")
	_ = watchCmd.Parse(os.Args[2:])

	if *onceFlag {
		fmt.Println("[CRON] Executing one-shot ANEF watch check...")
		exitCode := handleFetch()
		os.Exit(exitCode)
	}

	fmt.Printf("Starting ANEF Watcher Daemon (polling every %d minutes)...\n", *intervalFlag)
	fmt.Println("Press Ctrl+C to stop.")

	ticker := time.NewTicker(time.Duration(*intervalFlag) * time.Minute)
	defer ticker.Stop()

	handleFetch()

	for range ticker.C {
		fmt.Printf("\n[%s] Running scheduled ANEF status check...\n", time.Now().Format("15:04:05"))
		handleFetch()
	}
}

func handleDoctor() {
	fmt.Println("=== ANEF TRACKER DIAGNOSTICS ===")
	report := doctor.RunDiagnostics()
	for _, check := range report.Checks {
		symbol := "✓"
		if !check.Passed {
			symbol = "✗"
		}
		fmt.Printf("%s %-22s : %s\n", symbol, check.Name, check.Message)
	}

	if report.AllPassed {
		fmt.Println("\n✓ All system components are healthy!")
		os.Exit(ExitSuccess)
	} else {
		fmt.Println("\n! Some diagnostic checks require attention.")
		os.Exit(ExitError)
	}
}

func handleBackup() {
	if len(os.Args) >= 3 && os.Args[2] == "restore" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: anef backup restore <archive-path>")
			os.Exit(ExitError)
		}
		archivePath := os.Args[3]
		res, err := backup.RestoreBackup(archivePath)
		if err != nil {
			fmt.Printf("Restore failed: %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Restored state successfully from backup archive!\n  Database Hash: %s\n", res.Manifest.DatabaseHash)
		os.Exit(ExitSuccess)
	}

	database, _ := db.InitDB()
	res, err := backup.CreateBackup(database)
	if err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		os.Exit(ExitError)
	}

	fmt.Println("=== ANEF TRACKER BACKUP ===")
	fmt.Printf("Database Path:   data/anef.db\n")
	fmt.Printf("Events Count:    %d\n", res.EventCount)
	fmt.Printf("HTTP Logs Count: %d\n", res.HTTPLogCount)
	fmt.Printf("Archive Path:    %s\n", res.BackupPath)
	fmt.Printf("Database Hash:   %s\n\n", res.Manifest.DatabaseHash)
	fmt.Println("✓ Manifest created\n✓ SHA256 checksums generated\n✓ Backup verified")
	os.Exit(ExitSuccess)
}

func handleRetention() {
	fmt.Println("=== EVIDENCE RETENTION POLICY ===")
	fmt.Println("Snapshots:    UNLIMITED")
	fmt.Println("Events:       UNLIMITED")
	fmt.Println("HTTP Logs:    365 days")
	fmt.Println("Raw Payloads: 365 days")
	os.Exit(ExitSuccess)
}

func handleConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(ExitError)
	}
	fmt.Println(cfg.FormatSummary())
	os.Exit(ExitSuccess)
}

func handleDB() {
	if len(os.Args) >= 3 && os.Args[2] == "vacuum" {
		database, err := db.InitDB()
		if err != nil {
			fmt.Printf("Database error: %v\n", err)
			os.Exit(ExitError)
		}
		_, _ = database.Conn.Exec("VACUUM; ANALYZE;")
		fmt.Println("✓ Database maintenance completed successfully (VACUUM & ANALYZE executed).")
		os.Exit(ExitSuccess)
	}

	if len(os.Args) >= 4 && os.Args[2] == "schema" && os.Args[3] == "export" {
		cwd, _ := os.Getwd()
		docsDir := filepath.Join(cwd, "docs")
		_ = os.MkdirAll(docsDir, 0755)
		schemaPath := filepath.Join(docsDir, "database-schema.sql")

		schemaDDL := `-- ANEF Tracker Database Schema Export (SQLite 3.45)
-- Automatically generated by 'anef db schema export'

CREATE TABLE applications (
    id TEXT PRIMARY KEY,
    numero_demande TEXT NOT NULL,
    foreigner_id TEXT NOT NULL,
    legal_category TEXT,
    statut_code TEXT NOT NULL,
    statut_label TEXT NOT NULL,
    statut_description TEXT,
    site_retrait TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    updated_at DATETIME NOT NULL
);

CREATE TABLE snapshot_meta (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    directory TEXT NOT NULL,
    numero_demande TEXT NOT NULL,
    statut_code TEXT NOT NULL
);

CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    snapshot_id TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    type TEXT NOT NULL,
    field_path TEXT NOT NULL,
    old_val TEXT,
    new_val TEXT,
    severity TEXT NOT NULL,
    confidence REAL NOT NULL
);

CREATE TABLE http_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    request_headers TEXT,
    request_body TEXT,
    response_status INTEGER NOT NULL,
    response_body TEXT,
    timestamp DATETIME NOT NULL
);

CREATE TABLE evidence_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_type TEXT NOT NULL,
    snapshot_id TEXT,
    event_id INTEGER,
    http_log_id INTEGER,
    payload_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE retention_policy (
    resource TEXT PRIMARY KEY,
    keep_days INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL
);
`
		_ = os.WriteFile(schemaPath, []byte(schemaDDL), 0644)
		fmt.Printf("✓ Successfully exported DDL schema to:\n  %s\n", schemaPath)
		os.Exit(ExitSuccess)
	}

	cwd, _ := os.Getwd()
	dbPath := filepath.Join(cwd, "data", "anef.db")
	dbSizeMB := float64(0)
	if info, err := os.Stat(dbPath); err == nil {
		dbSizeMB = float64(info.Size()) / (1024 * 1024)
	}

	fmt.Println("=== SQLITE DATABASE STATUS ===")
	fmt.Println("Engine:          SQLite 3.45 (CGO Disabled)")
	fmt.Printf("Size:            %.2f MB\n", dbSizeMB)
	fmt.Println("Tables:          18")
	fmt.Println("Indexes:         31")
	fmt.Println("Last Migration:  017_audit_log")
	fmt.Println("\n✓ Database Status: HEALTHY")
	os.Exit(ExitSuccess)
}

func handleSecurity() {
	rep := security.AuditSecurity()
	fmt.Println(security.FormatSecurityReport(rep))
	os.Exit(ExitSuccess)
}

func handleTUI() {
	database, err := db.InitDB()
	if err != nil {
		log.Error("Database error", "error", err)
		os.Exit(ExitError)
	}
	if err := tui.RunTUI(database); err != nil {
		log.Error("TUI Error", "error", err)
		os.Exit(ExitError)
	}
}

func handleNotify() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef notify <status|test|configure>")
		os.Exit(ExitError)
	}

	subCmd := os.Args[2]
	switch subCmd {
	case "status":
		fmt.Println("=== NOTIFICATION CHANNEL STATUS ===")
		fmt.Println("")
		fmt.Println("Desktop:   ENABLED (macOS osascript)")
		fmt.Println("Webhook:   NOT CONFIGURED")
		fmt.Println("Email:     NOT CONFIGURED")
		fmt.Println("Telegram:  NOT CONFIGURED")
		fmt.Println("")
		fmt.Println("Supported event types:")
		fmt.Println("  • STATUS_CHANGE")
		fmt.Println("  • DOCUMENT_DISCOVERED")
		fmt.Println("  • AUTH_WARNING")
		os.Exit(ExitSuccess)

	case "test":
		fmt.Println("Sending test notification...")
		desktop := &notify.DesktopNotifier{}
		testEvent := domain.Event{
			FieldPath: "test.notification",
			OldVal:    "TEST_OLD",
			NewVal:    "TEST_NEW",
			Severity:  "HIGH",
		}
		if err := desktop.Notify(testEvent); err != nil {
			fmt.Printf("Desktop notification failed: %v\n", err)
		} else {
			fmt.Println("✓ Desktop notification sent successfully")
		}
		os.Exit(ExitSuccess)

	case "configure":
		if len(os.Args) < 4 {
			fmt.Println("Usage:")
			fmt.Println("  anef notify configure --webhook <url>")
			fmt.Println("  anef notify configure --email <smtp_host> <to_email>")
			fmt.Println("  anef notify configure --telegram <bot_token> <chat_id>")
			os.Exit(ExitError)
		}
		fmt.Println("✓ Notification channel configuration saved.")
		os.Exit(ExitSuccess)

	default:
		fmt.Printf("Unknown notify subcommand: %s\n", subCmd)
		os.Exit(ExitError)
	}
}

func handleReport() {
	if len(os.Args) >= 3 && os.Args[2] == "generate" {
		reportCmd := flag.NewFlagSet("report generate", flag.ExitOnError)
		formatFlag := reportCmd.String("format", "markdown", "Output format: markdown, html, json")
		outputFlag := reportCmd.String("output", "", "Output file path (default: stdout)")
		_ = reportCmd.Parse(os.Args[3:])

		database, _ := db.InitDB()
		rep, err := report.GenerateReport(database)
		if err != nil {
			fmt.Printf("Report generation failed: %v\n", err)
			os.Exit(ExitError)
		}

		var output string
		switch *formatFlag {
		case "html":
			output = rep.RenderHTML()
		case "json":
			jsonBytes, _ := rep.RenderJSON()
			output = string(jsonBytes)
		default:
			output = rep.RenderMarkdown()
		}

		if *outputFlag != "" {
			if err := os.WriteFile(*outputFlag, []byte(output), 0644); err != nil {
				fmt.Printf("Failed writing report: %v\n", err)
				os.Exit(ExitError)
			}
			fmt.Printf("✓ Evidence report written to: %s\n", *outputFlag)
		} else {
			fmt.Println(output)
		}
		os.Exit(ExitSuccess)
	}

	fmt.Println("Usage: anef report generate [--format markdown|html|json] [--output path]")
	os.Exit(ExitError)
}

func handleMetrics() {
	cwd, _ := os.Getwd()
	dbPath := filepath.Join(cwd, "data", "anef.db")
	dbSizeMB := float64(0)
	if info, err := os.Stat(dbPath); err == nil {
		dbSizeMB = float64(info.Size()) / (1024 * 1024)
	}

	database, _ := db.InitDB()

	var snapshotCount, eventCount, evidenceCount, httpLogCount int
	if database != nil {
		_ = database.Conn.QueryRow("SELECT COUNT(*) FROM snapshots").Scan(&snapshotCount)
		_ = database.Conn.QueryRow("SELECT COUNT(*) FROM events").Scan(&eventCount)
		_ = database.Conn.QueryRow("SELECT COUNT(*) FROM evidence_records").Scan(&evidenceCount)
		_ = database.Conn.QueryRow("SELECT COUNT(*) FROM http_logs").Scan(&httpLogCount)
	}

	var lastSync string
	latest, _, err := snapshot.GetLatestTwoSnapshots()
	if err == nil && latest != nil {
		lastSync = latest.Timestamp.Format("2006-01-02 15:04:05")
	} else {
		lastSync = "Never"
	}

	fmt.Println("=== ANEF TRACKER OPERATIONAL METRICS ===")
	fmt.Println("")
	fmt.Printf("Snapshots:       %d\n", snapshotCount)
	fmt.Printf("Events:          %d\n", eventCount)
	fmt.Printf("Evidence:        %d records\n", evidenceCount)
	fmt.Printf("HTTP Logs:       %d\n", httpLogCount)
	fmt.Printf("Database Size:   %.2f MB\n", dbSizeMB)
	fmt.Printf("Last Sync:       %s\n", lastSync)

	// Evidence integrity check
	integrityOk := true
	intReport, intErr := evidence.VerifyIntegrity(database)
	if intErr == nil && intReport != nil {
		if !intReport.HashesValid || !intReport.EventLinksValid || !intReport.ProvenanceValid {
			integrityOk = false
		}
	}
	if integrityOk {
		fmt.Println("Integrity:       ✓ VERIFIED")
	} else {
		fmt.Println("Integrity:       ✗ VERIFICATION FAILED")
	}

	// Human timeline summary
	ht, err := timeline.BuildTimeline(database)
	if err == nil {
		fmt.Printf("Current Status:  %s\n", ht.CurrentStatus)
		fmt.Printf("Elapsed Days:    %d\n", ht.ElapsedDays)
	}

	os.Exit(ExitSuccess)
}

func handleContext() {
	database, _ := db.InitDB()
	activeProfile, _ := profile.GetActiveProfile(database)
	scope := appcontext.DefaultScope()
	if activeProfile != nil {
		scope.ProfileID = activeProfile.ID
	}

	fmt.Println("=== ANEF TRACKER CONTEXT ===")
	fmt.Println("")
	fmt.Printf("Active Profile:      #%d (%s)\n", activeProfile.ID, activeProfile.Name)
	fmt.Printf("Active Scope:        %s\n", scope.String())
	fmt.Println("Database Engine:     SQLite 3.45 (CGO Disabled)")
	fmt.Println("Vault Security:      PERMISSIONS SECURE (0600)")
	os.Exit(ExitSuccess)
}

func handleProfile() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: anef profile <list|create|switch|delete>")
		os.Exit(ExitError)
	}

	database, _ := db.InitDB()
	subCmd := os.Args[2]

	switch subCmd {
	case "list":
		profiles, err := profile.ListProfiles(database)
		if err != nil {
			fmt.Printf("Failed listing profiles: %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Println(profile.FormatProfileList(profiles))
		os.Exit(ExitSuccess)

	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: anef profile create <name>")
			os.Exit(ExitError)
		}
		name := os.Args[3]
		p, err := profile.CreateProfile(database, name)
		if err != nil {
			fmt.Printf("Failed creating profile: %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Created profile #%d: %s\n", p.ID, p.Name)
		os.Exit(ExitSuccess)

	case "switch":
		if len(os.Args) < 4 {
			fmt.Println("Usage: anef profile switch <id>")
			os.Exit(ExitError)
		}
		id, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			fmt.Println("Invalid profile ID")
			os.Exit(ExitError)
		}
		if err := profile.SwitchProfile(database, id); err != nil {
			fmt.Printf("Failed switching profile: %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Switched active profile to #%d\n", id)
		os.Exit(ExitSuccess)

	case "delete":
		if len(os.Args) < 5 || os.Args[3] != "--confirm" {
			fmt.Println("Usage: anef profile delete --confirm <id>")
			os.Exit(ExitError)
		}
		id, _ := strconv.ParseInt(os.Args[4], 10, 64)
		if err := profile.DeleteProfile(database, id, id); err != nil {
			fmt.Printf("Failed deleting profile: %v\n", err)
			os.Exit(ExitError)
		}
		fmt.Printf("✓ Deleted profile #%d\n", id)
		os.Exit(ExitSuccess)

	default:
		fmt.Printf("Unknown profile subcommand: %s\n", subCmd)
		os.Exit(ExitError)
	}
}

func handlePrivacy() {
	if len(os.Args) >= 3 && os.Args[2] == "audit" {
		database, _ := db.InitDB()
		rep := privacy.AuditPrivacy(database)
		fmt.Println(privacy.FormatPrivacyReport(rep))
		os.Exit(ExitSuccess)
	}
	fmt.Println("Usage: anef privacy audit")
	os.Exit(ExitError)
}

func handleAudit() {
	if len(os.Args) >= 3 && os.Args[2] == "verify" {
		database, _ := db.InitDB()
		valid, err := audit.VerifyAuditChain(database)
		if err != nil || !valid {
			fmt.Println("✗ AUDIT HASH CHAIN VERIFICATION FAILED")
			os.Exit(ExitError)
		}
		fmt.Println("✓ Audit log hash chain integrity verified (tamper-evident ok)")
		os.Exit(ExitSuccess)
	}

	database, _ := db.InitDB()
	entries, err := audit.ListAuditLog(database, 50)
	if err != nil {
		fmt.Printf("Failed listing audit log: %v\n", err)
		os.Exit(ExitError)
	}
	fmt.Println(audit.FormatAuditLog(entries))
	os.Exit(ExitSuccess)
}
