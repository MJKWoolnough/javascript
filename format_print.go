package javascript

import (
	"io"
)

var (
	blockOpen                    = []byte{'{'}
	blockClose                   = []byte{'}'}
	commaSep                     = []byte{',', ' '}
	commaSepNL                   = []byte{',', '\n'}
	doubleNewLine                = []byte{'\n', '\n'}
	newLine                      = doubleNewLine[:1]
	conditionalStart             = []byte{' ', '?', ' '}
	conditionalSep               = []byte{' ', ':', ' '}
	labelPost                    = conditionalSep[1:]
	semiColon                    = []byte{';'}
	ifOpen                       = []byte{'i', 'f', ' ', '('}
	parenCloseSpace              = []byte{')', ' '}
	parenClose                   = parenCloseSpace[:1]
	elseOpen                     = []byte{' ', 'e', 'l', 's', 'e', ' '}
	doOpen                       = []byte{'d', 'o', ' '}
	doWhileOpen                  = []byte{' ', 'w', 'h', 'i', 'l', 'e', ' ', '('}
	doWhileClose                 = []byte{')', ';'}
	whileOpen                    = doWhileOpen[1:]
	forOpen                      = []byte{'f', 'o', 'r', ' ', '('}
	forAwaitOpen                 = []byte{'f', 'o', 'r', ' ', 'a', 'w', 'a', 'i', 't', ' ', '('}
	switchOpen                   = []byte{'s', 'w', 'i', 't', 'c', 'h', ' ', '('}
	switchClose                  = []byte{')', ' ', '{'}
	caseOpen                     = []byte{'c', 'a', 's', 'e', ' '}
	caseClose                    = labelPost[:1]
	defaultCase                  = []byte{'d', 'e', 'f', 'a', 'u', 'l', 't', ':'}
	withOpen                     = []byte{'w', 'i', 't', 'h', ' ', '('}
	forIn                        = []byte{' ', 'i', 'n', ' '}
	forOf                        = []byte{' ', 'o', 'f', ' '}
	varOpen                      = []byte{'v', 'a', 'r', ' '}
	letOpen                      = []byte{'l', 'e', 't', ' '}
	constOpen                    = []byte{'c', 'o', 'n', 's', 't', ' '}
	funcOpen                     = []byte{'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', ' '}
	asyncFuncOpen                = []byte{'a', 's', 'y', 'n', 'c', ' ', 'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', ' '}
	genFuncOpen                  = []byte{'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', '*', ' '}
	asyncGenFuncOpen             = []byte{'a', 's', 'y', 'n', 'c', ' ', 'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', '*', ' '}
	parenOpen                    = []byte{'('}
	tryOpen                      = []byte{'t', 'r', 'y', ' '}
	catchParenOpen               = []byte{' ', 'c', 'a', 't', 'c', 'h', ' ', '('}
	catchOpen                    = catchParenOpen[:7]
	finallyOpen                  = []byte{' ', 'f', 'i', 'n', 'a', 'l', 'l', 'y', ' '}
	classOpen                    = []byte{'c', 'l', 'a', 's', 's', ' '}
	extends                      = []byte{'e', 'x', 't', 'e', 'n', 'd', 's', ' '}
	assignment                   = []byte{' ', '=', ' '}
	assignmentMultiply           = []byte{' ', '*', '=', ' '}
	assignmentDivide             = []byte{' ', '/', '=', ' '}
	assignmentRemainder          = []byte{' ', '%', '=', ' '}
	assignmentAdd                = []byte{' ', '+', '=', ' '}
	assignmentSubtract           = []byte{' ', '-', '=', ' '}
	assignmentLeftShift          = []byte{' ', '<', '<', '=', ' '}
	assignmentSignRightShift     = []byte{' ', '>', '>', '=', ' '}
	assignmentZeroRightShift     = []byte{' ', '>', '>', '>', '=', ' '}
	assignmentAND                = []byte{' ', '&', '=', ' '}
	assignmentXOR                = []byte{' ', '^', '=', ' '}
	assignmentOR                 = []byte{' ', '|', '=', ' '}
	assignmentExponentiation     = []byte{' ', '*', '*', '=', ' '}
	assignmentLogicalAnd         = []byte{' ', '&', '&', '=', ' '}
	assignmentLogicalOr          = []byte{' ', '|', '|', '=', ' '}
	assignmentNullish            = []byte{' ', '?', '?', '=', ' '}
	yield                        = []byte{'y', 'i', 'e', 'l', 'd', ' '}
	delegate                     = []byte{'*', ' '}
	ellipsis                     = []byte{'.', '.', '.'}
	bracketOpen                  = []byte{'['}
	bracketClose                 = []byte{']'}
	methodAsyncGenerator         = []byte{'a', 's', 'y', 'n', 'c', ' ', '*', ' '}
	methodAsync                  = methodAsyncGenerator[0:6]
	methodGenerator              = methodAsyncGenerator[6:8]
	methodGet                    = []byte{'g', 'e', 't', ' '}
	methodSet                    = []byte{'s', 'e', 't', ' '}
	methodStatic                 = []byte{'s', 't', 'a', 't', 'i', 'c', ' '}
	arrow                        = []byte{'=', '>', ' '}
	news                         = []byte{'n', 'e', 'w', ' '}
	super                        = []byte{'s', 'u', 'p', 'e', 'r'}
	colonSep                     = []byte{':', ' '}
	logicalOR                    = []byte{' ', '|', '|', ' '}
	newTarget                    = []byte{'n', 'e', 'w', '.', 't', 'a', 'r', 'g', 'e', 't'}
	importMeta                   = []byte{'i', 'm', 'p', 'o', 'r', 't', '.', 'm', 'e', 't', 'a'}
	dot                          = ellipsis[:1]
	logicalAND                   = []byte{' ', '&', '&', ' '}
	this                         = []byte{'t', 'h', 'i', 's'}
	bitwiseOR                    = []byte{' ', '|', ' '}
	bitwiseXOR                   = []byte{' ', '^', ' '}
	bitwiseAND                   = []byte{' ', '&', ' '}
	equalityEqual                = []byte{' ', '=', '=', ' '}
	equalityNotEqual             = []byte{' ', '!', '=', ' '}
	equalityStrictEqual          = []byte{' ', '=', '=', '=', ' '}
	equalityStrictNotEqual       = []byte{' ', '!', '=', '=', ' '}
	relationshipLessThan         = []byte{' ', '<', ' '}
	relationshipGreaterThan      = []byte{' ', '>', ' '}
	relationshipLessThanEqual    = []byte{' ', '<', '=', ' '}
	relationshipGreaterThanEqual = []byte{' ', '>', '=', ' '}
	relationshipInstanceOf       = []byte{' ', 'i', 'n', 's', 't', 'a', 'n', 'c', 'e', 'o', 'f', ' '}
	relationshipIn               = forIn
	shiftLeft                    = []byte{' ', '<', '<', ' '}
	shiftRight                   = []byte{' ', '>', '>', ' '}
	shiftUnsignedRight           = []byte{' ', '>', '>', '>', ' '}
	additiveAdd                  = []byte{' ', '+', ' '}
	additiveMinus                = []byte{' ', '-', ' '}
	multiplicativeMultiply       = []byte{' ', '*', ' '}
	multiplicativeDivide         = []byte{' ', '/', ' '}
	multiplicativeRemainder      = []byte{' ', '%', ' '}
	exponentionation             = []byte{' ', '*', '*', ' '}
	unaryDelete                  = []byte{'d', 'e', 'l', 'e', 't', 'e', ' '}
	unaryVoid                    = []byte{'v', 'o', 'i', 'd', ' '}
	unaryTypeOf                  = []byte{'t', 'y', 'p', 'e', 'o', 'f', ' '}
	unaryAdd                     = []byte{'+'}
	unaryMinus                   = []byte{'-'}
	unaryBitwiseNot              = []byte{'~'}
	unaryLogicalNot              = []byte{'!'}
	unaryAwait                   = []byte{'a', 'w', 'a', 'i', 't', ' '}
	updateIncrement              = []byte{'+', '+'}
	updateDecrement              = []byte{'-', '-'}
	importc                      = []byte{'i', 'm', 'p', 'o', 'r', 't', ' '}
	from                         = []byte{' ', 'f', 'r', 'o', 'm', ' '}
	exportAll                    = exponentionation[1:2]
	exportd                      = []byte{'e', 'x', 'p', 'o', 'r', 't', ' ', 'd', 'e', 'f', 'a', 'u', 'l', 't', ' '}
	exportc                      = exportd[:7]
	namespaceImport              = []byte{'*', ' ', 'a', 's', ' '}
	as                           = namespaceImport[1:]
	importCall                   = []byte{'i', 'm', 'p', 'o', 'r', 't', '('}
	optionalChain                = []byte{'?', '.'}
	coalesceOperator             = []byte{' ', '?', '?', ' '}
	space                        = []byte{' '}
)

