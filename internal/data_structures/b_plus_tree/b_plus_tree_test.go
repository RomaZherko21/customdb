package b_plus_tree

import (
	"testing"
)

// TestNewBPlusTree проверяет создание нового дерева
func TestNewBPlusTree(t *testing.T) {
	tree := NewBPlusTree()
	if tree.root != nil {
		t.Error("Корень нового дерева должен быть nil")
	}
}

// TestInsertAndSearch проверяет вставку и поиск элементов
func TestInsertAndSearch(t *testing.T) {
	tree := NewBPlusTree()
	testData := []struct {
		key   int
		value string
	}{
		{10, "десять"},
		{20, "двадцать"},
		{5, "пять"},
		{15, "пятнадцать"},
		{25, "двадцать пять"},
	}

	// Вставляем элементы
	for _, data := range testData {
		tree.Insert(data.key, data.value)
	}

	// Проверяем поиск существующих элементов
	for _, data := range testData {
		value, found := tree.Search(data.key)
		if !found {
			t.Errorf("Ключ %d не найден после вставки", data.key)
		}
		if value != data.value {
			t.Errorf("Для ключа %d ожидалось значение %s, получено %v", data.key, data.value, value)
		}
	}

	// Проверяем поиск несуществующих элементов
	nonExistentKeys := []int{1, 7, 12, 17, 30}
	for _, key := range nonExistentKeys {
		if _, found := tree.Search(key); found {
			t.Errorf("Ключ %d не должен существовать", key)
		}
	}
}

// TestDelete проверяет удаление элементов
func TestDelete(t *testing.T) {
	tree := NewBPlusTree()
	testData := []struct {
		key   int
		value string
	}{
		{10, "десять"},
		{20, "двадцать"},
		{5, "пять"},
		{15, "пятнадцать"},
		{25, "двадцать пять"},
	}

	// Вставляем элементы
	for _, data := range testData {
		tree.Insert(data.key, data.value)
	}

	// Удаляем элементы и проверяем
	for i, data := range testData {
		// Проверяем, что элемент существует перед удалением
		if _, found := tree.Search(data.key); !found {
			t.Errorf("Ключ %d не найден перед удалением", data.key)
			continue
		}

		// Удаляем элемент
		if !tree.Delete(data.key) {
			t.Errorf("Не удалось удалить ключ %d", data.key)
			continue
		}

		// Проверяем, что элемент удален
		if _, found := tree.Search(data.key); found {
			t.Errorf("Ключ %d все еще существует после удаления", data.key)
		}

		// Проверяем, что остальные элементы на месте
		for j := i + 1; j < len(testData); j++ {
			remainingData := testData[j]
			if value, found := tree.Search(remainingData.key); !found {
				t.Errorf("Ключ %d не найден после удаления %d", remainingData.key, data.key)
			} else if value != remainingData.value {
				t.Errorf("Для ключа %d ожидалось значение %s, получено %v после удаления %d",
					remainingData.key, remainingData.value, value, data.key)
			}
		}
	}
}

// TestLeafNodeLinking проверяет связывание листовых узлов
func TestLeafNodeLinking(t *testing.T) {
	tree := NewBPlusTree()
	keys := []int{10, 20, 5, 15, 25}

	// Вставляем элементы
	for _, key := range keys {
		tree.Insert(key, key*2)
	}

	// Находим первый листовой узел
	var firstLeaf *Node
	node := tree.root
	for !node.isLeaf {
		node = node.children[0]
	}
	firstLeaf = node

	// Проверяем связывание листовых узлов
	current := firstLeaf
	expectedKeys := []int{5, 10, 15, 20, 25}
	i := 0

	for current != nil && i < len(expectedKeys) {
		// Проверяем, что ключи в узле отсортированы
		for j := 0; j < len(current.keys)-1; j++ {
			if current.keys[j] > current.keys[j+1] {
				t.Errorf("Ключи в листовом узле не отсортированы: %v", current.keys)
			}
		}

		// Проверяем, что все ключи в узле меньше следующего ключа
		if current.next != nil && len(current.next.keys) > 0 {
			if current.keys[len(current.keys)-1] >= current.next.keys[0] {
				t.Errorf("Нарушен порядок ключей между листовыми узлами: %v -> %v",
					current.keys, current.next.keys)
			}
		}

		// Проверяем значения
		for _, key := range current.keys {
			if key != expectedKeys[i] {
				t.Errorf("Ожидался ключ %d, получен %d", expectedKeys[i], key)
			}
			i++
		}

		current = current.next
	}

	if i != len(expectedKeys) {
		t.Errorf("Не все ключи найдены в листовых узлах. Ожидалось %d, получено %d",
			len(expectedKeys), i)
	}
}

