package main

import (
	"os"

	v1 "github.com/downloop/api/pkg/api/v1"
	"github.com/olekukonko/tablewriter"
)

func outputTable(resp interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	switch resp.(type) {
	case *v1.PostUsersResponse:
		resp.(*v1.PostUsersResponse).TableOutput(table)
	}
	table.Render()

	return nil
}
