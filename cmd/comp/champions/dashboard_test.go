package champions

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/shell"
)

func noopConn(_ string) ([]byte, error) { return []byte("[]"), nil }

func testTeams() []teamRow {
	return []teamRow{
		{
			Player:     "Alice",
			Record:     "7-1",
			Tournament: "Worlds 2026",
			Archetypes: []string{"Hyper Offense"},
			Pokemon:    []string{"Miraidon", "Flutter Mane", "Iron Hands", "Landorus", "Rillaboom", "Ogerpon"},
			WebURL:     "https://example.com/team/1",
		},
		{
			Player:     "Bob",
			Record:     "6-2",
			Tournament: "Regional Sao Paulo",
			Archetypes: nil,
			Pokemon:    []string{"Calyrex"},
			WebURL:     "",
		},
	}
}

func testCompInfo() []compInfoRow {
	return []compInfoRow{
		{
			Pokemon:         "Miraidon",
			WebURL:          "https://example.com/pokemon/1",
			CommonMoves:     []commonStat{{Name: "Protect", UsagePercent: 90.5}},
			CommonAbilities: []commonStat{{Name: "Hadron Engine", UsagePercent: 100}},
			CommonItems:     []commonStat{{Name: "Choice Specs", UsagePercent: 45.2}},
			CommonTeammates: []commonStat{{Name: "Flutter Mane", UsagePercent: 70.1}},
		},
		{
			Pokemon:         "Calyrex-Shadow",
			WebURL:          "",
			CommonMoves:     []commonStat{{Name: "Astral Barrage", UsagePercent: 99.1}},
			CommonAbilities: []commonStat{{Name: "As One", UsagePercent: 100}},
			CommonItems:     nil,
			CommonTeammates: []commonStat{{Name: "Miraidon", UsagePercent: 40.0}},
		},
	}
}

func testUsage() []usageRow {
	return []usageRow{
		{Rank: 1, Pokemon: "Basculegion", UsagePercent: 51.5},
		{Rank: 2, Pokemon: "Kingambit", UsagePercent: 40.69},
	}
}

func testSpeedTiers() []speedTierRow {
	return []speedTierRow{
		{Rank: 1, Pokemon: "Mega Aerodactyl", BaseSpe: 150, Neutral0: 170, Neutral252: 202, NegMin: 153, Max: 222, MaxScarf: 333, NeutralScarf: 303},
		{Rank: 9, Pokemon: "Aerodactyl", BaseSpe: 130, Neutral0: 150, Neutral252: 182, NegMin: 135, Max: 200, MaxScarf: 300, NeutralScarf: 273},
	}
}

func newTestDashboard() dashboardModel {
	return dashboardModel{
		conn:   noopConn,
		styles: shell.NewStyles(),
		width:  120,
		height: 40,
	}
}

func loadedTestDashboard() dashboardModel {
	m := newTestDashboard()
	nm, _ := m.Update(dataMsg{data: &dashboardData{CompInfo: testCompInfo(), Teams: testTeams(), Usage: testUsage(), SpeedTiers: testSpeedTiers()}})
	return nm.(dashboardModel)
}

func TestNewDashboard(t *testing.T) {
	m := newDashboard(noopConn)
	if m.styles == nil {
		t.Error("expected styles to be set")
	}
	if m.conn == nil {
		t.Error("expected conn to be set")
	}
	if m.data != nil || m.activeTab != 0 || m.goBack {
		t.Error("expected a clean initial model")
	}
}

func TestDashboard_Init_ReturnsCmd(t *testing.T) {
	if newTestDashboard().Init() == nil {
		t.Error("expected Init() to return a non-nil cmd")
	}
}

func TestDashboard_Update_Quit(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.KeyPressMsg
	}{
		{"ctrl+c", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}},
		{"esc", tea.KeyPressMsg{Code: tea.KeyEscape}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newM, cmd := newTestDashboard().Update(tt.msg)
			if newM.(dashboardModel).goBack {
				t.Error("quit should not set goBack")
			}
			if cmd == nil {
				t.Error("expected a quit command")
			}
		})
	}
}

