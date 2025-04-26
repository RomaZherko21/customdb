package lexer

import (
	"custom-database/internal/model"
	"fmt"
	"regexp"
	"strings"
)

// INSERT INTO films (code, title, did, date_prod, kind)
// VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');
func ParseInsertIntoCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return model.Table{}, fmt.Errorf("parseInsertIntoCommand(): not enough arguments")
	}

	tableName := parts[2]

	re := regexp.MustCompile(fmt.Sprintf(`%s \((.*)\)`, tableName))
	columnsSearch := re.FindStringSubmatch(input)
	if len(columnsSearch) < 1 {
		return model.Table{}, fmt.Errorf("parseInsertIntoCommand(): not found any columns")
	}

	columns := strings.Split(columnsSearch[1], ",")

	re = regexp.MustCompile(`VALUES \((.*)\)`)
	valuesSearch := re.FindStringSubmatch(input)
	if len(valuesSearch) < 1 {
		return model.Table{}, fmt.Errorf("parseInsertIntoCommand(): not found any values")
	}

	values := strings.Split(valuesSearch[1], ",")

	result := model.Table{
		TableName: tableName,
		Columns:   []model.Column{},
		Rows:      []model.Row{},
	}

	for _, column := range columns {
		column = strings.TrimSpace(column)
		column = strings.Trim(column, "()")
		column = strings.TrimSpace(column)

		result.Columns = append(result.Columns, model.Column{
			Name: column,
		})
	}

	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "()")
		value = strings.TrimSpace(value)

		result.Rows = append(result.Rows, model.Row{
			Values: append([]interface{}{}, value),
		})
	}

	fmt.Println("result", result)

	return result, nil
}
