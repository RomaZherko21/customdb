package lexer

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strings"
)

func ParseCreateTableCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return model.Table{}, fmt.Errorf("parseCreateTableCommand(): not enough arguments")
	}

	re := regexp.MustCompile(`\((.*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return model.Table{}, fmt.Errorf("parseCreateTableCommand(): not found any columns")
	}

	columns := strings.Split(matches[1], ",")

	result := model.Table{
		TableName: parts[2],
		Columns:   []model.Column{},
	}

	for _, column := range columns {
		column = strings.TrimSpace(column)
		column = strings.Trim(column, "()")
		column = strings.TrimSpace(column)

		columnParts := strings.Split(column, " ")
		if len(columnParts) != 2 {
			return model.Table{}, fmt.Errorf("parseCreateTableCommand(): invalid column definition")
		}

		column := model.Column{
			Name: columnParts[0],
			Type: model.DataType(columnParts[1]),
		}

		result.Columns = append(result.Columns, column)
	}

	return result, nil
}
