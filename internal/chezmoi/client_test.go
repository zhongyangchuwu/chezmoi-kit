package chezmoi

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestClientStatusRequestsAbsolutePathsAndParsesOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("M  /home/me/.zshrc\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	entries, err := client.Status([]string{".zshrc"})
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.outputCalls, [][]string{{"chezmoi", "status", "--path-style=absolute", ".zshrc"}}) {
		t.Fatalf("outputCalls = %#v", runner.outputCalls)
	}
	want := []StatusEntry{{Code: "M ", Path: "/home/me/.zshrc"}}
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
func TestClientOutputPreservesTrailingNewlines(t *testing.T) {
	runner := &fakeRunner{output: []byte("content\n\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	out, err := client.Output("cat", ".zshrc")
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	if string(out) != "content\n\n" {
		t.Fatalf("output = %q, want trailing newlines preserved", string(out))
	}
}

func TestClientOutputLimitCapsCapturedBytes(t *testing.T) {
	runner := &fakeRunner{output: []byte("abcdef")}
	client := Client{Binary: "chezmoi", Runner: runner}

	out, err := client.OutputLimit(3, "cat", ".zshrc")
	if err != nil {
		t.Fatalf("OutputLimit returned error: %v", err)
	}
	if string(out) != "abcd" {
		t.Fatalf("output = %q, want capped bytes plus sentinel", string(out))
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

func TestClientManagedFilesReturnsParsedOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte(".zshrc\n.gitconfig\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	files, err := client.ManagedFiles()
	if err != nil {
		t.Fatalf("ManagedFiles returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.outputCalls, [][]string{{"chezmoi", "managed"}}) {
		t.Fatalf("outputCalls = %#v", runner.outputCalls)
	}
	want := []string{".zshrc", ".gitconfig"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("files = %#v, want %#v", files, want)
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

func (f *fakeRunner) Run(command string, args []string, io RunnerIO) error {
	call := append([]string{command}, args...)
	f.runCalls = append(f.runCalls, call)
	if io.Stdout != nil {
		_, err := io.Stdout.Write(f.output)
		if err != nil {
			return err
		}
	}
	return f.err
}
