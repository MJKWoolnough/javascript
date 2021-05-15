package scope

import (
	"errors"
	"fmt"
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
				scope.addBinding(javascript.UnwrapConditional(s.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference, BindingRef)
				return scope, nil
			},
		},
		{ // 5
			`function a(){}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				if err := scope.setBinding(s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier, BindingHoistable); err != nil {
					return nil, err
				}
				scope.newFunctionScope(s.StatementList[0].Declaration.FunctionDeclaration)
				return scope, nil
			},
		},
		{ // 6
			`const a = () => true`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				if err := scope.setBinding(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier, BindingLexical); err != nil {
					return nil, err
				}
				scope.newArrowFunctionScope(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction)
				return scope, nil
			},
		},
		{ // 7
			`a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 8
			`{}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Bindings: make(map[string][]Binding),
				}
				scope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Statement.BlockStatement: &Scope{
						IsLexicalScope: true,
						Parent:         scope,
						Scopes:         make(map[fmt.Formatter]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				return scope, nil
			},
		},
		{ // 9
			`function a(){}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				fScope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
					Bindings: map[string][]Binding{
						"this":      []Binding{},
						"arguments": []Binding{},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration: fScope,
				}
				return scope, nil
			},
		},
		{ // 10
			`const a = () => true`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: {
						Parent:   scope,
						Scopes:   make(map[fmt.Formatter]*Scope),
						Bindings: make(map[string][]Binding),
					},
				}
				return scope, nil
			},
		},
		{ // 11
			`let a, a`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateBinding
			},
		},
		{ // 12
			`let a;a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 13
			`let a;{a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         make(map[fmt.Formatter]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 14
			`let a;{let a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         make(map[fmt.Formatter]*Scope),
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 15
			`{function a(){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
				}
				abind := []Binding{
					{
						BindingType: BindingHoistable,
						Scope:       bscope,
						Token:       s.StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
					},
				}
				bscope.Bindings = map[string][]Binding{"a": abind}
				bscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration: &Scope{
						Parent: bscope,
						Scopes: make(map[fmt.Formatter]*Scope),
						Bindings: map[string][]Binding{
							"this":      []Binding{},
							"arguments": []Binding{},
						},
					},
				}
				scope.Bindings = map[string][]Binding{"a": abind}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Statement.BlockStatement: bscope}
				return scope, nil
			},
		},
		{ // 16
			`let a;{function a(){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				bscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration: &Scope{
						Parent: bscope,
						Scopes: make(map[fmt.Formatter]*Scope),
						Bindings: map[string][]Binding{
							"this":      []Binding{},
							"arguments": []Binding{},
						},
					},
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
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
