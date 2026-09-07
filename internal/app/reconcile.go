package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

// ReconcileEntry is one non-script chezmoi status entry eligible for review.
type ReconcileEntry struct {
	Code string
	Path string
}

// SyncStatus separates file reconciliation from pending chezmoi scripts.
type SyncStatus struct {
	Entries []ReconcileEntry
	Scripts []ReconcileEntry
}

// TargetType is the rendered target-state type reported by chezmoi.
type TargetType string

const (
	TargetUnknown   TargetType = "unknown"
	TargetFile      TargetType = "file"
	TargetDirectory TargetType = "dir"
	TargetSymlink   TargetType = "symlink"
	TargetRemove    TargetType = "remove"
	TargetScript    TargetType = "script"
	TargetExternal  TargetType = "external"
)

// Review binds target metadata and an authoritative diff to one fingerprint.
type Review struct {
	Entry       ReconcileEntry
	Type        TargetType
	Template    bool
	Diff        string
	Fingerprint string
	Dirty       bool
}

// ActionResult preserves output from a successful non-interactive action.
type ActionResult struct {
	Stdout string
	Stderr string
}

func (r ActionResult) Notice() string {
	stdout := strings.TrimSpace(r.Stdout)
	stderr := strings.TrimSpace(r.Stderr)
	switch {
	case stdout != "" && stderr != "":
		return stdout + "\n" + stderr
	case stdout != "":
		return stdout
	default:
		return stderr
	}
}

func partitionSyncStatus(entries []chezmoi.StatusEntry) SyncStatus {
	status := SyncStatus{
		Entries: make([]ReconcileEntry, 0, len(entries)),
		Scripts: make([]ReconcileEntry, 0),
	}
	for _, entry := range entries {
		reconcileEntry := ReconcileEntry{Code: entry.Code, Path: entry.Path}
		if len(entry.Code) < 2 || entry.Code[1] == ' ' {
			continue
		}
		if entry.Code[1] == 'R' {
			status.Scripts = append(status.Scripts, reconcileEntry)
			continue
		}
		status.Entries = append(status.Entries, reconcileEntry)
	}
	return status
}
func findReconcileEntry(entries []ReconcileEntry, target string) (ReconcileEntry, bool) {
	for _, entry := range entries {
		if entry.Path == target {
			return entry, true
		}
	}
	return ReconcileEntry{}, false
}

func newReview(entry ReconcileEntry, targetType TargetType, template bool, diff string) Review {
	review := Review{
		Entry:    entry,
		Type:     targetType,
		Template: template,
		Diff:     diff,
		Dirty:    true,
	}
	review.Fingerprint = reviewFingerprint(review)
	return review
}

func cleanReview(target string) Review {
	review := Review{Entry: ReconcileEntry{Path: target}, Type: TargetUnknown}
	review.Fingerprint = reviewFingerprint(review)
	return review
}

func reviewFingerprint(review Review) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(review.Entry.Path))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(review.Entry.Code))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(review.Type))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(strconv.FormatBool(review.Template)))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write([]byte(review.Diff))
	return hex.EncodeToString(hash.Sum(nil))
}

// Allows reports whether the reviewed target supports an action.
func (r Review) Allows(kind ActionKind) bool {
	if !r.Dirty {
		return false
	}
	switch r.Type {
	case TargetFile:
		if r.Template && kind == ActionAdd {
			return false
		}
		return kind == ActionAdd || kind == ActionApply || kind == ActionMerge
	case TargetDirectory, TargetSymlink, TargetRemove:
		return kind == ActionApply
	default:
		return false
	}
}

// Action creates a pending action bound to this reviewed state.
func (r Review) Action(kind ActionKind) (Action, error) {
	if !r.Allows(kind) {
		return Action{}, fmt.Errorf("%s is not available for %s target %s", kind, r.Type, r.Entry.Path)
	}
	return Action{Target: r.Entry.Path, Kind: kind, Fingerprint: r.Fingerprint}, nil
}

func targetTypeFromChezmoi(value string) TargetType {
	switch value {
	case string(TargetFile):
		return TargetFile
	case string(TargetDirectory):
		return TargetDirectory
	case string(TargetSymlink):
		return TargetSymlink
	case string(TargetRemove):
		return TargetRemove
	case string(TargetScript):
		return TargetScript
	case string(TargetExternal):
		return TargetExternal
	default:
		return TargetUnknown
	}
}
