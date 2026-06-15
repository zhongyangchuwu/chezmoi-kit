package cli

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestRunDefaultsToReadOnlyStatus(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{}, service, strings.NewReader(""), &out, &out)

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
		sourceEntries: []sourceEntry{{Code: " M", Path: "dot_zshrc"}},
	}
	var out bytes.Buffer

	code := run([]string{"status"}, service, strings.NewReader(""), &out, &out)

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

func TestRunDiffForwardsToChezmoi(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"diff", ".zshrc"}, service, strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"diff", ".zshrc"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
	if service.statusCalls != 0 {
		t.Fatalf("statusCalls = %d, want 0", service.statusCalls)
	}
}

func TestRunSyncTUIFlagUsesTUI(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	code := run([]string{"sync", "--tui", ".zshrc"}, service, strings.NewReader("a"), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"add", "/home/me/.zshrc"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
	wantStatusArgs := [][]string{{".zshrc"}, {"/home/me/.zshrc"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}

func TestRunSyncPlainFlagKeepsPromptMode(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{Code: "MM", Path: "/home/me/.zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	code := run([]string{"sync", "--plain", ".zshrc"}, service, strings.NewReader("a\n"), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !strings.Contains(out.String(), "[d]iff [a]dd local") {
		t.Fatalf("output %q does not contain plain prompt", out.String())
	}
	if strings.Contains(out.String(), "cm sync") {
		t.Fatalf("output %q unexpectedly contains TUI title", out.String())
	}
}

func TestRunGitOpensSourceRepositoryWithLazygit(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"git"}, service, strings.NewReader(""), &out, &out)

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

			code := run(tt.args, service, strings.NewReader(""), &out, &out)

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

func TestRunVersionPrintsBuildInfo(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"version"}, service, strings.NewReader(""), &out, &out)

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

	code := run([]string{"completion", "bash"}, service, strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "cm") {
		t.Fatalf("completion output = %q, want generated script", out.String())
	}
}

type fakeService struct {
	entries       []chezmoi.StatusEntry
	statusResults [][]chezmoi.StatusEntry
	sourceEntries []sourceEntry
	statusCalls   int
	statusArgs    [][]string
	commands      [][]string
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

func (f *fakeService) SourceStatus() ([]sourceEntry, error) {
	return append([]sourceEntry(nil), f.sourceEntries...), nil
}

func (f *fakeService) Diff(targets []string) error {
	f.commands = append(f.commands, append([]string{"diff"}, targets...))
	return nil
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
