package card

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/styling"
)

func CardCommand() (string, error) {
	var output strings.Builder

	flag.Usage = func() {
		helpMessage := styling.HelpBorder.Render(
			"View data about cards from the TCG!\n\n",
			styling.StyleBold.Render("USAGE:"),
			fmt.Sprintf("\n\t%s %s %s", "poke-cli", styling.StyleBold.Render("card"), "[flag]"),
			"\n\n",
			styling.StyleBold.Render("FLAGS:"),
			fmt.Sprintf("\n\t%-30s %s", "-h, --help", "Prints out the help menu"),
		)
		output.WriteString(helpMessage)
	}

	if utils.CheckHelpFlag(&output, flag.Usage) {
		return output.String(), nil
	}

	flag.Parse()

	// Validate arguments
	if err := utils.ValidateCardArgs(os.Args); err != nil {
		output.WriteString(err.Error())
		return output.String(), err
	}

	seriesModel := SeriesList()
	// Program 1: Series selection
	finalModel, err := tea.NewProgram(seriesModel, tea.WithAltScreen()).Run()
	if err != nil {
		return "", fmt.Errorf("error running series selection program: %w", err)
	}

	result, ok := finalModel.(SeriesModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type from series selection: got %T, want SeriesModel", finalModel)
	}

	if result.SeriesID != "" {
		// Program 2: Sets selection
		setsModel, err := SetsList(result.SeriesID)

		if err != nil {
			return "", fmt.Errorf("error loading sets: %w", err)
		}

		finalSetsModel, err := tea.NewProgram(setsModel, tea.WithAltScreen()).Run()
		if err != nil {
			return "", fmt.Errorf("error running sets selection program: %w", err)
		}

		setsResult, ok := finalSetsModel.(SetsModel)
		if !ok {
			return "", fmt.Errorf("unexpected model type from sets selection: got %T, want SetsModel", finalSetsModel)
		}

		if setsResult.Quitting {
			return output.String(), nil
		}

		// Program 3: Cards display
		if setsResult.SetID != "" {
			cardsModel, err := CardsList(setsResult.SetID)
			if err != nil {
				return "", fmt.Errorf("error loading cards: %w", err)
			}

			for {
				finalCardsModel, err := tea.NewProgram(cardsModel, tea.WithAltScreen()).Run()
				if err != nil {
					return "", fmt.Errorf("error running cards program: %w", err)
				}

				cardsResult, ok := finalCardsModel.(CardsModel)
				if !ok {
					return "", fmt.Errorf("unexpected model type from cards display: got %T, want CardsModel", finalCardsModel)
				}

				if cardsResult.ViewImage {
					// Launch image viewer
					imageURL := cardsResult.ImageMap[cardsResult.SelectedOption]
					imageModel := ImageRenderer(cardsResult.SelectedOption, imageURL)
					_, err := tea.NewProgram(imageModel, tea.WithAltScreen()).Run()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Warning: image viewer error: %v\n", err)
					}

					// Re-launch cards with same state
					cardsResult.ViewImage = false
					cardsModel = cardsResult
				} else {
					break
				}
			}
		}
	}

	return output.String(), nil
}
