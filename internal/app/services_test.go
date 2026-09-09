package app

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/process"
	"github.com/zhongyangchuwu/cm/internal/report"
)

func TestServiceStatusReportUsesSemanticTokens(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("MM /home/me/.zshrc\n"), []byte("/home/me/src\n"), []byte(" M dot_zshrc\n")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, err := service.StatusReport(nil)
	if err != nil {
		t.Fatalf("StatusReport returned error: %v", err)
	}
	got := string(report.Plain(doc))
	for _, want := range []string{"local:", "! /home/me/.zshrc  differs from chezmoi", "  run cm sync", "chezmoi:", " M dot_zshrc"} {
		if !strings.Contains(got, want) {
			t.Fatalf("status report %q missing %q", got, want)
		}
	}
}

func TestServiceStatusReportSeparatesPendingScripts(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("MM /home/me/.zshrc\n R /home/me/install.sh\n"), []byte("/home/me/src\n"), nil},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, err := service.StatusReport(nil)
	if err != nil {
		t.Fatalf("StatusReport returned error: %v", err)
	}
	got := string(report.Plain(doc))
	for _, want := range []string{"local:", "automation:", "R /home/me/install.sh  apply would run this script", "cm sync does not execute scripts"} {
		if !strings.Contains(got, want) {
			t.Fatalf("status report %q missing %q", got, want)
		}
	}
}

func TestServiceSourceStatusUsesRunnerInSourceDir(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("/home/me/.local/share/chezmoi\n"), []byte(" M dot_zshrc\n")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	entries, err := service.SourceStatus()
	if err != nil {
		t.Fatalf("SourceStatus returned error: %v", err)
	}
	wantEntries := []chezmoi.StatusEntry{{Code: " M", Path: "dot_zshrc"}}
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

func TestServiceSourceStatusRejectsMalformedGitStatus(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("/home/me/.local/share/chezmoi\n"), []byte("warning\n")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	_, err := service.SourceStatus()
	if err == nil {
		t.Fatal("SourceStatus returned nil error, want malformed status error")
	}
	if got := err.Error(); got != "malformed git status line 1: \"warning\"" {
		t.Fatalf("error = %q", got)
	}
}

func TestServiceOpenSourceGitUsesRunnerIO(t *testing.T) {
	runner := &recordingRunner{outputs: [][]byte{[]byte("/home/me/src\n")}}
	stdin := bytes.NewBufferString("in")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	service := service{client: chezmoi.Client{Runner: runner, Stdin: stdin, Stdout: &stdout, Stderr: &stderr}}

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

func TestServiceDiffExpandsDirtyTargets(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{[]byte("MM /home/me/.zshrc\n M /home/me/.gitconfig\n")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, err := service.DiffReport(nil)
	if err != nil {
		t.Fatalf("DiffReport returned error: %v", err)
	}
	got := string(report.Plain(doc))
	for _, want := range []string{"no diff: /home/me/.zshrc", "no diff: /home/me/.gitconfig"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diff output %q does not contain %q", got, want)
		}
	}
	wantOutputCalls := [][]string{{"chezmoi", "status", "--include=all", "--exclude=none", "--path-style=absolute"}}
	if !reflect.DeepEqual(runner.outputCalls, wantOutputCalls) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, wantOutputCalls)
	}
	wantRunCalls := [][]string{
		{"chezmoi", "--color=false", "--no-pager", "--use-builtin-diff", "diff", "--include=all", "--exclude=none", "--reverse", "--script-contents=true", "/home/me/.zshrc"},
		{"chezmoi", "--color=false", "--no-pager", "--use-builtin-diff", "diff", "--include=all", "--exclude=none", "--reverse", "--script-contents=true", "/home/me/.gitconfig"},
	}
	if !reflect.DeepEqual(runner.runCalls, wantRunCalls) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRunCalls)
	}
}

func TestServiceExecuteNonInteractiveBuffersReAddApply(t *testing.T) {
	runner := &recordingRunner{}
	service := service{client: chezmoi.Client{Runner: runner}}

	if _, err := service.ExecuteNonInteractive(Action{Target: "/home/me/.zshrc", Kind: ActionAdd}); err != nil {
		t.Fatalf("ExecuteNonInteractive re-add returned error: %v", err)
	}
	if _, err := service.ExecuteNonInteractive(Action{Target: "/home/me/.gitconfig", Kind: ActionApply}); err != nil {
		t.Fatalf("ExecuteNonInteractive apply returned error: %v", err)
	}

	wantRuns := [][]string{
		{"chezmoi", "re-add", "/home/me/.zshrc"},
		{"chezmoi", "apply", "--force", "/home/me/.gitconfig"},
	}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
	for i, gotIO := range runner.runIO {
		if gotIO.Stdin != nil || gotIO.Stdout == nil || gotIO.Stderr == nil {
			t.Fatalf("runIO[%d] = %#v, want nil stdin and buffered stdout/stderr", i, gotIO)
		}
	}
}

