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
	m.rebuildWorkspaceIndex()
	if m.searchResults {
		m.entries = m.searchEntries()
	} else {
		m.entries = m.directoryEntries(m.currentDir)
	}
	m.cursor = m.cursorFor(m.currentDir, preserveTarget)
	m.resetPreviewPosition()
}

func (m *workspaceModel) rebuildWorkspaceIndex() {
	nodes := make(map[string]app.WorkspaceEntry, len(m.allEntries))
	aggregateStates := make(map[string]app.FileState)
	matchingNodes := make(map[string]bool)
	for _, original := range m.allEntries {
		relative := cleanWorkspaceRelative(original.RelativePath)
		if relative == "" {
			continue
		}
		entry := original
		entry.RelativePath = relative
		nodes[relative] = entry
		for ancestor := workspaceParent(relative); ancestor != ""; ancestor = workspaceParent(ancestor) {
			aggregateStates[ancestor] = aggregateDirectoryState(aggregateStates[ancestor], entry.State)
		}
		if workspaceEntryMatchesFilter(entry, m.filter) {
			for matched := relative; matched != ""; matched = workspaceParent(matched) {
				matchingNodes[matched] = true
			}
		}
	}
	for relative, state := range aggregateStates {
		if entry, exists := nodes[relative]; exists {
			entry.Type = app.TargetDirectory
			entry.State = aggregateDirectoryState(entry.State, state)
			nodes[relative] = entry
			continue
		}
		nodes[relative] = app.WorkspaceEntry{
			Path:         filepath.Join(m.snapshot.Root, filepath.FromSlash(relative)),
			RelativePath: relative,
			State:        state,
			Type:         app.TargetDirectory,
		}
	}
	m.workspaceNodes = nodes
	m.matchingNodes = matchingNodes
}

func (m workspaceModel) directoryEntries(directory string) []app.WorkspaceEntry {
	entries := make([]app.WorkspaceEntry, 0)
	for relative, entry := range m.workspaceNodes {
		if workspaceParent(relative) != directory || !m.matchingNodes[relative] {
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
	parent := workspaceParent(m.currentDir)
	entries := make([]app.WorkspaceEntry, 0)
	for relative, entry := range m.workspaceNodes {
		if workspaceParent(relative) == parent && (m.matchingNodes[relative] || relative == m.currentDir) {
			entries = append(entries, entry)
		}
	}
	sortWorkspaceEntries(entries)
	return entries
}

func (m workspaceModel) searchEntries() []app.WorkspaceEntry {
	query := strings.ToLower(strings.TrimSpace(m.fileQuery))
	if query == "" {
		return nil
	}
	entries := make([]app.WorkspaceEntry, 0)
	for relative, entry := range m.workspaceNodes {
		if !m.matchingNodes[relative] {
			continue
		}
		if strings.Contains(strings.ToLower(entry.RelativePath+"\x00"+entry.SourcePath), query) {
			entries = append(entries, entry)
		}
	}
	sortWorkspaceEntries(entries)
	return entries
}

func (m *workspaceModel) cursorFor(directory, target string) int {
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
	if m.searchResults {
		return 0
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
	if len(snapshot.Scopes) == 0 {
		return ""
	}
	directory := workspaceScopeDirectory(snapshot, snapshot.Scopes[0])
	for _, scope := range snapshot.Scopes[1:] {
		directory = commonWorkspaceDirectory(directory, workspaceScopeDirectory(snapshot, scope))
		if directory == "" {
			return ""
		}
	}
	return directory
}

func workspaceScopeDirectory(snapshot app.WorkspaceSnapshot, scope string) string {
	relative, err := filepath.Rel(snapshot.Root, scope)
	if err != nil || relative == "." || strings.HasPrefix(relative, "..") {
		return ""
	}
	relative = filepath.ToSlash(relative)
	exactNonDirectory := false
	for _, entry := range snapshot.Entries {
		entryRelative := cleanWorkspaceRelative(entry.RelativePath)
		if strings.HasPrefix(entryRelative, relative+"/") {
			return relative
		}
		if entryRelative == relative && entry.Type != app.TargetDirectory {
			exactNonDirectory = true
		}
	}
	if exactNonDirectory {
		return workspaceParent(relative)
	}
	return relative
}

func commonWorkspaceDirectory(left, right string) string {
	leftParts := splitWorkspacePath(left)
	rightParts := splitWorkspacePath(right)
	limit := min(len(leftParts), len(rightParts))
	common := leftParts[:0]
	for index := 0; index < limit && leftParts[index] == rightParts[index]; index++ {
		common = append(common, leftParts[index])
	}
	return strings.Join(common, "/")
}

func splitWorkspacePath(value string) []string {
	value = cleanWorkspaceRelative(value)
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func aggregateDirectoryState(current, candidate app.FileState) app.FileState {
	if directoryStateRank(candidate) > directoryStateRank(current) {
		return candidate
	}
	if current == "" {
		return app.FileClean
	}
	return current
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
	case app.FileClean:
		return 1
	default:
		return 0
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
