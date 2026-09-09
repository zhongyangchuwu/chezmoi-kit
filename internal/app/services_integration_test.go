package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/chezmoi"
)

type chezmoiIntegration struct {
	source      string
	destination string
	config      string
	root        string
	service     service
	wrapper     string
}

func newChezmoiIntegration(t *testing.T) *chezmoiIntegration {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("integration harness requires a POSIX shell")
	}

	binary, err := exec.LookPath("chezmoi")
	if err != nil {
		t.Skip("chezmoi is unavailable")
	}

	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	cache := filepath.Join(root, "cache")
	state := filepath.Join(root, "state", "chezmoistate.boltdb")
	config := filepath.Join(root, "config", "chezmoi.toml")
	home := filepath.Join(root, "home")
	for _, dir := range []string{source, destination, cache, filepath.Dir(state), filepath.Dir(config), home} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatalf("create %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(config, nil, 0o600); err != nil {
		t.Fatalf("create isolated config: %v", err)
	}

	wrapper := filepath.Join(root, "chezmoi")
	script := fmt.Sprintf("#!/bin/sh\nHOME=%q\nexport HOME\nexec %q --use-builtin-age=true --config=%q --source=%q --destination=%q --cache=%q --persistent-state=%q \"$@\"\n", home, binary, config, source, destination, cache, state)
	if err := os.WriteFile(wrapper, []byte(script), 0o700); err != nil {
		t.Fatalf("create chezmoi wrapper: %v", err)
	}

	return &chezmoiIntegration{
		source:      source,
		destination: destination,
		config:      config,
		root:        root,
		service:     service{client: chezmoi.Client{Binary: wrapper}},
		wrapper:     wrapper,
	}
}

func (h *chezmoiIntegration) run(t *testing.T, args ...string) {
	t.Helper()
	_ = h.output(t, args...)
}

func (h *chezmoiIntegration) output(t *testing.T, args ...string) []byte {
	t.Helper()
	out, err := exec.Command(h.wrapper, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("chezmoi %q: %v: %s", args, err, strings.TrimSpace(string(out)))
	}
	return out
}

func writeIntegrationFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("chmod %s: %v", path, err)
	}
}

func TestServiceDiffOutputReportsExecutableModeDrift(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".executable")
	writeIntegrationFile(t, filepath.Join(h.source, "executable_dot_executable"), "unchanged\n", 0o700)
	writeIntegrationFile(t, target, "unchanged\n", 0o600)

	diff, err := h.service.DiffOutput(target)
	if err != nil {
		t.Fatalf("DiffOutput returned error: %v", err)
	}
	if !strings.Contains(strings.ToLower(string(diff)), "mode") {
		t.Fatalf("DiffOutput did not report executable-bit drift")
	}
}

func TestServiceDiffOutputPreservesTargetToDestinationDirection(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".direction")
	writeIntegrationFile(t, target, "target state\n", 0o600)
	h.run(t, "add", target)
	writeIntegrationFile(t, target, "destination state\n", 0o600)

	diff, err := h.service.DiffOutput(target)
	if err != nil {
		t.Fatalf("DiffOutput returned error: %v", err)
	}
	got := string(diff)
	if !strings.Contains(got, "-target state") || !strings.Contains(got, "+destination state") {
		t.Fatalf("DiffOutput direction was not target to destination: %q", got)
	}
}

func TestServiceReviewComparesSymlinkTargets(t *testing.T) {
	h := newChezmoiIntegration(t)
	firstTarget := filepath.Join(h.destination, "symlink-target-one")
	secondTarget := filepath.Join(h.destination, "symlink-target-two")
	link := filepath.Join(h.destination, ".link")
	writeIntegrationFile(t, firstTarget, "same referent contents\n", 0o600)
	writeIntegrationFile(t, secondTarget, "same referent contents\n", 0o600)
	if err := os.Symlink(filepath.Base(firstTarget), link); err != nil {
		t.Fatalf("create original symlink: %v", err)
	}
	h.run(t, "add", link)
	if err := os.Remove(link); err != nil {
		t.Fatalf("remove original symlink: %v", err)
	}
	if err := os.Symlink(filepath.Base(secondTarget), link); err != nil {
		t.Fatalf("create changed symlink: %v", err)
	}

	review, err := h.service.Review(link)
	if err != nil {
		t.Fatalf("Review returned error: %v", err)
	}
	if review.Type != TargetSymlink {
		t.Fatalf("Review type = %q, want symlink", review.Type)
	}
	if !strings.Contains(review.Diff, filepath.Base(firstTarget)) || !strings.Contains(review.Diff, filepath.Base(secondTarget)) {
		t.Fatalf("Review did not identify both symlink targets")
	}
	if strings.Contains(review.Diff, "same referent contents") {
		t.Fatalf("Review compared symlink referent contents instead of link targets")
	}
}

