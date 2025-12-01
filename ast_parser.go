package javascript

import (
	"fmt"
	"strings"

	"vimagination.zapto.org/parser"
)

// Token represents a single parsed token with source positioning.
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}

func (t Token) IsTypescript() bool {
	return t.Type&tokenTypescript != 0
}

// Tokens is a collection of Token values.
type Tokens []Token

// Comments is a collection of Comment Tokens.
type Comments []*Token

type jsParser Tokens

// Tokeniser is an interface representing a tokeniser.
type Tokeniser interface {
	TokeniserState(parser.TokenFunc)
	Iter(func(parser.Token) bool)
	GetError() error
}

func newJSParser(t Tokeniser) (jsParser, error) {
	t.TokeniserState(new(jsTokeniser).inputElement)

	var (
		tokens             jsParser
		pos, line, linePos uint64
		err                error
	)

	for tk := range t.Iter {
		typ := tk.Type

		if typ >= tokenTypescript {
			typ = typ &^ tokenTypescript
		}

		tokens = append(tokens, Token{
			Token: parser.Token{
				Type: typ,
				Data: tk.Data,
			},
			Pos:     pos,
			Line:    line,
			LinePos: linePos,
		})

		switch typ {
		case parser.TokenError:
			err = Error{
				Err:     t.GetError(),
				Parsing: "Tokens",
				Token:   tokens[len(tokens)-1],
			}
		case TokenLineTerminator:
			var lastChar rune

			for _, c := range tk.Data {
				if lastChar != '\r' || c != '\n' {
					line++
				}

				lastChar = c
			}

			linePos = 0
		case TokenNoSubstitutionTemplate, TokenTemplateHead, TokenTemplateMiddle, TokenTemplateTail, TokenMultiLineComment:
			var (
				lastLT   int
				lastChar rune
			)

			for n, c := range tk.Data {
				if strings.ContainsRune(lineTerminators, c) {
					lastLT = n + 1
					linePos = 0

					if lastChar != '\r' || c != '\n' {
						line++
					}
				}

				lastChar = c
			}

			linePos += uint64(len(tk.Data) - lastLT)
		default:
			linePos += uint64(len(tk.Data))
		}

		pos += uint64(len(tk.Data))
	}

	return tokens[0:0:len(tokens)], err
}

func (j jsParser) NewGoal() jsParser {
	return j[len(j):]
}

func (j *jsParser) Score(k jsParser) {
	*j = (*j)[:len(*j)+len(k)]
}

func (j *jsParser) next() *Token {
	l := len(*j)
	if l == cap(*j) {
		return &(*j)[l-1]
	}

	*j = (*j)[:l+1]
	tk := (*j)[l]

	return &tk
}

func (j *jsParser) backup() {
	*j = (*j)[:len(*j)-1]
}

func (j *jsParser) Peek() parser.Token {
	tk := j.next().Token

	j.backup()

	return tk
}

func (j *jsParser) Accept(ts ...parser.TokenType) bool {
	tt := j.next().Type

	for _, pt := range ts {
		if pt == tt {
			return true
		}
	}

	j.backup()

	return false
}

func (j *jsParser) AcceptRun(ts ...parser.TokenType) parser.TokenType {
Loop:
	for {
		tt := j.next().Type

		for _, pt := range ts {
			if pt == tt {
				continue Loop
			}
		}

		j.backup()

		return tt
	}
}

func (j *jsParser) Skip() {
	j.next()
}

func (j *jsParser) Next() *Token {
	return j.next()
}

var depths = [...][2]parser.Token{
	{{Type: TokenPunctuator, Data: "["}, {Type: TokenPunctuator, Data: "]"}},
	{{Type: TokenPunctuator, Data: "("}, {Type: TokenPunctuator, Data: ")"}},
	{{Type: TokenPunctuator, Data: "{"}, {Type: TokenRightBracePunctuator, Data: "}"}},
}

func (j *jsParser) SkipDepth() bool {
	var (
		on    = -1
		depth = 1
	)

	for n, d := range depths {
		if j.AcceptToken(d[0]) {
			on = n

			break
		}
	}

	if on == -1 {
		return false
	}

	for depth > 0 {
		if j.AcceptToken(depths[on][0]) {
			depth++
		} else if j.AcceptToken(depths[on][1]) {
			depth--
		} else {
			j.Skip()
		}
	}

	return true
}

