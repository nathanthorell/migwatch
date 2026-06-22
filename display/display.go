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

	if result.Database != "" || result.Schema != "" {
		fmt.Println(styleMeta.Render(
			fmt.Sprintf("  database: %s   schema: %s", result.Database, result.Schema),
		))
	}

	fmt.Println()
}
