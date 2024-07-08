package beancode

import (
	"bytes"
)

func Unmarshal(data string, v any) error {
	rdr := bytes.NewReader([]byte(data))
	return NewDecoder(rdr).Decode(v)
}

func Marshal(data any) (string, error) {
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}