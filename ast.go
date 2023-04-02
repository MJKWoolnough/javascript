package javascript // import "vimagination.zapto.org/javascript"

import (
	"vimagination.zapto.org/parser"
)

// Script represents the top-level of a parsed javascript text
type Script struct {
	StatementList []StatementListItem
	Tokens        Tokens
}

// ParseScript parses a javascript input into an AST.
//
// It is recommended to use ParseModule instead of this function.
func ParseScript(t Tokeniser) (*Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	s := new(Script)
	if err := s.parse(&j); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Script) parse(j *jsParser) error {
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		si := len(s.StatementList)
		s.StatementList = append(s.StatementList, StatementListItem{})
		if err := s.StatementList[si].parse(&g, false, false, false); err != nil {
			return err
		}
		j.Score(g)
	}
	s.Tokens = j.ToTokens()
	return nil
}

func (j *jsParser) parseIdentifier(yield, await bool) *Token {
	if j.Accept(TokenIdentifier) || (!yield && j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "yield"}) || (!await && j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "await"}))) {
		return j.GetLastToken()
	}
	return nil
}

func (j *jsParser) parseSemicolon() bool {
	g := j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
		j.Score(g)
		return true
	} else if g.Peek().Type == TokenRightBracePunctuator {
		return true
	}
	g = j.NewGoal()
	g.AcceptRunWhitespaceNoNewLine()
	if t := g.Peek().Type; t == TokenLineTerminator || t == parser.TokenDone {
		return true
	}
	return false
}

// Declaration as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-Declaration
//
// Only one of ClassDeclaration, FunctionDeclaration or LexicalDeclaration must
// be non-nil
type Declaration struct {
	ClassDeclaration    *ClassDeclaration
	FunctionDeclaration *FunctionDeclaration
	LexicalDeclaration  *LexicalDeclaration
	Tokens              Tokens
}

func (d *Declaration) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	h := g.NewGoal()
	if h.SkipAbstract() {
		h.AcceptRunWhitespaceNoNewLine()
	}
	if tk := h.Peek(); tk == (parser.Token{Type: TokenKeyword, Data: "class"}) {
		d.ClassDeclaration = new(ClassDeclaration)
		if err := d.ClassDeclaration.parse(&h, yield, await, false); err != nil {
			return j.Error("Declaration", err)
		}
		g.Score(h)
	} else if tk == (parser.Token{Type: TokenKeyword, Data: "const"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "let"}) {
		d.LexicalDeclaration = new(LexicalDeclaration)
		if err := d.LexicalDeclaration.parse(&g, true, yield, await); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk == (parser.Token{Type: TokenIdentifier, Data: "async"}) || tk == (parser.Token{Type: TokenKeyword, Data: "function"}) {
		d.FunctionDeclaration = new(FunctionDeclaration)
		if err := d.FunctionDeclaration.parse(&g, yield, await, false); err != nil {
			return j.Error("Declaration", err)
		}
	} else {
		return j.Error("Declaration", ErrInvalidDeclaration)
	}
	j.Score(g)
	d.Tokens = j.ToTokens()
	return nil
}

// LetOrConst specifies whether a LexicalDeclaration is a let or const declaration
type LetOrConst bool

// Valid LetOrConst values
const (
	Let   LetOrConst = false
	Const LetOrConst = true
)

// LexicalDeclaration as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-LexicalDeclaration
type LexicalDeclaration struct {
	LetOrConst
	BindingList []LexicalBinding
	Tokens      Tokens
}

func (ld *LexicalDeclaration) parse(j *jsParser, in, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "let"}) {
		if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "const"}) {
			return j.Error("LexicalDeclaration", ErrInvalidLexicalDeclaration)
		}
		ld.LetOrConst = Const
	}
	for {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		lb := len(ld.BindingList)
		ld.BindingList = append(ld.BindingList, LexicalBinding{})
		if err := ld.BindingList[lb].parse(&g, in, yield, await); err != nil {
			return j.Error("LexicalDeclaration", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			if j.parseSemicolon() {
				break
			}
			return j.Error("LexicalDeclaration", ErrInvalidLexicalDeclaration)
		}
		j.Score(g)
	}
	ld.Tokens = j.ToTokens()
	return nil
}

// LexicalBinding as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-LexicalBinding
//
// Only one of BindingIdentifier, ArrayBindingPattern or ObjectBindingPattern
// must be non-nil. The Initializer is optional only for a BindingIdentifier.
type LexicalBinding struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}

