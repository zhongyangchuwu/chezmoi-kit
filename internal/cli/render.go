package cli

import (
	"io"
	"os"

	"github.com/zhongyangchuwu/cm/internal/report"
	"golang.org/x/term"
)

func renderReport(w io.Writer, doc report.Document) error {
	_, err := w.Write(report.ANSI(doc, report.Options{Color: report.ColorAuto, IsTTY: isTerminalWriter(w)}))
	return err
}

func isTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}
