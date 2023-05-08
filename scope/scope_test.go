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
				if err := scope.setBinding(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier, BindingLexicalConst); err != nil {
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
					Scopes: make(map[javascript.Type]*Scope),
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
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Statement.BlockStatement: {
						IsLexicalScope: true,
						Parent:         scope,
						Scopes:         make(map[javascript.Type]*Scope),
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
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				fScope := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
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
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: {
						Parent:   scope,
						Scopes:   make(map[javascript.Type]*Scope),
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
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
					Scopes:         make(map[javascript.Type]*Scope),
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
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
				bscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration: {
						Parent: bscope,
						Scopes: make(map[javascript.Type]*Scope),
						Bindings: map[string][]Binding{
							"this":      {},
							"arguments": {},
						},
					},
				}
				scope.Bindings = map[string][]Binding{"a": abind}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Statement.BlockStatement: bscope}
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
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				bscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.FunctionDeclaration: {
						Parent: bscope,
						Scopes: make(map[javascript.Type]*Scope),
						Bindings: map[string][]Binding{
							"this":      {},
							"arguments": {},
						},
					},
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				return scope, nil
			},
		},
		{ // 17
			`var a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Statement.BlockStatement: bscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": {
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": {
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
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": {
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
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.BlockStatement: bscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"b": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"b": {
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
					"c": {
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[1].SingleNameBinding,
						},
					},
					"d": {
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				ascope.Bindings = map[string][]Binding{
					"b": {
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
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"b": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
						},
					},
					"c": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingElementList[1].SingleNameBinding,
						},
					},
					"d": {
						{
							BindingType: BindingLexicalLet,
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"b": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
					"c": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
					},
					"d": {
						{
							BindingType: BindingLexicalLet,
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
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				ascope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingFunctionParam,
							Scope:       ascope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.ArrayLiteral).ElementList[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"b": {
						{
							BindingType: BindingFunctionParam,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[1].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.CallExpression).Arguments.ArgumentList[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"c": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
					"d": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.ArrayLiteral).ElementList[1].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this": {
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).This,
						},
					},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this": {},
					"arguments": {
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[1].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"b": {
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
					"c": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
					},
					"d": {
						{
							BindingType: BindingLexicalConst,
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
					"b": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[1].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingLexicalConst,
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
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
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
				iscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: {
						IsLexicalScope: true,
						Parent:         iscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
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
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
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
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a":         iscope.Bindings["a"],
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingBare,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.ForBindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 44
			`function a() {let a;for (const a of []){}}`,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor.ForBindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 45
			`function a() {for (var a of []){}}`,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.ForBindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a":         iscope.Bindings["a"],
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 46
			`function a() {for (let [a] of []){}}`,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.ForBindingPatternArray.BindingElementList[0].SingleNameBinding,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 47
			`function a() {for (let {a} of []){}}`,
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
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				iscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.Statement.BlockStatement: bscope}
				iscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       iscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor.ForBindingPatternObject.BindingPropertyList[0].BindingElement.SingleNameBinding,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.IterationStatementFor: iscope}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 48
			`function a() {var a;try{}catch(a){a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       tscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 49
			`function a() {var a;try{}catch(a){let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchParameterBindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 50
			`function a() {var a;try{}catch(a){a;var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
						{
							BindingType: BindingVar,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[1].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       tscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 51
			`function a() {var a;try{}catch(a){{let a}}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         tscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       bscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				tscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Statement.BlockStatement: bscope}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 52
			`function a() {var a;try{a}finally{}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock: tscope,
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.FinallyBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       tscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 53
			`function a() {var a;try{let a}finally{}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.TryBlock: tscope,
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.FinallyBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 54
			`function a() {try{var a}finally{}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock: tscope,
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.FinallyBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a":         tscope.Bindings["a"],
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 55
			`function a() {try{}catch(b){var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
					},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchBlock.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a":         tscope.Bindings["a"],
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 56
			`function a() {try{}catch(b){let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
					},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.CatchBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 57
			`function a() {try{}finally{var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.FinallyBlock.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.FinallyBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a":         tscope.Bindings["a"],
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 58
			`function a() {try{}finally{let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       tscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.FinallyBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         fscope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.TryStatement.FinallyBlock: tscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 59
			`function a() {let a;switch(a){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				sscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement: sscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 60
			`function a() {let a;switch(a){case a:}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				sscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
					Bindings:       make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement: sscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       sscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 61
			`function a() {let a;switch(a){case a:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				sscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				sscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       sscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       sscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement: sscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 62
			`function a() {let a;switch(a){case 1:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[1].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 63
			`function a() {let a;switch(a){case 1:let a;break;default:break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 64
			`function a() {let a;switch(a){default:break;case 1:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[1].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 65
			`function a() {let a;switch(a){default:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.DefaultClause[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 66
			`function a() {let a;switch(a){case 1:let a;break;default:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.DefaultClause[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 67
			`function a() {switch(0){case 1:{let a};case 2:{let a}}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				sscope := &Scope{
					IsLexicalScope: true,
					Parent:         fscope,
					Bindings:       make(map[string][]Binding),
				}
				bscopea := &Scope{
					IsLexicalScope: true,
					Parent:         sscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscopeb := &Scope{
					IsLexicalScope: true,
					Parent:         sscope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				bscopea.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       bscopea,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				bscopeb.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       bscopeb,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.SwitchStatement.CaseClauses[1].StatementList[0].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				sscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Statement.BlockStatement: bscopea,
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.SwitchStatement.CaseClauses[1].StatementList[0].Statement.BlockStatement: bscopeb,
				}
				fscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.SwitchStatement: sscope,
				}
				fscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 68
			`function a() {const a = () => this}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				ascope := &Scope{
					Parent:   fscope,
					Scopes:   make(map[javascript.Type]*Scope),
					Bindings: make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope}
				fscope.Bindings = map[string][]Binding{
					"this": {
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).This,
						},
					},
					"arguments": {},
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 69
			`function a() {const a = () => arguments}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
				}
				ascope := &Scope{
					Parent:   fscope,
					Scopes:   make(map[javascript.Type]*Scope),
					Bindings: make(map[string][]Binding),
				}
				fscope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope}
				fscope.Bindings = map[string][]Binding{
					"this": {},
					"arguments": {
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{s.StatementList[0].Declaration.FunctionDeclaration: fscope}
				scope.Bindings = map[string][]Binding{
					"a": {
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
		{ // 70
			`class a{}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 71
			`class a{}class b extends a{}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       s.StatementList[1].Declaration.ClassDeclaration.ClassHeritage.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
					"b": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[1].Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 72
			`class a{b(){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				mscope := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
				}
				mscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.ClassDeclaration.ClassBody[0].MethodDefinition: mscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 73
			`class a{b(){}a(){}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				mscopea := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
				}
				mscopea.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				mscopeb := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
				}
				mscopeb.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.ClassDeclaration.ClassBody[0].MethodDefinition: mscopea,
					s.StatementList[0].Declaration.ClassDeclaration.ClassBody[1].MethodDefinition: mscopeb,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 74
			`class a{b(){var a;let b}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				mscope := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
				}
				mscope.Bindings = map[string][]Binding{
					"this":      {},
					"arguments": {},
					"a": {
						{
							BindingType: BindingVar,
							Scope:       mscope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.FunctionBody.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingLexicalLet,
							Scope:       mscope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.FunctionBody.StatementList[1].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.ClassDeclaration.ClassBody[0].MethodDefinition: mscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 75
			`const a = (a) => a`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				ascope := &Scope{
					Parent: scope,
					Scopes: make(map[javascript.Type]*Scope),
				}
				ascope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingFunctionParam,
							Scope:       ascope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.FormalParameters.FormalParameterList[0].SingleNameBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       ascope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction: ascope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 76
			`function a() {b = 1}let b = 0;`,
			func(s *javascript.Script) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: map[javascript.Type]*Scope{},
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Declaration.FunctionDeclaration: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       s.StatementList[1].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingBare,
							Scope:       fscope,
							Token:       s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		source, err := javascript.ParseScript(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else {
			tscope, terr := test.Output(source)
			scope, err := ScriptScope(source, nil)
			if terr != nil && err != nil {
				if !errors.Is(terr, err) {
					t.Errorf("test %d: expecting error: %v\ngot: %v", n+1, terr, err)
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

func TestModuleScope(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output func(*javascript.Module) (*Scope, error)
	}{
		{ // 1
			`import {a as b} from './lib.js';let a = 0`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingImport,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ImportDeclaration.ImportClause.NamedImports.ImportList[0].ImportedBinding,
						},
					},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       m.ModuleListItems[1].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 2
			`export default class MyClass {}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"MyClass": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.DefaultClass.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 3
			`export default class MyClass {static INSTANCE = new MyClass()}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"MyClass": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.DefaultClass.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].ExportDeclaration.DefaultClass.ClassBody[0].FieldDefinition.Initializer.ConditionalExpression).(*javascript.MemberExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 4
			`export default function MyFunc() {}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: map[javascript.Type]*Scope{},
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[0].ExportDeclaration.DefaultFunction: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"MyFunc": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.DefaultFunction.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 5
			`export default function MyFunc() {MyFunc()}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: map[javascript.Type]*Scope{},
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[0].ExportDeclaration.DefaultFunction: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"MyFunc": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.DefaultFunction.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].ExportDeclaration.DefaultFunction.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 6
			`globalThis.console;window;let a = 1;{a;window}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         map[javascript.Type]*Scope{},
					Bindings:       map[string][]Binding{},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[3].StatementListItem.Statement.BlockStatement: bscope,
				}
				scope.Bindings = map[string][]Binding{
					"globalThis": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.MemberExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
					"window": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[3].StatementListItem.Statement.BlockStatement.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"a": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       m.ModuleListItems[2].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[3].StatementListItem.Statement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 7
			`{a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         map[javascript.Type]*Scope{},
					Bindings:       map[string][]Binding{},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[0].StatementListItem.Statement.BlockStatement: bscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 8
			`function b() {a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: map[javascript.Type]*Scope{},
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"b": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 9
			"let aValue = 1;{let bValue = 2;{aValue = 3}}",
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         map[javascript.Type]*Scope{},
					Bindings:       map[string][]Binding{},
				}
				bbscope := &Scope{
					IsLexicalScope: true,
					Parent:         bscope,
					Scopes:         map[javascript.Type]*Scope{},
					Bindings:       map[string][]Binding{},
				}
				bscope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[1].StatementListItem.Statement.BlockStatement.StatementList[1].Statement.BlockStatement: bbscope,
				}
				bscope.Bindings = map[string][]Binding{
					"bValue": {
						{
							BindingType: BindingLexicalLet,
							Scope:       bscope,
							Token:       m.ModuleListItems[1].StatementListItem.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[1].StatementListItem.Statement.BlockStatement: bscope,
				}
				scope.Bindings = map[string][]Binding{
					"aValue": {
						{
							BindingType: BindingLexicalLet,
							Scope:       scope,
							Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
						{
							BindingType: BindingBare,
							Scope:       bbscope,
							Token:       m.ModuleListItems[1].StatementListItem.Statement.BlockStatement.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 10
			"const aFunc = function b() {}",
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				fscope := &Scope{
					Parent: scope,
					Scopes: map[javascript.Type]*Scope{},
					Bindings: map[string][]Binding{
						"this":      {},
						"arguments": {},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.FunctionDeclaration): fscope,
				}
				scope.Bindings = map[string][]Binding{
					"aFunc": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
		{ // 10
			"const aClass = class b {}",
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"aClass": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
				}
				return scope, nil
			},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		source, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else {
			tscope, terr := test.Output(source)
			scope, err := ModuleScope(source, nil)
			if terr != nil && err != nil {
				if !errors.Is(terr, err) {
					t.Errorf("test %d: expecting error: %v\ngot: %v", n+1, terr, err)
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
