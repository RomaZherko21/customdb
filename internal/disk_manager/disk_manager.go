package disk_manager

func CreateTable(filename string) {
	createMetaFile(&MetaFile{
		Name: filename,
		Columns: []Column{
			{Name: "id", Type: ColumnTypeInt},
			{Name: "name", Type: ColumnTypeString},
			{Name: "is_admin", Type: ColumnTypeBoolean},
		},
	})
	createDataFile(filename)
}
