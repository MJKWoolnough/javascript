package javascript // import "vimagination.zapto.org/javascript"

import (
	"errors"

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
func ParseScript(t parser.Tokeniser) (*Script, error) {
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
	if j.Accept(TokenIdentifier) || (!yield && j.AcceptToken(parser.Token{TokenKeyword, "yield"}) || (!await && j.AcceptToken(parser.Token{TokenKeyword, "await"}))) {
		return j.GetLastToken()
	}
	return nil
}

func (j *jsParser) parseSemicolon() bool {
	g := j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
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
	if tk := g.Peek(); tk == (parser.Token{TokenKeyword, "class"}) {
		d.ClassDeclaration = new(ClassDeclaration)
		if err := d.ClassDeclaration.parse(&g, yield, await, false); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk == (parser.Token{TokenKeyword, "const"}) || tk == (parser.Token{TokenIdentifier, "let"}) {
		d.LexicalDeclaration = new(LexicalDeclaration)
		if err := d.LexicalDeclaration.parse(&g, true, yield, await); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk == (parser.Token{TokenIdentifier, "async"}) || tk == (parser.Token{TokenKeyword, "function"}) {
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
	if !j.AcceptToken(parser.Token{TokenIdentifier, "let"}) {
		if !j.AcceptToken(parser.Token{TokenKeyword, "const"}) {
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
		if !g.AcceptToken(parser.Token{TokenPunctuator, ","}) {
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
	if t := g.Peek(); t == (parser.Token{TokenPunctuator, "["}) {
		lb.ArrayBindingPattern = new(ArrayBindingPattern)
		if err := lb.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
	} else if t == (parser.Token{TokenPunctuator, "{"}) {
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
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
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
	if !j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		return j.Error("ArrayBindingPattern", ErrMissingOpeningBracket)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			ab.BindingElementList = append(ab.BindingElementList, BindingElement{})
			continue
		}
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			ab.BindingRestElement = new(BindingElement)
			if err := ab.BindingRestElement.parse(&h, nil, yield, await); err != nil {
				return g.Error("ArrayBindingPattern", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
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
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("ArrayBindingPattern", ErrMissingComma)
		}
	}
	ab.Tokens = j.ToTokens()
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
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return j.Error("ObjectBindingPattern", ErrMissingOpeningBrace)
	}
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		for {
			if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
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
			} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
				return j.Error("ObjectBindingPattern", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	ob.Tokens = j.ToTokens()
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
	if !h.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
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

// ArrayLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayLiteral
type ArrayLiteral struct {
	ElementList   []AssignmentExpression
	SpreadElement *AssignmentExpression
	Tokens        Tokens
}

func (al *ArrayLiteral) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		return j.Error("ArrayLiteral", ErrMissingOpeningBracket)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			al.ElementList = append(al.ElementList, AssignmentExpression{})
			continue
		}
		if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			j.AcceptRunWhitespace()
			g := j.NewGoal()
			al.SpreadElement = new(AssignmentExpression)
			if err := al.SpreadElement.parse(&g, true, yield, await); err != nil {
				return j.Error("ArrayLiteral", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return j.Error("ArrayLiteral", ErrMissingClosingBracket)
			}
			break
		}
		g := j.NewGoal()
		ae := len(al.ElementList)
		al.ElementList = append(al.ElementList, AssignmentExpression{})
		if err := al.ElementList[ae].parse(&g, true, yield, await); err != nil {
			return j.Error("ArrayLiteral", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
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
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
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
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("ObjectLiteral", ErrMissingComma)
		}
	}
	ol.Tokens = j.ToTokens()
	return nil
}

// PropertyDefinition as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-PropertyDefinition
//
// It is only valid for either IdentifierReference, IdentifierReference and
// AssignmentExpression (CoverInitializedName), PropertyName and
// AssignmentExpression (PropertyName: AssignmentExpression), MethodDefinition,
// or AssignmentExpression (...AssignmentExpression) to be non-nil.
type PropertyDefinition struct {
	IdentifierReference  *Token
	PropertyName         *PropertyName
	AssignmentExpression *AssignmentExpression
	MethodDefinition     *MethodDefinition
	Tokens               Tokens
}

func (pd *PropertyDefinition) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
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
			if h.AcceptToken(parser.Token{TokenPunctuator, "="}) {
				g.Score(h)
				g.AcceptRunWhitespace()
				pd.IdentifierReference = i
				h = g.NewGoal()
				pd.AssignmentExpression = new(AssignmentExpression)
				if err := pd.AssignmentExpression.parse(&h, true, yield, await); err != nil {
					return g.Error("PropertyDefinition", err)
				}
				g.Score(h)
			} else if t := h.Peek(); t.Type == TokenRightBracePunctuator || t == (parser.Token{TokenPunctuator, ","}) {
				pd.IdentifierReference = i
			}
		}
		if pd.IdentifierReference == nil {
			g = j.NewGoal()
			propertyName := true
			switch g.Peek() {
			case parser.Token{TokenPunctuator, "*"}:
				propertyName = false
			case parser.Token{TokenIdentifier, "async"}, parser.Token{TokenIdentifier, "get"}, parser.Token{TokenIdentifier, "set"}:
				g.Skip()
				g.AcceptRunWhitespace()
				if !g.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
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
				if h.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
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
				var pn *PropertyName
				if pd.PropertyName != nil && pd.PropertyName.ComputedPropertyName != nil {
					pn = pd.PropertyName
				} else {
					g = j.NewGoal()
				}
				pd.PropertyName = nil
				pd.MethodDefinition = new(MethodDefinition)
				if err := pd.MethodDefinition.parse(&g, pn, yield, await); err != nil {
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
// It is only valid for one of BindingIdentifier,
// CoverParenthesizedExpressionAndArrowParameterList (!Async), and
// FormalParameters (+Async) to be non-nil.
//
// Only one of AssignmentExpression or FunctionBody must be non-nil.
type ArrowFunction struct {
	Async                                             bool
	BindingIdentifier                                 *Token
	CoverParenthesizedExpressionAndArrowParameterList *CoverParenthesizedExpressionAndArrowParameterList
	FormalParameters                                  *FormalParameters
	AssignmentExpression                              *AssignmentExpression
	FunctionBody                                      *Block
	Tokens                                            Tokens
}

func (af *ArrowFunction) parse(j *jsParser, pe *PrimaryExpression, in, yield, await bool) error {
	if pe == nil {
		if !j.AcceptToken(parser.Token{TokenIdentifier, "async"}) {
			return j.Error("ArrowFunction", ErrInvalidAsyncArrowFunction)
		}
		af.Async = true
		j.AcceptRunWhitespaceNoNewLine()
		if j.Peek() == (parser.Token{TokenPunctuator, "("}) {
			g := j.NewGoal()
			af.FormalParameters = new(FormalParameters)
			if err := af.FormalParameters.parse(&g, false, true); err != nil {
				return j.Error("ArrowFunction", err)
			}
			j.Score(g)
		} else {
			g := j.NewGoal()
			if af.BindingIdentifier = g.parseIdentifier(yield, true); af.BindingIdentifier == nil {
				return j.Error("ArrowFunction", ErrNoIdentifier)
			}
			j.Score(g)
		}
	} else if pe.CoverParenthesizedExpressionAndArrowParameterList != nil {
		af.CoverParenthesizedExpressionAndArrowParameterList = pe.CoverParenthesizedExpressionAndArrowParameterList
	} else {
		af.BindingIdentifier = pe.IdentifierReference
	}
	j.AcceptRunWhitespaceNoNewLine()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "=>"}) {
		return j.Error("ArrowFunction", ErrMissingArrow)
	}
	j.AcceptRunWhitespace()
	if j.Peek() == (parser.Token{TokenPunctuator, "{"}) {
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

// Errors
var (
	ErrReservedIdentifier        = errors.New("reserved identifier")
	ErrNoIdentifier              = errors.New("missing identifier")
	ErrMissingFunction           = errors.New("missing function")
	ErrMissingOpeningParenthesis = errors.New("missing opening parenthesis")
	ErrMissingClosingParenthesis = errors.New("missing closing parenthesis")
	ErrMissingOpeningBrace       = errors.New("missing opening brace")
	ErrMissingClosingBrace       = errors.New("missing closing brace")
	ErrMissingOpeningBracket     = errors.New("missing opening bracket")
	ErrMissingClosingBracket     = errors.New("missing closing bracket")
	ErrMissingComma              = errors.New("missing comma")
	ErrMissingArrow              = errors.New("missing arrow")
	ErrMissingCaseClause         = errors.New("missing case clause")
	ErrMissingExpression         = errors.New("missing expression")
	ErrMissingCatchFinally       = errors.New("missing catch/finally block")
	ErrMissingSemiColon          = errors.New("missing semi-colon")
	ErrMissingColon              = errors.New("missing colon")
	ErrMissingInitializer        = errors.New("missing initializer")
	ErrInvalidStatementList      = errors.New("invalid statement list")
	ErrInvalidStatement          = errors.New("invalid statement")
	ErrInvalidDeclaration        = errors.New("invalid declaration")
	ErrInvalidLexicalDeclaration = errors.New("invalid lexical declaration")
	ErrInvalidAssignment         = errors.New("invalid assignment operator")
	ErrInvalidSuperProperty      = errors.New("invalid super property")
	ErrInvalidMetaProperty       = errors.New("invalid meta property")
	ErrInvalidTemplate           = errors.New("invalid template")
	ErrInvalidAsyncArrowFunction = errors.New("invalid async arrow function")
)
