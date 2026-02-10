package javascript

import (
	"vimagination.zapto.org/parser"
)

// FunctionType determines which type of function is specified by FunctionDeclaration
type FunctionType uint8

// Valid FunctionTypes
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
	Comments          [5]Comments
	Tokens            Tokens
}

func (fd *FunctionDeclaration) parse(j *jsParser, yield, await, def, export bool) error {
	if j.AcceptToken(parser.Token{Type: TokenIdentifier, Data: "async"}) {
		fd.Type = FunctionAsync
		fd.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

		j.AcceptRunWhitespaceNoNewLine()
	}

	if !j.AcceptToken(parser.Token{Type: TokenKeyword, Data: "function"}) {
		return j.Error("FunctionDeclaration", ErrInvalidFunction)
	}

	fd.Comments[1] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "*"}) {
		if fd.Type == FunctionAsync {
			fd.Type = FunctionAsyncGenerator
		} else {
			fd.Type = FunctionGenerator
		}

		fd.Comments[2] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
	}

	if bi := j.parseIdentifier(yield, await); bi == nil && !def {
		return j.Error("FunctionDeclaration", ErrNoIdentifier)
	} else {
		fd.BindingIdentifier = bi
		g := j.NewGoal()

		async := fd.Type == FunctionAsync || fd.Type == FunctionAsyncGenerator

		if g.SkipFunctionOverload(bi, yield, await, def, export, async) {
			h := g.NewGoal()

			h.AcceptRunWhitespace()

			for g.SkipFunctionOverload(bi, yield, await, def, export, async) {
				g.Score(h)

				h = g.NewGoal()

				h.AcceptRunWhitespace()
			}

			fd.Comments[3] = g.ToTypescriptComments()

			j.Score(g)
		}

		fd.Comments[3] = append(fd.Comments[3], j.AcceptRunWhitespaceComments()...)

		j.AcceptRunWhitespace()
	}

	g := j.NewGoal()

	if g.SkipGeneric() {
		fd.Comments[3] = append(fd.Comments[3], g.ToTypescriptComments()...)
		fd.Comments[3] = append(fd.Comments[3], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	g = j.NewGoal()

	if err := fd.FormalParameters.parse(&g, fd.Type == FunctionGenerator, fd.Type == FunctionAsync && await); err != nil {
		return j.Error("FunctionDeclaration", err)
	}

	j.Score(g)

	fd.Comments[4] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	g = j.NewGoal()

	if g.SkipReturnType() {
		fd.Comments[4] = append(fd.Comments[4], g.ToTypescriptComments()...)
		fd.Comments[4] = append(fd.Comments[4], g.AcceptRunWhitespaceComments()...)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

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
//
// Only one of BindingIdentifier, ArrayBindingPattern, or ObjectBindingPattern
// can be non-nil.
type FormalParameters struct {
	FormalParameterList  []BindingElement
	BindingIdentifier    *Token
	ArrayBindingPattern  *ArrayBindingPattern
	ObjectBindingPattern *ObjectBindingPattern
	Comments             [5]Comments
	Tokens               Tokens
}

func (fp *FormalParameters) parse(j *jsParser, yield, await bool) error {
	if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "("}) {
		return j.Error("FormalParameters", ErrMissingOpeningParenthesis)
	}

	fp.Comments[0] = j.AcceptRunWhitespaceNoNewlineComments()

	j.AcceptRunWhitespaceNoComment()

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	h := g.NewGoal()

	if h.SkipThisParam() {
		i := h.NewGoal()

		i.AcceptRunWhitespace()

		if i.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
			h.Score(i)
		}

		fp.Comments[0] = append(fp.Comments[0], h.ToTypescriptComments()...)
		fp.Comments[0] = append(fp.Comments[0], h.AcceptRunWhitespaceNoNewlineComments()...)

		g.Score(h)
		j.Score(g)
		j.AcceptRunWhitespaceNoComment()
	}

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
		for {
			g := j.NewGoal()

			g.AcceptRunWhitespace()

			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
				fp.Comments[1] = j.AcceptRunWhitespaceComments()

				j.AcceptRunWhitespace()
				j.Skip()

				fp.Comments[2] = j.AcceptRunWhitespaceComments()

				j.AcceptRunWhitespace()

				g = j.NewGoal()

				if t := g.Peek(); t == (parser.Token{Type: TokenPunctuator, Data: "["}) {
					fp.ArrayBindingPattern = new(ArrayBindingPattern)
					if err := fp.ArrayBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("FormalParameters", err)
					}
				} else if t == (parser.Token{Type: TokenPunctuator, Data: "{"}) {
					fp.ObjectBindingPattern = new(ObjectBindingPattern)
					if err := fp.ObjectBindingPattern.parse(&g, yield, await); err != nil {
						return j.Error("FormalParameters", err)
					}
				} else if fp.BindingIdentifier = g.parseIdentifier(yield, await); fp.BindingIdentifier == nil {
					return j.Error("FormalParameters", ErrNoIdentifier)
				}

				j.Score(g)

				fp.Comments[3] = j.AcceptRunWhitespaceNoNewlineComments()
				g = j.NewGoal()

				g.AcceptRunWhitespace()

				h := g.NewGoal()

				if h.SkipColonType() {
					fp.Comments[3] = append(fp.Comments[3], h.ToTypescriptComments()...)
					fp.Comments[3] = append(fp.Comments[3], h.AcceptRunWhitespaceNoNewlineComments()...)

					g.Score(h)
					j.Score(g)

					g = j.NewGoal()

					g.AcceptRunWhitespace()
				}

				if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
					return g.Error("FormalParameters", ErrMissingClosingParenthesis)
				}

				break
			}

			g = j.NewGoal()
			be := len(fp.FormalParameterList)

			fp.FormalParameterList = append(fp.FormalParameterList, BindingElement{})
			if err := fp.FormalParameterList[be].parse(&g, nil, yield, await); err != nil {
				return j.Error("FormalParameters", err)
			}

			j.Score(g)

			g = j.NewGoal()

			g.AcceptRunWhitespace()

			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ")"}) {
				break
			} else if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ","}) {
				return g.Error("FormalParameters", ErrMissingComma)
			}

			j.Score(g)

			j.AcceptRunWhitespaceNoComment()
		}
	}

	fp.Comments[4] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

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
	Comments             [2]Comments
	Tokens               Tokens
}

func (be *BindingElement) parse(j *jsParser, singleNameBinding *Token, yield, await bool) error {
	if singleNameBinding == nil {
		be.Comments[0] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
	}

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

	be.Comments[1] = j.AcceptRunWhitespaceCommentsInList()

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	h := g.NewGoal()

	if h.SkipOptionalColonType() {
		be.Comments[1] = append(be.Comments[1], h.ToTypescriptComments()...)

		g.Score(h)
		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespace()
	}

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
		j.Score(g)
		j.AcceptRunWhitespaceNoComment()

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
