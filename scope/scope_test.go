package scope

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

type indentPrinter struct {
	io.Writer
}

var indent = []byte{'	'}

func (i *indentPrinter) Write(p []byte) (int, error) {
	var (
		total int
		last  int
	)
	for n, c := range p {
		if c == '\n' {
			m, err := i.Writer.Write(p[last : n+1])
			total += m
			if err != nil {
				return total, err
			}
			_, err = i.Writer.Write(indent)
			if err != nil {
				return total, err
			}
			last = n + 1
		}
	}
	if last != len(p) {
		m, err := i.Writer.Write(p[last:])
		total += m
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func (i *indentPrinter) Print(args ...interface{}) {
	fmt.Fprint(i, args...)
}

func (i *indentPrinter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(i, format, args...)
}

func (i *indentPrinter) WriteString(s string) (int, error) {
	return i.Write([]byte(s))
}

func (s *Scope) Format(st fmt.State, _ rune) { s.printScope(&indentPrinter{st}) }

func (s *Scope) printScope(w *indentPrinter) {
	if s.IsLexicalScope {
		w.WriteString("LexicalScope {")
	} else {
		w.WriteString("Scope {")
	}
	pp := &indentPrinter{w}
	qq := &indentPrinter{pp}
	rr := &indentPrinter{qq}
	pp.WriteString("\nParent: ")
	if s.Parent == nil {
		pp.WriteString("nil")
	} else {
		pp.Printf("%p", s.Parent)
	}
	pp.WriteString("\nScopes: [")
	for t, scope := range s.Scopes {
		qq.WriteString("\n")
		qq.Printf("%p", t)
		qq.WriteString(": ")
		scope.printScope(qq)
	}
	pp.WriteString("\n]\nBindings: [")
	for ref, bindings := range s.Bindings {
		qq.WriteString("\n")
		qq.WriteString(ref)
		qq.WriteString(": [")
		for _, binding := range bindings {
			rr.WriteString("\n[")
			rr.WriteString("\n	BindingType: ")
			rr.Print(binding.BindingType)
			rr.WriteString("\n	Scope: ")
			rr.Printf("%p", binding.Scope)
			rr.WriteString("\n	Token: ")
			rr.Print(binding.Token)
			rr.WriteString("\n]")
		}
		qq.WriteString("\n]")
	}
	pp.WriteString("\n]")
	w.WriteString("\n}")
}

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
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
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
		{ // 17
			`var a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       scope,
							Token:       s.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 18
			`var a;a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       scope,
							Token:       s.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
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
		{ // 19
			`var a;{a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{}
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
							BindingType: BindingVar,
							Scope:       scope,
							Token:       s.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
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
		{ // 20
			`a;{var a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         make(map[fmt.Formatter]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
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
		{ // 21
			`var a;{var a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         make(map[fmt.Formatter]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       scope,
							Token:       s.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 22
			`function b(){var a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 23
			`function b() {var a;a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 24
			`function b() {var a;{a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 25
			`function b(){a;{var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 26
			`function b(){var a;{var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 27
			`function a(b, c, ...d) {b}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"b": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"c": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[1].SingleNameBinding,
						},
					},
					"d": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FunctionRestParameter.BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 28
			`const a = b => b`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				ascope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				ascope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       ascope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope,
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
				return scope, nil
			},
		},
		{ // 29
			`function a() {let [b, c, ...d] = [0, 1, 2, 3];d}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"b": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
						},
					},
					"c": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingElementList[1].SingleNameBinding,
						},
					},
					"d": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingRestElement.SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 30
			`function a() {let {z: b, c, ...d} = {};d}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"b": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
					"c": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
					},
					"d": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingRestProperty,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 31
			`function a(b) {let c = b => [b, d], d = 1;[d, b] = c(b)}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				ascope := &Scope{
					Parent: fscope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				ascope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       ascope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).ArrayLiteral.ElementList[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"b": []Binding{
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral.ElementList[1].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.CallExpression).Arguments.ArgumentList[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"c": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
					"d": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).ArrayLiteral.ElementList[1].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral.ElementList[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 32
			`function a() {return this}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this": []Binding{
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).This,
						},
					},
					"arguments": []Binding{},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 33
			`function a() {return arguments}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this": []Binding{},
					"arguments": []Binding{
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 34
			`let a;function b() {return a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[1].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"b": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[1].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 35
			`const {a: b, c, ...d} = {}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
					"c": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
					},
					"d": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingRestProperty,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 36
			`let {a, b} = {};({a, b} = {})`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).CoverParenthesizedExpressionAndArrowParameterList.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"b": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).CoverParenthesizedExpressionAndArrowParameterList.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral.PropertyDefinitionList[1].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 37
			`var c;{let c;{var c}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[1].Statement.BlockStatement.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
				}
			},
		},
		{ // 38
			`const {a: [b]} = {}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"b": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 39
			`var [{a}] = []`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[fmt.Formatter]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       scope,
							Token:       s.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 40
			`function a() {for (let a = 0; a < 2; a++){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{Parent: scope}
				iscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
				}
				iscope.Scopes = map[fmt.Formatter]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: &Scope{
						IsLexicalScope: true,
						Parent:         iscope,
						Scopes:         make(map[fmt.Formatter]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				iscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.InitLexical.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(javascript.WrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Conditional.Expressions[0].ConditionalExpression).(*javascript.RelationalExpression).RelationalExpression)).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Afterthought.Expressions[0].ConditionalExpression).(*javascript.UpdateExpression).LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 41
			`function a() {for (var a = 0; a < 2; a++){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				iscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         iscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingVar,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.InitVar[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(javascript.WrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Conditional.Expressions[0].ConditionalExpression).(*javascript.RelationalExpression).RelationalExpression)).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Afterthought.Expressions[0].ConditionalExpression).(*javascript.UpdateExpression).LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a":         iscope.Bindings["a"],
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 42
			`function a() {for (b = 0; b < 2; b++){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				iscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Bindings:       make(map[string][]Binding),
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         iscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				fscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
					"b": []Binding{
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.InitExpression.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(javascript.WrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Conditional.Expressions[0].ConditionalExpression).(*javascript.RelationalExpression).RelationalExpression)).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       iscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Afterthought.Expressions[0].ConditionalExpression).(*javascript.UpdateExpression).LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 43
			`function a() {let a;for (const a in {}){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				iscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         iscope,
					Scopes:         make(map[fmt.Formatter]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.ForBindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      []Binding{},
					"arguments": []Binding{},
					"a": []Binding{
						{
							BindingType: BindingLexical,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[fmt.Formatter]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": []Binding{
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
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
				t.Errorf("test %d: result did not match expected\nexpecting: %s\ngot: %s", n+1, tscope, scope)
			}
		}
	}
}
