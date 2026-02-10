package javascript // import "vimagination.zapto.org/javascript"

import "vimagination.zapto.org/parser"

// Script represents the top-level of a parsed JavaScript text
type Script struct {
	StatementList []StatementListItem
	Comments      [2]Comments
	Tokens        Tokens
}

// ParseScript parses a JavaScript input into an AST.
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

// ScriptToModule converts a Script type to a Module type
func ScriptToModule(s *Script) *Module {
	m := &Module{
		ModuleListItems: make([]ModuleItem, len(s.StatementList)),
		Comments:        s.Comments,
		Tokens:          s.Tokens,
	}

	for n := range s.StatementList {
		m.ModuleListItems[n] = ModuleItem{
			StatementListItem: &s.StatementList[n],
			Tokens:            s.StatementList[n].Tokens,
		}
	}

	return m
}

func (s *Script) parse(j *jsParser) error {
	s.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()
	g := j.NewGoal()

	for g.AcceptRunWhitespace() != parser.TokenDone {
		g = j.NewGoal()

		g.AcceptRunWhitespaceNoComment()
		j.Score(g)

		g = j.NewGoal()
		si := len(s.StatementList)
		s.StatementList = append(s.StatementList, StatementListItem{})

		if err := s.StatementList[si].parse(&g, false, false, false); err != nil {
			return err
		}

		j.Score(g)

		g = j.NewGoal()
	}

	s.Comments[1] = j.AcceptRunWhitespaceComments()
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

	t := g.Peek().Type

	return t == TokenLineTerminator || t == TokenSingleLineComment || t == parser.TokenDone
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
	Comments            Comments
	Tokens              Tokens
}

func (d *Declaration) parse(j *jsParser, yield, await, export bool) error {
	g := j.NewGoal()
	h := g.NewGoal()
	i := h.NewGoal()

	if i.SkipAbstract() {
		i.AcceptRunWhitespaceNoNewlineComments()
		h.Score(i)
		h.AcceptRunWhitespaceNoNewLine()
	}

	if tk := h.Peek(); tk == (parser.Token{Type: TokenKeyword, Data: "class"}) {
		d.Comments = i.ToTypescriptComments()
		d.ClassDeclaration = new(ClassDeclaration)

		if err := d.ClassDeclaration.parse(&h, yield, await, false); err != nil {
			return j.Error("Declaration", err)
		}

		g.Score(h)
	} else if tk = g.Peek(); tk == (parser.Token{Type: TokenKeyword, Data: "const"}) || tk == (parser.Token{Type: TokenIdentifier, Data: "let"}) {
		d.LexicalDeclaration = new(LexicalDeclaration)

		if err := d.LexicalDeclaration.parse(&g, true, yield, await); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk == (parser.Token{Type: TokenIdentifier, Data: "async"}) || tk == (parser.Token{Type: TokenKeyword, Data: "function"}) {
		d.FunctionDeclaration = new(FunctionDeclaration)

		if err := d.FunctionDeclaration.parse(&g, yield, await, false, export); err != nil {
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
		j.AcceptRunWhitespaceNoComment()

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
	Comments             [2]Comments
	Tokens               Tokens
}

func (lb *LexicalBinding) parse(j *jsParser, in, yield, await bool) error {
	lb.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

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

	lb.Comments[1] = j.AcceptRunWhitespaceComments()

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	h := g.NewGoal()

	if h.SkipColonType() {
		lb.Comments[1] = append(lb.Comments[1], h.ToTypescriptComments()...)
		lb.Comments[1] = append(lb.Comments[1], h.AcceptRunWhitespaceComments()...)

		g.Score(h)
		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()
	}

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		g.AcceptRunWhitespaceNoComment()
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

func (lb *LexicalBinding) hasFirstComment() bool {
	return len(lb.Comments[0]) > 0
}

// ArrayBindingPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayBindingPattern
type ArrayBindingPattern struct {
	BindingElementList []BindingElement
	BindingRestElement *BindingElement
	Comments           [3]Comments
	Tokens             Tokens
}

func (ab *ArrayBindingPattern) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return j.Error("ArrayBindingPattern", ErrMissingOpeningBracket)
	}

	ab.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			j.AcceptRunWhitespaceNoComment()

			g := j.NewGoal()
			ab.BindingElementList = append(ab.BindingElementList, BindingElement{
				Comments: [2]Comments{g.AcceptRunWhitespaceComments()},
				Tokens:   g.ToTokens(),
			})

			j.Score(g)
			j.AcceptRunWhitespace()
			j.Skip()

			continue
		}

		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
			ab.Comments[1] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()
			j.AcceptRunWhitespaceNoComment()

			g = j.NewGoal()
			ab.BindingRestElement = new(BindingElement)

			if err := ab.BindingRestElement.parse(&g, nil, yield, await); err != nil {
				return g.Error("ArrayBindingPattern", err)
			}

			j.Score(g)

			g = j.NewGoal()

			g.AcceptRunWhitespace()

			if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
				return g.Error("ArrayBindingPattern", ErrMissingClosingBracket)
			}

			break
		}

		g = j.NewGoal()

		be := len(ab.BindingElementList)
		ab.BindingElementList = append(ab.BindingElementList, BindingElement{})

		if err := ab.BindingElementList[be].parse(&g, nil, yield, await); err != nil {
			return j.Error("ArrayBindingPattern", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("ArrayBindingPattern", ErrMissingComma)
		}

		j.Score(g)
	}

	ab.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

	ab.Tokens = j.ToTokens()

	return nil
}

