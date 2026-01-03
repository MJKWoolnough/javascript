package minify

import (
	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

func removeDeadCode(m *javascript.Module) {
	c := changeTracker(true)

	for c {
		s, err := scope.ModuleScope(m, nil)
		if err != nil {
			return
		}

		c = false

		clearSinglesFromScope(s)

		walk.Walk(m, walk.HandlerFunc(c.deadWalker))
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

type changeTracker bool

func (c *changeTracker) deadWalker(t javascript.Type) error {
	var changed bool

	switch t := t.(type) {
	case *javascript.Module:
		changed = removeDeadCodeFromModule(t)
	case *javascript.Block:
		changed = blockAsModule(t, removeDeadCodeFromModule)
	case *javascript.Expression:
		changed = expressionsAsModule(t, removeDeadCodeFromModule)
	default:
		walk.Walk(t, walk.HandlerFunc(c.deadWalker))
	}

	if changed {
		*c = true
	}

	return nil
}

func removeDeadCodeFromModule(m *javascript.Module) bool {
	var changed bool

	for i := 0; i < len(m.ModuleListItems); i++ {
		switch sliBindable(m.ModuleListItems[i].StatementListItem) {
		case bindableConst, bindableLet:
			if ld := m.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration; removeDeadLexicalBindings(&m.ModuleListItems, i, &ld.BindingList, lexicalMaker(ld.LetOrConst)) {
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
				i--
			}
		case bindableFunction:
			if m.ModuleListItems[i].StatementListItem.Declaration.FunctionDeclaration.BindingIdentifier.Data == "" {
				m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
				i--
			}
		case bindableClass:
			if cd := m.ModuleListItems[i].StatementListItem.Declaration.ClassDeclaration; cd.BindingIdentifier.Data == "" {
				mis := extractStatementsFromClass(cd)
				m.ModuleListItems = append(m.ModuleListItems[:i], append(mis, m.ModuleListItems[i+1:]...)...)
				i--
			}
		}
	}

	return changed
}

func extractStatementsFromClass(cd *javascript.ClassDeclaration) []javascript.ModuleItem {
	var mis []javascript.ModuleItem

	if cd.ClassHeritage != nil {
		mis = append(mis, javascript.ModuleItem{
			StatementListItem: &javascript.StatementListItem{
				Statement: &javascript.Statement{
					ExpressionStatement: &javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(cd.ClassHeritage),
							},
						},
					},
				},
			},
		})
	}

	for _, ce := range cd.ClassBody {
		if ce.FieldDefinition != nil {
			fd := ce.FieldDefinition

			if fd.ClassElementName.PropertyName != nil && fd.ClassElementName.PropertyName.ComputedPropertyName != nil {
				mis = append(mis, javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{*fd.ClassElementName.PropertyName.ComputedPropertyName},
							},
						},
					},
				})
			}

			if fd.Initializer != nil {
				mis = append(mis, javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{*fd.Initializer},
							},
						},
					},
				})
			}
		} else if ce.MethodDefinition != nil {
			if md := ce.MethodDefinition; md.ClassElementName.PropertyName != nil && md.ClassElementName.PropertyName.ComputedPropertyName != nil {
				mis = append(mis, javascript.ModuleItem{
					StatementListItem: &javascript.StatementListItem{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{*md.ClassElementName.PropertyName.ComputedPropertyName},
							},
						},
					},
				})
			}
		} else if ce.ClassStaticBlock != nil {
		}
	}

	return nil
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
		if lb := sli.Declaration.LexicalDeclaration.BindingList[0]; lb.BindingIdentifier != nil {
			return lb.BindingIdentifier.Data == ""
		}
	case bindableVar:
		if vd := sli.Statement.VariableStatement.VariableDeclarationList[0]; vd.BindingIdentifier != nil {
			return vd.BindingIdentifier.Data == ""
		}
	case bindableClass:
		return sli.Declaration.ClassDeclaration.BindingIdentifier == nil || sli.Declaration.ClassDeclaration.BindingIdentifier.Data == ""
	case bindableFunction:
		return sli.Declaration.FunctionDeclaration.BindingIdentifier == nil || sli.Declaration.FunctionDeclaration.BindingIdentifier.Data == ""
	}

	return false
}
