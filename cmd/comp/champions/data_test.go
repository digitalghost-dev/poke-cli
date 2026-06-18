package champions

import (
	"errors"
	"testing"
)

func TestFetchDashboardData_Success(t *testing.T) {
	var capturedURL string
	conn := func(url string) ([]byte, error) {
		capturedURL = url
		return []byte(`[
			{"author":"Alice","record":"7-1","tournament":"Worlds 2026","archetypes":["Hyper Offense"],"pokemon":["Miraidon","Flutter Mane","Iron Hands"],"web_url":"https://example.com/1"},
			{"author":"Bob","record":"6-2","tournament":"Regional","archetypes":[],"pokemon":["Calyrex"],"web_url":""}
		]`), nil
	}

	msg := fetchDashboardData(conn)().(dataMsg)
	if msg.err != nil {
		t.Fatalf("unexpected error: %v", msg.err)
	}
	if msg.data == nil {
		t.Fatal("expected data, got nil")
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
	if capturedURL != topTeamsURL {
		t.Errorf("expected topTeamsURL to be fetched, got %q", capturedURL)
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
