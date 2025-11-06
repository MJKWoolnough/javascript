package javascript

func (s Script) printSource(w writer, v bool) {
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
	if v {
		s.Comments[0].printSource(w, true, false)
	}

	if s.Statement != nil {
		s.Statement.printSource(w, v)
	} else if s.Declaration != nil {
		s.Declaration.printSource(w, v)
	}

	if v {
		s.Comments[1].printSource(w, true, false)
	}
}

func (s Statement) printSource(w writer, v bool) {
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
			w.WriteString(s.LabelIdentifier.Data)
			w.WriteString(": ")
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
			w.WriteString("continue;")
		} else {
			w.WriteString("continue ")
			w.WriteString(s.LabelIdentifier.Data)
			w.PrintSemiColon()
		}
	case StatementBreak:
		if s.LabelIdentifier == nil {
			w.WriteString("break;")
		} else {
			w.WriteString("break ")
			w.WriteString(s.LabelIdentifier.Data)
			w.PrintSemiColon()
		}
	case StatementReturn:
		if s.ExpressionStatement == nil {
			w.WriteString("return;")
		} else {
			w.WriteString("return ")
			s.ExpressionStatement.printSource(w, v)
			w.PrintSemiColon()
		}
	case StatementThrow:
		if s.ExpressionStatement != nil {
			w.WriteString("throw ")
			s.ExpressionStatement.printSource(w, v)
			w.PrintSemiColon()
		}
	case StatementDebugger:
		w.WriteString("debugger;")
	}
}

func (d Declaration) printSource(w writer, v bool) {
	if d.ClassDeclaration != nil {
		d.ClassDeclaration.printSource(w, v)
	} else if d.FunctionDeclaration != nil {
		d.FunctionDeclaration.printSource(w, v)
	} else if d.LexicalDeclaration != nil {
		d.LexicalDeclaration.printSource(w, v)
	}
}

func (b Block) printSource(w writer, v bool) {
	w.WriteString("{")

	var lastLine uint64

	if v && len(b.Tokens) > 0 {
		lastLine = b.Tokens[0].Line
	}

	pp := w.Indent()

	for _, stmt := range b.StatementList {
		if v {
			if len(stmt.Tokens) > 0 {
				ll := stmt.Tokens[0].Line

				if ll > lastLine {
					pp.WriteString("\n")
				} else {
					pp.WriteString(" ")
				}

				lastLine = ll
			} else {
				pp.WriteString("\n")
			}
		} else {
			pp.WriteString("\n")
		}

		stmt.printSource(pp, v)
	}

	if len(b.StatementList) > 0 {
		if v && len(b.Tokens) > 0 {
			if b.Tokens[len(b.Tokens)-1].Line > lastLine {
				w.WriteString("\n")
			} else {
				pp.WriteString(" ")
			}
		} else {
			w.WriteString("\n")
		}
	}

	w.WriteString("}")
}

func (vs VariableStatement) printSource(w writer, v bool) {
	if len(vs.VariableDeclarationList) == 0 {
		return
	}

	w.WriteString("var ")

	var lastLine uint64

	if v && len(vs.Tokens) > 0 {
		lastLine = vs.Tokens[0].Line
	}

	for n, vd := range vs.VariableDeclarationList {
		if n > 0 {
			if v && len(vd.Tokens) > 0 {
				if ll := vd.Tokens[0].Line; ll > lastLine {
					lastLine = ll

					w.WriteString(",\n")
				} else {
					w.WriteString(", ")
				}
			} else {
				w.WriteString(", ")
			}
		}

		vd.printSource(w, v)
	}

	w.PrintSemiColon()
}

func (e Expression) printSource(w writer, v bool) {
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
				w.WriteString(",\n")
			} else {
				w.WriteString(", ")
			}
		} else {
			w.WriteString(", ")
		}

		ae.printSource(w, v)
	}
}

