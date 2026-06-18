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
	if len(capturedURLs) != 2 || capturedURLs[0] != compInfoURL || capturedURLs[1] != topTeamsURL {
		t.Errorf("expected compInfoURL then topTeamsURL, got %v", capturedURLs)
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