func TestDashboard_Update_Back(t *testing.T) {
	tm := teatest.NewTestModel(t, newTestDashboard(), teatest.WithInitialTermSize(120, 40))
	tm.Send(tea.KeyPressMsg{Code: 'b', Text: "b"})
	tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
	if !tm.FinalModel(t).(dashboardModel).goBack {
		t.Error("expected goBack=true after pressing b")
	}
}

func TestDashboard_Update_OpenWeb(t *testing.T) {
	newM, cmd := newTestDashboard().Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	if cmd == nil {
		t.Error("expected a command from pressing w")
	}
	if newM.(dashboardModel).goBack {
		t.Error("pressing w should not set goBack")
	}
}

func TestDashboard_Update_TabNavigation(t *testing.T) {
	m := newTestDashboard()
	for _, key := range []tea.KeyPressMsg{{Code: tea.KeyRight}, {Code: 'l', Text: "l"}, {Code: tea.KeyTab}} {
		nm, _ := m.Update(key)
		if nm.(dashboardModel).activeTab != 1 {
			t.Errorf("expected activeTab=1 after %v, got %d", key, nm.(dashboardModel).activeTab)
		}
	}

	m.activeTab = 1
	for _, key := range []tea.KeyPressMsg{{Code: tea.KeyLeft}, {Code: 'h', Text: "h"}, {Code: tea.KeyTab, Mod: tea.ModShift}} {
		nm, _ := m.Update(key)
		if nm.(dashboardModel).activeTab != 0 {
			t.Errorf("expected activeTab=0 after %v, got %d", key, nm.(dashboardModel).activeTab)
		}
	}
}

func TestDashboard_Update_TabNavigation_Clamps(t *testing.T) {
	m := newTestDashboard()
	nm, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if nm.(dashboardModel).activeTab != 0 {
		t.Errorf("expected activeTab to clamp at 0, got %d", nm.(dashboardModel).activeTab)
	}

	m.activeTab = len(tabs) - 1
	nm2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if nm2.(dashboardModel).activeTab != len(tabs)-1 {
		t.Errorf("expected activeTab to clamp at %d, got %d", len(tabs)-1, nm2.(dashboardModel).activeTab)
	}
}

func TestDashboard_Update_DataMsg_Success(t *testing.T) {
	m := loadedTestDashboard()
	if m.data == nil {
		t.Fatal("expected data to be set")
	}
	if len(m.teams.Rows()) != 2 {
		t.Errorf("expected 2 team rows, got %d", len(m.teams.Rows()))
	}
	if len(m.overview.Rows()) != 2 {
		t.Errorf("expected 2 overview rows, got %d", len(m.overview.Rows()))
	}
	if len(m.speed.Rows()) != 2 {
		t.Errorf("expected 2 speed rows, got %d", len(m.speed.Rows()))
	}
	if len(m.usage.Rows()) != 2 {
		t.Errorf("expected 2 usage rows, got %d", len(m.usage.Rows()))
	}
}

func TestDashboard_Update_OverviewTableNavigation(t *testing.T) {
	m := loadedTestDashboard()
	m.activeTab = 0
	nm, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if nm.(dashboardModel).overview.Cursor() != 1 {
		t.Errorf("expected overview cursor to advance to 1, got %d", nm.(dashboardModel).overview.Cursor())
	}
}

func TestDashboard_Update_DataMsg_Error(t *testing.T) {
	m := newTestDashboard()
	nm, _ := m.Update(dataMsg{err: errors.New("fetch failed")})
	if nm.(dashboardModel).err == nil {
		t.Error("expected err to be set")
	}
}

func TestDashboard_Update_WindowSize(t *testing.T) {
	loaded := loadedTestDashboard()
	nm, _ := loaded.Update(tea.WindowSizeMsg{Width: 160, Height: 50})
	result := nm.(dashboardModel)
	if result.width != 160 || result.height != 50 {
		t.Errorf("expected 160x50, got %dx%d", result.width, result.height)
	}
	if len(result.teams.Rows()) != 2 {
		t.Errorf("expected teams table rebuilt with 2 rows, got %d", len(result.teams.Rows()))
	}
}

