package shell

import (
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

type Styles struct {
	Doc         lipgloss.Style
	InactiveTab lipgloss.Style
	ActiveTab   lipgloss.Style
	Window      lipgloss.Style
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func NewStyles() *Styles {
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	s := new(Styles)
	s.Doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.InactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(styling.ThemeColor).
		Padding(0, 1)
	s.ActiveTab = s.InactiveTab.
		Border(activeTabBorder, true)
	s.Window = lipgloss.NewStyle().
		BorderForeground(styling.ThemeColor).
		Padding(2, 0).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	return s
}
