package commands

import (
	"log"

	"github.com/barelyhuman/musync/account"
	"github.com/barelyhuman/musync/sync"

	"github.com/urfave/cli/v2"
)

func Sync(c *cli.Context) {
	auth := account.AuthFromToken()

	syncer := sync.NewSyncer(
		c.String("source"),
		c.String("dest"),
		sync.WithAuthenticator(auth),
	)

	if err := syncer.Sync(); err != nil {
		log.Panic(err)
	}
}
