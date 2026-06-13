package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhongyangchuwu/cm/internal/build"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/reconcile"
	"github.com/zhongyangchuwu/cm/internal/ui"
)

type service interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	Diff(targets []string) error
	Add(target string) error
	Apply(target string) error
	Merge(target string) error
	AddTargets(targets []string) error
	ApplyTargets(targets []string) error
	MergeTargets(targets []string) error
}

func main() {
	client := chezmoi.Client{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	svc := chezmoiService{client: client}
	os.Exit(run(os.Args[1:], svc, os.Stdin, os.Stdout, os.Stderr))
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
	root.AddCommand(&cobra.Command{
		Use:   "sync [target...]",
		Short: "Interactively reconcile chezmoi changes",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ui.RunSync(svc, args, stdin, stdout)
		},
	})
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

func renderStatus(w io.Writer, svc service, targets []string) error {
	entries, err := svc.Status(targets)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		_, err = fmt.Fprintln(w, "clean")
		return err
	}
	for _, entry := range entries {
		status := string([]byte{byte(entry.LocalChange), byte(entry.TargetChange)})
		_, err := fmt.Fprintf(w, "%s %s  %s; run cm sync %s\n", status, entry.Path, reconcile.Describe(entry), entry.Path)
		if err != nil {
			return err
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
