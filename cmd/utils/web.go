// Opens the user's default browser.

package utils

import (
	"os/exec"
	"runtime"

	tea "charm.land/bubbletea/v2"
)

// Open returns a tea.Cmd that opens url in the user's default browser.
// It is best-effort: if no browser launcher is available on the system, it
// does nothing rather than surfacing an error into the TUI.
func Open(url string) tea.Cmd {
	return func() tea.Msg {
		var (
			browserCmd string
			openCmd    *exec.Cmd
		)

		switch runtime.GOOS {
		case "windows":
			browserCmd = "cmd"
			openCmd = exec.Command("cmd", "/c", "start", url) // #nosec G204
		case "darwin":
			browserCmd = "open"
			openCmd = exec.Command("open", url) // #nosec G204
		default:
			browserCmd = "xdg-open"
			openCmd = exec.Command("xdg-open", url) // #nosec G204
		}

		if _, err := exec.LookPath(browserCmd); err != nil {
			return nil
		}
		_ = openCmd.Start()
		return nil
	}
}
