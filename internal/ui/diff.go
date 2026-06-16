package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type diffState struct {
	content string
	loading bool
	err     error
}

func (m syncTUIModel) loadDiff(refresh bool) tea.Cmd {
	target := m.currentTarget()
	if target == "" {
		return nil
	}
	if state, ok := m.diffs[target]; ok && state.content != "" && !refresh {
		return nil
	}
	m.diffs[target] = diffState{loading: true}
	return func() tea.Msg {
		out, err := m.service.DiffOutput(target)
		return syncDiffMsg{target: target, diff: string(out), err: err}
	}
}

func (m *syncTUIModel) applyDiff(msg syncDiffMsg) {
	state := diffState{content: strings.TrimRight(msg.diff, "\n"), err: msg.err}
	m.diffs[msg.target] = state
	if msg.err != nil {
		m.message = msg.err.Error()
		return
	}
	m.message = ""
}

func (m syncTUIModel) currentDiffState() diffState {
	return m.diffs[m.currentTarget()]
}

func (m syncTUIModel) scrollDiff(delta int, height int) syncTUIModel {
	lines := splitLines(m.currentDiffState().content)
	maxScroll := len(lines) - height
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.diffScroll += delta
	if m.diffScroll < 0 {
		m.diffScroll = 0
	}
	if m.diffScroll > maxScroll {
		m.diffScroll = maxScroll
	}
	return m
}

func renderDiff(diff string) string {
	if diff == "" {
		return ""
	}

	lines := strings.Split(diff, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "@@"):
			lines[i] = diffHunkStyle.Render(line)
		case strings.HasPrefix(line, "diff "), strings.HasPrefix(line, "---"), strings.HasPrefix(line, "+++"):
			lines[i] = diffHeaderStyle.Render(line)
		case strings.HasPrefix(line, "+"):
			lines[i] = diffAddStyle.Render(line)
		case strings.HasPrefix(line, "-"):
			lines[i] = diffRemoveStyle.Render(line)
		case strings.HasPrefix(line, `\ No newline`):
			lines[i] = diffMetaStyle.Render(line)
		}
	}
	return strings.Join(lines, "\n")
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