// ObjectBindingPattern as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ObjectBindingPattern
type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	BindingRestProperty *Token
	Comments            [5]Comments
	Tokens              Tokens
}

func (ob *ObjectBindingPattern) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ObjectBindingPattern", ErrMissingOpeningBrace)
	}

	ob.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if !g.Accept(TokenRightBracePunctuator) {
		for {
			g = j.NewGoal()

			g.AcceptRunWhitespace()

			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				ob.Comments[1] = j.AcceptRunWhitespaceComments()

				j.AcceptRunWhitespace()
				j.Skip()

				ob.Comments[2] = j.AcceptRunWhitespaceComments()

				j.AcceptRunWhitespace()

				if ob.BindingRestProperty = j.parseIdentifier(yield, await); ob.BindingRestProperty == nil {
					return j.Error("ObjectBindingPattern", ErrNoIdentifier)
				}

				ob.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
				g = j.NewGoal()

				g.AcceptRunWhitespace()

				if !g.Accept(TokenRightBracePunctuator) {
					return g.Error("ObjectBindingPattern", ErrMissingClosingBrace)
				}

				break
			}

			j.AcceptRunWhitespaceNoComment()

			g = j.NewGoal()

			bp := len(ob.BindingPropertyList)
			ob.BindingPropertyList = append(ob.BindingPropertyList, BindingProperty{})

			if err := ob.BindingPropertyList[bp].parse(&g, yield, await); err != nil {
				return j.Error("ObjectBindingPattern", err)
			}

			j.Score(g)

			g = j.NewGoal()

			g.AcceptRunWhitespace()

			if g.Accept(TokenRightBracePunctuator) {
				break
			} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return g.Error("ObjectBindingPattern", ErrMissingComma)
			}

			j.Score(g)
		}
	}

	ob.Comments[4] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

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
	Comments       [2]Comments
	Tokens         Tokens
}

func (bp *BindingProperty) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	bp.Comments[0] = g.AcceptRunWhitespaceComments()

	g.AcceptRunWhitespace()

	h := g.NewGoal()

	if err := bp.PropertyName.parse(&h, yield, await); err != nil {
		return g.Error("BindingProperty", err)
	}

	g.Score(h)

	h = g.NewGoal()
	bp.Comments[1] = h.AcceptRunWhitespaceCommentsInList()

	h.AcceptRunWhitespace()

	var snb *Token

	if !h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		i := j.NewGoal()

		i.AcceptRunWhitespace()

		if bp.PropertyName.LiteralPropertyName == nil || i.parseIdentifier(yield, await) == nil {
			return h.Error("BindingProperty", ErrMissingColon)
		}

		lpn := *bp.PropertyName.LiteralPropertyName
		snb = &lpn
		bp.BindingElement.Comments = bp.Comments
		bp.Comments = [2]Comments{}
	} else {
		h.AcceptRunWhitespaceNoComment()
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
type VariableDeclaration = LexicalBinding

// ArrayElement is an element of ElementList in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ElementList
type ArrayElement struct {
	Spread               bool
	AssignmentExpression AssignmentExpression
	Comments             Comments
	Tokens               Tokens
}

func (ae *ArrayElement) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	h := j.NewGoal()

	h.AcceptRunWhitespace()

	if ae.Spread = h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}); ae.Spread {
		ae.Comments = g.AcceptRunWhitespaceComments()

		g.AcceptRunWhitespace()
		g.Skip()

		g.AcceptRunWhitespaceNoComment()
	}

	h = g.NewGoal()

	if err := ae.AssignmentExpression.parse(&h, true, yield, await); err != nil {
		return j.Error("ArrayElement", err)
	}

	g.Score(h)
	j.Score(g)

	ae.Tokens = j.ToTokens()

	return nil
}

