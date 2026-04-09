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

var indent = []byte{'\t'}

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
			} else if _, err = i.Writer.Write(indent); err != nil {
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

func (i *indentPrinter) Print(args ...any) {
	fmt.Fprint(i, args...)
}

func (i *indentPrinter) Printf(format string, args ...any) {
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

	pp.Printf("\nIsLexicalScope: %v", s.IsLexicalScope)

	if s.Scopes == nil {
		pp.WriteString("\nScopes: nil")
	} else if len(s.Scopes) == 0 {
		pp.WriteString("\nScopes: []")
	} else {
		pp.WriteString("\nScopes: [")

		for t, scope := range s.Scopes {
			qq.WriteString("\n")
			qq.Printf("%p", t)
			qq.WriteString(": ")
			scope.printScope(qq)
		}

		pp.WriteString("\n]")
	}

	if s.Bindings == nil {
		pp.WriteString("\nBindings: nil")
	} else if len(s.Bindings) == 0 {
		pp.WriteString("\nBindings: []")
	} else {

		pp.WriteString("\nBindings: [")

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
				rr.Printf("%+s", binding.Token)
				rr.WriteString("\n]")
			}

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
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.ArrayLiteral).ElementList[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
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
							BindingType: BindingBare,
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
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ConditionalExpression).(*javascript.ArrayLiteral).ElementList[1].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
						{
							BindingType: BindingBare,
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
							BindingType: BindingBare,
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
							BindingType: BindingBare,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[1].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 37
			`let {a, b} = {};({a, b = 1} = {})`,
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
							BindingType: BindingBare,
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
							BindingType: BindingBare,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[1].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[1].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 38
			`var c;{let c;{var c}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[1].Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[1].Statement.BlockStatement.StatementList[1].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
				}
			},
		},
		{ // 39
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
		{ // 40
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
		{ // 41
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
		{ // 42
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
		{ // 43
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
		{ // 44
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
		{ // 45
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
		{ // 46
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
		{ // 47
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
		{ // 48
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
		{ // 49
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
		{ // 50
			`function a() {var a;try{}catch(a){let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchParameterBindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.TryStatement.CatchBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 51
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
		{ // 52
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
		{ // 53
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
		{ // 54
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
		{ // 55
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
		{ // 56
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
		{ // 57
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
		{ // 58
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
		{ // 59
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
		{ // 60
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
		{ // 61
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
		{ // 62
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
		{ // 63
			`function a() {let a;switch(a){case 1:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[1].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 64
			`function a() {let a;switch(a){case 1:let a;break;default:break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 65
			`function a() {let a;switch(a){default:break;case 1:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[1].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 66
			`function a() {let a;switch(a){default:let a;break;case 2:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.DefaultClause[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.PostDefaultCaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 67
			`function a() {let a;switch(a){case 1:let a;break;default:let a}}`,
			func(s *javascript.Script) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.CaseClauses[0].StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   s.StatementList[0].Declaration.FunctionDeclaration.FunctionBody.StatementList[1].Statement.SwitchStatement.DefaultClause[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 68
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
		{ // 69
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
		{ // 70
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
		{ // 71
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
		{ // 72
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
		{ // 73
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
		{ // 74
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
		{ // 75
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
		{ // 76
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
		{ // 77
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
		{ // 78
			`const a = <b />`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 79
			`const a = <b:c />`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"c": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 80
			`const a = <b.c />`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 81
			`const a = <b c={d}/>`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
					"d": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).Attributes[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 82
			`const a = <b {...d}/>`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
					"d": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).Attributes[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 83
			`const a = <b><c /></b>`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
					"c": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).Children[0].JSXElement.ElementName.Identifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 84
			`const a = <b>{c}{...d}</b>`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).ElementName.Identifier,
						},
					},
					"c": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).Children[0].JSXChildExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"d": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXElement).Children[1].JSXChildExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 85
			`const a = <><b></b></>`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingLexicalConst,
							Scope:       scope,
							Token:       s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
						},
					},
					"b": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(s.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.JSXFragment).Children[0].JSXElement.ElementName.Identifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 86
			`try{}catch(a){var a}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         make(map[javascript.Type]*Scope),
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
						{
							BindingType: BindingVar,
							Scope:       tscope,
							Token:       s.StatementList[0].Statement.TryStatement.CatchBlock.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         scope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Statement.TryStatement.CatchBlock: tscope,
				}

				return scope, nil
			},
		},
		{ // 87
			`try{}catch(a){{var a}}`,
			func(s *javascript.Script) (*Scope, error) {
				scope := NewScope()
				tscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
				}
				bscope := &Scope{
					IsLexicalScope: true,
					Scopes:         make(map[javascript.Type]*Scope),
					Parent:         tscope,
				}
				tscope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingCatch,
							Scope:       tscope,
							Token:       s.StatementList[0].Statement.TryStatement.CatchParameterBindingIdentifier,
						},
						{
							BindingType: BindingVar,
							Scope:       bscope,
							Token:       s.StatementList[0].Statement.TryStatement.CatchBlock.StatementList[0].Statement.BlockStatement.StatementList[0].Statement.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
					},
				}
				bscope.Bindings = map[string][]Binding{
					"a": {tscope.Bindings["a"][1]},
				}
				tscope.Scopes = map[javascript.Type]*Scope{
					s.StatementList[0].Statement.TryStatement.CatchBlock.StatementList[0].Statement.BlockStatement: bscope,
				}
				scope.Scopes = map[javascript.Type]*Scope{
					&s.StatementList[0].Statement.TryStatement.TryBlock: {
						IsLexicalScope: true,
						Parent:         scope,
						Scopes:         make(map[javascript.Type]*Scope),
						Bindings:       make(map[string][]Binding),
					},
					s.StatementList[0].Statement.TryStatement.CatchBlock: tscope,
				}

				return scope, nil
			},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		source, err := javascript.ParseScript(javascript.AsJSX(&tk))
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
			`import * as a from './lib.js';export {a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingImport,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ImportDeclaration.NameSpaceImport,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       m.ModuleListItems[1].ExportDeclaration.ExportClause.ExportList[0].IdentifierName,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 3
			`import a from './lib.js';export {a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingImport,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       m.ModuleListItems[1].ExportDeclaration.ExportClause.ExportList[0].IdentifierName,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 4
			`export {default as a} from './lib.js';`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"default": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.ExportClause.ExportList[0].IdentifierName,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 5
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
		{ // 6
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
		{ // 7
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
		{ // 8
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
		{ // 9
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
		{ // 10
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
		{ // 11
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
		{ // 12
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
		{ // 13
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
		{ // 14
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
		{ // 15
			`export var a;a`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingVar,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 16
			`export class MyClass {}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"MyClass": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.Declaration.ClassDeclaration.BindingIdentifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 17
			`export function MyFunc() {MyFunc()}`,
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
					m.ModuleListItems[0].ExportDeclaration.Declaration.FunctionDeclaration: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"MyFunc": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].ExportDeclaration.Declaration.FunctionDeclaration.BindingIdentifier,
						},
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].ExportDeclaration.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 18
			`export default a`,
			func(m *javascript.Module) (*Scope, error) {
				scope := &Scope{
					Scopes: make(map[javascript.Type]*Scope),
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingRef,
							Scope:       scope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].ExportDeclaration.DefaultAssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 19
			`const a = 1, a = 2;`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 20
			`import a from './b';import a from './c';`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ImportDeclaration.ImportedDefaultBinding,
				}
			},
		},
		{ // 21
			`import a from './b';import * as a from './c';`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ImportDeclaration.NameSpaceImport,
				}
			},
		},
		{ // 22
			`import a from './b';import {a} from './c';`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ImportDeclaration.NamedImports.ImportList[0].ImportedBinding,
				}
			},
		},
		{ // 23
			`import a from './b';export var a;`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ExportDeclaration.VariableStatement.VariableDeclarationList[0].BindingIdentifier,
				}
			},
		},
		{ // 24
			`import a from './b';export let a;`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ExportDeclaration.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
				}
			},
		},
		{ // 25
			`import a from './b';export default function a(){};`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ExportDeclaration.DefaultFunction.BindingIdentifier,
				}
			},
		},
		{ // 26
			`import a from './b';export default class a{};`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ImportDeclaration.ImportedDefaultBinding,
					Duplicate:   m.ModuleListItems[1].ExportDeclaration.DefaultClass.BindingIdentifier,
				}
			},
		},
		{ // 27
			`export default () => {const a, a};`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].ExportDeclaration.DefaultAssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].ExportDeclaration.DefaultAssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 28
			`let a = () => b = false, b = true; do a();while(b)`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				ascope := NewScope()
				ascope.Parent = scope
				scope.Scopes[m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction] = ascope
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementDo.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
					},
					{
						BindingType: BindingBare,
						Scope:       ascope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementDo.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 29
			`do {let a, a} while(1)`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementDo.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementDo.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 30
			`do a; while((() => {let a, a})())`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IterationStatementDo.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IterationStatementDo.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 31
			`let a = () => b = false, b = true; while(b) a()`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				ascope := NewScope()
				ascope.Parent = scope
				scope.Scopes[m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction] = ascope
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementWhile.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
					},
					{
						BindingType: BindingBare,
						Scope:       ascope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementWhile.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 32
			`while (1) {let a, a}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementWhile.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementWhile.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 33
			`while((() => {let a, a})()) a`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IterationStatementWhile.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IterationStatementWhile.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 34
			`let a = () => b = false, b = true; with(b) a()`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				ascope := NewScope()
				ascope.Parent = scope
				scope.Scopes[m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction] = ascope
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.WithStatement.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
					},
					{
						BindingType: BindingBare,
						Scope:       ascope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.WithStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 35
			`with (1) {let a, a}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.WithStatement.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.WithStatement.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 36
			`with((() => {let a, a})()) a`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.WithStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.WithStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 37
			`c: function b() {a}`,
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
					m.ModuleListItems[0].StatementListItem.Statement.LabelledItemFunction: fscope,
				}
				scope.Bindings = map[string][]Binding{
					"a": {
						{
							BindingType: BindingRef,
							Scope:       fscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.LabelledItemFunction.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
					"b": {
						{
							BindingType: BindingHoistable,
							Scope:       scope,
							Token:       m.ModuleListItems[0].StatementListItem.Statement.LabelledItemFunction.BindingIdentifier,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 38
			`a: {b}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := new(Scope)
				bscope := &Scope{
					IsLexicalScope: true,
					Parent:         scope,
					Scopes:         map[javascript.Type]*Scope{},
					Bindings:       map[string][]Binding{},
				}
				scope.Scopes = map[javascript.Type]*Scope{
					m.ModuleListItems[0].StatementListItem.Statement.LabelledItemStatement.BlockStatement: bscope,
				}
				scope.Bindings = map[string][]Binding{
					"b": {
						{
							BindingType: BindingRef,
							Scope:       bscope,
							Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.LabelledItemStatement.BlockStatement.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
						},
					},
				}

				return scope, nil
			},
		},
		{ // 39
			`class a extends (() => {let a, a}) {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassHeritage.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassHeritage.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 40
			`class a {b(){let c, c}}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 41
			`let a = 1, b = 2;if (a) a = 2; else {let b = a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				bscope := NewScope()
				bscope.Parent = scope
				bscope.IsLexicalScope = true
				scope.Scopes[m.ModuleListItems[1].StatementListItem.Statement.IfStatement.ElseStatement.BlockStatement] = bscope
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IfStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
					{
						BindingType: BindingBare,
						Scope:       scope,
						Token:       m.ModuleListItems[1].StatementListItem.Statement.IfStatement.Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       bscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IfStatement.ElseStatement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
					},
				}
				bscope.Bindings["b"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       bscope,
						Token:       m.ModuleListItems[1].StatementListItem.Statement.IfStatement.ElseStatement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
				}

				return scope, nil
			},
		},
		{ // 42
			`if (1) {let a, a}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IfStatement.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IfStatement.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 43
			`if (1) {} else {let a, a}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IfStatement.ElseStatement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IfStatement.ElseStatement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 44
			`if((() => {let a, a})()) a`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IfStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.IfStatement.Expression.Expressions[0].ConditionalExpression).(*javascript.CallExpression).MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 45
			`function a(b, b) {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[1].SingleNameBinding,
				}
			},
		},
		{ // 46
			`for (var a = () => {const b, b};;){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitVar[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitVar[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 47
			`for (let a, a;;){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitLexical.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitLexical.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 48
			`for (a = () => {const b, b};;){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitExpression.Expressions[0].AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.InitExpression.Expressions[0].AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 49
			`for ((()=>{let a,a}) of b){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 50
			`const a = [];for (b of a){}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				lscope := scope.newLexicalScope(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementFor)
				lscope.newLexicalScope(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementFor.Statement.BlockStatement)
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalConst,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       lscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.IterationStatementFor.Of.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       lscope,
						Token:       m.ModuleListItems[1].StatementListItem.Statement.IterationStatementFor.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 51
			`for (const {a, a} of b){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.ForBindingPatternObject.BindingPropertyList[0].BindingElement.SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.ForBindingPatternObject.BindingPropertyList[1].BindingElement.SingleNameBinding,
				}
			},
		},
		{ // 52
			`for (const [a, a] of b){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.ForBindingPatternArray.BindingElementList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.ForBindingPatternArray.BindingElementList[1].SingleNameBinding,
				}
			},
		},
		{ // 53
			`let a;for (var a of b){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[1].StatementListItem.Statement.IterationStatementFor.ForBindingIdentifier,
				}
			},
		},
		{ // 54
			`for (; () => {const a, a};){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Conditional.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Conditional.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 55
			`for (;;() => {const a, a}){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Afterthought.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Afterthought.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 56
			`for (a in () => {const b, b}){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.In.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.In.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 57
			`for (a of () => {const b, b}){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Of.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Of.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 58
			`for (a of b){let c, c}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.IterationStatementFor.Statement.BlockStatement.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 59
			`switch(() => {const a, a}){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 60
			`switch(a) {case () => {const a, a}:}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.CaseClauses[0].Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.CaseClauses[0].Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 61
			`switch(a) {default:case () => {const a, a}:}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.PostDefaultCaseClauses[0].Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.SwitchStatement.PostDefaultCaseClauses[0].Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 62
			`try{let a, a}finally{}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.TryStatement.TryBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.TryStatement.TryBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 63
			`try{}catch([a]){a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.newLexicalScope(&m.ModuleListItems[0].StatementListItem.Statement.TryStatement.TryBlock)
				lscope := scope.newLexicalScope(m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchBlock)
				lscope.Bindings["a"] = []Binding{
					{
						BindingType: BindingCatch,
						Scope:       lscope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterArrayBindingPattern.BindingElementList[0].SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       lscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchBlock.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 64
			`try{}catch({a}){a}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.newLexicalScope(&m.ModuleListItems[0].StatementListItem.Statement.TryStatement.TryBlock)
				lscope := scope.newLexicalScope(m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchBlock)
				lscope.Bindings["a"] = []Binding{
					{
						BindingType: BindingCatch,
						Scope:       lscope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       lscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchBlock.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 65
			`try{}catch([a, a]){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterArrayBindingPattern.BindingElementList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterArrayBindingPattern.BindingElementList[1].SingleNameBinding,
				}
			},
		},
		{ // 66
			`try{}catch({a, a}){}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.TryStatement.CatchParameterObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
				}
			},
		},
		{ // 67
			`try{}finally{let a, a}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.TryStatement.FinallyBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.TryStatement.FinallyBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 68
			`class a{b = () => {let c, c}}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].FieldDefinition.Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].FieldDefinition.Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 69
			`let a;class b{static {let c = a}}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				sscope := scope.newLexicalScope(m.ModuleListItems[1].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].ClassStaticBlock)
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
					{
						BindingType: BindingRef,
						Scope:       sscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].ClassStaticBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].Initializer.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingHoistable,
						Scope:       scope,
						Token:       m.ModuleListItems[1].StatementListItem.Declaration.ClassDeclaration.BindingIdentifier,
					},
				}
				sscope.Bindings["c"] = []Binding{
					{
						BindingType: BindingLexicalLet,
						Scope:       sscope,
						Token:       m.ModuleListItems[1].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].ClassStaticBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					},
				}

				return scope, nil
			},
		},
		{ // 70
			`class a{static {let b, b}}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].ClassStaticBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].ClassStaticBlock.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 71
			`class a{[() => {let b, b}]}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].FieldDefinition.ClassElementName.PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].FieldDefinition.ClassElementName.PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 72
			`var {a} = {}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingVar,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					},
				}

				return scope, nil
			},
		},
		{ // 73
			`let a;var [a] = {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[1].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
				}
			},
		},
		{ // 74
			`let a;var {a} = {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[1].StatementListItem.Statement.VariableStatement.VariableDeclarationList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
				}
			},
		},
		{ // 75
			`a.b = [];a`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 76
			`a[() => {let b, b}] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 77
			`[a[() => {let b, b}]] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 78
			`({a = () => {let b, b}} = {})`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 79
			`a?.b[() => {let c,c}]`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalChain.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalChain.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 80
			`({a, ...b} = c); b`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingBare,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].PropertyName.LiteralPropertyName,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingBare,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["c"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 81
			`({a, ...b[() => {let b, b}]} = c)`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 82
			`({a: b[() => {let b, b}]} = c)`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.ParenthesizedExpression).Expressions[0].AssignmentPattern.ObjectAssignmentPattern.AssignmentPropertyList[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 83
			`[[a]] = b;a`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingBare,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[1].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 84
			`[[a[() => {let b, b}]]] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 85
			`[a = b] = []`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingBare,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].Initializer.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 86
			`[a = () => {let b, b}] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentElements[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 87
			`[...a] = []`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 88
			`[...a[() => {let b, b}]] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].AssignmentPattern.ArrayAssignmentPattern.AssignmentRestElement.NewExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 89
			`let {a, ...a} = {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingRestProperty,
				}
			},
		},
		{ // 90
			`const [a, ...a] = []`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ArrayBindingPattern.BindingRestElement.SingleNameBinding,
				}
			},
		},
		{ // 91
			`function a(b, ...[c, d]) {return b, c}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				fscope := scope.newFunctionScope(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration)
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingHoistable,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.BindingIdentifier,
					},
				}
				fscope.Bindings["b"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       fscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				fscope.Bindings["c"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       fscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[1].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				fscope.Bindings["d"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ArrayBindingPattern.BindingElementList[1].SingleNameBinding,
					},
				}

				return scope, nil
			},
		},
		{ // 92
			`function a(b, ...{c, d}) {return b, c}`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				fscope := scope.newFunctionScope(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration)
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingHoistable,
						Scope:       scope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.BindingIdentifier,
					},
				}
				fscope.Bindings["b"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       fscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				fscope.Bindings["c"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					},
					{
						BindingType: BindingRef,
						Scope:       fscope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FunctionBody.StatementList[0].Statement.ExpressionStatement.Expressions[1].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				fscope.Bindings["d"] = []Binding{
					{
						BindingType: BindingFunctionParam,
						Scope:       fscope,
						Token:       m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
					},
				}

				return scope, nil
			},
		},
		{ // 93
			`function a(b, ...[c, c]) {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ArrayBindingPattern.BindingElementList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ArrayBindingPattern.BindingElementList[1].SingleNameBinding,
				}
			},
		},
		{ // 94
			`function a(b, ...{c, c}) {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ObjectBindingPattern.BindingPropertyList[0].BindingElement.SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.ObjectBindingPattern.BindingPropertyList[1].BindingElement.SingleNameBinding,
				}
			},
		},
		{ // 95
			`function a(b, ...b) {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.FormalParameterList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.FunctionDeclaration.FormalParameters.BindingIdentifier,
				}
			},
		},
		{ // 96
			`class a {[() => {let b, b}](){}}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.ClassElementName.PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.ClassElementName.PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 97
			`class a {b(c, c){}}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.Params.FormalParameterList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.ClassDeclaration.ClassBody[0].MethodDefinition.Params.FormalParameterList[1].SingleNameBinding,
				}
			},
		},
		{ // 98
			`let a = () => {let b, b}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 99
			`a ?? b`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.WrapConditional(&m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.CoalesceExpressionHead.BitwiseORExpression)).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.WrapConditional(&m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.BitwiseORExpression)).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 100
			`a ?? (() => {let b, b})`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(javascript.WrapConditional(&m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.BitwiseORExpression)).(*javascript.ParenthesizedExpression).Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(javascript.WrapConditional(&m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.CoalesceExpression.BitwiseORExpression)).(*javascript.ParenthesizedExpression).Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 101
			`a ? b : c`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.WrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.LogicalORExpression)).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.True.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["c"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.False.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 102
			`a ? () => {let b, b}: c`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.True.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.True.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 103
			`a ? b: () => {let b, b}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.False.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression.False.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 104
			`let a = (b, b) => {}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.FormalParameters.FormalParameterList[0].SingleNameBinding,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.FormalParameters.FormalParameterList[1].SingleNameBinding,
				}
			},
		},
		{ // 105
			`let a = () => () => {let b,b}`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].Initializer.ArrowFunction.AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 106
			`import(a)`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).ImportCall.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 107
			`import(() => {let a, a})`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).ImportCall.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).ImportCall.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 108
			`a(b)[c]`,
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).CallExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).CallExpression.Arguments.ArgumentList[0].AssignmentExpression.ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}
				scope.Bindings["c"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 109
			`a(()=>{let b, b})[c]`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).CallExpression.Arguments.ArgumentList[0].AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).CallExpression.Arguments.ArgumentList[0].AssignmentExpression.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 110
			`a(b)[()=>{let c, c}]`,
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 111
			"a()`${b}`",
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).CallExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).TemplateLiteral.Expressions[0].Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 112
			"a()`${() => {let b, b}}`",
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).TemplateLiteral.Expressions[0].Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.CallExpression).TemplateLiteral.Expressions[0].Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 113
			"let {[() => {let a, a}]: b} = c",
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   m.ModuleListItems[0].StatementListItem.Declaration.LexicalDeclaration.BindingList[0].ObjectBindingPattern.BindingPropertyList[0].PropertyName.ComputedPropertyName.ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 114
			"a[() => {let b, b}]?.c",
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 115
			"a()?.[b]",
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).CallExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}
				scope.Bindings["b"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalChain.Expression.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 116
			"a()[() => {let b, b}]?.c",
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).CallExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).CallExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
			},
		},
		{ // 117
			"a?.b?.c",
			func(m *javascript.Module) (*Scope, error) {
				scope := NewScope()
				scope.Bindings["a"] = []Binding{
					{
						BindingType: BindingRef,
						Scope:       scope,
						Token:       javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalExpression.MemberExpression.PrimaryExpression.IdentifierReference,
					},
				}

				return scope, nil
			},
		},
		{ // 118
			"a[() => {let b, b}]?.c?.d",
			func(m *javascript.Module) (*Scope, error) {
				return nil, ErrDuplicateDeclaration{
					Declaration: javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[0].BindingIdentifier,
					Duplicate:   javascript.UnwrapConditional(m.ModuleListItems[0].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.OptionalExpression).OptionalExpression.MemberExpression.Expression.Expressions[0].ArrowFunction.FunctionBody.StatementList[0].Declaration.LexicalDeclaration.BindingList[1].BindingIdentifier,
				}
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

func TestFindIdentifier(t *testing.T) {
	tk := parser.NewStringTokeniser(`const a; {let b}`)

	if source, err := javascript.ParseModule(&tk); err != nil {
		t.Errorf("unexpected error parsing script: %s", err)
	} else if scope, err := ModuleScope(source, nil); err != nil {
		t.Errorf("unexpected error determining scope: %s", err)
	} else if scope.FindIdentifier("a") != scope {
		t.Errorf("test 1: didn't get expected scope")
	} else if scope.FindIdentifier("b") != nil {
		t.Errorf("test 2: didn't get expected scope")
	} else if inner := scope.Scopes[source.ModuleListItems[1].StatementListItem.Statement.BlockStatement]; inner.FindIdentifier("a") != scope {
		t.Errorf("test 3: didn't get expected scope")
	} else if inner.FindIdentifier("b") != inner {
		t.Errorf("test 4: didn't get expected scope")
	}
}

func TestRename(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		From, To      string
		Scope         func(*javascript.Module, *Scope) *Scope
		Renamed       bool
	}{
		{ // 1
			Input:   `const a = 1;`,
			Output:  `const a = 1;`,
			From:    "a",
			To:      "a",
			Scope:   func(_ *javascript.Module, s *Scope) *Scope { return s },
			Renamed: false,
		},
		{ // 2
			Input:   `const a = 1, b = 2;`,
			Output:  `const a = 1, b = 2;`,
			From:    "a",
			To:      "b",
			Scope:   func(_ *javascript.Module, s *Scope) *Scope { return s },
			Renamed: false,
		},
		{ // 3
			Input:   `const a = 1;`,
			Output:  `const b = 1;`,
			From:    "a",
			To:      "b",
			Scope:   func(_ *javascript.Module, s *Scope) *Scope { return s },
			Renamed: true,
		},
		{ // 4
			Input:   `const a = 1;{a}`,
			Output:  "const b = 1;\n\n{\n\tb;\n}",
			From:    "a",
			To:      "b",
			Scope:   func(_ *javascript.Module, s *Scope) *Scope { return s },
			Renamed: true,
		},
		{ // 5
			Input:   `const a = 1;{let a; a = 2}`,
			Output:  "const b = 1;\n\n{\n\tlet a;\n\ta = 2;\n}",
			From:    "a",
			To:      "b",
			Scope:   func(_ *javascript.Module, s *Scope) *Scope { return s },
			Renamed: true,
		},
		{ // 6
			Input:  `const a = 1;{let a; a = 2}`,
			Output: "const a = 1;\n\n{\n\tlet b;\n\tb = 2;\n}",
			From:   "a",
			To:     "b",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement]
			},
			Renamed: true,
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else if scope, err := ModuleScope(m, nil); err != nil {
			t.Errorf("test %d: unexpected error determining scope: %s", n+1, err)
		} else if renamed := test.Scope(m, scope).Rename(test.From, test.To); renamed != test.Renamed {
			t.Errorf("test %d: expecting renamed = %v, got %v", n+1, test.Renamed, renamed)
		} else if output := fmt.Sprintf("%s", m); output != test.Output {
			t.Errorf("test %d: expecting output = %q, got %q", n+1, test.Output, output)
		}
	}
}

func TestIdentifierInUse(t *testing.T) {
	for n, test := range [...]struct {
		Input      string
		Identifier string
		Scope      func(*javascript.Module, *Scope) *Scope
		InUse      bool
	}{
		{ // 1
			Input:      `const a = 1;`,
			Identifier: "a",
			Scope:      func(_ *javascript.Module, s *Scope) *Scope { return s },
			InUse:      true,
		},
		{ // 2
			Input:      `const a = 1;`,
			Identifier: "b",
			Scope:      func(_ *javascript.Module, s *Scope) *Scope { return s },
			InUse:      false,
		},
		{ // 3
			Input:      `const a = 1;{b}`,
			Identifier: "a",
			Scope:      func(_ *javascript.Module, s *Scope) *Scope { return s },
			InUse:      true,
		},
		{ // 4
			Input:      `const a = 1;{b}`,
			Identifier: "a",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement]
			},
			InUse: false,
		},
		{ // 5
			Input:      `const a = 1;{a}`,
			Identifier: "a",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement]
			},
			InUse: true,
		},
		{ // 6
			Input:      `const a = 1;{let b;{a}}`,
			Identifier: "a",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement]
			},
			InUse: true,
		},
		{ // 7
			Input:      `const a = 1;{let b;{a}}`,
			Identifier: "b",
			Scope: func(_ *javascript.Module, s *Scope) *Scope {
				return s
			},
			InUse: false,
		},
		{ // 8
			Input:      `const a = 1;{let b;{a}}`,
			Identifier: "b",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement]
			},
			InUse: true,
		},
		{ // 9
			Input:      `const a = 1;{let b;{a}}`,
			Identifier: "b",
			Scope: func(m *javascript.Module, s *Scope) *Scope {
				return s.Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement].Scopes[m.ModuleListItems[1].StatementListItem.Statement.BlockStatement.StatementList[1].Statement.BlockStatement]
			},
			InUse: false,
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)

		m, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing script: %s", n+1, err)
		} else if scope, err := ModuleScope(m, nil); err != nil {
			t.Errorf("test %d: unexpected error determining scope: %s", n+1, err)
		} else if inUse := test.Scope(m, scope).IdentifierInUse(test.Identifier); inUse != test.InUse {
			t.Errorf("test %d: expecting inUse = %v, got %v", n+1, test.InUse, inUse)
		}
	}
}
