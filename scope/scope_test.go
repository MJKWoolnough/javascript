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
		{ // 2
			`if(true) false;`,
			func(s *javascript.Script) (*Scope, error) {
				return NewScope(), nil
			},
		},
		{ // 3
			`{}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.newLexicalScope(s.StatementList[0].Statement.BlockStatement)
				return scope, nil
			},
		},
		{ // 4
			`a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.addBinding(javascript.UnwrapConditional(s.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference)
				return scope, nil
			},
		},
		{ // 5
			`function a(){}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				if err := scope.setBinding(s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier, true); err != nil {
					return nil, err
				}
				scope.newFunctionScope(s.StatementList[0].Declaration.FunctionDeclaration)
				return scope, nil
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
