package lexer

import "strings"

func trimParentheses(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "()")
	value = strings.TrimSpace(value)

	return value
}
