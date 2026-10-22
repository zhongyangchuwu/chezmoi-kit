package app

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/zhongyangchuwu/cm/internal/process"
)

// DestinationEditUnavailableReason explains why a workspace entry cannot be
// opened in the local destination editor. An empty result means it is allowed.
func DestinationEditUnavailableReason(entry WorkspaceEntry) string {
	if entry.Type == TargetSymlink {
		return "destination edit does not follow symlinks; select a regular file"
	}
	if entry.Type != TargetFile {
		return "select a managed regular file"
	}
	reason := SourceEditUnavailableReason(entry)
	return strings.ReplaceAll(reason, "source edit", "destination edit")
}

// MergeUnavailableReason explains why a workspace entry cannot be merged.
func MergeUnavailableReason(entry WorkspaceEntry) string {
	if entry.Type == TargetScript || entry.State == FileScript {
		return "scripts cannot be merged"
	}
	if entry.Type == TargetSymlink {
		return "merge does not follow symlinks; select a regular file"
	}
	if entry.Type != TargetFile {
		return "select a managed regular file"
	}
	reason := SourceEditUnavailableReason(entry)
	return strings.ReplaceAll(reason, "source edit", "merge")
}

// DestinationEditCommand opens an eligible destination file through VISUAL or
// EDITOR. It never invokes chezmoi apply, edit, or watch behavior.
func (s service) DestinationEditCommand(entry WorkspaceEntry) (TerminalCommand, error) {
	if reason := DestinationEditUnavailableReason(entry); reason != "" {
		return nil, fmt.Errorf("destination edit unavailable: %s", reason)
	}
	command, args, err := destinationEditor()
	if err != nil {
		return nil, err
	}
	return &runnerCommand{
		runner:  s.runner(),
		command: command,
		args:    append(args, entry.Path),
		io:      process.IO{Dir: s.client.Dir},
	}, nil
}

// MergeCommand opens chezmoi's configured three-way merge tool for one entry.
func (s service) MergeCommand(entry WorkspaceEntry) (TerminalCommand, error) {
	if reason := MergeUnavailableReason(entry); reason != "" {
		return nil, fmt.Errorf("merge unavailable: %s", reason)
	}
	return s.mergeTerminalCommand(entry.Path), nil
}

func destinationEditor() (string, []string, error) {
	for _, name := range []string{"VISUAL", "EDITOR"} {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			parts, err := splitEditorCommand(value)
			if err != nil {
				return "", nil, fmt.Errorf("parse %s: %w", name, err)
			}
			return parts[0], parts[1:], nil
		}
	}
	return "", nil, fmt.Errorf("destination editor is not configured; set VISUAL or EDITOR")
}

func splitEditorCommand(value string) ([]string, error) {
	var (
		args    []string
		word    []rune
		quote   rune
		escaped bool
		started bool
	)
	flush := func() {
		if started {
			args = append(args, string(word))
			word = word[:0]
			started = false
		}
	}
	for _, r := range value {
		switch {
		case escaped:
			word = append(word, r)
			escaped = false
			started = true
		case r == '\\' && quote != '\'':
			escaped = true
			started = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				word = append(word, r)
			}
			started = true
		case r == '\'' || r == '"':
			quote = r
			started = true
		case unicode.IsSpace(r):
			flush()
		default:
			word = append(word, r)
			started = true
		}
	}
	if escaped {
		return nil, fmt.Errorf("trailing escape")
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}
	flush()
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	return args, nil
}
