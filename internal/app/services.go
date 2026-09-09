package app

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/process"
	"github.com/zhongyangchuwu/cm/internal/report"
)

const maxReviewOutput int64 = 1 << 20

// StatusService builds the read-only reconciliation status use case.
type StatusService interface {
	StatusReport(targets []string) (report.Document, error)
}

// DiffService builds authoritative target-vs-local diff output.
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

// DoctorService builds read-only environment diagnostics.
type DoctorService interface {
	DoctorReport() (report.Document, bool)
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

// Action is a confirmed reconciliation action for one reviewed target state.
type Action struct {
	Target      string
	Kind        ActionKind
	Fingerprint string
}

// TerminalCommand is a command that temporarily owns the terminal.
type TerminalCommand interface {
	Run() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

// SyncService is the application boundary for reviewing and executing sync actions.
type SyncService interface {
	Status(targets []string) (SyncStatus, error)
	Review(target string) (Review, error)
	ExecuteNonInteractive(action Action) (ActionResult, error)
	TerminalCommand(action Action) (TerminalCommand, error)
}

// Services groups application use cases for adapters.
type Services struct {
	Status    StatusService
	Diff      DiffService
	Sync      SyncService
	Workspace WorkspaceService
	Target    TargetCommandService
	SourceGit SourceGitService
	Edit      EditService
	Doctor    DoctorService
}

// NewServices creates the default application service graph backed by chezmoi.
func NewServices(client chezmoi.Client) Services {
	s := service{client: client}
	return Services{
		Status:    s,
		Diff:      s,
		Sync:      s,
		Workspace: s,
		Target:    s,
		SourceGit: s,
		Edit:      s,
		Doctor:    s,
	}
}

type service struct {
	client chezmoi.Client
}

func (s service) Status(targets []string) (SyncStatus, error) {
	entries, err := s.client.Status(targets)
	if err != nil {
		return SyncStatus{}, err
	}
	return partitionSyncStatus(entries), nil
}

func (s service) StatusReport(targets []string) (report.Document, error) {
	status, err := s.Status(targets)
	if err != nil {
		return report.Document{}, err
	}
	sourceEntries, err := s.SourceStatus()
	if err != nil {
		return report.Document{}, err
	}

	var doc report.Document
	if len(status.Entries) == 0 && len(status.Scripts) == 0 && len(sourceEntries) == 0 {
		doc.Blocks = append(doc.Blocks, report.Paragraph(report.Text("clean")))
		return doc, nil
	}
	if len(status.Entries) > 0 {
		doc.Blocks = append(doc.Blocks,
			report.Heading(1, report.Strong("local:")),
			report.Paragraph(report.Text("  local config differs from chezmoi source")),
		)
		for _, entry := range status.Entries {
			doc.Blocks = append(doc.Blocks, report.Paragraph(
				report.Warning("!"),
				report.Text(" "),
				report.Path(entry.Path),
				report.Text("  differs from chezmoi"),
			))
		}
		doc.Blocks = append(doc.Blocks, report.Paragraph(report.Muted("  run cm sync")))
	}
	if len(status.Scripts) > 0 {
		if len(doc.Blocks) > 0 {
			doc.Blocks = append(doc.Blocks, report.Blank())
		}
		doc.Blocks = append(doc.Blocks,
			report.Heading(1, report.Strong("automation:")),
			report.Paragraph(report.Text("  chezmoi scripts are pending; cm sync does not execute scripts")),
		)
		for _, entry := range status.Scripts {
			doc.Blocks = append(doc.Blocks, report.Paragraph(
				report.Warning("R"),
				report.Text(" "),
				report.Path(entry.Path),
				report.Text("  apply would run this script"),
			))
		}
		doc.Blocks = append(doc.Blocks, report.Paragraph(report.Muted("  use chezmoi diff and chezmoi apply to review and run scripts")))
	}
	if len(sourceEntries) > 0 {
		if len(doc.Blocks) > 0 {
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

type doctorCheck struct {
	Name     string
	Status   string
	Required bool
	Message  string
}

func (s service) DoctorReport() (report.Document, bool) {
	checks := []doctorCheck{
		s.checkCommand("chezmoi", s.client.BinaryName(), []string{"--version"}, true),
		s.checkCommand("git", "git", []string{"--version"}, true),
	}

	sourceDir, sourceOK := s.checkSourceDir()
	checks = append(checks, sourceOK)
	if sourceOK.Status == "pass" && sourceDir != "" {
		checks = append(checks, s.checkGitRepository(sourceDir))
	}
	checks = append(checks, s.checkCommand("lazygit", "lazygit", []string{"--version"}, false))

	failed := false
	doc := report.Document{Blocks: []report.Block{report.Heading(1, report.Strong("doctor:"))}}
	for _, check := range checks {
		if check.Required && check.Status == "fail" {
			failed = true
		}
		doc.Blocks = append(doc.Blocks, report.Paragraph(
			doctorStatus(check.Status),
			report.Text(" "),
			report.Strong(check.Name+":"),
			report.Text(" "),
			report.Text(check.Message),
		))
	}
	return doc, failed
}

func (s service) checkCommand(name, command string, args []string, required bool) doctorCheck {
	_, err := s.runner().Output(command, args, process.IO{})
	if err != nil {
		status := "warn"
		if required {
			status = "fail"
		}
		return doctorCheck{Name: name, Status: status, Required: required, Message: err.Error()}
	}
	return doctorCheck{Name: name, Status: "pass", Required: required, Message: "available"}
}

func (s service) checkSourceDir() (string, doctorCheck) {
	sourceDir, err := s.sourceDir()
	if err != nil {
		return "", doctorCheck{Name: "source-path", Status: "fail", Required: true, Message: err.Error()}
	}
	if sourceDir == "" {
		return "", doctorCheck{Name: "source-path", Status: "fail", Required: true, Message: "empty chezmoi source path"}
	}
	return sourceDir, doctorCheck{Name: "source-path", Status: "pass", Required: true, Message: sourceDir}
}

func (s service) checkGitRepository(sourceDir string) doctorCheck {
	var stderr bytes.Buffer
	_, err := s.runner().Output("git", []string{"status", "--porcelain=v1"}, process.IO{Dir: sourceDir, Stderr: &stderr})
	if err != nil {
		return doctorCheck{Name: "source-git", Status: "fail", Required: true, Message: fmt.Sprintf("git status %s: %v%s", sourceDir, err, formatStderr(stderr.Bytes()))}
	}
	return doctorCheck{Name: "source-git", Status: "pass", Required: true, Message: "repository readable"}
}

func doctorStatus(status string) report.Inline {
	switch status {
	case "pass":
		return report.Status("pass")
	case "warn":
		return report.Warning("warn")
	default:
		return report.Warning("fail")
	}
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
		status, err := s.Status(nil)
		if err != nil {
			return report.Document{}, err
		}
		if len(status.Entries) == 0 {
			if len(status.Scripts) > 0 {
				return report.Document{Blocks: []report.Block{
					report.Paragraph(report.Text("no file diffs; chezmoi scripts are pending")),
					report.Paragraph(report.Muted("run chezmoi diff to review script contents")),
				}}, nil
			}
			return report.Document{Blocks: []report.Block{report.Paragraph(report.Text("clean"))}}, nil
		}
		targets = make([]string, 0, len(status.Entries))
		for _, entry := range status.Entries {
			targets = append(targets, entry.Path)
		}
	}

	doc := report.Document{Blocks: make([]report.Block, 0, len(targets))}
	for _, target := range targets {
		diffOutput, err := s.DiffOutput(target)
		if err != nil {
			return report.Document{}, err
		}
		doc.Blocks = append(doc.Blocks, report.DiffBlock(string(diffOutput)))
	}
	return doc, nil
}

func (s service) DiffOutput(target string) ([]byte, error) {
	out, err := s.authoritativeDiff(target)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return []byte(fmt.Sprintf("no diff: %s\n", target)), nil
	}
	return out, nil
}

func (s service) authoritativeDiff(target string) ([]byte, error) {
	out, err := s.client.AuthoritativeDiff(target, maxReviewOutput)
	if err != nil {
		return nil, err
	}
	if int64(len(out)) > maxReviewOutput {
		return nil, fmt.Errorf("diff output exceeds %d bytes for %s", maxReviewOutput, target)
	}
	return out, nil
}

func (s service) Review(target string) (Review, error) {
	status, err := s.Status([]string{target})
	if err != nil {
		return Review{}, err
	}
	if entry, ok := findReconcileEntry(status.Scripts, target); ok {
		return newReview(entry, TargetScript, false, ""), nil
	}
	entry, ok := findReconcileEntry(status.Entries, target)
	if !ok {
		if len(status.Entries) > 0 || len(status.Scripts) > 0 {
			return Review{}, fmt.Errorf("chezmoi status did not return requested target %s", target)
		}
		return cleanReview(target), nil
	}
	diffOutput, err := s.authoritativeDiff(entry.Path)
	if err != nil {
		return Review{}, err
	}
	if len(entry.Code) >= 2 && entry.Code[1] == 'D' {
		return newReview(entry, TargetRemove, false, string(diffOutput)), nil
	}
	metadata, err := s.client.TargetMetadata(entry.Path, maxReviewOutput)
	if err != nil {
		return Review{}, err
	}
	return newReview(entry, targetTypeFromChezmoi(metadata.Type), metadata.Template, string(diffOutput)), nil
}

func (s service) ExecuteNonInteractive(action Action) (ActionResult, error) {
	switch action.Kind {
	case ActionAdd:
		return s.runBuffered("re-add", action.Target)
	case ActionApply:
		return s.runBuffered("apply", "--force", action.Target)
	case ActionMerge:
		return ActionResult{}, fmt.Errorf("merge requires terminal execution")
	default:
		return ActionResult{}, fmt.Errorf("unknown reconcile action %d for %s", action.Kind, action.Target)
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
	if filepath.IsAbs(target) {
		return s.client.Run("edit", target)
	}
	targetDir, err := s.client.TargetDir()
	if err != nil {
		return err
	}
	return s.client.Run("edit", filepath.Join(targetDir, target))
}

func (s service) ManagedFiles() ([]string, error) {
	return s.client.ManagedFiles()
}

func (s service) runner() process.Runner {
	return s.client.ActiveRunner()
}

func (s service) runBuffered(args ...string) (ActionResult, error) {
	stdout, stderr, err := s.client.RunBuffered(args...)
	result := ActionResult{Stdout: string(stdout), Stderr: string(stderr)}
	if err != nil {
		return result, fmt.Errorf("%w%s", err, formatStderr(stderr))
	}
	return result, nil
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
