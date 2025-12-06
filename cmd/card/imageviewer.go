package card

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type ImageModel struct {
	CardName  string
	ImageURL  string
	Error     error
	Loading   bool
	Spinner   spinner.Model
	ImageData string
}

type imageReadyMsg struct {
	sixelData string
}

// fetchImageCmd downloads and renders the image asynchronously
func fetchImageCmd(imageURL string) tea.Cmd {
	return func() tea.Msg {
		sixelData, err := CardImage(imageURL)
		if err != nil {
			return imageReadyMsg{err.Error()}
		}
		return imageReadyMsg{sixelData: sixelData}
	}
}

func (m ImageModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchImageCmd(m.ImageURL),
	)
}

func (m ImageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case imageReadyMsg:
		m.Loading = false
		m.ImageData = msg.sixelData
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ImageModel) View() string {
	if m.Loading {
		return lipgloss.NewStyle().Padding(2).Render(
			m.Spinner.View() + "Loading image for \n" + m.CardName,
		)
	}
	return m.ImageData
}

func ImageRenderer(cardName string, imageURL string) ImageModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return ImageModel{
		CardName: cardName,
		ImageURL: imageURL,
		Loading:  true,
		Spinner:  s,
	}
}
