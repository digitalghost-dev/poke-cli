package card

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
)

var getCardData = connections.CallTCGData

type cardsModel struct {
	AllRows           []table.Row
	Choice            string
	Error             error
	IllustratorMap    map[string]string
	ImageMap          map[string]string
	Loading           bool
	PriceMap          map[string]string
	Quitting          bool
	RegulationMarkMap map[string]string
	Search            textinput.Model
	SelectedOption    string
	SeriesName        string
	SetID             string
	Spinner           spinner.Model
	Table             table.Model
	TableStyles       table.Styles
	ViewImage         bool
}

// Message type to carry fetched card data back to Update()
type cardDataMsg struct {
	allRows           []table.Row
	priceMap          map[string]string
	imageMap          map[string]string
	illustratorMap    map[string]string
	regulationMarkMap map[string]string
	err               error
}

var (
	activeTableSelectedBg   color.Color = styling.YellowColor
	inactiveTableSelectedBg color.Color = lipgloss.Color("#808080")
)

func cardTableStyles(selectedBg color.Color) table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.YellowColor).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(selectedBg)
	return s
}

func syncTableStylesForFocus(m *cardsModel) {
	if m.Search.Focused() {
		m.TableStyles = cardTableStyles(inactiveTableSelectedBg)
	} else {
		m.TableStyles = cardTableStyles(activeTableSelectedBg)
	}
	m.Table.SetStyles(m.TableStyles)
}

// fetchCardsCmd does the actual data fetching and returns a cardDataMsg
func fetchCardsCmd(setID string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/card_pricing_view?set_id=eq.%s&select=number_plus_name,market_price,image_url,illustrator,regulation_mark&order=localId", setID)
		body, err := getCardData(url)
		if err != nil {
			return cardDataMsg{err: err}
		}

		var allCards []cardData
		err = json.Unmarshal(body, &allCards)
		if err != nil {
			return cardDataMsg{err: err}
		}

		rows := make([]table.Row, len(allCards))
		priceMap := make(map[string]string)
		imageMap := make(map[string]string)
		illustratorMap := make(map[string]string)
		regulationMarkMap := make(map[string]string)

		for i, card := range allCards {
			rows[i] = []string{card.NumberPlusName}
			if card.MarketPrice != 0 {
				priceMap[card.NumberPlusName] = fmt.Sprintf("Price: $%.2f", card.MarketPrice)
			} else {
				priceMap[card.NumberPlusName] = "Pricing not available"
			}

			if card.Illustrator != "" {
				illustratorMap[card.NumberPlusName] = "Illustrator: " + card.Illustrator
			} else {
				illustratorMap[card.NumberPlusName] = "Illustrator not available"
			}

			if card.RegulationMark != "" {
				regulationMarkMap[card.NumberPlusName] = "Regulation: " + card.RegulationMark
			} else {
				regulationMarkMap[card.NumberPlusName] = "Regulation not available"
			}

			imageMap[card.NumberPlusName] = card.ImageURL
		}

		return cardDataMsg{
			allRows:           rows,
			priceMap:          priceMap,
			imageMap:          imageMap,
			illustratorMap:    illustratorMap,
			regulationMarkMap: regulationMarkMap,
		}
	}
}

func (m cardsModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchCardsCmd(m.SetID),
	)
}

