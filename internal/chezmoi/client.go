package chezmoi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RunnerIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
}

type Runner interface {
	Output(command string, args []string, io RunnerIO) ([]byte, error)
	Run(command string, args []string, io RunnerIO) error
}

type Client struct {
	Binary string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
	Runner Runner
}

func (c Client) Status(targets []string) ([]StatusEntry, error) {
	args := make([]string, 0, 1+len(targets))
	args = append(args, "status")
	args = append(args, targets...)

	out, err := c.Output(args...)
	if err != nil {
		return nil, err
	}
	return ParseStatus(out)
}

func (c Client) Output(args ...string) ([]byte, error) {
	runner := c.runner()
	out, err := runner.Output(c.binary(), args, c.runnerIO())
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c Client) Run(args ...string) error {
	runner := c.runner()
	return runner.Run(c.binary(), args, c.runnerIO())
}

func (c Client) binary() string {
	if c.Binary != "" {
		return c.Binary
	}
	return "chezmoi"
}

func (c Client) runner() Runner {
	if c.Runner != nil {
		return c.Runner
	}
	return execRunner{}
}

func (c Client) runnerIO() RunnerIO {
	return RunnerIO{
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Dir:    c.Dir,
	}
}

type execRunner struct{}

func (execRunner) Output(command string, args []string, io RunnerIO) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = io.Stdin
	cmd.Stderr = writerOrDefault(io.Stderr, os.Stderr)
	cmd.Dir = io.Dir

	out, err := cmd.Output()
	if errors.Is(err, exec.ErrNotFound) {
		return nil, fmt.Errorf("chezmoi not found in PATH")
	}
	if err != nil {
		return nil, err
	}
	return bytes.TrimRight(out, "\n"), nil
}

func (execRunner) Run(command string, args []string, io RunnerIO) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = readerOrDefault(io.Stdin, os.Stdin)
	cmd.Stdout = writerOrDefault(io.Stdout, os.Stdout)
	cmd.Stderr = writerOrDefault(io.Stderr, os.Stderr)
	cmd.Dir = io.Dir

	err := cmd.Run()
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("chezmoi not found in PATH")
	}
	return err
}

func readerOrDefault(r io.Reader, fallback io.Reader) io.Reader {
	if r != nil {
		return r
	}
	return fallback
}

func writerOrDefault(w io.Writer, fallback io.Writer) io.Writer {
	if w != nil {
		return w
	}
	return fallback
}
