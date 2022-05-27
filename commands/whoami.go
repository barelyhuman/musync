package commands

import (
	"fmt"

	"github.com/barelyhuman/musync/account"
	"github.com/barelyhuman/musync/utils"
	"github.com/urfave/cli/v2"
)

func WhoAMI(c *cli.Context) {
	auth := account.NewAuth(
		account.WithToken(),
	)
	writer := utils.NewWriter()
	if !auth.PingUser() {
		fmt.Println("Couldn't get user data, please log in again")
	}
	user, _ := auth.Client.CurrentUser()
	writer.Info("Logged in as: " + user.DisplayName)
}
