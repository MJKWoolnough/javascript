package minify

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

func removeDeadCode(m *javascript.Module) {
	for {
		s, err := scope.ModuleScope(m, nil)
		if err != nil {
			return
		}

		clearSinglesFromScope(s)

		walk.Walk(m, walk.HandlerFunc(deadWalker))
	}
}

func clearSinglesFromScope(s *scope.Scope) bool {
	changed := false
	for name, bindings := range s.Bindings {
		if name == "this" || name == "arguments" || len(bindings) != 1 {
			continue
		}
		bindings[0].Token.Data = ""
		changed = true
	}

	for _, cs := range s.Scopes {
		if clearSinglesFromScope(cs) {
			changed = true
		}
	}

	return changed
}

func deadWalker(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.Module:
		removeDeadCodeFromModule(t)
	case *javascript.Block:
		blockAsModule(t, removeDeadCodeFromModule)
	case *javascript.Expression:
		expressionsAsModule(t, removeDeadCodeFromModule)
	default:
		deadWalker(t)
	}
	return nil
}

func removeDeadCodeFromModule(m *javascript.Module) {
	for i := 0; i < len(m.ModuleListItems); i++ {
		switch sliBindable(m.ModuleListItems[i].StatementListItem) {
		case bindableConst, bindableLet:
			ld := m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration
			if removeDeadLexicalBindings(&m.ModuleListItems, i, &ld.BindingList, lexicalMaker(ld.LetOrConst)) {
				i--
			}
		case bindableVar:
			if removeDeadLexicalBindings(&m.ModuleListItems, i, &m.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList, variableMaker) {
				i--
			}
		case bindableBare:
			expr := m.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement
			if pe, ok := javascript.UnwrapConditional(expr.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression); ok && pe.IdentifierReference != nil && pe.IdentifierReference.Data == "" {
				expr.Expressions[0] = *expr.Expressions[0].AssignmentExpression
			}
		case bindableFunction:
			m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
			i--
		case bindableClass:
		}
	}
}

func lexicalMaker(LetOrConst javascript.LetOrConst) func([]javascript.LexicalBinding) javascript.ModuleItem {
	return func(lbs []javascript.LexicalBinding) javascript.ModuleItem {
		return javascript.ModuleItem{
			StatementListItem: &javascript.StatementListItem{
				Declaration: &javascript.Declaration{
					LexicalDeclaration: &javascript.LexicalDeclaration{
						LetOrConst:  LetOrConst,
						BindingList: lbs,
					},
				},
			},
		}
	}
}

func variableMaker(vds []javascript.VariableDeclaration) javascript.ModuleItem {
	return javascript.ModuleItem{
		StatementListItem: &javascript.StatementListItem{
			Statement: &javascript.Statement{
				VariableStatement: &javascript.VariableStatement{
					VariableDeclarationList: vds,
				},
			},
		},
	}
}

func removeDeadLexicalBindings(mlis *[]javascript.ModuleItem, pos int, lds *[]javascript.LexicalBinding, sm func([]javascript.LexicalBinding) javascript.ModuleItem) bool {
	for n, ld := range *lds {
		if ld.BindingIdentifier != nil && ld.BindingIdentifier.Data == "" {
			toAdd := make([]javascript.ModuleItem, 0, 3+len(*mlis)-pos)
			rest := (*lds)[n+1:]
			if n > 0 {
				toAdd = append(toAdd, (*mlis)[pos])
				*lds = (*lds)[:n]
			}
			if ld.Initializer != nil {
				toAdd = append(toAdd, javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{*ld.Initializer},
							},
						},
					},
				})
			}
			if len(rest) > 0 {
				toAdd = append(toAdd, sm(rest))
			}
			*mlis = append((*mlis)[:pos], append(toAdd, (*mlis)[pos+1:]...)...)
			return true
		}
	}
	return false
}

type bindable byte

const (
	bindableNone bindable = iota
	bindableConst
	bindableLet
	bindableVar
	bindableClass
	bindableFunction
	bindableBare
)

func sliBindable(sli *javascript.StatementListItem) bindable {
	if sli != nil {
		if sli.Declaration != nil {
			if sli.Declaration.LexicalDeclaration != nil {
				if sli.Declaration.LexicalDeclaration.LetOrConst == javascript.Const {
					return bindableConst
				}
				return bindableLet
			} else if sli.Declaration.ClassDeclaration != nil {
				return bindableClass
			} else if sli.Declaration.FunctionDeclaration != nil {
				return bindableFunction
			}
		}
		if sli.Statement != nil && sli.Statement.VariableStatement != nil {
			return bindableVar
		}
		if isStatementExpression(sli.Statement) {
			if sli.Statement.ExpressionStatement.Expressions[0].AssignmentOperator == javascript.AssignmentAssign {
				return bindableBare
			}
		}
	}
	return bindableNone
}

func removeDeadSLI(sli *javascript.StatementListItem) bool {
	switch sliBindable(sli) {
	case bindableConst, bindableLet:
		lb := sli.Declaration.LexicalDeclaration.BindingList[0]
		if lb.BindingIdentifier != nil {
			return lb.BindingIdentifier.Data == ""
		}
	case bindableVar:
		vd := sli.Statement.VariableStatement.VariableDeclarationList[0]
		if vd.BindingIdentifier != nil {
			return vd.BindingIdentifier.Data == ""
		}
	case bindableClass:
		return sli.Declaration.ClassDeclaration.BindingIdentifier == nil || sli.Declaration.ClassDeclaration.BindingIdentifier.Data == ""
	case bindableFunction:
		return sli.Declaration.FunctionDeclaration.BindingIdentifier == nil || sli.Declaration.FunctionDeclaration.BindingIdentifier.Data == ""
	}

	return false
}
