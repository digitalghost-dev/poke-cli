package card

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestItemFilterValue(t *testing.T) {
	testItem := item("Test Item")
	filterValue := testItem.FilterValue()

	if filterValue != "" {
		t.Errorf("Expected FilterValue to return empty string, got '%s'", filterValue)
	}
}

func TestItemDelegateHeight(t *testing.T) {
	delegate := itemDelegate{}
	height := delegate.Height()

	if height != 1 {
		t.Errorf("Expected Height to return 1, got %d", height)
	}
}

func TestItemDelegateSpacing(t *testing.T) {
	delegate := itemDelegate{}
	spacing := delegate.Spacing()

	if spacing != 0 {
		t.Errorf("Expected Spacing to return 0, got %d", spacing)
	}
}

func TestItemDelegateUpdate(t *testing.T) {
	delegate := itemDelegate{}
	cmd := delegate.Update(tea.KeyMsg{}, &list.Model{})

	if cmd != nil {
		t.Error("Expected Update to return nil, got non-nil value")
	}
}

func TestItemDelegateRender(t *testing.T) {
	delegate := itemDelegate{}

	items := []list.Item{
		item("First Item"),
		item("Second Item"),
		item("Third Item"),
	}

	l := list.New(items, delegate, 20, 10)

	var buf bytes.Buffer
	delegate.Render(&buf, l, 0, items[0])

	output := buf.String()
	if output == "" {
		t.Error("Expected non-empty output from Render")
	}
}

func TestItemDelegateRenderSelected(t *testing.T) {
	delegate := itemDelegate{}

	items := []list.Item{
		item("First Item"),
		item("Second Item"),
	}

	l := list.New(items, delegate, 20, 10)

	var buf bytes.Buffer
	delegate.Render(&buf, l, l.Index(), items[l.Index()])

	output := buf.String()
	if output == "" {
		t.Error("Expected non-empty output for selected item")
	}

	if len(output) == 0 {
		t.Error("Selected item should produce rendered output")
	}
}
