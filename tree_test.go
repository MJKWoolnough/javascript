package javascript

import (
	"testing"

	"vimagination.zapto.org/memio"
	"vimagination.zapto.org/parser"
)

func TestTree(t *testing.T) {
	var buf memio.Buffer
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
		{
			"const a = `Fifteen is ${a + b} and not ${2 * a + b}.`;",
			Tokens{
				Keyword("const"),
				Whitespace(" "),
				Identifier("a"),
				Whitespace(" "),
				Punctuator("="),
				Whitespace(" "),
				Template{
					TemplateStart("`Fifteen is ${"),
					Tokens{
						Tokens{
							Identifier("a"),
							Whitespace(" "),
							Punctuator("+"),
							Whitespace(" "),
							Identifier("b"),
						},
						TemplateMiddle("} and not ${"),
						Tokens{
							Number("2"),
							Whitespace(" "),
							Punctuator("*"),
							Whitespace(" "),
							Identifier("a"),
							Whitespace(" "),
							Punctuator("+"),
							Whitespace(" "),
							Identifier("b"),
						},
					},
					TemplateEnd("}.`"),
				},
				Punctuator(";"),
			},
		},
	} {
		out, err := Tree(parser.NewStringTokeniser(test.Input))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if !matchTokens(test.Output, out) {
			t.Errorf("test %d: bad match, expecting %v, got %v", n+1, test.Output, out)
		}
		buf = buf[:0]
		out.WriteTo(&buf)
		if string(buf) != test.Input {
			t.Errorf("test %d: output mismatch, expecting %s, got %s", n+1, test.Input, buf)
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
			if !ok || ta.TemplateStart != tb.TemplateStart || ta.TemplateEnd != tb.TemplateEnd || !matchTokens(ta.TemplateMiddle, tb.TemplateMiddle) {
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
