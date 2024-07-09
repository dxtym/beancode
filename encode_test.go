package beancode

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeEmpty(t *testing.T) {
	var buf bytes.Buffer
	input := ""
	want := errors.New("beancode: empty input")

	err := NewEncoder(&buf).Encode(input)
	require.EqualError(t, err, want.Error())
}

func TestEncodeInt(t *testing.T) {
	var buf bytes.Buffer
	input := 42
	want := "i42e"
	
	err := NewEncoder(&buf).Encode(input)
	require.NoError(t, err)
	require.Equal(t, want, buf.String())
}

func TestEncodeStr(t *testing.T) {
	var buf bytes.Buffer
	input := "foo"
	want := "3:foo"

	err := NewEncoder(&buf).Encode(input)
	require.NoError(t, err)
	require.Equal(t, want, buf.String())
}

func TestEncodeList(t *testing.T) {
	var buf bytes.Buffer
	input := []any{"foo", "bar", 42}
	want := "l3:foo3:bari42ee"

	err := NewEncoder(&buf).Encode(input)
	require.NoError(t, err)
	require.Equal(t, want, buf.String())
}

func TestEncodeDict(t *testing.T) {
	var buf bytes.Buffer
	input := map[string]any{
		"foo": "bar",
		"bar": []any{"foo", 42},
	}
	want := "d3:foo3:bar3:barl3:fooi42eee"

	err := NewEncoder(&buf).Encode(input)
	require.NoError(t, err)
	require.Equal(t, want, buf.String())
}