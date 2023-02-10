package data_structure

import "reflect"

var (
	_ Stack       = (*LinearList)(nil)
	_ Queue       = (*LinearList)(nil)
	_ Collections = (*LinearList)(nil)
)

type LinearList struct {
	check TypeCheck
	array  []Object
}

func (list *LinearList) Empty() bool {
	return len(list.array) == 0
}

func (list *LinearList) Length() int {
	return len(list.array)
}

func (list *LinearList) Add(val Object) error {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return err
	}
	list.array = append(list.array, val)
	return nil
}

func (list *LinearList) Remove(val Object) error {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return err
	}

	if list.Length() == 0 {
		return nil
	}

	for i, v := range list.array {
		if v.Equals(val) {
			list.array = append(list.array[:i], list.array[i+1])
			return nil
		}
	}
	return nil
}

func (list *LinearList) Contains(val Object) bool {
	if err := list.check.Check(reflect.TypeOf(val)); err != nil {
		return false
	}

	if list.Length() == 0 {
		return false
	}

	for _, v := range list.array {
		if v.Equals(val) {
			return true
		}
	}
	return false
}

func (list *LinearList) Poll() interface{} {
	if list.Empty() {
		return nil
	}
	node := list.array[0]
	list.array = list.array[1:]
	return node
}

func (list *LinearList) Pop() interface{} {
	if list.Empty() {
		return nil
	}
	node := list.array[len(list.array)-1]
	list.array = list.array[:len(list.array)-1]
	return node
}

func (list *LinearList) Peek() interface{} {
	if list.Empty() {
		return nil
	}
	return list.array[len(list.array)-1]
}

func (list *LinearList) Element() interface{} {
	if list.Empty() {
		return nil
	}
	return list.array[0]
}

func (list *LinearList) Push(val Object) error {
	return list.Add(val)
}
