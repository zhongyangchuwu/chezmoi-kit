package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
)

func (m syncTUIModel) executeActions(actions []reconcile.Action) tea.Cmd {
	return func() tea.Msg {
		if len(actions) == 0 {
			return executeMsg{}
		}
		targets := make([]string, 0, len(actions))
		for _, action := range actions {
			targets = append(targets, action.Target)
		}
		fresh, err := m.service.Status(targets)
		if err != nil {
			return executeMsg{err: err}
		}
		dirty := make(map[string]struct{}, len(fresh))
		for _, entry := range fresh {
			dirty[entry.Path] = struct{}{}
		}
		kept := actions[:0]
		for _, action := range actions {
			if _, ok := dirty[action.Target]; ok {
				kept = append(kept, action)
			}
		}
		if len(kept) == 0 {
			return executeMsg{skipped: len(actions)}
		}
		if err := m.service.Execute(kept); err != nil {
			return executeMsg{err: err}
		}
		return executeMsg{executed: kept, skipped: len(actions) - len(kept)}
	}
}

func actionCount(count int) string {
	if count == 1 {
		return "1 action"
	}
	return fmt.Sprintf("%d actions", count)
}

func actionLabel(kind reconcile.ActionKind) string {
	return kind.String()
}

func actionMarker(kind reconcile.ActionKind) string {
	return kind.Marker()
}
