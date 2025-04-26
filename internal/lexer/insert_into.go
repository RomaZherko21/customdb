package lexer

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strconv"
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

	if len(columns) != len(values) {
		return model.Table{}, fmt.Errorf("ParseInsertIntoCommand(): mismatched columns and values")
	}

	result := model.Table{
		TableName: tableName,
		Columns:   []model.Column{},
		Rows:      [][]interface{}{},
	}
	result.Rows = append(result.Rows, []interface{}{})

	for _, column := range columns {
		column = trimParentheses(column)

		result.Columns = append(result.Columns, model.Column{
			Name: column,
		})
	}

	for _, value := range values {
		value = trimParentheses(value)

		val, err := extractValue(value)
		if err != nil {
			return model.Table{}, fmt.Errorf("ParseInsertIntoCommand(): invalid value: %s", value)
		}

		result.Rows[0] = append(result.Rows[0], val)
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

	for i, value := range values {
		values[i] = strings.TrimSpace(value)
	}

	return values, nil
}

func extractValue(value string) (interface{}, error) {
	var err error

	if value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1], nil
	}
	if strings.ToUpper(value) == "NULL" {
		return nil, nil
	}
	if strings.ToUpper(value) == "TRUE" {
		return true, nil
	}
	if strings.ToUpper(value) == "FALSE" {
		return false, nil
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("extractValue(): invalid value: %s", value)
	}

	return val, nil
}
