package card

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var getSetsData = connections.CallTCGData

type SetsModel struct {
	Choice     string
	Loading    bool
	List       list.Model
	Quitting   bool
	SeriesName string
	SetID      string
	SetsIDMap  map[string]string // Maps set name -> set_id
	Spinner    spinner.Model
}

// Message type to carry fetched data back to Update()
type setsDataMsg struct {
	items     []list.Item
	setsIDMap map[string]string
	seriesID  string
	err       error
}

func fetchSetsCmd(seriesID string) tea.Cmd {
	return func() tea.Msg {
		body, err := getSetsData("https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/sets")
		if err != nil {
			return setsDataMsg{err: err}
		}

		var allSets []setData
		err = json.Unmarshal(body, &allSets)
		if err != nil {
			return setsDataMsg{err: err}
		}

		// Filter sets by series_id and build ID map
		var items []list.Item
		setsIDMap := make(map[string]string)
		for _, set := range allSets {
			if set.SeriesID == seriesID {
				items = append(items, item(set.SetName))
				setsIDMap[set.SetName] = set.SetID
			}
		}

		return setsDataMsg{
			items:     items,
			setsIDMap: setsIDMap,
			seriesID:  seriesID,
		}
	}
}

func (m SetsModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchSetsCmd(m.SeriesName),
	)
}

func (m SetsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(item)
			if ok {
				m.Choice = string(i)
				m.SetID = m.SetsIDMap[string(i)]
			}
			return m, tea.Quit
		}

	case setsDataMsg:
		// Data arrived - stop loading and build the list
		if msg.err != nil {
			m.Quitting = true
			return m, tea.Quit
		}

		const listWidth = 20
		const listHeight = 20

		l := list.New(msg.items, itemDelegate{}, listWidth, listHeight)
		l.Title = fmt.Sprintf("Pick a set from the %s series", msg.seriesID)
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		m.List = l
		m.SetsIDMap = msg.setsIDMap
		m.Loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		if !m.Loading {
			m.List.SetWidth(msg.Width)
		}
		return m, nil
	}

	// Only update the list if it's been initialized
	if !m.Loading {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SetsModel) View() string {
	if m.Choice != "" {
		return quitTextStyle.Render("Set selected:", m.Choice)
	}
	if m.Loading {
		return lipgloss.NewStyle().Padding(2).Render(
			m.Spinner.View() + "Loading sets...",
		)
	}
	if m.Quitting {
		return "\n  Quitting card search...\n\n"
	}

	return "\n" + m.List.View()
}

type setData struct {
	SeriesID          string `json:"series_id"`
	SetID             string `json:"set_id"`
	SetName           string `json:"set_name"`
	OfficialCardCount int    `json:"official_card_count"`
	TotalCardCount    int    `json:"total_card_count"`
	Logo              string `json:"logo"`
	Symbol            string `json:"symbol"`
}

// SetsList returns a minimal model - data fetching happens via Init()
func SetsList(seriesID string) (SetsModel, error) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return SetsModel{
		SeriesName: seriesID,
		Loading:    true,
		Spinner:    s,
	}, nil
}
