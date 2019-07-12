package javascript

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/parser"
)

type sourceFn struct {
	Source string
	Fn     func(*test, Tokens)
}

type test struct {
	Tokens                     jsParser
	Yield, Await, In, Def, Ret bool
	Output                     interface{}
	Err                        error
}

func doTests(t *testing.T, tests []sourceFn, fn func(*test) (interface{}, error)) {
	t.Helper()
	var err error
	for n, tt := range tests {
		var ts test
		ts.Tokens, err = newJSParser(parser.NewStringTokeniser(tt.Source))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		tt.Fn(&ts, Tokens(ts.Tokens[:cap(ts.Tokens)]))
		output, err := fn(&ts)
		if !reflect.DeepEqual(err, ts.Err) {
			t.Errorf("test %d: expecting error: %v, got %v", n+1, ts.Err, err)
		} else if ts.Output != nil && !reflect.DeepEqual(output, ts.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, ts.Output, output)
		}
	}
}

func TestIdentifier(t *testing.T) {
	doTests(t, []sourceFn{
		{`hello_world`, func(t *test, tk Tokens) {
			t.Output = &tk[0]
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Output = &tk[0]
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Await = true
			t.Output = &tk[0]
		}},
		{`yield`, func(t *test, tk Tokens) {
			t.Yield = true
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "Identifier",
				Token:   tk[0],
			}
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Output = &tk[0]
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Yield = true
			t.Output = &tk[0]
		}},
		{`await`, func(t *test, tk Tokens) {
			t.Await = true
			t.Err = Error{
				Err:     ErrMissingIdentifier,
				Parsing: "Identifier",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseIdentifier(t.Yield, t.Await)
	})
}

func TestScript(t *testing.T) {
	doTests(t, []sourceFn{
		{`
"use strict";

document.body.innerHTML  =
	"Hello, World";

function	runMe	(v) 
{
	let x = v * 2;
	alert(x);
}

for(
	var a = 1;
	a < 10;
	a ++
) {
	runMe ( a );
}
`, func(t *test, tk Tokens) {
			useStrict := makeConditionLiteral(tk, 1)
			helloWorld := makeConditionLiteral(tk, 13)
			v := makeConditionLiteral(tk, 34)
			two := makeConditionLiteral(tk, 38)
			multiply := wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &v.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeMultiply,
				ExponentiationExpression: two.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[34:39],
			})
			x := makeConditionLiteral(tk, 44)
			alert := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					CallExpression: &CallExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[42],
								Tokens:              tk[42:43],
							},
							Tokens: tk[42:43],
						},
						Arguments: &Arguments{
							ArgumentList: []AssignmentExpression{
								{
									ConditionalExpression: &x,
									Tokens:                tk[44:45],
								},
							},
							Tokens: tk[43:46],
						},
						Tokens: tk[42:46],
					},
					Tokens: tk[42:46],
				},
				Tokens: tk[42:46],
			})
			one := makeConditionLiteral(tk, 60)
			a := makeConditionLiteral(tk, 64)
			ten := makeConditionLiteral(tk, 68)
			aLessThanTen := wrapConditional(RelationalExpression{
				RelationalExpression: &a.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipLessThan,
				ShiftExpression:      ten.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[64:69],
			})
			aPlusPlus := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[72],
								Tokens:              tk[72:73],
							},
							Tokens: tk[72:73],
						},
						Tokens: tk[72:73],
					},
					Tokens: tk[72:73],
				},
				UpdateOperator: UpdatePostIncrement,
				Tokens:         tk[72:75],
			})
			argA := makeConditionLiteral(tk, 85)
			runMe := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					CallExpression: &CallExpression{
						MemberExpression: &MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[81],
								Tokens:              tk[81:82],
							},
							Tokens: tk[81:82],
						},
						Arguments: &Arguments{
							ArgumentList: []AssignmentExpression{
								{
									ConditionalExpression: &argA,
									Tokens:                tk[85:86],
								},
							},
							Tokens: tk[83:88],
						},
						Tokens: tk[81:88],
					},
					Tokens: tk[81:88],
				},
				Tokens: tk[81:88],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &useStrict,
										Tokens:                tk[1:2],
									},
								},
								Tokens: tk[1:2],
							},
							Tokens: tk[1:3],
						},
						Tokens: tk[1:3],
					},
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										LeftHandSideExpression: &LeftHandSideExpression{
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
														IdentifierName: &tk[6],
														Tokens:         tk[4:7],
													},
													IdentifierName: &tk[8],
													Tokens:         tk[4:9],
												},
												Tokens: tk[4:9],
											},
											Tokens: tk[4:9],
										},
										AssignmentOperator: AssignmentAssign,
										AssignmentExpression: &AssignmentExpression{
											ConditionalExpression: &helloWorld,
											Tokens:                tk[13:14],
										},
										Tokens: tk[4:14],
									},
								},
								Tokens: tk[4:14],
							},
							Tokens: tk[4:15],
						},
						Tokens: tk[4:15],
					},
					{
						Declaration: &Declaration{
							FunctionDeclaration: &FunctionDeclaration{
								BindingIdentifier: &tk[18],
								FormalParameters: FormalParameters{
									FormalParameterList: []BindingElement{
										{
											SingleNameBinding: &tk[21],
											Tokens:            tk[21:22],
										},
									},
									Tokens: tk[21:22],
								},
								FunctionBody: Block{
									StatementListItems: []StatementListItem{
										{
											Declaration: &Declaration{
												LexicalDeclaration: &LexicalDeclaration{
													LetOrConst: Let,
													BindingList: []LexicalBinding{
														{
															BindingIdentifier: &tk[30],
															Initializer: &AssignmentExpression{
																ConditionalExpression: &multiply,
																Tokens:                tk[34:39],
															},
															Tokens: tk[30:39],
														},
													},
													Tokens: tk[28:40],
												},
												Tokens: tk[28:40],
											},
											Tokens: tk[28:40],
										},
										{
											Statement: &Statement{
												ExpressionStatement: &Expression{
													Expressions: []AssignmentExpression{
														{
															ConditionalExpression: &alert,
															Tokens:                tk[42:46],
														},
													},
													Tokens: tk[42:46],
												},
												Tokens: tk[42:47],
											},
											Tokens: tk[42:47],
										},
									},
									Tokens: tk[25:49],
								},
								Tokens: tk[16:49],
							},
							Tokens: tk[16:49],
						},
						Tokens: tk[16:49],
					},
					{
						Statement: &Statement{
							IterationStatementFor: &IterationStatementFor{
								Type: ForNormalVar,
								InitVar: []VariableDeclaration{
									{
										BindingIdentifier: &tk[56],
										Initializer: &AssignmentExpression{
											ConditionalExpression: &one,
											Tokens:                tk[60:61],
										},
										Tokens: tk[56:61],
									},
								},
								Conditional: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &aLessThanTen,
											Tokens:                tk[64:69],
										},
									},
									Tokens: tk[64:69],
								},
								Afterthought: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &aPlusPlus,
											Tokens:                tk[72:75],
										},
									},
									Tokens: tk[72:75],
								},
								Statement: Statement{
									BlockStatement: &Block{
										StatementListItems: []StatementListItem{
											{
												Statement: &Statement{
													ExpressionStatement: &Expression{
														Expressions: []AssignmentExpression{
															{
																ConditionalExpression: &runMe,
																Tokens:                tk[81:88],
															},
														},
														Tokens: tk[81:88],
													},
													Tokens: tk[81:89],
												},
												Tokens: tk[81:89],
											},
										},
										Tokens: tk[78:91],
									},
									Tokens: tk[78:91],
								},
								Tokens: tk[50:91],
							},
							Tokens: tk[50:91],
						},
						Tokens: tk[50:91],
					},
				},
				Tokens: tk[:92],
			}
		}},
		{`if (typeof a === "b" && typeof c.d == "e") {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 5)
			litB := makeConditionLiteral(tk, 9)
			litE := makeConditionLiteral(tk, 21)
			typeOfA := wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryTypeOf},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[3:6],
			})
			CD := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							MemberExpression: &MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[15],
									Tokens:              tk[15:16],
								},
								Tokens: tk[15:16],
							},
							IdentifierName: &tk[17],
							Tokens:         tk[15:18],
						},
						Tokens: tk[15:18],
					},
					Tokens: tk[15:18],
				},
				Tokens: tk[15:18],
			})
			typeOfCD := wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryTypeOf},
				UpdateExpression: CD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[13:18],
			})
			AEquals := wrapConditional(EqualityExpression{
				EqualityExpression:   &typeOfA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityStrictEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[3:10],
			})
			CDEquals := wrapConditional(EqualityExpression{
				EqualityExpression:   &typeOfCD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityEqual,
				RelationalExpression: litE.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[13:22],
			})
			And := wrapConditional(LogicalANDExpression{
				LogicalANDExpression: &AEquals.LogicalORExpression.LogicalANDExpression,
				BitwiseORExpression:  CDEquals.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
				Tokens:               tk[3:22],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							IfStatement: &IfStatement{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &And,
											Tokens:                tk[3:22],
										},
									},
									Tokens: tk[3:22],
								},
								Statement: Statement{
									BlockStatement: &Block{
										Tokens: tk[24:26],
									},
									Tokens: tk[24:26],
								},
								Tokens: tk[:26],
							},
							Tokens: tk[:26],
						},
						Tokens: tk[:26],
					},
				},
				Tokens: tk[:26],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var s Script
		err := s.parse(&t.Tokens)
		return s, err
	})
}
