package tui

import (
	"fmt"
	"path"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type rect struct {
	width  int
	height int
}

func (m workspaceModel) View() tea.View {
	return tea.NewView(m.viewString())
}

func (m workspaceModel) viewString() string {
	if !m.isWorkspace() && len(m.entries) == 0 {
		if m.completed && m.message != "" {
			return m.message + "\n"
		}
		return "clean\n"
	}

	width := m.width
	if width <= 0 {
		width = 100
	}
	height := m.height
	if height <= 0 {
		height = 30
	}

	header := titleStyle.Render(truncate(m.title(), width))
	footer := m.footer()
	status := m.statusLine()
	messageHeight := 0
	if status != "" {
		messageHeight = lipgloss.Height(status)
	}
	bodyHeight := height - lipgloss.Height(header) - lipgloss.Height(footer) - messageHeight - 2
	if bodyHeight < 6 {
		bodyHeight = 6
	}

	var body string
	switch {
	case m.isWorkspace() && m.previewFull:
		body = m.renderMainPane(rect{width: width, height: bodyHeight})
	case m.isWorkspace() && width < 60:
		if m.focus == focusFiles {
			body = m.renderFilesPane(rect{width: width, height: bodyHeight})
		} else {
			body = m.renderMainPane(rect{width: width, height: bodyHeight})
		}
	default:
		leftWidth := clamp(width/3, 28, 48)
		if leftWidth > width-24 {
			leftWidth = width / 2
		}
		rightWidth := width - leftWidth - 1
		if rightWidth < 20 {
			rightWidth = 20
		}
		files := m.renderFilesPane(rect{width: leftWidth, height: bodyHeight})
		main := m.renderMainPane(rect{width: rightWidth, height: bodyHeight})
		body = lipgloss.JoinHorizontal(lipgloss.Top, files, main)
	}

	parts := []string{header, body}
	if status != "" {
		parts = append(parts, status)
	}
	parts = append(parts, footer)
	return strings.Join(parts, "\n") + "\n"
}

func (m workspaceModel) title() string {
	if !m.isWorkspace() {
		return "cm sync"
	}
	view := "tree"
	if m.flat {
		view = "flat"
	}
	return fmt.Sprintf("cm ui • %d/%d • %s • %s", len(m.entries), len(m.allEntries), m.filterLabel(), view)
}

func (m workspaceModel) statusLine() string {
	if m.search != searchNone {
		kind := "files"
		if m.search == searchPreview {
			kind = "preview"
		}
		return fmt.Sprintf("search %s: %s_", kind, m.searchInput)
	}
	if m.message != "" {
		return m.message
	}
	if m.isWorkspace() {
		parts := make([]string, 0, 2)
		if m.snapshot.Notice != "" {
			parts = append(parts, m.snapshot.Notice)
		}
		if m.fileQuery != "" {
			parts = append(parts, "path filter: "+m.fileQuery)
		}
		return strings.Join(parts, " • ")
	}
	return ""
}

func (m workspaceModel) renderFilesPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	lines := make([]string, 0, innerHeight)
	for _, entry := range visibleEntries(m.entries, m.cursor, innerHeight) {
		cursor := "  "
		if entry.index == m.cursor {
			cursor = "> "
		}
		line := cursor + m.entryLabel(entry.entry)
		lines = append(lines, truncate(line, innerWidth))
	}
	if len(lines) == 0 {
		lines = append(lines, "no entries")
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}

	style := paneStyle
	if m.focus == focusFiles && m.mode == modeReview {
		style = activePaneStyle
	}
	return style.Width(size.width).Height(size.height).Render(strings.Join(lines, "\n"))
}

func (m workspaceModel) entryLabel(entry app.WorkspaceEntry) string {
	pathLabel := entry.RelativePath
	if pathLabel == "" || !m.isWorkspace() {
		pathLabel = m.displayPath(entry.Path)
	}
	if !m.isWorkspace() {
		return m.pendingLabel(entry.Path) + " " + pathLabel
	}
	indent := ""
	treeMarker := ""
	if !m.flat {
		indent = strings.Repeat("  ", treeDepth(entry.RelativePath))
		if entry.Type == app.TargetDirectory {
			treeMarker = "▾ "
			if m.collapsed[entry.RelativePath] {
				treeMarker = "▸ "
			}
		}
	}
	displayName := pathLabel
	if !m.flat {
		displayName = path.Base(pathLabel)
	}
	return fmt.Sprintf("%s%s%s:%s %s%s", indent, treeMarker, stateMarker(entry.State), typeMarker(entry), displayName, attributeMarker(entry))
}

func stateMarker(state app.FileState) string {
	switch state {
	case app.FileDirty:
		return "D"
	case app.FileUnmanaged:
		return "U"
	case app.FileIgnored:
		return "I"
	case app.FileScript:
		return "R"
	case app.FileUninspected:
		return "?"
	default:
		return "C"
	}
}

func typeMarker(entry app.WorkspaceEntry) string {
	switch entry.Type {
	case app.TargetDirectory:
		return "d"
	case app.TargetSymlink:
		return "l"
	case app.TargetScript:
		return "s"
	case app.TargetRemove:
		return "x"
	case app.TargetExternal:
		return "e"
	case app.TargetFile:
		return "f"
	default:
		return "?"
	}
}

