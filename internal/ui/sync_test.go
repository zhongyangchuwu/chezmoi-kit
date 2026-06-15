package ui

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestRunSyncAddsAndRefreshesTarget(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, []string{".zshrc"}, strings.NewReader("a\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"add", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
	wantStatusArgs := [][]string{{".zshrc"}, {"/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}

func TestRunSyncInvalidInputReprompts(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: " M", Path: "/home/me/.gitconfig"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("x\np\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"apply", "/home/me/.gitconfig"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
	if !strings.Contains(out.String(), "unknown choice") {
		t.Fatalf("output %q does not contain unknown choice", out.String())
	}
}

func TestRunSyncQuitDoesNotMutate(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.config/nvim/init.lua"}},
		},
	}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("q\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	if len(service.commands) != 0 {
		t.Fatalf("commands = %#v, want none", service.commands)
	}
}

func TestRunSyncDiffDoesNotAdvanceEntry(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("d\nm\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"diff", "/home/me/.zshrc"}, {"merge", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
}

func TestSyncModelRendersEntriesAndCurrentDiff(t *testing.T) {
	model := newSyncModel([]chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: " M", Path: "/home/me/.gitconfig"},
	}, []string{"diff -- .zshrc"})

	got := model.viewString()
	for _, want := range []string{"cm sync", "> /home/me/.zshrc", "  /home/me/.gitconfig", "diff -- .zshrc", "add", "apply", "merge", "skip", "quit"} {
		if !strings.Contains(got, want) {
			t.Fatalf("view %q does not contain %q", got, want)
		}
	}
}

func TestSyncModelNavigationAndActions(t *testing.T) {
	model := newSyncModel([]chezmoi.StatusEntry{
		{Code: "MM", Path: "/home/me/.zshrc"},
		{Code: " M", Path: "/home/me/.gitconfig"},
	}, nil)

	model = model.moveDown()
	if got := model.current().Path; got != "/home/me/.gitconfig" {
		t.Fatalf("current after down = %q, want .gitconfig", got)
	}

	action, ok := model.actionForKey("p")
	if !ok {
		t.Fatal("p did not resolve to an action")
	}
	if action.kind != syncActionApply || action.target != "/home/me/.gitconfig" {
		t.Fatalf("action = %#v, want apply .gitconfig", action)
	}

	model = model.withEntryClean("/home/me/.gitconfig")
	if got := model.current().Path; got != "/home/me/.zshrc" {
		t.Fatalf("current after clean = %q, want .zshrc", got)
	}
	if len(model.entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(model.entries))
	}
}

func TestSyncModelIgnoresEnterAfterSingleKeyAction(t *testing.T) {
	model := newSyncTUIModel(nil, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Text: "d", Code: 'd'}))
	model = next.(syncTUIModel)
	next, _ = model.Update(syncActionMsg{action: syncAction{kind: syncActionDiff, target: "/home/me/.zshrc"}, diff: "diff output"})
	model = next.(syncTUIModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(syncTUIModel)

	if got := model.viewString(); strings.Contains(got, "unknown choice") {
		t.Fatalf("view %q contains unknown choice", got)
	}
}

func TestRunSyncTUIExecutesKeyActions(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSyncTUI(service, []string{".zshrc"}, strings.NewReader("a"), &out)
	if err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}

	wantCommands := [][]string{{"add", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
	wantStatusArgs := [][]string{{".zshrc"}, {"/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}
func TestRunSyncTUIAllowsActionAfterDiff(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
		diffOutput: "diff --git a/dot_zshrc b/dot_zshrc\n",
	}
	var out bytes.Buffer

	err := RunSyncTUI(service, []string{".zshrc"}, strings.NewReader("d\ra"), &out)
	if err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}

	wantCommands := [][]string{{"diff-output", "/home/me/.zshrc"}, {"add", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
	if strings.Contains(out.String(), "unknown choice") {
		t.Fatalf("output %q contains unknown choice", out.String())
	}
}

func TestSyncModelShowsDiffOutput(t *testing.T) {
	service := &fakeService{diffOutput: "diff --git a/dot_zshrc b/dot_zshrc\n"}
	model := newSyncTUIModel(service, []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}})
	action := syncAction{kind: syncActionDiff, target: "/home/me/.zshrc"}

	msg := model.runAction(action)().(syncActionMsg)
	if msg.err != nil {
		t.Fatalf("runAction returned error: %v", msg.err)
	}
	model = model.applyActionResult(msg.action, msg.diff)

	if got := model.viewString(); !strings.Contains(got, "diff --git a/dot_zshrc b/dot_zshrc") {
		t.Fatalf("view %q does not contain diff", got)
	}
}

func TestRunSyncTUIRendersCleanWithoutProgram(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	err := RunSyncTUI(service, nil, strings.NewReader(""), &out)
	if err != nil {
		t.Fatalf("RunSyncTUI returned error: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "clean" {
		t.Fatalf("output = %q, want clean", got)
	}
}

func TestRunSyncShowsSimplifiedPrompt(t *testing.T) {
	service := &fakeService{statusResults: [][]chezmoi.StatusEntry{{{Code: "MM", Path: "/home/me/.zshrc"}}}}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("q\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}
	got := out.String()
	for _, want := range []string{"! /home/me/.zshrc", "local differs from chezmoi", "a[p]ply chezmoi"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
		}
	}
	if strings.Contains(got, "recommended") || strings.Contains(got, "MM /home") || strings.Contains(got, "[p]apply") {
		t.Fatalf("output %q still contains obsolete prompt text", got)
	}
}

type fakeService struct {
	statusResults [][]chezmoi.StatusEntry
	statusArgs    [][]string
	commands      [][]string
	diffOutput    string
}

func (f *fakeService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	f.statusArgs = append(f.statusArgs, append([]string(nil), targets...))
	if len(f.statusResults) == 0 {
		return nil, nil
	}
	entries := f.statusResults[0]
	f.statusResults = f.statusResults[1:]
	return append([]chezmoi.StatusEntry(nil), entries...), nil
}

func (f *fakeService) Diff(targets []string) error {
	f.commands = append(f.commands, append([]string{"diff"}, targets...))
	return nil
}

func (f *fakeService) DiffOutput(targets []string) ([]byte, error) {
	f.commands = append(f.commands, append([]string{"diff-output"}, targets...))
	return []byte(f.diffOutput), nil
}

func (f *fakeService) Add(target string) error {
	f.commands = append(f.commands, []string{"add", target})
	return nil
}

func (f *fakeService) Apply(target string) error {
	f.commands = append(f.commands, []string{"apply", target})
	return nil
}

func (f *fakeService) Merge(target string) error {
	f.commands = append(f.commands, []string{"merge", target})
	return nil
}
