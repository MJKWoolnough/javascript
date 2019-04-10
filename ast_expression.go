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

func (j *jsParser) parseAssignmentOperator() (AssignmentOperator, error) {
	var ao AssignmentOperator
	switch j.Peek() {
	case parser.Token{TokenPunctuator, "="}:
		ao = AssignmentAssign
	case parser.Token{TokenPunctuator, "*="}:
		ao = AssignmentMultiply
	case parser.Token{TokenPunctuator, "/="}:
		ao = AssignmentDivide
	case parser.Token{TokenPunctuator, "%="}:
		ao = AssignmentRemainder
	case parser.Token{TokenPunctuator, "+="}:
		ao = AssignmentAdd
	case parser.Token{TokenPunctuator, "-="}:
		ao = AssignmentSubtract
	case parser.Token{TokenPunctuator, "<<="}:
		ao = AssignmentLeftShift
	case parser.Token{TokenPunctuator, ">>="}:
		ao = AssignmentSignPropagatinRightShift
	case parser.Token{TokenPunctuator, ">>>="}:
		ao = AssignmentZeroFillRightShift
	case parser.Token{TokenPunctuator, "&="}:
		ao = AssignmentBitwiseAND
	case parser.Token{TokenPunctuator, "^="}:
		ao = AssignmentBitwiseXOR
	case parser.Token{TokenPunctuator, "|="}:
		ao = AssignmentBitwiseOR
	case parser.Token{TokenPunctuator, "**="}:
		ao = AssignmentExponentiation
	default:
		return 0, ErrInvalidAssignment
	}
	j.Except()
	return ao, nil
}

type AssignmentExpression struct {
	ConditionalExpression  *ConditionalExpression
	ArrowFunction          *ArrowFunction
	LeftHandSideExpression *LeftHandSideExpression
	Yield                  bool
	Delegate               bool
	AssignmentOperator     AssignmentOperator
	AssignmentExpression   *AssignmentExpression
	Tokens                 []Token
}

func (j *jsParser) parseAssignmentExpression(in, yield, await bool) (AssignmentExpression, error) {
	var ae AssignmentExpression
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
		nae, err := g.parseAssignmentExpression(in, true, await)
		if err != nil {
			return err
		}
		j.Score(g)
		ae.AssignmentExpression = &nae
		return nil
	}, func(j *jsParser) error {
		af, err := j.parseArrowFunction(in, yield, await)
		if err != nil {
			return err
		}
		ae.ArrowFunction = &af
		return nil
	}, func(j *jsParser) error {
		ce, err := j.parseConditionalExpression(in, yield, await)
		if err != nil {
			return err
		}
		ae.ConditionalExpression = &ce
		return nil
	}, func(j *jsParser) error {
		lhs, err := j.parseLeftHandSideExpression(yield, await)
		if err != nil {
			return err
		}
		ae.LeftHandSideExpression = &lhs
		j.AcceptRunWhitespace()
		ae.AssignmentOperator, err = j.parseAssignmentOperator()
		if err != nil {
			return err
		}
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		nae, err := g.parseAssignmentExpression(in, yield, await)
		if err != nil {
			return err
		}
		j.Score(g)
		ae.AssignmentExpression = &nae
		return nil
	}); err != nil {
		return ae, err
	}
	ae.Tokens = j.ToTokens()
	return ae, nil
}

type LeftHandSideExpression struct {
	NewExpression  *NewExpression
	CallExpression *CallExpression
	Tokens         []Token
}

func (j *jsParser) parseLeftHandSideExpression(yield, await bool) (LeftHandSideExpression, error) {
	var lhs LeftHandSideExpression
	if err := j.FindGoal(func(j *jsParser) error {
		ne, err := j.parseNewExpression(yield, await)
		if err != nil {
			return err
		}
		lhs.NewExpression = &ne
		return nil
	}, func(j *jsParser) error {
		ce, err := j.parseCallExpression(yield, await)
		if err != nil {
			return err
		}
		lhs.CallExpression = &ce
		return nil
	}); err != nil {
		return lhs, err
	}
	lhs.Tokens = j.ToTokens()
	return lhs, nil
}

type Expression struct {
	Expressions []AssignmentExpression
	Tokens      []Token
}

