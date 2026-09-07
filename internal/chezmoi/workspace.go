package chezmoi

import (
	"encoding/json"
	"fmt"
	"sort"
)

// ManagedEntry is one path mapping returned by `chezmoi managed --path-style=all`.
type ManagedEntry struct {
	Relative       string
	Absolute       string
	SourceAbsolute string
	SourceRelative string
}

type managedPathMapping struct {
	Absolute       string `json:"absolute"`
	SourceAbsolute string `json:"sourceAbsolute"`
	SourceRelative string `json:"sourceRelative"`
}

func (c Client) ManagedInventory(scopes []string, limit int64) ([]ManagedEntry, error) {
	args := []string{"managed", "--include=all", "--exclude=none", "--path-style=all", "--format=json"}
	args = append(args, scopes...)
	out, err := c.boundedOutput(limit, args...)
	if err != nil {
		return nil, err
	}
	entries, err := ParseManagedInventory(out)
	if err != nil {
		return nil, fmt.Errorf("parse chezmoi managed inventory: %w", err)
	}
	return entries, nil
}

func (c Client) ManagedPathsByType(entryType string, scopes []string, limit int64) ([]string, error) {
	args := []string{"managed", "--include=" + entryType, "--exclude=none", "--path-style=absolute", "--nul-path-separator"}
	args = append(args, scopes...)
	out, err := c.boundedOutput(limit, args...)
	if err != nil {
		return nil, err
	}
	return ParseNULPaths(out), nil
}

func (c Client) IgnoredEntries(limit int64) ([]string, error) {
	out, err := c.boundedOutput(limit, "ignored", "--nul-path-separator")
	if err != nil {
		return nil, err
	}
	return ParseNULPaths(out), nil
}

func (c Client) UnmanagedEntries(scopes []string, limit int64) ([]string, error) {
	args := []string{"unmanaged", "--include=all", "--exclude=none", "--path-style=absolute", "--nul-path-separator"}
	args = append(args, scopes...)
	out, err := c.boundedOutput(limit, args...)
	if err != nil {
		return nil, err
	}
	return ParseNULPaths(out), nil
}

func (c Client) TargetContent(target string, skipSecrets bool, limit int64) ([]byte, error) {
	args := make([]string, 0, 3)
	if skipSecrets {
		args = append(args, "--skip-secrets")
	}
	args = append(args, "cat", target)
	return c.boundedOutput(limit, args...)
}

func (c Client) DecryptedSourceContent(sourcePath string, limit int64) ([]byte, error) {
	return c.boundedOutput(limit, "decrypt", sourcePath)
}

func (c Client) boundedOutput(limit int64, args ...string) ([]byte, error) {
	out, err := c.OutputLimit(limit, args...)
	if err != nil {
		return nil, err
	}
	if int64(len(out)) > limit {
		return nil, fmt.Errorf("chezmoi output exceeds %d bytes", limit)
	}
	return out, nil
}

func ParseManagedInventory(out []byte) ([]ManagedEntry, error) {
	var mappings map[string]managedPathMapping
	if err := json.Unmarshal(out, &mappings); err != nil {
		return nil, err
	}
	relativePaths := make([]string, 0, len(mappings))
	for relative := range mappings {
		relativePaths = append(relativePaths, relative)
	}
	sort.Strings(relativePaths)

	entries := make([]ManagedEntry, 0, len(relativePaths))
	for _, relative := range relativePaths {
		mapping := mappings[relative]
		if relative == "" || mapping.Absolute == "" {
			return nil, fmt.Errorf("managed entry %q has no absolute path", relative)
		}
		entries = append(entries, ManagedEntry{
			Relative:       relative,
			Absolute:       mapping.Absolute,
			SourceAbsolute: mapping.SourceAbsolute,
			SourceRelative: mapping.SourceRelative,
		})
	}
	return entries, nil
}
