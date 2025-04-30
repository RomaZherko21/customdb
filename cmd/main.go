package main

import (
	"bufio"
	"custom-database/config"
	"custom-database/internal/backend"
	"custom-database/internal/models"
	"custom-database/internal/parser"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// @title Custom Database API
// @version 1.0
// @description API для работы с кастомной базой данных
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	mode := flag.String("mode", "console", "Режим работы: console или http")
	// port := flag.String("port", cfg.Port, "Порт для HTTP сервера")
	flag.Parse()

	parser := parser.NewParser()
	mb, err := backend.NewMemoryBackend(cfg)
	if err != nil {
		log.Fatal("Error creating memory backend:", err)
	}

	// handlers := handlers.NewHttpHandlers(lexer)

	switch *mode {
	case "console":
		// console_mode.RunConsoleMode(lexer)
	case "http":
		// http_mode.RunHttpServer(handlers, *port)
	case "lex":
		newLexVersion(parser, mb)
	default:
		log.Fatal("Неизвестный режим работы. Используйте 'console' или 'http'")
	}
}

func newLexVersion(parser parser.ParserService, mb backend.MemoryBackendService) {
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
