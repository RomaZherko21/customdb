package lexer

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	Type  TokenType
	Value string
}

type LexerState struct {
	Keyword KeywordType
}

func Lexer(input string) ([]Token, error) {
	if len(input) == 0 {
		return nil, errors.New("Lexer(): empty input")
	}

	if input[len(input)-1] != SEMICOLON {
		return nil, errors.New("Lexer(): command must end with a semicolon")
	}

	var tokens []Token

	keyword, err := parseKeyword(input)
	if err != nil {
		return nil, err
	}

	fmt.Println(keyword)

	return tokens, nil
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
		return "", errors.New(fmt.Sprintf("Lexer(): unknown command 'CREATE %s'", subCommand))
	case "INSERT":
		subCommand := parts[1]
		if subCommand == "INTO" {
			return INSERT_INTO, nil
		}
		return "", errors.New(fmt.Sprintf("Lexer(): unknown command 'INSERT %s'", subCommand))
	}

	return "", errors.New(fmt.Sprintf("Lexer(): unknown command '%s'", command))
}

func parseCreateCommand(input string) (string, error) {
	parts := strings.Split(input, " ")
	tableName := parts[2]

	return tableName, nil
}
