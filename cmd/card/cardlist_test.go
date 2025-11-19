package card

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

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

func TestCardsModel_Update_SpaceBar(t *testing.T) {
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

	msg := tea.KeyMsg{Type: tea.KeySpace}
	newModel, cmd := model.Update(msg)

	resultModel := newModel.(CardsModel)

	if !resultModel.ViewImage {
		t.Error("ViewImage should be set to true when SPACE is pressed")
	}

	if cmd == nil {
		t.Error("Update with SPACE should return tea.Quit command")
	}
}

func TestCardsModel_Update_SelectionSync(t *testing.T) {
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
		Table:          tbl,
		SelectedOption: "",
	}

	// Simulate a key press that won't quit (e.g., arrow down)
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(msg)

	resultModel := newModel.(CardsModel)

	// The selected option should be updated to the current row
	if resultModel.SelectedOption == "" {
		t.Error("SelectedOption should be synced after table update")
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

func TestCardsModel_View_Normal(t *testing.T) {
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

	if result == "" {
		t.Error("View() should not return empty string in normal state")
	}

	// Check that it contains the key menu
	if !strings.Contains(result, "move up") {
		t.Error("View() should contain key menu instructions")
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
