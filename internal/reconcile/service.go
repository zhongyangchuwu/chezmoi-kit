package reconcile

import "github.com/zhongyangchuwu/cm/internal/chezmoi"

// ActionKind identifies a reconciliation action selected during review.
type ActionKind int

const (
	ActionAdd ActionKind = iota
	ActionApply
	ActionMerge
)

// Action is a confirmed reconciliation action for one managed target.
type Action struct {
	Target string
	Kind   ActionKind
}

// ReviewService is the boundary for reviewing and executing reconciliation actions.
// UI packages own presentation; this package owns the shared reconciliation contract.
type ReviewService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	Execute(actions []Action) error
}
