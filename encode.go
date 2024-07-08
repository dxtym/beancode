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

func (e *Encoder) Encode(v any) (error) {
	e.buf.Reset()

	switch v := v.(type) {
	case int:
		e.encodeInt(v)
	case string:
		e.encodeStr(v)
	case []any:
		e.encodeList(v)
	case map[string]any:
		e.encodeDict(v)
	default:
		// TODO: struct
	}

	_, err := e.w.Write(e.buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeInt(v any) {
	e.buf.WriteRune('i')
	
	if v, ok := v.(int); ok {
		e.buf.WriteString(strconv.Itoa(v))
	}
	
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeStr(v any) {
	length := len(v.(string))

	e.buf.WriteString(strconv.Itoa(length))
	e.buf.WriteRune(':')
	e.buf.WriteString(v.(string))
}

func (e *Encoder) encodeList(v any) {
	e.buf.WriteRune('l')
	e.w.Write(e.buf.Bytes())

	s := reflect.ValueOf(v)
    for i := 0; i < s.Len(); i++ {
        e.Encode(s.Index(i).Interface())
    }

	e.buf.Reset()
	e.buf.WriteRune('e')
}

func (e *Encoder) encodeDict(v any) {
	e.buf.WriteRune('d')
	e.w.Write(e.buf.Bytes())
	
	s := reflect.ValueOf(v)
	for _, key := range s.MapKeys() {
		e.Encode(key.Interface())
		val := s.MapIndex(key)
		e.Encode(val.Interface())
	}
	
	e.buf.Reset()
	e.buf.WriteRune('e')
}
