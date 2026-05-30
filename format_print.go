package javascript

func (s Script) printSource(w writer, v bool) {
	w.Start(s.Tokens)
	defer w.End()

	if v && len(s.Comments[0]) > 0 {
		s.Comments[0].printSource(w, true, true)
		w.WriteString("\n")
	}

	if len(s.StatementList) > 0 {
		s.StatementList[0].printSource(w, v)

		for _, stmt := range s.StatementList[1:] {
			w.WriteString("\n\n")
			stmt.printSource(w, v)
		}
	}

	if v && len(s.Comments[1]) > 0 {
		w.WriteString("\n")
		s.Comments[1].printSource(w, false, false)
	}
}

func (s StatementListItem) printSource(w writer, v bool) {
	w.Start(s.Tokens)
	defer w.End()

	if v {
		s.Comments[0].printSource(w, s.Statement != nil || s.Declaration != nil, false)
	}

	if s.Statement != nil {
		s.Statement.printSource(w, v)
	} else if s.Declaration != nil {
		s.Declaration.printSource(w, v)
	}

	if v {
		s.Comments[1].printSource(w, false, false)
	}
}

func (s Statement) printSource(w writer, v bool) {
	w.Start(s.Tokens)
	defer w.End()

	switch s.Type {
	case StatementNormal:
		if s.BlockStatement != nil {
			s.BlockStatement.printSource(w, v)
		} else if s.VariableStatement != nil {
			s.VariableStatement.printSource(w, v)
		} else if s.ExpressionStatement != nil {
			s.ExpressionStatement.printSource(w, v)
			w.PrintSemiColon()
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
			w.WriteToken(s.LabelIdentifier)

			if v {
				s.Comments[0].printSource(w, true, false)
			}

			w.WriteStringWithType(":", TokenPunctuator)
			w.WriteString(" ")

			if v {
				s.Comments[1].printSource(w, true, false)
			}

			if s.LabelledItemFunction != nil {
				s.LabelledItemFunction.printSource(w, v)
			} else if s.LabelledItemStatement != nil {
				s.LabelledItemStatement.printSource(w, v)
			}
		} else if s.TryStatement != nil {
			s.TryStatement.printSource(w, v)
		} else {
			w.WriteStringWithType(";", TokenPunctuator)
		}
	case StatementContinue, StatementBreak:
		if s.Type == StatementContinue {
			w.WriteStringWithType("continue", TokenKeyword)
		} else {
			w.WriteStringWithType("break", TokenKeyword)
		}

		if v {
			s.Comments[0].printSource(w, false, false)
		}

		if s.LabelIdentifier != nil {
			if !w.LastIsWhitespace() {
				w.WriteString(" ")
			}

			w.WriteToken(s.LabelIdentifier)
		}

		if v {
			s.Comments[1].printSource(w, false, false)
		}

		w.PrintSemiColon()
	case StatementReturn:
		if s.ExpressionStatement == nil {
			w.WriteStringWithType("return", TokenKeyword)

			if v {
				s.Comments[0].printSource(w, false, false)
			}

			w.PrintSemiColon()
		} else {
			w.WriteStringWithType("return", TokenKeyword)
			w.WriteString(" ")
			s.ExpressionStatement.printSource(w, v)
			w.PrintSemiColon()
		}
	case StatementThrow:
		if s.ExpressionStatement != nil {
			w.WriteStringWithType("throw", TokenKeyword)
			w.WriteString(" ")
			s.ExpressionStatement.printSource(w, v)
			w.PrintSemiColon()
		}
	case StatementDebugger:
		w.WriteStringWithType("debugger", TokenKeyword)

		if v {
			s.Comments[0].printSource(w, false, false)
		}

		w.PrintSemiColon()
	}
}

func (d Declaration) printSource(w writer, v bool) {
	w.Start(d.Tokens)
	defer w.End()

	if d.ClassDeclaration != nil {
		if v {
			d.Comments.printSource(w, true, false)
		}

		d.ClassDeclaration.printSource(w, v)
	} else if d.FunctionDeclaration != nil {
		d.FunctionDeclaration.printSource(w, v)
	} else if d.LexicalDeclaration != nil {
		d.LexicalDeclaration.printSource(w, v)
	}
}

func (b Block) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	if v && len(b.Comments[0]) > 0 {
		b.Comments[0].printSource(w, false, true)
	}

	ip := w.Indent()

	if len(b.StatementList) > 0 {
		for _, stmt := range b.StatementList {
			ip.WriteString("\n")
			stmt.printSource(ip, v)
		}
	}

	if v && len(b.Comments[1]) > 0 {
		w.WriteString("\n")
		b.Comments[1].printSource(w, false, true)
	} else if len(b.StatementList) > 0 && !w.LastIsWhitespace() {
		w.WriteString("\n")
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (vs VariableStatement) printSource(w writer, v bool) {
	w.Start(vs.Tokens)
	defer w.End()

	if len(vs.VariableDeclarationList) == 0 {
		return
	}

	w.WriteStringWithType("var", TokenKeyword)
	w.WriteString(" ")

	var lastLine uint64

	if v && len(vs.Tokens) > 0 {
		lastLine = vs.Tokens[0].Line
	}

	for n, vd := range vs.VariableDeclarationList {
		if n > 0 {
			w.WriteStringWithType(",", TokenPunctuator)
			if v && len(vd.Tokens) > 0 {
				if ll := vd.Tokens[0].Line; ll > lastLine {
					lastLine = ll

					w.WriteString("\n")
				} else {
					w.WriteString(" ")
				}
			} else {
				w.WriteString(" ")
			}
		}

		vd.printSource(w, v)
	}

	w.PrintSemiColon()
}

func (e Expression) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	if len(e.Expressions) == 0 {
		return
	}

	var lastLine uint64

	if v && len(e.Tokens) > 0 {
		lastLine = e.Tokens[0].Line
	}

	e.Expressions[0].printSource(w, v)

	for _, ae := range e.Expressions[1:] {
		w.WriteStringWithType(",", TokenPunctuator)
		if v && len(ae.Tokens) > 0 {
			if ll := ae.Tokens[0].Line; ll > lastLine {
				lastLine = ll
				w.WriteString("\n")
			} else {
				w.WriteString(" ")
			}
		} else {
			w.WriteString(" ")
		}

		ae.printSource(w, v)
	}
}

func (i IfStatement) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	w.WriteStringWithType("if", TokenKeyword)

	if v && len(i.Comments[0]) > 0 {
		i.Comments[0].printSource(w, true, false)
	} else {
		w.WriteString(" ")
	}

	w.WriteStringWithType("(", TokenPunctuator)

	if v {
		ip := w

		if hasSingleLineComment(i.Comments[1:3]) || i.Expression.hasSingleLineComment() {
			ip = w.Indent()

			i.Comments[1].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			i.Comments[1].printSource(w, true, false)
		}

		i.Expression.printSource(ip, v)

		if w != ip {
			w.WriteString("\n")
		}

		i.Comments[2].printSource(w, false, w != ip)
	} else {
		i.Expression.printSource(w, v)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")

	if v {
		i.Comments[3].printSource(w, true, false)
	}

	i.Statement.printSource(w, v)

	if i.ElseStatement != nil {
		if v && len(i.Comments[4]) > 0 {
			i.Comments[4].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}

		w.WriteStringWithType("else", TokenKeyword)
		w.WriteString(" ")

		if v {
			i.Comments[5].printSource(w, true, false)
		}

		i.ElseStatement.printSource(w, v)
	}
}

func (i IterationStatementDo) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	w.WriteStringWithType("do", TokenKeyword)
	w.WriteString(" ")

	if v {
		i.Comments[0].printSource(w, true, false)
	}

	i.Statement.printSource(w, v)

	if v {
		i.Comments[1].printSource(w, true, false)
	}

	if !w.LastIsWhitespace() {
		w.WriteString(" ")
	}

	w.WriteStringWithType("while", TokenKeyword)
	w.WriteString(" ")

	if v {
		i.Comments[2].printSource(w, true, false)
	}

	w.WriteStringWithType("(", TokenPunctuator)

	if v {
		ip := w

		if hasSingleLineComment(i.Comments[3:4]) || i.Expression.hasSingleLineComment() {
			ip = w.Indent()

			i.Comments[3].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			i.Comments[3].printSource(w, true, false)
		}

		i.Expression.printSource(ip, v)

		if w != ip {
			w.WriteString("\n")
		}

		i.Comments[4].printSource(w, false, w != ip)
	} else {
		i.Expression.printSource(w, v)
	}

	w.WriteStringWithType(")", TokenPunctuator)

	if v {
		i.Comments[5].printSource(w, true, false)
	}

	w.PrintSemiColon()
}

func (i IterationStatementWhile) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	w.WriteStringWithType("while", TokenKeyword)
	w.WriteString(" ")

	if v {
		i.Comments[0].printSource(w, true, false)
	}

	w.WriteStringWithType("(", TokenPunctuator)

	if v {
		ip := w

		if hasSingleLineComment(i.Comments[1:3]) || i.Expression.hasSingleLineComment() {
			ip = w.Indent()

			i.Comments[1].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			i.Comments[1].printSource(w, true, false)
		}

		i.Expression.printSource(ip, v)

		if w != ip {
			w.WriteString("\n")
		}

		i.Comments[2].printSource(w, false, w != ip)
	} else {
		i.Expression.printSource(w, v)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")

	if v {
		i.Comments[3].printSource(w, true, false)
	}

	i.Statement.printSource(w, v)
}

