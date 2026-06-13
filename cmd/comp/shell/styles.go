package shell

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
)

type Styles struct {
	Doc            lipgloss.Style
	InactiveTab    lipgloss.Style
	ActiveTab      lipgloss.Style
	Window         lipgloss.Style
	HighlightColor color.Color
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
	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(isDark)
	highlightColor := ld(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))

	s := new(Styles)
	s.Doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.InactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlightColor).
		Padding(0, 1)
	s.ActiveTab = s.InactiveTab.
		Border(activeTabBorder, true)
	s.Window = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(2, 0).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	s.HighlightColor = highlightColor
	return s
}
