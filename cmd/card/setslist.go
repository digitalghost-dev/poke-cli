package card

import (
	"encoding/json"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var getSetsData = connections.CallTCGData

type setsModel struct {
	Choice     string
	Error      error
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
				items = append(items, styling.Item(set.SetName))
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

func (m setsModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchSetsCmd(m.SeriesName),
	)
}

func (m setsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			if m.Error != nil {
				return m, nil
			}
			i, ok := m.List.SelectedItem().(styling.Item)
			if ok {
				m.Choice = string(i)
				m.SetID = m.SetsIDMap[string(i)]
			}
			return m, tea.Quit
		}

	// Once the data arrives, stop loading and build the list
	case setsDataMsg:
		if msg.err != nil {
			m.Error = msg.err
			m.Loading = false
			return m, nil
		}

		const listWidth = 20
		const listHeight = 20

		l := list.New(msg.items, styling.ItemDelegate{}, listWidth, listHeight)
		l.Title = "Choose a set!"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = styling.TitleStyle
		l.Styles.PaginationStyle = styling.PaginationStyle
		l.Styles.HelpStyle = styling.HelpStyle

		m.List = l
		m.SetsIDMap = msg.setsIDMap
		m.Loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		if !m.Loading && m.Error == nil {
			m.List.SetWidth(msg.Width)
		}
		return m, nil
	}

	// Only update the list if it's been initialized
	if !m.Loading && m.Error == nil {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m setsModel) View() tea.View {
	var content string
	if m.Error != nil {
		content = styling.ApiErrorStyle.Render(
			"Error loading sets from Supabase:\n" +
				m.Error.Error() + "\n\n" +
				"Press ctrl+c or esc to exit.",
		)
	} else if m.Choice != "" {
		content = styling.QuitTextStyle.Render("Set selected:", m.Choice)
	} else if m.Loading {
		content = lipgloss.NewStyle().Padding(2).Render(
			m.Spinner.View() + "Loading sets...",
		)
	} else if m.Quitting {
		content = "\n  Quitting card search...\n\n"
	} else {
		content = "\n" + m.List.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
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
func SetsList(seriesID string) (setsModel, error) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return setsModel{
		SeriesName: seriesID,
		Loading:    true,
		Spinner:    s,
	}, nil
}
