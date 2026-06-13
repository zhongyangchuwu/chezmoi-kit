package main

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

func TestRunStatusRendersRecommendation(t *testing.T) {
	service := &fakeService{entries: []chezmoi.StatusEntry{{LocalChange: chezmoi.ChangeModified, TargetChange: chezmoi.ChangeNone, Path: ".zshrc"}}}
	var out bytes.Buffer

	code := run([]string{"status"}, service, strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	got := out.String()
	for _, want := range []string{"M  .zshrc", "local changed", "cm sync .zshrc"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
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

type fakeService struct {
	entries     []chezmoi.StatusEntry
	statusCalls int
	statusArgs  [][]string
	commands    [][]string
}

func (f *fakeService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	f.statusCalls++
	f.statusArgs = append(f.statusArgs, append([]string(nil), targets...))
	return append([]chezmoi.StatusEntry(nil), f.entries...), nil
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