func (i IterationStatementFor) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

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

	w.WriteStringWithType("for", TokenKeyword)
	w.WriteString(" ")

	if v {
		i.Comments[0].printSource(w, true, false)
	}

	switch i.Type {
	case ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		w.WriteStringWithType("await", TokenKeyword)
		w.WriteString(" ")

		if v {
			i.Comments[1].printSource(w, true, false)
		}
	}

	w.WriteStringWithType("(", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(i.Comments[2:7]) {
			ip = w.Indent()
			sep = "\n"

			i.Comments[2].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			i.Comments[2].printSource(w, true, false)
		}

		i.Comments[3].printSource(ip, false, false)
	}

	switch i.Type {
	case ForNormal:
		ip.WriteStringWithType(";", TokenPunctuator)
	case ForNormalVar:
		ip.WriteStringWithType("var", TokenKeyword)
		ip.WriteString(" ")
		LexicalBinding(i.InitVar[0]).printSource(ip, v)

		for _, vd := range i.InitVar[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			vd.printSource(ip, v)
		}

		ip.WriteStringWithType(";", TokenPunctuator)
	case ForNormalLexicalDeclaration:
		i.InitLexical.printSource(ip, v)

		if ip.LastChar() == '\n' {
			ip.WriteStringWithType(";", TokenPunctuator)
		}
	case ForNormalExpression:
		i.InitExpression.printSource(ip, v)
		ip.WriteStringWithType(";", TokenPunctuator)
	case ForInLeftHandSide, ForOfLeftHandSide, ForAwaitOfLeftHandSide:
		i.LeftHandSideExpression.printSource(ip, v)
	default:
		switch i.Type {
		case ForInVar, ForOfVar, ForAwaitOfVar:
			ip.WriteStringWithType("var", TokenKeyword)
			ip.WriteString(" ")
		case ForInLet, ForOfLet, ForAwaitOfLet:
			ip.WriteStringWithType("let", TokenKeyword)
			ip.WriteString(" ")
		case ForInConst, ForOfConst, ForAwaitOfConst:
			ip.WriteStringWithType("const", TokenKeyword)
			ip.WriteString(" ")
		}

		if v {
			i.Comments[4].printSource(ip, true, false)
		}

		if i.ForBindingIdentifier != nil {
			ip.WriteToken(i.ForBindingIdentifier)
		} else if i.ForBindingPatternObject != nil {
			i.ForBindingPatternObject.printSource(ip, v)
		} else {
			i.ForBindingPatternArray.printSource(ip, v)
		}
	}

	switch i.Type {
	case ForNormal, ForNormalVar, ForNormalLexicalDeclaration, ForNormalExpression:
		if v {
			i.Comments[4].printSource(ip, true, true)
		}

		if i.Conditional != nil {
			if !v || !i.Conditional.hasFirstComment() {
				ip.WriteString(" ")
			}

			i.Conditional.printSource(ip, v)
		}

		ip.WriteStringWithType(";", TokenPunctuator)

		if v {
			i.Comments[5].printSource(ip, true, false)
		}

		if i.Afterthought != nil {
			if !v || !i.Afterthought.hasFirstComment() {
				ip.WriteString(" ")
			}

			i.Afterthought.printSource(ip, v)
		}
	case ForInLeftHandSide, ForInVar, ForInLet, ForInConst:
		if v && len(i.Comments[5]) > 0 {
			i.Comments[5].printSource(ip, true, false)
		}

		if !ip.LastIsWhitespace() {
			ip.WriteString(sep)
		}

		ip.WriteStringWithType("in", TokenKeyword)
		ip.WriteString(" ")
		i.In.printSource(ip, v)
	case ForOfLeftHandSide, ForOfVar, ForOfLet, ForOfConst, ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		if v && len(i.Comments[5]) > 0 {
			i.Comments[5].printSource(ip, true, false)
		} else {
			ip.WriteString(" ")
		}

		ip.WriteStringWithType("of", TokenKeyword)
		ip.WriteString(" ")
		i.Of.printSource(ip, v)
	}

	if v {
		if w != ip && (len(i.Comments[6]) > 0 || !ip.LastIsWhitespace()) {
			w.WriteString(sep)
		}

		i.Comments[6].printSource(w, true, false)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")

	if v {
		i.Comments[7].printSource(w, true, false)
	}

	i.Statement.printSource(w, v)
}

