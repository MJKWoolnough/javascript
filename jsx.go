package javascript

import (
	"strings"

	"vimagination.zapto.org/parser"
)

const jsxMarker = "X"

type jsx struct {
	Tokeniser
}

func (j *jsx) Iter(fn func(parser.Token) bool) {
	for tk := range j.Tokeniser.Iter {
		if tk.Type == parser.TokenDone {
			tk.Data += jsxMarker
		}

		if !fn(tk) {
			break
		}
	}
}

func (j *jsx) hasFlags() (bool, bool) {
	t, _ := tokeniserFlags(j.Tokeniser)

	return t, true
}

func AsJSX(t Tokeniser) Tokeniser {
	ts, _ := tokeniserFlags(t)
	jsx := &jsx{Tokeniser: t}

	jsx.TokeniserState((&jsTokeniser{isTypescript: ts, isJSX: true}).inputElement)

	return jsx
}

func (j *jsParser) IsJSX() bool {
	return strings.HasSuffix((*j)[:cap(*j)][cap(*j)-1].Data, jsxMarker)
}

type JSXElement struct {
	ElementName JSXElementName
	Attributes  []JSXAttribute
	SelfClosing bool
	Children    *JSXChildren
	Tokens      Tokens
}

func (je *JSXElement) parse(j *jsParser) error {
	j.Skip()

	g := j.NewGoal()

	if err := je.ElementName.parse(&g); err != nil {
		return j.Error("JSXElement", err)
	}

	j.Score(g)

	for {
		if j.AcceptRunWhitespace() == TokenDivPunctuator {
			je.SelfClosing = true

			break
		} else if j.Peek().Type == TokenJSXElementEnd {
			break
		}

		g = j.NewGoal()

		var a JSXAttribute

		if err := a.parse(&g); err != nil {
			return j.Error("JSXElement", err)
		}

		j.Score(g)

		je.Attributes = append(je.Attributes, a)
	}

	j.AcceptRunWhitespace()

	if !j.Accept(TokenJSXElementEnd) {
		return j.Error("JSXElement", ErrMissingTagClose)
	}

	if !je.SelfClosing {
		g = j.NewGoal()
		je.Children = new(JSXChildren)

		if err := je.Children.parse(&g); err != nil {
			return j.Error("JSXElement", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()
		j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<"})
		j.AcceptRunWhitespace()
		j.Accept(TokenDivPunctuator)
		j.AcceptRunWhitespace()

		g = j.NewGoal()

		var closing JSXElementName

		if err := closing.parse(&g); err != nil {
			return j.Error("JSXElement", err)
		}

		j.Score(g)

		if !je.ElementName.equal(&closing) {
			return j.Error("JSXElement", ErrInvalidClosingTag)
		}

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
			return j.Error("JSXElement", ErrMissingTagClose)
		}
	}

	je.Tokens = j.ToTokens()

	return nil
}

type JSXElementName struct {
	Namespace        *Token
	MemberExpression []*Token
	Identifier       *Token
	Tokens           Tokens
}

func (jn *JSXElementName) parse(j *jsParser) error {
	if !j.Accept(TokenJSXIdentifier) {
		return j.Error("JSXElementName", ErrMissingIdentifier)
	}

	jn.Identifier = j.GetLastToken()

	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		if !j.Accept(TokenJSXIdentifier) {
			return j.Error("JSXElementName", ErrMissingIdentifier)
		}

		jn.Namespace = jn.Identifier
		jn.Identifier = j.GetLastToken()
	} else {
		for j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			if !j.Accept(TokenJSXIdentifier) {
				return j.Error("JSXElementName", ErrMissingIdentifier)
			}

			jn.MemberExpression = append(jn.MemberExpression, jn.Identifier)
			jn.Identifier = j.GetLastToken()
		}
	}

	jn.Tokens = j.ToTokens()

	return nil
}

func (jn *JSXElementName) equal(cn *JSXElementName) bool {
	if !jn.Identifier.equal(cn.Identifier) || !jn.Namespace.equal(cn.Namespace) || len(jn.MemberExpression) != len(cn.MemberExpression) {
		return false
	}

	for n := range jn.MemberExpression {
		if !jn.MemberExpression[n].equal(cn.MemberExpression[n]) {
			return false
		}
	}

	return true
}

func (t *Token) equal(s *Token) bool {
	if t == nil {
		if s == nil {
			return true
		}

		return false
	} else if s == nil {
		return false
	}

	return t.Data == s.Data
}

type JSXAttribute struct{}

func (ja *JSXAttribute) parse(j *jsParser) error {
	return nil
}

type JSXChildren struct{}

func (jc *JSXChildren) parse(j *jsParser) error {
	return nil
}

type JSXFragment struct{}

func (jf *JSXFragment) parse(j *jsParser) error {
	return nil
}
