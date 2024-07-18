package beancode

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeInvalid(t *testing.T) {
	testCases := []struct {
		id    int
		name  string
		input string
		want  error
	}{
		{
			id:    0,
			name:  "invalid int",
			input: "3:foo",
			want:  fmt.Errorf("beancode: expected int, got %v", reflect.TypeFor[string]()),
		},
		{
			id:    1,
			name:  "invalid string",
			input: "i42e",
			want:  fmt.Errorf("beancode: expected string, got %v", reflect.TypeFor[int]()),
		},
		{
			id:    2,
			name:  "invalid array",
			input: "d3:fooi1e3:bari2e3:booi3ee",
			want:  fmt.Errorf("beancode: expected []any, got %v", reflect.TypeFor[map[string]any]()),
		},
		{
			id:    3,
			name:  "invalid map",
			input: "li1ei2ei3ee",
			want:  fmt.Errorf("beancode: expected map[string]any, got %v", reflect.TypeFor[[]any]()),
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
				got = new([]any)
			case 3:
				got = new(map[string]any)
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
			name:  "valid int",
			input: "i42e",
			want:  42,
		},
		{
			id:    1,
			name:  "valid string",
			input: "3:foo",
			want:  "foo",
		},
		{
			id:    2,
			name:  "valid array",
			input: "l3:foo3:bare",
			want:  []any{"foo", "bar"},
		},
		{
			id:    3,
			name:  "valid map",
			input: "d3:foo3:bar3:barl3:foo3:booee",
			want: map[string]any{
				"foo": "bar",
				"bar": []any{"foo", "boo"},
			},
		},
		{
			id:    4,
			name:  "valid struct",
			input: "d3:Foo3:bar3:Bari42ee",
			want: struct {
				Foo string `beancode:"Foo"`
				Bar int    `beancode:"Bar"`
			}{
				Foo: "bar",
				Bar: 42,
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
				got = new([]any)
			case 3:
				got = new(map[string]any)
			case 4:
				got = new(struct {
					Foo string `beancode:"Foo"`
					Bar int    `beancode:"Bar"`
				})
			}

			formatInput := bytes.NewReader([]byte(tc.input))
			err := NewDecoder(formatInput).Decode(got)
			require.NoError(t, err)

			// type assert and deref the result 
			switch tc.id {
			case 0:
				require.Equal(t, tc.want, *(got.(*int)))
			case 1:
				require.Equal(t, tc.want, *(got.(*string)))
			case 2:
				require.Equal(t, tc.want, *(got.(*[]any)))
			case 3:
				require.Equal(t, tc.want, *(got.(*map[string]any)))
			case 4:
				require.Equal(t, tc.want, *(got.(*struct {
					Foo string `beancode:"Foo"`
					Bar int    `beancode:"Bar"`
				})))
			}
		})
	}
}
