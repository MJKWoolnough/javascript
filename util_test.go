package javascript

import (
	"reflect"
	"testing"
)

func TestUnquote(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{ // 1
			"\"\"",
			"",
			nil,
		},
		{ // 2
			"''",
			"",
			nil,
		},
		{ // 3
			"\"a\"",
			"a",
			nil,
		},
		{ // 4
			"'a'",
			"a",
			nil,
		},
		{ // 5
			"\"\\'\\\"\\\\\\b\\f\\n\\r\\t\\v\"",
			"'\"\\\b\f\n\r\t\v",
			nil,
		},
		{ // 6
			"\"\x41\"",
			"A",
			nil,
		},
		{ // 7
			"",
			"",
			ErrInvalidQuoted,
		},
		{ // 8
			"\"\\x41\"",
			"A",
			nil,
		},
		{ // 9
			"\"\\x4G\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 10
			"\"\\xG1\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 11
			"\"\\u0041\"",
			"A",
			nil,
		},
		{ // 12
			"\"\\u004G\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 13
			"\"\\u00G1\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 14
			"\"\\u0G41\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 15
			"\"\\uG041\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 16
			"\"\\c\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 17
			"\"\n\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 18
			"\"\\0\"",
			"\000",
			nil,
		},
		{ // 19
			"\"\\01\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 20
			"\"\\u{41}\"",
			"A",
			nil,
		},
		{ // 21
			"\"\\u{}\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 22
			"\"\\u{41G}\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 23
			"\"\\u{41\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 24
			"'\\'\\\"\\\\\\b\\f\\n\\r\\t\\v'",
			"'\"\\\b\f\n\r\t\v",
			nil,
		},
		{ // 6
			"'\x41'",
			"A",
			nil,
		},
		{ // 8
			"'\\x41'",
			"A",
			nil,
		},
		{ // 9
			"'\\x4G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 10
			"'\\xG1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 11
			"'\\u0041'",
			"A",
			nil,
		},
		{ // 12
			"'\\u004G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 13
			"'\\u00G1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 14
			"'\\u0G41'",
			"",
			ErrInvalidQuoted,
		},
		{ // 15
			"'\\uG041'",
			"",
			ErrInvalidQuoted,
		},
		{ // 16
			"'\\c'",
			"",
			ErrInvalidQuoted,
		},
		{ // 17
			"'\n'",
			"",
			ErrInvalidQuoted,
		},
		{ // 18
			"'\\0'",
			"\000",
			nil,
		},
		{ // 19
			"'\\01'",
			"",
			ErrInvalidQuoted,
		},
		{ // 20
			"'\\u{41}'",
			"A",
			nil,
		},
		{ // 21
			"'\\u{}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 22
			"'\\u{41G}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 23
			"'\\u{41'",
			"",
			ErrInvalidQuoted,
		},
	} {
		o, err := Unquote(test.Input)
		if !reflect.DeepEqual(err, test.Err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.Err, err)
		} else if o != test.Output {
			t.Errorf("test %d: from %s, expecting output %q, got %q", n+1, test.Input, test.Output, o)
		}
	}
}
