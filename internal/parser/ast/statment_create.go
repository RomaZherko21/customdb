package ast

import "custom-database/internal/parser/lex"

type CreateTableStatement struct {
	Name lex.Token
	Cols *[]*columnDefinition
}

type columnDefinition struct {
	Name     lex.Token
	Datatype lex.Token
}

func parseCreateTableStatement(tokens []*lex.Token, initialCursor uint, delimiter lex.Token) (*CreateTableStatement, uint, bool) {
	cursor := initialCursor

	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.CreateKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if !lex.ExpectToken(tokens, cursor, lex.TokenFromKeyword(lex.TableKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	name, newCursor, ok := lex.ParseToken(tokens, cursor, lex.IdentifierToken)
	if !ok {
		lex.HelpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.LeftparenSymbol)) {
		lex.HelpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, lex.TokenFromSymbol(lex.RightparenSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.RightparenSymbol)) {
		lex.HelpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	return &CreateTableStatement{
		Name: *name,
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
			if !lex.ExpectToken(tokens, cursor, lex.TokenFromSymbol(lex.CommaSymbol)) {
				lex.HelpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for a column name
		columnName, newCursor, ok := lex.ParseToken(tokens, cursor, lex.IdentifierToken)
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		// Look for a column type
		columnType, newCursor, ok := lex.ParseToken(tokens, cursor, lex.KeywordToken)
		if !ok {
			lex.HelpMessage(tokens, cursor, "Expected column type")
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
