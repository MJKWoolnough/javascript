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
	for _, opt := range opts {
		m.Option |= opt
	}
	return m
}

type walker struct {
	*Minifier
}

func (w *walker) Handle(t javascript.Type) error {
	if err := walk.Walk(t, w); err != nil {
		return err
	}
	switch t := t.(type) {
	case *javascript.TemplateLiteral:
		w.minifyTemplates(t)
	case *javascript.PrimaryExpression:
		w.minifyLiterals(t)
		w.minifyNonHoistableNames(t)
	case *javascript.ArrowFunction:
		w.minifyArrowFunc(t)
		w.minifyLastReturnStatementInArrowFn(t)
		w.fixFirstArrowFuncExpression(t)
	case *javascript.Statement:
		w.minifyBlockToStatement(t)
		w.minifyIfToConditional(t)
		w.removeDebugger(t)
	case *javascript.PropertyName:
		w.minifyObjectKeys(t)
	case *javascript.AssignmentExpression:
		w.minifyFunctionExpressionAsArrowFunc(t)
		w.minifyAEParens(t)
	case *javascript.ParenthesizedExpression:
		w.minifyParenthsizedExpressionParens(t)
	case *javascript.Expression:
		w.minifyExpressionParens(t)
	case *javascript.Argument:
		w.minifyArgumentParens(t)
	case *javascript.MemberExpression:
		w.minifyMemberExpressionParens(t)
	case *javascript.CallExpression:
		w.minifyCallExpressionParens(t)
	case *javascript.LeftHandSideExpression:
		w.minifyLHSExpressionParens(t)
	case *javascript.Block:
		w.minifyEmptyStatementInBlock(t)
		w.minifyExpressionRunInBlock(t)
		w.fixFirstExpressionInBlock(t)
	case *javascript.Module:
		w.minifyEmptyStatementInModule(t)
		w.minifyExpressionRunInModule(t)
		w.fixFirstExpressionInModule(t)
	case *javascript.FunctionDeclaration:
		w.minifyLastReturnStatement(t)
	case *javascript.ConditionalExpression:
		w.minifyConditionExpressionParens(t)
	}

	return nil
}

func (m *Minifier) Process(jm *javascript.Module) {
	walk.Walk(jm, &walker{Minifier: m})
	if m.Has(RenameIdentifiers) {
		renameIdentifiers(jm)
	}
}

func minifyTemplate(t *javascript.Token) {
	if t != nil {
		str, err := javascript.UnquoteTemplate(t.Data)
		if err != nil {
			return
		}
		res := javascript.QuoteTemplate(str, javascript.TokenTypeToTemplateType(t.Type))
		if len(res) < len(t.Data) {
			t.Data = res
		}
	}
}

func (m *Minifier) minifyTemplates(t *javascript.TemplateLiteral) {
	if m.Has(Literals) {
		minifyTemplate(t.NoSubstitutionTemplate)
		minifyTemplate(t.TemplateHead)
		for _, m := range t.TemplateMiddleList {
			minifyTemplate(m)
		}
		minifyTemplate(t.TemplateTail)
	}
}

func (m *Minifier) minifyLiterals(pe *javascript.PrimaryExpression) {
	if m.Has(Literals) {
		if pe.Literal != nil {
			switch pe.Literal.Type {
			case javascript.TokenBooleanLiteral:
				switch pe.Literal.Data {
				case "true":
					pe.Literal.Data = "!0"
				case "false":
					pe.Literal.Data = "!1"
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
				}
			case javascript.TokenNumericLiteral:
				minifyNumbers(pe.Literal)
			}
		} else if pe.IdentifierReference != nil && pe.IdentifierReference.Data == "undefined" {
			pe.IdentifierReference.Data = "void 0"
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
			}
		}
	}
}

func minifyNumbers(nt *javascript.Token) {
	d := nt.Data
	d = strings.ReplaceAll(d, "_", "")
	if len(d) < len(nt.Data) {
		nt.Data = d
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
	}
}

