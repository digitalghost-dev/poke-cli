package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGolden(t *testing.T) {
	testDir := filepath.Join("..", "testdata")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create testdata directory: %v", err)
	}

	filename := "example.golden"
	expected := "hello golden file"

	fullPath := filepath.Join(testDir, filename)
	if err := os.WriteFile(fullPath, []byte(expected), 0644); err != nil {
		t.Fatalf("failed to write golden file: %v", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("failed to remove golden file: %v", err)
		}
	}(fullPath)

	actual := LoadGolden(t, filename)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
