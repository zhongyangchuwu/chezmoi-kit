package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"github.com/zhongyangchuwu/cm/internal/app"
)

type syncMode int

const (
	modeReview syncMode = iota
	modeConfirm
	modeExecuting
)

type runKind int

const (
	runSync runKind = iota
	runWorkspace
)

type syncFocus int

const (
	focusFiles syncFocus = iota
	focusDiff
)

type workspaceFilter int

const (
	filterAll workspaceFilter = iota
	filterManaged
	filterDirty
	filterUnmanaged
	filterIgnored
	filterScripts
)

type searchKind int

const (
	searchNone searchKind = iota
	searchFiles
	searchPreview
)

type workspaceModel struct {
	service          app.SyncService
	workspace        app.WorkspaceService
	runKind          runKind
	snapshot         app.WorkspaceSnapshot
	allEntries       []app.WorkspaceEntry
	entries          []app.WorkspaceEntry
	scriptCount      int
	searchOriginal   string
	searchTarget     string
	searchOriginDir  string
	searchResults    bool
	cursor           int
	pending          map[string]app.Action
	focus            syncFocus
	mode             syncMode
	diffs            map[string]diffState
	previewKind      app.PreviewKind
	previewFull      bool
	previewX         int
	diffScroll       int
	previewQuery     string
	previewMatches   []int
	previewEpoch     uint64
	previewMatch     int
	filter           workspaceFilter
	currentDir       string
	directoryCursors map[string]int
	workspaceNodes   map[string]app.WorkspaceEntry
	matchingNodes    map[string]bool
	fileQuery        string
	search           searchKind
	searchInput      string
	helpVisible      bool
	workspaceBusy    bool
	workspaceNotice  string
	width            int
	height           int
	homeDir          string
	styles           tuiStyles
	help             help.Model
	message          string
	err              error
	timing           *syncTimingLogger
	executionStart   time.Time
	executing        []app.Action
	executingIndex   int
	executedCount    int
	skippedCount     int
	deferredCount    int
	notices          []string
	revealedPreviews map[string]bool
	completed        bool
	stopped          bool
}

func newSyncTUIModel(service app.SyncService, status app.SyncStatus, timing ...*syncTimingLogger) workspaceModel {
	entries := make([]app.WorkspaceEntry, 0, len(status.Entries))
	for _, entry := range status.Entries {
		entries = append(entries, app.WorkspaceEntry{
			Path:         entry.Path,
			RelativePath: entry.Path,
			State:        app.FileDirty,
			Type:         app.TargetUnknown,
			Code:         entry.Code,
		})
	}
	m := newBaseModel(service, timing...)
	m.runKind = runSync
	m.entries = entries
	m.allEntries = append([]app.WorkspaceEntry(nil), entries...)
	m.scriptCount = len(status.Scripts)
	m.markCurrentLoading()
	return m
}

func newWorkspaceModel(service app.WorkspaceService, snapshot app.WorkspaceSnapshot, timing ...*syncTimingLogger) workspaceModel {
	m := newBaseModel(service, timing...)
	m.runKind = runWorkspace
	m.workspace = service
	m.snapshot = snapshot
	m.allEntries = append([]app.WorkspaceEntry(nil), snapshot.Entries...)
	for _, entry := range snapshot.Entries {
		if entry.State == app.FileScript {
			m.scriptCount++
		}
	}
	m.currentDir = initialWorkspaceDirectory(snapshot)
	m.rebuildEntries("")
	m.markCurrentLoading()
	return m
}

func newBaseModel(service app.SyncService, timing ...*syncTimingLogger) workspaceModel {
	m := workspaceModel{
		service:          service,
		pending:          make(map[string]app.Action),
		diffs:            make(map[string]diffState),
		previewKind:      app.PreviewDiff,
		filter:           filterAll,
		directoryCursors: make(map[string]int),
		workspaceNodes:   make(map[string]app.WorkspaceEntry),
		matchingNodes:    make(map[string]bool),
		revealedPreviews: make(map[string]bool),
		homeDir:          homeDir(),
		styles:           newTUIStyles(),
		help:             help.New(),
	}
	if len(timing) > 0 {
		m.timing = timing[0]
	}
	return m
}

