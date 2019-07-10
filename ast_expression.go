package javascript

import "vimagination.zapto.org/parser"

type AssignmentOperator uint8

const (
	AssignmentNone AssignmentOperator = iota
	AssignmentAssign
	AssignmentMultiply
	AssignmentDivide
	AssignmentRemainder
	AssignmentAdd
	AssignmentSubtract
	AssignmentLeftShift
	AssignmentSignPropagatinRightShift
	AssignmentZeroFillRightShift
	AssignmentBitwiseAND
	AssignmentBitwiseXOR
	AssignmentBitwiseOR
	AssignmentExponentiation
)

func (ao *AssignmentOperator) parse(j *jsParser) error {
	switch j.Peek() {
	case parser.Token{TokenPunctuator, "="}:
		*ao = AssignmentAssign
	case parser.Token{TokenPunctuator, "*="}:
		*ao = AssignmentMultiply
	case parser.Token{TokenDivPunctuator, "/="}:
		*ao = AssignmentDivide
	case parser.Token{TokenPunctuator, "%="}:
		*ao = AssignmentRemainder
	case parser.Token{TokenPunctuator, "+="}:
		*ao = AssignmentAdd
	case parser.Token{TokenPunctuator, "-="}:
		*ao = AssignmentSubtract
	case parser.Token{TokenPunctuator, "<<="}:
		*ao = AssignmentLeftShift
	case parser.Token{TokenPunctuator, ">>="}:
		*ao = AssignmentSignPropagatinRightShift
	case parser.Token{TokenPunctuator, ">>>="}:
		*ao = AssignmentZeroFillRightShift
	case parser.Token{TokenPunctuator, "&="}:
		*ao = AssignmentBitwiseAND
	case parser.Token{TokenPunctuator, "^="}:
		*ao = AssignmentBitwiseXOR
	case parser.Token{TokenPunctuator, "|="}:
		*ao = AssignmentBitwiseOR
	case parser.Token{TokenPunctuator, "**="}:
		*ao = AssignmentExponentiation
	default:
		return ErrInvalidAssignment
	}
	j.Except()
	return nil
}

type AssignmentExpression struct {
	ConditionalExpression  *ConditionalExpression
	ArrowFunction          *ArrowFunction
	LeftHandSideExpression *LeftHandSideExpression
	Yield                  bool
	Delegate               bool
	AssignmentOperator     AssignmentOperator
	AssignmentExpression   *AssignmentExpression
	Tokens                 Tokens
}

func (ae *AssignmentExpression) parse(j *jsParser, in, yield, await bool) error {
	if err := j.FindGoal(func(j *jsParser) error {
		if !yield || !j.AcceptToken(parser.Token{TokenKeyword, "yield"}) {
			return errNotApplicable
		}
		ae.Yield = true
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
			ae.Delegate = true
			j.AcceptRunWhitespace()
		}
		g := j.NewGoal()
		nae := newAssignmentExpression()
		if err := nae.parse(&g, in, true, await); err != nil {
			nae.clear()
			poolAssignmentExpression.Put(nae)
			return j.Error("AssignmentExpression.Yield", err)
		}
		j.Score(g)
		ae.AssignmentExpression = nae
		return nil
	}, func(j *jsParser) error {
		g := j.NewGoal()
		af := newArrowFunction()
		if err := af.parse(&g, in, yield, await); err != nil {
			af.clear()
			poolArrowFunction.Put(af)
			return j.Error("AssignmentExpression.ArrowFunction", err)
		}
		j.Score(g)
		ae.ArrowFunction = af
		return nil
	}, func(j *jsParser) error {
		g := j.NewGoal()
		lhs := newLeftHandSideExpression()
		if err := lhs.parse(&g, yield, await); err != nil {
			lhs.clear()
			poolLeftHandSideExpression.Put(lhs)
			return j.Error("AssignmentExpression.LeftHandSideExpression", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if err := ae.AssignmentOperator.parse(j); err != nil {
			lhs.clear()
			poolLeftHandSideExpression.Put(lhs)
			ae.AssignmentOperator = 0
			return j.Error("AssignmentExpression.LeftHandSideExpression", err)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		nae := newAssignmentExpression()
		if err := nae.parse(&g, in, yield, await); err != nil {
			lhs.clear()
			poolLeftHandSideExpression.Put(lhs)
			ae.AssignmentOperator = 0
			nae.clear()
			poolAssignmentExpression.Put(nae)
			return j.Error("AssignmentExpression.LeftHandSideExpression", err)
		}
		j.Score(g)
		ae.LeftHandSideExpression = lhs
		ae.AssignmentExpression = nae
		return nil
	}, func(j *jsParser) error {
		g := j.NewGoal()
		ce := newConditionalExpression()
		if err := ce.parse(&g, in, yield, await); err != nil {
			ce.clear()
			poolConditionalExpression.Put(ce)
			return j.Error("AssignmentExpression.ConditionalExpression", err)
		}
		j.Score(g)
		ae.ConditionalExpression = ce
		return nil
	}); err != nil {
		return err
	}
	ae.Tokens = j.ToTokens()
	return nil
}

type LeftHandSideExpression struct {
	NewExpression  *NewExpression
	CallExpression *CallExpression
	Tokens         Tokens
}

func (lhs *LeftHandSideExpression) parse(j *jsParser, yield, await bool) error {
	if err := j.FindGoal(func(j *jsParser) error {
		g := j.NewGoal()
		ce := newCallExpression()
		if err := ce.parse(&g, yield, await); err != nil {
			ce.clear()
			poolCallExpression.Put(ce)
			return j.Error("LeftHandSideExpression.CallExpression", err)
		}
		j.Score(g)
		lhs.CallExpression = ce
		return nil
	}, func(j *jsParser) error {
		g := j.NewGoal()
		ne := newNewExpression()
		if err := ne.parse(&g, yield, await); err != nil {
			ne.clear()
			poolNewExpression.Put(ne)
			return j.Error("LeftHandSideExpression.NewExpression", err)
		}
		j.Score(g)
		lhs.NewExpression = ne
		return nil
	}); err != nil {
		return err
	}
	lhs.Tokens = j.ToTokens()
	return nil
}

type Expression struct {
	Expressions []AssignmentExpression
	Tokens      Tokens
}

func (e *Expression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		var ae AssignmentExpression
		if err := ae.parse(&g, in, yield, await); err != nil {
			ae.clear()
			return j.Error("Expression", err)
		}
		j.Score(g)
		e.Expressions = append(e.Expressions, ae)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			break
		}
		j.Score(g)
	}
	e.Tokens = j.ToTokens()
	return nil
}

