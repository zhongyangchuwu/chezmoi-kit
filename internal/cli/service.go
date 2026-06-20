package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/diff"
	"github.com/zhongyangchuwu/cm/internal/process"
	"github.com/zhongyangchuwu/cm/internal/tui"
)

type statusService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	SourceStatus() ([]chezmoi.StatusEntry, error)
}

type diffService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
}

type targetCommandService interface {
	AddTargets(targets []string) error
	ApplyTargets(targets []string) error
	MergeTargets(targets []string) error
}

type sourceGitService interface {
	OpenSourceGit() error
}

type editService interface {
	EditTarget(target string) error
	ManagedFiles() ([]string, error)
}

type commandServices struct {
	Status    statusService
	Diff      diffService
	Sync      tui.ReviewService
	Target    targetCommandService
	SourceGit sourceGitService
	Edit      editService
}

func commandServicesFor(s chezmoiService) commandServices {
	return commandServices{
		Status:    s,
		Diff:      s,
		Sync:      s,
		Target:    s,
		SourceGit: s,
		Edit:      s,
	}
}

type chezmoiService struct {
	client chezmoi.Client
}

func (s chezmoiService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	return s.client.Status(targets)
}

func (s chezmoiService) SourceStatus() ([]chezmoi.StatusEntry, error) {
	sourceDir, err := s.sourceDir()
	if err != nil {
		return nil, err
	}
	if sourceDir == "" {
		return nil, nil
	}

	var stderr bytes.Buffer
	out, err := s.runner().Output("git", []string{"status", "--porcelain=v1"}, process.IO{Dir: sourceDir, Stderr: &stderr})
	if err != nil {
		return nil, fmt.Errorf("git status %s: %w%s", sourceDir, err, formatStderr(stderr.Bytes()))
	}
	return chezmoi.ParseGitStatus(out)
}

func (s chezmoiService) OpenSourceGit() error {
	sourceDir, err := s.sourceDir()
	if err != nil {
		return err
	}
	if sourceDir == "" {
		return fmt.Errorf("chezmoi source path is empty")
	}

	if err := s.runner().Run("lazygit", nil, process.IO{
		Dir:    sourceDir,
		Stdin:  s.client.Stdin,
		Stdout: s.client.Stdout,
		Stderr: s.client.Stderr,
	}); err != nil {
		return fmt.Errorf("lazygit %s: %w", sourceDir, err)
	}
	return nil
}

func (s chezmoiService) sourceDir() (string, error) {
	out, err := s.client.Output("source-path")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s chezmoiService) DiffOutput(target string) ([]byte, error) {
	differ := diff.Differ{Source: chezmoi.ContentLoader{Client: s.client}}
	return differ.Diff(target)
}

func (s chezmoiService) ExecuteNonInteractive(action tui.Action) error {
	switch action.Kind {
	case tui.ActionAdd:
		return s.runBuffered("add", action.Target)
	case tui.ActionApply:
		return s.runBuffered("apply", "--force", action.Target)
	case tui.ActionMerge:
		return fmt.Errorf("merge requires terminal execution")
	default:
		return fmt.Errorf("unknown reconcile action %d for %s", action.Kind, action.Target)
	}
}

func (s chezmoiService) TerminalCommand(action tui.Action) (tui.TerminalCommand, error) {
	if action.Kind != tui.ActionMerge {
		return nil, fmt.Errorf("%s does not require terminal execution", action.Kind)
	}
	return &runnerCommand{runner: s.runner(), command: s.client.BinaryName(), args: []string{"merge", action.Target}, io: process.IO{Dir: s.client.Dir}}, nil
}

type runnerCommand struct {
	runner  process.Runner
	command string
	args    []string
	io      process.IO
}

func (c *runnerCommand) Run() error {
	return c.runner.Run(c.command, c.args, c.io)
}

func (c *runnerCommand) SetStdin(r io.Reader) {
	if c.io.Stdin == nil {
		c.io.Stdin = r
	}
}

func (c *runnerCommand) SetStdout(w io.Writer) {
	if c.io.Stdout == nil {
		c.io.Stdout = w
	}
}

func (c *runnerCommand) SetStderr(w io.Writer) {
	if c.io.Stderr == nil {
		c.io.Stderr = w
	}
}

func (s chezmoiService) runner() process.Runner {
	return s.client.ActiveRunner()
}

func formatStderr(stderr []byte) string {
	stderr = bytes.TrimSpace(stderr)
	if len(stderr) == 0 {
		return ""
	}
	return ": " + string(stderr)
}

func (s chezmoiService) AddTargets(targets []string) error {
	return s.runTargets("add", targets)
}

func (s chezmoiService) ApplyTargets(targets []string) error {
	return s.runTargets("apply", targets)
}

func (s chezmoiService) MergeTargets(targets []string) error {
	return s.runTargets("merge", targets)
}

func (s chezmoiService) EditTarget(target string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return s.client.Run("edit", filepath.Join(home, target))
}

func (s chezmoiService) ManagedFiles() ([]string, error) {
	return s.client.ManagedFiles()
}

func (s chezmoiService) runBuffered(args ...string) error {
	_, stderr, err := s.client.RunBuffered(args...)
	if err != nil {
		return fmt.Errorf("%w%s", err, formatStderr(stderr))
	}
	return nil
}

func (s chezmoiService) runTargets(command string, targets []string) error {
	args := make([]string, 0, 1+len(targets))
	args = append(args, command)
	args = append(args, targets...)
	return s.client.Run(args...)
}
