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
	SuccessResponse outputtable
	Format          string
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
	table.SetHeader(w.SuccessResponse.tableHeaders())
	table.AppendBulk(w.SuccessResponse.tableData())
	table.Render()
}

func (w OutputWriter) WriteJSON() error {
	data, err := json.Marshal(w.SuccessResponse)
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

func (SessionResponse) tableHeaders() []string {
	return []string{"Start Time", "End Time"}
}

func (r SessionResponse) tableData() [][]string {
	return [][]string{[]string{r.Data.StartTime.String(), r.Data.EndTime.String()}}
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
