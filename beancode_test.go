package beancode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	var got map[string]any
	input := "d3:fool3:boo3:bare3:bood3:fooi100e3:bari100eee"
	want := map[string]any{
		"foo": []any{"boo", "bar"},
		"boo": map[string]any{
			"foo": 100,
			"bar": 100,
		},
	}

	err := Unmarshal(input, &got)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestMarshal(t *testing.T) {
	input := map[string]any{
		"foo": []any{"boo", "bar"},
		"boo": map[string]any{
			"foo": 100,
			"bar": 100,
		},
	}
	want := "d3:fool3:boo3:bare3:bood3:fooi100e3:bari100eee"

	val, err := Marshal(input)
	require.NoError(t, err)
	require.Equal(t, want, val)
}