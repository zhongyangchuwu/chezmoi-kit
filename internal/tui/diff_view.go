package tui

import (
	"strings"

	"github.com/zhongyangchuwu/cm/internal/report"
)

func renderPreviewLines(m workspaceModel, lines []string, state diffState) string {
	if len(lines) == 0 {
		return ""
	}
	styled := make([]string, len(lines))
	for i, line := range lines {
		styled[i] = m.renderPreviewLine(line, state)
	}
	return strings.Join(styled, "\n")
}

func (m workspaceModel) renderPreviewLine(line string, state diffState) string {
	if m.isWorkspace() {
		switch {
		case strings.HasPrefix(line, "view:"):
			return m.styles.section.Render(line)
		case strings.HasPrefix(line, "path:") || strings.HasPrefix(line, "source:"):
			return m.styles.muted.Render(line)
		case state.preview.Withheld && line == state.preview.Notice:
			return m.styles.withheld.Render(line)
		case line == "destination matches rendered target":
			return m.styles.clean.Render(line)
		}
	}
	return renderDiffLine(line, m.styles)
}

func renderDiffLine(line string, styles tuiStyles) string {
	switch report.ClassifyDiffLine(line).Kind {
	case report.DiffLineHunk:
		return styles.diffHunk.Render(line)
	case report.DiffLineHeader:
		return styles.diffHeader.Render(line)
	case report.DiffLineAdd:
		return styles.diffAdd.Render(line)
	case report.DiffLineRemove:
		return styles.diffRemove.Render(line)
	case report.DiffLineMeta:
		return styles.diffMeta.Render(line)
	default:
		return line
	}
}
