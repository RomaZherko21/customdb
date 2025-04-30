package lex

import (
	"fmt"
)

type lexer func(string, Cursor) (*Token, Cursor, bool)

func Lex(source string) ([]*Token, error) {
	tokens := []*Token{}
	cur := Cursor{}

lex:
	for cur.Pointer < uint(len(source)) {
		lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
		for _, lexFunc := range lexers {
			if token, newCursor, ok := lexFunc(source, cur); ok {
				cur = newCursor

				// Omit nil tokens for valid, but empty syntax like newlines
				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.Loc.Line, cur.Loc.Col)
	}

	return tokens, nil
}
