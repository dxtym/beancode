package beancode

import (
	"bytes"
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
		return &DecodeError{err.Error()}
	}

	if len(data) == 0 {
		return &DecodeError{"empty input"} 
	}

	var val any
	val, err = d.decode(data)
	if err != nil {
		return &DecodeError{err.Error()}
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
		return &InvalidTypeError{reflect.PointerTo(rv.Type()), rv.Type()}
	}

	rv = rv.Elem()
	rg := reflect.ValueOf(got)
    switch rv.Kind() {
	case reflect.Int:
		val, ok := got.(int)
		if !ok {
			return &InvalidTypeError{reflect.TypeFor[int](), rg.Type()}
		}
		rv.Set(reflect.ValueOf(val))
	case reflect.String:
		val, ok := got.(string)
		if !ok {
			return &InvalidTypeError{reflect.TypeFor[string](), rg.Type()}
		}
		rv.Set(reflect.ValueOf(val))
    case reflect.Slice:
        val, ok := got.([]any)
        if !ok {
            return &InvalidTypeError{reflect.TypeFor[[]any](), rg.Type()}
        }
        rv.Set(reflect.ValueOf(val))
    case reflect.Map:
        val, ok := got.(map[string]any)
        if !ok {
            return &InvalidTypeError{reflect.TypeFor[map[string]any](), rg.Type()}
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
		return 0, &DecodeError{"invalid decode format"}
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
			return nil, &DecodeError{"index out of bounds"}
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
			return nil, &DecodeError{"index out of bounds"}
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
		return "", &DecodeError{"invalid decode format"}
	}

	colon += d.idx
	length, err := strconv.Atoi(string(data[d.idx:colon]))
	if err != nil {
		return "", err
	}

	start := colon + 1
	end := start + length
	if end > len(data) {
		return "", &DecodeError{"index out of bounds"}
	}

	d.idx = end
	return string(data[start:end]), nil
}
