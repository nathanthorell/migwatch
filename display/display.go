package display

import (
	"fmt"

	"github.com/nathanthorell/migwatch/model"
)

func PrintBanner() {
	fmt.Print("\033[H\033[2J") // clears the screen
	fmt.Println(styleBanner.Render("migwatch") + styleDim.Render(" — database migration viewer"))
	fmt.Println()
}

func PrintEnvironmentHeader(result model.EnvironmentResult) {
	fmt.Println(styleEnvName.Render(result.Environment))

	if result.Database != "" {
		fmt.Println(styleMeta.Render("  database: " + result.Database))
	}

	fmt.Println()
}

func PrintSchemaLabel(schema string) {
	fmt.Println(styleDim.Render("  schema: " + schema))
}
