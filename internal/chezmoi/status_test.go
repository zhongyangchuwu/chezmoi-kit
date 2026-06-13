package chezmoi

import (
	"reflect"
	"testing"
)

func TestParseStatusKeepsRawCodeAndAbsolutePath(t *testing.T) {
	entries, err := ParseStatus([]byte("M  /home/me/.zshrc\n M /home/me/.gitconfig\nMM /home/me/.config/nvim/init.lua\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{Code: "M ", Path: "/home/me/.zshrc"},
		{Code: " M", Path: "/home/me/.gitconfig"},
		{Code: "MM", Path: "/home/me/.config/nvim/init.lua"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("entries = %#v, want %#v", entries, want)
	}
}

func TestParseStatusParsesAllKnownStatusCodesAsOutOfSyncEntries(t *testing.T) {
	entries, err := ParseStatus([]byte("A  /home/me/added\n D /home/me/deleted-by-apply\nD  /home/me/deleted-local\n R /home/me/run-script\n"))
	if err != nil {
		t.Fatalf("ParseStatus returned error: %v", err)
	}

	want := []StatusEntry{
		{Code: "A ", Path: "/home/me/added"},
		{Code: " D", Path: "/home/me/deleted-by-apply"},
		{Code: "D ", Path: "/home/me/deleted-local"},
		{Code: " R", Path: "/home/me/run-script"},
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
