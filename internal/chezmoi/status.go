package chezmoi

import (
	"bytes"
	"fmt"
)

type Change byte

const (
	ChangeNone     Change = ' '
	ChangeAdded    Change = 'A'
	ChangeDeleted  Change = 'D'
	ChangeModified Change = 'M'
	ChangeRun      Change = 'R'
)

type StatusEntry struct {
	LocalChange  Change
	TargetChange Change
	Path         string
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
		entries = append(entries, StatusEntry{
			LocalChange:  Change(line[0]),
			TargetChange: Change(line[1]),
			Path:         string(line[3:]),
		})
	}
	return entries, nil
}
