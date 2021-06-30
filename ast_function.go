package javascript

import (
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
// https://262.ecma-international.org/11.0/#prod-FunctionDeclaration
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
	if j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "async"}) {
		fd.Type = FunctionAsync
		j.AcceptRunWhitespaceNoNewLine()
	}
	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "function"}) {
		return j.Error("FunctionDeclaration", ErrInvalidFunction)
	}
	j.AcceptRunWhitespace()
	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
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
// https://262.ecma-international.org/11.0/#prod-FormalParameters
type FormalParameters struct {
	FormalParameterList  []BindingElement
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Tokens               Tokens
}

func (fp *FormalParameters) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("FormalParameters", ErrMissingOpeningParenthesis)
	}
	j.AcceptRunWhitespace()
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		for {
			g := j.NewGoal()
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				g.AcceptRunWhitespace()
				h := g.NewGoal()
				if t := h.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
					fp.ArrayBindingPattern = new(ArrayBindingPattern)
					if err := fp.ArrayBindingPattern.parse(&h, yield, await); err != nil {
						return g.Error("FormalParameters", err)
					}
				} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
					fp.ObjectBindingPattern = new(ObjectBindingPattern)
					if err := fp.ObjectBindingPattern.parse(&h, yield, await); err != nil {
						return g.Error("FormalParameters", err)
					}
				} else if fp.BindingIdentifier = h.parseIdentifier(yield, await); fp.BindingIdentifier == nil {
					return g.Error("FormalParameters", ErrNoIdentifier)
				}
				g.Score(h)
				j.Score(g)
				j.AcceptRunWhitespace()
				if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
					return j.Error("FormalParameters", ErrMissingClosingParenthesis)
				}
				break
			}
			h := g.NewGoal()
			be := len(fp.FormalParameterList)
			fp.FormalParameterList = append(fp.FormalParameterList, BindingElement{})
			if err := fp.FormalParameterList[be].parse(&h, nil, yield, await); err != nil {
				return g.Error("FormalParameters", err)
			}
			g.Score(h)
			j.Score(g)
			j.AcceptRunWhitespace()
			if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				break
			} else if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return j.Error("FormalParameters", ErrMissingComma)
			}
			j.AcceptRunWhitespace()
		}
	}
	fp.Tokens = j.ToTokens()
	return nil
}

// BindingElement as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-BindingElement
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

func (be *BindingElement) parse(j *jsParser, singleNameBinding *Token, yield, await bool) error {
	g := j.NewGoal()
	if singleNameBinding != nil {
		be.SingleNameBinding = singleNameBinding
	} else if t := g.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
		be.ArrayBindingPattern = new(ArrayBindingPattern)
		if err := be.ArrayBindingPattern.parse(&g, yield, await); err != nil {
			return j.Error("BindingElement", err)
		}
	} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
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
	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
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
