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
	lines := bytes.Split(out, []byte{'\n'})
	entries := make([]StatusEntry, 0, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		if len(line) < 4 || line[2] != ' ' {
			return nil, fmt.Errorf("malformed chezmoi status line %d: %q", i+1, line)
		}
		entries = append(entries, StatusEntry{Code: string(line[:2]), Path: string(line[3:])})
	}
	return entries, nil
}

func ParseManagedFiles(out []byte) []string {
	lines := bytes.Split(out, []byte{'\n'})
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			files = append(files, string(line))
		}
	}
	return files
}
