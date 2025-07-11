package speed

import (
	"errors"
	"github.com/charmbracelet/huh"
	"github.com/digitalghost-dev/poke-cli/connections"
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
		//pokemonTwo PokemonTwoDetails
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

	var output strings.Builder

	return output.String(), nil
}
