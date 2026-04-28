package types

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestDamageTable(t *testing.T) {
	originalStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	os.Stdout = w

	if err := DamageTable("fire", "type"); err != nil {
		t.Fatalf("DamageTable returned an error: %v", err)
	}

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
	output := styling.StripANSI(buf.String())

	// Step 7: Assert the output contains expected strings
	if !strings.Contains(output, "You selected the Fire type.") {
		t.Errorf("Expected output to contain Fire type header, got:\n%s", output)
	}

	if !strings.Contains(output, "Damage Chart:") {
		t.Errorf("Expected output to contain 'Damage Chart:', got:\n%s", output)
	}
}

func TestDamageTable_TypeNotFound(t *testing.T) {
	err := DamageTable("notatype", "type")
	if err == nil {
		t.Fatal("expected an error for unknown type, got nil")
	}
	actual := styling.StripANSI(err.Error())
	if !strings.Contains(actual, "Type not found") {
		t.Errorf("expected error to contain 'Type not found', got: %s", actual)
	}
	if !strings.Contains(actual, "Perhaps a typo?") {
		t.Errorf("expected error to contain 'Perhaps a typo?', got: %s", actual)
	}
}
