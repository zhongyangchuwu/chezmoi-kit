package tui

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func TestSyncModelTogglesReviewedPendingAction(t *testing.T) {
	target := "/home/me/.zshrc"
	model := newSyncTUIModel(nil, syncStatus(target))
	model.storeReview(dirtyReview(target, app.TargetFile, false, "review-1", "diff"))

	model = model.togglePending(app.ActionAdd)
	got, ok := model.pending[target]
	if !ok || got.Kind != app.ActionAdd || got.Fingerprint != "review-1" {
		t.Fatalf("pending = %#v", model.pending)
	}

	model = model.togglePending(app.ActionAdd)
	if len(model.pending) != 0 {
		t.Fatalf("pending = %#v, want empty", model.pending)
	}
}

func TestSyncModelRejectsActionInvalidForTargetType(t *testing.T) {
	target := "/home/me/.link"
	model := newSyncTUIModel(nil, syncStatus(target))
	model.storeReview(dirtyReview(target, app.TargetSymlink, false, "review-1", "diff"))

	model = model.togglePending(app.ActionAdd)
	if len(model.pending) != 0 {
		t.Fatalf("pending = %#v, want empty", model.pending)
	}
	if !strings.Contains(model.message, "not available") {
		t.Fatalf("message = %q", model.message)
	}
}

func TestHelpShowsOnlyActionsAllowedByCurrentReview(t *testing.T) {
	target := "/home/me/.link"
	model := newSyncTUIModel(nil, syncStatus(target))
	model.storeReview(dirtyReview(target, app.TargetSymlink, false, "review-1", "diff"))

	descriptions := make(map[string]bool)
	for _, binding := range defaultSyncKeys.reviewHelp(focusFiles, model.currentDiffState().review) {
		descriptions[binding.Help().Desc] = true
	}
	if !descriptions["apply"] {
		t.Fatalf("help descriptions = %#v, want apply", descriptions)
	}
	if descriptions["add"] || descriptions["merge"] {
		t.Fatalf("help descriptions = %#v, contains invalid actions", descriptions)
	}
}

func TestSyncModelPendingActionsKeepEntryOrder(t *testing.T) {
	first := "/home/me/.zshrc"
	second := "/home/me/.gitconfig"
	status := app.SyncStatus{Entries: []app.ReconcileEntry{{Code: "MM", Path: first}, {Code: "MM", Path: second}}}
	model := newSyncTUIModel(nil, status)
	model.storeReview(dirtyReview(first, app.TargetFile, false, "first", "diff"))
	model.storeReview(dirtyReview(second, app.TargetFile, false, "second", "diff"))
	model.cursor = 1
	model = model.togglePending(app.ActionMerge)
	model.cursor = 0
	model = model.togglePending(app.ActionAdd)

	want := []app.Action{
		{Target: first, Kind: app.ActionAdd, Fingerprint: "first"},
		{Target: second, Kind: app.ActionMerge, Fingerprint: "second"},
	}
	if got := model.pendingActions(); !reflect.DeepEqual(got, want) {
		t.Fatalf("pendingActions = %#v, want %#v", got, want)
	}
}

func TestExecuteCurrentActionSkipsAlreadyCleanTarget(t *testing.T) {
	target := "/home/me/.zshrc"
	service := &fakeReviewService{reviews: []app.Review{{Entry: app.ReconcileEntry{Path: target}, Dirty: false}}}
	action := app.Action{Target: target, Kind: app.ActionAdd, Fingerprint: "reviewed"}

	msg := executeNonInteractiveCmd(service, action, nil)().(executeMsg)
	if msg.err != nil || !msg.skipped || !msg.resolved || msg.executed {
		t.Fatalf("executeMsg = %#v", msg)
	}
	if len(service.executed) != 0 {
		t.Fatalf("executed = %#v, want none", service.executed)
	}
}

func TestExecuteCurrentActionDefersChangedReview(t *testing.T) {
	target := "/home/me/.zshrc"
	changed := dirtyReview(target, app.TargetFile, false, "changed", "new diff")
	service := &fakeReviewService{reviews: []app.Review{changed}}
	action := app.Action{Target: target, Kind: app.ActionAdd, Fingerprint: "reviewed"}

	msg := executeNonInteractiveCmd(service, action, nil)().(executeMsg)
	if msg.err != nil || !msg.deferred || msg.reason != "changed since review" {
		t.Fatalf("executeMsg = %#v", msg)
	}
	if len(service.executed) != 0 {
		t.Fatalf("executed = %#v, want none", service.executed)
	}
}

