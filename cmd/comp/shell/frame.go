package shell

import (
	"strings"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/styling"
)

const keyMenu = "← → (switch tab) • b (back) • w (web) • ctrl+c | esc (quit)"

var captionStyle = lipgloss.NewStyle().Foreground(styling.Gray).Italic(true)

func (s *Styles) Render(tabs []string, activeTab, width int, renderContent func(contentWidth int) string) string {
	doc := strings.Builder{}

	var renderedTabs []string
	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == activeTab
		if isActive {
			style = s.ActiveTab
		} else {
			style = s.InactiveTab
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	windowWidth := max(width-8, lipgloss.Width(row)-2)
	contentWidth := windowWidth - 2

	fillWidth := windowWidth - lipgloss.Width(row)
	if fillWidth > 0 {
		fill := lipgloss.NewStyle().Foreground(styling.ThemeColor).
			Render(strings.Repeat("─", fillWidth-1) + "┐")
		row = row + fill
	}

	content := renderContent(contentWidth)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(s.Window.Width(windowWidth).Render(content))
	doc.WriteString("\n")
	doc.WriteString(styling.KeyMenu.Render(keyMenu))

	return s.Doc.Render(doc.String())
}

func TableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.ThemeColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000")).
		Background(styling.ThemeColor)
	return s
}
