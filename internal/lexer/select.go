package lexer

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strings"
)

func ParseSelectCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 1 {
		return model.Table{}, fmt.Errorf("ParseSelectCommand(): not enough arguments")
	}

	tableName := strings.Trim(parts[len(parts)-1], ";")

	columns, err := extractSelectColumns(input)
	if err != nil {
		return model.Table{}, fmt.Errorf("ParseSelectCommand(): not found any columns")
	}

	return model.Table{
		TableName: tableName,
		Columns:   columns,
	}, nil
}

func extractSelectColumns(input string) ([]model.Column, error) {
	re := regexp.MustCompile(`SELECT\s+(.*)\s+FROM\s+(.*);`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("extractSelectColumns(): not found any columns")
	}

	columnsStr := strings.TrimSpace(matches[1])

	if columnsStr == "*" {
		return []model.Column{}, nil
	}

	columns := strings.Split(columnsStr, ",")

	result := make([]model.Column, 0)

	for _, column := range columns {
		column = strings.TrimSpace(column)

		columnArr := strings.Split(column, " ")
		if len(columnArr) > 1 {
			return nil, fmt.Errorf("extractSelectColumns(): invalid column name: %s", column)
		}

		result = append(result, model.Column{Name: column})
	}

	return result, nil
}
