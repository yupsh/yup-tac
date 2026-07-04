// Command yup-tac is the CLI wrapper around github.com/gloo-foo/cmd-tac.
package main

import (
	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-tac"
	urf "github.com/urfave/cli/v3"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const (
	name          = "tac"
	flagSeparator = "separator"
)

// synopsis is the multi-line --help usage block; urfave/cli indents it three
// spaces, so the lines stay flush-left.
const synopsis = `tac [OPTIONS] [FILE...]

Write each FILE to standard output, last line first.
With no FILE, or when FILE is -, read standard input.`

// spec declares the tac wrapper: a file-or-stdin filter with a separator flag.
var spec = clix.Spec{
	Name:     name,
	Summary:  "concatenate and print files in reverse",
	Synopsis: synopsis,
	Build:    build,
	Flags:    flags(),
}

// flags returns fresh flag instances. It is a constructor rather than a shared
// slice because urfave/cli flag structs retain per-parse state, so each parse
// (including tests) must build its own.
func flags() []urf.Flag {
	return []urf.Flag{
		&urf.StringFlag{
			Name:    flagSeparator,
			Aliases: []string{"s"},
			Usage:   "use STRING as the record separator instead of newline",
		},
	}
}

// build maps the invocation to tac's pipeline: a file-or-stdin source into the
// tac command configured by the flags.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	return clix.OperandsOrStdin(inv), command.Tac(options(inv.Args)...), nil
}

// options folds the parsed flags into tac's option values.
func options(c *urf.Command) []any {
	var opts []any
	if c.IsSet(flagSeparator) {
		opts = append(opts, command.TacSep(c.String(flagSeparator)))
	}
	return opts
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
