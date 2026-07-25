package main

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/report"
)

func main() {
	fmt.Println("=== ANEF Tracker Example: Evidence Report Generation ===")

	database, _ := db.InitDB()
	rep, err := report.GenerateReport(database)
	if err != nil {
		fmt.Printf("Report generation error: %v\n", err)
		return
	}

	fmt.Println(rep.RenderMarkdown())
}
