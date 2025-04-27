package ast

import (
	"custom-database/internal/lex"
	"errors"
	"fmt"
)

func Parse(source string) (*Ast, error) {
	tokens, err := lex.Lex(source)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		fmt.Println(t.Value, t.Kind)
	}

	a := Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, lex.TokenFromSymbol(lex.SemicolonSymbol))
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("Failed to parse, expected statement")
		}
		cursor = newCursor

		a.Statements = append(a.Statements, stmt)

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

func parseStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*Statement, uint, bool) {
	cursor := initialCursor

	// Look for a SELECT statement
	semicolonToken := lex.TokenFromSymbol(lex.SemicolonSymbol)
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

	return nil, initialCursor, false
}
