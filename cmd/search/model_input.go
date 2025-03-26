package search

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
	"strings"
)

// UpdateInput handles text input updates.
func UpdateInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.ShowResults {
			// If results are shown, pressing 'b' resets to search view
			if msg.String() == "b" {
				m.ShowResults = false
				m.TextInput.Reset()
				m.TextInput.Focus()
				m.WarningMessage = ""
				return m, textinput.Blink
			}
		} else {
			switch msg.Type {
			case tea.KeyEnter:
				searchTerm := m.TextInput.Value()
				_, endpoint := RenderInput(m)

				// checking for blank queries
				if strings.TrimSpace(searchTerm) == "" {
					errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nNo blank queries")
					m.WarningMessage = errMessage
					return m, nil
				}

				// Call PokéAPI
				result, err := query(endpoint, searchTerm)
				if err != nil {
					fmt.Printf("Error fetching search results: %v", err)
					return m, nil
				}

				// Format results
				var resultsDisplay string
				for _, r := range result.Results {
					resultsDisplay += fmt.Sprintf("%s %s\n", styling.ColoredBullet, r.Name)
				}

				m.SearchResults = resultsDisplay
				m.ShowResults = true
				return m, nil
			}
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

// RenderInput renders the input view.
func RenderInput(m Model) (string, string) {
	var msg string
	var endpoint string

	switch m.Choice {
	case 0:
		msg = "Enter a Pokémon name:"
		endpoint = "pokemon"
	case 1:
		msg = "Enter an Ability name:"
		endpoint = "ability"
	default:
		msg = "Enter your search query:"
	}

	if m.ShowResults {
		// Check if there are any results
		results := m.SearchResults
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
	if m.WarningMessage != "" {
		warning = "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render(m.WarningMessage)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s%s",
		msg,
		m.TextInput.View(),
		styling.KeyMenu.Render("Press Enter to confirm\nctrl+c | esc (quit)"),
		warning,
	), endpoint
}
