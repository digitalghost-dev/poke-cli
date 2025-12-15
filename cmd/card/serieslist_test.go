package card

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestSeriesModelInit(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)
	model := SeriesModel{List: l}

	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Expected Init() to return nil, got %v", cmd)
	}
}

func TestSeriesModelQuit(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)
	model := SeriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Test ctrl+c quit
	testModel.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(SeriesModel)

	if !final.Quitting {
		t.Errorf("Expected model to be quitting after ctrl+c")
	}
}

func TestSeriesModelEscQuit(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)
	model := SeriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Test esc quit
	testModel.Send(tea.KeyMsg{Type: tea.KeyEsc})
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(SeriesModel)

	if !final.Quitting {
		t.Errorf("Expected model to be quitting after esc")
	}
}

func TestSeriesModelSelection(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)
	model := SeriesModel{List: l}

	testModel := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(80, 24))

	// Navigate and select
	testModel.Send(tea.KeyMsg{Type: tea.KeyDown})  // Move to second item
	testModel.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Select it
	testModel.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))

	final := testModel.FinalModel(t).(SeriesModel)

	if final.Choice == "" {
		t.Errorf("Expected a choice to be made, got empty string")
	}
	if final.Choice != "Scarlet & Violet" {
		t.Errorf("Expected choice to be 'Scarlet & Violet', got '%s'", final.Choice)
	}
}

func TestSeriesModelWindowResize(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)
	model := SeriesModel{List: l}

	// Send window resize message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	finalModel := updatedModel.(SeriesModel)

	if finalModel.List.Width() != 100 {
		t.Errorf("Expected list width to be 100 after resize, got %d", finalModel.List.Width())
	}
}

func TestSeriesModelView(t *testing.T) {
	items := []list.Item{
		item("Mega Evolution"),
		item("Scarlet & Violet"),
		item("Sword & Shield"),
	}
	l := list.New(items, itemDelegate{}, 20, 12)

	// Test normal view
	model := SeriesModel{List: l}
	view := model.View()
	if view == "" {
		t.Errorf("Expected non-empty view, got empty string")
	}

	// Test quitting view
	model.Quitting = true
	view = model.View()
	if view != "\n  Quitting card search...\n\n" {
		t.Errorf("Expected quitting message, got '%s'", view)
	}

	// Test choice made view
	model.Quitting = false
	model.Choice = "Scarlet & Violet"
	view = model.View()
	if view == "" {
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

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify all three series are present
	expectedSeries := []string{"Mega Evolution", "Scarlet & Violet", "Sword & Shield"}
	for i, expected := range expectedSeries {
		itemStr := string(items[i].(item))
		if itemStr != expected {
			t.Errorf("Expected item %d to be '%s', got '%s'", i, expected, itemStr)
		}
	}
}
