package parser

import (
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"errors"
)

type ParserService interface {
	Parse(query string) (*ast.Ast, error)
}

type Parser struct {
}

func NewParser() ParserService {
	return &Parser{}
}

func (p *Parser) Parse(query string) (*ast.Ast, error) {
	tokens, err := lex.Lex(query)
	if err != nil {
		return nil, err
	}

	a := ast.Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		statement, newCursor, ok := ast.ParseStatement(tokens, cursor, lex.TokenFromSymbol(lex.SemicolonSymbol))
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("Failed to parse, expected statement")
		}
		cursor = newCursor

		a.Statements = append(a.Statements, statement)

		atLeastOneSemicolon := false
		for lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.SemicolonSymbol)) {
			cursor++
			atLeastOneSemicolon = true
		}

		if !atLeastOneSemicolon {
			lex.HelpMessage(tokens, cursor, "Expected semi-colon delimiter between statements")
			return nil, errors.New("Missing semi-colon between statements")
		}
	}

	return &a, nil
}
