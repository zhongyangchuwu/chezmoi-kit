package build

import "testing"

func TestCurrentReturnsRuntimeBuildInformation(t *testing.T) {
	info := Current()

	if info.Version == "" {
		t.Fatal("Version is empty")
	}
	if info.GoVersion == "" {
		t.Fatal("GoVersion is empty")
	}
}

func TestExplicitVersionOverridesRuntimeModuleVersion(t *testing.T) {
	previous := Version
	Version = "v9.9.9"
	t.Cleanup(func() { Version = previous })

	info := Current()

	if info.Version != "v9.9.9" {
		t.Fatalf("Version = %q, want explicit ldflag version", info.Version)
	}
}

func TestFormatDetailedIncludesCoreFields(t *testing.T) {
	info := Info{
		Version:   "v0.1.0",
		Commit:    "abc123",
		Time:      "2026-06-13T00:00:00Z",
		Modified:  "false",
		GoVersion: "go1.26.4",
	}

	got := info.FormatDetailed("cm")
	for _, want := range []string{"cm: v0.1.0", "commit: abc123", "built: 2026-06-13T00:00:00Z", "dirty: false", "go: go1.26.4"} {
		if !containsLine(got, want) {
			t.Fatalf("FormatDetailed() = %q, missing line %q", got, want)
		}
	}
}

func containsLine(s, line string) bool {
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			if s[start:i] == line {
				return true
			}
			start = i + 1
		}
	}
	return false
}
