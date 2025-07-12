package speed

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	xstrings "github.com/charmbracelet/x/exp/strings"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"math"
	"strconv"
	"strings"
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

func SpeedCommand() (string, error) {
	var (
		pokemon PokemonDetails
	)

	form := huh.NewForm(
		//TODO: add welcome screen
		// Page 2
		huh.NewGroup(
			huh.NewInput().
				Value(&pokemon.Name).
				Title("Enter the first Pokémon's name:").
				Placeholder("cacturne").
				Validate(func(s string) error {
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
					if num < 1 || num > 252 {
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
					if num < 1 || num > 31 {
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
	)

	err := form.Run()
	if err != nil {
		return "", err // return the error to the caller
	}

	// Initialize output builder
	var output strings.Builder

	// Define multiplier maps for speed calculation
	var modifierMultipliers = map[string]float64{
		"Choice Scarf": 1.5,
		"Tailwind":     2.0,
	}

	var abilityMultipliers = map[string]float64{
		"None":         1.0,
		"Swift Swim":   2.0,
		"Chlorophyll":  2.0,
		"Sand Rush":    2.0,
		"Slush Rush":   2.0,
		"Unburden":     2.0,
		"Quick Feet":   1.5,
		"Surge Surfer": 2.0,
	}

	// Calculate modifier multiplier
	modifierMultiplier := 1.0 // start with no change
	for _, mod := range pokemon.Modifier {
		if val, ok := modifierMultipliers[mod]; ok {
			modifierMultiplier *= val
		}
	}

	// Get ability multiplier from the map
	abilityMultiplier := abilityMultipliers[pokemon.Ability]

	// Define stage multipliers for speed calculation
	var stageMultipliers = map[int]float64{
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

	// Convert speed stage to integer and get the multiplier
	speedStageInt, err := strconv.Atoi(pokemon.SpeedStage)
	if err != nil {
		log.Fatalf("Invalid SpeedStage: %v", err)
	}
	stageMultiplier := stageMultipliers[speedStageInt]

	// Function to get the base speed stat from the API
	speedStat := func(name string) (string, error) {
		pokemonStruct, _, err := connections.PokemonApiCall("pokemon", name, connections.APIURL)
		if err != nil {
			return "", fmt.Errorf("API call failed: %w", err)
		}

		for _, stat := range pokemonStruct.Stats {
			if stat.Stat.Name == "speed" {
				return fmt.Sprintf("%d", stat.BaseStat), nil
			}
		}

		return "", errors.New("speed stat not found")
	}

	// Get the Pokémon's base speed
	speedStr, err := speedStat(pokemon.Name)
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

	// Define nature multipliers
	var natureMultipliers = map[string]float64{
		"+10%": 1.1,
		"0%":   1.0,
		"-10%": 0.9,
	}

	natureMultiplier := natureMultipliers[pokemon.Nature]

	// Format Pokémon name with title case
	chosenPokemon := cases.Title(language.English).String(pokemon.Name)

	// Calculate final speed using the formula:
	// (((2 x base + IV + (EV / 4)) x level / 100 + 5) * nature * modifier * ability * stage
	finalSpeed := float64(((2*baseSpeed+intIV+(intEV/4))*intLevel/100)+5) * natureMultiplier * abilityMultiplier * stageMultiplier * modifierMultiplier

	// Round down the final speed
	finalSpeedFloor := math.Floor(finalSpeed)
	finalSpeedStr := fmt.Sprintf("%.0f", finalSpeedFloor)

	// Format and return the output
	header := fmt.Sprintf("%s at level %s with selected options has a current speed of %s",
		styling.YellowAdaptive(chosenPokemon),
		styling.YellowAdaptive(pokemon.Level),
		styling.YellowAdaptive(finalSpeedStr),
	)
	body := fmt.Sprintf("EVs: %s\nIVs: %s\nModifiers: %s\nNature: %s\nAbility: %s\nSpeed Stage: %s\nBase Speed: %s",
		styling.YellowAdaptive(pokemon.SpeedEV),
		styling.YellowAdaptive(pokemon.SpeedIV),
		styling.YellowAdaptive(xstrings.EnglishJoin(pokemon.Modifier, true)),
		styling.YellowAdaptive(pokemon.Nature),
		styling.YellowAdaptive(pokemon.Ability),
		styling.YellowAdaptive(pokemon.SpeedStage),
		styling.YellowAdaptive(speedStr),
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