func (ae *ArrayElement) hasFirstComment() bool {
	if ae.Spread {
		return len(ae.Comments) > 0
	}

	return ae.AssignmentExpression.hasFirstComment()
}

// ArrayLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ArrayLiteral
type ArrayLiteral struct {
	ElementList []ArrayElement
	Comments    [2]Comments
	Tokens      Tokens
}

func (al *ArrayLiteral) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		return j.Error("ArrayLiteral", ErrMissingOpeningBracket)
	}

	al.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			j.Score(g)

			al.ElementList = append(al.ElementList, ArrayElement{})

			continue
		}

		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()

		var ae ArrayElement

		if err := ae.parse(&g, yield, await); err != nil {
			return j.Error("ArrayLiteral", err)
		}

		al.ElementList = append(al.ElementList, ae)

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("ArrayLiteral", ErrMissingComma)
		}

		j.Score(g)
	}

	al.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

	al.Tokens = j.ToTokens()

	return nil
}

// ObjectLiteral as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ObjectLiteral
type ObjectLiteral struct {
	PropertyDefinitionList []PropertyDefinition
	Comments               [2]Comments
	Tokens                 Tokens
}

func (ol *ObjectLiteral) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ObjectLiteral", ErrMissingOpeningBrace)
	}

	ol.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	for {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		}

		g = j.NewGoal()
		pd := len(ol.PropertyDefinitionList)
		ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, PropertyDefinition{})

		if err := ol.PropertyDefinitionList[pd].parse(&g, yield, await); err != nil {
			return j.Error("ObjectLiteral", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()

		if g.Accept(TokenRightBracePunctuator) {
			break
		} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			return g.Error("ObjectLiteral", ErrMissingComma)
		}

		j.Score(g)
	}

	ol.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

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
	Comments               [2]Comments
	Tokens                 Tokens
}

