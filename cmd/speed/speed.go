package speed

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DefaultSpeedStat is the default implementation of SpeedStatFunc
var DefaultSpeedStat SpeedStatFunc = func(name string) (string, error) {
	pokemonStruct, _, err := connections.PokemonApiCall("pokemon", name, connections.APIURL)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	for _, stat := range pokemonStruct.Stats {
		if stat.Stat.Name == "speed" {
			return strconv.Itoa(stat.BaseStat), nil
		}
	}

	return "", errors.New("speed stat not found")
}

var (
	pokemon            PokemonDetails
	output             strings.Builder
	abilityMultipliers = map[string]float64{
		"None":         1.0,
		"Swift Swim":   2.0,
		"Chlorophyll":  2.0,
		"Sand Rush":    2.0,
		"Slush Rush":   2.0,
		"Unburden":     2.0,
		"Quick Feet":   1.5,
		"Surge Surfer": 2.0,
	}
	modifierMultipliers = map[string]float64{
		"Choice Scarf": 1.5,
		"Tailwind":     2.0,
	}
	stageMultipliers = map[int]float64{
		-6: 0.25,
		-5: 2.0 / 7.0,
		-4: 1.0 / 3.0,
		-3: 0.4,
		-2: 0.5,
		-1: 2.0 / 3.0,
		0:  1.0,
		+1: 1.5,
		+2: 2.0,
		+3: 2.5,
		+4: 3.0,
		+5: 3.5,
		+6: 4.0,
	}
	natureMultipliers = map[string]float64{
		"+10%": 1.1,
		"0%":   1.0,
		"-10%": 0.9,
	}
)

type PokemonDetails struct {
	Name       string
	SpeedStage string
	Nature     string
	Level      string
	Modifier   []string
	Ability    string
	SpeedEV    string
	SpeedIV    string
}

// SpeedStatFunc is a function type for getting a Pokémon's base speed stat
type SpeedStatFunc func(name string) (string, error)

func SpeedCommand() (string, error) {
	// Reset the output string builder
	output.Reset()

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"Calculate the speed of a Pokémon.\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("speed"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		output.WriteString(helpMessage)
	}

	flag.Parse()

	// Handle help flag
	if len(os.Args) == 3 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		flag.Usage()
		return output.String(), nil
	}

	// Validate arguments
	if err := utils.ValidateSpeedArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	form := form()

	err := form.Run()
	if err != nil {
		return "", err
	}

	result, err := formula()
	if err != nil {
		return "", err
	}

	return result, nil
}

func form() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Speed Calculator").
			Description("This command will calculate the speed stat\nof a Pokémon during an in-game battle.").
			Next(true).
			NextLabel("Next"),
		),
		huh.NewGroup(
			huh.NewInput().
				Value(&pokemon.Name).
				Title("Enter the first Pokémon's name:").
				Placeholder("incineroar").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("input cannot be blank")
					}
					_, _, err := connections.PokemonApiCall("pokemon", s, connections.APIURL)
					if err != nil {
						return errors.New("not a valid Pokémon")
					}
					return nil
				}),
			huh.NewInput().
				Value(&pokemon.Level).
				Title("And its level:").
				Placeholder("50").
				Validate(func(s string) error {
					num, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("please enter a valid number")
					}
					if num < 1 || num > 100 {
						return errors.New("level must be between 1 and 100")
					}
					return nil
				}),
			huh.NewInput().
				Value(&pokemon.SpeedEV).
				Title("EV amount:").
				Description("Enter the Pokémon's EV level for the speed stat").
				Placeholder("252").
				Validate(func(s string) error {
					num, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("please enter a valid number")
					}
					if num < 0 || num > 252 {
						return errors.New("level must be between 0 and 252")
					}
					return nil
				}),
			huh.NewInput().
				Value(&pokemon.SpeedIV).
				Title("IV amount").
				Description("Enter the Pokémon's IV level for the speed stat").
				Placeholder("31").
				Validate(func(s string) error {
					num, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("please enter a valid number")
					}
					if num < 0 || num > 31 {
						return errors.New("level must be between 0 and 31")
					}
					return nil
				}),
			huh.NewMultiSelect[string]().
				Title("Modifiers").
				Description("Select any amount of options").
				Options(
					huh.NewOption("Choice Scarf", "Choice Scarf"),
					huh.NewOption("Tailwind", "Tailwind"),
				).Value(&pokemon.Modifier),
		),
		// Page 3
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Ability").
				Description("Select an ability in play").
				Value(&pokemon.Ability).
				Options(
					huh.NewOptions("None", "Swift Swim", "Chlorophyll", "Sand Rush", "Slush Rush", "Unburden", "Quick Feet", "Surge Surfer")...,
				),
			huh.NewSelect[string]().
				Options(
					huh.NewOptions("+10%", "0%", "-10%")...,
				).
				Title("Nature").
				Description("Nature benefit/detriment 1").
				Value(&pokemon.Nature).
				Validate(func(value string) error {
					if value == "" {
						return errors.New("please select a nature")
					}
					return nil
				}),
			huh.NewInput().
				Value(&pokemon.SpeedStage).
				Title("Speed Stage").
				Placeholder("ex: +6 or 0 or -3").
				Validate(func(s string) error {
					num, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("please enter a whole number between -6 and +6")
					}
					if num < -6 || num > 6 {
						return errors.New("level must be between -6 and +6")
					}
					return nil
				}),
		),
	).WithTheme(styling.FormTheme())

	return form
}

