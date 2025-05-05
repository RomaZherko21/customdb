package ast

import "custom-database/internal/parser/lex"

type columnDefinition struct {
	Name     lex.Token
	Datatype lex.Token
}

func parseCreateTableStatement(tokens []*lex.Token, initialCursor uint) (*CreateTableStatement, uint, bool) {
	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.CreateKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromKeyword(lex.TableKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	tableName, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.LeftparenSymbol)) {
		helpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, tokenFromSymbol(lex.RightparenSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.RightparenSymbol)) {
		helpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromSymbol(lex.SemicolonSymbol)) {
		helpMessage(tokens, cursor, "Expected semicolon")
		return nil, initialCursor, false
	}

	return &CreateTableStatement{
		Name: *tableName,
		Cols: cols,
	}, cursor, true
}

func parseColumnDefinitions(tokens []*lex.Token, initialCursor uint, endDelimiter lex.Token) (*[]*columnDefinition, uint, bool) {
	cursor := initialCursor

	cds := []*columnDefinition{}
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for a delimiter
		current := tokens[cursor]
		if endDelimiter.Equals(current) {
			break
		}

		// Look for a comma
		if len(cds) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(lex.CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for a column name
		columnName, newCursor, ok := parseToken(tokens, cursor, lex.IdentifierToken)
		if !ok {
			helpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		// Look for a column type
		columnType, newCursor, ok := parseToken(tokens, cursor, lex.KeywordToken)
		if !ok {
			helpMessage(tokens, cursor, "Expected column type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &columnDefinition{
			Name:     *columnName,
			Datatype: *columnType,
		})
	}

	return &cds, cursor, true
}
