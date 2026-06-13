package reconcile

import "github.com/zhongyangchuwu/cm/internal/chezmoi"

type Action int

const (
	ActionDiff Action = iota
	ActionAdd
	ActionApply
	ActionMerge
	ActionInspect
	ActionSkip
	ActionQuit
)

func Recommend(entry chezmoi.StatusEntry) Action {
	switch {
	case entry.LocalChange == chezmoi.ChangeModified && entry.TargetChange == chezmoi.ChangeNone:
		return ActionAdd
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeModified:
		return ActionApply
	case entry.LocalChange == chezmoi.ChangeModified && entry.TargetChange == chezmoi.ChangeModified:
		return ActionMerge
	case entry.LocalChange == chezmoi.ChangeDeleted && entry.TargetChange == chezmoi.ChangeNone:
		return ActionInspect
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeDeleted:
		return ActionApply
	case entry.LocalChange == chezmoi.ChangeAdded && entry.TargetChange == chezmoi.ChangeNone:
		return ActionAdd
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeAdded:
		return ActionApply
	default:
		return ActionDiff
	}
}

func Describe(entry chezmoi.StatusEntry) string {
	switch {
	case entry.LocalChange == chezmoi.ChangeModified && entry.TargetChange == chezmoi.ChangeNone:
		return "local changed"
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeModified:
		return "source changed"
	case entry.LocalChange == chezmoi.ChangeModified && entry.TargetChange == chezmoi.ChangeModified:
		return "both changed"
	case entry.LocalChange == chezmoi.ChangeDeleted && entry.TargetChange == chezmoi.ChangeNone:
		return "local deleted"
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeDeleted:
		return "source wants delete"
	case entry.LocalChange == chezmoi.ChangeAdded && entry.TargetChange == chezmoi.ChangeNone:
		return "local added"
	case entry.LocalChange == chezmoi.ChangeNone && entry.TargetChange == chezmoi.ChangeAdded:
		return "source adds target"
	default:
		return "inspect"
	}
}

func ActionName(action Action) string {
	switch action {
	case ActionDiff:
		return "diff"
	case ActionAdd:
		return "add"
	case ActionApply:
		return "apply"
	case ActionMerge:
		return "merge"
	case ActionInspect:
		return "inspect"
	case ActionSkip:
		return "skip"
	case ActionQuit:
		return "quit"
	default:
		return "unknown"
	}
}
