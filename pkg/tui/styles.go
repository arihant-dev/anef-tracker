package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Palette
	ColorBackground = lipgloss.Color("#1a1b26")
	ColorForeground = lipgloss.Color("#a9b1d6")
	ColorPrimary    = lipgloss.Color("#7aa2f7")
	ColorSecondary  = lipgloss.Color("#bb9af7")
	ColorSuccess    = lipgloss.Color("#9ece6a")
	ColorWarning    = lipgloss.Color("#e0af68")
	ColorDanger     = lipgloss.Color("#f7768e")
	ColorMuted      = lipgloss.Color("#565f89")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(ColorPrimary).
			Padding(0, 1)

	TabStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(ColorSecondary).
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	StatusBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1)

	BadgeLow = StatusBadgeStyle.
			Foreground(lipgloss.Color("#ffffff")).
			Background(ColorSuccess)

	BadgeMedium = StatusBadgeStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(ColorWarning)

	BadgeHigh = StatusBadgeStyle.
			Foreground(lipgloss.Color("#ffffff")).
			Background(ColorDanger)

	BadgeCritical = StatusBadgeStyle.
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff0055"))
)