type NewExpression struct {
	News             uint
	MemberExpression MemberExpression
	Tokens           Tokens
}

func (ne *NewExpression) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := ne.MemberExpression.parse(&g, yield, await); err == nil {
			j.Score(g)
			break
		} else if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
			return j.Error("NewExpression", err)
		}
		ne.MemberExpression.clear()
		ne.News++
		j.AcceptRunWhitespace()
	}
	ne.Tokens = j.ToTokens()
	return nil
}

type MemberExpression struct {
	MemberExpression  *MemberExpression
	PrimaryExpression *PrimaryExpression
	Expression        *Expression
	IdentifierName    *Token
	TemplateLiteral   *TemplateLiteral
	SuperProperty     bool
	MetaProperty      bool
	Arguments         *Arguments
	Tokens            Tokens
}

func (me *MemberExpression) parse(j *jsParser, yield, await bool) error {
	if err := j.FindGoal(
		func(j *jsParser) error {
			g := j.NewGoal()
			pe := newPrimaryExpression()
			if err := pe.parse(&g, yield, await); err != nil {
				pe.clear()
				poolPrimaryExpression.Put(pe)
				return j.Error("MemberExpression.PrimaryExpression", err)
			}
			j.Score(g)
			me.PrimaryExpression = pe
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
				g := j.NewGoal()
				e := newExpression()
				if err := e.parse(&g, true, yield, await); err != nil {
					e.clear()
					poolExpression.Put(e)
					return j.Error("MemberExpression.Super", err)
				}
				j.Score(g)
				if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
					e.clear()
					poolExpression.Put(e)
					return j.Error("MemberExpression.Super", ErrInvalidSuperProperty)
				}
				me.Expression = e
			} else if j.AcceptToken(parser.Token{TokenPunctuator, "."}) {
				if !j.Accept(TokenIdentifier, TokenKeyword) {
					return j.Error("MemberExpression.Super", ErrNoIdentifier)
				}
				me.IdentifierName = j.GetLastToken()
			} else {
				return j.Error("MemberExpression.Super", ErrInvalidSuperProperty)
			}
			me.SuperProperty = true
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, "."}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenIdentifier, "target"}) {
				return j.Error("MemberExpression.MetaProperty", ErrInvalidMetaProperty)
			}
			me.MetaProperty = true
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			nme := newMemberExpression()
			if err := nme.parse(&g, yield, await); err != nil {
				nme.clear()
				poolMemberExpression.Put(nme)
				return j.Error("MemberExpression.New", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			g = j.NewGoal()
			a := newArguments()
			if err := a.parse(&g, yield, await); err != nil {
				nme.clear()
				poolMemberExpression.Put(nme)
				a.clear()
				poolArguments.Put(a)
				return j.Error("MemberExpression.New", err)
			}
			j.Score(g)
			me.MemberExpression = nme
			me.Arguments = a
			return nil
		},
	); err != nil {
		return err
	}
	for {
		me.Tokens = j.ToTokens()
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		var (
			tl *TemplateLiteral
			i  *Token
			e  *Expression
		)
		switch tk := h.Peek(); tk.Type {
		case TokenNoSubstitutionTemplate, TokenTemplateHead:
			tl = newTemplateLiteral()
			if err := tl.parse(&h, yield, await); err != nil {
				tl.clear()
				poolTemplateLiteral.Put(tl)
				return g.Error("MemberExpression", err)
			}
		case TokenPunctuator:
			switch tk.Data {
			case ".":
				h.Except()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("MemberExpression", ErrMissingIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Except()
				h.AcceptRunWhitespace()
				i := h.NewGoal()
				e = newExpression()
				if err := e.parse(&i, true, yield, await); err != nil {
					e.clear()
					poolExpression.Put(e)
					return g.Error("MemberExpression", err)
				}
				h.Score(i)
				h.AcceptRunWhitespace()
				if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
					e.clear()
					poolExpression.Put(e)
					return g.Error("MemberExpression", ErrMissingClosingBracket)
				}
			default:
				return nil
			}
		default:
			return nil
		}
		g.Score(h)
		nme := newMemberExpression()
		*nme = *me
		*me = MemberExpression{
			MemberExpression: nme,
			Expression:       e,
			IdentifierName:   i,
			TemplateLiteral:  tl,
		}
		j.Score(g)
	}
}