func (m *Minifier) minifyArrowFunc(af *javascript.ArrowFunction) {
	if m.Has(ArrowFn) {
		if af.FormalParameters != nil && len(af.FormalParameters.FormalParameterList) == 1 && af.FormalParameters.ArrayBindingPattern == nil && af.FormalParameters.ObjectBindingPattern == nil && af.FormalParameters.BindingIdentifier == nil {
			if fp := af.FormalParameters.FormalParameterList[0]; fp.Initializer == nil && fp.SingleNameBinding != nil && fp.ArrayBindingPattern == nil && fp.ObjectBindingPattern == nil {
				af.BindingIdentifier = fp.SingleNameBinding
				af.FormalParameters = nil
			}
		}
		if af.FunctionBody != nil {
			if af.FormalParameters != nil {
				if len(af.FormalParameters.FormalParameterList) == 1 && af.FormalParameters.FormalParameterList[0].SingleNameBinding != nil && af.FormalParameters.FormalParameterList[0].Initializer == nil && af.FormalParameters.BindingIdentifier == nil && af.FormalParameters.ArrayBindingPattern == nil && af.FormalParameters.ObjectBindingPattern == nil {
					af.BindingIdentifier = af.FormalParameters.FormalParameterList[0].SingleNameBinding
					af.FormalParameters = nil
				}
			}
			expressions, hasReturn := statementsListItemsAsExpressionsAndReturn(af.FunctionBody.StatementList)
			if hasReturn {
				if len(expressions) == 1 {
					af.FunctionBody = nil
					af.AssignmentExpression = &expressions[0]
				} else if len(expressions) != 0 {
					af.AssignmentExpression = &javascript.AssignmentExpression{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: expressions,
							Tokens:      af.FunctionBody.Tokens,
						}),
						Tokens: af.FunctionBody.Tokens,
					}
					af.FunctionBody = nil
				}
			}
		}
	}
}

func (m *Minifier) minifyIfToConditional(s *javascript.Statement) {
	if m.Has(IfToConditional) && s.IfStatement != nil && s.IfStatement.ElseStatement != nil {
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
		} else if isStatementExpression(&s.IfStatement.Statement) {
			ifExpressions = s.IfStatement.Statement.ExpressionStatement.Expressions
		} else if s.IfStatement.Statement.BlockStatement != nil {
			ifExpressions, ifReturn = statementsListItemsAsExpressionsAndReturn(s.IfStatement.Statement.BlockStatement.StatementList)
		}
		if len(ifExpressions) == 0 {
			return
		}
		if isNonEmptyReturnStatement(s.IfStatement.ElseStatement) {
			elseReturn = true
			elseExpressions = s.IfStatement.ElseStatement.ExpressionStatement.Expressions
		} else if isStatementExpression(s.IfStatement.ElseStatement) {
			elseExpressions = s.IfStatement.ElseStatement.ExpressionStatement.Expressions
		} else if s.IfStatement.ElseStatement.BlockStatement != nil {
			elseExpressions, elseReturn = statementsListItemsAsExpressionsAndReturn(s.IfStatement.ElseStatement.BlockStatement.StatementList)
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
	return s.Declaration == nil && isStatementExpression(s.Statement)
}

func isStatementExpression(s *javascript.Statement) bool {
	return s != nil && s.Type == javascript.StatementNormal && s.ExpressionStatement != nil
}

func isEmptyStatement(s *javascript.Statement) bool {
	return s != nil && s.Type == javascript.StatementNormal && s.BlockStatement == nil && s.VariableStatement == nil && s.ExpressionStatement == nil && s.IfStatement == nil && s.IterationStatementDo == nil && s.IterationStatementFor == nil && s.IterationStatementWhile == nil && s.SwitchStatement == nil && s.WithStatement == nil && s.LabelledItemFunction == nil && s.LabelledItemStatement == nil && s.TryStatement == nil
}

func isHoistable(s *javascript.StatementListItem) bool {
	return (s.Statement != nil && (s.Statement.VariableStatement.VariableDeclarationList != nil || s.Statement.LabelledItemFunction != nil)) || (s.Declaration != nil && (s.Declaration.FunctionDeclaration != nil || s.Declaration.ClassDeclaration != nil))
}

func (m *Minifier) removeDebugger(s *javascript.Statement) {
	if m.Has(RemoveDebugger) && s.Type == javascript.StatementDebugger {
		s.Type = javascript.StatementNormal
	}
}

func (m *Minifier) minifyBlockToStatement(s *javascript.Statement) {
	if m.Has(BlocksToStatement) && s.BlockStatement != nil {
		if l := len(s.BlockStatement.StatementList); l == 1 {
			if s.BlockStatement.StatementList[0].Statement != nil {
				*s = *s.BlockStatement.StatementList[0].Statement
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
			}
		}
	}
}

func (m *Minifier) minifyObjectKeys(p *javascript.PropertyName) {
	if m.Has(Keys) {
		if ae := p.ComputedPropertyName; ae != nil && ae.AssignmentOperator == javascript.AssignmentNone && ae.ConditionalExpression != nil && !ae.Yield {
			pe, ok := javascript.UnwrapConditional(ae.ConditionalExpression).(*javascript.PrimaryExpression)
			if ok && pe.Literal != nil && pe.Literal.Type != javascript.TokenRegularExpressionLiteral {
				p.LiteralPropertyName = pe.Literal
				p.ComputedPropertyName = nil
			}
		}
		if p.LiteralPropertyName != nil && p.LiteralPropertyName.Type == javascript.TokenStringLiteral {
			key, err := javascript.Unquote(p.LiteralPropertyName.Data)
			if err == nil {
				if isIdentifier(key) {
					p.LiteralPropertyName.Data = key
					p.LiteralPropertyName.Type = javascript.TokenIdentifier // This type may not be technically correct, but should not matter.
				} else if isSimpleNumber(key) {
					p.LiteralPropertyName.Data = key
					p.LiteralPropertyName.Type = javascript.TokenNumericLiteral
				}
			}
		}
	}
}

func (m *Minifier) minifyNonHoistableNames(pe *javascript.PrimaryExpression) {
	if m.Has(RemoveExpressionNames) {
		if pe.FunctionExpression != nil {
			pe.FunctionExpression.BindingIdentifier = nil
		} else if pe.ClassExpression != nil {
			pe.ClassExpression.BindingIdentifier = nil
		}
	}
}

func (m *Minifier) minifyFunctionExpressionAsArrowFunc(ae *javascript.AssignmentExpression) {
	if m.Has(FunctionExpressionToArrowFunc) && ae.AssignmentOperator == javascript.AssignmentNone && ae.ConditionalExpression != nil {
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
			m.minifyArrowFunc(ae.ArrowFunction)
		}
	}
}

