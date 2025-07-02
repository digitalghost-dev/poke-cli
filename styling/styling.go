package styling

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"image/color"
	"regexp"
)

var (
	Green         = lipgloss.NewStyle().Foreground(lipgloss.Color("#38B000"))
	Red           = lipgloss.NewStyle().Foreground(lipgloss.Color("#D00000"))
	Gray          = lipgloss.Color("#777777")
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

// To avoid unnecessary dependencies, I adapted the MakeColor function from
// "github.com/lucasb-eyer/go-colorful" and implemented it using only the
// standard library. Since I only needed this function, importing the entire
// library was unnecessary.
type Color struct {
	R, G, B float64
}

// Implement the Go color.Color interface.
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
