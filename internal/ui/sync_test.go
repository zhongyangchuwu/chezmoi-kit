package ui

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/zhongyangchuwu/cm/internal/app"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/testutil"
)

func TestSyncModelTogglesPendingAction(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	model = model.togglePending(ActionAdd)
	if got, ok := model.pending["/home/me/.zshrc"]; !ok || got != ActionAdd {
		t.Fatalf("pending = %#v", model.pending)
	}

	model = model.togglePending(ActionAdd)
	if len(model.pending) != 0 {
		t.Fatalf("pending = %#v, want empty", model.pending)
	}
}

func TestSyncModelReplacesPendingAction(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	model = model.togglePending(ActionAdd)
	model = model.togglePending(ActionApply)

	if got := model.pending["/home/me/.zshrc"]; got != ActionApply {
		t.Fatalf("pending action = %v, want apply", got)
	}
}

func TestSyncModelPendingActionsKeepEntryOrder(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: "MM", Path: "/home/me/.gitconfig"},
	})
	model.cursor = 1
	model = model.togglePending(ActionMerge)
	model.cursor = 0
	model = model.togglePending(ActionAdd)

	want := []Action{
		{Target: "/home/me/.zshrc", Kind: ActionAdd},
		{Target: "/home/me/.gitconfig", Kind: ActionMerge},
	}
	if got := model.pendingActions(); !reflect.DeepEqual(got, want) {
		t.Fatalf("pendingActions = %#v, want %#v", got, want)
	}
}

func TestExecuteCurrentActionDropsCleanTarget(t *testing.T) {
	service := &fakeReviewService{
		statusResults: [][]chezmoi.StatusEntry{{}},
	}
	cmd := executeNonInteractiveCmd(service, Action{Target: "/home/me/.zshrc", Kind: ActionAdd}, nil)

	msg := cmd().(executeMsg)
	if msg.err != nil {
		t.Fatalf("executeActionCmd returned error: %v", msg.err)
	}
	if len(msg.executed) != 0 || msg.skipped != 1 {
		t.Fatalf("executeMsg = %#v, want one skipped clean target", msg)
	}
	if len(service.executed) != 0 {
		t.Fatalf("executed = %#v, want no execution", service.executed)
	}
	wantStatusArgs := [][]string{{"/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}

func TestExecuteCurrentActionWritesTimingLog(t *testing.T) {
	service := &fakeReviewService{statusResults: [][]chezmoi.StatusEntry{{{Code: "MM", Path: "/home/me/.zshrc"}}}}
	timing, err := newSyncTimingLogger(true)
	if err != nil {
		t.Fatalf("newSyncTimingLogger returned error: %v", err)
	}
	logPath := timing.Path()
	t.Cleanup(func() { _ = os.Remove(logPath) })
	cmd := executeNonInteractiveCmd(service, Action{Target: "/home/me/.zshrc", Kind: ActionAdd}, timing)

	msg := cmd().(executeMsg)
	if msg.err != nil {
		t.Fatalf("executeActionCmd returned error: %v", msg.err)
	}
	if err := timing.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	log := string(content)
	for _, want := range []string{"sync status preflight", "sync execute target", "target=/home/me/.zshrc", "action=add", "duration="} {
		if !strings.Contains(log, want) {
			t.Fatalf("log %q does not contain %q", log, want)
		}
	}
}

func TestRunSyncTUIDebugPrintsTempLogPath(t *testing.T) {
	service := &fakeReviewService{}
	var out bytes.Buffer
	var stderr bytes.Buffer

	err := RunSyncTUI(service, nil, strings.NewReader(""), &out, &app.Options{Debug: true, Stderr: &stderr})
	if err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}
	if strings.TrimSpace(out.String()) != "clean" {
		t.Fatalf("stdout = %q, want clean", out.String())
	}
	debug := stderr.String()
	if !strings.Contains(debug, "debug log: ") || !strings.Contains(debug, "debug log kept at: ") {
		t.Fatalf("stderr = %q, want start and end debug log paths", debug)
	}
	path := testutil.DebugLogPathFromLine(t, debug, "debug log: ")
	defer os.Remove(path)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}
	if !strings.Contains(string(content), "sync initial status") {
		t.Fatalf("debug log %q does not contain initial status: %q", path, string(content))
	}
}

func TestExecutionReturnsToReviewWithRemainingFiles(t *testing.T) {
	service := &fakeReviewService{statusResults: [][]chezmoi.StatusEntry{{{Code: "MM", Path: "/home/me/.zshrc"}}}}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: "MM", Path: "/home/me/.gitconfig"},
	})
	model.pending["/home/me/.zshrc"] = ActionAdd

	model, cmd := model.startExecution(model.pendingActions())
	if cmd == nil {
		t.Fatal("cmd is nil, want first action command")
	}
	msg := cmd().(executeMsg)
	updated, next := model.applyExecuteMsg(msg)
	model = updated

	if model.mode != modeReview {
		t.Fatalf("mode = %v, want review", model.mode)
	}
	if len(model.entries) != 1 || model.entries[0].Path != "/home/me/.gitconfig" {
		t.Fatalf("entries = %#v, want remaining gitconfig", model.entries)
	}
	if _, ok := model.pending["/home/me/.zshrc"]; ok {
		t.Fatalf("pending = %#v, want executed target removed", model.pending)
	}
	if next == nil {
		t.Fatal("next cmd is nil, want diff load for remaining file")
	}
}

