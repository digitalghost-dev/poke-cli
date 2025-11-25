package card

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type CardsModel struct {
	Choice         string
	IllustratorMap map[string]string
	ImageMap       map[string]string
	PriceMap       map[string]string
	Quitting       bool
	SelectedOption string
	SeriesName     string
	Table          table.Model
	ViewImage      bool
}

func (m CardsModel) Init() tea.Cmd {
	return nil
}

func (m CardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var bubbleCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case " ":
			m.ViewImage = true
			return m, tea.Quit
		}
	}

	m.Table, bubbleCmd = m.Table.Update(msg)

	// Keep the selected option in sync on every update
	if row := m.Table.SelectedRow(); len(row) > 0 {
		name := row[0]
		if name != m.SelectedOption {
			m.SelectedOption = name
		}
	}

	return m, bubbleCmd
}

func (m CardsModel) View() string {
	if m.Quitting {
		return "\n  Quitting card search...\n\n"
	}

	selectedCard := ""
	if row := m.Table.SelectedRow(); len(row) > 0 {
		cardName := row[0]
		price := m.PriceMap[cardName]
		if price == "" {
			price = "Price: Not available"
		}
		illustrator := m.IllustratorMap[cardName]
		selectedCard = cardName + "\n---\n" + price + "\n---\n" + illustrator
	}

	leftPanel := styling.TypesTableBorder.Render(m.Table.View())

	rightPanel := lipgloss.NewStyle().
		Width(40).
		Height(29).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		Padding(1).
		Render(selectedCard)

	screen := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return fmt.Sprintf("Highlight a card!\n%s\n%s",
		screen,
		styling.KeyMenu.Render("↑ (move up) • ↓ (move down)\nctrl+c | esc (quit)"))
}

type cardData struct {
	Illustrator    string  `json:"illustrator"`
	ImageURL       string  `json:"image_url"`
	MarketPrice    float64 `json:"market_price"`
	Name           string  `json:"name"`
	NumberPlusName string  `json:"number_plus_name"`
}

// CardsList creates and returns a new CardsModel with cards from a specific set
func CardsList(setID string) (CardsModel, error) {
	url := fmt.Sprintf("https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/card_pricing_view?set_id=eq.%s&select=number_plus_name,market_price,image_url,illustrator&order=localId", setID)
	body, err := CallCardData(url)
	if err != nil {
		return CardsModel{}, fmt.Errorf("failed to fetch card data: %w", err)
	}

	var allCards []cardData
	err = json.Unmarshal(body, &allCards)
	if err != nil {
		return CardsModel{}, fmt.Errorf("failed to unmarshal card data: %w", err)
	}

	// Extract card names and build table rows + price map
	rows := make([]table.Row, len(allCards))
	priceMap := make(map[string]string)
	imageMap := make(map[string]string)
	illustratorMap := make(map[string]string)
	for i, card := range allCards {
		rows[i] = []string{card.NumberPlusName}
		priceMap[card.NumberPlusName] = fmt.Sprintf("Price: $%.2f", card.MarketPrice)
		illustratorMap[card.NumberPlusName] = "Illustrator: " + card.Illustrator
		imageMap[card.NumberPlusName] = card.ImageURL
	}

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Card Name", Width: 35}}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(28),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(lipgloss.Color("#FFCC00"))
	t.SetStyles(s)

	return CardsModel{
		IllustratorMap: illustratorMap,
		ImageMap:       imageMap,
		PriceMap:       priceMap,
		Table:          t,
	}, nil
}

func CallCardData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("apikey", "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Authorization", "Bearer sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