func (i IfStatement) printSource(w writer, v bool) {
	w.WriteString("if (")

	if v {
		pp := w.Indent()

		var nl bool

		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			nl = true

			pp.WriteString("\n")
		}

		i.Expression.printSource(pp, true)

		if nl {
			w.WriteString("\n")
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.WriteString(") ")
	i.Statement.printSource(w, v)

	if i.ElseStatement != nil {
		w.WriteString(" else ")
		i.ElseStatement.printSource(w, v)
	}
}

func (i IterationStatementDo) printSource(w writer, v bool) {
	w.WriteString("do ")
	i.Statement.printSource(w, v)
	w.WriteString(" while (")

	if v {
		pp := w.Indent()

		var nl bool

		if len(i.Expression.Tokens) > 0 && len(i.Tokens) > 0 && i.Expression.Tokens[0].Line < i.Tokens[len(i.Tokens)-1].Line {
			nl = true

			pp.WriteString("\n")
		}

		i.Expression.printSource(pp, true)

		if nl {
			w.WriteString("\n")
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.WriteString(");")
}

func (i IterationStatementWhile) printSource(w writer, v bool) {
	w.WriteString("while (")

	if v {
		pp := w.Indent()

		var nl bool

		if len(i.Tokens) > 0 && len(i.Expression.Tokens) > 0 && i.Expression.Tokens[0].Line > i.Tokens[0].Line {
			pp.WriteString("\n")

			nl = true
		}

		i.Expression.printSource(pp, true)

		if nl {
			w.WriteString("\n")
		}
	} else {
		i.Expression.printSource(w, false)
	}

	w.WriteString(") ")
	i.Statement.printSource(w, v)
}

func (i IterationStatementFor) printSource(w writer, v bool) {
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
		w.WriteString("for await (")
	default:
		w.WriteString("for (")
	}

	pp := w.Indent()

	var lastLine uint64

	if v && len(i.Tokens) > 0 {
		lastLine = i.Tokens[0].Line
	}

	var endline bool

	switch i.Type {
	case ForNormal:
		w.WriteString(";")
	case ForNormalVar:
		if v && len(i.InitVar[0].Tokens) > 0 {
			if i.InitVar[0].Tokens[0].Line > lastLine {
				pp.WriteString("\n")
			}

			lastLine = i.InitVar[0].Tokens[len(i.InitVar[0].Tokens)-1].Line
		}

		pp.WriteString("var ")
		LexicalBinding(i.InitVar[0]).printSource(pp, v)

		for _, vd := range i.InitVar[1:] {
			if v && len(vd.Tokens) > 0 {
				if vd.Tokens[0].Line > lastLine {
					pp.WriteString(",\n")
				} else {
					pp.WriteString(", ")
				}
			} else {
				pp.WriteString(", ")
			}

			LexicalBinding(vd).printSource(pp, v)
		}

		w.WriteString(";")
	case ForNormalLexicalDeclaration:
		if v && len(i.InitLexical.Tokens) > 0 {
			if i.InitLexical.Tokens[0].Line > lastLine {
				endline = true

				pp.WriteString("\n")
			}

			lastLine = i.InitLexical.Tokens[len(i.InitLexical.Tokens)-1].Line
		}

		i.InitLexical.printSource(pp, v)
	case ForNormalExpression:
		if v && len(i.InitExpression.Tokens) > 0 {
			if i.InitExpression.Tokens[0].Line > lastLine {
				endline = true

				pp.WriteString("\n")
			}

			lastLine = i.InitExpression.Tokens[len(i.InitExpression.Tokens)-1].Line
		}

		i.InitExpression.printSource(pp, v)
		w.WriteString(";")
	case ForInLeftHandSide, ForOfLeftHandSide, ForAwaitOfLeftHandSide:
		if v {
			if len(i.LeftHandSideExpression.Tokens) > 0 {
				if i.LeftHandSideExpression.Tokens[0].Line > lastLine {
					endline = true

					pp.WriteString("\n")
				}

				lastLine = i.LeftHandSideExpression.Tokens[len(i.LeftHandSideExpression.Tokens)-1].Line
			}
		}

		i.LeftHandSideExpression.printSource(pp, v)
	default:
		switch i.Type {
		case ForInVar, ForOfVar, ForAwaitOfVar:
			w.WriteString("var ")
		case ForInLet, ForOfLet, ForAwaitOfLet:
			w.WriteString("let ")
		case ForInConst, ForOfConst, ForAwaitOfConst:
			w.WriteString("const ")
		}

		if i.ForBindingIdentifier != nil {
			w.WriteString(i.ForBindingIdentifier.Data)
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

					pp.WriteString("\n")
				} else {
					w.WriteString(" ")
				}

				lastLine = i.Conditional.Tokens[len(i.Conditional.Tokens)-1].Line
			} else {
				w.WriteString(" ")
			}

			i.Conditional.printSource(pp, v)
		}

		w.WriteString(";")

		if i.Afterthought != nil {
			if v && len(i.Afterthought.Tokens) > 0 {
				if i.Afterthought.Tokens[0].Line > lastLine {
					endline = true

					pp.WriteString("\n")
				} else {
					w.WriteString(" ")
				}
			} else {
				w.WriteString(" ")
			}

			i.Afterthought.printSource(pp, v)
		}
	case ForInLeftHandSide, ForInVar, ForInLet, ForInConst:
		w.WriteString(" in ")
		i.In.printSource(pp, v)
	case ForOfLeftHandSide, ForOfVar, ForOfLet, ForOfConst, ForAwaitOfLeftHandSide, ForAwaitOfVar, ForAwaitOfLet, ForAwaitOfConst:
		w.WriteString(" of ")
		i.Of.printSource(pp, v)
	}

	if endline {
		w.WriteString("\n")
	}

	w.WriteString(") ")
	i.Statement.printSource(w, v)
}

func (s SwitchStatement) printSource(w writer, v bool) {
	w.WriteString("switch (")

	if v {
		pp := w.Indent()

		var nl bool

		if len(s.Tokens) > 0 && len(s.Expression.Tokens) > 0 && s.Expression.Tokens[0].Line > s.Tokens[0].Line {
			nl = true

			pp.WriteString("\n")
		}

		s.Expression.printSource(pp, true)

		if nl {
			w.WriteString("\n")
		}
	} else {
		s.Expression.printSource(w, false)
	}

	w.WriteString(") {")

	if len(s.CaseClauses) > 0 || s.DefaultClause != nil || len(s.PostDefaultCaseClauses) > 0 {
		w.WriteString("\n")
	}

	for _, c := range s.CaseClauses {
		c.printSource(w, v)
		w.WriteString("\n")
	}

	if s.DefaultClause != nil {
		w.WriteString("default:")

		pp := w.Indent()

		for _, stmt := range s.DefaultClause {
			pp.WriteString("\n")
			stmt.printSource(pp, v)
		}

		w.WriteString("\n")
	}

	for _, c := range s.PostDefaultCaseClauses {
		c.printSource(w, v)
		w.WriteString("\n")
	}

	w.WriteString("}")
}

func (ws WithStatement) printSource(w writer, v bool) {
	w.WriteString("with (")

	if v {
		pp := w.Indent()

		var nl bool

		if len(ws.Tokens) > 0 && len(ws.Expression.Tokens) > 0 && ws.Expression.Tokens[0].Line > ws.Tokens[0].Line {
			nl = true

			pp.WriteString("\n")
		}

		ws.Expression.printSource(pp, true)

		if nl {
			w.WriteString("\n")
		}
	} else {
		ws.Expression.printSource(w, false)
	}

	w.WriteString(") ")
	ws.Statement.printSource(w, v)
}

func (f FunctionDeclaration) printSource(w writer, v bool) {
	switch f.Type {
	case FunctionNormal:
		w.WriteString("function ")
	case FunctionGenerator:
		w.WriteString("function* ")
	case FunctionAsync:
		w.WriteString("async function ")
	case FunctionAsyncGenerator:
		w.WriteString("async function* ")
	default:
		return
	}

	if f.BindingIdentifier != nil {
		w.WriteString(f.BindingIdentifier.Data)
	}

	f.FormalParameters.printSource(w.Indent(), v)
	f.FunctionBody.printSource(w, v)
}

func (t TryStatement) printSource(w writer, v bool) {
	w.WriteString("try ")
	t.TryBlock.printSource(w, v)

	if t.CatchBlock != nil {
		if t.CatchParameterBindingIdentifier != nil {
			w.WriteString(" catch (")
			w.WriteString(t.CatchParameterBindingIdentifier.Data)
			w.WriteString(") ")
		} else if t.CatchParameterArrayBindingPattern != nil {
			w.WriteString(" catch (")
			t.CatchParameterArrayBindingPattern.printSource(w, v)
			w.WriteString(") ")
		} else if t.CatchParameterObjectBindingPattern != nil {
			w.WriteString(" catch (")
			t.CatchParameterObjectBindingPattern.printSource(w, v)
			w.WriteString(") ")
		} else {
			w.WriteString(" catch ")
		}

		t.CatchBlock.printSource(w, v)
	}

	if t.FinallyBlock != nil {
		w.WriteString(" finally ")
		t.FinallyBlock.printSource(w, v)
	}
}