func TestSuccessfulExecutionRetainsStillDirtyTarget(t *testing.T) {
	target := "/home/me/.zshrc"
	before := dirtyReview(target, app.TargetFile, false, "reviewed", "old diff")
	after := dirtyReview(target, app.TargetFile, false, "after", "remaining diff")
	service := &fakeReviewService{
		reviews: []app.Review{before, after},
		result:  app.ActionResult{Stderr: "warning from chezmoi\n"},
	}
	action := app.Action{Target: target, Kind: app.ActionAdd, Fingerprint: before.Fingerprint}

	msg := executeNonInteractiveCmd(service, action, nil)().(executeMsg)
	if msg.err != nil || !msg.executed || !msg.deferred || msg.resolved {
		t.Fatalf("executeMsg = %#v", msg)
	}

	model := newSyncTUIModel(service, syncStatus(target))
	model.pending[target] = action
	model.executing = []app.Action{action}
	model.mode = modeExecuting
	updated, _ := model.applyExecuteMsg(msg)
	if len(updated.entries) != 1 || updated.entries[0].Path != target {
		t.Fatalf("entries = %#v, want unresolved target", updated.entries)
	}
	if !strings.Contains(updated.message, "target still differs after execution") || !strings.Contains(updated.message, "warning from chezmoi") {
		t.Fatalf("message = %q", updated.message)
	}
}

func TestExecutionNoticesPrefixEveryOutputLine(t *testing.T) {
	target := "/home/me/.zshrc"
	model := newSyncTUIModel(&fakeReviewService{}, syncStatus(target))
	model.executing = []app.Action{{Target: target, Kind: app.ActionApply}}
	model.mode = modeExecuting

	updated, _ := model.applyExecuteMsg(executeMsg{
		target:   target,
		executed: true,
		resolved: true,
		result: app.ActionResult{
			Stdout: "updated\nsecond update\n",
			Stderr: "warning\nsecond warning\n",
		},
	})

	prefix := updated.displayPath(target) + ": "
	for _, line := range []string{"updated", "second update", "warning", "second warning"} {
		if !strings.Contains(updated.message, prefix+line) {
			t.Fatalf("message does not attribute output line %q: %q", line, updated.message)
		}
	}
}

func TestExecutionQuitsOnlyAfterPostflightIsClean(t *testing.T) {
	target := "/home/me/.zshrc"
	before := dirtyReview(target, app.TargetFile, false, "reviewed", "diff")
	service := &fakeReviewService{
		reviews: []app.Review{before, {Entry: app.ReconcileEntry{Path: target}, Dirty: false}},
		result:  app.ActionResult{Stdout: "updated\n"},
	}
	model := newSyncTUIModel(service, syncStatus(target))
	model.storeReview(before)
	model = model.togglePending(app.ActionApply)

	model, cmd := model.startExecution(model.pendingActions())
	msg := cmd().(executeMsg)
	model, next := model.applyExecuteMsg(msg)

	if !model.completed || !model.stopped || len(model.entries) != 0 {
		t.Fatalf("model = %#v, want completed empty state", model)
	}
	if !strings.Contains(model.message, "sync complete: executed 1, skipped 0, deferred 0") || !strings.Contains(model.message, "updated") {
		t.Fatalf("message = %q", model.message)
	}
	if next == nil {
		t.Fatal("next cmd is nil, want quit command")
	}
}

func TestRunSyncTUIReportsScriptsWithoutOpeningFileReview(t *testing.T) {
	service := &fakeReviewService{statuses: []app.SyncStatus{{Scripts: []app.ReconcileEntry{{Code: " R", Path: "/home/me/install.sh"}}}}}
	var out bytes.Buffer

	if err := RunSyncTUI(service, nil, strings.NewReader(""), &out, &app.Options{}); err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}
	if got := out.String(); !strings.Contains(got, "no file changes to reconcile; 1 script pending") || !strings.Contains(got, "chezmoi apply") {
		t.Fatalf("output = %q", got)
	}
	if len(service.reviewArgs) != 0 || len(service.executed) != 0 {
		t.Fatalf("service calls = reviews %#v executed %#v", service.reviewArgs, service.executed)
	}
}

