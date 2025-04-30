package ast

import (
	"custom-database/internal/parser/lex"
)

func ParseStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*Statement, uint, bool) {
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
