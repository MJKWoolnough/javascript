package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Script struct {
	StatementList []StatementListItem
	Tokens        Tokens
}

func ParseScript(t parser.Tokeniser) (*Script, error) {
	j, err := newJSParser(t)
	if err != nil {
		return nil, err
	}
	s := newScript()
	if err := s.parse(&j); err != nil {
		s.clear()
		return nil, err
	}
	return s, nil
}

func (s *Script) parse(j *jsParser) error {
	for j.AcceptRunWhitespace() != parser.TokenDone {
		g := j.NewGoal()
		var si StatementListItem
		if err := si.parse(&g, false, false, false); err != nil {
			return err
		}
		j.Score(g)
		s.StatementList = append(s.StatementList, si)
	}
	s.Tokens = j.ToTokens()
	return nil
}

func (j *jsParser) parseIdentifier(yield, await bool) (*Token, error) {
	if j.Accept(TokenIdentifier) || (!yield && j.AcceptToken(parser.Token{TokenKeyword, "yield"}) || (!await && j.AcceptToken(parser.Token{TokenKeyword, "await"}))) {
		return j.GetLastToken(), nil
	}
	return nil, j.Error("Identifier", ErrNoIdentifier)
}

type Declaration struct {
	ClassDeclaration    *ClassDeclaration
	FunctionDeclaration *FunctionDeclaration
	LexicalDeclaration  *LexicalDeclaration
	Tokens              Tokens
}

func (d *Declaration) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenKeyword, "class"}) {
		d.ClassDeclaration = newClassDeclaration()
		if err := d.ClassDeclaration.parse(&g, yield, await, false); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk := g.Peek(); tk == (parser.Token{TokenKeyword, "const"}) || tk == (parser.Token{TokenIdentifier, "let"}) {
		d.LexicalDeclaration = newLexicalDeclaration()
		if err := d.LexicalDeclaration.parse(&g, true, yield, await); err != nil {
			return j.Error("Declaration", err)
		}
	} else if tk == (parser.Token{TokenIdentifier, "async"}) || tk == (parser.Token{TokenKeyword, "function"}) {
		d.FunctionDeclaration = newFunctionDeclaration()
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

type LetOrConst bool

const (
	Let   LetOrConst = false
	Const LetOrConst = true
)

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
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		}
		g := j.NewGoal()
		var lb LexicalBinding
		if err := lb.parse(&g, in, yield, await); err != nil {
			lb.clear()
			return j.Error("LexicalDeclaration", err)
		}
		j.Score(g)
		ld.BindingList = append(ld.BindingList, lb)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("LexicalDeclaration", ErrInvalidLexicalDeclaration)
		}
	}
	ld.Tokens = j.ToTokens()
	return nil
}

type LexicalBinding struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}

func (lb *LexicalBinding) parse(j *jsParser, in, yield, await bool) error {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		lb.ArrayBindingPattern = newArrayBindingPattern()
		if err := lb.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		lb.ObjectBindingPattern = newObjectBindingPattern()
		if err := lb.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
	} else {
		bi, err := g.parseIdentifier(yield, await)
		if err != nil {
			return j.Error("LexicalBinding", err)
		}
		lb.BindingIdentifier = bi
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		lb.Initializer = newAssignmentExpression()
		if err := lb.Initializer.parse(&g, in, yield, await); err != nil {
			return j.Error("LexicalBinding", err)
		}
		j.Score(g)
	}
	lb.Tokens = j.ToTokens()
	return nil
}

type ArrayBindingPattern struct {
	BindingElementList []BindingElement
	BindingRestElement *BindingElement
	Tokens             Tokens
}

func (ab *ArrayBindingPattern) parse(j *jsParser, yield, await bool) error {
	j.AcceptToken(parser.Token{TokenPunctuator, "["})
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
			ab.BindingRestElement = newBindingElement()
			if err := ab.BindingRestElement.parse(&g, yield, await); err != nil {
				return j.Error("ArrayBindingPattern", err)
			}
			j.Score(g)
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return j.Error("ArrayBindingPattern", ErrMissingClosingBracket)
			}
			break
		}
		var be BindingElement
		if err := be.parse(&g, yield, await); err != nil {
			be.clear()
			return j.Error("ArrayBindingPattern", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		ab.BindingElementList = append(ab.BindingElementList, be)
		if j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("ArrayBindingPattern", ErrMissingComma)
		}
	}
	ab.Tokens = j.ToTokens()
	return nil
}

