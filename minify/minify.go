package minify

import (
	"strconv"
	"strings"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/scope"
	"vimagination.zapto.org/javascript/walk"
)

type Minifier struct {
	Option
}

func New(opts ...Option) *Minifier {
	m := new(Minifier)

	if len(opts) == 0 {
		m.Option = Safe
	}

	for _, opt := range opts {
		m.Option |= opt
	}

	return m
}

type processor struct {
	*Minifier
	changed bool
}

func (p *processor) Handle(t javascript.Type) error {
	if err := walk.Walk(t, p); err != nil {
		return err
	}

	switch t := t.(type) {
	case *javascript.TemplateLiteral:
		p.minifyTemplates(t)
	case *javascript.PrimaryExpression:
		p.minifyLiterals(t)
		p.minifyNonHoistableNames(t)
	case *javascript.ArrowFunction:
		p.minifyArrowFunc(t)
		p.minifyLastReturnStatementInArrowFn(t)
		p.fixFirstArrowFuncExpression(t)
	case *javascript.Statement:
		p.minifyBlockToStatement(t)
		p.minifyIfToConditional(t)
		p.removeDebugger(t)
	case *javascript.PropertyName:
		p.minifyObjectKeys(t)
	case *javascript.AssignmentExpression:
		p.minifyFunctionExpressionAsArrowFunc(t)
		p.minifyAEParens(t)
	case *javascript.ParenthesizedExpression:
		p.minifyParenthsizedExpressionParens(t)
	case *javascript.Expression:
		p.minifyExpressionParens(t)
	case *javascript.Argument:
		p.minifyArgumentParens(t)
	case *javascript.MemberExpression:
		p.minifyMemberExpressionParens(t)
	case *javascript.CallExpression:
		p.minifyCallExpressionParens(t)
	case *javascript.LeftHandSideExpression:
		p.minifyLHSExpressionParens(t)
	case *javascript.Block:
		p.minifyRemoveDeadCode(t)
		blockAsModule(t, p.minifyEmptyStatement)
		blockAsModule(t, p.minifyExpressionRun)
		blockAsModule(t, p.fixFirstExpression)
		blockAsModule(t, p.minifyLexical)
		blockAsModule(t, p.minifyExpressionsBetweenLexicals)
	case *javascript.Module:
		p.minifyEmptyStatement(t)
		p.minifyExpressionRun(t)
		p.fixFirstExpression(t)
		p.minifyLexical(t)
		p.minifyExpressionsBetweenLexicals(t)
	case *javascript.FunctionDeclaration:
		p.minifyLastReturnStatement(t)
	case *javascript.ConditionalExpression:
		p.minifyConditionExpressionParens(t)
	}

	return nil
}

func (m *Minifier) Process(jm *javascript.Module) {
	p := &processor{changed: true, Minifier: m}

	for p.changed {
		p.changed = false

		if p.Has(RemoveDeadCode) {
			p.removeDeadCode(jm)
		}

		walk.Walk(jm, p)

		if p.Has(RenameIdentifiers) {
			renameIdentifiers(jm)
		}
	}
}

func (p *processor) minifyTemplate(t *javascript.Token) {
	if t != nil {
		str, err := javascript.UnquoteTemplate(t.Data)
		if err != nil {
			return
		}

		res := javascript.QuoteTemplate(str, javascript.TokenTypeToTemplateType(t.Type))
		if len(res) < len(t.Data) {
			t.Data = res
			p.changed = true
		}
	}
}

func (p *processor) minifyTemplates(t *javascript.TemplateLiteral) {
	if p.Has(Literals) {
		p.minifyTemplate(t.NoSubstitutionTemplate)
		p.minifyTemplate(t.TemplateHead)

		for _, m := range t.TemplateMiddleList {
			p.minifyTemplate(m)
		}

		p.minifyTemplate(t.TemplateTail)
	}
}

