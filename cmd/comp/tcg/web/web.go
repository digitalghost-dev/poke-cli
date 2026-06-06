// Package web opens the poke-cli Streamlit dashboard in the user's browser.
package web

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
			openCmd = exec.Command("cmd", "/c", "start", url) //nolint:gosec
		case "darwin":
			browserCmd = "open"
			openCmd = exec.Command("open", url)
		default:
			browserCmd = "xdg-open"
			openCmd = exec.Command("xdg-open", url)
		}

		if _, err := exec.LookPath(browserCmd); err != nil {
			return nil
		}
		_ = openCmd.Start()
		return nil
	}
}
