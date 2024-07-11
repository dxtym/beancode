package beancode

import (
	"fmt"
	"reflect"
)

type DecodeError struct {
	Msg string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("beancode: decode error, %v", e.Msg)
}

type InvalidTypeError struct {
	ExpectedType reflect.Type
	ActualType    reflect.Type
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("beancode: expected %v, got %v", e.ExpectedType, e.ActualType)
}

type EncodeError struct {
	Msg string
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("beancode: encode error, %v", e.Msg)
}
