package tui

import "strings"

func (m workspaceModel) renderWorkspaceHelp(width, height int) string {
	innerWidth := max(1, width-2)
	innerHeight := max(1, height-4)
	compact := innerHeight < 18 || innerWidth < 64
	lines := []string{m.styles.section.Render("cm ui · Quick help"), ""}
	if compact {
		lines = append(lines,
			"j/k move · h parent · l/enter directory · tab focus",
			"1-4 preview · / global path search · f filter · z full preview",
			"R reveals withheld sensitive content",
			"",
			m.workspaceStateLegend(),
			m.workspaceAttributeLegend(),
			"",
			"Esc, ?, or q closes help",
		)
	} else {
		lines = append(lines,
			m.styles.section.Render("Start here"),
			"1. j/k selects a direct child; l or Enter enters a directory; h returns.",
			"2. Tab moves to its preview. h/l then scroll the preview horizontally.",
			"3. / searches every inventory path; Enter locates the selected result.",
			"4. 1-4 selects diff, destination, target, or source; R reveals withheld content.",
			"",
			m.styles.section.Render("Directory browser"),
			"j/k move · h/left parent · l/right/enter directory · / search · f filter",
			"",
			m.styles.section.Render("Preview"),
			"j/k scroll · h/l horizontal · [/] hunks · z full screen",
			"",
			m.styles.section.Render("Legend"),
			m.workspaceStateLegend(),
			m.workspaceAttributeLegend(),
			"f file · d directory · l symlink · s script · x remove · e external",
			"",
			"Esc, ?, or q closes help",
		)
	}
	for i, line := range lines {
		lines[i] = truncate(line, innerWidth)
	}
	content := padLines(strings.Join(lines, "\n"), innerHeight)
	return m.styles.activePane.Width(width).Height(height).Render(content)
}

func (m workspaceModel) workspaceLegendParts() []string {
	return []string{
		m.styles.clean.Render("C") + " clean",
		m.styles.dirty.Render("D") + " dirty",
		m.styles.unmanaged.Render("U") + " unmanaged",
		m.styles.ignored.Render("I") + " ignored",
		m.styles.script.Render("R") + " script",
		m.styles.uninspected.Render("?") + " uninspected",
		m.styles.template.Render("[T]") + " template",
		m.styles.encrypted.Render("[E]") + " encrypted",
	}
}

func (m workspaceModel) workspaceStateLegend() string {
	return strings.Join(m.workspaceLegendParts()[:6], " · ")
}

func (m workspaceModel) workspaceAttributeLegend() string {
	return strings.Join(m.workspaceLegendParts()[6:], " · ")
}

func (m workspaceModel) workspaceLegend() string {
	return strings.Join(m.workspaceLegendParts(), " · ")
}

func (m workspaceModel) workspaceFooter() string {
	width := m.width
	if width <= 0 {
		width = 100
	}
	if m.helpVisible {
		return m.styles.help.Render(truncate("Esc/?/q close help", width))
	}
	if width < 60 {
		return m.styles.help.Render(truncate("tab focus · ? help · q quit", width))
	}
	var content string
	if m.focus == focusFiles {
		content = "j/k move · h parent · l/enter directory · / search · f filter · tab preview"
	} else {
		content = fmtPreviewFooter(m)
	}
	if width < 96 {
		if m.focus == focusFiles {
			content = "j/k move · h/l browse · / search · tab preview"
		} else {
			content = "j/k scroll · 1-4 view · / search · tab files"
		}
	}
	return m.styles.help.Render(truncate(content+" · ? help · q quit", width))
}

func fmtPreviewFooter(m workspaceModel) string {
	content := "j/k scroll · h/l horizontal · 1-4 view · z full · / search"
	if m.previewKind == "diff" {
		content += " · [/] hunks"
	}
	if m.previewQuery != "" {
		content += " · n/N matches"
	}
	if m.currentDiffState().preview.Withheld {
		content += " · R reveal"
	}
	return content
}
