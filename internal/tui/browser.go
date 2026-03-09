package tui

import (
	"os/exec"
	"runtime"
)

// OpenInBrowser opens the given file path in the system default browser.
// It uses platform-specific commands: open (macOS), xdg-open (Linux), start (Windows).
func OpenInBrowser(htmlPath string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", htmlPath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", htmlPath)
	default:
		cmd = exec.Command("xdg-open", htmlPath)
	}
	return cmd.Start()
}
