package reconcile

import (
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestRecommendMapsStatusToDefaultAction(t *testing.T) {
	tests := []struct {
		name  string
		entry chezmoi.StatusEntry
		want  Action
	}{
		{
			name:  "local modified recommends add",
			entry: entry(chezmoi.ChangeModified, chezmoi.ChangeNone),
			want:  ActionAdd,
		},
		{
			name:  "apply would modify recommends apply",
			entry: entry(chezmoi.ChangeNone, chezmoi.ChangeModified),
			want:  ActionApply,
		},
		{
			name:  "local drift plus apply pending recommends inspect",
			entry: entry(chezmoi.ChangeModified, chezmoi.ChangeModified),
			want:  ActionInspect,
		},
		{
			name:  "local deleted recommends inspect",
			entry: entry(chezmoi.ChangeDeleted, chezmoi.ChangeNone),
			want:  ActionInspect,
		},
		{
			name:  "source deleted recommends apply",
			entry: entry(chezmoi.ChangeNone, chezmoi.ChangeDeleted),
			want:  ActionApply,
		},
		{
			name:  "local added recommends add",
			entry: entry(chezmoi.ChangeAdded, chezmoi.ChangeNone),
			want:  ActionAdd,
		},
		{
			name:  "source added recommends apply",
			entry: entry(chezmoi.ChangeNone, chezmoi.ChangeAdded),
			want:  ActionApply,
		},
		{
			name:  "unknown recommends diff",
			entry: entry(chezmoi.ChangeRun, chezmoi.ChangeNone),
			want:  ActionDiff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Recommend(tt.entry); got != tt.want {
				t.Fatalf("Recommend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescribeExplainsCommonStatus(t *testing.T) {
	tests := []struct {
		name  string
		entry chezmoi.StatusEntry
		want  string
	}{
		{"local", entry(chezmoi.ChangeModified, chezmoi.ChangeNone), "local drift"},
		{"apply", entry(chezmoi.ChangeNone, chezmoi.ChangeModified), "apply pending"},
		{"both", entry(chezmoi.ChangeModified, chezmoi.ChangeModified), "local drift, apply pending"},
		{"local deleted", entry(chezmoi.ChangeDeleted, chezmoi.ChangeNone), "local deleted"},
		{"apply delete", entry(chezmoi.ChangeNone, chezmoi.ChangeDeleted), "apply would delete"},
		{"unknown", entry(chezmoi.ChangeRun, chezmoi.ChangeNone), "inspect"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Describe(tt.entry); got != tt.want {
				t.Fatalf("Describe() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActionName(t *testing.T) {
	if got := ActionName(ActionMerge); got != "merge" {
		t.Fatalf("ActionName(ActionMerge) = %q, want merge", got)
	}
}

func entry(local, target chezmoi.Change) chezmoi.StatusEntry {
	return chezmoi.StatusEntry{LocalChange: local, TargetChange: target, Path: ".zshrc"}
}