func TestServiceReviewRepresentsDirectoryModeDrift(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".private-directory")
	writeIntegrationFile(t, filepath.Join(target, "kept"), "unchanged\n", 0o600)
	if err := os.Chmod(target, 0o700); err != nil {
		t.Fatalf("make directory private: %v", err)
	}
	h.run(t, "add", target)
	if err := os.Chmod(target, 0o755); err != nil {
		t.Fatalf("change directory mode: %v", err)
	}

	review, err := h.service.Review(target)
	if err != nil {
		t.Fatalf("Review returned error for directory mode drift: %v", err)
	}
	if review.Type != TargetDirectory {
		t.Fatalf("Review type = %q, want directory", review.Type)
	}
	if !strings.Contains(strings.ToLower(review.Diff), "mode") {
		t.Fatalf("Review did not represent directory mode drift")
	}
}

func TestServiceReviewRepresentsRemoveEntry(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".obsolete")
	writeIntegrationFile(t, target, "remove me\n", 0o600)
	writeIntegrationFile(t, filepath.Join(h.source, "remove_dot_obsolete"), "", 0o600)

	review, err := h.service.Review(target)
	if err != nil {
		t.Fatalf("Review returned error for remove entry: %v", err)
	}
	if review.Type != TargetRemove {
		t.Fatalf("Review type = %q, want remove", review.Type)
	}
	if !review.Allows(ActionApply) || review.Allows(ActionAdd) || review.Allows(ActionMerge) {
		t.Fatalf("remove action matrix is unsafe")
	}
}

func TestServiceStatusSeparatesPendingScripts(t *testing.T) {
	h := newChezmoiIntegration(t)
	writeIntegrationFile(t, filepath.Join(h.source, "run_once_pending.sh"), "#!/bin/sh\nexit 0\n", 0o700)

	status, err := h.service.Status(nil)
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if len(status.Entries) != 0 {
		t.Fatalf("Status entries = %#v, want no ordinary entries", status.Entries)
	}
	if len(status.Scripts) != 1 || len(status.Scripts[0].Code) < 2 || status.Scripts[0].Code[1] != 'R' {
		t.Fatalf("Status scripts = %#v, want one pending script", status.Scripts)
	}
}

func TestServiceReviewDiffRendersTemplateContent(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".rendered-template")
	writeIntegrationFile(t, filepath.Join(h.source, "dot_rendered-template.tmpl"), "rendered={{ .chezmoi.os }}\n", 0o600)
	writeIntegrationFile(t, target, "rendered=local\n", 0o600)

	review, err := h.service.Review(target)
	if err != nil {
		t.Fatalf("Review returned error: %v", err)
	}
	if review.Type != TargetFile || !review.Template {
		t.Fatalf("Review metadata = type %q, template %t; want file template", review.Type, review.Template)
	}
	if !strings.Contains(review.Diff, "rendered="+runtime.GOOS) {
		t.Fatalf("Review diff did not contain rendered template output")
	}
	if strings.Contains(review.Diff, "{{ .chezmoi.os }}") {
		t.Fatalf("Review diff exposed the template source instead of rendered output")
	}
}

func TestServiceEditTargetUsesCustomDestinationNonInteractively(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".editable")
	writeIntegrationFile(t, target, "before\n", 0o600)
	h.run(t, "add", target)

	editor := filepath.Join(h.root, "editor")
	writeIntegrationFile(t, editor, "#!/bin/sh\nprintf 'edited-by-integration-test\\n' > \"$1\"\n", 0o700)
	t.Setenv("EDITOR", editor)
	t.Setenv("VISUAL", editor)
	if err := h.service.EditTarget(".editable"); err != nil {
		t.Fatalf("EditTarget returned error: %v", err)
	}
	h.run(t, "apply", "--force", target)

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read edited target: %v", err)
	}
	if string(content) != "edited-by-integration-test\n" {
		t.Fatalf("EditTarget did not edit the target resolved from the custom destination")
	}
}

func TestServiceSourceEditCommandDoesNotApplyDestination(t *testing.T) {
	h := newChezmoiIntegration(t)
	target := filepath.Join(h.destination, ".workspace-editable")
	writeIntegrationFile(t, target, "before\n", 0o600)
	h.run(t, "add", target)

	editor := filepath.Join(h.root, "workspace-editor")
	writeIntegrationFile(t, editor, "#!/bin/sh\nprintf 'edited-source\n' > \"$1\"\n", 0o700)
	t.Setenv("EDITOR", editor)
	t.Setenv("VISUAL", editor)
	command, err := h.service.SourceEditCommand(WorkspaceEntry{
		Path:       target,
		SourcePath: filepath.Join(h.source, "dot_workspace-editable"),
		State:      FileClean,
		Type:       TargetFile,
	})
	if err != nil {
		t.Fatalf("SourceEditCommand returned error: %v", err)
	}
	if err := command.Run(); err != nil {
		t.Fatalf("workspace source edit returned error: %v", err)
	}
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(content) != "before\n" {
		t.Fatalf("destination changed after workspace source edit: %q", content)
	}
	review, err := h.service.Review(target)
	if err != nil {
		t.Fatalf("Review returned error: %v", err)
	}
	if !review.Dirty || !strings.Contains(review.Diff, "edited-source") {
		t.Fatalf("review after source edit = %#v", review)
	}
}

