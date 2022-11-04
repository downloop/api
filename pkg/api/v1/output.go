package v1

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type outputtable interface {
	tableHeaders() []string
	tableData() [][]string
}

type OutputWriter struct {
	SuccessResponseType outputtable
	Format              string
}

func (w OutputWriter) Write() error {
	switch w.Format {
	case "table":
		w.WriteTable()
	case "json":
		return w.WriteJSON()
	}
	return nil
}

func (w OutputWriter) WriteTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(w.SuccessResponseType.tableHeaders())
	table.AppendBulk(w.SuccessResponseType.tableData())
	table.Render()
}

func (w OutputWriter) WriteJSON() error {
	data, err := json.Marshal(w.SuccessResponseType)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func (SessionResponseList) tableHeaders() []string {
	return []string{"Start Time", "End Time"}
}

func (r SessionResponseList) tableData() [][]string {
	var data [][]string
	for _, session := range r.Data {
		data = append(data, []string{session.StartTime.String(), session.EndTime.String()})
	}
	return data
}

func (UserResponseList) tableHeaders() []string {
	return []string{"Username"}
}

func (r UserResponseList) tableData() [][]string {
	var data [][]string
	for _, user := range r.Data {
		data = append(data, []string{user.Username})
	}
	return data
}