func (j *jsParser) AcceptToken(tk parser.Token) bool {
	if j.next().Token == tk {
		return true
	}

	j.backup()

	return false
}

func (j *jsParser) ToTokens() Tokens {
	return Tokens((*j)[:len(*j):len(*j)])
}

func (j jsParser) ToTypescriptComments() Comments {
	if len(j) == 0 {
		return nil
	}

	c := make(Comments, len(j))

	for n := range j {
		c[n] = &(j)[n]
		c[n].Type |= tokenTypescript
	}

	return c
}

func (j *jsParser) AcceptRunWhitespace() parser.TokenType {
	return j.AcceptRun(TokenWhitespace, TokenLineTerminator, TokenSingleLineComment, TokenMultiLineComment)
}

func (j *jsParser) AcceptRunWhitespaceNoNewLine() parser.TokenType {
	var tt parser.TokenType

	for {
		if tt = j.AcceptRun(TokenWhitespace); tt != TokenMultiLineComment {
			return tt
		}

		if strings.ContainsAny(j.Peek().Data, lineTerminators) {
			return tt
		}

		j.Skip()
	}
}

func (j *jsParser) AcceptRunWhitespaceNoComment() parser.TokenType {
	return j.AcceptRun(TokenWhitespace, TokenLineTerminator)
}

func (j *jsParser) AcceptRunWhitespaceComments() Comments {
	var c Comments

	g := j.NewGoal()

Loop:
	for {
		switch g.AcceptRunWhitespaceNoComment() {
		case TokenSingleLineComment, TokenMultiLineComment:
		default:
			break Loop
		}

		c = append(c, g.Next())

		j.Score(g)

		g = j.NewGoal()
	}

	return c
}

func (j *jsParser) AcceptRunWhitespaceNoNewLineNoComment() parser.TokenType {
	return j.AcceptRun(TokenWhitespace)
}

func (j *jsParser) AcceptRunWhitespaceNoNewlineComments() Comments {
	var c Comments

	g := j.NewGoal()

Loop:
	for {
		switch g.AcceptRunWhitespaceNoNewLineNoComment() {
		case TokenSingleLineComment, TokenMultiLineComment:
		default:
			break Loop
		}

		c = append(c, g.Next())

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespaceNoNewLineNoComment()

		if g.Accept(TokenLineTerminator) {
			if l := g.GetLastToken().Data; l != "\n" && l != "\r\n" {
				break
			}
		}
	}

	return c
}

func (j *jsParser) AcceptRunWhitespaceCommentsInList() Comments {
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.Accept(TokenPunctuator, TokenKeyword) {
		switch g.GetLastToken().Data {
		case ",", ".", "+", "-", "++", "--", "*", "**", "/", "%", "|", "||", "&", "&&", "^", "=", "==", "!=", "===", "!==", "in", "instanceof", "<<", "<", ">", "<=", "??", "?.", "?", ":", "(", "[", "else":
			return j.AcceptRunWhitespaceComments()
		}
	} else if g.Accept(TokenTemplateMiddle, TokenTemplateTail) {
		return j.AcceptRunWhitespaceComments()
	}

	return j.AcceptRunWhitespaceNoNewlineComments()
}

func (j *jsParser) GetLastToken() *Token {
	return &(*j)[len(*j)-1]
}

// Error is a parsing error with trace details.
type Error struct {
	Err     error
	Parsing string
	Token   Token
}

// Error returns the error string.
func (e Error) Error() string {
	return fmt.Sprintf("%s: error at position %d (%d:%d):\n%s", e.Parsing, e.Token.Pos+1, e.Token.Line+1, e.Token.LinePos+1, e.Err)
}

// Unwrap returns the wrapped error.
func (e Error) Unwrap() error {
	return e.Err
}

func (j *jsParser) Error(parsingFunc string, err error) error {
	tk := j.next()

	j.backup()

	return Error{
		Err:     err,
		Parsing: parsingFunc,
		Token:   *tk,
	}
}
