package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

func (m workspaceModel) startExecution(actions []app.Action) (workspaceModel, tea.Cmd) {
	m.mode = modeExecuting
	m.executing = append([]app.Action(nil), actions...)
	m.executingIndex = 0
	m.executedCount = 0
	m.skippedCount = 0
	m.deferredCount = 0
	m.notices = nil
	m.executionStart = time.Now()
	m.message = "executing " + actionCount(len(actions))
	m.timing.Info("sync execution start", "actions", len(actions))
	return m, m.executeCurrentAction()
}

func (m workspaceModel) executeCurrentAction() tea.Cmd {
	if m.executingIndex >= len(m.executing) {
		return nil
	}
	action := m.executing[m.executingIndex]
	if action.Kind == app.ActionMerge {
		return prepareTerminalActionCmd(m.service, action, m.timing)
	}
	return executeNonInteractiveCmd(m.service, action, m.timing)
}

func executeNonInteractiveCmd(service app.SyncService, action app.Action, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		preflight := reviewAction(service, action, timing)
		if preflight != nil {
			return *preflight
		}

		start := time.Now()
		result, err := service.ExecuteNonInteractive(action)
		timing.Info("sync execute target", "target", action.Target, "action", action.Kind.String(), "mode", "buffered", "duration", elapsed(start), "err", err)
		if err != nil {
			return executeMsg{target: action.Target, action: action, result: result, err: err}
		}
		return postflightAction(service, action, result, timing)
	}
}

func prepareTerminalActionCmd(service app.SyncService, action app.Action, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		preflight := reviewAction(service, action, timing)
		if preflight != nil {
			return *preflight
		}
		cmd, err := service.TerminalCommand(action)
		if err != nil {
			return terminalRequestMsg{err: err}
		}
		return terminalRequestMsg{action: action, cmd: cmd}
	}
}

func reviewAction(service app.SyncService, action app.Action, timing *syncTimingLogger) *executeMsg {
	start := time.Now()
	review, err := service.Review(action.Target)
	fingerprintMatch := err == nil && review.Fingerprint == action.Fingerprint
	timing.Info("sync review preflight", "target", action.Target, "action", action.Kind.String(), "dirty", review.Dirty, "fingerprint_match", fingerprintMatch, "duration", elapsed(start), "err", err)
	if err != nil {
		return &executeMsg{target: action.Target, action: action, err: err}
	}
	if !review.Dirty {
		return &executeMsg{target: action.Target, action: action, skipped: true, resolved: true}
	}
	if !fingerprintMatch {
		return &executeMsg{target: action.Target, action: action, deferred: true, review: review, reason: "changed since review"}
	}
	if !review.Allows(action.Kind) {
		return &executeMsg{target: action.Target, action: action, deferred: true, review: review, reason: "action no longer valid for target type"}
	}
	return nil
}

func postflightAction(service app.SyncService, action app.Action, result app.ActionResult, timing *syncTimingLogger) executeMsg {
	start := time.Now()
	review, err := service.Review(action.Target)
	timing.Info("sync review postflight", "target", action.Target, "action", action.Kind.String(), "dirty", review.Dirty, "duration", elapsed(start), "err", err)
	if err != nil {
		return executeMsg{target: action.Target, action: action, result: result, executed: true, err: err}
	}
	if review.Dirty {
		return executeMsg{target: action.Target, action: action, result: result, executed: true, deferred: true, review: review, reason: "target still differs after execution"}
	}
	return executeMsg{target: action.Target, action: action, result: result, executed: true, resolved: true}
}

func postflightActionCmd(service app.SyncService, action app.Action, result app.ActionResult, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		return postflightAction(service, action, result, timing)
	}
}

func (m workspaceModel) applyTerminalRequestMsg(msg terminalRequestMsg) (workspaceModel, tea.Cmd) {
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

func (m workspaceModel) applyTerminalExecuteMsg(msg terminalExecuteMsg) (workspaceModel, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}
	return m, postflightActionCmd(m.service, msg.action, app.ActionResult{}, m.timing)
}

func (m workspaceModel) applyExecuteMsg(msg executeMsg) (workspaceModel, tea.Cmd) {
	if msg.err != nil {
		m.err = msg.err
		return m, tea.Quit
	}
	if msg.executed {
		m.executedCount++
	}
	if msg.skipped {
		m.skippedCount++
	}
	if notice := msg.result.Notice(); notice != "" {
		m.notices = append(m.notices, prefixNoticeLines(m.displayPath(msg.target), notice))
	}
	if msg.resolved {
		m = m.removeEntry(msg.target)
	} else if msg.deferred {
		m.deferredCount++
		delete(m.pending, msg.target)
		if msg.review.Dirty {
			m.storeReview(msg.review)
		}
		reason := m.displayPath(msg.target) + ": " + msg.reason
		m.notices = append(m.notices, reason)
		m.message = reason
	}

	m.executingIndex++
	if m.executingIndex < len(m.executing) {
		m.message = fmt.Sprintf("executed %d, skipped %d, deferred %d; %s", m.executedCount, m.skippedCount, m.deferredCount, m.displayPath(m.executing[m.executingIndex].Target))
		return m, m.executeCurrentAction()
	}
	return m.finishExecution()
}

func (m workspaceModel) finishExecution() (workspaceModel, tea.Cmd) {
	m.executing = nil
	m.executingIndex = 0
	m.mode = modeReview
	m.timing.Info("sync execution finish", "executed", m.executedCount, "skipped", m.skippedCount, "deferred", m.deferredCount, "remaining", len(m.entries), "duration", elapsed(m.executionStart))
	if len(m.entries) == 0 {
		m.completed = true
		m.stopped = true
		m.message = fmt.Sprintf("sync complete: executed %d, skipped %d, deferred %d", m.executedCount, m.skippedCount, m.deferredCount)
		if m.scriptCount > 0 {
			m.message += fmt.Sprintf("; %s pending outside cm sync", scriptCount(m.scriptCount))
		}
		m.message = appendNotices(m.message, m.notices)
		return m, tea.Quit
	}
	m.message = fmt.Sprintf("executed %d, skipped %d, deferred %d; %s remaining", m.executedCount, m.skippedCount, m.deferredCount, fileCount(len(m.entries)))
	m.message = appendNotices(m.message, m.notices)
	m.executedCount = 0
	m.skippedCount = 0
	m.deferredCount = 0
	m.notices = nil
	m.executionStart = time.Time{}
	return m.startDiffLoad(false)
}

func prefixNoticeLines(target, notice string) string {
	lines := strings.Split(notice, "\n")
	for i, line := range lines {
		lines[i] = target + ": " + line
	}
	return strings.Join(lines, "\n")
}

func appendNotices(message string, notices []string) string {
	if len(notices) == 0 {
		return message
	}
	return message + "\n" + strings.Join(notices, "\n")
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

func scriptCount(count int) string {
	if count == 1 {
		return "1 script"
	}
	return fmt.Sprintf("%d scripts", count)
}
