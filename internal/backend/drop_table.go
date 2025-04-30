package backend

import "custom-database/internal/parser/ast"

func (mb *memoryBackend) dropTable(statement *ast.DropTableStatement) error {
	return mb.memoryStorage.DropTable(statement.Table.Value)
}
