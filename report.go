package sizeof

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// PrintReport prints SpaceUsageReport to the given io.Writer.
func PrintReport(r *SpaceUsageReport, w io.Writer) {
	printReport(r, w, "")
}

func printReport(r *SpaceUsageReport, w io.Writer, tab string) {
	fields := make([][2]string, 0, len(r.Values))

	for k, v := range r.Values {
		fields = append(fields, [2]string{k, v})
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i][0] < fields[j][0]
	})

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorder(false)
	table.SetCenterSeparator("|")

	// Header
	header := make([]string, len(fields))

	for i := range fields {
		header[i] = strings.ReplaceAll(fields[i][0], "_", "")
	}

	table.SetHeader(header)

	fs := make([]string, len(fields))
	for i := range fields {
		fs[i] = fields[i][1]
	}

	table.Append(fs)

	table.Render()

	strs := strings.Split(tableString.String(), "\n")
	for i := range strs {
		fmt.Fprintf(w, "%s%s\n", tab, strs[i])
	}

	for i := range r.Children {
		printReport(r.Children[i], w, fmt.Sprintf("\t%s", tab))
	}
}