func (m *Minifier) minifyExpressionParens(e *javascript.Expression) {
	if m.Has(UnwrapParens) {
		e.Expressions = m.minifyParens(e.Expressions)
	}
}

func (m *Minifier) minifyParenthsizedExpressionParens(pe *javascript.ParenthesizedExpression) {
	if m.Has(UnwrapParens) {
		pe.Expressions = m.minifyParens(pe.Expressions)
	}
}

func aeIsCE(ae *javascript.AssignmentExpression) bool {
	return ae != nil && ae.ConditionalExpression != nil && ae.AssignmentOperator == javascript.AssignmentNone && !ae.Yield
}

func aeAsParen(ae *javascript.AssignmentExpression) (*javascript.ParenthesizedExpression, bool) {
	if aeIsCE(ae) {
		pe, ok := javascript.UnwrapConditional(ae.ConditionalExpression).(*javascript.ParenthesizedExpression)
		return pe, ok
	}
	return nil, false
}

func (w *Minifier) minifyParens(e []javascript.AssignmentExpression) []javascript.AssignmentExpression {
	for i := 0; i < len(e); i++ {
		if pe, ok := aeAsParen(&e[i]); ok {
			add := make([]javascript.AssignmentExpression, 0, len(pe.Expressions)+len(e)-i-1)
			add = append(add, pe.Expressions...)
			add = append(add, e[i+1:]...)
			e = append(e[:i], add...)
			i += len(pe.Expressions) - 1
		}
	}
	return e
}

func (m *Minifier) minifyArgumentParens(a *javascript.Argument) {
	if m.Has(UnwrapParens) {
		if pe, ok := aeAsParen(&a.AssignmentExpression); ok && len(pe.Expressions) == 1 {
			a.AssignmentExpression = pe.Expressions[0]
		}
	}
}

func (m *Minifier) minifyAEParens(ae *javascript.AssignmentExpression) {
	if m.Has(UnwrapParens) {
		if pe, ok := aeAsParen(ae.AssignmentExpression); ok && len(pe.Expressions) == 1 {
			ae.AssignmentExpression = &pe.Expressions[0]
		}
	}
}

func meIsSinglePe(me *javascript.MemberExpression) bool {
	return me != nil && me.PrimaryExpression != nil && me.PrimaryExpression.ParenthesizedExpression != nil && len(me.PrimaryExpression.ParenthesizedExpression.Expressions) == 1 && aeIsCE(&me.PrimaryExpression.ParenthesizedExpression.Expressions[0])
}

