package data_structure

import (
	"errors"
	"reflect"
)

type Object interface {
	Hashcode() int64
	Equals(object Object) bool
}

type Collections interface {
	Empty() bool
	Length() int
	Add(val Object) error
	Remove(val Object) error
	Contains(val Object) bool
}

type TypeCheck struct {
	t reflect.Type
}

var TypeNotMatch = errors.New("collections type is not match target type")

func (check TypeCheck) Check(target reflect.Type) error {
	if !(check.t.String() == target.String()) {
		return TypeNotMatch
	}
	return nil
}

type Stack interface {
	Pop() interface{}
	Peek() interface{}
	Push(val Object) error
}

type Queue interface {
	Poll() interface{}
	Element() interface{}
	Push(val Object) error
}
