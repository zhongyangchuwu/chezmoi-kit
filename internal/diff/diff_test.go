package diff

import (
	"os"
	"strings"
	"testing"
)

func TestDiffBytesShowsChezmoiToLocalDirection(t *testing.T) {
	base := []byte("export EDITOR=nvim\nalias ll='ls -la'\n")
	local := []byte("export EDITOR=vim\nalias ll='ls -lah'\n")

	got := string(DiffBytes("/home/me/.zshrc", base, local, 0))
	for _, want := range []string{
		"--- chezmoi:/home/me/.zshrc",
		"+++ local:/home/me/.zshrc",
		"-export EDITOR=nvim",
		"+export EDITOR=vim",
		"-alias ll='ls -la'",
		"+alias ll='ls -lah'",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("diff %q does not contain %q", got, want)
		}
	}
}

func TestDiffBytesReportsIdenticalContent(t *testing.T) {
	got := string(DiffBytes("/home/me/.zshrc", []byte("same\n"), []byte("same\n"), 0))
	if got != "no diff: /home/me/.zshrc\n" {
		t.Fatalf("diff = %q", got)
	}
}

func TestDiffBytesKeepsMissingNewlineWarning(t *testing.T) {
	got := string(DiffBytes("/home/me/.zshrc", []byte("same\n"), []byte("same"), 0))
	if !strings.Contains(got, "\\ No newline at end of file") {
		t.Fatalf("diff %q does not contain missing newline warning", got)
	}
}

func TestDiffBytesSkipsBinaryContent(t *testing.T) {
	got := string(DiffBytes("/home/me/.ssh/key", []byte{'a', 0, 'b'}, []byte("plain"), 0))
	if got != "binary file differs: /home/me/.ssh/key\n" {
		t.Fatalf("diff = %q", got)
	}
}

func TestDiffBytesSkipsLargeContent(t *testing.T) {
	got := string(DiffBytes("/home/me/big", []byte("small"), []byte("large"), 3))
	if got != "file too large to diff: /home/me/big\n" {
		t.Fatalf("diff = %q", got)
	}
}

func TestDifferReadsSourceContent(t *testing.T) {
	source := fakeContentSource{
		base:  []byte("source\n"),
		local: []byte("local\n"),
	}
	differ := Differ{Source: &source}

	got, err := differ.Diff("/home/me/.zshrc")
	if err != nil {
		t.Fatalf("Diff returned error: %v", err)
	}
	if !strings.Contains(string(got), "-source") || !strings.Contains(string(got), "+local") {
		t.Fatalf("diff = %q", string(got))
	}
}
func TestDifferTreatsMissingLocalAsEmpty(t *testing.T) {
	source := fakeContentSource{
		base:     []byte("source\n"),
		localErr: os.ErrNotExist,
	}
	differ := Differ{Source: &source}

	got, err := differ.Diff("/home/me/.zshrc")
	if err != nil {
		t.Fatalf("Diff returned error: %v", err)
	}
	if !strings.Contains(string(got), "-source") {
		t.Fatalf("diff = %q", string(got))
	}
}

func TestDifferSkipsLocalReadWhenSourceIsTooLarge(t *testing.T) {
	source := fakeContentSource{
		base: []byte("too large"),
	}
	differ := Differ{Source: &source, MaxFileSize: 3}

	got, err := differ.Diff("/home/me/big")
	if err != nil {
		t.Fatalf("Diff returned error: %v", err)
	}
	if string(got) != "file too large to diff: /home/me/big\n" {
		t.Fatalf("diff = %q", string(got))
	}
	if source.localCalls != 0 {
		t.Fatalf("localCalls = %d, want 0", source.localCalls)
	}
	if source.targetLimit != 3 {
		t.Fatalf("targetLimit = %d, want 3", source.targetLimit)
	}
	if source.localLimit != 0 {
		t.Fatalf("localLimit = %d, want 0", source.localLimit)
	}
}

type fakeContentSource struct {
	base        []byte
	local       []byte
	localErr    error
	localCalls  int
	targetLimit int64
	localLimit  int64
}

func (f *fakeContentSource) TargetContent(_ string, limit int64) ([]byte, error) {
	f.targetLimit = limit
	return append([]byte(nil), f.base...), nil
}

func (f *fakeContentSource) LocalContent(_ string, limit int64) ([]byte, error) {
	f.localCalls++
	f.localLimit = limit
	if f.localErr != nil {
		return nil, f.localErr
	}
	return append([]byte(nil), f.local...), nil
}
