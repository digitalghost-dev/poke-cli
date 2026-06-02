package types

import (
	"strings"
	"testing"

	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestDamageTable(t *testing.T) {
	output, err := DamageTable("fire", "type")
	if err != nil {
		t.Fatalf("DamageTable returned an error: %v", err)
	}
	output = styling.StripANSI(output)

	if !strings.Contains(output, "You selected the Fire type.") {
		t.Errorf("Expected output to contain Fire type header, got:\n%s", output)
	}

	if !strings.Contains(output, "Damage Chart:") {
		t.Errorf("Expected output to contain 'Damage Chart:', got:\n%s", output)
	}
}

func TestDamageTable_TypeNotFound(t *testing.T) {
	_, err := DamageTable("notatype", "type")
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
