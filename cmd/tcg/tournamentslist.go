package tcg

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var getTournamentsData = connections.CallTCGData

type TournamentsModel struct {
	Choice     string
	Error      error
	List       list.Model
	Loading    bool
	Spinner    spinner.Model
	Tournament string
	Quitting   bool
}

type tournamentData struct {
	Location string `json:"location"`
	TextDate string `json:"text_date"`
}

type tournamentsDataMsg struct {
	items []list.Item
	err   error
}

func fetchTournaments(tournament string) tea.Cmd {
	return func() tea.Msg {
		url := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=location,text_date&rank=eq.1&order=start_date"
		body, err := getTournamentsData(url)
		if err != nil {
			return tournamentsDataMsg{err: err}
		}

		var allTournaments []tournamentData
		if err = json.Unmarshal(body, &allTournaments); err != nil {
			return tournamentsDataMsg{err: err}
		}

		var items []list.Item
		for _, t := range allTournaments {
			items = append(items, styling.Item(t.Location+" · "+t.TextDate))
		}

		return tournamentsDataMsg{items: items}
	}
}

func TournamentsList() TournamentsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return TournamentsModel{
		Loading: true,
		Spinner: s,
	}
}

func (m TournamentsModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchTournaments(m.Tournament),
	)
}

func (m TournamentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(styling.Item)
			if ok {
				m.Choice = string(i)
			}
			return m, tea.Quit
		}

	// Once the data arrives, stop loading and build the list
	case tournamentsDataMsg:
		if msg.err != nil {
			m.Error = msg.err
			m.Loading = false
			return m, nil
		}

		const listWidth = 20
		const listHeight = 16

		l := list.New(msg.items, styling.ItemDelegate{}, listWidth, listHeight)

		l.Title = "First, pick a series"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = styling.TitleStyle
		l.Styles.PaginationStyle = styling.PaginationStyle
		l.Styles.HelpStyle = styling.HelpStyle

		m.List = l
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

	if !m.Loading && m.Error == nil {
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m TournamentsModel) View() string {
	if m.Quitting {
		return "\n  Quitting...\n\n"
	}
	if m.Error != nil {
		return styling.ApiErrorStyle.Render(
			"Error loading tournaments from Supabase:\n" +
				m.Error.Error() + "\n\n" +
				"Press ctrl+c or esc to exit.",
		)
	}
	if m.Loading {
		return "\n  " + m.Spinner.View() + " Loading tournaments...\n\n"
	}
	if m.Choice != "" {
		return styling.QuitTextStyle.Render("Tournament selected:", m.Choice)
	}

	return "\n" + m.List.View()
}
