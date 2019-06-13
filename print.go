package javascript

import "io"

var (
	blockOpen    = []byte{'{'}
	blockClose   = []byte{'}'}
	commaSep     = []byte{',', ' '}
	commaSepNL   = []byte{',', ' ', '\n'}
	newLine      = []byte{'\n'}
	labelPost    = []byte{':', ' '}
	semiColon    = []byte{';'}
	ifOpen       = []byte{'i', 'f', ' ', '('}
	parenClose   = []byte{')', ' '}
	elseOpen     = []byte{' ', 'e', 'l', 's', 'e', ' '}
	doOpen       = []byte{'d', 'o', ' '}
	doWhileOpen  = []byte{' ', 'w', 'h', 'i', 'l', 'e', ' ', '('}
	doWhileClose = []byte{')', ';'}
	whileOpen    = doWhileOpen[1:]
	forOpen      = []byte{'f', 'o', 'r', ' ', '('}
	forAwaitOpen = []byte{'f', 'o', 'r', ' ', 'a', 'w', 'a', 'i', 't', ' ', '('}
	switchOpen   = []byte{'s', 'w', 'i', 't', 'c', 'h', ' ', '('}
	switchClose  = []byte{')', ' ', '{'}
	caseOpen     = []byte{'c', 'a', 's', 'e', ' '}
	caseClose    = labelPost[:1]
	withOpen     = []byte{'w', 'i', 't', 'h', ' ', '('}
	forIn        = []byte{' ', 'i', 'n', ' '}
	forOf        = []byte{' ', 'o', 'f', ' '}
	varOpen      = []byte{'v', 'a', 'r', ' '}
	letOpen      = []byte{'l', 'e', 't', ' '}
	constOpen    = []byte{'c', 'o', 'n', 's', 't', ' '}
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
	ip := indentPrinter{w}
	for n, stmt := range b.StatementListItems {
		if n > 0 {
			if v {
				if len(stmt.Tokens) > 0 {
					if ll := stmt.Tokens[0].Line; ll > lastLine {
						w.Write(newLine)
					} else {
						w.Write(space)
					}
				} else {
					w.Write(newLine)
				}
			} else {
				w.Write(space)
			}
		}
		stmt.printSource(&ip, v)
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
		i.InitVar.printSource(&pp, v)
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

}

func (ws WithStatement) printSource(w io.Writer, v bool) {

}

func (f FunctionDeclaration) printSource(w io.Writer, v bool) {

}

func (t TryStatement) printSource(w io.Writer, v bool) {

}

func (c ClassDeclaration) printSource(w io.Writer, v bool) {

}

func (l LexicalDeclaration) printSource(w io.Writer, v bool) {

}

func (vd VariableDeclaration) printSource(w io.Writer, v bool) {

}

func (a AssignmentExpression) printSource(w io.Writer, v bool) {

}

func (l LeftHandSideExpression) printSource(w io.Writer, v bool) {

}

func (o ObjectBindingPattern) printSource(w io.Writer, v bool) {

}

func (a ArrayBindingPattern) printSource(w io.Writer, v bool) {

}