func (s Script) printSource(w io.Writer, v bool) {
	if len(s.StatementList) > 0 {
		s.StatementList[0].printSource(w, v)

		for _, stmt := range s.StatementList[1:] {
			w.Write(doubleNewLine)
			stmt.printSource(w, v)
		}
	}
}

func (s StatementListItem) printSource(w io.Writer, v bool) {
	if s.Statement != nil {
		s.Statement.printSource(w, v)
	} else if s.Declaration != nil {
		s.Declaration.printSource(w, v)
	}
}

func (s Statement) printSource(w io.Writer, v bool) {
	switch s.Type {
	case StatementNormal:
		if s.BlockStatement != nil {
			s.BlockStatement.printSource(w, v)
		} else if s.VariableStatement != nil {
			s.VariableStatement.printSource(w, v)
		} else if s.ExpressionStatement != nil {
			s.ExpressionStatement.printSource(w, v)
			w.Write(semiColon)
		} else if s.IfStatement != nil {
			s.IfStatement.printSource(w, v)
		} else if s.IterationStatementDo != nil {
			s.IterationStatementDo.printSource(w, v)
		} else if s.IterationStatementWhile != nil {
			s.IterationStatementWhile.printSource(w, v)
		} else if s.IterationStatementFor != nil {
			s.IterationStatementFor.printSource(w, v)
		} else if s.SwitchStatement != nil {
			s.SwitchStatement.printSource(w, v)
		} else if s.WithStatement != nil {
			s.WithStatement.printSource(w, v)
		} else if s.LabelIdentifier != nil {
			io.WriteString(w, s.LabelIdentifier.Data)
			w.Write(labelPost)
			if s.LabelledItemFunction != nil {
				s.LabelledItemFunction.printSource(w, v)
			} else if s.LabelledItemStatement != nil {
				s.LabelledItemStatement.printSource(w, v)
			}
		} else if s.TryStatement != nil {
			s.TryStatement.printSource(w, v)
		}
	case StatementContinue:
		if s.LabelIdentifier == nil {
			io.WriteString(w, "continue;")
		} else {
			io.WriteString(w, "continue ")
			io.WriteString(w, s.LabelIdentifier.Data)
			w.Write(semiColon)
		}
	case StatementBreak:
		if s.LabelIdentifier == nil {
			io.WriteString(w, "break;")
		} else {
			io.WriteString(w, "break ")
			io.WriteString(w, s.LabelIdentifier.Data)
			w.Write(semiColon)
		}
	case StatementReturn:
		if s.ExpressionStatement == nil {
			io.WriteString(w, "return;")
		} else {
			io.WriteString(w, "return ")
			s.ExpressionStatement.printSource(w, v)
			w.Write(semiColon)
		}
	case StatementThrow:
		if s.ExpressionStatement != nil {
			io.WriteString(w, "throw ")
			s.ExpressionStatement.printSource(w, v)
			w.Write(semiColon)
		}
	case StatementDebugger:
		io.WriteString(w, "debugger;")
	}
}

func (d Declaration) printSource(w io.Writer, v bool) {
	if d.ClassDeclaration != nil {
		d.ClassDeclaration.printSource(w, v)
	} else if d.FunctionDeclaration != nil {
		d.FunctionDeclaration.printSource(w, v)
	} else if d.LexicalDeclaration != nil {
		d.LexicalDeclaration.printSource(w, v)
	}
}

func (b Block) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	var lastLine uint64

	if v && len(b.Tokens) > 0 {
		lastLine = b.Tokens[0].Line
	}

	pp := indentPrinter{w}

	for _, stmt := range b.StatementList {
		if v {
			if len(stmt.Tokens) > 0 {
				ll := stmt.Tokens[0].Line

				if ll > lastLine {
					pp.Write(newLine)
				} else {
					pp.Write(space)
				}

				lastLine = ll
			} else {
				pp.Write(newLine)
			}
		} else {
			pp.Write(newLine)
		}

		stmt.printSource(&pp, v)
	}

	if len(b.StatementList) > 0 {
		if v && len(b.Tokens) > 0 {
			if b.Tokens[len(b.Tokens)-1].Line > lastLine {
				w.Write(newLine)
			} else {
				pp.Write(space)
			}
		} else {
			w.Write(newLine)
		}
	}

	w.Write(blockClose)
}

func (vs VariableStatement) printSource(w io.Writer, v bool) {
	if len(vs.VariableDeclarationList) == 0 {
		return
	}

	io.WriteString(w, "var ")

	var lastLine uint64

	if v && len(vs.Tokens) > 0 {
		lastLine = vs.Tokens[0].Line
	}

	for n, vd := range vs.VariableDeclarationList {
		if n > 0 {
			if v && len(vd.Tokens) > 0 {
				if ll := vd.Tokens[0].Line; ll > lastLine {
					lastLine = ll

					w.Write(commaSepNL)
				} else {
					w.Write(commaSep)
				}
			} else {
				w.Write(commaSep)
			}
		}

		vd.printSource(w, v)
	}

	w.Write(semiColon)
}

func (e Expression) printSource(w io.Writer, v bool) {
	if len(e.Expressions) == 0 {
		return
	}

	var lastLine uint64

	if v && len(e.Tokens) > 0 {
		lastLine = e.Tokens[0].Line
	}

	e.Expressions[0].printSource(w, v)

	for _, ae := range e.Expressions[1:] {
		if v && len(ae.Tokens) > 0 {
			if ll := ae.Tokens[0].Line; ll > lastLine {
				lastLine = ll
				w.Write(commaSepNL)
			} else {
				w.Write(commaSep)
			}
		} else {
			w.Write(commaSep)
		}

		ae.printSource(w, v)
	}
}

