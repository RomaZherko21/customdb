package main

import (
	"bufio"
	// "custom-database/config"
	"custom-database/internal/backend"
	"custom-database/internal/parser"
	"custom-database/internal/parser/ast"
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
	// cfg, err := config.Load()
	// if err != nil {
	// 	log.Fatal("Error loading config:", err)
	// }

	mode := flag.String("mode", "console", "Режим работы: console или http")
	// port := flag.String("port", cfg.Port, "Порт для HTTP сервера")
	flag.Parse()

	// storage := storage.NewStorage(cfg)
	// executor := executor.NewExecutor(storage)
	// lexer := lexer.NewLexer(executor)

	// handlers := handlers.NewHttpHandlers(lexer)

	switch *mode {
	case "console":
		// console_mode.RunConsoleMode(lexer)
	case "http":
		// http_mode.RunHttpServer(handlers, *port)
	case "lex":
		newLexVersion()
	default:
		log.Fatal("Неизвестный режим работы. Используйте 'console' или 'http'")
	}
}

func newLexVersion() {
	mb := backend.NewMemoryBackend()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to gosql.")

	parser := parser.NewParser()
	for {
		fmt.Print("# ")
		text, err := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		result, err := parser.Parse(text)
		if err != nil {
			fmt.Println(err)
		}

		for _, stmt := range result.Statements {
			switch stmt.Kind {
			case ast.CreateTableKind:
				err = mb.CreateTable(result.Statements[0].CreateTableStatement)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("ok")
			case ast.DropTableKind:
				err = mb.DropTable(stmt.DropTableStatement)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("ok")
			case ast.InsertKind:
				err = mb.Insert(stmt.InsertStatement)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println("ok")
			case ast.SelectKind:
				results, err := mb.Select(stmt.SelectStatement)
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, col := range results.Columns {
					fmt.Printf("| %s ", col.Name)
				}
				fmt.Println("|")

				for i := 0; i < 20; i++ {
					fmt.Printf("=")
				}
				fmt.Println()

				for _, row := range results.Rows {
					fmt.Printf("|")

					for i, cell := range row {
						typ := results.Columns[i].Type
						s := ""
						switch typ {
						case backend.IntType:
							s = fmt.Sprintf("%d", cell.AsInt())
						case backend.TextType:
							s = cell.AsText()
						}

						fmt.Printf(" %s | ", s)
					}

					fmt.Println()
				}

				fmt.Println("ok")
			}
		}
	}
}
