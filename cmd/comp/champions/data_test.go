package champions

import (
	"errors"
	"testing"
)

func TestFetchDashboardData_Success(t *testing.T) {
	var capturedURLs []string
	conn := func(url string) ([]byte, error) {
		capturedURLs = append(capturedURLs, url)
		switch url {
		case compInfoURL:
			return []byte(`[
				{
					"pokemon":"Miraidon",
					"web_url":"https://example.com/pokemon/1",
					"common_moves":[{"name":"Protect","usage_percent":90.5}],
					"common_abilities":[{"name":"Hadron Engine","usage_percent":100}],
					"common_items":[{"name":"Choice Specs","usage_percent":45.2}],
					"common_teammates":[{"name":"Flutter Mane","usage_percent":70.1}]
				}
			]`), nil
		case topTeamsURL:
			return []byte(`[
				{"author":"Alice","record":"7-1","tournament":"Worlds 2026","archetypes":["Hyper Offense"],"pokemon":["Miraidon","Flutter Mane","Iron Hands"],"web_url":"https://example.com/1"},
				{"author":"Bob","record":"6-2","tournament":"Regional","archetypes":[],"pokemon":["Calyrex"],"web_url":""}
			]`), nil
		case usageURL:
			return []byte(`[
				{"rank":1,"pokemon":"Basculegion","usage_percent":51.5},
				{"rank":2,"pokemon":"Kingambit","usage_percent":40.69}
			]`), nil
		case speedTiersURL:
			return []byte(`[
				{"rank":1,"pokemon":"Mega Aerodactyl","base_spe":150,"neutral_0_sp":170,"neutral_32_sp":202,"neg_spe_0_sp":153,"max_speed":222,"max_scarf":333,"neutral_32_scarf":303}
			]`), nil
		default:
			t.Fatalf("unexpected URL %q", url)
			return nil, nil
		}
	}

	msg := fetchDashboardData(conn)().(dataMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.data == nil {
		t.Fatal("expected data, got nil")
	}
	if len(msg.data.CompInfo) != 1 {
		t.Fatalf("expected 1 comp info row, got %d", len(msg.data.CompInfo))
	}
	if len(msg.data.Teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(msg.data.Teams))
	}
	first := msg.data.Teams[0]
	if first.Player != "Alice" || first.Record != "7-1" || first.Tournament != "Worlds 2026" {
		t.Errorf("unexpected first team: %+v", first)
	}
	if len(first.Pokemon) != 3 || first.WebURL != "https://example.com/1" {
		t.Errorf("unexpected first team detail: %+v", first)
	}
	if len(msg.data.Usage) != 2 {
		t.Fatalf("expected 2 usage rows, got %d", len(msg.data.Usage))
	}
	if u := msg.data.Usage[0]; u.Pokemon != "Basculegion" || u.UsagePercent != 51.5 {
		t.Errorf("unexpected usage row: %+v", u)
	}
	if len(msg.data.SpeedTiers) != 1 {
		t.Fatalf("expected 1 speed tier row, got %d", len(msg.data.SpeedTiers))
	}
	if st := msg.data.SpeedTiers[0]; st.Pokemon != "Mega Aerodactyl" || st.BaseSpe != 150 || st.MaxScarf != 333 {
		t.Errorf("unexpected speed tier row: %+v", st)
	}
	want := []string{compInfoURL, topTeamsURL, usageURL, speedTiersURL}
	if len(capturedURLs) != len(want) {
		t.Fatalf("expected %d fetches, got %v", len(want), capturedURLs)
	}
	for i, u := range want {
		if capturedURLs[i] != u {
			t.Errorf("fetch %d = %q, want %q", i, capturedURLs[i], u)
		}
	}
}

func TestFetchDashboardData_ConnectionError(t *testing.T) {
	conn := func(_ string) ([]byte, error) { return nil, errors.New("refused") }
	msg := fetchDashboardData(conn)().(dataMsg)
	if msg.err == nil {
		t.Error("expected error")
	}
	if msg.data != nil {
		t.Error("expected nil data on connection error")
	}
}

func TestFetchDashboardData_InvalidJSON(t *testing.T) {
	conn := func(_ string) ([]byte, error) { return []byte("not json"), nil }
	msg := fetchDashboardData(conn)().(dataMsg)
	if msg.err == nil {
		t.Error("expected unmarshal error")
	}
	if msg.data != nil {
		t.Error("expected nil data on invalid json")
	}
}