type PrimaryExpression struct {
	This                                              bool
	IdentifierReference                               *Token
	Literal                                           *Token
	ArrayLiteral                                      *ArrayLiteral
	ObjectLiteral                                     *ObjectLiteral
	FunctionExpression                                *FunctionDeclaration
	ClassExpression                                   *ClassDeclaration
	TemplateLiteral                                   *TemplateLiteral
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	Tokens                                            Tokens
}

func (pe *PrimaryExpression) parse(j *jsParser, yield, await bool) error {
	if err := j.FindGoal(
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "this"}) {
				return errNotApplicable
			}
			pe.This = true
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			i, err := g.parseIdentifier(yield, await)
			if err != nil {
				return j.Error("PrimaryExpression.IdentifierReference", err)
			}
			j.Score(g)
			pe.IdentifierReference = i
			return nil
		},
		func(j *jsParser) error {
			if !j.Accept(TokenNullLiteral, TokenBooleanLiteral, TokenNumericLiteral, TokenStringLiteral, TokenRegularExpressionLiteral) {
				return errNotApplicable
			}
			pe.Literal = j.GetLastToken()
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			al := newArrayLiteral()
			if err := al.parse(&g, yield, await); err != nil {
				al.clear()
				poolArrayLiteral.Put(al)
				return j.Error("PrimaryExpression.ArrayLiteral", err)
			}
			j.Score(g)
			pe.ArrayLiteral = al
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			ol := newObjectLiteral()
			if err := ol.parse(&g, yield, await); err != nil {
				ol.clear()
				poolObjectLiteral.Put(ol)
				return j.Error("PrimaryExpression.ObjectLiteral", err)
			}
			j.Score(g)
			pe.ObjectLiteral = ol
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			fe := newFunctionDeclaration()
			if err := fe.parse(&g, false, false, true); err != nil {
				fe.clear()
				poolFunctionDeclaration.Put(fe)
				return j.Error("PrimaryExpression.FunctionDeclaration", err)
			}
			j.Score(g)
			pe.FunctionExpression = fe
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			ce := newClassDeclaration()
			if err := ce.parse(&g, yield, await, true); err != nil {
				ce.clear()
				poolClassDeclaration.Put(ce)
				return j.Error("PrimaryExpression.ClassExpression", err)
			}
			j.Score(g)
			pe.ClassExpression = ce
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			tl := newTemplateLiteral()
			if err := tl.parse(&g, yield, await); err != nil {
				tl.clear()
				poolTemplateLiteral.Put(tl)
				return j.Error("PrimaryExpression.TemplateLiteral", err)
			}
			j.Score(g)
			pe.TemplateLiteral = tl
			return nil
		},
		func(j *jsParser) error {
			g := j.NewGoal()
			cp := newCoverParenthesizedExpressionAndArrowParameterList()
			if err := cp.parse(&g, yield, await); err != nil {
				cp.clear()
				poolCoverParenthesizedExpressionAndArrowParameterList.Put(cp)
				return j.Error("PrimaryExpression.CoverParenthesizedExpressionAndArrowParameterList", err)
			}
			j.Score(g)
			pe.CoverParenthesizedExpressionAndArrowParameterList = cp
			return nil
		},
	); err != nil {
		return err
	}
	pe.Tokens = j.ToTokens()
	return nil
}

