package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

var cmdGenerate = &Command{
	Run:       runGenerate,
	UsageLine: "generate ",
	Short:     "Generate password from stdin or typing(default command)",
	Long: `

	`,
}

func runGenerate(args []string) int {
	var rawPassword []byte
	if len(args) != 0 && args[0] != "" {
		rawPassword = []byte(args[0])
	} else {
		var err error
		rawPassword, err = readNewPasswordFromStdin()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}
	}
	if len(rawPassword) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify a password")
	}

	c := sha512_crypt.New()
	v, err := c.Generate(rawPassword, []byte{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	fmt.Fprintln(os.Stdout, v)
	return ExitCodeOK
}

func readNewPasswordFromStdin() ([]byte, error) {
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
