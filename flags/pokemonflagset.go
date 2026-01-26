// This file holds all the flags used by the <pokemonName> subcommand

package flags

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/structs"
	"github.com/digitalghost-dev/poke-cli/styling"
	"github.com/disintegration/imaging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PokemonFlags struct {
	FlagSet        *flag.FlagSet
	Abilities      *bool
	ShortAbilities *bool
	Defense        *bool
	ShortDefense   *bool
	Image          *string
	ShortImage     *string
	Move           *bool
	ShortMove      *bool
	Stats          *bool
	ShortStats     *bool
	Types          *bool
	ShortTypes     *bool
}

func header(header string) string {
	var output strings.Builder

	HeaderBold := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styling.YellowColor).
		BorderTop(true).
		Bold(true).
		Render(header)

	output.WriteString(HeaderBold)

	return output.String()
}

func SetupPokemonFlagSet() *PokemonFlags {
	pf := &PokemonFlags{}
	pf.FlagSet = flag.NewFlagSet("pokeFlags", flag.ExitOnError)

	pf.Abilities = pf.FlagSet.Bool("abilities", false, "Print the Pokémon's abilities")
	pf.ShortAbilities = pf.FlagSet.Bool("a", false, "Print the Pokémon's abilities")

	pf.Defense = pf.FlagSet.Bool("defense", false, "Print the Pokémon's type defenses")
	pf.ShortDefense = pf.FlagSet.Bool("d", false, "Print the Pokémon's type defenses")

	pf.Image = pf.FlagSet.String("image", "", "Print the Pokémon's default sprite")
	pf.ShortImage = pf.FlagSet.String("i", "", "Print the Pokémon's default sprite")

	pf.Move = pf.FlagSet.Bool("moves", false, "Print the Pokémon's learnable moves")
	pf.ShortMove = pf.FlagSet.Bool("m", false, "Print the Pokémon's learnable moves")

	pf.Stats = pf.FlagSet.Bool("stats", false, "Print the Pokémon's base stats")
	pf.ShortStats = pf.FlagSet.Bool("s", false, "Print the Pokémon's base stats")

	pf.Types = pf.FlagSet.Bool("types", false, "Print the Pokémon's typing")
	pf.ShortTypes = pf.FlagSet.Bool("t", false, "Prints the Pokémon's typing")

	hintMessage := styling.StyleItalic.Render("options: [sm, md, lg]")

	pf.FlagSet.Usage = func() {
		helpMessage := styling.HelpBorder.Render("poke-cli pokemon <pokemon-name> [flags]\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-a, --abilities", "Prints the Pokémon's abilities."),
			fmt.Sprintf("\n\t%-30s %s", "-d, --defense", "Prints the Pokémon's type defenses."),
			fmt.Sprintf("\n\t%-30s %s", "-i=xx, --image=xx", "Prints out the Pokémon's default sprite."),
			fmt.Sprintf("\n\t%5s%-15s", "", hintMessage),
			fmt.Sprintf("\n\t%-30s %s", "-m, --moves", "Prints the Pokemon's learnable moves."),
			fmt.Sprintf("\n\t%-30s %s", "-s, --stats", "Prints the Pokémon's base stats."),
			fmt.Sprintf("\n\t%-30s %s", "-t, --types", "Prints the Pokémon's typing."),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints the help menu."),
		)
		fmt.Println(helpMessage)
	}

	return pf
}

