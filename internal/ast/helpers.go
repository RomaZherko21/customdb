package ast

import "custom-database/internal/lex"

func parseExpressions(tokens []*lex.Token, initialCursor uint, delimiters []lex.Token) (*[]*expression, uint, bool) {
	cursor := initialCursor

	exps := []*expression{}
outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for delimiter
		current := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.Equals(current) {
				break outer
			}
		}

		// Look for comma
		if len(exps) > 0 {
			if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.CommaSymbol)) {
				lex.HelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, lex.TokenFromSymbol(lex.CommaSymbol))
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*lex.Token, initialCursor uint, _ lex.Token) (*expression, uint, bool) {
	cursor := initialCursor

	kinds := []lex.TokenKind{lex.IdentifierToken, lex.NumericToken, lex.StringToken}
	for _, kind := range kinds {
		t, newCursor, ok := lex.ParseToken(tokens, cursor, kind)
		if ok {
			return &expression{
				literal: t,
				kind:    literalKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}
