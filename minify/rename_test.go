package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/parser"
)

func TestOrdererScope(t *testing.T) {
	for n, test := range [...]struct {
		Input    string
		Bindings []string
	}{
		{
			"",
			[]string{},
		},
		{
			"let a = 1",
			[]string{
				"a",
			},
		},
		{
			"let a = 1, b = 2;a",
			[]string{
				"a",
				"b",
			},
		},
		{
			"let a = 1, b = 2;b",
			[]string{
				"b",
				"a",
			},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if s, err := scope.ModuleScope(m, nil); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			bs := ordererScope(s)
			bindings := make([]string, len(bs))
			for n := range bs {
				bindings[n] = bs[n].Name
			}
			if !reflect.DeepEqual(bindings, test.Bindings) {
				t.Errorf("test %d: expecting bindings: %v, got %v", n+1, test.Bindings, bindings)
			}
		}
	}
}