func AbilitiesFlag(w io.Writer, endpoint string, pokemonName string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Abilities"))
	if err != nil {
		return err
	}

	for _, pokeAbility := range pokemonStruct.Abilities {
		formattedName := styling.CapitalizeResourceName(pokeAbility.Ability.Name)

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

func DefenseFlag(w io.Writer, endpoint string, pokemonName string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Type Defenses"))
	if err != nil {
		return err
	}

	allTypes := []string{"normal", "fire", "water", "electric", "grass", "ice",
		"fighting", "poison", "ground", "flying", "psychic",
		"bug", "rock", "ghost", "dragon", "dark", "steel", "fairy"}

	typeData := make(map[string]structs.TypesJSONStruct)
	for _, pokeType := range pokemonStruct.Types {
		typeStruct, _, _ := connections.TypesApiCall("type", pokeType.Type.Name, connections.APIURL)
		typeData[pokeType.Type.Name] = typeStruct
	}

	calculateTypeEffectiveness := func(attackingType string) float64 {
		totalEffectiveness := 1.0

		for _, pokeType := range pokemonStruct.Types {
			typeStruct := typeData[pokeType.Type.Name]
			effectiveness := 1.0

			// Check for double damage (weakness)
			for _, dmgType := range typeStruct.DamageRelations.DoubleDamageFrom {
				if dmgType.Name == attackingType {
					effectiveness = 2.0
					break
				}
			}

			// Check for half damage (resistance)
			for _, dmgType := range typeStruct.DamageRelations.HalfDamageFrom {
				if dmgType.Name == attackingType {
					effectiveness = 0.5
					break
				}
			}

			// Check for no damage (immunity)
			for _, dmgType := range typeStruct.DamageRelations.NoDamageFrom {
				if dmgType.Name == attackingType {
					effectiveness = 0.0
					break
				}
			}

			totalEffectiveness *= effectiveness
		}

		return totalEffectiveness
	}

	// Check for abilities that grant immunities or resistances
	checkAbilityEffects := func() {
		abilityImmunities := map[string][]string{
			"flash-fire":    {"fire"},
			"water-absorb":  {"water"},
			"storm-drain":   {"water"},
			"volt-absorb":   {"electric"},
			"motor-drive":   {"electric"},
			"lightning-rod": {"electric"},
			"sap-sipper":    {"grass"},
			"dry-skin":      {"water"},
			"levitate":      {"ground"},
			"earth-eater":   {"ground"},
		}

		abilityResistances := map[string][]string{
			"thick-fat": {"fire", "ice"},
			"heatproof": {"fire"},
		}

		for _, ability := range pokemonStruct.Abilities {
			abilityName := ability.Ability.Name
			formattedAbilityName := styling.CapitalizeResourceName(abilityName)

			if types, exists := abilityImmunities[abilityName]; exists {
				typeList := strings.Join(types, " and ")
				_, err := fmt.Fprintf(w, "%s, with the %s ability, grants it immunity to %s type moves.\n",
					cases.Title(language.English).String(pokemonName), formattedAbilityName, typeList)
				if err != nil {
					return
				}
			}

			if types, exists := abilityResistances[abilityName]; exists {
				typeList := strings.Join(types, " and ")
				_, err := fmt.Fprintf(w, "%s, with the %s ability, grants it resistance to %s type moves.\n",
					cases.Title(language.English).String(pokemonName), formattedAbilityName, typeList)
				if err != nil {
					return
				}
			}
		}
	}

	// Calculate effectiveness for all types
	typeEffectiveness := make(map[string]float64)
	for _, attackingType := range allTypes {
		typeEffectiveness[attackingType] = calculateTypeEffectiveness(attackingType)
	}

	var (
		immune          []string
		quarterDamage   []string
		halfDamage      []string
		normal          []string
		doubleDamage    []string
		quadrupleDamage []string
	)

	for typeName, eff := range typeEffectiveness {
		capitalizedType := cases.Title(language.English).String(typeName)
		switch eff {
		case 0.0:
			immune = append(immune, capitalizedType)
		case 0.25:
			quarterDamage = append(quarterDamage, capitalizedType)
		case 0.5:
			halfDamage = append(halfDamage, capitalizedType)
		case 1.0:
			normal = append(normal, capitalizedType)
		case 2.0:
			doubleDamage = append(doubleDamage, capitalizedType)
		case 4.0:
			quadrupleDamage = append(quadrupleDamage, capitalizedType)
		}
	}

	sort.Strings(immune)
	sort.Strings(quarterDamage)
	sort.Strings(halfDamage)
	sort.Strings(normal)
	sort.Strings(doubleDamage)
	sort.Strings(quadrupleDamage)

	if len(immune) > 0 {
		_, err := fmt.Fprintf(w, "Immune: %s\n", strings.Join(immune, ", "))
		if err != nil {
			return err
		}
	}
	if len(quarterDamage) > 0 {
		_, err := fmt.Fprintf(w, "0.25×   Damage: %s\n", strings.Join(quarterDamage, ", "))
		if err != nil {
			return err
		}
	}
	if len(halfDamage) > 0 {
		_, err := fmt.Fprintf(w, "0.5×    Damage: %s\n", strings.Join(halfDamage, ", "))
		if err != nil {
			return err
		}
	}
	if len(doubleDamage) > 0 {
		_, err := fmt.Fprintf(w, "2.0×    Damage: %s\n", strings.Join(doubleDamage, ", "))
		if err != nil {
			return err
		}
	}
	if len(quadrupleDamage) > 0 {
		_, err := fmt.Fprintf(w, "4.0×    Damage: %s\n", strings.Join(quadrupleDamage, ", "))
		if err != nil {
			return err
		}
	}

	// Add a newline before ability effects if there are any type effectiveness results
	if len(immune) > 0 || len(quarterDamage) > 0 || len(halfDamage) > 0 || len(doubleDamage) > 0 || len(quadrupleDamage) > 0 {
		_, err := fmt.Fprintln(w)
		if err != nil {
			return err
		}
	}

	checkAbilityEffects()

	return nil
}

func ImageFlag(w io.Writer, endpoint string, pokemonName string, size string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Image"))
	if err != nil {
		return err
	}

	// Anonymous function to transform the image to a string
	ToString := func(width int, height int, img image.Image) string {
		img = imaging.Resize(img, width, height, imaging.NearestNeighbor)
		b := img.Bounds()

		imageWidth := b.Max.X
		h := b.Max.Y

		rowCount := (h - 1) / 2
		if h%2 != 0 {
			rowCount++
		}
		estimatedSize := (imageWidth * rowCount * 55) + rowCount

		str := strings.Builder{}
		str.Grow(estimatedSize)

		// Cache for lipgloss styles to avoid recreating identical styles
		styleCache := make(map[string]lipgloss.Style)

		for heightCounter := 0; heightCounter < h-1; heightCounter += 2 {
			for x := 0; x < imageWidth; x++ {
				// Get the color of the current and next row's pixels
				c1, _ := styling.MakeColor(img.At(x, heightCounter))
				color1 := lipgloss.Color(c1.Hex())
				c2, _ := styling.MakeColor(img.At(x, heightCounter+1))
				color2 := lipgloss.Color(c2.Hex())

				styleKey := string(color1) + "_" + string(color2)
				style, exists := styleCache[styleKey]
				if !exists {
					style = lipgloss.NewStyle().Foreground(color1).Background(color2)
					styleCache[styleKey] = style
				}

				str.WriteString(style.Render("▀"))
			}

			str.WriteString("\n")
		}

		return str.String()
	}

	imageResp, err := http.Get(pokemonStruct.Sprites.FrontDefault)
	if err != nil {
		fmt.Println("Error downloading sprite image:", err)
		return err
	}
	defer imageResp.Body.Close()

	img, err := imaging.Decode(imageResp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return err
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
		errMessage := styling.ErrorBorder.Render(styling.ErrorColor.Render("✖ Error!"), "\nInvalid image size.\nValid sizes are: lg, md, sm")
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
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	_, err := fmt.Fprintln(w, header("Learnable Moves"))
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

	movesChan := make(chan MoveInfo)
	errorsChan := make(chan error)

	var wg sync.WaitGroup

	// Count eligible moves for concurrency
	eligibleMoves := 0
	for _, pokeMove := range pokemonStruct.Moves {
		for _, detail := range pokeMove.VersionGroupDetails {
			if detail.VersionGroup.Name != "scarlet-violet" || detail.MoveLearnedMethod.Name != "level-up" {
				continue
			}

			eligibleMoves++
			wg.Add(1)
			go func(moveName string, level int) {
				defer wg.Done()

				moveStruct, _, err := connections.MoveApiCall("move", moveName, connections.APIURL)
				if err != nil {
					errorsChan <- fmt.Errorf("error fetching move %s: %v", moveName, err)
					return
				}

				capitalizedMove := styling.CapitalizeResourceName(moveName)
				capitalizedType := cases.Title(language.English).String(moveStruct.Type.Name)

				movesChan <- MoveInfo{
					Accuracy: moveStruct.Accuracy,
					Level:    level,
					Name:     capitalizedMove,
					Power:    moveStruct.Power,
					Type:     capitalizedType,
				}
			}(pokeMove.Move.Name, detail.LevelLearnedAt)
		}
	}

	// Close channels when all goroutines are done
	go func() {
		wg.Wait()
		close(movesChan)
		close(errorsChan)
	}()

	// Collect results from channels
	movesOpen, errorsOpen := true, true
	for movesOpen || errorsOpen {
		select {
		case move, ok := <-movesChan:
			if !ok {
				movesOpen = false
				continue
			}
			moves = append(moves, move)
		case err, ok := <-errorsChan:
			if !ok {
				errorsOpen = false
				continue
			}
			log.Println(err)
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
		styledType := lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(styling.ColorMap[strings.ToLower(m.Type)])).
			Render(m.Type)

		rows = append(rows, []string{
			m.Name,
			strconv.Itoa(m.Level),
			styledType,
			strconv.Itoa(m.Accuracy),
			strconv.Itoa(m.Power),
		})
	}

	// Build and print table
	color := lipgloss.AdaptiveColor{Light: "#4B4B4B", Dark: "#D3D3D3"}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(color)).
		StyleFunc(func(row, column int) lipgloss.Style {
			var style lipgloss.Style

			switch column {
			case 0:
				style = style.Width(18)
			case 1:
				style = style.Width(8)
			case 2:
				style = style.Width(10)
			case 3:
				style = style.Width(10)
			case 4:
				style = style.Width(8)
			}

			return style
		}).
		Headers("Name", "Level", "Type", "Accuracy", "Power").
		Rows(rows...)

	_, err = fmt.Fprintln(w, t)
	if err != nil {
		return err
	}

	return nil
}

func StatsFlag(w io.Writer, endpoint string, pokemonName string) error {
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

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
	pokemonStruct, _, _ := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)

	// Print the header from header func
	_, err := fmt.Fprintln(w, header("Typing"))
	if err != nil {
		return err
	}

	for _, pokeType := range pokemonStruct.Types {
		colorHex, exists := styling.ColorMap[pokeType.Type.Name]
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

	fmt.Fprintln(w, styling.WarningBorder.Render(styling.WarningColor.Render("⚠ Warning!"), "\nThe '-t | --types' flag is deprecated\nand will be removed in v2.\n\nTyping is now included by default.\nYou no longer need this flag. "))

	return nil
}
