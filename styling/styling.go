package styling

import (
	"github.com/charmbracelet/lipgloss"
	"regexp"
)

var (
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("#38B000"))
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("#D00000"))
	Gray    = lipgloss.Color("#777777")
	KeyMenu = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))

	StyleBold      = lipgloss.NewStyle().Bold(true)
	StyleItalic    = lipgloss.NewStyle().Italic(true)
	StyleUnderline = lipgloss.NewStyle().Underline(true)
	HelpBorder     = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	ErrorColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	ErrorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	TypesTableBorder = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#FFCC00"))
	ColorMap = map[string]string{
		"normal":   "#B7B7A9",
		"fire":     "#FF4422",
		"water":    "#3499FF",
		"electric": "#FFCC33",
		"grass":    "#77CC55",
		"ice":      "#66CCFF",
		"fighting": "#BB5544",
		"poison":   "#AA5699",
		"ground":   "#DEBB55",
		"flying":   "#889AFF",
		"psychic":  "#FF5599",
		"bug":      "#AABC22",
		"rock":     "#BBAA66",
		"ghost":    "#6666BB",
		"dragon":   "#7766EE",
		"dark":     "#775544",
		"steel":    "#AAAABB",
		"fairy":    "#EE99EE",
	}
)

// GetTypeColor Helper function to get color for a given type name from colorMap
func GetTypeColor(typeName string) string {
	color := ColorMap[typeName]

	return color
}

// StripANSI function is used in tests to strip ANSI for plain text processing
func StripANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(input, "")
}