type ObjectBindingPattern struct {
	BindingPropertyList []BindingProperty
	BindingRestProperty *Token
	Tokens              Tokens
}

func (ob *ObjectBindingPattern) parse(j *jsParser, yield, await bool) error {
	j.AcceptToken(parser.Token{TokenPunctuator, "{"})
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		for {
			g := j.NewGoal()
			if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
				bi, err := g.parseIdentifier(yield, await)
				if err != nil {
					return j.Error("ObjectBindingPattern", err)
				}
				j.Score(g)
				ob.BindingRestProperty = bi
				j.AcceptRunWhitespace()
				if !j.Accept(TokenRightBracePunctuator) {
					return j.Error("ObjectBindingPattern", ErrMissingClosingBrace)
				}
				break
			}
			var bp BindingProperty
			if err := bp.parse(&g, yield, await); err != nil {
				bp.clear()
				return j.Error("ObjectBindingPattern", err)
			}
			j.Score(g)
			ob.BindingPropertyList = append(ob.BindingPropertyList, bp)
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

type BindingProperty struct {
	SingleNameBinding *Token
	Initializer       *AssignmentExpression
	PropertyName      *PropertyName
	BindingElement    *BindingElement
	Tokens            Tokens
}

func (bp *BindingProperty) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if i, err := g.parseIdentifier(yield, await); err == nil {
		h := g.NewGoal()
		h.AcceptRunWhitespace()
		if h.AcceptToken(parser.Token{TokenPunctuator, "="}) {
			g.Score(h)
			g.AcceptRunWhitespace()
			h = g.NewGoal()
			bp.Initializer = newAssignmentExpression()
			if err := bp.Initializer.parse(&h, true, yield, await); err != nil {
				return g.Error("BindingProperty", err)
			}
			g.Score(h)
			bp.SingleNameBinding = i
		} else if !h.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
			bp.SingleNameBinding = i
		}
	}
	if bp.SingleNameBinding == nil {
		g = j.NewGoal()
		bp.PropertyName = newPropertyName()
		if err := bp.PropertyName.parse(&g, yield, await); err != nil {
			return j.Error("BindingProperty", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
			return j.Error("BindingProperty", ErrMissingColon)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		bp.BindingElement = newBindingElement()
		if err := bp.BindingElement.parse(&g, yield, await); err != nil {
			return j.Error("BindingProperty", err)
		}
	}
	j.Score(g)
	bp.Tokens = j.ToTokens()
	return nil
}

type VariableDeclaration LexicalBinding

func (v *VariableDeclaration) parse(j *jsParser, in, yield, await bool) error {
	return ((*LexicalBinding)(v)).parse(j, in, yield, await)
}

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
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			g.AcceptRunWhitespace()
			al.SpreadElement = newAssignmentExpression()
			if err := al.SpreadElement.parse(&g, true, yield, await); err != nil {
				return j.Error("ArrayLiteral", err)
			}
			j.Score(g)
			if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
				return j.Error("ArrayLiteral", ErrMissingClosingBracket)
			}
			break
		}
		var ae AssignmentExpression
		if err := ae.parse(&g, true, yield, await); err != nil {
			ae.clear()
			return j.Error("ArrayLiteral", err)
		}
		j.Score(g)
		al.ElementList = append(al.ElementList, ae)
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
		var pd PropertyDefinition
		if err := pd.parse(&g, yield, await); err != nil {
			return j.Error("ObjectLiteral", err)
		}
		j.Score(g)
		ol.PropertyDefinitionList = append(ol.PropertyDefinitionList, pd)
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

type PropertyDefinition struct {
	IdentifierReference  *Token
	PropertyName         *PropertyName
	Spread               bool
	AssignmentExpression *AssignmentExpression
	MethodDefinition     *MethodDefinition
	Tokens               Tokens
}

