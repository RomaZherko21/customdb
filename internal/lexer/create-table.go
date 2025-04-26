package lexer

import (
	"errors"
	"regexp"
	"strings"
)

type Column struct {
	Name string
	Type ColumnType
}

type CreateTableCommand struct {
	TableName string
	Columns   []Column
}

func ParseCreateTableCommand(input string) (CreateTableCommand, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return CreateTableCommand{}, errors.New("parseCreateTableCommand(): not enough arguments")
	}

	re := regexp.MustCompile(`\((.*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return CreateTableCommand{}, errors.New("parseCreateTableCommand(): not found any columns")
	}

	columns := strings.Split(matches[1], ",")

	result := CreateTableCommand{
		TableName: parts[2],
		Columns:   []Column{},
	}

	for _, column := range columns {
		column = strings.TrimSpace(column)
		column = strings.Trim(column, "()")
		column = strings.TrimSpace(column)

		columnParts := strings.Split(column, " ")
		if len(columnParts) != 2 {
			return CreateTableCommand{}, errors.New("parseCreateTableCommand(): invalid column definition")
		}

		column := Column{
			Name: columnParts[0],
			Type: ColumnType(columnParts[1]),
		}

		result.Columns = append(result.Columns, column)
	}

	return result, nil
}
