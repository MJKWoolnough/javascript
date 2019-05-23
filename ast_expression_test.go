package javascript

import "testing"

func TestNewExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{`this`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						This:   true,
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`someIdentifier`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &IdentifierReference{Identifier: &tk[0]},
						Tokens:              tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`new someIdentifier`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				News: 1,
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &IdentifierReference{Identifier: &tk[2]},
						Tokens:              tk[2:3],
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`new new someIdentifier`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				News: 2,
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: &IdentifierReference{Identifier: &tk[4]},
						Tokens:              tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`null`, func(t *test, tk Tokens) {
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
		{`true`, func(t *test, tk Tokens) {
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
		{`false`, func(t *test, tk Tokens) {
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
		{`0`, func(t *test, tk Tokens) {
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
		{`"Hello"`, func(t *test, tk Tokens) {
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
		{`[]`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
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
			}
		}},
		{`{}`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
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
			}
		}},
		{`super.runMe`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					SuperProperty:  true,
					IdentifierName: &tk[2],
					Tokens:         tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`super[runMe]`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 2)
			t.Output = NewExpression{
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
			}
		}},
		{`this.key.field.next`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
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
			}
		}},
		{`new.target`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					MetaProperty: true,
					Tokens:       tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`new className()`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &IdentifierReference{Identifier: &tk[2]},
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
			}
		}},
		{`new new className()`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				News: 1,
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &IdentifierReference{Identifier: &tk[4]},
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
			}
		}},
		{`new new className()()`, func(t *test, tk Tokens) {
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					MemberExpression: &MemberExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &IdentifierReference{Identifier: &tk[4]},
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
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseNewExpression(t.Yield, t.Await)
	})
}

func TestAssignmentExpression(t *testing.T) {
	doTests(t, []sourceFn{
		{`yield 1`, func(t *test, tk Tokens) {}},
		{`yield 1`, func(t *test, tk Tokens) {
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
		{`yield *1`, func(t *test, tk Tokens) {
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
		{`a => {}`, func(t *test, tk Tokens) {
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &BindingIdentifier{Identifier: &tk[0]},
					FunctionBody: &Block{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`a => 1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 4)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					BindingIdentifier: &BindingIdentifier{Identifier: &tk[0]},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[4:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`async a => 1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 6)
			t.Output = AssignmentExpression{
				ArrowFunction: &ArrowFunction{
					Async:             true,
					BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[6:7],
					},
					Tokens: tk[:7],
				},
				Tokens: tk[:7],
			}
		}},
		{`() => {}`, func(t *test, tk Tokens) {
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
		{`(a) => b`, func(t *test, tk Tokens) {
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
		{`(a, b) => c`, func(t *test, tk Tokens) {
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
		{`(a, b, c) => d`, func(t *test, tk Tokens) {
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
		{`(a, ...b) => c`, func(t *test, tk Tokens) {
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
						BindingIdentifier: &BindingIdentifier{Identifier: &tk[5]},
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
		{`(a, ...[b]) => c`, func(t *test, tk Tokens) {
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
									SingleNameBinding: &BindingIdentifier{Identifier: &tk[6]},
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
		{`(a, ...{b}) => c`, func(t *test, tk Tokens) {
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
									SingleNameBinding: &BindingIdentifier{Identifier: &tk[6]},
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
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseAssignmentExpression(t.In, t.Yield, t.Await)
	})
}
