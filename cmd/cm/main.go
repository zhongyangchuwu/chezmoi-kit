package main

import (
	"fmt"
	"io"
	"os"

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
	command := "status"
	commandArgs := args
	if len(args) > 0 {
		command = args[0]
		commandArgs = args[1:]
	}

	var err error
	switch command {
	case "status":
		err = renderStatus(stdout, svc, commandArgs)
	case "diff":
		err = svc.Diff(commandArgs)
	case "add":
		err = svc.AddTargets(commandArgs)
	case "apply":
		err = svc.ApplyTargets(commandArgs)
	case "merge":
		err = svc.MergeTargets(commandArgs)
	case "sync":
		err = ui.RunSync(svc, commandArgs, stdin, stdout)
	case "help", "-h", "--help":
		renderUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", command)
		renderUsage(stderr)
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
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

func renderUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: cm [status|diff|sync|add|apply|merge] [target...]")
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
