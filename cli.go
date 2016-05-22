package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/STNS/libnss_stns/hash"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		s string
		c int
		m string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&s, "s", "", "String used to salt. The SNS is the user name")

	flags.IntVar(&c, "c", 0, "The number of times to stretching")

	flags.StringVar(&m, "m", "sha256", "Specifies the hash function(sha256/sha512)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if flags.Arg(0) == "" {
		fmt.Fprintln(cli.errStream, "Please specify a password")
		return ExitCodeError
	}

	fmt.Println(hash.Calculate(m, s != "", s, flags.Arg(0), c))

	return ExitCodeOK
}
