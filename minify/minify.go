package minify

import (
	"strconv"
	"strings"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
)

type Minifier struct {
	literals, numbers, arrowFn bool
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
		hasReturn := false
		expressions := make([]javascript.AssignmentExpression, 0)
		for _, s := range af.FunctionBody.StatementList {
			if s.Declaration != nil {
				return
			} else if s.Statement != nil {
				if s.Statement.Type == javascript.StatementReturn {
					hasReturn = true
					expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
					break
				} else if s.Statement.ExpressionStatement == nil {
					if isEmptyStatement(s.Statement) {
						continue
					}
					return
				}
				expressions = append(expressions, s.Statement.ExpressionStatement.Expressions...)
			}
		}
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

func isEmptyStatement(s *javascript.Statement) bool {
	return s.Type == javascript.StatementNormal && s.BlockStatement == nil && s.VariableStatement == nil && s.ExpressionStatement == nil && s.IfStatement == nil && s.IterationStatementDo == nil && s.IterationStatementFor == nil && s.IterationStatementWhile == nil && s.SwitchStatement == nil && s.WithStatement == nil && s.LabelledItemFunction == nil && s.LabelledItemStatement == nil && s.TryStatement == nil
}
