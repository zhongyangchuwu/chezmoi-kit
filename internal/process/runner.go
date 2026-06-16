package process

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
}

type Runner interface {
	Output(command string, args []string, io IO) ([]byte, error)
	Run(command string, args []string, io IO) error
}

type ExecRunner struct{}

func (ExecRunner) Output(command string, args []string, io IO) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = io.Stdin
	cmd.Stderr = writerOrDefault(io.Stderr, os.Stderr)
	cmd.Dir = io.Dir

	out, err := cmd.Output()
	if errors.Is(err, exec.ErrNotFound) {
		return nil, fmt.Errorf("%s not found in PATH", command)
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (ExecRunner) Run(command string, args []string, io IO) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = readerOrDefault(io.Stdin, os.Stdin)
	cmd.Stdout = writerOrDefault(io.Stdout, os.Stdout)
	cmd.Stderr = writerOrDefault(io.Stderr, os.Stderr)
	cmd.Dir = io.Dir

	err := cmd.Run()
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("%s not found in PATH", command)
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
