package shell

import (
	"strings"
	"testing"
)

func TestNewStyles(t *testing.T) {
	s := NewStyles()
	if s == nil {
		t.Fatal("expected non-nil styles")
	}
	if s.HighlightColor == nil {
		t.Error("expected highlight color to be set")
	}
}

func TestRender_ContainsTabsContentAndMenu(t *testing.T) {
	s := NewStyles()
	tabs := []string{"Overview", "Standings", "Usage", "Countries"}
	var gotWidth int
	body := s.Render(tabs, 0, 120, func(contentWidth int) string {
		gotWidth = contentWidth
		return "TAB-BODY"
	})

	for _, tab := range tabs {
		if !strings.Contains(body, tab) {
			t.Errorf("expected rendered frame to contain tab %q", tab)
		}
	}
	if !strings.Contains(body, "TAB-BODY") {
		t.Error("expected rendered frame to contain the content from the closure")
	}
	if !strings.Contains(body, "back") || !strings.Contains(body, "quit") {
		t.Error("expected the key menu in the rendered frame")
	}
	if gotWidth <= 0 {
		t.Errorf("expected a positive content width passed to renderContent, got %d", gotWidth)
	}
}

func TestRender_NarrowWidthDoesNotPanic(t *testing.T) {
	s := NewStyles()
	body := s.Render([]string{"A", "B"}, 1, 10, func(int) string { return "x" })
	if body == "" {
		t.Error("expected non-empty frame even at narrow width")
	}
}

func TestTableStyles(t *testing.T) {
	st := TableStyles()
	if st.Header.GetBold() != true {
		t.Error("expected header style to be bold")
	}
}
