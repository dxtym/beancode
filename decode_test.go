package beancode

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeEmpty(t *testing.T) {
	var got string
	input := ""
	want := &DecodeError{"empty input"}
	
	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.EqualError(t, err, want.Error())
}

func TestDecodeInvalidInt(t *testing.T) {
	var got int
	input := "3:foo"
	want := &InvalidTypeError{reflect.TypeFor[int](), reflect.TypeFor[string]()}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, want, err)
}

func TestDecodeInvalidString(t *testing.T) {
	var got string
	input := "i42e"
	want := &InvalidTypeError{reflect.TypeFor[string](), reflect.TypeFor[int]()}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, want, err)
}

func TestDecodeInvalidList(t *testing.T) {
	var got map[string]any
	input := "li1ei2ei3ee"
	want := &InvalidTypeError{reflect.TypeFor[map[string]any](), reflect.TypeFor[[]any]()}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, want, err)
}

func TestDecodeInvalidDict(t *testing.T) {
	var got []any
	input := "d3:fooi1e3:bari2e3:booi3ee"
	want := &InvalidTypeError{reflect.TypeFor[[]any](), reflect.TypeFor[map[string]any]()}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, want, err)
}

func TestDecodeInt(t *testing.T) {
	var got int
	input := "i42e"
	want := 42

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestDeCodeString(t *testing.T) {
	var got string
	input := "3:foo"
	want := "foo"

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestDecodeList(t *testing.T) {
	var got []any
	input := "l3:foo3:bari42ee"
	want := []any{"foo", "bar", 42}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestDecodeDict(t *testing.T) {
	var got map[string]any
	input := "d3:foo3:bar3:barl3:fooi42eee"
	want := map[string]any{
		"foo": "bar",
		"bar": []any{"foo", 42},
	}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
