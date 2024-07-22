package beancode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type DecodeError struct {
	Type string
	Err error
}

func (d *DecodeError) Error() string {
	return fmt.Sprintf("beancode: %s: %s", d.Type, d.Err.Error())
}

type Decoder struct {
	r   io.Reader
	idx int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(v any) error {
	data, err := io.ReadAll(d.r)
	if err != nil {
		return &DecodeError{
			Type: "internal error",
			Err: fmt.Errorf("cannot read from reader: %w", err),
		}
	}

	if len(data) == 0 {
		return &DecodeError{
			Type: "empty input",
			Err: fmt.Errorf("cannot decode empty input"),
		}
	}

	var val any
	val, err = d.decode(data)
	if err != nil {
		return err
	}

	return d.write(v, val)
}

func (d *Decoder) decode(data []byte) (any, error) {
	switch data[d.idx] {
	case 'i':
		return d.decodeInt(data)
	case 'l':
		return d.decodeList(data)
	case 'd':
		return d.decodeDict(data)
	default:
		return d.decodeStr(data)
	}
}


func (d *Decoder) decodeInt(data []byte) (int, error) {
	d.idx++
	end := bytes.IndexByte(data[d.idx:], 'e')
	if end == -1 {
		return 0, &DecodeError{
			Type: "invalid format",
			Err: fmt.Errorf("cannot decode invalid encoding"),
		}
	}

	end += d.idx
	val, err := strconv.Atoi(string(data[d.idx:end]))
	if err != nil {
		return 0, &DecodeError{
			Type: "internal error",
			Err: fmt.Errorf("cannot convert to str: %w", err),
		}
	}
	
	d.idx = end + 1
	return val, nil
}

func (d *Decoder) decodeStr(data []byte) (string, error) {
	colon := bytes.IndexByte(data[d.idx:], ':')
	if colon == -1 {
		return "", &DecodeError{
			Type: "invalid format",
			Err: fmt.Errorf("cannot decode invalid encoding"),
		}
	}

	colon += d.idx
	length, err := strconv.Atoi(string(data[d.idx:colon]))
	if err != nil {
		return "", &DecodeError{
			Type: "internal error",
			Err: fmt.Errorf("cannot convert to str: %w", err),
		}
	}

	start := colon + 1
	end := start + length
	if end > len(data) {
		return "", &DecodeError{
			Type: "invalid format",
			Err: fmt.Errorf("cannot decode invalid encoding"),
		}
	}

	d.idx = end
	return string(data[start:end]), nil
}

func (d *Decoder) decodeList(data []byte) ([]any, error) {
	buf := make([]any, 0)
	d.idx++

	for {
		if d.idx == len(data) {
			return nil, &DecodeError{
				Type: "invalid format",
				Err: fmt.Errorf("cannot decode invalid encoding"),
			}
		}

		if data[d.idx] == 'e' {
			d.idx++
			return buf, nil
		}

		val, err := d.decode(data)
		if err != nil {
			return nil, err
		}

		buf = append(buf, val)
	}
}

func (d *Decoder) decodeDict(data []byte) (map[string]any, error) {
	buf := make(map[string]any)
	d.idx++

	for {
		if d.idx == len(data) {
			return nil, &DecodeError{
				Type: "invalid format",
				Err: fmt.Errorf("cannot decode invalid encoding"),
			}
		}

		if data[d.idx] == 'e' {
			d.idx++
			return buf, nil
		}

		key, err := d.decodeStr(data)
		if err != nil {
			return nil, err
		}

		val, err := d.decode(data)
		if err != nil {
			return nil, err
		}

		buf[key] = val
	}
}


func (d *Decoder) write(v any, got any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &DecodeError{
			Type: "non-pointer or nil value",
			Err: fmt.Errorf("cannot decode to non-pointer or nil value"),
		}
	}

	rv = rv.Elem()
	rg := reflect.ValueOf(got)
	switch rv.Kind() {
	case reflect.Int:
		val, ok := got.(int)
		if !ok {
			return &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected int, got %v", rg.Type()),
			}
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.String:
		val, ok := got.(string)
		if !ok {
			return &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected string, got %v", rg.Type()),
			}
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.Slice:
		val, ok := got.([]any)
		if !ok {
			return &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected []any, got %v", rg.Type()),
			}
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.Map:
		val, ok := got.(map[string]any)
		if !ok {
			return &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected map[string]any, got %v", rg.Type()),
			}
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.Struct:
		// TODO: implement mapstructure myself 
		// TODO: add test afterwards
		if err := mapstructure.Decode(got, &v); err != nil {
			return &DecodeError{
				Type: "internal error",
				Err: fmt.Errorf("cannot decode to struct: %w", err),
			}
		}
	default:
		return  &EncodeError{
			Type: "unsupported type",
			Err: fmt.Errorf("cannot decode type %s", rv.Type().Name()),
		}
	}

	return nil
}