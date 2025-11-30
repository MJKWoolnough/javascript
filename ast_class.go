package javascript

import (
	"slices"
	"strings"

	"vimagination.zapto.org/parser"
)

// ClassDeclaration as defined in ECMA-262
// https://tc39.es/ecma262/#prod-ClassDeclaration
//
// Also covers ClassExpression when BindingIdentifier is nil.
type ClassDeclaration struct {
	BindingIdentifier *Token
	ClassHeritage     *LeftHandSideExpression
	ClassBody         []ClassElement
	Comments          [5]Comments
	Tokens            Tokens
}

func (cd *ClassDeclaration) parse(j *jsParser, yield, await, def bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "class"}) {
		return j.Error("ClassDeclaration", ErrInvalidClassDeclaration)
	}

	cd.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if cd.BindingIdentifier = j.parseIdentifier(yield, await); cd.BindingIdentifier == nil {
		if !def {
			return j.Error("ClassDeclaration", ErrNoIdentifier)
		}
	} else {
		cd.Comments[1] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
	}

	if g := j.NewGoal(); g.SkipGeneric() {
		cd.Comments[1] = append(cd.Comments[1], g.ToTypescriptComments()...)
		cd.Comments[1] = append(cd.Comments[1], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	if j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "extends"}) {
		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		cd.ClassHeritage = new(LeftHandSideExpression)
		if err := cd.ClassHeritage.parse(&g, yield, await); err != nil {
			return j.Error("ClassDeclaration", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()

		if g = j.NewGoal(); g.SkipTypeArguments() {
			cd.Comments[2] = g.ToTypescriptComments()
			cd.Comments[2] = append(cd.Comments[2], g.AcceptRunWhitespaceComments()...)

			j.Score(g)
			j.AcceptRunWhitespace()
		}
	}

	if g := j.NewGoal(); g.SkipHeritage() {
		cd.Comments[2] = append(cd.Comments[2], g.ToTypescriptComments()...)
		cd.Comments[2] = append(cd.Comments[2], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		return j.Error("ClassDeclaration", ErrMissingOpeningBrace)
	}

	cd.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	var c Comments

	g := j.NewGoal()

	for {
		h := g.NewGoal()

		h.AcceptRunWhitespace()

		if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
			c = append(c, g.AcceptRunWhitespaceComments()...)

			g.AcceptRunWhitespace()
			g.Skip()
			g.AcceptRunWhitespaceNoComment()

			if len(c) == 0 {
				j.Score(g)

				g = j.NewGoal()
			}

			continue
		} else if h.Accept(TokenRightBracePunctuator) {
			break
		}

		h = g.NewGoal()

		h.AcceptRunWhitespace()

		i := h.NewGoal()

		if i.SkipParameterProperties() {
			i.AcceptRunWhitespace()
		}

		if i.SkipAbstractField() || i.SkipIndexSignature() {
			i.parseSemicolon()

			c = append(c, i.ToTypescriptComments()...)

			h.Score(i)
			g.Score(h)

			continue
		}

		g.AcceptRunWhitespaceNoComment()

		if len(c) == 0 {
			j.Score(g)

			g = j.NewGoal()
		}

		md := len(cd.ClassBody)

		cd.ClassBody = append(cd.ClassBody, ClassElement{Comments: [3]Comments{c}})
		if err := cd.ClassBody[md].parse(&g, yield, await); err != nil {
			return j.Error("ClassDeclaration", err)
		}

		c = nil

		j.Score(g)
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()
	}

	j.Score(g)

	cd.Comments[4] = append(c, j.AcceptRunWhitespaceComments()...)

	j.AcceptRunWhitespace()
	j.Skip()

	cd.Tokens = j.ToTokens()

	return nil
}

// ClassElement as defined in ECMA-262
// https://tc39.es/ecma262/#prod-ClassElement
//
// Only one of MethodDefinition, FieldDefinition, or ClassStaticBlock must be
// non-nil.
//
// If ClassStaticBlock is non-nil, Static should be true
type ClassElement struct {
	Static           bool
	MethodDefinition *MethodDefinition
	FieldDefinition  *FieldDefinition
	ClassStaticBlock *Block
	Comments         [3]Comments
	Tokens           Tokens
}

func (ce *ClassElement) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.SkipParameterProperties() {
		ce.Comments[0] = append(ce.Comments[0], j.AcceptRunWhitespaceComments()...)

		j.AcceptRunWhitespace()

		g = j.NewGoal()

		g.SkipParameterProperties()

		ce.Comments[0] = append(ce.Comments[0], g.ToTypescriptComments()...)

		j.Score(g)
		j.AcceptRunWhitespaceNoComment()
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "static"}) {
		g.AcceptRunWhitespace()

		if tk := g.Peek(); (tk.Type != TokenPunctuator || (tk.Data != "=" && tk.Data != ";" && tk.Data != "(" && tk.Data != ":")) && tk.Type != TokenRightBracePunctuator {
			ce.Static = true
			ce.Comments[0] = append(ce.Comments[0], j.AcceptRunWhitespaceComments()...)

			j.AcceptRunWhitespace()
			j.Skip()
			j.AcceptRunWhitespaceNoComment()
		}
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if g.SkipReadOnly() {
		ce.Comments[1] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		g = j.NewGoal()

		g.SkipReadOnly()

		ce.Comments[1] = append(ce.Comments[1], g.ToTypescriptComments()...)

		j.Score(g)

		j.AcceptRunWhitespaceNoComment()
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if ce.Static && g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		ce.Comments[1] = append(ce.Comments[1], j.AcceptRunWhitespaceComments()...)

		j.AcceptRunWhitespace()

		g = j.NewGoal()

		ce.ClassStaticBlock = new(Block)
		if err := ce.ClassStaticBlock.parse(&g, false, true, false); err != nil {
			return j.Error("ClassElement", err)
		}

		j.Score(g)

		ce.Comments[2] = j.AcceptRunWhitespaceCommentsInList()
	} else {
		g = j.NewGoal()
		h := g.NewGoal()

		var (
			cen      ClassElementName
			isMethod bool
		)

		if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
			isMethod = true
		} else if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "async"}) {
			h.AcceptRunWhitespaceNoNewLine()

			tk := h.Peek()

			isMethod = tk == (parser.Token{Type: TokenPunctuator, Data: "*"}) || tk == (parser.Token{Type: TokenPunctuator, Data: "["}) || tk == (parser.Token{Type: TokenPunctuator, Data: "("}) || tk.Type == TokenIdentifier || tk.Type == TokenPrivateIdentifier
			if !isMethod && h.SkipGeneric() {
				h.AcceptRunWhitespace()

				tk := h.Peek()
				isMethod = tk == (parser.Token{Type: TokenPunctuator, Data: "("})
			}
		} else if h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "get"}) || h.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "set"}) {
			h.AcceptRunWhitespace()

			tk := h.Peek()
			isMethod = tk == (parser.Token{Type: TokenPunctuator, Data: "["}) || tk == (parser.Token{Type: TokenPunctuator, Data: "("}) || tk.Type == TokenIdentifier || tk.Type == TokenPrivateIdentifier
		} else {
			if err := cen.parse(&g, yield, await); err != nil {
				return j.Error("ClassElement", err)
			}

			h = g.NewGoal()

			h.AcceptRunWhitespace()

			if h.SkipGeneric() {
				h.AcceptRunWhitespace()
			}

			isMethod = h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("})
		}

		if isMethod {
			ce.MethodDefinition = &MethodDefinition{ClassElementName: cen}

			for {
				if h = g.NewGoal(); h.SkipMethodOverload(ce.Static, cen, yield, await) {
					g.Score(h)
				} else {
					break
				}
			}

			if err := ce.MethodDefinition.parse(&g, yield, await); err != nil {
				return j.Error("ClassElement", err)
			}
		} else {
			ce.FieldDefinition = &FieldDefinition{ClassElementName: cen}
			if err := ce.FieldDefinition.parse(&g, yield, await); err != nil {
				return j.Error("ClassElement", err)
			}

			if g.GetLastToken().Token != (parser.Token{Type: TokenPunctuator, Data: ";"}) {
				var hasNewline bool

			Loop:
				for _, tk := range slices.Backward(ce.FieldDefinition.Tokens) {
					switch tk.Type {
					case TokenSingleLineComment, TokenLineTerminator:
						hasNewline = true

						break Loop
					case TokenMultiLineComment:
						if strings.Contains(tk.Data, "\n") {
							hasNewline = true

							break Loop
						}
					case TokenWhitespace:
					default:
						break Loop
					}
				}

				h := g.NewGoal()

				h.AcceptRunWhitespaceNoNewLine()

				if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
					g.Score(h)
				} else if h.Accept(TokenLineTerminator, TokenSingleLineComment, TokenMultiLineComment) || hasNewline {
					h.AcceptRunWhitespace()

					if h.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ";"}) {
						g.Score(h)
					}
				} else if h.Peek() != (parser.Token{Type: TokenRightBracePunctuator, Data: "}"}) {
					return h.Error("ClassElement", ErrMissingSemiColon)
				}
			}
		}

		j.Score(g)
	}

	ce.Tokens = j.ToTokens()

	return nil
}