func (m *Minifier) minifyMemberExpressionParens(me *javascript.MemberExpression) {
	if m.Has(UnwrapParens) && meIsSinglePe(me) {
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

func (m *Minifier) minifyCallExpressionParens(ce *javascript.CallExpression) {
	if m.Has(UnwrapParens) && meIsSinglePe(ce.MemberExpression) {
		switch e := javascript.UnwrapConditional(ce.MemberExpression.PrimaryExpression.ParenthesizedExpression.Expressions[0].ConditionalExpression).(type) {
		case *javascript.CallExpression:
			ce.CallExpression = e
			ce.MemberExpression = nil
		}
	}
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

func (m *Minifier) minifyConditionExpressionParens(ce *javascript.ConditionalExpression) {
	if m.Has(UnwrapParens) {
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

func (m *Minifier) minifyLHSExpressionParens(lhs *javascript.LeftHandSideExpression) {
	if m.Has(UnwrapParens) && lhs.NewExpression != nil && lhs.NewExpression.News == 0 {
		ce := meAsCE(&lhs.NewExpression.MemberExpression)
		if ce != nil {
			lhs.CallExpression = ce
			lhs.NewExpression = nil
		}
	}
}

func (m *Minifier) minifyEmptyStatementInBlock(b *javascript.Block) {
	for i := 0; i < len(b.StatementList); i++ {
		if isEmptyStatement(b.StatementList[i].Statement) {
			b.StatementList = append(b.StatementList[:i], b.StatementList[i+1:]...)
			i--
		}
	}
}

func (m *Minifier) minifyEmptyStatementInModule(jm *javascript.Module) {
	for i := 0; i < len(jm.ModuleListItems); i++ {
		if jm.ModuleListItems[i].StatementListItem != nil && isEmptyStatement(jm.ModuleListItems[i].StatementListItem.Statement) {
			jm.ModuleListItems = append(jm.ModuleListItems[:i], jm.ModuleListItems[i+1:]...)
			i--
		}
	}
}

func removeLastReturnStatement(b *javascript.Block) {
	if len(b.StatementList) > 0 {
		s := b.StatementList[len(b.StatementList)-1].Statement
		if isReturnStatement(s) && s.ExpressionStatement == nil {
			b.StatementList = b.StatementList[:len(b.StatementList)-1]
		}
	}
}

func (m *Minifier) minifyLastReturnStatement(f *javascript.FunctionDeclaration) {
	if m.Has(RemoveLastEmptyReturn) {
		removeLastReturnStatement(&f.FunctionBody)
	}
}

func (m *Minifier) minifyLastReturnStatementInArrowFn(af *javascript.ArrowFunction) {
	if m.Has(RemoveLastEmptyReturn) && af.FunctionBody != nil {
		removeLastReturnStatement(af.FunctionBody)
	}
}

func (m *Minifier) minifyExpressionRunInBlock(b *javascript.Block) {
	if m.Has(CombineExpressionRuns) && len(b.StatementList) > 1 {
		lastWasExpression := isStatementExpression(b.StatementList[0].Statement)
		for i := 1; i < len(b.StatementList); i++ {
			isExpression := isStatementExpression(b.StatementList[i].Statement)
			if isExpression && lastWasExpression {
				e := b.StatementList[i-1].Statement.ExpressionStatement
				e.Expressions = append(e.Expressions, b.StatementList[i].Statement.ExpressionStatement.Expressions...)
				b.StatementList = append(b.StatementList[:i], b.StatementList[i+1:]...)
				i--
			} else {
				lastWasExpression = isExpression
			}
		}
	}
}

func isStatementListItemExpression(s *javascript.StatementListItem) bool {
	return s != nil && isStatementExpression(s.Statement)
}

func (m *Minifier) minifyExpressionRunInModule(jm *javascript.Module) {
	if m.Has(CombineExpressionRuns) && len(jm.ModuleListItems) > 1 {
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

func (m *Minifier) fixFirstExpressionInBlock(b *javascript.Block) {
	if m.Has(UnwrapParens) {
		for n := range b.StatementList {
			if isStatementListItemExpression(&b.StatementList[n]) {
				fixWrapping(b.StatementList[n].Statement)
			}
		}
	}
}

func (m *Minifier) fixFirstExpressionInModule(jm *javascript.Module) {
	if m.Has(UnwrapParens) {
		for n := range jm.ModuleListItems {
			if isStatementListItemExpression(jm.ModuleListItems[n].StatementListItem) {
				fixWrapping(jm.ModuleListItems[n].StatementListItem.Statement)
			}
		}
	}
}

func (m *Minifier) fixFirstArrowFuncExpression(af *javascript.ArrowFunction) {
	if m.Has(UnwrapParens) && af.AssignmentExpression != nil {
		fixWrapping(&javascript.Statement{
			ExpressionStatement: &javascript.Expression{
				Expressions: []javascript.AssignmentExpression{*af.AssignmentExpression},
			},
		})
	}
}
