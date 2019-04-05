package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var cmdUpdate = &Command{
	Run:       runUpdate,
	UsageLine: "update <username>",
	Short:     "Update password using STNS API",
	Long: ` Change the password of the specified user

Options
  endpoint   STNS API Endpoint URL
  user       STNS API Basic Authentication Password
  password   STNS API Basic Authentication Passowrd
  ca         STNS API TLS Authentication CA Certificate
  cert       STNS API TLS Authentication Certificate
  key        STNS API TLS Authentication Key
  insecure   Skil TLS Verify
`,
}

var (
	endpoint   *string
	user       *string
	password   *string
	ca         *string
	cert       *string
	key        *string
	skipVerify *bool
)

func init() {
	endpoint = cmdUpdate.Flag.String("endpoint", "http://localhost:1104/v1", "STNS API Endpoint URL")
	user = cmdUpdate.Flag.String("user", "", "STNS API Basic Authentication Password")
	password = cmdUpdate.Flag.String("password", "", "STNS API Basic Authentication Passowrd")
	ca = cmdUpdate.Flag.String("ca", "", "STNS API TLS Authentication CA Certificate")
	cert = cmdUpdate.Flag.String("cert", "", "STNS API TLS Authentication Certificate")
	key = cmdUpdate.Flag.String("key", "", "STNS API TLS Authentication Key")
	skipVerify = cmdUpdate.Flag.Bool("insecure", false, "Skil TLS Verify")
}
func runUpdate(args []string) int {
	var currentPassword, newPassword []byte
	var err error

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: stns-passwd update <username>\n\n")
		return ExitCodeError

	}
	currentPassword, err = readCurrentPasswordFromStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	newPassword, err = readNewPasswordFromStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	u := strings.TrimRight(*endpoint, "/") + fmt.Sprintf("/users/password/%s", args[0])

	params := &struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}{}

	params.CurrentPassword = string(currentPassword)
	params.NewPassword = string(newPassword)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	d, err := json.Marshal(params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	body := bytes.NewBuffer(d)

	req, err := http.NewRequest("PUT", u, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	if *user != "" && *password != "" {
		req.SetBasicAuth(*user, *password)
	}

	tc := &tls.Config{InsecureSkipVerify: !*skipVerify}
	if *cert != "" && *key != "" && *ca != "" {
		ce, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(*ca)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)

		tc.Certificates = []tls.Certificate{ce}
		tc.RootCAs = caPool

		tc.BuildNameToCertificate()
	}

	tr := &http.Transport{
		TLSClientConfig: tc,
		Dial: (&net.Dialer{
			Timeout: time.Duration(10) * time.Second,
		}).Dial,
	}

	tr.Proxy = http.ProxyFromEnvironment
	client := &http.Client{Transport: tr}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		fmt.Fprintln(os.Stderr, "Failed to change password. Maybe the current password is miss match")
		return ExitCodeError

	}
	return 0
}

func readCurrentPasswordFromStdin() ([]byte, error) {
	// Read password from terminal without echo back
	var err error
	fmt.Print("Current password: ")
	rawPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return nil, err
	}
	return rawPassword, nil
}
