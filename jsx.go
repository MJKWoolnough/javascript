package javascript

import "vimagination.zapto.org/parser"

type jsx struct {
	Tokeniser
}

func (j *jsx) hasFlags() (bool, bool) {
	t, _ := tokeniserFlags(j.Tokeniser)

	return t, true
}

// AsJSX converts the tokeniser to one that handles the JSX extentions to
// JavaScript.
//
// Can be combined with AsTSX to read TSX.
func AsJSX(t Tokeniser) Tokeniser {
	ts, _ := tokeniserFlags(t)
	jsx := &jsx{Tokeniser: t}

	jsx.TokeniserState((&jsTokeniser{isTypescript: ts, isJSX: true}).inputElement)

	return jsx
}

// JSXElement as defined in:
// https://facebook.github.io/jsx/#prod-JSXElement
type JSXElement struct {
	ElementName JSXElementName
	Attributes  []JSXAttribute
	SelfClosing bool
	Children    []JSXChild
	Comments    [4]Comments
	Tokens      Tokens
}

func (je *JSXElement) parse(j *jsParser, yield, await bool) error {
	j.Skip()

	je.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()

	g := j.NewGoal()

	if err := je.ElementName.parse(&g); err != nil {
		return j.Error("JSXElement", err)
	}

	j.Score(g)

	for {
		j.AcceptRunWhitespace()

		if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
			je.Comments[3] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			je.SelfClosing = true

			if !j.AcceptToken(parser.Token{Type: TokenJSXElementEnd, Data: ">"}) {
				return j.Error("JSXElement", ErrMissingTagClose)
			}

			break
		} else if j.Accept(TokenJSXElementEnd) {
			j.AcceptRunWhitespace()

			break
		}

		g = j.NewGoal()

		var a JSXAttribute

		if err := a.parse(&g, yield, await); err != nil {
			return j.Error("JSXElement", err)
		}

		j.Score(g)

		je.Attributes = append(je.Attributes, a)
	}

	if !je.SelfClosing {
		for {
			g = j.NewGoal()

			if g.Accept(TokenJSXElementStart) {
				je.Comments[1] = g.AcceptRunWhitespaceComments()

				g.AcceptRunWhitespace()

				if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
					je.Comments[2] = g.AcceptRunWhitespaceComments()

					g.AcceptRunWhitespace()
					j.Score(g)

					break
				}
			}

			g = j.NewGoal()

			var child JSXChild

			if err := child.parse(&g, yield, await); err != nil {
				return j.Error("JSXElement", err)
			}

			je.Children = append(je.Children, child)

			j.Score(g)
			j.AcceptRunWhitespace()
		}

		g = j.NewGoal()

		var closing JSXElementName

		if err := closing.parse(&g); err != nil {
			return j.Error("JSXElement", err)
		}

		j.Score(g)

		if !je.ElementName.equal(&closing) {
			return j.Error("JSXElement", ErrInvalidClosingTag)
		}

		je.Comments[3] = closing.Comments[2]

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenJSXElementEnd, Data: ">"}) {
			return j.Error("JSXElement", ErrMissingTagClose)
		}
	}

	je.Tokens = j.ToTokens()

	return nil
}

// JSXElementName as defined in:
// https://facebook.github.io/jsx/#prod-JSXElementName
//
// Identifier must be defined, and only one of Namespace and MemberExpression
// can be non-nil.
type JSXElementName struct {
	Namespace        *Token
	Identifier       *Token
	MemberExpression []CommentsToken
	Comments         [3]Comments
	Tokens           Tokens
}

