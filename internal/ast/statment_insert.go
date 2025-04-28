package ast

import "custom-database/internal/lex"

type InsertStatement struct {
	Table  lex.Token
	Values *[]*Expression
}

func parseInsertStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*InsertStatement, uint, bool) {
	cursor := initialCursor

	// Look for INSERT
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.InsertKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for INTO
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.IntoKeyword)) {
		lex.HelpMessage(tokens, cursor, "Expected into")
		return nil, initialCursor, false
	}
	cursor++

	// Look for table name
	table, newCursor, ok := lex.ParseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		lex.HelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for VALUES
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.ValuesKeyword)) {
		lex.HelpMessage(tokens, cursor, "Expected VALUES")
		return nil, initialCursor, false
	}
	cursor++

	// Look for left paren
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.LeftparenSymbol)) {
		lex.HelpMessage(tokens, cursor, "Expected left paren")
		return nil, initialCursor, false
	}
	cursor++

	// Look for expression list
	values, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{lex.TokenFromSymbol(lex.RightparenSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for right paren
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.RightparenSymbol)) {
		lex.HelpMessage(tokens, cursor, "Expected right paren")
		return nil, initialCursor, false
	}
	cursor++

	return &InsertStatement{
		Table:  *table,
		Values: values,
	}, cursor, true
}
