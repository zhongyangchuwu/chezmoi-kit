package ui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

func (m syncTUIModel) startExecution(actions []Action) (syncTUIModel, tea.Cmd) {
	m.mode = modeExecuting
	m.executing = append([]Action(nil), actions...)
	m.executingIndex = 0
	m.executedCount = 0
	m.skippedCount = 0
	m.executionStart = time.Now()
	m.message = "executing " + actionCount(len(actions))
	m.timing.Info("sync execution start", "actions", len(actions))
	return m, m.executeCurrentAction()
}

func (m syncTUIModel) executeCurrentAction() tea.Cmd {
	if m.executingIndex >= len(m.executing) {
		return nil
	}
	action := m.executing[m.executingIndex]
	if action.Kind == ActionMerge {
		return prepareTerminalActionCmd(m.service, action, m.timing)
	}
	return executeNonInteractiveCmd(m.service, action, m.timing)
}

func executeNonInteractiveCmd(service ReviewService, action Action, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		dirty, err := service.Status([]string{action.Target})
		timing.Info("sync status preflight", "target", action.Target, "action", action.Kind.String(), "dirty", len(dirty) > 0, "duration", elapsed(start), "err", err)
		if err != nil {
			return executeMsg{err: err}
		}
		if len(dirty) == 0 {
			return executeMsg{target: action.Target, skipped: 1}
		}
		start = time.Now()
		err = service.ExecuteNonInteractive(action)
		timing.Info("sync execute target", "target", action.Target, "action", action.Kind.String(), "mode", "buffered", "duration", elapsed(start), "err", err)
		if err != nil {
			return executeMsg{err: err}
		}
		return executeMsg{target: action.Target, executed: []Action{action}}
	}
}

func prepareTerminalActionCmd(service ReviewService, action Action, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		dirty, err := service.Status([]string{action.Target})
		timing.Info("sync status preflight", "target", action.Target, "action", action.Kind.String(), "dirty", len(dirty) > 0, "duration", elapsed(start), "err", err)
		if err != nil {
			return terminalRequestMsg{err: err}
		}
		if len(dirty) == 0 {
			return executeMsg{target: action.Target, skipped: 1}
		}
		cmd, err := service.TerminalCommand(action)
		if err != nil {
			return terminalRequestMsg{err: err}
		}
		return terminalRequestMsg{action: action, cmd: cmd}
	}
}

func (m syncTUIModel) applyTerminalRequestMsg(msg terminalRequestMsg) (syncTUIModel, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}
	start := time.Now()
	return m, tea.Exec(msg.cmd, func(err error) tea.Msg {
		m.timing.Info("sync execute target", "target", msg.action.Target, "action", msg.action.Kind.String(), "mode", "terminal", "duration", elapsed(start), "err", err)
		return terminalExecuteMsg{action: msg.action, err: err}
	})
}

func (m syncTUIModel) applyTerminalExecuteMsg(msg terminalExecuteMsg) (syncTUIModel, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}
	return m.applyExecuteMsg(executeMsg{target: msg.action.Target, executed: []Action{msg.action}})
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
	m.timing.Info("sync execution finish", "executed", m.executedCount, "skipped", m.skippedCount, "remaining", len(m.entries), "duration", elapsed(m.executionStart))
	if len(m.entries) == 0 {
		m.completed = true
		m.stopped = true
		m.message = fmt.Sprintf("sync complete: executed %d, skipped %d", m.executedCount, m.skippedCount)
		return m, tea.Quit
	}
	m.message = fmt.Sprintf("executed %d, skipped %d; %s remaining", m.executedCount, m.skippedCount, fileCount(len(m.entries)))
	m.executedCount = 0
	m.skippedCount = 0
	m.executionStart = time.Time{}
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