func (p *processor) minifyLiterals(pe *javascript.PrimaryExpression) {
	if p.Has(Literals) {
		if pe.Literal != nil {
			switch pe.Literal.Type {
			case javascript.TokenBooleanLiteral:
				switch pe.Literal.Data {
				case "true":
					pe.Literal.Data = "!0"
					p.changed = true
				case "false":
					pe.Literal.Data = "!1"
					p.changed = true
				}
			case javascript.TokenStringLiteral:
				str, err := javascript.Unquote(pe.Literal.Data)
				if err != nil {
					return
				}

				tmpl := javascript.QuoteTemplate(str, javascript.TemplateNoSubstitution)
				if len(tmpl) < len(pe.Literal.Data) {
					pe.TemplateLiteral = &javascript.TemplateLiteral{
						NoSubstitutionTemplate: pe.Literal,
						Tokens:                 pe.Tokens,
					}
					pe.Literal = nil
					pe.TemplateLiteral.NoSubstitutionTemplate.Type = javascript.TokenNoSubstitutionTemplate
					pe.TemplateLiteral.NoSubstitutionTemplate.Data = tmpl
					p.changed = true
				}
			case javascript.TokenNumericLiteral:
				p.minifyNumbers(pe.Literal)
			}
		} else if pe.IdentifierReference != nil && pe.IdentifierReference.Data == "undefined" {
			pe.IdentifierReference.Data = "void 0"
			p.changed = true
		} else if pe.TemplateLiteral != nil && pe.TemplateLiteral.NoSubstitutionTemplate != nil {
			str, err := javascript.UnquoteTemplate(pe.TemplateLiteral.NoSubstitutionTemplate.Data)
			if err != nil {
				return
			}

			asStrLit := strconv.Quote(str)
			if len(asStrLit) < len(pe.TemplateLiteral.NoSubstitutionTemplate.Data) {
				tk := pe.TemplateLiteral.NoSubstitutionTemplate
				tk.Data = asStrLit
				tk.Type = javascript.TokenStringLiteral
				pe.TemplateLiteral = nil
				pe.Literal = tk
				p.changed = true
			}
		}
	}
}

func (p *processor) minifyNumbers(nt *javascript.Token) {
	d := strings.ReplaceAll(nt.Data, "_", "")

	if len(d) < len(nt.Data) {
		nt.Data = d
		p.changed = true
	}

	if !strings.HasSuffix(d, "n") {
		var n float64
		if strings.HasPrefix(d, "0o") || strings.HasPrefix(d, "0O") {
			h, err := strconv.ParseUint(d[2:], 8, 64)
			if err != nil {
				return
			}

			n = float64(h)
		} else if strings.HasPrefix(d, "0b") || strings.HasPrefix(d, "0B") {
			h, err := strconv.ParseUint(d[2:], 2, 64)
			if err != nil {
				return
			}

			n = float64(h)
		} else if strings.HasPrefix(d, "0x") || strings.HasPrefix(d, "0X") {
			h, err := strconv.ParseUint(d[2:], 16, 64)
			if err != nil {
				return
			}

			n = float64(h)
		} else {
			f, err := strconv.ParseFloat(nt.Data, 64)
			if err != nil {
				return
			}

			n = f
		}

		d = strconv.FormatFloat(n, 'f', -1, 64)

		if strings.HasSuffix(d, "000") {
			var e uint64

			for strings.HasSuffix(d, "0") {
				d = d[:len(d)-1]
				e++
			}

			d += "e" + strconv.FormatUint(e, 10)
		} else if strings.HasPrefix(d, "0.00") {
			for strings.HasSuffix(d, "0") {
				d = d[:len(d)-1]
			}

			d = d[2:]
			e := uint64(len(d))

			for strings.HasPrefix(d, "0") {
				d = d[1:]
			}

			d = d + "e-" + strconv.FormatUint(e, 10)
		} else if !strings.Contains(d, ".") && n >= 999999999999 {
			d = "0x" + strconv.FormatUint(uint64(n), 16)
		} else {
			d = strconv.FormatFloat(n, 'f', -1, 64)
		}
	}

	if len(d) < len(nt.Data) {
		nt.Data = d
		p.changed = true
	}
}

