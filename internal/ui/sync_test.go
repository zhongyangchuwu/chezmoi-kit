package ui

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

func TestRunSyncAddsAndRefreshesTarget(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{LocalChange: chezmoi.ChangeModified, TargetChange: chezmoi.ChangeNone, Path: ".zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, []string{".zshrc"}, strings.NewReader("a\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"add", ".zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
	wantStatusArgs := [][]string{{".zshrc"}, {".zshrc"}}
	if !reflect.DeepEqual(service.statusArgs, wantStatusArgs) {
		t.Fatalf("statusArgs = %#v, want %#v", service.statusArgs, wantStatusArgs)
	}
}

func TestRunSyncInvalidInputReprompts(t *testing.T) {
	service := &fakeService{
		statusResults: [][]chezmoi.StatusEntry{
			{{LocalChange: chezmoi.ChangeNone, TargetChange: chezmoi.ChangeModified, Path: ".gitconfig"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("x\np\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"apply", ".gitconfig"}}
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
			{{LocalChange: chezmoi.ChangeModified, TargetChange: chezmoi.ChangeModified, Path: ".config/nvim/init.lua"}},
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
			{{LocalChange: chezmoi.ChangeModified, TargetChange: chezmoi.ChangeModified, Path: ".zshrc"}},
			nil,
		},
	}
	var out bytes.Buffer

	err := RunSync(service, nil, strings.NewReader("d\nm\n"), &out)
	if err != nil {
		t.Fatalf("RunSync returned error: %v", err)
	}

	wantCommands := [][]string{{"diff", ".zshrc"}, {"merge", ".zshrc"}}
	if !reflect.DeepEqual(service.commands, wantCommands) {
		t.Fatalf("commands = %#v, want %#v", service.commands, wantCommands)
	}
}

type fakeService struct {
	statusResults [][]chezmoi.StatusEntry
	statusArgs    [][]string
	commands      [][]string
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
