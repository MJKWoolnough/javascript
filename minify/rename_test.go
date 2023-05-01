package minify

import (
	"fmt"
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
		{
			"function b() { a()} function c() {b()} function a(){} b()",
			[]string{
				"b",
				"a",
				"c",
			},
		},
		{
			"window",
			[]string{},
		},
		{
			"console.log(1);function a() {console.log(2)}",
			[]string{"a"},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if s, err := scope.ModuleScope(m, nil); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			bs := orderedScope(s)
			bindings := make([]string, 0, len(bs))
			for _, b := range bs {
				if !b.NameSet {
					bindings = append(bindings, b.Name)
				}
			}
			if !reflect.DeepEqual(bindings, test.Bindings) {
				t.Errorf("test %d: expecting bindings: %v, got %v", n+1, test.Bindings, bindings)
			}
		}
	}
}

func init() {
	startChars = []byte{'_', '$'}
	extraChars = []byte{'a', 'b'}
}

func TestUniqueName(t *testing.T) {
	used := make(map[string]struct{})
	for n, next := range []string{"_", "$", "_a", "_b", "$a", "$b", "_aa"} {
		if name := makeUniqueName(used); name != next {
			t.Errorf("test %d: expecting name %s, got %s", n+1, next, name)
		}
		used[next] = struct{}{}
	}
}

func TestRename(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
	}{
		{},
		{
			"let value = 1;",
			"let _ = 1;",
		},
		{
			"let value = 1, anotherValue = 2;",
			"let $ = 1, _ = 2;",
		},
		{
			"const bValue = 1;function aFunction(aValue, bValue) {aValue}",
			"const $ = 1;\n\nfunction _(_, $) {\n	_;\n}",
		},
		{
			"const aValue = 1;aValue;aValue;function aFunction(aValue, bValue){aValue}",
			"const _ = 1;\n\n_;\n\n_;\n\nfunction $(_, $) {\n	_;\n}",
		},
		{
			"const aValue = 1;{const aValue = 2;{const aValue = 3}}",
			"const _ = 1;\n\n{\n	const _ = 2;\n	{\n		const _ = 3;\n	}\n}",
		},
		{
			"let aValue = 1;{let bValue = 2;{aValue = 3}}",
			"let _ = 1;\n\n{\n	let $ = 2;\n	{\n		_ = 3;\n	}\n}",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if err = renameIdentifiers(m); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if str := fmt.Sprintf("%s", m); str != test.Output {
			t.Errorf("test %d: expecting output %s, got %s", n+1, test.Output, str)
		}
	}
}
