package card

import (
	"testing"
	"time"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func TestSeriesModelInit(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)
	model := seriesModel{List: l}

	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Expected Init() to return nil, got %v", cmd)
	}
}

func TestSeriesModelQuit(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)
	model := seriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Test ctrl+c quit
	testModel.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(seriesModel)

	if !final.Quitting {
		t.Errorf("Expected model to be quitting after ctrl+c")
	}
}

func TestSeriesModelEscQuit(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)
	model := seriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Test esc quit
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyEscape})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(seriesModel)

	if !final.Quitting {
		t.Errorf("Expected model to be quitting after esc")
	}
}

func TestSeriesModelSelection(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)
	model := seriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Navigate and select
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyDown})  // Move to second item
	testModel.Send(tea.KeyPressMsg{Code: tea.KeyEnter}) // Select it
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(seriesModel)

	if final.Choice == "" {
		t.Errorf("Expected a choice to be made, got empty string")
	}
	if final.Choice != "Scarlet & Violet" {
		t.Errorf("Expected choice to be 'Scarlet & Violet', got '%s'", final.Choice)
	}
}

func TestSeriesModelWindowResize(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)
	model := seriesModel{List: l}

	// Send window resize message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	finalModel := updatedModel.(seriesModel)

	if finalModel.List.Width() != 100 {
		t.Errorf("Expected list width to be 100 after resize, got %d", finalModel.List.Width())
	}
}

func TestSeriesModelView(t *testing.T) {
	items := []list.Item{
		styling.Item("Mega Evolution"),
		styling.Item("Scarlet & Violet"),
		styling.Item("Sword & Shield"),
	}
	l := list.New(items, styling.ItemDelegate{}, 20, 12)

	// Test normal view
	model := seriesModel{List: l}
	view := model.View()
	if view.Content == "" {
		t.Errorf("Expected non-empty view, got empty string")
	}

	// Test quitting view
	model.Quitting = true
	view = model.View()
	if view.Content != "\n  Quitting card search...\n\n" {
		t.Errorf("Expected quitting message, got '%s'", view.Content)
	}

	// Test choice made view
	model.Quitting = false
	model.Choice = "Scarlet & Violet"
	view = model.View()
	if view.Content == "" {
		t.Errorf("Expected non-empty view for choice, got empty string")
	}
}

func TestSeriesList(t *testing.T) {
	model := SeriesList()
	items := model.List.Items()

	// Check that list has 3 items
	if items == nil {
		t.Error("SeriesList() should create a list with items")
	}

	if len(items) != 4 {
		t.Errorf("Expected 4 items, got %d", len(items))
	}

	// Verify all four series are present
	expectedSeries := []string{"Mega Evolution", "Scarlet & Violet", "Sword & Shield", "Sun & Moon"}
	for i, expected := range expectedSeries {
		itemStr := string(items[i].(styling.Item))
		if itemStr != expected {
			t.Errorf("Expected item %d to be '%s', got '%s'", i, expected, itemStr)
		}
	}
}
