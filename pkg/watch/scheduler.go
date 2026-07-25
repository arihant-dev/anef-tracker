package watch

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"time"
)

type SchedulerConfig struct {
	IntervalMinutes int
	MaxRuns         int
}

type Scheduler struct {
	Config   SchedulerConfig
	Database *db.DB
	RunCount int
}

func NewScheduler(database *db.DB, intervalMinutes int) *Scheduler {
	return &Scheduler{
		Config: SchedulerConfig{
			IntervalMinutes: intervalMinutes,
			MaxRuns:         0, // 0 = unlimited
		},
		Database: database,
	}
}

func (s *Scheduler) Start(fetchFn func() int) {
	fmt.Printf("[SCHEDULER] Starting watch scheduler (interval: %d minutes)\n", s.Config.IntervalMinutes)

	// Execute immediately on start
	exitCode := fetchFn()
	s.RunCount++
	status := "SUCCESS"
	changes := 0
	if exitCode == 3 { // ExitChangeDetect
		changes = 1
	}
	if exitCode == 1 {
		status = "ERROR"
	}
	_ = RecordWatchRun(s.Database, status, changes)

	ticker := time.NewTicker(time.Duration(s.Config.IntervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("\n[%s] Scheduled watch run #%d...\n", time.Now().Format("15:04:05"), s.RunCount+1)
		exitCode := fetchFn()
		s.RunCount++

		status := "SUCCESS"
		changes := 0
		if exitCode == 3 {
			changes = 1
		}
		if exitCode == 1 {
			status = "ERROR"
		}
		_ = RecordWatchRun(s.Database, status, changes)

		if s.Config.MaxRuns > 0 && s.RunCount >= s.Config.MaxRuns {
			fmt.Printf("[SCHEDULER] Reached max runs (%d). Stopping.\n", s.Config.MaxRuns)
			return
		}
	}
}
