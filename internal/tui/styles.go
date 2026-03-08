package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary  = lipgloss.Color("39")  // bright blue
	colorMuted    = lipgloss.Color("240") // dim grey
	colorDone     = lipgloss.Color("241") // dimmer grey
	colorCritical = lipgloss.Color("196") // red
	colorHigh     = lipgloss.Color("214") // orange
	colorNormal   = lipgloss.Color("252") // near-white
	colorLow      = lipgloss.Color("240") // dim
	colorBorder   = lipgloss.Color("238")
	colorSelected = lipgloss.Color("39")
	colorHeader   = lipgloss.Color("39")
)

var (
	styleBase = lipgloss.NewStyle()

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	styleHeader = lipgloss.NewStyle().
			Foreground(colorHeader).
			Bold(true)

	styleStatusBar = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleKeys = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleKeyHint = lipgloss.NewStyle().
			Foreground(colorPrimary)

	styleSelected = lipgloss.NewStyle().
			Foreground(colorSelected).
			Bold(true)

	styleDone = lipgloss.NewStyle().
			Foreground(colorDone)

	styleTitle = lipgloss.NewStyle().
			Foreground(colorNormal).
			Bold(true)

	styleMuted = lipgloss.NewStyle().
			Foreground(colorMuted)

	styleLabel = lipgloss.NewStyle().
			Foreground(colorMuted).
			Width(10)

	styleValue = lipgloss.NewStyle().
			Foreground(colorNormal)

	styleDivider = lipgloss.NewStyle().
			Foreground(colorBorder)
)

// priorityStyle returns a style coloured for the given priority string.
func priorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "critical":
		return lipgloss.NewStyle().Foreground(colorCritical)
	case "high":
		return lipgloss.NewStyle().Foreground(colorHigh)
	case "low":
		return lipgloss.NewStyle().Foreground(colorLow)
	default:
		return lipgloss.NewStyle().Foreground(colorNormal)
	}
}
