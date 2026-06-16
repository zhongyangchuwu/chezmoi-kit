package cli

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

func renderStatus(w io.Writer, svc statusService, targets []string) error {
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