func (lb *LexicalBinding) parse(j *jsParser, in, yield, await bool) error {
	g := j.NewGoal()
	if t := g.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
		lb.ArrayBindingPattern = new(ArrayBindingPattern)
		if err := lb.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
	} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		lb.ObjectBindingPattern = new(ObjectBindingPattern)
		if err := lb.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
	} else if lb.BindingIdentifier = g.parseIdentifier(yield, await); lb.BindingIdentifier == nil {
		return j.Error("LexicalBinding", ErrNoIdentifier)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.SkipColonType() {
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
	}
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		lb.Initializer = new(AssignmentExpression)
		if err := lb.Initializer.parse(&g, in, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
		j.Score(g)
	} else if lb.BindingIdentifier == nil {
		return j.Error("LexicalBinding", ErrMissingInitializer)
	}
	lb.Tokens = j.ToTokens()
	return nil
}

// ArrayBindingPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayBindingPattern
type ArrayBindingPattern struct {
	BindingElementList []BindingElement
	BindingRestElement *BindingElement
	Tokens             Tokens
}

func (ab *ArrayBindingPattern) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return j.Error("ArrayBindingPattern", ErrMissingOpeningBracket)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			ab.BindingElementList = append(ab.BindingElementList, BindingElement{})
			continue
		}
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			ab.BindingRestElement = new(BindingElement)
			if err := ab.BindingRestElement.parse(&h, nil, yield, await); err != nil {
				return g.Error("ArrayBindingPattern", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return j.Error("ArrayBindingPattern", ErrMissingClosingBracket)
			}
			break
		}
		be := len(ab.BindingElementList)
		ab.BindingElementList = append(ab.BindingElementList, BindingElement{})
		if err := ab.BindingElementList[be].parse(&g, nil, yield, await); err != nil {
			return j.Error("ArrayBindingPattern", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return j.Error("ArrayBindingPattern", ErrMissingComma)
		}
	}
	ab.Tokens = j.ToTokens()
	return nil
}

func (ab *ArrayBindingPattern) from(al *ArrayLiteral) error {
	ab.BindingElementList = make([]BindingElement, 0, len(al.ElementList))
	hasSpread := false
	for _, ae := range al.ElementList {
		if hasSpread {
			z := jsParser(al.Tokens[:0])
			return z.Error("ArrayBindingPattern", ErrBadRestElement)
		} else if ae.Spread {
			hasSpread = true
			ab.BindingRestElement = new(BindingElement)
			if err := ab.BindingRestElement.from(&ae.AssignmentExpression); err != nil {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayBindingPattern", err)
			}
		} else {
			var be BindingElement
			if err := be.from(&ae.AssignmentExpression); err != nil {
				z := jsParser(al.Tokens[:0])
				return z.Error("ArrayBindingPattern", err)
			}
			ab.BindingElementList = append(ab.BindingElementList, be)
		}
	}
	ab.Tokens = al.Tokens
	return nil
}

func (ab *ArrayBindingPattern) fromAP(ap *ArrayAssignmentPattern) error {
	ab.BindingElementList = make([]BindingElement, len(ap.AssignmentElements))
	for n, ae := range ap.AssignmentElements {
		if err := ab.BindingElementList[n].fromAP(ae.DestructuringAssignmentTarget.LeftHandSideExpression, ae.DestructuringAssignmentTarget.AssignmentPattern); err != nil {
			z := jsParser(ap.Tokens[:0])
			return z.Error("ArrayBindingPattern", err)
		}
		ab.BindingElementList[n].Initializer = ae.Initializer
		ab.BindingElementList[n].Tokens = ae.Tokens
	}
	if ap.AssignmentRestElement != nil {
		ab.BindingRestElement = new(BindingElement)
		if err := ab.BindingRestElement.fromAP(ap.AssignmentRestElement, nil); err != nil {
			z := jsParser(ap.Tokens[:0])
			return z.Error("ArrayBindingPattern", err)
		}
	}
	ab.Tokens = ap.Tokens
	return nil
}

// ObjectBindingPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ObjectBindingPattern
type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	BindingRestProperty *Token
	Tokens              Tokens
}

