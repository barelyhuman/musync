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

	clientId := utils.GetEnvDefault("MUSYNC_CLIENT_ID", c.String("clientid"))
	clientSecret := utils.GetEnvDefault("MUSYNC_CLIENT_SECRET", c.String("clientsecret"))
	port := utils.GetEnvDefault("MUSYNC_CALLBACK_PORT", c.String("port"))

	authenticator := account.NewAuth(account.AuthConfig{
		RedirectURL: "http://localhost:" + port + "/callback",
	})

	fmt.Println(authenticator.GetAuthURL(
		clientId,
		clientSecret,
	))

	go authenticator.StartVerificationServer(port, onComplete)

	<-onComplete

	err := authenticator.SaveToken()
	if err != nil {
		log.Panic(err)
	}
}
