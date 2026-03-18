package tcg

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type tournamentsModel struct {
	choice   string
	error    error
	list     list.Model
	loading  bool
	spinner  spinner.Model
	quitting bool
}

type tournamentData struct {
	Location string `json:"location"`
	TextDate string `json:"text_date"`
}

type tournamentsDataMsg struct {
	items []list.Item
	err   error
}

func fetchTournaments() tea.Cmd {
	return func() tea.Msg {
		url := "https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/standings?select=location,text_date&rank=eq.1&order=start_date"
		body, err := supabaseConn(url)
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

func tournamentsList() tournamentsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return tournamentsModel{
		loading: true,
		spinner: s,
	}
}

func (m tournamentsModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchTournaments(),
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
			i, ok := m.list.SelectedItem().(styling.Item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}

	// Once the data arrives, stop loading and build the list
	case tournamentsDataMsg:
		if msg.err != nil {
			m.error = msg.err
			m.loading = false
			return m, nil
		}

		const listWidth = 20
		const listHeight = 16

		l := list.New(msg.items, styling.ItemDelegate{}, listWidth, listHeight)

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
	if m.choice != "" {
		return styling.QuitTextStyle.Render("Tournament selected:", m.choice)
	}

	return "\n" + m.list.View()
}
