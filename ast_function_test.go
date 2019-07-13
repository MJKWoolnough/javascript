package javascript

import (
	"testing"
)

func TestParseFunction(t *testing.T) {
	doTests(t, []sourceFn{
		{`function nameHere(){}`, func(ft *test, t Tokens) { // 1
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					Tokens: t[3:5],
				},
				FunctionBody: Block{
					Tokens: t[5:7],
				},
				Tokens: t[:7],
			}
		}},
		{`async function nameHere(){}`, func(ft *test, t Tokens) { // 2
			ft.Output = FunctionDeclaration{
				Type:              FunctionAsync,
				BindingIdentifier: &t[4],
				FormalParameters: FormalParameters{
					Tokens: t[5:7],
				},
				FunctionBody: Block{
					Tokens: t[7:9],
				},
				Tokens: t[:9],
			}
		}},
		{`function *nameHere(){}`, func(ft *test, t Tokens) { // 3
			ft.Output = FunctionDeclaration{
				Type:              FunctionGenerator,
				BindingIdentifier: &t[3],
				FormalParameters: FormalParameters{
					Tokens: t[4:6],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function (){}`, func(ft *test, t Tokens) { // 4
			ft.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   t[2],
			}
		}},
		{`async function (){}`, func(ft *test, t Tokens) { // 5
			ft.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   t[4],
			}
		}},
		{`function *(){}`, func(ft *test, t Tokens) { // 6
			ft.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   t[3],
			}
		}},
		{`function (){}`, func(ft *test, t Tokens) { // 7
			ft.Def = true
			ft.Output = FunctionDeclaration{
				FormalParameters: FormalParameters{
					Tokens: t[2:4],
				},
				FunctionBody: Block{
					Tokens: t[4:6],
				},
				Tokens: t[:6],
			}
		}},
		{`async function (){}`, func(ft *test, t Tokens) { // 8
			ft.Def = true
			ft.Output = FunctionDeclaration{
				Type: FunctionAsync,
				FormalParameters: FormalParameters{
					Tokens: t[4:6],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function *(){}`, func(ft *test, t Tokens) { // 9
			ft.Def = true
			ft.Output = FunctionDeclaration{
				Type: FunctionGenerator,
				FormalParameters: FormalParameters{
					Tokens: t[3:5],
				},
				FunctionBody: Block{
					Tokens: t[5:7],
				},
				Tokens: t[:7],
			}
		}},
		{`function myFunc(a){}`, func(ft *test, t Tokens) { // 10
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &t[4],
							Tokens:            t[4:5],
						},
					},
					Tokens: t[3:6],
				},
				FunctionBody: Block{
					Tokens: t[6:8],
				},
				Tokens: t[:8],
			}
		}},
		{`function myFunc(aye, bee){}`, func(ft *test, t Tokens) { // 11
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
					Tokens: t[3:9],
				},
				FunctionBody: Block{
					Tokens: t[9:11],
				},
				Tokens: t[:11],
			}
		}},
		{`function myFunc(aye, be, sea, ...dee){}`, func(ft *test, t Tokens) { // 12
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
					Tokens: t[3:16],
				},
				FunctionBody: Block{
					Tokens: t[16:18],
				},
				Tokens: t[:18],
			}
		}},
		{`function myFunc(...aye){}`, func(ft *test, t Tokens) { // 13
			ft.Output = FunctionDeclaration{
				BindingIdentifier: &t[2],
				FormalParameters: FormalParameters{
					FunctionRestParameter: &FunctionRestParameter{
						BindingIdentifier: &t[5],
						Tokens:            t[5:6],
					},
					Tokens: t[3:7],
				},
				FunctionBody: Block{
					Tokens: t[7:9],
				},
				Tokens: t[:9],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var fd FunctionDeclaration
		err := fd.parse(&t.Tokens, t.Yield, t.Await, t.Def)
		return fd, err
	})
}
