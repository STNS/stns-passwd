package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"github.com/mattn/go-isatty"
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

	var rawPassword []byte
	if flags.Arg(0) != "" {
		rawPassword = []byte(flags.Arg(0))
	} else {
		var err error
		rawPassword, err = readPasswordFromStdin()
		if err != nil {
			fmt.Fprintln(cli.errStream, err)
			return ExitCodeError
		}
	}
	if len(rawPassword) == 0 {
		fmt.Fprintln(cli.errStream, "Please specify a password")
		return ExitCodeError
	}

	c := sha512_crypt.New()
	v, err := c.Generate(rawPassword, []byte{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	fmt.Println(v)
	return ExitCodeOK
}

func readPasswordFromStdin() ([]byte, error) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		// Read password from terminal without echo back
		var err error
		fmt.Print("Enter password: ")
		rawPassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return nil, err
		}
		fmt.Print("Retype password: ")
		verify, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Println()
		if !bytes.Equal(rawPassword, verify) {
			return nil, errors.New("Sorry, passwords do not match")
		}
		return rawPassword, nil
	} else {
		// Read password from stdin (not a terminal)
		s := bufio.NewScanner(os.Stdin)
		if s.Scan() {
			return s.Bytes(), nil
		}
		return nil, s.Err()
	}

}
