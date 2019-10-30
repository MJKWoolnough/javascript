package javascript

import (
	"errors"

	"vimagination.zapto.org/parser"
)

// ClassDeclaration as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-ClassDeclaration
//
// Also covers ClassExpression when BindingIdentifier is nil.
type ClassDeclaration struct {
	BindingIdentifier *Token
	ClassHeritage     *LeftHandSideExpression
	ClassBody         []MethodDefinition
	Tokens            Tokens
}

func (cd *ClassDeclaration) parse(j *jsParser, yield, await, def bool) error {
	if !j.AcceptToken(parser.Token{TokenKeyword, "class"}) {
		return j.Error("ClassDeclaration", ErrInvalidClassDeclaration)
	}
	j.AcceptRunWhitespace()
	if cd.BindingIdentifier = j.parseIdentifier(yield, await); cd.BindingIdentifier == nil {
		if !def {
			return j.Error("ClassDeclaration", ErrNoIdentifier)
		}
	} else {
		j.AcceptRunWhitespace()
	}
	if j.AcceptToken(parser.Token{TokenKeyword, "extends"}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		cd.ClassHeritage = new(LeftHandSideExpression)
		if err := cd.ClassHeritage.parse(&g, yield, await); err != nil {
			return j.Error("ClassDeclaration", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return j.Error("ClassDeclaration", ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			continue
		} else if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		md := len(cd.ClassBody)
		cd.ClassBody = append(cd.ClassBody, MethodDefinition{})
		if err := cd.ClassBody[md].parse(&g, nil, yield, await); err != nil {
			return j.Error("ClassDeclaration", err)
		}
		j.Score(g)
	}
	cd.Tokens = j.ToTokens()
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
	MethodStatic
	MethodStaticGenerator
	MethodStaticAsync
	MethodStaticAsyncGenerator
	MethodStaticGetter
	MethodStaticSetter
)

// MethodDefinition as specified in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-MethodDefinition
//
// Static methods from ClassElement are parsed here with the `static` prefix
type MethodDefinition struct {
	Type         MethodType
	PropertyName PropertyName
	Params       FormalParameters
	FunctionBody Block
	Tokens       Tokens
}

func (md *MethodDefinition) parse(j *jsParser, pn *PropertyName, yield, await bool) error {
	var prev MethodType
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenIdentifier, "static"}) {
		md.Type = MethodStatic
		g.AcceptRunWhitespace()
	}
	switch g.Peek() {
	case parser.Token{TokenIdentifier, "get"}:
		j.Score(g)
		g = j.NewGoal()
		g.Skip()
		g.AcceptRunWhitespace()
		prev = md.Type
		switch md.Type {
		case MethodNormal:
			md.Type = MethodGetter
		case MethodStatic:
			md.Type = MethodStaticGetter
		}
	case parser.Token{TokenIdentifier, "set"}:
		j.Score(g)
		g = j.NewGoal()
		g.Skip()
		g.AcceptRunWhitespace()
		prev = md.Type
		switch md.Type {
		case MethodNormal:
			md.Type = MethodSetter
		case MethodStatic:
			md.Type = MethodStaticSetter
		}
	case parser.Token{TokenIdentifier, "async"}:
		j.Score(g)
		g = j.NewGoal()
		g.Skip()
		if t := g.AcceptRunWhitespaceNoNewLine(); t == TokenLineTerminator || t == TokenSingleLineComment || t == TokenMultiLineComment {
			g = j.NewGoal()
			break
		}
		prev = md.Type
		switch md.Type {
		case MethodNormal:
			md.Type = MethodAsync
		case MethodStatic:
			md.Type = MethodStaticAsync
		}
		fallthrough
	default:
		if g.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
			j.Score(g)
			j.AcceptRunWhitespace()
			g = j.NewGoal()
			switch md.Type {
			case MethodNormal:
				md.Type = MethodGenerator
			case MethodStatic:
				md.Type = MethodStaticGenerator
			case MethodAsync:
				md.Type = MethodAsyncGenerator
			case MethodStaticAsync:
				md.Type = MethodStaticAsyncGenerator
			}
			prev = md.Type
		}
	}
	if g.Peek() == (parser.Token{TokenPunctuator, "("}) {
		md.Type = prev
	} else {
		j.Score(g)
	}
	g = j.NewGoal()
	if pn != nil {
		md.PropertyName = *pn
	} else if err := md.PropertyName.parse(&g, yield, await); err != nil {
		return j.Error("MethodDefinition", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	switch md.Type {
	case MethodGetter, MethodStaticGetter:
		g := j.NewGoal()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			return g.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return g.Error("MethodDefinition", ErrMissingClosingParenthesis)
		}
		md.Params.Tokens = g.ToTokens()
		j.Score(g)
	case MethodSetter, MethodStaticSetter:
		g := j.NewGoal()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			return g.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}
		g.AcceptRunWhitespace()
		md.Params.FormalParameterList = make([]BindingElement, 1)
		h := g.NewGoal()
		if err := md.Params.FormalParameterList[0].parse(&h, false, false); err != nil {
			return j.Error("MethodDefinition", err)
		}
		g.Score(h)
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return g.Error("MethodDefinition", ErrMissingClosingParenthesis)
		}
		md.Params.Tokens = g.ToTokens()
		j.Score(g)
	default:
		g = j.NewGoal()
		if err := md.Params.parse(&g, md.Type == MethodGenerator, md.Type == MethodAsync); err != nil {
			return j.Error("MethodDefinition", err)
		}
		j.Score(g)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := md.FunctionBody.parse(&g, md.Type == MethodGenerator, md.Type == MethodAsync, true); err != nil {
		return j.Error("MethodDefinition", err)
	}
	j.Score(g)
	md.Tokens = j.ToTokens()
	return nil
}

// PropertyName as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-PropertyName
//
// Only one of LiteralPropertyName or ComputedPropertyName must be non-nil.
type PropertyName struct {
	LiteralPropertyName  *Token
	ComputedPropertyName *AssignmentExpression
	Tokens               Tokens
}

func (pn *PropertyName) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		pn.ComputedPropertyName = new(AssignmentExpression)
		if err := pn.ComputedPropertyName.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyName", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, "]"}) {
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

// Errors
var (
	ErrInvalidMethodName       = errors.New("invalid method name")
	ErrInvalidPropertyName     = errors.New("invalid property name")
	ErrInvalidClassDeclaration = errors.New("invalid class declaration")
)
