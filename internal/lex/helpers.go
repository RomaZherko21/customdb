package lex

import (
	"fmt"
	"strings"
)

// longestMatch iterates through a source string starting at the given
// cursor to find the longest matching substring among the provided
// options
func longestMatch(source string, ic Cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.Pointer < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.Pointer]))...)
		cur.Pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.Pointer-ic.Pointer]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

func TokenFromKeyword(k keyword) Token {
	return Token{
		Kind:  KeywordToken,
		Value: string(k),
	}
}

func TokenFromSymbol(s symbol) Token {
	return Token{
		Kind:  SymbolToken,
		Value: string(s),
	}
}

func ParseToken(tokens []*Token, initialCursor uint, kind TokenKind) (*Token, uint, bool) {
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

func ExpectToken(tokens []*Token, cursor uint, t Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}

func HelpMessage(tokens []*Token, cursor uint, msg string) {
	var c *Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, msg, c.Value)
}