func (p *processor) minifyArrowFunc(af *javascript.ArrowFunction) {
	if p.Has(ArrowFn) {
		if af.FormalParameters != nil && len(af.FormalParameters.FormalParameterList) == 1 && af.FormalParameters.ArrayBindingPattern == nil && af.FormalParameters.ObjectBindingPattern == nil && af.FormalParameters.BindingIdentifier == nil {
			if fp := af.FormalParameters.FormalParameterList[0]; fp.Initializer == nil && fp.SingleNameBinding != nil && fp.ArrayBindingPattern == nil && fp.ObjectBindingPattern == nil {
				af.BindingIdentifier = fp.SingleNameBinding
				af.FormalParameters = nil
				p.changed = true
			}
		}

		if af.FunctionBody != nil {
			if af.FormalParameters != nil {
				if len(af.FormalParameters.FormalParameterList) == 1 && af.FormalParameters.FormalParameterList[0].SingleNameBinding != nil && af.FormalParameters.FormalParameterList[0].Initializer == nil && af.FormalParameters.BindingIdentifier == nil && af.FormalParameters.ArrayBindingPattern == nil && af.FormalParameters.ObjectBindingPattern == nil {
					af.BindingIdentifier = af.FormalParameters.FormalParameterList[0].SingleNameBinding
					af.FormalParameters = nil
					p.changed = true
				}
			}

			expressions, hasReturn := statementsListItemsAsExpressionsAndReturn(af.FunctionBody.StatementList)
			if hasReturn {
				if len(expressions) == 1 {
					af.FunctionBody = nil
					af.AssignmentExpression = &expressions[0]
					p.changed = true
				} else if len(expressions) != 0 {
					af.AssignmentExpression = &javascript.AssignmentExpression{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: expressions,
							Tokens:      af.FunctionBody.Tokens,
						}),
						Tokens: af.FunctionBody.Tokens,
					}
					af.FunctionBody = nil
					p.changed = true
				}
			}
		}
	}
}

func (p *processor) minifyIfToConditional(s *javascript.Statement) {
	if p.Has(IfToConditional) && s.IfStatement != nil && s.IfStatement.ElseStatement != nil {
		last := s.IfStatement.Expression.Expressions[len(s.IfStatement.Expression.Expressions)-1]

		if last.AssignmentOperator != javascript.AssignmentNone || last.ConditionalExpression == nil || last.ArrowFunction != nil || last.AssignmentExpression != nil || last.AssignmentPattern != nil || last.Delegate || last.LeftHandSideExpression != nil || last.Yield {
			return
		}

		var (
			ifExpressions, elseExpressions []javascript.AssignmentExpression
			ifReturn, elseReturn           bool
		)

		if isNonEmptyReturnStatement(&s.IfStatement.Statement) {
			ifReturn = true
			ifExpressions = s.IfStatement.Statement.ExpressionStatement.Expressions
			p.changed = true
		} else if isStatementExpression(&s.IfStatement.Statement) {
			ifExpressions = s.IfStatement.Statement.ExpressionStatement.Expressions
			p.changed = true
		} else if s.IfStatement.Statement.BlockStatement != nil {
			ifExpressions, ifReturn = statementsListItemsAsExpressionsAndReturn(s.IfStatement.Statement.BlockStatement.StatementList)
			p.changed = true
		}

		if len(ifExpressions) == 0 {
			return
		}

		if isNonEmptyReturnStatement(s.IfStatement.ElseStatement) {
			elseReturn = true
			elseExpressions = s.IfStatement.ElseStatement.ExpressionStatement.Expressions
			p.changed = true
		} else if isStatementExpression(s.IfStatement.ElseStatement) {
			elseExpressions = s.IfStatement.ElseStatement.ExpressionStatement.Expressions
			p.changed = true
		} else if s.IfStatement.ElseStatement.BlockStatement != nil {
			elseExpressions, elseReturn = statementsListItemsAsExpressionsAndReturn(s.IfStatement.ElseStatement.BlockStatement.StatementList)
			p.changed = true
		}

		if ifReturn != elseReturn {
			return
		}

		if len(elseExpressions) == 0 {
			return
		} else if len(elseExpressions) == 1 {
			last.ConditionalExpression.False = &elseExpressions[0]
		} else {
			last.ConditionalExpression.False = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
					Expressions: elseExpressions,
					Tokens:      s.IfStatement.ElseStatement.Tokens,
				}),
				Tokens: s.IfStatement.ElseStatement.Tokens,
			}
		}

		if len(ifExpressions) == 1 {
			last.ConditionalExpression.True = &ifExpressions[0]
		} else {
			last.ConditionalExpression.True = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
					Expressions: ifExpressions,
					Tokens:      s.IfStatement.Statement.Tokens,
				}),
				Tokens: s.IfStatement.Statement.Tokens,
			}
		}

		if ifReturn {
			s.Type = javascript.StatementReturn
		}

		s.ExpressionStatement = &s.IfStatement.Expression
		s.IfStatement = nil
		p.changed = true
	}
}

