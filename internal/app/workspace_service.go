package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	maxInventoryOutput  int64 = 4 << 20
	maxWorkspaceEntries       = 10_000
	maxUnmanagedQueries       = 2_000
)

type managedAttributes struct {
	directories map[string]bool
	symlinks    map[string]bool
	scripts     map[string]bool
	removes     map[string]bool
	externals   map[string]bool
	templates   map[string]bool
	encrypted   map[string]bool
}

func (s service) Inventory(scopes []string) (WorkspaceSnapshot, error) {
	root, err := s.client.TargetDir()
	if err != nil {
		return WorkspaceSnapshot{}, err
	}
	normalizedScopes, err := normalizeWorkspaceScopes(root, scopes)
	if err != nil {
		return WorkspaceSnapshot{}, err
	}
	managed, err := s.client.ManagedInventory(normalizedScopes, maxInventoryOutput)
	if err != nil {
		return WorkspaceSnapshot{}, err
	}
	statusEntries, err := s.client.StatusForInventory(normalizedScopes)
	if err != nil {
		return WorkspaceSnapshot{}, err
	}
	attributes, err := s.managedAttributes(normalizedScopes)
	if err != nil {
		return WorkspaceSnapshot{}, err
	}
	ignored, err := s.client.IgnoredEntries(maxInventoryOutput)
	if err != nil {
		return WorkspaceSnapshot{}, err
	}

	entries := make(map[string]WorkspaceEntry, len(managed)+len(ignored))
	statusByPath := make(map[string]string, len(statusEntries))
	for _, statusEntry := range statusEntries {
		statusByPath[statusEntry.Path] = statusEntry.Code
	}
	for _, managedEntry := range managed {
		code := statusByPath[managedEntry.Absolute]
		entries[managedEntry.Absolute] = WorkspaceEntry{
			Path:         managedEntry.Absolute,
			RelativePath: managedEntry.Relative,
			SourcePath:   managedEntry.SourceAbsolute,
			State:        managedState(code, attributes.scripts[managedEntry.Absolute], attributes.templates[managedEntry.Absolute], attributes.encrypted[managedEntry.Absolute]),
			Type:         managedTargetType(managedEntry.Absolute, attributes),
			Template:     attributes.templates[managedEntry.Absolute],
			Encrypted:    attributes.encrypted[managedEntry.Absolute],
			Code:         code,
		}
	}

	for _, relative := range ignored {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if !pathWithinScopes(path, normalizedScopes) {
			continue
		}
		if _, exists := entries[path]; exists {
			continue
		}
		entries[path] = WorkspaceEntry{
			Path:         path,
			RelativePath: filepath.ToSlash(relative),
			State:        FileIgnored,
			Type:         localTargetType(path),
		}
	}

	if len(normalizedScopes) > 0 {
		unmanaged, err := s.discoverUnmanaged(root, normalizedScopes)
		if err != nil {
			return WorkspaceSnapshot{}, err
		}
		for _, entry := range unmanaged {
			if _, exists := entries[entry.Path]; !exists {
				entries[entry.Path] = entry
			}
		}
	}

	snapshot := WorkspaceSnapshot{
		Root:   root,
		Scopes: append([]string(nil), normalizedScopes...),
	}
	if len(normalizedScopes) == 0 {
		snapshot.Notice = "unmanaged discovery disabled; pass one or more paths to cm ui"
	}
	snapshot.Entries = make([]WorkspaceEntry, 0, len(entries))
	for _, entry := range entries {
		snapshot.Entries = append(snapshot.Entries, entry)
	}
	sort.Slice(snapshot.Entries, func(i, j int) bool {
		if snapshot.Entries[i].RelativePath == snapshot.Entries[j].RelativePath {
			return snapshot.Entries[i].State < snapshot.Entries[j].State
		}
		return snapshot.Entries[i].RelativePath < snapshot.Entries[j].RelativePath
	})
	return snapshot, nil
}

func (s service) managedAttributes(scopes []string) (managedAttributes, error) {
	attributes := managedAttributes{}
	queries := []struct {
		entryType string
		dest      *map[string]bool
	}{
		{entryType: "dirs", dest: &attributes.directories},
		{entryType: "symlinks", dest: &attributes.symlinks},
		{entryType: "scripts", dest: &attributes.scripts},
		{entryType: "remove", dest: &attributes.removes},
		{entryType: "externals", dest: &attributes.externals},
		{entryType: "templates", dest: &attributes.templates},
		{entryType: "encrypted", dest: &attributes.encrypted},
	}
	for _, query := range queries {
		paths, err := s.client.ManagedPathsByType(query.entryType, scopes, maxInventoryOutput)
		if err != nil {
			return managedAttributes{}, err
		}
		set := make(map[string]bool, len(paths))
		for _, path := range paths {
			set[path] = true
		}
		*query.dest = set
	}
	return attributes, nil
}

