package ast

import (
	"custom-database/internal/parser/lex"
	"errors"
)

type AstService interface {
	Parse(query string) (*Ast, error)
}

type ast struct {
}

func NewAst() AstService {
	return &ast{}
}

func (s *ast) Parse(query string) (*Ast, error) {
	lexer := lex.NewLexer()
	tokens, err := lexer.Lex(query)
	if err != nil {
		return nil, err
	}

	result := Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		statement, newCursor, ok := parseStatement(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("Failed to parse, expected statement")
		}
		cursor = newCursor

		result.Statements = append(result.Statements, statement)

		atLeastOneSemicolon := false
		for expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
			cursor++
			atLeastOneSemicolon = true
		}

		if !atLeastOneSemicolon {
			helpMessage(tokens, cursor, "Expected semi-colon delimiter between statements")
			return nil, errors.New("Missing semi-colon between statements")
		}
	}

	return &result, nil
}

func parseStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*Statement, uint, bool) {
	cursor := initialCursor
	semicolonToken := tokenFromSymbol(lex.SemicolonSymbol)

	// Look for a SELECT statement
	slct, newCursor, ok := parseSelectStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:            SelectKind,
			SelectStatement: slct,
		}, newCursor, true
	}

	// Look for a INSERT statement
	inst, newCursor, ok := parseInsertStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:            InsertKind,
			InsertStatement: inst,
		}, newCursor, true
	}

	// Look for a CREATE statement
	crtTbl, newCursor, ok := parseCreateTableStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:                 CreateTableKind,
			CreateTableStatement: crtTbl,
		}, newCursor, true
	}

	// Look for a DROP statement
	dropTbl, newCursor, ok := parseDropTableStatement(tokens, cursor, semicolonToken)
	if ok {
		return &Statement{
			Kind:               DropTableKind,
			DropTableStatement: dropTbl,
		}, newCursor, true
	}

	return nil, initialCursor, false
}
