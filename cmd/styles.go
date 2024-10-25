package cmd

import "github.com/charmbracelet/lipgloss"

// This file holds all lipgloss stylization variables in one spot since they
// are used throughout the package and don't need to be redeclared.
var (
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