type CoverParenthesizedExpressionAndArrowParameterList struct {
	Expressions          []AssignmentExpression
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (cp *CoverParenthesizedExpressionAndArrowParameterList) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		for {
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
					cp.ArrayBindingPattern = newArrayBindingPattern()
					if err := cp.ArrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
					cp.ObjectBindingPattern = newObjectBindingPattern()
					if err := cp.ObjectBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else {
					bi, err := g.parseIdentifier(yield, await)
					if err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
					cp.BindingIdentifier = bi
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
					return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrMissingClosingParenthesis)
				}
				break
			}
			g := j.NewGoal()
			var e AssignmentExpression
			if err := e.parse(&g, true, yield, await); err != nil {
				e.clear()
				return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
			}
			j.Score(g)
			cp.Expressions = append(cp.Expressions, e)
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	cp.Tokens = j.ToTokens()
	return nil
}

type Arguments struct {
	ArgumentList   []AssignmentExpression
	SpreadArgument *AssignmentExpression
	Tokens         Tokens
}

func (a *Arguments) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("Arguments", ErrMissingOpeningParenthesis)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		}
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			g.AcceptRunWhitespace()
			a.SpreadArgument = newAssignmentExpression()
			if err := a.SpreadArgument.parse(&g, true, yield, await); err != nil {
				return j.Error("Arguments", err)
			}
			j.Score(g)
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return j.Error("Arguments", ErrMissingClosingParenthesis)
			}
			break
		}
		var ae AssignmentExpression
		if err := ae.parse(&g, true, yield, await); err != nil {
			ae.clear()
			return j.Error("Arguments", err)
		}
		j.Score(g)
		a.ArgumentList = append(a.ArgumentList, ae)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("Arguments", ErrMissingComma)
		}
	}
	a.Tokens = j.ToTokens()
	return nil
}

type CallExpression struct {
	MemberExpression *MemberExpression
	SuperCall        bool
	ImportCall       *AssignmentExpression
	CallExpression   *CallExpression
	Arguments        *Arguments
	Expression       *Expression
	IdentifierName   *Token
	TemplateLiteral  *TemplateLiteral
	Tokens           Tokens
}

func (ce *CallExpression) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
		ce.SuperCall = true
	} else if j.AcceptToken(parser.Token{TokenKeyword, "import"}) {
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			return j.Error("CallExpression", ErrMissingOpeningParenthesis)
		}
		g := j.NewGoal()
		ce.ImportCall = newAssignmentExpression()
		if err := ce.ImportCall.parse(&g, true, yield, await); err != nil {
			return j.Error("CallExpression", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return j.Error("CallExpression", ErrMissingClosingParenthesis)
		}
		ce.Tokens = j.ToTokens()
		return nil
	} else {
		g := j.NewGoal()
		ce.MemberExpression = newMemberExpression()
		if err := ce.MemberExpression.parse(&g, yield, await); err != nil {
			return j.Error("CallExpression", err)
		}
		j.Score(g)
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	ce.Arguments = newArguments()
	if err := ce.Arguments.parse(&g, yield, await); err != nil {
		return j.Error("CallExpression", err)
	}
	j.Score(g)
	for {
		ce.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		var (
			tl *TemplateLiteral
			a  *Arguments
			i  *Token
			e  *Expression
		)
		switch tk := h.Peek(); tk.Type {
		case TokenNoSubstitutionTemplate, TokenTemplateHead:
			tl = newTemplateLiteral()
			if err := tl.parse(&h, yield, await); err != nil {
				return g.Error("CallExpression", err)
			}
		case TokenPunctuator:
			switch tk.Data {
			case "(":
				a = newArguments()
				if err := a.parse(&h, yield, await); err != nil {
					a.clear()
					poolArguments.Put(a)
					return g.Error("CallExpression", err)
				}
			case ".":
				h.Except()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("CallExpression", ErrMissingIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Except()
				h.AcceptRunWhitespace()
				i := h.NewGoal()
				e = newExpression()
				if err := e.parse(&i, true, yield, await); err != nil {
					e.clear()
					poolExpression.Put(e)
					return g.Error("CallExpression", err)
				}
				h.Score(i)
				h.AcceptRunWhitespace()
				if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
					e.clear()
					poolExpression.Put(e)
					return g.Error("CallExpression", ErrMissingClosingBracket)
				}
			default:
				return nil
			}
		default:
			return nil
		}
		g.Score(h)
		nce := newCallExpression()
		*nce = *ce
		*ce = CallExpression{
			CallExpression:  nce,
			Expression:      e,
			Arguments:       a,
			IdentifierName:  i,
			TemplateLiteral: tl,
		}
		j.Score(g)
	}
}
