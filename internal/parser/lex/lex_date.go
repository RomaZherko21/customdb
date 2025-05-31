package lex

import "time"

func lexDate(source string, ic Cursor) (*Token, Cursor, bool) {
	newToken, newCursor, ok := lexCharacterDelimited(source, ic, '\'')
	if !ok {
		return nil, ic, false
	}

	// проверяем, что дата в формате 'YYYY-MM-DD HH:MM:SS'
	_, err := time.Parse("2006-01-02 15:04:05", newToken.Value)
	if err != nil {
		return nil, ic, false
	}

	return &Token{
		Value: newToken.Value,
		Kind:  DateToken,
	}, newCursor, true
}
