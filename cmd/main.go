package main

import (
	"bufio"
	"custom-database/internal/executor"
	"custom-database/internal/lexer"
	"custom-database/internal/storage"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в консоль! Введите 'exit' для выхода.")

	storage := storage.NewStorage()
	executor := executor.NewExecutor(storage)
	lexer := lexer.NewLexer(executor)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			continue
		}

		// Удаляем символ новой строки
		input = strings.TrimSpace(input)

		// Проверяем команду выхода
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
