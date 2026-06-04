package tcg

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

func newTestModel() model {
	return model{
		conn:       noopConn,
		tabs:       []string{"Overview", "Standings", "Decks", "Countries"},
		styles:     newStyles(),
		tournament: "London",
		width:      120,
		height:     40,
	}
}

func loadedTestModel() model {
	m := newTestModel()
	items := []standingRows{
		{Rank: 1, Name: "Ash", Points: 47, Deck: "gardevoir", PlayerCountry: "USA", ISOCode: "US", PlayerQty: 500, TextDate: "Jan 10", Type: "Regional"},
		{Rank: 2, Name: "Misty", Points: 44, Deck: "dragapult", PlayerCountry: "Japan", ISOCode: "JP", PlayerQty: 500, TextDate: "Jan 10", Type: "Regional"},
	}
	newModel, _ := m.Update(standingsDataMsg{items: items})
	return newModel.(model)
}

// Init

func TestDashboardModel_Init_ReturnsCmd(t *testing.T) {
	m := newTestModel()
	if m.Init() == nil {
		t.Error("expected Init() to return a non-nil cmd")
	}
}

// Update — key messages

func TestDashboardModel_Update_Quit(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
	}{
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}},
		{"esc", tea.KeyPressMsg{Code: tea.KeyEscape}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 40))
			tm.Send(tt.msg)
			tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
		})
	}
}

func TestDashboardModel_Update_Back(t *testing.T) {
	m := newTestModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 40))
	tm.Send(tea.KeyPressMsg{Code: 'b', Text: "b"})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	final := tm.FinalModel(t).(model)
	if !final.goBack {
		t.Error("expected goBack=true after pressing b")
	}
}

func TestDashboardModel_Update_TabNavigation(t *testing.T) {
	m := newTestModel()
	// right moves forward
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if newM.(model).activeTab != 1 {
		t.Errorf("expected activeTab=1 after right, got %d", newM.(model).activeTab)
	}
	// left moves back
	newM2, _ := newM.(model).Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if newM2.(model).activeTab != 0 {
		t.Errorf("expected activeTab=0 after left, got %d", newM2.(model).activeTab)
	}
}

func TestDashboardModel_Update_TabNavigation_Clamps(t *testing.T) {
	m := newTestModel()
	// can't go below 0
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if newM.(model).activeTab != 0 {
		t.Errorf("expected activeTab to clamp at 0, got %d", newM.(model).activeTab)
	}
	// can't exceed last tab
	m.activeTab = 3
	newM2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if newM2.(model).activeTab != 3 {
		t.Errorf("expected activeTab to clamp at 3, got %d", newM2.(model).activeTab)
	}
}

// Update — standingsDataMsg

func TestDashboardModel_Update_StandingsDataMsg_Success(t *testing.T) {
	m := loadedTestModel()
	if len(m.standings) != 2 {
		t.Errorf("expected 2 standings rows, got %d", len(m.standings))
	}
	if m.winner != "Ash" {
		t.Errorf("expected winner=Ash, got %q", m.winner)
	}
	if m.winningDeck != "gardevoir" {
		t.Errorf("expected winningDeck=gardevoir, got %q", m.winningDeck)
	}
	if m.totalPlayers != 500 {
		t.Errorf("expected totalPlayers=500, got %d", m.totalPlayers)
	}
	if len(m.countryStats) == 0 {
		t.Error("expected countryStats to be populated")
	}
	if len(m.deckStats) == 0 {
		t.Error("expected deckStats to be populated")
	}
}

func TestDashboardModel_Update_StandingsDataMsg_Error(t *testing.T) {
	m := newTestModel()
	newM, _ := m.Update(standingsDataMsg{err: errors.New("fetch failed")})
	result := newM.(model)
	if result.err == nil {
		t.Error("expected err to be set")
	}
}

func TestDashboardModel_Update_EmptyPlayerCountry_Skipped(t *testing.T) {
	m := newTestModel()
	items := []standingRows{
		{Rank: 1, Name: "Ash", PlayerCountry: ""},
		{Rank: 2, Name: "Misty", PlayerCountry: "Japan"},
	}
	newM, _ := m.Update(standingsDataMsg{items: items})
	result := newM.(model)
	if len(result.countryStats) != 1 {
		t.Errorf("expected 1 countryStats entry (empty country skipped), got %d", len(result.countryStats))
	}
}

// Update — WindowSizeMsg

func TestDashboardModel_Update_WindowSize(t *testing.T) {
	m := newTestModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 160, Height: 50})
	result := newM.(model)
	if result.width != 160 {
		t.Errorf("expected width=160, got %d", result.width)
	}
	if result.height != 50 {
		t.Errorf("expected height=50, got %d", result.height)
	}
}

// View

func TestDashboardModel_View_NilStyles(t *testing.T) {
	m := model{}
	if m.View().Content != "" {
		t.Error("expected empty string when styles is nil")
	}
}

func TestDashboardModel_View_ContainsTabs(t *testing.T) {
	m := newTestModel()
	view := m.View()
	for _, tab := range []string{"Overview", "Standings", "Decks", "Countries"} {
		if !strings.Contains(view.Content, tab) {
			t.Errorf("expected view to contain tab %q", tab)
		}
	}
}

func TestDashboardModel_View_LoadingState(t *testing.T) {
	m := newTestModel()
	view := m.View()
	if !strings.Contains(view.Content, "Loading") {
		t.Error("expected loading message before data arrives")
	}
}

func TestDashboardModel_View_FetchError(t *testing.T) {
	m := newTestModel()
	m.err = errors.New("network error")
	view := m.View()
	if !strings.Contains(view.Content, "fetch error") {
		t.Errorf("expected fetch error in view, got: %s", view.Content)
	}
}

func TestDashboardModel_View_AllTabs(t *testing.T) {
	m := loadedTestModel()
	for tab := 0; tab <= 3; tab++ {
		m.activeTab = tab
		view := m.View()
		if view.Content == "" {
			t.Errorf("expected non-empty view for tab %d", tab)
		}
	}
}
