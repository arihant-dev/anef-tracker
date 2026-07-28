package doctor

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/evidence"
	"github.com/arihant-dev/anef-tracker/pkg/knowledge"
	"github.com/arihant-dev/anef-tracker/pkg/schema"
	"github.com/arihant-dev/anef-tracker/pkg/session"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
)

type CheckResult struct {
	Name    string
	Passed  bool
	Message string
}

type DoctorReport struct {
	Timestamp time.Time
	Checks    []CheckResult
	AllPassed bool
}

func RunDiagnostics() *DoctorReport {
	report := &DoctorReport{
		Timestamp: time.Now(),
		AllPassed: true,
	}

	// 1. SQLite DB check
	database, err := db.InitDB()
	if err != nil {
		report.addCheck("Database Connection", false, fmt.Sprintf("SQLite error: %v", err))
	} else {
		report.addCheck("Database Connection", true, fmt.Sprintf("SQLite database healthy and migrated (%s)", schema.CurrentSchemaMigration))
	}

	// 2. Session File check
	sess, err := session.LoadSession()
	if err != nil {
		report.addCheck("Session File", false, fmt.Sprintf("Session error or missing: %v", err))
	} else {
		report.addCheck("Session File", true, fmt.Sprintf("Active session found for user login %s", sess.User))

		// 3. Token Validity check
		if sess.IsExpired() {
			report.addCheck("Token Validity", false, "JWT token expired. Use 'anef login' or refresh token")
		} else {
			rem := time.Until(sess.ExpiresAt)
			report.addCheck("Token Validity", true, fmt.Sprintf("Token valid (expires in %s)", rem.Round(time.Second)))
		}
	}

	// 4. Provider Reachability check
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://administration-etrangers-en-france.interieur.gouv.fr/usagers/")
	if err != nil {
		report.addCheck("ANEF Connectivity", false, fmt.Sprintf("Network reachability failed: %v", err))
	} else {
		resp.Body.Close()
		report.addCheck("ANEF Connectivity", true, fmt.Sprintf("HTTP %d reachability confirmed", resp.StatusCode))
	}

	// 5. Snapshot Storage check
	latest, _, err := snapshot.GetLatestTwoSnapshots()
	if err != nil || latest == nil {
		report.addCheck("Snapshot Storage", false, "No snapshots recorded yet")
	} else {
		report.addCheck("Snapshot Storage", true, fmt.Sprintf("Latest snapshot %s verified", latest.SnapshotID))
	}

	// 6. Evidence Integrity check
	evReport, err := evidence.VerifyIntegrity(database)
	if err == nil && evReport.HashesValid {
		report.addCheck("Evidence Integrity", true, fmt.Sprintf("%d evidence records verified", evReport.EventsChecked+evReport.HTTPLogsChecked))
	} else {
		report.addCheck("Evidence Integrity", false, "Discrepancies identified in evidence hashes")
	}

	// 7. Graph Validation check
	if database != nil {
		gb := knowledge.NewGraphBuilder(database)
		g, _ := gb.BuildFromDB()
		valRep := g.ValidateGraph()
		if valRep.Valid {
			report.addCheck("Graph Validation", true, fmt.Sprintf("%d nodes and %d edges consistent", valRep.TotalNodes, valRep.TotalEdges))
		} else {
			report.addCheck("Graph Validation", false, "Orphan nodes or broken edges identified")
		}
		database.Conn.Close()
	}

	// 8. Backup Availability check
	cwd, _ := os.Getwd()
	backupDir := filepath.Join(cwd, "backups")
	if entries, err := os.ReadDir(backupDir); err == nil && len(entries) > 0 {
		report.addCheck("Backup Availability", true, fmt.Sprintf("%d backup archives available", len(entries)))
	} else {
		report.addCheck("Backup Availability", true, "No backup archives created yet (use 'anef backup create')")
	}

	return report
}

func (r *DoctorReport) addCheck(name string, passed bool, msg string) {
	if !passed {
		r.AllPassed = false
	}
	r.Checks = append(r.Checks, CheckResult{
		Name:    name,
		Passed:  passed,
		Message: msg,
	})
}
