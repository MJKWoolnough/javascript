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
					BindingIdentifier: &t[14],
					Tokens:            t[3:16],
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
					BindingIdentifier: &t[5],
					Tokens:            t[3:7],
				},
				FunctionBody: Block{
					Tokens: t[7:9],
				},
				Tokens: t[:9],
			}
		}},
	}, func(t *test) (Type, error) {
		var fd FunctionDeclaration

		err := fd.parse(&t.Tokens, t.Yield, t.Await, t.Def)

		return fd, err
	})
}

func TestFunctionDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidFunction,
				Parsing: "FunctionDeclaration",
				Token:   tk[0],
			}
		}},
		{"function", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   tk[1],
			}
		}},
		{"async function", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   tk[3],
			}
		}},
		{"async\nfunction", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrInvalidFunction,
				Parsing: "FunctionDeclaration",
				Token:   tk[1],
			}
		}},
		{"function*", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   tk[2],
			}
		}},
		{"async function*", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FunctionDeclaration",
				Token:   tk[4],
			}
		}},
		{"function", func(t *test, tk Tokens) { // 7
			t.Def = true
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "FormalParameters",
					Token:   tk[1],
				},
				Parsing: "FunctionDeclaration",
				Token:   tk[1],
			}
		}},
		{"function\na", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "FormalParameters",
					Token:   tk[3],
				},
				Parsing: "FunctionDeclaration",
				Token:   tk[3],
			}
		}},
		{"function\na\n()", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[6],
				},
				Parsing: "FunctionDeclaration",
				Token:   tk[6],
			}
		}},
		{"function\na\n()\n{}", func(t *test, tk Tokens) { // 10
			t.Output = FunctionDeclaration{
				BindingIdentifier: &tk[2],
				FormalParameters: FormalParameters{
					Tokens: tk[4:6],
				},
				FunctionBody: Block{
					Tokens: tk[7:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"async function\na\n()\n{}", func(t *test, tk Tokens) { // 11
			t.Output = FunctionDeclaration{
				Type:              FunctionAsync,
				BindingIdentifier: &tk[4],
				FormalParameters: FormalParameters{
					Tokens: tk[6:8],
				},
				FunctionBody: Block{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"function\n*\na\n()\n{}", func(t *test, tk Tokens) { // 12
			t.Output = FunctionDeclaration{
				Type:              FunctionGenerator,
				BindingIdentifier: &tk[4],
				FormalParameters: FormalParameters{
					Tokens: tk[6:8],
				},
				FunctionBody: Block{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"async function\n*\na\n()\n{}", func(t *test, tk Tokens) { // 13
			t.Output = FunctionDeclaration{
				Type:              FunctionAsyncGenerator,
				BindingIdentifier: &tk[6],
				FormalParameters: FormalParameters{
					Tokens: tk[8:10],
				},
				FunctionBody: Block{
					Tokens: tk[11:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"function\na\n()\n{}", func(t *test, tk Tokens) { // 14
			t.Def = true
			t.Output = FunctionDeclaration{
				BindingIdentifier: &tk[2],
				FormalParameters: FormalParameters{
					Tokens: tk[4:6],
				},
				FunctionBody: Block{
					Tokens: tk[7:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"async function\na\n()\n{}", func(t *test, tk Tokens) { // 15
			t.Def = true
			t.Output = FunctionDeclaration{
				Type:              FunctionAsync,
				BindingIdentifier: &tk[4],
				FormalParameters: FormalParameters{
					Tokens: tk[6:8],
				},
				FunctionBody: Block{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"function\n*\na\n()\n{}", func(t *test, tk Tokens) { // 16
			t.Def = true
			t.Output = FunctionDeclaration{
				Type:              FunctionGenerator,
				BindingIdentifier: &tk[4],
				FormalParameters: FormalParameters{
					Tokens: tk[6:8],
				},
				FunctionBody: Block{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"async function\n*\na\n()\n{}", func(t *test, tk Tokens) { // 17
			t.Def = true
			t.Output = FunctionDeclaration{
				Type:              FunctionAsyncGenerator,
				BindingIdentifier: &tk[6],
				FormalParameters: FormalParameters{
					Tokens: tk[8:10],
				},
				FunctionBody: Block{
					Tokens: tk[11:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"function\n()\n{}", func(t *test, tk Tokens) { // 18
			t.Def = true
			t.Output = FunctionDeclaration{
				FormalParameters: FormalParameters{
					Tokens: tk[2:4],
				},
				FunctionBody: Block{
					Tokens: tk[5:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"async function\n()\n{}", func(t *test, tk Tokens) { // 19
			t.Def = true
			t.Output = FunctionDeclaration{
				Type: FunctionAsync,
				FormalParameters: FormalParameters{
					Tokens: tk[4:6],
				},
				FunctionBody: Block{
					Tokens: tk[7:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"function\n*\n()\n{}", func(t *test, tk Tokens) { // 20
			t.Def = true
			t.Output = FunctionDeclaration{
				Type: FunctionGenerator,
				FormalParameters: FormalParameters{
					Tokens: tk[4:6],
				},
				FunctionBody: Block{
					Tokens: tk[7:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"async function\n*\n()\n{}", func(t *test, tk Tokens) { // 21
			t.Def = true
			t.Output = FunctionDeclaration{
				Type: FunctionAsyncGenerator,
				FormalParameters: FormalParameters{
					Tokens: tk[6:8],
				},
				FunctionBody: Block{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
	}, func(t *test) (Type, error) {
		var fd FunctionDeclaration

		err := fd.parse(&t.Tokens, t.Yield, t.Await, t.Def)

		return fd, err
	})
}

func TestFormalParameters(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "FormalParameters",
				Token:   tk[0],
			}
		}},
		{"(\n)", func(t *test, tk Tokens) { // 2
			t.Output = FormalParameters{
				Tokens: tk[:3],
			}
		}},
		{"(\n...\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "FormalParameters",
				Token:   tk[4],
			}
		}},
		{"(\n...\na\n)", func(t *test, tk Tokens) { // 4
			t.Output = FormalParameters{
				BindingIdentifier: &tk[4],
				Tokens:            tk[:7],
			}
		}},
		{"(\n...\na\nb)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "FormalParameters",
				Token:   tk[6],
			}
		}},
		{"(\n,)", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[2],
				},
				Parsing: "FormalParameters",
				Token:   tk[2],
			}
		}},
		{"(\na\n)", func(t *test, tk Tokens) { // 7
			t.Output = FormalParameters{
				FormalParameterList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"(\na\nb)", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "FormalParameters",
				Token:   tk[4],
			}
		}},
		{"(\na\n,\nb\n)", func(t *test, tk Tokens) { // 9
			t.Output = FormalParameters{
				FormalParameterList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						SingleNameBinding: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"(\na\n,\n...\nb\n)", func(t *test, tk Tokens) { // 10
			t.Output = FormalParameters{
				FormalParameterList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				BindingIdentifier: &tk[8],
				Tokens:            tk[:11],
			}
		}},
		{"(...[])", func(t *test, tk Tokens) { // 11
			t.Output = FormalParameters{
				ArrayBindingPattern: &ArrayBindingPattern{
					Tokens: tk[2:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"(...{})", func(t *test, tk Tokens) { // 12
			t.Output = FormalParameters{
				ObjectBindingPattern: &ObjectBindingPattern{
					Tokens: tk[2:4],
				},
				Tokens: tk[:5],
			}
		}},
		{`(...[!])`, func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[3],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[3],
				},
				Parsing: "FormalParameters",
				Token:   tk[2],
			}
		}},
		{`(...{!})`, func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[3],
						},
						Parsing: "BindingProperty",
						Token:   tk[3],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[3],
				},
				Parsing: "FormalParameters",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var fp FormalParameters

		err := fp.parse(&t.Tokens, t.Yield, t.Await)

		return fp, err
	})
}

func TestBindingElement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "BindingElement",
				Token:   tk[0],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 2
			t.Output = BindingElement{
				SingleNameBinding: &tk[0],
				Tokens:            tk[:1],
			}
		}},
		{"a\n=\n1", func(t *test, tk Tokens) { // 3
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = BindingElement{
				SingleNameBinding: &tk[0],
				Initializer: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"[a]\n=\n", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     assignmentError(tk[6]),
				Parsing: "BindingElement",
				Token:   tk[6],
			}
		}},
		{"[]", func(t *test, tk Tokens) { // 5
			t.Output = BindingElement{
				ArrayBindingPattern: &ArrayBindingPattern{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"[]\n=\na", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 5)
			t.Output = BindingElement{
				ArrayBindingPattern: &ArrayBindingPattern{
					Tokens: tk[:2],
				},
				Initializer: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[5:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"{a}\n=\n", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     assignmentError(tk[6]),
				Parsing: "BindingElement",
				Token:   tk[6],
			}
		}},
		{"{}", func(t *test, tk Tokens) { // 8
			t.Output = BindingElement{
				ObjectBindingPattern: &ObjectBindingPattern{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"{}\n=\na", func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 5)
			t.Output = BindingElement{
				ObjectBindingPattern: &ObjectBindingPattern{
					Tokens: tk[:2],
				},
				Initializer: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[5:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`[!]`, func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[1],
				},
				Parsing: "BindingElement",
				Token:   tk[0],
			}
		}},
		{`{!}`, func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[1],
						},
						Parsing: "BindingProperty",
						Token:   tk[1],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[1],
				},
				Parsing: "BindingElement",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var be BindingElement

		err := be.parse(&t.Tokens, nil, t.Yield, t.Await)

		return be, err
	})
}
