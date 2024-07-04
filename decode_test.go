package beancode

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeEmpty(t *testing.T) {
	var got string
	input := ""
	output := errors.New("beancode: empty input")
	
	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, output, err)
}

func TestDecodeInvalidInt(t *testing.T) {
	var got int
	input := "3:foo"
	output := errors.New("beancode: invalid type")

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, output, err)
}

func TestDecodeInvalidString(t *testing.T) {
	var got string
	input := "i42e"
	output := errors.New("beancode: invalid type")

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, output, err)
}

func TestDecodeInvalidList(t *testing.T) {
	var got map[string]any
	input := "li1ei2ei3ee"
	output := errors.New("beancode: invalid type")

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, output, err)
}

func TestDecodeInvalidDict(t *testing.T) {
	var got []any
	input := "d3:fooi1e3:bari2e3:booi3ee"
	output := errors.New("beancode: invalid type")

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.Error(t, err)
	require.Equal(t, output, err)
}

func TestDecodeInt(t *testing.T) {
	var got int
	input := "i42e"
	output := 42

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, output, got)
}

func TestDeCodeString(t *testing.T) {
	var got string
	input := "3:foo"
	output := "foo"

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, output, got)
}

func TestDecodeList(t *testing.T) {
	var got []any
	input := "l3:foo3:bari42ee"
	output := []any{"foo", "bar", 42}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, output, got)
}

func TestDecodeDict(t *testing.T) {
	var got map[string]any
	input := "d3:foo3:bar3:barl3:fooi42eee"
	output := map[string]any{
		"foo": "bar",
		"bar": []any{"foo", 42},
	}

	formatInput := bytes.NewReader([]byte(input))
	err := NewDecoder(formatInput).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, output, got)
}
