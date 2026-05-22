package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	valueStyle = lipgloss.NewStyle().
			Bold(true)

	summaryStyle = lipgloss.NewStyle().
			Foreground(colorMuted)
)

func RenderPodHeader(podName, namespace string) string {
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		titleStyle.Render("pvu"),
		"  ",
		labelStyle.Render("pod: "),
		valueStyle.Render(podName),
		"   ",
		labelStyle.Render("namespace: "),
		valueStyle.Render(namespace),
	)

	return headerBoxStyle.Render(header)
}