func (c ClassDeclaration) printSource(w writer, v bool) {
	w.WriteString("class ")

	if c.BindingIdentifier != nil {
		w.WriteString(c.BindingIdentifier.Data)
		w.WriteString(" ")
	}

	if c.ClassHeritage != nil {
		w.WriteString("extends ")
		c.ClassHeritage.printSource(w, v)
		w.WriteString(" ")
	}

	w.WriteString("{")

	if len(c.ClassBody) > 0 {
		pp := w.Indent()

		for _, ce := range c.ClassBody {
			pp.WriteString("\n")
			ce.printSource(pp, v)
		}

		w.WriteString("\n")
	}

	w.WriteString("}")
}

func (l LexicalDeclaration) printSource(w writer, v bool) {
	if len(l.BindingList) == 0 {
		return
	}

	if l.LetOrConst == Let {
		w.WriteString("let ")
	} else if l.LetOrConst == Const {
		w.WriteString("const ")
	}

	l.BindingList[0].printSource(w, v)

	for _, lb := range l.BindingList[1:] {
		if v {
			w.WriteString(",\n")
		} else {
			w.WriteString(", ")
		}

		lb.printSource(w, v)
	}

	w.PrintSemiColon()
}

func (l LexicalBinding) printSource(w writer, v bool) {
	if l.BindingIdentifier != nil {
		w.WriteString(l.BindingIdentifier.Data)
	} else if l.ArrayBindingPattern != nil {
		l.ArrayBindingPattern.printSource(w, v)
	} else if l.ObjectBindingPattern != nil {
		l.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if l.Initializer != nil {
		w.WriteString(" = ")
		l.Initializer.printSource(w, v)
	}
}

func (a AssignmentExpression) printSource(w writer, v bool) {
	if a.Yield && a.AssignmentExpression != nil {
		w.WriteString("yield ")

		if a.Delegate {
			w.WriteString("* ")
		}

		a.AssignmentExpression.printSource(w, v)
	} else if a.ArrowFunction != nil {
		a.ArrowFunction.printSource(w, v)
	} else if a.LeftHandSideExpression != nil && a.AssignmentExpression != nil {
		ao := " = "

		switch a.AssignmentOperator {
		case AssignmentAssign:
		case AssignmentMultiply:
			ao = " *= "
		case AssignmentDivide:
			ao = " /= "
		case AssignmentRemainder:
			ao = " %= "
		case AssignmentAdd:
			ao = " += "
		case AssignmentSubtract:
			ao = " -= "
		case AssignmentLeftShift:
			ao = " <<= "
		case AssignmentSignPropagatingRightShift:
			ao = " >>= "
		case AssignmentZeroFillRightShift:
			ao = " >>>= "
		case AssignmentBitwiseAND:
			ao = " &= "
		case AssignmentBitwiseXOR:
			ao = " ^= "
		case AssignmentBitwiseOR:
			ao = " |= "
		case AssignmentExponentiation:
			ao = " **= "
		case AssignmentLogicalAnd:
			ao = " &&= "
		case AssignmentLogicalOr:
			ao = " ||= "
		case AssignmentNullish:
			ao = " ??= "
		default:
			return
		}

		a.LeftHandSideExpression.printSource(w, v)
		w.WriteString(ao)
		a.AssignmentExpression.printSource(w, v)
	} else if a.AssignmentPattern != nil && a.AssignmentExpression != nil && a.AssignmentOperator == AssignmentAssign {
		a.AssignmentPattern.printSource(w, v)
		w.WriteString(" = ")
		a.AssignmentExpression.printSource(w, v)
	} else if a.ConditionalExpression != nil {
		a.ConditionalExpression.printSource(w, v)
	}
}

func (l LeftHandSideExpression) printSource(w writer, v bool) {
	if l.NewExpression != nil {
		l.NewExpression.printSource(w, v)
	} else if l.CallExpression != nil {
		l.CallExpression.printSource(w, v)
	} else if l.OptionalExpression != nil {
		l.OptionalExpression.printSource(w, v)
	}
}

func (a AssignmentPattern) printSource(w writer, v bool) {
	if a.ArrayAssignmentPattern != nil {
		a.ArrayAssignmentPattern.printSource(w, v)
	} else if a.ObjectAssignmentPattern != nil {
		a.ObjectAssignmentPattern.printSource(w, v)
	}
}

func (a ArrayAssignmentPattern) printSource(w writer, v bool) {
	w.WriteString("[")

	for n, ae := range a.AssignmentElements {
		if n > 0 {
			w.WriteString(", ")
		}

		ae.printSource(w, v)
	}

	if a.AssignmentRestElement != nil {
		if len(a.AssignmentElements) > 0 {
			w.WriteString(", ")
		}

		w.WriteString("...")
		a.AssignmentRestElement.printSource(w, v)
	}

	w.WriteString("]")
}

func (o ObjectAssignmentPattern) printSource(w writer, v bool) {
	w.WriteString("{")

	for n, ap := range o.AssignmentPropertyList {
		if n > 0 {
			w.WriteString(", ")
		}

		ap.printSource(w, v)
	}

	if o.AssignmentRestElement != nil {
		if len(o.AssignmentPropertyList) > 0 {
			w.WriteString(", ")
		}

		w.WriteString("...")
		o.AssignmentRestElement.printSource(w, v)
	}

	w.WriteString("}")
}

func (a AssignmentElement) printSource(w writer, v bool) {
	a.DestructuringAssignmentTarget.printSource(w, v)

	if a.Initializer != nil {
		w.WriteString(" = ")
		a.Initializer.printSource(w, v)
	}
}

func (a AssignmentProperty) printSource(w writer, v bool) {
	a.PropertyName.printSource(w, v)

	if a.DestructuringAssignmentTarget != nil {
		if !v && a.PropertyName.LiteralPropertyName != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.CallExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.OptionalExpression == nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News == 0 && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil && a.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference.Data == a.PropertyName.LiteralPropertyName.Data {
			return
		}

		w.WriteString(": ")
		a.DestructuringAssignmentTarget.printSource(w, v)
	}

	if a.Initializer != nil {
		w.WriteString(" = ")
		a.Initializer.printSource(w, v)
	}
}

func (d DestructuringAssignmentTarget) printSource(w writer, v bool) {
	if d.LeftHandSideExpression != nil {
		d.LeftHandSideExpression.printSource(w, v)
	} else if d.AssignmentPattern != nil {
		d.AssignmentPattern.printSource(w, v)
	}
}

func (o ObjectBindingPattern) printSource(w writer, v bool) {
	w.WriteString("{")

	for n, bp := range o.BindingPropertyList {
		if n > 0 {
			w.WriteString(", ")
		}

		bp.printSource(w, v)
	}

	if o.BindingRestProperty != nil {
		if len(o.BindingPropertyList) > 0 {
			w.WriteString(", ")
		}

		w.WriteString("...")
		w.WriteString(o.BindingRestProperty.Data)
	}

	w.WriteString("}")
}

func (a ArrayBindingPattern) printSource(w writer, v bool) {
	w.WriteString("[")

	for n, be := range a.BindingElementList {
		if n > 0 {
			w.WriteString(", ")
		}

		be.printSource(w, v)
	}

	if a.BindingRestElement != nil {
		if len(a.BindingElementList) > 0 {
			w.WriteString(", ")
		}

		w.WriteString("...")
		a.BindingRestElement.printSource(w, v)
	}

	w.WriteString("]")
}

func (c CaseClause) printSource(w writer, v bool) {
	w.WriteString("case ")
	c.Expression.printSource(w, v)
	w.WriteString(":")

	pp := w.Indent()

	for _, stmt := range c.StatementList {
		pp.WriteString("\n")
		stmt.printSource(pp, v)
	}
}

func (f FormalParameters) printSource(w writer, v bool) {
	w.WriteString("(")

	if len(f.FormalParameterList) > 0 {
		f.FormalParameterList[0].printSource(w, v)

		for _, be := range f.FormalParameterList[1:] {
			w.WriteString(", ")
			be.printSource(w, v)
		}

		if f.BindingIdentifier != nil || f.ArrayBindingPattern != nil || f.ObjectBindingPattern != nil {
			w.WriteString(", ")
		}
	}

	if f.BindingIdentifier != nil {
		w.WriteString("...")
		w.WriteString(f.BindingIdentifier.Data)
	} else if f.ArrayBindingPattern != nil {
		w.WriteString("...")
		f.ArrayBindingPattern.printSource(w, v)
	} else if f.ObjectBindingPattern != nil {
		w.WriteString("...")
		f.ObjectBindingPattern.printSource(w, v)
	}

	w.WriteString(") ")
}

func (m MethodDefinition) printSource(w writer, v bool) {
	switch m.Type {
	case MethodNormal:
	case MethodGenerator:
		w.WriteString("* ")
	case MethodAsync:
		w.WriteString("async ")
	case MethodAsyncGenerator:
		w.WriteString("async * ")
	case MethodGetter:
		w.WriteString("get ")
	case MethodSetter:
		w.WriteString("set ")
	default:
		return
	}

	m.ClassElementName.printSource(w, v)
	m.Params.printSource(w, v)
	m.FunctionBody.printSource(w, v)
}

func (ce ClassElement) printSource(w writer, v bool) {
	if ce.Static {
		w.WriteString("static ")
	}

	if ce.MethodDefinition != nil {
		ce.MethodDefinition.printSource(w, v)
	} else if ce.FieldDefinition != nil {
		ce.FieldDefinition.printSource(w, v)
	} else if ce.ClassStaticBlock != nil {
		ce.ClassStaticBlock.printSource(w, v)
	}
}

func (fd FieldDefinition) printSource(w writer, v bool) {
	fd.ClassElementName.printSource(w, v)

	if fd.Initializer != nil {
		w.WriteString(" = ")
		fd.Initializer.printSource(w, v)
	}

	w.PrintSemiColon()
}

func (cen ClassElementName) printSource(w writer, v bool) {
	if cen.PropertyName != nil {
		cen.PropertyName.printSource(w, v)
	} else if cen.PrivateIdentifier != nil {
		w.WriteString(cen.PrivateIdentifier.Data)
	}
}

func (c ConditionalExpression) printSource(w writer, v bool) {
	if c.LogicalORExpression != nil {
		c.LogicalORExpression.printSource(w, v)
	} else if c.CoalesceExpression != nil {
		c.CoalesceExpression.printSource(w, v)
	}

	if c.True != nil && c.False != nil {
		w.WriteString(" ? ")
		c.True.printSource(w, v)
		w.WriteString(" : ")
		c.False.printSource(w, v)
	}
}

func (a ArrowFunction) printSource(w writer, v bool) {
	if a.FunctionBody == nil && a.AssignmentExpression == nil || (a.BindingIdentifier == nil && a.FormalParameters == nil) {
		return
	}

	if a.Async {
		w.WriteString("async ")
	}

	if a.BindingIdentifier != nil {
		w.WriteString(a.BindingIdentifier.Data)
		w.WriteString(" ")
	} else if a.FormalParameters != nil {
		a.FormalParameters.printSource(w, v)
	}

	w.WriteString("=> ")

	if a.FunctionBody != nil {
		a.FunctionBody.printSource(w, v)
	} else {
		a.AssignmentExpression.printSource(w, v)
	}
}

func (n NewExpression) printSource(w writer, v bool) {
	for i := uint(0); i < n.News; i++ {
		w.WriteString("new ")
	}

	n.MemberExpression.printSource(w, v)
}

func (c CallExpression) printSource(w writer, v bool) {
	if c.SuperCall && c.Arguments != nil {
		w.WriteString("super")
		c.Arguments.printSource(w, v)
	} else if c.ImportCall != nil {
		w.WriteString("import(")
		c.ImportCall.printSource(w, v)
		w.WriteString(")")
	} else if c.MemberExpression != nil && c.Arguments != nil {
		c.MemberExpression.printSource(w, v)
		c.Arguments.printSource(w, v)
	} else if c.CallExpression != nil {
		if c.Arguments != nil {
			c.CallExpression.printSource(w, v)
			c.Arguments.printSource(w, v)
		} else if c.Expression != nil {
			c.CallExpression.printSource(w, v)
			w.WriteString("[")
			c.Expression.printSource(w, v)
			w.WriteString("]")
		} else if c.IdentifierName != nil {
			c.CallExpression.printSource(w, v)

			if v && len(c.CallExpression.Tokens) > 0 && c.IdentifierName.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteString(".")
			w.WriteString(c.IdentifierName.Data)
		} else if c.TemplateLiteral != nil {
			c.CallExpression.printSource(w, v)
			c.TemplateLiteral.printSource(w, v)
		} else if c.PrivateIdentifier != nil {
			c.CallExpression.printSource(w, v)

			if v && len(c.CallExpression.Tokens) > 0 && c.PrivateIdentifier.Line > c.CallExpression.Tokens[len(c.CallExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteString(".")
			w.WriteString(c.PrivateIdentifier.Data)
		}
	}
}

func (b BindingProperty) printSource(w writer, v bool) {
	if !v && b.PropertyName.LiteralPropertyName != nil && b.BindingElement.SingleNameBinding != nil && b.PropertyName.LiteralPropertyName.Data == b.BindingElement.SingleNameBinding.Data {
		b.BindingElement.printSource(w, v)
	} else {
		b.PropertyName.printSource(w, v)
		w.WriteString(": ")
		b.BindingElement.printSource(w, v)
	}
}

func (b BindingElement) printSource(w writer, v bool) {
	if b.SingleNameBinding != nil {
		w.WriteString(b.SingleNameBinding.Data)
	} else if b.ArrayBindingPattern != nil {
		b.ArrayBindingPattern.printSource(w, v)
	} else if b.ObjectBindingPattern != nil {
		b.ObjectBindingPattern.printSource(w, v)
	} else {
		return
	}

	if b.Initializer != nil {
		w.WriteString(" = ")
		b.Initializer.printSource(w, v)
	}
}

func (p PropertyName) printSource(w writer, v bool) {
	if p.LiteralPropertyName != nil {
		w.WriteString(p.LiteralPropertyName.Data)
	} else if p.ComputedPropertyName != nil {
		ip := w.Indent()

		w.WriteString("[")

		if v {
			p.Comments[0].printSource(w, false, true)

			if p.ComputedPropertyName.hasFirstComment() {
				ip.WriteString("\n")
			}
		}

		p.ComputedPropertyName.printSource(ip, v)

		if v && len(p.Comments[1]) > 0 {
			w.WriteString("\n")
			p.Comments[1].printSource(w, false, false)
		}

		w.WriteString("]")
	}
}

func (l LogicalORExpression) printSource(w writer, v bool) {
	if l.LogicalORExpression != nil {
		l.LogicalORExpression.printSource(w, v)
		w.WriteString(" || ")
	}

	l.LogicalANDExpression.printSource(w, v)
}

func (c ParenthesizedExpression) printSource(w writer, v bool) {
	w.WriteString("(")

	ip := w.Indent()

	if v && len(c.Comments[0]) > 0 {
		c.Comments[0].printSource(w, true, true)
		ip.WriteString("\n")
	}

	if len(c.Expressions) > 0 {
		c.Expressions[0].printSource(ip, v)

		for _, e := range c.Expressions[1:] {
			ip.WriteString(", ")
			e.printSource(ip, v)
		}
	}

	if v && len(c.Comments[1]) > 0 {
		ip.WriteString("\n")
		c.Comments[1].printSource(w, true, false)
	}

	w.WriteString(")")
}

func (m MemberExpression) printSource(w writer, v bool) {
	if v {
		m.Comments[0].printSource(w, true, false)
	}

	if m.MemberExpression != nil {
		if m.Arguments != nil {
			w.WriteString("new ")
			if v {
				m.Comments[1].printSource(w, true, false)
			}

			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v {
				m.Comments[2].printSource(w, true, false)
			}

			m.Arguments.printSource(w.Indent(), v)
		} else if m.Expression != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			w.WriteString("[")

			ip := w.Indent()

			if v && len(m.Comments[1]) > 0 {
				m.Comments[1].printSource(w, true, false)
				ip.WriteString("\n")
			}

			m.Expression.printSource(ip, v)

			if v && len(m.Comments[2]) > 0 {
				w.WriteString("\n")
				m.Comments[2].printSource(w, true, true)
			}

			w.WriteString("]")
		} else if m.IdentifierName != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v && len(m.MemberExpression.Tokens) > 0 && m.IdentifierName.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteString(".")

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteString(m.IdentifierName.Data)
		} else if m.PrivateIdentifier != nil {
			m.MemberExpression.printSource(w, v)

			if v && m.MemberExpression.Comments[4].LastIsMulti() {
				w.WriteString(" ")
			}

			if v && len(m.MemberExpression.Tokens) > 0 && m.PrivateIdentifier.Line > m.MemberExpression.Tokens[len(m.MemberExpression.Tokens)-1].Line {
				w.WriteString("\n")
			}

			w.WriteString(".")

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteString(m.PrivateIdentifier.Data)
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
			w.WriteString("super")

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteString("[")

			ip := w.Indent()

			if v && len(m.Comments[2]) > 0 {
				m.Comments[2].printSource(w, true, false)
				ip.WriteString("\n")
			}

			m.Expression.printSource(ip, v)

			if v && len(m.Comments[3]) > 0 {
				w.WriteString("\n")
				m.Comments[3].printSource(w, false, false)
			}

			w.WriteString("]")
		} else if m.IdentifierName != nil {
			w.WriteString("super")

			if v {
				m.Comments[1].printSource(w, true, false)
			}

			w.WriteString(".")

			if v {
				m.Comments[2].printSource(w, true, false)
			}

			w.WriteString(m.IdentifierName.Data)
		}
	} else if m.NewTarget {
		w.WriteString("new")

		if v {
			m.Comments[1].printSource(w, true, false)
		}

		w.WriteString(".")

		if v {
			m.Comments[2].printSource(w, true, false)
		}

		w.WriteString("target")
	} else if m.ImportMeta {
		w.WriteString("import")

		if v {
			m.Comments[1].printSource(w, true, false)
		}

		w.WriteString(".")

		if v {
			m.Comments[2].printSource(w, true, false)
		}

		w.WriteString("meta")
	}

	if v {
		m.Comments[4].printSource(w, false, false)
	}
}

func (a Argument) printSource(w writer, v bool) {
	if a.Spread {
		if v {
			a.Comments.printSource(w, true, false)
		}

		w.WriteString("...")
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a Arguments) printSource(w writer, v bool) {
	w.WriteString("(")

	ip := w.Indent()

	if v && len(a.Comments[0]) > 0 {
		a.Comments[0].printSource(w, false, true)
	}

	if len(a.ArgumentList) > 0 {
		if v && a.ArgumentList[0].hasFirstComment() {
			ip.WriteString("\n")
		}

		a.ArgumentList[0].printSource(ip, v)

		for _, ae := range a.ArgumentList[1:] {
			ip.WriteString(", ")
			ae.printSource(ip, v)
		}
	}

	if v && len(a.Comments[1]) > 0 {
		w.WriteString("\n")
		a.Comments[1].printSource(w, false, true)
	}

	w.WriteString(")")
}

func (t TemplateLiteral) printSource(w writer, v bool) {
	x := w.Underlying()

	if t.NoSubstitutionTemplate != nil {
		if len(t.NoSubstitutionTemplate.Data) > 0 {
			w.WriteString(t.NoSubstitutionTemplate.Data[:1])
			x.WriteString(t.NoSubstitutionTemplate.Data[1:])
		}
	} else if t.TemplateHead != nil && t.TemplateTail != nil && len(t.Expressions) == len(t.TemplateMiddleList)+1 {
		if len(t.TemplateHead.Data) > 0 {
			w.WriteString(t.TemplateHead.Data[:1])
			x.WriteString(t.TemplateHead.Data[1:])
			t.Expressions[0].printSource(w, v)

			for n, e := range t.Expressions[1:] {
				x.WriteString(t.TemplateMiddleList[n].Data)
				e.printSource(w, v)
			}

			x.WriteString(t.TemplateTail.Data)
		}
	}
}

func (l LogicalANDExpression) printSource(w writer, v bool) {
	if l.LogicalANDExpression != nil {
		l.LogicalANDExpression.printSource(w, v)
		w.WriteString(" && ")
	}

	l.BitwiseORExpression.printSource(w, v)
}

func (p PrimaryExpression) printSource(w writer, v bool) {
	if p.This != nil {
		w.WriteString("this")
	} else if p.IdentifierReference != nil {
		w.WriteString(p.IdentifierReference.Data)
	} else if p.Literal != nil {
		w.WriteString(p.Literal.Data)
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

func (b BitwiseORExpression) printSource(w writer, v bool) {
	if b.BitwiseORExpression != nil {
		b.BitwiseORExpression.printSource(w, v)
		w.WriteString(" | ")
	}

	b.BitwiseXORExpression.printSource(w, v)
}

func (a ArrayElement) printSource(w writer, v bool) {
	if v {
		a.Comments.printSource(w, true, false)
	}

	if a.Spread {
		w.WriteString("...")
	}

	a.AssignmentExpression.printSource(w, v)
}

func (a ArrayLiteral) printSource(w writer, v bool) {
	w.WriteString("[")

	ip := w.Indent()

	if v && len(a.Comments[0]) > 0 {
		a.Comments[0].printSource(w, false, true)
	}

	if len(a.ElementList) > 0 {
		if v && a.ElementList[0].hasFirstComment() {
			ip.WriteString("\n")
		}

		a.ElementList[0].printSource(ip, v)

		for _, ae := range a.ElementList[1:] {
			ip.WriteString(", ")
			ae.printSource(ip, v)
		}
	}

	if v && len(a.Comments[1]) > 0 {
		w.WriteString("\n")
		a.Comments[1].printSource(w, false, false)
	}

	w.WriteString("]")
}

func (o ObjectLiteral) printSource(w writer, v bool) {
	w.WriteString("{")

	if len(o.PropertyDefinitionList) > 0 {
		var lastLine uint64

		x := w

		if v && len(o.Tokens) > 0 {
			lastLine = o.Tokens[0].Line
			x = w.Indent()
		}

		for n, pd := range o.PropertyDefinitionList {
			if n > 0 {
				if v && len(pd.Tokens) > 0 {
					if ll := pd.Tokens[0].Line; ll > lastLine {
						lastLine = ll

						x.WriteString(",\n")
					} else {
						x.WriteString(", ")
					}
				} else {
					x.WriteString(", ")
				}
			} else if v && len(pd.Tokens) > 0 {
				if ll := pd.Tokens[0].Line; ll > lastLine {
					lastLine = ll

					x.WriteString("\n")
				}
			}

			pd.printSource(x, v)
		}
		if v && len(o.Tokens) > 0 {
			if ll := o.Tokens[len(o.Tokens)-1].Line; ll > lastLine {
				w.WriteString("\n")
			}
		}
	}

	w.WriteString("}")
}

func (b BitwiseXORExpression) printSource(w writer, v bool) {
	if b.BitwiseXORExpression != nil {
		b.BitwiseXORExpression.printSource(w, v)
		w.WriteString(" ^ ")
	}

	b.BitwiseANDExpression.printSource(w, v)
}

func (p PropertyDefinition) printSource(w writer, v bool) {
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
					w.WriteString(" = ")
				} else {
					w.WriteString(": ")
				}

				p.AssignmentExpression.printSource(w, v)
			}
		} else {
			w.WriteString("...")
			p.AssignmentExpression.printSource(w, v)
		}
	} else if p.MethodDefinition != nil {
		p.MethodDefinition.printSource(w, v)
	}
}