// TestInsertWithDuplicates проверяет вставку дубликатов
func TestInsertWithDuplicates(t *testing.T) {
	tree := NewBPlusTree()

	// Вставляем ключ с первым значением
	tree.Insert(10, "первое")

	// Вставляем тот же ключ с другим значением
	tree.Insert(10, "второе")

	// Проверяем, что сохранилось последнее значение
	if value, found := tree.Search(10); !found {
		t.Error("Ключ 10 не найден после вставки дубликата")
	} else if value != "второе" {
		t.Errorf("Ожидалось значение 'второе', получено %v", value)
	}
}

// TestEmptyTreeOperations проверяет операции с пустым деревом
func TestEmptyTreeOperations(t *testing.T) {
	tree := NewBPlusTree()

	// Проверяем поиск в пустом дереве
	if _, found := tree.Search(1); found {
		t.Error("Поиск в пустом дереве должен возвращать false")
	}

	// Проверяем удаление из пустого дерева
	if tree.Delete(1) {
		t.Error("Удаление из пустого дерева должно возвращать false")
	}
}

// TestInsertAndDeleteSequence проверяет последовательность вставок и удалений
func TestInsertAndDeleteSequence(t *testing.T) {
	tree := NewBPlusTree()

	// Последовательность операций: вставка -> проверка -> удаление -> проверка
	operations := []struct {
		key   int
		value string
	}{
		{5, "пять"},
		{3, "три"},
		{7, "семь"},
		{1, "один"},
		{9, "девять"},
	}

	// Выполняем операции
	for i, op := range operations {
		// Вставляем элемент
		tree.Insert(op.key, op.value)

		// Проверяем, что элемент вставлен
		if value, found := tree.Search(op.key); !found {
			t.Errorf("Ключ %d не найден после вставки", op.key)
		} else if value != op.value {
			t.Errorf("Для ключа %d ожидалось значение %s, получено %v", op.key, op.value, value)
		}

		// Удаляем элемент
		if !tree.Delete(op.key) {
			t.Errorf("Не удалось удалить ключ %d", op.key)
		}

		// Проверяем, что элемент удален
		if _, found := tree.Search(op.key); found {
			t.Errorf("Ключ %d все еще существует после удаления", op.key)
		}

		// Проверяем, что предыдущие элементы удалены
		for j := 0; j < i; j++ {
			prevOp := operations[j]
			if _, found := tree.Search(prevOp.key); found {
				t.Errorf("Ключ %d все еще существует после удаления %d", prevOp.key, op.key)
			}
		}
	}
}

// TestTreeBalancing проверяет балансировку дерева
func TestTreeBalancing(t *testing.T) {
	tree := NewBPlusTree()

	// Вставляем элементы в порядке, который может вызвать перебалансировку
	keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 20, 14, 4029, 16, 190, 18, 19, 20, -30}
	for _, key := range keys {
		tree.Insert(key, key*2)
	}

	// Проверяем, что все элементы доступны
	for _, key := range keys {
		if value, found := tree.Search(key); !found {
			t.Errorf("Ключ %d не найден после вставки", key)
		} else if value != key*2 {
			t.Errorf("Для ключа %d ожидалось значение %d, получено %v", key, key*2, value)
		}
	}

	// Проверяем, что все листовые узлы находятся на одном уровне
	var checkLeafLevel func(*Node, int, *int) bool
	checkLeafLevel = func(node *Node, currentLevel int, leafLevel *int) bool {
		if node.isLeaf {
			if *leafLevel == -1 {
				*leafLevel = currentLevel
			} else if *leafLevel != currentLevel {
				return false
			}
			return true
		}
		for _, child := range node.children {
			if !checkLeafLevel(child, currentLevel+1, leafLevel) {
				return false
			}
		}
		return true
	}

	leafLevel := -1
	if !checkLeafLevel(tree.root, 0, &leafLevel) {
		t.Error("Листовые узлы находятся на разных уровнях")
	}
}