func (i IfStatement) printSource(w io.Writer, v bool) {
	w.Write(ifOpen)

	if v {
		pp := indentPrinter{w}

		var nl bool

		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			nl = true

			pp.Write(newLine)
		}

		i.Expression.printSource(&pp, true)

		if nl {
			w.Write(newLine)
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.Write(parenCloseSpace)
	i.Statement.printSource(w, v)

	if i.ElseStatement != nil {
		w.Write(elseOpen)
		i.ElseStatement.printSource(w, v)
	}
}

func (i IterationStatementDo) printSource(w io.Writer, v bool) {
	w.Write(doOpen)
	i.Statement.printSource(w, v)
	w.Write(doWhileOpen)

	if v {
		pp := indentPrinter{w}

		var nl bool

		if len(i.Expression.Tokens) > 0 && len(i.Tokens) > 0 && i.Expression.Tokens[0].Line < i.Tokens[len(i.Tokens)-1].Line {
			nl = true

			pp.Write(newLine)
		}

		i.Expression.printSource(&pp, true)

		if nl {
			w.Write(newLine)
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.Write(doWhileClose)
}

func (i IterationStatementWhile) printSource(w io.Writer, v bool) {
	w.Write(whileOpen)

	if v {
		pp := indentPrinter{w}

		var nl bool

		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			pp.Write(newLine)

			nl = true
		}

		i.Expression.printSource(&pp, true)

		if nl {
			w.Write(newLine)
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.Write(parenCloseSpace)
	i.Statement.printSource(w, v)
}

func (i IterationStatementFor) printSource(w io.Writer, v bool) {
	switch i.Type {
	case ForNormal:
		if i.InitVar != nil || i.InitLexical != nil || i.InitExpression != nil {
			return
		}
	case ForNormalVar:
		if len(i.InitVar) == 0 {
			return
		}
	case ForNormalLexicalDeclaration:
		if i.InitLexical == nil {
			return
		}
	case ForNormalExpression:
		if i.InitExpression == nil {
			return
		}
	case ForInLeftHandSide, ForOfLeftHandSide, ForAwaitOfLeftHandSide:
		if i.LeftHandSideExpression == nil {
			return
		}
	case ForInVar, ForOfVar, ForAwaitOfVar, ForInLet, ForOfLet, ForAwaitOfLet, ForInConst, ForOfConst, ForAwaitOfConst:
		if i.ForBindingIdentifier == nil && i.ForBindingPatternObject == nil && i.ForBindingPatternArray == nil {
			return
		}
	default:
		return
	}

	switch i.Type {
	case ForInLeftHandSide, ForInVar, ForInLet, ForInConst:
		if i.In == nil {
			return
		}
	case ForOfLeftHandSide, ForOfVar, ForOfLet, ForOfConst, ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		if i.Of == nil {
			return
		}
	}

	switch i.Type {
	case ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		w.Write(forAwaitOpen)
	default:
		w.Write(forOpen)
	}

	pp := indentPrinter{w}

	var lastLine uint64

	if v && len(i.Tokens) > 0 {
		lastLine = i.Tokens[0].Line
	}

	var endline bool

	switch i.Type {
	case ForNormal:
		w.Write(semiColon)
	case ForNormalVar:
		if v && len(i.InitVar[0].Tokens) > 0 {
			if i.InitVar[0].Tokens[0].Line > lastLine {
				pp.Write(newLine)
			}

			lastLine = i.InitVar[0].Tokens[len(i.InitVar[0].Tokens)-1].Line
		}

		w.Write(varOpen)
		LexicalBinding(i.InitVar[0]).printSource(&pp, v)

		for _, vd := range i.InitVar[1:] {
			if v && len(vd.Tokens) > 0 {
				if vd.Tokens[0].Line > lastLine {
					pp.Write(commaSepNL)
				} else {
					pp.Write(commaSep)
				}
			} else {
				pp.Write(commaSep)
			}

			LexicalBinding(vd).printSource(&pp, v)
		}

		w.Write(semiColon)
	case ForNormalLexicalDeclaration:
		if v && len(i.InitLexical.Tokens) > 0 {
			if i.InitLexical.Tokens[0].Line > lastLine {
				endline = true

				pp.Write(newLine)
			}

			lastLine = i.InitLexical.Tokens[len(i.InitLexical.Tokens)-1].Line
		}

		i.InitLexical.printSource(&pp, v)
	case ForNormalExpression:
		if v && len(i.InitExpression.Tokens) > 0 {
			if i.InitExpression.Tokens[0].Line > lastLine {
				endline = true

				pp.Write(newLine)
			}

			lastLine = i.InitExpression.Tokens[len(i.InitExpression.Tokens)-1].Line
		}

		i.InitExpression.printSource(&pp, v)
		w.Write(semiColon)
	case ForInLeftHandSide, ForOfLeftHandSide, ForAwaitOfLeftHandSide:
		if v {
			if len(i.LeftHandSideExpression.Tokens) > 0 {
				if i.LeftHandSideExpression.Tokens[0].Line > lastLine {
					endline = true

					pp.Write(newLine)
				}

				lastLine = i.LeftHandSideExpression.Tokens[len(i.LeftHandSideExpression.Tokens)-1].Line
			}
		}

		i.LeftHandSideExpression.printSource(&pp, v)
	default:
		switch i.Type {
		case ForInVar, ForOfVar, ForAwaitOfVar:
			w.Write(varOpen)
		case ForInLet, ForOfLet, ForAwaitOfLet:
			w.Write(letOpen)
		case ForInConst, ForOfConst, ForAwaitOfConst:
			w.Write(constOpen)
		}

		if i.ForBindingIdentifier != nil {
			io.WriteString(w, i.ForBindingIdentifier.Data)
		} else if i.ForBindingPatternObject != nil {
			i.ForBindingPatternObject.printSource(w, v)
		} else {
			i.ForBindingPatternArray.printSource(w, v)
		}
	}

	switch i.Type {
	case ForNormal, ForNormalVar, ForNormalLexicalDeclaration, ForNormalExpression:
		if i.Conditional != nil {
			if v && len(i.Conditional.Tokens) > 0 {
				if i.Conditional.Tokens[0].Line > lastLine {
					endline = true

					pp.Write(newLine)
				} else {
					w.Write(space)
				}

				lastLine = i.Conditional.Tokens[len(i.Conditional.Tokens)-1].Line
			} else {
				w.Write(space)
			}

			i.Conditional.printSource(&pp, v)
		}

		w.Write(semiColon)

		if i.Afterthought != nil {
			if v && len(i.Afterthought.Tokens) > 0 {
				if i.Afterthought.Tokens[0].Line > lastLine {
					endline = true

					pp.Write(newLine)
				} else {
					w.Write(space)
				}
			} else {
				w.Write(space)
			}

			i.Afterthought.printSource(&pp, v)
		}
	case ForInLeftHandSide, ForInVar, ForInLet, ForInConst:
		w.Write(forIn)
		i.In.printSource(&pp, v)
	case ForOfLeftHandSide, ForOfVar, ForOfLet, ForOfConst, ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		w.Write(forOf)
		i.Of.printSource(&pp, v)
	}

	if endline {
		w.Write(newLine)
	}

	w.Write(parenCloseSpace)
	i.Statement.printSource(w, v)
}

func (s SwitchStatement) printSource(w io.Writer, v bool) {
	w.Write(switchOpen)

	if v {
		pp := indentPrinter{w}

		var nl bool

		if len(s.Tokens) > 0 && len(s.Expression.Tokens) > 0 && s.Expression.Tokens[0].Line > s.Tokens[0].Line {
			nl = true

			pp.Write(newLine)
		}

		s.Expression.printSource(&pp, true)

		if nl {
			w.Write(newLine)
		}
	} else {
		s.Expression.printSource(w, false)
	}

	w.Write(switchClose)

	if len(s.CaseClauses) > 0 || s.DefaultClause != nil || len(s.PostDefaultCaseClauses) > 0 {
		w.Write(newLine)
	}

	for _, c := range s.CaseClauses {
		c.printSource(w, v)
		w.Write(newLine)
	}

	if s.DefaultClause != nil {
		w.Write(defaultCase)

		pp := indentPrinter{w}

		for _, stmt := range s.DefaultClause {
			pp.Write(newLine)
			stmt.printSource(&pp, v)
		}

		w.Write(newLine)
	}

	for _, c := range s.PostDefaultCaseClauses {
		c.printSource(w, v)
		w.Write(newLine)
	}

	w.Write(blockClose)
}

func (ws WithStatement) printSource(w io.Writer, v bool) {
	w.Write(withOpen)

	if v {
		pp := indentPrinter{w}

		var nl bool

		if len(ws.Tokens) > 0 && len(ws.Expression.Tokens) > 0 && ws.Expression.Tokens[0].Line > ws.Tokens[0].Line {
			nl = true

			pp.Write(newLine)
		}

		ws.Expression.printSource(&pp, true)

		if nl {
			w.Write(newLine)
		}
	} else {
		ws.Expression.printSource(w, false)
	}

	w.Write(parenCloseSpace)
	ws.Statement.printSource(w, v)
}

func (f FunctionDeclaration) printSource(w io.Writer, v bool) {
	switch f.Type {
	case FunctionNormal:
		w.Write(funcOpen)
	case FunctionGenerator:
		w.Write(genFuncOpen)
	case FunctionAsync:
		w.Write(asyncFuncOpen)
	case FunctionAsyncGenerator:
		w.Write(asyncGenFuncOpen)
	default:
		return
	}

	if f.BindingIdentifier != nil {
		io.WriteString(w, f.BindingIdentifier.Data)
	}

	f.FormalParameters.printSource(&indentPrinter{w}, v)
	f.FunctionBody.printSource(w, v)
}

func (t TryStatement) printSource(w io.Writer, v bool) {
	w.Write(tryOpen)
	t.TryBlock.printSource(w, v)

	if t.CatchBlock != nil {
		if t.CatchParameterBindingIdentifier != nil {
			w.Write(catchParenOpen)
			io.WriteString(w, t.CatchParameterBindingIdentifier.Data)
			w.Write(parenCloseSpace)
		} else if t.CatchParameterArrayBindingPattern != nil {
			w.Write(catchParenOpen)
			t.CatchParameterArrayBindingPattern.printSource(w, v)
			w.Write(parenCloseSpace)
		} else if t.CatchParameterObjectBindingPattern != nil {
			w.Write(catchParenOpen)
			t.CatchParameterObjectBindingPattern.printSource(w, v)
			w.Write(parenCloseSpace)
		} else {
			w.Write(catchOpen)
		}

		t.CatchBlock.printSource(w, v)
	}

	if t.FinallyBlock != nil {
		w.Write(finallyOpen)
		t.FinallyBlock.printSource(w, v)
	}
}

func (c ClassDeclaration) printSource(w io.Writer, v bool) {
	w.Write(classOpen)

	if c.BindingIdentifier != nil {
		io.WriteString(w, c.BindingIdentifier.Data)
		w.Write(space)
	}

	if c.ClassHeritage != nil {
		w.Write(extends)
		c.ClassHeritage.printSource(w, v)
		w.Write(space)
	}

	w.Write(blockOpen)

	if len(c.ClassBody) > 0 {
		pp := indentPrinter{w}

		for _, ce := range c.ClassBody {
			pp.Write(newLine)
			ce.printSource(&pp, v)
		}

		w.Write(newLine)
	}

	w.Write(blockClose)
}

func (l LexicalDeclaration) printSource(w io.Writer, v bool) {
	if len(l.BindingList) == 0 {
		return
	}

	if l.LetOrConst == Let {
		w.Write(letOpen)
	} else if l.LetOrConst == Const {
		w.Write(constOpen)
	}

	l.BindingList[0].printSource(w, v)

	for _, lb := range l.BindingList[1:] {
		if v {
			w.Write(commaSepNL)
		} else {
			w.Write(commaSep)
		}

		lb.printSource(w, v)
	}

	w.Write(semiColon)
}

func (l LexicalBinding) printSource(w io.Writer, v bool) {
	if l.BindingIdentifier != nil {
		io.WriteString(w, l.BindingIdentifier.Data)
	} else if l.ArrayBindingPattern != nil {
		l.ArrayBindingPattern.printSource(w, v)
	} else if l.ObjectBindingPattern != nil {
		l.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if l.Initializer != nil {
		w.Write(assignment)
		l.Initializer.printSource(w, v)
	}
}

func (a AssignmentExpression) printSource(w io.Writer, v bool) {
	if a.Yield && a.AssignmentExpression != nil {
		w.Write(yield)

		if a.Delegate {
			w.Write(delegate)
		}

		a.AssignmentExpression.printSource(w, v)
	} else if a.ArrowFunction != nil {
		a.ArrowFunction.printSource(w, v)
	} else if a.LeftHandSideExpression != nil && a.AssignmentExpression != nil {
		ao := assignment

		switch a.AssignmentOperator {
		case AssignmentAssign:
		case AssignmentMultiply:
			ao = assignmentMultiply
		case AssignmentDivide:
			ao = assignmentDivide
		case AssignmentRemainder:
			ao = assignmentRemainder
		case AssignmentAdd:
			ao = assignmentAdd
		case AssignmentSubtract:
			ao = assignmentSubtract
		case AssignmentLeftShift:
			ao = assignmentLeftShift
		case AssignmentSignPropagatingRightShift:
			ao = assignmentSignRightShift
		case AssignmentZeroFillRightShift:
			ao = assignmentZeroRightShift
		case AssignmentBitwiseAND:
			ao = assignmentAND
		case AssignmentBitwiseXOR:
			ao = assignmentXOR
		case AssignmentBitwiseOR:
			ao = assignmentOR
		case AssignmentExponentiation:
			ao = assignmentExponentiation
		case AssignmentLogicalAnd:
			ao = assignmentLogicalAnd
		case AssignmentLogicalOr:
			ao = assignmentLogicalOr
		case AssignmentNullish:
			ao = assignmentNullish
		default:
			return
		}

		a.LeftHandSideExpression.printSource(w, v)
		w.Write(ao)
		a.AssignmentExpression.printSource(w, v)
	} else if a.AssignmentPattern != nil && a.AssignmentExpression != nil && a.AssignmentOperator == AssignmentAssign {
		a.AssignmentPattern.printSource(w, v)
		w.Write(assignment)
		a.AssignmentExpression.printSource(w, v)
	} else if a.ConditionalExpression != nil {
		a.ConditionalExpression.printSource(w, v)
	}
}

func (l LeftHandSideExpression) printSource(w io.Writer, v bool) {
	if l.NewExpression != nil {
		l.NewExpression.printSource(w, v)
	} else if l.CallExpression != nil {
		l.CallExpression.printSource(w, v)
	} else if l.OptionalExpression != nil {
		l.OptionalExpression.printSource(w, v)
	}
}

func (a AssignmentPattern) printSource(w io.Writer, v bool) {
	if a.ArrayAssignmentPattern != nil {
		a.ArrayAssignmentPattern.printSource(w, v)
	} else if a.ObjectAssignmentPattern != nil {
		a.ObjectAssignmentPattern.printSource(w, v)
	}
}

func (a ArrayAssignmentPattern) printSource(w io.Writer, v bool) {
	w.Write(bracketOpen)

	for n, ae := range a.AssignmentElements {
		if n > 0 {
			w.Write(commaSep)
		}

		ae.printSource(w, v)
	}

	if a.AssignmentRestElement != nil {
		if len(a.AssignmentElements) > 0 {
			w.Write(commaSep)
		}

		w.Write(ellipsis)
		a.AssignmentRestElement.printSource(w, v)
	}

	w.Write(bracketClose)
}

func (o ObjectAssignmentPattern) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	for n, ap := range o.AssignmentPropertyList {
		if n > 0 {
			w.Write(commaSep)
		}

		ap.printSource(w, v)
	}

	if o.AssignmentRestElement != nil {
		if len(o.AssignmentPropertyList) > 0 {
			w.Write(commaSep)
		}

		w.Write(ellipsis)
		o.AssignmentRestElement.printSource(w, v)
	}

	w.Write(blockClose)
}

func (a AssignmentElement) printSource(w io.Writer, v bool) {
	a.DestructuringAssignmentTarget.printSource(w, v)

	if a.Initializer != nil {
		w.Write(assignment)
		a.Initializer.printSource(w, v)
	}
}

func (a AssignmentProperty) printSource(w io.Writer, v bool) {
	a.PropertyName.printSource(w, v)

	if a.DestructuringAssignmentTarget != nil {
		if !v && a.PropertyName.LiteralPropertyName != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.CallExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.OptionalExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News == 0 && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference.Data == a.PropertyName.LiteralPropertyName.Data {
			return
		}

		w.Write(colonSep)
		a.DestructuringAssignmentTarget.printSource(w, v)
	}

	if a.Initializer != nil {
		w.Write(assignment)
		a.Initializer.printSource(w, v)
	}
}

func (d DestructuringAssignmentTarget) printSource(w io.Writer, v bool) {
	if d.LeftHandSideExpression != nil {
		d.LeftHandSideExpression.printSource(w, v)
	} else if d.AssignmentPattern != nil {
		d.AssignmentPattern.printSource(w, v)
	}
}

func (o ObjectBindingPattern) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	for n, bp := range o.BindingPropertyList {
		if n > 0 {
			w.Write(commaSep)
		}

		bp.printSource(w, v)
	}

	if o.BindingRestProperty != nil {
		if len(o.BindingPropertyList) > 0 {
			w.Write(commaSep)
		}

		w.Write(ellipsis)
		io.WriteString(w, o.BindingRestProperty.Data)
	}

	w.Write(blockClose)
}

func (a ArrayBindingPattern) printSource(w io.Writer, v bool) {
	w.Write(bracketOpen)

	for n, be := range a.BindingElementList {
		if n > 0 {
			w.Write(commaSep)
		}

		be.printSource(w, v)
	}

	if a.BindingRestElement != nil {
		if len(a.BindingElementList) > 0 {
			w.Write(commaSep)
		}

		w.Write(ellipsis)
		a.BindingRestElement.printSource(w, v)
	}

	w.Write(bracketClose)
}

func (c CaseClause) printSource(w io.Writer, v bool) {
	w.Write(caseOpen)
	c.Expression.printSource(w, v)
	w.Write(caseClose)

	pp := indentPrinter{w}

	for _, stmt := range c.StatementList {
		pp.Write(newLine)
		stmt.printSource(&pp, v)
	}
}

func (f FormalParameters) printSource(w io.Writer, v bool) {
	w.Write(parenOpen)

	if len(f.FormalParameterList) > 0 {
		f.FormalParameterList[0].printSource(w, v)

		for _, be := range f.FormalParameterList[1:] {
			w.Write(commaSep)
			be.printSource(w, v)
		}

		if f.BindingIdentifier != nil || f.ArrayBindingPattern != nil || f.ObjectBindingPattern != nil {
			w.Write(commaSep)
		}
	}

	if f.BindingIdentifier != nil {
		w.Write(ellipsis)
		io.WriteString(w, f.BindingIdentifier.Data)
	} else if f.ArrayBindingPattern != nil {
		w.Write(ellipsis)
		f.ArrayBindingPattern.printSource(w, v)
	} else if f.ObjectBindingPattern != nil {
		w.Write(ellipsis)
		f.ObjectBindingPattern.printSource(w, v)
	}

	w.Write(parenCloseSpace)
}

func (m MethodDefinition) printSource(w io.Writer, v bool) {
	switch m.Type {
	case MethodNormal:
	case MethodGenerator:
		w.Write(methodGenerator)
	case MethodAsync:
		w.Write(methodAsync)
	case MethodAsyncGenerator:
		w.Write(methodAsyncGenerator)
	case MethodGetter:
		w.Write(methodGet)
	case MethodSetter:
		w.Write(methodSet)
	default:
		return
	}

	m.ClassElementName.printSource(w, v)
	m.Params.printSource(w, v)
	m.FunctionBody.printSource(w, v)
}

func (ce ClassElement) printSource(w io.Writer, v bool) {
	if ce.Static {
		w.Write(methodStatic)
	}

	if ce.MethodDefinition != nil {
		ce.MethodDefinition.printSource(w, v)
	} else if ce.FieldDefinition != nil {
		ce.FieldDefinition.printSource(w, v)
	} else if ce.ClassStaticBlock != nil {
		ce.ClassStaticBlock.printSource(w, v)
	}
}

func (fd FieldDefinition) printSource(w io.Writer, v bool) {
	fd.ClassElementName.printSource(w, v)

	if fd.Initializer != nil {
		w.Write(assignment)
		fd.Initializer.printSource(w, v)
	}

	w.Write(semiColon)
}

func (cen ClassElementName) printSource(w io.Writer, v bool) {
	if cen.PropertyName != nil {
		cen.PropertyName.printSource(w, v)
	} else if cen.PrivateIdentifier != nil {
		io.WriteString(w, cen.PrivateIdentifier.Data)
	}
}

func (c ConditionalExpression) printSource(w io.Writer, v bool) {
	if c.LogicalORExpression != nil {
		c.LogicalORExpression.printSource(w, v)
	} else if c.CoalesceExpression != nil {
		c.CoalesceExpression.printSource(w, v)
	}

	if c.True != nil && c.False != nil {
		w.Write(conditionalStart)
		c.True.printSource(w, v)
		w.Write(conditionalSep)
		c.False.printSource(w, v)
	}
}

func (a ArrowFunction) printSource(w io.Writer, v bool) {
	if a.FunctionBody == nil && a.AssignmentExpression == nil || (a.BindingIdentifier == nil && a.FormalParameters == nil) {
		return
	}

	if a.Async {
		w.Write(methodAsync)
	}

	if a.BindingIdentifier != nil {
		io.WriteString(w, a.BindingIdentifier.Data)
		w.Write(space)
	} else if a.FormalParameters != nil {
		a.FormalParameters.printSource(w, v)
	}

	w.Write(arrow)

	if a.FunctionBody != nil {
		a.FunctionBody.printSource(w, v)
	} else {
		a.AssignmentExpression.printSource(w, v)
	}
}

func (n NewExpression) printSource(w io.Writer, v bool) {
	for i := uint(0); i < n.News; i++ {
		w.Write(news)
	}

	n.MemberExpression.printSource(w, v)
}

func (c CallExpression) printSource(w io.Writer, v bool) {
	if c.SuperCall && c.Arguments != nil {
		w.Write(super)
		c.Arguments.printSource(w, v)
	} else if c.ImportCall != nil {
		w.Write(importCall)
		c.ImportCall.printSource(w, v)
		w.Write(parenClose)
	} else if c.MemberExpression != nil && c.Arguments != nil {
		c.MemberExpression.printSource(w, v)
		c.Arguments.printSource(w, v)
	} else if c.CallExpression != nil {
		if c.Arguments != nil {
			c.CallExpression.printSource(w, v)
			c.Arguments.printSource(w, v)
		} else if c.Expression != nil {
			c.CallExpression.printSource(w, v)
			w.Write(bracketOpen)
			c.Expression.printSource(w, v)
			w.Write(bracketClose)
		} else if c.IdentifierName != nil {
			c.CallExpression.printSource(w, v)

			if v && len(c.CallExpression.Tokens) > 0 && c.IdentifierName.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.Write(newLine)
			}

			w.Write(dot)
			io.WriteString(w, c.IdentifierName.Data)
		} else if c.TemplateLiteral != nil {
			c.CallExpression.printSource(w, v)
			c.TemplateLiteral.printSource(w, v)
		} else if c.PrivateIdentifier != nil {
			c.CallExpression.printSource(w, v)

			if v && len(c.CallExpression.Tokens) > 0 && c.PrivateIdentifier.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.Write(newLine)
			}

			w.Write(dot)
			io.WriteString(w, c.PrivateIdentifier.Data)
		}
	}
}

func (b BindingProperty) printSource(w io.Writer, v bool) {
	if !v && b.PropertyName.LiteralPropertyName != nil && b.BindingElement.SingleNameBinding != nil && b.PropertyName.LiteralPropertyName.Data == b.BindingElement.SingleNameBinding.Data {
		b.BindingElement.printSource(w, v)
	} else {
		b.PropertyName.printSource(w, v)
		w.Write(colonSep)
		b.BindingElement.printSource(w, v)
	}
}

func (b BindingElement) printSource(w io.Writer, v bool) {
	if b.SingleNameBinding != nil {
		io.WriteString(w, b.SingleNameBinding.Data)
	} else if b.ArrayBindingPattern != nil {
		b.ArrayBindingPattern.printSource(w, v)
	} else if b.ObjectBindingPattern != nil {
		b.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if b.Initializer != nil {
		w.Write(assignment)
		b.Initializer.printSource(w, v)
	}
}

func (p PropertyName) printSource(w io.Writer, v bool) {
	if p.LiteralPropertyName != nil {
		io.WriteString(w, p.LiteralPropertyName.Data)
	} else if p.ComputedPropertyName != nil {
		w.Write(bracketOpen)
		p.ComputedPropertyName.printSource(w, v)
		w.Write(bracketClose)
	}
}

func (l LogicalORExpression) printSource(w io.Writer, v bool) {
	if l.LogicalORExpression != nil {
		l.LogicalORExpression.printSource(w, v)
		w.Write(logicalOR)
	}

	l.LogicalANDExpression.printSource(w, v)
}

func (c ParenthesizedExpression) printSource(w io.Writer, v bool) {
	w.Write(parenOpen)

	if len(c.Expressions) > 0 {
		c.Expressions[0].printSource(w, v)

		for _, e := range c.Expressions[1:] {
			w.Write(commaSep)
			e.printSource(w, v)
		}
	}

	w.Write(parenClose)
}

func (m MemberExpression) printSource(w io.Writer, v bool) {
	if m.MemberExpression != nil {
		if m.Arguments != nil {
			w.Write(news)
			m.MemberExpression.printSource(w, v)
			m.Arguments.printSource(w, v)
		} else if m.Expression != nil {
			m.MemberExpression.printSource(w, v)
			w.Write(bracketOpen)
			m.Expression.printSource(w, v)
			w.Write(bracketClose)
		} else if m.IdentifierName != nil {
			m.MemberExpression.printSource(w, v)

			if v && len(m.MemberExpression.Tokens) > 0 && m.IdentifierName.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.Write(newLine)
			}

			w.Write(dot)
			io.WriteString(w, m.IdentifierName.Data)
		} else if m.PrivateIdentifier != nil {
			m.MemberExpression.printSource(w, v)

			if v && len(m.MemberExpression.Tokens) > 0 && m.PrivateIdentifier.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.Write(newLine)
			}

			w.Write(dot)
			io.WriteString(w, m.PrivateIdentifier.Data)
		} else if m.TemplateLiteral != nil {
			m.MemberExpression.printSource(w, v)
			m.TemplateLiteral.printSource(w, v)
		}
	} else if m.PrimaryExpression != nil {
		m.PrimaryExpression.printSource(w, v)
	} else if m.SuperProperty {
		if m.Expression != nil {
			w.Write(super)
			w.Write(bracketOpen)
			m.Expression.printSource(w, v)
			w.Write(bracketClose)
		} else if m.IdentifierName != nil {
			w.Write(super)
			w.Write(dot)
			io.WriteString(w, m.IdentifierName.Data)
		}
	} else if m.NewTarget {
		w.Write(newTarget)
	} else if m.ImportMeta {
		w.Write(importMeta)
	}
}

