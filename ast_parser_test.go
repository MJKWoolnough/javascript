package javascript

import (
	"fmt"
	"reflect"
	"testing"

	"vimagination.zapto.org/parser"
)

func TestNewJSParser(t *testing.T) {
	for n, test := range [...]struct {
		Source   string
		JSParser jsParser
		Err      error
	}{
		{"", jsParser{Token{Token: parser.Token{Type: parser.TokenDone}}}, nil},
		{
			"\"use strict\";\n\nvar hello = `World\n!`;",
			jsParser{
				{
					parser.Token{Type: TokenStringLiteral, Data: "\"use strict\""},
					0, 0, 0,
				},
				{
					parser.Token{Type: TokenPunctuator, Data: ";"},
					12, 0, 12,
				},
				{
					parser.Token{Type: TokenLineTerminator, Data: "\n\n"},
					13, 0, 13,
				},
				{
					parser.Token{Type: TokenKeyword, Data: "var"},
					15, 2, 0,
				},
				{
					parser.Token{Type: TokenWhitespace, Data: " "},
					18, 2, 3,
				},
				{
					parser.Token{Type: TokenIdentifier, Data: "hello"},
					19, 2, 4,
				},
				{
					parser.Token{Type: TokenWhitespace, Data: " "},
					24, 2, 9,
				},
				{
					parser.Token{Type: TokenPunctuator, Data: "="},
					25, 2, 10,
				},
				{
					parser.Token{Type: TokenWhitespace, Data: " "},
					26, 2, 11,
				},
				{
					parser.Token{Type: TokenNoSubstitutionTemplate, Data: "`World\n!`"},
					27, 2, 12,
				},
				{
					parser.Token{Type: TokenPunctuator, Data: ";"},
					36, 3, 2,
				},
				{
					parser.Token{Type: parser.TokenDone, Data: ""},
					37, 3, 3,
				},
			},
			nil,
		},
		{
			"@",
			jsParser{
				{
					parser.Token{Type: parser.TokenError, Data: "invalid character: @"},
					0, 0, 0,
				},
			},
			Error{
				Err:     fmt.Errorf("%w: %s", ErrInvalidCharacter, "@"),
				Parsing: "Tokens",
				Token: Token{
					parser.Token{Type: parser.TokenError, Data: "invalid character: @"},
					0, 0, 0,
				},
			},
		},
	} {
		j, err := newJSParser(makeTokeniser(parser.NewStringTokeniser(test.Source)))
		if !reflect.DeepEqual(err, test.Err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.Err, err)
		}
		for m, tk := range j[:cap(j)] {
			tkp := j.Peek()
			tkn := j.next()
			tkl := j.GetLastToken()
			if tkn != tk {
				t.Errorf("test %d.%d.1: expecting %v, got %v", n+1, m+1, tk, tkn)
			} else if tkp != tkn.Token {
				t.Errorf("test %d.%d.2: expecting to Peek %v, got %v", n+1, m+1, tkn.Token, tkp)
			} else if *tkl != tkn {
				t.Errorf("test %d.%d.3: expectign to GetLast %v, got %v", n+1, m+1, tkn, *tkl)
			}
		}
		if test.Err == nil {
			if tk := j.next(); tk.Type != parser.TokenDone {
				t.Errorf("test %d: expecting TokenDone, got %v", cap(j)+1, tk)
			}
		}
	}
	pErr := Error{
		Err:     fmt.Errorf("%w: %s", ErrInvalidCharacter, "@"),
		Parsing: "Tokens",
		Token: Token{
			Token: parser.Token{
				Type: parser.TokenError,
				Data: "invalid character: @",
			},
			Pos:     0,
			Line:    0,
			LinePos: 0,
		},
	}
	if _, err := ParseScript(makeTokeniser(parser.NewStringTokeniser("@"))); !reflect.DeepEqual(err, pErr) {
		t.Errorf("Script token error test: expecting %s, got %s", pErr, err)
	}
	if _, err := ParseModule(makeTokeniser(parser.NewStringTokeniser("@"))); !reflect.DeepEqual(err, pErr) {
		t.Errorf("Module token error test: expecting %s, got %s", pErr, err)
	}
	tk := Token{
		Token: parser.Token{
			Type: TokenPunctuator,
			Data: "?",
		},
		Pos:     0,
		Line:    0,
		LinePos: 0,
	}
	sErr := Error{
		Err: Error{
			Err: Error{
				Err:     assignmentError(tk),
				Parsing: "Expression",
				Token:   tk,
			},
			Parsing: "Statement",
			Token:   tk,
		},
		Parsing: "StatementListItem",
		Token:   tk,
	}
	if _, err := ParseScript(makeTokeniser(parser.NewStringTokeniser("?"))); !reflect.DeepEqual(err, sErr) {
		t.Errorf("Script error test: expecting %s, got %s", sErr, err)
	}
	mErr := Error{
		Err:     sErr,
		Parsing: "ModuleItem",
		Token:   tk,
	}
	if _, err := ParseModule(makeTokeniser(parser.NewStringTokeniser("?"))); !reflect.DeepEqual(err, mErr) {
		t.Errorf("Module error test: expecting %s, got %s", mErr, err)
	}
	fErr := Error{
		Err:     errorStr("TEST"),
		Parsing: "FAUX",
		Token: Token{
			Pos:     1,
			Line:    2,
			LinePos: 3,
		},
	}
	const e = "FAUX: error at position 2 (3:4):\nTEST"
	if str := fErr.Error(); str != e {
		t.Errorf("error test: expecting %q, got %q", e, str)
	}
}

type errorStr string

func (e errorStr) Error() string {
	return string(e)
}
