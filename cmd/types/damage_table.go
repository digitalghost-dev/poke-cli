package types

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/styling"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

// DamageTable Function to display type details after a type is selected
func DamageTable(typesName string, endpoint string) {
	// Setting up variables to style the list
	var columnWidth = 11
	var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	var list = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(subtle).MarginRight(2).Height(8)
	var listHeader = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(subtle).MarginRight(2).Render
	var listItem = lipgloss.NewStyle().Render
	var docStyle = lipgloss.NewStyle().Padding(1, 1, 1, 1)

	typesStruct, typeName, _ := connections.TypesApiCall(endpoint, typesName, connections.APIURL)

	// Format selected type
	selectedType := cases.Title(language.English).String(typeName)
	coloredType := lipgloss.NewStyle().Foreground(lipgloss.Color(styling.GetTypeColor(typeName))).Render(selectedType)

	fmt.Printf("You selected the %s type.\nNumber of Pokémon with type: %d\nNumber of moves with type: %d\n", coloredType, len(typesStruct.Pokemon), len(typesStruct.Moves))
	fmt.Println("----------")
	fmt.Println(styling.StyleBold.Render("Damage Chart:"))

	physicalWidth, _, _ := term.GetSize(uintptr(int(os.Stdout.Fd())))
	doc := strings.Builder{}

	// Helper function to build list items
	buildListItems := func(items []struct{ Name, URL string }) string {
		var itemList []string
		for _, item := range items {
			color := styling.GetTypeColor(item.Name)
			coloredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			coloredItem := coloredStyle.Render(cases.Title(language.English).String(item.Name))
			itemList = append(itemList, listItem(coloredItem))
		}
		return lipgloss.JoinVertical(lipgloss.Left, itemList...)
	}

	// Render lists based on Damage Relations
	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Weakness"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.DoubleDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x2 Dmg"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.DoubleDamageTo)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Resists"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.HalfDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x0.5 Dmg"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.HalfDamageTo)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Immune"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.NoDamageFrom)),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("x0 Dmg"),
				buildListItems([]struct{ Name, URL string }(typesStruct.DamageRelations.NoDamageTo)),
			),
		),
	)

	// Append lists to document
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	// Print the rendered document
	fmt.Println(docStyle.Render(doc.String()))
}
