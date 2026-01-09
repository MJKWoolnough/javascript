package minify

import (
	"slices"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

func (p *processor) removeDeadCode(m *javascript.Module) {
	s, err := scope.ModuleScope(m, nil)
	if err != nil {
		return
	}

	p.clearSinglesFromScope(s)
	p.deadWalker(m)
}

func (p *processor) clearSinglesFromScope(s *scope.Scope) {
	for name, bindings := range s.Bindings {
		if name == "this" || name == "arguments" || len(bindings) != 1 || bindings[0].BindingType == scope.BindingRef {
			continue
		}

		bindings[0].Token.Data = ""
		p.changed = true
	}

	for _, cs := range s.Scopes {
		p.clearSinglesFromScope(cs)
	}
}

func (p *processor) deadWalker(t javascript.Type) error {
	walk.Walk(t, walk.HandlerFunc(p.deadWalker))

	switch t := t.(type) {
	case *javascript.Module:
		if removeDeadCodeFromModule(t) {
			p.changed = true
		}
	case *javascript.Block:
		if blockAsModule(t, removeDeadCodeFromModule) {
			p.changed = true
		}
	case *javascript.Expression:
		if expressionsAsModule(t, removeDeadCodeFromModule) {
			p.changed = true
		}
	case *javascript.Statement:
		if t.ExpressionStatement != nil {
			if newExpressions := removeDeadExpressions(t.ExpressionStatement.Expressions); len(newExpressions) != len(t.ExpressionStatement.Expressions) {
				t.ExpressionStatement.Expressions = newExpressions
				p.changed = true
			}
		} else if t.IfStatement != nil && isEmptyStatement(&t.IfStatement.Statement) {
			t.ExpressionStatement = &t.IfStatement.Expression
			t.IfStatement = nil
			p.changed = true
		}
	case *javascript.ParenthesizedExpression:
		if newExpressions := removeDeadExpressions(t.Expressions[:len(t.Expressions)-1]); len(newExpressions) != len(t.Expressions)-1 {
			t.Expressions = append(newExpressions, t.Expressions[len(t.Expressions)-1])
			p.changed = true
		}
	case *javascript.IfStatement:
		if isEmptyStatement(&t.Statement) {
			if t.ElseStatement != nil {
				t.Expression.Expressions = []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.UnaryExpression{
							UnaryOperators: []javascript.UnaryOperatorComments{
								{
									UnaryOperator: javascript.UnaryLogicalNot,
								},
							},
							UpdateExpression: javascript.UpdateExpression{
								LeftHandSideExpression: &javascript.LeftHandSideExpression{
									NewExpression: &javascript.NewExpression{
										MemberExpression: javascript.MemberExpression{
											PrimaryExpression: &javascript.PrimaryExpression{
												ParenthesizedExpression: &javascript.ParenthesizedExpression{
													Expressions: t.Expression.Expressions,
												},
											},
										},
									},
								},
							},
						}),
					},
				}
				t.Statement = *t.ElseStatement
				t.ElseStatement = nil
				p.changed = true
			}
		}
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
				changed = true
			}
		case bindableVar:
			if removeDeadLexicalBindings(&m.ModuleListItems, i, &m.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList, variableMaker) {
				i--
				changed = true
			}
		case bindableBare:
			expr := m.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement

			if pe, ok := javascript.UnwrapConditional(expr.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression); ok && pe.IdentifierReference != nil && pe.IdentifierReference.Data == "" {
				expr.Expressions[0] = *expr.Expressions[0].AssignmentExpression
				i--
				changed = true
			}
		case bindableFunction:
			if m.ModuleListItems[i].StatementListItem.Declaration.FunctionDeclaration.BindingIdentifier.Data == "" {
				m.ModuleListItems = slices.Delete(m.ModuleListItems, i, i+1)
				i--
				changed = true
			}
		case bindableClass:
			if cd := m.ModuleListItems[i].StatementListItem.Declaration.ClassDeclaration; cd.BindingIdentifier.Data == "" {
				mis := extractStatementsFromClass(cd)
				m.ModuleListItems = append(m.ModuleListItems[:i], append(mis, m.ModuleListItems[i+1:]...)...)
				i--
				changed = true
			}
		default:
			if isSLIExpression(m.ModuleListItems[i].StatementListItem) && len(m.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement.Expressions) == 0 {
				m.ModuleListItems = slices.Delete(m.ModuleListItems, i, i+1)
				i--
				changed = true
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

		if isNonEmptyStatementExpression(sli.Statement) {
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

func removeDeadExpressions(expressions []javascript.AssignmentExpression) []javascript.AssignmentExpression {
	for n := 0; n < len(expressions); n++ {
		if expressions[n].ConditionalExpression == nil {
			continue
		}

		switch javascript.UnwrapConditional(expressions[n].ConditionalExpression).(type) {
		case *javascript.PrimaryExpression:
			expressions = slices.Delete(expressions, n, n+1)
			n--
		}
	}

	return expressions
}
