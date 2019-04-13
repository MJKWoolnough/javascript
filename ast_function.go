package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type FunctionType uint8

const (
	FunctionNormal FunctionType = iota
	FunctionGenerator
	FunctionAsync
)

type FunctionDeclaration struct {
	Type              FunctionType
	BindingIdentifier *BindingIdentifier
	FormalParameters  FormalParameters
	FunctionBody      Block
	Tokens            Tokens
}

func (j *jsParser) parseFunctionDeclaration(yield, await, def bool) (FunctionDeclaration, error) {
	var fd FunctionDeclaration
	if j.AcceptToken(parser.Token{TokenKeyword, "async"}) {
		fd.Type = FunctionAsync
		j.AcceptRunWhitespaceNoNewLine()
	}
	if !j.AcceptToken(parser.Token{TokenKeyword, "function"}) {
		return fd, j.Error(ErrInvalidFunction)
	}
	j.AcceptRunWhitespace()
	if fd.Type == 0 && j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		fd.Type = FunctionGenerator
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	bi, err := g.parseBindingIdentifier(yield, await)
	if err != nil {
		if !def {
			return fd, j.Error(err)
		}
	} else {
		j.Score(g)
		fd.BindingIdentifier = &bi
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return fd, j.Error(ErrMissingOpeningParentheses)
	}
	g = j.NewGoal()
	fd.FormalParameters, err = g.parseFormalParameters(fd.Type == FunctionGenerator, fd.Type == FunctionAsync && await)
	if err != nil {
		return fd, j.Error(err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return fd, j.Error(ErrMissingClosingParentheses)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	fd.FunctionBody, err = j.parseBlock(fd.Type == FunctionGenerator, fd.Type == FunctionAsync, true)
	if err != nil {
		return fd, j.Error(err)
	}
	j.Score(g)
	fd.Tokens = j.ToTokens()
	return fd, nil
}

type FormalParameters struct {
	FormalParameterList   []BindingElement
	FunctionRestParameter *FunctionRestParameter
	Tokens                Tokens
}

func (j *jsParser) parseFormalParameters(yield, await bool) (FormalParameters, error) {
	var fp FormalParameters
	for {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		}
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			fr, err := g.parseFunctionRestParameter(yield, await)
			if err != nil {
				return fp, j.Error(err)
			}
			j.Score(g)
			fp.FunctionRestParameter = &fr
			break
		}
		be, err := g.parseBindingElement(yield, await)
		if err != nil {
			return fp, err
		}
		j.Score(g)
		fp.FormalParameterList = append(fp.FormalParameterList, be)
		j.AcceptRunWhitespace()
		if j.Peek() == (parser.Token{TokenPunctuator, ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return fp, j.Error(ErrInvalidFormalParameterList)
		}
	}
	fp.Tokens = j.ToTokens()
	return fp, nil
}

type BindingElement struct {
	SingleNameBinding    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}

func (j *jsParser) parseBindingElement(yield, await bool) (BindingElement, error) {
	var be BindingElement
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		ab, err := g.parseArrayBindingPattern(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.ArrayBindingPattern = &ab
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ob, err := g.parseObjectBindingPattern(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.ObjectBindingPattern = &ob
	} else {
		bi, err := g.parseBindingIdentifier(yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		be.SingleNameBinding = &bi
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		ae, err := g.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return be, j.Error(err)
		}
		j.Score(g)
		be.Initializer = &ae
	}
	be.Tokens = j.ToTokens()
	return be, nil
}

type FunctionRestParameter struct {
	BindingIdentifier    *BindingIdentifier
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (j *jsParser) parseFunctionRestParameter(yield, await bool) (FunctionRestParameter, error) {
	var fr FunctionRestParameter
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		ab, err := g.parseArrayBindingPattern(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.ArrayBindingPattern = &ab
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		ob, err := g.parseObjectBindingPattern(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.ObjectBindingPattern = &ob
	} else {
		bi, err := g.parseBindingIdentifier(yield, await)
		if err != nil {
			return fr, j.Error(err)
		}
		fr.BindingIdentifier = &bi
	}
	j.Score(g)
	fr.Tokens = j.ToTokens()
	return fr, nil
}

const (
	ErrInvalidFunction errors.Error = "invalid function"
)
