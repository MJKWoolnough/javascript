package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

// AssignmentOperator specifies the type of assignment in AssignmentExpression
type AssignmentOperator uint8

// Valid AssignmentOperator's
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
	j.Skip()
	return nil
}

// AssignmentExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-AssignmentExpression
//
// It is only valid for one of ConditionalExpression, ArrowFunction,
// LeftHandSideExpression to be non-nil.
//
// If LeftHandSideExpression is non-nil, then AssignmentOperator must not be
// AssignmentNone and AssignmentExpression must be non-nil.
//
// If Yield is true, AssignmentExpression must be non-nil.
//
// If AssignmentOperator is AssignmentNone LeftHandSideExpression must be nil.
//
// If LeftHandSideExpression is nil and Yield is false, AssignmentExpression
// must be nil.
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
		if lhs := ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs != nil && len(ae.ConditionalExpression.Tokens) == len(lhs.Tokens) {
			h := g.NewGoal()
			if lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.CoverParenthesizedExpressionAndArrowParameterList != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil) {
				h.AcceptRunWhitespaceNoNewLine()
				if h.Peek() == (parser.Token{TokenPunctuator, "=>"}) {
					ae.ConditionalExpression = nil
					ae.ArrowFunction = new(ArrowFunction)
					if err := ae.ArrowFunction.parse(&g, lhs.NewExpression.MemberExpression.PrimaryExpression, in, yield, await); err != nil {
						return j.Error("AssignmentExpression", err)
					}
				}
			}
			if ae.ConditionalExpression != nil {
				h.AcceptRunWhitespace()
				if err := ae.AssignmentOperator.parse(&h); err == nil {
					g.Score(h)
					g.AcceptRunWhitespace()
					ae.ConditionalExpression = nil
					ae.LeftHandSideExpression = lhs
					h = g.NewGoal()
					ae.AssignmentExpression = new(AssignmentExpression)
					if err := ae.AssignmentExpression.parse(&h, in, yield, await); err != nil {
						return g.Error("AssignmentExpression", err)
					}
					g.Score(h)
				}
			}
		}
		j.Score(g)
	}
	ae.Tokens = j.ToTokens()
	return nil
}

// LeftHandSideExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-LeftHandSideExpression
//
// It is only valid for one of NewExpression or CallExpression to be non-nil.
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
		if lhs.NewExpression.News == 0 {
			h := g.NewGoal()
			h.AcceptRunWhitespace()
			if h.Peek() == (parser.Token{TokenPunctuator, "("}) {
				lhs.CallExpression = new(CallExpression)
				if err := lhs.CallExpression.parse(&g, &lhs.NewExpression.MemberExpression, yield, await); err != nil {
					return j.Error("LeftHandSideExpression", err)
				}
				lhs.NewExpression = nil
			}
		}
		j.Score(g)
	}
	lhs.Tokens = j.ToTokens()
	return nil
}

// Expression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-Expression
//
// Expressions must have a length of at least one to be valid.
type Expression struct {
	Expressions []AssignmentExpression
	Tokens      Tokens
}

func (e *Expression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		ae := len(e.Expressions)
		e.Expressions = append(e.Expressions, AssignmentExpression{})
		if err := e.Expressions[ae].parse(&g, in, yield, await); err != nil {
			return j.Error("Expression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			break
		}
		g.AcceptRunWhitespace()
		j.Score(g)
	}
	e.Tokens = j.ToTokens()
	return nil
}

// NewExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-NewExpression
//
// The News field is a count of the number of 'new' keywords that proceed the
// MemberExpression
type NewExpression struct {
	News             uint
	MemberExpression MemberExpression
	Tokens           Tokens
}

