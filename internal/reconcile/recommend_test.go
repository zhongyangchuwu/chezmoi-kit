package reconcile

import (
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestDescribeUsesSingleLocalReconciliationState(t *testing.T) {
	got := Describe(chezmoi.StatusEntry{Code: "MM", Path: "/home/me/.zshrc"})
	if got != "differs from chezmoi" {
		t.Fatalf("Describe() = %q, want differs from chezmoi", got)
	}
}
