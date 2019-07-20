package javascript

import "testing"

func TestLeftHandSideExpressionOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`this`, func(t *test, tk Tokens) { // 1
			t.Output = LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							This:   true,
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
										This:   true,
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
						MetaProperty: true,
						Tokens:       tk[:3],
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
	}, func(t *test) (interface{}, error) {
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
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
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
			litA := makeConditionLiteral(tk, 1)
			litB := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
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
			litA := makeConditionLiteral(tk, 1)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 9)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
							},
							{
								ConditionalExpression: &litB,
								Tokens:                tk[4:5],
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
			litA := makeConditionLiteral(tk, 1)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 7)
			litD := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
							},
							{
								ConditionalExpression: &litB,
								Tokens:                tk[4:5],
							},
							{
								ConditionalExpression: &litC,
								Tokens:                tk[7:8],
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
			litA := makeConditionLiteral(tk, 1)
			litC := makeConditionLiteral(tk, 10)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
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
			litA := makeConditionLiteral(tk, 1)
			litC := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
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
			litA := makeConditionLiteral(tk, 1)
			litC := makeConditionLiteral(tk, 12)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[1:2],
							},
						},
						ObjectBindingPattern: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
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
	}, func(t *test) (interface{}, error) {
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
		{"yield\n,", func(t *test, tk Tokens) { // 4
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "AssignmentExpression",
				Token:   tk[2],
			}
		}},
		{"yield\n1", func(t *test, tk Tokens) { // 5
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
		{"yield\n*\n", func(t *test, tk Tokens) { // 6
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "AssignmentExpression",
				Token:   tk[4],
			}
		}},
		{"yield\n*\n*", func(t *test, tk Tokens) { // 7
			t.Yield = true
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "AssignmentExpression",
				Token:   tk[4],
			}
		}},
		{"yield\n*\n1", func(t *test, tk Tokens) { // 8
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
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "ArrowFunction",
					Token:   tk[1],
				},
				Parsing: "AssignmentExpression",
				Token:   tk[0],
			}
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
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
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
					CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
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
	}, func(t *test) (interface{}, error) {
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
		{`super`, func(t *test, tk Tokens) { // 6
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
		{`import`, func(t *test, tk Tokens) { // 7
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
		{",", func(t *test, tk Tokens) { // 8
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
		{"a", func(t *test, tk Tokens) { // 9
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
		{"a\n()", func(t *test, tk Tokens) { // 10
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
		{"a\n(,)", func(t *test, tk Tokens) { // 11
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
	}, func(t *test) (interface{}, error) {
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
	}, func(t *test) (interface{}, error) {
		var e Expression
		err := e.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return e, err
	})
}

func TestNewExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) {
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
		{",", func(t *test, tk Tokens) {
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
		{"new\n,", func(t *test, tk Tokens) {
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
				Parsing: "NewExpression",
				Token:   tk[2],
			}
		}},
		{"1", func(t *test, tk Tokens) {
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
		{"new\na", func(t *test, tk Tokens) {
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
		{"new\na\n()", func(t *test, tk Tokens) {
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
		{"new\nnew\na\n()", func(t *test, tk Tokens) {
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
		{"new new\nnew\na\n()", func(t *test, tk Tokens) {
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
	}, func(t *test) (interface{}, error) {
		var ne NewExpression
		err := ne.parse(&t.Tokens, t.Yield, t.Await)
		return ne, err
	})
}
