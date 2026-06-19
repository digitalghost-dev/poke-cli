package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestOpenDoesNothingWhenLauncherIsUnavailable(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	cmd := Open("https://example.com")
	if cmd == nil {
		t.Fatal("Open returned nil Cmd")
	}

	if msg := cmd(); msg != nil {
		t.Fatalf("expected nil message, got %#v", msg)
	}
}

func TestOpenStartsLauncherWithURL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses a POSIX shell script as a fake launcher")
	}

	dir := t.TempDir()
	marker := filepath.Join(dir, "launched")
	launcher := filepath.Join(dir, browserLauncherName())
	script := "#!/bin/sh\nprintf '%s' \"$1\" > \"$OPEN_TEST_MARKER\"\n"

	if err := os.WriteFile(launcher, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake launcher: %v", err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("OPEN_TEST_MARKER", marker)

	url := "https://example.com/cards?q=pikachu"
	if msg := Open(url)(); msg != nil {
		t.Fatalf("expected nil message, got %#v", msg)
	}

	got := waitForLauncherOutput(t, marker)
	if got != url {
		t.Fatalf("launcher received %q, want %q", got, url)
	}
}

func browserLauncherName() string {
	if runtime.GOOS == "darwin" {
		return "open"
	}
	return "xdg-open"
}

func waitForLauncherOutput(t *testing.T, marker string) string {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		got, err := os.ReadFile(marker)
		if err == nil {
			return string(got)
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("fake launcher did not write marker file %q", marker)
	return ""
}
