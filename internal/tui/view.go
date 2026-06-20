package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

type rect struct {
	width  int
	height int
}

func (m syncTUIModel) View() tea.View {
	return tea.NewView(m.viewString())
}

func (m syncTUIModel) viewString() string {
	if len(m.entries) == 0 {
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

	header := titleStyle.Render("cm sync")
	footer := m.footer()
	messageHeight := 0
	if m.message != "" {
		messageHeight = 1
	}
	bodyHeight := height - lipgloss.Height(header) - lipgloss.Height(footer) - messageHeight - 2
	if bodyHeight < 6 {
		bodyHeight = 6
	}

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
	body := lipgloss.JoinHorizontal(lipgloss.Top, files, main)

	parts := []string{header, body}
	if m.message != "" {
		parts = append(parts, m.message)
	}
	parts = append(parts, footer)
	return strings.Join(parts, "\n") + "\n"
}

func (m syncTUIModel) renderFilesPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	lines := make([]string, 0, innerHeight)
	for _, entry := range visibleEntries(m.entries, m.cursor, innerHeight) {
		cursor := "  "
		if entry.index == m.cursor {
			cursor = "> "
		}
		line := cursor + m.pendingLabel(entry.path) + " " + m.displayPath(entry.path)
		lines = append(lines, truncate(line, innerWidth))
	}
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}

	style := paneStyle
	if m.focus == focusFiles && m.mode == modeReview {
		style = activePaneStyle
	}
	return style.Width(size.width - 2).Height(size.height - 2).Render(strings.Join(lines, "\n"))
}

func (m syncTUIModel) renderMainPane(size rect) string {
	if m.mode == modeConfirm || m.mode == modeExecuting {
		return m.renderConfirmPane(size)
	}
	return m.renderDiffPane(size)
}

func (m syncTUIModel) renderDiffPane(size rect) string {
	innerWidth := size.width - 2
	innerHeight := size.height - 2
	state := m.currentDiffState()
	var content string
	switch {
	case state.loading:
		content = "loading diff..."
	case state.err != nil:
		content = state.err.Error()
	case len(state.lines) == 0:
		content = "no diff"
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
			visible[i] = truncate(line, innerWidth)
		}
		content = renderDiffLines(visible)
	}

	content = padLines(content, innerHeight)
	style := paneStyle
	if m.focus == focusDiff && m.mode == modeReview {
		style = activePaneStyle
	}
	return style.Width(size.width - 2).Height(size.height - 2).Render(content)
}

func (m syncTUIModel) renderConfirmPane(size rect) string {
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
		lines = append(lines, "", fmt.Sprintf("executed: %d", m.executedCount), fmt.Sprintf("skipped: %d", m.skippedCount))
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
	return activePaneStyle.Width(size.width - 2).Height(size.height - 2).Render(content)
}

func (m syncTUIModel) footer() string {
	if m.mode == modeExecuting {
		return helpStyle.Render("executing...")
	}
	if m.mode == modeConfirm {
		return helpStyle.Render(m.help.ShortHelpView(defaultSyncKeys.confirmHelp()))
	}
	focus := "files"
	if m.focus == focusDiff {
		focus = "diff"
	}
	return helpStyle.Render(focus + " • " + m.help.ShortHelpView(defaultSyncKeys.reviewHelp(m.focus)))
}

type visibleEntry struct {
	index int
	path  string
}

func visibleEntries(entries []chezmoi.StatusEntry, cursor int, height int) []visibleEntry {
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
		visible = append(visible, visibleEntry{index: i, path: entries[i].Path})
	}
	return visible
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(s) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return s[:width-1] + "…"
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
