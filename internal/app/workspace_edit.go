package app

import (
	"fmt"

	"github.com/zhongyangchuwu/cm/internal/process"
)

// SourceEditUnavailableReason explains why a workspace entry cannot be edited
// through chezmoi source state. An empty result means source editing is allowed.
func SourceEditUnavailableReason(entry WorkspaceEntry) string {
	switch entry.Type {
	case TargetDirectory:
		return "select a managed file or symlink"
	case TargetScript:
		return "scripts cannot be edited through workspace source edit"
	case TargetRemove, TargetExternal, TargetUnknown:
		if entry.State != FileUnmanaged && entry.State != FileIgnored {
			return fmt.Sprintf("source edit unsupported for %s target", entry.Type)
		}
	}
	if entry.SourcePath == "" {
		switch entry.State {
		case FileUnmanaged:
			return "unmanaged target has no chezmoi source"
		case FileIgnored:
			return "ignored target has no chezmoi source"
		default:
			return "target has no chezmoi source mapping"
		}
	}
	if entry.State == FileScript {
		return "scripts cannot be edited through workspace source edit"
	}
	if entry.Type == TargetFile || entry.Type == TargetSymlink {
		return ""
	}
	return fmt.Sprintf("source edit unsupported for %s target", entry.Type)
}

// SourceEditCommand opens an eligible workspace entry in chezmoi's configured
// editor without allowing inherited edit configuration to apply or watch changes.
func (s service) SourceEditCommand(entry WorkspaceEntry) (TerminalCommand, error) {
	if reason := SourceEditUnavailableReason(entry); reason != "" {
		return nil, fmt.Errorf("source edit unavailable: %s", reason)
	}
	return &runnerCommand{
		runner:  s.runner(),
		command: s.client.BinaryName(),
		args:    []string{"edit", "--apply=false", "--watch=false", entry.Path},
		io:      process.IO{Dir: s.client.Dir},
	}, nil
}
