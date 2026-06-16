package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type diffState struct {
	lines   []string
	loading bool
	err     error
}

func (m syncTUIModel) startDiffLoad(refresh bool) (syncTUIModel, tea.Cmd) {
	target := m.currentTarget()
	if target == "" {
		return m, nil
	}
	if state, ok := m.diffs[target]; ok && len(state.lines) > 0 && !refresh {
		return m, nil
	}
	m.diffs[target] = diffState{loading: true}
	return m, loadDiffCmd(m.service, target)
}

func loadDiffCmd(service interface{ DiffOutput(string) ([]byte, error) }, target string) tea.Cmd {
	return func() tea.Msg {
		out, err := service.DiffOutput(target)
		return syncDiffMsg{target: target, diff: string(out), err: err}
	}
}

func (m *syncTUIModel) applyDiff(msg syncDiffMsg) {
	state := diffState{lines: splitLines(strings.TrimRight(msg.diff, "\n")), err: msg.err}
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
	lines := m.currentDiffState().lines
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
