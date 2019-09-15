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
		{
			"\"\"",
			"",
			nil,
		},
		{
			"''",
			"",
			nil,
		},
		{
			"\"a\"",
			"a",
			nil,
		},
		{
			"'a'",
			"a",
			nil,
		},
		{
			"\"\\\"'\\n\\t\"",
			"\"'\n\t",
			nil,
		},
		{
			"'\x41'",
			"A",
			nil,
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
