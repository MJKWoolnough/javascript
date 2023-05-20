package javascript

import (
	"errors"
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
		{ // 25
			"'\x41'",
			"A",
			nil,
		},
		{ // 26
			"'\\x41'",
			"A",
			nil,
		},
		{ // 27
			"'\\x4G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 28
			"'\\xG1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 29
			"'\\u0041'",
			"A",
			nil,
		},
		{ // 30
			"'\\u004G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 31
			"'\\u00G1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 32
			"'\\u0G41'",
			"",
			ErrInvalidQuoted,
		},
		{ // 33
			"'\\uG041'",
			"",
			ErrInvalidQuoted,
		},
		{ // 34
			"'\\c'",
			"",
			ErrInvalidQuoted,
		},
		{ // 35
			"'\n'",
			"",
			ErrInvalidQuoted,
		},
		{ // 36
			"'\\0'",
			"\000",
			nil,
		},
		{ // 37
			"'\\01'",
			"",
			ErrInvalidQuoted,
		},
		{ // 38
			"'\\u{41}'",
			"A",
			nil,
		},
		{ // 39
			"'\\u{}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 40
			"'\\u{41G}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 41
			"'\\u{41'",
			"",
			ErrInvalidQuoted,
		},
	} {
		o, err := Unquote(test.Input)
		if !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.Err, err)
		} else if o != test.Output {
			t.Errorf("test %d: from %s, expecting output %q, got %q", n+1, test.Input, test.Output, o)
		}
	}
}
