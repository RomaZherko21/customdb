package lexer

import (
	"custom-database/internal/executor"
	"fmt"
	"strings"
)

type Lexer interface {
	ParseQuery(input string) ([][]interface{}, error)
}

type lexer struct {
	exec executor.Executor
}

func NewLexer(exec executor.Executor) Lexer {
	return &lexer{
		exec: exec,
	}
}

func (l *lexer) ParseQuery(input string) ([][]interface{}, error) {
	err := validateQuery(input)
	if err != nil {
		return nil, fmt.Errorf("ParseQuery(): %w", err)
	}

	keyword, err := parseKeyword(input)
	if err != nil {
		return nil, fmt.Errorf("ParseQuery(): %w", err)
	}

	switch keyword {
	case CREATE_TABLE:
		parsed, err := ParseCreateTableCommand(input)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery(): %w", err)
		}

		return nil, l.exec.CreateTable(parsed)
	case INSERT_INTO:
		parsed, err := ParseInsertIntoCommand(input)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery(): %w", err)
		}

		return nil, l.exec.InsertInto(parsed)
	case SELECT:
		parsed, err := ParseSelectCommand(input)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery(): %w", err)
		}

		result, err := l.exec.Select(parsed)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery(): %w", err)
		}

		return result, nil
	}

	return nil, nil
}

func validateQuery(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("validateQuery(): empty input")
	}

	if input[len(input)-1] != SEMICOLON {
		return fmt.Errorf("validateQuery(): command must end with a semicolon")
	}

	return nil
}

func parseKeyword(input string) (KeywordType, error) {
	parts := strings.Split(input, " ")
	command := parts[0]

	switch strings.ToUpper(command) {
	case "SELECT":
		return SELECT, nil
	case "CREATE":
		subCommand := strings.ToUpper(parts[1])

		if subCommand == "TABLE" {
			return CREATE_TABLE, nil
		}
		return "", fmt.Errorf("parseKeyword(): unknown command 'CREATE %s'", subCommand)
	case "INSERT":
		subCommand := strings.ToUpper(parts[1])
		if subCommand == "INTO" {
			return INSERT_INTO, nil
		}
		return "", fmt.Errorf("parseKeyword(): unknown command 'INSERT %s'", subCommand)
	}

	return "", fmt.Errorf("parseKeyword(): unknown command '%s'", command)
}
