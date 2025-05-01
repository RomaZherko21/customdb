package ast

import "custom-database/internal/parser/lex"

func parseSelectStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*SelectStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	slct := SelectStatement{}

	exps, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{tokenFromKeyword(lex.FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	slct.SelectedColumns = *exps
	cursor = newCursor

	if expectToken(tokens, cursor, tokenFromKeyword(lex.FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
		if !ok {
			helpMessage(tokens, cursor, "Expected FROM lex.Token")
			return nil, initialCursor, false
		}

		slct.From = *from
		cursor = newCursor
	}

	return &slct, cursor, true
}
