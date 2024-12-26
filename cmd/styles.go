package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"regexp"
)

// This file holds all lipgloss stylization variables in one spot since they
// are used throughout the package and don't need to be redeclared.
var (
	green       = lipgloss.NewStyle().Foreground(lipgloss.Color("#38B000"))
	red         = lipgloss.NewStyle().Foreground(lipgloss.Color("#D00000"))
	gray        = lipgloss.Color("#777777")
	keyMenu     = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))
	errorColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	errorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	typesTableBorder = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#FFCC00"))
	styleBold   = lipgloss.NewStyle().Bold(true)
	styleItalic = lipgloss.NewStyle().Italic(true)
	colorMap    = map[string]string{
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

// Helper function to get color for a given type name from colorMap
func getTypeColor(typeName string) string {
	color := colorMap[typeName]

	return color
}

// stripANSI function is used in tests to strip ANSI for plain text processing
func stripANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(input, "")
}
