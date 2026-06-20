package tui

import "strings"

func renderDiffLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	styled := make([]string, len(lines))
	for i, line := range lines {
		styled[i] = renderDiffLine(line)
	}
	return strings.Join(styled, "\n")
}

func renderDiffLine(line string) string {
	switch {
	case strings.HasPrefix(line, "@@"):
		return diffHunkStyle.Render(line)
	case strings.HasPrefix(line, "diff "), strings.HasPrefix(line, "---"), strings.HasPrefix(line, "+++"):
		return diffHeaderStyle.Render(line)
	case strings.HasPrefix(line, "+"):
		return diffAddStyle.Render(line)
	case strings.HasPrefix(line, "-"):
		return diffRemoveStyle.Render(line)
	case strings.HasPrefix(line, `\ No newline`):
		return diffMetaStyle.Render(line)
	default:
		return line
	}
}
