// Package scope parses out a scope tree for a javascript module or script
package scope // import "vimagination.zapto.org/javascript/scope"

import "vimagination.zapto.org/javascript"

// Binding represents a single instance of a bound name
type Binding struct {
	*Scope
	*Token
}

// Scope represents a single level of variable scope
type Scope struct {
	IsLexicalScope bool
	Parent         *Scope
	Scopes         []Scope
	Bindings       map[string][]Binding
}

func (s *Scope) getFunctionScope() *Scope {
	for s.IsLexicalScope && s.Parent != nil {
		s = s.Parent
	}
	return s
}

// NewScope returns a init'd Scope type
func NewScope() *Scope {
	return &Scope{
		Bindings: make(map[string]Binding),
	}
}

func newFunctionScope(parent *Scope) *Scope {
	return &Scope{
		Parent: &parent,
		Bindings: map[string]Binding{
			"this":      []Binding{},
			"arguments": []Binding{},
		},
	}
}

// ModuleScope parses out the scope tree for a javascript Module
func ModuleScope(m *javascript.Module, global *Scope) *Scope {
	if global == nil {
		global = NewScope()
	}
	for _, i := range m.ModuleListItems {
		if i.ImportDeclaration != nil {

		} else if i.StatementListItem != nil {
			processStatementListItem(i.StatementListItem, global)
		} else if i.ExportDeclaration != nil {

		}
	}
	return global
}

// ScriptScope parses out the scope tree for a javascript script
func ScriptScope(s *javascript.Script, global *Scope) *Scope {
	if global == nil {
		global = NewScope()
	}
	for _, i := range s.StatementList {
		processStatementListItem(i, global)
	}
	return global
}

func processStatementListItem(s *javascript.StatementListItem, scope *Scope) {
	if s.Statement != nil {
		processStatement(s.Statement, scope)
	} else if s.Declaration != nil {
		processDeclaration(s.Declaration, scope)
	}
}

func processStatement(s *javascript.Statement, scope *Scope) {

}

func processDeclaration(d *javascript.Declaration, scope *Scope) {

}