func (ob *ObjectBindingPattern) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ObjectBindingPattern", ErrMissingOpeningBrace)
	}
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		for {
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				j.AcceptRunWhitespace()
				if ob.BindingRestProperty = j.parseIdentifier(yield, await); ob.BindingRestProperty == nil {
					return j.Error("ObjectBindingPattern", ErrNoIdentifier)
				}
				j.AcceptRunWhitespace()
				if !j.Accept(TokenRightBracePunctuator) {
					return j.Error("ObjectBindingPattern", ErrMissingClosingBrace)
				}
				break
			}
			g := j.NewGoal()
			bp := len(ob.BindingPropertyList)
			ob.BindingPropertyList = append(ob.BindingPropertyList, BindingProperty{})
			if err := ob.BindingPropertyList[bp].parse(&g, yield, await); err != nil {
				return j.Error("ObjectBindingPattern", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if j.Accept(TokenRightBracePunctuator) {
				break
			} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return j.Error("ObjectBindingPattern", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	ob.Tokens = j.ToTokens()
	return nil
}

func (ob *ObjectBindingPattern) from(ol *ObjectLiteral) error {
	for _, pd := range ol.PropertyDefinitionList {
		if pd.AssignmentExpression == nil {
			z := jsParser(ol.Tokens[:0])
			return z.Error("ObjectBindingPattern", ErrNoIdentifier)
		}
		bp := BindingProperty{
			Tokens: pd.Tokens,
		}
		if pd.PropertyName != nil {
			bp.PropertyName = *pd.PropertyName
			if pd.IsCoverInitializedName {
				bp.BindingElement.SingleNameBinding = pd.PropertyName.LiteralPropertyName
				bp.BindingElement.Initializer = pd.AssignmentExpression
				bp.BindingElement.Tokens = pd.Tokens
			} else if err := bp.BindingElement.from(pd.AssignmentExpression); err != nil {
				z := jsParser(ol.Tokens[:0])
				return z.Error("ObjectBindingPattern", err)
			}
		} else {
			if pe, ok := UnwrapConditional(pd.AssignmentExpression.ConditionalExpression).(*PrimaryExpression); ok && pe.IdentifierReference != nil {
				ob.BindingRestProperty = pe.IdentifierReference
				break
			} else {
				z := jsParser(ol.Tokens[:0])
				return z.Error("ObjectBindingPattern", ErrNoIdentifier)
			}
		}
		ob.BindingPropertyList = append(ob.BindingPropertyList, bp)
	}
	ob.Tokens = ol.Tokens
	return nil
}

func (ob *ObjectBindingPattern) fromAP(op *ObjectAssignmentPattern) error {
	ob.BindingPropertyList = make([]BindingProperty, len(op.AssignmentPropertyList))
	for n := range op.AssignmentPropertyList {
		if err := ob.BindingPropertyList[n].fromAP(&op.AssignmentPropertyList[n]); err != nil {
			z := jsParser(op.Tokens[:0])
			return z.Error("ObjectBindingPattern", err)
		}
	}
	if op.AssignmentRestElement != nil {
		if op.AssignmentRestElement.CallExpression != nil || op.AssignmentRestElement.OptionalExpression != nil || op.AssignmentRestElement.NewExpression.News != 0 || op.AssignmentRestElement.NewExpression.MemberExpression.PrimaryExpression == nil || op.AssignmentRestElement.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference == nil {
			z := jsParser(op.Tokens[:0])
			return z.Error("ObjectBindingPattern", ErrBadRestElement)
		}
		ob.BindingRestProperty = op.AssignmentRestElement.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference
	}
	ob.Tokens = op.Tokens
	return nil
}

// BindingProperty as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-BindingProperty
//
// A SingleNameBinding, with or without an initializer, is cloned into the
// Property Name and Binding Element. This allows the Binding Element
// Identifier to be modified while keeping the correct Property Name
type BindingProperty struct {
	PropertyName   PropertyName
	BindingElement BindingElement
	Tokens         Tokens
}

func (bp *BindingProperty) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if err := bp.PropertyName.parse(&g, yield, await); err != nil {
		return j.Error("BindingProperty", err)
	}
	h := g.NewGoal()
	h.AcceptRunWhitespace()
	var snb *Token
	if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		i := j.NewGoal()
		if bp.PropertyName.LiteralPropertyName == nil || i.parseIdentifier(yield, await) == nil {
			return h.Error("BindingProperty", ErrMissingColon)
		}
		lpn := *bp.PropertyName.LiteralPropertyName
		snb = &lpn
	} else {
		h.AcceptRunWhitespace()
		g.Score(h)
		j.Score(g)
		g = j.NewGoal()
	}
	if err := bp.BindingElement.parse(&g, snb, yield, await); err != nil {
		return j.Error("BindingProperty", err)
	}
	j.Score(g)
	bp.Tokens = j.ToTokens()
	return nil
}

