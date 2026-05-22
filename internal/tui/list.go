package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	listTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "33", Dark: "39"})

	listHeaderBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.AdaptiveColor{Light: "7", Dark: "240"}).
				Padding(0, 1)

	listSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "8", Dark: "245"})

	listBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "255", Dark: "230"}).
			Background(lipgloss.AdaptiveColor{Light: "33", Dark: "62"}).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "255", Dark: "230"}).
			Background(lipgloss.AdaptiveColor{Light: "33", Dark: "62"}).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "249"}).
			Padding(0, 1)

	filterBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "8", Dark: "240"}).
			Padding(0, 1)

	filterActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.AdaptiveColor{Light: "33", Dark: "62"}).
				Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "8", Dark: "243"})

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "160", Dark: "203"})
)

func (m Model) View() string {
	switch m.screen {
	case screenDetail:
		return m.viewDetail()
	default:
		return m.viewList()
	}
}

func (m Model) viewList() string {
	if m.loading {
		return listTitleStyle.Render("pvu") + "\n\nLoading pods..."
	}

	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	var b strings.Builder

	b.WriteString(m.renderListHeader())
	b.WriteString("\n\n")
	b.WriteString(m.renderFilterBar())
	b.WriteString("\n\n")

	if len(m.filteredPods) == 0 {
		if len(m.pods) == 0 {
			b.WriteString(helpStyle.Render("No pods found in this namespace."))
		} else {
			b.WriteString(helpStyle.Render("No matching pods."))
		}
		b.WriteString("\n\n")
		b.WriteString(m.renderListFooter())
		return b.String()
	}

	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")
	b.WriteString(m.renderListFooter())

	return b.String()
}

func (m Model) renderListHeader() string {
	countLabel := fmt.Sprintf("%d pods", len(m.filteredPods))
	if m.filterInput.Value() != "" {
		countLabel = fmt.Sprintf("%d/%d pods", len(m.filteredPods), len(m.pods))
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		listTitleStyle.Render("pvu"),
		"  ",
		listSubtitleStyle.Render(fmt.Sprintf("namespace: %s", m.namespace)),
		"  ",
		listBadgeStyle.Render(countLabel),
	)

	return listHeaderBoxStyle.Width(m.contentWidth()).Render(header)
}

func (m Model) renderFilterBar() string {
	if m.filtering {
		return filterActiveStyle.Width(m.contentWidth()).Render(m.filterInput.View())
	}

	if m.filterInput.Value() != "" {
		return filterBoxStyle.Width(m.contentWidth()).Render(m.filterInput.View())
	}

	return filterBoxStyle.Width(m.contentWidth()).Render(
		helpStyle.Render("/ filter pods"),
	)
}

func (m Model) renderListContent(width int) string {
	if len(m.filteredPods) == 0 {
		return ""
	}

	var b strings.Builder

	rowWidth := width - 4
	if rowWidth < 10 {
		rowWidth = 10
	}

	for i, pod := range m.filteredPods {
		line := truncate(pod.Name, rowWidth)
		if i == m.selected {
			b.WriteString(selectedStyle.Width(width).Render("› " + line))
		} else {
			b.WriteString(normalStyle.Width(width).Render("  " + line))
		}
		if i < len(m.filteredPods)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m Model) renderListFooter() string {
	if m.filtering {
		return helpStyle.Render("type to filter • ↑/k up • ↓/j down • enter open • esc clear/back • q quit")
	}

	if m.filterInput.Value() != "" {
		return helpStyle.Render("/ filter • ↑/k up • ↓/j down • enter open • esc clear/back • q quit")
	}

	return helpStyle.Render("/ filter • ↑/k up • ↓/j down • enter open • q quit")
}
