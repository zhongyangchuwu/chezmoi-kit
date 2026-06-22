package cli

import (
	"io"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func renderStatus(w io.Writer, svc app.StatusService, targets []string, opts renderOptions) error {
	doc, err := svc.StatusReport(targets)
	if err != nil {
		return err
	}
	return renderReport(w, doc, opts)
}
