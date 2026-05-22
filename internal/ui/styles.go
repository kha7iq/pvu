package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorAccent = lipgloss.AdaptiveColor{Light: "33", Dark: "39"}

	colorMuted  = lipgloss.AdaptiveColor{Light: "8", Dark: "245"}
	colorBorder = lipgloss.AdaptiveColor{Light: "7", Dark: "240"}

	colorOK       = lipgloss.AdaptiveColor{Light: "28", Dark: "77"}
	colorWarning  = lipgloss.AdaptiveColor{Light: "130", Dark: "214"}
	colorCritical = lipgloss.AdaptiveColor{Light: "160", Dark: "203"}
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Foreground(colorBorder)

	evenRowStyle = lipgloss.NewStyle()

	oddRowStyle = lipgloss.NewStyle()

	usageNormalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorOK)

	usageWarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning)

	usageCriticalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorCritical)

	emptyStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(colorMuted)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCritical)
)
