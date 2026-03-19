package tcg

import (
	"errors"
	"strings"
	"testing"
)

func TestFetchStandings_ConnectionError(t *testing.T) {
	orig := supabaseConn
	defer func() { supabaseConn = orig }()

	supabaseConn = func(_ string) ([]byte, error) {
		return nil, errors.New("connection refused")
	}

	cmd := fetchData("London")
	msg := cmd()

	result, ok := msg.(standingsDataMsg)
	if !ok {
		t.Fatalf("expected standingsDataMsg, got %T", msg)
	}
	if result.err == nil {
		t.Error("expected error, got nil")
	}
	if result.items != nil {
		t.Error("expected nil items on error")
	}
}

func TestFetchStandings_InvalidJSON(t *testing.T) {
	orig := supabaseConn
	defer func() { supabaseConn = orig }()

	supabaseConn = func(_ string) ([]byte, error) {
		return []byte("not json"), nil
	}

	cmd := fetchData("London")
	msg := cmd()

	result, ok := msg.(standingsDataMsg)
	if !ok {
		t.Fatalf("expected standingsDataMsg, got %T", msg)
	}
	if result.err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}

func TestFetchStandings_Success(t *testing.T) {
	orig := supabaseConn
	defer func() { supabaseConn = orig }()

	supabaseConn = func(_ string) ([]byte, error) {
		return []byte(`[{"rank":1,"name":"Alice","player_country":"USA"},{"rank":2,"name":"Bob","player_country":"Japan"}]`), nil
	}

	cmd := fetchData("London")
	msg := cmd()

	result, ok := msg.(standingsDataMsg)
	if !ok {
		t.Fatalf("expected standingsDataMsg, got %T", msg)
	}
	if result.err != nil {
		t.Errorf("expected no error, got %v", result.err)
	}
	if len(result.items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.items))
	}
	if result.items[0].Name != "Alice" {
		t.Errorf("expected first item name to be Alice, got %q", result.items[0].Name)
	}
}

func TestFetchStandings_URLEncoding(t *testing.T) {
	orig := supabaseConn
	defer func() { supabaseConn = orig }()

	var capturedURL string
	supabaseConn = func(url string) ([]byte, error) {
		capturedURL = url
		return []byte(`[]`), nil
	}

	cmd := fetchData("São Paulo")
	cmd()

	if !strings.Contains(capturedURL, "S%C3%A3o") {
		t.Errorf("expected URL-encoded tournament name in URL, got %q", capturedURL)
	}
}

func TestCountryFlag(t *testing.T) {
	tests := []struct {
		isoCode string
		want    string
	}{
		{"gb", "🇬🇧"},
		{"GB", "🇬🇧"},
		{"us", "🇺🇸"},
		{"jp", "🇯🇵"},
		{"", ""},
		{"x", ""},
		{"abc", ""},
	}

	for _, tt := range tests {
		t.Run(tt.isoCode, func(t *testing.T) {
			got := countryFlag(tt.isoCode)
			if got != tt.want {
				t.Errorf("countryFlag(%q) = %q, want %q", tt.isoCode, got, tt.want)
			}
		})
	}
}