func (bp *BindingProperty) fromAP(ap *AssignmentProperty) error {
	bp.PropertyName = ap.PropertyName
	if ap.DestructuringAssignmentTarget != nil {
		if ap.DestructuringAssignmentTarget.LeftHandSideExpression != nil {
			if ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression == nil || ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.News != 0 || ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression == nil || ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference == nil {
				y := jsParser(ap.Tokens[:0])
				z := jsParser(ap.DestructuringAssignmentTarget.Tokens[:0])
				return y.Error("BindingProperty", z.Error("BindingElement", ErrNoIdentifier))
			}
			bp.BindingElement.SingleNameBinding = ap.DestructuringAssignmentTarget.LeftHandSideExpression.NewExpression.MemberExpression.PrimaryExpression.IdentifierReference
		} else if ap.DestructuringAssignmentTarget.AssignmentPattern.ArrayAssignmentPattern != nil {
			bp.BindingElement.ArrayBindingPattern = new(ArrayBindingPattern)
			if err := bp.BindingElement.ArrayBindingPattern.fromAP(ap.DestructuringAssignmentTarget.AssignmentPattern.ArrayAssignmentPattern); err != nil {
				z := jsParser(ap.Tokens[:0])
				return z.Error("BindingProperty", err)
			}
		} else {
			bp.BindingElement.ObjectBindingPattern = new(ObjectBindingPattern)
			if err := bp.BindingElement.ObjectBindingPattern.fromAP(ap.DestructuringAssignmentTarget.AssignmentPattern.ObjectAssignmentPattern); err != nil {
				z := jsParser(ap.Tokens[:0])
				return z.Error("BindingProperty", err)
			}
		}
		for n := range ap.Tokens {
			if &ap.Tokens[n] == &ap.DestructuringAssignmentTarget.Tokens[0] {
				bp.BindingElement.Tokens = ap.Tokens[n:]
				break
			}
		}
	} else {
		bp.BindingElement.SingleNameBinding = bp.PropertyName.LiteralPropertyName
		bp.BindingElement.Tokens = ap.Tokens
	}
	bp.BindingElement.Initializer = ap.Initializer
	bp.Tokens = ap.Tokens
	return nil
}

// VariableDeclaration as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-VariableDeclaration
type VariableDeclaration LexicalBinding

func (v *VariableDeclaration) parse(j *jsParser, in, yield, await bool) error {
	if err := ((*LexicalBinding)(v)).parse(j, in, yield, await); err != nil {
		errr := err.(Error)
		errr.Parsing = "VariableDeclaration"
		return errr
	}
	return nil
}

// ArrayElement is an element of ElementList in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ElementList
type ArrayElement struct {
	Spread bool
	AssignmentExpression
	Tokens Tokens
}

func (ae *ArrayElement) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	ae.Spread = g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."})
	g.AcceptRunWhitespace()
	h := g.NewGoal()
	if err := ae.AssignmentExpression.parse(&h, true, yield, await); err != nil {
		return j.Error("ArrayElement", err)
	}
	g.Score(h)
	j.Score(g)
	ae.Tokens = j.ToTokens()
	return nil
}

// ArrayLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayLiteral
type ArrayLiteral struct {
	ElementList []ArrayElement
	Tokens      Tokens
}

func (al *ArrayLiteral) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return j.Error("ArrayLiteral", ErrMissingOpeningBracket)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			al.ElementList = append(al.ElementList, ArrayElement{})
			continue
		}
		g := j.NewGoal()
		var ae ArrayElement
		if err := ae.parse(&g, yield, await); err != nil {
			return j.Error("ArrayLiteral", err)
		}
		al.ElementList = append(al.ElementList, ae)
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return j.Error("ArrayLiteral", ErrMissingComma)
		}
	}
	al.Tokens = j.ToTokens()
	return nil
}

// ObjectLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ObjectLiteral
type ObjectLiteral struct {
	PropertyDefinitionList []PropertyDefinition
	Tokens                 Tokens
}

