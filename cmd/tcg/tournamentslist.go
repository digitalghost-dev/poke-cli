package tcg

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type tournamentsModel struct {
	conn        func(string) ([]byte, error)
	tournaments []tournamentData
	selected    *tournamentData
	error       error
	list        list.Model
	loading     bool
	spinner     spinner.Model
	quitting    bool
}

type tournamentData struct {
	Location string `json:"location"`
	TextDate string `json:"text_date"`
}

type tournamentsDataMsg struct {
	tournaments []tournamentData
	err         error
}

func fetchTournaments(conn func(string) ([]byte, error)) tea.Cmd {
	return func() tea.Msg {
		endpoint := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=location,text_date&rank=eq.1&order=start_date.desc"
		body, err := conn(endpoint)
		if err != nil {
			return tournamentsDataMsg{err: err}
		}

		var allTournaments []tournamentData
		if err = json.Unmarshal(body, &allTournaments); err != nil {
			return tournamentsDataMsg{err: err}
		}

		return tournamentsDataMsg{tournaments: allTournaments}
	}
}

func tournamentsList(conn func(string) ([]byte, error)) tournamentsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return tournamentsModel{
		conn:    conn,
		loading: true,
		spinner: s,
	}
}

func (m tournamentsModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchTournaments(m.conn),
	)
}

func (m tournamentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			idx := m.list.Index()
			if idx >= 0 && idx < len(m.tournaments) {
				m.selected = &m.tournaments[idx]
			}
			return m, tea.Quit
		}

	case tournamentsDataMsg:
		if msg.err != nil {
			m.error = msg.err
			m.loading = false
			return m, nil
		}

		m.tournaments = msg.tournaments

		var items []list.Item
		for _, t := range msg.tournaments {
			items = append(items, styling.Item(t.Location+" · "+t.TextDate))
		}

		const listWidth = 40
		const listHeight = 16

		l := list.New(items, styling.ItemDelegate{}, listWidth, listHeight)

		l.Title = "First, pick a tournament"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = styling.TitleStyle
		l.Styles.PaginationStyle = styling.PaginationStyle
		l.Styles.HelpStyle = styling.HelpStyle

		m.list = l
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		if !m.loading && m.error == nil {
			m.list.SetWidth(msg.Width)
		}
		return m, nil
	}

	if !m.loading && m.error == nil {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m tournamentsModel) View() string {
	if m.quitting {
		return "\n  Quitting...\n\n"
	}
	if m.error != nil {
		return styling.ApiErrorStyle.Render(
			"Error loading tournaments from Supabase:\n" +
				m.error.Error() + "\n\n" +
				"Press ctrl+c or esc to exit.",
		)
	}
	if m.loading {
		return "\n  " + m.spinner.View() + " Loading tournaments...\n\n"
	}
	if m.selected != nil {
		return styling.QuitTextStyle.Render("Tournament selected:", m.selected.Location+" · "+m.selected.TextDate)
	}

	return "\n" + m.list.View()
}
