package data_structure

import (
	"fmt"
	"reflect"
	"testing"
)

type Int int

func (m Int) Equals(x Object) bool {
	return m == x
}

func (m Int) Hashcode() int64 {
	return int64(m)
}

func TestLinkedListStack(t *testing.T) {
	var stack Stack = &LinkedList{
		check: TypeCheck{t: reflect.TypeOf(Int(0))},
	}

	stack.Push(Int(1))
	stack.Push(Int(2))
	fmt.Println(stack.Pop().(*LinkedNode).Value())
	fmt.Println(stack.Pop().(*LinkedNode).Value())
}

func TestLinkedListQueue(t *testing.T) {
	var queue Queue = &LinkedList{
		check: TypeCheck{t: reflect.TypeOf(Int(0))},
	}

	queue.Push(Int(1))
	queue.Push(Int(2))
	fmt.Println(queue.Poll().(*LinkedNode).Value())
	fmt.Println(queue.Poll().(*LinkedNode).Value())
}

func TestLinearListStack(t *testing.T) {
	var stack Stack = &LinearList{
		check: TypeCheck{t: reflect.TypeOf(Int(0))},
	}

	stack.Push(Int(1))
	stack.Push(Int(2))
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
}

func TestLinearListQueue(t *testing.T) {
	var queue Queue = &LinearList{
		check: TypeCheck{t: reflect.TypeOf(Int(0))},
	}

	queue.Push(Int(1))
	queue.Push(Int(2))
	fmt.Println(queue.Poll())
	fmt.Println(queue.Poll())
}
