package minify

import (
	"vimagination.zapto.org/javascript"
)

func blockAsModule(b *javascript.Block, fn func(*javascript.Module)) {
	if len(b.StatementList) == 0 {
		return
	}
	m := javascript.ScriptToModule(&javascript.Script{
		StatementList: b.StatementList,
	})
	fn(m)
	b.StatementList = make([]javascript.StatementListItem, len(m.ModuleListItems))
	for n, mi := range m.ModuleListItems {
		b.StatementList[n] = *mi.StatementListItem
	}
}

func expressionsAsModule(e *javascript.Expression, fn func(*javascript.Module)) {
	if len(e.Expressions) == 0 {
		return
	}
	m := &javascript.Module{
		ModuleListItems: make([]javascript.ModuleItem, len(e.Expressions)),
	}
	for n := range e.Expressions {
		m.ModuleListItems[n] = javascript.ModuleItem{
			StatementListItem: &javascript.StatementListItem{
				Statement: &javascript.Statement{
					ExpressionStatement: &javascript.Expression{
						Expressions: e.Expressions[n : n+1],
					},
				},
			},
		}
	}
	fn(m)
	e.Expressions = make([]javascript.AssignmentExpression, len(m.ModuleListItems))
	for n := range m.ModuleListItems {
		e.Expressions[n] = m.ModuleListItems[n].StatementListItem.Statement.ExpressionStatement.Expressions[0]
	}
}

func isReturnStatement(s *javascript.Statement) bool {
	return s != nil && s.Type == javascript.StatementReturn
}

func isNonEmptyReturnStatement(s *javascript.Statement) bool {
	return isReturnStatement(s) && s.ExpressionStatement != nil
}

func statementsListItemsAsExpressionsAndReturn(sli []javascript.StatementListItem) ([]javascript.AssignmentExpression, bool) {
	var (
		expressions []javascript.AssignmentExpression
		hasReturn   bool
	)
	for n := range sli {
		s := &sli[n]
		if hasReturn {
			if isHoistable(s) {
				return nil, true
			}
		} else if isNonEmptyReturnStatement(s.Statement) {
			expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
			hasReturn = true
		} else if !isSLIExpression(s) {
			if isEmptyStatement(s.Statement) {
				continue
			}
			return nil, false
		} else {
			expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
		}
	}
	return expressions, hasReturn
}

func isSLIExpression(s *javascript.StatementListItem) bool {
	return s != nil && s.Declaration == nil && isStatementExpression(s.Statement)
}

func isStatementExpression(s *javascript.Statement) bool {
	return s != nil && s.Type == javascript.StatementNormal && s.ExpressionStatement != nil
}

func isEmptyStatement(s *javascript.Statement) bool {
	return s != nil && s.Type == javascript.StatementNormal && s.BlockStatement == nil && s.VariableStatement == nil && s.ExpressionStatement == nil && s.IfStatement == nil && s.IterationStatementDo == nil && s.IterationStatementFor == nil && s.IterationStatementWhile == nil && s.SwitchStatement == nil && s.WithStatement == nil && s.LabelledItemFunction == nil && s.LabelledItemStatement == nil && s.TryStatement == nil
}

func isHoistable(s *javascript.StatementListItem) bool {
	return s != nil && ((s.Statement != nil && (s.Statement.VariableStatement != nil || s.Statement.LabelledItemFunction != nil)) || (s.Declaration != nil && (s.Declaration.FunctionDeclaration != nil || s.Declaration.ClassDeclaration != nil)))
}

func aeIsCE(ae *javascript.AssignmentExpression) bool {
	return ae != nil && ae.ConditionalExpression != nil && ae.AssignmentOperator == javascript.AssignmentNone && !ae.Yield
}

func aeAsParen(ae *javascript.AssignmentExpression) *javascript.ParenthesizedExpression {
	if aeIsCE(ae) {
		pe, ok := javascript.UnwrapConditional(ae.ConditionalExpression).(*javascript.ParenthesizedExpression)
		if ok {
			return pe
		}
	}
	return nil
}

func meIsSinglePe(me *javascript.MemberExpression) bool {
	return me != nil && me.PrimaryExpression != nil && me.PrimaryExpression.ParenthesizedExpression != nil && len(me.PrimaryExpression.ParenthesizedExpression.Expressions) == 1 && aeIsCE(&me.PrimaryExpression.ParenthesizedExpression.Expressions[0])
}

func meAsCE(me *javascript.MemberExpression) *javascript.CallExpression {
	var ce *javascript.CallExpression
	if meIsSinglePe(me) {
		ce, _ = javascript.UnwrapConditional(me.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression).(*javascript.CallExpression)
	} else if me.MemberExpression != nil && me.Arguments == nil {
		ce = meAsCE(me.MemberExpression)
		if ce != nil {
			ce = &javascript.CallExpression{
				CallExpression:    ce,
				IdentifierName:    me.IdentifierName,
				Expression:        me.Expression,
				TemplateLiteral:   me.TemplateLiteral,
				PrivateIdentifier: me.PrivateIdentifier,
				Tokens:            me.Tokens,
			}
		}
	}
	return ce
}

func isStatementListItemExpression(s *javascript.StatementListItem) bool {
	return s != nil && isStatementExpression(s.Statement)
}

