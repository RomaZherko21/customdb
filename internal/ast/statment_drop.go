package ast

import "custom-database/internal/lex"

type DropTableStatement struct {
	Table lex.Token
}

func parseDropTableStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*DropTableStatement, uint, bool) {
	cursor := initialCursor

	// Look for DROP
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.DropKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for TABLE
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.TableKeyword)) {
		lex.HelpMessage(tokens, cursor, "Expected table")
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

	return &DropTableStatement{
		Table: *table,
	}, cursor, true
}