func (p *processor) removeDebugger(s *javascript.Statement) {
	if p.Has(RemoveDebugger) && s.Type == javascript.StatementDebugger {
		s.Type = javascript.StatementNormal
		p.changed = true
	}
}

func (p *processor) minifyBlockToStatement(s *javascript.Statement) {
	if p.Has(BlocksToStatement) && s.BlockStatement != nil {
		if l := len(s.BlockStatement.StatementList); l == 1 {
			if s.BlockStatement.StatementList[0].Statement != nil {
				*s = *s.BlockStatement.StatementList[0].Statement
				p.changed = true
			}
		} else if l > 1 {
			if expressions, hasReturn := statementsListItemsAsExpressionsAndReturn(s.BlockStatement.StatementList); !hasReturn && len(expressions) > 0 {
				*s = javascript.Statement{
					ExpressionStatement: &javascript.Expression{
						Expressions: expressions,
						Tokens:      s.Tokens,
					},
					Tokens: s.Tokens,
				}
				p.changed = true
			}
		}
	}
}

func (p *processor) minifyObjectKeys(pn *javascript.PropertyName) {
	if p.Has(Keys) {
		if ae := pn.ComputedPropertyName; ae != nil && ae.AssignmentOperator == javascript.AssignmentNone && ae.ConditionalExpression != nil && !ae.Yield {
			pe, ok := javascript.UnwrapConditional(ae.ConditionalExpression).(*javascript.PrimaryExpression)
			if ok && pe.Literal != nil && pe.Literal.Type != javascript.TokenRegularExpressionLiteral {
				pn.LiteralPropertyName = pe.Literal
				pn.ComputedPropertyName = nil
			}
		}

		if pn.LiteralPropertyName != nil && pn.LiteralPropertyName.Type == javascript.TokenStringLiteral {
			key, err := javascript.Unquote(pn.LiteralPropertyName.Data)
			if err == nil {
				if isIdentifier(key) {
					pn.LiteralPropertyName.Data = key
					pn.LiteralPropertyName.Type = javascript.TokenIdentifier // This type may not be technically correct, but should not matter.
				} else if isSimpleNumber(key) {
					pn.LiteralPropertyName.Data = key
					pn.LiteralPropertyName.Type = javascript.TokenNumericLiteral
				}
			}
		}
	}
}

func (p *processor) minifyNonHoistableNames(pe *javascript.PrimaryExpression) {
	if p.Has(RemoveExpressionNames) {
		if pe.FunctionExpression != nil {
			pe.FunctionExpression.BindingIdentifier = nil
		} else if pe.ClassExpression != nil {
			pe.ClassExpression.BindingIdentifier = nil
		}
	}
}

func (p *processor) minifyFunctionExpressionAsArrowFunc(ae *javascript.AssignmentExpression) {
	if p.Has(FunctionExpressionToArrowFunc) && ae.AssignmentOperator == javascript.AssignmentNone && ae.ConditionalExpression != nil {
		if fe, ok := javascript.UnwrapConditional(ae.ConditionalExpression).(*javascript.FunctionDeclaration); ok && (fe.Type == javascript.FunctionAsync || fe.Type == javascript.FunctionNormal) {
			s, err := scope.ScriptScope(&javascript.Script{
				StatementList: fe.FunctionBody.StatementList,
			}, nil)
			if err != nil {
				return
			}

			_, hasArguments := s.Bindings["arguments"]
			_, hasThis := s.Bindings["this"]

			if hasArguments || hasThis {
				return
			}

			ae.ArrowFunction = &javascript.ArrowFunction{
				Async:            fe.Type == javascript.FunctionAsync,
				FormalParameters: &fe.FormalParameters,
				FunctionBody:     &fe.FunctionBody,
				Tokens:           fe.Tokens,
			}
			ae.ConditionalExpression = nil

			p.minifyArrowFunc(ae.ArrowFunction)
		}
	}
}

