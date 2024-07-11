package beancode

import (
	"fmt"
	"reflect"
)

type DecodeError struct {
	Msg string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("beancode: %v", e.Msg)
}

type OutOfBoundsError struct {
	Current int
}

func (e *OutOfBoundsError) Error() string {
	return fmt.Sprintf("beancode: index %v out of bounds", e.Current)
}

type InvalidTypeError struct {
	ExpectedType reflect.Type
	ActualType    reflect.Type
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("beancode: expected %v, got %v", e.ExpectedType, e.ActualType)
}