// ClassElementName as defined in ECMA-262
// https://tc39.es/ecma262/#prod-ClassElementName
//
// Only one of PropertyName or PrivateIdentifier must be non-nil
type ClassElementName struct {
	PropertyName      *PropertyName
	PrivateIdentifier *Token
	Comments          [2]Comments
	Tokens            Tokens
}

func (cen *ClassElementName) parse(j *jsParser, yield, await bool) error {
	cen.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if j.Accept(TokenPrivateIdentifier) {
		cen.PrivateIdentifier = j.GetLastToken()
	} else {
		g := j.NewGoal()

		cen.PropertyName = new(PropertyName)
		if err := cen.PropertyName.parse(&g, yield, await); err != nil {
			return j.Error("ClassElementName", err)
		}

		j.Score(g)
	}

	cen.Comments[1] = j.AcceptRunWhitespaceCommentsInList()
	cen.Tokens = j.ToTokens()

	return nil
}

// FieldDefinition as defined in ECMA-262
// https://tc39.es/ecma262/#prod-FieldDefinition
type FieldDefinition struct {
	ClassElementName ClassElementName
	Initializer      *AssignmentExpression
	Comments         Comments
	Tokens           Tokens
}

func (fd *FieldDefinition) parse(j *jsParser, yield, await bool) error {
	if len(fd.ClassElementName.Tokens) == 0 {
		fd.ClassElementName.parse(j, yield, await)
	}

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	h := g.NewGoal()

	if h.SkipForce() {
		i := h.NewGoal()

		i.AcceptRunWhitespace()

		if i.SkipColonType() {
			h.Score(i)
		}

		fd.Comments = h.ToTypescriptComments()

		g.Score(h)
		j.Score(g)

		fd.Comments = append(fd.Comments, j.AcceptRunWhitespaceCommentsInList()...)

		g = j.NewGoal()

		g.AcceptRunWhitespace()
	} else if h.SkipOptionalColonType() {
		fd.Comments = h.ToTypescriptComments()

		g.Score(h)
		j.Score(g)

		fd.Comments = append(fd.Comments, j.AcceptRunWhitespaceCommentsInList()...)

		g = j.NewGoal()

		g.AcceptRunWhitespace()
	}

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		g.AcceptRunWhitespaceNoComment()

		h := g.NewGoal()

		fd.Initializer = new(AssignmentExpression)
		if err := fd.Initializer.parse(&h, true, yield, await); err != nil {
			return g.Error("FieldDefinition", err)
		}

		g.Score(h)
		j.Score(g)
	}

	fd.Tokens = j.ToTokens()

	return nil
}

