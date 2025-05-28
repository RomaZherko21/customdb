package b_plus_tree

const (
	// Минимальная степень B+ дерева
	minDegree = 2
	// Максимальное количество ключей в узле
	maxKeys = 2*minDegree - 1
	// Максимальное количество дочерних узлов
	maxChildren = 2 * minDegree
)

type BPlusTree struct {
	root *Node
}

type Node struct {
	keys     []int
	values   []interface{} // только для листовых узлов
	children []*Node
	isLeaf   bool
	next     *Node // указатель на следующий листовой узел
}

func NewBPlusTree() *BPlusTree {
	return &BPlusTree{
		root: nil,
	}
}

// Создание нового узла
func newNode() *Node {
	return &Node{
		keys:     make([]int, 0, maxKeys),
		values:   make([]interface{}, 0, maxKeys),
		children: make([]*Node, 0, maxChildren),
		isLeaf:   true,
	}
}

// Insert добавляет новую пару ключ-значение в B+ дерево
func (t *BPlusTree) Insert(key int, value interface{}) {
	if t.root == nil {
		t.root = newNode()
		t.root.keys = append(t.root.keys, key)
		t.root.values = append(t.root.values, value)
		return
	}

	if len(t.root.keys) == maxKeys {
		newRoot := newNode()
		newRoot.isLeaf = false
		newRoot.children = append(newRoot.children, t.root)
		t.root = newRoot
		t.splitChild(newRoot, 0)
	}
	t.insertNonFull(t.root, key, value)
}

// insertNonFull вставляет ключ в неполный узел
func (t *BPlusTree) insertNonFull(node *Node, key int, value interface{}) {
	if node.isLeaf {
		// Найти позицию для вставки
		i := 0
		for i < len(node.keys) && node.keys[i] < key {
			i++
		}
		// Если ключ уже есть — обновить значение
		if i < len(node.keys) && node.keys[i] == key {
			node.values[i] = value
			return
		}
		// Вставить ключ и значение
		node.keys = append(node.keys, 0)
		node.values = append(node.values, nil)
		copy(node.keys[i+1:], node.keys[i:])
		copy(node.values[i+1:], node.values[i:])
		node.keys[i] = key
		node.values[i] = value
		return
	}
	// Внутренний узел: найти потомка
	i := 0
	for i < len(node.keys) && key >= node.keys[i] {
		i++
	}
	// Если потомок переполнен — разделить и пересчитать индекс
	if len(node.children[i].keys) == maxKeys {
		t.splitChild(node, i)
		if key >= node.keys[i] {
			i++
		}
	}
	t.insertNonFull(node.children[i], key, value)
}

// splitChild разделяет заполненный дочерний узел
func (t *BPlusTree) splitChild(parent *Node, childIndex int) {
	child := parent.children[childIndex]
	newNode := newNode()
	newNode.isLeaf = child.isLeaf

	mid := len(child.keys) / 2

	if child.isLeaf {
		// Новый лист: копируем вторую половину ключей и значений
		newNode.keys = append(newNode.keys, child.keys[mid:]...)
		newNode.values = append(newNode.values, child.values[mid:]...)
		child.keys = child.keys[:mid]
		child.values = child.values[:mid]
		// Связный список листьев
		newNode.next = child.next
		child.next = newNode
		// В родителя — первый ключ нового листа
		parent.keys = append(parent.keys, 0)
		copy(parent.keys[childIndex+1:], parent.keys[childIndex:])
		parent.keys[childIndex] = newNode.keys[0]
		// Для листа у newNode нет детей
		if len(parent.children) == 0 {
			parent.children = []*Node{child, newNode}
		} else {
			parent.children = append(parent.children, nil)
			copy(parent.children[childIndex+2:], parent.children[childIndex+1:])
			parent.children[childIndex+1] = newNode
		}
	} else {
		// Внутренний: средний ключ поднимается, дети делятся
		midKey := child.keys[mid]
		newNode.keys = append(newNode.keys, child.keys[mid+1:]...)
		newNode.children = append(newNode.children, child.children[mid+1:]...)
		child.keys = child.keys[:mid] // средний ключ не остается в child
		child.children = child.children[:mid+1]
		// В родителя — средний ключ
		parent.keys = append(parent.keys, 0)
		copy(parent.keys[childIndex+1:], parent.keys[childIndex:])
		parent.keys[childIndex] = midKey
		parent.children = append(parent.children, nil)
		copy(parent.children[childIndex+2:], parent.children[childIndex+1:])
		parent.children[childIndex+1] = newNode
	}
}

// Search ищет значение по ключу
func (t *BPlusTree) Search(key int) (interface{}, bool) {
	if t.root == nil {
		return nil, false
	}
	return t.searchNode(t.root, key)
}

// searchNode рекурсивно ищет ключ в узле
func (t *BPlusTree) searchNode(node *Node, key int) (interface{}, bool) {
	if node.isLeaf {
		for i, k := range node.keys {
			if k == key {
				return node.values[i], true
			}
		}
		return nil, false
	}
	i := 0
	for i < len(node.keys) && key >= node.keys[i] {
		i++
	}
	return t.searchNode(node.children[i], key)
}

