package chezmoi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TargetMetadata describes the target-state properties needed to gate sync actions.
type TargetMetadata struct {
	Type     string
	Template bool
}

// AuthoritativeDiff returns chezmoi's builtin target-to-destination diff.
func (c Client) AuthoritativeDiff(target string, limit int64) ([]byte, error) {
	return c.OutputLimit(
		limit,
		"--color=false",
		"--no-pager",
		"--use-builtin-diff",
		"diff",
		"--include=all",
		"--exclude=none",
		"--reverse",
		"--script-contents=true",
		target,
	)
}

// TargetMetadata loads one target's type and template state from chezmoi.
func (c Client) TargetMetadata(target string, limit int64) (TargetMetadata, error) {
	out, err := c.OutputLimit(limit, "dump", "--include=all", "--exclude=none", "--format=json", "--recursive=false", target)
	if err != nil {
		return TargetMetadata{}, err
	}
	if int64(len(out)) > limit {
		return TargetMetadata{}, fmt.Errorf("target metadata exceeds %d bytes for %s", limit, target)
	}
	targetType, err := ParseDumpTargetType(out)
	if err != nil {
		return TargetMetadata{}, fmt.Errorf("parse target metadata for %s: %w", target, err)
	}
	metadata := TargetMetadata{Type: targetType}
	if targetType != "file" && targetType != "symlink" {
		return metadata, nil
	}

	templateOut, err := c.Output(
		"managed",
		"--include=templates",
		"--exclude=none",
		"--path-style=absolute",
		"--nul-path-separator",
		target,
	)
	if err != nil {
		return TargetMetadata{}, err
	}
	for _, path := range ParseNULPaths(templateOut) {
		if path == target {
			metadata.Template = true
			break
		}
	}
	return metadata, nil
}

// TargetDir returns chezmoi's configured destination directory.
func (c Client) TargetDir() (string, error) {
	out, err := c.Output("target-path")
	if err != nil {
		return "", err
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("chezmoi target path is empty")
	}
	return dir, nil
}

// ParseDumpTargetType parses a single-entry chezmoi dump document.
func ParseDumpTargetType(out []byte) (string, error) {
	var entries map[string]struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(out, &entries); err != nil {
		return "", err
	}
	if len(entries) != 1 {
		return "", fmt.Errorf("expected one target entry, got %d", len(entries))
	}
	for _, entry := range entries {
		if entry.Type == "" {
			return "", fmt.Errorf("target type is empty")
		}
		return entry.Type, nil
	}
	return "", fmt.Errorf("target metadata contained no entries")
}
