package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v1.0.0"
	Commit    = "dev"
	BuildDate = "2026-07-28"
)

func GetVersionInfo() string {
	return fmt.Sprintf("ANEF Tracker\n\nVersion:    %s\nCommit:     %s\nBuild Date: %s\nGo Version: %s\nPlatform:   %s/%s",
		Version, Commit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
