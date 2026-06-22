package cli

import (
	"io"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func renderDiff(w io.Writer, svc app.DiffService, targets []string, opts renderOptions) error {
	doc, err := svc.DiffReport(targets)
	if err != nil {
		return err
	}
	return renderReport(w, doc, opts)
}
