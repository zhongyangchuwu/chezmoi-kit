package ui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
)

func (m syncTUIModel) updateReview(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Key().Code == tea.KeyEnter || msg.Key().Code == tea.KeyReturn:
		if len(m.pending) == 0 {
			m.message = "no pending actions"
			return m, nil
		}
		m.mode = modeConfirm
		m.message = "confirm pending actions"
		return m, nil
	case msg.Key().Code == tea.KeyTab:
		m = m.toggleFocus()
		if m.focus == focusDiff {
			return m.startDiffLoad(false)
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, defaultSyncKeys.Up):
		return m.handleUp()
	case key.Matches(msg, defaultSyncKeys.Down):
		return m.handleDown()
	case key.Matches(msg, defaultSyncKeys.Diff):
		return m.startDiffLoad(true)
	case key.Matches(msg, defaultSyncKeys.Add):
		return m.togglePending(reconcile.ActionAdd), nil
	case key.Matches(msg, defaultSyncKeys.Apply):
		return m.togglePending(reconcile.ActionApply), nil
	case key.Matches(msg, defaultSyncKeys.Merge):
		return m.togglePending(reconcile.ActionMerge), nil
	case key.Matches(msg, defaultSyncKeys.Skip):
		return m.clearPending(), nil
	default:
		return m, nil
	}
}

func (m syncTUIModel) updateConfirm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, defaultSyncKeys.Back):
		m.mode = modeReview
		m.message = ""
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, defaultSyncKeys.Execute):
		actions := m.pendingActions()
		return m.startExecution(actions)
	default:
		return m, nil
	}
}

func (m syncTUIModel) updateExecuting(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m syncTUIModel) handleUp() (tea.Model, tea.Cmd) {
	if m.focus == focusDiff {
		return m.scrollDiff(-1, m.diffPaneHeight()), nil
	}
	var moved bool
	m, moved = m.moveUp()
	if moved {
		return m.startDiffLoad(false)
	}
	return m, nil
}

func (m syncTUIModel) handleDown() (tea.Model, tea.Cmd) {
	if m.focus == focusDiff {
		return m.scrollDiff(1, m.diffPaneHeight()), nil
	}
	var moved bool
	m, moved = m.moveDown()
	if moved {
		return m.startDiffLoad(false)
	}
	return m, nil
}

func (m syncTUIModel) diffPaneHeight() int {
	height := m.height
	if height <= 0 {
		height = 30
	}
	paneHeight := height - 6
	if paneHeight < 1 {
		return 1
	}
	return paneHeight
}
