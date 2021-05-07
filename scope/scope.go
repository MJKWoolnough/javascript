// Package scope parses out a scope tree for a javascript module or script
package scope // import "vimagination.zapto.org/javascript/scope"

import (
	"errors"
	"fmt"

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
	Scopes         map[fmt.Formatter]*Scope
	Bindings       map[string][]Binding
}

func (s *Scope) getFunctionScope() *Scope {
	for s.IsLexicalScope && s.Parent != nil {
		s = s.Parent
	}
	return s
}

func (s *Scope) setBinding(t *javascript.Token, hoist bool) error {
	name := t.Data
	if _, ok := s.Bindings[name]; ok {
		return ErrDuplicateBinding
	}
	binding := Binding{Token: t, Scope: s}
	s.Bindings[name] = []Binding{binding}
	if hoist && s.IsLexicalScope {
		s = s.getFunctionScope()
		if _, ok := s.Bindings[name]; ok {
			return ErrDuplicateBinding
		}
		s.Bindings[name] = []Binding{binding}
	}
	return nil
}

func (s *Scope) addBinding(t *javascript.Token) {
	name := t.Data
	binding := Binding{Token: t, Scope: s}
	for {
		if bs, ok := s.Bindings[name]; ok {
			s.Bindings[name] = append(bs, binding)
			return
		}
		if s.Parent == nil {
			s.Bindings[name] = []Binding{binding}
		}
		s = s.Parent
	}
}

// NewScope returns a init'd Scope type
func NewScope() *Scope {
	return &Scope{
		Scopes:   make(map[fmt.Formatter]*Scope),
		Bindings: make(map[string][]Binding),
	}
}

func (s *Scope) newFunctionScope(js fmt.Formatter) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}
	ns := &Scope{
		Parent: s,
		Scopes: make(map[fmt.Formatter]*Scope),
		Bindings: map[string][]Binding{
			"this":      []Binding{},
			"arguments": []Binding{},
		},
	}
	s.Scopes[js] = ns
	return ns
}

func (s *Scope) newArrowFunctionScope(js fmt.Formatter) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}
	ns := &Scope{
		Parent:   s,
		Scopes:   make(map[fmt.Formatter]*Scope),
		Bindings: make(map[string][]Binding),
	}
	s.Scopes[js] = ns
	return ns
}

func (s *Scope) newLexicalScope(js fmt.Formatter) *Scope {
	if ns, ok := s.Scopes[js]; ok {
		return ns
	}
	ns := &Scope{
		Parent:         s,
		IsLexicalScope: true,
		Scopes:         make(map[fmt.Formatter]*Scope),
		Bindings:       make(map[string][]Binding),
	}
	s.Scopes[js] = ns
	return ns
}

// ModuleScope parses out the scope tree for a javascript Module
func ModuleScope(m *javascript.Module, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}
	if err := processModule(m, global, true); err != nil {
		return nil, err
	}
	processModule(m, global, false)
	return global, nil
}

// ScriptScope parses out the scope tree for a javascript script
func ScriptScope(s *javascript.Script, global *Scope) (*Scope, error) {
	if global == nil {
		global = NewScope()
	}
	for _, i := range s.StatementList {
		if err := processStatementListItem(&i, global, true); err != nil {
			return nil, err
		}
	}
	for _, i := range s.StatementList {
		processStatementListItem(&i, global, false)
	}
	return global, nil
}

func processModule(m *javascript.Module, global *Scope, set bool) error {
	for _, i := range m.ModuleListItems {
		if i.ImportDeclaration != nil && i.ImportDeclaration.ImportClause != nil {
			if set {
				if i.ImportDeclaration.ImportedDefaultBinding != nil {
					if err := global.setBinding(i.ImportDeclaration.ImportedDefaultBinding, false); err != nil {
						return err
					}
				}
				if i.ImportDeclaration.NameSpaceImport != nil {
					if err := global.setBinding(i.ImportDeclaration.NameSpaceImport, false); err != nil {
						return err
					}
				}
				if i.ImportDeclaration.NamedImports != nil {
					for _, is := range i.ImportDeclaration.NamedImports.ImportList {
						if is.IdentifierName == nil {
							return ErrInvalidImport
						}
						var tk = is.IdentifierName
						if is.ImportedBinding != nil {
							tk = is.ImportedBinding
						}
						if err := global.setBinding(tk, false); err != nil {
							return err
						}
					}
				}
			}
		} else if i.StatementListItem != nil {
			if err := processStatementListItem(i.StatementListItem, global, set); err != nil {
				return err
			}
		} else if i.ExportDeclaration != nil {

		}
	}
	return nil
}