func (ol *ObjectLiteral) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ObjectLiteral", ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		pd := len(ol.PropertyDefinitionList)
		ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, PropertyDefinition{})
		if err := ol.PropertyDefinitionList[pd].parse(&g, yield, await); err != nil {
			return j.Error("ObjectLiteral", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.Accept(TokenRightBracePunctuator) {
			break
		} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return j.Error("ObjectLiteral", ErrMissingComma)
		}
	}
	ol.Tokens = j.ToTokens()
	return nil
}

// PropertyDefinition as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-PropertyDefinition
//
// One, and only one, of AssignmentExpression or MethodDefinition must be
// non-nil.
//
// It is only valid for PropertyName to be non-nil when AssignmentExpression is
// also non-nil.
//
// The IdentifierReference is stored within PropertyName.
type PropertyDefinition struct {
	IsCoverInitializedName bool
	PropertyName           *PropertyName
	AssignmentExpression   *AssignmentExpression
	MethodDefinition       *MethodDefinition
	Tokens                 Tokens
}

func (pd *PropertyDefinition) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		pd.AssignmentExpression = new(AssignmentExpression)
		if err := pd.AssignmentExpression.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyDefinition", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		if i := g.parseIdentifier(yield, await); i != nil {
			h := g.NewGoal()
			h.AcceptRunWhitespace()
			if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
				pd.PropertyName = &PropertyName{
					LiteralPropertyName: i,
					Tokens:              g.ToTokens(),
				}
				g.Score(h)
				g.AcceptRunWhitespace()
				h = g.NewGoal()
				pd.AssignmentExpression = new(AssignmentExpression)
				if err := pd.AssignmentExpression.parse(&h, true, yield, await); err != nil {
					return g.Error("PropertyDefinition", err)
				}
				g.Score(h)
				pd.IsCoverInitializedName = true
			} else if t := h.Peek(); t.Type == TokenRightBracePunctuator || t == (parser.Token{Type: TokenPunctuator, Data: ","}) {
				pd.PropertyName = &PropertyName{
					LiteralPropertyName: i,
					Tokens:              g.ToTokens(),
				}
				pd.AssignmentExpression = &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: &Token{
							Token:   i.Token,
							Pos:     i.Pos,
							Line:    i.Line,
							LinePos: i.LinePos,
						},
						Tokens: pd.PropertyName.Tokens,
					}),
					Tokens: pd.PropertyName.Tokens,
				}
			}
		}
		if pd.PropertyName == nil {
			g = j.NewGoal()
			propertyName := true
			switch g.Peek() {
			case parser.Token{Type: TokenPunctuator, Data: "*"}:
				propertyName = false
			case parser.Token{Type: TokenIdentifier, Data: "async"}, parser.Token{Type: TokenIdentifier, Data: "get"}, parser.Token{Type: TokenIdentifier, Data: "set"}:
				g.Skip()
				g.AcceptRunWhitespace()
				if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
					propertyName = false
				}
			}
			g = j.NewGoal()
			if propertyName {
				pd.PropertyName = new(PropertyName)
				if err := pd.PropertyName.parse(&g, yield, await); err != nil {
					return j.Error("PropertyDefinition", err)
				}
				h := g.NewGoal()
				h.AcceptRunWhitespace()
				if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
					h.AcceptRunWhitespace()
					i := h.NewGoal()
					pd.AssignmentExpression = new(AssignmentExpression)
					if err := pd.AssignmentExpression.parse(&i, true, yield, await); err != nil {
						return h.Error("PropertyDefinition", err)
					}
					h.Score(i)
					g.Score(h)
				} else {
					propertyName = false
				}
			}
			if !propertyName {
				pd.MethodDefinition = new(MethodDefinition)
				if pd.PropertyName != nil && pd.PropertyName.ComputedPropertyName != nil {
					pd.MethodDefinition.ClassElementName.PropertyName = pd.PropertyName
					pd.MethodDefinition.ClassElementName.Tokens = pd.PropertyName.Tokens
				} else {
					g = j.NewGoal()
				}
				pd.PropertyName = nil
				if err := pd.MethodDefinition.parse(&g, yield, await); err != nil {
					return j.Error("PropertyDefinition", err)
				}
			}
		}
		j.Score(g)
	}
	pd.Tokens = j.ToTokens()
	return nil
}

// TemplateLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-TemplateLiteral
//
// If NoSubstitutionTemplate is non-nil it is only valid for TemplateHead,
// Expressions, TemplateMiddleList, and TemplateTail to be nil.
//
// If NoSubstitutionTemplate is nil, the TemplateHead, Expressions, and
// TemplateTail must be non-nil. TemplateMiddleList must have a length of one
// less than the length of Expressions.
type TemplateLiteral struct {
	NoSubstitutionTemplate *Token
	TemplateHead           *Token
	Expressions            []Expression
	TemplateMiddleList     []*Token
	TemplateTail           *Token
	Tokens                 Tokens
}

func (tl *TemplateLiteral) parse(j *jsParser, yield, await bool) error {
	if j.Accept(TokenNoSubstitutionTemplate) {
		tl.NoSubstitutionTemplate = j.GetLastToken()
	} else if !j.Accept(TokenTemplateHead) {
		return j.Error("TemplateLiteral", ErrInvalidTemplate)
	} else {
		tl.TemplateHead = j.GetLastToken()
		for {
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			e := len(tl.Expressions)
			tl.Expressions = append(tl.Expressions, Expression{})
			if err := tl.Expressions[e].parse(&g, true, yield, await); err != nil {
				return j.Error("TemplateLiteral", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if j.Accept(TokenTemplateTail) {
				tl.TemplateTail = j.GetLastToken()
				break
			} else if !j.Accept(TokenTemplateMiddle) {
				return j.Error("TemplateLiteral", ErrInvalidTemplate)
			}
			tl.TemplateMiddleList = append(tl.TemplateMiddleList, j.GetLastToken())
		}
	}
	tl.Tokens = j.ToTokens()
	return nil
}

// ArrowFunction as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrowFunction
//
// Also includes AsyncArrowFunction.
//
// Only one of BindingIdentifier or FormalParameters must be non-nil.
//
// Only one of AssignmentExpression or FunctionBody must be non-nil.
type ArrowFunction struct {
	Async                bool
	BindingIdentifier    *Token
	FormalParameters     *FormalParameters
	AssignmentExpression *AssignmentExpression
	FunctionBody         *Block
	Tokens               Tokens
}

func (af *ArrowFunction) parse(j *jsParser, pe *PrimaryExpression, in, yield, await bool) error {
	if pe == nil {
		if !j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "async"}) {
			return j.Error("ArrowFunction", ErrInvalidAsyncArrowFunction)
		}
		af.Async = true
		j.AcceptRunWhitespaceNoNewLine()
		if j.SkipGeneric() {
			j.AcceptRunWhitespaceNoNewLine()
		}
		if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
			g := j.NewGoal()
			af.FormalParameters = new(FormalParameters)
			if err := af.FormalParameters.parse(&g, false, true); err != nil {
				return j.Error("ArrowFunction", err)
			}
			j.Score(g)
			j.AcceptRunWhitespaceNoNewLine()
			j.SkipReturnType()
		} else {
			g := j.NewGoal()
			if af.BindingIdentifier = g.parseIdentifier(yield, true); af.BindingIdentifier == nil {
				return j.Error("ArrowFunction", ErrNoIdentifier)
			}
			j.Score(g)
		}
	} else if pe.ParenthesizedExpression != nil {
		af.FormalParameters = new(FormalParameters)
		if err := af.FormalParameters.from(pe.ParenthesizedExpression); err != nil {
			z := jsParser(pe.ParenthesizedExpression.Tokens[:0])
			return z.Error("ArrowFunction", err)
		}
		j.AcceptRunWhitespaceNoNewLine()
		j.SkipReturnType()
	} else {
		af.BindingIdentifier = pe.IdentifierReference
	}
	j.AcceptRunWhitespaceNoNewLine()
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "=>"}) {
		return j.Error("ArrowFunction", ErrMissingArrow)
	}
	j.AcceptRunWhitespace()
	if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		g := j.NewGoal()
		af.FunctionBody = new(Block)
		if err := af.FunctionBody.parse(&g, false, af.Async, true); err != nil {
			return j.Error("ArrowFunction", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		af.AssignmentExpression = new(AssignmentExpression)
		if err := af.AssignmentExpression.parse(&g, in, false, af.Async); err != nil {
			return j.Error("ArrowFunction", err)
		}
		j.Score(g)
	}
	af.Tokens = j.ToTokens()
	return nil
}