func TestServiceExecuteNonInteractiveReturnsSuccessfulOutput(t *testing.T) {
	runner := &recordingRunner{runStdouts: [][]byte{[]byte("updated\n")}, runStderrs: [][]byte{[]byte("warning\n")}}
	service := service{client: chezmoi.Client{Runner: runner}}

	result, err := service.ExecuteNonInteractive(Action{Target: "/home/me/.zshrc", Kind: ActionAdd})
	if err != nil {
		t.Fatalf("ExecuteNonInteractive returned error: %v", err)
	}
	if result.Stdout != "updated\n" || result.Stderr != "warning\n" {
		t.Fatalf("result = %#v", result)
	}
}

func TestServiceTerminalCommandBuildsMergeCommand(t *testing.T) {
	runner := &recordingRunner{}
	service := service{client: chezmoi.Client{Runner: runner, Dir: "/work"}}

	cmd, err := service.TerminalCommand(Action{Target: "/home/me/.tmux.conf", Kind: ActionMerge})
	if err != nil {
		t.Fatalf("TerminalCommand returned error: %v", err)
	}
	if err := cmd.Run(); err != nil {
		t.Fatalf("terminal command Run returned error: %v", err)
	}

	wantRuns := [][]string{{"chezmoi", "merge", "/home/me/.tmux.conf"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
	if got := runner.runIO[0].Dir; got != "/work" {
		t.Fatalf("Dir = %q, want /work", got)
	}
}

func TestServiceSourceEditCommandForcesNoApplyOrWatch(t *testing.T) {
	runner := &recordingRunner{}
	service := service{client: chezmoi.Client{Runner: runner, Dir: "/work"}}
	entry := WorkspaceEntry{
		Path:       "/home/me/.zshrc",
		SourcePath: "/home/me/source/dot_zshrc",
		State:      FileDirty,
		Type:       TargetFile,
	}

	command, err := service.SourceEditCommand(entry)
	if err != nil {
		t.Fatalf("SourceEditCommand returned error: %v", err)
	}
	if err := command.Run(); err != nil {
		t.Fatalf("source edit command Run returned error: %v", err)
	}
	if want := [][]string{{"chezmoi", "edit", "--apply=false", "--watch=false", "/home/me/.zshrc"}}; !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
	if got := runner.runIO[0].Dir; got != "/work" {
		t.Fatalf("Dir = %q, want /work", got)
	}
}

func TestSourceEditUnavailableReasonRejectsUnsupportedWorkspaceEntries(t *testing.T) {
	tests := []struct {
		name  string
		entry WorkspaceEntry
		want  string
	}{
		{"unmanaged", WorkspaceEntry{State: FileUnmanaged, Type: TargetFile}, "unmanaged target has no chezmoi source"},
		{"directory", WorkspaceEntry{SourcePath: "/source", Type: TargetDirectory}, "select a managed file or symlink"},
		{"script", WorkspaceEntry{SourcePath: "/source", State: FileScript, Type: TargetScript}, "scripts cannot be edited through workspace source edit"},
		{"external", WorkspaceEntry{SourcePath: "/source", Type: TargetExternal}, "source edit unsupported for external target"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := SourceEditUnavailableReason(test.entry); got != test.want {
				t.Fatalf("SourceEditUnavailableReason() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSourceEditUnavailableReasonAllowsManagedFileAndSymlinkVariants(t *testing.T) {
	for _, entry := range []WorkspaceEntry{
		{SourcePath: "/source/template", State: FileDirty, Type: TargetFile, Template: true},
		{SourcePath: "/source/encrypted", State: FileUninspected, Type: TargetFile, Encrypted: true},
		{SourcePath: "/source/link", State: FileClean, Type: TargetSymlink},
	} {
		if reason := SourceEditUnavailableReason(entry); reason != "" {
			t.Fatalf("SourceEditUnavailableReason(%#v) = %q, want allowed", entry, reason)
		}
	}
}

func TestServiceEditTargetUsesConfiguredDestination(t *testing.T) {
	runner := &recordingRunner{outputs: [][]byte{[]byte("/srv/dotfiles-home\n")}}
	service := service{client: chezmoi.Client{Runner: runner}}

	if err := service.EditTarget(".zshrc"); err != nil {
		t.Fatalf("EditTarget returned error: %v", err)
	}

	wantOutputs := [][]string{{"chezmoi", "target-path"}}
	if !reflect.DeepEqual(runner.outputCalls, wantOutputs) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, wantOutputs)
	}
	wantRuns := [][]string{{"chezmoi", "edit", "/srv/dotfiles-home/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
}

func TestServiceEditTargetPreservesAbsolutePath(t *testing.T) {
	runner := &recordingRunner{}
	service := service{client: chezmoi.Client{Runner: runner}}

	if err := service.EditTarget("/srv/dotfiles-home/.zshrc"); err != nil {
		t.Fatalf("EditTarget returned error: %v", err)
	}
	if len(runner.outputCalls) != 0 {
		t.Fatalf("outputCalls = %#v, want none", runner.outputCalls)
	}
	wantRuns := [][]string{{"chezmoi", "edit", "/srv/dotfiles-home/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
}

func TestServiceDirectTargetWrappersForwardTargets(t *testing.T) {
	runner := &recordingRunner{}
	service := service{client: chezmoi.Client{Runner: runner}}

	if err := service.AddTargets([]string{".zshrc"}); err != nil {
		t.Fatalf("AddTargets returned error: %v", err)
	}
	if err := service.ApplyTargets([]string{".gitconfig"}); err != nil {
		t.Fatalf("ApplyTargets returned error: %v", err)
	}
	if err := service.MergeTargets([]string{".config/nvim/init.lua"}); err != nil {
		t.Fatalf("MergeTargets returned error: %v", err)
	}

	wantRuns := [][]string{
		{"chezmoi", "add", ".zshrc"},
		{"chezmoi", "apply", ".gitconfig"},
		{"chezmoi", "merge", ".config/nvim/init.lua"},
	}
	if !reflect.DeepEqual(runner.runCalls, wantRuns) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRuns)
	}
}

func TestServiceDoctorReportPassesWithOptionalWarning(t *testing.T) {
	runner := &recordingRunner{
		outputs: [][]byte{
			[]byte("chezmoi version v2\n"),
			[]byte("git version 2\n"),
			[]byte("/home/me/src\n"),
			[]byte(""),
		},
		outputErrors: []error{nil, nil, nil, nil, errors.New("lazygit not found in PATH")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, failed := service.DoctorReport()

	if failed {
		t.Fatalf("DoctorReport failed, want optional warning only")
	}
	got := string(report.Plain(doc))
	for _, want := range []string{"doctor:", "pass chezmoi: available", "pass git: available", "pass source-path: /home/me/src", "pass source-git: repository readable", "warn lazygit: lazygit not found in PATH"} {
		if !strings.Contains(got, want) {
			t.Fatalf("DoctorReport() = %q, missing %q", got, want)
		}
	}
	if len(runner.runCalls) != 0 {
		t.Fatalf("runCalls = %#v, want none", runner.runCalls)
	}
}

func TestServiceDoctorReportFailsForRequiredTool(t *testing.T) {
	runner := &recordingRunner{
		outputs:      [][]byte{[]byte("git version 2\n"), []byte("/home/me/src\n"), []byte(""), []byte("lazygit version\n")},
		outputErrors: []error{errors.New("chezmoi not found in PATH")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, failed := service.DoctorReport()

	if !failed {
		t.Fatalf("DoctorReport failed = false, want true")
	}
	got := string(report.Plain(doc))
	if !strings.Contains(got, "fail chezmoi: chezmoi not found in PATH") {
		t.Fatalf("DoctorReport() = %q, want chezmoi failure", got)
	}
}

func TestServiceDoctorReportFailsForSourceGit(t *testing.T) {
	runner := &recordingRunner{
		outputs:      [][]byte{[]byte("chezmoi version v2\n"), []byte("git version 2\n"), []byte("/home/me/src\n"), []byte("lazygit version\n")},
		outputErrors: []error{nil, nil, nil, errors.New("not a git repository")},
	}
	service := service{client: chezmoi.Client{Runner: runner}}

	doc, failed := service.DoctorReport()

	if !failed {
		t.Fatalf("DoctorReport failed = false, want true")
	}
	got := string(report.Plain(doc))
	if !strings.Contains(got, "fail source-git: git status /home/me/src: not a git repository") {
		t.Fatalf("DoctorReport() = %q, want source git failure", got)
	}
}

type recordingRunner struct {
	outputs      [][]byte
	outputErrors []error
	runStdouts   [][]byte
	runStderrs   [][]byte
	runErrors    []error
	outputCalls  [][]string
	runCalls     [][]string
	outputIO     []process.IO
	runIO        []process.IO
}

func (r *recordingRunner) Output(command string, args []string, io process.IO) ([]byte, error) {
	r.outputCalls = append(r.outputCalls, append([]string{command}, args...))
	r.outputIO = append(r.outputIO, io)
	var err error
	if len(r.outputErrors) > 0 {
		err = r.outputErrors[0]
		r.outputErrors = r.outputErrors[1:]
	}
	if len(r.outputs) == 0 {
		return nil, err
	}
	out := bytes.Clone(r.outputs[0])
	r.outputs = r.outputs[1:]
	return out, err
}

func (r *recordingRunner) Run(command string, args []string, io process.IO) error {
	r.runCalls = append(r.runCalls, append([]string{command}, args...))
	r.runIO = append(r.runIO, io)
	var stdout, stderr []byte
	if len(r.runStdouts) > 0 {
		stdout = r.runStdouts[0]
		r.runStdouts = r.runStdouts[1:]
	}
	if len(r.runStderrs) > 0 {
		stderr = r.runStderrs[0]
		r.runStderrs = r.runStderrs[1:]
	}
	if io.Stdout != nil {
		_, _ = io.Stdout.Write(stdout)
	}
	if io.Stderr != nil {
		_, _ = io.Stderr.Write(stderr)
	}
	if len(r.runErrors) == 0 {
		return nil
	}
	err := r.runErrors[0]
	r.runErrors = r.runErrors[1:]
	return err
}
