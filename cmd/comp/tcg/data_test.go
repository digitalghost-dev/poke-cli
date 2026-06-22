package tcg

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestDecode_Success(t *testing.T) {
	body := []byte(`[
		{"rank":1,"name":"Ash","points":47,"record":"15 - 1 - 0","deck":"gardevoir","player_country":"USA","type":"Regional","text_date":"Jan 10","player_quantity":500,"location":"London"},
		{"rank":2,"name":"Misty","points":44,"deck":"dragapult","player_country":"Japan"}
	]`)
	d, err := decode(body)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(d.TableRows) != 2 {
		t.Errorf("expected 2 table rows, got %d", len(d.TableRows))
	}
	if len(d.Countries) != 2 {
		t.Errorf("expected 2 country tallies, got %d", len(d.Countries))
	}

	overview := d.Overview(120, lipgloss.Color("#7D56F4"))
	for _, s := range []string{"London", "Regional", "Ash", "gardevoir", "500"} {
		if !strings.Contains(overview, s) {
			t.Errorf("expected overview to contain %q", s)
		}
	}

	if d.Extra.NameHeader != "Deck" || d.Extra.CountHeader != "Players" {
		t.Errorf("unexpected Extra headers: %q / %q", d.Extra.NameHeader, d.Extra.CountHeader)
	}
	decks := map[string]int{}
	for _, it := range d.Extra.Items {
		decks[it.Label] = it.Count
	}
	if decks["gardevoir"] != 1 || decks["dragapult"] != 1 {
		t.Errorf("expected Decks tallies for gardevoir+dragapult, got %v", decks)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	if _, err := decode([]byte("not json")); err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}

func TestDecode_EmptyCountrySkipped(t *testing.T) {
	body := []byte(`[{"rank":1,"name":"Ash","player_country":""},{"rank":2,"name":"Misty","player_country":"Japan"}]`)
	d, err := decode(body)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(d.Countries) != 1 {
		t.Errorf("expected empty country skipped (1 tally), got %d", len(d.Countries))
	}
}

func TestStandingsColumns_HasDeck(t *testing.T) {
	cols := standingsColumns(120)
	if len(cols) != 8 {
		t.Fatalf("expected 8 columns, got %d", len(cols))
	}
	found := false
	for _, c := range cols {
		if c.Title == "Deck" {
			found = true
		}
	}
	if !found {
		t.Error("expected a Deck column in the TCG standings table")
	}
}

func TestSpec_URLs(t *testing.T) {
	s := Spec()
	if !strings.Contains(s.ListURL, "comp_tcg_standings_view") {
		t.Errorf("expected TCG view in list URL, got %q", s.ListURL)
	}
	durl := s.DashboardURL("São Paulo")
	if !strings.Contains(durl, "comp_tcg_standings_view") {
		t.Errorf("expected TCG view in dashboard URL, got %q", durl)
	}
	if !strings.Contains(durl, "S%C3%A3o") {
		t.Errorf("expected URL-encoded location, got %q", durl)
	}
}
