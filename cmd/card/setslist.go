package card

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SetsModel struct {
	List       list.Model
	Choice     string
	SetID      string
	Quitting   bool
	SeriesName string
	setsIDMap  map[string]string // Maps set name -> set_id
}

func (m SetsModel) Init() tea.Cmd {
	return nil
}

func (m SetsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.List.SelectedItem().(item)
			if ok {
				m.Choice = string(i)
				m.SetID = m.setsIDMap[string(i)] // Look up the set_id
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m SetsModel) View() string {
	if m.Quitting {
		return "\n  Quitting card search...\n\n"
	}
	if m.Choice != "" {
		return quitTextStyle.Render("Set selected:", m.Choice)
	}

	return "\n" + m.List.View()
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

// creating a function variable to swap the implementation in tests
var getSetsData = callSetsData

func SetsList(seriesID string) (SetsModel, error) {
	body, err := getSetsData("https://uoddayfnfkebrijlpfbh.supabase.co/rest/v1/sets")
	if err != nil {
		return SetsModel{}, fmt.Errorf("error getting sets data: %v", err)
	}
	var allSets []setData

	err = json.Unmarshal(body, &allSets)
	if err != nil {
		return SetsModel{}, fmt.Errorf("error parsing sets data: %v", err)
	}

	// Filter sets by series_id and build ID map
	var items []list.Item
	setsIDMap := make(map[string]string)
	for _, set := range allSets {
		if set.SeriesID == seriesID {
			items = append(items, item(set.SetName))
			setsIDMap[set.SetName] = set.SetID // Map name -> ID
		}
	}

	const listWidth = 20
	const listHeight = 20

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = fmt.Sprintf("Pick a set from the %s series", seriesID)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SetsModel{
			List:       l,
			SeriesName: seriesID,
			setsIDMap:  setsIDMap,
		},
		nil
}

func callSetsData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("apikey", "sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Authorization", "Bearer sb_publishable_oondaaAIQC-wafhEiNgpSQ_reRiEp7j")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
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
