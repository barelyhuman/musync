package commands

import (
	"fmt"
	"log"

	"github.com/barelyhuman/musync/account"
	"github.com/barelyhuman/musync/utils"
	"github.com/urfave/cli/v2"
)

func Login(c *cli.Context) {
	onComplete := make(chan bool)
	writer := utils.NewWriter()

	clientId := utils.GetEnvDefault("MUSYNC_CLIENT_ID", c.String("clientid"))
	clientSecret := utils.GetEnvDefault("MUSYNC_CLIENT_SECRET", c.String("clientsecret"))
	port := utils.GetEnvDefault("MUSYNC_CALLBACK_PORT", c.String("port"))

	authenticator := account.NewAuth(
		account.WithoutToken("http://localhost:" + port + "/callback"),
	)

	u := authenticator.GetAuthURL(
		clientId,
		clientSecret,
	)
	writer.Info(
		"Please open the below url in your browser to log into the service:",
	)
	fmt.Printf("%s\n\n", u)

	go authenticator.StartVerificationServer(port, onComplete)

	<-onComplete

	err := authenticator.SaveToken()
	if err != nil {
		log.Panic(err)
	}
}
