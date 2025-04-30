package ast

import (
	"custom-database/internal/parser/lex"
	"fmt"
)

func parseExpressions(tokens []*lex.Token, initialCursor uint, delimiters []lex.Token) (*[]*Expression, uint, bool) {
	cursor := initialCursor

	exps := []*Expression{}
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
			if !expectToken(tokens, cursor, tokenFromSymbol(lex.CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(lex.CommaSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*lex.Token, initialCursor uint, _ lex.Token) (*Expression, uint, bool) {
	cursor := initialCursor

	kinds := []lex.TokenKind{lex.IdentifierToken, lex.NumericToken, lex.StringToken}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &Expression{
				Literal: t,
				Kind:    LiteralKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func tokenFromKeyword(k lex.Keyword) lex.Token {
	return lex.Token{
		Kind:  lex.KeywordToken,
		Value: string(k),
	}
}

func tokenFromSymbol(s lex.Symbol) lex.Token {
	return lex.Token{
		Kind:  lex.SymbolToken,
		Value: string(s),
	}
}

func parseToken(tokens []*lex.Token, initialCursor uint, kind lex.TokenKind) (*lex.Token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	current := tokens[cursor]
	if current.Kind == kind {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func expectToken(tokens []*lex.Token, cursor uint, t lex.Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func helpMessage(tokens []*lex.Token, cursor uint, msg string) {
	var c *lex.Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, msg, c.Value)
}
