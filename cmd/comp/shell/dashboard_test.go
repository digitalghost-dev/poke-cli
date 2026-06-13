package shell

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

func newTestDashboard() dashboardModel {
	return dashboardModel{
		spec:       testSpec(),
		conn:       noopConn,
		styles:     NewStyles(),
		tournament: "London",
		width:      120,
		height:     40,
	}
}

func loadedTestDashboard() dashboardModel {
	m := newTestDashboard()
	nm, _ := m.Update(dataMsg{decoded: testDecoded()})
	return nm.(dashboardModel)
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
			tm := teatest.NewTestModel(t, newTestDashboard(), teatest.WithInitialTermSize(120, 40))
			tm.Send(tt.msg)
			tm.WaitFinished(t, teatest.WithFinalTimeout(300*time.Millisecond))
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

func TestDashboard_Update_TabNavigation(t *testing.T) {
	m := newTestDashboard()
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if newM.(dashboardModel).activeTab != 1 {
		t.Errorf("expected activeTab=1 after right, got %d", newM.(dashboardModel).activeTab)
	}
	newM2, _ := newM.(dashboardModel).Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if newM2.(dashboardModel).activeTab != 0 {
		t.Errorf("expected activeTab=0 after left, got %d", newM2.(dashboardModel).activeTab)
	}
}

func TestDashboard_Update_TabNavigation_Clamps(t *testing.T) {
	m := newTestDashboard()
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if newM.(dashboardModel).activeTab != 0 {
		t.Errorf("expected activeTab to clamp at 0, got %d", newM.(dashboardModel).activeTab)
	}
	m.activeTab = 3
	newM2, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if newM2.(dashboardModel).activeTab != 3 {
		t.Errorf("expected activeTab to clamp at 3, got %d", newM2.(dashboardModel).activeTab)
	}
}

func TestDashboard_Update_DataMsg_Success(t *testing.T) {
	m := loadedTestDashboard()
	if m.decoded == nil {
		t.Fatal("expected decoded to be set")
	}
	if len(m.table.Rows()) != 2 {
		t.Errorf("expected 2 table rows, got %d", len(m.table.Rows()))
	}
}

func TestDashboard_Update_DataMsg_Error(t *testing.T) {
	m := newTestDashboard()
	newM, _ := m.Update(dataMsg{err: errors.New("fetch failed")})
	if newM.(dashboardModel).err == nil {
		t.Error("expected err to be set")
	}
}

func TestDashboard_Update_WindowSize(t *testing.T) {
	m := newTestDashboard()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 160, Height: 50})
	result := newM.(dashboardModel)
	if result.width != 160 || result.height != 50 {
		t.Errorf("expected 160x50, got %dx%d", result.width, result.height)
	}
}

func TestDashboard_View_NilStyles(t *testing.T) {
	if (dashboardModel{}).View().Content != "" {
		t.Error("expected empty view when styles is nil")
	}
}

func TestDashboard_View_ContainsTabs(t *testing.T) {
	view := newTestDashboard().View()
	for _, tab := range []string{"Overview", "Standings", "Extra", "Countries"} {
		if !strings.Contains(view.Content, tab) {
			t.Errorf("expected view to contain tab %q", tab)
		}
	}
}

func TestDashboard_View_LoadingState(t *testing.T) {
	if !strings.Contains(newTestDashboard().View().Content, "Loading") {
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

func TestDashboard_View_AllTabs(t *testing.T) {
	m := loadedTestDashboard()
	wants := map[int]string{0: "OVERVIEW-BODY", 2: "EXTRA-CAPTION", 3: "USA"}
	for tab := 0; tab <= 3; tab++ {
		m.activeTab = tab
		content := m.View().Content
		if content == "" {
			t.Errorf("expected non-empty view for tab %d", tab)
		}
		if want, ok := wants[tab]; ok && !strings.Contains(content, want) {
			t.Errorf("expected tab %d to render %q", tab, want)
		}
	}
}
