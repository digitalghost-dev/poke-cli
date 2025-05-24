package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func LoadGolden(t *testing.T, filename string) string {
	t.Helper()

	// Get the current working directory (e.g., cmd/types/)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Fatal("could not find project root (go.mod)")
		}
		projectRoot = parent
	}

	testdataPath := filepath.Join(projectRoot, "testdata")

	joinedPath := filepath.Join(testdataPath, filename)
	cleanPath := filepath.Clean(joinedPath)

	testdataPathClean := filepath.Clean(testdataPath)
	if !strings.HasPrefix(cleanPath, testdataPathClean+string(os.PathSeparator)) {
		t.Fatalf("invalid filename: %q", filename)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	return string(content)
}
