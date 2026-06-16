package chezmoi

import (
	"io"

	"github.com/zhongyangchuwu/cm/internal/process"
)

type RunnerIO = process.IO

type Runner = process.Runner

type Client struct {
	Binary string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
	Runner Runner
}

func (c Client) Status(targets []string) ([]StatusEntry, error) {
	args := make([]string, 0, 2+len(targets))
	args = append(args, "status", "--path-style=absolute")
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
	return process.ExecRunner{}
}

func (c Client) runnerIO() RunnerIO {
	return RunnerIO{
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Dir:    c.Dir,
	}
}
