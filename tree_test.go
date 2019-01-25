package javascript

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTree(t *testing.T) {
	for n, test := range [...]struct {
		Input  string
		Output Tokens
	}{
		{"", Tokens{}},
		{"\"a\"", Tokens{String("\"a\"")}},
		{
			"\"use strict\";\nfunction a(b, c, d) {\n	alert(1);\n}",
			Tokens{
				String("\"use strict\""),
				Punctuator(";"),
				LineTerminators("\n"),
				Keyword("function"),
				Whitespace(" "),
				Identifier("a"),
				Punctuator("("),
				Tokens{
					Identifier("b"),
					Punctuator(","),
					Whitespace(" "),
					Identifier("c"),
					Punctuator(","),
					Whitespace(" "),
					Identifier("d"),
				},
				Punctuator(")"),
				Whitespace(" "),
				Punctuator("{"),
				Tokens{
					LineTerminators("\n"),
					Whitespace("	"),
					Identifier("alert"),
					Punctuator("("),
					Tokens{
						Number("1"),
					},
					Punctuator(")"),
					Punctuator(";"),
					LineTerminators("\n"),
				},
				Punctuator("}"),
			},
		},
	} {
		out, err := Tree(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if !matchTokens(test.Output, out) {
			t.Errorf("test %d: bad match, expecting %v, got %v", n+1, test.Output, out)
		}
	}
}

func matchTokens(a, b Tokens) bool {
	if len(a) != len(b) {
		return false
	}
	for n, ta := range a {
		switch ta := ta.(type) {
		case Tokens:
			tb, ok := b[n].(Tokens)
			if !ok || !matchTokens(ta, tb) {
				return false
			}
		case Template:
			tb, ok := b[n].(Template)
			if !ok || !matchTokens(Tokens(ta), Tokens(tb)) {
				return false
			}
		default:
			if ta != b[n] {
				return false
			}
		}
	}
	return true
}
