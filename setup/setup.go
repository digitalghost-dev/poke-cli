package setup

import (
	_ "embed"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/flags"
	"github.com/digitalghost-dev/poke-cli/styling"
)

//go:embed snorlax.txt
var snorlax string

const (
	rowTheme = iota
	rowCacheWarn
	rowReleases
	rowSave
	rowCount
)

const releasesURL = "https://github.com/digitalghost-dev/poke-cli/releases"

var themeChoices = []string{flags.ThemeYellow, flags.ThemeRed, flags.ThemeBlue}

type model struct {
	cfg     flags.Config
	cursor  int
	cacheOK bool
	width   int
	height  int
	saved   bool
}

func newModel(cfg flags.Config) model {
	_, err := exec.LookPath("poke-cache")
	return model{cfg: cfg, cacheOK: err == nil}
}

func Run(cfg flags.Config) (flags.Config, bool, error) {
	out, err := tea.NewProgram(newModel(cfg)).Run()
	if err != nil {
		return cfg, false, err
	}
	final, ok := out.(model)
	if !ok || !final.saved {
		styling.ApplyTheme(cfg.Display.Theme)
		return cfg, false, nil
	}
	styling.ApplyTheme(final.cfg.Display.Theme)
	return final.cfg, true, nil
}

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case "down", "j":
		if m.cursor < rowCount-1 {
			m.cursor++
		}
		return m, nil
	case "left", "right":
		m.adjust(msg.String() == "right")
		return m, nil
	case "enter", "space":
		switch m.cursor {
		case rowReleases:
			return m, utils.Open(releasesURL)
		case rowSave:
			m.saved = true
			return m, tea.Quit
		default:
			m.adjust(true)
		}
		return m, nil
	}
	return m, nil
}

func (m *model) adjust(forward bool) {
	switch m.cursor {
	case rowTheme:
		m.cfg.Display.Theme = cycle(themeChoices, m.cfg.Display.Theme, forward)
		styling.ApplyTheme(m.cfg.Display.Theme)
	case rowCacheWarn:
		m.cfg.Cache.ShowWarning = !m.cfg.Cache.ShowWarning
	}
}

func cycle(opts []string, cur string, forward bool) string {
	i := 0
	for j, o := range opts {
		if o == cur {
			i = j
			break
		}
	}
	if forward {
		i = (i + 1) % len(opts)
	} else {
		i = (i - 1 + len(opts)) % len(opts)
	}
	return opts[i]
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) View() tea.View {
	welcomeMessage := "Welcome! Please choose some quick prefrences.\n\n"

	row := func(label, value string, focused int) string {
		cursor := "  "
		line := lipgloss.NewStyle().Render(label + "  " + value)
		if m.cursor == focused {
			cursor = "> "
			line = styling.Yellow.Render(label + "  " + value)
		}
		return cursor + line
	}

	settings := lipgloss.JoinVertical(lipgloss.Left,
		row("Theme", m.cfg.Display.Theme, rowTheme),
		row("Cache warning", onOff(m.cfg.Cache.ShowWarning), rowCacheWarn),
		"",
		row("Open releases page", "", rowReleases),
		row("Save & quit", "", rowSave),
	)

	cache := "not found"
	if m.cacheOK {
		cache = styling.Green.Render("✓ installed")
	}
	status := "poke-cache: " + cache

	panel := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).
		BorderForeground(styling.YellowColor).Padding(1, 2)

	body := lipgloss.JoinVertical(lipgloss.Left, panel.Render(settings), panel.Render(status))

	content := body
	if m.width == 0 || m.width >= lipgloss.Width(snorlax)+lipgloss.Width(body)+4 {
		content = lipgloss.JoinHorizontal(lipgloss.Top, body, "  ", snorlax)
	}
	content += "\n" + styling.KeyMenu.Render("↑/↓ move • ←/→ change • enter select • esc quit")

	v := tea.NewView(welcomeMessage + content)
	v.AltScreen = true
	return v
}

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
