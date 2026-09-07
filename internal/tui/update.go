package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/app"
)

func (m workspaceModel) updateReview(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.isWorkspace() {
		return m.updateWorkspaceReview(msg)
	}
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
		return m.togglePending(app.ActionAdd), nil
	case key.Matches(msg, defaultSyncKeys.Apply):
		return m.togglePending(app.ActionApply), nil
	case key.Matches(msg, defaultSyncKeys.Merge):
		return m.togglePending(app.ActionMerge), nil
	case key.Matches(msg, defaultSyncKeys.Skip):
		return m.clearPending(), nil
	default:
		return m, nil
	}
}

func (m workspaceModel) updateWorkspaceReview(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, defaultSyncKeys.Quit):
		return m, tea.Quit
	case msg.Key().Code == tea.KeyTab:
		if m.previewFull {
			m.previewFull = false
		}
		m = m.toggleFocus()
		return m.startDiffLoad(false)
	case key.Matches(msg, defaultSyncKeys.Up):
		return m.handleUp()
	case key.Matches(msg, defaultSyncKeys.Down):
		return m.handleDown()
	case key.Matches(msg, defaultSyncKeys.PageUp):
		if m.focus == focusDiff {
			return m.scrollDiff(-max(1, m.diffPaneHeight()/2), m.diffPaneHeight()), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.PageDown):
		if m.focus == focusDiff {
			return m.scrollDiff(max(1, m.diffPaneHeight()/2), m.diffPaneHeight()), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Left):
		if m.focus == focusDiff {
			return m.scrollPreviewHorizontal(-4, m.previewPaneWidth()), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Right):
		if m.focus == focusDiff {
			return m.scrollPreviewHorizontal(4, m.previewPaneWidth()), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.Diff):
		return m.startDiffLoad(true)
	case key.Matches(msg, defaultSyncKeys.Tree):
		oldTarget := m.currentTarget()
		m.toggleTreeMode()
		view := "tree"
		if m.flat {
			view = "flat"
		}
		m.message = "view: " + view
		return m.loadIfSelectionChanged(oldTarget)
	case key.Matches(msg, defaultSyncKeys.Filter):
		oldTarget := m.currentTarget()
		m.cycleFilter()
		m.message = "filter: " + m.filterLabel()
		return m.loadIfSelectionChanged(oldTarget)
	case key.Matches(msg, defaultSyncKeys.ToggleDirectory):
		oldTarget := m.currentTarget()
		m.toggleCurrentDirectory()
		return m.loadIfSelectionChanged(oldTarget)
	case key.Matches(msg, defaultSyncKeys.Search):
		return m.beginSearch(), nil
	case key.Matches(msg, defaultSyncKeys.Full):
		m.previewFull = !m.previewFull
		if m.previewFull {
			m.focus = focusDiff
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.ViewDiff):
		return m.setPreviewKind(app.PreviewDiff)
	case key.Matches(msg, defaultSyncKeys.ViewDestination):
		return m.setPreviewKind(app.PreviewDestination)
	case key.Matches(msg, defaultSyncKeys.ViewTarget):
		return m.setPreviewKind(app.PreviewTarget)
	case key.Matches(msg, defaultSyncKeys.ViewSource):
		return m.setPreviewKind(app.PreviewSource)
	case key.Matches(msg, defaultSyncKeys.Reveal):
		return m.revealCurrentPreview()
	case key.Matches(msg, defaultSyncKeys.HunkPrevious):
		if m.focus == focusDiff && m.previewKind == app.PreviewDiff {
			return m.jumpHunk(-1), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.HunkNext):
		if m.focus == focusDiff && m.previewKind == app.PreviewDiff {
			return m.jumpHunk(1), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.MatchPrevious):
		if m.focus == focusDiff {
			return m.movePreviewMatch(-1), nil
		}
		return m, nil
	case key.Matches(msg, defaultSyncKeys.MatchNext):
		if m.focus == focusDiff {
			return m.movePreviewMatch(1), nil
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m workspaceModel) beginSearch() workspaceModel {
	m.searchTarget = m.currentTarget()
	if m.focus == focusDiff {
		m.search = searchPreview
		m.searchInput = m.previewQuery
		m.searchOriginal = m.previewQuery
	} else {
		m.search = searchFiles
		m.searchInput = m.fileQuery
		m.searchOriginal = m.fileQuery
	}
	return m
}

func (m workspaceModel) updateSearch(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.Key().Code {
	case tea.KeyEnter:
		search := m.search
		m.search = searchNone
		m.message = ""
		if search == searchFiles {
			return m.startDiffLoad(false)
		}
		return m, nil
	case tea.KeyEscape:
		oldTarget := m.currentTarget()
		if m.search == searchFiles {
			m.fileQuery = m.searchOriginal
			m.rebuildEntries(m.searchTarget)
		} else {
			m.previewQuery = m.searchOriginal
			m.updatePreviewMatches()
		}
		m.search = searchNone
		m.searchInput = ""
		m.message = ""
		return m.loadIfSelectionChanged(oldTarget)
	case tea.KeyBackspace, tea.KeyDelete:
		m.searchInput = trimLastRune(m.searchInput)
	default:
		if msg.Key().Text != "" && msg.Key().Mod == 0 {
			m.searchInput += msg.Key().Text
		} else {
			return m, nil
		}
	}

	oldTarget := m.currentTarget()
	if m.search == searchFiles {
		m.fileQuery = m.searchInput
		m.rebuildEntries(oldTarget)
		return m, nil
	}
	m.previewQuery = m.searchInput
	m.updatePreviewMatches()
	return m, nil
}

func trimLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}
	return string(runes[:len(runes)-1])
}

func (m workspaceModel) loadIfSelectionChanged(oldTarget string) (workspaceModel, tea.Cmd) {
	if m.currentTarget() != oldTarget {
		return m.startDiffLoad(false)
	}
	return m, nil
}

func (m workspaceModel) updateConfirm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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

func (m workspaceModel) updateExecuting(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m workspaceModel) handleUp() (tea.Model, tea.Cmd) {
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

func (m workspaceModel) handleDown() (tea.Model, tea.Cmd) {
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

func (m workspaceModel) diffPaneHeight() int {
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

func (m workspaceModel) previewPaneWidth() int {
	width := m.width
	if width <= 0 {
		width = 100
	}
	if m.previewFull {
		return max(1, width-2)
	}
	leftWidth := clamp(width/3, 28, 48)
	if leftWidth > width-24 {
		leftWidth = width / 2
	}
	return max(1, width-leftWidth-3)
}
