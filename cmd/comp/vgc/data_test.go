package vgc

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestDecode_Success(t *testing.T) {
	body := []byte(`[
		{"rank":1,"name":"Arsal","points":45,"player_country":"United States","type":"Regional","text_date":"May 29-31, 2026","player_quantity":1013,"location":"Indianapolis",
		 "team":[{"name":"Venusaur"},{"name":"Landorus [Therian Forme]"},{"name":"Iron Crown"}]},
		{"rank":2,"name":"Wolfe","points":42,"player_country":"United States",
		 "team":[{"name":"Sneasler"},{"name":"Landorus [Therian Forme]"},{"name":"Iron Crown"}]}
	]`)
	d, err := decode(body)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(d.TableRows) != 2 {
		t.Errorf("expected 2 table rows, got %d", len(d.TableRows))
	}
	if len(d.Countries) != 1 {
		t.Errorf("expected 1 country tally, got %d", len(d.Countries))
	}

	overview := d.Overview(120, lipgloss.Color("#7D56F4"))
	for _, s := range []string{"Indianapolis", "Regional", "Arsal", "Venusaur", "1,013", "Winning Team"} {
		if !strings.Contains(overview, s) {
			t.Errorf("expected overview to contain %q", s)
		}
	}

	if !strings.Contains(overview, "Landorus") || strings.Contains(overview, "Therian Forme") {
		t.Error("expected Landorus forme bracket stripped in the winner team box")
	}

	if d.Extra.NameHeader != "Pokémon" || d.Extra.CountHeader != "Teams" {
		t.Errorf("unexpected Extra headers: %q / %q", d.Extra.NameHeader, d.Extra.CountHeader)
	}
	usage := map[string]int{}
	for _, it := range d.Extra.Items {
		usage[it.Label] = it.Count
	}
	if usage["Iron Crown"] != 2 || usage["Landorus [Therian Forme]"] != 2 {
		t.Errorf("expected Usage tallies to keep full forme names, got %v", usage)
	}
}

func TestDecode_TeamJSONB(t *testing.T) {
	body := []byte(`[{"rank":1,"name":"Arsal","team":[
		{"id":"10021","name":"Landorus [Therian Forme]","item":"Choice Band","ability":"Intimidate","teratype":"Steel","badges":["Stomping Tantrum","U-turn","Earthquake","Rock Slide"]}
	]}]`)
	d, err := decode(body)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	usage := map[string]int{}
	for _, it := range d.Extra.Items {
		usage[it.Label] = it.Count
	}
	if usage["Landorus [Therian Forme]"] != 1 {
		t.Errorf("expected forme name kept in usage tallies, got %v", usage)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	if _, err := decode([]byte("not json")); err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}

func TestStandingsColumns_NoDeck(t *testing.T) {
	cols := standingsColumns(120)
	if len(cols) != 7 {
		t.Fatalf("expected 7 columns, got %d", len(cols))
	}
	for _, c := range cols {
		if c.Title == "Deck" {
			t.Error("VGC standings should not have a Deck column")
		}
	}
}

func TestSpec_URLs(t *testing.T) {
	s := Spec()
	if !strings.Contains(s.ListURL, "comp_vgc_standings_view") {
		t.Errorf("expected VGC view in list URL, got %q", s.ListURL)
	}
	durl := s.DashboardURL("São Paulo")
	if !strings.Contains(durl, "comp_vgc_standings_view") || !strings.Contains(durl, "team") {
		t.Errorf("expected VGC view + team column in dashboard URL, got %q", durl)
	}
	if !strings.Contains(durl, "S%C3%A3o") {
		t.Errorf("expected URL-encoded location, got %q", durl)
	}
}
