package reconcile

import "github.com/zhongyangchuwu/cm/internal/chezmoi"

// Service is the boundary for reconciling managed local files with chezmoi state.
// UI packages own presentation; this package owns the shared reconciliation contract.
type Service interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	Add(target string) error
	Apply(target string) error
	Merge(target string) error
}
