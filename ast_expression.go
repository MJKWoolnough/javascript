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
	if yield && j.AcceptToken(parser.Token{TokenKeyword, "yield"}) {
		ae.Yield = true
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
			ae.Delegate = true
			j.AcceptRunWhitespace()
		}
		g := j.NewGoal()
		ae.AssignmentExpression = new(AssignmentExpression)
		if err := ae.AssignmentExpression.parse(&g, in, true, await); err != nil {
			return j.Error("AssignmentExpression", err)
		}
		j.Score(g)
	} else if j.Peek() == (parser.Token{TokenIdentifier, "async"}) { // TODO: Combine with next branch
		g := j.NewGoal()
		ae.ArrowFunction = new(ArrowFunction)
		if err := ae.ArrowFunction.parse(&g, nil, in, yield, await); err != nil {
			return j.Error("AssignmentExpression", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		ae.ConditionalExpression = new(ConditionalExpression)
		if err := ae.ConditionalExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("AssignmentExpression", err)
		}
		j.Score(g)
		if lhs := ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs != nil && len(ae.ConditionalExpression.Tokens) == len(lhs.Tokens) {
			g := j.NewGoal()
			if lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.CoverParenthesizedExpressionAndArrowParameterList != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil) {
				g.AcceptRunWhitespaceNoNewLine()
				if g.Peek() == (parser.Token{TokenPunctuator, "=>"}) {
					ae.ConditionalExpression = nil
					ae.ArrowFunction = new(ArrowFunction)
					if err := ae.ArrowFunction.parse(j, lhs.NewExpression.MemberExpression.PrimaryExpression, in, yield, await); err != nil {
						return g.Error("AssignmentExpression", err)
					}
				}
			}
			if ae.ConditionalExpression != nil {
				g.AcceptRunWhitespace()
				if err := ae.AssignmentOperator.parse(&g); err == nil {
					j.Score(g)
					j.AcceptRunWhitespace()
					ae.ConditionalExpression = nil
					ae.LeftHandSideExpression = lhs
					g = j.NewGoal()
					ae.AssignmentExpression = new(AssignmentExpression)
					if err := ae.AssignmentExpression.parse(&g, in, yield, await); err != nil {
						return j.Error("AssignmentExpression", err)
					}
					j.Score(g)
				}
			}
		}
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
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenKeyword, "super"}) || g.AcceptToken(parser.Token{TokenKeyword, "import"}) {
		g.AcceptRunWhitespace()
		if g.Peek() == (parser.Token{TokenPunctuator, "("}) {
			g = j.NewGoal()
			lhs.CallExpression = new(CallExpression)
			if err := lhs.CallExpression.parse(&g, nil, yield, await); err != nil {
				return j.Error("LeftHandSideExpression", err)
			}
			j.Score(g)
		}
	}
	if lhs.CallExpression == nil {
		g = j.NewGoal()
		lhs.NewExpression = new(NewExpression)
		if err := lhs.NewExpression.parse(&g, yield, await); err != nil {
			return j.Error("LeftHandSideExpression", err)
		}
		j.Score(g)
		if lhs.NewExpression.News == 0 {
			g = j.NewGoal()
			g.AcceptRunWhitespace()
			if g.Peek() == (parser.Token{TokenPunctuator, "("}) {
				lhs.CallExpression = new(CallExpression)
				if err := lhs.CallExpression.parse(j, &lhs.NewExpression.MemberExpression, yield, await); err != nil {
					return j.Error("LeftHandSideExpression", err)
				}
				lhs.NewExpression = nil
			}
		}
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
		ne.MemberExpression = MemberExpression{}
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
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenKeyword, "super"}) {
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			me.Expression = new(Expression)
			if err := me.Expression.parse(&h, true, yield, await); err != nil {
				return g.Error("MemberExpression", err)
			}
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return g.Error("MemberExpression", ErrInvalidSuperProperty)
			}
		} else if g.AcceptToken(parser.Token{TokenPunctuator, "."}) {
			g.AcceptRunWhitespace()
			if !g.Accept(TokenIdentifier, TokenKeyword) {
				return g.Error("MemberExpression", ErrNoIdentifier)
			}
			me.IdentifierName = g.GetLastToken()
		} else {
			return g.Error("MemberExpression", ErrInvalidSuperProperty)
		}
		me.SuperProperty = true
	} else if g.AcceptToken(parser.Token{TokenKeyword, "new"}) {
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, "."}) {
			g.AcceptRunWhitespace()
			if !g.AcceptToken(parser.Token{TokenIdentifier, "target"}) {
				return j.Error("MemberExpression", ErrInvalidMetaProperty)
			}
			me.MetaProperty = true
		} else {
			h := g.NewGoal()
			me.MemberExpression = new(MemberExpression)
			if err := me.MemberExpression.parse(&h, yield, await); err != nil {
				return g.Error("MemberExpression", err)
			}
			g.Score(h)
			g.AcceptRunWhitespace()
			h = g.NewGoal()
			me.Arguments = new(Arguments)
			if err := me.Arguments.parse(&h, yield, await); err != nil {
				return g.Error("MemberExpression", err)
			}
			g.Score(h)
		}
	} else {
		me.PrimaryExpression = new(PrimaryExpression)
		if err := me.PrimaryExpression.parse(&g, yield, await); err != nil {
			return j.Error("MemberExpression", err)
		}
	}
	j.Score(g)
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
			tl = new(TemplateLiteral)
			if err := tl.parse(&h, yield, await); err != nil {
				return g.Error("MemberExpression", err)
			}
		case TokenPunctuator:
			switch tk.Data {
			case ".":
				h.Except()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("MemberExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Except()
				h.AcceptRunWhitespace()
				i := h.NewGoal()
				e = new(Expression)
				if err := e.parse(&i, true, yield, await); err != nil {
					return g.Error("MemberExpression", err)
				}
				h.Score(i)
				h.AcceptRunWhitespace()
				if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
					return g.Error("MemberExpression", ErrMissingClosingBracket)
				}
			default:
				return nil
			}
		default:
			return nil
		}
		g.Score(h)
		nme := new(MemberExpression)
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
	if j.AcceptToken(parser.Token{TokenKeyword, "this"}) {
		pe.This = true
	} else if j.Accept(TokenNullLiteral, TokenBooleanLiteral, TokenNumericLiteral, TokenStringLiteral, TokenRegularExpressionLiteral) {
		pe.Literal = j.GetLastToken()
	} else if t := j.Peek(); t == (parser.Token{TokenPunctuator, "["}) {
		g := j.NewGoal()
		pe.ArrayLiteral = new(ArrayLiteral)
		if err := pe.ArrayLiteral.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{TokenPunctuator, "{"}) {
		g := j.NewGoal()
		pe.ObjectLiteral = new(ObjectLiteral)
		if err := pe.ObjectLiteral.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{TokenIdentifier, "async"}) || t == (parser.Token{TokenKeyword, "function"}) {
		g := j.NewGoal()
		pe.FunctionExpression = new(FunctionDeclaration)
		if err := pe.FunctionExpression.parse(&g, false, false, true); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{TokenKeyword, "class"}) {
		g := j.NewGoal()
		pe.ClassExpression = new(ClassDeclaration)
		if err := pe.ClassExpression.parse(&g, yield, await, true); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t.Type == TokenNoSubstitutionTemplate || t.Type == TokenTemplateHead {
		g := j.NewGoal()
		pe.TemplateLiteral = new(TemplateLiteral)
		if err := pe.TemplateLiteral.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{TokenPunctuator, "("}) {
		g := j.NewGoal()
		pe.CoverParenthesizedExpressionAndArrowParameterList = new(CoverParenthesizedExpressionAndArrowParameterList)
		if err := pe.CoverParenthesizedExpressionAndArrowParameterList.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		if pe.IdentifierReference = g.parseIdentifier(yield, await); pe.IdentifierReference == nil {
			return j.Error("PrimaryExpression", ErrNoIdentifier)
		}
		j.Score(g)
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
					cp.ArrayBindingPattern = new(ArrayBindingPattern)
					if err := cp.ArrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
					cp.ObjectBindingPattern = new(ObjectBindingPattern)
					if err := cp.ObjectBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if cp.BindingIdentifier = g.parseIdentifier(yield, await); cp.BindingIdentifier == nil {
					return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrNoIdentifier)
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
			a.SpreadArgument = new(AssignmentExpression)
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

func (ce *CallExpression) parse(j *jsParser, me *MemberExpression, yield, await bool) error {
	if me == nil {
		if j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
			ce.SuperCall = true
		} else if j.AcceptToken(parser.Token{TokenKeyword, "import"}) {
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
				return j.Error("CallExpression", ErrMissingOpeningParenthesis)
			}
			g := j.NewGoal()
			ce.ImportCall = new(AssignmentExpression)
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
		}
	} else {
		ce.MemberExpression = me
	}
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	ce.Arguments = new(Arguments)
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
			tl = new(TemplateLiteral)
			if err := tl.parse(&h, yield, await); err != nil {
				return g.Error("CallExpression", err)
			}
		case TokenPunctuator:
			switch tk.Data {
			case "(":
				a = new(Arguments)
				if err := a.parse(&h, yield, await); err != nil {
					return g.Error("CallExpression", err)
				}
			case ".":
				h.Except()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("CallExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Except()
				h.AcceptRunWhitespace()
				i := h.NewGoal()
				e = new(Expression)
				if err := e.parse(&i, true, yield, await); err != nil {
					return g.Error("CallExpression", err)
				}
				h.Score(i)
				h.AcceptRunWhitespace()
				if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
					return g.Error("CallExpression", ErrMissingClosingBracket)
				}
			default:
				return nil
			}
		default:
			return nil
		}
		g.Score(h)
		nce := new(CallExpression)
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
