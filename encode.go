package beancode

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
)

type Encoder struct {
	w io.Writer
	buf bytes.Buffer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		return &EncodeError{"empty input"}
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
	default:
		e.encodeStruct(rv)
	}

	_, err := e.w.Write(e.buf.Bytes())
	if err != nil {
		return &EncodeError{err.Error()}
	}
	return nil
}

func (e *Encoder) encodeInt(rv reflect.Value) {
	e.buf.WriteRune('i')
	e.buf.WriteString(strconv.Itoa(int(rv.Int())))
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeStr(rv reflect.Value) {
	e.buf.WriteString(strconv.Itoa(rv.Len()))
	e.buf.WriteRune(':')
	e.buf.WriteString(rv.String())
}

func (e *Encoder) encodeList(rv reflect.Value) {
	e.buf.WriteRune('l')
	e.w.Write(e.buf.Bytes())

    for i := 0; i < rv.Len(); i++ {
        e.Encode(rv.Index(i).Interface())
    }

	e.buf.Reset()
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeDict(rv reflect.Value) {
	e.buf.WriteRune('d')
	e.w.Write(e.buf.Bytes())
	
	for _, key := range rv.MapKeys() {
		e.Encode(key.Interface())
		val := rv.MapIndex(key)
		e.Encode(val.Interface())
	}
	
	e.buf.Reset()
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeStruct(rv reflect.Value) {
	e.buf.WriteRune('d')
	e.w.Write(e.buf.Bytes())

	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		e.Encode(f.Tag.Get("beancode"))
		e.Encode(rv.FieldByName(f.Name).Interface())
	}

	e.buf.Reset()
	e.buf.WriteRune('e')
}