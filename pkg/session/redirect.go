package session

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens the specified URL in the default system browser.
func OpenBrowser(targetURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{targetURL}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", targetURL}
	default:
		cmd = "xdg-open"
		args = []string{targetURL}
	}

	return exec.Command(cmd, args...).Start()
}
