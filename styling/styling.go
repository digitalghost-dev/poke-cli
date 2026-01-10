package styling

import (
	"fmt"
	"image/color"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	HyphenHint = "Use a hyphen when typing a name with a space."
)

var (
	Green         = lipgloss.NewStyle().Foreground(lipgloss.Color("#38B000"))
	Red           = lipgloss.NewStyle().Foreground(lipgloss.Color("#D00000"))
	Gray          = lipgloss.Color("#777777")
	Yellow        = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#E1AD01", Dark: "#FFDE00"})
	ColoredBullet = lipgloss.NewStyle().
			SetString("•").
			Foreground(lipgloss.Color("#FFCC00"))
	CheckboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00"))
	KeyMenu       = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))

	DocsLink = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#E1AD01", Dark: "#FFCC00"}).
			Render("\x1b]8;;https://docs.poke-cli.com\x1b\\docs.poke-cli.com\x1b]8;;\x1b\\")

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
	WarningColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00"))
	WarningBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF8C00"))
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
	typeColor := ColorMap[typeName]

	return typeColor
}

// StripANSI function is used in tests to strip ANSI for plain text processing
func StripANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(input, "")
}

// smallWords are words that should remain lowercase in titles (unless first word)
var smallWords = map[string]bool{
	"of":  true,
	"the": true,
	"to":  true,
	"as":  true,
}

// CapitalizeResourceName converts hyphenated resource names to title case
// Example: "strong-jaw" -> "Strong Jaw", "sword-of-ruin" -> "Sword of Ruin"
func CapitalizeResourceName(name string) string {
	caser := cases.Title(language.English)

	name = strings.ReplaceAll(name, "-", " ")
	words := strings.Split(name, " ")

	for i, word := range words {
		if _, found := smallWords[strings.ToLower(word)]; found && i != 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = caser.String(word)
		}
	}

	return strings.Join(words, " ")
}

// Color To avoid unnecessary dependencies, I adapted the MakeColor function from
// "github.com/lucasb-eyer/go-colorful" and implemented it using only the
// standard library. Since I only needed this function, importing the entire
// library was unnecessary.
type Color struct {
	R, G, B float64
}

// RGBA Implement the Go color.Color interface.
func (col Color) RGBA() (uint32, uint32, uint32, uint32) {
	return uint32(col.R*65535.0 + 0.5), uint32(col.G*65535.0 + 0.5), uint32(col.B*65535.0 + 0.5), 0xFFFF
}

// MakeColor constructs a Color from a color.Color.
func MakeColor(c color.Color) (Color, bool) {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return Color{}, false
	}

	// Undo alpha pre-multiplication
	return Color{
		R: float64(r) / float64(a),
		G: float64(g) / float64(a),
		B: float64(b) / float64(a),
	}, true
}

// Hex returns the hex representation of the color, like "#ff0080".
func (col Color) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x",
		uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
}

func FormTheme() *huh.Theme {
	var (
		yellow   = lipgloss.Color("#FFDE00")
		blue     = lipgloss.Color("#3B4CCA")
		red      = lipgloss.Color("#D00000")
		black    = lipgloss.Color("#000000")
		normalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
	)
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(blue).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(blue).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(blue)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(red)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(yellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(yellow)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(red)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(red)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(red).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(black).Background(yellow)
	t.Focused.Next = t.Focused.FocusedButton

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(yellow)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(red)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
