package scope

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func TestScriptScope(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output func(*javascript.Script) (*Scope, error)
	}{
		{ // 1
			``,
			func(s *javascript.Script) (*Scope, error) {
				return NewScope(), nil
			},
		},
	} {
		source, err := javascript.ParseScript(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else {
			tscope, terr := test.Output(source)
			scope, err := ScriptScope(source, nil)
			if terr != nil && err != nil {
				if !errors.Is(terr, err) {
					t.Errorf("test %d: expecting error: %s\ngot: %s", n+1, terr, err)
				}
			} else if terr != nil {
				t.Errorf("test %d: received no error when expecting: %s", n+1, terr)
			} else if err != nil {
				t.Errorf("test %d: receieved error when expecting none: %s", n+1, err)
			} else if !reflect.DeepEqual(scope, tscope) {
				t.Errorf("test %d: result did not match expected", n+1)
			}
		}
	}
}