func (p *processor) minifyExpressionParens(e *javascript.Expression) {
	if p.Has(UnwrapParens) {
		e.Expressions = p.minifyParens(e.Expressions)
	}
}

func (p *processor) minifyParenthsizedExpressionParens(pe *javascript.ParenthesizedExpression) {
	if p.Has(UnwrapParens) {
		pe.Expressions = p.minifyParens(pe.Expressions)
	}
}

func (p *Minifier) minifyParens(e []javascript.AssignmentExpression) []javascript.AssignmentExpression {
	for i := 0; i < len(e); i++ {
		if pe := aeAsParen(&e[i]); pe != nil {
			add := make([]javascript.AssignmentExpression, 0, len(pe.Expressions)+len(e)-i-1)
			add = append(add, pe.Expressions...)
			add = append(add, e[i+1:]...)
			e = append(e[:i], add...)
			i += len(pe.Expressions) - 1
		}
	}

	return e
}

func (p *processor) minifyArgumentParens(a *javascript.Argument) {
	if p.Has(UnwrapParens) {
		if pe := aeAsParen(&a.AssignmentExpression); pe != nil && len(pe.Expressions) == 1 {
			a.AssignmentExpression = pe.Expressions[0]
		}
	}
}

func (p *processor) minifyAEParens(ae *javascript.AssignmentExpression) {
	if p.Has(UnwrapParens) {
		if pe := aeAsParen(ae.AssignmentExpression); pe != nil && len(pe.Expressions) == 1 {
			ae.AssignmentExpression = &pe.Expressions[0]
		}
	}
}

func (p *processor) minifyMemberExpressionParens(me *javascript.MemberExpression) {
	if p.Has(UnwrapParens) && meIsSinglePe(me) {
		switch e := javascript.UnwrapConditional(me.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression).(type) {
		case *javascript.PrimaryExpression:
			me.PrimaryExpression = e
		case *javascript.ArrayLiteral:
			me.PrimaryExpression = &javascript.PrimaryExpression{
				ArrayLiteral: e,
				Tokens:       e.Tokens,
			}
		case *javascript.ObjectLiteral:
			me.PrimaryExpression = &javascript.PrimaryExpression{
				ObjectLiteral: e,
				Tokens:        e.Tokens,
			}
		case *javascript.FunctionDeclaration:
			me.PrimaryExpression = &javascript.PrimaryExpression{
				FunctionExpression: e,
				Tokens:             e.Tokens,
			}
		case *javascript.ClassDeclaration:
			me.PrimaryExpression = &javascript.PrimaryExpression{
				ClassExpression: e,
				Tokens:          e.Tokens,
			}
		case *javascript.TemplateLiteral:
			me.PrimaryExpression = &javascript.PrimaryExpression{
				TemplateLiteral: e,
				Tokens:          e.Tokens,
			}
		case *javascript.MemberExpression:
			*me = *e
		}
	}
}

func (p *processor) minifyCallExpressionParens(ce *javascript.CallExpression) {
	if p.Has(UnwrapParens) && meIsSinglePe(ce.MemberExpression) {
		switch e := javascript.UnwrapConditional(ce.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression).(type) {
		case *javascript.CallExpression:
			ce.CallExpression = e
			ce.MemberExpression = nil
		}
	}
}

