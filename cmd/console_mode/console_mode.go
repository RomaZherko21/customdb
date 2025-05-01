package console_mode

import (
	"bufio"
	"custom-database/internal/backend"
	"custom-database/internal/models"
	"custom-database/internal/parser"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func RunConsoleMode(parser parser.ParserService, mb backend.MemoryBackendService) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to gosql.")

	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		result, err := parser.Parse(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		results, err := mb.ExecuteStatement(result)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if results != nil {
			printTable(results)
			continue
		}

		fmt.Println("ok")
	}
}

func printTable(results *models.Table) error {
	if len(results.Rows) == 0 {
		fmt.Println("(no results)")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	header := []string{}
	for _, col := range results.Columns {
		header = append(header, fmt.Sprintf("%s", col.Name))
	}
	table.SetHeader(header)
	table.SetAutoFormatHeaders(false)

	rows := [][]string{}
	for _, result := range results.Rows {
		row := []string{}
		for i, cell := range result {
			typ := results.Columns[i].Type
			r := ""
			switch typ {
			case models.IntType:
				i := cell.AsInt()
				r = fmt.Sprintf("%d", i)
			case models.TextType:
				r = cell.AsText()
			}

			row = append(row, r)
		}

		rows = append(rows, row)
	}

	table.SetBorder(true)
	table.AppendBulk(rows)
	table.Render()

	fmt.Printf("(%d rows)\n", len(rows))

	return nil
}
