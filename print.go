package javascript

import "io"

var (
	blockOpen                  = []byte{'{'}
	blockClose                 = []byte{'}'}
	commaSep                   = []byte{',', ' '}
	commaSepNL                 = []byte{',', '\n'}
	newLine                    = []byte{'\n'}
	conditionalStart           = []byte{' ', '?', ' '}
	conditionalSep             = []byte{' ', ':', ' '}
	labelPost                  = conditionalSep[1:]
	semiColon                  = []byte{';'}
	ifOpen                     = []byte{'i', 'f', ' ', '('}
	parenClose                 = []byte{')', ' '}
	elseOpen                   = []byte{' ', 'e', 'l', 's', 'e', ' '}
	doOpen                     = []byte{'d', 'o', ' '}
	doWhileOpen                = []byte{' ', 'w', 'h', 'i', 'l', 'e', ' ', '('}
	doWhileClose               = []byte{')', ';'}
	whileOpen                  = doWhileOpen[1:]
	forOpen                    = []byte{'f', 'o', 'r', ' ', '('}
	forAwaitOpen               = []byte{'f', 'o', 'r', ' ', 'a', 'w', 'a', 'i', 't', ' ', '('}
	switchOpen                 = []byte{'s', 'w', 'i', 't', 'c', 'h', ' ', '('}
	switchClose                = []byte{')', ' ', '{'}
	caseOpen                   = []byte{'c', 'a', 's', 'e', ' '}
	caseClose                  = labelPost[:1]
	defaultCase                = []byte{'d', 'e', 'f', 'a', 'u', 'l', 't', ':', '\n'}
	withOpen                   = []byte{'w', 'i', 't', 'h', ' ', '('}
	forIn                      = []byte{' ', 'i', 'n', ' '}
	forOf                      = []byte{' ', 'o', 'f', ' '}
	varOpen                    = []byte{'v', 'a', 'r', ' '}
	letOpen                    = []byte{'l', 'e', 't', ' '}
	constOpen                  = []byte{'c', 'o', 'n', 's', 't', ' '}
	funcOpen                   = []byte{'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', ' '}
	asyncFuncOpen              = []byte{'a', 's', 'y', 'n', 'c', ' ', 'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', ' '}
	genFuncOpen                = []byte{'f', 'u', 'n', 'c', 't', 'i', 'o', 'n', '*', ' '}
	parenOpen                  = []byte{'('}
	tryOpen                    = []byte{'t', 'r', 'y', ' '}
	catchParenOpen             = []byte{' ', 'c', 'a', 't', 'c', 'h', ' ', '('}
	catchOpen                  = catchParenOpen[:7]
	finallyOpen                = []byte{' ', 'f', 'i', 'n', 'a', 'l', 'l', 'y', ' '}
	classOpen                  = []byte{'c', 'l', 'a', 's', 's', ' '}
	extends                    = []byte{'e', 'x', 't', 'e', 'n', 'd', 's', ' '}
	assignment                 = []byte{' ', '=', ' '}
	assignmentMultiply         = []byte{' ', '*', '=', ' '}
	assignmentDivide           = []byte{' ', '/', '=', ' '}
	assignmentRemainder        = []byte{' ', '%', '=', ' '}
	assignmentAdd              = []byte{' ', '+', '=', ' '}
	assignmentSubtract         = []byte{' ', '-', '=', ' '}
	assignmentLeftShift        = []byte{' ', '<', '<', '=', ' '}
	assignmentSignRightShift   = []byte{' ', '>', '>', '=', ' '}
	assignmentZeroRightShift   = []byte{' ', '>', '>', '>', '=', ' '}
	assignmentAND              = []byte{' ', '&', '=', ' '}
	assignmentXOR              = []byte{' ', '^', '=', ' '}
	assignmentOR               = []byte{' ', '|', '=', ' '}
	assignmentExponentiation   = []byte{' ', '*', '*', '=', ' '}
	yield                      = []byte{'y', 'i', 'e', 'l', 'd', ' '}
	delegate                   = []byte{'*', ' '}
	ellipsis                   = []byte{'.', '.', '.'}
	bracketOpen                = []byte{'['}
	bracketClose               = []byte{']'}
	methodStaticAsyncGenerator = []byte{'s', 't', 'a', 't', 'i', 'c', ' ', 'a', 's', 'y', 'n', 'c', ' ', '*', ' '}
	methodAsyncGenerator       = methodStaticAsyncGenerator[7:]
	methodStatic               = methodStaticAsyncGenerator[:7]
	methodAsync                = methodStaticAsyncGenerator[7:13]
	methodStaticAsync          = methodStaticAsyncGenerator[:13]
	methodGenerator            = methodStaticAsyncGenerator[13:15]
	methodStaticGenerator      = []byte{'s', 't', 'a', 't', 'i', 'c', ' ', '*', ' '}
	methodStaticGet            = []byte{'s', 't', 'a', 't', 'i', 'c', ' ', 'g', 'e', 't', ' '}
	methodStaticSet            = []byte{'s', 't', 'a', 't', 'i', 'c', ' ', 's', 'e', 't', ' '}
	methodGet                  = methodStaticGet[8:]
	methodSet                  = methodStaticSet[8:]
)

