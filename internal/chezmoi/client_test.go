package chezmoi

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestClientStatusExecutesChezmoiStatusAndParsesOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("M  .zshrc\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	entries, err := client.Status([]string{".zshrc"})
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.outputCalls, [][]string{{"chezmoi", "status", ".zshrc"}}) {
		t.Fatalf("outputCalls = %#v", runner.outputCalls)
	}
	want := []StatusEntry{{LocalChange: ChangeModified, TargetChange: ChangeNone, Path: ".zshrc"}}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestClientRunForwardsCommand(t *testing.T) {
	runner := &fakeRunner{}
	client := Client{Binary: "chezmoi", Runner: runner}

	if err := client.Run("diff", ".zshrc"); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.runCalls, [][]string{{"chezmoi", "diff", ".zshrc"}}) {
		t.Fatalf("runCalls = %#v", runner.runCalls)
	}
}

func TestClientStatusWrapsRunnerError(t *testing.T) {
	runner := &fakeRunner{err: errors.New("boom")}
	client := Client{Binary: "chezmoi", Runner: runner}

	_, err := client.Status(nil)
	if err == nil {
		t.Fatal("Status returned nil error")
	}
}

type fakeRunner struct {
	output      []byte
	err         error
	outputCalls [][]string
	runCalls    [][]string
}

func (f *fakeRunner) Output(command string, args []string, _ RunnerIO) ([]byte, error) {
	call := append([]string{command}, args...)
	f.outputCalls = append(f.outputCalls, call)
	return bytes.Clone(f.output), f.err
}

func (f *fakeRunner) Run(command string, args []string, _ RunnerIO) error {
	call := append([]string{command}, args...)
	f.runCalls = append(f.runCalls, call)
	return f.err
}
