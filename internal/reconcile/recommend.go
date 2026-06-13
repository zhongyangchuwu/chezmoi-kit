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
	case entry.LocalChange != chezmoi.ChangeNone && entry.TargetChange != chezmoi.ChangeNone:
		return ActionInspect
	case entry.LocalChange == chezmoi.ChangeModified:
		return ActionAdd
	case entry.LocalChange == chezmoi.ChangeAdded:
		return ActionAdd
	case entry.LocalChange == chezmoi.ChangeDeleted:
		return ActionInspect
	case entry.TargetChange == chezmoi.ChangeModified:
		return ActionApply
	case entry.TargetChange == chezmoi.ChangeAdded:
		return ActionApply
	case entry.TargetChange == chezmoi.ChangeDeleted:
		return ActionApply
	default:
		return ActionDiff
	}
}

func Describe(entry chezmoi.StatusEntry) string {
	local := localDescription(entry.LocalChange)
	apply := applyDescription(entry.TargetChange)
	switch {
	case local != "" && apply != "":
		return local + ", " + apply
	case local != "":
		return local
	case apply != "":
		return apply
	default:
		return "inspect"
	}
}

func localDescription(change chezmoi.Change) string {
	switch change {
	case chezmoi.ChangeModified:
		return "local drift"
	case chezmoi.ChangeAdded:
		return "local added"
	case chezmoi.ChangeDeleted:
		return "local deleted"
	default:
		return ""
	}
}

func applyDescription(change chezmoi.Change) string {
	switch change {
	case chezmoi.ChangeModified:
		return "apply pending"
	case chezmoi.ChangeAdded:
		return "apply would add"
	case chezmoi.ChangeDeleted:
		return "apply would delete"
	case chezmoi.ChangeRun:
		return "apply would run"
	default:
		return ""
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