func (s Script) printSource(w io.Writer, v bool) {
	for _, stmt := range s.StatementList {
		stmt.printSource(w, v)
		w.Write(newLine)
	}
}

func (s StatementListItem) printSource(w io.Writer, v bool) {
	if s.Statement != nil {
		s.Statement.printSource(w, v)
	} else if s.Declaration != nil {
		s.Declaration.printSource(w, v)
	}
	w.Write(newLine)
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
			io.WriteString(w, s.LabelIdentifier.Identifier.Data)
			w.Write(labelPost)
			if s.LabelledItemFunction != nil {
				s.LabelledItemFunction.printSource(w, v)
			} else if s.LabelledItemStatement != nil {
				s.LabelledItemStatement.printSource(w, v)
			}
		} else if s.TryStatement != nil {
			s.TryStatement.printSource(w, v)
		} else if s.DebuggerStatement != nil {
			io.WriteString(w, "debugger;")
		}
	case StatementContinue:
		if s.LabelIdentifier == nil {
			io.WriteString(w, "continue;")
		} else {
			io.WriteString(w, "continue ")
			io.WriteString(w, s.LabelIdentifier.Identifier.Data)
			w.Write(semiColon)
		}
	case StatementBreak:
		if s.LabelIdentifier == nil {
			io.WriteString(w, "break;")
		} else {
			io.WriteString(w, "break ")
			io.WriteString(w, s.LabelIdentifier.Identifier.Data)
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
	for _, stmt := range b.StatementListItems {
		if v {
			if len(stmt.Tokens) > 0 {
				if ll := stmt.Tokens[0].Line; ll > lastLine {
					pp.Write(newLine)
				} else {
					pp.Write(space)
				}
			} else {
				pp.Write(newLine)
			}
		} else {
			pp.Write(newLine)
		}
		stmt.printSource(&pp, v)
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
		LexicalBinding(vd).printSource(w, v)
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
	for n, ae := range e.Expressions {
		if n > 0 {
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
		}
		ae.printSource(w, v)
	}
}

func (i IfStatement) printSource(w io.Writer, v bool) {
	w.Write(ifOpen)
	if v {
		pp := indentPrinter{w}
		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			pp.Write(newLine)
		}
		i.Expression.printSource(&pp, true)
		if len(i.Expression.Tokens) > 0 && i.Expression.Tokens[len(i.Expression.Tokens)-1].Line > i.Expression.Tokens[0].Line {
			w.Write(newLine)
		}
	} else {
		i.Expression.printSource(w, false)
	}
	w.Write(parenClose)
	i.Statement.printSource(w, v)
	if i.ElseStatement != nil {
		w.Write(elseOpen)
		i.Expression.printSource(w, v)
	}
}

