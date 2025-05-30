package ast

import "custom-database/internal/parser/lex"

func parseDropTableStatement(tokens []*lex.Token, initialCursor uint) (*DropTableStatement, uint, bool) {
	cursor := initialCursor

	// Look for DROP
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.DropKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for TABLE
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.TableKeyword)) {
		helpMessage(tokens, cursor, "Expected table")
		return nil, initialCursor, false
	}
	cursor++

	// Look for table name
	table, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
		helpMessage(tokens, cursor, "Expected semicolon")
		return nil, initialCursor, false
	}

	return &DropTableStatement{
		Table: *table,
	}, cursor, true
}
