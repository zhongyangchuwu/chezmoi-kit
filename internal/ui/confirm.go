package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
)

func (m syncTUIModel) startExecution(actions []reconcile.Action) (syncTUIModel, tea.Cmd) {
	m.mode = modeExecuting
	m.executing = append([]reconcile.Action(nil), actions...)
	m.executingIndex = 0
	m.executedCount = 0
	m.skippedCount = 0
	m.message = "executing " + actionCount(len(actions))
	return m, m.executeCurrentAction()
}

func (m syncTUIModel) executeCurrentAction() tea.Cmd {
	if m.executingIndex >= len(m.executing) {
		return nil
	}
	return executeActionCmd(m.service, m.executing[m.executingIndex])
}

func executeActionCmd(service reconcile.ReviewService, action reconcile.Action) tea.Cmd {
	return func() tea.Msg {
		dirty, err := service.Status([]string{action.Target})
		if err != nil {
			return executeMsg{err: err}
		}
		if len(dirty) == 0 {
			return executeMsg{target: action.Target, skipped: 1}
		}
		if err := service.ExecuteOne(action); err != nil {
			return executeMsg{err: err}
		}
		return executeMsg{target: action.Target, executed: []reconcile.Action{action}}
	}
}

func (m syncTUIModel) applyExecuteMsg(msg executeMsg) (syncTUIModel, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}
	if len(msg.executed) > 0 {
		m.executedCount += len(msg.executed)
		for _, action := range msg.executed {
			m = m.removeEntry(action.Target)
		}
	} else if msg.skipped > 0 {
		m = m.removeEntry(msg.target)
	}
	m.skippedCount += msg.skipped
	m.executingIndex++
	if m.executingIndex < len(m.executing) {
		m.message = fmt.Sprintf("executed %d, skipped %d; %s", m.executedCount, m.skippedCount, m.displayPath(m.executing[m.executingIndex].Target))
		return m, m.executeCurrentAction()
	}
	return m.finishExecution()
}

func (m syncTUIModel) finishExecution() (syncTUIModel, tea.Cmd) {
	m.executing = nil
	m.executingIndex = 0
	m.mode = modeReview
	if len(m.entries) == 0 {
		m.completed = true
		m.stopped = true
		m.message = fmt.Sprintf("sync complete: executed %d, skipped %d", m.executedCount, m.skippedCount)
		return m, tea.Quit
	}
	m.message = fmt.Sprintf("executed %d, skipped %d; %s remaining", m.executedCount, m.skippedCount, fileCount(len(m.entries)))
	m.executedCount = 0
	m.skippedCount = 0
	return m.startDiffLoad(false)
}

func actionCount(count int) string {
	if count == 1 {
		return "1 action"
	}
	return fmt.Sprintf("%d actions", count)
}

func fileCount(count int) string {
	if count == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", count)
}

func actionLabel(kind reconcile.ActionKind) string {
	return kind.String()
}

func actionMarker(kind reconcile.ActionKind) string {
	return kind.Marker()
}
