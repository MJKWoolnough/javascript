package javascript

import (
	"errors"

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
	case parser.Token{Type: TokenPunctuator, Data: "="}:
		*ao = AssignmentAssign
	case parser.Token{Type: TokenPunctuator, Data: "*="}:
		*ao = AssignmentMultiply
	case parser.Token{Type: TokenDivPunctuator, Data: "/="}:
		*ao = AssignmentDivide
	case parser.Token{Type: TokenPunctuator, Data: "%="}:
		*ao = AssignmentRemainder
	case parser.Token{Type: TokenPunctuator, Data: "+="}:
		*ao = AssignmentAdd
	case parser.Token{Type: TokenPunctuator, Data: "-="}:
		*ao = AssignmentSubtract
	case parser.Token{Type: TokenPunctuator, Data: "<<="}:
		*ao = AssignmentLeftShift
	case parser.Token{Type: TokenPunctuator, Data: ">>="}:
		*ao = AssignmentSignPropagatinRightShift
	case parser.Token{Type: TokenPunctuator, Data: ">>>="}:
		*ao = AssignmentZeroFillRightShift
	case parser.Token{Type: TokenPunctuator, Data: "&="}:
		*ao = AssignmentBitwiseAND
	case parser.Token{Type: TokenPunctuator, Data: "^="}:
		*ao = AssignmentBitwiseXOR
	case parser.Token{Type: TokenPunctuator, Data: "|="}:
		*ao = AssignmentBitwiseOR
	case parser.Token{Type: TokenPunctuator, Data: "**="}:
		*ao = AssignmentExponentiation
	default:
		return ErrInvalidAssignment
	}
	j.Skip()
	return nil
}

// AssignmentExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AssignmentExpression
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
	done := false
	if yield && j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "yield"}) {
		ae.Yield = true
		j.AcceptRunWhitespaceNoNewLine()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
			ae.Delegate = true
			j.AcceptRunWhitespace()
		}
		g := j.NewGoal()
		ae.AssignmentExpression = new(AssignmentExpression)
		if err := ae.AssignmentExpression.parse(&g, in, true, await); err != nil {
			return j.Error("AssignmentExpression", err)
		}
		j.Score(g)
		done = true
	} else if j.Peek() == (parser.Token{Type: TokenIdentifier, Data: "async"}) {
		g := j.NewGoal()
		g.Skip()
		g.AcceptRunWhitespaceNoNewLine()
		if t := g.Peek().Type; t == TokenPunctuator || t == TokenIdentifier {
			g := j.NewGoal()
			ae.ArrowFunction = new(ArrowFunction)
			if err := ae.ArrowFunction.parse(&g, nil, in, yield, await); err != nil {
				return j.Error("AssignmentExpression", err)
			}
			j.Score(g)
			done = true
		}
	}
	if !done {
		g := j.NewGoal()
		ae.ConditionalExpression = new(ConditionalExpression)
		if err := ae.ConditionalExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("AssignmentExpression", err)
		}
		if ae.ConditionalExpression.LogicalORExpression != nil {
			if lhs := ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs != nil && len(ae.ConditionalExpression.Tokens) == len(lhs.Tokens) {
				h := g.NewGoal()
				if lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.CoverParenthesizedExpressionAndArrowParameterList != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil) {
					h.AcceptRunWhitespaceNoNewLine()
					if h.Peek() == (parser.Token{Type: TokenPunctuator, Data: "=>"}) {
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
		}
		j.Score(g)
	}
	ae.Tokens = j.ToTokens()
	return nil
}

// LeftHandSideExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-LeftHandSideExpression
//
// It is only valid for one of NewExpression, CallExpression or
// OptionalExpression to be non-nil.
//
// Includes OptionalExpression as per TC39 (2020-03)
type LeftHandSideExpression struct {
	NewExpression      *NewExpression
	CallExpression     *CallExpression
	OptionalExpression *OptionalExpression
	Tokens             Tokens
}

