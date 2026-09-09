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
	if m.isWorkspace() && m.helpVisible {
		return m.renderWorkspaceHelp(width, height) + "\n"
	}

	header := m.styles.title.Render(truncate(m.title(), width))
	footer := m.footer()
	status := m.renderStatusLine()
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
	case m.isWorkspace() && width < 100:
		filesWidth := clamp(width*2/5, 28, 42)
		files := m.renderFilesPane(rect{width: filesWidth, height: bodyHeight})
		main := m.renderMainPane(rect{width: width - filesWidth - 1, height: bodyHeight})
		body = lipgloss.JoinHorizontal(lipgloss.Top, files, main)
	case m.isWorkspace():
		parentWidth := clamp(width/5, 20, 30)
		filesWidth := clamp(width/3, 30, 44)
		parent := m.renderParentPane(rect{width: parentWidth, height: bodyHeight})
		files := m.renderFilesPane(rect{width: filesWidth, height: bodyHeight})
		main := m.renderMainPane(rect{width: width - parentWidth - filesWidth - 2, height: bodyHeight})
		body = lipgloss.JoinHorizontal(lipgloss.Top, parent, files, main)
	default:
		leftWidth := clamp(width/3, 28, 48)
		if leftWidth > width-24 {
			leftWidth = width / 2
		}
		files := m.renderFilesPane(rect{width: leftWidth, height: bodyHeight})
		main := m.renderMainPane(rect{width: width - leftWidth - 1, height: bodyHeight})
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
	view := m.directoryLabel()
	if m.searchResults {
		view = "search"
	}
	return fmt.Sprintf("cm ui • %d/%d • %s • %s", len(m.entries), len(m.workspaceNodes), m.filterLabel(), view)
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
	if m.workspaceNotice != "" {
		return m.workspaceNotice
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

func (m workspaceModel) renderStatusLine() string {
	status := m.statusLine()
	if status == "" {
		return ""
	}
	if m.search != searchNone {
		return m.styles.loading.Render(status)
	}
	state := m.currentDiffState()
	if state.preview.Withheld && status == state.preview.Notice {
		return m.styles.withheld.Render(status)
	}
	if strings.HasPrefix(status, "error:") {
		return m.styles.error.Render(status)
	}
	return m.styles.muted.Render(status)
}
func (m workspaceModel) renderFilesPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	legendHeight := 0
	if m.isWorkspace() && innerHeight >= 8 {
		legendHeight = 1
	}
	listHeight := max(1, innerHeight-1-legendHeight)
	lines := []string{m.styles.section.Render(truncate(m.filesPaneTitle(), innerWidth))}
	for _, entry := range visibleEntries(m.entries, m.cursor, listHeight) {
		cursor := "  "
		if entry.index == m.cursor {
			cursor = "> "
		}
		line := cursor + m.entryLabel(entry.entry)
		if entry.index == m.cursor {
			if m.focus == focusFiles && m.mode == modeReview {
				line = m.styles.selectedActive.Render(line)
			} else {
				line = m.styles.selected.Render(line)
			}
		}
		lines = append(lines, truncate(line, innerWidth))
	}
	if len(lines) == 1 {
		lines = append(lines, m.styles.muted.Render("no entries"))
	}
	if legendHeight > 0 {
		lines = append(lines, truncate(m.workspaceLegend(), innerWidth))
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}

	style := m.styles.pane
	if m.focus == focusFiles && m.mode == modeReview {
		style = m.styles.activePane
	}
	return style.Width(size.width).Height(size.height).Render(strings.Join(lines, "\n"))
}

func (m workspaceModel) renderParentPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	lines := []string{m.styles.section.Render(truncate("Parent", innerWidth))}
	entries := m.parentEntries()
	cursor := 0
	for index, entry := range entries {
		if entry.RelativePath == m.currentDir {
			cursor = index
			break
		}
	}
	for _, entry := range visibleEntries(entries, cursor, max(1, innerHeight-1)) {
		line := "  " + m.entryLabel(entry.entry)
		if entry.entry.RelativePath == m.currentDir {
			line = m.styles.selected.Render("· " + m.entryLabel(entry.entry))
		}
		lines = append(lines, truncate(line, innerWidth))
	}
	if len(lines) == 1 {
		lines = append(lines, m.styles.muted.Render("root"))
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}
	return m.styles.pane.Width(size.width).Height(size.height).Render(strings.Join(lines, "\n"))
}

func (m workspaceModel) entryLabel(entry app.WorkspaceEntry) string {
	pathLabel := entry.RelativePath
	if pathLabel == "" || !m.isWorkspace() {
		pathLabel = m.displayPath(entry.Path)
	}
	if !m.isWorkspace() {
		return m.pendingLabel(entry.Path) + " " + pathLabel
	}
	displayName := path.Base(pathLabel)
	if entry.Type == app.TargetDirectory {
		displayName += "/"
	}
	state := m.styles.stateStyle(entry.State).Render(stateMarker(entry.State))
	targetType := m.styles.typeStyle(entry.Type).Render(typeMarker(entry))
	name := m.styles.typeStyle(entry.Type).Render(displayName)
	return fmt.Sprintf("%s:%s %s%s", state, targetType, name, m.attributeMarker(entry))
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

func (m workspaceModel) attributeMarker(entry app.WorkspaceEntry) string {
	var marker strings.Builder
	if entry.Template {
		marker.WriteString(" " + m.styles.template.Render("[T]"))
	}
	if entry.Encrypted {
		marker.WriteString(" " + m.styles.encrypted.Render("[E]"))
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
	contentHeight := innerHeight
	var heading string
	if m.isWorkspace() {
		previewLabel := string(m.previewKind)
		if m.current().Type == app.TargetDirectory {
			previewLabel = "directory"
		}
		heading = m.styles.section.Render(truncate("Preview · "+previewLabel, innerWidth))
		contentHeight = max(1, innerHeight-1)
	}
	state := m.currentDiffState()
	var content string
	switch {
	case len(m.entries) == 0:
		content = m.styles.muted.Render("no entries match the current view")
	case state.loading:
		content = m.styles.loading.Render("loading preview...")
	case state.err != nil:
		content = m.styles.error.Render("error: " + state.err.Error())
	case len(state.lines) == 0:
		content = m.styles.muted.Render("no preview")
	default:
		lines := state.lines
		if m.diffScroll > len(lines) {
			lines = nil
		} else {
			lines = lines[m.diffScroll:]
		}
		if len(lines) > contentHeight {
			lines = lines[:contentHeight]
		}
		visible := make([]string, len(lines))
		for i, line := range lines {
			visible[i] = cropLine(line, m.previewX, innerWidth)
		}
		content = renderPreviewLines(m, visible, state)
	}

	content = padLines(content, contentHeight)
	if heading != "" {
		content = heading + "\n" + content
	}
	style := m.styles.pane
	if m.focus == focusDiff && m.mode == modeReview {
		style = m.styles.activePane
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
	lines := []string{m.styles.section.Render(title), ""}
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
	return m.styles.activePane.Width(size.width).Height(size.height).Render(content)
}

func (m workspaceModel) footer() string {
	if m.isWorkspace() {
		return m.workspaceFooter()
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
	return m.styles.help.Render(content)
}

func (m workspaceModel) filesPaneTitle() string {
	if !m.isWorkspace() {
		return "Files"
	}
	if m.searchResults {
		return fmt.Sprintf("Search · %s", m.filterLabel())
	}
	return fmt.Sprintf("Current · %s · %s", m.directoryLabel(), m.filterLabel())
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
	if cursor < 0 {
		cursor = 0
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
