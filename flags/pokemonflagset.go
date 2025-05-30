// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/disintegration/imaging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func header(header string) string {
	var output strings.Builder

	HeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFCC00")).
		BorderTop(true).
		Bold(true).
		Render(header)

	output.WriteString(HeaderBold)

	return output.String()
}

func SetupPokemonFlagSet() (*flag.FlagSet, *bool, *bool, *string, *string, *bool, *bool, *bool, *bool, *bool, *bool) {
	pokeFlags := flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	abilitiesFlag := pokeFlags.Bool("abilities", false, "Print the Pokémon's abilities")
	shortAbilitiesFlag := pokeFlags.Bool("a", false, "Print the Pokémon's abilities")

	imageFlag := pokeFlags.String("image", "", "Print the Pokémon's default sprite")
	shortImageFlag := pokeFlags.String("i", "", "Print the Pokémon's default sprite")

	moveFlag := pokeFlags.Bool("moves", false, "Print the Pokémon's learnable moves")
	shortMoveFlag := pokeFlags.Bool("m", false, "Print the Pokémon's learnable moves")

	statsFlag := pokeFlags.Bool("stats", false, "Print the Pokémon's base stats")
	shortStatsFlag := pokeFlags.Bool("s", false, "Print the Pokémon's base stats")

	typesFlag := pokeFlags.Bool("types", false, "Print the Pokémon's typing")
	shortTypesFlag := pokeFlags.Bool("t", false, "Prints the Pokémon's typing")

	hintMessage := styling.StyleItalic.Render("options: [sm, md, lg]")

	pokeFlags.Usage = func() {
		helpMessage := styling.HelpBorder.Render("poke-cli pokemon <pokemon-name> [flags]\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-i=xx, --image=xx", "Prints out the Pokémon's default sprite."),
			fmt.Sprintf("\n\t%5s%-15s", "", hintMessage),
			fmt.Sprintf("\n\t%-30s %s", "-m, --moves", "Prints the Pokemon's learnable moves."),
			fmt.Sprintf("\n\t%-30s %s", "-s, --stats", "Prints the Pokémon's base stats."),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	return pokeFlags, abilitiesFlag, shortAbilitiesFlag, imageFlag, shortImageFlag, moveFlag, shortMoveFlag, statsFlag, shortStatsFlag, typesFlag, shortTypesFlag
}

func AbilitiesFlag(w io.Writer, endpoint string, pokemonName string) error {
	pokemonStruct, _, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Abilities"))
	if err != nil {
		return err
	}

	// Anonymous function to format ability names
	formatAbilityName := func(name string) string {
		exceptions := map[string]bool{
			"of":  true,
			"the": true,
			"to":  true,
			"as":  true,
		}

		name = strings.ReplaceAll(name, "-", " ")
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

		switch pokeAbility.Slot {
		case 1, 2:
			_, err := fmt.Fprintf(w, "Ability %d: %s\n", pokeAbility.Slot, formattedName)
			if err != nil {
				return err
			}
		default:
			_, err := fmt.Fprintf(w, "Hidden Ability: %s\n", formattedName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ImageFlag(w io.Writer, endpoint string, pokemonName string, size string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Image"))
	if err != nil {
		return err
	}

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
				c1, _ := styling.MakeColor(img.At(x, heightCounter))
				color1 := lipgloss.Color(c1.Hex())
				c2, _ := styling.MakeColor(img.At(x, heightCounter+1))
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
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("Error!"), "\nInvalid image size.\nValid sizes are: lg, md, sm")
		return fmt.Errorf("%s", errMessage)
	}

	imgStr := ToString(dimensions[0], dimensions[1], img)
	_, err = fmt.Fprint(w, imgStr)
	if err != nil {
		return err
	}

	return nil
}

func MovesFlag(w io.Writer, endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	_, err := fmt.Fprintln(w, header("Moves"))
	if err != nil {
		return err
	}

	type MoveInfo struct {
		Accuracy int
		Level    int
		Name     string
		Power    int
		Type     string
	}

	var moves []MoveInfo

	for _, pokeMove := range pokemonStruct.Moves {
		for _, detail := range pokeMove.VersionGroupDetails {
			if detail.VersionGroup.Name != "scarlet-violet" || detail.MoveLearnedMethod.Name != "level-up" {
				continue
			}

			moveName := pokeMove.Move.Name
			moveStruct, _, err := connections.MoveApiCall("move", moveName, baseURL)
			if err != nil {
				log.Printf("Error fetching move %s: %v", moveName, err)
				continue
			}

			moves = append(moves, MoveInfo{
				Accuracy: moveStruct.Accuracy,
				Level:    detail.LevelLearnedAt,
				Name:     moveName,
				Power:    moveStruct.Power,
				Type:     moveStruct.Type.Name,
			})
		}
	}

	if len(moves) == 0 {
		fmt.Fprintln(w, "No level-up moves found for Scarlet & Violet.")
		return nil
	}

	// Sort by level
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].Level < moves[j].Level
	})

	// Convert to table rows
	var rows [][]string
	for _, m := range moves {
		rows = append(rows, []string{
			m.Name,
			m.Type,
			strconv.Itoa(m.Accuracy),
			strconv.Itoa(m.Level),
			strconv.Itoa(m.Power),
		})
	}

	// Build and print table
	color := lipgloss.AdaptiveColor{Light: "#4B4B4B", Dark: "#D3D3D3"}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(color)).
		Headers("Type", "Name", "Accuracy", "Level", "Power").
		Rows(rows...)

	fmt.Fprintln(w, t)
	return nil
}

func StatsFlag(w io.Writer, endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Base Stats"))
	if err != nil {
		return err
	}

	// Anonymous function to map stat values to specific categories
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
		_, err := fmt.Fprintf(w, "%-10s %s %d\n", label, coloredBar, value)
		if err != nil {
			return
		}
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

	// Find the maximum stat value
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

	_, err = fmt.Fprintf(w, "%-10s %d\n", "Total", totalBaseStats)
	if err != nil {
		return err
	}

	return nil
}

func TypesFlag(w io.Writer, endpoint string, pokemonName string) error {
	baseURL := "https://pokeapi.co/api/v2/"
	pokemonStruct, _, _, _, _, _ := connections.PokemonApiCall(endpoint, pokemonName, baseURL)

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

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Typing"))
	if err != nil {
		return err
	}

	for _, pokeType := range pokemonStruct.Types {
		colorHex, exists := colorMap[pokeType.Type.Name]
		if exists {
			color := lipgloss.Color(colorHex)
			style := lipgloss.NewStyle().Bold(true).Foreground(color)
			styledName := style.Render(cases.Title(language.English).String(pokeType.Type.Name))
			_, err := fmt.Fprintf(w, "Type %d: %s\n", pokeType.Slot, styledName)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "Type %d: %s\n", pokeType.Slot, cases.Title(language.English).String(pokeType.Type.Name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
