package data_structure

import "reflect"

type LinkedNode struct {
	val     Object
	next    *LinkedNode
	preview *LinkedNode
}

func (node *LinkedNode) Next() *LinkedNode {
	return node.next
}

func (node *LinkedNode) Preview() *LinkedNode {
	return node.preview
}

func (node *LinkedNode) Value() interface{} {
	return node.val
}

var (
	_ Stack       = (*LinkedList)(nil)
	_ Queue       = (*LinkedList)(nil)
	_ Collections = (*LinkedList)(nil)
)

type LinkedList struct {
	length int
	check  TypeCheck
	head   *LinkedNode
	tail   *LinkedNode
}

func (list *LinkedList) Empty() bool {
	return list.length == 0
}

func (list *LinkedList) Length() int {
	return list.length
}

func (list *LinkedList) Add(val Object) error {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return err
	}

	node := &LinkedNode{val: val}
	if list.head == nil {
		list.head = node
		list.tail = node
		return nil
	}

	node.next = list.head
	if list.head != nil {
		list.head.preview = node
	}

	if list.tail == nil {
		list.tail = node
	}
	list.head = node
	list.length++
	return nil
}

func (list *LinkedList) Remove(val Object) error {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return err
	}

	if list.Length() == 0 {
		return nil
	}

	q := list.head
	for ; q != nil; q = q.next {
		if q.val.Equals(val) {
			break
		}
	}

	// 没有找到
	if q == nil {
		return nil
	}

	// 删除头指针
	if q.preview == nil {
		list.head = q.next
		if q.next != nil {
			q.next.preview = nil
		}
	}

	// 删除尾指针
	if q.next == nil {
		list.tail = q.preview
		if q.preview != nil {
			q.preview.next = nil
		}
	}

	if q.preview != nil && q.next != nil {
		preview, next := q.preview, q.next
		preview.next = next
		next.preview = preview
	}

	q.next = nil
	q.preview = nil
	return nil
}

func (list *LinkedList) Contains(val Object) bool {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return false
	}

	if list.Length() == 0 {
		return false
	}

	q := list.head
	for ; q != nil; q = q.next {
		if q.val.Equals(val) {
			return true
		}
	}
	return false
}

func (list *LinkedList) Poll() interface{} {
	if list.Empty() {
		return nil
	}
	node := list.tail
	_ = list.Remove(node.val)
	return node
}

func (list *LinkedList) Pop() interface{} {
	if list.Empty() {
		return nil
	}
	node := list.head
	_ = list.Remove(node.val)
	return node
}

func (list *LinkedList) Peek() interface{} {
	if list.Empty() {
		return nil
	}
	return list.head
}

func (list *LinkedList) Element() interface{} {
	if list.Empty() {
		return nil
	}
	return list.tail
}

func (list *LinkedList) Push(val Object) error {
	return list.Add(val)
}
