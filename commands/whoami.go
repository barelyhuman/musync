package commands

import (
	"fmt"

	"github.com/barelyhuman/musync/account"
	"github.com/urfave/cli/v2"
)

func WhoAMI(c *cli.Context) {
	auth := account.AuthFromToken()
	if auth.PingUser() {
		user, _ := auth.Client.CurrentUser()
		fmt.Println("Logged in as:", user.DisplayName)
	}
}
