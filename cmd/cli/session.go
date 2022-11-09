package main

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/urfave/cli/v2"
)

func sessionCreate(c *cli.Context) error {
	session := v1.SessionPost{
		StartTime: time.Now(),
	}

	client, err := newClient()
	if err != nil {
		return err
	}
	resp, err := client.PostSessionsWithResponse(context.Background(), session)
	if err != nil {
		return err
	}
	o := v1.OutputWriter{
		SuccessResponse: resp.JSON201,
		Format:          c.Generic("format").(*EnumValue).String(),
	}
	err = o.Write()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func sessionsGet(c *cli.Context) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	resp, err := client.GetSessionsWithResponse(context.Background())
	if err != nil {
		return err
	}
	o := v1.OutputWriter{
		SuccessResponse: resp.JSON200,
		Format:          c.Generic("format").(*EnumValue).String(),
	}
	err = o.Write()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