func processStatementListItem(s *javascript.StatementListItem, scope *Scope, set bool) error {
	if s.Statement != nil {
		return processStatement(s.Statement, scope, set)
	} else if s.Declaration != nil {
		return processDeclaration(s.Declaration, scope, set)
	}
	return nil
}

func processStatement(s *javascript.Statement, scope *Scope, set bool) error {
	if s.BlockStatement != nil {
		return processBlockStatement(s.BlockStatement, scope.newLexicalScope(s.BlockStatement), set)
	} else if s.VariableStatement != nil {
		return processVariableStatement(s.VariableStatement, scope, set)
	} else if s.ExpressionStatement != nil {
		return processExpression(s.ExpressionStatement, scope, set)
	} else if s.IfStatement != nil {
		return processIfStatement(s.IfStatement, scope, set)
	} else if s.IterationStatementDo != nil {
		return processIterationStatementDo(s.IterationStatementDo, scope, set)
	} else if s.IterationStatementWhile != nil {
		return processIterationStatementWhile(s.IterationStatementWhile, scope, set)
	} else if s.IterationStatementFor != nil {
		return processIterationStatementFor(s.IterationStatementFor, scope, set)
	} else if s.SwitchStatement != nil {
		return processSwitchStatement(s.SwitchStatement, scope, set)
	} else if s.WithStatement != nil {
		return processWithStatement(s.WithStatement, scope, set)
	} else if s.LabelIdentifier != nil {
		if s.LabelledItemFunction != nil {
			return processFunctionDeclaration(s.LabelledItemFunction, scope, set)
		} else if s.LabelledItemStatement != nil {
			return processStatement(s.LabelledItemStatement, scope, set)
		}
	} else if s.TryStatement != nil {
		return processTryStatement(s.TryStatement, scope, set)
	}
	return nil
}

