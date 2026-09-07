package chezmoi

import (
	"bytes"
	"errors"
	"io"

	"github.com/zhongyangchuwu/cm/internal/process"
)

type Client struct {
	Binary string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
	Runner process.Runner
}

func (c Client) Status(targets []string) ([]StatusEntry, error) {
	return c.status(targets, false)
}

func (c Client) StatusForInventory(targets []string) ([]StatusEntry, error) {
	return c.status(targets, true)
}

func (c Client) status(targets []string, skipSecrets bool) ([]StatusEntry, error) {
	args := make([]string, 0, 5+len(targets))
	if skipSecrets {
		args = append(args, "--skip-secrets")
	}
	args = append(args, "status", "--include=all", "--exclude=none", "--path-style=absolute")
	args = append(args, targets...)

	out, err := c.Output(args...)
	if err != nil {
		return nil, err
	}
	return ParseStatus(out)
}

func (c Client) Output(args ...string) ([]byte, error) {
	runner := c.ActiveRunner()
	out, err := runner.Output(c.binary(), args, c.runnerIO())
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c Client) OutputLimit(limit int64, args ...string) ([]byte, error) {
	if limit < 0 {
		limit = 0
	}
	var out limitedBuffer
	out.limit = limit
	io := c.runnerIO()
	io.Stdout = &out
	err := c.ActiveRunner().Run(c.binary(), args, io)
	if errors.Is(err, errLimitExceeded) {
		return out.Bytes(), nil
	}
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c Client) Run(args ...string) error {
	runner := c.ActiveRunner()
	return runner.Run(c.binary(), args, c.runnerIO())
}

func (c Client) RunBuffered(args ...string) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	io := process.IO{Stdout: &stdout, Stderr: &stderr, Dir: c.Dir}
	err := c.ActiveRunner().Run(c.binary(), args, io)
	return stdout.Bytes(), stderr.Bytes(), err
}

func (c Client) ManagedFiles() ([]string, error) {
	out, err := c.Output("managed", "--nul-path-separator")
	if err != nil {
		return nil, err
	}
	return ParseNULPaths(out), nil
}

func (c Client) BinaryName() string {
	return c.binary()
}

func (c Client) binary() string {
	if c.Binary != "" {
		return c.Binary
	}
	return "chezmoi"
}

func (c Client) ActiveRunner() process.Runner {
	if c.Runner != nil {
		return c.Runner
	}
	return process.ExecRunner{}
}

func (c Client) runnerIO() process.IO {
	return process.IO{
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Dir:    c.Dir,
	}
}

var errLimitExceeded = errors.New("output limit exceeded")

type limitedBuffer struct {
	bytes.Buffer
	limit int64
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	remaining := b.limit + 1 - int64(b.Len())
	if remaining <= 0 {
		return 0, errLimitExceeded
	}
	if int64(len(p)) <= remaining {
		_, _ = b.Buffer.Write(p)
		return len(p), nil
	}
	_, _ = b.Buffer.Write(p[:remaining])
	return int(remaining), errLimitExceeded
}
