package chezmoi

import (
	"bytes"
	"fmt"
)

type StatusEntry struct {
	Code string
	Path string
}

func ParseStatus(out []byte) ([]StatusEntry, error) {
	return parseStatusLines(out, "chezmoi")
}

func ParseGitStatus(out []byte) ([]StatusEntry, error) {
	return parseStatusLines(out, "git")
}

func parseStatusLines(out []byte, source string) ([]StatusEntry, error) {
	lines := bytes.Split(out, []byte{'\n'})
	entries := make([]StatusEntry, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		if len(line) < 4 || line[2] != ' ' {
			return nil, fmt.Errorf("malformed %s status line %d: %q", source, i+1, line)
		}
		entries = append(entries, StatusEntry{Code: string(line[:2]), Path: string(line[3:])})
	}
	return entries, nil
}

func ParseNULPaths(out []byte) []string {
	paths := bytes.Split(out, []byte{0})
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if len(path) > 0 {
			result = append(result, string(path))
		}
	}
	return result
}
