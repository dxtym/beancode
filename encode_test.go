package beancode

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeInvalid(t *testing.T) {
	testCases := []struct {
		id int
		name  string
		input any
		want  error
	}{
		{
			id:   0,
			name:  "zero value",
			input: "",
			want:  fmt.Errorf("beancode: zero value"),
		},
		{
			id:   1,
			name:  "unsupported type",
			input: 42.0,
			want:  fmt.Errorf("beancode: unsupported type float64"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := NewEncoder(&buf).Encode(tc.input)
			require.EqualError(t, err, tc.want.Error())
		})
	}
}

func TestEncodeValid(t *testing.T) {
	testCases := []struct {
		id int
		name  string
		input any
		want  string
	}{
		{
			id: 0,
			name:  "int",
			input: 42,
			want:  "i42e",
		},
		{
			id: 1,
			name:  "string",
			input: "foo",
			want:  "3:foo",
		},
		{
			id: 2,
			name:  "list",
			input: []any{"foo", "bar"},
			want:  "l3:foo3:bare",
		},
		{
			id: 3,
			name: "dict",
			input: map[string]any{
				"foo": "bar",
				"bar": []any{"foo", "boo"},
			},
			want: "d3:foo3:bar3:barl3:foo3:booee",
		},
		{
			id: 4,
			name: "struct",
			input: struct {
				Foo string `beancode:"Foo"`
				Bar int    `beancode:"Bar"`
			}{
				Foo: "bar",
				Bar: 42,
			},
			want: "d3:Foo3:bar3:Bari42ee",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := NewEncoder(&buf).Encode(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, buf.String())
		})
	}
}
