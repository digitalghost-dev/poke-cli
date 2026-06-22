package styling

import (
	"fmt"
	"image/color"
	"os"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"golang.org/x/term"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const HyphenHint = "Use a hyphen when typing a name with a space."

var palettes = map[string]string{
	"yellow": "#E1AD01",
	"red":    "#f00000",
	"blue":   "#3B4CCA",
}

var accent string

var (
	YellowColor     color.Color
	YellowAdaptive  color.Color
	YellowAdaptive2 color.Color
)

var (
	Green         = lipgloss.NewStyle().Foreground(lipgloss.Color("#38B000"))
	Red           = lipgloss.NewStyle().Foreground(lipgloss.Color("#D00000"))
	Gray          = lipgloss.Color("#777777")
	Yellow        lipgloss.Style
	ColoredBullet lipgloss.Style
	CheckboxStyle lipgloss.Style
	KeyMenu       = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))

	DocsLink string

	StyleBold      = lipgloss.NewStyle().Bold(true)
	StyleItalic    = lipgloss.NewStyle().Italic(true)
	StyleUnderline = lipgloss.NewStyle().Underline(true)
	HelpBorder     lipgloss.Style
	ErrorColor     = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	ErrorBorder    = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	ApiErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(1, 2)
	WarningColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8C00"))
	WarningBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF8C00"))
	TypesTableBorder lipgloss.Style
	ColorMap         = map[string]string{
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

func ApplyTheme(name string) {
	hex, ok := palettes[name]
	if !ok {
		hex = palettes["yellow"]
	}
	accent = hex
	c := lipgloss.Color(hex)

	YellowColor = c
	YellowAdaptive = c
	YellowAdaptive2 = c
	Yellow = lipgloss.NewStyle().Foreground(c)
	ColoredBullet = lipgloss.NewStyle().SetString("•").Foreground(c)
	CheckboxStyle = lipgloss.NewStyle().Foreground(c)
	HelpBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(c)
	TypesTableBorder = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(c)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(c)
	DocsLink = lipgloss.NewStyle().
		Foreground(c).
		Render("\x1b]8;;https://docs.poke-cli.com\x1b\\docs.poke-cli.com\x1b]8;;\x1b\\")
}

func init() {
	ApplyTheme("yellow")
}

func HasDarkBackground() bool {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) { // #nosec G115
		return true
	}
	return lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
}

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
