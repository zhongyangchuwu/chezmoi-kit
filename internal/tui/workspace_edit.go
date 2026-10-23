package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type workspaceHandoffKind string

const (
	workspaceSourceEdit      workspaceHandoffKind = "source edit"
	workspaceDestinationEdit workspaceHandoffKind = "destination edit"
	workspaceMerge           workspaceHandoffKind = "merge"
	workspaceSourceGit       workspaceHandoffKind = "source git"
)

func (k workspaceHandoffKind) normalized() workspaceHandoffKind {
	if k == "" {
		return workspaceSourceEdit
	}
	return k
}

func (k workspaceHandoffKind) openingMessage() string {
	switch k.normalized() {
	case workspaceDestinationEdit:
		return "opening destination editor..."
	case workspaceMerge:
		return "opening three-way merge: destination · source · rendered target..."
	case workspaceSourceGit:
		return "opening lazygit in chezmoi source repository..."
	default:
		return "opening source editor..."
	}
}

func (k workspaceHandoffKind) closedMessage(err error) string {
	var tool string
	switch k.normalized() {
	case workspaceDestinationEdit:
		tool = "destination editor"
	case workspaceMerge:
		tool = "merge tool"
	case workspaceSourceGit:
		tool = "lazygit"
	default:
		tool = "source editor"
	}
	if err != nil {
		return tool + " exited with error: " + err.Error()
	}
	return tool + " closed"
}

func (k workspaceHandoffKind) refreshedNotice(outcome string) string {
	if k.normalized() == workspaceSourceEdit {
		return outcome + " · workspace refreshed · destination unchanged"
	}
	return outcome + " · workspace refreshed"
}

func (m workspaceModel) startWorkspaceSourceEdit() (workspaceModel, tea.Cmd) {
	return m.startWorkspaceHandoff(workspaceSourceEdit, m.current(), app.SourceEditUnavailableReason(m.current()))
}

func (m workspaceModel) startWorkspaceDestinationEdit() (workspaceModel, tea.Cmd) {
	return m.startWorkspaceHandoff(workspaceDestinationEdit, m.current(), app.DestinationEditUnavailableReason(m.current()))
}

func (m workspaceModel) startWorkspaceMerge() (workspaceModel, tea.Cmd) {
	return m.startWorkspaceHandoff(workspaceMerge, m.current(), app.MergeUnavailableReason(m.current()))
}

func (m workspaceModel) startWorkspaceSourceGit() (workspaceModel, tea.Cmd) {
	return m.startWorkspaceHandoff(workspaceSourceGit, m.current(), "")
}

func (m workspaceModel) startWorkspaceHandoff(kind workspaceHandoffKind, entry app.WorkspaceEntry, reason string) (workspaceModel, tea.Cmd) {
	if reason != "" {
		m.message = string(kind) + " unavailable: " + reason
		return m, nil
	}
	if m.workspace == nil {
		m.message = string(kind) + " unavailable: workspace service is unavailable"
		return m, nil
	}
	m.workspaceBusy = true
	m.workspaceNotice = ""
	m.message = kind.openingMessage()
	return m, prepareWorkspaceHandoffCmd(m.workspace, entry, kind, m.timing)
}

func prepareWorkspaceHandoffCmd(service app.WorkspaceService, entry app.WorkspaceEntry, kind workspaceHandoffKind, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		var (
			command app.TerminalCommand
			err     error
		)
		switch kind.normalized() {
		case workspaceDestinationEdit:
			command, err = service.DestinationEditCommand(entry)
		case workspaceMerge:
			command, err = service.MergeCommand(entry)
		case workspaceSourceGit:
			command, err = service.SourceGitCommand()
		default:
			command, err = service.SourceEditCommand(entry)
		}
		timing.Info("workspace "+string(kind)+" prepare", "target", entry.Path, "duration", elapsed(start), "err", err)
		return workspaceHandoffRequestMsg{entry: entry, kind: kind, command: command, err: err}
	}
}

func refreshWorkspaceCmd(service app.WorkspaceService, scopes []string, target string, kind workspaceHandoffKind, toolErr error, timing *syncTimingLogger) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()
		snapshot, err := service.Inventory(scopes)
		timing.Info("workspace refresh after handoff", "target", target, "kind", kind, "scopes", len(scopes), "entries", len(snapshot.Entries), "duration", elapsed(start), "err", err)
		return workspaceRefreshMsg{snapshot: snapshot, target: target, kind: kind, toolErr: toolErr, err: err}
	}
}

func (m workspaceModel) applyWorkspaceHandoffRequest(msg workspaceHandoffRequestMsg) (workspaceModel, tea.Cmd) {
	kind := msg.kind.normalized()
	if msg.err != nil {
		m.workspaceBusy = false
		m.message = string(kind) + " failed: " + msg.err.Error()
		return m, nil
	}
	start := time.Now()
	return m, tea.Exec(msg.command, func(err error) tea.Msg {
		m.timing.Info("workspace "+string(kind), "target", msg.entry.Path, "duration", elapsed(start), "err", err)
		return workspaceHandoffDoneMsg{target: msg.entry.Path, kind: kind, err: err}
	})
}

func (m workspaceModel) applyWorkspaceHandoffDone(msg workspaceHandoffDoneMsg) (workspaceModel, tea.Cmd) {
	m.message = "refreshing workspace..."
	return m, refreshWorkspaceCmd(m.workspace, append([]string(nil), m.snapshot.Scopes...), msg.target, msg.kind, msg.err, m.timing)
}

func (m workspaceModel) applyWorkspaceRefresh(msg workspaceRefreshMsg) (workspaceModel, tea.Cmd) {
	m.workspaceBusy = false
	kind := msg.kind.normalized()
	outcome := kind.closedMessage(msg.toolErr)
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
	m.workspaceNotice = kind.refreshedNotice(outcome)
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

func (m workspaceModel) workspaceActionFooters() []string {
	actions := make([]string, 0, 4)
	entry := m.current()
	if app.SourceEditUnavailableReason(entry) == "" {
		actions = append(actions, "e source")
	}
	if app.DestinationEditUnavailableReason(entry) == "" {
		actions = append(actions, "E destination")
	}
	if app.MergeUnavailableReason(entry) == "" {
		actions = append(actions, "m merge")
	}
	if m.workspace != nil {
		actions = append(actions, "g git")
	}
	footers := make([]string, len(actions))
	for i := range actions {
		footers[i] = " · " + strings.Join(actions[:i+1], " · ")
	}
	return footers
}
