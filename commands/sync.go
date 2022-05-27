package commands

import (
	"log"

	"github.com/barelyhuman/musync/account"
	"github.com/barelyhuman/musync/sync"
	"github.com/barelyhuman/musync/utils"

	"github.com/urfave/cli/v2"
)

func Sync(c *cli.Context) {
	auth := account.NewAuth(
		account.WithToken(),
	)

	writer := utils.NewWriter(
		utils.WithSilent(c.Bool("quiet")),
	)

	syncer := sync.NewSyncer(
		c.String("source"),
		c.String("dest"),
		sync.WithAuthenticator(auth),
		sync.WithWriter(writer),
	)

	if err := syncer.Sync(); err != nil {
		log.Panic(err)
	}
}
