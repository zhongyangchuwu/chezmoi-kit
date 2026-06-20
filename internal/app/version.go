package app

import (
	"fmt"
	"runtime/debug"

	"github.com/zhongyangchuwu/cm/internal/report"
)

var Version = "dev"

type Info struct {
	Version   string
	Commit    string
	Time      string
	Modified  string
	GoVersion string
}

func Current() Info {
	info := Info{Version: Version}
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	info.GoVersion = buildInfo.GoVersion
	if Version == "dev" && buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		info.Version = buildInfo.Main.Version
	}
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.Commit = setting.Value
		case "vcs.time":
			info.Time = setting.Value
		case "vcs.modified":
			info.Modified = setting.Value
		}
	}
	return info
}

func (i Info) FormatShort(name string) string {
	return fmt.Sprintf("%s %s", name, i.Version)
}

func (i Info) Report(name string) report.Document {
	return report.Document{Blocks: []report.Block{
		report.Paragraph(report.Strong(name+":"), report.Text(" "), report.Status(i.Version)),
		report.Paragraph(report.Muted("commit:"), report.Text(" "), report.Code(fallback(i.Commit, "unknown"))),
		report.Paragraph(report.Muted("built:"), report.Text(" "), report.Text(fallback(i.Time, "unknown"))),
		report.Paragraph(report.Muted("dirty:"), report.Text(" "), report.Text(fallback(i.Modified, "unknown"))),
		report.Paragraph(report.Muted("go:"), report.Text(" "), report.Text(fallback(i.GoVersion, "unknown"))),
	}}
}

func fallback(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
