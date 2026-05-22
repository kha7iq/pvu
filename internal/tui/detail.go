package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kha7iq/pvu/internal/ui"
)

var (
	detailSummaryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "8", Dark: "245"})
)

func (m Model) viewDetail() string {
	if m.loading {
		return "Loading pod details..."
	}

	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n" +
			helpStyle.Render("b back • q quit")
	}

	var b strings.Builder

	b.WriteString(ui.RenderPodHeader(m.detail.PodName, m.detail.Namespace))
	b.WriteString("\n\n")
	b.WriteString(detailSummaryStyle.Render(ui.RenderSummary(m.detail.Volumes)))
	b.WriteString("\n\n")

	if ui.UseCompactLayout(m.contentWidth()) {
		b.WriteString(ui.RenderCompactVolumeList(m.detail.Volumes))
	} else {
		b.WriteString(ui.RenderVolumeTable(m.detail.Volumes))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("b back • q quit"))

	return b.String()
}
