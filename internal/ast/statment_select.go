package ast

import "custom-database/internal/lex"

type SelectStatement struct {
	item []*expression
	from lex.Token
}

func parseSelectStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*SelectStatement, uint, bool) {
	cursor := initialCursor
	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	slct := SelectStatement{}

	exps, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{lex.TokenFromKeyword(lex.FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	slct.item = *exps
	cursor = newCursor

	if lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.FromKeyword)) {
		cursor++

		from, newCursor, ok := lex.ParseToken(tokens, cursor, lex.IdentifierToken)
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected FROM lex.Token")
			return nil, initialCursor, false
		}

		slct.from = *from
		cursor = newCursor
	}

	return &slct, cursor, true
}
