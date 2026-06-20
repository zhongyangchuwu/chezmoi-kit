package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/zhongyangchuwu/cm/internal/app"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
	"github.com/zhongyangchuwu/cm/internal/tui"
)

func Main(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	client := chezmoi.Client{Stdin: stdin, Stdout: stdout, Stderr: stderr}
	return run(args, app.NewServices(client), stdin, stdout, stderr)
}

func run(args []string, services commandServices, stdin io.Reader, stdout, stderr io.Writer) int {
	options := &app.Options{Stderr: stderr}
	cmd := newRootCommand(services, stdin, stdout, stderr, options)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func newRootCommand(services commandServices, stdin io.Reader, stdout, stderr io.Writer, options *app.Options) *cobra.Command {
	info := app.Current()
	root := &cobra.Command{
		Use:           "cm",
		Short:         "Chezmoi reconciliation manager",
		Version:       info.FormatShort("cm"),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderStatus(stdout, services.Status, args)
		},
	}
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().BoolVar(&options.Debug, "debug", false, "write debug diagnostics to a temporary log file")

	root.AddCommand(&cobra.Command{
		Use:   "status [target...]",
		Short: "Show chezmoi reconciliation status",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderStatus(stdout, services.Status, args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "diff [target...]",
		Short: "Show internal sync diff",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderDiff(stdout, services.Diff, args)
		},
	})
	syncCmd := &cobra.Command{
		Use:   "sync [target...]",
		Short: "Interactively reconcile chezmoi changes",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.RunSyncTUI(services.Sync, args, stdin, stdout, options)
		},
	}
	root.AddCommand(syncCmd)
	root.AddCommand(&cobra.Command{
		Use:   "add [target...]",
		Short: "Accept local files into chezmoi source state",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return services.Target.AddTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "apply [target...]",
		Short: "Apply chezmoi target state locally",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return services.Target.ApplyTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "merge [target...]",
		Short: "Open chezmoi merge for targets",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return services.Target.MergeTargets(args)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "edit <target>",
		Short: "Edit a chezmoi managed file in your configured editor",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			files, err := services.Edit.ManagedFiles()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return files, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return services.Edit.EditTarget(args[0])
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "git",
		Short: "Open lazygit in the chezmoi source repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return services.SourceGit.OpenSourceGit()
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print build version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderReport(stdout, info.Report("cm"))
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
