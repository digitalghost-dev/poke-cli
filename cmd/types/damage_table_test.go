package types

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDamageTable(t *testing.T) {
	originalStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stdout = w

	DamageTable("fire", "type")

	err = w.Close()
	if err != nil {
		t.Fatalf("Failed to close pipe: %v", err)
	}
	os.Stdout = originalStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	output := buf.String()

	// Step 7: Assert the output contains expected strings
	if !strings.Contains(output, "You selected the Fire type.") {
		t.Errorf("Expected output to contain Fire type header, got:\n%s", output)
	}

	if !strings.Contains(output, "Damage Chart:") {
		t.Errorf("Expected output to contain 'Damage Chart:', got:\n%s", output)
	}
}
