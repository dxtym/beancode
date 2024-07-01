package beancode

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	testcases := []struct {
		name  string
		input string
		want  any
	}{
		{
			name:  "Integer",
			input: "i42e",
			want:  42,
		},
		{
			name:  "String",
			input: "3:foo",
			want:  "foo",
		},
		{
			name:  "List",
			input: "l3:foo3:bari42ee",
			want:  []any{"foo", "bar", 42},
		},
		{
			name:  "Dict",
			input: "d3:foo3:bar3:barl3:fooi42eee",
			want: map[string]any{
				"foo": "bar",
				"bar": []any{"foo", 42},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			typeof := reflect.TypeOf(tc.input)
			got := reflect.New(typeof).Interface()

			formatInput := bytes.NewReader([]byte(tc.input))
			err := NewDecoder(formatInput).Decode(&got)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)

			t.Log(got)
		})
	}
}
