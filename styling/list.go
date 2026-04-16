package styling

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle lipgloss.Style
	PaginationStyle   lipgloss.Style
	HelpStyle         lipgloss.Style
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func init() {
	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(isDark)
	defaults := list.DefaultStyles(isDark)
	PaginationStyle = defaults.PaginationStyle.PaddingLeft(4)
	HelpStyle = defaults.HelpStyle.PaddingLeft(4).PaddingBottom(1)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(ld(lipgloss.Color(DarkYellow), lipgloss.Color(LightYellow)))
}

type Item string

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := string(i)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
