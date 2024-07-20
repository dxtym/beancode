package beancode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// better error handling
type EncodeError struct {
	Type string
	Err error
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("beancode: %s - %s", e.Type, e.Err.Error())
}

type Encoder struct {
	w   io.Writer
	buf bytes.Buffer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.IsZero() || !rv.IsValid() {
		return &EncodeError{
			Type: "zero or nil value", 
			Err: fmt.Errorf("cannot encode zero or nil value"),
		}
	}

	e.buf.Reset()
	switch rv.Kind() {
	case reflect.Int:
		e.encodeInt(rv)
	case reflect.String:
		e.encodeStr(rv)
	case reflect.Slice:
		e.encodeList(rv)
	case reflect.Map:
		e.encodeDict(rv)
	case reflect.Struct:
		e.encodeStruct(rv)
	default:
		return &EncodeError{
			Type: "unsupported type",
			Err: fmt.Errorf("cannot encode type %s", rv.Type().Name()),
		}
	}

	_, err := e.w.Write(e.buf.Bytes())
	if err != nil {
		return &EncodeError{
			Type: "internal error",
			Err: fmt.Errorf("cannot write to writer: %w", err),
		}
	}

	return nil
}

func (e *Encoder) encodeInt(rv reflect.Value) {
	e.buf.WriteRune('i')
	
	if v, ok := rv.Interface().(int); ok {
		e.buf.WriteString(strconv.Itoa(v))
	}
	
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeStr(rv reflect.Value) {
	e.buf.WriteString(strconv.Itoa(rv.Len()))
	e.buf.WriteRune(':')
	e.buf.WriteString(rv.String())
}

func (e *Encoder) encodeList(rv reflect.Value) {
	e.w.Write([]byte{'l'})
	
	for i := 0; i < rv.Len(); i++ {
		e.Encode(rv.Index(i).Interface())
	}

	e.buf.Reset()
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeDict(rv reflect.Value) {
	e.w.Write([]byte{'d'}) 

	iter := rv.MapRange() // to keep the order of keys
	for iter.Next() {
		e.Encode(iter.Key().Interface())
		e.Encode(iter.Value().Interface())
	}
	
	e.buf.Reset()
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeStruct(rv reflect.Value) {
	e.w.Write([]byte{'d'})
	
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		e.Encode(f.Tag.Get("beancode"))
		e.Encode(rv.FieldByName(f.Name).Interface())
	}
	
	e.buf.Reset()
	e.buf.WriteRune('e')
}
