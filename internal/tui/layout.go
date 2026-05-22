package tui

func (m Model) contentWidth() int {
	if m.width <= 0 {
		return 80
	}
	if m.width > 120 {
		return 120
	}
	return m.width - 2
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}
