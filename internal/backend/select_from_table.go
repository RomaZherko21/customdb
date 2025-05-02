package backend

import (
	"bytes"
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"encoding/binary"
	"fmt"
	"slices"
)

func (mb *memoryBackend) selectFromTable(statement *ast.SelectStatement) (*models.Table, error) {
	// table, err := mb.memoryStorage.Select(statement.From.Value)
	// if err != nil {
	// 	return nil, err
	// }

	table, err := mb.persistentStorage.Select(statement.From.Value)
	if err != nil {
		return nil, err
	}

	// WHERE
	if statement.Where != nil && len(statement.Where) > 0 {
		table.Rows = mb.filterRows(table.Columns, table.Rows, statement.Where)
	}
	//

	//// FILTERING
	if statement.SelectedColumns != nil && len(statement.SelectedColumns) != 0 {
		selectedColumnNames := []string{}
		for _, value := range statement.SelectedColumns {
			selectedColumnNames = append(selectedColumnNames, value.Literal.Value)
		}

		for i, row := range table.Rows {
			resultRow := []interface{}{}
			for i, column := range table.Columns {
				if slices.Contains(selectedColumnNames, column.Name) {
					resultRow = append(resultRow, row[i])
				}
			}

			table.Rows[i] = resultRow
		}
		newColumns := []models.Column{}
		for _, column := range table.Columns {
			if slices.Contains(selectedColumnNames, column.Name) {
				newColumns = append(newColumns, column)
			}
		}
		table.Columns = newColumns
	}

	////

	rows := [][]models.Cell{}

	for _, row := range table.Rows {
		newRow := []models.Cell{}
		for i, cell := range row {
			column := table.Columns[i]
			var memoryCell MemoryCell

			if column.Type == models.IntType {
				buf := new(bytes.Buffer)
				err := binary.Write(buf, binary.BigEndian, int32(cell.(float64)))
				if err != nil {
					return nil, fmt.Errorf("failed to convert int: %w", err)
				}
				memoryCell = MemoryCell(buf.Bytes())
			}

			if column.Type == models.TextType {
				memoryCell = MemoryCell(cell.(string))
			}

			newRow = append(newRow, memoryCell)
		}
		rows = append(rows, newRow)
	}

	return &models.Table{
		Name:    table.Name,
		Columns: table.Columns,
		Rows:    rows,
	}, nil
}

func (mb *memoryBackend) filterRows(columns []models.Column, rows [][]interface{}, whereClause []*ast.WhereClause) [][]interface{} {
	filteredRows := [][]interface{}{}
	for _, row := range rows {
		if mb.filterRow(columns, row, whereClause) {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

func (mb *memoryBackend) filterRow(columns []models.Column, row []interface{}, whereClause []*ast.WhereClause) bool {
	if len(whereClause) == 0 {
		return true
	}

	// Вычисляем все условия
	conditions := make([]bool, 0)
	logicalOps := make([]lex.Token, 0)

	for _, clause := range whereClause {
		if clause.Operator.Kind == lex.LogicalOperatorToken {
			logicalOps = append(logicalOps, clause.Operator)
			continue
		}

		leftVal, err := mb.getValueFromExpression(columns, row, clause.Left)
		if err != nil {
			return false
		}

		rightVal, err := mb.getValueFromExpression(columns, row, clause.Right)
		if err != nil {
			return false
		}

		condition := mb.evaluateCondition(leftVal, rightVal, clause.Operator.Value)
		conditions = append(conditions, condition)
	}

	// Применяем логические операторы с учетом приоритета
	result := conditions[0]
	for i, op := range logicalOps {
		if i >= len(conditions)-1 {
			break
		}

		switch op.Value {
		case string(lex.AndOperator):
			result = result && conditions[i+1]
		case string(lex.OrOperator):
			result = result || conditions[i+1]
		}
	}

	return result
}

func (mb *memoryBackend) getValueFromExpression(columns []models.Column, row []interface{}, expr *ast.Expression) (interface{}, error) {
	if expr.Literal.Kind == lex.IdentifierToken {
		// Ищем индекс колонки
		for i, column := range columns {
			if column.Name == expr.Literal.Value {
				return row[i], nil
			}
		}
		return nil, fmt.Errorf("column not found: %s", expr.Literal.Value)
	}
	return expr.Literal.Value, nil
}

func (mb *memoryBackend) evaluateCondition(left, right interface{}, operator string) bool {
	switch operator {
	case string(lex.EqualOperator):
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	case string(lex.NotEqualOperator):
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right)
	case string(lex.GreaterThanOperator):
		return fmt.Sprintf("%v", left) > fmt.Sprintf("%v", right)
	case string(lex.LessThanOperator):
		return fmt.Sprintf("%v", left) < fmt.Sprintf("%v", right)
	default:
		return false
	}
}