func (m cardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var bubbleCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "esc":
			// If in the search bar, exit search mode instead of quitting.
			if m.Search.Focused() {
				m.Search.Blur()
				m.Table.Focus()
				syncTableStylesForFocus(&m)
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit
		case "?":
			if !m.Search.Focused() {
				// Sync the selected option before quitting to ensure the correct card is shown
				if row := m.Table.SelectedRow(); len(row) > 0 {
					m.SelectedOption = row[0]
				}
				m.ViewImage = true
				return m, tea.Quit
			}
		case "tab":
			if m.Search.Focused() {
				m.Search.Blur()
				m.Table.Focus()
			} else {
				m.Table.Blur()
				m.Search.Focus()
			}
			syncTableStylesForFocus(&m)
			return m, nil
		}

	case cardDataMsg:
		// Data arrived - stop loading and build the table
		if msg.err != nil {
			m.Error = msg.err
			m.Loading = false
			return m, nil
		}

		ti := textinput.New()
		ti.Placeholder = "type name..."
		ti.Prompt = "🔎 "
		ti.CharLimit = 24
		ti.SetWidth(30)
		ti.Blur()

		t := table.New(
			table.WithColumns([]table.Column{{Title: "Card Name", Width: 35}}),
			table.WithRows(msg.allRows),
			table.WithFocused(true),
			table.WithHeight(25),
			table.WithWidth(35),
		)

		styles := cardTableStyles(activeTableSelectedBg)
		t.SetStyles(styles)

		m.AllRows = msg.allRows
		m.PriceMap = msg.priceMap
		m.ImageMap = msg.imageMap
		m.IllustratorMap = msg.illustratorMap
		m.RegulationMarkMap = msg.regulationMarkMap
		m.Search = ti
		m.Table = t
		m.TableStyles = styles
		m.Loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	// Only handle search/table updates when not loading
	if !m.Loading {
		if m.Search.Focused() {
			prev := m.Search.Value()
			m.Search, bubbleCmd = m.Search.Update(msg)
			if m.Search.Value() != prev {
				applyFilter(&m)
			}
		} else {
			m.Table, bubbleCmd = m.Table.Update(msg)
		}

		// Keep the selected option in sync on every update
		if row := m.Table.SelectedRow(); len(row) > 0 {
			name := row[0]
			if name != m.SelectedOption {
				m.SelectedOption = name
			}
		}
	}

	return m, bubbleCmd
}

func applyFilter(m *cardsModel) {
	q := strings.TrimSpace(strings.ToLower(m.Search.Value()))
	if q == "" {
		m.Table.SetRows(m.AllRows)
		m.Table.SetCursor(0)
		return
	}

	filtered := make([]table.Row, 0, len(m.AllRows))
	for _, r := range m.AllRows {
		if len(r) == 0 {
			continue
		}
		if strings.Contains(strings.ToLower(r[0]), q) {
			filtered = append(filtered, r)
		}
	}

	m.Table.SetRows(filtered)
	m.Table.SetCursor(0)
}

func (m cardsModel) View() tea.View {
	var content string
	if m.Quitting {
		content = "\n  Quitting card search...\n\n"
	} else if m.Error != nil {
		content = styling.ApiErrorStyle.Render(
			"Error loading cards from Supabase:\n" +
				m.Error.Error() + "\n\n" +
				"Press ctrl+c or esc to exit.",
		)
	} else if m.Loading {
		content = lipgloss.NewStyle().Padding(2).Render(
			m.Spinner.View() + " Loading cards...",
		)
	} else {
		selectedCard := ""
		if row := m.Table.SelectedRow(); len(row) > 0 {
			cardName := row[0]
			price := m.PriceMap[cardName]
			if price == "" {
				price = "Price: Not available"
			}
			illustrator := m.IllustratorMap[cardName]
			regulationMark := m.RegulationMarkMap[cardName]
			selectedCard = cardName + "\n---\n" + price + "\n---\n" + illustrator + "\n---\n" + regulationMark
		}

		leftContent := lipgloss.JoinVertical(lipgloss.Left, m.Search.View(), m.Table.View())
		leftPanel := styling.TypesTableBorder.Render(leftContent)

		rightPanel := lipgloss.NewStyle().
			Width(42).
			Height(29).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styling.YellowColor).
			Padding(1).
			Render(selectedCard)

		screen := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

		content = fmt.Sprintf(
			"Highlight a card!\n\nNote: Prices are for normal variations of cards.\n%s\n%s",
			screen,
			styling.KeyMenu.Render("↑ (move up)\n↓ (move down)\n? (view image)\ntab (toggle search)\nctrl+c | esc (quit)"))
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

type cardData struct {
	Illustrator    string  `json:"illustrator"`
	ImageURL       string  `json:"image_url"`
	MarketPrice    float64 `json:"market_price"`
	Name           string  `json:"name"`
	NumberPlusName string  `json:"number_plus_name"`
	RegulationMark string  `json:"regulation_mark"`
}

// CardsList returns a minimal model - data fetching happens via Init()
func CardsList(setID string) (cardsModel, error) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return cardsModel{
		SetID:   setID,
		Loading: true,
		Spinner: s,
	}, nil
}
