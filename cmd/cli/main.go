package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const server = "http://localhost:8080"

func main() {
	app := &cli.App{
		Name: "cli",
		Commands: []*cli.Command{
			{
				Name: "auth",
				Subcommands: []*cli.Command{
					{
						Name:   "login",
						Action: authenticateUser,
					},
				},
			},
			{
				Name: "user",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Action: userCreate,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "username",
								Required: true,
							},
						},
					},
					{
						Name:   "get",
						Action: usersGet,
					},
				},
			},
			{
				Name: "session",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Action: sessionCreate,
					},
					{
						Name:   "get",
						Action: sessionsGet,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.GenericFlag{
				Name: "format, f",
				Value: &EnumValue{
					Enum:    []string{"json", "table", "xml"},
					Default: "json",
				},
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
