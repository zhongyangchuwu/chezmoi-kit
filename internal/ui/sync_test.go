package ui

import (
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"

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

	msg := model.executeActions(model.pendingActions())().(executeMsg)
	if msg.err != nil {
		t.Fatalf("executeActions returned error: %v", msg.err)
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

func TestExecutingModeIgnoresQuitKey(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})
	model.mode = modeExecuting
	model.message = "executing 1 action"

	updated, cmd := model.Update(keyPress('q'))
	got := updated.(syncTUIModel)

	if cmd != nil {
		t.Fatal("cmd is non-nil, want executing mode to ignore quit")
	}
	if got.mode != modeExecuting || got.message != "executing 1 action" {
		t.Fatalf("model = %#v, want unchanged executing state", got)
	}
}

func TestExecutingHelpDoesNotAdvertiseQuit(t *testing.T) {
	if got := defaultSyncKeys.executingHelp(); len(got) != 0 {
		t.Fatalf("executing help = %#v, want empty", got)
	}
}

func TestMoveSelectionDoesNotLoadDiffUntilDiffFocus(t *testing.T) {
	service := &fakeReviewService{diffOutput: "diff"}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: "MM", Path: "/home/me/.gitconfig"},
	})

	updated, cmd := model.handleDown()
	model = updated.(syncTUIModel)
	if model.cursor != 1 {
		t.Fatalf("cursor = %d, want 1", model.cursor)
	}
	if cmd != nil {
		t.Fatal("cmd is non-nil, want lazy diff load")
	}

	updated, cmd = model.updateReview(keyPress(tea.KeyTab))
	model = updated.(syncTUIModel)
	if model.focus != focusDiff {
		t.Fatalf("focus = %v, want diff", model.focus)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want diff load")
	}
	msg := cmd().(syncDiffMsg)
	if msg.target != "/home/me/.gitconfig" {
		t.Fatalf("diff target = %q", msg.target)
	}
}

func TestDiffScrollOnlyMovesInDiffFocus(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})
	model.height = 8
	model.diffs["/home/me/.zshrc"] = diffState{lines: []string{"1", "2", "3", "4", "5", "6", "7", "8"}}
	model.focus = focusDiff

	updated, _ := model.handleDown()
	model = updated.(syncTUIModel)
	if model.diffScroll != 1 {
		t.Fatalf("diffScroll = %d, want 1", model.diffScroll)
	}
}

func TestFocusToggle(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})
	model = model.toggleFocus()
	if model.focus != focusDiff {
		t.Fatalf("focus = %v, want diff", model.focus)
	}
	model = model.toggleFocus()
	if model.focus != focusFiles {
		t.Fatalf("focus = %v, want files", model.focus)
	}
}

func keyPress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
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
