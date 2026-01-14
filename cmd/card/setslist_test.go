package card

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestSetsModel_Init(t *testing.T) {
	model := SetsModel{
		SeriesName: "sv",
		Quitting:   false,
	}

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestSetsModel_Update_EscKey(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
		item("Paldea Evolved"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:       l,
		SeriesName: "sv",
		Quitting:   false,
	}

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(SetsModel)
	if !ok {
		t.Fatalf("expected SetsModel, got %T", newModel)
	}

	if !resultModel.Quitting {
		t.Error("Quitting should be set to true when ESC is pressed")
	}

	if cmd == nil {
		t.Error("Update with ESC should return tea.Quit command")
	}
}

func TestSetsModel_Update_CtrlC(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:       l,
		SeriesName: "sv",
		Quitting:   false,
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(SetsModel)
	if !ok {
		t.Fatalf("expected SetsModel, got %T", newModel)
	}

	if !resultModel.Quitting {
		t.Error("Quitting should be set to true when Ctrl+C is pressed")
	}

	if cmd == nil {
		t.Error("Update with Ctrl+C should return tea.Quit command")
	}
}

func TestSetsModel_Update_WindowSizeMsg(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:       l,
		SeriesName: "sv",
	}

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(SetsModel)
	if !ok {
		t.Fatalf("expected SetsModel, got %T", newModel)
	}

	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}

	if resultModel.Quitting {
		t.Error("WindowSizeMsg should not set Quitting to true")
	}
}

func TestSetsModel_View_Quitting(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:     l,
		Quitting: true,
	}

	result := model.View()

	if !strings.Contains(result, "Quitting card search") {
		t.Errorf("View() when quitting should contain 'Quitting card search', got: %s", result)
	}
}

func TestSetsModel_View_ChoiceMade(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:   l,
		Choice: "Scarlet & Violet",
	}

	result := model.View()

	if !strings.Contains(result, "Set selected: Scarlet & Violet") {
		t.Errorf("View() with choice should contain 'Set selected: Scarlet & Violet', got: %s", result)
	}
}

func TestSetsModel_View_Normal(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	model := SetsModel{
		List:     l,
		Quitting: false,
		Choice:   "",
	}

	result := model.View()

	if result == "" {
		t.Error("View() should not return empty string in normal state")
	}
}

func TestSetsModel_Update_EnterKey(t *testing.T) {
	items := []list.Item{
		item("Scarlet & Violet"),
		item("Paldea Evolved"),
	}
	l := list.New(items, itemDelegate{}, 20, 20)

	setsIDMap := map[string]string{
		"Scarlet & Violet": "sv01",
		"Paldea Evolved":   "sv02",
	}

	model := SetsModel{
		List:      l,
		setsIDMap: setsIDMap,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Update with Enter should return tea.Quit command")
	}
}

func TestCallSetsData_SendsHeadersAndReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("apikey"); got != testSupabaseKey {
			t.Fatalf("missing or wrong apikey header: %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+testSupabaseKey {
			t.Fatalf("missing or wrong Authorization header: %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("missing or wrong Content-Type header: %q", got)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, err := callSetsData(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestCallSetsData_Non200Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := callSetsData(srv.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
	if !strings.Contains(err.Error(), "unexpected status code: 500") {
		t.Fatalf("error should mention status code, got: %v", err)
	}
}

func TestCallSetsData_BadURL(t *testing.T) {
	_, err := callSetsData("http://%41:80/") // invalid URL host
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}

func TestSetsList_Success(t *testing.T) {
	original := getSetsData
	defer func() { getSetsData = original }()

	getSetsData = func(url string) ([]byte, error) {
		json := `[
			{"series_id":"sv","set_id":"sv01","set_name":"Scarlet & Violet","official_card_count":198,"total_card_count":258,"logo":"https://example.com/sv01.png","symbol":"https://example.com/sv01-symbol.png"},
			{"series_id":"sv","set_id":"sv02","set_name":"Paldea Evolved","official_card_count":193,"total_card_count":279,"logo":"https://example.com/sv02.png","symbol":"https://example.com/sv02-symbol.png"},
			{"series_id":"swsh","set_id":"swsh01","set_name":"Sword & Shield","official_card_count":202,"total_card_count":216,"logo":"https://example.com/swsh01.png","symbol":"https://example.com/swsh01-symbol.png"}
		]`
		return []byte(json), nil
	}

	model, err := SetsList("sv")
	if err != nil {
		t.Fatalf("SetsList returned error: %v", err)
	}

	// Should only have 2 sets (filtered by series_id "sv")
	if model.SeriesName != "sv" {
		t.Errorf("expected SeriesName 'sv', got %s", model.SeriesName)
	}

	// Check setsIDMap has correct mappings
	if model.setsIDMap["Scarlet & Violet"] != "sv01" {
		t.Errorf("expected setsIDMap['Scarlet & Violet'] = 'sv01', got %s", model.setsIDMap["Scarlet & Violet"])
	}
	if model.setsIDMap["Paldea Evolved"] != "sv02" {
		t.Errorf("expected setsIDMap['Paldea Evolved'] = 'sv02', got %s", model.setsIDMap["Paldea Evolved"])
	}

	// swsh set should not be in the map
	if _, exists := model.setsIDMap["Sword & Shield"]; exists {
		t.Error("Sword & Shield should not be in setsIDMap (different series)")
	}

	if model.View() == "" {
		t.Error("model view should render")
	}
}

func TestSetsList_FetchError(t *testing.T) {
	original := getSetsData
	defer func() { getSetsData = original }()

	getSetsData = func(url string) ([]byte, error) {
		return nil, errors.New("network error")
	}

	_, err := SetsList("sv")
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
	if !strings.Contains(err.Error(), "error getting sets data") {
		t.Errorf("error should mention 'error getting sets data', got: %v", err)
	}
}

func TestSetsList_BadJSON(t *testing.T) {
	original := getSetsData
	defer func() { getSetsData = original }()

	getSetsData = func(url string) ([]byte, error) {
		return []byte("not-json"), nil
	}

	_, err := SetsList("sv")
	if err == nil {
		t.Fatal("expected error for bad JSON payload")
	}
	if !strings.Contains(err.Error(), "error parsing sets data") {
		t.Errorf("error should mention 'error parsing sets data', got: %v", err)
	}
}

func TestSetsList_EmptyResult(t *testing.T) {
	original := getSetsData
	defer func() { getSetsData = original }()

	getSetsData = func(url string) ([]byte, error) {
		return []byte("[]"), nil
	}

	model, err := SetsList("sv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(model.setsIDMap) != 0 {
		t.Errorf("expected empty setsIDMap, got %d entries", len(model.setsIDMap))
	}

	if model.View() == "" {
		t.Error("expected view to render with empty data")
	}
}

func TestSetsList_NoMatchingSeries(t *testing.T) {
	original := getSetsData
	defer func() { getSetsData = original }()

	getSetsData = func(url string) ([]byte, error) {
		json := `[
			{"series_id":"swsh","set_id":"swsh01","set_name":"Sword & Shield","official_card_count":202,"total_card_count":216,"logo":"","symbol":""}
		]`
		return []byte(json), nil
	}

	model, err := SetsList("sv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No sets match "sv" series
	if len(model.setsIDMap) != 0 {
		t.Errorf("expected empty setsIDMap when no series match, got %d entries", len(model.setsIDMap))
	}
}
