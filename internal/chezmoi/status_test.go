package chezmoi

import (
	"reflect"
	"testing"
)

func TestParseStatusParsesTwoColumnEntries(t *testing.T) {
	entries, err := ParseStatus([]byte("M  .zshrc\n M .gitconfig\nMM .config/nvim/init.lua\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{LocalChange: ChangeModified, TargetChange: ChangeNone, Path: ".zshrc"},
		{LocalChange: ChangeNone, TargetChange: ChangeModified, Path: ".gitconfig"},
		{LocalChange: ChangeModified, TargetChange: ChangeModified, Path: ".config/nvim/init.lua"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestParseStatusParsesAllKnownStatusCodes(t *testing.T) {
	entries, err := ParseStatus([]byte("A  added\n D deleted-by-source\nD  deleted-local\n R run-script\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{LocalChange: ChangeAdded, TargetChange: ChangeNone, Path: "added"},
		{LocalChange: ChangeNone, TargetChange: ChangeDeleted, Path: "deleted-by-source"},
		{LocalChange: ChangeDeleted, TargetChange: ChangeNone, Path: "deleted-local"},
		{LocalChange: ChangeNone, TargetChange: ChangeRun, Path: "run-script"},
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
