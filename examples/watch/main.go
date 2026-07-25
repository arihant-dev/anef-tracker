package main

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/watch"
)

func main() {
	fmt.Println("=== ANEF Tracker Example: Watch Daemon Scheduler ===")

	database, _ := db.InitDB()
	scheduler := watch.NewScheduler(database, 360) // 6 hours

	fmt.Printf("Configured Scheduler interval: %d minutes\n", scheduler.Config.IntervalMinutes)
	_ = watch.RecordWatchRun(database, "SUCCESS", 0)
	fmt.Println("Recorded test watch run entry successfully.")
}
