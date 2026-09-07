package chezmoi

import (
	"reflect"
	"testing"
)

func TestParseManagedInventoryKeepsPathMappings(t *testing.T) {
	entries, err := ParseManagedInventory([]byte(`{
		".config/app/config": {
			"absolute": "/home/me/.config/app/config",
			"sourceAbsolute": "/source/private_dot_config/app/config.tmpl",
			"sourceRelative": "private_dot_config/app/config.tmpl"
		}
	}`))
	if err != nil {
		t.Fatalf("ParseManagedInventory() error = %v", err)
	}
	want := []ManagedEntry{{
		Relative:       ".config/app/config",
		Absolute:       "/home/me/.config/app/config",
		SourceAbsolute: "/source/private_dot_config/app/config.tmpl",
		SourceRelative: "private_dot_config/app/config.tmpl",
	}}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("ParseManagedInventory() = %#v, want %#v", entries, want)
	}
}

func TestParseManagedInventoryRejectsMissingAbsolutePath(t *testing.T) {
	_, err := ParseManagedInventory([]byte(`{".zshrc":{"sourceAbsolute":"/source/dot_zshrc"}}`))
	if err == nil {
		t.Fatal("ParseManagedInventory() returned nil error")
	}
}

func TestClientManagedInventoryUsesAllPathMappings(t *testing.T) {
	runner := &fakeRunner{output: []byte(`{".zshrc":{"absolute":"/home/me/.zshrc","sourceAbsolute":"/source/dot_zshrc","sourceRelative":"dot_zshrc"}}`)}
	client := Client{Binary: "chezmoi", Runner: runner}

	entries, err := client.ManagedInventory([]string{"/home/me/.config"}, 1024)
	if err != nil {
		t.Fatalf("ManagedInventory() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("ManagedInventory() = %#v", entries)
	}
	want := [][]string{{"chezmoi", "managed", "--include=all", "--exclude=none", "--path-style=all", "--format=json", "/home/me/.config"}}
	if !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
}

func TestClientManagedPathsByTypeUsesNULOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("/home/me/.zshrc\x00")}
	client := Client{Binary: "chezmoi", Runner: runner}

	paths, err := client.ManagedPathsByType("templates", nil, 1024)
	if err != nil {
		t.Fatalf("ManagedPathsByType() error = %v", err)
	}
	if !reflect.DeepEqual(paths, []string{"/home/me/.zshrc"}) {
		t.Fatalf("paths = %#v", paths)
	}
	want := [][]string{{"chezmoi", "managed", "--include=templates", "--exclude=none", "--path-style=absolute", "--nul-path-separator"}}
	if !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
}

func TestClientInventoryPathCommandsAreNULDelimited(t *testing.T) {
	runner := &fakeRunner{outputs: [][]byte{[]byte(".ignored\x00"), []byte("/home/me/new\x00")}}
	client := Client{Binary: "chezmoi", Runner: runner}

	ignored, err := client.IgnoredEntries(1024)
	if err != nil {
		t.Fatalf("IgnoredEntries() error = %v", err)
	}
	unmanaged, err := client.UnmanagedEntries([]string{"/home/me/.config"}, 1024)
	if err != nil {
		t.Fatalf("UnmanagedEntries() error = %v", err)
	}
	if !reflect.DeepEqual(ignored, []string{".ignored"}) || !reflect.DeepEqual(unmanaged, []string{"/home/me/new"}) {
		t.Fatalf("ignored=%#v unmanaged=%#v", ignored, unmanaged)
	}
	want := [][]string{
		{"chezmoi", "ignored", "--nul-path-separator"},
		{"chezmoi", "unmanaged", "--include=all", "--exclude=none", "--path-style=absolute", "--nul-path-separator", "/home/me/.config"},
	}
	if !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
}

func TestClientContentCommandsRequireExplicitRevealChoice(t *testing.T) {
	runner := &fakeRunner{outputs: [][]byte{[]byte("target\n"), []byte("source\n")}}
	client := Client{Binary: "chezmoi", Runner: runner}

	if _, err := client.TargetContent("/home/me/.zshrc", true, 1024); err != nil {
		t.Fatalf("TargetContent() error = %v", err)
	}
	if _, err := client.DecryptedSourceContent("/source/encrypted_dot_zshrc.age", 1024); err != nil {
		t.Fatalf("DecryptedSourceContent() error = %v", err)
	}
	want := [][]string{
		{"chezmoi", "--skip-secrets", "cat", "/home/me/.zshrc"},
		{"chezmoi", "decrypt", "/source/encrypted_dot_zshrc.age"},
	}
	if !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
}
