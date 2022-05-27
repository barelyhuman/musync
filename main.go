package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/barelyhuman/musync/commands"
	"github.com/urfave/cli/v2"
)

//go:embed .commitlog.release
var version string

func main() {

	app := &cli.App{
		Name:  "musync",
		Usage: "Sync spotify playlists",
		Action: func(c *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login into your spotify account",
				Action: func(c *cli.Context) error {
					commands.Login(c)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "port",
						Usage:    "port to run the server on, please make sure you match this on the spotify portal, flag takes priority over MUSYNC_CLIENT_ID",
						Aliases:  []string{"p"},
						Value:    "8080",
						Required: false,
					},
					&cli.StringFlag{
						Name:  "clientid",
						Usage: "The client id from the spotify developer console, flag takes priority over MUSYNC_CLIENT_ID",
					},
					&cli.StringFlag{
						Name:  "clientsecret",
						Usage: "the client secret from the spotify developer console, flag takes priority over MUSYNC_CLIENT_SECRET",
					},
				},
			},

			{
				Name:  "whoami",
				Usage: "show the current logged in account",
				Action: func(c *cli.Context) error {
					commands.WhoAMI(c)
					return nil
				},
			},
			{
				Name:  "logout",
				Usage: "logout from your spotify account",
				Action: func(c *cli.Context) error {
					commands.Logout(c)
					return nil
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "source",
						Aliases: []string{"s"},
						Value:   "lib",
						Usage:   "Source playlist, type 'lib' to use your music library as the source",
					},
					&cli.StringFlag{
						Name:     "dest",
						Aliases:  []string{"d"},
						Usage:    "Destination playlists, playlist to transfer everything into",
						Required: true,
					},
				},
				Usage: "start a sync task between playlists / library",
				Action: func(c *cli.Context) error {
					commands.Sync(c)
					return nil
				},
			},
		},
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	app.Version = version

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
