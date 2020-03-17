package common

import (
	"os"

	"github.com/kataras/tablewriter"
)

func Table(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(header)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetColWidth(40)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	return table
}
