package javascript

import (
	"testing"
)

func TestParseFunction(t *testing.T) {
	doTests(t, []sourceFn{
		{`function nameHere(){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					Tokens: t[4:4],
				},
				FunctionBody: Block{
					Tokens: t[5:7],
				},
				Tokens: t[:7],
			}
		}},
		{`async function nameHere(){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				Type:              FunctionAsync,
				BindingIdentifier: &t[4],
				FormalParameters: FormalParameters{
					Tokens: t[6:6],
				},
				FunctionBody: Block{
					Tokens: t[7:9],
				},
				Tokens: t[:9],
			}
		}},
		{`function *nameHere(){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				Type:              FunctionGenerator,
				BindingIdentifier: &t[3],
				FormalParameters: FormalParameters{
					Tokens: t[5:5],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function (){}`, func(ft *test, t Tokens) {
			ft.Err = Error{
				Err: Error{
					Err:     ErrMissingIdentifier,
					Parsing: "Identifier",
					Token:   t[2],
				},
				Parsing: "FunctionDeclaration",
				Token:   t[2],
			}
		}},
		{`async function (){}`, func(ft *test, t Tokens) {
			ft.Err = Error{
				Err: Error{
					Err:     ErrMissingIdentifier,
					Parsing: "Identifier",
					Token:   t[4],
				},
				Parsing: "FunctionDeclaration",
				Token:   t[4],
			}
		}},
		{`function *(){}`, func(ft *test, t Tokens) {
			ft.Err = Error{
				Err: Error{
					Err:     ErrMissingIdentifier,
					Parsing: "Identifier",
					Token:   t[3],
				},
				Parsing: "FunctionDeclaration",
				Token:   t[3],
			}
		}},
		{`function (){}`, func(ft *test, t Tokens) {
			ft.Def = true
			ft.Output = FunctionDeclaration{
				FormalParameters: FormalParameters{
					Tokens: t[3:3],
				},
				FunctionBody: Block{
					Tokens: t[4:6],
				},
				Tokens: t[:6],
			}
		}},
		{`async function (){}`, func(ft *test, t Tokens) {
			ft.Def = true
			ft.Output = FunctionDeclaration{
				Type: FunctionAsync,
				FormalParameters: FormalParameters{
					Tokens: t[5:5],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function *(){}`, func(ft *test, t Tokens) {
			ft.Def = true
			ft.Output = FunctionDeclaration{
				Type: FunctionGenerator,
				FormalParameters: FormalParameters{
					Tokens: t[4:4],
				},
				FunctionBody: Block{
					Tokens: t[5:7],
				},
				Tokens: t[:7],
			}
		}},
		{`function myFunc(a){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &t[4],
							Tokens:            t[4:5],
						},
					},
					Tokens: t[4:5],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function myFunc(aye, bee){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &t[4],
							Tokens:            t[4:5],
						},
						{
							SingleNameBinding: &t[7],
							Tokens:            t[7:8],
						},
					},
					Tokens: t[4:8],
				},
				FunctionBody: Block{
					Tokens: t[9:11],
				},
				Tokens: t[:11],
			}
		}},
		{`function myFunc(aye, be, sea, ...dee){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &t[4],
							Tokens:            t[4:5],
						},
						{
							SingleNameBinding: &t[7],
							Tokens:            t[7:8],
						},
						{
							SingleNameBinding: &t[10],
							Tokens:            t[10:11],
						},
					},
					FunctionRestParameter: &FunctionRestParameter{
						BindingIdentifier: &t[14],
						Tokens:            t[14:15],
					},
					Tokens: t[4:15],
				},
				FunctionBody: Block{
					Tokens: t[16:18],
				},
				Tokens: t[:18],
			}
		}},
		{`function myFunc(...aye){}`, func(ft *test, t Tokens) {
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FunctionRestParameter: &FunctionRestParameter{
						BindingIdentifier: &t[5],
						Tokens:            t[5:6],
					},
					Tokens: t[4:6],
				},
				FunctionBody: Block{
					Tokens: t[7:9],
				},
				Tokens: t[:9],
			}
		}},
	}, func(test *test) (interface{}, error) {
		return test.Tokens.parseFunctionDeclaration(test.Yield, test.Await, test.Def)
	})
}
