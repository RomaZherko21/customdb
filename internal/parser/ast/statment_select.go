package ast

import "custom-database/internal/parser/lex"

// SELECT id, name FROM users WHERE id = 5 AND name = 'John';
func parseSelectStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*SelectStatement, uint, bool) {
	statement := SelectStatement{}

	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	exps, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{tokenFromKeyword(lex.FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	statement.SelectedColumns = *exps

	// FROM table не всегда будет. Может быть такой запрос SELECT 2 + 2;
	if expectToken(tokens, cursor, tokenFromKeyword(lex.FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
		if !ok {
			helpMessage(tokens, cursor, "Expected FROM lex.Token")
			return nil, initialCursor, false
		}

		statement.From = *from
		cursor = newCursor

		where, newCursor, ok := parseWhereClause(tokens, cursor, delimiter)
		if ok {
			statement.Where = where
			cursor = newCursor
		}
	}

	return &statement, cursor, true
}

// SELECT id, name FROM users WHERE id = 5 AND name = 'John';
func parseWhereClause(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) ([]*WhereClause, uint, bool) {
	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.WhereKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// delimiter может быть ; или GROUP BY или ORDER BY или LIMIT или OFFSET ??
	whereExps, newCursor, ok := parseWhereExpression(tokens, cursor, []lex.Token{delimiter})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	return whereExps, cursor, true
}

func parseWhereExpression(tokens []*lex.Token, initialCursor uint, delimiters []lex.Token) ([]*WhereClause, uint, bool) {
	result := []*WhereClause{}
	cursor := initialCursor

	clause := WhereClause{}
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		currentToken := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.Equals(currentToken) {
				break
			}
		}

		if currentToken.Kind == lex.IdentifierToken {
			clause.Left = &Expression{Literal: currentToken}
			cursor++
			continue
		}

		// if currentToken.Kind == lex.OperatorToken {
		// 	clause.Op = currentToken
		// 	cursor++
		// }

	}

	return result, cursor, true
}