func (lhs *LeftHandSideExpression) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "super"}) || g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "import"}) {
		g.AcceptRunWhitespace()
		if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
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
			if h.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
				lhs.CallExpression = new(CallExpression)
				if err := lhs.CallExpression.parse(&g, &lhs.NewExpression.MemberExpression, yield, await); err != nil {
					return j.Error("LeftHandSideExpression", err)
				}
				lhs.NewExpression = nil
			}
		}
		j.Score(g)
	}
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "?."}) {
		if lhs.CallExpression != nil {
			g = j.NewGoal()
			j.AcceptRunWhitespace()
			lhs.OptionalExpression = new(OptionalExpression)
			if err := lhs.OptionalExpression.parse(j, yield, await, nil, lhs.CallExpression); err != nil {
				return g.Error("LeftHandSideExpression", err)
			}
			lhs.CallExpression = nil
		} else if lhs.NewExpression.News == 0 {
			g = j.NewGoal()
			j.AcceptRunWhitespace()
			lhs.OptionalExpression = new(OptionalExpression)
			if err := lhs.OptionalExpression.parse(j, yield, await, &lhs.NewExpression.MemberExpression, nil); err != nil {
				return g.Error("LeftHandSideExpression", err)
			}
			lhs.NewExpression = nil
		}
	}
	lhs.Tokens = j.ToTokens()
	return nil
}

// OptionalExpression as defined in TC39
// https://tc39.es/ecma262/#prod-OptionalExpression
//
// It is only valid for one of NewExpression, CallExpression or
// OptionalExpression to be non-nil.
type OptionalExpression struct {
	MemberExpression   *MemberExpression
	CallExpression     *CallExpression
	OptionalExpression *OptionalExpression
	OptionalChain      OptionalChain
	Tokens             Tokens
}

func (oe *OptionalExpression) parse(j *jsParser, yield, await bool, me *MemberExpression, ce *CallExpression) error {
	if me != nil {
		oe.MemberExpression = me
	} else {
		oe.CallExpression = ce
	}
	g := j.NewGoal()
	if err := oe.OptionalChain.parse(&g, yield, await); err != nil {
		return j.Error("OptionalExpression", err)
	}
	for {
		j.Score(g)
		oe.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if g.Peek() != (parser.Token{Type: TokenPunctuator, Data: "?."}) {
			break
		}
		noe := new(OptionalExpression)
		*noe = *oe
		*oe = OptionalExpression{
			OptionalExpression: noe,
		}
		h := g.NewGoal()
		if err := oe.OptionalChain.parse(&h, yield, await); err != nil {
			return g.Error("OptionalExpression", err)
		}
		g.Score(h)
	}
	return nil
}

// OptionalChain as defined in TC39
// https://tc39.es/ecma262/#prod-OptionalExpression
//
// It is only valid for one of Arguments, Expression, IdentifierName, or
// TemplateLiteral to be non-nil.
type OptionalChain struct {
	OptionalChain   *OptionalChain
	Arguments       *Arguments
	Expression      *Expression
	IdentifierName  *Token
	TemplateLiteral *TemplateLiteral
	Tokens          Tokens
}

