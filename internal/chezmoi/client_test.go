package chezmoi

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/process"
)

func TestClientStatusRequestsAbsolutePathsAndParsesOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("M  /home/me/.zshrc\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	entries, err := client.Status([]string{".zshrc"})
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.outputCalls, [][]string{{"chezmoi", "status", "--include=all", "--exclude=none", "--path-style=absolute", ".zshrc"}}) {
		t.Fatalf("outputCalls = %#v", runner.outputCalls)
	}
	want := []StatusEntry{{Code: "M ", Path: "/home/me/.zshrc"}}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestClientStatusForInventorySkipsSecretTemplates(t *testing.T) {
	runner := &fakeRunner{}
	client := Client{Binary: "chezmoi", Runner: runner}

	if _, err := client.StatusForInventory(nil); err != nil {
		t.Fatalf("StatusForInventory returned error: %v", err)
	}
	want := [][]string{{"chezmoi", "--skip-secrets", "status", "--include=all", "--exclude=none", "--path-style=absolute"}}
	if !reflect.DeepEqual(runner.outputCalls, want) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, want)
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

func TestClientRunBufferedCapturesOutputAndDisablesStdin(t *testing.T) {
	runner := &fakeRunner{output: []byte("stdout")}
	client := Client{Binary: "chezmoi", Runner: runner, Stdin: bytes.NewBufferString("ignored"), Dir: "/work"}

	stdout, stderr, err := client.RunBuffered("add", ".zshrc")
	if err != nil {
		t.Fatalf("RunBuffered returned error: %v", err)
	}
	if string(stdout) != "stdout" || len(stderr) != 0 {
		t.Fatalf("stdout=%q stderr=%q", stdout, stderr)
	}
	if !reflect.DeepEqual(runner.runCalls, [][]string{{"chezmoi", "add", ".zshrc"}}) {
		t.Fatalf("runCalls = %#v", runner.runCalls)
	}
	gotIO := runner.runIO[0]
	if gotIO.Stdin != nil || gotIO.Stdout == nil || gotIO.Stderr == nil || gotIO.Dir != "/work" {
		t.Fatalf("run IO = %#v, want buffered output, nil stdin, preserved dir", gotIO)
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

func TestClientManagedFilesReturnsNULDelimitedOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte(".zshrc\x00.gitconfig\x00")}
	client := Client{Binary: "chezmoi", Runner: runner}

	files, err := client.ManagedFiles()
	if err != nil {
		t.Fatalf("ManagedFiles returned error: %v", err)
	}

	if !reflect.DeepEqual(runner.outputCalls, [][]string{{"chezmoi", "managed", "--nul-path-separator"}}) {
		t.Fatalf("outputCalls = %#v", runner.outputCalls)
	}
	want := []string{".zshrc", ".gitconfig"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("files = %#v, want %#v", files, want)
	}
}

type fakeRunner struct {
	output      []byte
	outputs     [][]byte
	err         error
	outputCalls [][]string
	runCalls    [][]string
	runIO       []process.IO
}

func (f *fakeRunner) Output(command string, args []string, _ process.IO) ([]byte, error) {
	call := append([]string{command}, args...)
	f.outputCalls = append(f.outputCalls, call)
	if len(f.outputs) > 0 {
		out := bytes.Clone(f.outputs[0])
		f.outputs = f.outputs[1:]
		return out, f.err
	}
	return bytes.Clone(f.output), f.err
}

func (f *fakeRunner) Run(command string, args []string, io process.IO) error {
	call := append([]string{command}, args...)
	f.runCalls = append(f.runCalls, call)
	f.runIO = append(f.runIO, io)
	out := f.output
	if len(f.outputs) > 0 {
		out = f.outputs[0]
		f.outputs = f.outputs[1:]
	}
	if io.Stdout != nil {
		_, err := io.Stdout.Write(out)
		if err != nil {
			return err
		}
	}
	return f.err
}
