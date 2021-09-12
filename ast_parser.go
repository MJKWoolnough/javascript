package javascript

import (
	"fmt"
	"strings"

	"vimagination.zapto.org/parser"
)

// Token represents a single parsed token with source positioning
type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}

// Tokens is a collection of Token values
type Tokens []Token

type jsParser Tokens

func newJSParser(t parser.Tokeniser) (jsParser, error) {
	t.TokeniserState(new(jsTokeniser).inputElement)
	var (
		tokens             jsParser
		pos, line, linePos uint64
	)
	for {
		tk, _ := t.GetToken()
		tokens = append(tokens, Token{
			Token:   tk,
			Pos:     pos,
			Line:    line,
			LinePos: linePos,
		})
		switch tk.Type {
		case parser.TokenDone:
			return tokens[0:0:len(tokens)], nil
		case parser.TokenError:
			return nil, Error{
				Err:     t.Err,
				Parsing: "Tokens",
				Token:   tokens[len(tokens)-1],
			}
		default:
			switch tk.Type {
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
	}
}

func (j jsParser) NewGoal() jsParser {
	return j[len(j):]
}

func (j *jsParser) Score(k jsParser) {
	*j = (*j)[:len(*j)+len(k)]
}

func (j *jsParser) next() Token {
	l := len(*j)
	if l == cap(*j) {
		return (*j)[l-1]
	}
	*j = (*j)[:l+1]
	tk := (*j)[l]
	return tk
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

func (j *jsParser) ExceptRun(ts ...parser.TokenType) parser.TokenType {
	for {
		tt := j.next().Type
		for _, pt := range ts {
			if pt == tt || tt < 0 {
				j.backup()
				return tt
			}
		}
	}
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

func (j *jsParser) AcceptRunWhitespace() parser.TokenType {
	return j.AcceptRun(TokenWhitespace, TokenLineTerminator, TokenSingleLineComment, TokenMultiLineComment)
}

func (j *jsParser) AcceptRunWhitespaceNoNewLine() parser.TokenType {
	var tt parser.TokenType
	for {
		tt = j.AcceptRun(TokenWhitespace)
		if tt != TokenMultiLineComment {
			return tt
		}
		if strings.ContainsAny(j.Peek().Data, lineTerminators) {
			return tt
		}
		j.Skip()
	}
}

func (j *jsParser) GetLastToken() *Token {
	return &(*j)[len(*j)-1]
}

// Error is a parsing error with trace details
type Error struct {
	Err     error
	Parsing string
	Token   Token
}

// Error returns the error string
func (e Error) Error() string {
	return fmt.Sprintf("%s: error at position %d (%d:%d):\n%s", e.Parsing, e.Token.Pos+1, e.Token.Line+1, e.Token.LinePos+1, e.Err)
}

// Unwrap returns the wrapped error
func (e Error) Unwrap() error {
	return e.Err
}

func (j *jsParser) Error(parsingFunc string, err error) error {
	tk := j.next()
	j.backup()
	return Error{
		Err:     err,
		Parsing: parsingFunc,
		Token:   tk,
	}
}
