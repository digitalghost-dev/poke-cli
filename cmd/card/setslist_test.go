package card

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestSetsModel_Init(t *testing.T) {
	model, _ := SetsList("sv")

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return commands (spinner tick + fetch)")
	}
}

func TestSetsModel_Update_EscKey(t *testing.T) {
	items := []list.Item{
		styling.Item("Scarlet & Violet"),
		styling.Item("Paldea Evolved"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
		List:       l,
		SeriesName: "sv",
		Quitting:   false,
	}

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(setsModel)
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
		styling.Item("Scarlet & Violet"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
		List:       l,
		SeriesName: "sv",
		Quitting:   false,
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(setsModel)
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
		styling.Item("Scarlet & Violet"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
		List:       l,
		SeriesName: "sv",
	}

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, cmd := model.Update(msg)

	resultModel, ok := newModel.(setsModel)
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
		styling.Item("Scarlet & Violet"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
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
		styling.Item("Scarlet & Violet"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
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
		styling.Item("Scarlet & Violet"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	model := setsModel{
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
		styling.Item("Scarlet & Violet"),
		styling.Item("Paldea Evolved"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 20)

	setsIDMap := map[string]string{
		"Scarlet & Violet": "sv01",
		"Paldea Evolved":   "sv02",
	}

	model := setsModel{
		List:      l,
		SetsIDMap: setsIDMap,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Update with Enter should return tea.Quit command")
	}
}

func TestSetsList_Success(t *testing.T) {
	model, err := SetsList("sv")
	if err != nil {
		t.Fatalf("SetsList returned error: %v", err)
	}

	// SetsList now returns minimal model with Loading=true
	if model.SeriesName != "sv" {
		t.Errorf("expected SeriesName 'sv', got %s", model.SeriesName)
	}

	if !model.Loading {
		t.Error("expected Loading to be true")
	}

	// View should show loading spinner
	if model.View() == "" {
		t.Error("model view should render loading state")
	}
}

func TestSetsDataMsg_PopulatesModel(t *testing.T) {
	// Start with a loading model
	model, _ := SetsList("sv")

	// Simulate receiving data via setsDataMsg
	msg := setsDataMsg{
		items: []list.Item{
			styling.Item("Scarlet & Violet"),
			styling.Item("Paldea Evolved"),
		},
		setsIDMap: map[string]string{
			"Scarlet & Violet": "sv01",
			"Paldea Evolved":   "sv02",
		},
		seriesID: "sv",
	}

	newModel, _ := model.Update(msg)
	resultModel := newModel.(setsModel)

	if resultModel.Loading {
		t.Error("Loading should be false after receiving data")
	}

	if resultModel.SetsIDMap["Scarlet & Violet"] != "sv01" {
		t.Errorf("expected SetsIDMap['Scarlet & Violet'] = 'sv01', got %s", resultModel.SetsIDMap["Scarlet & Violet"])
	}
	if resultModel.SetsIDMap["Paldea Evolved"] != "sv02" {
		t.Errorf("expected SetsIDMap['Paldea Evolved'] = 'sv02', got %s", resultModel.SetsIDMap["Paldea Evolved"])
	}
}

func TestSetsDataMsg_Error_StoresError(t *testing.T) {
	model, _ := SetsList("sv")

	// Simulate receiving an error via setsDataMsg
	msg := setsDataMsg{
		err: errors.New("network error"),
	}

	newModel, cmd := model.Update(msg)
	resultModel := newModel.(setsModel)

	if resultModel.Error == nil {
		t.Error("Error should be set when error received")
	}

	if resultModel.Error.Error() != "network error" {
		t.Errorf("Expected error message 'network error', got '%s'", resultModel.Error.Error())
	}

	if resultModel.Loading {
		t.Error("Loading should be false after error")
	}

	if cmd != nil {
		t.Error("Should not return a command on error (stays in view to show error)")
	}
}

func TestSetsDataMsg_EmptyResult(t *testing.T) {
	model, _ := SetsList("sv")

	// Simulate receiving empty data
	msg := setsDataMsg{
		items:     []list.Item{},
		setsIDMap: map[string]string{},
		seriesID:  "sv",
	}

	newModel, _ := model.Update(msg)
	resultModel := newModel.(setsModel)

	if resultModel.Loading {
		t.Error("Loading should be false after receiving data")
	}

	if len(resultModel.SetsIDMap) != 0 {
		t.Errorf("expected empty SetsIDMap, got %d entries", len(resultModel.SetsIDMap))
	}
}
