package beancode

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Decoder struct {
	r     io.Reader
	idx int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(v any) error {
	data, err := io.ReadAll(d.r)
	if err != nil {
		return fmt.Errorf("beancode: %v", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("beancode: empty input")
	}

	var val any
	val, err = d.decode(data)
	if err != nil {
		return fmt.Errorf("beancode: %v", err)
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

func (d *Decoder) write(v any, got any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("beancode: expected pointer, got %v", rv.Type())
	}

	rv = rv.Elem()
	rg := reflect.ValueOf(got)
    switch rv.Kind() {
	case reflect.Int:
		val, ok := got.(int)
		if !ok {
			return fmt.Errorf("beancode: expected int, got %v", rg.Type())
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.String:
		val, ok := got.(string)
		if !ok {
			return fmt.Errorf("beancode: expected string, got %v", rg.Type())
		}
		rv.Set(reflect.ValueOf(val))
    case reflect.Slice:
        val, ok := got.([]any)
        if !ok {
            return fmt.Errorf("beancode: expected []any, got %v", rg.Type())
        }
        rv.Set(reflect.ValueOf(val))
    case reflect.Map:
        val, ok := got.(map[string]any)
        if !ok {
            return fmt.Errorf("beancode: expected map[string]any, got %v", rg.Type())
        }
        rv.Set(reflect.ValueOf(val))
	default:
		// TODO: struct
		
		
    }

    return nil
}

func (d *Decoder) decodeInt(data []byte) (int, error) {
	d.idx++
	end := bytes.IndexByte(data[d.idx:], 'e')
	if end == -1 {
		return 0, fmt.Errorf("beancode: invalid decode format")
	}

	end += d.idx
	val, err := strconv.Atoi(string(data[d.idx:end]))
	if err != nil {
		return 0, err
	}

	d.idx = end + 1
	return val, nil
}

func (d *Decoder) decodeList(data []byte) ([]any, error) {
	got := make([]any, 0)
	d.idx++

	for {
		if d.idx == len(data) {
			return nil, fmt.Errorf("beancode: index out of bounds")
		}
		if data[d.idx] == 'e' {
			d.idx++
			return got, nil
		}

		val, err := d.decode(data)
		if err != nil {
			return nil, err
		}

		got = append(got, val)
	}
}

func (d *Decoder) decodeDict(data []byte) (map[string]any, error) {
	got := make(map[string]any)
	d.idx++

	for {
		if d.idx == len(data) {
			return nil, fmt.Errorf("beancode: index out of bounds")
		}
		if data[d.idx] == 'e' {
			d.idx++
			return got, nil
		}

		key, err := d.decodeStr(data)
		if err != nil {
			return nil, err
		}

		val, err := d.decode(data)
		if err != nil {
			return nil, err
		}
		
		got[key] = val
	}
}

func (d *Decoder) decodeStr(data []byte) (string, error) {
	colon := bytes.IndexByte(data[d.idx:], ':')
	if colon == -1 {
		return "", fmt.Errorf("beancode: invalid decode format")
	}

	colon += d.idx
	length, err := strconv.Atoi(string(data[d.idx:colon]))
	if err != nil {
		return "", err
	}

	start := colon + 1
	end := start + length
	if end > len(data) {
		return "", fmt.Errorf("beancode: index out of bounds")
	}

	d.idx = end
	return string(data[start:end]), nil
}
