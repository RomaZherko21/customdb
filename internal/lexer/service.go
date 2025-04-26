package lexer

import (
	"custom-database/internal/executor"
	"fmt"
	"strings"
)

type Lexer interface {
	ParseQuery(input string) error
}

type lexer struct {
	exec executor.Executor
}

func NewLexer(exec executor.Executor) Lexer {
	return &lexer{
		exec: exec,
	}
}

func (l *lexer) ParseQuery(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("ParseQuery(): empty input")
	}

	if input[len(input)-1] != SEMICOLON {
		return fmt.Errorf("ParseQuery(): command must end with a semicolon")
	}

	keyword, err := parseKeyword(input)
	if err != nil {
		return err
	}

	switch keyword {
	case CREATE_TABLE:
		result, err := ParseCreateTableCommand(input)
		if err != nil {
			return fmt.Errorf("ParseQuery(): %w", err)
		}

		return l.exec.CreateTable(result)
	case INSERT_INTO:
		result, err := ParseInsertIntoCommand(input)
		if err != nil {
			return fmt.Errorf("ParseQuery(): %w", err)
		}

		return l.exec.InsertInto(result)
	}

	return nil
}

func parseKeyword(input string) (KeywordType, error) {
	parts := strings.Split(input, " ")
	command := parts[0]

	switch command {
	case "SELECT":
		return SELECT, nil
	case "CREATE":
		subCommand := parts[1]

		if subCommand == "TABLE" {
			return CREATE_TABLE, nil
		}
		return "", fmt.Errorf("Lexer(): unknown command 'CREATE %s'", subCommand)
	case "INSERT":
		subCommand := parts[1]
		if subCommand == "INTO" {
			return INSERT_INTO, nil
		}
		return "", fmt.Errorf("Lexer(): unknown command 'INSERT %s'", subCommand)
	}

	return "", fmt.Errorf("Lexer(): unknown command '%s'", command)
}