func TestExecutionQuitsWithCompleteMessageWhenAllFilesDone(t *testing.T) {
	service := &fakeReviewService{statusResults: [][]chezmoi.StatusEntry{{{Code: "MM", Path: "/home/me/.zshrc"}}}}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})
	model.pending["/home/me/.zshrc"] = ActionApply

	model, cmd := model.startExecution(model.pendingActions())
	msg := cmd().(executeMsg)
	updated, next := model.applyExecuteMsg(msg)
	model = updated

	if !model.completed || !model.stopped {
		t.Fatalf("completed=%v stopped=%v, want true", model.completed, model.stopped)
	}
	if len(model.entries) != 0 {
		t.Fatalf("entries = %#v, want none", model.entries)
	}
	if model.message != "sync complete: executed 1, skipped 0" {
		t.Fatalf("message = %q", model.message)
	}
	if got := model.viewString(); got != "sync complete: executed 1, skipped 0\n" {
		t.Fatalf("view = %q", got)
	}
	if next == nil {
		t.Fatal("next cmd is nil, want quit command")
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

func TestInitLoadsCurrentDiffByDefault(t *testing.T) {
	service := &fakeReviewService{diffOutput: "diff"}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	cmd := model.Init()
	if cmd == nil {
		t.Fatal("cmd is nil, want initial diff load")
	}
	if state := model.currentDiffState(); !state.loading {
		t.Fatalf("initial diff state = %#v, want loading", state)
	}
	msg := cmd().(syncDiffMsg)
	if msg.target != "/home/me/.zshrc" {
		t.Fatalf("diff target = %q", msg.target)
	}
}

func TestMoveSelectionLoadsDiffByDefault(t *testing.T) {
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
	if cmd == nil {
		t.Fatal("cmd is nil, want default diff load")
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
	diffArgs      []string
	executed      []Action
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

func (f *fakeReviewService) DiffOutput(target string) ([]byte, error) {
	f.diffArgs = append(f.diffArgs, target)
	return []byte(f.diffOutput), nil
}

func (f *fakeReviewService) ExecuteNonInteractive(action Action) error {
	f.executed = append(f.executed, action)
	return nil
}
func (f *fakeReviewService) TerminalCommand(action Action) (TerminalCommand, error) {
	return testutil.TerminalCommand{RunFunc: func() error {
		f.executed = append(f.executed, action)
		return nil
	}}, nil
}