func leftMostLHS(c javascript.ConditionalWrappable) *javascript.LeftHandSideExpression {
	for {
		switch t := c.(type) {
		case *javascript.ConditionalExpression:
			if t.CoalesceExpression != nil {
				ce := t.CoalesceExpression
				for ce.CoalesceExpressionHead != nil {
					ce = ce.CoalesceExpressionHead
				}
				c = &ce.BitwiseORExpression
			} else if t.LogicalORExpression != nil {
				c = t.LogicalORExpression
			} else {
				return nil
			}
		case *javascript.LogicalORExpression:
			if t.LogicalORExpression != nil {
				c = t.LogicalORExpression
			} else {
				c = &t.LogicalANDExpression
			}
		case *javascript.LogicalANDExpression:
			if t.LogicalANDExpression != nil {
				c = t.LogicalANDExpression
			} else {
				c = &t.BitwiseORExpression
			}
		case *javascript.BitwiseORExpression:
			if t.BitwiseORExpression != nil {
				c = t.BitwiseORExpression
			} else {
				c = &t.BitwiseXORExpression
			}
		case *javascript.BitwiseXORExpression:
			if t.BitwiseXORExpression != nil {
				c = t.BitwiseXORExpression
			} else {
				c = &t.BitwiseANDExpression
			}
		case *javascript.BitwiseANDExpression:
			if t.BitwiseANDExpression != nil {
				c = t.BitwiseANDExpression
			} else {
				c = &t.EqualityExpression
			}
		case *javascript.EqualityExpression:
			if t.EqualityExpression != nil {
				c = t.EqualityExpression
			} else {
				c = &t.RelationalExpression
			}
		case *javascript.RelationalExpression:
			if t.RelationalExpression != nil {
				c = t.RelationalExpression
			} else {
				c = &t.ShiftExpression
			}
		case *javascript.ShiftExpression:
			if t.ShiftExpression != nil {
				c = t.ShiftExpression
			} else {
				c = &t.AdditiveExpression
			}
		case *javascript.AdditiveExpression:
			if t.AdditiveExpression != nil {
				c = t.AdditiveExpression
			} else {
				c = &t.MultiplicativeExpression
			}
		case *javascript.MultiplicativeExpression:
			if t.MultiplicativeExpression != nil {
				c = t.MultiplicativeExpression
			} else {
				c = &t.ExponentiationExpression
			}
		case *javascript.ExponentiationExpression:
			if t.ExponentiationExpression != nil {
				c = t.ExponentiationExpression
			} else {
				c = &t.UnaryExpression.UpdateExpression
			}
		case *javascript.UpdateExpression:
			if t.LeftHandSideExpression != nil {
				return t.LeftHandSideExpression
			}
			c = &t.UnaryExpression.UpdateExpression
		}
	}
}

func fixWrapping(s *javascript.Statement) {
	ae := &s.ExpressionStatement.Expressions[0]
	if aeIsCE(ae) {
		if lhs := leftMostLHS(ae.ConditionalExpression); lhs != nil && lhs.NewExpression != nil && lhs.NewExpression.News == 0 {
			me := &lhs.NewExpression.MemberExpression
			for me.MemberExpression != nil {
				me = me.MemberExpression
			}
			if me.PrimaryExpression.ObjectLiteral != nil || me.PrimaryExpression.FunctionExpression != nil || me.PrimaryExpression.ClassExpression != nil {
				me.PrimaryExpression = &javascript.PrimaryExpression{
					ParenthesizedExpression: &javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(me.PrimaryExpression),
								Tokens:                me.PrimaryExpression.Tokens,
							},
						},
						Tokens: me.PrimaryExpression.Tokens,
					},
					Tokens: me.PrimaryExpression.Tokens,
				}
			}
		}
		switch javascript.UnwrapConditional(ae.ConditionalExpression).(type) {
		case *javascript.ObjectLiteral, *javascript.FunctionDeclaration, *javascript.ClassDeclaration:
			ae.ConditionalExpression = javascript.WrapConditional(&javascript.ParenthesizedExpression{
				Expressions: []javascript.AssignmentExpression{*ae},
				Tokens:      ae.Tokens,
			})
		}
	}
}

func scoreCE(ce javascript.ConditionalWrappable) int {
	switch ce.(type) {
	case *javascript.LogicalORExpression:
		return 1
	case *javascript.LogicalANDExpression:
		return 2
	case *javascript.BitwiseORExpression:
		return 3
	case *javascript.BitwiseXORExpression:
		return 4
	case *javascript.BitwiseANDExpression:
		return 5
	case *javascript.EqualityExpression:
		return 6
	case *javascript.RelationalExpression:
		return 7
	case *javascript.ShiftExpression:
		return 8
	case *javascript.AdditiveExpression:
		return 9
	case *javascript.MultiplicativeExpression:
		return 10
	case *javascript.ExponentiationExpression:
		return 11
	case *javascript.UnaryExpression:
		return 12
	case *javascript.UpdateExpression:
		return 13
	}
	return -1
}

func isConditionalWrappingAConditional(w javascript.ConditionalWrappable, below javascript.ConditionalWrappable) *javascript.ConditionalExpression {
	pe, ok := javascript.UnwrapConditional(javascript.WrapConditional(w)).(*javascript.ParenthesizedExpression)
	if !ok || len(pe.Expressions) != 1 || !aeIsCE(&pe.Expressions[0]) {
		return nil
	}
	uw := javascript.UnwrapConditional(pe.Expressions[0].ConditionalExpression)
	if scoreCE(uw) < scoreCE(below) {
		return nil
	}
	return pe.Expressions[0].ConditionalExpression
}

func removeLastReturnStatement(b *javascript.Block) {
	if len(b.StatementList) > 0 {
		s := b.StatementList[len(b.StatementList)-1].Statement
		if isReturnStatement(s) && s.ExpressionStatement == nil {
			b.StatementList = b.StatementList[:len(b.StatementList)-1]
		}
	}
}
