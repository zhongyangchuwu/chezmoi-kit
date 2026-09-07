package app

import (
	"reflect"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestPartitionSyncStatusSeparatesScripts(t *testing.T) {
	got := partitionSyncStatus([]chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: " R", Path: "/home/me/install.sh"},
		{Code: " D", Path: "/home/me/.old"},
		{Code: "M ", Path: "/home/me/already-matches-target"},
	})
	wantEntries := []ReconcileEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: " D", Path: "/home/me/.old"},
	}
	wantScripts := []ReconcileEntry{{Code: " R", Path: "/home/me/install.sh"}}
	if !reflect.DeepEqual(got.Entries, wantEntries) || !reflect.DeepEqual(got.Scripts, wantScripts) {
		t.Fatalf("partitionSyncStatus() = %#v", got)
	}
}

func TestReviewAllowsOnlySafeActions(t *testing.T) {
	tests := []struct {
		name    string
		review  Review
		allowed []ActionKind
		blocked []ActionKind
	}{
		{
			name:    "regular file",
			review:  newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, false, "diff"),
			allowed: []ActionKind{ActionAdd, ActionApply, ActionMerge},
		},
		{
			name:    "template file",
			review:  newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, true, "diff"),
			allowed: []ActionKind{ActionApply, ActionMerge},
			blocked: []ActionKind{ActionAdd},
		},
		{
			name:    "symlink",
			review:  newReview(ReconcileEntry{Code: " M", Path: "/link"}, TargetSymlink, false, "diff"),
			allowed: []ActionKind{ActionApply},
			blocked: []ActionKind{ActionAdd, ActionMerge},
		},
		{
			name:    "directory",
			review:  newReview(ReconcileEntry{Code: " M", Path: "/dir"}, TargetDirectory, false, "diff"),
			allowed: []ActionKind{ActionApply},
			blocked: []ActionKind{ActionAdd, ActionMerge},
		},
		{
			name:    "remove",
			review:  newReview(ReconcileEntry{Code: " D", Path: "/old"}, TargetRemove, false, "diff"),
			allowed: []ActionKind{ActionApply},
			blocked: []ActionKind{ActionAdd, ActionMerge},
		},
		{
			name:    "unknown",
			review:  newReview(ReconcileEntry{Code: "MM", Path: "/unknown"}, TargetUnknown, false, "diff"),
			blocked: []ActionKind{ActionAdd, ActionApply, ActionMerge},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, kind := range tt.allowed {
				if !tt.review.Allows(kind) {
					t.Fatalf("Allows(%s) = false", kind)
				}
			}
			for _, kind := range tt.blocked {
				if tt.review.Allows(kind) {
					t.Fatalf("Allows(%s) = true", kind)
				}
			}
		})
	}
}

func TestReviewFingerprintChangesWithReviewedState(t *testing.T) {
	base := newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, false, "diff")
	tests := []Review{
		newReview(ReconcileEntry{Code: " M", Path: "/file"}, TargetFile, false, "diff"),
		newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetSymlink, false, "diff"),
		newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, true, "diff"),
		newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, false, "changed diff"),
	}
	for _, changed := range tests {
		if changed.Fingerprint == base.Fingerprint {
			t.Fatalf("fingerprint did not change: base=%#v changed=%#v", base, changed)
		}
	}
}

func TestReviewActionCarriesFingerprint(t *testing.T) {
	review := newReview(ReconcileEntry{Code: "MM", Path: "/file"}, TargetFile, false, "diff")
	action, err := review.Action(ActionAdd)
	if err != nil {
		t.Fatalf("Action() error = %v", err)
	}
	if action.Target != "/file" || action.Kind != ActionAdd || action.Fingerprint != review.Fingerprint {
		t.Fatalf("Action() = %#v", action)
	}
}
