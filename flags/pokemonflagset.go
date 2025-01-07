// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"image"
	"net/http"
	"os"
	"strings"
)

var (
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	styleBold   = lipgloss.NewStyle().Bold(true)
	errorColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	errorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	styleItalic = lipgloss.NewStyle().Italic(true)
)

func header(header string) {
	HeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render(header)

	fmt.Println(HeaderBold)
}

func SetupPokemonFlagSet() (*flag.FlagSet, *bool, *bool, *string, *string, *bool, *bool, *bool, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	abilitiesFlag := pokeFlags.Bool("abilities", false, "Print the Pokémon's abilities")
	shortAbilitiesFlag := pokeFlags.Bool("a", false, "Print the Pokémon's abilities")

	imageFlag := pokeFlags.String("image", "", "Print the Pokémon's default sprite")
	shortImageFlag := pokeFlags.String("i", "", "Print the Pokémon's default sprite")

	statsFlag := pokeFlags.Bool("stats", false, "Print the Pokémon's base stats")
	shortStatsFlag := pokeFlags.Bool("s", false, "Print the Pokémon's base stats")

	typesFlag := pokeFlags.Bool("types", false, "Print the Pokémon's typing")
	shortTypesFlag := pokeFlags.Bool("t", false, "Prints the Pokémon's typing")

	hintMessage := styleItalic.Render("options: [sm, md, lg]")

	pokeFlags.Usage = func() {
		helpMessage := helpBorder.Render("poke-cli pokemon <pokemon-name> [flags]\n\n",
			styleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-i=xx, --image=xx", "Prints out the Pokémon's default sprite."),
			fmt.Sprintf("\n\t%5s%-15s", "", hintMessage),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	return pokeFlags, abilitiesFlag, shortAbilitiesFlag, imageFlag, shortImageFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag
}

func AbilitiesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	header("Abilities")

	// Anonymous function to format ability names
	formatAbilityName := func(name string) string {
		exceptions := map[string]bool{
			"of":  true,
			"the": true,
			"to":  true,
			"as":  true,
		}

		name = strings.Replace(name, "-", " ", -1)
		words := strings.Split(name, " ")
		titleCaser := cases.Title(language.English)

		// Process each word
		for i, word := range words {
			if _, found := exceptions[strings.ToLower(word)]; found && i != 0 {
				words[i] = strings.ToLower(word)
			} else {
				words[i] = titleCaser.String(word)
			}
		}
		return strings.Join(words, " ")
	}

	for _, pokeAbility := range pokemonStruct.Abilities {
		formattedName := formatAbilityName(pokeAbility.Ability.Name)
		if pokeAbility.Slot == 1 {
			fmt.Printf("Ability %d: %s\n", pokeAbility.Slot, formattedName)
		} else if pokeAbility.Slot == 2 {
			fmt.Printf("Ability %d: %s\n", pokeAbility.Slot, formattedName)
		} else {
			fmt.Printf("Hidden Ability: %s\n", formattedName)
		}
	}

	return nil
}

func ImageFlag(endpoint string, pokemonName string, size string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	header("Image")

	// Anonymous function to transform the image to a string
	// ToString generates an ASCII representation of the image with color
	ToString := func(width int, height int, img image.Image) string {
		// Resize the image to the specified width, preserving aspect ratio
		img = imaging.Resize(img, width, height, imaging.NearestNeighbor)
		b := img.Bounds()
		imageWidth := b.Max.X - 2 // Adjust width to exclude margins
		h := b.Max.Y - 4          // Adjust height to exclude margins
		str := strings.Builder{}

		// Loop through the image pixels, two rows at a time
		for heightCounter := 2; heightCounter < h; heightCounter += 2 {
			for x := 1; x < imageWidth; x++ {
				// Get the color of the current and next row's pixels
				c1, _ := colorful.MakeColor(img.At(x, heightCounter))
				color1 := lipgloss.Color(c1.Hex())
				c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
				color2 := lipgloss.Color(c2.Hex())

				// Render the half-block character with the two colors
				str.WriteString(lipgloss.NewStyle().
					Foreground(color1).
					Background(color2).
					Render("▀"))
			}

			// Add a newline after each row
			str.WriteString("\n")
		}

		return str.String()
	}

	imageResp, err := http.Get(pokemonStruct.Sprites.FrontDefault)
	if err != nil {
		fmt.Println("Error downloading sprite image:", err)
		os.Exit(1)
	}
	defer imageResp.Body.Close()

	img, err := imaging.Decode(imageResp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		os.Exit(1)
	}

	// Define size map
	sizeMap := map[string][2]int{
		"lg": {120, 120},
		"md": {90, 90},
		"sm": {55, 55},
	}

	// Validate size
	dimensions, exists := sizeMap[strings.ToLower(size)]
	if !exists {
		errMessage := errorBorder.Render(errorColor.Render("Error!"), "\nInvalid image size. Valid sizes are: lg, md, sm")
		return fmt.Errorf("%s", errMessage)
	}

	imgStr := ToString(dimensions[0], dimensions[1], img)
	fmt.Println(imgStr)

	return nil
}

func StatsFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	header("Base Stats")

	// Anonymous function to map stat values to categories
	getStatCategory := func(value int) string {
		switch {
		case value < 20:
			return "lowest"
		case value < 60:
			return "lower"
		case value < 90:
			return "low"
		case value < 120:
			return "high"
		case value < 150:
			return "higher"
		default:
			return "highest"
		}
	}

	// Helper function to print the bar for a stat
	printBar := func(label string, value, maxWidth, maxValue int, style lipgloss.Style) {
		scaledValue := (value * maxWidth) / maxValue
		bar := strings.Repeat("▇", scaledValue)
		coloredBar := style.Render(bar)
		fmt.Printf("%-10s %s %d\n", label, coloredBar, value)
	}

	// Mapping from API stat names to custom display names
	nameMapping := map[string]string{
		"hp":              "HP",
		"attack":          "Atk",
		"defense":         "Def",
		"special-attack":  "Sp. Atk",
		"special-defense": "Sp. Def",
		"speed":           "Speed",
	}

	statColorMap := map[string]string{
		"lowest":  "#F34444",
		"lower":   "#FF7F0F",
		"low":     "#FFDD57",
		"high":    "#A0E515",
		"higher":  "#22C65A",
		"highest": "#00C2B8",
	}

	// Find the maxium stat value
	maxValue := 0
	for _, stat := range pokemonStruct.Stats {
		if stat.BaseStat > maxValue {
			maxValue = stat.BaseStat
		}
	}

	maxWidth := 45

	// Print bars for each stat
	for _, stat := range pokemonStruct.Stats {
		apiName := stat.Stat.Name
		customName, exists := nameMapping[apiName]
		if !exists {
			continue
		}

		category := getStatCategory(stat.BaseStat)
		color := statColorMap[category]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

		printBar(customName, stat.BaseStat, maxWidth, maxValue, style)
	}

	totalBaseStats := 0
	for _, stat := range pokemonStruct.Stats {
		totalBaseStats += stat.BaseStat
	}

	fmt.Printf("%-10s %d\n", "Total", totalBaseStats)

	return nil
}

func TypesFlag(endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	colorMap := map[string]string{
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

	header("Typing")

	for _, pokeType := range pokemonStruct.Types {
		colorHex, exists := colorMap[pokeType.Type.Name]
		if exists {
			color := lipgloss.Color(colorHex)
			style := lipgloss.NewStyle().Bold(true).Foreground(color)
			styledName := style.Render(cases.Title(language.English).String(pokeType.Type.Name)) // Apply styling here
			fmt.Printf("Type %d: %s\n", pokeType.Slot, styledName)                               // Interpolate styled text
		} else {
			fmt.Printf("Type %d: %s\n", pokeType.Slot, cases.Title(language.English).String(pokeType.Type.Name))
		}
	}

	return nil
}