func (p *processor) minifyConditionExpressionParens(ce *javascript.ConditionalExpression) {
	if p.Has(UnwrapParens) {
		w := javascript.UnwrapConditional(ce)
		switch w := w.(type) {
		case *javascript.LogicalORExpression:
			if ce := isConditionalWrappingAConditional(w.LogicalORExpression, w); ce != nil {
				w.LogicalORExpression = ce.LogicalORExpression
			}
		case *javascript.LogicalANDExpression:
			if ce := isConditionalWrappingAConditional(w.LogicalANDExpression, w); ce != nil {
				w.LogicalANDExpression = &ce.LogicalORExpression.LogicalANDExpression
			}
		case *javascript.BitwiseORExpression:
			if ce := isConditionalWrappingAConditional(w.BitwiseORExpression, w); ce != nil {
				w.BitwiseORExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression
			}
		case *javascript.BitwiseXORExpression:
			if ce := isConditionalWrappingAConditional(w.BitwiseXORExpression, w); ce != nil {
				w.BitwiseXORExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression
			}
		case *javascript.BitwiseANDExpression:
			if ce := isConditionalWrappingAConditional(w.BitwiseANDExpression, w); ce != nil {
				w.BitwiseANDExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression
			}
		case *javascript.EqualityExpression:
			if ce := isConditionalWrappingAConditional(w.EqualityExpression, w); ce != nil {
				w.EqualityExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression
			}
		case *javascript.RelationalExpression:
			if ce := isConditionalWrappingAConditional(w.RelationalExpression, w); ce != nil {
				w.RelationalExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression
			}
		case *javascript.ShiftExpression:
			if ce := isConditionalWrappingAConditional(w.ShiftExpression, w); ce != nil {
				w.ShiftExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression
			}
		case *javascript.AdditiveExpression:
			if ce := isConditionalWrappingAConditional(w.AdditiveExpression, w); ce != nil {
				w.AdditiveExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression
			}
		case *javascript.MultiplicativeExpression:
			if ce := isConditionalWrappingAConditional(w.MultiplicativeExpression, w); ce != nil {
				w.MultiplicativeExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression
			}
		case *javascript.ExponentiationExpression:
			if ce := isConditionalWrappingAConditional(w.ExponentiationExpression, w); ce != nil {
				w.ExponentiationExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression
			}
		case *javascript.UpdateExpression:
			if w.UnaryExpression != nil {
				if ce := isConditionalWrappingAConditional(w.UnaryExpression, w); ce != nil {
					w.UnaryExpression = &ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression
				}
			}
		}
	}
}

func (p *processor) minifyLHSExpressionParens(lhs *javascript.LeftHandSideExpression) {
	if p.Has(UnwrapParens) && lhs.NewExpression != nil && len(lhs.NewExpression.News) == 0 {
		if ce := meAsCE(&lhs.NewExpression.MemberExpression); ce != nil {
			lhs.CallExpression = ce
			lhs.NewExpression = nil
		}
	}
}

func (p *processor) minifyEmptyStatement(jm *javascript.Module) bool {
	for i := 0; i < len(jm.ModuleListItems); i++ {
		if jm.ModuleListItems[i].StatementListItem != nil && isEmptyStatement(jm.ModuleListItems[i].StatementListItem.Statement) {
			jm.ModuleListItems = append(jm.ModuleListItems[:i], jm.ModuleListItems[i+1:]...)
			i--
		}
	}

	return false
}

func (p *processor) minifyLastReturnStatement(f *javascript.FunctionDeclaration) {
	if p.Has(RemoveLastEmptyReturn) {
		removeLastReturnStatement(&f.FunctionBody)
	}
}

func (p *processor) minifyLastReturnStatementInArrowFn(af *javascript.ArrowFunction) {
	if p.Has(RemoveLastEmptyReturn) && af.FunctionBody != nil {
		removeLastReturnStatement(af.FunctionBody)
	}
}

func (p *processor) minifyExpressionRun(jm *javascript.Module) bool {
	if p.Has(CombineExpressionRuns) && len(jm.ModuleListItems) > 1 {
		lastWasExpression := isStatementListItemExpression(jm.ModuleListItems[0].StatementListItem)
		for i := 1; i < len(jm.ModuleListItems); i++ {
			isExpression := isStatementListItemExpression(jm.ModuleListItems[i].StatementListItem)
			if isExpression && lastWasExpression {
				e := jm.ModuleListItems[i-1].StatementListItem.Statement.ExpressionStatement
				e.Expressions = append(e.Expressions, jm.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement.Expressions...)
				jm.ModuleListItems = append(jm.ModuleListItems[:i], jm.ModuleListItems[i+1:]...)
				i--
			} else {
				lastWasExpression = isExpression
			}
		}
	}

	return false
}

func (p *processor) fixFirstExpression(jm *javascript.Module) bool {
	if p.Has(UnwrapParens) {
		for n := range jm.ModuleListItems {
			if isStatementListItemExpression(jm.ModuleListItems[n].StatementListItem) {
				fixWrapping(jm.ModuleListItems[n].StatementListItem.Statement)
			}
		}
	}

	return false
}

