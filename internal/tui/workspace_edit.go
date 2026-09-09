package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

func (m workspaceModel) startWorkspaceSourceEdit() (workspaceModel, tea.Cmd) {
	entry := m.current()
	if reason := app.SourceEditUnavailableReason(entry); reason != "" {
		m.message = "source edit unavailable: " + reason
		return m, nil
	}
	if m.workspace == nil {
		m.message = "source edit unavailable: workspace service is unavailable"
		return m, nil
	}
	m.workspaceBusy = true
	m.workspaceNotice = ""
	m.message = "opening source editor..."
	return m, prepareWorkspaceSourceEditCmd(m.workspace, entry, m.timing)
}

func prepareWorkspaceSourceEditCmd(service app.WorkspaceService, entry app.WorkspaceEntry, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		command, err := service.SourceEditCommand(entry)
		timing.Info("workspace source edit prepare", "target", entry.Path, "duration", elapsed(start), "err", err)
		return workspaceHandoffRequestMsg{entry: entry, command: command, err: err}
	}
}

func refreshWorkspaceCmd(service app.WorkspaceService, scopes []string, target string, editorErr error, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		snapshot, err := service.Inventory(scopes)
		timing.Info("workspace refresh after handoff", "target", target, "scopes", len(scopes), "entries", len(snapshot.Entries), "duration", elapsed(start), "err", err)
		return workspaceRefreshMsg{snapshot: snapshot, target: target, editorErr: editorErr, err: err}
	}
}

func (m workspaceModel) applyWorkspaceHandoffRequest(msg workspaceHandoffRequestMsg) (workspaceModel, tea.Cmd) {
	if msg.err != nil {
		m.workspaceBusy = false
		m.message = "source edit failed: " + msg.err.Error()
		return m, nil
	}
	start := time.Now()
	return m, tea.Exec(msg.command, func(err error) tea.Msg {
		m.timing.Info("workspace source edit", "target", msg.entry.Path, "duration", elapsed(start), "err", err)
		return workspaceHandoffDoneMsg{target: msg.entry.Path, err: err}
	})
}

func (m workspaceModel) applyWorkspaceHandoffDone(msg workspaceHandoffDoneMsg) (workspaceModel, tea.Cmd) {
	m.message = "refreshing workspace..."
	return m, refreshWorkspaceCmd(m.workspace, append([]string(nil), m.snapshot.Scopes...), msg.target, msg.err, m.timing)
}

func (m workspaceModel) applyWorkspaceRefresh(msg workspaceRefreshMsg) (workspaceModel, tea.Cmd) {
	m.workspaceBusy = false
	outcome := "editor closed"
	if msg.editorErr != nil {
		outcome = "editor exited with error: " + msg.editorErr.Error()
	}
	if msg.err != nil {
		m.clearWorkspacePreviewState()
		m.message = outcome + " · workspace refresh failed: " + msg.err.Error()
		return m, nil
	}

	previousDirectory := m.currentDir
	m.snapshot = msg.snapshot
	m.allEntries = append([]app.WorkspaceEntry(nil), msg.snapshot.Entries...)
	m.scriptCount = 0
	for _, entry := range m.allEntries {
		if entry.State == app.FileScript {
			m.scriptCount++
		}
	}
	m.clearWorkspacePreviewState()
	m.rebuildWorkspaceIndex()
	m.currentDir = m.closestWorkspaceDirectory(previousDirectory)
	m.rebuildEntries(msg.target)
	m.workspaceNotice = outcome + " · workspace refreshed · destination unchanged"
	m.message = ""
	return m.startDiffLoad(false)
}

func (m *workspaceModel) clearWorkspacePreviewState() {
	m.previewEpoch++
	m.diffs = make(map[string]diffState)
	m.revealedPreviews = make(map[string]bool)
	m.previewQuery = ""
	m.previewMatches = nil
	m.previewMatch = 0
	m.diffScroll = 0
	m.previewX = 0
}

func (m workspaceModel) closestWorkspaceDirectory(directory string) string {
	for directory != "" {
		entry, ok := m.workspaceNodes[directory]
		if ok && entry.Type == app.TargetDirectory {
			return directory
		}
		directory = workspaceParent(directory)
	}
	return ""
}

func (m workspaceModel) sourceEditFooter() string {
	if reason := app.SourceEditUnavailableReason(m.current()); reason == "" {
		return " · e edit source"
	}
	return ""
}
