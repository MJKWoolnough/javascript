package javascript

import "testing"

func TestClassDeclarationOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`class myClass{}`, func(t *test, tk Tokens) { // 1
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				Tokens:            tk[:5],
			}
		}},
		{`class myClass extends OtherClass{}`, func(t *test, tk Tokens) { // 2
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassHeritage: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[6],
								Tokens:              tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[:9],
			}
		}},
		{`class myClass {constructor(){}}`, func(t *test, tk Tokens) { // 3
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[5],
							Tokens:              tk[5:6],
						},
						Params: FormalParameters{
							Tokens: tk[6:8],
						},
						FunctionBody: Block{
							Tokens: tk[8:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`class myClass {method(arg1, arg2){}}`, func(t *test, tk Tokens) { // 4
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[5],
							Tokens:              tk[5:6],
						},
						Params: FormalParameters{
							FormalParameterList: []BindingElement{
								{
									SingleNameBinding: &tk[7],
									Tokens:            tk[7:8],
								},
								{
									SingleNameBinding: &tk[10],
									Tokens:            tk[10:11],
								},
							},
							Tokens: tk[6:12],
						},
						FunctionBody: Block{
							Tokens: tk[12:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{`class myClass {set method(...args){}}`, func(t *test, tk Tokens) { // 5
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						Type: MethodSetter,
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[7],
							Tokens:              tk[7:8],
						},
						Params: FormalParameters{
							FunctionRestParameter: &FunctionRestParameter{
								BindingIdentifier: &tk[10],
								Tokens:            tk[10:11],
							},
							Tokens: tk[8:12],
						},
						FunctionBody: Block{
							Tokens: tk[12:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{`class myClass {get value(){}}`, func(t *test, tk Tokens) { // 6
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						Type: MethodGetter,
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[7],
							Tokens:              tk[7:8],
						},
						FunctionBody: Block{
							Tokens: tk[10:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{`` +
			`class myClass {
				get value(){}
				set value(v){}
				static hello(){}
			}`, func(t *test, tk Tokens) { // 7
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						Type: MethodGetter,
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[9],
							Tokens:              tk[9:10],
						},
						FunctionBody: Block{
							Tokens: tk[12:14],
						},
						Tokens: tk[7:14],
					},
					{
						Type: MethodSetter,
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[18],
							Tokens:              tk[18:19],
						},
						Params: FormalParameters{
							FormalParameterList: []BindingElement{
								{
									SingleNameBinding: &tk[20],
									Tokens:            tk[20:21],
								},
							},
							Tokens: tk[19:22],
						},
						FunctionBody: Block{
							Tokens: tk[22:24],
						},
						Tokens: tk[16:24],
					},
					{
						Type: MethodStatic,
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[28],
							Tokens:              tk[28:29],
						},
						Params: FormalParameters{
							Tokens: tk[29:31],
						},
						FunctionBody: Block{
							Tokens: tk[31:33],
						},
						Tokens: tk[26:33],
					},
				},
				Tokens: tk[:36],
			}
		}},
		{`class{}`, func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ClassDeclaration",
				Token:   tk[1],
			}
		}},
		{`class{}`, func(t *test, tk Tokens) { // 9
			t.Def = true
			t.Output = ClassDeclaration{
				Tokens: tk[:3],
			}
		}},
		{`class beep`, func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ClassDeclaration",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var cd ClassDeclaration
		err := cd.parse(&t.Tokens, t.Yield, t.Await, t.Def)
		return cd, err
	})
}