func attributeMarker(entry app.WorkspaceEntry) string {
	var marker strings.Builder
	if entry.Template {
		marker.WriteString(" [T]")
	}
	if entry.Encrypted {
		marker.WriteString(" [E]")
	}
	return marker.String()
}

func (m workspaceModel) renderMainPane(size rect) string {
	if m.mode == modeConfirm || m.mode == modeExecuting {
		return m.renderConfirmPane(size)
	}
	return m.renderDiffPane(size)
}

func (m workspaceModel) renderDiffPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	state := m.currentDiffState()
	var content string
	switch {
	case len(m.entries) == 0:
		content = "no entries match the current view"
	case state.loading:
		content = "loading preview..."
	case state.err != nil:
		content = state.err.Error()
	case len(state.lines) == 0:
		content = "no preview"
	default:
		lines := state.lines
		if m.diffScroll > len(lines) {
			lines = nil
		} else {
			lines = lines[m.diffScroll:]
		}
		if len(lines) > innerHeight {
			lines = lines[:innerHeight]
		}
		visible := make([]string, len(lines))
		for i, line := range lines {
			visible[i] = cropLine(line, m.previewX, innerWidth)
		}
		content = renderDiffLines(visible)
	}

	content = padLines(content, innerHeight)
	style := paneStyle
	if m.focus == focusDiff && m.mode == modeReview {
		style = activePaneStyle
	}
	return style.Width(size.width).Height(size.height).Render(content)
}

func (m workspaceModel) renderConfirmPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	title := "Confirm actions"
	if m.mode == modeExecuting {
		title = "Executing actions"
	}
	lines := []string{sectionStyle.Render(title), ""}
	if m.mode == modeExecuting {
		if m.executingIndex < len(m.executing) {
			action := m.executing[m.executingIndex]
			lines = append(lines, truncate(fmt.Sprintf("%d/%d %s %s", m.executingIndex+1, len(m.executing), action.Kind.String(), m.displayPath(action.Target)), innerWidth))
		}
		lines = append(lines, "", fmt.Sprintf("executed: %d", m.executedCount), fmt.Sprintf("skipped: %d", m.skippedCount), fmt.Sprintf("deferred: %d", m.deferredCount))
	} else {
		for _, action := range m.pendingActions() {
			lines = append(lines, truncate(fmt.Sprintf("%s %s", action.Kind.String(), m.displayPath(action.Target)), innerWidth))
		}
		if len(lines) == 2 {
			lines = append(lines, "no pending actions")
		}
	}
	content := strings.Join(lines, "\n")
	content = padLines(content, innerHeight)
	return activePaneStyle.Width(size.width).Height(size.height).Render(content)
}

func (m workspaceModel) footer() string {
	if m.isWorkspace() {
		content := "files • tab preview • j/k move • t view • f filter • / search"
		if m.focus == focusDiff {
			content = fmt.Sprintf("preview • tab files • j/k scroll • h/l x=%d • 1-4 view • z full • / search", m.previewX)
			if m.previewKind == app.PreviewDiff {
				content += " • [/] hunks"
			}
			if m.previewQuery != "" {
				content += " • n/N matches"
			}
			if m.currentDiffState().preview.Withheld {
				content += " • R reveal"
			}
		}
		content += " • q quit"
		width := m.width
		if width <= 0 {
			width = 100
		}
		return helpStyle.Render(truncate(content, width))
	}
	var content string
	if m.mode == modeExecuting {
		content = "executing..."
	} else if m.mode == modeConfirm {
		content = m.help.ShortHelpView(defaultSyncKeys.confirmHelp())
	} else {
		focus := "files"
		if m.focus == focusDiff {
			focus = "diff"
		}
		content = focus + " • " + m.help.ShortHelpView(defaultSyncKeys.reviewHelp(m.focus, m.currentDiffState().review))
	}
	if m.scriptCount > 0 {
		content += " • " + scriptCount(m.scriptCount) + " outside cm sync"
	}
	return helpStyle.Render(content)
}

type visibleEntry struct {
	index int
	entry app.WorkspaceEntry
}

func visibleEntries(entries []app.WorkspaceEntry, cursor int, height int) []visibleEntry {
	if height <= 0 || len(entries) == 0 {
		return nil
	}
	if height > len(entries) {
		height = len(entries)
	}
	start := cursor - height/2
	if start < 0 {
		start = 0
	}
	if start+height > len(entries) {
		start = len(entries) - height
	}
	visible := make([]visibleEntry, 0, height)
	for i := start; i < start+height; i++ {
		visible = append(visible, visibleEntry{index: i, entry: entries[i]})
	}
	return visible
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return ansi.Truncate(s, width, "…")
}

func cropLine(s string, offset int, width int) string {
	if width <= 0 {
		return ""
	}
	if offset > 0 {
		s = ansi.TruncateLeft(s, offset, "")
	}
	return ansi.Truncate(s, width, "")
}

func padLines(content string, height int) string {
	lines := splitLines(content)
	for len(lines) < height {
		lines = append(lines, "")
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	return strings.Join(lines, "\n")
}
