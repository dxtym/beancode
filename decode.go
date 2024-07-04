package beancode

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

type Decoder struct {
	r     io.Reader
	index int
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(v any) error {
	data, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("beancode: empty input")
	}

	var val any
	for d.index < len(data) {
		val, err = d.decode(data)
		if err != nil {
			return err
		}
	}

	return d.write(v, val)
}

func (d *Decoder) decode(data []byte) (any, error) {
	switch data[d.index] {
	case 'i':
		return d.decodeInt(data)
	case 'l':
		return d.decodeList(data)
	case 'd':
		return d.decodeDict(data)
	default:
		return d.decodeString(data)
	}
}

func (d *Decoder) write(v any, got any) error {
    switch v := v.(type) {
	case *int:
		val, ok := got.(int)
		if !ok {
			return errors.New("beancode: invalid type")
		}
		*v = val
	case *string:
		val, ok := got.(string)
		if !ok {
			return errors.New("beancode: invalid type")
		}
		*v = val
    case *[]any:
        val, ok := got.([]any)
        if !ok {
            return errors.New("beancode: invalid type")
        }
        *v = val
    case *map[string]any:
        val, ok := got.(map[string]any)
        if !ok {
            return errors.New("beancode: invalid type")
        }
        *v = val
    }

    return nil
}

func (d *Decoder) decodeInt(data []byte) (int, error) {
	d.index++
	end := bytes.IndexByte(data[d.index:], 'e')
	if end == -1 {
		return 0, errors.New("beancode: invalid int")
	}

	end += d.index
	val, err := strconv.Atoi(string(data[d.index:end]))
	if err != nil {
		return 0, err
	}

	d.index = end + 1
	return val, nil
}

func (d *Decoder) decodeList(data []byte) ([]any, error) {
	got := make([]any, 0)
	d.index++

	for {
		if d.index == len(data) {
			return nil, errors.New("beancode: out of bounds")
		}
		if data[d.index] == 'e' {
			d.index++
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
	d.index++

	for {
		if d.index == len(data) {
			return nil, errors.New("beancode: out of bounds")
		}
		if data[d.index] == 'e' {
			d.index++
			return got, nil
		}
		key, err := d.decodeString(data)
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

func (d *Decoder) decodeString(data []byte) (string, error) {
	colon := bytes.IndexByte(data[d.index:], ':')
	if colon == -1 {
		return "", errors.New("beancode: invalid string")
	}

	colon += d.index
	length, err := strconv.Atoi(string(data[d.index:colon]))
	if err != nil {
		return "", err
	}

	start := colon + 1
	end := start + length
	if end > len(data) {
		return "", errors.New("beancode: invalid string length")
	}

	d.index = end
	return string(data[start:end]), nil
}
