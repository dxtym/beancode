package beancode

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: update testcase based on new err types
func TestDecodeInvalid(t *testing.T) {
	testCases := []struct {
		id    int
		name  string
		input string
		want  error
	}{
		{
			id:    0,
			name:  "int",
			input: "3:foo",
			want:  &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected int, got %v", reflect.TypeFor[string]()),
			},
		},
		{
			id:    1,
			name:  "string",
			input: "i42e",
			want:  &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected string, got %v", reflect.TypeFor[int]()),
			},
		},
		{
			id:    2,
			name:  "array",
			input: "d3:fooi1e3:bari2e3:booi3ee",
			want:  &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected []any, got %v", reflect.TypeFor[map[string]any]()),
			},
		},
		{
			id:    3,
			name:  "map",
			input: "li1ei2ei3ee",
			want:  &DecodeError{
				Type: "type mismatch",
				Err: fmt.Errorf("expected map[string]any, got %v",  reflect.TypeFor[[]any]()),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			switch tc.id {
			case 0:
				got = new(int)
			case 1:
				got = new(string)
			case 2:
				got = &[]any{}
			case 3:
				got = &map[string]any{}
			}

			formatInput := bytes.NewReader([]byte(tc.input))
			err := NewDecoder(formatInput).Decode(got)
			require.Error(t, err)
			require.Equal(t, tc.want, err)
		})
	}
}

func TestDecodeValid(t *testing.T) {
	testCases := []struct {
		id    int
		name  string
		input string
		want  any
	}{
		{
			id:    0,
			name:  "int",
			input: "i42e",
			want:  42,
		},
		{
			id:    1,
			name:  "string",
			input: "3:foo",
			want:  "foo",
		},
		{
			id:    2,
			name:  "array",
			input: "l3:foo3:bare",
			want:  []any{"foo", "bar"},
		},
		{
			id:    3,
			name:  "map",
			input: "d3:foo3:bar3:barl3:foo3:booee",
			want: map[string]any{
				"foo": "bar",
				"bar": []any{"foo", "boo"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			switch tc.id {
			case 0:
				got = new(int)
			case 1:
				got = new(string)
			case 2:
				got = &[]any{}
			case 3:
				got = &map[string]any{}
			}

			formatInput := bytes.NewReader([]byte(tc.input))
			err := NewDecoder(formatInput).Decode(got)
			require.NoError(t, err)

			switch v := got.(type) {
			case *int:
				require.Equal(t, tc.want, *v)
			case *string:
				require.Equal(t, tc.want, *v)
			case *[]any:
				require.Equal(t, tc.want, *v)
			case *map[string]any:
				require.Equal(t, tc.want, *v)
			}
		})
	}
}