func (i IterationStatementDo) printSource(w io.Writer, v bool) {
	w.Write(doOpen)
	i.Statement.printSource(w, v)
	w.Write(doWhileOpen)
	if v {
		pp := indentPrinter{w}
		if len(i.Statement.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Statement.Tokens[len(i.Statement.Tokens)-1].Line {
			pp.Write(newLine)
		}
		i.Expression.printSource(&pp, true)
		if len(i.Expression.Tokens) > 0 && i.Expression.Tokens[len(i.Expression.Tokens)-1].Line > i.Expression.Tokens[0].Line {
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
		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			pp.Write(newLine)
		}
		i.Expression.printSource(&pp, true)
		if len(i.Expression.Tokens) > 0 && i.Expression.Tokens[len(i.Expression.Tokens)-1].Line > i.Expression.Tokens[0].Line {
			w.Write(newLine)
		}
	} else {
		i.Expression.printSource(w, false)
	}
	w.Write(parenClose)
	i.Statement.printSource(w, v)
}

func (i IterationStatementFor) printSource(w io.Writer, v bool) {
	switch i.Type {
	case ForNormal:
		if i.InitVar != nil || i.InitLexical != nil || i.InitExpression != nil {
			return
		}
	case ForNormalVar:
		if i.InitVar == nil {
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
	switch i.Type {
	case ForNormal:
		w.Write(semiColon)
	case ForNormalVar:
		if v && len(i.InitVar.Tokens) > 0 {
			if i.InitVar.Tokens[0].Line > lastLine {
				pp.Write(newLine)
			}
			lastLine = i.InitVar.Tokens[len(i.InitVar.Tokens)-1].Line
		}
		(*LexicalBinding)(i.InitVar).printSource(&pp, v)
	case ForNormalLexicalDeclaration:
		i.InitLexical.printSource(w, v)
		if v && len(i.InitLexical.Tokens) > 0 {
			if i.InitLexical.Tokens[0].Line > lastLine {
				pp.Write(newLine)
			}
			lastLine = i.InitLexical.Tokens[len(i.InitLexical.Tokens)-1].Line
		}
		i.InitLexical.printSource(&pp, v)
	case ForNormalExpression:
		if v && len(i.InitLexical.Tokens) > 0 {
			if i.InitExpression.Tokens[0].Line > lastLine {
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
					pp.Write(newLine)
				}
				lastLine = i.LeftHandSideExpression.Tokens[len(i.LeftHandSideExpression.Tokens)-1].Line
			}
		}
		i.LeftHandSideExpression.printSource(&pp, v)
		w.Write(semiColon)
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
			io.WriteString(w, i.ForBindingIdentifier.Identifier.Data)
		} else if i.ForBindingPatternObject != nil {
			i.ForBindingPatternObject.printSource(w, v)
		} else {
			i.ForBindingPatternArray.printSource(w, v)
		}
	}
	switch i.Type {
	case ForNormal, ForNormalVar, ForNormalLexicalDeclaration, ForNormalExpression:
		if i.Conditional != nil {
			w.Write(space)
			if v && len(i.Conditional.Tokens) > 0 {
				if i.Conditional.Tokens[0].Line > lastLine {
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
					pp.Write(newLine)
				} else {
					w.Write(space)
				}
				lastLine = i.Afterthought.Tokens[len(i.Conditional.Tokens)-1].Line
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
	w.Write(parenClose)
	i.Statement.printSource(w, v)
}

func (s SwitchStatement) printSource(w io.Writer, v bool) {
	w.Write(switchOpen)
	if v {
		pp := indentPrinter{w}
		if len(s.Tokens) > 0 && len(s.Expression.Tokens) > 0 && s.Expression.Tokens[0].Line > s.Tokens[0].Line {
			pp.Write(newLine)
		}
		s.Expression.printSource(&pp, true)
		if len(s.Expression.Tokens) > 0 && s.Expression.Tokens[len(s.Expression.Tokens)-1].Line > s.Expression.Tokens[0].Line {
			w.Write(newLine)
		}
	} else {
		s.Expression.printSource(w, false)
	}
	w.Write(switchClose)
	for _, c := range s.CaseClauses {
		c.printSource(w, v)
	}
	if s.DefaultClause != nil {
		w.Write(defaultCase)
		pp := indentPrinter{w}
		for _, stmt := range s.DefaultClause {
			stmt.printSource(&pp, v)
			w.Write(newLine)
		}
	}
	for _, c := range s.PostDefaultCaseClauses {
		c.printSource(w, v)
	}
	w.Write(blockClose)
}

func (ws WithStatement) printSource(w io.Writer, v bool) {
	w.Write(withOpen)
	if v {
		pp := indentPrinter{w}
		if len(ws.Tokens) > 0 && len(ws.Expression.Tokens) > 0 && ws.Expression.Tokens[0].Line > ws.Tokens[0].Line {
			pp.Write(newLine)
		}
		ws.Expression.printSource(&pp, true)
		if len(ws.Expression.Tokens) > 0 && ws.Expression.Tokens[len(ws.Expression.Tokens)-1].Line > ws.Expression.Tokens[0].Line {
			w.Write(newLine)
		}
	} else {
		ws.Expression.printSource(w, false)
	}
	w.Write(parenClose)
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
	default:
		return
	}
	if f.BindingIdentifier != nil {
		io.WriteString(w, f.BindingIdentifier.Identifier.Data)
	}
	w.Write(parenOpen)
	f.FormalParameters.printSource(&indentPrinter{w}, v)
	w.Write(parenClose)
	f.FunctionBody.printSource(w, v)
}

func (t TryStatement) printSource(w io.Writer, v bool) {
	w.Write(tryOpen)
	t.TryBlock.printSource(w, v)
	if t.CatchBlock != nil {
		if t.CatchParameterBindingIdentifier != nil {
			w.Write(catchParenOpen)
			io.WriteString(w, t.CatchParameterBindingIdentifier.Identifier.Data)
			w.Write(parenClose)
		} else if t.CatchParameterArrayBindingPattern != nil {
			w.Write(catchParenOpen)
			t.CatchParameterArrayBindingPattern.printSource(w, v)
			w.Write(parenClose)
		} else if t.CatchParameterObjectBindingPattern != nil {
			w.Write(catchParenOpen)
			t.CatchParameterObjectBindingPattern.printSource(w, v)
			w.Write(parenClose)
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
		io.WriteString(w, c.BindingIdentifier.Identifier.Data)
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
		pp.Write(newLine)
		for _, md := range c.ClassBody {
			md.printSource(&pp, v)
			pp.Write(newLine)
		}
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
		io.WriteString(w, l.BindingIdentifier.Identifier.Data)
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
		case AssignmentSignPropagatinRightShift:
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
		default:
			return
		}
		a.LeftHandSideExpression.printSource(w, v)
		w.Write(ao)
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
		io.WriteString(w, o.BindingRestProperty.Identifier.Data)
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
	for n, be := range f.FormalParameterList {
		if n > 0 {
			w.Write(commaSep)
		}
		be.printSource(w, v)
	}
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
	case MethodStatic:
		w.Write(methodStatic)
	case MethodStaticGenerator:
		w.Write(methodStaticGenerator)
	case MethodStaticAsync:
		w.Write(methodStaticAsync)
	case MethodStaticAsyncGenerator:
		w.Write(methodStaticAsyncGenerator)
	case MethodStaticGetter:
		w.Write(methodStaticGet)
	case MethodStaticSetter:
		w.Write(methodStaticSet)
	default:
		return
	}
	m.PropertyName.printSource(w, v)
	w.Write(parenOpen)
	m.Params.printSource(w, v)
	w.Write(parenClose)
	m.FunctionBody.printSource(w, v)
}

func (c ConditionalExpression) printSource(w io.Writer, v bool) {
	c.LogicalORExpression.printSource(w, v)
	if c.True != nil && c.False != nil {
		w.Write(conditionalStart)
		c.True.printSource(w, v)
		w.Write(conditionalSep)
		c.False.printSource(w, v)
	}
}

func (a ArrowFunction) printSource(w io.Writer, v bool) {

}

func (n NewExpression) printSource(w io.Writer, v bool) {

}

func (c CallExpression) printSource(w io.Writer, v bool) {

}

func (b BindingProperty) printSource(w io.Writer, v bool) {

}

func (b BindingElement) printSource(w io.Writer, v bool) {

}

func (p PropertyName) printSource(w io.Writer, v bool) {

}

func (l LogicalORExpression) printSource(w io.Writer, v bool) {

}
