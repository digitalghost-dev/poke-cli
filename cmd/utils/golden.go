package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func LoadGolden(t *testing.T, filename string) string {
	t.Helper()

	basePath := filepath.Join("../..", "testdata")

	// gosec
	joinedPath := filepath.Join(basePath, filename)
	cleanPath := filepath.Clean(joinedPath)

	// Ensure the cleaned path is still within basePath
	basePathClean := filepath.Clean(basePath)
	if !strings.HasPrefix(cleanPath, basePathClean+string(os.PathSeparator)) {
		t.Fatalf("invalid filename: %q", filename)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	return string(content)
}