func (s service) discoverUnmanaged(root string, scopes []string) ([]WorkspaceEntry, error) {
	initial, err := s.client.UnmanagedEntries(scopes, maxInventoryOutput)
	if err != nil {
		return nil, err
	}
	queue := append([]string(nil), initial...)
	seen := make(map[string]bool, len(queue))
	entries := make([]WorkspaceEntry, 0, len(queue))
	queries := 0
	for len(queue) > 0 {
		path := filepath.Clean(queue[0])
		queue = queue[1:]
		if seen[path] {
			continue
		}
		seen[path] = true
		if len(entries) >= maxWorkspaceEntries {
			return nil, fmt.Errorf("unmanaged discovery exceeds %d entries; narrow the cm ui path", maxWorkspaceEntries)
		}
		entryType := localTargetType(path)
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil || relative == "." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || relative == ".." {
			return nil, fmt.Errorf("unmanaged path %s is outside destination %s", path, root)
		}
		entries = append(entries, WorkspaceEntry{
			Path:         path,
			RelativePath: filepath.ToSlash(relative),
			State:        FileUnmanaged,
			Type:         entryType,
		})
		if entryType != TargetDirectory {
			continue
		}
		children, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("read unmanaged directory %s: %w", path, err)
		}
		if len(children) == 0 {
			continue
		}
		if len(children) > maxWorkspaceEntries-len(seen) {
			return nil, fmt.Errorf("unmanaged directory %s contains too many entries; narrow the cm ui path", path)
		}
		const batchSize = 256
		for start := 0; start < len(children); start += batchSize {
			queries++
			if queries > maxUnmanagedQueries {
				return nil, fmt.Errorf("unmanaged discovery exceeds %d queries; narrow the cm ui path", maxUnmanagedQueries)
			}
			end := min(start+batchSize, len(children))
			childPaths := make([]string, 0, end-start)
			for _, child := range children[start:end] {
				childPaths = append(childPaths, filepath.Join(path, child.Name()))
			}
			unmanagedChildren, err := s.client.UnmanagedEntries(childPaths, maxInventoryOutput)
			if err != nil {
				return nil, err
			}
			for _, child := range unmanagedChildren {
				if filepath.Clean(child) != path && !seen[filepath.Clean(child)] {
					queue = append(queue, child)
				}
			}
		}
	}
	return entries, nil
}

