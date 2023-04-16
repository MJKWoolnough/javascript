package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func cloneToken(tk javascript.Token) javascript.Token {
	return javascript.Token{
		Token: parser.Token{
			Type: tk.Type,
			Data: tk.Data,
		},
		Pos:     tk.Pos,
		Line:    tk.Line,
		LinePos: tk.LinePos,
	}
}

func TestMinify(t *testing.T) {
	for n, test := range [...]struct {
		Options []Option
		Input   string
		Output  func(tk javascript.Tokens) *javascript.Module
	}{
		{
			[]Option{Literals()},
			"let a = false, b = true, c = undefined;",
			func(tk javascript.Tokens) *javascript.Module {
				f := cloneToken(tk[6])
				f.Data = "!1"
				t := cloneToken(tk[13])
				t.Data = "!0"
				u := cloneToken(tk[20])
				u.Data = "void 0"
				return &javascript.Module{
					ModuleListItems: []javascript.ModuleItem{
						{
							StatementListItem: &javascript.StatementListItem{
								Declaration: &javascript.Declaration{
									LexicalDeclaration: &javascript.LexicalDeclaration{
										BindingList: []javascript.LexicalBinding{
											{
												BindingIdentifier: &tk[2],
												Initializer: &javascript.AssignmentExpression{
													ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
														Literal: &f,
														Tokens:  []javascript.Token{f},
													}),
													Tokens: []javascript.Token{f},
												},
												Tokens: tk[2:7],
											},
											{
												BindingIdentifier: &tk[9],
												Initializer: &javascript.AssignmentExpression{
													ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
														Literal: &t,
														Tokens:  []javascript.Token{t},
													}),
													Tokens: []javascript.Token{t},
												},
												Tokens: tk[9:14],
											},
											{
												BindingIdentifier: &tk[16],
												Initializer: &javascript.AssignmentExpression{
													ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
														IdentifierReference: &u,
														Tokens:              []javascript.Token{u},
													}),
													Tokens: []javascript.Token{u},
												},
												Tokens: tk[16:21],
											},
										},
										Tokens: tk[:22],
									},
									Tokens: tk[:22],
								},
								Tokens: tk[:22],
							},
							Tokens: tk[:22],
						},
					},
					Tokens: tk[:22],
				}
			},
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		m := New(test.Options...)
		out, err := javascript.ParseModule(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error parsing test input: %s", n+1, err)
		} else {
			m.Process(out)
			if expected := test.Output(out.Tokens); !reflect.DeepEqual(out, expected) {
				t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, out, expected)
			}
		}
	}
}
