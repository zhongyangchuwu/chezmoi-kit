package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zhongyangchuwu/cm/internal/build"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/syncdiff"
	"github.com/zhongyangchuwu/cm/internal/ui"
	"golang.org/x/term"
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

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	client := chezmoi.Client{Stdin: stdin, Stdout: stdout, Stderr: stderr}
	svc := chezmoiService{client: client}
	return run(args, svc, stdin, stdout, stderr)
}

func run(args []string, svc service, stdin io.Reader, stdout, stderr io.Writer) int {
	cmd := newRootCommand(svc, stdin, stdout, stderr)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func newRootCommand(svc service, stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	info := build.Current()
	root := &cobra.Command{
		Use:           "cm",
		Short:         "Chezmoi reconciliation manager",
		Version:       info.FormatShort("cm"),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderStatus(stdout, svc, args)
		},
	}
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)

	root.AddCommand(&cobra.Command{
		Use:   "status [target...]",
		Short: "Show chezmoi reconciliation status",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderStatus(stdout, svc, args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "diff [target...]",
		Short: "Show chezmoi diff",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Diff(args)
		},
	})
	syncCmd := &cobra.Command{
		Use:   "sync [target...]",
		Short: "Interactively reconcile chezmoi changes",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			useTUI, err := cmd.Flags().GetBool("tui")
			if err != nil {
				return err
			}
			plain, err := cmd.Flags().GetBool("plain")
			if err != nil {
				return err
			}
			if useTUI && plain {
				return fmt.Errorf("--tui and --plain are mutually exclusive")
			}
			if useTUI || (!plain && isTerminal(stdin) && isTerminal(stdout)) {
				return ui.RunSyncTUI(svc, args, stdin, stdout)
			}
			return ui.RunSync(svc, args, stdin, stdout)
		},
	}
	syncCmd.Flags().Bool("tui", false, "run sync in terminal UI mode")
	syncCmd.Flags().Bool("plain", false, "run sync in plain prompt mode")
	root.AddCommand(syncCmd)
	root.AddCommand(&cobra.Command{
		Use:   "add [target...]",
		Short: "Accept local files into chezmoi source state",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.AddTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "apply [target...]",
		Short: "Apply chezmoi target state locally",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.ApplyTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "merge [target...]",
		Short: "Open chezmoi merge for targets",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.MergeTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "git",
		Short: "Open lazygit in the chezmoi source repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.OpenSourceGit()
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprint(stdout, info.FormatDetailed("cm"))
			return err
		},
	})
	root.AddCommand(newCompletionCommand(root, stdout))
	return root
}

func newCompletionCommand(root *cobra.Command, stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(stdout)
			case "zsh":
				return root.GenZshCompletion(stdout)
			case "fish":
				return root.GenFishCompletion(stdout, true)
			case "powershell":
				return root.GenPowerShellCompletion(stdout)
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}
}

func isTerminal(v any) bool {
	file, ok := v.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}

func renderStatus(w io.Writer, svc service, targets []string) error {
	entries, err := svc.Status(targets)
	if err != nil {
		return err
	}
	sourceEntries, err := svc.SourceStatus()
	if err != nil {
		return err
	}
	if len(entries) == 0 && len(sourceEntries) == 0 {
		_, err = fmt.Fprintln(w, "clean")
		return err
	}
	if len(entries) > 0 {
		if _, err := color.New(color.FgCyan, color.Bold).Fprintln(w, "local:"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "  local config differs from chezmoi source"); err != nil {
			return err
		}
		for _, entry := range entries {
			if _, err := color.New(color.FgYellow, color.Bold).Fprint(w, "!"); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, " %s  differs from chezmoi\n", entry.Path); err != nil {
				return err
			}
		}
		if _, err := color.New(color.FgHiBlack).Fprintln(w, "  run cm sync"); err != nil {
			return err
		}
	}
	if len(sourceEntries) > 0 {
		if len(entries) > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := color.New(color.FgCyan, color.Bold).Fprintln(w, "chezmoi:"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "  source repository has git changes"); err != nil {
			return err
		}
		for _, entry := range sourceEntries {
			if _, err := fmt.Fprintf(w, "%s %s\n", entry.Code, entry.Path); err != nil {
				return err
			}
		}
	}
	return nil
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
	differ := syncdiff.Differ{Source: syncdiff.ChezmoiSource{Target: s.client}}
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
	return s.client.Run(append([]string{"diff"}, targets...)...)
}

func (s chezmoiService) Add(target string) error {
	return s.client.Run("add", target)
}

func (s chezmoiService) Apply(target string) error {
	return s.client.Run("apply", target)
}

func (s chezmoiService) Merge(target string) error {
	return s.client.Run("merge", target)
}

func (s chezmoiService) AddTargets(targets []string) error {
	return s.client.Run(append([]string{"add"}, targets...)...)
}

func (s chezmoiService) ApplyTargets(targets []string) error {
	return s.client.Run(append([]string{"apply"}, targets...)...)
}

func (s chezmoiService) MergeTargets(targets []string) error {
	return s.client.Run(append([]string{"merge"}, targets...)...)
}
