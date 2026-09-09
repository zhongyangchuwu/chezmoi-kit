package cli

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/zhongyangchuwu/cm/internal/app"
	"github.com/zhongyangchuwu/cm/internal/report"
)

func TestRunDefaultsToReadOnlyStatus(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !reflect.DeepEqual(service.statusReportArgs, [][]string{nil}) {
		t.Fatalf("statusReportArgs = %#v", service.statusReportArgs)
	}
	if len(service.commands) != 0 {
		t.Fatalf("mutating/diff commands were called: %#v", service.commands)
	}
	if got := strings.TrimSpace(out.String()); got != "clean" {
		t.Fatalf("output = %q, want clean", got)
	}
}

func TestRunStatusWritesAppReportOutput(t *testing.T) {
	service := &fakeService{statusReport: "app status\n"}
	var out bytes.Buffer

	code := run([]string{"status"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if out.String() != "app status\n" {
		t.Fatalf("output = %q", out.String())
	}
	if !reflect.DeepEqual(service.statusReportArgs, [][]string{nil}) {
		t.Fatalf("statusReportArgs = %#v", service.statusReportArgs)
	}
}

func TestRunDiffWritesAppReportOutput(t *testing.T) {
	service := &fakeService{diffOutput: "internal diff\n"}
	var out bytes.Buffer

	code := run([]string{"diff", ".zshrc"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if out.String() != "internal diff\n" {
		t.Fatalf("output = %q", out.String())
	}
	if !reflect.DeepEqual(service.diffArgs, [][]string{{".zshrc"}}) {
		t.Fatalf("diffArgs = %#v", service.diffArgs)
	}
}

func TestRunReportOutputFlagSelectsMarkdown(t *testing.T) {
	service := &fakeService{diffOutput: "--- source\n+++ local\n-old\n+new\n"}
	var out bytes.Buffer

	code := run([]string{"diff", "--output", "markdown", ".zshrc"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	want := "```diff\n--- source\n+++ local\n-old\n+new\n```\n"
	if out.String() != want {
		t.Fatalf("output = %q, want %q", out.String(), want)
	}
	if !reflect.DeepEqual(service.diffArgs, [][]string{{".zshrc"}}) {
		t.Fatalf("diffArgs = %#v", service.diffArgs)
	}
}

func TestRunReportColorFlagRequestsANSIOnNonTTY(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	service := &fakeService{}
	var out bytes.Buffer
	code := run([]string{"version", "--color", "always"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("output = %q, want ANSI escapes", out.String())
	}
}

func TestRunReportFlagsRejectInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "output", args: []string{"status", "--output", "json"}, want: `unsupported output "json"`},
		{name: "color", args: []string{"status", "--color", "sometimes"}, want: `unsupported color "sometimes"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeService{}
			var out bytes.Buffer

			code := run(tt.args, testServices(service), strings.NewReader(""), &out, &out)

			if code != 1 {
				t.Fatalf("run exit code = %d, want 1", code)
			}
			if !strings.Contains(out.String(), tt.want) {
				t.Fatalf("output = %q, want substring %q", out.String(), tt.want)
			}
		})
	}
}

func TestRunSyncExecutesConfirmedReviewedAction(t *testing.T) {
	target := "/home/me/.zshrc"
	dirty := app.Review{
		Entry:       app.ReconcileEntry{Code: "MM", Path: target},
		Type:        app.TargetFile,
		Diff:        "diff",
		Fingerprint: "reviewed",
		Dirty:       true,
	}
	reviewReady := make(chan struct{})
	service := &fakeService{
		statusResults: []app.SyncStatus{{Entries: []app.ReconcileEntry{{Code: "MM", Path: target}}}},
		reviews:       []app.Review{dirty, dirty, {Entry: app.ReconcileEntry{Path: target}, Dirty: false}},
		reviewReady:   reviewReady,
	}
	var out bytes.Buffer
	input := &gatedReader{ready: reviewReady, reader: strings.NewReader("a\ry")}

	code := run([]string{"sync", ".zshrc"}, testServices(service), input, &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	wantActions := []app.Action{{Target: target, Kind: app.ActionAdd, Fingerprint: "reviewed"}}
	if !reflect.DeepEqual(service.executed, wantActions) {
		t.Fatalf("executed = %#v, want %#v", service.executed, wantActions)
	}
}

func TestRunWorkspaceKeepsCleanManagedEntriesVisible(t *testing.T) {
	target := "/home/me/.config/app/config"
	service := &fakeService{workspaceSnapshot: app.WorkspaceSnapshot{
		Root: "/home/me",
		Entries: []app.WorkspaceEntry{{
			Path:         target,
			RelativePath: ".config/app/config",
			State:        app.FileClean,
			Type:         app.TargetFile,
		}},
	}}
	var out bytes.Buffer

	code := run([]string{"ui", ".config/app"}, testServices(service), strings.NewReader("q"), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !reflect.DeepEqual(service.inventoryScopes, [][]string{{".config/app"}}) {
		t.Fatalf("inventory scopes = %#v", service.inventoryScopes)
	}
	if !strings.Contains(out.String(), "cm ui") || !strings.Contains(out.String(), "config") {
		t.Fatalf("workspace output = %q", out.String())
	}
	if len(service.executed) != 0 {
		t.Fatalf("workspace executed mutations: %#v", service.executed)
	}
}

func TestRunSyncDebugWritesTempLogPathToStderr(t *testing.T) {
	service := &fakeService{}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"sync", "--debug"}, testServices(service), strings.NewReader(""), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; stderr %q", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "clean" {
		t.Fatalf("stdout = %q, want clean", stdout.String())
	}
	debug := stderr.String()
	if !strings.Contains(debug, "debug log: ") || !strings.Contains(debug, "debug log kept at: ") {
		t.Fatalf("stderr = %q, want debug log path messages", debug)
	}
	path := debugLogPathFromLine(t, debug, "debug log: ")
	defer os.Remove(path)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}
	if !strings.Contains(string(content), "sync initial status") {
		t.Fatalf("debug log %q does not contain initial status: %q", path, string(content))
	}
}

func TestRunGitOpensSourceRepositoryWithLazygit(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"git"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"git"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
	if service.statusCalls != 0 {
		t.Fatalf("statusCalls = %d, want 0", service.statusCalls)
	}
}

func TestRunMutatingWrappersForwardToChezmoi(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want [][]string
	}{
		{name: "add", args: []string{"add", ".zshrc"}, want: [][]string{{"add", ".zshrc"}}},
		{name: "apply", args: []string{"apply", ".gitconfig"}, want: [][]string{{"apply", ".gitconfig"}}},
		{name: "merge", args: []string{"merge", ".config/nvim/init.lua"}, want: [][]string{{"merge", ".config/nvim/init.lua"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeService{}
			var out bytes.Buffer

			code := run(tt.args, testServices(service), strings.NewReader(""), &out, &out)

			if code != 0 {
				t.Fatalf("run exit code = %d, want 0", code)
			}
			if !reflect.DeepEqual(service.commands, tt.want) {
				t.Fatalf("commands = %#v, want %#v", service.commands, tt.want)
			}
			if service.statusCalls != 0 {
				t.Fatalf("statusCalls = %d, want 0", service.statusCalls)
			}
		})
	}
}

func TestRunEditForwardsTargetToChezmoi(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"edit", ".zshrc"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !reflect.DeepEqual(service.commands, [][]string{{"edit", ".zshrc"}}) {
		t.Fatalf("commands = %#v", service.commands)
	}
}

func TestRunVersionPrintsBuildInfo(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"version"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	got := out.String()
	for _, want := range []string{"cm:", "go:"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %q", got, want)
		}
	}
}

func TestRunDoctorWritesReportOutput(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"doctor", "--output", "markdown"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; output %q", code, out.String())
	}
	if !strings.Contains(out.String(), "# **doctor:**") || !strings.Contains(out.String(), "pass") {
		t.Fatalf("output = %q, want doctor markdown", out.String())
	}
}

func TestRunDoctorFailsWhenRequiredCheckFails(t *testing.T) {
	service := &fakeService{doctorFailed: true}
	var out bytes.Buffer

	code := run([]string{"doctor"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 1 {
		t.Fatalf("run exit code = %d, want 1", code)
	}
	if !strings.Contains(out.String(), "fail required") || !strings.Contains(out.String(), "doctor found failed checks") {
		t.Fatalf("output = %q, want failed doctor report and error", out.String())
	}
}

func TestRunCompletionPrintsShellScript(t *testing.T) {
	service := &fakeService{}
	var out bytes.Buffer

	code := run([]string{"completion", "bash"}, testServices(service), strings.NewReader(""), &out, &out)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "cm") {
		t.Fatalf("completion output = %q, want generated script", out.String())
	}
}

func testServices(service *fakeService) commandServices {
	return commandServices{
		Status:    service,
		Diff:      service,
		Sync:      service,
		Workspace: service,
		Target:    service,
		SourceGit: service,
		Edit:      service,
		Doctor:    service,
	}
}

type fakeService struct {
	status            app.SyncStatus
	statusResults     []app.SyncStatus
	reviews           []app.Review
	reviewReady       chan struct{}
	reviewReadyOnce   sync.Once
	statusReport      string
	statusCalls       int
	statusArgs        [][]string
	statusReportArgs  [][]string
	diffArgs          [][]string
	commands          [][]string
	diffOutput        string
	doctorFailed      bool
	executed          []app.Action
	workspaceSnapshot app.WorkspaceSnapshot
	inventoryScopes   [][]string
}

func (f *fakeService) Status(targets []string) (app.SyncStatus, error) {
	f.statusCalls++
	f.statusArgs = append(f.statusArgs, append([]string(nil), targets...))
	if len(f.statusResults) > 0 {
		status := f.statusResults[0]
		f.statusResults = f.statusResults[1:]
		return status, nil
	}
	return f.status, nil
}

func (f *fakeService) StatusReport(targets []string) (report.Document, error) {
	f.statusReportArgs = append(f.statusReportArgs, append([]string(nil), targets...))
	if f.statusReport != "" {
		return report.Document{Blocks: []report.Block{report.CodeBlock("", f.statusReport)}}, nil
	}
	return report.Document{Blocks: []report.Block{report.Paragraph(report.Text("clean"))}}, nil
}

func (f *fakeService) DiffReport(targets []string) (report.Document, error) {
	f.diffArgs = append(f.diffArgs, append([]string(nil), targets...))
	return report.Document{Blocks: []report.Block{report.DiffBlock(f.diffOutput)}}, nil
}

func (f *fakeService) DoctorReport() (report.Document, bool) {
	status := report.Status("pass")
	text := "optional"
	if f.doctorFailed {
		status = report.Warning("fail")
		text = "required"
	}
	return report.Document{Blocks: []report.Block{
		report.Heading(1, report.Strong("doctor:")),
		report.Paragraph(status, report.Text(" "), report.Text(text)),
	}}, f.doctorFailed
}

func (f *fakeService) Review(target string) (app.Review, error) {
	if f.reviewReady != nil {
		f.reviewReadyOnce.Do(func() { close(f.reviewReady) })
	}
	if len(f.reviews) == 0 {
		return app.Review{Entry: app.ReconcileEntry{Path: target}}, nil
	}
	review := f.reviews[0]
	f.reviews = f.reviews[1:]
	return review, nil
}

func (f *fakeService) Inventory(scopes []string) (app.WorkspaceSnapshot, error) {
	f.inventoryScopes = append(f.inventoryScopes, append([]string(nil), scopes...))
	return f.workspaceSnapshot, nil
}

func (f *fakeService) Preview(entry app.WorkspaceEntry, kind app.PreviewKind, reveal bool) (app.WorkspacePreview, error) {
	return app.WorkspacePreview{Entry: entry, Kind: kind, Notice: "destination matches rendered target"}, nil
}

func (f *fakeService) SourceEditCommand(entry app.WorkspaceEntry) (app.TerminalCommand, error) {
	return terminalCommand{}, nil
}

func (f *fakeService) ExecuteNonInteractive(action app.Action) (app.ActionResult, error) {
	f.executed = append(f.executed, action)
	return app.ActionResult{}, nil
}

func (f *fakeService) TerminalCommand(action app.Action) (app.TerminalCommand, error) {
	return terminalCommand{run: func() error {
		f.executed = append(f.executed, action)
		return nil
	}}, nil
}

type gatedReader struct {
	once   sync.Once
	ready  <-chan struct{}
	reader io.Reader
}

func (r *gatedReader) Read(p []byte) (int, error) {
	r.once.Do(func() { <-r.ready })
	return r.reader.Read(p)
}

type terminalCommand struct {
	run func() error
}

func (c terminalCommand) Run() error {
	if c.run == nil {
		return nil
	}
	return c.run()
}

func (terminalCommand) SetStdin(io.Reader)  {}
func (terminalCommand) SetStdout(io.Writer) {}
func (terminalCommand) SetStderr(io.Writer) {}

func debugLogPathFromLine(t *testing.T, output, prefix string) string {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	t.Fatalf("output %q missing prefix %q", output, prefix)
	return ""
}

func (f *fakeService) AddTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"add"}, targets...))
	return nil
}

func (f *fakeService) ApplyTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"apply"}, targets...))
	return nil
}

func (f *fakeService) MergeTargets(targets []string) error {
	f.commands = append(f.commands, append([]string{"merge"}, targets...))
	return nil
}

func (f *fakeService) OpenSourceGit() error {
	f.commands = append(f.commands, []string{"git"})
	return nil
}

func (f *fakeService) EditTarget(target string) error {
	f.commands = append(f.commands, []string{"edit", target})
	return nil
}

func (f *fakeService) ManagedFiles() ([]string, error) {
	return []string{".zshrc", ".gitconfig"}, nil
}
