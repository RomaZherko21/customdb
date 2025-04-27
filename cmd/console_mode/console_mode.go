package console_mode

import (
	"bufio"
	"custom-database/internal/lexer"
	"fmt"
	"os"
	"strings"
)

func RunConsoleMode(lexer lexer.Lexer) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в консоль! Введите 'exit' для выхода.")

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

		result, err := lexer.ParseQuery(input)
		if err != nil {
			fmt.Println("Main():", err)
			continue
		}

		if result == nil {
			fmt.Println("Query executed successfully")
			continue
		}

		fmt.Println(result)
	}
}
