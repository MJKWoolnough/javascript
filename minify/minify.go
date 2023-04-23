package minify

import (
	"strconv"
	"strings"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

type Minifier struct {
	literals, numbers, arrowFn, ifToConditional bool
}

func New(opts ...Option) *Minifier {
	m := new(Minifier)
	for _, opt := range opts {
		opt(m)
	}
	return m
}

type walker struct {
	*Minifier
}

func (w *walker) Handle(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.PrimaryExpression:
		w.minifyLiterals(t)
		w.minifyNumbers(t)
	case *javascript.ArrowFunction:
		w.minifyArrowFunc(t)
	case *javascript.Statement:
		w.minifyIfToConditional(t)
	}
	return walk.Walk(t, w)
}

func (m *Minifier) Process(jm *javascript.Module) {
	walk.Walk(jm, &walker{Minifier: m})
}

func (m *Minifier) minifyLiterals(pe *javascript.PrimaryExpression) {
	if m.literals {
		if pe.Literal != nil {
			switch pe.Literal.Data {
			case "true":
				pe.Literal.Data = "!0"
			case "false":
				pe.Literal.Data = "!1"
			}
		} else if pe.IdentifierReference != nil && pe.IdentifierReference.Data == "undefined" {
			pe.IdentifierReference.Data = "void 0"
		}
	}
}

func (m *Minifier) minifyNumbers(pe *javascript.PrimaryExpression) {
	if m.numbers && pe.Literal != nil && pe.Literal.Type == javascript.TokenNumericLiteral {
		d := pe.Literal.Data
		d = strings.ReplaceAll(d, "_", "")
		if len(d) < len(pe.Literal.Data) {
			pe.Literal.Data = d
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
				f, err := strconv.ParseFloat(pe.Literal.Data, 64)
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
		if len(d) < len(pe.Literal.Data) {
			pe.Literal.Data = d
		}
	}
}

func (m *Minifier) minifyArrowFunc(af *javascript.ArrowFunction) {
	if m.arrowFn && af.FunctionBody != nil {
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
				}
				af.FunctionBody = nil
			}
		}
	}
}

func (m *Minifier) minifyIfToConditional(s *javascript.Statement) {
	if s.IfStatement != nil && s.IfStatement.ElseStatement != nil {
		last := s.IfStatement.Expression.Expressions[len(s.IfStatement.Expression.Expressions)-1]
		if last.AssignmentOperator != javascript.AssignmentNone || last.ConditionalExpression == nil || last.ArrowFunction != nil || last.AssignmentExpression != nil || last.AssignmentPattern != nil || last.Delegate || last.LeftHandSideExpression != nil || last.Yield {
			return
		}
		var (
			ifExpressions, elseExpressions []javascript.AssignmentExpression
			ifReturn, elseReturn           bool
		)
		if isReturnStatement(&s.IfStatement.Statement) {
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
		if isReturnStatement(s.IfStatement.ElseStatement) {
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
				}),
			}
		}
		if len(ifExpressions) == 1 {
			last.ConditionalExpression.True = &ifExpressions[0]
		} else {
			last.ConditionalExpression.True = &javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
					Expressions: ifExpressions,
				}),
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
	return s != nil && s.Type == javascript.StatementReturn && s.ExpressionStatement != nil
}

func statementsListItemsAsExpressionsAndReturn(sli []javascript.StatementListItem) ([]javascript.AssignmentExpression, bool) {
	var (
		expressions []javascript.AssignmentExpression
		hasReturn   bool
	)
	for n := range sli {
		s := &sli[n]
		if isReturnStatement(s.Statement) {
			expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
			hasReturn = true
			break
		} else if !isSLIExpression(s) {
			if s.Statement != nil && isEmptyStatement(s.Statement) {
				continue
			}
			return nil, false
		}
		expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
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
	return s.Type == javascript.StatementNormal && s.BlockStatement == nil && s.VariableStatement == nil && s.ExpressionStatement == nil && s.IfStatement == nil && s.IterationStatementDo == nil && s.IterationStatementFor == nil && s.IterationStatementWhile == nil && s.SwitchStatement == nil && s.WithStatement == nil && s.LabelledItemFunction == nil && s.LabelledItemStatement == nil && s.TryStatement == nil
}
