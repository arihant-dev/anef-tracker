package commands

type CommandItem struct {
	Key         string
	Title       string
	Description string
	TargetTab   int
}

var Catalog = []CommandItem{
	{Key: "fetch", Title: "Fetch Latest Application Status", Description: "Poll live ANEF REST API for updates", TargetTab: 0},
	{Key: "status", Title: "Open Overview Dashboard", Description: "View current residence permit status summary", TargetTab: 0},
	{Key: "timeline", Title: "Open Progress Timeline", Description: "View application progress milestones", TargetTab: 1},
	{Key: "diff", Title: "View Snapshot Structural Diff", Description: "Compare latest snapshot changes", TargetTab: 2},
	{Key: "events", Title: "Open Recorded Events Log", Description: "View status & schema discovery events", TargetTab: 3},
	{Key: "documents", Title: "Open Document Inventory", Description: "View uploaded files & attestations", TargetTab: 4},
	{Key: "api", Title: "Open API Explorer", Description: "Inspect observed REST endpoints", TargetTab: 5},
	{Key: "schema", Title: "Open Schema Intelligence", Description: "View registered payload fields & types", TargetTab: 6},
	{Key: "replay", Title: "Open HTTP Replay Engine", Description: "Replay requests & verify hash matching", TargetTab: 7},
	{Key: "logs", Title: "Open HTTP Traffic Logs", Description: "View sanitized HTTP requests & latencies", TargetTab: 8},
	{Key: "workflow", Title: "Open Reconstructed Workflow Graph", Description: "View deterministic workflow state machine", TargetTab: 9},
	{Key: "analytics", Title: "Open Historical Duration Analytics", Description: "View status duration statistics & medians", TargetTab: 10},
	{Key: "graph", Title: "Open Knowledge Graph", Description: "Explore discovered API & field graph relationships", TargetTab: 11},
	{Key: "export", Title: "Export View Data", Description: "Export events/schemas in CSV, JSON, YAML", TargetTab: -1},
	{Key: "quit", Title: "Quit Application", Description: "Exit ANEF Tracker TUI", TargetTab: -1},
}
