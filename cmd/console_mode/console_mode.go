package console_mode

import (
	"custom-database/internal/backend"
	"custom-database/internal/models"
	"custom-database/internal/parser"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"github.com/olekukonko/tablewriter"
)

func RunConsoleMode(parser parser.ParserService, mb backend.MemoryBackendService) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "# ",
		HistoryFile:     "/tmp/tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("Welcome to custom-database.")

repl:
	for {
		fmt.Print("# ")
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue repl
			}
		} else if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading line:", err)
			continue repl
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "quit" || trimmed == "exit" || trimmed == "\\q" {
			break
		}

		isDebug := strings.HasPrefix(trimmed, "ddl")
		if isDebug {
			tableName := strings.TrimSpace(trimmed[len("ddl"):])
			line = fmt.Sprintf("SELECT * FROM %s;", tableName)
		}

		result, err := parser.Parse(line)
		if err != nil {
			fmt.Println(err)
			continue repl
		}

		results, err := mb.ExecuteStatement(result)
		if err != nil {
			fmt.Println(err)
			continue repl
		}

		if isDebug {
			debugTable(results)
			continue repl
		}

		if results != nil {
			printTable(results)
			continue repl
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
				if cell.IsNull() {
					r = "null"
				} else {
					i := cell.AsInt()
					r = fmt.Sprintf("%d", i)
				}
			case models.TextType:
				if cell.IsNull() {
					r = "null"
				} else {
					r = cell.AsText()
				}
			case models.BoolType:
				if cell.IsNull() {
					r = "null"
				} else {
					b := cell.AsBoolean()
					r = fmt.Sprintf("%t", b)
				}
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

func debugTable(results *models.Table) {

	if results == nil {
		fmt.Println("Did not find any relation.")
		return
	}

	fmt.Printf("Table \"%s\"\n", results.Name)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Column", "Type", "Nullable"})
	table.SetAutoFormatHeaders(false)
	table.SetBorder(true)

	rows := [][]string{}
	for _, c := range results.Columns {
		typeString := "integer"
		switch c.Type {
		case models.TextType:
			typeString = "text"
		case models.BoolType:
			typeString = "boolean"
		}
		nullable := ""
		rows = append(rows, []string{c.Name, typeString, nullable})
	}

	table.AppendBulk(rows)
	table.Render()

	fmt.Println("")
}
