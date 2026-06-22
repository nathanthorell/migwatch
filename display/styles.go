package display

import "github.com/charmbracelet/lipgloss"

var (
	styleHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleFailed  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleBorder  = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	styleEnvName = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleMeta    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleBanner  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
)

func statusSymbol(success bool) string {
	if success {
		return "✓"
	}
	return "✗"
}
