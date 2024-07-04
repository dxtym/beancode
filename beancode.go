package beancode

import (
	"bytes"
)

func Unmarshal(data string, v any) error {
	rdr := bytes.NewReader([]byte(data))
	return NewDecoder(rdr).Decode(v)
}