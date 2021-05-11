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

func (s *Scope) setBinding(t *javascript.Token, hoist bool) error {
	name := t.Data
	if _, ok := s.Bindings[name]; ok {
		return ErrDuplicateBinding
	}
	binding := Binding{Token: t, Scope: s}
	s.Bindings[name] = []Binding{binding}
	if hoist && s.IsLexicalScope {
		for s.IsLexicalScope && s.Parent != nil {
			s = s.Parent
		}
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
	for n := range s.StatementList {
		if err := processStatementListItem(&s.StatementList[n], global, true); err != nil {
			return nil, err
		}
	}
	for n := range s.StatementList {
		processStatementListItem(&s.StatementList[n], global, false)
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
			return processFunctionDeclaration(s.LabelledItemFunction, scope, set, false)
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
		if err := processClassDeclaration(d.ClassDeclaration, scope, set, false); err != nil {
			return err
		}
	} else if d.FunctionDeclaration != nil {
		if err := processFunctionDeclaration(d.FunctionDeclaration, scope, set, false); err != nil {
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
	for n := range b.StatementList {
		if err := processStatementListItem(&b.StatementList[n], scope, set); err != nil {
			return nil
		}
	}
	return nil
}

func processVariableStatement(v *javascript.VariableStatement, scope *Scope, set bool) error {
	for n := range v.VariableDeclarationList {
		if err := processVariableDeclaration(&v.VariableDeclarationList[n], scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processExpression(e *javascript.Expression, scope *Scope, set bool) error {
	for n := range e.Expressions {
		if err := processAssignmentExpression(&e.Expressions[n], scope, set); err != nil {
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
		for n := range f.InitVar {
			if err := processVariableDeclaration(&f.InitVar[n], scope, set); err != nil {
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
		if f.ForBindingPatternObject != nil {
			if err := processObjectBindingPattern(f.ForBindingPatternObject, scope, set, false, true); err != nil {
				return err
			}
		} else if f.ForBindingPatternArray != nil {
			if err := processArrayBindingPattern(f.ForBindingPatternArray, scope, set, false, true); err != nil {
				return err
			}
		} else if f.ForBindingIdentifier != nil && !set {
			scope.Parent.addBinding(f.ForBindingIdentifier)
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
	for n := range s.CaseClauses {
		c := &s.CaseClauses[n]
		if err := processExpression(&c.Expression, scope, set); err != nil {
			return err
		}
		for m := range c.StatementList {
			if err := processStatementListItem(&c.StatementList[m], scope, set); err != nil {
				return err
			}
		}
	}
	for n := range s.DefaultClause {
		if err := processStatementListItem(&s.DefaultClause[n], scope, set); err != nil {
			return err
		}
	}
	for n := range s.PostDefaultCaseClauses {
		c := &s.CaseClauses[n]
		if err := processExpression(&c.Expression, scope, set); err != nil {
			return err
		}
		for m := range c.StatementList {
			if err := processStatementListItem(&c.StatementList[m], scope, set); err != nil {
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

func processFunctionDeclaration(f *javascript.FunctionDeclaration, scope *Scope, set, expression bool) error {
	if f.BindingIdentifier != nil && set && !expression {
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

func processClassDeclaration(c *javascript.ClassDeclaration, scope *Scope, set, expression bool) error {
	if c.BindingIdentifier != nil && set && !expression {
		if err := scope.setBinding(c.BindingIdentifier, false); err != nil {
			return err
		}
	}
	if c.ClassHeritage != nil {
		if err := processLeftHandSideExpression(c.ClassHeritage, scope, set); err != nil {
			return err
		}
	}
	for n := range c.ClassBody {
		if err := processMethodDefinition(&c.ClassBody[n], scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processLexicalDeclaration(l *javascript.LexicalDeclaration, scope *Scope, set bool) error {
	for n := range l.BindingList {
		if err := processLexicalBinding(&l.BindingList[n], scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processVariableDeclaration(v *javascript.VariableDeclaration, scope *Scope, set bool) error {
	if v.ArrayBindingPattern != nil {
		if err := processArrayBindingPattern(v.ArrayBindingPattern, scope, set, true, false); err != nil {
			return err
		}
	} else if v.ObjectBindingPattern != nil {
		if err := processObjectBindingPattern(v.ObjectBindingPattern, scope, set, true, false); err != nil {
			return err
		}
	} else if v.BindingIdentifier != nil && set {
		if err := scope.setBinding(v.BindingIdentifier, true); err != nil {
			return err
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

func processObjectBindingPattern(o *javascript.ObjectBindingPattern, scope *Scope, set, hoist, bare bool) error {
	for n := range o.BindingPropertyList {
		if err := processBindingProperty(&o.BindingPropertyList[n], scope, set, hoist, bare); err != nil {
			return err
		}
	}
	if o.BindingRestProperty != nil {
		if bare {
			scope.addBinding(o.BindingRestProperty)
		} else if set {
			if err := scope.setBinding(o.BindingRestProperty, hoist); err != nil {
				return err
			}
		}
	}
	return nil
}

func processArrayBindingPattern(a *javascript.ArrayBindingPattern, scope *Scope, set, hoist, bare bool) error {
	for n := range a.BindingElementList {
		if err := processBindingElement(&a.BindingElementList[n], scope, set, hoist, bare); err != nil {
			return err
		}
	}
	if a.BindingRestElement != nil {
		if err := processBindingElement(a.BindingRestElement, scope, set, hoist, bare); err != nil {
			return err
		}
	}
	return nil
}

func processFormalParameters(f *javascript.FormalParameters, scope *Scope, set bool) error {
	for n := range f.FormalParameterList {
		if err := processBindingElement(&f.FormalParameterList[n], scope, set, false, false); err != nil {
			return err
		}
	}
	if f.FunctionRestParameter != nil {
		if err := processFunctionRestParameter(f.FunctionRestParameter, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processMethodDefinition(m *javascript.MethodDefinition, scope *Scope, set bool) error {
	if err := processPropertyName(&m.PropertyName, scope, set); err != nil {
		return err
	}
	scope = scope.Parent.newFunctionScope(m)
	if err := processFormalParameters(&m.Params, scope, set); err != nil {
		return err
	}
	if err := processBlockStatement(&m.FunctionBody, scope, set); err != nil {
		return err
	}
	return nil
}

func processLexicalBinding(l *javascript.LexicalBinding, scope *Scope, set bool) error {
	if l.ArrayBindingPattern != nil {
		if err := processArrayBindingPattern(l.ArrayBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if l.ObjectBindingPattern != nil {
		if err := processObjectBindingPattern(l.ObjectBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if l.BindingIdentifier != nil && set {
		if err := scope.setBinding(l.BindingIdentifier, false); err != nil {
			return err
		}
	}
	if l.Initializer != nil {
		if err := processAssignmentExpression(l.Initializer, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processConditionalExpression(c *javascript.ConditionalExpression, scope *Scope, set bool) error {
	if c.LogicalORExpression != nil {
		if err := processLogicalORExpression(c.LogicalORExpression, scope, set); err != nil {
			return err
		}
	} else if c.CoalesceExpression != nil {
		if err := processCoalesceExpression(c.CoalesceExpression, scope, set); err != nil {
			return err
		}
	}
	if c.True != nil {
		if err := processAssignmentExpression(c.True, scope, set); err != nil {
			return err
		}
	}
	if c.False != nil {
		if err := processAssignmentExpression(c.False, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processArrowFunction(a *javascript.ArrowFunction, scope *Scope, set bool) error {
	scope = scope.newArrowFunctionScope(a)
	if a.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := processCoverParenthesizedExpressionAndArrowParameterList(a.CoverParenthesizedExpressionAndArrowParameterList, scope, set); err != nil {
			return err
		}
	} else if a.FormalParameters != nil {
		if err := processFormalParameters(a.FormalParameters, scope, set); err != nil {
			return err
		}
	} else if a.BindingIdentifier != nil && set {
		if err := scope.setBinding(a.BindingIdentifier, false); err != nil {
			return err
		}
	}
	if a.AssignmentExpression != nil {
		if err := processAssignmentExpression(a.AssignmentExpression, scope, set); err != nil {
			return err
		}
	} else if a.FunctionBody != nil {
		if err := processBlockStatement(a.FunctionBody, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processNewExpression(n *javascript.NewExpression, scope *Scope, set bool) error {
	if err := processMemberExpression(&n.MemberExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processCallExpression(c *javascript.CallExpression, scope *Scope, set bool) error {
	if c.MemberExpression != nil {
		if err := processMemberExpression(c.MemberExpression, scope, set); err != nil {
			return err
		}
	} else if c.ImportCall != nil {
		if err := processAssignmentExpression(c.ImportCall, scope, set); err != nil {
			return err
		}
	} else if c.CallExpression != nil {
		if err := processCallExpression(c.CallExpression, scope, set); err != nil {
			return err
		}
	}
	if c.Arguments != nil {
		if err := processArguments(c.Arguments, scope, set); err != nil {
			return err
		}
	} else if c.Expression != nil {
		if err := processExpression(c.Expression, scope, set); err != nil {
			return err
		}
	} else if c.TemplateLiteral != nil {
		if err := processTemplateLiteral(c.TemplateLiteral, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processOptionalExpression(o *javascript.OptionalExpression, scope *Scope, set bool) error {
	if o.MemberExpression != nil {
		if err := processMemberExpression(o.MemberExpression, scope, set); err != nil {
			return err
		}
	} else if o.CallExpression != nil {
		if err := processCallExpression(o.CallExpression, scope, set); err != nil {
			return err
		}
	} else if o.OptionalExpression != nil {
		if err := processOptionalExpression(o.OptionalExpression, scope, set); err != nil {
			return err
		}
	}
	if err := processOptionalChain(&o.OptionalChain, scope, set); err != nil {
		return err
	}
	return nil
}

func processBindingProperty(b *javascript.BindingProperty, scope *Scope, set, hoist, bare bool) error {
	if err := processPropertyName(&b.PropertyName, scope, set); err != nil {
		return err
	}
	if err := processBindingElement(&b.BindingElement, scope, set, hoist, bare); err != nil {
		return err
	}
	return nil
}

func processBindingElement(b *javascript.BindingElement, scope *Scope, set, hoist, bare bool) error {
	if b.SingleNameBinding != nil {
		if bare {
			scope.addBinding(b.SingleNameBinding)
		} else if set {
			if err := scope.setBinding(b.SingleNameBinding, hoist); err != nil {
				return err
			}
		}
	} else if b.ArrayBindingPattern != nil {
		if err := processArrayBindingPattern(b.ArrayBindingPattern, scope, set, hoist, bare); err != nil {
			return err
		}
	} else if b.ObjectBindingPattern != nil {
		if err := processObjectBindingPattern(b.ObjectBindingPattern, scope, set, hoist, bare); err != nil {
			return err
		}
	}
	if b.Initializer != nil {
		if err := processAssignmentExpression(b.Initializer, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processFunctionRestParameter(f *javascript.FunctionRestParameter, scope *Scope, set bool) error {
	if f.ArrayBindingPattern != nil {
		if err := processArrayBindingPattern(f.ArrayBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if f.ObjectBindingPattern != nil {
		if err := processObjectBindingPattern(f.ObjectBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if f.BindingIdentifier != nil && set {
		if err := scope.setBinding(f.BindingIdentifier, false); err != nil {
			return err
		}
	}
	return nil
}

func processPropertyName(p *javascript.PropertyName, scope *Scope, set bool) error {
	if p.ComputedPropertyName != nil {
		if err := processAssignmentExpression(p.ComputedPropertyName, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processLogicalORExpression(l *javascript.LogicalORExpression, scope *Scope, set bool) error {
	if l.LogicalORExpression != nil {
		if err := processLogicalORExpression(l.LogicalORExpression, scope, set); err != nil {
			return err
		}
	}
	if err := processLogicalANDExpression(&l.LogicalANDExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processCoalesceExpression(c *javascript.CoalesceExpression, scope *Scope, set bool) error {
	if c.CoalesceExpressionHead != nil {
		if err := processCoalesceExpression(c.CoalesceExpressionHead, scope, set); err != nil {
			return err
		}
	}
	if err := processBitwiseORExpression(&c.BitwiseORExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processCoverParenthesizedExpressionAndArrowParameterList(c *javascript.CoverParenthesizedExpressionAndArrowParameterList, scope *Scope, set bool) error {
	for n := range c.Expressions {
		if err := processAssignmentExpression(&c.Expressions[n], scope, set); err != nil {
			return err
		}
	}
	if c.ArrayBindingPattern != nil {
		if err := processArrayBindingPattern(c.ArrayBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if c.ObjectBindingPattern != nil {
		if err := processObjectBindingPattern(c.ObjectBindingPattern, scope, set, false, false); err != nil {
			return err
		}
	} else if c.BindingIdentifier != nil && set {
		if err := scope.setBinding(c.BindingIdentifier, false); err != nil {
			return err
		}
	}
	return nil
}

func processMemberExpression(m *javascript.MemberExpression, scope *Scope, set bool) error {
	if m.PrimaryExpression != nil {
		if err := processPrimaryExpression(m.PrimaryExpression, scope, set); err != nil {
			return err
		}
	} else if m.MemberExpression != nil {
		if err := processMemberExpression(m.MemberExpression, scope, set); err != nil {
			return err
		}
		if m.Expression != nil {
			if err := processExpression(m.Expression, scope, set); err != nil {
				return err
			}
		} else if m.TemplateLiteral != nil {
			if err := processTemplateLiteral(m.TemplateLiteral, scope, set); err != nil {
				return err
			}
		} else if m.Arguments != nil {
			if err := processArguments(m.Arguments, scope, set); err != nil {
				return err
			}
		}
	}
	return nil
}

func processArguments(a *javascript.Arguments, scope *Scope, set bool) error {
	for n := range a.ArgumentList {
		if err := processAssignmentExpression(&a.ArgumentList[n], scope, set); err != nil {
			return err
		}
	}
	if a.SpreadArgument != nil {
		if err := processAssignmentExpression(a.SpreadArgument, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processTemplateLiteral(t *javascript.TemplateLiteral, scope *Scope, set bool) error {
	for n := range t.Expressions {
		if err := processExpression(&t.Expressions[n], scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processOptionalChain(o *javascript.OptionalChain, scope *Scope, set bool) error {
	if o.OptionalChain != nil {
		if err := processOptionalChain(o.OptionalChain, scope, set); err != nil {
			return err
		}
	}
	if o.Arguments != nil {
		if err := processArguments(o.Arguments, scope, set); err != nil {
			return err
		}
	} else if o.Expression != nil {
		if err := processExpression(o.Expression, scope, set); err != nil {
			return err
		}
	} else if o.TemplateLiteral != nil {
		if err := processTemplateLiteral(o.TemplateLiteral, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processLogicalANDExpression(l *javascript.LogicalANDExpression, scope *Scope, set bool) error {
	if l.LogicalANDExpression != nil {
		if err := processLogicalANDExpression(l.LogicalANDExpression, scope, set); err != nil {
			return err
		}
	}
	if err := processBitwiseORExpression(&l.BitwiseORExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processBitwiseORExpression(b *javascript.BitwiseORExpression, scope *Scope, set bool) error {
	if b.BitwiseORExpression != nil {
		if err := processBitwiseORExpression(b.BitwiseORExpression, scope, set); err != nil {
			return err
		}
	}
	if err := processBitwiseXORExpression(&b.BitwiseXORExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processPrimaryExpression(p *javascript.PrimaryExpression, scope *Scope, set bool) error {
	if p.ArrayLiteral != nil {
		if err := processArrayLiteral(p.ArrayLiteral, scope, set); err != nil {
			return err
		}
	} else if p.ObjectLiteral != nil {
		if err := processObjectLiteral(p.ObjectLiteral, scope, set); err != nil {
			return err
		}
	} else if p.FunctionExpression != nil {
		if err := processFunctionDeclaration(p.FunctionExpression, scope, set, true); err != nil {
			return err
		}
	} else if p.ClassExpression != nil {
		if err := processClassDeclaration(p.ClassExpression, scope, set, true); err != nil {
			return err
		}
	} else if p.TemplateLiteral != nil {
		if err := processTemplateLiteral(p.TemplateLiteral, scope, set); err != nil {
			return err
		}
	} else if p.CoverParenthesizedExpressionAndArrowParameterList != nil {
		if err := processCoverParenthesizedExpressionAndArrowParameterList(p.CoverParenthesizedExpressionAndArrowParameterList, scope, set); err != nil {
			return err
		}
	} else if p.IdentifierReference != nil && !set {
		scope.addBinding(p.IdentifierReference)
	}
	return nil
}

func processBitwiseXORExpression(b *javascript.BitwiseXORExpression, scope *Scope, set bool) error {
	if b.BitwiseXORExpression != nil {
		if err := processBitwiseXORExpression(b.BitwiseXORExpression, scope, set); err != nil {
			return err
		}
	}
	if err := processBitwiseANDExpression(&b.BitwiseANDExpression, scope, set); err != nil {
		return err
	}
	return nil
}

func processArrayLiteral(a *javascript.ArrayLiteral, scope *Scope, set bool) error {
	for n := range a.ElementList {
		if err := processAssignmentExpression(&a.ElementList[n], scope, set); err != nil {
			return err
		}
	}
	if a.SpreadElement != nil {
		if err := processAssignmentExpression(a.SpreadElement, scope, set); err != nil {
			return err
		}
	}
	return nil
}

func processObjectLiteral(o *javascript.ObjectLiteral, scope *Scope, set bool) error {
	return nil
}

func processBitwiseANDExpression(b *javascript.BitwiseANDExpression, scope *Scope, set bool) error {
	return nil
}

// Errors
var (
	ErrDuplicateBinding = errors.New("duplicate binding")
	ErrInvalidImport    = errors.New("invalid import")
)