func (pd *PropertyDefinition) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		pd.AssignmentExpression = newAssignmentExpression()
		if err := pd.AssignmentExpression.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyDefinition", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		if i, err := g.parseIdentifier(yield, await); err == nil {
			h := g.NewGoal()
			h.AcceptRunWhitespace()
			if h.AcceptToken(parser.Token{TokenPunctuator, "="}) {
				g.Score(h)
				g.AcceptRunWhitespace()
				pd.IdentifierReference = i
				h = g.NewGoal()
				pd.AssignmentExpression = newAssignmentExpression()
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
				g.Except()
				g.AcceptRunWhitespace()
				if !g.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
					propertyName = false
				}
			}
			g = j.NewGoal()
			if propertyName {
				pd.PropertyName = newPropertyName()
				if err := pd.PropertyName.parse(&g, yield, await); err != nil {
					return j.Error("PropertyDefinition", err)
				}
				h := g.NewGoal()
				h.AcceptRunWhitespace()
				if h.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
					h.AcceptRunWhitespace()
					i := h.NewGoal()
					pd.AssignmentExpression = newAssignmentExpression()
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
				pd.MethodDefinition = newMethodDefinition()
				if err := pd.MethodDefinition.parse(&g, pd.PropertyName, yield, await); err != nil {
					j.Error("PropertyDefinition", err)
				}
				pd.PropertyName = nil
			}
		}
		j.Score(g)
	}
	pd.Tokens = j.ToTokens()
	return nil
}

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
			var e Expression
			if err := e.parse(&g, true, yield, await); err != nil {
				e.clear()
				return j.Error("TemplateLiteral", err)
			}
			j.Score(g)
			tl.Expressions = append(tl.Expressions, e)
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
			j.Error("ArrowFunction", ErrInvalidAsyncArrowFunction)
		}
		af.Async = true
		j.AcceptRunWhitespaceNoNewLine()
		if j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			g := j.NewGoal()
			af.FormalParameters = newFormalParameters()
			if err := af.FormalParameters.parse(&g, false, true); err != nil {
				return j.Error("ArrowFunction", err)
			}
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return j.Error("ArrowFunction", ErrMissingClosingParenthesis)
			}
		} else {
			g := j.NewGoal()
			bi, err := g.parseIdentifier(yield, true)
			if err != nil {
				return j.Error("ArrowFunction", err)
			}
			j.Score(g)
			af.BindingIdentifier = bi
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
		af.FunctionBody = newBlock()
		if err := af.FunctionBody.parse(&g, false, af.Async, true); err != nil {
			return j.Error("ArrowFunction", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		af.AssignmentExpression = newAssignmentExpression()
		if err := af.AssignmentExpression.parse(&g, in, false, af.Async); err != nil {
			return j.Error("ArrowFunction", err)
		}
		j.Score(g)
	}
	af.Tokens = j.ToTokens()
	return nil
}

var (
	ErrReservedIdentifier          = errors.New("reserved identifier")
	ErrNoIdentifier                = errors.New("missing identifier")
	ErrMissingFunction             = errors.New("missing function")
	ErrMissingOpeningParenthesis   = errors.New("missing opening parenthesis")
	ErrMissingClosingParenthesis   = errors.New("missing closing parenthesis")
	ErrMissingOpeningBrace         = errors.New("missing opening brace")
	ErrMissingClosingBrace         = errors.New("missing closing brace")
	ErrMissingOpeningBracket       = errors.New("missing opening bracket")
	ErrMissingClosingBracket       = errors.New("missing closing bracket")
	ErrMissingComma                = errors.New("missing comma")
	ErrMissingArrow                = errors.New("missing arrow")
	ErrMissingCaseClause           = errors.New("missing case clause")
	ErrMissingExpression           = errors.New("missing expression")
	ErrMissingCatchFinally         = errors.New("missing catch/finally block")
	ErrMissingSemiColon            = errors.New("missing semi-colon")
	ErrMissingColon                = errors.New("missing colon")
	ErrInvalidStatementList        = errors.New("invalid statement list")
	ErrInvalidStatement            = errors.New("invalid statement")
	ErrInvalidFormalParameterList  = errors.New("invalid formal parameter list")
	ErrInvalidDeclaration          = errors.New("invalid declaration")
	ErrInvalidLexicalDeclaration   = errors.New("invalid lexical declaration")
	ErrInvalidAssignment           = errors.New("invalid assignment operator")
	ErrInvalidSuperProperty        = errors.New("invalid super property")
	ErrInvalidMetaProperty         = errors.New("invalid meta property")
	ErrInvalidTemplate             = errors.New("invalid template")
	ErrInvalidIterationStatementDo = errors.New("invalid do interation statement")
	ErrInvalidForLoop              = errors.New("invalid for loop")
	ErrInvalidAsyncArrowFunction   = errors.New("invalid async arrow function")
)
