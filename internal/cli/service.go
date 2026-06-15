package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/syncdiff"
)

type service interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	SourceStatus() ([]sourceEntry, error)
	Diff(targets []string) error
	Add(target string) error
	Apply(target string) error
	Merge(target string) error
	AddTargets(targets []string) error
	ApplyTargets(targets []string) error
	MergeTargets(targets []string) error
	OpenSourceGit() error
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

	cmd := exec.Command("git", "-C", sourceDir, "status", "--porcelain=v1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
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

	cmd := exec.Command("lazygit")
	cmd.Dir = sourceDir
	cmd.Stdin = s.client.Stdin
	cmd.Stdout = s.client.Stdout
	cmd.Stderr = s.client.Stderr
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("lazygit not found in PATH")
		}
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

func (s chezmoiService) DiffOutput(targets []string) ([]byte, error) {
	if len(targets) != 1 {
		return nil, fmt.Errorf("diff output requires exactly one target, got %d", len(targets))
	}
	differ := syncdiff.Differ{Source: syncdiff.ChezmoiContentLoader{Client: s.client}}
	return differ.Diff(targets[0])
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

func (s chezmoiService) Diff(targets []string) error {
	return s.runTargets("diff", targets)
}

func (s chezmoiService) Add(target string) error {
	return s.runTarget("add", target)
}

func (s chezmoiService) Apply(target string) error {
	return s.runTarget("apply", target)
}

func (s chezmoiService) Merge(target string) error {
	return s.runTarget("merge", target)
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

func (s chezmoiService) runTarget(command, target string) error {
	return s.runTargets(command, []string{target})
}

func (s chezmoiService) runTargets(command string, targets []string) error {
	args := make([]string, 0, 1+len(targets))
	args = append(args, command)
	args = append(args, targets...)
	return s.client.Run(args...)
}
