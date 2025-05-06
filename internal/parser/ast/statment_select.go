package ast

import (
	"custom-database/internal/parser/lex"
	"strconv"
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

	limit, newCursor, ok := parseLimit(tokens, cursor)
	if !ok {
		helpMessage(tokens, cursor, "Expected limit value")
		return nil, initialCursor, false
	}
	statement.Limit = limit
	cursor = newCursor

	offset, newCursor, ok := parseOffset(tokens, cursor)
	if !ok {
		helpMessage(tokens, cursor, "Expected offset value")
		return nil, initialCursor, false
	}
	statement.Offset = offset
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
		helpMessage(tokens, cursor, "Expected semicolon")
		return nil, initialCursor, false
	}

	return statement, cursor, true
}

func parseLimit(tokens []*lex.Token, initialCursor uint) (int, uint, bool) {
	limit := 0

	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.LimitKeyword)) {
		return 0, initialCursor, true
	}
	cursor++

	limitToken, newCursor, ok := parseToken(tokens, cursor, lex.NumericToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected limit value")
		return 0, initialCursor, false
	}
	cursor = newCursor
	limit, err := strconv.Atoi(limitToken.Value)
	if err != nil {
		helpMessage(tokens, cursor, "Expected limit value")
		return 0, initialCursor, false
	}

	return limit, newCursor, true
}

func parseOffset(tokens []*lex.Token, initialCursor uint) (int, uint, bool) {
	offset := 0

	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.OffsetKeyword)) {
		return 0, initialCursor, true
	}
	cursor++

	offsetToken, newCursor, ok := parseToken(tokens, cursor, lex.NumericToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected offset value")
		return 0, initialCursor, false
	}
	cursor = newCursor
	offset, err := strconv.Atoi(offsetToken.Value)
	if err != nil {
		helpMessage(tokens, cursor, "Expected offset value")
		return 0, initialCursor, false
	}

	return offset, newCursor, true
}
