package main

import (
	"context"
	"log"
	"os"
	"fmt"

	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/urfave/cli/v2"
)

const server = "http://localhost:8080"

func main() {
	app := &cli.App{
		Name: "cli",
		Commands: []*cli.Command{
			{
				Name: "user",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Action: createUser,
					},
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

func createUser(c *cli.Context) error {
	username := "craig"
	user := v1.UserPost{
		Username: username,
	}
	client, err := v1.NewClientWithResponses(server)
	resp, err := client.PostUsersWithResponse(context.Background(), user)
	if err != nil {
		return err
	}
	fmt.Printf("RESP %+v %d", string(resp.Body), resp.JSONDefault.Code)
	outputTable(resp)
	return err
}
