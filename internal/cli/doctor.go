package cli

import (
	"fmt"
	"io"

	"github.com/zhongyangchuwu/cm/internal/app"
)

func renderDoctor(w io.Writer, svc app.DoctorService, opts renderOptions) error {
	doc, failed := svc.DoctorReport()
	if err := renderReport(w, doc, opts); err != nil {
		return err
	}
	if failed {
		return fmt.Errorf("doctor found failed checks")
	}
	return nil
}