func (oc *OptionalChain) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?."}) {
		return j.Error("OptionalChain", ErrMissingOptional)
	}
	j.AcceptRunWhitespace()
	if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
		g := j.NewGoal()
		oc.Arguments = new(Arguments)
		if err := oc.Arguments.parse(&g, yield, await); err != nil {
			return j.Error("OptionalChain", err)
		}
		j.Score(g)
	} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		oc.Expression = new(Expression)
		if err := oc.Expression.parse(&g, true, yield, await); err != nil {
			return j.Error("OptionalChain", err)
		}
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return g.Error("OptionalChain", ErrMissingClosingBracket)
		}
		j.Score(g)
	} else if j.Accept(TokenIdentifier, TokenKeyword) {
		oc.IdentifierName = j.GetLastToken()
	} else if t := j.Peek().Type; t == TokenNoSubstitutionTemplate || t == TokenTemplateHead {
		g := j.NewGoal()
		oc.TemplateLiteral = new(TemplateLiteral)
		if err := oc.TemplateLiteral.parse(&g, yield, await); err != nil {
			return j.Error("OptionalChain", err)
		}
		j.Score(g)
	} else {
		return j.Error("OptionalChain", ErrInvalidOptionalChain)
	}
	for {
		oc.Tokens = j.ToTokens()
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		var (
			arguments       *Arguments
			expression      *Expression
			identifierName  *Token
			templateLiteral *TemplateLiteral
		)
		if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
			h := g.NewGoal()
			arguments = new(Arguments)
			if err := arguments.parse(&h, yield, await); err != nil {
				return g.Error("OptionalChain", err)
			}
			g.Score(h)
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			expression = new(Expression)
			if err := expression.parse(&h, true, yield, await); err != nil {
				return g.Error("OptionalChain", err)
			}
			h.AcceptRunWhitespace()
			if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return h.Error("OptionalChain", ErrMissingClosingBracket)
			}
			g.Score(h)
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			g.AcceptRunWhitespace()
			if !g.Accept(TokenIdentifier, TokenKeyword) {
				return g.Error("OptionalChain", ErrNoIdentifier)
			}
			identifierName = g.GetLastToken()
		} else if t := g.Peek().Type; t == TokenNoSubstitutionTemplate || t == TokenTemplateHead {
			h := g.NewGoal()
			templateLiteral = new(TemplateLiteral)
			if err := templateLiteral.parse(&h, yield, await); err != nil {
				return g.Error("OptionalChain", err)
			}
			g.Score(h)
		} else {
			break
		}
		noc := new(OptionalChain)
		*noc = *oc
		*oc = OptionalChain{
			Arguments:       arguments,
			Expression:      expression,
			IdentifierName:  identifierName,
			TemplateLiteral: templateLiteral,
			OptionalChain:   noc,
		}
		j.Score(g)
	}
	return nil
}

// Expression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-Expression
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
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			break
		}
		g.AcceptRunWhitespace()
		j.Score(g)
	}
	e.Tokens = j.ToTokens()
	return nil
}

// NewExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-NewExpression
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
			if ne.MemberExpression.MemberExpression == nil || !h.AcceptToken(parser.Token{Type: TokenKeyword, Data: "new"}) {
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
// https://262.ecma-international.org/11.0/#prod-MemberExpression
//
// If PrimaryExpression is nil, SuperProperty is true, NewTarget is true, or
// ImportMeta is true, Expression, IdentifierName, TemplateLiteral, and
// Arguments must be nil.
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
	NewTarget         bool
	ImportMeta        bool
	Arguments         *Arguments
	Tokens            Tokens
}

