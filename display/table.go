package display

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nathanthorell/migwatch/model"
)

func PrintMigrationTable(result model.EnvironmentResult) {
	PrintEnvironmentHeader(result)

	if result.Error != nil {
		fmt.Println(styleFailed.Render("  error: " + result.Error.Error()))
		fmt.Println()
		return
	}

	if len(result.Migrations) == 0 {
		fmt.Println(styleDim.Render("  no migrations found"))
		fmt.Println()
		return
	}

	rows := make([][]string, len(result.Migrations))
	for i, m := range result.Migrations {
		rows[i] = []string{
			strconv.Itoa(m.InstalledRank),
			m.Version,
			m.Description,
			m.InstalledOn.Format(time.DateTime),
			strconv.Itoa(m.ExecutionTime) + "ms",
			statusSymbol(m.Success),
		}
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(styleBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styleHeader
			}
			if col == 5 {
				if rows[row][col] == "✓" {
					return styleSuccess
				}
				return styleFailed
			}
			return lipgloss.NewStyle()
		}).
		Headers("Rank", "Version", "Description", "Installed On", "Duration", "Status").
		Rows(rows...)

	fmt.Println(t.Render())
	fmt.Println()
}
