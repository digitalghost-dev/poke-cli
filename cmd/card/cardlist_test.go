package card

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCardsModel_Init(t *testing.T) {
	model, _ := CardsList("sv01")

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return commands (spinner tick + fetch)")
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

func TestCardsList_ReturnsLoadingModel(t *testing.T) {
	model, err := CardsList("set123")
	if err != nil {
		t.Fatalf("CardsList returned error: %v", err)
	}

	if model.SetID != "set123" {
		t.Errorf("expected SetID 'set123', got %s", model.SetID)
	}

	if !model.Loading {
		t.Error("expected Loading to be true")
	}

	// View should show loading spinner
	view := model.View()
	if !strings.Contains(view, "Loading cards") {
		t.Errorf("expected view to show loading state, got: %s", view)
	}
}

func TestCardDataMsg_PopulatesModel(t *testing.T) {
	model, _ := CardsList("set123")

	// Simulate receiving data via cardDataMsg
	msg := cardDataMsg{
		allRows: []table.Row{
			{"001/198 - Bulbasaur"},
			{"002/198 - Ivysaur"},
		},
		priceMap: map[string]string{
			"001/198 - Bulbasaur": "Price: $1.50",
			"002/198 - Ivysaur":   "Pricing not available",
		},
		imageMap: map[string]string{
			"001/198 - Bulbasaur": "https://example.com/bulba.png",
			"002/198 - Ivysaur":   "https://example.com/ivy.png",
		},
		illustratorMap: map[string]string{
			"001/198 - Bulbasaur": "Illustrator: Ken Sugimori",
			"002/198 - Ivysaur":   "Illustrator not available",
		},
		regulationMarkMap: map[string]string{},
	}

	newModel, _ := model.Update(msg)
	resultModel := newModel.(CardsModel)

	if resultModel.Loading {
		t.Error("Loading should be false after receiving data")
	}

	// PriceMap expectations
	if got := resultModel.PriceMap["001/198 - Bulbasaur"]; got != "Price: $1.50" {
		t.Errorf("unexpected price for Bulbasaur: %s", got)
	}
	if got := resultModel.PriceMap["002/198 - Ivysaur"]; got != "Pricing not available" {
		t.Errorf("unexpected price for Ivysaur: %s", got)
	}

	// IllustratorMap expectations
	if got := resultModel.IllustratorMap["001/198 - Bulbasaur"]; got != "Illustrator: Ken Sugimori" {
		t.Errorf("unexpected illustrator for Bulbasaur: %s", got)
	}
	if got := resultModel.IllustratorMap["002/198 - Ivysaur"]; got != "Illustrator not available" {
		t.Errorf("unexpected illustrator for Ivysaur: %s", got)
	}

	// Image map
	if resultModel.ImageMap["001/198 - Bulbasaur"] != "https://example.com/bulba.png" {
		t.Errorf("unexpected image url for Bulbasaur: %s", resultModel.ImageMap["001/198 - Bulbasaur"])
	}
	if resultModel.ImageMap["002/198 - Ivysaur"] != "https://example.com/ivy.png" {
		t.Errorf("unexpected image url for Ivysaur: %s", resultModel.ImageMap["002/198 - Ivysaur"])
	}
}

func TestCardDataMsg_Error_QuitsModel(t *testing.T) {
	model, _ := CardsList("set123")

	// Simulate receiving an error via cardDataMsg
	msg := cardDataMsg{
		err: errors.New("network error"),
	}

	newModel, cmd := model.Update(msg)
	resultModel := newModel.(CardsModel)

	if !resultModel.Quitting {
		t.Error("Quitting should be true when error received")
	}

	if cmd == nil {
		t.Error("Should return tea.Quit command on error")
	}
}

func TestCardDataMsg_EmptyResult(t *testing.T) {
	model, _ := CardsList("set123")

	// Simulate receiving empty data
	msg := cardDataMsg{
		allRows:           []table.Row{},
		priceMap:          map[string]string{},
		imageMap:          map[string]string{},
		illustratorMap:    map[string]string{},
		regulationMarkMap: map[string]string{},
	}

	newModel, _ := model.Update(msg)
	resultModel := newModel.(CardsModel)

	if resultModel.Loading {
		t.Error("Loading should be false after receiving data")
	}

	if len(resultModel.PriceMap) != 0 || len(resultModel.IllustratorMap) != 0 || len(resultModel.ImageMap) != 0 {
		t.Errorf("expected empty maps, got price:%d illus:%d image:%d", len(resultModel.PriceMap), len(resultModel.IllustratorMap), len(resultModel.ImageMap))
	}
}

