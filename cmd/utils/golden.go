package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func LoadGolden(t *testing.T, filename string) string {
	t.Helper()
	goldenPath := filepath.Join("..", "testdata", filename)
	content, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}
	return string(content)
}
