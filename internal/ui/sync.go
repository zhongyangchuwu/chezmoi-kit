package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

type SyncService interface {
	Status(targets []string) ([]chezmoi.StatusEntry, error)
	Diff(targets []string) error
	Add(target string) error
	Apply(target string) error
	Merge(target string) error
}

func RunSync(service SyncService, targets []string, input io.Reader, output io.Writer) error {
	entries, err := service.Status(targets)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		_, err = fmt.Fprintln(output, "clean")
		return err
	}

	reader := bufio.NewReader(input)
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		for {
			if err := renderEntry(output, entry); err != nil {
				return err
			}

			choice, err := readChoice(reader)
			if err != nil {
				return err
			}

			switch choice {
			case "d":
				if err := service.Diff([]string{entry.Path}); err != nil {
					return err
				}
				continue
			case "a":
				if err := service.Add(entry.Path); err != nil {
					return err
				}
			case "p":
				if err := service.Apply(entry.Path); err != nil {
					return err
				}
			case "m":
				if err := service.Merge(entry.Path); err != nil {
					return err
				}
			case "s":
				break
			case "q":
				return nil
			default:
				if _, err := fmt.Fprintln(output, "unknown choice"); err != nil {
					return err
				}
				continue
			}

			if choice == "s" {
				break
			}

			fresh, err := service.Status([]string{entry.Path})
			if err != nil {
				return err
			}
			if len(fresh) == 0 {
				break
			}
			entry = fresh[0]
		}
	}
	return nil
}

func renderEntry(w io.Writer, entry chezmoi.StatusEntry) error {
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	if _, err := color.New(color.FgYellow, color.Bold).Fprint(w, "!"); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, " %s\nlocal differs from chezmoi\n[d]iff [a]dd local a[p]ply chezmoi [m]erge [s]kip [q]uit\n> ", entry.Path)
	return err
}

func readChoice(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.ToLower(strings.TrimSpace(line)), nil
}
