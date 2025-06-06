package data

import "errors"

func (ds *dataService) ParsePageHeader(pageID uint32) (*PageHeader, error) {
	start := uint32(ds.metaDataSpace) + PAGE_SIZE*(pageID-1)
	end := start + PAGE_HEADER_SIZE

	pageData, err := ds.ReadFileRange(start, end)
	if err != nil {
		return nil, err
	}

	result, err := deserializePageHeader(pageData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds *dataService) ParsePageHeaders() ([]*PageHeader, error) {
	result := make([]*PageHeader, 0)
	for i := 1; i <= int(ds.meta.PageCount); i++ {
		start := uint32(ds.metaDataSpace) + PAGE_SIZE*(uint32(i)-1)
		end := start + PAGE_HEADER_SIZE

		pageData, err := ds.ReadFileRange(start, end)
		if err != nil {
			return nil, err
		}

		pageHeader, err := deserializePageHeader(pageData)
		if err != nil {
			return nil, err
		}

		result = append(result, pageHeader)
	}

	return result, nil
}

func (ds *dataService) ParsePageSlots(pageID uint32) ([]PageSlot, error) {
	start := uint32(ds.metaDataSpace) + PAGE_SIZE*(pageID-1) + PAGE_HEADER_SIZE
	end := start + SLOTS_SPACE

	pageData, err := ds.ReadFileRange(start, end)
	if err != nil {
		return nil, err
	}

	result, err := deserializePageSlots(pageData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds *dataService) ParseDataRow(pageID uint32, slotID uint16) ([]DataCell, error) {
	pageSlots, err := ds.ParsePageSlots(pageID)
	if err != nil {
		return nil, err
	}

	slotOffsetStart := 0
	slotOffsetEnd := 0
	for _, slot := range pageSlots {
		if slot.SlotId == slotID {
			slotOffsetStart = int(slot.Offset)
			slotOffsetEnd = int(slot.Offset) + int(slot.RowSize)
			break
		}
	}

	if slotOffsetStart == 0 || slotOffsetEnd == 0 {
		return nil, errors.New("slot not found")
	}

	start := uint32(ds.metaDataSpace) + PAGE_SIZE*(pageID-1) + uint32(slotOffsetStart)
	end := start + uint32(slotOffsetEnd)

	pageData, err := ds.ReadFileRange(start, end)
	if err != nil {
		return nil, err
	}

	result, err := deserializeDataRow(pageData, ds.meta.Columns)
	if err != nil {
		return nil, err
	}

	return result, nil
}
