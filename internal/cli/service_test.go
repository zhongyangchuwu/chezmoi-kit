package cli

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
)

func TestChezmoiServiceSourceStatusUsesRunnerInSourceDir(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("/home/me/.local/share/chezmoi\n"), []byte(" M dot_zshrc\n")},
	}
	service := chezmoiService{client: chezmoi.Client{Runner: runner}}

	entries, err := service.SourceStatus()
	if err != nil {
		t.Fatalf("SourceStatus returned error: %v", err)
	}
	wantEntries := []sourceEntry{{Code: " M", Path: "dot_zshrc"}}
	if !reflect.DeepEqual(entries, wantEntries) {
		t.Fatalf("entries = %#v, want %#v", entries, wantEntries)
	}
	wantCalls := [][]string{{"chezmoi", "source-path"}, {"git", "status", "--porcelain=v1"}}
	if !reflect.DeepEqual(runner.outputCalls, wantCalls) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, wantCalls)
	}
	if got := runner.outputIO[1].Dir; got != "/home/me/.local/share/chezmoi" {
		t.Fatalf("git dir = %q", got)
	}
}

func TestChezmoiServiceOpenSourceGitUsesRunnerIO(t *testing.T) {
	runner := &recordingRunner{outputs: [][]byte{[]byte("/home/me/src\n")}}
	stdin := bytes.NewBufferString("in")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	service := chezmoiService{client: chezmoi.Client{Runner: runner, Stdin: stdin, Stdout: &stdout, Stderr: &stderr}}

	if err := service.OpenSourceGit(); err != nil {
		t.Fatalf("OpenSourceGit returned error: %v", err)
	}
	wantRuns := [][]string{{"lazygit"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
	gotIO := runner.runIO[0]
	if gotIO.Dir != "/home/me/src" || gotIO.Stdin != stdin || gotIO.Stdout != &stdout || gotIO.Stderr != &stderr {
		t.Fatalf("run IO = %#v", gotIO)
	}
}

func TestChezmoiServiceExecuteForcesConfirmedApply(t *testing.T) {
	runner := &recordingRunner{}
	service := chezmoiService{client: chezmoi.Client{Runner: runner}}

	err := service.Execute([]reconcile.Action{{Target: "/home/me/.zshrc", Kind: reconcile.ActionApply}})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	wantRuns := [][]string{{"chezmoi", "apply", "--force", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
}

func TestChezmoiServiceEditTargetRunsEditWithAbsolutePath(t *testing.T) {
	t.Setenv("HOME", "/home/testuser")
	runner := &recordingRunner{}
	service := chezmoiService{client: chezmoi.Client{Runner: runner}}

	if err := service.EditTarget(".zshrc"); err != nil {
		t.Fatalf("EditTarget returned error: %v", err)
	}

	wantRuns := [][]string{{"chezmoi", "edit", "/home/testuser/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
}

type recordingRunner struct {
	outputs     [][]byte
	outputCalls [][]string
	runCalls    [][]string
	outputIO    []chezmoi.RunnerIO
	runIO       []chezmoi.RunnerIO
}

func (r *recordingRunner) Output(command string, args []string, io chezmoi.RunnerIO) ([]byte, error) {
	r.outputCalls = append(r.outputCalls, append([]string{command}, args...))
	r.outputIO = append(r.outputIO, io)
	if len(r.outputs) == 0 {
		return nil, nil
	}
	out := bytes.Clone(r.outputs[0])
	r.outputs = r.outputs[1:]
	return out, nil
}

func (r *recordingRunner) Run(command string, args []string, io chezmoi.RunnerIO) error {
	r.runCalls = append(r.runCalls, append([]string{command}, args...))
	r.runIO = append(r.runIO, io)
	return nil
}
