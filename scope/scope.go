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

func newLexicalScope(parent *Scope) *Scope {
	return &Scope{
		IsLexicalScope: true,
		Bindings:       make(map[string][]Binding),
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
	if s.BlockStatement != nil {
		return processBlockStatement(s.BlockStatement, scope)
	} else if s.VariableStatement != nil {
		return processVariableStatement(s.VariableStatement, scope)
	} else if s.ExpressionStatement != nil {
		return processExpressionStatement(s.ExpressionStatement, scope)
	} else if s.IfStatement != nil {
		return processIfStatement(s.IfStatement, scope)
	} else if s.IterationStatementDo != nil {
		return processIterationStatementDo(s.IterationStatementDo, scope)
	} else if s.IterationStatementWhile != nil {
		return processIterationStatementWhile(s.IterationStatementWhile, scope)
	} else if s.IterationStatementFor != nil {
		return processIterationStatementFor(s.IterationStatementFor, scope)
	} else if s.SwitchStatement != nil {
		return processSwitchStatement(s.SwitchStatement, scope)
	} else if s.WithStatement != nil {
		return processWithStatement(s.WithStatement, scope)
	} else if s.LabelIdentifier != nil {
		if s.LabelledItemFunction != nil {
			return processFunctionDeclaration(s.LabelledItemFunction, scope)
		} else if s.LabelledItemStatement != nil {
			return processStatement(s.LabelledItemStatement, scope)
		}
	} else if s.TryStatement != nil {
		return processTryStatement(s.TryStatement, scope)
	}
	return nil
}

func processDeclaration(d *javascript.Declaration, scope *Scope) error {
	if d.ClassDeclaration != nil {
		if err := processClassDeclaration(d.ClassDeclaration, scope*Scope); err != nil {
			return err
		}
	} else if d.FunctionDeclaration != nil {
		if err := processFunctionDeclaration(d.FunctionDeclaration, scope*Scope); err != nil {
			return err
		}
	} else if d.LexicalDeclaration != nil {
		if err := processLexicalDeclaration(d.LexicalDeclaration, scope*Scope); err != nil {
			return err
		}
	}
	return nil
}

func processBlockStatement(d *javascript.Block, scope *Scope) error {
	return nil
}

func processVariableStatement(d *javascript.VariableStatement, scope *Scope) error {
	return nil
}

func processExpressionStatement(d *javascript.Expression, scope *Scope) error {
	return nil
}

func processIfStatement(d *javascript.IfStatement, scope *Scope) error {
	return nil
}

func processIterationStatementDo(d *javascript.IterationStatementDo, scope *Scope) error {
	return nil
}

func processIterationStatementWhile(d *javascript.IterationStatementWhile, scope *Scope) error {
	return nil
}

func processIterationStatementFor(d *javascript.IterationStatementFor, scope *Scope) error {
	return nil
}

func processSwitchStatement(d *javascript.SwitchStatement, scope *Scope) error {
	return nil
}

func processWithStatement(d *javascript.WithStatement, scope *Scope) error {
	return nil
}

func processFunctionDeclaration(d *javascript.FunctionDeclaration, scope *Scope) error {
	return nil
}

func processTryStatement(d *javascript.TryStatement, scope *Scope) error {
	return nil
}

func processClassDeclaration(d *javascript.ClassDeclaration, scope *Scope) error {
	return nil
}

func processLexicalDeclaration(d *javascript.LexicalDeclaration, scope *Scope) error {
	return nil
}

// Errors
var (
	ErrDuplicateBinding = errors.New("duplicate binding")
	ErrInvalidImport    = errors.New("invalid import")
)
