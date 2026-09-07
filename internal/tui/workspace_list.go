package tui

import (
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func (m *workspaceModel) rebuildEntries(preserveTarget string) {
	if !m.isWorkspace() {
		return
	}
	if preserveTarget == "" {
		preserveTarget = m.currentTarget()
	}
	if m.searchResults {
		m.entries = m.searchEntries()
	} else {
		m.entries = m.directoryEntries(m.currentDir)
	}
	m.cursor = m.cursorFor(m.currentDir, preserveTarget)
	m.resetPreviewPosition()
}

func (m workspaceModel) directoryEntries(directory string) []app.WorkspaceEntry {
	nodes := m.directoryNodes()
	entries := make([]app.WorkspaceEntry, 0)
	for relative, entry := range nodes {
		if workspaceParent(relative) != directory || !m.entryOrDescendantMatches(relative) {
			continue
		}
		entries = append(entries, entry)
	}
	sortWorkspaceEntries(entries)
	return entries
}

func (m workspaceModel) parentEntries() []app.WorkspaceEntry {
	if m.currentDir == "" {
		return nil
	}
	return m.directoryEntries(workspaceParent(m.currentDir))
}

func (m workspaceModel) searchEntries() []app.WorkspaceEntry {
	query := strings.ToLower(strings.TrimSpace(m.fileQuery))
	if query == "" {
		return nil
	}
	nodes := m.directoryNodes()
	entries := make([]app.WorkspaceEntry, 0)
	for relative, entry := range nodes {
		if !m.entryOrDescendantMatches(relative) {
			continue
		}
		if strings.Contains(strings.ToLower(entry.RelativePath+"\x00"+entry.SourcePath), query) {
			entries = append(entries, entry)
		}
	}
	sortWorkspaceEntries(entries)
	return entries
}
func (m workspaceModel) directoryNodes() map[string]app.WorkspaceEntry {
	nodes := make(map[string]app.WorkspaceEntry, len(m.allEntries))
	states := make(map[string][]app.FileState)
	for _, entry := range m.allEntries {
		relative := cleanWorkspaceRelative(entry.RelativePath)
		if relative == "" {
			continue
		}
		nodes[relative] = entry
		for ancestor := path.Dir(relative); ancestor != "."; ancestor = path.Dir(ancestor) {
			states[ancestor] = append(states[ancestor], entry.State)
		}
	}
	for relative, states := range states {
		if _, exists := nodes[relative]; exists {
			continue
		}
		nodes[relative] = app.WorkspaceEntry{
			Path:         filepath.Join(m.snapshot.Root, filepath.FromSlash(relative)),
			RelativePath: relative,
			State:        aggregateDirectoryState(states),
			Type:         app.TargetDirectory,
		}
	}
	return nodes
}

func (m workspaceModel) entryOrDescendantMatches(relative string) bool {
	for _, entry := range m.allEntries {
		if entry.RelativePath == relative || strings.HasPrefix(entry.RelativePath, relative+"/") {
			if workspaceEntryMatchesFilter(entry, m.filter) {
				return true
			}
		}
	}
	return false
}

func (m workspaceModel) cursorFor(directory, target string) int {
	if len(m.entries) == 0 {
		return 0
	}
	if target != "" {
		for index, entry := range m.entries {
			if entry.Path == target {
				m.directoryCursors[directory] = index
				return index
			}
		}
	}
	if cursor := m.directoryCursors[directory]; cursor >= 0 && cursor < len(m.entries) {
		return cursor
	}
	return 0
}

func (m *workspaceModel) replaceWorkspaceEntry(updated app.WorkspaceEntry) {
	if updated.Path == "" {
		return
	}
	currentTarget := m.currentTarget()
	for i, entry := range m.allEntries {
		if entry.Path == updated.Path && entry.State != updated.State {
			m.allEntries[i] = updated
			m.rebuildEntries(currentTarget)
			return
		}
	}
}

func (m *workspaceModel) rememberCursor() {
	if m.isWorkspace() && !m.searchResults && len(m.entries) > 0 {
		m.directoryCursors[m.currentDir] = m.cursor
	}
}

func (m *workspaceModel) enterCurrentDirectory() bool {
	entry := m.current()
	if m.searchResults {
		m.searchResults = false
		m.fileQuery = ""
		if entry.Type == app.TargetDirectory {
			m.currentDir = entry.RelativePath
		} else {
			m.currentDir = workspaceParent(entry.RelativePath)
		}
		m.rebuildEntries(entry.Path)
		return true
	}
	if entry.Type != app.TargetDirectory {
		return false
	}
	m.rememberCursor()
	m.currentDir = entry.RelativePath
	m.rebuildEntries("")
	return true
}

func (m *workspaceModel) leaveCurrentDirectory() bool {
	if m.searchResults || m.currentDir == "" {
		return false
	}
	child := m.currentDir
	m.rememberCursor()
	m.currentDir = workspaceParent(child)
	m.rebuildEntries(filepath.Join(m.snapshot.Root, filepath.FromSlash(child)))
	return true
}

func (m *workspaceModel) cycleFilter() {
	target := m.currentTarget()
	m.filter = (m.filter + 1) % (filterScripts + 1)
	m.rebuildEntries(target)
}

func workspaceEntryMatchesFilter(entry app.WorkspaceEntry, filter workspaceFilter) bool {
	switch filter {
	case filterManaged:
		return entry.State == app.FileClean || entry.State == app.FileDirty || entry.State == app.FileScript || entry.State == app.FileUninspected
	case filterDirty:
		return entry.State == app.FileDirty
	case filterUnmanaged:
		return entry.State == app.FileUnmanaged
	case filterIgnored:
		return entry.State == app.FileIgnored
	case filterScripts:
		return entry.State == app.FileScript
	default:
		return true
	}
}

func (m workspaceModel) filterLabel() string {
	switch m.filter {
	case filterManaged:
		return "managed"
	case filterDirty:
		return "dirty"
	case filterUnmanaged:
		return "unmanaged"
	case filterIgnored:
		return "ignored"
	case filterScripts:
		return "scripts"
	default:
		return "all"
	}
}

func (m workspaceModel) directoryLabel() string {
	if m.currentDir == "" {
		return "/"
	}
	return m.currentDir + "/"
}

func initialWorkspaceDirectory(snapshot app.WorkspaceSnapshot) string {
	if len(snapshot.Scopes) != 1 {
		return ""
	}
	relative, err := filepath.Rel(snapshot.Root, snapshot.Scopes[0])
	if err != nil || relative == "." || strings.HasPrefix(relative, "..") {
		return ""
	}
	relative = filepath.ToSlash(relative)
	for _, entry := range snapshot.Entries {
		if entry.RelativePath == relative && entry.Type != app.TargetDirectory {
			return workspaceParent(relative)
		}
	}
	return relative
}

func aggregateDirectoryState(states []app.FileState) app.FileState {
	best := app.FileClean
	for _, state := range states {
		if directoryStateRank(state) > directoryStateRank(best) {
			best = state
		}
	}
	return best
}

func directoryStateRank(state app.FileState) int {
	switch state {
	case app.FileDirty:
		return 6
	case app.FileUninspected:
		return 5
	case app.FileUnmanaged:
		return 4
	case app.FileIgnored:
		return 3
	case app.FileScript:
		return 2
	default:
		return 1
	}
}

func cleanWorkspaceRelative(relative string) string {
	relative = path.Clean(strings.Trim(relative, "/"))
	if relative == "." {
		return ""
	}
	return relative
}

func workspaceParent(relative string) string {
	parent := path.Dir(relative)
	if parent == "." {
		return ""
	}
	return parent
}

func sortWorkspaceEntries(entries []app.WorkspaceEntry) {
	sort.Slice(entries, func(i, j int) bool {
		left, right := entries[i], entries[j]
		if left.Type == app.TargetDirectory && right.Type != app.TargetDirectory {
			return true
		}
		if left.Type != app.TargetDirectory && right.Type == app.TargetDirectory {
			return false
		}
		return left.RelativePath < right.RelativePath
	})
}
