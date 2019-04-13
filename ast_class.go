package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type ClassDeclaration struct {
	BindingIdentifier *BindingIdentifier
	Extends           *LeftHandSideExpression
	ClassBody         ClassBody
	Tokens            Tokens
}

func (j *jsParser) parseClassDeclaration(yield, await, def bool) (ClassDeclaration, error) {
	var cd ClassDeclaration
	j.AcceptToken(parser.Token{TokenKeyword, "class"})
	j.AcceptRunWhitespace()
	g := j.NewGoal()
	bi, err := g.parseBindingIdentifier(yield, await)
	if err != nil {
		if !def {
			return cd, j.Error(err)
		}
	} else {
		j.Score(g)
		cd.BindingIdentifier = &bi
		j.AcceptRunWhitespace()
	}
	if j.AcceptToken(parser.Token{TokenKeyword, "extends"}) {
		g = j.NewGoal()
		lhs, err := g.parseLeftHandSideExpression(yield, await)
		if err != nil {
			return cd, j.Error(err)
		}
		j.Score(g)
		cd.Extends = &lhs
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return cd, j.Error(ErrMissingSemiColon)
	}
	g = j.NewGoal()
	cd.ClassBody, err = g.parseClassBody(yield, await)
	if err != nil {
		return cd, j.Error(err)
	}
	j.Score(g)
	if !j.Accept(TokenRightBracePunctuator) {
		return cd, j.Error(ErrMissingClosingBrace)
	}
	cd.Tokens = j.ToTokens()
	return cd, nil
}

type ClassBody struct {
	Methods       []MethodDefinition
	StaticMethods []MethodDefinition
	Tokens        Tokens
}

func (j *jsParser) parseClassBody(yield, await bool) (ClassBody, error) {
	var cb ClassBody
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ";"}) {
			continue
		} else if j.Peek().Type == TokenRightBracePunctuator {
			break
		}
		g := j.NewGoal()
		static := g.AcceptToken(parser.Token{TokenIdentifier, "static"})
		g.AcceptRunWhitespace()
		md, err := g.parseMethodDefinition(yield, await)
		if err != nil {
			return cb, j.Error(err)
		}
		j.Score(g)
		if static {
			cb.StaticMethods = append(cb.StaticMethods, md)
		} else {
			cb.Methods = append(cb.Methods, md)
		}
	}
	cb.Tokens = j.ToTokens()
	return cb, nil
}

type MethodType uint8

const (
	MethodNormal MethodType = iota
	MethodGenerator
	MethodAsync
	MethodGetter
	MethodSetter
)

type MethodDefinition struct {
	PropertyName PropertyName
	Params       FormalParameters
	FunctionBody StatementList
	Type         MethodType
	Tokens       Tokens
}

func (j *jsParser) parseMethodDefinition(yield, await bool) (MethodDefinition, error) {
	var md MethodDefinition
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		md.Type = MethodGenerator
		j.AcceptRunWhitespace()
	} else if j.AcceptToken(parser.Token{TokenIdentifier, "async"}) {
		md.Type = MethodAsync
		j.AcceptRunWhitespaceNoNewLine()
	} else if j.AcceptToken(parser.Token{TokenIdentifier, "get"}) {
		md.Type = MethodGetter
		j.AcceptRunWhitespace()
	} else if j.AcceptToken(parser.Token{TokenIdentifier, "set"}) {
		md.Type = MethodSetter
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	var err error
	if md.PropertyName, err = g.parsePropertyName(yield, await); err != nil {
		return md, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return md, j.Error(ErrMissingOpeningParentheses)
	}
	j.AcceptRunWhitespace()
	if md.Type != MethodGetter {
		g = j.NewGoal()
		if md.Params, err = g.parseFormalParameters(md.Type == MethodGenerator, md.Type == MethodAsync); err != nil {
			return md, j.Error(err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return md, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		return md, j.Error(ErrMissingOpeningBrace)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if md.FunctionBody, err = g.parseStatementList(md.Type == MethodGenerator, md.Type == MethodAsync, true); err != nil {
		return md, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.Accept(TokenRightBracePunctuator) {
		return md, j.Error(ErrMissingClosingBrace)
	}
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
			return pn, j.Error(err)
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
