package chezmoi

import (
	"reflect"
	"testing"
)

func TestParseStatusParsesTwoColumnEntries(t *testing.T) {
	entries, err := ParseStatus([]byte("M  /home/me/.zshrc\n M /home/me/.gitconfig\nMM /home/me/.config/nvim/init.lua\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{LocalChange: ChangeModified, TargetChange: ChangeNone, Path: "/home/me/.zshrc"},
		{LocalChange: ChangeNone, TargetChange: ChangeModified, Path: "/home/me/.gitconfig"},
		{LocalChange: ChangeModified, TargetChange: ChangeModified, Path: "/home/me/.config/nvim/init.lua"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestParseStatusParsesAllKnownStatusCodes(t *testing.T) {
	entries, err := ParseStatus([]byte("A  /home/me/added\n D /home/me/deleted-by-apply\nD  /home/me/deleted-local\n R /home/me/run-script\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{LocalChange: ChangeAdded, TargetChange: ChangeNone, Path: "/home/me/added"},
		{LocalChange: ChangeNone, TargetChange: ChangeDeleted, Path: "/home/me/deleted-by-apply"},
		{LocalChange: ChangeDeleted, TargetChange: ChangeNone, Path: "/home/me/deleted-local"},
		{LocalChange: ChangeNone, TargetChange: ChangeRun, Path: "/home/me/run-script"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestParseStatusRejectsMalformedLine(t *testing.T) {
	_, err := ParseStatus([]byte("M\n"))
	if err == nil {
		t.Fatal("ParseStatus returned nil error for malformed line")
	}
}