func (p *processor) fixFirstArrowFuncExpression(af *javascript.ArrowFunction) {
	if p.Has(UnwrapParens) && af.AssignmentExpression != nil {
		fixWrapping(&javascript.Statement{
			ExpressionStatement: &javascript.Expression{
				Expressions: []javascript.AssignmentExpression{*af.AssignmentExpression},
			},
		})
	}
}

func (p *processor) minifyRemoveDeadCode(b *javascript.Block) {
	if p.Has(RemoveDeadCode) {
		retPos := -1
		for n := range b.StatementList {
			if isReturnStatement(b.StatementList[n].Statement) {
				retPos = n

				break
			}
		}

		if retPos >= 0 {
			for i := retPos + 1; i < len(b.StatementList); i++ {
				if !isHoistable(&b.StatementList[i]) {
					b.StatementList = append(b.StatementList[:i], b.StatementList[i+1:]...)
					i--
				}
			}
		}
	}
}

func (p *processor) minifyLexical(jm *javascript.Module) bool {
	if p.Has(MergeLexical) {
		last := bindableNone

		for i := 0; i < len(jm.ModuleListItems); i++ {
			next := sliBindable(jm.ModuleListItems[i].StatementListItem)
			if last == next {
				switch next {
				case bindableConst, bindableLet:
					ld := jm.ModuleListItems[i-1].StatementListItem.Declaration.LexicalDeclaration
					ld.BindingList = append(ld.BindingList, jm.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList...)
					jm.ModuleListItems = append(jm.ModuleListItems[:i], jm.ModuleListItems[i+1:]...)
					i--
				case bindableVar:
					vs := jm.ModuleListItems[i-1].StatementListItem.Statement.VariableStatement
					vs.VariableDeclarationList = append(vs.VariableDeclarationList, jm.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList...)
					jm.ModuleListItems = append(jm.ModuleListItems[:i], jm.ModuleListItems[i+1:]...)
					i--
				}
			}

			last = next
		}
	}

	return false
}

func (p *processor) minifyExpressionsBetweenLexicals(jm *javascript.Module) bool {
	if p.Has(MergeLexical) && p.Has(CombineExpressionRuns) {
		for i := 2; i < len(jm.ModuleListItems); i++ {
			if last := sliBindable(jm.ModuleListItems[i-2].StatementListItem); (last == bindableLet || last == bindableConst || last == bindableVar) && isStatementListItemExpression(jm.ModuleListItems[i-1].StatementListItem) && sliBindable(jm.ModuleListItems[i].StatementListItem) == last {
				var flbs, lbs []javascript.LexicalBinding

				if last == bindableVar {
					flbs = jm.ModuleListItems[i-2].StatementListItem.Statement.VariableStatement.VariableDeclarationList
					lbs = jm.ModuleListItems[i].StatementListItem.Statement.VariableStatement.VariableDeclarationList
				} else {
					flbs = jm.ModuleListItems[i-2].StatementListItem.Declaration.LexicalDeclaration.BindingList
					lbs = jm.ModuleListItems[i].StatementListItem.Declaration.LexicalDeclaration.BindingList
				}

				back := 0

				for n := range lbs {
					if pe := aeAsParen(lbs[n].Initializer); pe != nil {
						pe.Expressions = append(jm.ModuleListItems[i-1].StatementListItem.Statement.ExpressionStatement.Expressions, pe.Expressions...)
						back = 1

						break
					} else if !isSimpleAE(lbs[n].Initializer) {
						lbs[n].Initializer = &javascript.AssignmentExpression{
							ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
								Expressions: append(jm.ModuleListItems[i-1].StatementListItem.Statement.ExpressionStatement.Expressions, *lbs[n].Initializer),
							}),
						}
						back = 1

						break
					}
				}

				flbs = append(flbs, lbs...)

				if last == bindableVar {
					jm.ModuleListItems[i-2].StatementListItem.Statement.VariableStatement.VariableDeclarationList = flbs
				} else {
					jm.ModuleListItems[i-2].StatementListItem.Declaration.LexicalDeclaration.BindingList = flbs
				}

				jm.ModuleListItems = append(jm.ModuleListItems[:i-back], jm.ModuleListItems[i+1:]...)
				i--
			}
		}
	}

	return false
}
