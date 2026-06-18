package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/process"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
	"github.com/zhongyangchuwu/cm/internal/syncdiff"
)

type statusService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	SourceStatus() ([]sourceEntry, error)
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
	Sync      reconcile.ReviewService
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

type sourceEntry struct {
	Code string
	Path string
}

type chezmoiService struct {
	client chezmoi.Client
}

func (s chezmoiService) Status(targets []string) ([]chezmoi.StatusEntry, error) {
	return s.client.Status(targets)
}

func (s chezmoiService) SourceStatus() ([]sourceEntry, error) {
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
	return parseSourceStatus(out), nil
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
	differ := syncdiff.Differ{Source: chezmoi.ContentLoader{Client: s.client}}
	return differ.Diff(target)
}

func (s chezmoiService) Execute(actions []reconcile.Action) error {
	var addTargets []string
	var applyTargets []string
	var mergeTargets []string

	for _, action := range actions {
		switch action.Kind {
		case reconcile.ActionAdd:
			addTargets = append(addTargets, action.Target)
		case reconcile.ActionApply:
			applyTargets = append(applyTargets, action.Target)
		case reconcile.ActionMerge:
			mergeTargets = append(mergeTargets, action.Target)
		default:
			return fmt.Errorf("unknown reconcile action %d for %s", action.Kind, action.Target)
		}
	}

	if len(addTargets) > 0 {
		if err := s.runTargets("add", addTargets); err != nil {
			return fmt.Errorf("add %s: %w", targetSummary(addTargets), err)
		}
	}
	if len(applyTargets) > 0 {
		args := make([]string, 0, 2+len(applyTargets))
		args = append(args, "apply", "--force")
		args = append(args, applyTargets...)
		if err := s.client.Run(args...); err != nil {
			return fmt.Errorf("apply %s: %w", targetSummary(applyTargets), err)
		}
	}
	for _, target := range mergeTargets {
		if err := s.runTarget("merge", target); err != nil {
			return fmt.Errorf("merge %s: %w", target, err)
		}
	}
	return nil
}

func targetSummary(targets []string) string {
	if len(targets) == 1 {
		return targets[0]
	}
	return fmt.Sprintf("%d targets", len(targets))
}

func (s chezmoiService) runner() process.Runner {
	return s.client.ActiveRunner()
}

func parseSourceStatus(out []byte) []sourceEntry {
	lines := bytes.Split(bytes.TrimRight(out, "\n"), []byte{'\n'})
	if len(lines) == 1 && len(lines[0]) == 0 {
		return nil
	}
	entries := make([]sourceEntry, 0, len(lines))
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		entries = append(entries, sourceEntry{Code: string(line[:2]), Path: string(line[3:])})
	}
	return entries
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

func (s chezmoiService) runTarget(command, target string) error {
	return s.runTargets(command, []string{target})
}

func (s chezmoiService) runTargets(command string, targets []string) error {
	args := make([]string, 0, 1+len(targets))
	args = append(args, command)
	args = append(args, targets...)
	return s.client.Run(args...)
}
