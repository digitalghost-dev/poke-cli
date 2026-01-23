package card

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// testSupabaseKey is the publishable API key used in tests.
// Extracted to a constant for easier maintenance if the key changes.
const testSupabaseKey = "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j"

func TestCardsModel_Init(t *testing.T) {
	model := CardsModel{
		SeriesName: "sv",
	}

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestCardsModel_Update_EscKey(t *testing.T) {
	rows := []table.Row{
		{"001/198 - Bulbasaur"},
		{"002/198 - Ivysaur"},
	}
	columns := []table.Column{
		{Title: "Card Name", Width: 35},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	model := CardsModel{
		Table:    tbl,
		Quitting: false,
	}

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := model.Update(msg)

	resultModel := newModel.(CardsModel)

	if !resultModel.Quitting {
		t.Error("Quitting should be set to true when ESC is pressed")
	}

	if cmd == nil {
		t.Error("Update with ESC should return tea.Quit command")
	}
}

func TestCardsModel_Update_CtrlC(t *testing.T) {
	rows := []table.Row{
		{"001/198 - Bulbasaur"},
	}
	columns := []table.Column{
		{Title: "Card Name", Width: 35},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	model := CardsModel{
		Table:    tbl,
		Quitting: false,
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(msg)

	resultModel := newModel.(CardsModel)

	if !resultModel.Quitting {
		t.Error("Quitting should be set to true when Ctrl+C is pressed")
	}

	if cmd == nil {
		t.Error("Update with Ctrl+C should return tea.Quit command")
	}
}

func TestCardsModel_Update_TabTogglesSearchFocusAndTableSelectedBackground(t *testing.T) {
	rows := []table.Row{{"001/198 - Bulbasaur"}}
	columns := []table.Column{{Title: "Card Name", Width: 35}}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	search := textinput.New()
	search.Blur()

	initialStyles := cardTableStyles(activeTableSelectedBg)
	tbl.SetStyles(initialStyles)

	model := CardsModel{
		Search:      search,
		Table:       tbl,
		TableStyles: initialStyles,
	}

	// Tab into the search bar.
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	m1 := newModel.(CardsModel)
	if !m1.Search.Focused() {
		t.Fatal("expected search to be focused after tab")
	}

	bg1 := m1.TableStyles.Selected.GetBackground()
	r1, g1, b1, a1 := bg1.RGBA()
	re, ge, be, ae := inactiveTableSelectedBg.RGBA()
	if r1 != re || g1 != ge || b1 != be || a1 != ae {
		t.Fatalf("expected selected background to be gray when searching; got RGBA(%d,%d,%d,%d)", r1, g1, b1, a1)
	}

	// Tab back to the table.
	newModel2, _ := m1.Update(tea.KeyMsg{Type: tea.KeyTab})
	m2 := newModel2.(CardsModel)
	if m2.Search.Focused() {
		t.Fatal("expected search to be blurred after tabbing back")
	}

	bg2 := m2.TableStyles.Selected.GetBackground()
	r2, g2, b2, a2 := bg2.RGBA()
	re2, ge2, be2, ae2 := activeTableSelectedBg.RGBA()
	if r2 != re2 || g2 != ge2 || b2 != be2 || a2 != ae2 {
		t.Fatalf("expected selected background to be yellow when table is focused; got RGBA(%d,%d,%d,%d)", r2, g2, b2, a2)
	}
}

func TestCardsModel_Update_ViewImageKey_QuestionMark(t *testing.T) {
	rows := []table.Row{
		{"001/198 - Bulbasaur"},
	}
	columns := []table.Column{
		{Title: "Card Name", Width: 35},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	model := CardsModel{
		Table:     tbl,
		ViewImage: false,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, cmd := model.Update(msg)

	resultModel := newModel.(CardsModel)

	if !resultModel.ViewImage {
		t.Error("ViewImage should be set to true when '?' is pressed")
	}

	if cmd == nil {
		t.Error("Update with '?' should return tea.Quit command")
	}
}

func TestCardsModel_Update_ViewImageKey_DoesNotOverrideSearch(t *testing.T) {
	rows := []table.Row{{"001/198 - Bulbasaur"}}
	columns := []table.Column{{Title: "Card Name", Width: 35}}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	search := textinput.New()
	search.Focus()

	model := CardsModel{
		Search:    search,
		Table:     tbl,
		ViewImage: false,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := model.Update(msg)
	resultModel := newModel.(CardsModel)

	if resultModel.ViewImage {
		t.Fatal("expected ViewImage to remain false when typing '?' in the search field")
	}
	if resultModel.Quitting {
		t.Fatal("expected Quitting to remain false when typing in the search field")
	}
	if got := resultModel.Search.Value(); got != "?" {
		t.Fatalf("expected search input to receive '?'; got %q", got)
	}
}

func TestCardsModel_View_Quitting(t *testing.T) {
	model := CardsModel{
		Quitting: true,
	}

	result := model.View()

	if !strings.Contains(result, "Quitting card search") {
		t.Errorf("View() when quitting should contain 'Quitting card search', got: %s", result)
	}
}

func TestCardsModel_View_PriceDisplay(t *testing.T) {
	rows := []table.Row{
		{"001/198 - Bulbasaur"},
	}
	columns := []table.Column{
		{Title: "Card Name", Width: 35},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	priceMap := map[string]string{
		"001/198 - Bulbasaur": "Price: $1.50",
	}

	model := CardsModel{
		Table:    tbl,
		PriceMap: priceMap,
		Quitting: false,
	}

	result := model.View()

	// The view should include the card name
	if !strings.Contains(result, "001/198 - Bulbasaur") {
		t.Error("View() should display selected card name")
	}
}

func TestCardsModel_View_MissingPrice(t *testing.T) {
	rows := []table.Row{
		{"001/198 - Bulbasaur"},
	}
	columns := []table.Column{
		{Title: "Card Name", Width: 35},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Empty price map - simulates missing price data
	priceMap := map[string]string{}

	model := CardsModel{
		Table:    tbl,
		PriceMap: priceMap,
		Quitting: false,
	}

	result := model.View()

	// Should show "Price: Not available" when price is missing
	if !strings.Contains(result, "Price: Not available") {
		t.Error("View() should display 'Price: Not available' for cards without pricing")
	}
}

func TestCardsList_SuccessAndFallbacks(t *testing.T) {
	// Save and restore getCardData stub
	original := getCardData
	defer func() { getCardData = original }()

	var capturedURL string
	getCardData = func(url string) ([]byte, error) {
		capturedURL = url
		// Return two cards: one with price + illustrator, one with fallbacks
		json := `[
            {"number_plus_name":"001/198 - Bulbasaur","market_price":1.5,"image_url":"https://example.com/bulba.png","illustrator":"Ken Sugimori"},
            {"number_plus_name":"002/198 - Ivysaur","market_price":0,"image_url":"https://example.com/ivy.png","illustrator":""}
        ]`
		return []byte(json), nil
	}

	model, err := CardsList("set123")
	if err != nil {
		t.Fatalf("CardsList returned error: %v", err)
	}

	// URL should target the correct set id and select fields
	if !strings.Contains(capturedURL, "set_id=eq.set123") {
		t.Errorf("expected URL to contain set_id filter, got %s", capturedURL)
	}
	if !strings.Contains(capturedURL, "select=number_plus_name,market_price,image_url,illustrator") {
		t.Errorf("expected URL to contain select fields, got %s", capturedURL)
	}

	// PriceMap expectations
	if got := model.PriceMap["001/198 - Bulbasaur"]; got != "Price: $1.50" {
		t.Errorf("unexpected price for Bulbasaur: %s", got)
	}
	if got := model.PriceMap["002/198 - Ivysaur"]; got != "Pricing not available" {
		t.Errorf("unexpected price for Ivysaur: %s", got)
	}

	// IllustratorMap expectations
	if got := model.IllustratorMap["001/198 - Bulbasaur"]; got != "Illustrator: Ken Sugimori" {
		t.Errorf("unexpected illustrator for Bulbasaur: %s", got)
	}
	if got := model.IllustratorMap["002/198 - Ivysaur"]; got != "Illustrator not available" {
		t.Errorf("unexpected illustrator for Ivysaur: %s", got)
	}

	// Image map
	if model.ImageMap["001/198 - Bulbasaur"] != "https://example.com/bulba.png" {
		t.Errorf("unexpected image url for Bulbasaur: %s", model.ImageMap["001/198 - Bulbasaur"])
	}
	if model.ImageMap["002/198 - Ivysaur"] != "https://example.com/ivy.png" {
		t.Errorf("unexpected image url for Ivysaur: %s", model.ImageMap["002/198 - Ivysaur"])
	}

	if row := model.Table.SelectedRow(); len(row) == 0 {
		if model.View() == "" {
			t.Error("model view should render even if no row is selected")
		}
	}
}

func TestCardsList_FetchError(t *testing.T) {
	original := getCardData
	defer func() { getCardData = original }()

	getCardData = func(url string) ([]byte, error) {
		return nil, errors.New("network error")
	}

	_, err := CardsList("set123")
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
}

func TestCardsList_BadJSON(t *testing.T) {
	original := getCardData
	defer func() { getCardData = original }()

	getCardData = func(url string) ([]byte, error) {
		return []byte("not-json"), nil
	}

	_, err := CardsList("set123")
	if err == nil {
		t.Fatal("expected error for bad JSON payload")
	}
}

func TestCardsList_EmptyResult(t *testing.T) {
	original := getCardData
	defer func() { getCardData = original }()

	getCardData = func(url string) ([]byte, error) {
		return []byte("[]"), nil
	}

	model, err := CardsList("set123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(model.PriceMap) != 0 || len(model.IllustratorMap) != 0 || len(model.ImageMap) != 0 {
		t.Errorf("expected empty maps, got price:%d illus:%d image:%d", len(model.PriceMap), len(model.IllustratorMap), len(model.ImageMap))
	}

	if model.View() == "" {
		t.Error("expected view to render with empty data")
	}
}

func TestCallCardData_SendsHeadersAndReturnsBody(t *testing.T) {
	// Start a test HTTP server that validates headers and returns a body
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

	body, err := CallCardData(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestCallCardData_Non200Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := CallCardData(srv.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
	if !strings.Contains(err.Error(), "unexpected status code: 500") {
		t.Fatalf("error should mention status code, got: %v", err)
	}
}

func TestCallCardData_BadURL(t *testing.T) {
	_, err := CallCardData("http://%41:80/") // invalid URL host
	if err == nil {
		t.Fatal("expected error for bad URL")
	}
}
