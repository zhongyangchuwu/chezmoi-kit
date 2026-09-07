package chezmoi

import (
	"reflect"
	"testing"
)

func TestParseDumpTargetTypeRequiresOneEntry(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "file", input: `{ ".zshrc": {"type":"file"} }`, want: "file"},
		{name: "directory", input: `{ ".config": {"type":"dir"} }`, want: "dir"},
		{name: "empty", input: `{}`, wantErr: true},
		{name: "multiple", input: `{ "a": {"type":"file"}, "b": {"type":"file"} }`, wantErr: true},
		{name: "malformed", input: `{`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDumpTargetType([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseDumpTargetType() returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseDumpTargetType() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseDumpTargetType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseNULPathsPreservesNewlines(t *testing.T) {
	got := ParseNULPaths([]byte(".zshrc\x00.config/name\nwith-newline\x00"))
	want := []string{".zshrc", ".config/name\nwith-newline"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseNULPaths() = %#v, want %#v", got, want)
	}
}

func TestClientAuthoritativeDiffForcesBuiltinReverseOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("diff output\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	out, err := client.AuthoritativeDiff("/home/me/.zshrc", 1024)
	if err != nil {
		t.Fatalf("AuthoritativeDiff() error = %v", err)
	}
	if string(out) != "diff output\n" {
		t.Fatalf("AuthoritativeDiff() = %q", out)
	}
	want := [][]string{{"chezmoi", "--color=false", "--no-pager", "--use-builtin-diff", "diff", "--include=all", "--exclude=none", "--reverse", "--script-contents=true", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, want) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, want)
	}
}

func TestClientTargetMetadataUsesDumpAndTemplateFilter(t *testing.T) {
	runner := &fakeRunner{outputs: [][]byte{
		[]byte(`{ ".zshrc": {"type":"file"} }`),
		[]byte("/home/me/.zshrc\x00"),
	}}
	client := Client{Binary: "chezmoi", Runner: runner}

	metadata, err := client.TargetMetadata("/home/me/.zshrc", 1024)
	if err != nil {
		t.Fatalf("TargetMetadata() error = %v", err)
	}
	if metadata.Type != "file" || !metadata.Template {
		t.Fatalf("TargetMetadata() = %#v", metadata)
	}
	wantRunCalls := [][]string{{"chezmoi", "dump", "--include=all", "--exclude=none", "--format=json", "--recursive=false", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(runner.runCalls, wantRunCalls) {
		t.Fatalf("runCalls = %#v, want %#v", runner.runCalls, wantRunCalls)
	}
	wantOutputCalls := [][]string{{"chezmoi", "managed", "--include=templates", "--exclude=none", "--path-style=absolute", "--nul-path-separator", "/home/me/.zshrc"}}
	if !reflect.DeepEqual(runner.outputCalls, wantOutputCalls) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, wantOutputCalls)
	}
}

func TestClientTargetMetadataSkipsTemplateQueryForDirectory(t *testing.T) {
	runner := &fakeRunner{output: []byte(`{ ".config": {"type":"dir"} }`)}
	client := Client{Binary: "chezmoi", Runner: runner}

	metadata, err := client.TargetMetadata("/home/me/.config", 1024)
	if err != nil {
		t.Fatalf("TargetMetadata() error = %v", err)
	}
	if metadata.Type != "dir" || metadata.Template {
		t.Fatalf("TargetMetadata() = %#v", metadata)
	}
	if len(runner.outputCalls) != 0 {
		t.Fatalf("outputCalls = %#v, want no template query", runner.outputCalls)
	}
}

func TestClientTargetDirUsesChezmoiTargetPath(t *testing.T) {
	runner := &fakeRunner{output: []byte("/srv/dotfiles-home\n")}
	client := Client{Binary: "chezmoi", Runner: runner}

	dir, err := client.TargetDir()
	if err != nil {
		t.Fatalf("TargetDir() error = %v", err)
	}
	if dir != "/srv/dotfiles-home" {
		t.Fatalf("TargetDir() = %q", dir)
	}
	want := [][]string{{"chezmoi", "target-path"}}
	if !reflect.DeepEqual(runner.outputCalls, want) {
		t.Fatalf("outputCalls = %#v, want %#v", runner.outputCalls, want)
	}
}