func (jn *JSXElementName) parse(j *jsParser) error {
	if !j.Accept(TokenJSXIdentifier) {
		return j.Error("JSXElementName", ErrMissingIdentifier)
	}

	jn.Identifier = j.GetLastToken()

	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
		jn.Comments[0] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()
		j.Skip()

		jn.Comments[1] = j.AcceptRunWhitespaceComments()
		j.AcceptRunWhitespace()

		if !j.Accept(TokenJSXIdentifier) {
			return j.Error("JSXElementName", ErrMissingIdentifier)
		}

		jn.Namespace = jn.Identifier
		jn.Identifier = j.GetLastToken()
	} else {
		for g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "."}) {
			var ct CommentsToken

			ct.Comments[0] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()

			ct.Comments[1] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			if !j.Accept(TokenJSXIdentifier) {
				return j.Error("JSXElementName", ErrMissingIdentifier)
			}

			ct.Token = j.GetLastToken()

			jn.MemberExpression = append(jn.MemberExpression, ct)

			g = j.NewGoal()

			g.AcceptRunWhitespace()
		}
	}

	jn.Comments[2] = j.AcceptRunWhitespaceComments()
	jn.Tokens = j.ToTokens()

	return nil
}

func (jn *JSXElementName) equal(cn *JSXElementName) bool {
	if !jn.Identifier.equal(cn.Identifier) || !jn.Namespace.equal(cn.Namespace) || len(jn.MemberExpression) != len(cn.MemberExpression) {
		return false
	}

	for n := range jn.MemberExpression {
		if !jn.MemberExpression[n].equal(cn.MemberExpression[n].Token) {
			return false
		}
	}

	return true
}

func (t *Token) equal(s *Token) bool {
	if t == nil {
		return s == nil
	} else if s == nil {
		return false
	}

	return t.Data == s.Data
}

type CommentsToken struct {
	Comments [2]Comments
	*Token
}

// JSXAttribute as defined in:
// https://facebook.github.io/jsx/#prod-JSXAttributes
//
// Namespace can only be non-nil if Identifier is non-nil.
//
// One, and only one of JSXString, JSXFragment, JSXElement, and
// AssignmentExpression must be non-nil.
type JSXAttribute struct {
	Namespace            *Token
	Identifier           *Token
	JSXString            *Token
	JSXFragment          *JSXFragment
	JSXElement           *JSXElement
	AssignmentExpression *AssignmentExpression
	Comments             [5]Comments
	Tokens               Tokens
}

func (ja *JSXAttribute) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		ja.Comments[0] = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}) {
			return j.Error("JSXAttribute", ErrMissingSpread)
		}

		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()
		ja.AssignmentExpression = new(AssignmentExpression)

		if err := ja.AssignmentExpression.parse(&g, false, yield, await); err != nil {
			return j.Error("JSXAttribute", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
			return j.Error("JSXAttribute", ErrMissingClosingBrace)
		}
	} else {
		if !j.Accept(TokenJSXIdentifier) {
			return j.Error("JSXAttribute", ErrMissingIdentifier)
		}

		ja.Identifier = j.GetLastToken()

		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
			ja.Comments[0] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()

			ja.Comments[1] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			if !j.Accept(TokenJSXIdentifier) {
				return j.Error("JSXAttribute", ErrMissingIdentifier)
			}

			ja.Namespace = ja.Identifier
			ja.Identifier = j.GetLastToken()

			g = j.NewGoal()

			g.AcceptRunWhitespace()
		}

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
			ja.Comments[2] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()

			ja.Comments[3] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			if j.Accept(TokenJSXString) {
				ja.JSXString = j.GetLastToken()
			} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
				j.AcceptRunWhitespaceNoComment()

				g := j.NewGoal()
				ja.AssignmentExpression = new(AssignmentExpression)

				if err := ja.AssignmentExpression.parse(&g, false, yield, await); err != nil {
					return j.Error("JSXAttribute", err)
				}

				j.Score(g)
				j.AcceptRunWhitespace()

				if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
					return j.Error("JSXAttribute", ErrMissingClosingBrace)
				}
			} else if tk := j.Peek(); tk.Type == TokenJSXElementStart {
				g := j.NewGoal()

				g.Skip()

				if g.AcceptRunWhitespace() == TokenJSXElementEnd {
					g = j.NewGoal()
					ja.JSXFragment = new(JSXFragment)

					if err := ja.JSXFragment.parse(&g, yield, await); err != nil {
						return j.Error("JSXAttribute", err)
					}

					j.Score(g)
				} else {
					g = j.NewGoal()
					ja.JSXElement = new(JSXElement)

					if err := ja.JSXElement.parse(&g, yield, await); err != nil {
						return j.Error("JSXAttribute", err)
					}

					j.Score(g)
				}
			} else {
				return j.Error("JSXAttribute", ErrMissingAttribute)
			}
		} else if tk := j.Peek(); tk.Type != TokenWhitespace && tk.Type != TokenLineTerminator && tk.Type != TokenJSXElementEnd && tk != (parser.Token{Type: TokenPunctuator, Data: "/"}) {
			return j.Error("JSXAttribute", ErrMissingEquals)
		}
	}

	ja.Comments[4] = j.AcceptRunWhitespaceComments()
	ja.Tokens = j.ToTokens()

	return nil
}