func TestRunSyncTUIDebugWritesReviewTimingWithoutContents(t *testing.T) {
	service := &fakeReviewService{}
	var out bytes.Buffer
	var stderr bytes.Buffer

	err := RunSyncTUI(service, nil, strings.NewReader(""), &out, &app.Options{Debug: true, Stderr: &stderr})
	if err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}
	path := debugLogPathFromLine(t, stderr.String(), "debug log: ")
	defer os.Remove(path)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}
	if !strings.Contains(string(content), "sync initial status") {
		t.Fatalf("debug log = %q", content)
	}
}

func TestExecutingModeIgnoresQuitKey(t *testing.T) {
	model := newSyncTUIModel(nil, syncStatus("/home/me/.zshrc"))
	model.mode = modeExecuting
	model.message = "executing 1 action"

	updated, cmd := model.Update(keyPress('q'))
	got := updated.(workspaceModel)
	if cmd != nil || got.mode != modeExecuting || got.message != "executing 1 action" {
		t.Fatalf("model = %#v cmd=%v", got, cmd)
	}
}

func TestInitLoadsCurrentReviewByDefault(t *testing.T) {
	target := "/home/me/.zshrc"
	service := &fakeReviewService{reviews: []app.Review{dirtyReview(target, app.TargetFile, false, "review", "diff")}}
	model := newSyncTUIModel(service, syncStatus(target))

	cmd := model.Init()
	if cmd == nil || !model.currentDiffState().loading {
		t.Fatalf("cmd=%v state=%#v", cmd, model.currentDiffState())
	}
	msg := cmd().(syncReviewMsg)
	if msg.target != target || msg.review.Fingerprint != "review" {
		t.Fatalf("syncDiffMsg = %#v", msg)
	}
}

func TestReviewThatBecameCleanQuitsWhenNoEntriesRemain(t *testing.T) {
	target := "/home/me/.zshrc"
	model := newSyncTUIModel(nil, syncStatus(target))

	updated, cmd := model.applyReview(syncReviewMsg{target: target, review: app.Review{Entry: app.ReconcileEntry{Path: target}}})
	if !updated.completed || !updated.stopped || len(updated.entries) != 0 {
		t.Fatalf("model = %#v, want completed clean state", updated)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want quit command")
	}
}
func TestReviewThatBecameCleanLoadsNextEntry(t *testing.T) {
	first := "/home/me/.zshrc"
	second := "/home/me/.gitconfig"
	status := app.SyncStatus{Entries: []app.ReconcileEntry{{Code: "MM", Path: first}, {Code: "MM", Path: second}}}
	service := &fakeReviewService{reviews: []app.Review{dirtyReview(second, app.TargetFile, false, "second", "diff")}}
	model := newSyncTUIModel(service, status)

	updated, cmd := model.applyReview(syncReviewMsg{target: first, review: app.Review{Entry: app.ReconcileEntry{Path: first}}})
	if len(updated.entries) != 1 || updated.entries[0].Path != second {
		t.Fatalf("entries = %#v, want second entry", updated.entries)
	}
	if cmd == nil {
		t.Fatal("cmd is nil, want next review load")
	}
	msg := cmd().(syncReviewMsg)
	if msg.target != second {
		t.Fatalf("next review target = %q, want %q", msg.target, second)
	}
}

func TestMoveSelectionLoadsReviewByDefault(t *testing.T) {
	first := "/home/me/.zshrc"
	second := "/home/me/.gitconfig"
	status := app.SyncStatus{Entries: []app.ReconcileEntry{{Code: "MM", Path: first}, {Code: "MM", Path: second}}}
	service := &fakeReviewService{reviews: []app.Review{dirtyReview(second, app.TargetFile, false, "second", "diff")}}
	model := newSyncTUIModel(service, status)

	updated, cmd := model.handleDown()
	model = updated.(workspaceModel)
	if model.cursor != 1 || cmd == nil {
		t.Fatalf("cursor=%d cmd=%v", model.cursor, cmd)
	}
	msg := cmd().(syncReviewMsg)
	if msg.target != second {
		t.Fatalf("target = %q", msg.target)
	}
}

