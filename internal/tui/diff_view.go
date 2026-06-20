package tui

import (
	"strings"

	"github.com/zhongyangchuwu/cm/internal/report"
)

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
	switch report.ClassifyDiffLine(line).Kind {
	case report.DiffLineHunk:
		return diffHunkStyle.Render(line)
	case report.DiffLineHeader:
		return diffHeaderStyle.Render(line)
	case report.DiffLineAdd:
		return diffAddStyle.Render(line)
	case report.DiffLineRemove:
		return diffRemoveStyle.Render(line)
	case report.DiffLineMeta:
		return diffMetaStyle.Render(line)
	default:
		return line
	}
}
