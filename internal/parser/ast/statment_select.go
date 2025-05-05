package ast

import (
	"custom-database/internal/parser/lex"
)

func parseSelectStatement(tokens []*lex.Token, initialCursor uint) (*SelectStatement, uint, bool) {
	statement := &SelectStatement{
		SelectedColumns: []*Expression{},
		Where:           &WhereClause{},
	}

	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if expectToken(tokens, cursor, tokenFromSymbol(lex.AsteriskSymbol)) {
		statement.SelectedColumns = []*Expression{}
		cursor++
	} else {
		exps, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{tokenFromKeyword(lex.FromKeyword), tokenFromSymbol(lex.SemicolonSymbol)})
		if !ok {
			return nil, initialCursor, false
		}
		cursor = newCursor
		statement.SelectedColumns = *exps
	}

	// Парсим FROM (опционально)
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.FromKeyword)) {
		return statement, cursor, true
	}
	cursor++

	from, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name after FROM")
		return nil, initialCursor, false
	}
	statement.From = *from
	cursor = newCursor

	// Парсим WHERE (опционально)
	where, newCursor, ok := parseWhereClause(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol))
	if !ok {
		helpMessage(tokens, cursor, "Invalid WHERE clause")
		return nil, initialCursor, false
	}
	statement.Where = where
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
		helpMessage(tokens, cursor, "Expected semicolon")
		return nil, initialCursor, false
	}

	return statement, cursor, true
}
