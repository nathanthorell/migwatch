package display

import "github.com/charmbracelet/lipgloss"

var (
	styleHeader      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleSuccess     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleFailed      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
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

func statusSymbol(success bool) string {
	if success {
		return "✓"
	}
	return "✗"
}
