package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в консоль! Введите 'exit' для выхода.")

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

		parseCommand(input)
	}
}

// commands
const SELECT = "SELECT"
const CREATE = "CREATE"
const INSERT = "INSERT"

// instances
const TABLE = "TABLE"

func parseCommand(input string) {
	if len(input) == 0 {
		return
	}

	if input[len(input)-1] != ';' {
		fmt.Println("Ошибка: команда должна заканчиваться точкой с запятой")
		return
	}

	parts := strings.Split(input, " ")
	command := parts[0]

	switch command {
	case SELECT:
		parseSelectCommand(input)
	case CREATE:
		parseCreateCommand(input)
	case INSERT:
		parseInsertCommand(input)
	}
}

func parseSelectCommand(input string) {
	parts := strings.Split(input, " ")
	tableName := parts[1]

	fmt.Printf("Выбрали таблицу: %s\n", tableName)
}

// CREATE TABLE users (id INT, name TEXT);
func parseCreateCommand(input string) {
	parts := strings.Split(input, " ")
	instance := parts[1]

	switch instance {
	case TABLE:
		createTable(input)
	}
}

func createTable(input string) {
	parts := strings.Split(input, " ")
	tableName := parts[2]

	// как получить все что в скобках например из строки // CREATE TABLE users (id INT, name TEXT);?
	re := regexp.MustCompile(`\((.*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		fmt.Println("Ошибка: не найдены столбцы")
		return
	}

	columns := strings.Split(matches[1], ",")
	for _, column := range columns {
		column = strings.TrimSpace(column)
		column = strings.Trim(column, "()")
		column = strings.TrimSpace(column)
		fmt.Println("Столбец:", column)
	}

	fmt.Printf("Создали таблицу: %s\n", tableName)
}

func parseInsertCommand(input string) {
	parts := strings.Split(input, " ")
	tableName := parts[2]

	fmt.Printf("Вставили в таблицу: %s\n", tableName)
}