func (s service) Preview(entry WorkspaceEntry, kind PreviewKind, reveal bool) (WorkspacePreview, error) {
	preview := WorkspacePreview{Entry: entry, Kind: kind}
	switch kind {
	case PreviewDiff:
		if entry.State == FileUnmanaged || entry.State == FileIgnored {
			preview.Notice = "no target diff for " + string(entry.State) + " entry"
			return preview, nil
		}
		if entry.State == FileScript {
			preview.Notice = "script execution is outside ordinary file reconciliation; use source or target view"
			return preview, nil
		}
		if entry.State == FileUninspected && !reveal {
			preview.Withheld = true
			preview.Notice = "authoritative diff withheld; press R to inspect this sensitive target"
			return preview, nil
		}
		if entry.State == FileDirty {
			review, err := s.Review(entry.Path)
			if err != nil {
				return WorkspacePreview{}, err
			}
			preview.Review = review
			preview.Content = review.Diff
			return preview, nil
		}
		out, err := s.authoritativeDiff(entry.Path)
		if err != nil {
			return WorkspacePreview{}, err
		}
		if len(out) == 0 {
			preview.Notice = "destination matches rendered target"
			if entry.State == FileUninspected {
				preview.Entry.State = FileClean
			}
		} else {
			preview.Content = string(out)
			if entry.State == FileUninspected {
				preview.Entry.State = FileDirty
			}
		}
		return preview, nil
	case PreviewDestination:
		content, err := readWorkspacePath(entry.Path, maxReviewOutput)
		if err != nil {
			if os.IsNotExist(err) {
				preview.Notice = "destination does not exist"
				return preview, nil
			}
			return WorkspacePreview{}, fmt.Errorf("read destination %s: %w", entry.Path, err)
		}
		preview.Content = content
		return preview, nil
	case PreviewTarget:
		if entry.State == FileUnmanaged || entry.State == FileIgnored {
			preview.Notice = "entry has no rendered target"
			return preview, nil
		}
		switch entry.Type {
		case TargetDirectory:
			preview.Notice = "directory target has no content view; use diff for target metadata"
			return preview, nil
		case TargetRemove:
			preview.Notice = "target state removes the destination entry"
			return preview, nil
		case TargetExternal:
			preview.Notice = "external target content is not loaded by cm ui"
			return preview, nil
		}
		if (entry.Template || entry.Encrypted) && !reveal {
			preview.Withheld = true
			preview.Notice = "rendered target withheld; press R to reveal this target"
			return preview, nil
		}
		out, err := s.client.TargetContent(entry.Path, false, maxReviewOutput)
		if err != nil {
			return WorkspacePreview{}, fmt.Errorf("read rendered target %s: %w", entry.Path, err)
		}
		preview.Content = string(out)
		return preview, nil
	case PreviewSource:
		if entry.SourcePath == "" {
			preview.Notice = "source path is unavailable for this entry"
			return preview, nil
		}
		if entry.Encrypted && !reveal {
			preview.Withheld = true
			preview.Notice = "encrypted source withheld; press R to decrypt and reveal this target"
			return preview, nil
		}
		if entry.Encrypted {
			out, err := s.client.DecryptedSourceContent(entry.SourcePath, maxReviewOutput)
			if err != nil {
				return WorkspacePreview{}, fmt.Errorf("decrypt source %s: %w", entry.SourcePath, err)
			}
			preview.Content = string(out)
			return preview, nil
		}
		content, err := readWorkspacePath(entry.SourcePath, maxReviewOutput)
		if err != nil {
			return WorkspacePreview{}, fmt.Errorf("read source %s: %w", entry.SourcePath, err)
		}
		preview.Content = content
		return preview, nil
	default:
		return WorkspacePreview{}, fmt.Errorf("unsupported preview kind %q", kind)
	}
}

func normalizeWorkspaceScopes(root string, scopes []string) ([]string, error) {
	if len(scopes) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(scopes))
	seen := make(map[string]bool, len(scopes))
	for _, scope := range scopes {
		if !filepath.IsAbs(scope) {
			scope = filepath.Join(root, scope)
		}
		scope = filepath.Clean(scope)
		relative, err := filepath.Rel(root, scope)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("workspace scope %s is outside destination %s", scope, root)
		}
		if !seen[scope] {
			seen[scope] = true
			result = append(result, scope)
		}
	}
	return result, nil
}

func pathWithinScopes(path string, scopes []string) bool {
	if len(scopes) == 0 {
		return true
	}
	for _, scope := range scopes {
		relative, err := filepath.Rel(scope, path)
		if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func managedState(code string, script bool, template bool, encrypted bool) FileState {
	if script || (len(code) >= 2 && code[1] == 'R') {
		return FileScript
	}
	if encrypted {
		return FileUninspected
	}
	if len(code) >= 2 && code[1] != ' ' {
		return FileDirty
	}
	if template && code == "" {
		return FileUninspected
	}
	return FileClean
}

func managedTargetType(path string, attributes managedAttributes) TargetType {
	switch {
	case attributes.directories[path]:
		return TargetDirectory
	case attributes.symlinks[path]:
		return TargetSymlink
	case attributes.scripts[path]:
		return TargetScript
	case attributes.removes[path]:
		return TargetRemove
	case attributes.externals[path]:
		return TargetExternal
	default:
		return TargetFile
	}
}

func localTargetType(path string) TargetType {
	info, err := os.Lstat(path)
	if err != nil {
		return TargetUnknown
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return TargetSymlink
	case info.IsDir():
		return TargetDirectory
	default:
		return TargetFile
	}
}

func readWorkspacePath(path string, limit int64) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		return link + "\n", nil
	}
	if info.IsDir() {
		return fmt.Sprintf("directory\nmode: %04o\n", info.Mode().Perm()), nil
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return "", err
	}
	if int64(len(content)) > limit {
		return "", fmt.Errorf("content exceeds %d bytes", limit)
	}
	return string(content), nil
}
