package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kless/osutil/user/crypt/sha512_crypt"
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
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

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

	c := sha512_crypt.New()
	v, err := c.Generate([]byte(flags.Arg(0)), []byte{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	fmt.Println(v)
	return ExitCodeOK
}
