package javascript

import (
	"fmt"
	"runtime"
	"strings"

	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

type Token struct {
	parser.Token
	Pos, Line, LinePos uint64
}

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
			return nil, tokens.Error(t.Err)
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

func (j *jsParser) Except(ts ...parser.TokenType) bool {
	tt := j.next().Type
	for _, pt := range ts {
		if pt == tt {
			j.backup()
			return false
		}
	}
	return true
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
		j.Except()
	}
}

func (j *jsParser) GetLastToken() *Token {
	if len(*j) == 0 {
		return nil
	}
	return &(*j)[len(*j)-1]
}

type Error struct {
	Err     error
	Parsing string
	Token   Token
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: error at position %d (%d:%d):\n%s", e.Parsing, e.Token.Pos, e.Token.Line, e.Token.LinePos, e.Err)
}

func (e Error) getLastPos() uint64 {
	if e, ok := e.Err.(Error); ok {
		return e.getLastPos()
	}
	return e.Token.Pos
}

func (j *jsParser) Error(err error) error {
	const prefix = "vimagination.zapto.org/javascript.(*jsParser).parse"
	ptrs := make([]uintptr, 1000)
	var parsing string
	for _, ptr := range ptrs[:runtime.Callers(2, ptrs)] {
		fn := runtime.FuncForPC(ptr).Name()
		if strings.HasPrefix(fn, prefix) {
			parsing = fn[len(prefix):]
			break
		}
	}
	return Error{
		Err:     err,
		Parsing: parsing,
		Token:   j.next(),
	}
}

func (j *jsParser) FindGoal(fns ...func(*jsParser) error) error {
	var (
		err     error
		lastPos uint64
	)
	for _, fn := range fns {
		g := j.NewGoal()
		if errr := fn(&g); errr == nil {
			j.Score(g)
			return nil
		} else if p := g.next().Pos; errr != errNotApplicable && (err == nil || lastPos < p) {
			err = g.Error(errr)
			lastPos = p
		}
	}
	return j.Error(err)
}

const (
	errNotApplicable errors.Error = ""
)