func (m workspaceModel) isWorkspace() bool {
	return m.runKind == runWorkspace
}

func (m workspaceModel) current() app.WorkspaceEntry {
	if len(m.entries) == 0 {
		return app.WorkspaceEntry{}
	}
	return m.entries[m.cursor]
}

func (m workspaceModel) currentTarget() string {
	return m.current().Path
}

func (m *workspaceModel) markCurrentLoading() {
	if m.service == nil || m.currentTarget() == "" {
		return
	}
	m.diffs[m.previewStateKey()] = diffState{loading: true}
}

func (m workspaceModel) moveUp() (workspaceModel, bool) {
	if m.cursor > 0 {
		m.cursor--
		m.rememberCursor()
		m.resetPreviewPosition()
		return m, true
	}
	return m, false
}

func (m workspaceModel) moveDown() (workspaceModel, bool) {
	if m.cursor < len(m.entries)-1 {
		m.cursor++
		m.rememberCursor()
		m.resetPreviewPosition()
		return m, true
	}
	return m, false
}

func (m *workspaceModel) resetPreviewPosition() {
	m.diffScroll = 0
	m.previewX = 0
	m.previewMatches = nil
	m.previewMatch = 0
}

func (m workspaceModel) toggleFocus() workspaceModel {
	if m.focus == focusFiles {
		m.focus = focusDiff
	} else {
		m.focus = focusFiles
	}
	m.message = ""
	return m
}

func (m workspaceModel) togglePending(kind app.ActionKind) workspaceModel {
	if m.isWorkspace() {
		m.message = "cm ui is read-only in this phase; use cm sync to reconcile"
		return m
	}
	target := m.currentTarget()
	state, ok := m.diffs[m.previewStateKey()]
	if !ok || state.loading {
		m.message = "review still loading for " + m.displayPath(target)
		return m
	}
	if state.err != nil {
		m.message = state.err.Error()
		return m
	}
	action, err := state.review.Action(kind)
	if err != nil {
		m.message = err.Error()
		return m
	}
	if current, ok := m.pending[target]; ok && current.Kind == kind && current.Fingerprint == action.Fingerprint {
		delete(m.pending, target)
		m.message = "cleared " + m.displayPath(target)
		return m
	}
	m.pending[target] = action
	m.message = kind.String() + " " + m.displayPath(target)
	return m
}

func (m workspaceModel) clearPending() workspaceModel {
	target := m.currentTarget()
	delete(m.pending, target)
	m.message = "skipped " + m.displayPath(target)
	return m
}

func (m workspaceModel) removeEntry(target string) workspaceModel {
	delete(m.pending, target)
	for key := range m.diffs {
		if key == target || strings.HasPrefix(key, target+"\x00") {
			delete(m.diffs, key)
		}
	}
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
		m.resetPreviewPosition()
		return m
	}
	return m
}

func (m workspaceModel) pendingActions() []app.Action {
	actions := make([]app.Action, 0, len(m.pending))
	for _, entry := range m.entries {
		action, ok := m.pending[entry.Path]
		if !ok {
			continue
		}
		actions = append(actions, action)
	}
	return actions
}

func (m workspaceModel) pendingLabel(target string) string {
	action, ok := m.pending[target]
	if !ok {
		return "[ ]"
	}
	return "[" + action.Kind.Marker() + "]"
}

func (m workspaceModel) previewStateKey() string {
	if !m.isWorkspace() {
		return m.currentTarget()
	}
	return fmt.Sprintf("%s\x00%s\x00%t", m.currentTarget(), m.previewKind, m.currentPreviewRevealed())
}

func (m workspaceModel) revealKey() string {
	return fmt.Sprintf("%s\x00%s", m.currentTarget(), m.previewKind)
}

func (m workspaceModel) currentPreviewRevealed() bool {
	return m.revealedPreviews[m.revealKey()]
}