func (b BitwiseANDExpression) printSource(w writer, v bool) {
	if b.BitwiseANDExpression != nil {
		b.BitwiseANDExpression.printSource(w, v)
		w.WriteString(" & ")
	}

	b.EqualityExpression.printSource(w, v)
}

func (e EqualityExpression) printSource(w writer, v bool) {
	if e.EqualityExpression != nil {
		var eo string

		switch e.EqualityOperator {
		case EqualityEqual:
			eo = " == "
		case EqualityNotEqual:
			eo = " != "
		case EqualityStrictEqual:
			eo = " === "
		case EqualityStrictNotEqual:
			eo = " !== "
		default:
			return
		}

		e.EqualityExpression.printSource(w, v)
		w.WriteString(eo)
	}

	e.RelationalExpression.printSource(w, v)
}

func (r RelationalExpression) printSource(w writer, v bool) {
	if r.PrivateIdentifier != nil {
		w.WriteString(r.PrivateIdentifier.Data)
		w.WriteString(" in ")
	} else if r.RelationalExpression != nil {
		var ro string

		switch r.RelationshipOperator {
		case RelationshipLessThan:
			ro = " < "
		case RelationshipGreaterThan:
			ro = " > "
		case RelationshipLessThanEqual:
			ro = " <= "
		case RelationshipGreaterThanEqual:
			ro = " >= "
		case RelationshipInstanceOf:
			ro = " instanceof "
		case RelationshipIn:
			ro = " in "
		default:
			return
		}

		r.RelationalExpression.printSource(w, v)
		w.WriteString(ro)
	}

	r.ShiftExpression.printSource(w, v)
}

