package javascript

import (
	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

// FunctionType determines which type of function is specified by FunctionDeclaration
type FunctionType uint8

// Valid FunctionType's
const (
	FunctionNormal FunctionType = iota
	FunctionGenerator
	FunctionAsync
	FunctionAsyncGenerator
)

// FunctionDeclaration as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-FunctionDeclaration
//
// Also parses FunctionExpression, for when BindingIdentifier is nil.
//
// Include TC39 proposal for async generator functions
// https://github.com/tc39/proposal-async-iteration#async-generator-functions
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
	if j.AcceptToken(parser.Token{TokenPunctuator, "*"}) {
		if fd.Type == FunctionAsync {
			fd.Type = FunctionAsyncGenerator
		} else {
			fd.Type = FunctionGenerator
		}
		j.AcceptRunWhitespace()
	}
	if bi := j.parseIdentifier(yield, await); bi == nil {
		if !def {
			return j.Error("FunctionDeclaration", ErrNoIdentifier)
		}
	} else {
		fd.BindingIdentifier = bi
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	if err := fd.FormalParameters.parse(&g, fd.Type == FunctionGenerator, fd.Type == FunctionAsync && await); err != nil {
		return j.Error("FunctionDeclaration", err)
	}
	j.Score(g)
	j.AcceptRunWhitespace()
	g = j.NewGoal()
	if err := fd.FunctionBody.parse(&g, fd.Type == FunctionGenerator, fd.Type == FunctionAsync, true); err != nil {
		return j.Error("FunctionDeclaration", err)
	}
	j.Score(g)
	fd.Tokens = j.ToTokens()
	return nil
}

// FormalParameters as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-FormalParameters
type FormalParameters struct {
	FormalParameterList   []BindingElement
	FunctionRestParameter *FunctionRestParameter
	Tokens                Tokens
}

func (fp *FormalParameters) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{TokenPunctuator, "("}) {
		return j.Error("FormalParameters", ErrMissingOpeningParenthesis)
	}
	for {
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		}
		g := j.NewGoal()
		if g.AcceptToken(parser.Token{TokenPunctuator, "..."}) {
			g.AcceptRunWhitespace()
			h := g.NewGoal()
			fp.FunctionRestParameter = new(FunctionRestParameter)
			if err := fp.FunctionRestParameter.parse(&h, yield, await); err != nil {
				return j.Error("FormalParameters", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if !j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
				return j.Error("FormalParameters", ErrMissingClosingParenthesis)
			}
			break
		}
		h := g.NewGoal()
		be := len(fp.FormalParameterList)
		fp.FormalParameterList = append(fp.FormalParameterList, BindingElement{})
		if err := fp.FormalParameterList[be].parse(&h, yield, await); err != nil {
			return g.Error("FormalParameters", err)
		}
		g.Score(h)
		j.Score(g)
		j.AcceptRunWhitespace()
		if j.AcceptToken(parser.Token{TokenPunctuator, ")"}) {
			break
		} else if !j.AcceptToken(parser.Token{TokenPunctuator, ","}) {
			return j.Error("FormalParameters", ErrInvalidFormalParameterList)
		}
	}
	fp.Tokens = j.ToTokens()
	return nil
}

// BindingElement as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-BindingElement
//
// Only one of SingleNameBinding, ArrayBindingPattern, or ObjectBindingPattern
// must be non-nil.
//
// The Initializer is optional.
type BindingElement struct {
	SingleNameBinding    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Initializer          *AssignmentExpression
	Tokens               Tokens
}

func (be *BindingElement) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if t := g.Peek(); t == (parser.Token{TokenPunctuator, "["}) {
		be.ArrayBindingPattern = new(ArrayBindingPattern)
		if err := be.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
	} else if t == (parser.Token{TokenPunctuator, "{"}) {
		be.ObjectBindingPattern = new(ObjectBindingPattern)
		if err := be.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
	} else if be.SingleNameBinding = g.parseIdentifier(yield, await); be.SingleNameBinding == nil {
		return j.Error("BindingElement", ErrNoIdentifier)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "="}) {
		g.AcceptRunWhitespace()
		j.Score(g)
		g = j.NewGoal()
		be.Initializer = new(AssignmentExpression)
		if err := be.Initializer.parse(&g, true, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
		j.Score(g)
	}
	be.Tokens = j.ToTokens()
	return nil
}

// FunctionRestParameter as defined in ECMA-262
// https://www.ecma-international.org/ecma-262/#prod-FunctionRestParameter
//
// Only one of BindingIdentifier, ArrayBindingPattern, or ObjectBindingPattern
// must be non-nil.
type FunctionRestParameter struct {
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (fr *FunctionRestParameter) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()
	if t := g.Peek(); t == (parser.Token{TokenPunctuator, "["}) {
		fr.ArrayBindingPattern = new(ArrayBindingPattern)
		if err := fr.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("FunctionRestParameter", err)
		}
	} else if t == (parser.Token{TokenPunctuator, "{"}) {
		fr.ObjectBindingPattern = new(ObjectBindingPattern)
		if err := fr.ObjectBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("FunctionRestParameter", err)
		}
	} else if fr.BindingIdentifier = g.parseIdentifier(yield, await); fr.BindingIdentifier == nil {
		return j.Error("FunctionRestParameter", ErrNoIdentifier)
	}
	j.Score(g)
	fr.Tokens = j.ToTokens()
	return nil
}

// Errors
var (
	ErrInvalidFunction = errors.New("invalid function")
)
