package javascript

import (
	"testing"
)

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
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Params: FormalParameters{
								Tokens: tk[6:8],
							},
							FunctionBody: Block{
								Tokens: tk[8:10],
							},
							Tokens: tk[5:10],
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
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
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
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{`class myClass {set method(args){}}`, func(t *test, tk Tokens) { // 5
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &tk[9],
										Tokens:            tk[9:10],
									},
								},
								Tokens: tk[8:11],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`class myClass {get value(){}}`, func(t *test, tk Tokens) { // 6
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[10:12],
							},
							Tokens: tk[5:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{`class myClass {
	get value(){}
	set value(v){}
	static hello(){}
}`, func(t *test, tk Tokens) { // 7
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[7:14],
						},
						Tokens: tk[7:14],
					},
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[18],
									Tokens:              tk[18:19],
								},
								Tokens: tk[18:19],
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
						Tokens: tk[16:24],
					},
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[28],
									Tokens:              tk[28:29],
								},
								Tokens: tk[28:29],
							},
							Params: FormalParameters{
								Tokens: tk[29:31],
							},
							FunctionBody: Block{
								Tokens: tk[31:33],
							},
							Tokens: tk[28:33],
						},
						Tokens: tk[26:33],
					},
				},
				Tokens: tk[:35],
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
	}, func(t *test) (Type, error) {
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
		{"class\na\n{\nb(){}\n}", func(t *test, tk Tokens) { // 15
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[6],
									Tokens:              tk[6:7],
								},
								Tokens: tk[6:7],
							},
							Params: FormalParameters{
								Tokens: tk[7:9],
							},
							FunctionBody: Block{
								Tokens: tk[9:11],
							},
							Tokens: tk[6:11],
						},
						Tokens: tk[6:11],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class\na\n{\n;\na(){}\n;\n}", func(t *test, tk Tokens) { // 16
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[8],
									Tokens:              tk[8:9],
								},
								Tokens: tk[8:9],
							},
							Params: FormalParameters{
								Tokens: tk[9:11],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[8:13],
						},
						Tokens: tk[8:13],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class\na\n{\na(){}\nb(){}\n}", func(t *test, tk Tokens) { // 17
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[6],
									Tokens:              tk[6:7],
								},
								Tokens: tk[6:7],
							},
							Params: FormalParameters{
								Tokens: tk[7:9],
							},
							FunctionBody: Block{
								Tokens: tk[9:11],
							},
							Tokens: tk[6:11],
						},
						Tokens: tk[6:11],
					},
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[12],
									Tokens:              tk[12:13],
								},
								Tokens: tk[12:13],
							},
							Params: FormalParameters{
								Tokens: tk[13:15],
							},
							FunctionBody: Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[12:17],
						},
						Tokens: tk[12:17],
					},
				},
				Tokens: tk[:19],
			}
		}},
		{"class\na\n{\n;\na(){}\n;\nb(){}\n;\n}", func(t *test, tk Tokens) { // 18
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[8],
									Tokens:              tk[8:9],
								},
								Tokens: tk[8:9],
							},
							Params: FormalParameters{
								Tokens: tk[9:11],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[8:13],
						},
						Tokens: tk[8:13],
					},
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[16],
									Tokens:              tk[16:17],
								},
								Tokens: tk[16:17],
							},
							Params: FormalParameters{
								Tokens: tk[17:19],
							},
							FunctionBody: Block{
								Tokens: tk[19:21],
							},
							Tokens: tk[16:21],
						},
						Tokens: tk[16:21],
					},
				},
				Tokens: tk[:25],
			}
		}},
		{"class a {static b() {} }", func(t *test, tk Tokens) { // 19
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[7:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static [\"b\"]() {} }", func(t *test, tk Tokens) { // 20
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											Literal: &tk[8],
											Tokens:  tk[8:9],
										}),
										Tokens: tk[8:9],
									},
									Tokens: tk[7:10],
								},
								Tokens: tk[7:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class a {static #b() {} }", func(t *test, tk Tokens) { // 21
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[7:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static b}", func(t *test, tk Tokens) { // 22
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static b = 1}", func(t *test, tk Tokens) { // 23
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static [b]}", func(t *test, tk Tokens) { // 24
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[8],
											Tokens:              tk[8:9],
										}),
										Tokens: tk[8:9],
									},
									Tokens: tk[7:10],
								},
								Tokens: tk[7:10],
							},
							Tokens: tk[7:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {static [b] = 1}", func(t *test, tk Tokens) { // 25
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[8],
											Tokens:              tk[8:9],
										}),
										Tokens: tk[8:9],
									},
									Tokens: tk[7:10],
								},
								Tokens: tk[7:10],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[13],
									Tokens:  tk[13:14],
								}),
								Tokens: tk[13:14],
							},
							Tokens: tk[7:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static #b}", func(t *test, tk Tokens) { // 26
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static #b = 1}", func(t *test, tk Tokens) { // 27
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static b;}", func(t *test, tk Tokens) { // 28
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:9],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{"class a {static b = 1 ;}", func(t *test, tk Tokens) { // 29
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static [b]\n;}", func(t *test, tk Tokens) { // 30
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[8],
											Tokens:              tk[8:9],
										}),
										Tokens: tk[8:9],
									},
									Tokens: tk[7:10],
								},
								Tokens: tk[7:10],
							},
							Tokens: tk[7:10],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static [b] = 1;;}", func(t *test, tk Tokens) { // 31
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[8],
											Tokens:              tk[8:9],
										}),
										Tokens: tk[8:9],
									},
									Tokens: tk[7:10],
								},
								Tokens: tk[7:10],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[13],
									Tokens:  tk[13:14],
								}),
								Tokens: tk[13:14],
							},
							Tokens: tk[7:14],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class a {static #b \n;}", func(t *test, tk Tokens) { // 32
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:11],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{"class a {static #b = 1\n ;}", func(t *test, tk Tokens) { // 33
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static () {}}", func(t *test, tk Tokens) { // 34
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Params: FormalParameters{
								Tokens: tk[7:9],
							},
							FunctionBody: Block{
								Tokens: tk[10:12],
							},
							Tokens: tk[5:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static}", func(t *test, tk Tokens) { // 35
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {static = 1}", func(t *test, tk Tokens) { // 36
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {static\n;}", func(t *test, tk Tokens) { // 37
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static = 1\n;}", func(t *test, tk Tokens) { // 38
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static\n{}}", func(t *test, tk Tokens) { // 39
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						ClassStaticBlock: &Block{
							Tokens: tk[7:9],
						},
						Tokens: tk[5:9],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{"class a {static{b;c}}", func(t *test, tk Tokens) { // 40
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						ClassStaticBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: WrapConditional(&PrimaryExpression{
														IdentifierReference: &tk[7],
														Tokens:              tk[7:8],
													}),
													Tokens: tk[7:8],
												},
											},
											Tokens: tk[7:8],
										},
										Tokens: tk[7:9],
									},
									Tokens: tk[7:9],
								},
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: WrapConditional(&PrimaryExpression{
														IdentifierReference: &tk[9],
														Tokens:              tk[9:10],
													}),
													Tokens: tk[9:10],
												},
											},
											Tokens: tk[9:10],
										},
										Tokens: tk[9:10],
									},
									Tokens: tk[9:10],
								},
							},
							Tokens: tk[6:11],
						},
						Tokens: tk[5:11],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{"class a {static async}", func(t *test, tk Tokens) { // 41
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static async;}", func(t *test, tk Tokens) { // 42
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:9],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{"class a {static async\n= 1}", func(t *test, tk Tokens) { // 43
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static async = 1;}", func(t *test, tk Tokens) { // 44
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {static\n*\nb() {}}", func(t *test, tk Tokens) { // 45
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static * async() {}}", func(t *test, tk Tokens) { // 46
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static\n* get() {}}", func(t *test, tk Tokens) { // 47
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static *\nset() {}}", func(t *test, tk Tokens) { // 48
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static\nasync *\nb() {}}", func(t *test, tk Tokens) { // 49
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[11],
									Tokens:              tk[11:12],
								},
								Tokens: tk[11:12],
							},
							Params: FormalParameters{
								Tokens: tk[12:14],
							},
							FunctionBody: Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[7:17],
						},
						Tokens: tk[5:17],
					},
				},
				Tokens: tk[:18],
			}
		}},
		{"class a {static async * async() {}}", func(t *test, tk Tokens) { // 50
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[11],
									Tokens:              tk[11:12],
								},
								Tokens: tk[11:12],
							},
							Params: FormalParameters{
								Tokens: tk[12:14],
							},
							FunctionBody: Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[7:17],
						},
						Tokens: tk[5:17],
					},
				},
				Tokens: tk[:18],
			}
		}},
		{"class a {static\nasync * get() {}}", func(t *test, tk Tokens) { // 51
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[11],
									Tokens:              tk[11:12],
								},
								Tokens: tk[11:12],
							},
							Params: FormalParameters{
								Tokens: tk[12:14],
							},
							FunctionBody: Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[7:17],
						},
						Tokens: tk[5:17],
					},
				},
				Tokens: tk[:18],
			}
		}},
		{"class a {static async #b(){}}", func(t *test, tk Tokens) { // 52
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsync,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[9],
								Tokens:            tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[7:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static async*#b(){}}", func(t *test, tk Tokens) { // 53
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[9],
								Tokens:            tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[7:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {static\nget\nb()\n{}}", func(t *test, tk Tokens) { // 54
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static\nset\nb(c)\n{}}", func(t *test, tk Tokens) { // 55
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &tk[11],
										Tokens:            tk[11:12],
									},
								},
								Tokens: tk[10:13],
							},
							FunctionBody: Block{
								Tokens: tk[14:16],
							},
							Tokens: tk[7:16],
						},
						Tokens: tk[5:16],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class a {static\nget\n#b()\n{}}", func(t *test, tk Tokens) { // 56
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[9],
								Tokens:            tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[7:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {static\nset\n#b(c)\n{}}", func(t *test, tk Tokens) { // 57
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[9],
								Tokens:            tk[9:10],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &tk[11],
										Tokens:            tk[11:12],
									},
								},
								Tokens: tk[10:13],
							},
							FunctionBody: Block{
								Tokens: tk[14:16],
							},
							Tokens: tk[7:16],
						},
						Tokens: tk[5:16],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"class a {static\nget}", func(t *test, tk Tokens) { // 58
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static\nset}", func(t *test, tk Tokens) { // 59
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[5:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {static\nget\n= 1}", func(t *test, tk Tokens) { // 60
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {static\nset\n= 1}", func(t *test, tk Tokens) { // 61
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						Static: true,
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[7:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {#b() {} }", func(t *test, tk Tokens) { // 62
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[5],
								Tokens:            tk[5:6],
							},
							Params: FormalParameters{
								Tokens: tk[6:8],
							},
							FunctionBody: Block{
								Tokens: tk[9:11],
							},
							Tokens: tk[5:11],
						},
						Tokens: tk[5:11],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {b}", func(t *test, tk Tokens) { // 63
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {b = 1}", func(t *test, tk Tokens) { // 64
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {#b}", func(t *test, tk Tokens) { // 65
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[5],
								Tokens:            tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {#b = 1}", func(t *test, tk Tokens) { // 66
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[5],
								Tokens:            tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {[b]\n;}", func(t *test, tk Tokens) { // 67
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[6],
											Tokens:              tk[6:7],
										}),
										Tokens: tk[6:7],
									},
									Tokens: tk[5:8],
								},
								Tokens: tk[5:8],
							},
							Tokens: tk[5:8],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {[b] = 1;;}", func(t *test, tk Tokens) { // 68
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[6],
											Tokens:              tk[6:7],
										}),
										Tokens: tk[6:7],
									},
									Tokens: tk[5:8],
								},
								Tokens: tk[5:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[11],
									Tokens:  tk[11:12],
								}),
								Tokens: tk[11:12],
							},
							Tokens: tk[5:12],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {#b \n;}", func(t *test, tk Tokens) { // 69
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[5],
								Tokens:            tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:9],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{"class a {#b = 1\n ;}", func(t *test, tk Tokens) { // 70
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[5],
								Tokens:            tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {async}", func(t *test, tk Tokens) { // 71
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {async\n= 1}", func(t *test, tk Tokens) { // 72
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {async = 1;}", func(t *test, tk Tokens) { // 73
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:11],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{"class a {*\nb() {}}", func(t *test, tk Tokens) { // 74
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {* async() {}}", func(t *test, tk Tokens) { // 75
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {* get() {}}", func(t *test, tk Tokens) { // 76
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {*\nset() {}}", func(t *test, tk Tokens) { // 77
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {async *\nb() {}}", func(t *test, tk Tokens) { // 78
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[5:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {async * async() {}}", func(t *test, tk Tokens) { // 79
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[5:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {async * get() {}}", func(t *test, tk Tokens) { // 80
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[9],
									Tokens:              tk[9:10],
								},
								Tokens: tk[9:10],
							},
							Params: FormalParameters{
								Tokens: tk[10:12],
							},
							FunctionBody: Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[5:15],
						},
						Tokens: tk[5:15],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"class a {async #b(){}}", func(t *test, tk Tokens) { // 81
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodAsync,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[10:12],
							},
							Tokens: tk[5:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {async*#b(){}}", func(t *test, tk Tokens) { // 82
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodAsyncGenerator,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[10:12],
							},
							Tokens: tk[5:12],
						},
						Tokens: tk[5:12],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{"class a {get\nb()\n{}}", func(t *test, tk Tokens) { // 83
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {set\nb(c)\n{}}", func(t *test, tk Tokens) { // 84
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &tk[9],
										Tokens:            tk[9:10],
									},
								},
								Tokens: tk[8:11],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[5:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {get\n#b()\n{}}", func(t *test, tk Tokens) { // 85
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodGetter,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Params: FormalParameters{
								Tokens: tk[8:10],
							},
							FunctionBody: Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[5:13],
						},
						Tokens: tk[5:13],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{"class a {set\n#b(c)\n{}}", func(t *test, tk Tokens) { // 86
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						MethodDefinition: &MethodDefinition{
							Type: MethodSetter,
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[7],
								Tokens:            tk[7:8],
							},
							Params: FormalParameters{
								FormalParameterList: []BindingElement{
									{
										SingleNameBinding: &tk[9],
										Tokens:            tk[9:10],
									},
								},
								Tokens: tk[8:11],
							},
							FunctionBody: Block{
								Tokens: tk[12:14],
							},
							Tokens: tk[5:14],
						},
						Tokens: tk[5:14],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{"class a {get}", func(t *test, tk Tokens) { // 87
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {set}", func(t *test, tk Tokens) { // 88
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"class a {get\n= 1}", func(t *test, tk Tokens) { // 89
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {set\n= 1}", func(t *test, tk Tokens) { // 90
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"class a {set = 1\n#set = 2;a\n=\n3\nb = 4;;static async c(){}static\nasync d(){}#e(){}static{}#f=5;g=6\nh=7;i\n[j]}", func(t *test, tk Tokens) { // 91
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[9],
									Tokens:  tk[9:10],
								}),
								Tokens: tk[9:10],
							},
							Tokens: tk[5:10],
						},
						Tokens: tk[5:10],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[11],
								Tokens:            tk[11:12],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[15],
									Tokens:  tk[15:16],
								}),
								Tokens: tk[15:16],
							},
							Tokens: tk[11:16],
						},
						Tokens: tk[11:17],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[17],
									Tokens:              tk[17:18],
								},
								Tokens: tk[17:18],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[21],
									Tokens:  tk[21:22],
								}),
								Tokens: tk[21:22],
							},
							Tokens: tk[17:22],
						},
						Tokens: tk[17:22],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[23],
									Tokens:              tk[23:24],
								},
								Tokens: tk[23:24],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[27],
									Tokens:  tk[27:28],
								}),
								Tokens: tk[27:28],
							},
							Tokens: tk[23:28],
						},
						Tokens: tk[23:29],
					},
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsync,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[34],
									Tokens:              tk[34:35],
								},
								Tokens: tk[34:35],
							},
							Params: FormalParameters{
								Tokens: tk[35:37],
							},
							FunctionBody: Block{
								Tokens: tk[37:39],
							},
							Tokens: tk[32:39],
						},
						Tokens: tk[30:39],
					},
					{
						Static: true,
						MethodDefinition: &MethodDefinition{
							Type: MethodAsync,
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[43],
									Tokens:              tk[43:44],
								},
								Tokens: tk[43:44],
							},
							Params: FormalParameters{
								Tokens: tk[44:46],
							},
							FunctionBody: Block{
								Tokens: tk[46:48],
							},
							Tokens: tk[41:48],
						},
						Tokens: tk[39:48],
					},
					{
						MethodDefinition: &MethodDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[48],
								Tokens:            tk[48:49],
							},
							Params: FormalParameters{
								Tokens: tk[49:51],
							},
							FunctionBody: Block{
								Tokens: tk[51:53],
							},
							Tokens: tk[48:53],
						},
						Tokens: tk[48:53],
					},
					{
						Static: true,
						ClassStaticBlock: &Block{
							Tokens: tk[54:56],
						},
						Tokens: tk[53:56],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PrivateIdentifier: &tk[56],
								Tokens:            tk[56:57],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[58],
									Tokens:  tk[58:59],
								}),
								Tokens: tk[58:59],
							},
							Tokens: tk[56:59],
						},
						Tokens: tk[56:60],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[60],
									Tokens:              tk[60:61],
								},
								Tokens: tk[60:61],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[62],
									Tokens:  tk[62:63],
								}),
								Tokens: tk[62:63],
							},
							Tokens: tk[60:63],
						},
						Tokens: tk[60:63],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[64],
									Tokens:              tk[64:65],
								},
								Tokens: tk[64:65],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									Literal: &tk[66],
									Tokens:  tk[66:67],
								}),
								Tokens: tk[66:67],
							},
							Tokens: tk[64:67],
						},
						Tokens: tk[64:68],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[68],
									Tokens:              tk[68:69],
								},
								Tokens: tk[68:69],
							},
							Tokens: tk[68:69],
						},
						Tokens: tk[68:69],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									ComputedPropertyName: &AssignmentExpression{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[71],
											Tokens:              tk[71:72],
										}),
										Tokens: tk[71:72],
									},
									Tokens: tk[70:73],
								},
								Tokens: tk[70:73],
							},
							Tokens: tk[70:73],
						},
						Tokens: tk[70:73],
					},
				},
				Tokens: tk[:74],
			}
		}},
		{"class a {static{function}}", func(t *test, tk Tokens) { // 92
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrInvalidStatement,
								Parsing: "Statement",
								Token:   tk[7],
							},
							Parsing: "StatementListItem",
							Token:   tk[7],
						},
						Parsing: "Block",
						Token:   tk[7],
					},
					Parsing: "ClassElement",
					Token:   tk[6],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
		{"class a {a b}", func(t *test, tk Tokens) { // 93
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingSemiColon,
					Parsing: "ClassElement",
					Token:   tk[7],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
		{"class a {a/* */b}", func(t *test, tk Tokens) { // 94
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingSemiColon,
					Parsing: "ClassElement",
					Token:   tk[7],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
		{"class a {a/*\n*/b}", func(t *test, tk Tokens) { // 95
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[7],
									Tokens:              tk[7:8],
								},
								Tokens: tk[7:8],
							},
							Tokens: tk[7:8],
						},
						Tokens: tk[7:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"class a {a//\nc}", func(t *test, tk Tokens) { // 96
			t.Output = ClassDeclaration{
				BindingIdentifier: &tk[2],
				ClassBody: []ClassElement{
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[5],
									Tokens:              tk[5:6],
								},
								Tokens: tk[5:6],
							},
							Tokens: tk[5:6],
						},
						Tokens: tk[5:6],
					},
					{
						FieldDefinition: &FieldDefinition{
							ClassElementName: ClassElementName{
								PropertyName: &PropertyName{
									LiteralPropertyName: &tk[8],
									Tokens:              tk[8:9],
								},
								Tokens: tk[8:9],
							},
							Tokens: tk[8:9],
						},
						Tokens: tk[8:9],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{"class a {b()}", func(t *test, tk Tokens) { // 97
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningBrace,
							Parsing: "Block",
							Token:   tk[8],
						},
						Parsing: "MethodDefinition",
						Token:   tk[8],
					},
					Parsing: "ClassElement",
					Token:   tk[5],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
		{"class a {+}", func(t *test, tk Tokens) { // 98
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[5],
						},
						Parsing: "ClassElementName",
						Token:   tk[5],
					},
					Parsing: "ClassElement",
					Token:   tk[5],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
		{"class a {b = +}", func(t *test, tk Tokens) { // 99
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err: Error{
																					Err: Error{
																						Err: Error{
																							Err: Error{
																								Err: Error{
																									Err:     ErrNoIdentifier,
																									Parsing: "PrimaryExpression",
																									Token:   tk[10],
																								},
																								Parsing: "MemberExpression",
																								Token:   tk[10],
																							},
																							Parsing: "NewExpression",
																							Token:   tk[10],
																						},
																						Parsing: "LeftHandSideExpression",
																						Token:   tk[10],
																					},
																					Parsing: "UpdateExpression",
																					Token:   tk[10],
																				},
																				Parsing: "UnaryExpression",
																				Token:   tk[10],
																			},
																			Parsing: "ExponentiationExpression",
																			Token:   tk[9],
																		},
																		Parsing: "MultiplicativeExpression",
																		Token:   tk[9],
																	},
																	Parsing: "AdditiveExpression",
																	Token:   tk[9],
																},
																Parsing: "ShiftExpression",
																Token:   tk[9],
															},
															Parsing: "RelationalExpression",
															Token:   tk[9],
														},
														Parsing: "EqualityExpression",
														Token:   tk[9],
													},
													Parsing: "BitwiseANDExpression",
													Token:   tk[9],
												},
												Parsing: "BitwiseXORExpression",
												Token:   tk[9],
											},
											Parsing: "BitwiseORExpression",
											Token:   tk[9],
										},
										Parsing: "LogicalANDExpression",
										Token:   tk[9],
									},
									Parsing: "LogicalORExpression",
									Token:   tk[9],
								},
								Parsing: "ConditionalExpression",
								Token:   tk[9],
							},
							Parsing: "AssignmentExpression",
							Token:   tk[9],
						},
						Parsing: "FieldDefinition",
						Token:   tk[9],
					},
					Parsing: "ClassElement",
					Token:   tk[5],
				},
				Parsing: "ClassDeclaration",
				Token:   tk[5],
			}
		}},
	}, func(t *test) (Type, error) {
		var cd ClassDeclaration

		err := cd.parse(&t.Tokens, t.Yield, t.Await, t.Def)

		return cd, err
	})
}

