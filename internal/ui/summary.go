package ui

import (
	"fmt"
	"strings"

	"github.com/kha7iq/pvu/internal/models"
)

func RenderSummary(volumes []models.VolumeRow) string {
	total := len(volumes)
	warn := 0
	critical := 0

	for _, v := range volumes {
		switch {
		case v.UsagePercent >= 90:
			critical++
		case v.UsagePercent >= 75:
			warn++
		}
	}

	parts := []string{fmt.Sprintf("%d PVCs", total)}
	if warn > 0 {
		parts = append(parts, fmt.Sprintf("%d warning", warn))
	}
	if critical > 0 {
		parts = append(parts, fmt.Sprintf("%d critical", critical))
	}

	return strings.Join(parts, "  •  ")
}