func (s ShiftExpression) printSource(w writer, v bool) {
	if s.ShiftExpression != nil {
		var so string

		switch s.ShiftOperator {
		case ShiftLeft:
			so = " << "
		case ShiftRight:
			so = " >> "
		case ShiftUnsignedRight:
			so = " >>> "
		default:
			return
		}

		s.ShiftExpression.printSource(w, v)
		w.WriteString(so)
	}

	s.AdditiveExpression.printSource(w, v)
}

func (a AdditiveExpression) printSource(w writer, v bool) {
	if a.AdditiveExpression != nil {
		var ao string

		switch a.AdditiveOperator {
		case AdditiveAdd:
			ao = " + "
		case AdditiveMinus:
			ao = " - "
		default:
			return
		}

		a.AdditiveExpression.printSource(w, v)
		w.WriteString(ao)
	}

	a.MultiplicativeExpression.printSource(w, v)
}

func (m MultiplicativeExpression) printSource(w writer, v bool) {
	if m.MultiplicativeExpression != nil {
		var mo string

		switch m.MultiplicativeOperator {
		case MultiplicativeMultiply:
			mo = " * "
		case MultiplicativeDivide:
			mo = " / "
		case MultiplicativeRemainder:
			mo = " % "
		default:
			return
		}

		m.MultiplicativeExpression.printSource(w, v)
		w.WriteString(mo)
	}

	m.ExponentiationExpression.printSource(w, v)
}

