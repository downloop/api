package main

import (
	"context"
	"fmt"
	"io/ioutil"

	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/urfave/cli/v2"
)

func userCreate(c *cli.Context) error {

	username := c.String("username")
	user := v1.UserPost{
		Username: username,
	}

	client, err := newClient()
	if err != nil {
		return err
	}
	resp, err := client.PostUsers(context.Background(), user)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response: %s\n", body)
	return nil
}

func usersGet(c *cli.Context) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	resp, err := client.GetUsersWithResponse(context.Background(), &v1.GetUsersParams{})
	if err != nil {
		return err
	}
	o := v1.OutputWriter{
		SuccessResponseType: resp.JSON200,
		Format:              c.Generic("format").(*EnumValue).String(),
	}
	err = o.Write()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