func TestDashboard_Update_WindowSize_BeforeData(t *testing.T) {
	nm, _ := newTestDashboard().Update(tea.WindowSizeMsg{Width: 160, Height: 50})
	result := nm.(dashboardModel)
	if result.width != 160 || result.height != 50 {
		t.Errorf("expected 160x50, got %dx%d", result.width, result.height)
	}
}

func TestDashboard_Update_TeamsTableNavigation(t *testing.T) {
	m := loadedTestDashboard()
	m.activeTab = 2
	nm, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if nm.(dashboardModel).teams.Cursor() != 1 {
		t.Errorf("expected table cursor to advance to 1, got %d", nm.(dashboardModel).teams.Cursor())
	}
}

func TestDashboard_Update_UsageTableNavigation(t *testing.T) {
	m := loadedTestDashboard()
	m.activeTab = 1
	nm, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if nm.(dashboardModel).usage.Cursor() != 1 {
		t.Errorf("expected usage cursor to advance to 1, got %d", nm.(dashboardModel).usage.Cursor())
	}
}

func TestDashboard_Update_SpeedTableNavigation(t *testing.T) {
	m := loadedTestDashboard()
	m.activeTab = 3
	nm, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if nm.(dashboardModel).speed.Cursor() != 1 {
		t.Errorf("expected speed cursor to advance to 1, got %d", nm.(dashboardModel).speed.Cursor())
	}
}

func TestDashboard_Update_UnhandledKey(t *testing.T) {
	m := newTestDashboard()
	nm, cmd := m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	if cmd != nil {
		t.Error("expected no command for an unhandled key without data")
	}
	if nm.(dashboardModel).activeTab != 0 {
		t.Error("unhandled key should not change the active tab")
	}
}

func TestDashboard_View_NilStyles(t *testing.T) {
	if (dashboardModel{}).View().Content != "" {
		t.Error("expected empty view when styles is nil")
	}
}

func TestDashboard_View_Loading(t *testing.T) {
	v := newTestDashboard().View()
	if !v.AltScreen {
		t.Error("expected AltScreen enabled")
	}
	if !strings.Contains(v.Content, "Loading") {
		t.Error("expected loading message before data arrives")
	}
}

func TestDashboard_View_FetchError(t *testing.T) {
	m := newTestDashboard()
	m.err = errors.New("network error")
	if !strings.Contains(m.View().Content, "fetch error") {
		t.Error("expected fetch error in view")
	}
}

func TestDashboard_View_ContainsTabs(t *testing.T) {
	content := newTestDashboard().View().Content
	for _, tab := range tabs {
		if !strings.Contains(content, tab) {
			t.Errorf("expected view to contain tab %q", tab)
		}
	}
}

func TestDashboard_View_AllTabs(t *testing.T) {
	m := loadedTestDashboard()
	for tab := range len(tabs) {
		m.activeTab = tab
		if m.View().Content == "" {
			t.Errorf("expected non-empty view for tab %d", tab)
		}
	}
}

func TestDashboard_RenderTab(t *testing.T) {
	loaded := loadedTestDashboard()

	tests := []struct {
		name     string
		model    dashboardModel
		tab      int
		contains string
	}{
		{"error", func() dashboardModel { m := newTestDashboard(); m.err = errors.New("boom"); return m }(), 0, "fetch error"},
		{"loading", newTestDashboard(), 0, "Loading"},
		{"overview", loaded, 0, "Miraidon"},
		{"usage", loaded, 1, "Basculegion"},
		{"top teams", loaded, 2, "Alice"},
		{"speed tiers", loaded, 3, "Mega Aerodactyl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.activeTab = tt.tab
			out := tt.model.renderTab(contentWidth(tt.model.width))
			if !strings.Contains(out, tt.contains) {
				t.Errorf("renderTab tab %d = %q, want substring %q", tt.tab, out, tt.contains)
			}
		})
	}
}

func TestContentWidth(t *testing.T) {
	tests := []struct {
		width int
		want  int
	}{
		{120, 110},
		{50, 40},
		{10, 40},
		{0, 40},
	}
	for _, tt := range tests {
		if got := contentWidth(tt.width); got != tt.want {
			t.Errorf("contentWidth(%d) = %d, want %d", tt.width, got, tt.want)
		}
	}
}