// Delete удаляет ключ из B+ дерева
func (t *BPlusTree) Delete(key int) bool {
	if t.root == nil {
		return false
	}

	// Сначала проверяем, существует ли ключ
	_, exists := t.Search(key)
	if !exists {
		return false
	}

	success := t.deleteKey(t.root, key)

	// Если корень пуст и не является листом, обновляем корень
	if len(t.root.keys) == 0 && !t.root.isLeaf {
		t.root = t.root.children[0]
	}

	return success
}

// fill заполняет дочерний узел, который имеет меньше t-1 ключей
func (t *BPlusTree) fill(parent *Node, idx int) {
	if idx != 0 && len(parent.children[idx-1].keys) >= minDegree {
		t.borrowFromPrev(parent, idx)
	} else if idx != len(parent.children)-1 && len(parent.children[idx+1].keys) >= minDegree {
		t.borrowFromNext(parent, idx)
	} else {
		if idx != len(parent.children)-1 {
			t.merge(parent, idx)
		} else {
			t.merge(parent, idx-1)
		}
	}
}

// Обновляем метод deleteKey, добавляя вызов fill
func (t *BPlusTree) deleteKey(node *Node, key int) bool {
	if node.isLeaf {
		idx := -1
		for i, k := range node.keys {
			if k == key {
				idx = i
				break
			}
		}
		if idx == -1 {
			return false
		}

		// Удаляем ключ и значение
		copy(node.keys[idx:], node.keys[idx+1:])
		copy(node.values[idx:], node.values[idx+1:])
		node.keys = node.keys[:len(node.keys)-1]
		node.values = node.values[:len(node.values)-1]

		return true
	}

	// Во внутреннем узле ищем подходящего потомка
	i := 0
	for i < len(node.keys) && key >= node.keys[i] {
		i++
	}
	if i > 0 {
		i--
	}

	// Проверяем, нужно ли заполнить дочерний узел
	if len(node.children[i].keys) < minDegree {
		t.fill(node, i)
		// После fill индекс может измениться
		if i > 0 && key < node.keys[i-1] {
			i--
		}
	}

	// Если нашли ключ во внутреннем узле, заменяем его на преемника
	if i < len(node.keys) && node.keys[i] == key {
		successor := t.getSuccessor(node.children[i+1])
		node.keys[i] = successor
		return t.deleteKey(node.children[i+1], successor)
	}

	return t.deleteKey(node.children[i], key)
}

// getSuccessor находит преемника ключа
func (t *BPlusTree) getSuccessor(node *Node) int {
	for !node.isLeaf {
		node = node.children[0]
	}
	return node.keys[0]
}

// borrowFromPrev заимствует ключ у предыдущего дочернего узла
func (t *BPlusTree) borrowFromPrev(parent *Node, idx int) {
	child := parent.children[idx]
	sibling := parent.children[idx-1]

	if child.isLeaf {
		// Для листового узла
		child.keys = append([]int{sibling.keys[len(sibling.keys)-1]}, child.keys...)
		child.values = append([]interface{}{sibling.values[len(sibling.values)-1]}, child.values...)
		sibling.keys = sibling.keys[:len(sibling.keys)-1]
		sibling.values = sibling.values[:len(sibling.values)-1]
		parent.keys[idx-1] = child.keys[0]
	} else {
		// Для внутреннего узла
		child.keys = append([]int{parent.keys[idx-1]}, child.keys...)
		parent.keys[idx-1] = sibling.keys[len(sibling.keys)-1]
		sibling.keys = sibling.keys[:len(sibling.keys)-1]

		if !child.isLeaf {
			child.children = append([]*Node{sibling.children[len(sibling.children)-1]}, child.children...)
			sibling.children = sibling.children[:len(sibling.children)-1]
		}
	}
}

// borrowFromNext заимствует ключ у следующего дочернего узла
func (t *BPlusTree) borrowFromNext(parent *Node, idx int) {
	child := parent.children[idx]
	sibling := parent.children[idx+1]

	if child.isLeaf {
		// Для листового узла
		child.keys = append(child.keys, sibling.keys[0])
		child.values = append(child.values, sibling.values[0])
		sibling.keys = sibling.keys[1:]
		sibling.values = sibling.values[1:]
		parent.keys[idx] = sibling.keys[0]
	} else {
		// Для внутреннего узла
		child.keys = append(child.keys, parent.keys[idx])
		parent.keys[idx] = sibling.keys[0]
		sibling.keys = sibling.keys[1:]

		if !child.isLeaf {
			child.children = append(child.children, sibling.children[0])
			sibling.children = sibling.children[1:]
		}
	}
}

// merge объединяет два дочерних узла node[idx] и node[idx+1]
func (t *BPlusTree) merge(parent *Node, idx int) {
	child := parent.children[idx]
	sibling := parent.children[idx+1]

	if child.isLeaf {
		// Для листовых узлов
		child.keys = append(child.keys, sibling.keys...)
		child.values = append(child.values, sibling.values...)
		child.next = sibling.next
	} else {
		// Для внутренних узлов
		child.keys = append(child.keys, parent.keys[idx])
		child.keys = append(child.keys, sibling.keys...)
		child.children = append(child.children, sibling.children...)
	}

	// Удаляем ключ из родителя и sibling из его детей
	copy(parent.keys[idx:], parent.keys[idx+1:])
	copy(parent.children[idx+1:], parent.children[idx+2:])
	parent.keys = parent.keys[:len(parent.keys)-1]
	parent.children = parent.children[:len(parent.children)-1]

	// Освобождаем память
	sibling.keys = nil
	sibling.values = nil
	sibling.children = nil
}
