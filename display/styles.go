package display

import "github.com/charmbracelet/lipgloss"

var (
	styleHeader      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleSuccess     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleFailed      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleWarning     = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleDim         = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleBorder      = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	styleEnvName     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	styleMeta        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleBannerTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleBannerSub   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleBannerBox   = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("238")).
				Padding(1, 2)
)

// valuePalette[0] = all-match color (green); rest are distinct mismatch colors
var valuePalette = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("10")),  // green - all match
	lipgloss.NewStyle().Foreground(lipgloss.Color("11")),  // yellow
	lipgloss.NewStyle().Foreground(lipgloss.Color("14")),  // cyan
	lipgloss.NewStyle().Foreground(lipgloss.Color("13")),  // magenta
	lipgloss.NewStyle().Foreground(lipgloss.Color("12")),  // blue
	lipgloss.NewStyle().Foreground(lipgloss.Color("208")), // orange
}

// ColorMap assigns a palette style to each unique non-empty value.
// Single unique value all get valuePalette[0] (green).
// Multiple unique values each gets a distinct mismatch color.
func ColorMap(values []string) map[string]lipgloss.Style {
	var unique []string
	seen := make(map[string]bool)
	for _, v := range values {
		if v != "" && !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	m := make(map[string]lipgloss.Style)
	if len(unique) == 1 {
		m[unique[0]] = valuePalette[0]
		return m
	}
	for i, v := range unique {
		idx := (i % (len(valuePalette) - 1)) + 1
		m[v] = valuePalette[idx]
	}
	return m
}

func statusSymbol(success bool) string {
	if success {
		return "✓"
	}
	return "✗"
}