func (e ExponentiationExpression) printSource(w writer, v bool) {
	if e.ExponentiationExpression != nil {
		e.ExponentiationExpression.printSource(w, v)
		w.WriteString(" ** ")
	}

	e.UnaryExpression.printSource(w, v)
}

func (u UnaryExpression) printSource(w writer, v bool) {
	for _, uo := range u.UnaryOperators {
		switch uo {
		case UnaryDelete:
			w.WriteString("delete ")
		case UnaryVoid:
			w.WriteString("void ")
		case UnaryTypeOf:
			w.WriteString("typeof ")
		case UnaryAdd:
			w.WriteString("+")
		case UnaryMinus:
			w.WriteString("-")
		case UnaryBitwiseNot:
			w.WriteString("~")
		case UnaryLogicalNot:
			w.WriteString("!")
		case UnaryAwait:
			w.WriteString("await ")
		}
	}

	u.UpdateExpression.printSource(w, v)
}

func (u UpdateExpression) printSource(w writer, v bool) {
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
			w.WriteString(uo)
		}
	} else if u.UnaryExpression != nil {
		switch u.UpdateOperator {
		case UpdatePreIncrement:
			w.WriteString("++")
		case UpdatePreDecrement:
			w.WriteString("--")
		default:
			return
		}

		u.UnaryExpression.printSource(w, v)
	}
}

