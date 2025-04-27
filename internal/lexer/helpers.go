package lexer

import (
	"custom-database/internal/model"
	"strings"
)

func trimParentheses(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "()")
	value = strings.TrimSpace(value)

	return value
}

func toJson(table *model.Table) string {

	json := ""

	return string(json)
}
