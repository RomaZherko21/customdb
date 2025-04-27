package lex

import "strings"

func lexIdentifier(source string, ic Cursor) (*Token, Cursor, bool) {
	// Handle separately if is a double-quoted identifier
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	c := source[cur.Pointer]
	// Other characters count too, big ignoring non-ascii for now
	isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	if !isAlphabetical {
		return nil, ic, false
	}
	cur.Pointer++
	cur.Loc.Col++

	value := []byte{c}
	for ; cur.Pointer < uint(len(source)); cur.Pointer++ {
		c = source[cur.Pointer]

		// Other characters count too, big ignoring non-ascii for now
		isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNumeric := c >= '0' && c <= '9'
		if isAlphabetical || isNumeric || c == '$' || c == '_' {
			value = append(value, c)
			cur.Loc.Col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &Token{
		// Unquoted dentifiers are case-insensitive
		Value: strings.ToLower(string(value)),
		Loc:   ic.Loc,
		Kind:  IdentifierToken,
	}, cur, true
}
