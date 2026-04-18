package card

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type imageModel struct {
	CardName  string
	ImageURL  string
	Error     error
	Loading   bool
	Spinner   spinner.Model
	ImageData string
	Protocol  string
}

type imageReadyMsg struct {
	imageData string
	protocol  string
	err       error
}

// fetchImageCmd downloads and renders the image asynchronously
func fetchImageCmd(imageURL string) tea.Cmd {
	return func() tea.Msg {
		imageData, protocol, err := CardImage(imageURL)
		if err != nil {
			return imageReadyMsg{err: err}
		}
		return imageReadyMsg{
			imageData: imageData,
			protocol:  protocol,
		}
	}
}

func (m imageModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		fetchImageCmd(m.ImageURL),
	)
}

func (m imageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case imageReadyMsg:
		m.Loading = false
		if msg.err != nil {
			m.Error = msg.err
			m.ImageData = ""
		} else {
			m.Error = nil
			m.ImageData = msg.imageData
			m.Protocol = msg.protocol
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m imageModel) View() tea.View {
	var content string
	if m.Loading {
		content = lipgloss.NewStyle().Padding(2).Render(
			m.Spinner.View() + "Loading image for \n" + m.CardName,
		)
	} else if m.Error != nil {
		// Styling the error message with padding for better readability
		content = lipgloss.NewStyle().
			Padding(2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(styling.YellowColor).
			Render(styling.Red.Render(m.Error.Error()))
	} else {
		content = m.ImageData
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func ImageRenderer(cardName string, imageURL string) imageModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styling.Yellow

	return imageModel{
		CardName: cardName,
		ImageURL: imageURL,
		Loading:  true,
		Spinner:  s,
	}
}
