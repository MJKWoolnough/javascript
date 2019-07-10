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
	j.AcceptToken(parser.Token{TokenKeyword, "class"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseIdentifier(yield, await)
	if err != nil {
		if !def {
			return j.Error("ClassDeclaration", err)
		}
	} else {
		j.Score(g)
		cd.BindingIdentifier = bi
		j.AcceptRunWhitespace()
	}
	if j.AcceptToken(parser.Token{TokenKeyword, "extends"}) {
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		cd.ClassHeritage = newLeftHandSideExpression()
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
		var md MethodDefinition
		if err := md.parse(&g, yield, await); err != nil {
			md.clear()
			return j.Error("ClassDeclaration", err)
		}
		j.Score(g)
		cd.ClassBody = append(cd.ClassBody, md)
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

func (md *MethodDefinition) parse(j *jsParser, yield, await bool) error {
	static := j.AcceptToken(parser.Token{TokenIdentifier, "static"})
	j.AcceptRunWhitespace()
	async := j.AcceptToken(parser.Token{TokenIdentifier, "async"})
	if async {
		if static {
			md.Type = MethodStaticAsync
		} else {
			md.Type = MethodAsync
		}
		j.AcceptRunWhitespaceNoNewLine()
	} else if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
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
		md.Type = MethodStatic
	}
	g := j.NewGoal()
	if err := md.PropertyName.parse(&g, yield, await); err != nil {
		return j.Error("MethodDefinition", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("MethodDefinition", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if md.Type != MethodGetter {
		g = j.NewGoal()
		if err := md.Params.parse(&g, md.Type == MethodGenerator, md.Type == MethodAsync); err != nil {
			return j.Error("MethodDefinition", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("MethodDefinition", ErrMissingClosingParenthesis)
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
	if j.Accept(TokenIdentifier, TokenStringLiteral, TokenNumericLiteral) {
		pn.LiteralPropertyName = j.GetLastToken()
	} else {
		g := j.NewGoal()
		pn.ComputedPropertyName = newAssignmentExpression()
		if err := pn.ComputedPropertyName.parse(&g, true, yield, await); err != nil {
			return j.Error("PropertyName", err)
		}
		j.Score(g)
	}
	pn.Tokens = j.ToTokens()
	return nil
}

var (
	ErrInvalidMethodName = errors.New("invalid method name")
)
