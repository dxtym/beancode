package beancode

import (
	"bytes"
)

// Unmarshaler is the interface implemented by types
// that can unmarshal a Bencode of themselves.
type Unmarshaler interface {
	Unmarshal(data string, v any) error
}

// Marshaler is the interface implemented by types
// that can marshal themselves to valid Bencode.
type Marshaler interface {
	Marshaler(data any) (string, error)
}

// Unmarshal parses the Bencoded data and stores
// the result in the value pointed by v.
func Unmarshal(data string, v any) error {
	rdr := bytes.NewReader([]byte(data))
	return NewDecoder(rdr).Decode(v)
}

// Marshal converts the data to corresponding
// value of Bencode encoding and returns it.
func Marshal(data any) (string, error) {
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}