func (pd *PropertyDefinition) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
		pd.Comments[0] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
		j.Skip()
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()
		pd.AssignmentExpression = new(AssignmentExpression)

		if err := pd.AssignmentExpression.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyDefinition", err)
		}

		j.Score(g)
	} else {
		g = j.NewGoal()
		pd.Comments[0] = g.AcceptRunWhitespaceComments()

		g.AcceptRunWhitespace()

		h := g.NewGoal()

		if i := h.parseIdentifier(yield, await); i != nil {
			g.Score(h)

			pd.Comments[1] = g.AcceptRunWhitespaceCommentsInList()
			k := g.NewGoal()

			k.AcceptRunWhitespace()

			if k.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
				pd.PropertyName = &PropertyName{
					LiteralPropertyName: i,
					Tokens:              h.ToTokens(),
				}

				g.Score(k)
				g.AcceptRunWhitespaceNoComment()

				h = g.NewGoal()
				pd.AssignmentExpression = new(AssignmentExpression)

				if err := pd.AssignmentExpression.parse(&h, true, yield, await); err != nil {
					return g.Error("PropertyDefinition", err)
				}

				g.Score(h)

				pd.IsCoverInitializedName = true
			} else if t := k.Peek(); t.Type == TokenRightBracePunctuator || t == (parser.Token{Type: TokenPunctuator, Data: ","}) {
				pd.PropertyName = &PropertyName{
					LiteralPropertyName: i,
					Tokens:              h.ToTokens(),
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
			} else {
				pd.Comments = [2]Comments{}
			}
		} else {
			pd.Comments = [2]Comments{}
		}

		if pd.PropertyName == nil {
			g = j.NewGoal()
			propertyName := true

			g.AcceptRunWhitespace()

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
				pd.Comments[0] = g.AcceptRunWhitespaceComments()

				g.AcceptRunWhitespace()

				h := g.NewGoal()
				pd.PropertyName = new(PropertyName)

				if err := pd.PropertyName.parse(&h, yield, await); err != nil {
					return g.Error("PropertyDefinition", err)
				}

				g.Score(h)

				pd.Comments[1] = g.AcceptRunWhitespaceCommentsInList()
				h = g.NewGoal()

				h.AcceptRunWhitespace()

				if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
					h.AcceptRunWhitespaceNoComment()

					i := h.NewGoal()
					pd.AssignmentExpression = new(AssignmentExpression)

					if err := pd.AssignmentExpression.parse(&i, true, yield, await); err != nil {
						return h.Error("PropertyDefinition", err)
					}

					h.Score(i)
					g.Score(h)
				} else {
					propertyName = false
					pd.Comments = [2]Comments{}
				}
			}

			if !propertyName {
				pd.MethodDefinition = new(MethodDefinition)

				if pd.PropertyName != nil && pd.PropertyName.ComputedPropertyName != nil {
					pd.MethodDefinition.ClassElementName.PropertyName = pd.PropertyName
					pd.MethodDefinition.ClassElementName.Tokens = pd.PropertyName.Tokens
					pd.MethodDefinition.ClassElementName.Comments = pd.Comments
					pd.Comments = [2]Comments{}
				} else {
					g = j.NewGoal()
				}

				pd.PropertyName = nil

				if err := pd.MethodDefinition.parse(&g, false, false, yield, await); err != nil {
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
			j.AcceptRunWhitespaceNoComment()

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
	Comments             [5]Comments
	Tokens               Tokens
}

func (af *ArrowFunction) parse(j *jsParser, in, yield, await bool) error {
	af.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	af.Async = j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "async"})
	af.Comments[1] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoNewLine()

	g := j.NewGoal()

	if g.SkipGeneric() {
		af.Comments[1] = append(af.Comments[1], g.ToTypescriptComments()...)
		af.Comments[1] = append(af.Comments[1], g.AcceptRunWhitespaceNoNewlineComments()...)

		j.Score(g)
		j.AcceptRunWhitespaceNoNewLine()
	}

	if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
		g = j.NewGoal()
		af.FormalParameters = new(FormalParameters)

		if err := af.FormalParameters.parse(&g, yield, await); err != nil {
			return j.Error("ArrowFunction", err)
		}

		j.Score(g)

		af.Comments[2] = j.AcceptRunWhitespaceNoNewlineComments()

		j.AcceptRunWhitespaceNoNewLine()

		g = j.NewGoal()

		if g.SkipReturnType() {
			af.Comments[2] = append(af.Comments[2], g.ToTypescriptComments()...)
			af.Comments[2] = append(af.Comments[2], g.AcceptRunWhitespaceNoNewlineComments()...)

			j.Score(g)
			j.AcceptRunWhitespaceNoNewLine()
		}
	} else {
		g = j.NewGoal()

		if af.BindingIdentifier = g.parseIdentifier(yield, await); af.BindingIdentifier == nil {
			return j.Error("ArrowFunction", ErrNoIdentifier)
		}

		j.Score(g)

		af.Comments[2] = j.AcceptRunWhitespaceNoNewlineComments()
	}

	j.AcceptRunWhitespaceNoNewLine()

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "=>"}) {
		return j.Error("ArrowFunction", ErrMissingArrow)
	}

	af.Comments[3] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if j.Peek() == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
		g := j.NewGoal()
		af.FunctionBody = new(Block)

		if err := af.FunctionBody.parse(&g, false, af.Async, true); err != nil {
			return j.Error("ArrowFunction", err)
		}

		j.Score(g)

		af.Comments[4] = j.AcceptRunWhitespaceCommentsInList()
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

func (af *ArrowFunction) hasFirstComment() bool {
	return len(af.Comments[0]) > 0
}

func (af *ArrowFunction) hasLastComment() bool {
	if af.FunctionBody != nil {
		return len(af.Comments[0]) > 0
	}

	if af.AssignmentExpression != nil {
		return af.AssignmentExpression.hasLastComment()
	}

	return false
}
