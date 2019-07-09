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

func (j *jsParser) parseClassDeclaration(yield, await, def bool) (ClassDeclaration, error) {
	var cd ClassDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "class"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseIdentifier(yield, await)
	if err != nil {
		if !def {
			return cd, j.Error("ClassDeclaration", err)
		}
	} else {
		j.Score(g)
		cd.BindingIdentifier = bi
		j.AcceptRunWhitespace()
	}
	if j.AcceptToken(parser.Token{TokenKeyword, "extends"}) {
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		lhs, err := g.parseLeftHandSideExpression(yield, await)
		if err != nil {
			return cd, j.Error("ClassDeclaration", err)
		}
		j.Score(g)
		cd.ClassHeritage = &lhs
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return cd, j.Error("ClassDeclaration", ErrMissingOpeningBrace)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			continue
		} else if j.Accept(TokenRightBracePunctuator) {
			break
		}
		g := j.NewGoal()
		md, err := g.parseMethodDefinition(yield, await)
		if err != nil {
			return cd, j.Error("ClassDeclaration", err)
		}
		j.Score(g)
		cd.ClassBody = append(cd.ClassBody, md)
	}
	cd.Tokens = j.ToTokens()
	return cd, nil
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

func (j *jsParser) parseMethodDefinition(yield, await bool) (MethodDefinition, error) {
	var md MethodDefinition
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
	var err error
	if md.PropertyName, err = g.parsePropertyName(yield, await); err != nil {
		return md, j.Error("MethodDefinition", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return md, j.Error("MethodDefinition", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if md.Type != MethodGetter {
		g = j.NewGoal()
		if md.Params, err = g.parseFormalParameters(md.Type == MethodGenerator, md.Type == MethodAsync); err != nil {
			return md, j.Error("MethodDefinition", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return md, j.Error("MethodDefinition", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if md.FunctionBody, err = g.parseBlock(md.Type == MethodGenerator, md.Type == MethodAsync, true); err != nil {
		return md, j.Error("MethodDefinition", err)
	}
	j.Score(g)
	md.Tokens = j.ToTokens()
	return md, nil
}

type PropertyName struct {
	LiteralPropertyName  *Token
	ComputedPropertyName *AssignmentExpression
	Tokens               Tokens
}

func (j *jsParser) parsePropertyName(yield, await bool) (PropertyName, error) {
	var pn PropertyName
	if j.Accept(TokenIdentifier, TokenStringLiteral, TokenNumericLiteral) {
		pn.LiteralPropertyName = j.GetLastToken()
	} else {
		g := j.NewGoal()
		cp, err := g.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return pn, j.Error("PropertyName", err)
		}
		j.Score(g)
		pn.ComputedPropertyName = &cp
	}
	pn.Tokens = j.ToTokens()
	return pn, nil
}

const (
	ErrInvalidMethodName errors.Error = "invalid method name"
)
