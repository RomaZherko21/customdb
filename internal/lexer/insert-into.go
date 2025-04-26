package lexer

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strings"
)

func ParseInsertIntoCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return model.Table{}, fmt.Errorf("ParseInsertIntoCommand(): not enough arguments")
	}

	tableName := parts[2]

	columns, err := extractColumns(input)
	if err != nil {
		return model.Table{}, fmt.Errorf("ParseInsertIntoCommand(): not found any columns")
	}

	values, err := extractValues(input)
	if err != nil {
		return model.Table{}, fmt.Errorf("ParseInsertIntoCommand(): not found any values")
	}

	result := model.Table{
		TableName: tableName,
		Columns:   []model.Column{},
		Rows:      [][]interface{}{},
	}

	for _, column := range columns {
		column = trimParentheses(column)

		result.Columns = append(result.Columns, model.Column{
			Name: column,
		})
	}

	for _, value := range values {
		value = trimParentheses(value)

		result.Rows = append(result.Rows, []interface{}{})
		result.Rows[0] = append(result.Rows[0], value)
	}

	return result, nil
}

func extractColumns(input string) ([]string, error) {
	re := regexp.MustCompile(`INTO\s+\w+\s+\((.*?)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return nil, fmt.Errorf("not found any columns")
	}

	columns := strings.Split(matches[1], ",")

	for i, col := range columns {
		columns[i] = strings.TrimSpace(col)
	}

	return columns, nil
}

func extractValues(input string) ([]string, error) {
	re := regexp.MustCompile(`VALUES\s+\((.*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("not found any values")
	}

	values := strings.Split(matches[1], ",")

	return values, nil
}
