package shell

import (
	"encoding/json"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/comp/web"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type TournamentRef struct {
	Location string `json:"location"`
	TextDate string `json:"text_date"`
}

type pickerModel struct {
	conn        ConnFunc
	listURL     string
	tournaments []TournamentRef
	selected    *TournamentRef
	error       error
	list        list.Model
	loading     bool
	spinner     spinner.Model
	quitting    bool
}

type tournamentsDataMsg struct {
	tournaments []TournamentRef
	err         error
}

func fetchTournaments(listURL string, conn ConnFunc) tea.Cmd {
	return func() tea.Msg {
		body, err := conn(listURL)
		if err != nil {
			return tournamentsDataMsg{err: err}
		}

		var all []TournamentRef
		if err = json.Unmarshal(body, &all); err != nil {
			return tournamentsDataMsg{err: err}
		}

		return tournamentsDataMsg{tournaments: all}
	}
}

func newPicker(spec Spec, conn ConnFunc) pickerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return pickerModel{
		conn:    conn,
		listURL: spec.ListURL,
		loading: true,
		spinner: s,
	}
}

func (m pickerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchTournaments(m.listURL, m.conn),
	)
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "w":
			return m, web.Open("https://web.poke-cli.com/")
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

		webBinding := key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "web"),
		)
		l.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{webBinding} }
		l.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{webBinding} }

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

func (m pickerModel) View() tea.View {
	var content string
	if m.quitting {
		content = "\n  Quitting...\n\n"
	} else if m.error != nil {
		content = styling.ApiErrorStyle.Render(
			"Error loading tournaments from Supabase:\n" +
				m.error.Error() + "\n\n" +
				"Press ctrl+c or esc to exit.",
		)
	} else if m.loading {
		content = "\n  " + m.spinner.View() + " Loading tournaments...\n\n"
	} else if m.selected != nil {
		content = styling.QuitTextStyle.Render("Tournament selected:", m.selected.Location+" · "+m.selected.TextDate)
	} else {
		content = "\n" + m.list.View()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