func formula() (string, error) {
	modifierMultiplier := 1.0 // start with no change
	for _, mod := range pokemon.Modifier {
		if val, ok := modifierMultipliers[mod]; ok {
			modifierMultiplier *= val
		}
	}

	abilityMultiplier := abilityMultipliers[pokemon.Ability]

	speedStageInt, err := strconv.Atoi(pokemon.SpeedStage)
	if err != nil {
		log.Fatalf("Invalid SpeedStage: %v", err)
	}
	stageMultiplier := stageMultipliers[speedStageInt]

	// Get the Pokémon's base speed using the DefaultSpeedStat function
	speedStr, err := DefaultSpeedStat(pokemon.Name)
	if err != nil {
		return "", err
	}

	baseSpeed, err := strconv.Atoi(speedStr)
	if err != nil {
		return "", fmt.Errorf("failed to convert speed to int: %w", err)
	}

	intIV, err := strconv.Atoi(pokemon.SpeedIV)
	if err != nil {
		return "", fmt.Errorf("failed to convert speed to int: %w", err)
	}
	intEV, err := strconv.Atoi(pokemon.SpeedEV)
	if err != nil {
		return "", fmt.Errorf("failed to convert speed to int: %w", err)
	}

	intLevel, err := strconv.Atoi(pokemon.Level)
	if err != nil {
		return "", fmt.Errorf("failed to convert speed to int: %w", err)
	}

	natureMultiplier := natureMultipliers[pokemon.Nature]

	chosenPokemon := cases.Title(language.English).String(pokemon.Name)

	// Calculate final speed using the formula:
	// (((2 x base + IV + (EV / 4)) x level / 100 + 5) * nature * modifier * ability * stage
	finalSpeed := float64(((2*baseSpeed+intIV+(intEV/4))*intLevel/100)+5) * natureMultiplier * abilityMultiplier * stageMultiplier * modifierMultiplier

	// Round down the final speed
	finalSpeedFloor := math.Floor(finalSpeed)
	finalSpeedStr := fmt.Sprintf("%.0f", finalSpeedFloor)

	header := fmt.Sprintf("%s at level %s with selected options has a current speed of %s.",
		styling.Yellow.Render(chosenPokemon),
		styling.Yellow.Render(pokemon.Level),
		styling.Yellow.Render(finalSpeedStr),
	)
	body := fmt.Sprintf("EVs: %s\nIVs: %s\nModifiers: %s\nNature: %s\nAbility: %s\nSpeed Stage: %s\nBase Speed: %s",
		styling.Yellow.Render(pokemon.SpeedEV),
		styling.Yellow.Render(pokemon.SpeedIV),
		styling.Yellow.Render(xstrings.EnglishJoin(pokemon.Modifier, true)),
		styling.Yellow.Render(pokemon.Nature),
		styling.Yellow.Render(pokemon.Ability),
		styling.Yellow.Render(pokemon.SpeedStage),
		styling.Yellow.Render(speedStr),
	)

	docStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#444", Dark: "#EEE"}).
		Width(32)

	fullDoc := lipgloss.JoinVertical(lipgloss.Top, header, "---", body)
	output.WriteString(docStyle.Render(fullDoc))

	return output.String(), nil
}
