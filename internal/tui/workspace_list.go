package tui

import (
	"path"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func (m *workspaceModel) rebuildEntries(preserveTarget string) {
	oldCursor := m.cursor
	if preserveTarget == "" {
		preserveTarget = m.currentTarget()
	}
	entries := make([]app.WorkspaceEntry, 0, len(m.allEntries))
	query := strings.ToLower(strings.TrimSpace(m.fileQuery))
	for _, entry := range m.allEntries {
		if !workspaceEntryMatchesFilter(entry, m.filter) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(entry.RelativePath+"\x00"+entry.SourcePath), query) {
			continue
		}
		if !m.flat && m.hiddenByCollapsedParent(entry.RelativePath) {
			continue
		}
		entries = append(entries, entry)
	}
	m.entries = entries
	m.cursor = 0
	for i, entry := range entries {
		if entry.Path == preserveTarget {
			m.cursor = i
			m.resetPreviewPosition()
			return
		}
	}
	if len(entries) == 0 {
		m.cursor = 0
		m.resetPreviewPosition()
		return
	}
	if oldCursor >= len(entries) {
		oldCursor = len(entries) - 1
	}
	if oldCursor < 0 {
		oldCursor = 0
	}
	m.cursor = oldCursor
	m.resetPreviewPosition()
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

func (m workspaceModel) hiddenByCollapsedParent(relativePath string) bool {
	parent := path.Dir(relativePath)
	for parent != "." && parent != "/" {
		if m.collapsed[parent] {
			return true
		}
		parent = path.Dir(parent)
	}
	return false
}

func (m *workspaceModel) toggleTreeMode() {
	target := m.currentTarget()
	m.flat = !m.flat
	m.rebuildEntries(target)
}

func (m *workspaceModel) toggleCurrentDirectory() {
	entry := m.current()
	if m.flat || entry.Type != app.TargetDirectory {
		return
	}
	m.collapsed[entry.RelativePath] = !m.collapsed[entry.RelativePath]
	m.rebuildEntries(entry.Path)
}

func (m *workspaceModel) cycleFilter() {
	target := m.currentTarget()
	m.filter = (m.filter + 1) % (filterScripts + 1)
	m.rebuildEntries(target)
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

func treeDepth(relativePath string) int {
	relativePath = strings.Trim(relativePath, "/")
	if relativePath == "" {
		return 0
	}
	return strings.Count(relativePath, "/")
}
