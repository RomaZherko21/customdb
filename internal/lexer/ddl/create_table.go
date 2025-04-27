package ddl

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strings"
)

func ParseCreateTableCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return model.Table{}, fmt.Errorf("ParseCreateTableCommand(): not enough arguments")
	}

	columns, err := extractColumnsWithTypes(input)
	if err != nil {
		return model.Table{}, fmt.Errorf("ParseCreateTableCommand(): %w", err)
	}

	return model.Table{
		TableName: parts[2],
		Columns:   columns,
	}, nil
}

func extractColumnsWithTypes(input string) ([]model.Column, error) {
	re := regexp.MustCompile(`\((.*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("ParseCreateTableCommand(): not found any columns")
	}

	columns := strings.Split(matches[1], ",")

	for i, column := range columns {
		columns[i] = strings.TrimSpace(column)
	}

	result := []model.Column{}

	for _, column := range columns {
		columnParts := strings.Split(column, " ")
		if len(columnParts) != 2 {
			return nil, fmt.Errorf("ParseCreateTableCommand(): invalid column definition")
		}

		result = append(result, model.Column{Name: columnParts[0], Type: model.DataType(columnParts[1])})
	}

	return result, nil
}