func TestDiffScrollOnlyMovesInDiffFocus(t *testing.T) {
	target := "/home/me/.zshrc"
	model := newSyncTUIModel(nil, syncStatus(target))
	model.diffs[target] = diffState{lines: []string{"1", "2", "3", "4", "5", "6", "7", "8"}}
	model.height = 8
	model.focus = focusDiff

	updated, _ := model.handleDown()
	if got := updated.(workspaceModel).diffScroll; got != 1 {
		t.Fatalf("diffScroll = %d, want 1", got)
	}
}

func TestFocusToggle(t *testing.T) {
	model := newSyncTUIModel(nil, syncStatus("/home/me/.zshrc"))
	model = model.toggleFocus()
	if model.focus != focusDiff {
		t.Fatalf("focus = %v, want diff", model.focus)
	}
	model = model.toggleFocus()
	if model.focus != focusFiles {
		t.Fatalf("focus = %v, want files", model.focus)
	}
}

func TestTruncatePreservesUTF8AndDisplayWidth(t *testing.T) {
	for _, tt := range []struct {
		input string
		width int
	}{{"配置文件路径", 7}, {"配置", 1}, {"abcdef", 4}} {
		got := truncate(tt.input, tt.width)
		if !utf8.ValidString(got) || lipgloss.Width(got) > tt.width {
			t.Fatalf("truncate(%q, %d) = %q width=%d", tt.input, tt.width, got, lipgloss.Width(got))
		}
	}
}

func TestRenderedPanesRespectDisplayWidth(t *testing.T) {
	target := "/home/me/配置文件路径"
	model := newSyncTUIModel(&fakeReviewService{}, syncStatus(target))
	model.diffs[target] = diffState{lines: []string{"+配置文件路径很长"}}

	for name, rendered := range map[string]string{
		"files": model.renderFilesPane(rect{width: 12, height: 4}),
		"diff":  model.renderDiffPane(rect{width: 12, height: 4}),
	} {
		for _, line := range strings.Split(rendered, "\n") {
			if lipgloss.Width(line) > 12 {
				t.Fatalf("%s line width = %d, line %q", name, lipgloss.Width(line), line)
			}
		}
	}
}

func syncStatus(target string) app.SyncStatus {
	return app.SyncStatus{Entries: []app.ReconcileEntry{{Code: "MM", Path: target}}}
}

func dirtyReview(target string, targetType app.TargetType, template bool, fingerprint, diff string) app.Review {
	return app.Review{
		Entry:       app.ReconcileEntry{Code: "MM", Path: target},
		Type:        targetType,
		Template:    template,
		Diff:        diff,
		Fingerprint: fingerprint,
		Dirty:       true,
	}
}

func keyPress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

type fakeReviewService struct {
	statuses   []app.SyncStatus
	reviews    []app.Review
	reviewArgs []string
	result     app.ActionResult
	executeErr error
	executed   []app.Action
}

func (f *fakeReviewService) Status(targets []string) (app.SyncStatus, error) {
	if len(f.statuses) == 0 {
		return app.SyncStatus{}, nil
	}
	status := f.statuses[0]
	f.statuses = f.statuses[1:]
	return status, nil
}

func (f *fakeReviewService) Review(target string) (app.Review, error) {
	f.reviewArgs = append(f.reviewArgs, target)
	if len(f.reviews) == 0 {
		return app.Review{Entry: app.ReconcileEntry{Path: target}, Dirty: false}, nil
	}
	review := f.reviews[0]
	f.reviews = f.reviews[1:]
	return review, nil
}

func (f *fakeReviewService) ExecuteNonInteractive(action app.Action) (app.ActionResult, error) {
	f.executed = append(f.executed, action)
	return f.result, f.executeErr
}

func (f *fakeReviewService) TerminalCommand(action app.Action) (app.TerminalCommand, error) {
	return terminalCommand{run: func() error {
		f.executed = append(f.executed, action)
		return nil
	}}, nil
}

type terminalCommand struct {
	run func() error
}

func (c terminalCommand) Run() error {
	if c.run == nil {
		return nil
	}
	return c.run()
}

func (terminalCommand) SetStdin(io.Reader)  {}
func (terminalCommand) SetStdout(io.Writer) {}
func (terminalCommand) SetStderr(io.Writer) {}

func debugLogPathFromLine(t *testing.T, output, prefix string) string {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	t.Fatalf("output %q missing prefix %q", output, prefix)
	return ""
}
