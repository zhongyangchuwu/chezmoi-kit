package cli

import (
	"fmt"
	"io"
)

func renderDiff(w io.Writer, svc service, targets []string) error {
	if len(targets) == 0 {
		entries, err := svc.Status(nil)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			_, err = fmt.Fprintln(w, "clean")
			return err
		}
		targets = make([]string, 0, len(entries))
		for _, entry := range entries {
			targets = append(targets, entry.Path)
		}
	}

	for _, target := range targets {
		diff, err := svc.DiffOutput(target)
		if err != nil {
			return err
		}
		if _, err := w.Write(diff); err != nil {
			return err
		}
	}
	return nil
}
