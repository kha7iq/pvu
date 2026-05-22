package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/kha7iq/pvu/internal/models"
)

func RenderVolumeTable(volumes []models.VolumeRow) string {
	if len(volumes) == 0 {
		return emptyStyle.Render("No PersistentVolumeClaims (PVCs) attached to this pod.")
	}

	sorted := make([]models.VolumeRow, len(volumes))
	copy(sorted, volumes)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].UsagePercent > sorted[j].UsagePercent
	})

	rows := make([][]string, 0, len(sorted))
	for _, vol := range sorted {
		rows = append(rows, []string{
			vol.VolumeName,
			vol.PVCName,
			formatBytes(vol.CapacityBytes),
			formatBytes(vol.UsedBytes),
			formatBytes(vol.AvailableBytes),
			renderUsageCell(vol.UsagePercent),
		})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(borderStyle).
		Headers("POD VOLUME", "PVC", "CAPACITY", "USED", "AVAILABLE", "USAGE").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerCellStyle(col)
			}
			return bodyCellStyle(row, col)
		})

	return t.String()
}

func renderUsageCell(usage float64) string {
	value := fmt.Sprintf("%.2f%%", usage)

	switch {
	case usage >= 90:
		return usageCriticalStyle.Width(7).Align(lipgloss.Right).Render(value)
	case usage >= 75:
		return usageWarnStyle.Width(7).Align(lipgloss.Right).Render(value)
	default:
		return usageNormalStyle.Width(7).Align(lipgloss.Right).Render(value)
	}
}

func headerCellStyle(col int) lipgloss.Style {
	style := headerStyle

	switch col {
	case 0:
		return style.Width(12)
	case 1:
		return style.Width(25)
	case 2:
		return style.Width(10).Align(lipgloss.Right)
	case 3:
		return style.Width(10).Align(lipgloss.Right)
	case 4:
		return style.Width(10).Align(lipgloss.Right)
	case 5:
		return style.Width(7).Align(lipgloss.Right)
	default:
		return style
	}
}

func bodyCellStyle(row, col int) lipgloss.Style {
	var style lipgloss.Style
	if row%2 == 0 {
		style = evenRowStyle
	} else {
		style = oddRowStyle
	}

	switch col {
	case 0:
		return style.Width(12)
	case 1:
		return style.Width(25)
	case 2:
		return style.Width(10).Align(lipgloss.Right)
	case 3:
		return style.Width(10).Align(lipgloss.Right)
	case 4:
		return style.Width(10).Align(lipgloss.Right)
	case 5:
		return style.Width(7).Align(lipgloss.Right)
	default:
		return style
	}
}

func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
		PB = 1024 * TB
	)

	switch {
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/float64(PB))
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
