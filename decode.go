<<<<<<< HEAD
package main
=======
package beancode
>>>>>>> 86e3eec (feat: write to var and unit tests)

import (
	"bytes"
	"errors"
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> 86e3eec (feat: write to var and unit tests)
	"io"
	"strconv"
)

type Decoder struct {
<<<<<<< HEAD
	r io.Reader
=======
	r     io.Reader
>>>>>>> 86e3eec (feat: write to var and unit tests)
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
		return errors.New("bencode: empty input")
	}

<<<<<<< HEAD
	d.index = 0
	var got []any
	for d.index < len(data) {
		val, err := d.decode(data)
		if err != nil {
			return err
		}
		got = append(got, val)
	}
	fmt.Print(got)
	return nil // TODO: write to the v
=======
	var val any
	for d.index < len(data) {
		val, err = d.decode(data)
		if err != nil {
			return err
		}
	}

	return d.write(v, val)
>>>>>>> 86e3eec (feat: write to var and unit tests)
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

<<<<<<< HEAD
// parse int
=======
func (d *Decoder) write(v, got any) error {
	switch v := v.(type) {
	case *any:
		*v = got
	case *[]any:
		val, ok := got.([]any)
		if !ok {
			return errors.New("bencode: invalid type")
		}
		*v = val
	case *map[string]any:
		val, ok := got.(map[string]any)
		if !ok {
			return errors.New("bencode: invalid type")
		}
		*v = val
	default:
		// TODO: handle struct
	}

	return nil
}

>>>>>>> 86e3eec (feat: write to var and unit tests)
func (d *Decoder) decodeInt(data []byte) (int, error) {
	d.index++
	end := bytes.IndexByte(data[d.index:], 'e')
	if end == -1 {
		return 0, errors.New("bencode: invalid int")
	}

	end += d.index
	val, err := strconv.Atoi(string(data[d.index:end]))
	if err != nil {
		return 0, err
	}

	d.index = end + 1
	return val, nil
}

<<<<<<< HEAD
// parse list
=======
>>>>>>> 86e3eec (feat: write to var and unit tests)
func (d *Decoder) decodeList(data []byte) ([]any, error) {
	got := make([]any, 0)
	d.index++

	for {
		if d.index == len(data) {
<<<<<<< HEAD
			return nil, errors.New("bencode: reached out of bounds")
=======
			return nil, errors.New("bencode: out of bounds")
>>>>>>> 86e3eec (feat: write to var and unit tests)
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

<<<<<<< HEAD
// parse dictionary
=======
>>>>>>> 86e3eec (feat: write to var and unit tests)
func (d *Decoder) decodeDict(data []byte) (map[string]any, error) {
	got := make(map[string]any)
	d.index++

	for {
		if d.index == len(data) {
<<<<<<< HEAD
			return nil, errors.New("bencode: reached out of bounds")
=======
			return nil, errors.New("bencode: out of bounds")
>>>>>>> 86e3eec (feat: write to var and unit tests)
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

<<<<<<< HEAD
// parse string
=======
>>>>>>> 86e3eec (feat: write to var and unit tests)
func (d *Decoder) decodeString(data []byte) (string, error) {
	colon := bytes.IndexByte(data[d.index:], ':')
	if colon == -1 {
		return "", errors.New("bencode: invalid string")
	}

<<<<<<< HEAD
	colon += d.index 
=======
	colon += d.index
>>>>>>> 86e3eec (feat: write to var and unit tests)
	length, err := strconv.Atoi(string(data[d.index:colon]))
	if err != nil {
		return "", err
	}

	start := colon + 1
	end := start + length
	if end > len(data) {
		return "", errors.New("bencode: invalid string length")
	}
<<<<<<< HEAD
	d.index = end
	return string(data[start:end]), nil
}


func main() {
	test := []byte("d2:mei1e3:youl3:one3:twoee")
	d := NewDecoder(bytes.NewReader(test))

	var i int
	d.Decode(&i)
}
=======

	d.index = end
	return string(data[start:end]), nil
}
>>>>>>> 86e3eec (feat: write to var and unit tests)