func (ne *NewExpression) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if err := ne.MemberExpression.parse(&g, yield, await); err != nil {
		h := j.NewGoal()
		for {
			if ne.MemberExpression.MemberExpression == nil || !h.AcceptToken(parser.Token{TokenKeyword, "new"}) {
				return j.Error("NewExpression", err)
			}
			ne.MemberExpression = *ne.MemberExpression.MemberExpression
			ne.News++
			if ne.MemberExpression.Tokens != nil {
				break
			}
			h.AcceptRunWhitespace()
		}
	}
	j.Score(g)
	ne.Tokens = j.ToTokens()
	return nil
}

// MemberExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-MemberExpression
//
// If PrimaryExpression is nil, SuperProperty is true, or MetaProperty = is,
// Expression, IdentifierName, TemplateLiteral, and Arguments must be nil.
//
// If Expression, IdentifierName, TemplateLiteral, or Arguments is non-nil,
// then MemberExpression must be non-nil.
//
// It is only valid if one of Expression, IdentifierName, TemplateLiteral, and
// Arguments is non-nil.
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
			g.AcceptRunWhitespace()
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
				h.backup()
				err = j.Error("MemberExpression", err)
				g.Score(h)
				j.Score(g)
				return err
			}
			g.Score(h)
			h = g.NewGoal()
			h.AcceptRunWhitespace()
			i := h.NewGoal()
			me.Arguments = new(Arguments)
			if err := me.Arguments.parse(&i, yield, await); err != nil {
				j.Score(g)
				return g.Error("MemberExpression", err)
			}
			h.Score(i)
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
				h.Skip()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("MemberExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Skip()
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

// PrimaryExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-PrimaryExpression
//
// It is only valid is one IdentifierReference, Literal, ArrayLiteral,
// ObjectLiteral, FunctionExpression, ClassExpression, TemplateLiteral, or
// CoverParenthesizedExpressionAndArrowParameterList is non-nil or This is true.
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

// CoverParenthesizedExpressionAndArrowParameterList as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-CoverParenthesizedExpressionAndArrowParameterList
//
// It is valid for only one of BindingIdentifier, ArrayBindingPattern, and
// ObjectBindingPattern to be non-nil
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
				if t := g.Peek(); t == (parser.Token{TokenPunctuator, "["}) {
					cp.ArrayBindingPattern = new(ArrayBindingPattern)
					if err := cp.ArrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if t == (parser.Token{TokenPunctuator, "{"}) {
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
			e := len(cp.Expressions)
			cp.Expressions = append(cp.Expressions, AssignmentExpression{})
			if err := cp.Expressions[e].parse(&g, true, yield, await); err != nil {
				return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
			}
			j.Score(g)
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

// Arguments as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-Arguments
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
			h := g.NewGoal()
			a.SpreadArgument = new(AssignmentExpression)
			if err := a.SpreadArgument.parse(&h, true, yield, await); err != nil {
				return j.Error("Arguments", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return j.Error("Arguments", ErrMissingClosingParenthesis)
			}
			break
		}
		ae := len(a.ArgumentList)
		a.ArgumentList = append(a.ArgumentList, AssignmentExpression{})
		if err := a.ArgumentList[ae].parse(&g, true, yield, await); err != nil {
			return j.Error("Arguments", err)
		}
		j.Score(g)
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

// CallExpression as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-CallExpression
//
// Includes the TC39 proposal for the dynamic import function call
// https://github.com/tc39/proposal-dynamic-import/#import
//
// It is only valid for one of MemberExpression, ImportCall, or CallExpression
// to be non-nil or SuperCall to be true.
//
// If MemberExpression is non-nil, or SuperCall is true, Arguments must be
// non-nil.
//
// If CallExpression is non-nil, only one of Arguments, Expression,
// IdentifierName, or TemplateLiteral must be non-nil.
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
			j.AcceptRunWhitespace()
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
		} else {
			return j.Error("CallExpression", ErrInvalidCallExpression)
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
				h.Skip()
				h.AcceptRunWhitespace()
				if !h.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("CallExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
			case "[":
				h.Skip()
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

// Errors
var (
	ErrInvalidCallExpression = errors.Error("invalid CallExpression")
)
