package minify

import (
	"fmt"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/parser"
)

func TestOrderedScope(t *testing.T) {
	for n, test := range [...]struct {
		Input    string
		Bindings []string
	}{
		{ // 1
			"",
			[]string{},
		},
		{ // 2
			"let a = 1",
			[]string{
				"a",
			},
		},
		{ // 3
			"let a = 1, b = 2;a",
			[]string{
				"a",
				"b",
			},
		},
		{ // 4
			"let a = 1, b = 2;b",
			[]string{
				"b",
				"a",
			},
		},
		{ // 5
			"function b() { a()} function c() {b()} function a(){} b()",
			[]string{
				"b",
				"a",
				"c",
			},
		},
		{ // 6
			"window",
			[]string{},
		},
		{ // 7
			"console.log(1);function a() {console.log(2)}",
			[]string{"a"},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		if m, err := javascript.ParseModule(&tk); err != nil {
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
		{ // 1
		},
		{ // 2
			"let value = 1;",
			"let _ = 1;",
		},
		{ // 3
			"let value = 1, anotherValue = 2;",
			"let $ = 1, _ = 2;",
		},
		{ // 4
			"const bValue = 1;function aFunction(aValue, bValue) {aValue}",
			"const $ = 1;\n\nfunction _(_, $) {\n	_;\n}",
		},
		{ // 5
			"const aValue = 1;aValue;aValue;function aFunction(aValue, bValue){aValue}",
			"const _ = 1;\n\n_;\n\n_;\n\nfunction $(_, $) {\n	_;\n}",
		},
		{ // 6
			"const aValue = 1;{const aValue = 2;{const aValue = 3}}",
			"const _ = 1;\n\n{\n	const _ = 2;\n	{\n		const _ = 3;\n	}\n}",
		},
		{ // 7
			"let aValue = 1;{let bValue = 2;{aValue = 3}}",
			"let _ = 1;\n\n{\n	let $ = 2;\n	{\n		_ = 3;\n	}\n}",
		},
		{ // 8
			"function aFunction(){}",
			"function _() {}",
		},
		{ // 9
			"class aClass {}",
			"class _ {}",
		},
		{ // 10
			"class aClass {}\nclass bClass extends aClass {}",
			"class _ {}\n\nclass $ extends _ {}",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		if m, err := javascript.ParseModule(&tk); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if err = renameIdentifiers(m); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if str := fmt.Sprintf("%s", m); str != test.Output {
			t.Errorf("test %d: expecting output:\n%s\n, got:\n%s", n+1, test.Output, str)
		}
	}
}