// JSXChild as defined in:
// https://facebook.github.io/jsx/#prod-JSXChild
//
// One, and only one of JSXText, JSXElement, JSXFragment, and JSXChildExpression
// must be non-nil.
//
// Spread can only be true if JSXChildExpression is non-nil.
type JSXChild struct {
	JSXText            *Token
	JSXElement         *JSXElement
	JSXFragment        *JSXFragment
	Spread             bool
	JSXChildExpression *AssignmentExpression
	Comments           Comments
	Tokens             Tokens
}

func (jc *JSXChild) parse(j *jsParser, yield, await bool) error {
	if j.Accept(TokenJSXText) {
		jc.JSXText = j.GetLastToken()
	} else if j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "{"}) {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
			jc.Comments = j.AcceptRunWhitespaceComments()
		} else {
			if jc.Spread = g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "..."}); jc.Spread {
				jc.Comments = j.AcceptRunWhitespaceComments()

				j.AcceptRunWhitespace()
				j.Skip()
			}

			j.AcceptRunWhitespaceNoComment()

			g = j.NewGoal()
			jc.JSXChildExpression = new(AssignmentExpression)

			if err := jc.JSXChildExpression.parse(&g, false, yield, await); err != nil {
				return j.Error("JSXChild", err)
			}

			j.Score(g)
		}

		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "}"}) {
			return j.Error("JSXChild", ErrMissingClosingBrace)
		}
	} else {
		g := j.NewGoal()

		g.Skip()

		if g.AcceptRunWhitespace() == TokenJSXElementEnd {
			g = j.NewGoal()
			jc.JSXFragment = new(JSXFragment)

			if err := jc.JSXFragment.parse(&g, yield, await); err != nil {
				return j.Error("JSXChild", err)
			}

			j.Score(g)
		} else {
			g = j.NewGoal()
			jc.JSXElement = new(JSXElement)

			if err := jc.JSXElement.parse(&g, yield, await); err != nil {
				return j.Error("JSXChild", err)
			}

			j.Score(g)
		}
	}

	jc.Tokens = j.ToTokens()

	return nil
}

// JSXFragment as defined in:
// https://facebook.github.io/jsx/#prod-JSXFragment
type JSXFragment struct {
	Children []JSXChild
	Comments [3]Comments
	Tokens   Tokens
}

func (jf *JSXFragment) parse(j *jsParser, yield, await bool) error {
	j.Skip()

	jf.Comments[0] = j.AcceptRunWhitespaceComments()

	j.AcceptRunWhitespace()
	j.Skip()

	for {
		j.AcceptRunWhitespace()

		g := j.NewGoal()

		if g.Accept(TokenJSXElementStart) {
			jf.Comments[1] = g.AcceptRunWhitespaceComments()

			g.AcceptRunWhitespace()

			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "/"}) {
				jf.Comments[2] = g.AcceptRunWhitespaceComments()

				g.AcceptRunWhitespace()
				j.Score(g)

				break
			}
		}

		g = j.NewGoal()

		var child JSXChild

		if err := child.parse(&g, yield, await); err != nil {
			return j.Error("JSXFragment", err)
		}

		jf.Children = append(jf.Children, child)

		j.Score(g)
		j.AcceptRunWhitespace()
	}

	if !j.Accept(TokenJSXElementEnd) {
		return j.Error("JSXFragment", ErrMissingTagClose)
	}

	jf.Tokens = j.ToTokens()

	return nil
}
