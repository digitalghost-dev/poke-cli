package card

import (
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
