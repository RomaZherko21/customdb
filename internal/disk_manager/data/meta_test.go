package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeAndDeserializeMetaData(t *testing.T) {
	t.Run("сериализация и десериализация простой таблицы", func(t *testing.T) {
		// Подготовка тестовых данных
		metaData := &MetaData{
			Name:      "test_table",
			PageCount: 1,
			Columns: []Column{
				{
					Name:       "id",
					Type:       TypeInt32,
					IsNullable: false,
				},
				{
					Name:       "name",
					Type:       TypeText,
					IsNullable: true,
				},
			},
		}

		// Сериализация
		serialized := serializeMetaData(metaData)
		assert.NotEmpty(t, serialized, "сериализованные данные не должны быть пустыми")

		// Десериализация
		deserialized, offset := deserializeMetaData(serialized)
		assert.NotNil(t, deserialized, "десериализованные данные не должны быть nil")
		assert.Greater(t, offset, 0, "offset должен быть больше 0")

		// Проверка полей
		assert.Equal(t, metaData.Name, deserialized.Name, "имя таблицы должно совпадать")
		assert.Equal(t, metaData.PageCount, deserialized.PageCount, "количество страниц должно совпадать")
		assert.Len(t, deserialized.Columns, len(metaData.Columns), "количество колонок должно совпадать")

		// Проверка первой колонки
		assert.Equal(t, "id", deserialized.Columns[0].Name, "имя первой колонки должно совпадать")
		assert.Equal(t, TypeInt32, deserialized.Columns[0].Type, "тип первой колонки должен совпадать")
		assert.False(t, deserialized.Columns[0].IsNullable, "первая колонка не должна быть nullable")

		// Проверка второй колонки
		assert.Equal(t, "name", deserialized.Columns[1].Name, "имя второй колонки должно совпадать")
		assert.Equal(t, TypeText, deserialized.Columns[1].Type, "тип второй колонки должен совпадать")
		assert.True(t, deserialized.Columns[1].IsNullable, "вторая колонка должна быть nullable")
	})

	t.Run("сериализация и десериализация таблицы без колонок", func(t *testing.T) {
		metaData := &MetaData{
			Name:      "empty_table",
			PageCount: 1,
			Columns:   []Column{},
		}

		serialized := serializeMetaData(metaData)
		assert.NotEmpty(t, serialized, "сериализованные данные не должны быть пустыми")

		deserialized, offset := deserializeMetaData(serialized)
		assert.NotNil(t, deserialized, "десериализованные данные не должны быть nil")
		assert.Greater(t, offset, 0, "offset должен быть больше 0")

		assert.Equal(t, metaData.Name, deserialized.Name, "имя таблицы должно совпадать")
		assert.Equal(t, metaData.PageCount, deserialized.PageCount, "количество страниц должно совпадать")
		assert.Empty(t, deserialized.Columns, "список колонок должен быть пустым")
	})

	t.Run("сериализация и десериализация таблицы с разными типами данных", func(t *testing.T) {
		metaData := &MetaData{
			Name:      "complex_table",
			PageCount: 5,
			Columns: []Column{
				{
					Name:       "int_col",
					Type:       TypeInt32,
					IsNullable: false,
				},
				{
					Name:       "string_col",
					Type:       TypeText,
					IsNullable: true,
				},
				{
					Name:       "uint_col",
					Type:       TypeUint32,
					IsNullable: false,
				},
				{
					Name:       "bool_col",
					Type:       TypeBoolean,
					IsNullable: true,
				},
			},
		}

		serialized := serializeMetaData(metaData)
		assert.NotEmpty(t, serialized, "сериализованные данные не должны быть пустыми")

		deserialized, offset := deserializeMetaData(serialized)
		assert.NotNil(t, deserialized, "десериализованные данные не должны быть nil")
		assert.Greater(t, offset, 0, "offset должен быть больше 0")

		assert.Equal(t, metaData.Name, deserialized.Name, "имя таблицы должно совпадать")
		assert.Equal(t, metaData.PageCount, deserialized.PageCount, "количество страниц должно совпадать")
		assert.Len(t, deserialized.Columns, len(metaData.Columns), "количество колонок должно совпадать")

		// Проверка всех колонок
		assert.Equal(t, "int_col", deserialized.Columns[0].Name)
		assert.Equal(t, TypeInt32, deserialized.Columns[0].Type)
		assert.False(t, deserialized.Columns[0].IsNullable)

		assert.Equal(t, "string_col", deserialized.Columns[1].Name)
		assert.Equal(t, TypeText, deserialized.Columns[1].Type)
		assert.True(t, deserialized.Columns[1].IsNullable)

		assert.Equal(t, "uint_col", deserialized.Columns[2].Name)
		assert.Equal(t, TypeUint32, deserialized.Columns[2].Type)
		assert.False(t, deserialized.Columns[2].IsNullable)

		assert.Equal(t, "bool_col", deserialized.Columns[3].Name)
		assert.Equal(t, TypeBoolean, deserialized.Columns[3].Type)
		assert.True(t, deserialized.Columns[3].IsNullable)
	})
}
