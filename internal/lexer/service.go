package lexer

import (
	"custom-database/internal/executor"
	"custom-database/internal/model"
	"fmt"
	"strings"

	"custom-database/internal/lexer/ddl"
	"custom-database/internal/lexer/dml"
)

type Lexer interface {
	ParseQuery(input string) (*model.Table, error)
}

type lexer struct {
	exec executor.Executor
}

func NewLexer(exec executor.Executor) Lexer {
	return &lexer{
		exec: exec,
	}
}

func (l *lexer) ParseQuery(input string) (*model.Table, error) {
	err := validateQuery(input)
	if err != nil {
		return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
	}

	keyword, err := parseKeyword(input)
	if err != nil {
		return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
	}

	switch keyword {
	case CREATE_TABLE:
		parsed, err := ddl.ParseCreateTableCommand(input)
		if err != nil {
			return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
		}

		return &model.Table{}, l.exec.CreateTable(parsed)
	case DROP_TABLE:
		parsed, err := ddl.ParseDropTableCommand(input)
		if err != nil {
			return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
		}

		return &model.Table{}, l.exec.DropTable(parsed)
	case INSERT_INTO:
		parsed, err := dml.ParseInsertIntoCommand(input)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery(): %w", err)
		}

		return nil, l.exec.InsertInto(parsed)
	case SELECT_ACTION:
		parsed, err := dml.ParseSelectCommand(input)
		if err != nil {
			return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
		}

		result, err := l.exec.Select(parsed)
		if err != nil {
			return &model.Table{}, fmt.Errorf("ParseQuery(): %w", err)
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

func parseKeyword(input string) (SqlActionType, error) {
	parts := strings.Split(input, " ")
	command := parts[0]

	switch SqlKeywordType(strings.ToUpper(command)) {
	case SELECT:
		return SELECT_ACTION, nil
	case CREATE:
		subCommand := SqlKeywordType(strings.ToUpper(parts[1]))

		if subCommand == TABLE {
			return CREATE_TABLE, nil
		}
		return "", fmt.Errorf("parseKeyword(): unknown command 'CREATE %s'", subCommand)
	case DROP:
		subCommand := SqlKeywordType(strings.ToUpper(parts[1]))

		if subCommand == TABLE {
			return DROP_TABLE, nil
		}
		return "", fmt.Errorf("parseKeyword(): unknown command 'DROP %s'", subCommand)
	case INSERT:
		subCommand := SqlKeywordType(strings.ToUpper(parts[1]))
		if subCommand == INTO {
			return INSERT_INTO, nil
		}
		return "", fmt.Errorf("parseKeyword(): unknown command 'INSERT %s'", subCommand)
	}

	return "", fmt.Errorf("parseKeyword(): unknown command '%s'", command)
}
