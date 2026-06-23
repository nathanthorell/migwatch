package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/nathanthorell/migwatch/model"
)

func PrintBanner(summaries []model.EnvironmentSummary) {
	fmt.Print("\033[H\033[2J") // clears the screen

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	lines := []string{
		styleBannerTitle.Render("migwatch") + styleBannerSub.Render(" — database migration viewer"),
	}

	if len(summaries) > 0 {
		lines = append(lines, "")
		for _, s := range summaries {
			lines = append(lines, renderSummaryLine(s, termWidth))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	fmt.Println(styleBannerBox.Render(content))
	fmt.Println()
}

func renderSummaryLine(s model.EnvironmentSummary, termWidth int) string {
	envLabel := styleMeta.Render("Env: ") + styleBannerTitle.Render(s.Label)
	connParts := []string{
		"Driver: " + string(s.Driver),
		"Server: " + s.Host,
		"Database: " + s.Database,
	}
	connInfo := styleMeta.Render(strings.Join(connParts, " | "))

	singleLine := lipgloss.JoinHorizontal(lipgloss.Top, envLabel, "   ", connInfo)

	// box has border (1 each side) + padding (2 each side) = 6 chars overhead
	if lipgloss.Width(singleLine) <= termWidth-6 {
		return singleLine
	}

	return lipgloss.JoinVertical(lipgloss.Left, envLabel, "  "+connInfo)
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
