package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v0.9.0"
	Commit    = "dev"
	BuildDate = "2026-07-25"
)

func GetVersionInfo() string {
	return fmt.Sprintf("ANEF Tracker\n\nVersion:    %s\nCommit:     %s\nBuild Date: %s\nGo Version: %s\nPlatform:   %s/%s",
		Version, Commit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
