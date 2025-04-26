package main

import (
	"bufio"
	"custom-database/config"
	"custom-database/internal/executor"
	"custom-database/internal/lexer"
	"custom-database/internal/storage"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в консоль! Введите 'exit' для выхода.")

	storage := storage.NewStorage(cfg)
	executor := executor.NewExecutor(storage)
	lexer := lexer.NewLexer(executor)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("До свидания!")
			return
		}

		err = lexer.ParseQuery(input)
		if err != nil {
			fmt.Println("Main():", err)
			continue
		}
	}
}
