package app

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
	"github.com/zhongyangchuwu/cm/internal/report"
)

// StatusService builds the read-only reconciliation status use case.
type StatusService interface {
	StatusReport(targets []string) (report.Document, error)
}

// DiffService builds internal target-vs-local diff output.
type DiffService interface {
	DiffReport(targets []string) (report.Document, error)
}

// TargetCommandService executes direct chezmoi target wrappers.
type TargetCommandService interface {
	AddTargets(targets []string) error
	ApplyTargets(targets []string) error
	MergeTargets(targets []string) error
}

// SourceGitService opens the chezmoi source repository in a terminal UI.
type SourceGitService interface {
	OpenSourceGit() error
}

// EditService edits chezmoi-managed targets and exposes completion candidates.
type EditService interface {
	EditTarget(target string) error
	ManagedFiles() ([]string, error)
}

// ActionKind identifies a reconciliation action selected during review.
type ActionKind int

const (
	ActionAdd ActionKind = iota
	ActionApply
	ActionMerge
)

func (k ActionKind) String() string {
	switch k {
	case ActionAdd:
		return "add"
	case ActionApply:
		return "apply"
	case ActionMerge:
		return "merge"
	default:
		return "unknown"
	}
}

func (k ActionKind) Marker() string {
	switch k {
	case ActionAdd:
		return "A"
	case ActionApply:
		return "P"
	case ActionMerge:
		return "M"
	default:
		return "?"
	}
}

// Action is a confirmed reconciliation action for one managed target.
type Action struct {
	Target string
	Kind   ActionKind
}

// TerminalCommand is a reconciliation command that temporarily owns the terminal.
type TerminalCommand interface {
	Run() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

// SyncService is the application boundary for reviewing and executing sync actions.
type SyncService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	DiffOutput(target string) ([]byte, error)
	ExecuteNonInteractive(action Action) error
	TerminalCommand(action Action) (TerminalCommand, error)
}

// Services groups application use cases for adapters.
type Services struct {
	Status    StatusService
	Diff      DiffService
	Sync      SyncService
	Target    TargetCommandService
	SourceGit SourceGitService
	Edit      EditService
}

// NewServices creates the default application service graph backed by chezmoi.
func NewServices(client chezmoi.Client) Services {
	s := service{client: client}
	return Services{
		Status:    s,
		Diff:      s,
		Sync:      s,
		Target:    s,
		SourceGit: s,
		Edit:      s,
	}
}

type service struct {
	client chezmoi.Client
}

func (s service) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	return s.client.Status(targets)
}

func (s service) StatusReport(targets []string) (report.Document, error) {
	entries, err := s.Status(targets)
	if err != nil {
		return report.Document{}, err
	}
	sourceEntries, err := s.SourceStatus()
	if err != nil {
		return report.Document{}, err
	}

	var doc report.Document
	if len(entries) == 0 && len(sourceEntries) == 0 {
		doc.Blocks = append(doc.Blocks, report.Paragraph(report.Text("clean")))
		return doc, nil
	}
	if len(entries) > 0 {
		doc.Blocks = append(doc.Blocks,
			report.Heading(1, report.Strong("local:")),
			report.Paragraph(report.Text("  local config differs from chezmoi source")),
		)
		for _, entry := range entries {
			doc.Blocks = append(doc.Blocks, report.Paragraph(
				report.Warning("!"),
				report.Text(" "),
				report.Path(entry.Path),
				report.Text("  differs from chezmoi"),
			))
		}
		doc.Blocks = append(doc.Blocks, report.Paragraph(report.Muted("  run cm sync")))
	}
	if len(sourceEntries) > 0 {
		if len(entries) > 0 {
			doc.Blocks = append(doc.Blocks, report.Blank())
		}
		doc.Blocks = append(doc.Blocks,
			report.Heading(1, report.Strong("chezmoi:")),
			report.Paragraph(report.Text("  source repository has git changes")),
		)
		for _, entry := range sourceEntries {
			doc.Blocks = append(doc.Blocks, report.Paragraph(
				report.Status(entry.Code),
				report.Text(" "),
				report.Path(entry.Path),
			))
		}
	}
	return doc, nil
}

func (s service) SourceStatus() ([]chezmoi.StatusEntry, error) {
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

func (s service) OpenSourceGit() error {
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

func (s service) sourceDir() (string, error) {
	out, err := s.client.Output("source-path")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s service) DiffReport(targets []string) (report.Document, error) {
	if len(targets) == 0 {
		entries, err := s.Status(nil)
		if err != nil {
			return report.Document{}, err
		}
		if len(entries) == 0 {
			return report.Document{Blocks: []report.Block{report.Paragraph(report.Text("clean"))}}, nil
		}
		targets = make([]string, 0, len(entries))
		for _, entry := range entries {
			targets = append(targets, entry.Path)
		}
	}

	doc := report.Document{Blocks: make([]report.Block, 0, len(targets))}
	for _, target := range targets {
		diff, err := s.DiffOutput(target)
		if err != nil {
			return report.Document{}, err
		}
		doc.Blocks = append(doc.Blocks, report.DiffBlock(string(diff)))
	}
	return doc, nil
}

func (s service) DiffOutput(target string) ([]byte, error) {
	differ := diff.Differ{Source: chezmoi.ContentLoader{Client: s.client}}
	return differ.Diff(target)
}

func (s service) ExecuteNonInteractive(action Action) error {
	switch action.Kind {
	case ActionAdd:
		return s.runBuffered("add", action.Target)
	case ActionApply:
		return s.runBuffered("apply", "--force", action.Target)
	case ActionMerge:
		return fmt.Errorf("merge requires terminal execution")
	default:
		return fmt.Errorf("unknown reconcile action %d for %s", action.Kind, action.Target)
	}
}

func (s service) TerminalCommand(action Action) (TerminalCommand, error) {
	if action.Kind != ActionMerge {
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

func (s service) AddTargets(targets []string) error {
	return s.runTargets("add", targets)
}

func (s service) ApplyTargets(targets []string) error {
	return s.runTargets("apply", targets)
}

func (s service) MergeTargets(targets []string) error {
	return s.runTargets("merge", targets)
}

func (s service) EditTarget(target string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return s.client.Run("edit", filepath.Join(home, target))
}

func (s service) ManagedFiles() ([]string, error) {
	return s.client.ManagedFiles()
}

func (s service) runner() process.Runner {
	return s.client.ActiveRunner()
}

func (s service) runBuffered(args ...string) error {
	_, stderr, err := s.client.RunBuffered(args...)
	if err != nil {
		return fmt.Errorf("%w%s", err, formatStderr(stderr))
	}
	return nil
}

func (s service) runTargets(command string, targets []string) error {
	args := make([]string, 0, 1+len(targets))
	args = append(args, command)
	args = append(args, targets...)
	return s.client.Run(args...)
}

func formatStderr(stderr []byte) string {
	stderr = bytes.TrimSpace(stderr)
	if len(stderr) == 0 {
		return ""
	}
	return ": " + string(stderr)
}