func (me *MemberExpression) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	e := false
	if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "super"}) {
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			me.Expression = new(Expression)
			if err := me.Expression.parse(&h, true, yield, await); err != nil {
				return g.Error("MemberExpression", err)
			}
			g.Score(h)
			g.AcceptRunWhitespace()
			if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return g.Error("MemberExpression", ErrInvalidSuperProperty)
			}
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			g.AcceptRunWhitespace()
			if !g.Accept(TokenIdentifier, TokenKeyword) {
				return g.Error("MemberExpression", ErrNoIdentifier)
			}
			me.IdentifierName = g.GetLastToken()
		} else {
			return g.Error("MemberExpression", ErrInvalidSuperProperty)
		}
		me.SuperProperty = true
	} else if g.Peek() == (parser.Token{Type: TokenKeyword, Data: "new"}) || g.Peek() == (parser.Token{Type: TokenKeyword, Data: "import"}) {
		var isNew bool
		if g.Peek().Data == "new" {
			isNew = true
		}
		g.Skip()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			g.AcceptRunWhitespace()
			var id string
			if isNew {
				id = "target"
				me.NewTarget = true
			} else {
				id = "meta"
				me.ImportMeta = true
			}
			if !g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: id}) {
				return j.Error("MemberExpression", ErrInvalidMetaProperty)
			}
		} else if isNew {
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
		} else {
			g = j.NewGoal()
			e = true
		}
	} else {
		e = true
	}
	if e {
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
				if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
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
// https://262.ecma-international.org/11.0/#prod-PrimaryExpression
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
	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"}) {
		pe.This = true
	} else if j.Accept(TokenNullLiteral, TokenBooleanLiteral, TokenNumericLiteral, TokenStringLiteral, TokenRegularExpressionLiteral) {
		pe.Literal = j.GetLastToken()
	} else if t := j.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
		g := j.NewGoal()
		pe.ArrayLiteral = new(ArrayLiteral)
		if err := pe.ArrayLiteral.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		g := j.NewGoal()
		pe.ObjectLiteral = new(ObjectLiteral)
		if err := pe.ObjectLiteral.parse(&g, yield, await); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{Type: TokenIdentifier, Data: "async"}) || t == (parser.Token{Type: TokenKeyword, Data: "function"}) {
		g := j.NewGoal()
		pe.FunctionExpression = new(FunctionDeclaration)
		if err := pe.FunctionExpression.parse(&g, false, false, true); err != nil {
			return j.Error("PrimaryExpression", err)
		}
		j.Score(g)
	} else if t == (parser.Token{Type: TokenKeyword, Data: "class"}) {
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
	} else if t == (parser.Token{Type: TokenPunctuator, Data: "("}) {
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
// https://262.ecma-international.org/11.0/#prod-CoverParenthesizedExpressionAndArrowParameterList
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
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		for {
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if t := g.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
					cp.ArrayBindingPattern = new(ArrayBindingPattern)
					if err := cp.ArrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
					cp.ObjectBindingPattern = new(ObjectBindingPattern)
					if err := cp.ObjectBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("CoverParenthesizedExpressionAndArrowParameterList", err)
					}
				} else if cp.BindingIdentifier = g.parseIdentifier(yield, await); cp.BindingIdentifier == nil {
					return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrNoIdentifier)
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
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
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return j.Error("CoverParenthesizedExpressionAndArrowParameterList", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	cp.Tokens = j.ToTokens()
	return nil
}

// Arguments as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-Arguments
type Arguments struct {
	ArgumentList   []AssignmentExpression
	SpreadArgument *AssignmentExpression
	Tokens         Tokens
}

func (a *Arguments) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("Arguments", ErrMissingOpeningParenthesis)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			break
		}
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			a.SpreadArgument = new(AssignmentExpression)
			if err := a.SpreadArgument.parse(&h, true, yield, await); err != nil {
				return j.Error("Arguments", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
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
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return j.Error("Arguments", ErrMissingComma)
		}
	}
	a.Tokens = j.ToTokens()
	return nil
}

// CallExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-CallExpression
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
	skip := false
	if me == nil {
		if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "super"}) {
			ce.SuperCall = true
		} else if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "import"}) {
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
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
			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				return j.Error("CallExpression", ErrMissingClosingParenthesis)
			}
			skip = true
		} else {
			return j.Error("CallExpression", ErrInvalidCallExpression)
		}
	} else {
		ce.MemberExpression = me
	}
	if !skip {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		ce.Arguments = new(Arguments)
		if err := ce.Arguments.parse(&g, yield, await); err != nil {
			return j.Error("CallExpression", err)
		}
		j.Score(g)
	}
	for {
		ce.Tokens = j.ToTokens()
		g := j.NewGoal()
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
					return h.Error("CallExpression", ErrNoIdentifier)
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
				if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
					return h.Error("CallExpression", ErrMissingClosingBracket)
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
	ErrInvalidCallExpression = errors.New("invalid CallExpression")
	ErrMissingOptional       = errors.New("missing optional chain punctuator")
	ErrInvalidOptionalChain  = errors.New("invalid OptionalChain")
)