func (a Argument) printSource(w io.Writer, v bool) {
	if a.Spread {
		w.Write(ellipsis)
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a Arguments) printSource(w io.Writer, v bool) {
	w.Write(parenOpen)

	if len(a.ArgumentList) > 0 {
		a.ArgumentList[0].printSource(w, v)

		for _, ae := range a.ArgumentList[1:] {
			w.Write(commaSep)
			ae.printSource(w, v)
		}
	}

	w.Write(parenClose)
}

func (t TemplateLiteral) printSource(w io.Writer, v bool) {
	x := w

	for {
		j, ok := x.(*indentPrinter)
		if !ok {
			break
		}

		x = j.Writer
	}

	if t.NoSubstitutionTemplate != nil {
		io.WriteString(x, t.NoSubstitutionTemplate.Data)
	} else if t.TemplateHead != nil && t.TemplateTail != nil && len(t.Expressions) == len(t.TemplateMiddleList)+1 {
		io.WriteString(x, t.TemplateHead.Data)
		t.Expressions[0].printSource(w, v)

		for n, e := range t.Expressions[1:] {
			io.WriteString(x, t.TemplateMiddleList[n].Data)
			e.printSource(w, v)
		}

		io.WriteString(x, t.TemplateTail.Data)
	}
}

func (l LogicalANDExpression) printSource(w io.Writer, v bool) {
	if l.LogicalANDExpression != nil {
		l.LogicalANDExpression.printSource(w, v)
		w.Write(logicalAND)
	}

	l.BitwiseORExpression.printSource(w, v)
}

func (p PrimaryExpression) printSource(w io.Writer, v bool) {
	if p.This != nil {
		w.Write(this)
	} else if p.IdentifierReference != nil {
		io.WriteString(w, p.IdentifierReference.Data)
	} else if p.Literal != nil {
		io.WriteString(w, p.Literal.Data)
	} else if p.ArrayLiteral != nil {
		p.ArrayLiteral.printSource(w, v)
	} else if p.ObjectLiteral != nil {
		p.ObjectLiteral.printSource(w, v)
	} else if p.FunctionExpression != nil {
		p.FunctionExpression.printSource(w, v)
	} else if p.ClassExpression != nil {
		p.ClassExpression.printSource(w, v)
	} else if p.TemplateLiteral != nil {
		p.TemplateLiteral.printSource(w, v)
	} else if p.ParenthesizedExpression != nil {
		p.ParenthesizedExpression.printSource(w, v)
	}
}

func (b BitwiseORExpression) printSource(w io.Writer, v bool) {
	if b.BitwiseORExpression != nil {
		b.BitwiseORExpression.printSource(w, v)
		w.Write(bitwiseOR)
	}

	b.BitwiseXORExpression.printSource(w, v)
}

func (a ArrayElement) printSource(w io.Writer, v bool) {
	if a.Spread {
		w.Write(ellipsis)
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a ArrayLiteral) printSource(w io.Writer, v bool) {
	w.Write(bracketOpen)

	if len(a.ElementList) > 0 {
		a.ElementList[0].printSource(w, v)

		for _, ae := range a.ElementList[1:] {
			w.Write(commaSep)
			ae.printSource(w, v)
		}
	}

	w.Write(bracketClose)
}

func (o ObjectLiteral) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	if len(o.PropertyDefinitionList) > 0 {
		var lastLine uint64

		x := w

		if v && len(o.Tokens) > 0 {
			lastLine = o.Tokens[0].Line
			x = &indentPrinter{w}
		}

		for n, pd := range o.PropertyDefinitionList {
			if n > 0 {
				if v && len(pd.Tokens) > 0 {
					if ll := pd.Tokens[0].Line; ll > lastLine {
						lastLine = ll

						x.Write(commaSepNL)
					} else {
						x.Write(commaSep)
					}
				} else {
					x.Write(commaSep)
				}
			} else if v && len(pd.Tokens) > 0 {
				if ll := pd.Tokens[0].Line; ll > lastLine {
					lastLine = ll

					x.Write(newLine)
				}
			}

			pd.printSource(x, v)
		}
		if v && len(o.Tokens) > 0 {
			if ll := o.Tokens[len(o.Tokens)-1].Line; ll > lastLine {
				w.Write(newLine)
			}
		}
	}

	w.Write(blockClose)
}

func (b BitwiseXORExpression) printSource(w io.Writer, v bool) {
	if b.BitwiseXORExpression != nil {
		b.BitwiseXORExpression.printSource(w, v)
		w.Write(bitwiseXOR)
	}

	b.BitwiseANDExpression.printSource(w, v)
}

func (p PropertyDefinition) printSource(w io.Writer, v bool) {
	if p.AssignmentExpression != nil {
		if p.PropertyName != nil {
			p.PropertyName.printSource(w, v)

			done := false

			if !v && !p.IsCoverInitializedName && p.PropertyName.LiteralPropertyName != nil && p.AssignmentExpression.ConditionalExpression != nil {
				c := UnwrapConditional(p.AssignmentExpression.ConditionalExpression)

				if pe, ok := c.(*PrimaryExpression); ok && pe.IdentifierReference != nil {
					done = pe.IdentifierReference.Type == p.PropertyName.LiteralPropertyName.Type && pe.IdentifierReference.Data == p.PropertyName.LiteralPropertyName.Data
				}
			}

			if !done {
				if p.IsCoverInitializedName {
					w.Write(assignment)
				} else {
					w.Write(colonSep)
				}

				p.AssignmentExpression.printSource(w, v)
			}
		} else {
			w.Write(ellipsis)
			p.AssignmentExpression.printSource(w, v)
		}
	} else if p.MethodDefinition != nil {
		p.MethodDefinition.printSource(w, v)
	}
}

func (b BitwiseANDExpression) printSource(w io.Writer, v bool) {
	if b.BitwiseANDExpression != nil {
		b.BitwiseANDExpression.printSource(w, v)
		w.Write(bitwiseAND)
	}

	b.EqualityExpression.printSource(w, v)
}

func (e EqualityExpression) printSource(w io.Writer, v bool) {
	if e.EqualityExpression != nil {
		var eo []byte

		switch e.EqualityOperator {
		case EqualityEqual:
			eo = equalityEqual
		case EqualityNotEqual:
			eo = equalityNotEqual
		case EqualityStrictEqual:
			eo = equalityStrictEqual
		case EqualityStrictNotEqual:
			eo = equalityStrictNotEqual
		default:
			return
		}

		e.EqualityExpression.printSource(w, v)
		w.Write(eo)
	}

	e.RelationalExpression.printSource(w, v)
}

func (r RelationalExpression) printSource(w io.Writer, v bool) {
	if r.PrivateIdentifier != nil {
		io.WriteString(w, r.PrivateIdentifier.Data)
		w.Write(relationshipIn)
	} else if r.RelationalExpression != nil {
		var ro []byte

		switch r.RelationshipOperator {
		case RelationshipLessThan:
			ro = relationshipLessThan
		case RelationshipGreaterThan:
			ro = relationshipGreaterThan
		case RelationshipLessThanEqual:
			ro = relationshipLessThanEqual
		case RelationshipGreaterThanEqual:
			ro = relationshipGreaterThanEqual
		case RelationshipInstanceOf:
			ro = relationshipInstanceOf
		case RelationshipIn:
			ro = relationshipIn
		default:
			return
		}

		r.RelationalExpression.printSource(w, v)
		w.Write(ro)
	}

	r.ShiftExpression.printSource(w, v)
}

func (s ShiftExpression) printSource(w io.Writer, v bool) {
	if s.ShiftExpression != nil {
		var so []byte

		switch s.ShiftOperator {
		case ShiftLeft:
			so = shiftLeft
		case ShiftRight:
			so = shiftRight
		case ShiftUnsignedRight:
			so = shiftUnsignedRight
		default:
			return
		}

		s.ShiftExpression.printSource(w, v)
		w.Write(so)
	}

	s.AdditiveExpression.printSource(w, v)
}

func (a AdditiveExpression) printSource(w io.Writer, v bool) {
	if a.AdditiveExpression != nil {
		var ao []byte

		switch a.AdditiveOperator {
		case AdditiveAdd:
			ao = additiveAdd
		case AdditiveMinus:
			ao = additiveMinus
		default:
			return
		}

		a.AdditiveExpression.printSource(w, v)
		w.Write(ao)
	}

	a.MultiplicativeExpression.printSource(w, v)
}

func (m MultiplicativeExpression) printSource(w io.Writer, v bool) {
	if m.MultiplicativeExpression != nil {
		var mo []byte

		switch m.MultiplicativeOperator {
		case MultiplicativeMultiply:
			mo = multiplicativeMultiply
		case MultiplicativeDivide:
			mo = multiplicativeDivide
		case MultiplicativeRemainder:
			mo = multiplicativeRemainder
		default:
			return
		}

		m.MultiplicativeExpression.printSource(w, v)
		w.Write(mo)
	}

	m.ExponentiationExpression.printSource(w, v)
}

func (e ExponentiationExpression) printSource(w io.Writer, v bool) {
	if e.ExponentiationExpression != nil {
		e.ExponentiationExpression.printSource(w, v)
		w.Write(exponentionation)
	}

	e.UnaryExpression.printSource(w, v)
}

func (u UnaryExpression) printSource(w io.Writer, v bool) {
	for _, uo := range u.UnaryOperators {
		switch uo {
		case UnaryDelete:
			w.Write(unaryDelete)
		case UnaryVoid:
			w.Write(unaryVoid)
		case UnaryTypeOf:
			w.Write(unaryTypeOf)
		case UnaryAdd:
			w.Write(unaryAdd)
		case UnaryMinus:
			w.Write(unaryMinus)
		case UnaryBitwiseNot:
			w.Write(unaryBitwiseNot)
		case UnaryLogicalNot:
			w.Write(unaryLogicalNot)
		case UnaryAwait:
			w.Write(unaryAwait)
		}
	}

	u.UpdateExpression.printSource(w, v)
}

func (u UpdateExpression) printSource(w io.Writer, v bool) {
	if u.LeftHandSideExpression != nil {
		var uo []byte

		switch u.UpdateOperator {
		case UpdatePostIncrement:
			uo = updateIncrement
		case UpdatePostDecrement:
			uo = updateDecrement
		case UpdatePreIncrement, UpdatePreDecrement:
			return
		default:
		}

		u.LeftHandSideExpression.printSource(w, v)

		if len(uo) > 0 {
			w.Write(uo)
		}
	} else if u.UnaryExpression != nil {
		switch u.UpdateOperator {
		case UpdatePreIncrement:
			w.Write(updateIncrement)
		case UpdatePreDecrement:
			w.Write(updateDecrement)
		default:
			return
		}

		u.UnaryExpression.printSource(w, v)
	}
}

func (m Module) printSource(w io.Writer, v bool) {
	if len(m.ModuleListItems) > 0 {
		m.ModuleListItems[0].printSource(w, v)

		for _, mi := range m.ModuleListItems[1:] {
			w.Write(doubleNewLine)
			mi.printSource(w, v)
		}
	}
}

func (m ModuleItem) printSource(w io.Writer, v bool) {
	if m.ImportDeclaration != nil {
		m.ImportDeclaration.printSource(w, v)
	} else if m.ExportDeclaration != nil {
		m.ExportDeclaration.printSource(w, v)
	} else if m.StatementListItem != nil {
		m.StatementListItem.printSource(w, v)
	}
}

func (i ImportDeclaration) printSource(w io.Writer, v bool) {
	if i.ImportClause == nil && i.FromClause.ModuleSpecifier == nil {
		return
	}

	w.Write(importc)

	if i.ImportClause != nil {
		i.ImportClause.printSource(w, v)
		i.FromClause.printSource(w, v)
	} else if i.FromClause.ModuleSpecifier != nil {
		io.WriteString(w, i.FromClause.ModuleSpecifier.Data)
	}

	if i.WithClause != nil {
		w.Write(space)
		i.WithClause.printSource(w, v)
	}

	w.Write(semiColon)
}

func (e ExportDeclaration) printSource(w io.Writer, v bool) {
	if e.FromClause != nil {
		w.Write(exportc)

		if e.ExportClause != nil {
			e.ExportClause.printSource(w, v)
		} else {
			w.Write(exportAll)

			if e.ExportFromClause != nil {
				w.Write(as)
				io.WriteString(w, e.ExportFromClause.Data)
			}
		}

		e.FromClause.printSource(w, v)
		w.Write(semiColon)
	} else if e.ExportClause != nil {
		w.Write(exportc)
		e.ExportClause.printSource(w, v)
		w.Write(semiColon)
	} else if e.VariableStatement != nil {
		w.Write(exportc)
		e.VariableStatement.printSource(w, v)
	} else if e.Declaration != nil {
		w.Write(exportc)
		e.Declaration.printSource(w, v)
	} else if e.DefaultFunction != nil {
		w.Write(exportd)
		e.DefaultFunction.printSource(w, v)
	} else if e.DefaultClass != nil {
		w.Write(exportd)
		e.DefaultClass.printSource(w, v)
	} else if e.DefaultAssignmentExpression != nil {
		w.Write(exportd)
		e.DefaultAssignmentExpression.printSource(w, v)
		w.Write(semiColon)
	}
}

func (wc WithClause) printSource(w io.Writer, v bool) {
	w.Write(withOpen[:5])
	w.Write(blockOpen)

	if len(wc.WithEntries) > 0 {
		wc.WithEntries[0].printSource(w, v)

		for _, we := range wc.WithEntries[1:] {
			w.Write(commaSep)
			we.printSource(w, v)
		}
	}

	w.Write(blockClose)
}

func (we WithEntry) printSource(w io.Writer, v bool) {
}

func (i ImportClause) printSource(w io.Writer, v bool) {
	if i.ImportedDefaultBinding != nil {
		io.WriteString(w, i.ImportedDefaultBinding.Data)

		if i.NameSpaceImport != nil || i.NamedImports != nil {
			w.Write(commaSep)
		}
	}

	if i.NameSpaceImport != nil {
		w.Write(namespaceImport)
		io.WriteString(w, i.NameSpaceImport.Data)
	} else if i.NamedImports != nil {
		i.NamedImports.printSource(w, v)
	}
}

func (f FromClause) printSource(w io.Writer, v bool) {
	if f.ModuleSpecifier == nil {
		return
	}

	w.Write(from)
	io.WriteString(w, f.ModuleSpecifier.Data)
}

func (e ExportClause) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	if len(e.ExportList) > 0 {
		e.ExportList[0].printSource(w, v)

		for _, es := range e.ExportList[1:] {
			w.Write(commaSep)
			es.printSource(w, v)
		}
	}

	w.Write(blockClose)
}

