package data

import (
	"custom-database/internal/disk_manager/meta"
	"fmt"
	"os"
	"path/filepath"
)

type fileConnection struct {
	lastPageID uint32
	file       *os.File
	columns    []meta.Column
}

const INITIAL_PAGE_ID = 1

func NewFileConnection(isNewFile bool, filename string, filePath string, columns []meta.Column) (*fileConnection, error) {
	filePath = filepath.Join(filePath, filename+".data")

	if isNewFile {
		fc := &fileConnection{
			lastPageID: INITIAL_PAGE_ID,
		}

		if err := fc.createDataFile(filePath); err != nil {
			return nil, fmt.Errorf("NewFileConnection(): CreateDataFile: %w", err)
		}
		return fc, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("NewFileConnection(): os.Open: %w", err)
	}

	return &fileConnection{
		lastPageID: INITIAL_PAGE_ID,
		file:       file,
		columns:    columns,
	}, nil
}

func (fc *fileConnection) Close() error {
	if fc.file == nil {
		return nil
	}

	return fc.file.Close()
}

func (fc *fileConnection) createDataFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("createDataFile(): table already exists: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("createDataFile(): os.Create: %w", err)
	}
	fc.file = file

	page := fc.newPage(INITIAL_PAGE_ID)
	data := fc.serializePage(page)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("createDataFile(): file.Write: %w", err)
	}

	return nil
}

func (fc *fileConnection) newPage(pageID uint32) *Page {
	return &Page{
		Header: PageHeader{
			PageId:   pageID,
			PageSize: PAGE_SIZE,
		},
		Slots: make([]PageSlot, MAX_SLOTS),
		Data:  make([]DataRow, 0),
	}
}
