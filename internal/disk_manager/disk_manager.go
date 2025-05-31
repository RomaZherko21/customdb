package disk_manager

import (
	"custom-database/internal/disk_manager/data"
	"custom-database/internal/disk_manager/meta"
)

func CreateTable(filename string, columns []meta.Column) {
	meta.CreateMetaFile(&meta.MetaFile{
		Name:    filename,
		Columns: columns,
	})
	data.CreateDataFile(filename)
}
