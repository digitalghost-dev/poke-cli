package search

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/styling"
	"log"
)

// UpdateInput handles text input updates.
func UpdateInput(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter: // User presses enter after input
			searchTerm := m.TextInput.Value()
			_, endpoint := RenderInput(m) // Get endpoint from RenderInput()

			// Call PokéAPI with the search term
			result, err := query(endpoint, searchTerm)
			if err != nil {
				log.Printf("Error fetching search results: %v", err)
				m.Quitting = true
				return m, tea.Quit
			}

			// Format results
			var resultsDisplay string
			for _, r := range result.Results {
				resultsDisplay += fmt.Sprintf("- %s\n", r.Name)
			}

			// Store results in the model (so they can be rendered)
			m.SearchResults = resultsDisplay
			m.ShowResults = true // Switch to results view
			return m, nil
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
		help := styling.KeyMenu.Render("\nctrl+c | esc (quit)")
		// Show the search results instead of the input field
		return fmt.Sprintf("Search Results:\n\n%s\n%s", m.SearchResults, help), endpoint
	}

	return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			msg,
			m.TextInput.View(),
			styling.KeyMenu.Render("Press Enter to confirm"),
		),
		endpoint
}
