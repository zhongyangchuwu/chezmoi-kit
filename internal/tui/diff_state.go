package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type diffState struct {
	lines   []string
	loading bool
	err     error
}

func (m syncTUIModel) startDiffLoad(refresh bool) (syncTUIModel, tea.Cmd) {
	target := m.currentTarget()
	if target == "" || m.service == nil {
		return m, nil
	}
	if state, ok := m.diffs[target]; ok && len(state.lines) > 0 && !refresh {
		return m, nil
	}
	m.diffs[target] = diffState{loading: true}
	return m, loadDiffCmd(m.service, target, m.timing)
}

func loadDiffCmd(service interface{ DiffOutput(string) ([]byte, error) }, target string, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		out, err := service.DiffOutput(target)
		timing.Info("sync diff", "target", target, "bytes", len(out), "duration", elapsed(start), "err", err)
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

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
