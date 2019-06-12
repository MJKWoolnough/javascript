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
