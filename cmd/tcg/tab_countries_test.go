package tcg

import (
	"strings"
	"testing"
)

func TestCountriesContent(t *testing.T) {
	result := countriesContent([]countryStats{{Country: "USA", Total: 10}}, 80)
	if !strings.Contains(result, "USA") {
		t.Error("expected output to contain country name")
	}
}
