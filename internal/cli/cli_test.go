package cli

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/testutil"
	"github.com/zhongyangchuwu/cm/internal/ui"
)

func TestRunDefaultsToReadOnlyStatus(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if service.statusCalls != 1 {
		t.Fatalf("statusCalls = %d, want 1", service.statusCalls)
	}
	if len(service.commands) != 0 {
		t.Fatalf("mutating/diff commands were called: %#v", service.commands)
	}
	if got := strings.TrimSpace(out.String()); got != "clean" {
		t.Fatalf("output = %q, want clean", got)
	}
}

func TestRunStatusRendersSimplifiedLocalAndSourceGitStatus(t *testing.T) {
	service := &fakeService{
		entries:       []chezmoi.StatusEntry{{Code: "MM", Path: "/home/me/.zshrc"}},
		sourceEntries: []chezmoi.StatusEntry{{Code: " M", Path: "dot_zshrc"}},
	}
	var out bytes.Buffer

	code := run([]string{"status"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	got := out.String()
	for _, want := range []string{"local:", "! /home/me/.zshrc", "differs from chezmoi", "run cm sync", "chezmoi:", " M dot_zshrc"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
		}
	}
	for _, notWant := range []string{"MM /home/me/.zshrc", "local drift", "apply pending", "run cm sync /home/me/.zshrc"} {
		if strings.Contains(got, notWant) {
			t.Fatalf("output %q unexpectedly contains %q", got, notWant)
		}
	}
}

func TestRunDiffUsesInternalDiff(t *testing.T) {
	service := &fakeService{diffOutput: "internal diff\n"}
	var out bytes.Buffer

	code := run([]string{"diff", ".zshrc"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if out.String() != "internal diff\n" {
		t.Fatalf("output = %q", out.String())
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"diff-output", ".zshrc"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
}

func TestRunSyncExecutesConfirmedTUIActions(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			{{Code: "MM", Path: "/home/me/.zshrc"}},
		},
	}
	var out bytes.Buffer

	code := run([]string{"sync", ".zshrc"}, testServices(service), strings.NewReader("a\ry"), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	wantActions := []ui.Action{{Target: "/home/me/.zshrc", Kind: ui.ActionAdd}}
	if !reflect.DeepEqual(service.executed, wantActions) {
		t.Fatalf("executed = %#v, want %#v", service.executed, wantActions)
	}
}

func TestRunSyncDebugWritesTempLogPathToStderr(t *testing.T) {
	service := &fakeService{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"sync", "--debug"}, testServices(service), strings.NewReader(""), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; stderr %q", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "clean" {
		t.Fatalf("stdout = %q, want clean", stdout.String())
	}
	debug := stderr.String()
	if !strings.Contains(debug, "debug log: ") || !strings.Contains(debug, "debug log kept at: ") {
		t.Fatalf("stderr = %q, want debug log path messages", debug)
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

func TestRunGitOpensSourceRepositoryWithLazygit(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"git"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"git"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
	if service.statusCalls != 0 {
		t.Fatalf("statusCalls = %d, want 0", service.statusCalls)
	}
}

func TestRunMutatingWrappersForwardToChezmoi(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want [][]string
	}{
		{name: "add", args: []string{"add", ".zshrc"}, want: [][]string{{"add", ".zshrc"}}},
		{name: "apply", args: []string{"apply", ".gitconfig"}, want: [][]string{{"apply", ".gitconfig"}}},
		{name: "merge", args: []string{"merge", ".config/nvim/init.lua"}, want: [][]string{{"merge", ".config/nvim/init.lua"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeService{}
			var out bytes.Buffer

			code := run(tt.args, testServices(service), strings.NewReader(""), &out, &out)

			if code != 0 {
				t.Fatalf("run exit code = %d, want 0", code)
			}
			if !reflect.DeepEqual(service.commands, tt.want) {
				t.Fatalf("commands = %#v, want %#v", service.commands, tt.want)
			}
			if service.statusCalls != 0 {
				t.Fatalf("statusCalls = %d, want 0", service.statusCalls)
			}
		})
	}
}

func TestRunEditForwardsTargetToChezmoi(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"edit", ".zshrc"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"edit", ".zshrc"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
}

func TestRunVersionPrintsBuildInfo(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"version"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	got := out.String()
	for _, want := range []string{"cm:", "go:"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
		}
	}
}

func TestRunCompletionPrintsShellScript(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"completion", "bash"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "cm") {
		t.Fatalf("completion output = %q, want generated script", out.String())
	}
}

func testServices(service *fakeService) commandServices {
	return commandServices{
		Status:    service,
		Diff:      service,
		Sync:      service,
		Target:    service,
		SourceGit: service,
		Edit:      service,
	}
}

type fakeService struct {
	entries       []chezmoi.StatusEntry
	statusResults [][]chezmoi.StatusEntry
	sourceEntries []chezmoi.StatusEntry
	statusCalls   int
	statusArgs    [][]string
	commands      [][]string
	diffOutput    string
	executed      []ui.Action
}

func (f *fakeService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	f.statusCalls++
	f.statusArgs = append(f.statusArgs, append([]string(nil), targets...))
	if len(f.statusResults) > 0 {
		entries := f.statusResults[0]
		f.statusResults = f.statusResults[1:]
		return append([]chezmoi.StatusEntry(nil), entries...), nil
	}
	return append([]chezmoi.StatusEntry(nil), f.entries...), nil
}

func (f *fakeService) SourceStatus() ([]chezmoi.StatusEntry, error) {
	return append([]chezmoi.StatusEntry(nil), f.sourceEntries...), nil
}

func (f *fakeService) DiffOutput(target string) ([]byte, error) {
	f.commands = append(f.commands, []string{"diff-output", target})
	return []byte(f.diffOutput), nil
}

func (f *fakeService) ExecuteNonInteractive(action ui.Action) error {
	f.executed = append(f.executed, action)
	return nil
}
func (f *fakeService) TerminalCommand(action ui.Action) (ui.TerminalCommand, error) {
	return testutil.TerminalCommand{RunFunc: func() error {
		f.executed = append(f.executed, action)
		return nil
	}}, nil
}

func (f *fakeService) AddTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"add"}, targets...))
	return nil
}

func (f *fakeService) ApplyTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"apply"}, targets...))
	return nil
}

func (f *fakeService) MergeTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"merge"}, targets...))
	return nil
}

func (f *fakeService) OpenSourceGit() error {
	f.commands = append(f.commands, []string{"git"})
	return nil
}

func (f *fakeService) EditTarget(target string) error {
	f.commands = append(f.commands, []string{"edit", target})
	return nil
}

func (f *fakeService) ManagedFiles() ([]string, error) {
	return []string{".zshrc", ".gitconfig"}, nil
}
