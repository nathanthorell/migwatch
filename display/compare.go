package display

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nathanthorell/migwatch/model"
)

func PrintCompareTable(results []model.EnvironmentResult) {
	if len(results) == 0 {
		fmt.Println(styleDim.Render("  no environments to compare"))
		return
	}

	headers := []string{"Migration"}
	for _, r := range results {
		headers = append(headers, r.Environment)
	}

	// latest versioned migration per env color-coded by version value
	latestValues := make([]string, len(results))
	for i, r := range results {
		if r.Error == nil {
			latestValues[i] = latestVersion(r.Migrations)
		}
	}
	latestColors := ColorMap(latestValues)

	latestRow := []string{"Latest Version"}
	for i, r := range results {
		if r.Error != nil {
			latestRow = append(latestRow, styleFailed.Render("error"))
			continue
		}
		v := latestValues[i]
		if v == "" {
			latestRow = append(latestRow, styleDim.Render("—"))
		} else {
			latestRow = append(latestRow, latestColors[v].Render(v))
		}
	}

	repeatableOrder, repeatableIndex := collectRepeatables(results)

	var rows [][]string
	rows = append(rows, latestRow)

	for _, desc := range repeatableOrder {
		// collect raw checksum strings per env position for color assignment
		rawValues := make([]string, len(results))
		type cellMeta struct {
			value   string
			failed  bool
			missing bool
		}
		cells := make([]cellMeta, len(results))

		for i, r := range results {
			if r.Error != nil {
				cells[i] = cellMeta{value: "error", failed: true}
				continue
			}
			idx, ok := repeatableIndex[r.Environment][desc]
			if !ok {
				cells[i] = cellMeta{missing: true}
				continue
			}
			m := r.Migrations[idx]
			cs := strconv.Itoa(m.Checksum)
			cells[i] = cellMeta{value: cs, failed: !m.Success}
			if !cells[i].failed {
				rawValues[i] = cs
			}
		}

		colors := ColorMap(rawValues)

		row := []string{desc}
		for i, c := range cells {
			_ = i
			switch {
			case c.failed:
				row = append(row, styleFailed.Render(c.value))
			case c.missing:
				row = append(row, styleDim.Render("—"))
			default:
				row = append(row, colors[c.value].Render(c.value))
			}
		}
		rows = append(rows, row)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(styleBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styleHeader
			}
			if col == 0 {
				return styleDim
			}
			return lipgloss.NewStyle()
		}).
		Headers(headers...).
		Rows(rows...)

	fmt.Println(t.Render())
	fmt.Println()
}

func latestVersion(migrations []model.Migration) string {
	var latest *model.Migration
	for i := range migrations {
		m := &migrations[i]
		if m.Version == "" || !m.Success {
			continue
		}
		if latest == nil || m.InstalledRank > latest.InstalledRank {
			latest = m
		}
	}
	if latest == nil {
		return ""
	}
	return "V" + latest.Version
}

func collectRepeatables(results []model.EnvironmentResult) ([]string, map[string]map[string]int) {
	seen := make(map[string]bool)
	var order []string
	index := make(map[string]map[string]int)

	for _, r := range results {
		index[r.Environment] = make(map[string]int)
		for i, m := range r.Migrations {
			if m.Version != "" {
				continue
			}
			index[r.Environment][m.Description] = i
			if !seen[m.Description] {
				seen[m.Description] = true
				order = append(order, m.Description)
			}
		}
	}

	return order, index
}
