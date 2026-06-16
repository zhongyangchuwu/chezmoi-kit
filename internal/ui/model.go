package ui

import (
	"charm.land/bubbles/v2/help"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
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
	service    reconcile.ReviewService
	entries    []chezmoi.StatusEntry
	cursor     int
	pending    map[string]reconcile.ActionKind
	focus      syncFocus
	mode       syncMode
	diffs      map[string]diffState
	diffScroll int
	width      int
	height     int
	homeDir    string
	help       help.Model
	message    string
	err        error
}

func newSyncTUIModel(service reconcile.ReviewService, entries []chezmoi.StatusEntry) syncTUIModel {
	m := syncTUIModel{
		service: service,
		entries: append([]chezmoi.StatusEntry(nil), entries...),
		pending: make(map[string]reconcile.ActionKind),
		diffs:   make(map[string]diffState),
		homeDir: homeDir(),
		help:    help.New(),
	}
	if target := m.currentTarget(); target != "" {
		m.diffs[target] = diffState{loading: true}
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

func (m syncTUIModel) togglePending(kind reconcile.ActionKind) syncTUIModel {
	target := m.currentTarget()
	if current, ok := m.pending[target]; ok && current == kind {
		delete(m.pending, target)
		display := m.displayPath(target)
		m.message = "cleared " + display
		return m
	}
	m.pending[target] = kind
	m.message = actionLabel(kind) + " " + m.displayPath(target)
	return m
}

func (m syncTUIModel) clearPending() syncTUIModel {
	target := m.currentTarget()
	delete(m.pending, target)
	m.message = "skipped " + m.displayPath(target)
	return m
}

func (m syncTUIModel) pendingActions() []reconcile.Action {
	actions := make([]reconcile.Action, 0, len(m.pending))
	for _, entry := range m.entries {
		kind, ok := m.pending[entry.Path]
		if !ok {
			continue
		}
		actions = append(actions, reconcile.Action{Target: entry.Path, Kind: kind})
	}
	return actions
}

func (m syncTUIModel) pendingLabel(target string) string {
	kind, ok := m.pending[target]
	if !ok {
		return "[ ]"
	}
	return "[" + actionMarker(kind) + "]"
}
