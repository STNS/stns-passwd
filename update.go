package main

var cmdUpdate = &Command{
	Run:       runUpdate,
	UsageLine: "update ",
	Short:     "Update password using STNS API",
	Long: `

	`,
}
var (
	endpoint = cmdUpdate.Flag.String("endpoint", "http://localhost:1104/v1", "STNS API Endpoint URL")
	user     = cmdUpdate.Flag.String("user", "", "STNS API Basic Authentication Password")
	password = cmdUpdate.Flag.String("password", "", "STNS API Basic Authentication Passowrd")
	cert     = cmdUpdate.Flag.String("cert", "", "STNS API TLS Authentication Certificate")
	key      = cmdUpdate.Flag.String("key", "", "STNS API TLS Authentication Key")
)

func runUpdate(args []string) int {

	return 0
}
