package ddl

import (
	"custom-database/internal/model"
	"fmt"
	"strings"
)

func ParseDropTableCommand(input string) (model.Table, error) {
	parts := strings.Split(input, " ")

	if len(parts) <= 2 {
		return model.Table{}, fmt.Errorf("ParseDropTableCommand(): not enough arguments")
	}

	if len(parts) > 3 {
		return model.Table{}, fmt.Errorf("ParseDropTableCommand(): too many arguments")
	}

	tableName := strings.Trim(parts[len(parts)-1], ";")

	return model.Table{
		TableName: tableName,
	}, nil
}
