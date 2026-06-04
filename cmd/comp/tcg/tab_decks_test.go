package tcg

import (
	"strings"
	"testing"
)

func TestDecksContent(t *testing.T) {
	result := decksContent([]deckStats{{Deck: "gardevoir", Total: 10}}, 80)
	if !strings.Contains(result, "gardevoir") {
		t.Error("expected output to contain deck name")
	}
}
