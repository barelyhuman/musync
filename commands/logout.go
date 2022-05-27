package commands

import (
	"fmt"

	"github.com/barelyhuman/musync/account"
	"github.com/urfave/cli/v2"
)

func Logout(c *cli.Context) {
	storage := account.NewTokenStorage()
	storage.ClearToken()
	fmt.Printf("Logged out!")
}
