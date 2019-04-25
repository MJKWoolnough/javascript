package javascript

import "testing"

func TestClassDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{`class myClass{}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Tokens: tk[3:3],
				},
				Tokens: tk[:5],
			}
		}},
		{`class myClass extends OtherClass{}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				Extends: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &IdentifierReference{
									Identifier: &tk[6],
								},
								Tokens: tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[6:7],
				},
				ClassBody: ClassBody{
					Tokens: tk[7:7],
				},
				Tokens: tk[:9],
			}
		}},
		{`class myClass {constructor(){}}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
						{
							PropertyName: PropertyName{
								LiteralPropertyName: &tk[5],
								Tokens:              tk[5:6],
							},
							Params: FormalParameters{
								Tokens: tk[7:7],
							},
							FunctionBody: Block{
								Tokens: tk[8:10],
							},
							Tokens: tk[5:10],
						},
					},
					Tokens: tk[5:10],
				},
				Tokens: tk[:11],
			}
		}},
		{`class myClass {method(arg1, arg2){}}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
						{
							PropertyName: PropertyName{
								LiteralPropertyName: &tk[5],
								Tokens:              tk[5:6],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &BindingIdentifier{Identifier: &tk[7]},
										Tokens:            tk[7:8],
									},
									{
										SingleNameBinding: &BindingIdentifier{Identifier: &tk[10]},
										Tokens:            tk[10:11],
									},
								},
								Tokens: tk[7:11],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[5:14],
						},
					},
					Tokens: tk[5:14],
				},
				Tokens: tk[:15],
			}
		}},
		{`class myClass {set method(...args){}}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
						{
							Type: MethodSetter,
							PropertyName: PropertyName{
								LiteralPropertyName: &tk[7],
								Tokens:              tk[7:8],
							},
							Params: FormalParameters{
								FunctionRestParameter: &FunctionRestParameter{
									BindingIdentifier: &BindingIdentifier{Identifier: &tk[10]},
									Tokens:            tk[10:11],
								},
								Tokens: tk[9:11],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[5:14],
						},
					},
					Tokens: tk[5:14],
				},
				Tokens: tk[:15],
			}
		}},
		{`class myClass {get value(){}}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
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
					Tokens: tk[5:12],
				},
				Tokens: tk[:13],
			}
		}},
		{`` +
			`class myClass {
	get value(){}
	set value(v){}
	static hello(){}
}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
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
										SingleNameBinding: &BindingIdentifier{Identifier: &tk[20]},
										Tokens:            tk[20:21],
									},
								},
								Tokens: tk[20:21],
							},
							FunctionBody: Block{
								Tokens: tk[22:24],
							},
							Tokens: tk[16:24],
						},
					},
					StaticMethods: []MethodDefinition{
						{
							PropertyName: PropertyName{
								LiteralPropertyName: &tk[28],
								Tokens:              tk[28:29],
							},
							Params: FormalParameters{
								Tokens: tk[30:30],
							},
							FunctionBody: Block{
								Tokens: tk[31:33],
							},
							Tokens: tk[26:33],
						},
					},
					Tokens: tk[5:34],
				},
				Tokens: tk[:35],
			}
		}},
		{`class{}`, func(t *test, tk Tokens) {
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingIdentifier,
					Parsing: "Identifier",
					Token:   tk[1],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[1],
			}
		}},
		{`class{}`, func(t *test, tk Tokens) {
			t.Def = true
			t.Output = ClassDeclaration{
				ClassBody: ClassBody{
					Tokens: tk[2:2],
				},
				Tokens: tk[:3],
			}
		}},
		{`class beep`, func(t *test, tk Tokens) {
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ClassDeclaration",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseClassDeclaration(t.Yield, t.Await, t.Def)
	})
}
