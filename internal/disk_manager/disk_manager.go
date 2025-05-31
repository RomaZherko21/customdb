package disk_manager

func CreateTable(filename string) {
	createMetaFile(&MetaFile{
		Name: filename,
		Columns: []Column{
			{Name: "id", Type: TypeInt32},
			{Name: "name", Type: TypeText},
			{Name: "is_admin", Type: TypeBoolean},
		},
	})
	createDataFile(filename)
}