// MethodType determines the prefixes for MethodDefinition
type MethodType uint8

// Valid MethodType's
const (
	MethodNormal MethodType = iota
	MethodGenerator
	MethodAsync
	MethodAsyncGenerator
	MethodGetter
	MethodSetter
)

// MethodDefinition as specified in ECMA-262
// https://tc39.es/ecma262/#prod-MethodDefinition
type MethodDefinition struct {
	Type             MethodType
	ClassElementName ClassElementName
	Params           FormalParameters
	FunctionBody     Block
	Comments         [4]Comments
	Tokens           Tokens
}

func (md *MethodDefinition) parse(j *jsParser, yield, await bool) error {
	var prev MethodType

	if len(md.ClassElementName.Tokens) == 0 {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		switch g.Peek() {
		case parser.Token{Type: TokenIdentifier, Data: "get"}:
			g.Skip()

			md.Type = MethodGetter
		case parser.Token{Type: TokenIdentifier, Data: "set"}:
			g.Skip()

			md.Type = MethodSetter
		case parser.Token{Type: TokenIdentifier, Data: "async"}:
			g.Skip()

			if t := g.AcceptRunWhitespaceNoNewLine(); t == TokenLineTerminator || t == TokenSingleLineComment || t == TokenMultiLineComment || g.SkipGeneric() {
				break
			}

			g.AcceptRunWhitespace()

			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
				md.Type = MethodAsyncGenerator
				prev = MethodAsyncGenerator
			} else {
				md.Type = MethodAsync
			}
		case parser.Token{Type: TokenPunctuator, Data: "*"}:
			g.Skip()

			md.Type = MethodGenerator
			prev = MethodGenerator
		}

		g.AcceptRunWhitespace()

		if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "("}) {
			md.Type = prev
		} else if md.Type != MethodNormal {
			md.Comments[0] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()
			j.AcceptRunWhitespaceNoComment()

			if md.Type == MethodAsyncGenerator {
				md.Comments[1] = append(md.Comments[1], j.AcceptRunWhitespaceComments()...)

				j.AcceptRunWhitespace()
				j.Skip()
				j.AcceptRunWhitespaceNoComment()
			}
		}

		g = j.NewGoal()

		if err := md.ClassElementName.parse(&g, yield, await); err != nil {
			return j.Error("MethodDefinition", err)
		}

		j.Score(g)
	}

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if g.SkipGeneric() {
		md.ClassElementName.Comments[1] = append(md.ClassElementName.Comments[1], g.ToTypescriptComments()...)
		md.ClassElementName.Comments[1] = append(md.ClassElementName.Comments[1], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	switch md.Type {
	case MethodGetter:
		g := j.NewGoal()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
			return g.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}

		md.Params.Comments[0] = g.AcceptRunWhitespaceNoNewlineComments()

		g.AcceptRunWhitespaceNoComment()

		md.Params.Comments[4] = g.AcceptRunWhitespaceComments()

		g.AcceptRunWhitespace()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			return g.Error("MethodDefinition", ErrMissingClosingParenthesis)
		}

		md.Params.Tokens = g.ToTokens()

		j.Score(g)
	case MethodSetter:
		g := j.NewGoal()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
			return g.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}

		md.Params.Comments[0] = g.AcceptRunWhitespaceNoNewlineComments()

		g.AcceptRunWhitespaceNoComment()

		md.Params.FormalParameterList = make([]BindingElement, 1)
		h := g.NewGoal()

		if err := md.Params.FormalParameterList[0].parse(&h, nil, false, false); err != nil {
			return j.Error("MethodDefinition", err)
		}

		g.Score(h)

		md.Params.Comments[4] = g.AcceptRunWhitespaceComments()

		g.AcceptRunWhitespace()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
			return g.Error("MethodDefinition", ErrMissingClosingParenthesis)
		}

		md.Params.Tokens = g.ToTokens()

		j.Score(g)
	default:
		g := j.NewGoal()

		if err := md.Params.parse(&g, md.Type == MethodGenerator, md.Type == MethodAsync); err != nil {
			return j.Error("MethodDefinition", err)
		}

		j.Score(g)
	}

	md.Comments[2] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if g.SkipReturnType() {
		md.Comments[2] = append(md.Comments[2], g.ToTypescriptComments()...)
		md.Comments[2] = append(md.Comments[2], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	g = j.NewGoal()
	if err := md.FunctionBody.parse(&g, md.Type == MethodGenerator, md.Type == MethodAsync, true); err != nil {
		return j.Error("MethodDefinition", err)
	}

	j.Score(g)

	md.Comments[3] = j.AcceptRunWhitespaceCommentsInList()
	md.Tokens = j.ToTokens()

	return nil
}

// PropertyName as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-PropertyName
//
// Only one of LiteralPropertyName or ComputedPropertyName must be non-nil.
type PropertyName struct {
	LiteralPropertyName  *Token
	ComputedPropertyName *AssignmentExpression
	Comments             [2]Comments
	Tokens               Tokens
}

func (pn *PropertyName) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "["}) {
		pn.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		pn.ComputedPropertyName = new(AssignmentExpression)
		if err := pn.ComputedPropertyName.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyName", err)
		}

		j.Score(g)

		pn.Comments[1] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "]"}) {
			return j.Error("PropertyName", ErrMissingClosingBracket)
		}
	} else if j.Accept(TokenIdentifier, TokenKeyword, TokenStringLiteral, TokenNumericLiteral) {
		pn.LiteralPropertyName = j.GetLastToken()
	} else {
		return j.Error("PropertyName", ErrInvalidPropertyName)
	}

	pn.Tokens = j.ToTokens()

	return nil
}
