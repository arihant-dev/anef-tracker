package main

import (
	"fmt"
	appcontext "github.com/arihant-dev/anef-tracker/pkg/context"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/profile"
)

func main() {
	fmt.Println("=== ANEF Tracker Example: Basic Profile & Context ===")

	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Database init failed: %v\n", err)
		return
	}

	active, err := profile.GetActiveProfile(database)
	if err != nil {
		fmt.Printf("GetActiveProfile failed: %v\n", err)
		return
	}

	scope := appcontext.DefaultScope()
	if active != nil {
		scope.ProfileID = active.ID
	}

	fmt.Printf("Active Profile : %s (#%d)\n", active.Name, active.ID)
	fmt.Printf("Scope Context  : %s\n", scope.String())
}
