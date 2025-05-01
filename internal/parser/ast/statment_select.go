package ast

import "custom-database/internal/parser/lex"

func parseSelectStatement(tokens []*lex.Token, initialCursor uint) (*SelectStatement, uint, bool) {
	statement := &SelectStatement{
		SelectedColumns: []*Expression{},
		Where:           []*WhereClause{},
	}

	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	exps, newCursor, ok := parseExpressions(tokens, cursor, []lex.Token{tokenFromKeyword(lex.FromKeyword), tokenFromSymbol(lex.SemicolonSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor
	statement.SelectedColumns = *exps

	// Парсим FROM (опционально)
	if !expectToken(tokens, cursor, tokenFromKeyword(lex.FromKeyword)) {
		return statement, cursor, true
	}
	cursor++

	from, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name after FROM")
		return nil, initialCursor, false
	}
	statement.From = *from
	cursor = newCursor

	// Парсим WHERE (опционально)
	where, newCursor, ok := parseWhereClause(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol))
	if !ok {
		helpMessage(tokens, cursor, "Invalid WHERE clause")
		return nil, initialCursor, false
	}
	statement.Where = where
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
		helpMessage(tokens, cursor, "Expected semicolon")
		return nil, initialCursor, false
	}

	return statement, cursor, true
}

func parseWhereClause(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) ([]*WhereClause, uint, bool) {
	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.WhereKeyword)) {
		return nil, initialCursor, true
	}
	cursor++

	whereExps, newCursor, ok := parseWhereExpression(tokens, cursor, []lex.Token{delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	return whereExps, newCursor, true
}

func parseWhereExpression(tokens []*lex.Token, initialCursor uint, delimiters []lex.Token) ([]*WhereClause, uint, bool) {
	result := []*WhereClause{}
	cursor := initialCursor

loop:
	for cursor < uint(len(tokens)) {
		currentToken := tokens[cursor]

		if isDelimiter(currentToken, delimiters) {
			break loop
		}

		clause := &WhereClause{}

		if currentToken.Kind == lex.LogicalOperatorToken {
			clause.Operator = *currentToken
			cursor++
			if !isValidCursor(tokens, cursor) {
				return nil, initialCursor, false
			}
			result = append(result, clause)
			continue
		}

		if isOperand(currentToken) {
			clause.Left = &Expression{Literal: currentToken}
			cursor++
			if !isValidCursor(tokens, cursor) {
				return nil, initialCursor, false
			}
			currentToken = tokens[cursor]
		}

		if currentToken.Kind == lex.MathOperatorToken {
			clause.Operator = *currentToken
			cursor++
			if !isValidCursor(tokens, cursor) {
				return nil, initialCursor, false
			}
			currentToken = tokens[cursor]
		}

		if isOperand(currentToken) {
			clause.Right = &Expression{Literal: currentToken}
			cursor++
		}

		if !isValidWhereClause(clause) {
			helpMessage(tokens, cursor, "Invalid WHERE clause: missing operand or operator")
			return nil, initialCursor, false
		}

		result = append(result, clause)
	}

	return result, cursor, true
}

func isDelimiter(token *lex.Token, delimiters []lex.Token) bool {
	for _, delimiter := range delimiters {
		if delimiter.Equals(token) {
			return true
		}
	}
	return false
}

func isValidCursor(tokens []*lex.Token, cursor uint) bool {
	if cursor >= uint(len(tokens)) {
		helpMessage(tokens, cursor, "Unexpected end of input")
		return false
	}
	return true
}

func isOperand(token *lex.Token) bool {
	return token.Kind == lex.IdentifierToken ||
		token.Kind == lex.StringToken ||
		token.Kind == lex.NumericToken
}

func isValidWhereClause(clause *WhereClause) bool {
	return clause.Left != nil &&
		clause.Operator.Value != "" &&
		clause.Right != nil
}
