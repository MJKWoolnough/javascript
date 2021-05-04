// Package scope parses out a scope tree for a javascript module or script
package scope // import "vimagination.zapto.org/javascript/scope"

import (
	"errors"

	"vimagination.zapto.org/javascript"
)

// Binding represents a single instance of a bound name
type Binding struct {
	*Scope
	*javascript.Token
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

func (s *Scope) setBinding(name string, binding Binding) error {
	if _, ok := s.Bindings[name]; ok {
		return ErrDuplicateBinding
	}
	s.Bindings[name] = []Binding{binding}
	return nil
}

func (s *Scope) addBinding(name string, binding Binding) error {
	for {
		if bs, ok := s.Bindings[name]; ok {
			s.Bindings[name] = append(bs, binding)
			return nil
		}
		if s.Parent == nil {
			return s.setBinding(name, binding)
		}
		s = s.Parent
	}
}

// NewScope returns a init'd Scope type
func NewScope() *Scope {
	return &Scope{
		Bindings: make(map[string][]Binding),
	}
}

func newFunctionScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Bindings: map[string][]Binding{
			"this":      []Binding{},
			"arguments": []Binding{},
		},
	}
}

// ModuleScope parses out the scope tree for a javascript Module
func ModuleScope(m *javascript.Module, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}
	for _, i := range m.ModuleListItems {
		if i.ImportDeclaration != nil && i.ImportDeclaration.ImportClause != nil {
			if i.ImportDeclaration.ImportedDefaultBinding != nil {
				if err := global.setBinding(i.ImportDeclaration.ImportedDefaultBinding.Data, Binding{Token: i.ImportDeclaration.ImportedDefaultBinding, Scope: global}); err != nil {
					return nil, err
				}
			}
			if i.ImportDeclaration.NameSpaceImport != nil {
				if err := global.setBinding(i.ImportDeclaration.NameSpaceImport.Data, Binding{Token: i.ImportDeclaration.NameSpaceImport, Scope: global}); err != nil {
					return nil, err
				}
			}
			if i.ImportDeclaration.NamedImports != nil {
				for _, is := range i.ImportDeclaration.NamedImports.ImportList {
					if is.IdentifierName == nil {
						return nil, ErrInvalidImport
					}
					name := is.IdentifierName.Data
					if is.ImportedBinding != nil {
						name = is.ImportedBinding.Data
					}
					if err := global.setBinding(name, Binding{Token: is.IdentifierName, Scope: global}); err != nil {
						return nil, err
					}
				}
			}
		} else if i.StatementListItem != nil {
			if err := processStatementListItem(i.StatementListItem, global); err != nil {
				return nil, err
			}
		} else if i.ExportDeclaration != nil {

		}
	}
	return global, nil
}

// ScriptScope parses out the scope tree for a javascript script
func ScriptScope(s *javascript.Script, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}
	for _, i := range s.StatementList {
		if err := processStatementListItem(&i, global); err != nil {
			return nil, err
		}
	}
	return global, nil
}

func processStatementListItem(s *javascript.StatementListItem, scope *Scope) error {
	if s.Statement != nil {
		return processStatement(s.Statement, scope)
	} else if s.Declaration != nil {
		return processDeclaration(s.Declaration, scope)
	}
	return nil
}

func processStatement(s *javascript.Statement, scope *Scope) error {
	return nil
}

func processDeclaration(d *javascript.Declaration, scope *Scope) error {
	return nil
}

// Errors
var (
	ErrDuplicateBinding = errors.New("duplicate binding")
	ErrInvalidImport    = errors.New("invalid import")
)
