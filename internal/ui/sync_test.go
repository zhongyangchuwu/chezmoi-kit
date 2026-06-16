package ui

import (
	"reflect"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
)

func TestSyncModelTogglesPendingAction(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	model = model.togglePending(reconcile.ActionAdd)
	if got, ok := model.pending["/home/me/.zshrc"]; !ok || got != reconcile.ActionAdd {
		t.Fatalf("pending = %#v", model.pending)
	}

	model = model.togglePending(reconcile.ActionAdd)
	if len(model.pending) != 0 {
		t.Fatalf("pending = %#v, want empty", model.pending)
	}
}

func TestSyncModelReplacesPendingAction(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	model = model.togglePending(reconcile.ActionAdd)
	model = model.togglePending(reconcile.ActionApply)

	if got := model.pending["/home/me/.zshrc"]; got != reconcile.ActionApply {
		t.Fatalf("pending action = %v, want apply", got)
	}
}

func TestSyncModelPendingActionsKeepEntryOrder(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: "MM", Path: "/home/me/.gitconfig"},
	})
	model.cursor = 1
	model = model.togglePending(reconcile.ActionMerge)
	model.cursor = 0
	model = model.togglePending(reconcile.ActionAdd)

	want := []reconcile.Action{
		{Target: "/home/me/.zshrc", Kind: reconcile.ActionAdd},
		{Target: "/home/me/.gitconfig", Kind: reconcile.ActionMerge},
	}
	if got := model.pendingActions(); !reflect.DeepEqual(got, want) {
		t.Fatalf("pendingActions = %#v, want %#v", got, want)
	}
}

func TestExecutePendingPreflightDropsCleanTargets(t *testing.T) {
	service := &fakeReviewService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.gitconfig"}},
		},
	}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: "MM", Path: "/home/me/.gitconfig"},
	})
	model.pending["/home/me/.zshrc"] = reconcile.ActionAdd
	model.pending["/home/me/.gitconfig"] = reconcile.ActionMerge

	msg := model.executePending()().(executeMsg)
	if msg.err != nil {
		t.Fatalf("executePending returned error: %v", msg.err)
	}
	wantExecuted := []reconcile.Action{{Target: "/home/me/.gitconfig", Kind: reconcile.ActionMerge}}
	if !reflect.DeepEqual(service.executed, wantExecuted) {
		t.Fatalf("executed = %#v, want %#v", service.executed, wantExecuted)
	}
	wantStatusArgs := [][]string{{"/home/me/.zshrc", "/home/me/.gitconfig"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}

type fakeReviewService struct {
	statusResults [][]chezmoi.StatusEntry
	statusArgs    [][]string
	diffOutput    string
	executed      []reconcile.Action
}

func (f *fakeReviewService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	f.statusArgs = append(f.statusArgs, append([]string(nil), targets...))
	if len(f.statusResults) == 0 {
		return nil, nil
	}
	entries := f.statusResults[0]
	f.statusResults = f.statusResults[1:]
	return append([]chezmoi.StatusEntry(nil), entries...), nil
}

func (f *fakeReviewService) DiffOutput(string) ([]byte, error) {
	return []byte(f.diffOutput), nil
}

func (f *fakeReviewService) Execute(actions []reconcile.Action) error {
	f.executed = append(f.executed, actions...)
	return nil
}