func TestMethodDefinition(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidPropertyName,
						Parsing: "PropertyName",
						Token:   tk[0],
					},
					Parsing: "ClassElementName",
					Token:   tk[0],
				},
				Parsing: "MethodDefinition",
				Token:   tk[0],
			}
		}},
		{"get\n", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidPropertyName,
						Parsing: "PropertyName",
						Token:   tk[2],
					},
					Parsing: "ClassElementName",
					Token:   tk[2],
				},
				Parsing: "MethodDefinition",
				Token:   tk[2],
			}
		}},
		{"set\n", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidPropertyName,
						Parsing: "PropertyName",
						Token:   tk[2],
					},
					Parsing: "ClassElementName",
					Token:   tk[2],
				},
				Parsing: "MethodDefinition",
				Token:   tk[2],
			}
		}},
		{"get\na\n", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "MethodDefinition",
				Token:   tk[4],
			}
		}},
		{"a\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "FormalParameters",
					Token:   tk[2],
				},
				Parsing: "MethodDefinition",
				Token:   tk[2],
			}
		}},
		{"get\na\n(\na\n)", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "MethodDefinition",
				Token:   tk[6],
			}
		}},
		{"get\na\n(\n)\n", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[8],
				},
				Parsing: "MethodDefinition",
				Token:   tk[8],
			}
		}},
		{"set\na\n(\n)", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[6],
				},
				Parsing: "MethodDefinition",
				Token:   tk[4],
			}
		}},
		{"set\na\n(\nb\n)\n", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[10],
				},
				Parsing: "MethodDefinition",
				Token:   tk[10],
			}
		}},
		{"get\na\n(\n)\n{}", func(t *test, tk Tokens) { // 10
			t.Output = MethodDefinition{
				Type: MethodGetter,
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Params: FormalParameters{
					Tokens: tk[4:7],
				},
				FunctionBody: Block{
					Tokens: tk[8:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"set\na\n(\nb\n)\n{}", func(t *test, tk Tokens) { // 11
			t.Output = MethodDefinition{
				Type: MethodSetter,
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Params: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				FunctionBody: Block{
					Tokens: tk[10:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"a\n", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "FormalParameters",
					Token:   tk[2],
				},
				Parsing: "MethodDefinition",
				Token:   tk[2],
			}
		}},
		{"get\n()", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[4],
				},
				Parsing: "MethodDefinition",
				Token:   tk[4],
			}
		}},
		{"set\n()", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[4],
				},
				Parsing: "MethodDefinition",
				Token:   tk[4],
			}
		}},
		{"a\n(\n)\n", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[6],
				},
				Parsing: "MethodDefinition",
				Token:   tk[6],
			}
		}},
		{"a\n(\nb\n)\n", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[8],
				},
				Parsing: "MethodDefinition",
				Token:   tk[8],
			}
		}},
		{"a\n(\nb\n,\nc\n)\n", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[12],
				},
				Parsing: "MethodDefinition",
				Token:   tk[12],
			}
		}},
		{"a\n(\n)\n{}", func(t *test, tk Tokens) { // 18
			t.Output = MethodDefinition{
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Params: FormalParameters{
					Tokens: tk[2:5],
				},
				FunctionBody: Block{
					Tokens: tk[6:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"a\n(\nb\n)\n{}", func(t *test, tk Tokens) { // 19
			t.Output = MethodDefinition{
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Params: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &tk[4],
							Tokens:            tk[4:5],
						},
					},
					Tokens: tk[2:7],
				},
				FunctionBody: Block{
					Tokens: tk[8:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"a\n(\nb\n,\nc\n)\n{}", func(t *test, tk Tokens) { // 20
			t.Output = MethodDefinition{
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Params: FormalParameters{
					FormalParameterList: []BindingElement{
						{
							SingleNameBinding: &tk[4],
							Tokens:            tk[4:5],
						},
						{
							SingleNameBinding: &tk[8],
							Tokens:            tk[8:9],
						},
					},
					Tokens: tk[2:11],
				},
				FunctionBody: Block{
					Tokens: tk[12:14],
				},
				Tokens: tk[:14],
			}
		}},
		{"async a\n(\n)\n{}", func(t *test, tk Tokens) { // 21
			t.Output = MethodDefinition{
				Type: MethodAsync,
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Params: FormalParameters{
					Tokens: tk[4:7],
				},
				FunctionBody: Block{
					Tokens: tk[8:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"*\na\n(\n)\n{}", func(t *test, tk Tokens) { // 22
			t.Output = MethodDefinition{
				Type: MethodGenerator,
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Params: FormalParameters{
					Tokens: tk[4:7],
				},
				FunctionBody: Block{
					Tokens: tk[8:10],
				},
				Tokens: tk[:10],
			}
		}},
		{"async *\na\n(\n)\n{}", func(t *test, tk Tokens) { // 23
			t.Output = MethodDefinition{
				Type: MethodAsyncGenerator,
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[4],
						Tokens:              tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Params: FormalParameters{
					Tokens: tk[6:9],
				},
				FunctionBody: Block{
					Tokens: tk[10:12],
				},
				Tokens: tk[:12],
			}
		}},
		{"async\n(){}", func(t *test, tk Tokens) { // 24
			t.Output = MethodDefinition{
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Params: FormalParameters{
					Tokens: tk[2:4],
				},
				FunctionBody: Block{
					Tokens: tk[4:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"set a b", func(t *test, tk Tokens) { // 25
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "MethodDefinition",
				Token:   tk[4],
			}
		}},
		{"set a(b+c)", func(t *test, tk Tokens) { // 26
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "MethodDefinition",
				Token:   tk[5],
			}
		}},
		{"get(){}", func(t *test, tk Tokens) { // 27
			t.Output = MethodDefinition{
				ClassElementName: ClassElementName{
					PropertyName: &PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Params: FormalParameters{
					Tokens: tk[1:3],
				},
				FunctionBody: Block{
					Tokens: tk[3:5],
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
		var md MethodDefinition

		err := md.parse(&t.Tokens, t.Yield, t.Await)

		return md, err
	})
}

func TestPropertyName(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidPropertyName,
				Parsing: "PropertyName",
				Token:   tk[0],
			}
		}},
		{"[\n]", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "PropertyName",
				Token:   tk[2],
			}
		}},
		{"[\na\n]", func(t *test, tk Tokens) { // 3
			litA := makeConditionLiteral(tk, 2)
			t.Output = PropertyName{
				ComputedPropertyName: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 4
			t.Output = PropertyName{
				LiteralPropertyName: &tk[0],
				Tokens:              tk[:1],
			}
		}},
		{`43`, func(t *test, tk Tokens) { // 5
			t.Output = PropertyName{
				LiteralPropertyName: &tk[0],
				Tokens:              tk[:1],
			}
		}},
		{`"43"`, func(t *test, tk Tokens) { // 6
			t.Output = PropertyName{
				LiteralPropertyName: &tk[0],
				Tokens:              tk[:1],
			}
		}},
		{`null`, func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrInvalidPropertyName,
				Parsing: "PropertyName",
				Token:   tk[0],
			}
		}},
		{`[a, b]`, func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "PropertyName",
				Token:   tk[2],
			}
		}},
		{`await`, func(t *test, tk Tokens) { // 9
			t.Output = PropertyName{
				LiteralPropertyName: &tk[0],
				Tokens:              tk[:1],
			}
		}},
	}, func(t *test) (Type, error) {
		var pn PropertyName

		err := pn.parse(&t.Tokens, t.Yield, t.Await)

		return pn, err
	})
}
