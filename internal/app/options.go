package app

import "io"

// Options holds process-wide CLI flags and runtime sinks shared by commands.
type Options struct {
	Debug  bool
	Stderr io.Writer
}

func NormalizeOptions(opts *Options) *Options {
	if opts == nil {
		return &Options{Stderr: io.Discard}
	}
	if opts.Stderr == nil {
		opts.Stderr = io.Discard
	}
	return opts
}
