package javascript

import "vimagination.zapto.org/parser"

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
	AssignmentSignPropagatingRightShift
	AssignmentZeroFillRightShift
	AssignmentBitwiseAND
	AssignmentBitwiseXOR
	AssignmentBitwiseOR
	AssignmentExponentiation
	AssignmentLogicalAnd
	AssignmentLogicalOr
	AssignmentNullish
)

func (ao *AssignmentOperator) parse(j *jsParser) error {
	g := j.NewGoal()
	h := j.NewGoal()
	i := j.NewGoal()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		*ao = AssignmentAssign
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*="}) {
		*ao = AssignmentMultiply
	} else if g.AcceptToken(parser.Token{Type: TokenDivPunctuator, Data: "/="}) {
		*ao = AssignmentDivide
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "%="}) {
		*ao = AssignmentRemainder
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "+="}) {
		*ao = AssignmentAdd
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "-="}) {
		*ao = AssignmentSubtract
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<<="}) {
		*ao = AssignmentLeftShift
	} else if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		*ao = AssignmentSignPropagatingRightShift
		g.Score(h)
	} else if i.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && i.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && i.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && i.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		*ao = AssignmentZeroFillRightShift
		g.Score(i)
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&="}) {
		*ao = AssignmentBitwiseAND
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^="}) {
		*ao = AssignmentBitwiseXOR
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|="}) {
		*ao = AssignmentBitwiseOR
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "**="}) {
		*ao = AssignmentExponentiation
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&="}) {
		*ao = AssignmentLogicalAnd
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||="}) {
		*ao = AssignmentLogicalOr
	} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "??="}) {
		*ao = AssignmentNullish
	} else {
		return ErrInvalidAssignment
	}
	j.Score(g)
	return nil
}

// AssignmentExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AssignmentExpression
//
// It is only valid for one of ConditionalExpression, ArrowFunction,
// LeftHandSideExpression, and AssignmentPattern to be non-nil.
//
// If LeftHandSideExpression, or AssignmentPattern are non-nil, then
// AssignmentOperator must not be AssignmentNone and AssignmentExpression must
// be non-nil.
//
// If LeftHandSideArray, or LeftHandSideObject are non-nil, AssignmentOperator
// must be AssignmentAssign.
//
// If Yield is true, AssignmentExpression must be non-nil.
//
// It is only valid for Delagate to be true if Yield is also true.
//
// If AssignmentOperator is AssignmentNone LeftHandSideExpression must be nil.
//
// If LeftHandSideExpression, and AssignmentPattern are nil and Yield is false,
// AssignmentExpression must be nil.
type AssignmentExpression struct {
	ConditionalExpression  *ConditionalExpression
	ArrowFunction          *ArrowFunction
	LeftHandSideExpression *LeftHandSideExpression
	AssignmentPattern      *AssignmentPattern
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
		if g.SkipGeneric() {
			g.AcceptRunWhitespaceNoNewLine()
		}
		if t := g.Peek().Type; t == TokenPunctuator || t == TokenIdentifier {
			g := j.NewGoal()
			ae.ArrowFunction = new(ArrowFunction)
			if err := ae.ArrowFunction.parse(&g, nil, in, yield, await); err != nil {
				return j.Error("AssignmentExpression", err)
			}
			j.Score(g)
			done = true
		}
	} else {
		g := j.NewGoal()
		if g.SkipGeneric() {
			g.AcceptRunWhitespaceNoNewLine()
			h := g.NewGoal()
			var pe ParenthesizedExpression
			if pe.parse(&h, yield, await) == nil {
				g.Score(h)
				ae.ArrowFunction = new(ArrowFunction)
				if err := ae.ArrowFunction.parse(&g, &PrimaryExpression{
					ParenthesizedExpression: &pe,
					Tokens:                  pe.Tokens,
				}, in, yield, await); err != nil {
					return g.Error("AssignmentExpression", err)
				}
				j.Score(g)
				done = true
			}
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
				if lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference != nil) {
					h.AcceptRunWhitespaceNoNewLine()
					if lhs.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression != nil && h.SkipReturnType() {
						h.AcceptRunWhitespaceNoNewLine()
					}
					if h.Peek() == (parser.Token{Type: TokenPunctuator, Data: "=>"}) {
						ae.ConditionalExpression = nil
						ae.ArrowFunction = new(ArrowFunction)
						if err := ae.ArrowFunction.parse(&g, lhs.NewExpression.MemberExpression.PrimaryExpression, in, yield, await); err != nil {
							return j.Error("AssignmentExpression", err)
						}
					} else if cpe := lhs.NewExpression.MemberExpression.PrimaryExpression.ParenthesizedExpression; cpe != nil && (cpe.bindingIdentifier != nil || cpe.arrayBindingPattern != nil || cpe.objectBindingPattern != nil) {
						return h.Error("AssignmentExpression", ErrMissingArrow)
					}
				}
				if ae.ConditionalExpression != nil {
					h.AcceptRunWhitespace()
					if err := ae.AssignmentOperator.parse(&h); err == nil {
						g.Score(h)
						g.AcceptRunWhitespace()
						ae.ConditionalExpression = nil
						ae.LeftHandSideExpression = lhs
						if ae.AssignmentOperator == AssignmentAssign && lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral != nil) {
							ae.AssignmentPattern = new(AssignmentPattern)
							if err := ae.AssignmentPattern.from(lhs.NewExpression.MemberExpression.PrimaryExpression); err != nil {
								z := jsParser(lhs.Tokens[:0])
								return z.Error("AssignmentExpression", err)
							}
							ae.LeftHandSideExpression = nil
						}
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
			if h.SkipTypeArguments() {
				h.AcceptRunWhitespace()
			}
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

// IsSimple returns whether or not the LeftHandSideExpression is classed as 'simple'
func (lhs *LeftHandSideExpression) IsSimple() bool {
	return lhs.OptionalExpression == nil && (lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.IsSimple() || lhs.CallExpression != nil && lhs.CallExpression.IsSimple())
}

// AssignmentPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AssignmentPattern
//
// Only one of ObjectAssignmentPattern or ArrayAssignmentPattern must be
// non-nil
type AssignmentPattern struct {
	ObjectAssignmentPattern *ObjectAssignmentPattern
	ArrayAssignmentPattern  *ArrayAssignmentPattern
	Tokens                  Tokens
}

func (a *AssignmentPattern) from(p *PrimaryExpression) error {
	if p.ArrayLiteral != nil {
		a.ArrayAssignmentPattern = new(ArrayAssignmentPattern)
		if err := a.ArrayAssignmentPattern.from(p.ArrayLiteral); err != nil {
			z := jsParser(p.ArrayLiteral.Tokens[:0])
			return z.Error("AssignmentPattern", err)
		}
	} else {
		a.ObjectAssignmentPattern = new(ObjectAssignmentPattern)
		if err := a.ObjectAssignmentPattern.from(p.ObjectLiteral); err != nil {
			z := jsParser(p.ObjectLiteral.Tokens[:0])
			return z.Error("AssignmentPattern", err)
		}
	}
	a.Tokens = p.Tokens
	return nil
}

// ObjectAssignmentPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ObjectAssignmentPattern
type ObjectAssignmentPattern struct {
	AssignmentPropertyList []AssignmentProperty
	AssignmentRestElement  *LeftHandSideExpression
	Tokens                 Tokens
}

func (o *ObjectAssignmentPattern) from(ol *ObjectLiteral) error {
	o.AssignmentPropertyList = make([]AssignmentProperty, len(ol.PropertyDefinitionList))
	for n := range ol.PropertyDefinitionList {
		pd := &ol.PropertyDefinitionList[n]
		if pd.PropertyName == nil && pd.AssignmentExpression != nil {
			if n == len(ol.PropertyDefinitionList)-1 {
				o.AssignmentPropertyList = o.AssignmentPropertyList[:n:n]
				var dat DestructuringAssignmentTarget
				if pd.AssignmentExpression.AssignmentOperator != AssignmentNone {
					z := jsParser(pd.AssignmentExpression.Tokens[:0])
					return z.Error("ObjectAssignmentPattern", ErrInvalidAssignment)
				}
				if err := dat.from(pd.AssignmentExpression); err != nil {
					z := jsParser(pd.AssignmentExpression.Tokens[:0])
					return z.Error("ObjectAssignmentPattern", err)
				}
				if dat.AssignmentPattern != nil {
					z := jsParser(dat.Tokens[:0])
					return z.Error("ObjectAssignmentPattern", ErrBadRestElement)
				}
				o.AssignmentRestElement = dat.LeftHandSideExpression
				break
			}
		}
		if err := o.AssignmentPropertyList[n].from(pd); err != nil {
			z := jsParser(pd.Tokens[:0])
			return z.Error("ObjectAssignmentPattern", err)
		}
	}
	o.Tokens = ol.Tokens
	return nil
}

// AssignmentProperty as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AssignmentProperty
type AssignmentProperty struct {
	PropertyName                  PropertyName
	DestructuringAssignmentTarget *DestructuringAssignmentTarget
	Initializer                   *AssignmentExpression
	Tokens                        Tokens
}

func (a *AssignmentProperty) from(pd *PropertyDefinition) error {
	if pd.MethodDefinition != nil || pd.PropertyName == nil {
		z := jsParser(pd.Tokens[:0])
		return z.Error("AssignmentProperty", ErrInvalidAssignmentProperty)
	}
	if pd.PropertyName.LiteralPropertyName == nil {
		z := jsParser(pd.Tokens[:0])
		return z.Error("AssignmentProperty", z.Error("PropertyName", ErrNotSimple))
	}
	a.PropertyName = *pd.PropertyName
	if pd.IsCoverInitializedName {
		a.Initializer = pd.AssignmentExpression
	} else {
		a.DestructuringAssignmentTarget = new(DestructuringAssignmentTarget)
		if err := a.DestructuringAssignmentTarget.from(pd.AssignmentExpression); err != nil {
			z := jsParser(pd.Tokens[:0])
			return z.Error("AssignmentProperty", err)
		}
		a.Initializer = pd.AssignmentExpression.AssignmentExpression
	}
	a.Tokens = pd.Tokens
	return nil
}

// DestructuringAssignmentTarget as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-DestructuringAssignmentTarget
//
// Only one of LeftHandSideExpression or AssignmentPattern must be non-nil
type DestructuringAssignmentTarget struct {
	LeftHandSideExpression *LeftHandSideExpression
	AssignmentPattern      *AssignmentPattern
	Tokens                 Tokens
}

func (d *DestructuringAssignmentTarget) from(ae *AssignmentExpression) error {
	if ae.LeftHandSideExpression != nil {
		d.LeftHandSideExpression = ae.LeftHandSideExpression
		d.Tokens = ae.LeftHandSideExpression.Tokens
	} else if ae.AssignmentPattern != nil {
		d.AssignmentPattern = ae.AssignmentPattern
		d.Tokens = ae.AssignmentPattern.Tokens
	} else if ae.ConditionalExpression == nil {
		z := jsParser(ae.Tokens[:0])
		return z.Error("DestructuringAssignmentTarget", ErrInvalidDestructuringAssignmentTarget)
	} else {
		switch UnwrapConditional(ae.ConditionalExpression).(type) {
		case *ArrayLiteral, *ObjectLiteral:
			d.AssignmentPattern = new(AssignmentPattern)
			if err := d.AssignmentPattern.from(ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression); err != nil {
				z := jsParser(ae.Tokens[:0])
				return z.Error("DestructuringAssignmentTarget", err)
			}
		case *CallExpression, *MemberExpression, *PrimaryExpression:
			d.LeftHandSideExpression = ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression
			if !d.LeftHandSideExpression.IsSimple() {
				z := jsParser(ae.Tokens[:0])
				return z.Error("DestructuringAssignmentTarget", ErrInvalidDestructuringAssignmentTarget)
			}
		default:
			z := jsParser(ae.Tokens[:0])
			return z.Error("DestructuringAssignmentTarget", ErrInvalidDestructuringAssignmentTarget)
		}
		d.Tokens = ae.ConditionalExpression.Tokens
	}
	return nil
}

// AssignmentElement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AssignmentElement
type AssignmentElement struct {
	DestructuringAssignmentTarget DestructuringAssignmentTarget
	Initializer                   *AssignmentExpression
	Tokens                        Tokens
}

func (a *AssignmentElement) from(ae *AssignmentExpression) error {
	switch ae.AssignmentOperator {
	case AssignmentNone, AssignmentAssign:
		if err := a.DestructuringAssignmentTarget.from(ae); err != nil {
			z := jsParser(ae.Tokens[:0])
			return z.Error("AssignmentElement", err)
		}
		a.Initializer = ae.AssignmentExpression
	default:
		z := jsParser(ae.Tokens[:0])
		return z.Error("AssignmentElement", ErrInvalidAssignment)
	}
	a.Tokens = ae.Tokens
	return nil
}

// ArrayAssignmentPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayAssignmentPattern
type ArrayAssignmentPattern struct {
	AssignmentElements    []AssignmentElement
	AssignmentRestElement *LeftHandSideExpression
	Tokens                Tokens
}

func (a *ArrayAssignmentPattern) from(al *ArrayLiteral) error {
	a.AssignmentElements = make([]AssignmentElement, 0, len(al.ElementList))
	hasSpread := false
	for _, ae := range al.ElementList {
		if hasSpread {
			z := jsParser(al.Tokens[:0])
			return z.Error("ArrayAssignmentPattern", ErrBadRestElement)
		} else if ae.Spread {
			hasSpread = true
			var dat DestructuringAssignmentTarget
			if ae.AssignmentExpression.AssignmentOperator != AssignmentNone {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayAssignmentPattern", ErrInvalidAssignment)
			}
			if err := dat.from(&ae.AssignmentExpression); err != nil {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayAssignmentPattern", err)
			}
			if dat.AssignmentPattern != nil {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayAssignmentPattern", ErrBadRestElement)
			}
			a.AssignmentRestElement = dat.LeftHandSideExpression
		} else if len(ae.Tokens) > 0 {
			var e AssignmentElement
			if err := e.from(&ae.AssignmentExpression); err != nil {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayAssignmentPattern", err)
			}
			a.AssignmentElements = append(a.AssignmentElements, e)
		} else {
			a.AssignmentElements = append(a.AssignmentElements, AssignmentElement{})
		}
	}
	a.Tokens = al.Tokens
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
// It is only valid for one of Arguments, Expression, IdentifierName,
// TemplateLiteral, or PrivateIdentifier to be non-nil.
type OptionalChain struct {
	OptionalChain     *OptionalChain
	Arguments         *Arguments
	Expression        *Expression
	IdentifierName    *Token
	TemplateLiteral   *TemplateLiteral
	PrivateIdentifier *Token
	Tokens            Tokens
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
	} else if j.Accept(TokenPrivateIdentifier) {
		oc.PrivateIdentifier = j.GetLastToken()
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
			arguments         *Arguments
			expression        *Expression
			identifierName    *Token
			templateLiteral   *TemplateLiteral
			privateIdentifier *Token
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
			if g.Accept(TokenPrivateIdentifier) {
				privateIdentifier = g.GetLastToken()
			} else {
				if !g.Accept(TokenIdentifier, TokenKeyword) {
					return g.Error("OptionalChain", ErrNoIdentifier)
				}
				identifierName = g.GetLastToken()
			}
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
			Arguments:         arguments,
			Expression:        expression,
			IdentifierName:    identifierName,
			TemplateLiteral:   templateLiteral,
			PrivateIdentifier: privateIdentifier,
			OptionalChain:     noc,
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
// https://tc39.es/ecma262/#prod-MemberExpression
//
// If PrimaryExpression is nil, SuperProperty is true, NewTarget is true, or
// ImportMeta is true, Expression, IdentifierName, TemplateLiteral, Arguments
// and PrivateIdentifier must be nil.
//
// If Expression, IdentifierName, TemplateLiteral, Arguments, or
// PrivateIdentifier is non-nil, then MemberExpression must be non-nil.
//
// It is only valid if one of Expression, IdentifierName, TemplateLiteral,
// Arguments, and PrivateIdentifier is non-nil.
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
	PrivateIdentifier *Token
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
				err = j.Error("MemberExpression", err)
				g.Score(h)
				j.Score(g)
				return err
			}
			g.Score(h)
			h = g.NewGoal()
			h.AcceptRunWhitespace()
			if h.SkipTypeArguments() {
				h.AcceptRunWhitespace()
			}
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
		h := g.NewGoal()
		h.AcceptRunWhitespaceNoNewLine()
		if h.SkipForce() {
			g.Score(h)
		}
	}
	j.Score(g)
	for {
		me.Tokens = j.ToTokens()
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		var (
			tl   *TemplateLiteral
			i, p *Token
			e    *Expression
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
				if !h.Accept(TokenIdentifier, TokenKeyword, TokenPrivateIdentifier) {
					return g.Error("MemberExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
				if i.Type == TokenPrivateIdentifier {
					p = i
					i = nil
				}
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
			i := h.NewGoal()
			i.AcceptRunWhitespaceNoNewLine()
			if i.SkipForce() {
				h.Score(i)
			}
		default:
			return nil
		}
		g.Score(h)
		nme := new(MemberExpression)
		*nme = *me
		*me = MemberExpression{
			MemberExpression:  nme,
			Expression:        e,
			IdentifierName:    i,
			TemplateLiteral:   tl,
			PrivateIdentifier: p,
		}
		j.Score(g)
	}
}

// IsSimple returns whether or not the MemberExpression is classed as 'simple'
func (me *MemberExpression) IsSimple() bool {
	return me.Expression != nil || me.IdentifierName != nil || me.SuperProperty || me.PrivateIdentifier != nil || (me.PrimaryExpression != nil && me.PrimaryExpression.IsSimple())
}

// PrimaryExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-PrimaryExpression
//
// It is only valid is one IdentifierReference, Literal, ArrayLiteral,
// ObjectLiteral, FunctionExpression, ClassExpression, TemplateLiteral, or
// ParenthesizedExpression is non-nil or This is true.
type PrimaryExpression struct {
	This                    *Token
	IdentifierReference     *Token
	Literal                 *Token
	ArrayLiteral            *ArrayLiteral
	ObjectLiteral           *ObjectLiteral
	FunctionExpression      *FunctionDeclaration
	ClassExpression         *ClassDeclaration
	TemplateLiteral         *TemplateLiteral
	ParenthesizedExpression *ParenthesizedExpression
	Tokens                  Tokens
}

func (pe *PrimaryExpression) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if g.SkipAbstract() {
		g.AcceptRunWhitespaceNoNewLine()
	}
	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "this"}) {
		pe.This = j.GetLastToken()
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
	} else if t == (parser.Token{Type: TokenKeyword, Data: "class"}) || g.Peek() == (parser.Token{Type: TokenKeyword, Data: "class"}) {
		g := j.NewGoal()
		if g.SkipAbstract() {
			g.AcceptRunWhitespaceNoNewLine()
		}
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
		pe.ParenthesizedExpression = new(ParenthesizedExpression)
		if err := pe.ParenthesizedExpression.parse(&g, yield, await); err != nil {
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

// IsSimple returns whether or not the PrimaryExpression is classed as 'simple'
func (pe *PrimaryExpression) IsSimple() bool {
	return pe.IdentifierReference != nil && pe.IdentifierReference.Data != "eval" && pe.IdentifierReference.Data != "arguments"
}

// ParenthesizedExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ParenthesizedExpression
//
// It is valid for only one of BindingIdentifier, ArrayBindingPattern, and
// ObjectBindingPattern to be non-nil
type ParenthesizedExpression struct {
	Expressions          []AssignmentExpression
	bindingIdentifier    *Token
	arrayBindingPattern  *ArrayBindingPattern
	objectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (cp *ParenthesizedExpression) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("ParenthesizedExpression", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		for {
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				j.AcceptRunWhitespace()
				g := j.NewGoal()
				if t := g.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
					cp.arrayBindingPattern = new(ArrayBindingPattern)
					if err := cp.arrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("ParenthesizedExpression", err)
					}
				} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
					cp.objectBindingPattern = new(ObjectBindingPattern)
					if err := cp.objectBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("ParenthesizedExpression", err)
					}
				} else if cp.bindingIdentifier = g.parseIdentifier(yield, await); cp.bindingIdentifier == nil {
					return j.Error("ParenthesizedExpression", ErrNoIdentifier)
				}
				j.Score(g)
				j.AcceptRunWhitespace()
				if j.SkipOptionalColonType() {
					j.AcceptRunWhitespace()
				}
				if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
					return j.Error("ParenthesizedExpression", ErrMissingClosingParenthesis)
				}
				break
			}
			g := j.NewGoal()
			e := len(cp.Expressions)
			cp.Expressions = append(cp.Expressions, AssignmentExpression{})
			if err := cp.Expressions[e].parse(&g, true, yield, await); err != nil {
				return j.Error("ParenthesizedExpression", err)
			}
			if ae := &cp.Expressions[e]; ae.AssignmentOperator == AssignmentNone && g.SkipOptionalColonType() {
				g.AcceptRunWhitespace()
				if ae.ConditionalExpression != nil && ae.ConditionalExpression.LogicalORExpression != nil {
					if lhs := ae.ConditionalExpression.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression; lhs != nil && len(ae.ConditionalExpression.Tokens) == len(lhs.Tokens) {
						if ae.AssignmentOperator.parse(&g) == nil {
							g.AcceptRunWhitespace()
							ae.ConditionalExpression = nil
							ae.LeftHandSideExpression = lhs
							if ae.AssignmentOperator == AssignmentAssign && lhs.NewExpression != nil && lhs.NewExpression.News == 0 && lhs.NewExpression.MemberExpression.PrimaryExpression != nil && (lhs.NewExpression.MemberExpression.PrimaryExpression.ArrayLiteral != nil || lhs.NewExpression.MemberExpression.PrimaryExpression.ObjectLiteral != nil) {
								ae.AssignmentPattern = new(AssignmentPattern)
								if err := ae.AssignmentPattern.from(lhs.NewExpression.MemberExpression.PrimaryExpression); err != nil {
									z := jsParser(lhs.Tokens[:0])
									return z.Error("AssignmentExpression", err)
								}
								ae.LeftHandSideExpression = nil
							}
							h := g.NewGoal()
							ae.AssignmentExpression = new(AssignmentExpression)
							if err := ae.AssignmentExpression.parse(&h, true, yield, await); err != nil {
								return g.Error("AssignmentExpression", err)
							}
							g.Score(h)
							ae.Tokens = g.ToTokens()
						}
					}
				}
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return j.Error("ParenthesizedExpression", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	cp.Tokens = j.ToTokens()
	return nil
}

// Arguments as defined in TC39
// https://tc39.es/ecma262/#prod-Arguments
type Arguments struct {
	ArgumentList []Argument
	Tokens       Tokens
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
		ae := len(a.ArgumentList)
		a.ArgumentList = append(a.ArgumentList, Argument{})
		g := j.NewGoal()
		if err := a.ArgumentList[ae].parse(&g, yield, await); err != nil {
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

// Argument is an item in an ArgumentList and contains the spread information
// and the AssignementExpression
type Argument struct {
	Spread               bool
	AssignmentExpression AssignmentExpression
	Tokens               Tokens
}

func (a *Argument) parse(j *jsParser, yield, await bool) error {
	a.Spread = j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	if err := a.AssignmentExpression.parse(&g, true, yield, await); err != nil {
		return j.Error("Argument", err)
	}
	j.Score(g)
	a.Tokens = j.ToTokens()
	return nil
}

// CallExpression as defined in ECMA-262
// https://tc39.es/ecma262/#prod-CallExpression
//
// It is only valid for one of MemberExpression, ImportCall, or CallExpression
// to be non-nil or SuperCall to be true.
//
// If MemberExpression is non-nil, or SuperCall is true, Arguments must be
// non-nil.
//
// If CallExpression is non-nil, only one of Arguments, Expression,
// IdentifierName, TemplateLiteral, or PrivateIdentifier must be non-nil.
type CallExpression struct {
	MemberExpression  *MemberExpression
	SuperCall         bool
	ImportCall        *AssignmentExpression
	CallExpression    *CallExpression
	Arguments         *Arguments
	Expression        *Expression
	IdentifierName    *Token
	TemplateLiteral   *TemplateLiteral
	PrivateIdentifier *Token
	Tokens            Tokens
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
		if j.SkipTypeArguments() {
			j.AcceptRunWhitespace()
		}
		g := j.NewGoal()
		ce.Arguments = new(Arguments)
		if err := ce.Arguments.parse(&g, yield, await); err != nil {
			return j.Error("CallExpression", err)
		}
		h := g.NewGoal()
		h.AcceptRunWhitespaceNoNewLine()
		if h.SkipForce() {
			g.Score(h)
		}
		j.Score(g)
	}
	for {
		ce.Tokens = j.ToTokens()
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		h := g.NewGoal()
		var (
			tl   *TemplateLiteral
			a    *Arguments
			i, p *Token
			e    *Expression
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
				if !h.Accept(TokenIdentifier, TokenKeyword, TokenPrivateIdentifier) {
					return h.Error("CallExpression", ErrNoIdentifier)
				}
				i = h.GetLastToken()
				if i.Type == TokenPrivateIdentifier {
					p = i
					i = nil
				}
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
			i := h.NewGoal()
			i.AcceptRunWhitespaceNoNewLine()
			if i.SkipForce() {
				h.Score(i)
			}
		default:
			return nil
		}
		g.Score(h)
		nce := new(CallExpression)
		*nce = *ce
		*ce = CallExpression{
			CallExpression:    nce,
			Expression:        e,
			Arguments:         a,
			IdentifierName:    i,
			TemplateLiteral:   tl,
			PrivateIdentifier: p,
		}
		j.Score(g)
	}
}

// IsSimple returns whether or not the CallExpression is classed as 'simple'
func (ce *CallExpression) IsSimple() bool {
	return ce.Expression != nil || ce.IdentifierName != nil || ce.PrivateIdentifier != nil
}
