package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

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

type MethodType uint8

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

type MethodDefinition struct {
	Type         MethodType
	PropertyName PropertyName
	Params       FormalParameters
	FunctionBody Block
	Tokens       Tokens
}

func (md *MethodDefinition) parse(j *jsParser, pn *PropertyName, yield, await bool) error {
	static := j.AcceptToken(parser.Token{TokenIdentifier, "static"})
	j.AcceptRunWhitespace()
	async := j.AcceptToken(parser.Token{TokenIdentifier, "async"})
	if async {
		md.Type = MethodAsync
		j.AcceptRunWhitespaceNoNewLine()
	}
	if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		if static {
			if async {
				md.Type = MethodStaticAsyncGenerator
			} else {
				md.Type = MethodStaticGenerator
			}
		} else {
			if async {
				md.Type = MethodAsyncGenerator
			} else {
				md.Type = MethodGenerator
			}
		}
		j.AcceptRunWhitespace()
	} else if j.AcceptToken(parser.Token{TokenIdentifier, "get"}) {
		if static {
			md.Type = MethodStaticGetter
		} else {
			md.Type = MethodGetter
		}
		j.AcceptRunWhitespace()
	} else if j.AcceptToken(parser.Token{TokenIdentifier, "set"}) {
		if static {
			md.Type = MethodStaticSetter
		} else {
			md.Type = MethodSetter
		}
		j.AcceptRunWhitespace()
	} else if static {
		if async {
			md.Type = MethodStaticAsync
		} else {
			md.Type = MethodStatic
		}
	}
	g := j.NewGoal()
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
			return j.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return j.Error("MethodDefinition", ErrMissingClosingParenthesis)
		}
		md.Params.Tokens = g.ToTokens()
		j.Score(g)
	case MethodSetter, MethodStaticSetter:
		g := j.NewGoal()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "("}) {
			return j.Error("MethodDefinition", ErrMissingOpeningParenthesis)
		}
		md.Params.FormalParameterList = make([]BindingElement, 1)
		h := g.NewGoal()
		if err := md.Params.FormalParameterList[0].parse(&h, false, false); err != nil {
			return j.Error("MethodDefinition", err)
		}
		g.Score(h)
		if !g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			return j.Error("MethodDefinition", ErrMissingClosingParenthesis)
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
	} else if j.Accept(TokenIdentifier, TokenStringLiteral, TokenNumericLiteral) {
		pn.LiteralPropertyName = j.GetLastToken()
	} else {
		return j.Error("PropertyName", ErrInvalidPropertyName)
	}
	pn.Tokens = j.ToTokens()
	return nil
}

var (
	ErrInvalidMethodName       = errors.New("invalid method name")
	ErrInvalidPropertyName     = errors.New("invalid property name")
	ErrInvalidClassDeclaration = errors.New("invalid class declaration")
)