func (m Module) printSource(w writer, v bool) {
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
	if m.ImportDeclaration != nil {
		m.ImportDeclaration.printSource(w, v)
	} else if m.ExportDeclaration != nil {
		m.ExportDeclaration.printSource(w, v)
	} else if m.StatementListItem != nil {
		m.StatementListItem.printSource(w, v)
	}
}

func (i ImportDeclaration) printSource(w writer, v bool) {
	if v {
		i.Comments[0].printSource(w, true, false)
	}

	if i.ImportClause == nil && i.FromClause.ModuleSpecifier == nil {
		return
	}

	w.WriteString("import ")

	if i.ImportClause != nil {
		i.ImportClause.printSource(w, v)
		i.FromClause.printSource(w, v)
	} else if i.FromClause.ModuleSpecifier != nil {
		if v {
			i.Comments[1].printSource(w, true, false)
		}

		w.WriteString(i.FromClause.ModuleSpecifier.Data)
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
	if e.FromClause != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		if e.ExportClause != nil {
			e.ExportClause.printSource(w, v)
		} else {
			w.WriteString("*")

			if v && len(e.Comments[2]) > 0 {
				w.WriteString(" ")
				e.Comments[2].printSource(w, false, false)
			}

			if e.ExportFromClause != nil {
				w.WriteString(" as ")

				if v && len(e.Comments[3]) > 0 {
					e.Comments[3].printSource(w, true, false)
				}

				w.WriteString(e.ExportFromClause.Data)
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

		w.WriteString("export ")

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

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		e.VariableStatement.printSource(w, v)
	} else if e.Declaration != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		e.Declaration.printSource(w, v)
	} else if e.DefaultFunction != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteString("default ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		e.DefaultFunction.printSource(w, v)
	} else if e.DefaultClass != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteString("default ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		e.DefaultClass.printSource(w, v)
	} else if e.DefaultAssignmentExpression != nil {
		if v {
			e.Comments[0].printSource(w, true, false)
		}

		w.WriteString("export ")

		if v {
			e.Comments[1].printSource(w, true, false)
		}

		w.WriteString("default ")

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
	if v {
		wc.Comments[0].printSource(w, true, false)
	}

	w.WriteString("with ")

	if v {
		wc.Comments[1].printSource(w, true, false)
	}

	w.WriteString("{")

	if v {
		wc.Comments[2].printSource(w, false, true)
	}

	if len(wc.WithEntries) > 0 {
		ip := w.Indent()

		if v && len(wc.WithEntries[0].Comments[0]) > 0 {
			ip.WriteString("\n")
		}

		wc.WithEntries[0].printSource(ip, v)

		for _, we := range wc.WithEntries[1:] {
			ip.WriteString(", ")
			we.printSource(ip, v)
		}
	}

	if v && len(wc.Comments[3]) > 0 {
		w.WriteString("\n")
		wc.Comments[3].printSource(w, false, false)
	}

	w.WriteString("}")
}

func (we WithEntry) printSource(w writer, v bool) {
	if we.AttributeKey == nil || we.Value == nil {
		return
	}

	if v {
		we.Comments[0].printSource(w, true, false)
	}

	w.WriteString(we.AttributeKey.Data)

	if v {
		we.Comments[1].printSource(w, false, false)
	}

	w.WriteString(": ")

	if v {
		we.Comments[2].printSource(w, true, false)
	}

	w.WriteString(we.Value.Data)

	if v {
		we.Comments[3].printSource(w, false, false)
	}
}

func (i ImportClause) printSource(w writer, v bool) {
	if v {
		i.Comments[0].printSource(w, true, false)
	}

	if i.ImportedDefaultBinding != nil {
		w.WriteString(i.ImportedDefaultBinding.Data)

		if v {
			i.Comments[1].printSource(w, true, false)
		}

		if i.NameSpaceImport != nil || i.NamedImports != nil {
			w.WriteString(", ")

			if v {
				i.Comments[2].printSource(w, true, false)
			}
		}
	}

	if i.NameSpaceImport != nil {
		w.WriteString("*")

		if v && len(i.Comments[3]) > 0 {
			i.Comments[3].printSource(w, true, false)
		} else {
			w.WriteString(" ")
		}

		w.WriteString("as ")

		if v {
			i.Comments[4].printSource(w, true, false)
		}

		w.WriteString(i.NameSpaceImport.Data)
	} else if i.NamedImports != nil {
		i.NamedImports.printSource(w, v)
	}

	if v {
		i.Comments[5].printSource(w, false, false)
	}
}

func (f FromClause) printSource(w writer, v bool) {
	if f.ModuleSpecifier == nil {
		return
	}

	w.WriteString(" from ")

	if v {
		f.Comments.printSource(w, true, false)
	}

	w.WriteString(f.ModuleSpecifier.Data)
}

func (e ExportClause) printSource(w writer, v bool) {
	w.WriteString("{")

	if v {
		e.Comments[0].printSource(w, false, true)
	}

	if len(e.ExportList) > 0 {
		ip := w.Indent()

		if v && len(e.ExportList[0].Comments[0]) > 0 {
			ip.WriteString("\n")
		}

		e.ExportList[0].printSource(ip, v)

		for _, es := range e.ExportList[1:] {
			ip.WriteString(", ")
			es.printSource(ip, v)
		}
	}

	if v && len(e.Comments[1]) > 0 {
		w.WriteString("\n")
		e.Comments[1].printSource(w, false, false)
	}

	w.WriteString("}")
}

func (n NamedImports) printSource(w writer, v bool) {
	w.WriteString("{")

	if v && len(n.Comments[0]) > 0 {
		n.Comments[0].printSource(w, false, true)
	}

	if len(n.ImportList) > 0 {
		ip := w.Indent()

		if v && len(n.ImportList[0].Comments[0]) > 0 {
			ip.WriteString("\n")
		}

		n.ImportList[0].printSource(ip, v)

		for _, is := range n.ImportList[1:] {
			ip.WriteString(", ")
			is.printSource(ip, v)
		}
	}

	if v && len(n.Comments[1]) > 0 {
		w.WriteString("\n")
		n.Comments[1].printSource(w, false, false)
	}

	w.WriteString("}")
}

func (e ExportSpecifier) printSource(w writer, v bool) {
	if v {
		e.Comments[0].printSource(w, true, false)
	}

	if e.IdentifierName == nil {
		return
	}

	w.WriteString(e.IdentifierName.Data)

	if e.EIdentifierName != nil && (e.EIdentifierName.Type != e.IdentifierName.Type || e.EIdentifierName.Data != e.IdentifierName.Data || v) {
		if v {
			e.Comments[1].printSource(w, false, false)
		}

		w.WriteString(" as ")

		if v {
			e.Comments[2].printSource(w, true, false)
		}

		w.WriteString(e.EIdentifierName.Data)
	}

	if v {
		e.Comments[3].printSource(w, false, false)
	}
}

func (i ImportSpecifier) printSource(w writer, v bool) {
	if i.ImportedBinding == nil {
		return
	}

	if v && len(i.Comments[0]) > 0 {
		i.Comments[0].printSource(w, true, false)
	}

	if i.IdentifierName != nil && (i.IdentifierName.Type != i.ImportedBinding.Type || i.IdentifierName.Data != i.ImportedBinding.Data || v) {
		w.WriteString(i.IdentifierName.Data)

		if v {
			i.Comments[1].printSource(w, false, false)
		}

		w.WriteString(" as ")

		if v {
			i.Comments[2].printSource(w, true, false)
		}
	}

	w.WriteString(i.ImportedBinding.Data)

	if v {
		i.Comments[3].printSource(w, false, false)
	}
}

func (oe OptionalExpression) printSource(w writer, v bool) {
	if oe.MemberExpression != nil {
		oe.MemberExpression.printSource(w, v)
	} else if oe.CallExpression != nil {
		oe.CallExpression.printSource(w, v)
	} else if oe.OptionalExpression != nil {
		oe.OptionalExpression.printSource(w, v)
	}

	oe.OptionalChain.printSource(w, v)
}

func (oe OptionalChain) printSource(w writer, v bool) {
	if oe.OptionalChain != nil {
		oe.OptionalChain.printSource(w, v)
	} else {
		w.WriteString("?.")
	}

	if oe.Arguments != nil {
		oe.Arguments.printSource(w, v)
	} else if oe.Expression != nil {
		w.WriteString("[")
		oe.Expression.printSource(w, v)
		w.WriteString("]")
	} else if oe.IdentifierName != nil {
		if oe.OptionalChain != nil {
			w.WriteString(".")
		}

		w.WriteString(oe.IdentifierName.Data)
	} else if oe.TemplateLiteral != nil {
		oe.TemplateLiteral.printSource(w, v)
	} else if oe.PrivateIdentifier != nil {
		if oe.OptionalChain != nil {
			w.WriteString(".")
		}

		w.WriteString(oe.PrivateIdentifier.Data)
	}
}

func (ce CoalesceExpression) printSource(w writer, v bool) {
	if ce.CoalesceExpressionHead != nil {
		ce.CoalesceExpressionHead.printSource(w, v)
		w.WriteString(" ?? ")
	}

	ce.BitwiseORExpression.printSource(w, v)
}
