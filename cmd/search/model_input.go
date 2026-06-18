package search

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

// UpdateInput handles text input updates.
func UpdateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.showResults {
			// If results are shown, pressing 'b' resets to search view
			if msg.String() == "b" {
				m.showResults = false
				m.textInput.Reset()
				m.textInput.Focus()
				m.warningMessage = ""
				return m, textinput.Blink
			}
		} else {
			switch msg.Code {
			case tea.KeyEnter:
				searchTerm := m.textInput.Value()
				_, endpoint := RenderInput(m)

				// checking for blank queries
				if strings.TrimSpace(searchTerm) == "" {
					m.warningMessage = utils.FormatError("No blank queries")
					return m, nil
				}

				// Call PokéAPI
				result, err := query(endpoint, searchTerm)
				if err != nil {
					m.warningMessage = utils.FormatError(fmt.Sprintf("Error fetching search results: %v", err))
					return m, nil
				}

				// Format results
				var sb strings.Builder
				for _, r := range result.Results {
					sb.WriteString(styling.ColoredBullet.String() + " " + r.Name + "\n")
				}
				resultsDisplay := sb.String()

				m.searchResults = resultsDisplay
				m.showResults = true
				return m, nil
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// RenderInput renders the input view.
func RenderInput(m model) (string, string) {
	var msg string
	var endpoint string

	switch m.choice {
	case 0:
		msg = "Enter an Ability name:"
		endpoint = "ability"
	case 1:
		msg = "Enter a Move name:"
		endpoint = "move"
	case 2:
		msg = "Enter a Pokémon name:"
		endpoint = "pokemon"
	default:
		msg = "Enter your search query:"
	}

	if m.showResults {
		// Check if there are any results
		results := m.searchResults
		if strings.TrimSpace(results) == "" {
			results = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Render("No results found.")
		}

		return fmt.Sprintf(
			"Search Results:\n\n%s\n\n%s",
			results,
			styling.KeyMenu.Render("Press 'b' to search again\nenter (select) • ctrl+c | esc (quit)"),
		), endpoint
	}

	warning := ""
	if m.warningMessage != "" {
		warning = "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render(m.warningMessage)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s%s",
		msg,
		m.textInput.View(),
		styling.KeyMenu.Render("Press Enter to confirm\nctrl+c | esc (quit)"),
		warning,
	), endpoint
}