func (n NamedImports) printSource(w io.Writer, v bool) {
	w.Write(blockOpen)

	if len(n.ImportList) > 0 {
		n.ImportList[0].printSource(w, v)

		for _, is := range n.ImportList[1:] {
			w.Write(commaSep)
			is.printSource(w, v)
		}
	}

	w.Write(blockClose)
}

func (e ExportSpecifier) printSource(w io.Writer, v bool) {
	if e.IdentifierName == nil {
		return
	}

	io.WriteString(w, e.IdentifierName.Data)

	if e.EIdentifierName != nil && (e.EIdentifierName.Type != e.IdentifierName.Type || e.EIdentifierName.Data != e.IdentifierName.Data || v) {
		w.Write(as)
		io.WriteString(w, e.EIdentifierName.Data)
	}
}

func (i ImportSpecifier) printSource(w io.Writer, v bool) {
	if i.ImportedBinding == nil {
		return
	}

	if i.IdentifierName != nil && (i.IdentifierName.Type != i.ImportedBinding.Type || i.IdentifierName.Data != i.ImportedBinding.Data || v) {
		io.WriteString(w, i.IdentifierName.Data)
		w.Write(as)
	}

	io.WriteString(w, i.ImportedBinding.Data)
}

func (oe OptionalExpression) printSource(w io.Writer, v bool) {
	if oe.MemberExpression != nil {
		oe.MemberExpression.printSource(w, v)
	} else if oe.CallExpression != nil {
		oe.CallExpression.printSource(w, v)
	} else if oe.OptionalExpression != nil {
		oe.OptionalExpression.printSource(w, v)
	}

	oe.OptionalChain.printSource(w, v)
}

func (oe OptionalChain) printSource(w io.Writer, v bool) {
	if oe.OptionalChain != nil {
		oe.OptionalChain.printSource(w, v)
	} else {
		w.Write(optionalChain)
	}

	if oe.Arguments != nil {
		oe.Arguments.printSource(w, v)
	} else if oe.Expression != nil {
		w.Write(bracketOpen)
		oe.Expression.printSource(w, v)
		w.Write(bracketClose)
	} else if oe.IdentifierName != nil {
		if oe.OptionalChain != nil {
			w.Write(dot)
		}

		io.WriteString(w, oe.IdentifierName.Data)
	} else if oe.TemplateLiteral != nil {
		oe.TemplateLiteral.printSource(w, v)
	} else if oe.PrivateIdentifier != nil {
		if oe.OptionalChain != nil {
			w.Write(dot)
		}

		io.WriteString(w, oe.PrivateIdentifier.Data)
	}
}

func (ce CoalesceExpression) printSource(w io.Writer, v bool) {
	if ce.CoalesceExpressionHead != nil {
		ce.CoalesceExpressionHead.printSource(w, v)
		w.Write(coalesceOperator)
	}

	ce.BitwiseORExpression.printSource(w, v)
}