func processDeclaration(d *javascript.Declaration, scope *Scope, set bool) error {
	if d.ClassDeclaration != nil {
		if err := processClassDeclaration(d.ClassDeclaration, scope, set); err != nil {
			return err
		}
	} else if d.FunctionDeclaration != nil {
		if err := processFunctionDeclaration(d.FunctionDeclaration, scope, set); err != nil {
			return err
		}
	} else if d.LexicalDeclaration != nil {
		if err := processLexicalDeclaration(d.LexicalDeclaration, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processBlockStatement(b *javascript.Block, scope *Scope, set bool) error {
	for _, sli := range b.StatementList {
		if err := processStatementListItem(&sli, scope, set); err != nil {
			return nil
		}
	}
	return nil
}

func processVariableStatement(v *javascript.VariableStatement, scope *Scope, set bool) error {
	for _, vs := range v.VariableDeclarationList {
		if err := processVariableDeclaration(vs, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processExpression(e *javascript.Expression, scope *Scope, set bool) error {
	for _, ae := range e.Expressions {
		if err := processAssignmentExpression(&ae, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processIfStatement(i *javascript.IfStatement, scope *Scope, set bool) error {
	if err := processExpression(&i.Expression, scope, set); err != nil {
		return err
	}
	if err := processStatement(&i.Statement, scope, set); err != nil {
		return err
	}
	if i.ElseStatement != nil {
		if err := processStatement(i.ElseStatement, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processIterationStatementDo(d *javascript.IterationStatementDo, scope *Scope, set bool) error {
	if err := processStatement(&d.Statement, scope, set); err != nil {
		return err
	}
	if err := processExpression(&d.Expression, scope, set); err != nil {
		return err
	}
	return nil
}

func processIterationStatementWhile(w *javascript.IterationStatementWhile, scope *Scope, set bool) error {
	if err := processExpression(&w.Expression, scope, set); err != nil {
		return err
	}
	if err := processStatement(&w.Statement, scope, set); err != nil {
		return err
	}
	return nil
}

func processIterationStatementFor(f *javascript.IterationStatementFor, scope *Scope, set bool) error {
	switch f.Type {
	case javascript.ForNormal:
	case javascript.ForNormalVar:
		for _, v := range f.InitVar {
			if err := processVariableDeclaration(v, scope, set); err != nil {
				return err
			}
		}
	case javascript.ForNormalLexicalDeclaration:
		if f.InitLexical != nil {
			scope = scope.newLexicalScope(f.InitLexical)
			if err := processLexicalDeclaration(f.InitLexical, scope, set); err != nil {
				return err
			}
		}
	case javascript.ForNormalExpression:
		if f.InitExpression != nil {
			if err := processExpression(f.InitExpression, scope, set); err != nil {
				return err
			}
		}
	case javascript.ForInLeftHandSide, javascript.ForOfLeftHandSide, javascript.ForAwaitOfLeftHandSide:
		if f.LeftHandSideExpression != nil {
			if err := processLeftHandSideExpression(f.LeftHandSideExpression, scope, set); err != nil {
				return err
			}
		}
	default:
		if f.ForBindingIdentifier != nil {
			if !set {
				scope.Parent.addBinding(f.ForBindingIdentifier)
			}
		} else if f.ForBindingPatternObject != nil && set {
			if err := processObjectBindingPattern(f.ForBindingPatternObject, scope, false, true); err != nil {
				return err
			}
		} else if f.ForBindingPatternArray != nil && set {
			if err := processArrayBindingPattern(f.ForBindingPatternArray, scope, false, true); err != nil {
				return err
			}
		}
	}
	switch f.Type {
	case javascript.ForNormal, javascript.ForNormalVar, javascript.ForNormalLexicalDeclaration, javascript.ForNormalExpression:
		if f.Conditional != nil {
			if err := processExpression(f.Conditional, scope, set); err != nil {
				return err
			}
		}
		if f.Afterthought != nil {
			if err := processExpression(f.Afterthought, scope, set); err != nil {
				return err
			}
		}
	case javascript.ForInLeftHandSide, javascript.ForInVar, javascript.ForInLet, javascript.ForInConst:
		if f.In != nil {
			if err := processExpression(f.In, scope, set); err != nil {
				return err
			}
		}
	case javascript.ForOfLeftHandSide, javascript.ForOfVar, javascript.ForOfLet, javascript.ForOfConst, javascript.ForAwaitOfLeftHandSide, javascript.ForAwaitOfVar, javascript.ForAwaitOfLet, javascript.ForAwaitOfConst:
		if f.Of != nil {
			if err := processAssignmentExpression(f.Of, scope, set); err != nil {
				return err
			}
		}
	}
	if err := processStatement(&f.Statement, scope, set); err != nil {
		return err
	}
	return nil
}

func processSwitchStatement(s *javascript.SwitchStatement, scope *Scope, set bool) error {
	if err := processExpression(&s.Expression, scope, set); err != nil {
		return err
	}
	scope = scope.newLexicalScope(s)
	for _, c := range s.CaseClauses {
		if err := processExpression(&c.Expression, scope, set); err != nil {
			return err
		}
		for _, sli := range c.StatementList {
			if err := processStatementListItem(&sli, scope, set); err != nil {
				return err
			}
		}
	}
	for _, sli := range s.DefaultClause {
		if err := processStatementListItem(&sli, scope, set); err != nil {
			return err
		}
	}
	for _, c := range s.PostDefaultCaseClauses {
		if err := processExpression(&c.Expression, scope, set); err != nil {
			return err
		}
		for _, sli := range c.StatementList {
			if err := processStatementListItem(&sli, scope, set); err != nil {
				return err
			}
		}
	}
	return nil
}

func processWithStatement(w *javascript.WithStatement, scope *Scope, set bool) error {
	processExpression(&w.Expression, scope, set)
	processStatement(&w.Statement, scope, set)
	return nil
}

func processFunctionDeclaration(f *javascript.FunctionDeclaration, scope *Scope, set bool) error {
	if f.BindingIdentifier != nil && set {
		if err := scope.setBinding(f.BindingIdentifier, true); err != nil {
			return err
		}
	}
	scope = scope.newFunctionScope(f)
	if err := processFormalParameters(&f.FormalParameters, scope, set); err != nil {
		return err
	}
	if err := processBlockStatement(&f.FunctionBody, scope, set); err != nil {
		return err
	}
	return nil
}

func processTryStatement(t *javascript.TryStatement, scope *Scope, set bool) error {
	if err := processBlockStatement(&t.TryBlock, scope.newLexicalScope(t.TryBlock), set); err != nil {
		return err
	}
	if t.CatchBlock != nil {
		if err := processBlockStatement(t.CatchBlock, scope.newLexicalScope(t.CatchBlock), set); err != nil {
			return err
		}
	}
	if t.FinallyBlock != nil {
		if err := processBlockStatement(t.FinallyBlock, scope.newLexicalScope(t.FinallyBlock), set); err != nil {
			return err
		}
	}
	return nil
}

func processClassDeclaration(c *javascript.ClassDeclaration, scope *Scope, set bool) error {
	if c.BindingIdentifier != nil && set {
		if err := scope.setBinding(c.BindingIdentifier, false); err != nil {
			return err
		}
	}
	if c.ClassHeritage != nil {
		if err := processLeftHandSideExpression(c.ClassHeritage, scope, set); err != nil {
			return err
		}
	}
	for _, md := range c.ClassBody {
		if err := processMethodDefinition(md, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processLexicalDeclaration(l *javascript.LexicalDeclaration, scope *Scope, set bool) error {
	for _, lb := range l.BindingList {
		if err := processLexicalBinding(lb, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processVariableDeclaration(v javascript.VariableDeclaration, scope *Scope, set bool) error {
	if set {
		if v.BindingIdentifier != nil {
			if err := scope.setBinding(v.BindingIdentifier, true); err != nil {
				return err
			}
		} else if v.ArrayBindingPattern != nil {
			if err := processArrayBindingPattern(v.ArrayBindingPattern, scope, true, false); err != nil {
				return err
			}
		} else if v.ObjectBindingPattern != nil {
			if err := processObjectBindingPattern(v.ObjectBindingPattern, scope, true, false); err != nil {
				return err
			}
		}
	}
	if v.Initializer != nil {
		if err := processAssignmentExpression(v.Initializer, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processAssignmentExpression(a *javascript.AssignmentExpression, scope *Scope, set bool) error {
	if a.ConditionalExpression != nil {
		if err := processConditionalExpression(a.ConditionalExpression, scope, set); err != nil {
			return err
		}
	} else if a.ArrowFunction != nil {
		if err := processArrowFunction(a.ArrowFunction, scope, set); err != nil {
			return err
		}
	} else if a.LeftHandSideExpression != nil {
		if err := processLeftHandSideExpression(a.LeftHandSideExpression, scope, set); err != nil {
			return err
		}
	}
	if a.AssignmentExpression != nil {
		if err := processAssignmentExpression(a.AssignmentExpression, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processLeftHandSideExpression(l *javascript.LeftHandSideExpression, scope *Scope, set bool) error {
	if l.NewExpression != nil {
		if err := processNewExpression(l.NewExpression, scope, set); err != nil {
			return err
		}
	} else if l.CallExpression != nil {
		if err := processCallExpression(l.CallExpression, scope, set); err != nil {
			return err
		}
	} else if l.OptionalExpression != nil {
		if err := processOptionalExpression(l.OptionalExpression, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processObjectBindingPattern(o *javascript.ObjectBindingPattern, scope *Scope, hoist, bare bool) error {
	for _, bp := range o.BindingPropertyList {
		if err := processBindingProperty(bp, scope, hoist, bare); err != nil {
			return err
		}
	}
	if o.BindingRestProperty != nil {
		if bare {
			scope.addBinding(o.BindingRestProperty)
		} else {
			scope.setBinding(o.BindingRestProperty, hoist)
		}
	}
	return nil
}

func processArrayBindingPattern(a *javascript.ArrayBindingPattern, scope *Scope, hoist, bare bool) error {
	for _, be := range a.BindingElementList {
		if err := processBindingElement(&be, scope, hoist, bare); err != nil {
			return err
		}
	}
	if a.BindingRestElement != nil {
		if err := processBindingElement(a.BindingRestElement, scope, hoist, bare); err != nil {
			return err
		}
	}
	return nil
}

func processFormalParameters(f *javascript.FormalParameters, scope *Scope, set bool) error {
	return nil
}

func processMethodDefinition(m javascript.MethodDefinition, scope *Scope, set bool) error {
	return nil
}

func processLexicalBinding(l javascript.LexicalBinding, scope *Scope, set bool) error {
	return nil
}

func processConditionalExpression(c *javascript.ConditionalExpression, scope *Scope, set bool) error {
	return nil
}

func processArrowFunction(a *javascript.ArrowFunction, scope *Scope, set bool) error {
	return nil
}

func processNewExpression(n *javascript.NewExpression, scope *Scope, set bool) error {
	return nil
}

func processCallExpression(c *javascript.CallExpression, scope *Scope, set bool) error {
	return nil
}

func processOptionalExpression(o *javascript.OptionalExpression, scope *Scope, set bool) error {
	return nil
}

func processBindingProperty(b javascript.BindingProperty, scope *Scope, hoist, bare bool) error {
	return nil
}

func processBindingElement(b *javascript.BindingElement, scope *Scope, hoist, bare bool) error {
	return nil
}

// Errors
var (
	ErrDuplicateBinding = errors.New("duplicate binding")
	ErrInvalidImport    = errors.New("invalid import")
)
