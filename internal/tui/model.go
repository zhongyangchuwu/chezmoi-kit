package tui

import (
	"time"

	"charm.land/bubbles/v2/help"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

type syncMode int

const (
	modeReview syncMode = iota
	modeConfirm
	modeExecuting
)

type syncFocus int

const (
	focusFiles syncFocus = iota
	focusDiff
)

type syncTUIModel struct {
	service        ReviewService
	entries        []chezmoi.StatusEntry
	cursor         int
	pending        map[string]ActionKind
	focus          syncFocus
	mode           syncMode
	diffs          map[string]diffState
	diffScroll     int
	width          int
	height         int
	homeDir        string
	help           help.Model
	message        string
	err            error
	timing         *syncTimingLogger
	executionStart time.Time
	executing      []Action
	executingIndex int
	executedCount  int
	skippedCount   int
	completed      bool
	stopped        bool
}

func newSyncTUIModel(service ReviewService, entries []chezmoi.StatusEntry, timing ...*syncTimingLogger) syncTUIModel {
	m := syncTUIModel{
		service: service,
		entries: append([]chezmoi.StatusEntry(nil), entries...),
		pending: make(map[string]ActionKind),
		diffs:   make(map[string]diffState),
		homeDir: homeDir(),
		help:    help.New(),
	}
	if len(timing) > 0 {
		m.timing = timing[0]
	}
	if service != nil && len(m.entries) > 0 {
		m.diffs[m.currentTarget()] = diffState{loading: true}
	}
	return m
}

func (m syncTUIModel) current() chezmoi.StatusEntry {
	if len(m.entries) == 0 {
		return chezmoi.StatusEntry{}
	}
	return m.entries[m.cursor]
}

func (m syncTUIModel) currentTarget() string {
	return m.current().Path
}

func (m syncTUIModel) moveUp() (syncTUIModel, bool) {
	if m.cursor > 0 {
		m.cursor--
		m.diffScroll = 0
		return m, true
	}
	return m, false
}

func (m syncTUIModel) moveDown() (syncTUIModel, bool) {
	if m.cursor < len(m.entries)-1 {
		m.cursor++
		m.diffScroll = 0
		return m, true
	}
	return m, false
}

func (m syncTUIModel) toggleFocus() syncTUIModel {
	if m.focus == focusFiles {
		m.focus = focusDiff
	} else {
		m.focus = focusFiles
	}
	m.message = ""
	return m
}

func (m syncTUIModel) togglePending(kind ActionKind) syncTUIModel {
	target := m.currentTarget()
	if current, ok := m.pending[target]; ok && current == kind {
		delete(m.pending, target)
		display := m.displayPath(target)
		m.message = "cleared " + display
		return m
	}
	m.pending[target] = kind
	m.message = kind.String() + " " + m.displayPath(target)
	return m
}

func (m syncTUIModel) clearPending() syncTUIModel {
	target := m.currentTarget()
	delete(m.pending, target)
	m.message = "skipped " + m.displayPath(target)
	return m
}
func (m syncTUIModel) removeEntry(target string) syncTUIModel {
	delete(m.pending, target)
	delete(m.diffs, target)
	for i, entry := range m.entries {
		if entry.Path != target {
			continue
		}
		m.entries = append(m.entries[:i], m.entries[i+1:]...)
		if m.cursor >= len(m.entries) {
			m.cursor = len(m.entries) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.diffScroll = 0
		return m
	}
	return m
}

func (m syncTUIModel) pendingActions() []Action {
	actions := make([]Action, 0, len(m.pending))
	for _, entry := range m.entries {
		kind, ok := m.pending[entry.Path]
		if !ok {
			continue
		}
		actions = append(actions, Action{Target: entry.Path, Kind: kind})
	}
	return actions
}

func (m syncTUIModel) pendingLabel(target string) string {
	kind, ok := m.pending[target]
	if !ok {
		return "[ ]"
	}
	return "[" + kind.Marker() + "]"
}
