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

func assignmentError(tk Token) Error {
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
																				Err:     ErrNoIdentifier,
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
	}, func(t *test) (interface{}, error) {
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
						Tokens: tk[4:4],
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
	}, func(t *test) (interface{}, error) {
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
		{`const a`, func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrInvalidLexicalDeclaration,
				Parsing: "LexicalDeclaration",
				Token:   tk[3],
			}
		}},
		{"const\na;", func(t *test, tk Tokens) { // 4
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
		{"const\na,\nb;", func(t *test, tk Tokens) { // 5
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
		{"const \n a;", func(t *test, tk Tokens) { // 6
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
		{"const \n a, \n b;", func(t *test, tk Tokens) { // 7
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
		{`let`, func(t *test, tk Tokens) { // 8
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
		{`let a`, func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrInvalidLexicalDeclaration,
				Parsing: "LexicalDeclaration",
				Token:   tk[3],
			}
		}},
		{"let\na;", func(t *test, tk Tokens) { // 10
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
		{"let\na,\nb;", func(t *test, tk Tokens) { // 11
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
		{"let \n a;", func(t *test, tk Tokens) { // 12
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
		{"let \n a, \n b;", func(t *test, tk Tokens) { // 13
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
	}, func(t *test) (interface{}, error) {
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
		{`{a}`, func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingInitializer,
				Parsing: "LexicalBinding",
				Token:   tk[3],
			}
		}},
		{"{a}\n=\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     assignmentError(tk[6]),
				Parsing: "LexicalBinding",
				Token:   tk[6],
			}
		}},
		{"{}\n=\na", func(t *test, tk Tokens) { // 6
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
	}, func(t *test) (interface{}, error) {
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
	}, func(t *test) (interface{}, error) {
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
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
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
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						SingleNameBinding: &tk[5],
						Tokens:            tk[5:6],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{"{\na,\n...\nb\n}", func(t *test, tk Tokens) { // 11
			t.Output = ObjectBindingPattern{
				BindingPropertyList: []BindingProperty{
					{
						SingleNameBinding: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				BindingRestProperty: &tk[7],
				Tokens:              tk[:10],
			}
		}},
	}, func(t *test) (interface{}, error) {
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
				SingleNameBinding: &tk[0],
				Tokens:            tk[:1],
			}
		}},
		{"a\n=\n", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "BindingProperty",
				Token:   tk[4],
			}
		}},
		{"a\n=\n1", func(t *test, tk Tokens) { // 5
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = BindingProperty{
				SingleNameBinding: &tk[0],
				Initializer: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
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
				PropertyName: &PropertyName{
					LiteralPropertyName: &tk[0],
					Tokens:              tk[:1],
				},
				BindingElement: &BindingElement{
					SingleNameBinding: &tk[4],
					Tokens:            tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (interface{}, error) {
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
	}, func(t *test) (interface{}, error) {
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
			t.Output = ObjectLiteral{
				PropertyDefinitionList: []PropertyDefinition{
					{
						IdentifierReference: &tk[2],
						Tokens:              tk[2:3],
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
			t.Output = ObjectLiteral{
				PropertyDefinitionList: []PropertyDefinition{
					{
						IdentifierReference: &tk[2],
						Tokens:              tk[2:3],
					},
					{
						IdentifierReference: &tk[6],
						Tokens:              tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
	}, func(t *test) (interface{}, error) {
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
			t.Output = PropertyDefinition{
				IdentifierReference: &tk[0],
				Tokens:              tk[:1],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 5
			t.Output = PropertyDefinition{
				IdentifierReference: &tk[2],
				Tokens:              tk[2:3],
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
				IdentifierReference: &tk[0],
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
					Err:     ErrMissingOpeningParenthesis,
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
						Tokens: tk[3:3],
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
						Tokens: tk[5:5],
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
						Tokens: tk[5:5],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"async\na\n(){}", func(t *test, tk Tokens) { // 15
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodAsync,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[5:5],
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
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"set\na\n(){}", func(t *test, tk Tokens) { // 17
			t.Output = PropertyDefinition{
				MethodDefinition: &MethodDefinition{
					Type: MethodSetter,
					PropertyName: PropertyName{
						LiteralPropertyName: &tk[2],
						Tokens:              tk[2:3],
					},
					Params: FormalParameters{
						Tokens: tk[5:5],
					},
					FunctionBody: Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
	}, func(t *test) (interface{}, error) {
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
	}, func(t *test) (interface{}, error) {
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
				Token:   tk[4],
			}
		}},
		{"async (\n...\na\nb)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "ArrowFunction",
				Token:   tk[8],
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
				CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
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
					Tokens: tk[3:3],
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
				CoverParenthesizedExpressionAndArrowParameterList: &CoverParenthesizedExpressionAndArrowParameterList{
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
					Tokens: tk[3:3],
				},
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[7:8],
				},
				Tokens: tk[:8],
			}
		}},
	}, func(t *test) (interface{}, error) {
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
