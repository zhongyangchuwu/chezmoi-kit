package build

import (
	"fmt"
	"runtime/debug"
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
	if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
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

func (i Info) FormatDetailed(name string) string {
	return fmt.Sprintf("%s: %s\ncommit: %s\nbuilt: %s\ndirty: %s\ngo: %s\n", name, i.Version, fallback(i.Commit, "unknown"), fallback(i.Time, "unknown"), fallback(i.Modified, "unknown"), fallback(i.GoVersion, "unknown"))
}

func fallback(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
