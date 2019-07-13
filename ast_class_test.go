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

func TestClassDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidClassDeclaration,
				Parsing: "ClassDeclaration",
				Token:   tk[0],
			}
		}},
		{"class\n", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ClassDeclaration",
				Token:   tk[2],
			}
		}},
		{"class\n", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ClassDeclaration",
				Token:   tk[2],
			}
			t.Def = true
		}},
		{"class\na", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ClassDeclaration",
				Token:   tk[3],
			}
		}},
		{"class\na\nextends\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "PrimaryExpression",
								Token:   tk[6],
							},
							Parsing: "MemberExpression",
							Token:   tk[6],
						},
						Parsing: "NewExpression",
						Token:   tk[6],
					},
					Parsing: "LeftHandSideExpression",
					Token:   tk[6],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[6],
			}
		}},
		{"class\nextends\n", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "PrimaryExpression",
								Token:   tk[4],
							},
							Parsing: "MemberExpression",
							Token:   tk[4],
						},
						Parsing: "NewExpression",
						Token:   tk[4],
					},
					Parsing: "LeftHandSideExpression",
					Token:   tk[4],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[4],
			}
			t.Def = true
		}},
		{"class\na\nextends\nb\n", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ClassDeclaration",
				Token:   tk[8],
			}
		}},
		{"class\na\n{\n}", func(t *test, tk Tokens) { // 8
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				Tokens:            tk[:7],
			}
		}},
		{"class\na\nextends\nb\n{\n}", func(t *test, tk Tokens) { // 9
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
				Tokens: tk[:11],
			}
		}},
		{"class\nextends\na\n{\n}", func(t *test, tk Tokens) { // 10
			t.Output = ClassDeclaration{
				ClassHeritage: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[4],
								Tokens:              tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:9],
			}
			t.Def = true
		}},
		{"class\na\n{\n}", func(t *test, tk Tokens) { // 11
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				Tokens:            tk[:7],
			}
			t.Def = true
		}},
		{"class\na\n{\n;\n}", func(t *test, tk Tokens) { // 12
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				Tokens:            tk[:9],
			}
			t.Def = true
		}},
		{"class\n{\n}", func(t *test, tk Tokens) { // 13
			t.Output = ClassDeclaration{
				Tokens: tk[:5],
			}
			t.Def = true
		}},
		{"class\n{\n;\n}", func(t *test, tk Tokens) { // 14
			t.Output = ClassDeclaration{
				Tokens: tk[:7],
			}
			t.Def = true
		}},
		{"class\na\n{\na\n}", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningParenthesis,
						Parsing: "FormalParameters",
						Token:   tk[8],
					},
					Parsing: "MethodDefinition",
					Token:   tk[8],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[6],
			}
		}},
		{"class\na\n{\nb(){}\n}", func(t *test, tk Tokens) { // 16
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[6],
							Tokens:              tk[6:7],
						},
						Params: FormalParameters{
							Tokens: tk[7:9],
						},
						FunctionBody: Block{
							Tokens: tk[9:11],
						},
						Tokens: tk[6:11],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class\na\n{\n;\na(){}\n;\n}", func(t *test, tk Tokens) { // 17
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[8],
							Tokens:              tk[8:9],
						},
						Params: FormalParameters{
							Tokens: tk[9:11],
						},
						FunctionBody: Block{
							Tokens: tk[11:13],
						},
						Tokens: tk[8:13],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class\na\n{\na(){}\nb(){}\n}", func(t *test, tk Tokens) { // 18
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[6],
							Tokens:              tk[6:7],
						},
						Params: FormalParameters{
							Tokens: tk[7:9],
						},
						FunctionBody: Block{
							Tokens: tk[9:11],
						},
						Tokens: tk[6:11],
					},
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[12],
							Tokens:              tk[12:13],
						},
						Params: FormalParameters{
							Tokens: tk[13:15],
						},
						FunctionBody: Block{
							Tokens: tk[15:17],
						},
						Tokens: tk[12:17],
					},
				},
				Tokens: tk[:19],
			}
		}},
		{"class\na\n{\n;\na(){}\n;\nb(){}\n;\n}", func(t *test, tk Tokens) { // 19
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []MethodDefinition{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[8],
							Tokens:              tk[8:9],
						},
						Params: FormalParameters{
							Tokens: tk[9:11],
						},
						FunctionBody: Block{
							Tokens: tk[11:13],
						},
						Tokens: tk[8:13],
					},
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[16],
							Tokens:              tk[16:17],
						},
						Params: FormalParameters{
							Tokens: tk[17:19],
						},
						FunctionBody: Block{
							Tokens: tk[19:21],
						},
						Tokens: tk[16:21],
					},
				},
				Tokens: tk[:25],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var cd ClassDeclaration
		err := cd.parse(&t.Tokens, t.Yield, t.Await, t.Def)
		return cd, err
	})
}
