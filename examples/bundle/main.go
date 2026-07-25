package main

import (
	"fmt"
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/export"
)

func main() {
	fmt.Println("=== ANEF Tracker Example: Redacted Evidence Bundle Exporter ===")

	database, _ := db.InitDB()
	scope := appcontext.DefaultScope()

	res, err := export.CreateEvidenceBundle(database, scope, true) // redact=true
	if err != nil {
		fmt.Printf("Bundle export error: %v\n", err)
		return
	}

	fmt.Printf("✓ Created redacted evidence bundle ZIP archive:\n  Path: %s\n  Hash: %s\n",
		res.ArchivePath, res.DatabaseHash)
}
