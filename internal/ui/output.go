package ui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/kha7iq/pvu/internal/models"
)

func RenderPodVolumeView(view models.PodVolumeView) string {
	var b strings.Builder

	b.WriteString(RenderPodHeader(view.PodName, view.Namespace))
	b.WriteString("\n\n")
	b.WriteString(summaryStyle.Render(RenderSummary(view.Volumes)))
	b.WriteString("\n\n")

	if terminalWidth() < 90 {
		b.WriteString(RenderCompactVolumeList(view.Volumes))
	} else {
		b.WriteString(RenderVolumeTable(view.Volumes))
	}

	b.WriteString("\n")
	return b.String()
}

func RenderError(err error) string {
	return errorStyle.Render(fmt.Sprintf("Error: %v", err))
}

func terminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80
	}
	return width
}