func (j *jsParser) parseExpression(in, yield, await bool) (Expression, error) {
	var e Expression
	for {
		g := j.NewGoal()
		ae, err := g.parseAssignmentExpression(in, yield, await)
		if err != nil {
			return e, j.Error(err)
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
	return e, nil
}

type NewExpression struct {
	News             uint
	MemberExpression MemberExpression
	Tokens           []Token
}

func (j *jsParser) parseNewExpression(yield, await bool) (NewExpression, error) {
	var (
		ne  NewExpression
		err error
	)
	for {
		if !j.AcceptToken(parser.Token{TokenKeyword, "new"}) {
			break
		}
		ne.News++
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	ne.MemberExpression, err = g.parseMemberExpression(yield, await)
	if err != nil {
		return ne, j.Error(err)
	}
	j.Score(g)
	ne.Tokens = j.ToTokens()
	return ne, nil
}

type MemberExpression struct {
	MemberExpression    *MemberExpression
	PrimaryExpression   *PrimaryExpression
	Expression          *Expression
	IdentifierName      *Token
	TemplateLiteral     *TemplateLiteral
	SuperProperty       bool
	MetaProperty        bool
	NewTarget           bool
	NewMemberExpression *MemberExpression
	Arguments           *Arguments
	Tokens              []Token
}

func (j *jsParser) parseMemberExpression(yield, await bool) (MemberExpression, error) {
	var me MemberExpression
	if err := j.FindGoal(
		func(j *jsParser) error {
			pe, err := j.parserPrimaryExpression(yield, await)
			if err != nil {
				return err
			}
			me.PrimaryExpression = &pe
			return nil
		},
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
				return errNotApplicable
			}
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
				g := j.NewGoal()
				e, err := g.parseExpression(true, yield, await)
				if err != nil {
					return err
				}
				j.Score(g)
				me.Expression = &e
			} else if j.AcceptToken(parser.Token{TokenPunctuator, "."}) {
				if !j.Accept(TokenIdentifier) {
					return ErrNoIdentifier
				}
				me.IdentifierName = j.GetLastToken()
			} else {
				return ErrInvalidSuperProperty
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
				return ErrInvalidMetaProperty
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
			nme, err := g.parseMemberExpression(yield, await)
			if err != nil {
				return err
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			g = j.NewGoal()
			a, err := g.parseArguments(yield, await)
			if err != nil {
				return err
			}
			j.Score(g)
			me.NewMemberExpression = &nme
			me.Arguments = &a
			return nil
		},
	); err != nil {
		return me, err
	}
Loop:
	for {
		var nme MemberExpression
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		switch tk := h.Peek(); tk.Type {
		case TokenNoSubstitutionTemplate, TokenTemplateHead:
			tl, err := h.parserTemplateLiteral(yield, await)
			if err != nil {
				return me, g.Error(err)
			}
			nme.TemplateLiteral = &tl
		case TokenIdentifier:
			h.Except()
			nme.IdentifierName = h.GetLastToken()
		case TokenPunctuator:
			if tk.Data != "[" {
				break Loop
			}
			h.Except()
			h.AcceptRunWhitespace()
			i := h.NewGoal()
			e, err := i.parseExpression(true, yield, await)
			if err != nil {
				return me, g.Error(err)
			}
			h.Score(i)
			h.AcceptRunWhitespace()
			if !h.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return me, g.Error(ErrMissingClosingBracket)
			}
			nme.Expression = &e
		default:
			break Loop
		}
		g.Score(h)
		ome := me
		ome.Tokens = j.ToTokens()
		nme.MemberExpression = &ome
		me = nme
		j.Score(g)
	}
	me.Tokens = j.ToTokens()
	return me, nil
}

type PrimaryExpression struct {
	This                                              bool
	IdentifierReference                               *IdentifierReference
	Literal                                           *Token
	ArrayLiteral                                      *ArrayLiteral
	ObjectLiteral                                     *ObjectLiteral
	FunctionExpression                                *FunctionDeclaration
	ClassExpression                                   *ClassDeclaration
	TemplateLiteral                                   *TemplateLiteral
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	Tokens                                            []Token
}

func (j *jsParser) parserPrimaryExpression(yield, await bool) (PrimaryExpression, error) {
	var pe PrimaryExpression
	if err := j.FindGoal(
		func(j *jsParser) error {
			if !j.AcceptToken(parser.Token{TokenKeyword, "this"}) {
				return errNotApplicable
			}
			pe.This = true
			return nil
		},
		func(j *jsParser) error {
			i, err := j.parseIdentifierReference(yield, await)
			if err != nil {
				return err
			}
			pe.IdentifierReference = &i
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
			al, err := j.parseArrayLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.ArrayLiteral = &al
			return nil
		},
		func(j *jsParser) error {
			ol, err := j.parseObjectLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.ObjectLiteral = &ol
			return nil
		},
		func(j *jsParser) error {
			fe, err := j.parseFunctionDeclaration(false, false, true)
			if err != nil {
				return err
			}
			pe.FunctionExpression = &fe
			return nil
		},
		func(j *jsParser) error {
			ce, err := j.parseClassDeclaration(yield, await, true)
			if err != nil {
				return err
			}
			pe.ClassExpression = &ce
			return nil
		},
		func(j *jsParser) error {
			tl, err := j.parserTemplateLiteral(yield, await)
			if err != nil {
				return err
			}
			pe.TemplateLiteral = &tl
			return nil
		},
		func(j *jsParser) error {
			cp, err := j.parseCoverParenthesizedExpressionAndArrowParameterList(yield, await)
			if err != nil {
				return err
			}
			pe.CoverParenthesizedExpressionAndArrowParameterList = &cp
			return nil
		},
	); err != nil {
		return pe, err
	}
	pe.Tokens = j.ToTokens()
	return pe, nil
}

type CoverParenthesizedExpressionAndArrowParameterList struct {
	Expressions          []Expression
	BindingIdentifier    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               []Token
}

func (j *jsParser) parseCoverParenthesizedExpressionAndArrowParameterList(yield, await bool) (CoverParenthesizedExpressionAndArrowParameterList, error) {
	var cp CoverParenthesizedExpressionAndArrowParameterList
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return cp, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		for {
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
					ab, err := j.parseArrayBindingPattern(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.ArrayBindingPattern = &ab

				} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
					ob, err := j.parseObjectBindingPattern(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.ObjectBindingPattern = &ob
				} else {
					bi, err := g.parseBindingIdentifier(yield, await)
					if err != nil {
						return cp, j.Error(err)
					}
					cp.BindingIdentifier = &bi
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
					return cp, j.Error(ErrMissingClosingParentheses)
				}
				break
			}
			g := j.NewGoal()
			e, err := g.parseExpression(true, yield, await)
			if err != nil {
				return cp, j.Error(err)
			}
			j.Score(g)
			cp.Expressions = append(cp.Expressions, e)
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return cp, j.Error(ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	cp.Tokens = j.ToTokens()
	return cp, nil
}

type Arguments struct {
	ArgumentList   []AssignmentExpression
	SpreadArgument *AssignmentExpression
	Tokens         []Token
}

func (j *jsParser) parseArguments(yield, await bool) (Arguments, error) {
	var a Arguments
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return a, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		for {
			var spread bool
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				spread = true
				j.AcceptRunWhitespace()
			}
			g := j.NewGoal()
			ae, err := g.parseAssignmentExpression(true, yield, await)
			if err != nil {
				return a, j.Error(err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if spread {
				a.SpreadArgument = &ae
				if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
					return a, j.Error(err)
				}
				break
			}
			a.ArgumentList = append(a.ArgumentList, ae)
			if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return a, j.Error(ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	a.Tokens = j.ToTokens()
	return a, nil
}

type CallExpression struct {
	CoverCallExpressionAndAsyncArrowHead *MemberExpression
	SuperCall                            bool
	CallExpression                       *CallExpression
	Arguments                            *Arguments
	Expression                           *Expression
	IdentifierName                       *Token
	TemplateLiteral                      *TemplateLiteral
	Tokens                               []Token
}

func (j *jsParser) parseCallExpression(yield, await bool) (CallExpression, error) {
	var ce CallExpression
	if j.AcceptToken(parser.Token{TokenKeyword, "super"}) {
		ce.SuperCall = true
	} else {
		g := j.NewGoal()
		me, err := g.parseMemberExpression(yield, await)
		if err != nil {
			return ce, j.Error(err)
		}
		ce.CoverCallExpressionAndAsyncArrowHead = &me
	}
	j.AcceptRunWhitespace()
	a, err := j.parseArguments(yield, await)
	if err != nil {
		return ce, err
	}
	ce.Arguments = &a
	for {
		oce := ce
		nce := CallExpression{
			CallExpression: &oce,
		}
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if tk := g.Peek(); tk == (parser.Token{TokenPunctuator, "("}) {
			h := g.NewGoal()
			a, err := h.parseArguments(yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			nce.Arguments = &a
		} else if tk == (parser.Token{TokenPunctuator, "["}) {
			g.Except()
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			e, err := h.parseExpression(true, yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			if !g.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return ce, j.Error(ErrMissingClosingBracket)
			}
			nce.Expression = &e
		} else if tk == (parser.Token{TokenPunctuator, "."}) {
			g.Except()
			g.AcceptRunWhitespace()
			if !g.Accept(TokenIdentifier, TokenKeyword) {
				return ce, j.Error(ErrNoIdentifier)
			}
			nce.IdentifierName = g.GetLastToken()
		} else if tk.Type == TokenTemplateHead || tk.Type == TokenNoSubstitutionTemplate {
			h := g.NewGoal()
			tl, err := h.parserTemplateLiteral(yield, await)
			if err != nil {
				return ce, j.Error(err)
			}
			g.Score(h)
			nce.TemplateLiteral = &tl
		} else {
			break
		}
		j.Score(g)
		nce.Tokens = j.ToTokens()
		ce = nce
	}
	ce.Tokens = j.ToTokens()
	return ce, nil
}