func (s SwitchStatement) printSource(w writer, v bool) {
	w.Start(s.Tokens)
	defer w.End()

	w.WriteStringWithType("switch", TokenKeyword)
	w.WriteString(" ")

	if v {
		s.Comments[0].printSource(w, true, false)
	}

	w.WriteStringWithType("(", TokenPunctuator)

	if v {
		ip := w

		if s.Comments[1].hasSingleLineComment() || s.Comments[2].hasSingleLineComment() || s.Expression.hasSingleLineComment() {
			ip = w.Indent()

			s.Comments[1].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			s.Comments[1].printSource(w, true, false)
		}

		s.Expression.printSource(ip, v)

		if w != ip {
			w.WriteString("\n")
		}

		s.Comments[2].printSource(w, false, w != ip)
	} else {
		s.Expression.printSource(w, v)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")

	if v {
		s.Comments[3].printSource(w, true, false)
	}

	w.WriteStringWithType("{", TokenPunctuator)

	if v {
		s.Comments[4].printSource(w, false, true)
	}

	for _, c := range s.CaseClauses {
		w.WriteString("\n")
		c.printSource(w, v)
	}

	if s.DefaultClause != nil {
		w.WriteString("\n")
		if v {
			s.Comments[5].printSource(w, false, true)
		}

		w.WriteStringWithType("default", TokenKeyword)

		if v {
			s.Comments[6].printSource(w, false, false)
		}

		w.WriteStringWithType(":", TokenPunctuator)

		if v {
			s.Comments[7].printSource(w, false, true)
		}

		ip := w.Indent()

		for _, stmt := range s.DefaultClause {
			ip.WriteString("\n")
			stmt.printSource(ip, v)
		}
	}

	for _, c := range s.PostDefaultCaseClauses {
		w.WriteString("\n")
		c.printSource(w, v)
	}

	if v && len(s.Comments[8]) > 0 {
		w.WriteString("\n")
		s.Comments[8].printSource(w, false, true)
	}

	if c := w.LastChar(); c != '\n' && c != '{' {
		w.WriteString("\n")
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (ws WithStatement) printSource(w writer, v bool) {
	w.Start(ws.Tokens)
	defer w.End()

	w.WriteStringWithType("with", TokenKeyword)
	w.WriteString(" ")

	if v {
		ws.Comments[0].printSource(w, true, false)
	}

	w.WriteStringWithType("(", TokenPunctuator)

	if v {
		ip := w

		if ws.Comments[1].hasSingleLineComment() || ws.Comments[2].hasSingleLineComment() || ws.Expression.hasSingleLineComment() {
			ip = w.Indent()

			ws.Comments[1].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			ws.Comments[1].printSource(w, true, false)
		}

		ws.Expression.printSource(ip, v)

		if w != ip {
			w.WriteString("\n")
		}

		ws.Comments[2].printSource(w, false, w != ip)
	} else {
		ws.Expression.printSource(w, v)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")

	if v {
		ws.Comments[3].printSource(w, true, false)
	}

	ws.Statement.printSource(w, v)
}

func (f FunctionDeclaration) printSource(w writer, v bool) {
	w.Start(f.Tokens)
	defer w.End()

	switch f.Type {
	case FunctionNormal:
		w.WriteStringWithType("function", TokenKeyword)
		w.WriteString(" ")

		if v {
			f.Comments[1].printSource(w, true, false)
		}
	case FunctionGenerator:
		w.WriteStringWithType("function", TokenKeyword)

		if v {
			f.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("*", TokenPunctuator)

		if v && len(f.Comments[2]) > 0 {
			f.Comments[2].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}
	case FunctionAsync:
		w.WriteStringWithType("async", TokenIdentifier)
		w.WriteString(" ")

		if v {
			f.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("function", TokenKeyword)
		w.WriteString(" ")

		if v {
			f.Comments[1].printSource(w, true, false)
		}
	case FunctionAsyncGenerator:
		w.WriteStringWithType("async", TokenIdentifier)
		w.WriteString(" ")

		if v {
			f.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("function", TokenKeyword)

		if v {
			f.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("*", TokenPunctuator)

		if v && len(f.Comments[2]) > 0 {
			f.Comments[2].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}
	default:
		return
	}

	if f.BindingIdentifier != nil {
		w.WriteToken(f.BindingIdentifier)
	}

	if v {
		f.Comments[3].printSource(w, true, false)
	}

	f.FormalParameters.printSource(w, v)

	if v {
		f.Comments[4].printSource(w, true, false)
	}

	f.FunctionBody.printSource(w, v)
}

func (t TryStatement) printSource(w writer, v bool) {
	w.Start(t.Tokens)
	defer w.End()

	w.WriteStringWithType("try", TokenKeyword)
	w.WriteString(" ")

	if v {
		t.Comments[0].printSource(w, true, false)
	}

	t.TryBlock.printSource(w, v)

	if t.CatchBlock != nil {
		w.WriteString(" ")

		if v {
			t.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("catch", TokenKeyword)
		w.WriteString(" ")

		if v {
			t.Comments[2].printSource(w, true, false)
		}

		if t.CatchParameterBindingIdentifier != nil || t.CatchParameterArrayBindingPattern != nil || t.CatchParameterObjectBindingPattern != nil {
			w.WriteStringWithType("(", TokenPunctuator)

			ip := w

			if v {
				if hasSingleLineComment(t.Comments[3:7]) || t.CatchParameterBindingIdentifier.hasSingleLineComment() || t.CatchParameterArrayBindingPattern.hasSingleLineComment() || t.CatchParameterObjectBindingPattern.hasSingleLineComment() {
					ip = w.Indent()

					t.Comments[3].printSource(w, false, true)
					ip.WriteString("\n")
				} else {
					t.Comments[3].printSource(w, true, false)
				}

				if len(t.Comments[4]) > 0 {
					t.Comments[4].printSource(ip, true, w != ip)
				}
			}

			if t.CatchParameterBindingIdentifier != nil {
				ip.WriteToken(t.CatchParameterBindingIdentifier)
			} else if t.CatchParameterArrayBindingPattern != nil {
				t.CatchParameterArrayBindingPattern.printSource(ip, v)
			} else if t.CatchParameterObjectBindingPattern != nil {
				t.CatchParameterObjectBindingPattern.printSource(ip, v)
			}

			if v {
				t.Comments[5].printSource(ip, false, true)

				if w != ip {
					w.WriteString("\n")
				}

				t.Comments[6].printSource(w, false, w != ip)
			}

			w.WriteStringWithType(")", TokenPunctuator)
			w.WriteString(" ")

			if v {
				t.Comments[7].printSource(w, true, false)
			}
		}

		t.CatchBlock.printSource(w, v)
	}

	if t.FinallyBlock != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		if v {
			t.Comments[8].printSource(w, true, false)
		}

		w.WriteStringWithType("finally", TokenKeyword)
		w.WriteString(" ")

		if v {
			t.Comments[9].printSource(w, true, false)
		}

		t.FinallyBlock.printSource(w, v)
	}
}

func (c ClassDeclaration) printSource(w writer, v bool) {
	w.Start(c.Tokens)
	defer w.End()

	w.WriteStringWithType("class", TokenKeyword)
	w.WriteString(" ")

	if v {
		c.Comments[0].printSource(w, true, false)
	}

	if c.BindingIdentifier != nil {
		w.WriteToken(c.BindingIdentifier)
		w.WriteString(" ")

		if v {
			c.Comments[1].printSource(w, true, false)
		}
	}

	if c.ClassHeritage != nil {
		w.WriteStringWithType("extends", TokenKeyword)
		w.WriteString(" ")
		c.ClassHeritage.printSource(w, v)
		w.WriteString(" ")
	}

	if v {
		c.Comments[2].printSource(w, true, false)
	}

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w.Indent()

	if v {
		c.Comments[3].printSource(w, false, true)
	}

	if len(c.ClassBody) > 0 {
		for _, ce := range c.ClassBody {
			if !v && ce.FieldDefinition == nil && ce.MethodDefinition == nil && ce.ClassStaticBlock == nil {
				continue
			}

			ip.WriteString("\n")
			ce.printSource(ip, v)
		}

		if w.LastChar() != '\n' {
			w.WriteString("\n")
		}
	}

	if v && len(c.Comments[4]) > 0 {
		w.WriteString("\n")
		c.Comments[4].printSource(w, false, true)
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (l LexicalDeclaration) printSource(w writer, v bool) {
	w.Start(l.Tokens)
	defer w.End()

	if len(l.BindingList) == 0 {
		return
	}

	switch l.LetOrConst {
	case Let:
		w.WriteStringWithType("let", TokenKeyword)
		w.WriteString(" ")
	case Const:
		w.WriteStringWithType("const", TokenKeyword)
		w.WriteString(" ")
	}

	l.BindingList[0].printSource(w, v)

	for _, lb := range l.BindingList[1:] {
		w.WriteStringWithType(",", TokenPunctuator)

		if v {
			w.WriteString("\n")
		} else {
			w.WriteString(" ")
		}

		lb.printSource(w, v)
	}

	w.PrintSemiColon()
}

func (l LexicalBinding) printSource(w writer, v bool) {
	w.Start(l.Tokens)
	defer w.End()

	if v {
		l.Comments[0].printSource(w, true, false)
	}

	if l.BindingIdentifier != nil {
		w.WriteToken(l.BindingIdentifier)
	} else if l.ArrayBindingPattern != nil {
		l.ArrayBindingPattern.printSource(w, v)
	} else if l.ObjectBindingPattern != nil {
		l.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if v {
		l.Comments[1].printSource(w, false, false)
	}

	if l.Initializer != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		l.Initializer.printSource(w, v)
	}
}

func (a AssignmentExpression) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if a.Yield && a.AssignmentExpression != nil {
		if v {
			a.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("yield", TokenKeyword)
		w.WriteString(" ")

		if a.Delegate {
			if v {
				a.Comments[1].printSource(w, true, false)
			}

			w.WriteStringWithType("*", TokenPunctuator)
			w.WriteString(" ")
		}

		a.AssignmentExpression.printSource(w, v)
	} else if a.ArrowFunction != nil {
		a.ArrowFunction.printSource(w, v)
	} else if a.LeftHandSideExpression != nil && a.AssignmentExpression != nil {
		ao := a.AssignmentOperator.String()
		if ao == "" || ao == unknown {
			return
		}

		a.LeftHandSideExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(ao, TokenPunctuator)
		w.WriteString(" ")
		a.AssignmentExpression.printSource(w, v)
	} else if a.AssignmentPattern != nil && a.AssignmentExpression != nil && a.AssignmentOperator == AssignmentAssign {
		a.AssignmentPattern.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		a.AssignmentExpression.printSource(w, v)
	} else if a.ConditionalExpression != nil {
		a.ConditionalExpression.printSource(w, v)
	}
}

func (l LeftHandSideExpression) printSource(w writer, v bool) {
	w.Start(l.Tokens)
	defer w.End()

	if l.NewExpression != nil {
		l.NewExpression.printSource(w, v)
	} else if l.CallExpression != nil {
		l.CallExpression.printSource(w, v)
	} else if l.OptionalExpression != nil {
		l.OptionalExpression.printSource(w, v)
	}

	if v {
		l.Comments.printSource(w, false, false)
	}
}

func (a AssignmentPattern) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if v {
		a.Comments[0].printSource(w, true, false)
	}

	if a.ArrayAssignmentPattern != nil {
		a.ArrayAssignmentPattern.printSource(w, v)
	} else if a.ObjectAssignmentPattern != nil {
		a.ObjectAssignmentPattern.printSource(w, v)
	}

	if v {
		a.Comments[1].printSource(w, true, false)
	}
}

func (a ArrayAssignmentPattern) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	w.WriteStringWithType("[", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if a.hasSingleLineComment() {
			sep = "\n"
			ip = w.Indent()

			a.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			a.Comments[0].printSource(w, true, false)
		}
	}

	if len(a.AssignmentElements) > 0 {
		a.AssignmentElements[0].printSource(ip, v)

		for _, ae := range a.AssignmentElements[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			ae.printSource(ip, v)
		}
	}

	if a.AssignmentRestElement != nil {
		if len(a.AssignmentElements) > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		if v {
			a.Comments[1].printSource(ip, true, false)
		}

		ip.WriteStringWithType("...", TokenPunctuator)
		a.AssignmentRestElement.printSource(ip, v)
	}

	if v {
		if w != ip && (len(a.Comments[2]) > 0 || !ip.LastIsWhitespace()) {
			w.WriteString("\n")
		}

		a.Comments[2].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("]", TokenPunctuator)
}

func (o ObjectAssignmentPattern) printSource(w writer, v bool) {
	w.Start(o.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if o.hasSingleLineComment() {
			sep = "\n"
			ip = w.Indent()

			o.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			o.Comments[0].printSource(w, true, false)
		}
	}

	if len(o.AssignmentPropertyList) > 0 {
		o.AssignmentPropertyList[0].printSource(ip, v)

		for _, ap := range o.AssignmentPropertyList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			ap.printSource(ip, v)
		}
	}

	if o.AssignmentRestElement != nil {
		if len(o.AssignmentPropertyList) > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		if v {
			o.Comments[1].printSource(ip, true, false)
		}

		ip.WriteStringWithType("...", TokenPunctuator)
		o.AssignmentRestElement.printSource(ip, v)
	}

	if v {
		if w != ip && (len(o.Comments[2]) > 0 || !ip.LastIsWhitespace()) {
			w.WriteString("\n")
		}

		o.Comments[2].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (a AssignmentElement) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	a.DestructuringAssignmentTarget.printSource(w, v)

	if a.Initializer != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		a.Initializer.printSource(w, v)
	}
}

func (a AssignmentProperty) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if v {
		a.Comments[0].printSource(w, true, false)
	}

	a.PropertyName.printSource(w, v)

	if v {
		a.Comments[1].printSource(w, true, false)
	}

	if a.DestructuringAssignmentTarget != nil {
		if v || a.DestructuringAssignmentTarget.LeftHandSideExpression == nil || a.PropertyName.LiteralPropertyName != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.CallExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.OptionalExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression != nil && len(a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News) == 0 && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference.Data != a.PropertyName.LiteralPropertyName.Data {
			w.WriteStringWithType(":", TokenPunctuator)
			w.WriteString(" ")
			a.DestructuringAssignmentTarget.printSource(w, v)
		}
	}

	if a.Initializer != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		a.Initializer.printSource(w, v)
	}
}

func (d DestructuringAssignmentTarget) printSource(w writer, v bool) {
	w.Start(d.Tokens)
	defer w.End()

	if d.LeftHandSideExpression != nil {
		d.LeftHandSideExpression.printSource(w, v)
	} else if d.AssignmentPattern != nil {
		d.AssignmentPattern.printSource(w, v)
	}
}

func (o ObjectBindingPattern) printSource(w writer, v bool) {
	w.Start(o.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if o.hasSingleLineComment() {
			sep = "\n"
			ip = w.Indent()

			o.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			o.Comments[0].printSource(w, true, false)
		}
	}

	for n, bp := range o.BindingPropertyList {
		if n > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		bp.printSource(ip, v)
	}

	if o.BindingRestProperty != nil {
		if len(o.BindingPropertyList) > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		if v {
			o.Comments[1].printSource(ip, true, false)
		}

		ip.WriteStringWithType("...", TokenPunctuator)

		if v {
			o.Comments[2].printSource(ip, true, false)
		}

		ip.WriteToken(o.BindingRestProperty)

		if v {
			o.Comments[3].printSource(ip, true, false)
		}
	}

	if v {
		if len(o.Comments[4]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		o.Comments[4].printSource(w, false, false)
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (a ArrayBindingPattern) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	w.WriteStringWithType("[", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if a.hasSingleLineComment() {
			sep = "\n"
			ip = w.Indent()

			a.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			a.Comments[0].printSource(w, false, true)
		}
	}

	for n, be := range a.BindingElementList {
		if n > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		be.printSource(ip, v)
	}

	if a.BindingRestElement != nil {
		if len(a.BindingElementList) > 0 {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}

		if v {
			a.Comments[1].printSource(ip, true, false)
		}

		ip.WriteStringWithType("...", TokenPunctuator)
		a.BindingRestElement.printSource(ip, v)
	}

	if v {
		if len(a.Comments[2]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		a.Comments[2].printSource(w, false, true)
	}

	w.WriteStringWithType("]", TokenPunctuator)
}

func (c CaseClause) printSource(w writer, v bool) {
	w.Start(c.Tokens)
	defer w.End()

	if v {
		c.Comments[0].printSource(w, false, true)
	}

	w.WriteStringWithType("case", TokenKeyword)
	w.WriteString(" ")
	c.Expression.printSource(w, v)
	w.WriteStringWithType(":", TokenPunctuator)

	if v {
		c.Comments[1].printSource(w, false, true)
	}

	ip := w.Indent()

	for _, stmt := range c.StatementList {
		ip.WriteString("\n")
		stmt.printSource(ip, v)
	}
}

func (f FormalParameters) printSource(w writer, v bool) {
	w.Start(f.Tokens)
	defer w.End()

	w.WriteStringWithType("(", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if f.hasSingleLineComment() {
			sep = "\n"
			ip = w.Indent()

			f.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			f.Comments[0].printSource(w, len(f.FormalParameterList) > 0 || f.BindingIdentifier != nil || f.ArrayBindingPattern != nil || f.ObjectBindingPattern != nil, false)
		}
	}

	if len(f.FormalParameterList) > 0 {
		f.FormalParameterList[0].printSource(ip, v)

		for _, be := range f.FormalParameterList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			be.printSource(ip, v)
		}

		if f.BindingIdentifier != nil || f.ArrayBindingPattern != nil || f.ObjectBindingPattern != nil {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
		}
	}

	if f.BindingIdentifier != nil {
		if v {
			f.Comments[1].printSource(ip, false, true)
		}

		ip.WriteStringWithType("...", TokenPunctuator)

		if v {
			f.Comments[2].printSource(ip, false, true)
		}

		ip.WriteToken(f.BindingIdentifier)

		if v {
			f.Comments[3].printSource(ip, false, true)
		}
	} else if f.ArrayBindingPattern != nil {
		if v {
			f.Comments[1].printSource(ip, false, true)
		}

		ip.WriteStringWithType("...", TokenPunctuator)

		if v {
			f.Comments[2].printSource(ip, false, true)
		}

		f.ArrayBindingPattern.printSource(ip, v)

		if v {
			f.Comments[3].printSource(ip, false, true)
		}
	} else if f.ObjectBindingPattern != nil {
		if v {
			f.Comments[1].printSource(ip, false, true)
		}

		ip.WriteStringWithType("...", TokenPunctuator)

		if v {
			f.Comments[2].printSource(ip, false, true)
		}

		f.ObjectBindingPattern.printSource(ip, v)

		if v {
			f.Comments[3].printSource(ip, false, true)
		}
	}

	if v {
		if w != ip && (len(f.FormalParameterList) > 0 || f.BindingIdentifier != nil || f.ArrayBindingPattern != nil || f.ObjectBindingPattern != nil) {
			w.WriteString("\n")
		}

		f.Comments[4].printSource(w, false, false)
	}

	w.WriteStringWithType(")", TokenPunctuator)
	w.WriteString(" ")
}

func (m MethodDefinition) printSource(w writer, v bool) {
	w.Start(m.Tokens)
	defer w.End()

	switch m.Type {
	case MethodNormal:
	case MethodGenerator:
		if v {
			m.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("*", TokenPunctuator)
		w.WriteString(" ")
	case MethodAsync:
		if v {
			m.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("async", TokenIdentifier)
		w.WriteString(" ")
	case MethodAsyncGenerator:
		if v {
			m.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("async", TokenIdentifier)
		w.WriteString(" ")

		if v {
			m.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("*", TokenPunctuator)
		w.WriteString(" ")
	case MethodGetter:
		if v {
			m.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("get", TokenIdentifier)
		w.WriteString(" ")
	case MethodSetter:
		if v {
			m.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("set", TokenIdentifier)
		w.WriteString(" ")
	default:
		return
	}

	m.ClassElementName.printSource(w, v)
	m.Params.printSource(w, v)

	if v {
		m.Comments[2].printSource(w, true, false)
	}

	m.FunctionBody.printSource(w, v)

	if v {
		m.Comments[3].printSource(w, true, false)
	}
}

func (ce ClassElement) printSource(w writer, v bool) {
	w.Start(ce.Tokens)
	defer w.End()

	if v {
		ce.Comments[0].printSource(w, false, true)
	}

	if ce.Static {
		w.WriteStringWithType("static", TokenIdentifier)
		w.WriteString(" ")
	}

	if v {
		ce.Comments[1].printSource(w, true, false)
	}

	if ce.MethodDefinition != nil {
		ce.MethodDefinition.printSource(w, v)
	} else if ce.FieldDefinition != nil {
		ce.FieldDefinition.printSource(w, v)
		w.PrintSemiColon()
	} else if ce.ClassStaticBlock != nil {
		ce.ClassStaticBlock.printSource(w, v)
	} else if v {
		w.WriteStringWithType(";", TokenPunctuator)
	}

	if v {
		ce.Comments[2].printSource(w, false, false)
	}
}

func (fd FieldDefinition) printSource(w writer, v bool) {
	w.Start(fd.Tokens)
	defer w.End()

	fd.ClassElementName.printSource(w, v)

	if v {
		fd.Comments.printSource(w, false, false)
	}

	if fd.Initializer != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		fd.Initializer.printSource(w, v)
	}
}

func (cen ClassElementName) printSource(w writer, v bool) {
	w.Start(cen.Tokens)
	defer w.End()

	if v {
		cen.Comments[0].printSource(w, true, false)
	}

	if cen.PropertyName != nil {
		cen.PropertyName.printSource(w, v)
	} else if cen.PrivateIdentifier != nil {
		w.WriteToken(cen.PrivateIdentifier)
	}

	if v {
		cen.Comments[1].printSource(w, false, false)
	}
}

func (c ConditionalExpression) printSource(w writer, v bool) {
	w.Start(c.Tokens)
	defer w.End()

	if c.LogicalORExpression != nil {
		c.LogicalORExpression.printSource(w, v)
	} else if c.CoalesceExpression != nil {
		c.CoalesceExpression.printSource(w, v)
	}

	if c.True != nil && c.False != nil {
		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("?", TokenPunctuator)
		w.WriteString(" ")
		c.True.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(":", TokenPunctuator)
		w.WriteString(" ")
		c.False.printSource(w, v)
	}
}

func (a ArrowFunction) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if a.FunctionBody == nil && a.AssignmentExpression == nil || a.BindingIdentifier == nil && a.FormalParameters == nil {
		return
	}

	if v {
		a.Comments[0].printSource(w, false, true)
	}

	if a.Async {
		w.WriteStringWithType("async", TokenIdentifier)
		w.WriteString(" ")
	}

	if v {
		a.Comments[1].printSource(w, true, false)
	}

	if a.BindingIdentifier != nil {
		w.WriteToken(a.BindingIdentifier)
		w.WriteString(" ")
	} else if a.FormalParameters != nil {
		a.FormalParameters.printSource(w, v)
	}

	if v {
		a.Comments[2].printSource(w, true, false)
	}

	w.WriteStringWithType("=>", TokenPunctuator)
	w.WriteString(" ")

	if v {
		a.Comments[3].printSource(w, true, false)
	}

	if a.FunctionBody != nil {
		a.FunctionBody.printSource(w, v)
	} else {
		a.AssignmentExpression.printSource(w, v)
	}

	if v {
		a.Comments[4].printSource(w, false, true)
	}
}

func (n NewExpression) printSource(w writer, v bool) {
	w.Start(n.Tokens)
	defer w.End()

	for _, c := range n.News {
		if v {
			c.printSource(w, true, false)
		}

		w.WriteStringWithType("new", TokenKeyword)
		w.WriteString(" ")
	}

	n.MemberExpression.printSource(w, v)
}

func (c CallExpression) printSource(w writer, v bool) {
	w.Start(c.Tokens)
	defer w.End()

	if v {
		c.Comments[0].printSource(w, true, false)
	}

	if c.SuperCall && c.Arguments != nil {
		w.WriteStringWithType("super", TokenKeyword)

		if v {
			c.Comments[1].printSource(w, true, false)
		}

		c.Arguments.printSource(w, v)
	} else if c.ImportCall != nil {
		w.WriteStringWithType("import", TokenKeyword)

		if v {
			c.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("(", TokenPunctuator)

		ip := w

		if v {
			if c.Comments[2].hasSingleLineComment() || c.Comments[3].hasSingleLineComment() || c.ImportCall.hasSingleLineComment() {
				ip = w.Indent()

				c.Comments[2].printSource(w, false, true)
				ip.WriteString("\n")
			} else {
				c.Comments[2].printSource(w, true, false)
			}
		}

		c.ImportCall.printSource(ip, v)

		if v {
			if w != ip {
				w.WriteString("\n")
			}

			c.Comments[3].printSource(w, true, w != ip)
		}

		w.WriteStringWithType(")", TokenPunctuator)
	} else if c.MemberExpression != nil && c.Arguments != nil {
		c.MemberExpression.printSource(w, v)

		if v {
			c.Comments[1].printSource(w, true, false)
		}

		c.Arguments.printSource(w, v)
	} else if c.CallExpression != nil {
		if c.Arguments != nil {
			c.CallExpression.printSource(w, v)
			c.Arguments.printSource(w, v)
		} else if c.Expression != nil {
			c.CallExpression.printSource(w, v)
			w.WriteStringWithType("[", TokenPunctuator)

			ip := w

			if v {
				if c.Comments[2].hasSingleLineComment() || c.Comments[3].hasSingleLineComment() || c.ImportCall.hasSingleLineComment() {
					ip = w.Indent()

					c.Comments[2].printSource(w, false, true)
					ip.WriteString("\n")
				} else {
					c.Comments[2].printSource(w, true, false)
				}
			}

			c.Expression.printSource(ip, v)

			if v {
				if w != ip {
					w.WriteString("\n")
				}

				c.Comments[3].printSource(w, true, w != ip)
			}

			w.WriteStringWithType("]", TokenPunctuator)
		} else if c.IdentifierName != nil {
			c.CallExpression.printSource(w, v)

			if v && w.LastChar() != '\n' && len(c.CallExpression.Tokens) > 0 && c.IdentifierName.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				c.Comments[2].printSource(w, true, false)
			}

			w.WriteToken(c.IdentifierName)
		} else if c.TemplateLiteral != nil {
			c.CallExpression.printSource(w, v)
			c.TemplateLiteral.printSource(w, v)
		} else if c.PrivateIdentifier != nil {
			c.CallExpression.printSource(w, v)

			if v && w.LastChar() != '\n' && len(c.CallExpression.Tokens) > 0 && c.PrivateIdentifier.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				c.Comments[2].printSource(w, true, false)
			}

			w.WriteToken(c.PrivateIdentifier)
		}
	}

	if v {
		c.Comments[4].printSource(w, false, false)
	}
}

func (b BindingProperty) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	if !v && b.PropertyName.LiteralPropertyName != nil && b.BindingElement.SingleNameBinding != nil && b.PropertyName.LiteralPropertyName.Data == b.BindingElement.SingleNameBinding.Data {
		b.BindingElement.printSource(w, v)
	} else {
		if v {
			b.Comments[0].printSource(w, true, false)
		}

		b.PropertyName.printSource(w, v)

		if v {
			b.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType(":", TokenPunctuator)
		w.WriteString(" ")
		b.BindingElement.printSource(w, v)
	}
}

func (b BindingElement) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	if v {
		b.Comments[0].printSource(w, true, false)
	}

	if b.SingleNameBinding != nil {
		w.WriteToken(b.SingleNameBinding)
	} else if b.ArrayBindingPattern != nil {
		b.ArrayBindingPattern.printSource(w, v)
	} else if b.ObjectBindingPattern != nil {
		b.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if v && len(b.Comments[1]) > 0 {
		b.Comments[1].printSource(w, b.Initializer != nil, false)
	} else if b.Initializer != nil {
		w.WriteString(" ")
	}

	if b.Initializer != nil {
		w.WriteStringWithType("=", TokenPunctuator)
		w.WriteString(" ")
		b.Initializer.printSource(w, v)
	}
}

func (p PropertyName) printSource(w writer, v bool) {
	w.Start(p.Tokens)
	defer w.End()

	if p.LiteralPropertyName != nil {
		w.WriteStringWithType("", tokenColonSplit)
		w.WriteToken(p.LiteralPropertyName)
	} else if p.ComputedPropertyName != nil {
		w.WriteStringWithType("[", TokenPunctuator)

		ip := w

		if v {
			if hasSingleLineComment(p.Comments[:]) {
				ip = w.Indent()

				p.Comments[0].printSource(w, false, true)
				ip.WriteString("\n")
			} else {
				p.Comments[0].printSource(w, true, false)
			}
		}

		p.ComputedPropertyName.printSource(ip, v)

		if v {
			if len(p.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
				w.WriteString("\n")
			}

			p.Comments[1].printSource(w, false, w != ip)
		}

		w.WriteStringWithType("]", TokenPunctuator)
	}
}

func (l LogicalORExpression) printSource(w writer, v bool) {
	w.Start(l.Tokens)
	defer w.End()

	if l.LogicalORExpression != nil {
		l.LogicalORExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("||", TokenPunctuator)
		w.WriteString(" ")
	}

	l.LogicalANDExpression.printSource(w, v)
}

func (c ParenthesizedExpression) printSource(w writer, v bool) {
	w.Start(c.Tokens)
	defer w.End()

	w.WriteStringWithType("(", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(c.Comments[:]) || hasSingleLineComment(c.Expressions) {
			ip = w.Indent()
			sep = "\n"

			c.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			c.Comments[0].printSource(w, true, false)
		}
	}

	if len(c.Expressions) > 0 {
		c.Expressions[0].printSource(ip, v)

		for _, e := range c.Expressions[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			e.printSource(ip, v)
		}
	}

	if v {
		if len(c.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		c.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType(")", TokenPunctuator)
}

func (m MemberExpression) printSource(w writer, v bool) {
	w.Start(m.Tokens)
	defer w.End()

	if v {
		m.Comments[0].printSource(w, true, false)
	}

	if m.MemberExpression != nil {
		if m.Arguments != nil {
			w.WriteStringWithType("new", TokenKeyword)
			w.WriteString(" ")

			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v {
				m.Comments[2].printSource(w, true, false)
			}

			m.Arguments.printSource(w, v)
		} else if m.Expression != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			w.WriteStringWithType("[", TokenPunctuator)

			ip := w

			if v {
				if hasSingleLineComment(m.Comments[1:3]) || m.Expression.hasSingleLineComment() {
					ip = w.Indent()

					m.Comments[1].printSource(w, false, true)
					ip.WriteString("\n")
				} else {
					m.Comments[1].printSource(w, true, false)
				}
			}

			m.Expression.printSource(ip, v)

			if v {
				if len(m.Comments[2]) > 0 || w != ip && !w.LastIsWhitespace() {
					w.WriteString("\n")
				}

				m.Comments[2].printSource(w, false, w != ip)
			}

			w.WriteStringWithType("]", TokenPunctuator)
		} else if m.IdentifierName != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v && len(m.MemberExpression.Tokens) > 0 && m.IdentifierName.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteToken(m.IdentifierName)
		} else if m.PrivateIdentifier != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v && len(m.MemberExpression.Tokens) > 0 && m.PrivateIdentifier.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteToken(m.PrivateIdentifier)
		} else if m.TemplateLiteral != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			m.TemplateLiteral.printSource(w, v)
		}
	} else if m.PrimaryExpression != nil {
		m.PrimaryExpression.printSource(w, v)
	} else if m.SuperProperty {
		if m.Expression != nil {
			w.WriteStringWithType("super", TokenKeyword)

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteStringWithType("[", TokenPunctuator)

			ip := w

			if v {
				if hasSingleLineComment(m.Comments[2:4]) || m.Expression.hasSingleLineComment() {
					ip = w.Indent()

					m.Comments[2].printSource(w, false, true)
					ip.WriteString("\n")
				} else {
					m.Comments[2].printSource(w, true, false)
				}
			}

			m.Expression.printSource(ip, v)

			if v {
				if len(m.Comments[3]) > 0 || w != ip && !w.LastIsWhitespace() {
					w.WriteString("\n")
				}

				m.Comments[3].printSource(w, false, w != ip)
			}

			w.WriteStringWithType("]", TokenPunctuator)
		} else if m.IdentifierName != nil {
			w.WriteStringWithType("super", TokenKeyword)

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				m.Comments[2].printSource(w, true, false)
			}

			w.WriteToken(m.IdentifierName)
		}
	} else if m.NewTarget {
		w.WriteStringWithType("new", TokenKeyword)

		if v {
			m.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType(".", TokenPunctuator)

		if v {
			m.Comments[2].printSource(w, true, false)
		}

		w.WriteStringWithType("target", TokenIdentifier)
	} else if m.ImportMeta {
		w.WriteStringWithType("import", TokenKeyword)

		if v {
			m.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType(".", TokenPunctuator)

		if v {
			m.Comments[2].printSource(w, true, false)
		}

		w.WriteStringWithType("meta", TokenKeyword)
	}

	if v {
		m.Comments[4].printSource(w, false, false)
	}
}

func (a Argument) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if a.Spread {
		if v {
			a.Comments.printSource(w, true, false)
		}

		w.WriteStringWithType("...", TokenPunctuator)
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a Arguments) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	w.WriteStringWithType("(", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if a.hasSingleLineComment() {
			ip = w.Indent()
			sep = "\n"

			a.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			a.Comments[0].printSource(w, true, false)
		}
	}

	if len(a.ArgumentList) > 0 {
		a.ArgumentList[0].printSource(ip, v)

		for _, ae := range a.ArgumentList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			ae.printSource(ip, v)
		}
	}

	if v {
		if len(a.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		a.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType(")", TokenPunctuator)
}

func (t TemplateLiteral) printSource(w writer, v bool) {
	w.Start(t.Tokens)
	defer w.End()

	if t.NoSubstitutionTemplate != nil {
		w.WriteToken(t.NoSubstitutionTemplate)
	} else if t.TemplateHead != nil && t.TemplateTail != nil && len(t.Expressions) == len(t.TemplateMiddleList)+1 {
		w.WriteToken(t.TemplateHead)
		t.Expressions[0].printSource(w, v)

		for n, e := range t.Expressions[1:] {
			w.WriteToken(t.TemplateMiddleList[n])
			e.printSource(w, v)
		}

		w.WriteToken(t.TemplateTail)
	}
}

func (l LogicalANDExpression) printSource(w writer, v bool) {
	w.Start(l.Tokens)
	defer w.End()

	if l.LogicalANDExpression != nil {
		l.LogicalANDExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("&&", TokenPunctuator)
		w.WriteString(" ")
	}

	l.BitwiseORExpression.printSource(w, v)
}

func (p PrimaryExpression) printSource(w writer, v bool) {
	w.Start(p.Tokens)
	defer w.End()

	if p.This != nil {
		w.WriteStringWithType("this", TokenKeyword)
	} else if p.IdentifierReference != nil {
		w.WriteToken(p.IdentifierReference)
	} else if p.Literal != nil {
		w.WriteToken(p.Literal)
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
	} else if p.JSXElement != nil {
		p.JSXElement.printSource(w, v)
	} else if p.JSXFragment != nil {
		p.JSXFragment.printSource(w, v)
	}
}

func (b BitwiseORExpression) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	if b.BitwiseORExpression != nil {
		b.BitwiseORExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("|", TokenPunctuator)
		w.WriteString(" ")
	}

	b.BitwiseXORExpression.printSource(w, v)
}

func (a ArrayElement) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if v {
		a.Comments.printSource(w, true, false)
	}

	if a.Spread {
		w.WriteStringWithType("...", TokenPunctuator)
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a ArrayLiteral) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	w.WriteStringWithType("[", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(a.Comments[:]) || hasSingleLineComment(a.ElementList) {
			ip = w.Indent()
			sep = "\n"

			a.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			a.Comments[0].printSource(w, true, false)
		}
	}

	if len(a.ElementList) > 0 {
		a.ElementList[0].printSource(ip, v)

		for _, ae := range a.ElementList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			ae.printSource(ip, v)
		}
	}

	if v {
		if len(a.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		a.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("]", TokenPunctuator)
}

func (o ObjectLiteral) printSource(w writer, v bool) {
	w.Start(o.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(o.Comments[:]) || hasSingleLineComment(o.PropertyDefinitionList) {
			ip = w.Indent()
			sep = "\n"

			o.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			o.Comments[0].printSource(w, true, false)
		}
	}

	if len(o.PropertyDefinitionList) > 0 {
		o.PropertyDefinitionList[0].printSource(ip, v)

		for _, pd := range o.PropertyDefinitionList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			pd.printSource(ip, v)
		}
	}

	if v {
		if len(o.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		o.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("", tokenPossibleTrailingComma)
	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (b BitwiseXORExpression) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	if b.BitwiseXORExpression != nil {
		b.BitwiseXORExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("^", TokenPunctuator)
		w.WriteString(" ")
	}

	b.BitwiseANDExpression.printSource(w, v)
}

func (p PropertyDefinition) printSource(w writer, v bool) {
	w.Start(p.Tokens)
	defer w.End()

	if p.AssignmentExpression != nil {
		if p.PropertyName != nil {
			if v {
				p.Comments[0].printSource(w, true, false)
			}

			p.PropertyName.printSource(w, v)

			if v {
				p.Comments[1].printSource(w, true, false)
			}

			done := false

			if !v && !p.IsCoverInitializedName && p.PropertyName.LiteralPropertyName != nil && p.AssignmentExpression.ConditionalExpression != nil {
				c := UnwrapConditional(p.AssignmentExpression.ConditionalExpression)

				if pe, ok := c.(*PrimaryExpression); ok && pe.IdentifierReference != nil {
					done = pe.IdentifierReference.Type == p.PropertyName.LiteralPropertyName.Type && pe.IdentifierReference.Data == p.PropertyName.LiteralPropertyName.Data
				}
			}

			if !done {
				if p.IsCoverInitializedName {
					if !v || len(p.Comments[1]) == 0 {
						w.WriteString(" ")
					}

					w.WriteStringWithType("=", TokenPunctuator)
					w.WriteString(" ")
				} else {
					w.WriteStringWithType(":", TokenPunctuator)
					w.WriteString(" ")
				}

				p.AssignmentExpression.printSource(w, v)
			}
		} else {
			if v {
				p.Comments[0].printSource(w, true, false)
			}

			w.WriteStringWithType("...", TokenPunctuator)
			p.AssignmentExpression.printSource(w, v)
		}
	} else if p.MethodDefinition != nil {
		p.MethodDefinition.printSource(w, v)
	}
}

func (b BitwiseANDExpression) printSource(w writer, v bool) {
	w.Start(b.Tokens)
	defer w.End()

	if b.BitwiseANDExpression != nil {
		b.BitwiseANDExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("&", TokenPunctuator)
		w.WriteString(" ")
	}

	b.EqualityExpression.printSource(w, v)
}

func (e EqualityExpression) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	if e.EqualityExpression != nil {
		eo := e.EqualityOperator.String()
		if eo == "" || eo == unknown {
			return
		}

		e.EqualityExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(eo, TokenPunctuator)
		w.WriteString(" ")
	}

	e.RelationalExpression.printSource(w, v)
}

func (r RelationalExpression) printSource(w writer, v bool) {
	w.Start(r.Tokens)
	defer w.End()

	if r.PrivateIdentifier != nil {
		if v {
			r.Comments[0].printSource(w, true, false)
		}

		w.WriteToken(r.PrivateIdentifier)

		if v && len(r.Comments[1]) > 0 {
			r.Comments[1].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}

		w.WriteStringWithType("in", TokenKeyword)
		w.WriteString(" ")
	} else if r.RelationalExpression != nil {
		ro := r.RelationshipOperator.String()
		if ro == "" || ro == unknown {
			return
		}

		r.RelationalExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(ro, r.RelationshipOperator.tokenType())
		w.WriteString(" ")
	}

	r.ShiftExpression.printSource(w, v)
}

func (s ShiftExpression) printSource(w writer, v bool) {
	w.Start(s.Tokens)
	defer w.End()

	if s.ShiftExpression != nil {
		so := s.ShiftOperator.String()
		if so == "" || so == unknown {
			return
		}

		s.ShiftExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(so, TokenPunctuator)
		w.WriteString(" ")
	}

	s.AdditiveExpression.printSource(w, v)
}

func (a AdditiveExpression) printSource(w writer, v bool) {
	w.Start(a.Tokens)
	defer w.End()

	if a.AdditiveExpression != nil {
		ao := a.AdditiveOperator.String()
		if ao == "" || ao == unknown {
			return
		}
		a.AdditiveExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(ao, TokenPunctuator)
		w.WriteString(" ")
	}

	a.MultiplicativeExpression.printSource(w, v)
}

func (m MultiplicativeExpression) printSource(w writer, v bool) {
	w.Start(m.Tokens)
	defer w.End()

	if m.MultiplicativeExpression != nil {
		mo := m.MultiplicativeOperator.String()
		if mo == "" || mo == unknown {
			return
		}

		m.MultiplicativeExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType(mo, TokenPunctuator)
		w.WriteString(" ")
	}

	m.ExponentiationExpression.printSource(w, v)
}

func (e ExponentiationExpression) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	if e.ExponentiationExpression != nil {
		e.ExponentiationExpression.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("**", TokenPunctuator)
		w.WriteString(" ")
	}

	e.UnaryExpression.printSource(w, v)
}

func (u UnaryOperatorComments) printSource(w writer, v bool) {
	if v {
		u.Comments.printSource(w, true, false)
	}

	uo := u.String()
	if uo == "" || uo == unknown {
		return
	}

	w.WriteStringWithType(uo, u.tokenType())

	if len(uo) > 1 {
		w.WriteString(" ")
	}
}

func (u UnaryExpression) printSource(w writer, v bool) {
	w.Start(u.Tokens)
	defer w.End()

	for _, uo := range u.UnaryOperators {
		uo.printSource(w, v)
	}

	u.UpdateExpression.printSource(w, v)
}

func (u UpdateExpression) printSource(w writer, v bool) {
	w.Start(u.Tokens)
	defer w.End()

	if u.LeftHandSideExpression != nil {
		var uo string

		switch u.UpdateOperator {
		case UpdatePostIncrement:
			uo = "++"
		case UpdatePostDecrement:
			uo = "--"
		case UpdatePreIncrement, UpdatePreDecrement:
			return
		default:
		}

		u.LeftHandSideExpression.printSource(w, v)

		if len(uo) > 0 {
			w.WriteStringWithType(uo, TokenPunctuator)
		}

		if v {
			u.Comments.printSource(w, false, false)
		}
	} else if u.UnaryExpression != nil {
		if v {
			u.Comments.printSource(w, true, false)
		}

		switch u.UpdateOperator {
		case UpdatePreIncrement:
			w.WriteStringWithType("++", TokenPunctuator)
		case UpdatePreDecrement:
			w.WriteStringWithType("--", TokenPunctuator)
		default:
			return
		}

		u.UnaryExpression.printSource(w, v)
	}
}

func (m Module) printSource(w writer, v bool) {
	w.Start(m.Tokens)
	defer w.End()

	if v && len(m.Comments[0]) > 0 {
		m.Comments[0].printSource(w, false, true)
		w.WriteString("\n")
	}

	if len(m.ModuleListItems) > 0 {
		m.ModuleListItems[0].printSource(w, v)

		for _, mi := range m.ModuleListItems[1:] {
			w.WriteString("\n\n")
			mi.printSource(w, v)
		}
	}

	if v && len(m.Comments[1]) > 0 {
		w.WriteString("\n")
		m.Comments[1].printSource(w, false, false)
	}
}

func (m ModuleItem) printSource(w writer, v bool) {
	w.Start(m.Tokens)
	defer w.End()

	if m.ImportDeclaration != nil {
		m.ImportDeclaration.printSource(w, v)
	} else if m.ExportDeclaration != nil {
		m.ExportDeclaration.printSource(w, v)
	} else if m.StatementListItem != nil {
		m.StatementListItem.printSource(w, v)
	}
}

func (i ImportDeclaration) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	if v {
		i.Comments[0].printSource(w, true, false)
	}

	if i.ImportClause == nil && i.FromClause.ModuleSpecifier == nil {
		return
	}

	w.WriteStringWithType("import", TokenKeyword)
	w.WriteString(" ")

	if i.ImportClause != nil {
		i.ImportClause.printSource(w, v)
		i.FromClause.printSource(w, v)
	} else if i.FromClause.ModuleSpecifier != nil {
		if v {
			i.Comments[1].printSource(w, true, false)
		}

		w.WriteToken(i.FromClause.ModuleSpecifier)
	}

	if i.WithClause != nil {
		w.WriteString(" ")
		i.WithClause.printSource(w, v)
	}

	if v {
		i.Comments[2].printSource(w, false, false)
	}

	w.PrintSemiColon()

	if v {
		i.Comments[3].printSource(w, false, false)
	}
}

func (e ExportDeclaration) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	if e.FromClause != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		if e.ExportClause != nil {
			e.ExportClause.printSource(w, v)
		} else {
			w.WriteStringWithType("*", TokenPunctuator)

			if v && len(e.Comments[2]) > 0 {
				w.WriteString(" ")
				e.Comments[2].printSource(w, false, false)
			}

			if e.ExportFromClause != nil {
				if !w.LastIsWhitespace() {
					w.WriteString(" ")
				}

				w.WriteStringWithType("as", TokenIdentifier)
				w.WriteString(" ")

				if v && len(e.Comments[3]) > 0 {
					e.Comments[3].printSource(w, true, false)
				}

				w.WriteToken(e.ExportFromClause)
			}
		}

		if v {
			e.Comments[4].printSource(w, false, false)
		}

		e.FromClause.printSource(w, v)

		if v {
			e.Comments[5].printSource(w, false, false)
		}

		w.PrintSemiColon()
	} else if e.ExportClause != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		e.ExportClause.printSource(w, v)

		if v {
			e.Comments[5].printSource(w, false, false)
		}

		w.PrintSemiColon()
	} else if e.VariableStatement != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		e.VariableStatement.printSource(w, v)
	} else if e.Declaration != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		e.Declaration.printSource(w, v)
	} else if e.DefaultFunction != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("default", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		e.DefaultFunction.printSource(w, v)
	} else if e.DefaultClass != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("default", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		e.DefaultClass.printSource(w, v)
	} else if e.DefaultAssignmentExpression != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteStringWithType("export", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteStringWithType("default", TokenKeyword)
		w.WriteString(" ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		e.DefaultAssignmentExpression.printSource(w, v)

		if v {
			e.Comments[5].printSource(w, false, false)
		}

		w.PrintSemiColon()
	}

	if v {
		e.Comments[6].printSource(w, false, true)
	}
}

func (wc WithClause) printSource(w writer, v bool) {
	w.Start(wc.Tokens)
	defer w.End()

	if v {
		wc.Comments[0].printSource(w, true, false)
	}

	w.WriteStringWithType("with", TokenKeyword)
	w.WriteString(" ")

	if v {
		wc.Comments[1].printSource(w, true, false)
	}

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(wc.Comments[2:4]) || hasSingleLineComment(wc.WithEntries) {
			ip = w.Indent()
			sep = "\n"

			wc.Comments[2].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			wc.Comments[2].printSource(w, true, false)
		}
	}

	if len(wc.WithEntries) > 0 {
		wc.WithEntries[0].printSource(ip, v)

		for _, we := range wc.WithEntries[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			we.printSource(ip, v)
		}
	}

	if v {
		if len(wc.Comments[3]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		wc.Comments[3].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("", tokenPossibleTrailingComma)
	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (we WithEntry) printSource(w writer, v bool) {
	w.Start(we.Tokens)
	defer w.End()

	if we.AttributeKey == nil || we.Value == nil {
		return
	}

	if v {
		we.Comments[0].printSource(w, true, false)
	}

	w.WriteToken(we.AttributeKey)

	if v {
		we.Comments[1].printSource(w, false, false)
	}

	w.WriteStringWithType(":", TokenPunctuator)
	w.WriteString(" ")

	if v {
		we.Comments[2].printSource(w, true, false)
	}

	w.WriteToken(we.Value)

	if v {
		we.Comments[3].printSource(w, false, false)
	}
}

func (i ImportClause) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	if v {
		i.Comments[0].printSource(w, true, false)
	}

	if i.ImportedDefaultBinding != nil {
		w.WriteToken(i.ImportedDefaultBinding)

		if v {
			i.Comments[1].printSource(w, true, false)
		}

		if i.NameSpaceImport != nil || i.NamedImports != nil {
			w.WriteStringWithType(",", TokenPunctuator)
			w.WriteString(" ")

			if v {
				i.Comments[2].printSource(w, true, false)
			}
		}
	}

	if i.NameSpaceImport != nil {
		w.WriteStringWithType("*", TokenPunctuator)

		if v && len(i.Comments[3]) > 0 {
			i.Comments[3].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}

		w.WriteStringWithType("as", TokenIdentifier)
		w.WriteString(" ")

		if v {
			i.Comments[4].printSource(w, true, false)
		}

		w.WriteToken(i.NameSpaceImport)
	} else if i.NamedImports != nil {
		i.NamedImports.printSource(w, v)
	}

	if v {
		i.Comments[5].printSource(w, false, false)
	}
}

func (f FromClause) printSource(w writer, v bool) {
	w.Start(f.Tokens)
	defer w.End()

	if f.ModuleSpecifier == nil {
		return
	}

	if !w.LastIsWhitespace() {
		w.WriteString(" ")
	}

	w.WriteStringWithType("from", TokenIdentifier)
	w.WriteString(" ")

	if v {
		f.Comments.printSource(w, true, false)
	}

	w.WriteToken(f.ModuleSpecifier)
}

func (e ExportClause) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(e.Comments[:]) || hasSingleLineComment(e.ExportList) {
			ip = w.Indent()
			sep = "\n"

			e.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			e.Comments[0].printSource(w, true, false)
		}
	}

	if len(e.ExportList) > 0 {
		e.ExportList[0].printSource(ip, v)

		for _, es := range e.ExportList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			es.printSource(ip, v)
		}
	}

	if v {
		if len(e.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		e.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (n NamedImports) printSource(w writer, v bool) {
	w.Start(n.Tokens)
	defer w.End()

	w.WriteStringWithType("{", TokenPunctuator)

	ip := w
	sep := " "

	if v {
		if hasSingleLineComment(n.Comments[:]) || hasSingleLineComment(n.ImportList) {
			ip = w.Indent()
			sep = "\n"

			n.Comments[0].printSource(w, false, true)
			ip.WriteString("\n")
		} else {
			n.Comments[0].printSource(w, true, false)
		}
	}

	if len(n.ImportList) > 0 {
		n.ImportList[0].printSource(ip, v)

		for _, is := range n.ImportList[1:] {
			ip.WriteStringWithType(",", TokenPunctuator)
			ip.WriteString(sep)
			is.printSource(ip, v)
		}
	}

	if v {
		if len(n.Comments[1]) > 0 || w != ip && !w.LastIsWhitespace() {
			w.WriteString("\n")
		}

		n.Comments[1].printSource(w, false, w != ip)
	}

	w.WriteStringWithType("}", TokenRightBracePunctuator)
}

func (e ExportSpecifier) printSource(w writer, v bool) {
	w.Start(e.Tokens)
	defer w.End()

	if v {
		e.Comments[0].printSource(w, true, false)
	}

	if e.IdentifierName == nil {
		return
	}

	w.WriteStringWithType("", tokenColonSplit)
	w.WriteToken(e.IdentifierName)

	if e.EIdentifierName != nil && (e.EIdentifierName.Type != e.IdentifierName.Type || e.EIdentifierName.Data != e.IdentifierName.Data || v) {
		if v {
			e.Comments[1].printSource(w, false, false)
		}

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("as", TokenIdentifier)
		w.WriteString(" ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		w.WriteToken(e.EIdentifierName)
	}

	if v {
		e.Comments[3].printSource(w, false, false)
	}
}

func (i ImportSpecifier) printSource(w writer, v bool) {
	w.Start(i.Tokens)
	defer w.End()

	if i.ImportedBinding == nil {
		return
	}

	if v && len(i.Comments[0]) > 0 {
		i.Comments[0].printSource(w, true, false)
	}

	if i.IdentifierName != nil && (i.IdentifierName.Type != i.ImportedBinding.Type || i.IdentifierName.Data != i.ImportedBinding.Data || v) {
		w.WriteStringWithType("", tokenColonSplit)
		w.WriteToken(i.IdentifierName)

		if v {
			i.Comments[1].printSource(w, false, false)
		}

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("as", TokenIdentifier)
		w.WriteString(" ")

		if v {
			i.Comments[2].printSource(w, true, false)
		}
	}

	w.WriteToken(i.ImportedBinding)

	if v {
		i.Comments[3].printSource(w, false, false)
	}
}

func (oe OptionalExpression) printSource(w writer, v bool) {
	w.Start(oe.Tokens)
	defer w.End()

	if oe.MemberExpression != nil {
		oe.MemberExpression.printSource(w, v)
	} else if oe.CallExpression != nil {
		oe.CallExpression.printSource(w, v)
	} else if oe.OptionalExpression != nil {
		oe.OptionalExpression.printSource(w, v)
	}

	oe.OptionalChain.printSource(w, v)
}

func (oc OptionalChain) printSource(w writer, v bool) {
	w.Start(oc.Tokens)
	defer w.End()

	if oc.OptionalChain != nil {
		oc.OptionalChain.printSource(w, v)
	} else {
		w.WriteStringWithType("?.", TokenPunctuator)
	}

	if v {
		oc.Comments[0].printSource(w, true, false)
	}

	if oc.Arguments != nil {
		oc.Arguments.printSource(w, v)
	} else if oc.Expression != nil {
		w.WriteStringWithType("[", TokenPunctuator)

		ip := w

		if v {
			if oc.Comments[1].hasSingleLineComment() || oc.Comments[2].hasSingleLineComment() || oc.Expression.hasSingleLineComment() {
				ip = w.Indent()

				oc.Comments[1].printSource(w, false, true)
				ip.WriteString("\n")
			} else {
				oc.Comments[1].printSource(w, true, false)
			}
		}

		oc.Expression.printSource(ip, v)
		if v {
			if len(oc.Comments[2]) > 0 || w != ip && !w.LastIsWhitespace() {
				w.WriteString("\n")
			}

			oc.Comments[2].printSource(w, false, w != ip)
		}

		w.WriteStringWithType("]", TokenPunctuator)
	} else if oc.IdentifierName != nil {
		if oc.OptionalChain != nil {
			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				oc.Comments[1].printSource(w, true, false)
			}
		}

		w.WriteToken(oc.IdentifierName)
	} else if oc.TemplateLiteral != nil {
		oc.TemplateLiteral.printSource(w, v)
	} else if oc.PrivateIdentifier != nil {
		if oc.OptionalChain != nil {
			w.WriteStringWithType(".", TokenPunctuator)

			if v {
				oc.Comments[1].printSource(w, true, false)
			}
		}

		w.WriteToken(oc.PrivateIdentifier)
	}

	if v {
		oc.Comments[3].printSource(w, false, false)
	}
}

func (ce CoalesceExpression) printSource(w writer, v bool) {
	w.Start(ce.Tokens)
	defer w.End()

	if ce.CoalesceExpressionHead != nil {
		ce.CoalesceExpressionHead.printSource(w, v)

		if !w.LastIsWhitespace() {
			w.WriteString(" ")
		}

		w.WriteStringWithType("??", TokenPunctuator)
		w.WriteString(" ")
	}

	ce.BitwiseORExpression.printSource(w, v)
}

func (ja *JSXAttribute) printSource(w writer, v bool) {
	w.Start(ja.Tokens)
	defer w.End()

	if ja.Identifier != nil {
		if ja.Namespace != nil {
			w.WriteToken(ja.Namespace)
			w.WriteStringWithType(":", TokenPunctuator)
		}

		w.WriteToken(ja.Identifier)

		if ja.JSXString != nil {
			w.WriteStringWithType("=", TokenPunctuator)
			w.WriteToken(ja.JSXString)
		} else if ja.JSXElement != nil {
			w.WriteStringWithType("=", TokenPunctuator)
			ja.JSXElement.printSource(w, v)
		} else if ja.JSXFragment != nil {
			w.WriteStringWithType("=", TokenPunctuator)
			ja.JSXFragment.printSource(w, v)
		} else if ja.AssignmentExpression != nil {
			w.WriteStringWithType("=", TokenPunctuator)
			w.WriteStringWithType("{", TokenPunctuator)

			ip := w

			if v && (ja.Comments.hasSingleLineComment() || ja.AssignmentExpression.hasSingleLineComment()) {
				ip = w.Indent()

				ip.WriteString("\n")
				ja.Comments.printSource(ip, false, true)
			}

			ja.AssignmentExpression.printSource(ip, v)
			w.WriteStringWithType("}", TokenRightBracePunctuator)
		}
	} else if ja.AssignmentExpression != nil {
		w.WriteStringWithType("{", TokenPunctuator)

		ip := w

		if v && (ja.Comments.hasSingleLineComment() || ja.AssignmentExpression.hasSingleLineComment()) {
			ip = w.Indent()

			if ja.Comments.hasSingleLineComment() {
				ip.WriteString("\n")
				ja.Comments.printSource(ip, false, true)
			}
		}

		ip.WriteStringWithType("...", TokenPunctuator)
		ja.AssignmentExpression.printSource(ip, v)

		if v && w != ip && w.LastChar() != '\n' {
			w.WriteString("\n")
		}

		w.WriteStringWithType("}", TokenRightBracePunctuator)
	}
}

func (jc *JSXChild) printSource(w writer, v bool) {
	w.Start(jc.Tokens)
	defer w.End()

	if jc.JSXText != nil {
		w.WriteToken(jc.JSXText)
	} else if jc.JSXElement != nil {
		jc.JSXElement.printSource(w, v)
	} else if jc.JSXFragment != nil {
		jc.JSXFragment.printSource(w, v)
	} else {
		ip := w

		w.WriteStringWithType("{", TokenPunctuator)

		if v && (jc.Comments.hasSingleLineComment() || jc.JSXChildExpression.hasSingleLineComment()) {
			ip = w.Indent()

			ip.WriteString("\n")
			jc.Comments.printSource(ip, false, true)
		}

		if jc.JSXChildExpression != nil {
			if jc.Spread {
				ip.WriteStringWithType("...", TokenPunctuator)
			}

			jc.JSXChildExpression.printSource(ip, v)
		}

		if v && w != ip && w.LastChar() != '\n' {
			w.WriteString("\n")
		}

		w.WriteStringWithType("}", TokenRightBracePunctuator)
	}
}

func (je *JSXElement) printSource(w writer, v bool) {
	w.Start(je.Tokens)
	defer w.End()

	w.WriteStringWithType("<", TokenJSXElementStart)
	je.ElementName.printSource(w, v)

	for _, attr := range je.Attributes {
		w.WriteString(" ")
		attr.printSource(w, v)
	}

	if je.SelfClosing {
		w.WriteString(" ")
		w.WriteStringWithType("/", TokenPunctuator)
		w.WriteStringWithType(">", TokenJSXElementEnd)
	} else {
		w.WriteStringWithType(">", TokenJSXElementEnd)

		if len(je.Children) > 1 {
			ip := w.Indent()

			for _, child := range je.Children {
				ip.WriteString("\n")
				child.printSource(ip, v)
			}

			w.WriteString("\n")
		} else if len(je.Children) == 1 {
			je.Children[0].printSource(w, v)
		}

		w.WriteStringWithType("<", TokenJSXElementStart)
		w.WriteStringWithType("/", TokenPunctuator)
		je.ElementName.printSource(w, v)
		w.WriteStringWithType(">", TokenJSXElementEnd)
	}
}

func (jn *JSXElementName) printSource(w writer, v bool) {
	if jn.Identifier != nil {
		if jn.Namespace != nil {
			w.WriteStringWithType(jn.Namespace.Data, jn.Namespace.Type)
			w.WriteStringWithType(":", TokenPunctuator)
			w.WriteStringWithType(jn.Identifier.Data, jn.Identifier.Type)
		} else {
			w.WriteStringWithType(jn.Identifier.Data, jn.Identifier.Type)

			for _, m := range jn.MemberExpression {
				w.WriteStringWithType(".", TokenPunctuator)
				w.WriteStringWithType(m.Data, m.Type)
			}
		}
	}
}

func (jf *JSXFragment) printSource(w writer, v bool) {
	w.Start(jf.Tokens)
	defer w.End()

	w.WriteStringWithType("<", TokenJSXElementStart)
	w.WriteStringWithType(">", TokenJSXElementEnd)

	if len(jf.Children) > 1 {
		ip := w.Indent()

		for _, child := range jf.Children {
			ip.WriteString("\n")
			child.printSource(ip, v)
		}

		w.WriteString("\n")
	} else if len(jf.Children) == 1 {
		jf.Children[0].printSource(w, v)
	}

	w.WriteStringWithType("<", TokenJSXElementStart)
	w.WriteStringWithType("/", TokenPunctuator)
	w.WriteStringWithType(">", TokenJSXElementEnd)
}
