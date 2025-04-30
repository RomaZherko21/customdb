package console_mode

import (
	"bufio"
	"custom-database/internal/backend"
	"custom-database/internal/models"
	"custom-database/internal/parser"
	"fmt"
	"os"
	"strings"
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

func printTable(table *models.Table) {
	for _, col := range table.Columns {
		fmt.Printf("| %s ", col.Name)
	}
	fmt.Println("|")

	for i := 0; i < 20; i++ {
		fmt.Printf("=")
	}
	fmt.Println()

	for _, row := range table.Rows {
		fmt.Printf("|")

		for i, cell := range row {
			typ := table.Columns[i].Type
			s := ""
			switch typ {
			case models.IntType:
				s = fmt.Sprintf("%d", cell.AsInt())
			case models.TextType:
				s = cell.AsText()
			}

			fmt.Printf(" %s | ", s)
		}

		fmt.Println()
	}
}