func TestServiceSourceEditCommandEditsTemplateAndEncryptedSourcesWithoutApply(t *testing.T) {
	h := newChezmoiIntegration(t)
	identity := filepath.Join(h.root, "age-identity.txt")
	h.run(t, "age-keygen", "--output", identity)
	recipient := strings.TrimSpace(string(h.output(t, "age-keygen", "--convert", identity)))
	if recipient == "" {
		t.Fatal("age-keygen returned an empty recipient")
	}
	writeIntegrationFile(t, h.config, fmt.Sprintf("encryption = \"age\"\n[age]\nidentity = %q\nrecipient = %q\n", identity, recipient), 0o600)

	editor := filepath.Join(h.root, "workspace-editor")
	writeIntegrationFile(t, editor, "#!/bin/sh\nprintf 'edited-source\\n' > \"$1\"\n", 0o700)
	t.Setenv("EDITOR", editor)
	t.Setenv("VISUAL", editor)

	templateTarget := filepath.Join(h.destination, ".template")
	writeIntegrationFile(t, templateTarget, "template-before\n", 0o600)
	h.run(t, "add", "--template", templateTarget)
	templateSource := strings.TrimSpace(string(h.output(t, "source-path", templateTarget)))

	encryptedTarget := filepath.Join(h.destination, ".secret")
	writeIntegrationFile(t, encryptedTarget, "secret-before\n", 0o600)
	h.run(t, "add", "--encrypt", encryptedTarget)
	encryptedSource := strings.TrimSpace(string(h.output(t, "source-path", encryptedTarget)))

	for _, entry := range []WorkspaceEntry{
		{Path: templateTarget, SourcePath: templateSource, State: FileClean, Type: TargetFile, Template: true},
		{Path: encryptedTarget, SourcePath: encryptedSource, State: FileUninspected, Type: TargetFile, Encrypted: true},
	} {
		command, err := h.service.SourceEditCommand(entry)
		if err != nil {
			t.Fatalf("SourceEditCommand(%s) returned error: %v", entry.Path, err)
		}
		if err := command.Run(); err != nil {
			t.Fatalf("workspace source edit %s returned error: %v", entry.Path, err)
		}
	}

	for _, target := range []string{templateTarget, encryptedTarget} {
		content, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read destination %s: %v", target, err)
		}
		if strings.Contains(string(content), "edited-source") {
			t.Fatalf("destination changed after workspace source edit: %s = %q", target, content)
		}
		if got := string(h.output(t, "cat", target)); got != "edited-source\n" {
			t.Fatalf("rendered source after edit %s = %q", target, got)
		}
	}
	encryptedBytes, err := os.ReadFile(encryptedSource)
	if err != nil {
		t.Fatalf("read encrypted source: %v", err)
	}
	if strings.Contains(string(encryptedBytes), "edited-source") {
		t.Fatal("encrypted source contains plaintext editor content")
	}
}

func TestServiceReAddPreservesEncryptedSourceAttribute(t *testing.T) {
	h := newChezmoiIntegration(t)
	identity := filepath.Join(h.root, "age-identity.txt")
	h.run(t, "age-keygen", "--output", identity)
	recipient := strings.TrimSpace(string(h.output(t, "age-keygen", "--convert", identity)))
	if recipient == "" {
		t.Fatal("age-keygen returned an empty recipient")
	}
	config := fmt.Sprintf("encryption = \"age\"\n[age]\nidentity = %q\nrecipient = %q\n", identity, recipient)
	writeIntegrationFile(t, h.config, config, 0o600)

	target := filepath.Join(h.destination, ".secret")
	writeIntegrationFile(t, target, "first secret\n", 0o600)
	h.run(t, "add", "--encrypt", target)

	entries, err := os.ReadDir(h.source)
	if err != nil {
		t.Fatalf("read source directory: %v", err)
	}
	if len(entries) != 1 || !strings.HasPrefix(entries[0].Name(), "encrypted_") {
		t.Fatalf("encrypted source entries = %#v", entries)
	}
	sourceName := entries[0].Name()

	writeIntegrationFile(t, target, "second secret\n", 0o600)
	if _, err := h.service.ExecuteNonInteractive(Action{Target: target, Kind: ActionAdd}); err != nil {
		t.Fatalf("ExecuteNonInteractive re-add returned error: %v", err)
	}
	entries, err = os.ReadDir(h.source)
	if err != nil {
		t.Fatalf("read source directory after re-add: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != sourceName || !strings.HasPrefix(entries[0].Name(), "encrypted_") {
		t.Fatalf("re-add changed encrypted source identity: %#v", entries)
	}

	writeIntegrationFile(t, target, "stale local\n", 0o600)
	h.run(t, "apply", "--force", target)
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read applied encrypted target: %v", err)
	}
	if string(content) != "second secret\n" {
		t.Fatalf("applied content = %q, want updated encrypted content", content)
	}
}
