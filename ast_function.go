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
	BindingIdentifier *Token
	FormalParameters  FormalParameters
	FunctionBody      Block
	Tokens            Tokens
}

func (fd *FunctionDeclaration) parse(j *jsParser, yield, await, def bool) error {
	if j.AcceptToken(parser.Token{TokenIdentifier, "async"}) {
		fd.Type = FunctionAsync
		j.AcceptRunWhitespaceNoNewLine()
	}
	if !j.AcceptToken(parser.Token{TokenKeyword, "function"}) {
		return j.Error("FunctionDeclaration", ErrInvalidFunction)
	}
	j.AcceptRunWhitespace()
	if fd.Type == 0 && j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		fd.Type = FunctionGenerator
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	bi, err := g.parseIdentifier(yield, await)
	if err != nil {
		if !def {
			return j.Error("FunctionDeclaration", err)
		}
	} else {
		j.Score(g)
		fd.BindingIdentifier = bi
		j.AcceptRunWhitespace()
	}
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("FunctionDeclaration", ErrMissingOpeningParenthesis)
	}
	g = j.NewGoal()
	if err := fd.FormalParameters.parse(&g, fd.Type == FunctionGenerator, fd.Type == FunctionAsync && await); err != nil {
		return j.Error("FunctionDeclaration", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
		return j.Error("FunctionDeclaration", ErrMissingClosingParenthesis)
	}
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := fd.FunctionBody.parse(&g, fd.Type == FunctionGenerator, fd.Type == FunctionAsync, true); err != nil {
		return j.Error("FunctionDeclaration", err)
	}
	j.Score(g)
	fd.Tokens = j.ToTokens()
	return nil
}

type FormalParameters struct {
	FormalParameterList   []BindingElement
	FunctionRestParameter *FunctionRestParameter
	Tokens                Tokens
}

func (fp *FormalParameters) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		g.AcceptRunWhitespace()
		if g.Peek() == (parser.Token{TokenPunctuator, ")"}) {
			break
		}
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			h := g.NewGoal()
			fp.FunctionRestParameter = newFunctionRestParameter()
			if err := fp.FunctionRestParameter.parse(&h, yield, await); err != nil {
				return j.Error("FormalParameters", err)
			}
			g.Score(h)
			j.Score(g)
			break
		}
		h := g.NewGoal()
		var be BindingElement
		if err := be.parse(&h, yield, await); err != nil {
			return g.Error("FormalParameters", err)
		}
		g.Score(h)
		j.Score(g)
		fp.FormalParameterList = append(fp.FormalParameterList, be)
		j.AcceptRunWhitespace()
		if j.Peek() == (parser.Token{TokenPunctuator, ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("FormalParameters", ErrInvalidFormalParameterList)
		}
	}
	fp.Tokens = j.ToTokens()
	return nil
}

type BindingElement struct {
	SingleNameBinding    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}

func (be *BindingElement) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		be.ArrayBindingPattern = newArrayBindingPattern()
		if err := be.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		be.ObjectBindingPattern = newObjectBindingPattern()
		if err := be.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
	} else {
		bi, err := g.parseIdentifier(yield, await)
		if err != nil {
			return j.Error("BindingElement", err)
		}
		be.SingleNameBinding = bi
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		be.Initializer = newAssignmentExpression()
		if err := be.Initializer.parse(&g, true, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
		j.Score(g)
	}
	be.Tokens = j.ToTokens()
	return nil
}

type FunctionRestParameter struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (fr *FunctionRestParameter) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if g.AcceptToken(parser.Token{TokenPunctuator, "["}) {
		fr.ArrayBindingPattern = newArrayBindingPattern()
		if err := fr.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("FunctionRestParameter", err)
		}
	} else if g.AcceptToken(parser.Token{TokenPunctuator, "{"}) {
		fr.ObjectBindingPattern = newObjectBindingPattern()
		if err := fr.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("FunctionRestParameter", err)
		}
	} else {
		bi, err := g.parseIdentifier(yield, await)
		if err != nil {
			return j.Error("FunctionRestParameter", err)
		}
		fr.BindingIdentifier = bi
	}
	j.Score(g)
	fr.Tokens = j.ToTokens()
	return nil
}

var (
	ErrInvalidFunction = errors.New("invalid function")
)
