package tui

import (
	"io"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

// ActionKind identifies a reconciliation action selected during review.
type ActionKind int

const (
	ActionAdd ActionKind = iota
	ActionApply
	ActionMerge
)

func (k ActionKind) String() string {
	switch k {
	case ActionAdd:
		return "add"
	case ActionApply:
		return "apply"
	case ActionMerge:
		return "merge"
	default:
		return "unknown"
	}
}

func (k ActionKind) Marker() string {
	switch k {
	case ActionAdd:
		return "A"
	case ActionApply:
		return "P"
	case ActionMerge:
		return "M"
	default:
		return "?"
	}
}

// Action is a confirmed reconciliation action for one managed target.
type Action struct {
	Target string
	Kind   ActionKind
}

// TerminalCommand is a reconciliation command that temporarily owns the terminal.
type TerminalCommand interface {
	Run() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

// ReviewService is the boundary for reviewing and executing sync actions.
type ReviewService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	ExecuteNonInteractive(action Action) error
	TerminalCommand(action Action) (TerminalCommand, error)
}
