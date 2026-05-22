package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/kha7iq/pvu/internal/models"
)

var (
	compactBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	compactTitleStyle = lipgloss.NewStyle().
				Bold(true)

	compactLabelStyle = lipgloss.NewStyle().
				Foreground(colorMuted)

	compactValueStyle = lipgloss.NewStyle()
)

func RenderCompactVolumeList(volumes []models.VolumeRow) string {
	if len(volumes) == 0 {
		return emptyStyle.Render("No PersistentVolumeClaims (PVCs) attached to this pod.")
	}

	sorted := make([]models.VolumeRow, len(volumes))
	copy(sorted, volumes)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].UsagePercent > sorted[j].UsagePercent
	})

	blocks := make([]string, 0, len(sorted))
	for _, vol := range sorted {
		blocks = append(blocks, renderCompactVolume(vol))
	}

	return strings.Join(blocks, "\n\n")
}

func renderCompactVolume(vol models.VolumeRow) string {
	usage := fmt.Sprintf("%.2f%%", vol.UsagePercent)

	usageRendered := usageNormalStyle.Render(usage)
	switch {
	case vol.UsagePercent >= 90:
		usageRendered = usageCriticalStyle.Render(usage)
	case vol.UsagePercent >= 75:
		usageRendered = usageWarnStyle.Render(usage)
	}

	body := strings.Join([]string{
		compactTitleStyle.Render(vol.VolumeName),
		renderCompactRow("PVC", vol.PVCName),
		renderCompactRow("Capacity", formatBytes(vol.CapacityBytes)),
		renderCompactRow("Used", formatBytes(vol.UsedBytes)),
		renderCompactRow("Available", formatBytes(vol.AvailableBytes)),
		renderCompactRow("Usage", usageRendered),
	}, "\n")

	return compactBoxStyle.Render(body)
}

func renderCompactRow(label, value string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		compactLabelStyle.Width(10).Render(label),
		compactValueStyle.Render(value),
	)
}
