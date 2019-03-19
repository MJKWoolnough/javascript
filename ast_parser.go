package javascript

import (
	"fmt"
	"strings"

	"vimagination.zapto.org/parser"
)

type TokenPos struct {
	parser.Token
	Pos, Line, LinePos uint64
}

type jsParser []TokenPos

func newJSParser(t parser.Tokeniser) (jsParser, error) {
	t.TokeniserState(new(jsTokeniser).inputElement)
	var (
		tokens             jsParser
		pos, line, linePos uint64
	)
	for {
		tk, _ := t.GetToken()
		tokens = append(tokens, TokenPos{
			Token:   tk,
			Pos:     pos,
			Line:    line,
			LinePos: linePos,
		})
		switch tk, _ := t.GetToken(); tk.Type {
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

func (j *jsParser) next() TokenPos {
	l := len(*j)
	*j = (*j)[:l+1]
	tk := (*j)[l]
	if tk.Type == parser.TokenDone {
		*j = (*j)[:l]
	}
	return tk
}

func (j *jsParser) backup() {
	*j = (*j)[:len(*j)-1]
}

func (j *jsParser) Peek() TokenPos {
	tk := j.next()
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

type Error struct {
	Err error
	TokenPos
}

func (e Error) Error() string {
	return fmt.Sprintf("error at position %d (%d:%d):\n%s", e.Pos, e.Line, e.LinePos, e.Err)
}

func (j *jsParser) Error(err error) error {
	return Error{
		Err:      err,
		TokenPos: j.Peek(),
	}
}
