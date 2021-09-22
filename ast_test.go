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

func doTests(t *testing.T, tests []sourceFn, fn func(*test) (Type, error)) {
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

func assignmentError(tk Token) Error {
	return assignmentCustomError(tk, ErrNoIdentifier)
}

func assignmentCustomError(tk Token, err error) Error {
	return Error{
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
																				Err:     err,
																				Parsing: "PrimaryExpression",
																				Token:   tk,
																			},
																			Parsing: "MemberExpression",
																			Token:   tk,
																		},
																		Parsing: "NewExpression",
																		Token:   tk,
																	},
																	Parsing: "LeftHandSideExpression",
																	Token:   tk,
																},
																Parsing: "UpdateExpression",
																Token:   tk,
															},
															Parsing: "UnaryExpression",
															Token:   tk,
														},
														Parsing: "ExponentiationExpression",
														Token:   tk,
													},
													Parsing: "MultiplicativeExpression",
													Token:   tk,
												},
												Parsing: "AdditiveExpression",
												Token:   tk,
											},
											Parsing: "ShiftExpression",
											Token:   tk,
										},
										Parsing: "RelationalExpression",
										Token:   tk,
									},
									Parsing: "EqualityExpression",
									Token:   tk,
								},
								Parsing: "BitwiseANDExpression",
								Token:   tk,
							},
							Parsing: "BitwiseXORExpression",
							Token:   tk,
						},
						Parsing: "BitwiseORExpression",
						Token:   tk,
					},
					Parsing: "LogicalANDExpression",
					Token:   tk,
				},
				Parsing: "LogicalORExpression",
				Token:   tk,
			},
			Parsing: "ConditionalExpression",
			Token:   tk,
		},
		Parsing: "AssignmentExpression",
		Token:   tk,
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
			t.Output = nil
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
			t.Output = nil
		}},
		{`for`, func(t *test, tk Tokens) {
			t.Output = nil
		}},
		{`"for"`, func(t *test, tk Tokens) {
			t.Output = nil
		}},
		{`1`, func(t *test, tk Tokens) {
			t.Output = nil
		}},
		{`+`, func(t *test, tk Tokens) {
			t.Output = nil
		}},
	}, func(t *test) (Type, error) {
		return t.Tokens.parseIdentifier(t.Yield, t.Await), nil
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
`, func(t *test, tk Tokens) { // 1
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
									Tokens: tk[20:23],
								},
								FunctionBody: Block{
									StatementList: []StatementListItem{
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
										StatementList: []StatementListItem{
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
		{`if (typeof a === "b" && typeof c.d == "e") {}`, func(t *test, tk Tokens) { // 2
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
		{"{ 1 2 } 3", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingSemiColon,
								Parsing: "Statement",
								Token:   tk[3],
							},
							Parsing: "StatementListItem",
							Token:   tk[2],
						},
						Parsing: "Block",
						Token:   tk[2],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"{ 1\n2 } 3", func(t *test, tk Tokens) { // 4
			lit1 := makeConditionLiteral(tk, 2)
			lit2 := makeConditionLiteral(tk, 4)
			lit3 := makeConditionLiteral(tk, 8)
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							BlockStatement: &Block{
								StatementList: []StatementListItem{
									{
										Statement: &Statement{
											ExpressionStatement: &Expression{
												Expressions: []AssignmentExpression{
													{
														ConditionalExpression: &lit1,
														Tokens:                tk[2:3],
													},
												},
												Tokens: tk[2:3],
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:3],
									},
									{
										Statement: &Statement{
											ExpressionStatement: &Expression{
												Expressions: []AssignmentExpression{
													{
														ConditionalExpression: &lit2,
														Tokens:                tk[4:5],
													},
												},
												Tokens: tk[4:5],
											},
											Tokens: tk[4:5],
										},
										Tokens: tk[4:5],
									},
								},
								Tokens: tk[:7],
							},
							Tokens: tk[:7],
						},
						Tokens: tk[:7],
					},
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &lit3,
										Tokens:                tk[8:9],
									},
								},
								Tokens: tk[8:9],
							},
							Tokens: tk[8:9],
						},
						Tokens: tk[8:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"function fn() {\nreturn\na + b\n}", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 10)
			litB := makeConditionLiteral(tk, 14)
			add := wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[10:15],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Declaration: &Declaration{
							FunctionDeclaration: &FunctionDeclaration{
								BindingIdentifier: &tk[2],
								FormalParameters: FormalParameters{
									Tokens: tk[3:5],
								},
								FunctionBody: Block{
									StatementList: []StatementListItem{
										{
											Statement: &Statement{
												Type:   StatementReturn,
												Tokens: tk[8:9],
											},
											Tokens: tk[8:9],
										},
										{
											Statement: &Statement{
												ExpressionStatement: &Expression{
													Expressions: []AssignmentExpression{
														{
															ConditionalExpression: &add,
															Tokens:                tk[10:15],
														},
													},
													Tokens: tk[10:15],
												},
												Tokens: tk[10:15],
											},
											Tokens: tk[10:15],
										},
									},
									Tokens: tk[6:17],
								},
								Tokens: tk[:17],
							},
							Tokens: tk[:17],
						},
						Tokens: tk[:17],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"a = b\n++c", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 7)
			pa := wrapConditional(UpdateExpression{
				UpdateOperator:  UpdatePreIncrement,
				UnaryExpression: &litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[6:8],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
										AssignmentOperator:     AssignmentAssign,
										AssignmentExpression: &AssignmentExpression{
											ConditionalExpression: &litB,
											Tokens:                tk[4:5],
										},
										Tokens: tk[:5],
									},
								},
								Tokens: tk[:5],
							},
							Tokens: tk[:5],
						},
						Tokens: tk[:5],
					},
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &pa,
										Tokens:                tk[6:8],
									},
								},
								Tokens: tk[6:8],
							},
							Tokens: tk[6:8],
						},
						Tokens: tk[6:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{"if (a > b)\nelse d = e", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     assignmentError(tk[10]),
								Parsing: "Expression",
								Token:   tk[10],
							},
							Parsing: "Statement",
							Token:   tk[10],
						},
						Parsing: "IfStatement",
						Token:   tk[10],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"if\n(a\n>\nb)\nc\nelse\nd\n=\ne", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 3)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 10)
			litD := makeConditionLiteral(tk, 14)
			litE := makeConditionLiteral(tk, 18)
			ab := wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipGreaterThan,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[3:8],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							IfStatement: &IfStatement{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &ab,
											Tokens:                tk[3:8],
										},
									},
									Tokens: tk[3:8],
								},
								Statement: Statement{
									ExpressionStatement: &Expression{
										Expressions: []AssignmentExpression{
											{
												ConditionalExpression: &litC,
												Tokens:                tk[10:11],
											},
										},
										Tokens: tk[10:11],
									},
									Tokens: tk[10:11],
								},
								ElseStatement: &Statement{
									ExpressionStatement: &Expression{
										Expressions: []AssignmentExpression{
											{
												LeftHandSideExpression: litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
												AssignmentOperator:     AssignmentAssign,
												AssignmentExpression: &AssignmentExpression{
													ConditionalExpression: &litE,
													Tokens:                tk[18:19],
												},
												Tokens: tk[14:19],
											},
										},
										Tokens: tk[14:19],
									},
									Tokens: tk[14:19],
								},
								Tokens: tk[:19],
							},
							Tokens: tk[:19],
						},
						Tokens: tk[:19],
					},
				},
				Tokens: tk[:19],
			}
		}},
		{"a = b + c\n(d + e).print()", func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litD := makeConditionLiteral(tk, 11)
			litE := makeConditionLiteral(tk, 15)
			de := wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litE.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[11:16],
			})
			bc := wrapConditional(AdditiveExpression{
				AdditiveExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:   AdditiveAdd,
				MultiplicativeExpression: wrapConditional(UpdateExpression{
					LeftHandSideExpression: &LeftHandSideExpression{
						CallExpression: &CallExpression{
							CallExpression: &CallExpression{
								CallExpression: &CallExpression{
									MemberExpression: &MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: &tk[8],
											Tokens:              tk[8:9],
										},
										Tokens: tk[8:9],
									},
									Arguments: &Arguments{
										ArgumentList: []AssignmentExpression{
											{
												ConditionalExpression: &de,
												Tokens:                tk[11:16],
											},
										},
										Tokens: tk[10:17],
									},
									Tokens: tk[8:17],
								},
								IdentifierName: &tk[18],
								Tokens:         tk[8:19],
							},
							Arguments: &Arguments{
								Tokens: tk[19:21],
							},
							Tokens: tk[8:21],
						},
						Tokens: tk[8:21],
					},
					Tokens: tk[8:21],
				}).LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens: tk[4:21],
			})
			t.Output = Script{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
										AssignmentOperator:     AssignmentAssign,
										AssignmentExpression: &AssignmentExpression{
											ConditionalExpression: &bc,
											Tokens:                tk[4:21],
										},
										Tokens: tk[:21],
									},
								},
								Tokens: tk[:21],
							},
							Tokens: tk[:21],
						},
						Tokens: tk[:21],
					},
				},
				Tokens: tk[:21],
			}
		}},
		{"await a", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingSemiColon,
					Parsing: "Statement",
					Token:   tk[1],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var s Script
		err := s.parse(&t.Tokens)
		return s, err
	})
}

func TestDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{`class`, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "ClassDeclaration",
					Token:   tk[1],
				},
				Parsing: "Declaration",
				Token:   tk[0],
			}
		}},
		{`class a{}`, func(t *test, tk Tokens) { // 2
			t.Output = Declaration{
				ClassDeclaration: &ClassDeclaration{
					BindingIdentifier: &tk[2],
					Tokens:            tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`const`, func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[1],
					},
					Parsing: "LexicalDeclaration",
					Token:   tk[1],
				},
				Parsing: "Declaration",
				Token:   tk[0],
			}
		}},
		{`const a = 1;`, func(t *test, tk Tokens) { // 4
			lit1 := makeConditionLiteral(tk, 6)
			t.Output = Declaration{
				LexicalDeclaration: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[2],
							Initializer: &AssignmentExpression{
								ConditionalExpression: &lit1,
								Tokens:                tk[6:7],
							},
							Tokens: tk[2:7],
						},
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`let`, func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[1],
					},
					Parsing: "LexicalDeclaration",
					Token:   tk[1],
				},
				Parsing: "Declaration",
				Token:   tk[0],
			}
		}},
		{`let a = 1;`, func(t *test, tk Tokens) { // 4
			lit1 := makeConditionLiteral(tk, 6)
			t.Output = Declaration{
				LexicalDeclaration: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[2],
							Initializer: &AssignmentExpression{
								ConditionalExpression: &lit1,
								Tokens:                tk[6:7],
							},
							Tokens: tk[2:7],
						},
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`function`, func(t *test, tk Tokens) { // 7 {
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "FunctionDeclaration",
					Token:   tk[1],
				},
				Parsing: "Declaration",
				Token:   tk[0],
			}
		}},
		{`function a(){}`, func(t *test, tk Tokens) { // 8
			t.Output = Declaration{
				FunctionDeclaration: &FunctionDeclaration{
					BindingIdentifier: &tk[2],
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
		{`wrong`, func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrInvalidDeclaration,
				Parsing: "Declaration",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var d Declaration
		err := d.parse(&t.Tokens, t.Yield, t.Await)
		return d, err
	})
}

func TestLexicalDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{`wrong`, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidLexicalDeclaration,
				Parsing: "LexicalDeclaration",
				Token:   tk[0],
			}
		}},
		{`const`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "LexicalBinding",
					Token:   tk[1],
				},
				Parsing: "LexicalDeclaration",
				Token:   tk[1],
			}
		}},
		{"const\na", func(t *test, tk Tokens) { // 3
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"const\na\n", func(t *test, tk Tokens) { // 4
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"const\na;", func(t *test, tk Tokens) { // 5
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:4],
			}
		}},
		{"const\na,\nb;", func(t *test, tk Tokens) { // 6
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						BindingIdentifier: &tk[5],
						Tokens:            tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"const \n a;", func(t *test, tk Tokens) { // 7
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[4],
						Tokens:            tk[4:5],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{"const \n a, \n b;", func(t *test, tk Tokens) { // 8
			t.Output = LexicalDeclaration{
				LetOrConst: Const,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[4],
						Tokens:            tk[4:5],
					},
					{
						BindingIdentifier: &tk[9],
						Tokens:            tk[9:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`let`, func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "LexicalBinding",
					Token:   tk[1],
				},
				Parsing: "LexicalDeclaration",
				Token:   tk[1],
			}
		}},
		{"let\na", func(t *test, tk Tokens) { // 10
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"let\na\n", func(t *test, tk Tokens) { // 11
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"let\na;", func(t *test, tk Tokens) { // 12
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:4],
			}
		}},
		{"let\na,\nb;", func(t *test, tk Tokens) { // 13
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						BindingIdentifier: &tk[5],
						Tokens:            tk[5:6],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"let \n a;", func(t *test, tk Tokens) { // 14
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[4],
						Tokens:            tk[4:5],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{"let \n a, \n b;", func(t *test, tk Tokens) { // 15
			t.Output = LexicalDeclaration{
				LetOrConst: Let,
				BindingList: []LexicalBinding{
					{
						BindingIdentifier: &tk[4],
						Tokens:            tk[4:5],
					},
					{
						BindingIdentifier: &tk[9],
						Tokens:            tk[9:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"let\na b;", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err:     ErrInvalidLexicalDeclaration,
				Parsing: "LexicalDeclaration",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) { // 13
		var ld LexicalDeclaration
		err := ld.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return ld, err
	})
}

func TestLexicalBinding(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "LexicalBinding",
				Token:   tk[0],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 2
			t.Output = LexicalBinding{
				BindingIdentifier: &tk[0],
				Tokens:            tk[:1],
			}
		}},
		{"a\n=\n1", func(t *test, tk Tokens) { // 3
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = LexicalBinding{
				BindingIdentifier: &tk[0],
				Initializer: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`[a]`, func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingInitializer,
				Parsing: "LexicalBinding",
				Token:   tk[3],
			}
		}},
		{"[a]\n=\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     assignmentError(tk[6]),
				Parsing: "LexicalBinding",
				Token:   tk[6],
			}
		}},
		{"[]\n=\na", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 5)
			t.Output = LexicalBinding{
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
		{`{a}`, func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingInitializer,
				Parsing: "LexicalBinding",
				Token:   tk[3],
			}
		}},
		{"{a}\n=\n", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     assignmentError(tk[6]),
				Parsing: "LexicalBinding",
				Token:   tk[6],
			}
		}},
		{"{}\n=\na", func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 5)
			t.Output = LexicalBinding{
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
				Parsing: "LexicalBinding",
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
				Parsing: "LexicalBinding",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var lb LexicalBinding
		err := lb.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return lb, err
	})
}

func TestArrayBindingPattern(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBracket,
				Parsing: "ArrayBindingPattern",
				Token:   tk[0],
			}
		}},
		{`{}`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningBracket,
				Parsing: "ArrayBindingPattern",
				Token:   tk[0],
			}
		}},
		{"[\n]", func(t *test, tk Tokens) { // 3
			t.Output = ArrayBindingPattern{
				Tokens: tk[:3],
			}
		}},
		{"[\nnull\n]", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[2],
				},
				Parsing: "ArrayBindingPattern",
				Token:   tk[2],
			}
		}},
		{"[\n,\n]", func(t *test, tk Tokens) { // 5
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{},
				},
				Tokens: tk[:5],
			}
		}},
		{"[\n,\n]", func(t *test, tk Tokens) { // 6
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{},
				},
				Tokens: tk[:5],
			}
		}},
		{"[\n,\n,\n]", func(t *test, tk Tokens) { // 7
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{},
					{},
				},
				Tokens: tk[:7],
			}
		}},
		{"[\n...\na\n]", func(t *test, tk Tokens) { // 8
			t.Output = ArrayBindingPattern{
				BindingRestElement: &BindingElement{
					SingleNameBinding: &tk[4],
					Tokens:            tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"[\n...\nnull\n]", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[4],
				},
				Parsing: "ArrayBindingPattern",
				Token:   tk[4],
			}
		}},
		{"[\n...\na\n,]", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "ArrayBindingPattern",
				Token:   tk[6],
			}
		}},
		{"[\na\n]", func(t *test, tk Tokens) { // 11
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"[\na\n,\n]", func(t *test, tk Tokens) { // 12
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"[\na\n,\n,\n]", func(t *test, tk Tokens) { // 13
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
					{},
				},
				Tokens: tk[:9],
			}
		}},
		{"[\n,\na\n]", func(t *test, tk Tokens) { // 14
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{},
					{
						SingleNameBinding: &tk[4],
						Tokens:            tk[4:5],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"[\n,\nnull\n]", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[4],
				},
				Parsing: "ArrayBindingPattern",
				Token:   tk[4],
			}
		}},
		{"[\na,\n,\n...\nb\n]", func(t *test, tk Tokens) { // 16
			t.Output = ArrayBindingPattern{
				BindingElementList: []BindingElement{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
					{},
				},
				BindingRestElement: &BindingElement{
					SingleNameBinding: &tk[9],
					Tokens:            tk[9:10],
				},
				Tokens: tk[:12],
			}
		}},
		{"[\na\nb\n]", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "ArrayBindingPattern",
				Token:   tk[4],
			}
		}},
	}, func(t *test) (Type, error) {
		var ab ArrayBindingPattern
		err := ab.parse(&t.Tokens, t.Yield, t.Await)
		return ab, err
	})
}

func TestObjectBindingPattern(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ObjectBindingPattern",
				Token:   tk[0],
			}
		}},
		{`[]`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ObjectBindingPattern",
				Token:   tk[0],
			}
		}},
		{"{\n}", func(t *test, tk Tokens) { // 3
			t.Output = ObjectBindingPattern{
				Tokens: tk[:3],
			}
		}},
		{"{\n...\nnull\n}", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ObjectBindingPattern",
				Token:   tk[4],
			}
		}},
		{"{\n...\na\n}", func(t *test, tk Tokens) { // 5
			t.Output = ObjectBindingPattern{
				BindingRestProperty: &tk[4],
				Tokens:              tk[:7],
			}
		}},
		{"{\n...\na\nb\n}", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingClosingBrace,
				Parsing: "ObjectBindingPattern",
				Token:   tk[6],
			}
		}},
		{"{\nnull\n}", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidPropertyName,
						Parsing: "PropertyName",
						Token:   tk[2],
					},
					Parsing: "BindingProperty",
					Token:   tk[2],
				},
				Parsing: "ObjectBindingPattern",
				Token:   tk[2],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 8
			t.Output = ObjectBindingPattern{
				BindingPropertyList: []BindingProperty{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[2],
							Tokens:              tk[2:3],
						},
						BindingElement: BindingElement{
							SingleNameBinding: &tk[2],
							Tokens:            tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"{\na\nb\n}", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "ObjectBindingPattern",
				Token:   tk[4],
			}
		}},
		{"{\na,\nb\n}", func(t *test, tk Tokens) { // 10
			t.Output = ObjectBindingPattern{
				BindingPropertyList: []BindingProperty{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[2],
							Tokens:              tk[2:3],
						},
						BindingElement: BindingElement{
							SingleNameBinding: &tk[2],
							Tokens:            tk[2:3],
						},
						Tokens: tk[2:3],
					},
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[5],
							Tokens:              tk[5:6],
						},
						BindingElement: BindingElement{
							SingleNameBinding: &tk[5],
							Tokens:            tk[5:6],
						},
						Tokens: tk[5:6],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{"{\na,\n...\nb\n}", func(t *test, tk Tokens) { // 11
			t.Output = ObjectBindingPattern{
				BindingPropertyList: []BindingProperty{
					{
						PropertyName: PropertyName{
							LiteralPropertyName: &tk[2],
							Tokens:              tk[2:3],
						},
						BindingElement: BindingElement{
							SingleNameBinding: &tk[2],
							Tokens:            tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				BindingRestProperty: &tk[7],
				Tokens:              tk[:10],
			}
		}},
	}, func(t *test) (Type, error) {
		var ob ObjectBindingPattern
		err := ob.parse(&t.Tokens, t.Yield, t.Await)
		return ob, err
	})
}

func TestBindingProperty(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidPropertyName,
					Parsing: "PropertyName",
					Token:   tk[0],
				},
				Parsing: "BindingProperty",
				Token:   tk[0],
			}
		}},
		{`null`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidPropertyName,
					Parsing: "PropertyName",
					Token:   tk[0],
				},
				Parsing: "BindingProperty",
				Token:   tk[0],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 3
			t.Output = BindingProperty{
				PropertyName: PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				BindingElement: BindingElement{
					SingleNameBinding: &tk[0],
					Tokens:            tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"a\n=\n", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "BindingElement",
					Token:   tk[4],
				},
				Parsing: "BindingProperty",
				Token:   tk[0],
			}
		}},
		{"a\n=\n1", func(t *test, tk Tokens) { // 5
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = BindingProperty{
				PropertyName: PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				BindingElement: BindingElement{
					SingleNameBinding: &tk[0],
					Initializer: &AssignmentExpression{
						ConditionalExpression: &lit1,
						Tokens:                tk[4:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a\n:\n", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[4],
				},
				Parsing: "BindingProperty",
				Token:   tk[4],
			}
		}},
		{"'a'", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "BindingProperty",
				Token:   tk[1],
			}
		}},
		{"a\n:\n''", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "BindingElement",
					Token:   tk[4],
				},
				Parsing: "BindingProperty",
				Token:   tk[4],
			}
		}},
		{"a\n:\nb", func(t *test, tk Tokens) { // 9
			t.Output = BindingProperty{
				PropertyName: PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				BindingElement: BindingElement{
					SingleNameBinding: &tk[4],
					Tokens:            tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
		var bp BindingProperty
		err := bp.parse(&t.Tokens, t.Yield, t.Await)
		return bp, err
	})
}

func TestArrayLiteral(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBracket,
				Parsing: "ArrayLiteral",
				Token:   tk[0],
			}
		}},
		{`{}`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningBracket,
				Parsing: "ArrayLiteral",
				Token:   tk[0],
			}
		}},
		{"[\n]", func(t *test, tk Tokens) { // 3
			t.Output = ArrayLiteral{
				Tokens: tk[:3],
			}
		}},
		{"[\n,\n]", func(t *test, tk Tokens) { // 4
			t.Output = ArrayLiteral{
				ElementList: []AssignmentExpression{
					{},
				},
				Tokens: tk[:5],
			}
		}},
		{"[\n,\n,\n]", func(t *test, tk Tokens) { // 5
			t.Output = ArrayLiteral{
				ElementList: []AssignmentExpression{
					{},
					{},
				},
				Tokens: tk[:7],
			}
		}},
		{"[\n...\n]", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "ArrayLiteral",
				Token:   tk[4],
			}
		}},
		{"[\n...\na\n]", func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 4)
			t.Output = ArrayLiteral{
				SpreadElement: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:7],
			}
		}},
		{"[\n...\na\nb]", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingClosingBracket,
				Parsing: "ArrayLiteral",
				Token:   tk[6],
			}
		}},
		{"[\n*\n]", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "ArrayLiteral",
				Token:   tk[2],
			}
		}},
		{"[\na\n]", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 2)
			t.Output = ArrayLiteral{
				ElementList: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"[\na\nb\n]", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "ArrayLiteral",
				Token:   tk[4],
			}
		}},
		{"[\na\n,\nb\n]", func(t *test, tk Tokens) { // 12
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			t.Output = ArrayLiteral{
				ElementList: []AssignmentExpression{
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
		{"[\na\n,\n,\nb\n,\n,\n...\nc\n]", func(t *test, tk Tokens) { // 12
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			litC := makeConditionLiteral(tk, 16)
			t.Output = ArrayLiteral{
				ElementList: []AssignmentExpression{
					{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
					{},
					{
						ConditionalExpression: &litB,
						Tokens:                tk[8:9],
					},
					{},
				},
				SpreadElement: &AssignmentExpression{
					ConditionalExpression: &litC,
					Tokens:                tk[16:17],
				},
				Tokens: tk[:19],
			}
		}},
	}, func(t *test) (Type, error) {
		var al ArrayLiteral
		err := al.parse(&t.Tokens, t.Yield, t.Await)
		return al, err
	})
}

func TestObjectLiteral(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ObjectLiteral",
				Token:   tk[0],
			}
		}},
		{`[]`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "ObjectLiteral",
				Token:   tk[0],
			}
		}},
		{"{\n}", func(t *test, tk Tokens) { // 3
			t.Output = ObjectLiteral{
				Tokens: tk[:3],
			}
		}},
		{"{\n,\n}", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidPropertyName,
						Parsing: "PropertyName",
						Token:   tk[2],
					},
					Parsing: "PropertyDefinition",
					Token:   tk[2],
				},
				Parsing: "ObjectLiteral",
				Token:   tk[2],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 2)
			t.Output = ObjectLiteral{
				PropertyDefinitionList: []PropertyDefinition{
					{
						PropertyName: &PropertyName{
							LiteralPropertyName: &tk[2],
							Tokens:              tk[2:3],
						},
						AssignmentExpression: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"{\n...a\nb\n}", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "ObjectLiteral",
				Token:   tk[5],
			}
		}},
		{"{\na\n,\nb\n}", func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			t.Output = ObjectLiteral{
				PropertyDefinitionList: []PropertyDefinition{
					{
						PropertyName: &PropertyName{
							LiteralPropertyName: &tk[2],
							Tokens:              tk[2:3],
						},
						AssignmentExpression: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
						Tokens: tk[2:3],
					},
					{
						PropertyName: &PropertyName{
							LiteralPropertyName: &tk[6],
							Tokens:              tk[6:7],
						},
						AssignmentExpression: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[6:7],
						},
						Tokens: tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
	}, func(t *test) (Type, error) {
		var ol ObjectLiteral
		err := ol.parse(&t.Tokens, t.Yield, t.Await)
		return ol, err
	})
}

func TestPropertyDefinition(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidPropertyName,
					Parsing: "PropertyName",
					Token:   tk[0],
				},
				Parsing: "PropertyDefinition",
				Token:   tk[0],
			}
		}},
		{`...`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     assignmentError(tk[1]),
				Parsing: "PropertyDefinition",
				Token:   tk[1],
			}
		}},
		{"...\na", func(t *test, tk Tokens) { // 3
			litA := makeConditionLiteral(tk, 2)
			t.Output = PropertyDefinition{
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"a\n,", func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 0)
			t.Output = PropertyDefinition{
				PropertyName: &PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 2)
			t.Output = PropertyDefinition{
				PropertyName: &PropertyName{
					LiteralPropertyName: &tk[2],
					Tokens:              tk[2:3],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litA,
					Tokens:                tk[2:3],
				},
				Tokens: tk[2:3],
			}
			t.Tokens = jsParser(tk[2:2])
		}},
		{"a\n=\n", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "PropertyDefinition",
				Token:   tk[4],
			}
		}},
		{"a\n=\nb", func(t *test, tk Tokens) { // 7
			litB := makeConditionLiteral(tk, 4)
			t.Output = PropertyDefinition{
				IsCoverInitializedName: true,
				PropertyName: &PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a\n:\n", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "PropertyDefinition",
				Token:   tk[4],
			}
		}},
		{"a\n:\nb", func(t *test, tk Tokens) { // 9
			litB := makeConditionLiteral(tk, 4)
			t.Output = PropertyDefinition{
				PropertyName: &PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"[\na\n]\n:\nb", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			t.Output = PropertyDefinition{
				PropertyName: &PropertyName{
					ComputedPropertyName: &AssignmentExpression{
						ConditionalExpression: &litA,
						Tokens:                tk[2:3],
					},
					Tokens: tk[:5],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"[\na\n]\n", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningParenthesis,
						Parsing: "FormalParameters",
						Token:   tk[6],
					},
					Parsing: "MethodDefinition",
					Token:   tk[6],
				},
				Parsing: "PropertyDefinition",
				Token:   tk[0],
			}
		}},
		{"a\n(){}", func(t *test, tk Tokens) { // 12
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[0],
						Tokens:              tk[:1],
					},
					Params: FormalParameters{
						Tokens: tk[2:4],
					},
					FunctionBody: Block{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"static\na\n(){}", func(t *test, tk Tokens) { // 13
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodStatic,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[4:6],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"*\na\n(){}", func(t *test, tk Tokens) { // 14
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodGenerator,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[4:6],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"async a\n(){}", func(t *test, tk Tokens) { // 15
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodAsync,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[4:6],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"get\na\n(){}", func(t *test, tk Tokens) { // 16
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodGetter,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[4:6],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"set\na\n(b){}", func(t *test, tk Tokens) { // 17
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodSetter,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						FormalParameterList: []BindingElement{
							{
								SingleNameBinding: &tk[5],
								Tokens:            tk[5:6],
							},
						},
						Tokens: tk[4:7],
					},
					FunctionBody: Block{
						Tokens: tk[7:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"[\n\"a\"\n]\n()\n{}", func(t *test, tk Tokens) { // 18
			litA := makeConditionLiteral(tk, 2)
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					PropertyName: PropertyName{
						ComputedPropertyName: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
						Tokens: tk[:5],
					},
					Params: FormalParameters{
						Tokens: tk[6:8],
					},
					FunctionBody: Block{
						Tokens: tk[9:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"async/* */a (){}", func(t *test, tk Tokens) { // 19
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodAsync,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[4:6],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"async/*\n*/a (){}", func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningParenthesis,
						Parsing: "FormalParameters",
						Token:   tk[2],
					},
					Parsing: "MethodDefinition",
					Token:   tk[2],
				},
				Parsing: "PropertyDefinition",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var pd PropertyDefinition
		err := pd.parse(&t.Tokens, t.Yield, t.Await)
		return pd, err
	})
}

func TestTemplateLiteral(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidTemplate,
				Parsing: "TemplateLiteral",
				Token:   tk[0],
			}
		}},
		{"``", func(t *test, tk Tokens) { // 2
			t.Output = TemplateLiteral{
				NoSubstitutionTemplate: &tk[0],
				Tokens:                 tk[:1],
			}
		}},
		{"`${\na\n}`", func(t *test, tk Tokens) { // 3
			litA := makeConditionLiteral(tk, 2)
			t.Output = TemplateLiteral{
				TemplateHead: &tk[0],
				Expressions: []Expression{
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
				},
				TemplateTail: &tk[4],
				Tokens:       tk[:5],
			}
		}},
		{"`${\na\nb\n}`", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrInvalidTemplate,
				Parsing: "TemplateLiteral",
				Token:   tk[4],
			}
		}},
		{"`${\na\n}${\nb\n}`", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			t.Output = TemplateLiteral{
				TemplateHead: &tk[0],
				Expressions: []Expression{
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
				},
				TemplateMiddleList: []*Token{
					&tk[4],
				},
				TemplateTail: &tk[8],
				Tokens:       tk[:9],
			}
		}},
		{"`${\na\n}${\nb\n}${\nc\n}`", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			litC := makeConditionLiteral(tk, 10)
			t.Output = TemplateLiteral{
				TemplateHead: &tk[0],
				Expressions: []Expression{
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litC,
								Tokens:                tk[10:11],
							},
						},
						Tokens: tk[10:11],
					},
				},
				TemplateMiddleList: []*Token{
					&tk[4],
					&tk[8],
				},
				TemplateTail: &tk[12],
				Tokens:       tk[:13],
			}
		}},
		{"`${,}`", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[1]),
					Parsing: "Expression",
					Token:   tk[1],
				},
				Parsing: "TemplateLiteral",
				Token:   tk[1],
			}
		}},
	}, func(t *test) (Type, error) {
		var tl TemplateLiteral
		err := tl.parse(&t.Tokens, t.Yield, t.Await)
		return tl, err
	})
}

func TestArrowFunction(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidAsyncArrowFunction,
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"async (\nnull\n)", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[4],
					},
					Parsing: "FormalParameters",
					Token:   tk[4],
				},
				Parsing: "ArrowFunction",
				Token:   tk[2],
			}
		}},
		{"async (\n...\na\nb)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingClosingParenthesis,
					Parsing: "FormalParameters",
					Token:   tk[8],
				},
				Parsing: "ArrowFunction",
				Token:   tk[2],
			}
		}},
		{"async", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ArrowFunction",
				Token:   tk[1],
			}
		}},
		{"async\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ArrowFunction",
				Token:   tk[1],
			}
		}},
		{"a\n", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[1],
			}
		}},
		{"a ", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[2],
			}
		}},
		{"(a)\n", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[3],
			}
		}},
		{"(a) ", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[4],
			}
		}},
		{"async a\n", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[3],
			}
		}},
		{"async a ", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[4],
			}
		}},
		{"async (a)\n", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[5],
			}
		}},
		{"async (a) ", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err:     ErrMissingArrow,
				Parsing: "ArrowFunction",
				Token:   tk[6],
			}
		}},
		{"a=>{:}", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     assignmentError(tk[3]),
								Parsing: "Expression",
								Token:   tk[3],
							},
							Parsing: "Statement",
							Token:   tk[3],
						},
						Parsing: "StatementListItem",
						Token:   tk[3],
					},
					Parsing: "Block",
					Token:   tk[3],
				},
				Parsing: "ArrowFunction",
				Token:   tk[2],
			}
		}},
		{"a =>\n{}", func(t *test, tk Tokens) { // 15
			t.Output = ArrowFunction{
				BindingIdentifier: &tk[0],
				FunctionBody: &Block{
					Tokens: tk[4:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"() =>\n{}", func(t *test, tk Tokens) { // 16
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					Tokens: tk[:2],
				},
				FunctionBody: &Block{
					Tokens: tk[5:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"async a =>\n{}", func(t *test, tk Tokens) { // 17
			t.Output = ArrowFunction{
				Async:             true,
				BindingIdentifier: &tk[2],
				FunctionBody: &Block{
					Tokens: tk[6:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"async () =>\n{}", func(t *test, tk Tokens) { // 18
			t.Output = ArrowFunction{
				Async: true,
				FormalParameters: &FormalParameters{
					Tokens: tk[2:4],
				},
				FunctionBody: &Block{
					Tokens: tk[7:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"a=>:", func(t *test, tk Tokens) { // 19
			t.Err = Error{
				Err:     assignmentError(tk[2]),
				Parsing: "ArrowFunction",
				Token:   tk[2],
			}
		}},
		{"a =>\nb", func(t *test, tk Tokens) { // 20
			litB := makeConditionLiteral(tk, 4)
			t.Output = ArrowFunction{
				BindingIdentifier: &tk[0],
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"() =>\nb", func(t *test, tk Tokens) { // 21
			litB := makeConditionLiteral(tk, 5)
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					Tokens: tk[:2],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[5:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"async a =>\nb", func(t *test, tk Tokens) { // 22
			litB := makeConditionLiteral(tk, 6)
			t.Output = ArrowFunction{
				Async:             true,
				BindingIdentifier: &tk[2],
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[6:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"async () =>\nb", func(t *test, tk Tokens) { // 23
			litB := makeConditionLiteral(tk, 7)
			t.Output = ArrowFunction{
				Async: true,
				FormalParameters: &FormalParameters{
					Tokens: tk[2:4],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[7:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"([a]) => a", func(t *test, tk Tokens) { // 24
			litAb := makeConditionLiteral(tk, 8)
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								BindingElementList: []BindingElement{
									{
										SingleNameBinding: &tk[2],
										Tokens:            tk[2:3],
									},
								},
								Tokens: tk[1:4],
							},
							Tokens: tk[1:4],
						},
					},
					Tokens: tk[:5],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litAb,
					Tokens:                tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"({a}) => a", func(t *test, tk Tokens) { // 25
			litAb := makeConditionLiteral(tk, 8)
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								BindingPropertyList: []BindingProperty{
									{
										PropertyName: PropertyName{
											LiteralPropertyName: &tk[2],
											Tokens:              tk[2:3],
										},
										BindingElement: BindingElement{
											SingleNameBinding: &tk[2],
											Tokens:            tk[2:3],
										},
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[1:4],
							},
							Tokens: tk[1:4],
						},
					},
					Tokens: tk[:5],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litAb,
					Tokens:                tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"([a, {b}, {c: d}, {e, f, ...g}, {h: {i = j}}, ...k], {l: [m, n] = [o, p]}, ...q) => r", func(t *test, tk Tokens) { // 26
			litJ := makeConditionLiteral(tk, 39)
			litR := makeConditionLiteral(tk, 77)
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								BindingElementList: []BindingElement{
									{
										SingleNameBinding: &tk[2],
										Tokens:            tk[2:3],
									},
									{
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
										Tokens: tk[5:8],
									},
									{
										ObjectBindingPattern: &ObjectBindingPattern{
											BindingPropertyList: []BindingProperty{
												{
													PropertyName: PropertyName{
														LiteralPropertyName: &tk[11],
														Tokens:              tk[11:12],
													},
													BindingElement: BindingElement{
														SingleNameBinding: &tk[14],
														Tokens:            tk[14:15],
													},
													Tokens: tk[11:15],
												},
											},
											Tokens: tk[10:16],
										},
										Tokens: tk[10:16],
									},
									{
										ObjectBindingPattern: &ObjectBindingPattern{
											BindingPropertyList: []BindingProperty{
												{
													PropertyName: PropertyName{
														LiteralPropertyName: &tk[19],
														Tokens:              tk[19:20],
													},
													BindingElement: BindingElement{
														SingleNameBinding: &tk[19],
														Tokens:            tk[19:20],
													},
													Tokens: tk[19:20],
												},
												{
													PropertyName: PropertyName{
														LiteralPropertyName: &tk[22],
														Tokens:              tk[22:23],
													},
													BindingElement: BindingElement{
														SingleNameBinding: &tk[22],
														Tokens:            tk[22:23],
													},
													Tokens: tk[22:23],
												},
											},
											BindingRestProperty: &tk[26],
											Tokens:              tk[18:28],
										},
										Tokens: tk[18:28],
									},
									{
										ObjectBindingPattern: &ObjectBindingPattern{
											BindingPropertyList: []BindingProperty{
												{
													PropertyName: PropertyName{
														LiteralPropertyName: &tk[31],
														Tokens:              tk[31:32],
													},
													BindingElement: BindingElement{
														ObjectBindingPattern: &ObjectBindingPattern{
															BindingPropertyList: []BindingProperty{
																{
																	PropertyName: PropertyName{
																		LiteralPropertyName: &tk[35],
																		Tokens:              tk[35:36],
																	},
																	BindingElement: BindingElement{
																		SingleNameBinding: &tk[35],
																		Initializer: &AssignmentExpression{
																			ConditionalExpression: &litJ,
																			Tokens:                tk[39:40],
																		},
																		Tokens: tk[35:40],
																	},
																	Tokens: tk[35:40],
																},
															},
															Tokens: tk[34:41],
														},
														Tokens: tk[34:41],
													},
													Tokens: tk[31:41],
												},
											},
											Tokens: tk[30:42],
										},
										Tokens: tk[30:42],
									},
								},
								BindingRestElement: &BindingElement{
									SingleNameBinding: &tk[45],
									Tokens:            tk[45:46],
								},
								Tokens: tk[1:47],
							},
							Tokens: tk[1:47],
						},
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								BindingPropertyList: []BindingProperty{
									{
										PropertyName: PropertyName{
											LiteralPropertyName: &tk[50],
											Tokens:              tk[50:51],
										},
										BindingElement: BindingElement{
											ArrayBindingPattern: &ArrayBindingPattern{
												BindingElementList: []BindingElement{
													{
														SingleNameBinding: &tk[54],
														Tokens:            tk[54:55],
													},
													{
														SingleNameBinding: &tk[57],
														Tokens:            tk[57:58],
													},
												},
												Tokens: tk[53:59],
											},
											Initializer: &AssignmentExpression{
												ConditionalExpression: WrapConditional(&ArrayLiteral{
													ElementList: []AssignmentExpression{
														{
															ConditionalExpression: WrapConditional(&PrimaryExpression{
																IdentifierReference: &tk[63],
																Tokens:              tk[63:64],
															}),
															Tokens: tk[63:64],
														},
														{
															ConditionalExpression: WrapConditional(&PrimaryExpression{
																IdentifierReference: &tk[66],
																Tokens:              tk[66:67],
															}),
															Tokens: tk[66:67],
														},
													},
													Tokens: tk[62:68],
												}),
												Tokens: tk[62:68],
											},
											Tokens: tk[53:68],
										},
										Tokens: tk[50:68],
									},
								},
								Tokens: tk[49:69],
							},
							Tokens: tk[49:69],
						},
					},
					BindingIdentifier: &tk[72],
					Tokens:            tk[:74],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litR,
					Tokens:                tk[77:78],
				},
				Tokens: tk[:78],
			}
		}},
		{"([, a]) => b", func(t *test, tk Tokens) { // 27
			litB := makeConditionLiteral(tk, 10)
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								BindingElementList: []BindingElement{
									{},
									{
										SingleNameBinding: &tk[4],
										Tokens:            tk[4:5],
									},
								},
								Tokens: tk[1:6],
							},
							Tokens: tk[1:6],
						},
					},
					Tokens: tk[0:7],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[10:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"(()=>a) => b", func(t *test, tk Tokens) { // 28
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"(a/=b) => b", func(t *test, tk Tokens) { // 29
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidAssignment,
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"([a /= b]) => 1", func(t *test, tk Tokens) { // 30
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrInvalidAssignment,
								Parsing: "BindingElement",
								Token:   tk[2],
							},
							Parsing: "ArrayBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"([...a /= b]) => 1", func(t *test, tk Tokens) { // 31
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrInvalidAssignment,
								Parsing: "BindingElement",
								Token:   tk[3],
							},
							Parsing: "ArrayBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"({a() {}}) => 1", func(t *test, tk Tokens) { // 32
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "ObjectBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"({a: b /= 1}) => 1", func(t *test, tk Tokens) { // 33
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrInvalidAssignment,
								Parsing: "BindingElement",
								Token:   tk[5],
							},
							Parsing: "ObjectBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"({...a=>1}) => 1", func(t *test, tk Tokens) { // 34
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrNoIdentifier,
							Parsing: "ObjectBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"(``) => 1", func(t *test, tk Tokens) { // 35
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"(`` = 1) => 1", func(t *test, tk Tokens) { // 36
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"([`` = 1] = 1) => 1", func(t *test, tk Tokens) { // 37
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "BindingElement",
								Token:   tk[2],
							},
							Parsing: "ArrayBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"({} = {}) => 1", func(t *test, tk Tokens) { // 38
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								BindingPropertyList: []BindingProperty{},
								Tokens:              tk[1:3],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&ObjectLiteral{
									Tokens: tk[6:8],
								}),
								Tokens: tk[6:8],
							},
							Tokens: tk[1:8],
						},
					},
					Tokens: tk[:9],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						Literal: &tk[12],
						Tokens:  tk[12:13],
					}),
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"([...a] = []) => 1", func(t *test, tk Tokens) { // 39
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								BindingElementList: []BindingElement{},
								BindingRestElement: &BindingElement{
									SingleNameBinding: &tk[3],
									Tokens:            tk[3:4],
								},
								Tokens: tk[1:5],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&ArrayLiteral{
									Tokens: tk[8:10],
								}),
								Tokens: tk[8:10],
							},
							Tokens: tk[1:10],
						},
					},
					Tokens: tk[:11],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						Literal: &tk[14],
						Tokens:  tk[14:15],
					}),
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"([...a.b] = []) => 1", func(t *test, tk Tokens) { // 40
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "BindingElement",
								Token:   tk[3],
							},
							Parsing: "ArrayBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
		{"({a} = {}) => 1", func(t *test, tk Tokens) { // 41
			t.Output = ArrowFunction{
				FormalParameters: &FormalParameters{
					FormalParameterList: []BindingElement{
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								BindingPropertyList: []BindingProperty{
									{
										PropertyName: PropertyName{
											LiteralPropertyName: &tk[2],
											Tokens:              tk[2:3],
										},
										BindingElement: BindingElement{
											SingleNameBinding: &tk[2],
											Tokens:            tk[2:3],
										},
										Tokens: tk[2:3],
									},
								},
								Tokens: tk[1:4],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: WrapConditional(&ObjectLiteral{
									Tokens: tk[7:9],
								}),
								Tokens: tk[7:9],
							},
							Tokens: tk[1:9],
						},
					},
					Tokens: tk[:10],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						Literal: &tk[13],
						Tokens:  tk[13:14],
					}),
					Tokens: tk[13:14],
				},
				Tokens: tk[:14],
			}
		}},
		{"({a: a.b} = {}) => 1", func(t *test, tk Tokens) { // 42
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrNoIdentifier,
									Parsing: "BindingElement",
									Token:   tk[5],
								},
								Parsing: "BindingProperty",
								Token:   tk[2],
							},
							Parsing: "ObjectBindingPattern",
							Token:   tk[1],
						},
						Parsing: "BindingElement",
						Token:   tk[1],
					},
					Parsing: "FormalParameters",
					Token:   tk[0],
				},
				Parsing: "ArrowFunction",
				Token:   tk[0],
			}
		}},
	}, func(t *test) (Type, error) {
		var (
			pe  PrimaryExpression
			af  ArrowFunction
			err error
		)
		g := t.Tokens.NewGoal()
		if err = pe.parse(&g, t.Yield, t.Await); err == nil {
			t.Tokens.Score(g)
			err = af.parse(&t.Tokens, &pe, t.In, t.Yield, t.Await)
		} else {
			err = af.parse(&t.Tokens, nil, t.In, t.Yield, t.Await)
		}
		return af, err
	})
}
