package card

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
)

func CardCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		output.WriteString(
			utils.GenerateHelpMessage(
				utils.HelpConfig{
					Description: "Get details about a specific card.",
					CmdName:     "card",
				},
			),
		)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	// Validate arguments
	if err := utils.ValidateArgs(os.Args, utils.Validator{MaxArgs: 3, CmdName: "card", RequireName: false, HasFlags: false}); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	// Program 1: Series selection
	finalModel, err := tea.NewProgram(SeriesList(), tea.WithAltScreen()).Run()
	if err != nil {
		return "", fmt.Errorf("error running series selection program: %w", err)
	}

	result, ok := finalModel.(seriesModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type from series selection: got %T, want seriesModel", finalModel)
	}

	if result.SeriesID != "" {
		// Program 2: Sets selection
		setsMdl, err := SetsList(result.SeriesID)
		if err != nil {
			return "", fmt.Errorf("error loading sets: %w", err)
		}

		finalSetsModel, err := tea.NewProgram(setsMdl, tea.WithAltScreen()).Run()
		if err != nil {
			return "", fmt.Errorf("error running sets selection program: %w", err)
		}

		setsResult, ok := finalSetsModel.(setsModel)
		if !ok {
			return "", fmt.Errorf("unexpected model type from sets selection: got %T, want setsModel", finalSetsModel)
		}

		if setsResult.Quitting {
			return output.String(), nil
		}

		// Program 3: Cards display
		if setsResult.SetID != "" {
			cardsMdl, err := CardsList(setsResult.SetID)
			if err != nil {
				return "", fmt.Errorf("error loading cards: %w", err)
			}

			for {
				finalCardsModel, err := tea.NewProgram(cardsMdl, tea.WithAltScreen()).Run()
				if err != nil {
					return "", fmt.Errorf("error running cards program: %w", err)
				}

				cardsResult, ok := finalCardsModel.(cardsModel)
				if !ok {
					return "", fmt.Errorf("unexpected model type from cards display: got %T, want cardsModel", finalCardsModel)
				}

				if cardsResult.ViewImage {
					// Launch image viewer
					imageURL := cardsResult.ImageMap[cardsResult.SelectedOption]
					_, err := tea.NewProgram(ImageRenderer(cardsResult.SelectedOption, imageURL), tea.WithAltScreen()).Run()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Warning: image viewer error: %v\n", err)
					}

					// Re-launch cards with same state
					cardsResult.ViewImage = false
					cardsMdl = cardsResult
				} else {
					break
				}
			}
		}
	}

	return output.String(), nil
}
