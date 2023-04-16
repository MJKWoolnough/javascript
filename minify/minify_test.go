package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func makeToken(typ parser.TokenType, data string) *javascript.Token {
	return &javascript.Token{
		Token: parser.Token{
			Type: typ,
			Data: data,
		},
	}
}

func TestTransforms(t *testing.T) {
	for n, test := range [...]struct {
		Options       []Option
		Input, Output javascript.Type
	}{
		{
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "false"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!1"),
			},
		},
		{
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "true"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!0"),
			},
		},
		{
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "undefined"),
			},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "void 0"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e3"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123450000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "12345e4"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e-3"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.00123400"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1234e-6"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xff"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "255"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e12"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xe8d4a51001"),
			},
		},
		{
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0o7"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "7"),
			},
		},
	} {
		w := walker{New(test.Options...)}
		w.Handle(test.Input)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, test.Output, test.Input)
		}
	}
}
