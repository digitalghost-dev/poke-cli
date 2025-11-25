package card

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ImageModel struct {
	CardName string
	ImageURL string
	Error    error
}

func (m ImageModel) Init() tea.Cmd {
	return nil
}

func (m ImageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ImageModel) View() string {
	return m.ImageURL
}

func ImageRenderer(cardName string, imageURL string) ImageModel {
	imageData, err := CardImage(imageURL)

	return ImageModel{
		CardName: cardName,
		ImageURL: imageData,
		Error:    err,
	}
}
