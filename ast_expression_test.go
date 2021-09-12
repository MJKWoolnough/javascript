package javascript

import "testing"

func TestLeftHandSideExpressionOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`this`, func(t *test, tk Tokens) { // 1
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							This:   &tk[0],
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`someIdentifier`, func(t *test, tk Tokens) { // 2
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[0],
							Tokens:              tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`new someIdentifier`, func(t *test, tk Tokens) { // 3
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					News: 1,
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[2],
							Tokens:              tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`new new someIdentifier`, func(t *test, tk Tokens) { // 4
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					News: 2,
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[4],
							Tokens:              tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`null`, func(t *test, tk Tokens) { // 5
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[0],
							Tokens:  tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`true`, func(t *test, tk Tokens) { // 6
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[0],
							Tokens:  tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`false`, func(t *test, tk Tokens) { // 7
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[0],
							Tokens:  tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`0`, func(t *test, tk Tokens) { // 8
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[0],
							Tokens:  tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`"Hello"`, func(t *test, tk Tokens) { // 9
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[0],
							Tokens:  tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`[]`, func(t *test, tk Tokens) { // 10
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							ArrayLiteral: &ArrayLiteral{
								Tokens: tk[:2],
							},
							Tokens: tk[:2],
						},
						Tokens: tk[:2],
					},
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`{}`, func(t *test, tk Tokens) { // 11
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							ObjectLiteral: &ObjectLiteral{
								Tokens: tk[:2],
							},
							Tokens: tk[:2],
						},
						Tokens: tk[:2],
					},
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`super.runMe`, func(t *test, tk Tokens) { // 12
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						SuperProperty:  true,
						IdentifierName: &tk[2],
						Tokens:         tk[:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`super[runMe]`, func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 2)
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						SuperProperty: true,
						Expression: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`this.key.field.next`, func(t *test, tk Tokens) { // 14
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						MemberExpression: &MemberExpression{
							MemberExpression: &MemberExpression{
								MemberExpression: &MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										This:   &tk[0],
										Tokens: tk[:1],
									},
									Tokens: tk[:1],
								},
								IdentifierName: &tk[2],
								Tokens:         tk[:3],
							},
							IdentifierName: &tk[4],
							Tokens:         tk[:5],
						},
						IdentifierName: &tk[6],
						Tokens:         tk[:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`new.target`, func(t *test, tk Tokens) { // 15
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						NewTarget: true,
						Tokens:    tk[:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`new className()`, func(t *test, tk Tokens) { // 16
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[2],
								Tokens:              tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Arguments: &Arguments{
							Tokens: tk[3:5],
						},
						Tokens: tk[:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`new new className()`, func(t *test, tk Tokens) { // 17
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					News: 1,
					MemberExpression: MemberExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[4],
								Tokens:              tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Arguments: &Arguments{
							Tokens: tk[5:7],
						},
						Tokens: tk[2:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`new new className()()`, func(t *test, tk Tokens) { // 18
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						MemberExpression: &MemberExpression{
							MemberExpression: &MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[4],
									Tokens:              tk[4:5],
								},
								Tokens: tk[4:5],
							},
							Arguments: &Arguments{
								Tokens: tk[5:7],
							},
							Tokens: tk[2:7],
						},
						Arguments: &Arguments{
							Tokens: tk[7:9],
						},
						Tokens: tk[:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{`call()`, func(t *test, tk Tokens) { // 19
			t.Output = LeftHandSideExpression{
				CallExpression: &CallExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[0],
							Tokens:              tk[:1],
						},
						Tokens: tk[:1],
					},
					Arguments: &Arguments{
						Tokens: tk[1:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a\n?.\nb\n.\nc\n?.\nd", func(t *test, tk Tokens) { // 20
			t.Output = LeftHandSideExpression{
				OptionalExpression: &OptionalExpression{
					OptionalExpression: &OptionalExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						OptionalChain: OptionalChain{
							OptionalChain: &OptionalChain{
								IdentifierName: &tk[4],
								Tokens:         tk[2:5],
							},
							IdentifierName: &tk[8],
							Tokens:         tk[2:9],
						},
						Tokens: tk[:9],
					},
					OptionalChain: OptionalChain{
						IdentifierName: &tk[12],
						Tokens:         tk[10:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"a\n()\n?.\nb\n.\nc\n?.\nd", func(t *test, tk Tokens) { // 21
			t.Output = LeftHandSideExpression{
				OptionalExpression: &OptionalExpression{
					OptionalExpression: &OptionalExpression{
						CallExpression: &CallExpression{
							MemberExpression: &MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[0],
									Tokens:              tk[:1],
								},
								Tokens: tk[:1],
							},
							Arguments: &Arguments{
								Tokens: tk[2:4],
							},
							Tokens: tk[:4],
						},
						OptionalChain: OptionalChain{
							OptionalChain: &OptionalChain{
								IdentifierName: &tk[7],
								Tokens:         tk[5:8],
							},
							IdentifierName: &tk[11],
							Tokens:         tk[5:12],
						},
						Tokens: tk[:12],
					},
					OptionalChain: OptionalChain{
						IdentifierName: &tk[15],
						Tokens:         tk[13:16],
					},
					Tokens: tk[:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"a()?.", func(t *test, tk Tokens) { // 22
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidOptionalChain,
						Parsing: "OptionalChain",
						Token:   tk[4],
					},
					Parsing: "OptionalExpression",
					Token:   tk[3],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[3],
			}
		}},
		{"a?.", func(t *test, tk Tokens) { // 22
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidOptionalChain,
						Parsing: "OptionalChain",
						Token:   tk[2],
					},
					Parsing: "OptionalExpression",
					Token:   tk[1],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var lhs LeftHandSideExpression
		err := lhs.parse(&t.Tokens, t.Yield, t.Await)
		return lhs, err
	})
}

func TestAssignmentExpressionOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`yield 1`, func(t *test, tk Tokens) {}}, // 1
		{`yield 1`, func(t *test, tk Tokens) { // 2
			t.Yield = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = AssignmentExpression{
				Yield: true,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`yield *1`, func(t *test, tk Tokens) { // 3
			t.Yield = true
			litA := makeConditionLiteral(tk, 3)
			t.Output = AssignmentExpression{
				Yield:    true,
				Delegate: true,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[3:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`a => {}`, func(t *test, tk Tokens) { // 4
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &tk[0],
					FunctionBody: &Block{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`a => 1`, func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &tk[0],
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[4:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`async a => 1`, func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					Async:             true,
					BindingIdentifier: &tk[2],
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[6:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`() => {}`, func(t *test, tk Tokens) { // 7
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						Tokens: tk[:2],
					},
					FunctionBody: &Block{
						Tokens: tk[5:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`(a) => b`, func(t *test, tk Tokens) { // 8
			litB := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
						},
						Tokens: tk[:3],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litB,
						Tokens:                tk[6:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`(a, b) => c`, func(t *test, tk Tokens) { // 9
			litC := makeConditionLiteral(tk, 9)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
							{
								SingleNameBinding: &tk[4],
								Tokens:            tk[4:5],
							},
						},
						Tokens: tk[:6],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litC,
						Tokens:                tk[9:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`(a, b, c) => d`, func(t *test, tk Tokens) { // 10
			litD := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
							{
								SingleNameBinding: &tk[4],
								Tokens:            tk[4:5],
							},
							{
								SingleNameBinding: &tk[7],
								Tokens:            tk[7:8],
							},
						},
						Tokens: tk[:9],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litD,
						Tokens:                tk[12:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`(a, ...b) => c`, func(t *test, tk Tokens) { // 11
			litC := makeConditionLiteral(tk, 10)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
						},
						BindingIdentifier: &tk[5],
						Tokens:            tk[:7],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litC,
						Tokens:                tk[10:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{`(a, ...[b]) => c`, func(t *test, tk Tokens) { // 12
			litC := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
						},
						ArrayBindingPattern: &ArrayBindingPattern{
							BindingElementList: []BindingElement{
								{
									SingleNameBinding: &tk[6],
									Tokens:            tk[6:7],
								},
							},
							Tokens: tk[5:8],
						},
						Tokens: tk[:9],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litC,
						Tokens:                tk[12:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`(a, ...{b}) => c`, func(t *test, tk Tokens) { // 13
			litC := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[1],
								Tokens:            tk[1:2],
							},
						},
						ObjectBindingPattern: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									PropertyName: PropertyName{
										LiteralPropertyName: &tk[6],
										Tokens:              tk[6:7],
									},
									BindingElement: BindingElement{
										SingleNameBinding: &tk[6],
										Tokens:            tk[6:7],
									},
									Tokens: tk[6:7],
								},
							},
							Tokens: tk[5:8],
						},
						Tokens: tk[:9],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litC,
						Tokens:                tk[12:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`a = 1`, func(t *test, tk Tokens) { // 14
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a *= 1`, func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentMultiply,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a /= 1`, func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentDivide,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a %= 1`, func(t *test, tk Tokens) { // 17
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentRemainder,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a += 1`, func(t *test, tk Tokens) { // 18
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentAdd,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a -= 1`, func(t *test, tk Tokens) { // 19
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentSubtract,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a <<= 1`, func(t *test, tk Tokens) { // 20
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentLeftShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a >>= 1`, func(t *test, tk Tokens) { // 21
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentSignPropagatinRightShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a >>>= 1`, func(t *test, tk Tokens) { // 22
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentZeroFillRightShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a &= 1`, func(t *test, tk Tokens) { // 23
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseAND,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a ^= 1`, func(t *test, tk Tokens) { // 24
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseXOR,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a |= 1`, func(t *test, tk Tokens) { // 25
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseOR,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`a **= 1`, func(t *test, tk Tokens) { // 26
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentExponentiation,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`import(a)`, func(t *test, tk Tokens) { // 27
			litA := makeConditionLiteral(tk, 2)
			call := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					CallExpression: &CallExpression{
						ImportCall: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
						Tokens: tk[0:4],
					},
					Tokens: tk[0:4],
				},
				Tokens: tk[0:4],
			})
			t.Output = AssignmentExpression{
				ConditionalExpression: &call,
				Tokens:                tk[0:4],
			}
		}},
	}, func(t *test) (Type, error) {
		var ae AssignmentExpression
		err := ae.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return ae, err
	})
}

func TestAssignmentExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = assignmentError(tk[0])
		}},
		{"yield", func(t *test, tk Tokens) { // 2
			litYield := makeConditionLiteral(tk, 0)
			t.Output = AssignmentExpression{
				ConditionalExpression: &litYield,
				Tokens:                tk[:1],
			}
		}},
		{"yield", func(t *test, tk Tokens) { // 3
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[1]),
				Parsing: "AssignmentExpression",
				Token:   tk[1],
			}
		}},
		{"yield ,", func(t *test, tk Tokens) { // 4
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "AssignmentExpression",
				Token:   tk[2],
			}
		}},
		{"yield 1", func(t *test, tk Tokens) { // 5
			lit1 := makeConditionLiteral(tk, 2)
			t.Yield = true
			t.Output = AssignmentExpression{
				Yield: true,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"yield *\n", func(t *test, tk Tokens) { // 6
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "AssignmentExpression",
				Token:   tk[4],
			}
		}},
		{"yield *\n*", func(t *test, tk Tokens) { // 7
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "AssignmentExpression",
				Token:   tk[4],
			}
		}},
		{"yield *\n1", func(t *test, tk Tokens) { // 8
			lit1 := makeConditionLiteral(tk, 4)
			t.Yield = true
			t.Output = AssignmentExpression{
				Yield:    true,
				Delegate: true,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"async", func(t *test, tk Tokens) { // 9
			t.Err = assignmentCustomError(tk[0], Error{
				Err:     ErrInvalidFunction,
				Parsing: "FunctionDeclaration",
				Token:   tk[1],
			})
		}},
		{"async a => {}", func(t *test, tk Tokens) { // 10
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					Async:             true,
					BindingIdentifier: &tk[2],
					FunctionBody: &Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"() => {}", func(t *test, tk Tokens) { // 11
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						Tokens: tk[:2],
					},
					FunctionBody: &Block{
						Tokens: tk[5:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"a => {}", func(t *test, tk Tokens) { // 12
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &tk[0],
					FunctionBody: &Block{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"() => a", func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 5)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					FormalParameters: &FormalParameters{
						Tokens: tk[:2],
					},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[5:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"a => b", func(t *test, tk Tokens) { // 14
			litB := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &tk[0],
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litB,
						Tokens:                tk[4:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a =>", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[3]),
					Parsing: "ArrowFunction",
					Token:   tk[3],
				},
				Parsing: "AssignmentExpression",
				Token:   tk[0],
			}
		}},
		{"1", func(t *test, tk Tokens) { // 16
			lit1 := makeConditionLiteral(tk, 0)
			t.Output = AssignmentExpression{
				ConditionalExpression: &lit1,
				Tokens:                tk[:1],
			}
		}},
		{",", func(t *test, tk Tokens) { // 17
			t.Err = assignmentError(tk[0])
		}},
		{"a\n=", func(t *test, tk Tokens) { // 18
			t.Err = Error{
				Err:     assignmentError(tk[3]),
				Parsing: "AssignmentExpression",
				Token:   tk[3],
			}
		}},
		{"a\n=\n1", func(t *test, tk Tokens) { // 19
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"[a] = 1", func(t *test, tk Tokens) { // 20
			lit1 := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				AssignmentPattern: &AssignmentPattern{
					ArrayAssignmentPattern: &ArrayAssignmentPattern{
						AssignmentElements: []AssignmentElement{
							{
								DestructuringAssignmentTarget: DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[1],
													Tokens:              tk[1:2],
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
										Tokens: tk[1:2],
									},
									Tokens: tk[1:2],
								},
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[:3],
					},
					Tokens: tk[:3],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[6:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"{a} = 1", func(t *test, tk Tokens) { // 21
			lit1 := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				AssignmentPattern: &AssignmentPattern{
					ObjectAssignmentPattern: &ObjectAssignmentPattern{
						AssignmentPropertyList: []AssignmentProperty{
							{
								PropertyName: PropertyName{
									LiteralPropertyName: &tk[1],
									Tokens:              tk[1:2],
								},
								DestructuringAssignmentTarget: &DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[1],
													Tokens:              tk[1:2],
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
										Tokens: tk[1:2],
									},
									Tokens: tk[1:2],
								},
								Tokens: tk[1:2],
							},
						},
						Tokens: tk[:3],
					},
					Tokens: tk[:3],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[6:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"a *= 1", func(t *test, tk Tokens) { // 22
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentMultiply,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a /= 1", func(t *test, tk Tokens) { // 23
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentDivide,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a %= 1", func(t *test, tk Tokens) { // 24
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentRemainder,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a += 1", func(t *test, tk Tokens) { // 25
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentAdd,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a -= 1", func(t *test, tk Tokens) { // 26
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentSubtract,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a <<= 1", func(t *test, tk Tokens) { // 27
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentLeftShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a >>= 1", func(t *test, tk Tokens) { // 28
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentSignPropagatinRightShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a >>>= 1", func(t *test, tk Tokens) { // 29
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentZeroFillRightShift,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a &= 1", func(t *test, tk Tokens) { // 30
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseAND,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a ^= 1", func(t *test, tk Tokens) { // 31
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseXOR,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a |= 1", func(t *test, tk Tokens) { // 32
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentBitwiseOR,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a **= 1", func(t *test, tk Tokens) { // 33
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentExponentiation,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a &&= 1", func(t *test, tk Tokens) { // 34
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentLogicalAnd,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a ||= 1", func(t *test, tk Tokens) { // 35
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentLogicalOr,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a ??= 1", func(t *test, tk Tokens) { // 36
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				AssignmentOperator: AssignmentNullish,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"[a, b] = [b, a]", func(t *test, tk Tokens) { // 37
			t.Output = AssignmentExpression{
				AssignmentPattern: &AssignmentPattern{
					ArrayAssignmentPattern: &ArrayAssignmentPattern{
						AssignmentElements: []AssignmentElement{
							{
								DestructuringAssignmentTarget: DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[1],
													Tokens:              tk[1:2],
												},
												Tokens: tk[1:2],
											},
											Tokens: tk[1:2],
										},
										Tokens: tk[1:2],
									},
									Tokens: tk[1:2],
								},
								Tokens: tk[1:2],
							},
							{
								DestructuringAssignmentTarget: DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
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
									Tokens: tk[4:5],
								},
								Tokens: tk[4:5],
							},
						},
						Tokens: tk[:6],
					},
					Tokens: tk[:6],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&ArrayLiteral{
						ElementList: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								}),
								Tokens: tk[10:11],
							},
							{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									IdentifierReference: &tk[13],
									Tokens:              tk[13:14],
								}),
								Tokens: tk[13:14],
							},
						},
						Tokens: tk[9:15],
					}),
					Tokens: tk[9:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"[a.b, a.c] = [a.c, a.b]", func(t *test, tk Tokens) { // 38
			t.Output = AssignmentExpression{
				AssignmentPattern: &AssignmentPattern{
					ArrayAssignmentPattern: &ArrayAssignmentPattern{
						AssignmentElements: []AssignmentElement{
							{
								DestructuringAssignmentTarget: DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[1],
														Tokens:              tk[1:2],
													},
													Tokens: tk[1:2],
												},
												IdentifierName: &tk[3],
												Tokens:         tk[1:4],
											},
											Tokens: tk[1:4],
										},
										Tokens: tk[1:4],
									},
									Tokens: tk[1:4],
								},
								Tokens: tk[1:4],
							},
							{
								DestructuringAssignmentTarget: DestructuringAssignmentTarget{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[6],
														Tokens:              tk[6:7],
													},
													Tokens: tk[6:7],
												},
												IdentifierName: &tk[8],
												Tokens:         tk[6:9],
											},
											Tokens: tk[6:9],
										},
										Tokens: tk[6:9],
									},
									Tokens: tk[6:9],
								},
								Tokens: tk[6:9],
							},
						},
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				AssignmentOperator: AssignmentAssign,
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&ArrayLiteral{
						ElementList: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									MemberExpression: &MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: &tk[14],
											Tokens:              tk[14:15],
										},
										Tokens: tk[14:15],
									},
									IdentifierName: &tk[16],
									Tokens:         tk[14:17],
								}),
								Tokens: tk[14:17],
							},
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									MemberExpression: &MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: &tk[19],
											Tokens:              tk[19:20],
										},
										Tokens: tk[19:20],
									},
									IdentifierName: &tk[21],
									Tokens:         tk[19:22],
								}),
								Tokens: tk[19:22],
							},
						},
						Tokens: tk[13:23],
					}),
					Tokens: tk[13:23],
				},
				Tokens: tk[:23],
			}
		}},
		{"async (1) => 1", func(t *test, tk Tokens) { // 38
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "BindingElement",
							Token:   tk[3],
						},
						Parsing: "FormalParameters",
						Token:   tk[3],
					},
					Parsing: "ArrowFunction",
					Token:   tk[2],
				},
				Parsing: "AssignmentExpression",
				Token:   tk[0],
			}
		}},
		{"(...a)", func(t *test, tk Tokens) { // 39
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "AssignmentExpression",
				Token:   tk[4],
			}
		}},
		{"[1] =", func(t *test, tk Tokens) { // 40
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrInvalidDestructuringAssignmentTarget,
								Parsing: "DestructuringAssignmentTarget",
								Token:   tk[1],
							},
							Parsing: "AssignmentElement",
							Token:   tk[1],
						},
						Parsing: "ArrayAssignmentPattern",
						Token:   tk[0],
					},
					Parsing: "AssignmentPattern",
					Token:   tk[0],
				},
				Parsing: "AssignmentExpression",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var ae AssignmentExpression
		err := ae.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return ae, err
	})
}

func TestLeftHandSideExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "PrimaryExpression",
							Token:   tk[0],
						},
						Parsing: "MemberExpression",
						Token:   tk[0],
					},
					Parsing: "NewExpression",
					Token:   tk[0],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{"super\n(,)", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[3]),
						Parsing: "Arguments",
						Token:   tk[3],
					},
					Parsing: "CallExpression",
					Token:   tk[2],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{"super\n()", func(t *test, tk Tokens) { // 3
			t.Output = LeftHandSideExpression{
				CallExpression: &CallExpression{
					SuperCall: true,
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"import\n(,)", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[3]),
					Parsing: "CallExpression",
					Token:   tk[3],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{"import\n(a)", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 3)
			t.Output = LeftHandSideExpression{
				CallExpression: &CallExpression{
					ImportCall: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[3:4],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`import(a).then(b)`, func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			t.Output = LeftHandSideExpression{
				CallExpression: &CallExpression{
					CallExpression: &CallExpression{
						CallExpression: &CallExpression{
							ImportCall: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
							Tokens: tk[0:4],
						},
						IdentifierName: &tk[5],
						Tokens:         tk[0:6],
					},
					Arguments: &Arguments{
						ArgumentList: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[7:8],
							},
						},
						Tokens: tk[6:9],
					},
					Tokens: tk[0:9],
				},
				Tokens: tk[0:9],
			}
		}},
		{`super`, func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidSuperProperty,
						Parsing: "MemberExpression",
						Token:   tk[1],
					},
					Parsing: "NewExpression",
					Token:   tk[0],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{`import`, func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "PrimaryExpression",
							Token:   tk[0],
						},
						Parsing: "MemberExpression",
						Token:   tk[0],
					},
					Parsing: "NewExpression",
					Token:   tk[0],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{",", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "PrimaryExpression",
							Token:   tk[0],
						},
						Parsing: "MemberExpression",
						Token:   tk[0],
					},
					Parsing: "NewExpression",
					Token:   tk[0],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 10
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[0],
							Tokens:              tk[:1],
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"a\n()", func(t *test, tk Tokens) { // 11
			t.Output = LeftHandSideExpression{
				CallExpression: &CallExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[0],
							Tokens:              tk[:1],
						},
						Tokens: tk[:1],
					},
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"a\n(,)", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[3]),
						Parsing: "Arguments",
						Token:   tk[3],
					},
					Parsing: "CallExpression",
					Token:   tk[2],
				},
				Parsing: "LeftHandSideExpression",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var lhs LeftHandSideExpression
		err := lhs.parse(&t.Tokens, t.Yield, t.Await)
		return lhs, err
	})
}

func TestExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     assignmentError(tk[0]),
				Parsing: "Expression",
				Token:   tk[0],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 2
			litA := makeConditionLiteral(tk, 0)
			t.Output = Expression{
				Expressions: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[:1],
					},
				},
				Tokens: tk[:1],
			}
		}},
		{"a\n,\n", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "Expression",
				Token:   tk[4],
			}
		}},
		{"a\n,\nb", func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = Expression{
				Expressions: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[:1],
					},
					{
						ConditionalExpression: &litB,
						Tokens:                tk[4:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
		var e Expression
		err := e.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return e, err
	})
}

func TestNewExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "PrimaryExpression",
						Token:   tk[0],
					},
					Parsing: "MemberExpression",
					Token:   tk[0],
				},
				Parsing: "NewExpression",
				Token:   tk[0],
			}
		}},
		{",", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "PrimaryExpression",
						Token:   tk[0],
					},
					Parsing: "MemberExpression",
					Token:   tk[0],
				},
				Parsing: "NewExpression",
				Token:   tk[0],
			}
		}},
		{"new\n,", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "PrimaryExpression",
							Token:   tk[2],
						},
						Parsing: "MemberExpression",
						Token:   tk[2],
					},
					Parsing: "MemberExpression",
					Token:   tk[0],
				},
				Parsing: "NewExpression",
				Token:   tk[0],
			}
		}},
		{"1", func(t *test, tk Tokens) { // 4
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						Literal: &tk[0],
						Tokens:  tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"new\na", func(t *test, tk Tokens) { // 5
			t.Output = NewExpression{
				News: 1,
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &tk[2],
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"new\na\n()", func(t *test, tk Tokens) { // 6
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[2],
							Tokens:              tk[2:3],
						},
						Tokens: tk[2:3],
					},
					Arguments: &Arguments{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"new\nnew\na\n()", func(t *test, tk Tokens) { // 7
			t.Output = NewExpression{
				News: 1,
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[4],
							Tokens:              tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Arguments: &Arguments{
						Tokens: tk[6:8],
					},
					Tokens: tk[2:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"new new\nnew\na\n()", func(t *test, tk Tokens) { // 8
			t.Output = NewExpression{
				News: 2,
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[6],
							Tokens:              tk[6:7],
						},
						Tokens: tk[6:7],
					},
					Arguments: &Arguments{
						Tokens: tk[8:10],
					},
					Tokens: tk[4:10],
				},
				Tokens: tk[:10],
			}
		}},
	}, func(t *test) (Type, error) {
		var ne NewExpression
		err := ne.parse(&t.Tokens, t.Yield, t.Await)
		return ne, err
	})
}

func TestMemberExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "PrimaryExpression",
					Token:   tk[0],
				},
				Parsing: "MemberExpression",
				Token:   tk[0],
			}
		}},
		{"super\n[\n,\n]", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "MemberExpression",
				Token:   tk[4],
			}
		}},
		{"super\n[\n1\n]", func(t *test, tk Tokens) { // 3
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = MemberExpression{
				SuperProperty: true,
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &lit1,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"super\n[\n1\n2\n]", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrInvalidSuperProperty,
				Parsing: "MemberExpression",
				Token:   tk[6],
			}
		}},
		{"super\n.\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "MemberExpression",
				Token:   tk[4],
			}
		}},
		{"super\n.\n1", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "MemberExpression",
				Token:   tk[4],
			}
		}},
		{"super\n.\na", func(t *test, tk Tokens) { // 7
			t.Output = MemberExpression{
				SuperProperty:  true,
				IdentifierName: &tk[4],
				Tokens:         tk[:5],
			}
		}},
		{"super\n", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrInvalidSuperProperty,
				Parsing: "MemberExpression",
				Token:   tk[2],
			}
		}},
		{"new", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "PrimaryExpression",
						Token:   tk[1],
					},
					Parsing: "MemberExpression",
					Token:   tk[1],
				},
				Parsing: "MemberExpression",
				Token:   tk[0],
			}
		}},
		{"new\n.\n", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrInvalidMetaProperty,
				Parsing: "MemberExpression",
				Token:   tk[0],
			}
		}},
		{"new\n.\ntarget", func(t *test, tk Tokens) { // 11
			t.Output = MemberExpression{
				NewTarget: true,
				Tokens:    tk[:5],
			}
		}},
		{"new\n,", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "PrimaryExpression",
						Token:   tk[2],
					},
					Parsing: "MemberExpression",
					Token:   tk[2],
				},
				Parsing: "MemberExpression",
				Token:   tk[0],
			}
		}},
		{"new\n1\n", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "Arguments",
					Token:   tk[4],
				},
				Parsing: "MemberExpression",
				Token:   tk[3],
			}
		}},
		{"new\n1\n()", func(t *test, tk Tokens) { // 14
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						Literal: &tk[2],
						Tokens:  tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Arguments: &Arguments{
					Tokens: tk[4:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"new\nnew\n1\n()", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "Arguments",
					Token:   tk[8],
				},
				Parsing: "MemberExpression",
				Token:   tk[8],
			}
		}},
		{"new\nnew\n1\n()\n()", func(t *test, tk Tokens) { // 16
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[4],
							Tokens:  tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Arguments: &Arguments{
						Tokens: tk[6:8],
					},
					Tokens: tk[2:8],
				},
				Arguments: &Arguments{
					Tokens: tk[9:11],
				},
				Tokens: tk[:11],
			}
		}},
		{",", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "PrimaryExpression",
					Token:   tk[0],
				},
				Parsing: "MemberExpression",
				Token:   tk[0],
			}
		}},
		{"1", func(t *test, tk Tokens) { // 18
			t.Output = MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					Literal: &tk[0],
					Tokens:  tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 18
			t.Output = MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					IdentifierReference: &tk[0],
					Tokens:              tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"a\n`${\n1\n1\n}`", func(t *test, tk Tokens) { // 19
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidTemplate,
					Parsing: "TemplateLiteral",
					Token:   tk[6],
				},
				Parsing: "MemberExpression",
				Token:   tk[2],
			}
		}},
		{"a\n``", func(t *test, tk Tokens) { // 20
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[2],
					Tokens:                 tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a\n.\n", func(t *test, tk Tokens) { // 21
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "MemberExpression",
				Token:   tk[2],
			}
		}},
		{"a\n.\nb", func(t *test, tk Tokens) { // 22
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				IdentifierName: &tk[4],
				Tokens:         tk[:5],
			}
		}},
		{"a\n[\n]", func(t *test, tk Tokens) { // 23
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "MemberExpression",
				Token:   tk[2],
			}
		}},
		{"a\n[\n1\n2\n]", func(t *test, tk Tokens) { // 24
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "MemberExpression",
				Token:   tk[2],
			}
		}},
		{"a\n[\n1\n]", func(t *test, tk Tokens) { // 25
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &tk[0],
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &lit1,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"a\n.\nb\n[\nc\n]\n``", func(t *test, tk Tokens) { // 26
			litC := makeConditionLiteral(tk, 8)
			t.Output = MemberExpression{
				MemberExpression: &MemberExpression{
					MemberExpression: &MemberExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[0],
								Tokens:              tk[:1],
							},
							Tokens: tk[:1],
						},
						IdentifierName: &tk[4],
						Tokens:         tk[:5],
					},
					Expression: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litC,
								Tokens:                tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					Tokens: tk[:11],
				},
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[12],
					Tokens:                 tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"import . meta", func(t *test, tk Tokens) { // 27
			t.Output = MemberExpression{
				ImportMeta: true,
				Tokens:     tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
		var me MemberExpression
		err := me.parse(&t.Tokens, t.Yield, t.Await)
		return me, err
	})
}

func TestPrimaryExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{`this`, func(t *test, tk Tokens) { // 2
			t.Output = PrimaryExpression{
				This:   &tk[0],
				Tokens: tk[:1],
			}
		}},
		{`null`, func(t *test, tk Tokens) { // 3
			t.Output = PrimaryExpression{
				Literal: &tk[0],
				Tokens:  tk[:1],
			}
		}},
		{`true`, func(t *test, tk Tokens) { // 4
			t.Output = PrimaryExpression{
				Literal: &tk[0],
				Tokens:  tk[:1],
			}
		}},
		{`1.234`, func(t *test, tk Tokens) { // 5
			t.Output = PrimaryExpression{
				Literal: &tk[0],
				Tokens:  tk[:1],
			}
		}},
		{`"string"`, func(t *test, tk Tokens) { // 6
			t.Output = PrimaryExpression{
				Literal: &tk[0],
				Tokens:  tk[:1],
			}
		}},
		{`/a/`, func(t *test, tk Tokens) { // 7
			t.Output = PrimaryExpression{
				Literal: &tk[0],
				Tokens:  tk[:1],
			}
		}},
		{`[yield]`, func(t *test, tk Tokens) { // 8
			t.Yield = true
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[2]),
						Parsing: "AssignmentExpression",
						Token:   tk[2],
					},
					Parsing: "ArrayLiteral",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{`[]`, func(t *test, tk Tokens) { // 9
			t.Output = PrimaryExpression{
				ArrayLiteral: &ArrayLiteral{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`{,}`, func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[1],
						},
						Parsing: "PropertyDefinition",
						Token:   tk[1],
					},
					Parsing: "ObjectLiteral",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{`{}`, func(t *test, tk Tokens) { // 11
			t.Output = PrimaryExpression{
				ObjectLiteral: &ObjectLiteral{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"async", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidFunction,
					Parsing: "FunctionDeclaration",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"async function(){}", func(t *test, tk Tokens) { // 13
			t.Output = PrimaryExpression{
				FunctionExpression: &FunctionDeclaration{
					Type: FunctionAsync,
					FormalParameters: FormalParameters{
						Tokens: tk[3:5],
					},
					FunctionBody: Block{
						Tokens: tk[5:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"function", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningParenthesis,
						Parsing: "FormalParameters",
						Token:   tk[1],
					},
					Parsing: "FunctionDeclaration",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"function(){}", func(t *test, tk Tokens) { // 15
			t.Output = PrimaryExpression{
				FunctionExpression: &FunctionDeclaration{
					FormalParameters: FormalParameters{
						Tokens: tk[1:3],
					},
					FunctionBody: Block{
						Tokens: tk[3:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"class", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "ClassDeclaration",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"class{}", func(t *test, tk Tokens) { // 17
			t.Output = PrimaryExpression{
				ClassExpression: &ClassDeclaration{
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"``", func(t *test, tk Tokens) { // 18
			t.Output = PrimaryExpression{
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[0],
					Tokens:                 tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"`${1 1}`", func(t *test, tk Tokens) { // 19
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidTemplate,
					Parsing: "TemplateLiteral",
					Token:   tk[3],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"(,)", func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[1]),
					Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
					Token:   tk[1],
				},
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"()", func(t *test, tk Tokens) { // 21
			t.Output = PrimaryExpression{
				CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{".", func(t *test, tk Tokens) { // 22
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "PrimaryExpression",
				Token:   tk[0],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 23
			t.Output = PrimaryExpression{
				IdentifierReference: &tk[0],
				Tokens:              tk[:1],
			}
		}},
	}, func(t *test) (Type, error) {
		var pe PrimaryExpression
		err := pe.parse(&t.Tokens, t.Yield, t.Await)
		return pe, err
	})
}

func TestCoverParenthesizedExpressionAndArrowParameterList(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[0],
			}
		}},
		{"(\n)", func(t *test, tk Tokens) { // 2
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				Tokens: tk[:3],
			}
		}},
		{"(\n...\n[{,}])", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrInvalidPropertyName,
									Parsing: "PropertyName",
									Token:   tk[6],
								},
								Parsing: "BindingProperty",
								Token:   tk[6],
							},
							Parsing: "ObjectBindingPattern",
							Token:   tk[6],
						},
						Parsing: "BindingElement",
						Token:   tk[5],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[5],
				},
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[4],
			}
		}},
		{"(\n...\n[a]\n)", func(t *test, tk Tokens) { // 4
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				arrayBindingPattern: &ArrayBindingPattern{
					BindingElementList: []BindingElement{
						{
							SingleNameBinding: &tk[5],
							Tokens:            tk[5:6],
						},
					},
					Tokens: tk[4:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"(\n...\n{,})", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[5],
						},
						Parsing: "BindingProperty",
						Token:   tk[5],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[5],
				},
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[4],
			}
		}},
		{"(\n...\n{}\n)", func(t *test, tk Tokens) { // 6
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				objectBindingPattern: &ObjectBindingPattern{
					Tokens: tk[4:6],
				},
				Tokens: tk[:8],
			}
		}},
		{"(\n...\n1\n)", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[4],
			}
		}},
		{"(\n...\na\n)", func(t *test, tk Tokens) { // 8
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				bindingIdentifier: &tk[4],
				Tokens:            tk[:7],
			}
		}},
		{"(\n...\na\n,)", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[6],
			}
		}},
		{"(\n,)", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[2],
			}
		}},
		{"(\n1\n)", func(t *test, tk Tokens) { // 11
			lit1 := makeConditionLiteral(tk, 2)
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				Expressions: []AssignmentExpression{
					{
						ConditionalExpression: &lit1,
						Tokens:                tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"(\n1\n2)", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "CoverParenthesizedExpressionAndArrowParameterList",
				Token:   tk[4],
			}
		}},
		{"(\n1\n,\n2\n)", func(t *test, tk Tokens) { // 13
			lit1 := makeConditionLiteral(tk, 2)
			lit2 := makeConditionLiteral(tk, 6)
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				Expressions: []AssignmentExpression{
					{
						ConditionalExpression: &lit1,
						Tokens:                tk[2:3],
					},
					{
						ConditionalExpression: &lit2,
						Tokens:                tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"(\n1\n,\n...\na\n)", func(t *test, tk Tokens) { // 14
			lit1 := makeConditionLiteral(tk, 2)
			t.Output = CoverParenthesizedExpressionAndArrowParameterList{
				Expressions: []AssignmentExpression{
					{
						ConditionalExpression: &lit1,
						Tokens:                tk[2:3],
					},
				},
				bindingIdentifier: &tk[8],
				Tokens:            tk[:11],
			}
		}},
	}, func(t *test) (Type, error) {
		var c CoverParenthesizedExpressionAndArrowParameterList
		err := c.parse(&t.Tokens, t.Yield, t.Await)
		return c, err
	})
}

func TestArguments(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "Arguments",
				Token:   tk[0],
			}
		}},
		{"(\n)", func(t *test, tk Tokens) { // 2
			t.Output = Arguments{
				Tokens: tk[:3],
			}
		}},
		{"(\n...\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "Arguments",
				Token:   tk[2],
			}
		}},
		{"(\n...\na\n)", func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 4)
			t.Output = Arguments{
				SpreadArgument: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"(\n...\na\nb)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "Arguments",
				Token:   tk[6],
			}
		}},
		{"(\n,)", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "Arguments",
				Token:   tk[2],
			}
		}},
		{"(\na\n)", func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 2)
			t.Output = Arguments{
				ArgumentList: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"(\na\nb)", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "Arguments",
				Token:   tk[4],
			}
		}},
		{"(\na\n,\nb\n)", func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			t.Output = Arguments{
				ArgumentList: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
					{
						ConditionalExpression: &litB,
						Tokens:                tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"(\na\n,\n...\nb\n)", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			t.Output = Arguments{
				ArgumentList: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
				},
				SpreadArgument: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[8:9],
				},
				Tokens: tk[:11],
			}
		}},
		{"(async function(){})", func(t *test, tk Tokens) { // 11
			t.Output = Arguments{
				ArgumentList: []AssignmentExpression{
					{
						ConditionalExpression: WrapConditional(PrimaryExpression{
							FunctionExpression: &FunctionDeclaration{
								Type: FunctionAsync,
								FormalParameters: FormalParameters{
									Tokens: tk[4:6],
								},
								FunctionBody: Block{
									Tokens: tk[6:8],
								},
								Tokens: tk[1:8],
							},
							Tokens: tk[1:8],
						}),
						Tokens: tk[1:8],
					},
				},
				Tokens: tk[:9],
			}
		}},
	}, func(t *test) (Type, error) {
		var a Arguments
		err := a.parse(&t.Tokens, t.Yield, t.Await)
		return a, err
	})
}

func TestCallExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidCallExpression,
				Parsing: "CallExpression",
				Token:   tk[0],
			}
		}},
		{"super", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "Arguments",
					Token:   tk[1],
				},
				Parsing: "CallExpression",
				Token:   tk[1],
			}
		}},
		{"super\n()", func(t *test, tk Tokens) { // 3
			t.Output = CallExpression{
				SuperCall: true,
				Arguments: &Arguments{
					Tokens: tk[2:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"import", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "CallExpression",
				Token:   tk[1],
			}
		}},
		{"import\n(\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "CallExpression",
				Token:   tk[4],
			}
		}},
		{"import\n(\na\n)", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 4)
			t.Output = CallExpression{
				ImportCall: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"import\n(\na\n,\nb\n)", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "CallExpression",
				Token:   tk[6],
			}
		}},
		{"super\n()\n`${1 1}`", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidTemplate,
					Parsing: "TemplateLiteral",
					Token:   tk[8],
				},
				Parsing: "CallExpression",
				Token:   tk[5],
			}
		}},
		{"super\n()\n``", func(t *test, tk Tokens) { // 9
			t.Output = CallExpression{
				CallExpression: &CallExpression{
					SuperCall: true,
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[5],
					Tokens:                 tk[5:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"super\n()\n(,)", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[6]),
					Parsing: "Arguments",
					Token:   tk[6],
				},
				Parsing: "CallExpression",
				Token:   tk[5],
			}
		}},
		{"super\n()\n()", func(t *test, tk Tokens) { // 11
			t.Output = CallExpression{
				CallExpression: &CallExpression{
					SuperCall: true,
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				Arguments: &Arguments{
					Tokens: tk[5:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"super\n()\n.\n", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "CallExpression",
				Token:   tk[7],
			}
		}},
		{"super\n()\n.\na", func(t *test, tk Tokens) { // 13
			t.Output = CallExpression{
				CallExpression: &CallExpression{
					SuperCall: true,
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				IdentifierName: &tk[7],
				Tokens:         tk[:8],
			}
		}},
		{"super\n()\n[\n]", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[7]),
					Parsing: "Expression",
					Token:   tk[7],
				},
				Parsing: "CallExpression",
				Token:   tk[5],
			}
		}},
		{"super\n()\n[\na\n]", func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 7)
			t.Output = CallExpression{
				CallExpression: &CallExpression{
					SuperCall: true,
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[7:8],
						},
					},
					Tokens: tk[7:8],
				},
				Tokens: tk[:10],
			}
		}},
		{"super\n()\n[\na\nb]", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "CallExpression",
				Token:   tk[9],
			}
		}},
		{"super\n()\n``\n()\n.\na\n[\nb\n]", func(t *test, tk Tokens) { // 17
			litB := makeConditionLiteral(tk, 16)
			t.Output = CallExpression{
				CallExpression: &CallExpression{
					CallExpression: &CallExpression{
						CallExpression: &CallExpression{
							CallExpression: &CallExpression{
								SuperCall: true,
								Arguments: &Arguments{
									Tokens: tk[2:4],
								},
								Tokens: tk[:4],
							},
							TemplateLiteral: &TemplateLiteral{
								NoSubstitutionTemplate: &tk[5],
								Tokens:                 tk[5:6],
							},
							Tokens: tk[:6],
						},
						Arguments: &Arguments{
							Tokens: tk[7:9],
						},
						Tokens: tk[:9],
					},
					IdentifierName: &tk[12],
					Tokens:         tk[:13],
				},
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[16:17],
						},
					},
					Tokens: tk[16:17],
				},
				Tokens: tk[:19],
			}
		}},
	}, func(t *test) (Type, error) {
		var ce CallExpression
		err := ce.parse(&t.Tokens, nil, t.Yield, t.Await)
		return ce, err
	})
}

func TestOptionalChain(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOptional,
				Parsing: "OptionalChain",
				Token:   tk[0],
			}
		}},
		{"?.", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrInvalidOptionalChain,
				Parsing: "OptionalChain",
				Token:   tk[1],
			}
		}},
		{"?.\n", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrInvalidOptionalChain,
				Parsing: "OptionalChain",
				Token:   tk[2],
			}
		}},
		{"?.\n()", func(t *test, tk Tokens) { // 4
			t.Output = OptionalChain{
				Arguments: &Arguments{
					Tokens: tk[2:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"?.\n[\na\n]", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 4)
			t.Output = OptionalChain{
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"?.\na", func(t *test, tk Tokens) { // 6
			t.Output = OptionalChain{
				IdentifierName: &tk[2],
				Tokens:         tk[:3],
			}
		}},
		{"?.\n``", func(t *test, tk Tokens) { // 7
			t.Output = OptionalChain{
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[2],
					Tokens:                 tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"?.\n()\n``", func(t *test, tk Tokens) { // 8
			t.Output = OptionalChain{
				OptionalChain: &OptionalChain{
					Arguments: &Arguments{
						Tokens: tk[2:4],
					},
					Tokens: tk[:4],
				},
				TemplateLiteral: &TemplateLiteral{
					NoSubstitutionTemplate: &tk[5],
					Tokens:                 tk[5:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"?.\n``\n[\na\n]", func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 6)
			t.Output = OptionalChain{
				OptionalChain: &OptionalChain{
					TemplateLiteral: &TemplateLiteral{
						NoSubstitutionTemplate: &tk[2],
						Tokens:                 tk[2:3],
					},
					Tokens: tk[:3],
				},
				Expression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"?.\n[\na\n]\n.\nb", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 4)
			t.Output = OptionalChain{
				OptionalChain: &OptionalChain{
					Expression: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[4:5],
							},
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[:7],
				},
				IdentifierName: &tk[10],
				Tokens:         tk[:11],
			}
		}},
		{"?.\na\n()", func(t *test, tk Tokens) { // 11
			t.Output = OptionalChain{
				OptionalChain: &OptionalChain{
					IdentifierName: &tk[2],
					Tokens:         tk[:3],
				},
				Arguments: &Arguments{
					Tokens: tk[4:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"?.(.)", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[2]),
					Parsing: "Arguments",
					Token:   tk[2],
				},
				Parsing: "OptionalChain",
				Token:   tk[1],
			}
		}},
		{"?.[]", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[2]),
					Parsing: "Expression",
					Token:   tk[2],
				},
				Parsing: "OptionalChain",
				Token:   tk[2],
			}
		}},
		{"?.[1 1]", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "OptionalChain",
				Token:   tk[4],
			}
		}},
		{"?.`${}`", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[2]),
						Parsing: "Expression",
						Token:   tk[2],
					},
					Parsing: "TemplateLiteral",
					Token:   tk[2],
				},
				Parsing: "OptionalChain",
				Token:   tk[1],
			}
		}},
		{"?.a(.)", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[3]),
					Parsing: "Arguments",
					Token:   tk[3],
				},
				Parsing: "OptionalChain",
				Token:   tk[2],
			}
		}},
		{"?.a[]", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[3]),
					Parsing: "Expression",
					Token:   tk[3],
				},
				Parsing: "OptionalChain",
				Token:   tk[3],
			}
		}},
		{"?.a[1 1]", func(t *test, tk Tokens) { // 18
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "OptionalChain",
				Token:   tk[5],
			}
		}},
		{"?.a.", func(t *test, tk Tokens) { // 19
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "OptionalChain",
				Token:   tk[3],
			}
		}},
		{"?.a`${}`", func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[3]),
						Parsing: "Expression",
						Token:   tk[3],
					},
					Parsing: "TemplateLiteral",
					Token:   tk[3],
				},
				Parsing: "OptionalChain",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var oc OptionalChain
		err := oc.parse(&t.Tokens, t.Yield, t.Await)
		return oc, err
	})
}

func TestOptionalExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOptional,
					Parsing: "OptionalChain",
					Token:   tk[0],
				},
				Parsing: "OptionalExpression",
				Token:   tk[0],
			}
		}},
		{"?.", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidOptionalChain,
					Parsing: "OptionalChain",
					Token:   tk[1],
				},
				Parsing: "OptionalExpression",
				Token:   tk[0],
			}
		}},
		{"?.\na", func(t *test, tk Tokens) { // 3
			t.Output = OptionalExpression{
				OptionalChain: OptionalChain{
					IdentifierName: &tk[2],
					Tokens:         tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"?.\na\n?.", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidOptionalChain,
					Parsing: "OptionalChain",
					Token:   tk[5],
				},
				Parsing: "OptionalExpression",
				Token:   tk[4],
			}
		}},
		{"?.\na\n?.\nb", func(t *test, tk Tokens) { // 5
			t.Output = OptionalExpression{
				OptionalExpression: &OptionalExpression{
					OptionalChain: OptionalChain{
						IdentifierName: &tk[2],
						Tokens:         tk[:3],
					},
					Tokens: tk[:3],
				},
				OptionalChain: OptionalChain{
					IdentifierName: &tk[6],
					Tokens:         tk[4:7],
				},
				Tokens: tk[:7],
			}
		}},
	}, func(t *test) (Type, error) {
		var oe OptionalExpression
		err := oe.parse(&t.Tokens, t.Yield, t.Await, nil, nil)
		return oe, err
	})
}
