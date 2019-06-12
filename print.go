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
