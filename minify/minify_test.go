package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/parser"
)

func makeToken(typ parser.TokenType, data string) *javascript.Token {
	return &javascript.Token{
		Token: parser.Token{
			Type: typ,
			Data: data,
		},
	}
}

func TestTransforms(t *testing.T) {
	for n, test := range [...]struct {
		Options       []Option
		Input, Output javascript.Type
	}{
		{ // 1
			[]Option{Literals},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "false"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!1"),
			},
		},
		{ // 2
			[]Option{Literals},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "true"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!0"),
			},
		},
		{ // 3
			[]Option{Literals},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "undefined"),
			},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "void 0"),
			},
		},
		{ // 4
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
		},
		{ // 5
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e3"),
			},
		},
		{ // 6
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123450000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "12345e4"),
			},
		},
		{ // 7
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
		},
		{ // 8
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e-3"),
			},
		},
		{ // 9
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.00123400"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1234e-6"),
			},
		},
		{ // 10
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xff"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "255"),
			},
		},
		{ // 11
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
		},
		{ // 12
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e12"),
			},
		},
		{ // 13
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xe8d4a51001"),
			},
		},
		{ // 14
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0o7"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "7"),
			},
		},
		{ // 15
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0b10"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "2"),
			},
		},
		{ // 16
			[]Option{Numbers},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123_456"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123456"),
			},
		},
		{ // 17
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									Type: javascript.StatementReturn,
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "1"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					AssignmentExpression: &javascript.AssignmentExpression{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
				},
			},
		},
		{ // 18
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
													MemberExpression: &javascript.MemberExpression{
														PrimaryExpression: &javascript.PrimaryExpression{
															IdentifierReference: makeToken(javascript.TokenIdentifier, "m"),
														},
													},
													Arguments: &javascript.Arguments{},
												}),
											},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{},
							},
							{
								Statement: &javascript.Statement{
									Type: javascript.StatementReturn,
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "2"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					AssignmentExpression: &javascript.AssignmentExpression{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
										MemberExpression: &javascript.MemberExpression{
											PrimaryExpression: &javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "m"),
											},
										},
										Arguments: &javascript.Arguments{},
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "2"),
									}),
								},
							},
						}),
					},
				},
			},
		},
		{ // 19
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "1"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "1"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 20
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									IterationStatementWhile: &javascript.IterationStatementWhile{
										Expression: javascript.Expression{
											Expressions: []javascript.AssignmentExpression{
												{
													ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
														MemberExpression: &javascript.MemberExpression{
															PrimaryExpression: &javascript.PrimaryExpression{
																IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
															},
														},
														Arguments: &javascript.Arguments{},
													}),
												},
											},
										},
										Statement: javascript.Statement{
											BlockStatement: &javascript.Block{},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "1"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									IterationStatementWhile: &javascript.IterationStatementWhile{
										Expression: javascript.Expression{
											Expressions: []javascript.AssignmentExpression{
												{
													ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
														MemberExpression: &javascript.MemberExpression{
															PrimaryExpression: &javascript.PrimaryExpression{
																IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
															},
														},
														Arguments: &javascript.Arguments{},
													}),
												},
											},
										},
										Statement: javascript.Statement{
											BlockStatement: &javascript.Block{},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "1"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 21
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody:      &javascript.Block{},
				},
			},
		},
		{ // 22
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "b"),
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "b"),
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
		},
		{ // 23
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
								Initializer: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
								Initializer: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
		},
		{ // 24
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					},
					FunctionBody: &javascript.Block{},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					},
					FunctionBody: &javascript.Block{},
				},
			},
		},
		{ // 25
			[]Option{IfToConditional},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
					ElseStatement: &javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: &javascript.ConditionalExpression{
								LogicalORExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}).LogicalORExpression,
								True: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
								False: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
									}),
								},
							},
						},
					},
				},
			},
		},
		{ // 26
			[]Option{IfToConditional},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
			},
		},
		{ // 27
			[]Option{IfToConditional},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
								}),
							},
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
								}),
							},
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "d"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "e"),
									}),
								},
							},
						},
					},
					ElseStatement: &javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "f"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
							}),
						},
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
							}),
						},
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
							}),
						},
						{
							ConditionalExpression: &javascript.ConditionalExpression{
								LogicalORExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "d"),
								}).LogicalORExpression,
								True: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "e"),
									}),
								},
								False: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "f"),
									}),
								},
							},
						},
					},
				},
			},
		},
		{ // 28
			[]Option{IfToConditional},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
									}),
								},
							},
						},
					},
					ElseStatement: &javascript.Statement{
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "d"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "e"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: &javascript.ConditionalExpression{
								LogicalORExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}).LogicalORExpression,
								True: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
												}),
											},
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
												}),
											},
										},
									}),
								},
								False: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "d"),
												}),
											},
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "e"),
												}),
											},
										},
									}),
								},
							},
						},
					},
				},
			},
		},
		{ // 29
			[]Option{IfToConditional},
			&javascript.Statement{
				IfStatement: &javascript.IfStatement{
					Expression: javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
					Statement: javascript.Statement{
						Type: javascript.StatementReturn,
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
					ElseStatement: &javascript.Statement{
						Type: javascript.StatementReturn,
						ExpressionStatement: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				Type: javascript.StatementReturn,
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: &javascript.ConditionalExpression{
								LogicalORExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}).LogicalORExpression,
								True: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
								False: &javascript.AssignmentExpression{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
									}),
								},
							},
						},
					},
				},
			},
		},
		{ // 30
			[]Option{ArrowFn, IfToConditional},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ArrowFunction: &javascript.ArrowFunction{
								BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
								FunctionBody: &javascript.Block{
									StatementList: []javascript.StatementListItem{
										{
											Statement: &javascript.Statement{
												IfStatement: &javascript.IfStatement{
													Expression: javascript.Expression{
														Expressions: []javascript.AssignmentExpression{
															{
																ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
																	IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
																}),
															},
														},
													},
													Statement: javascript.Statement{
														Type: javascript.StatementReturn,
														ExpressionStatement: &javascript.Expression{
															Expressions: []javascript.AssignmentExpression{
																{
																	ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
																		IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
																	}),
																},
															},
														},
													},
													ElseStatement: &javascript.Statement{
														Type: javascript.StatementReturn,
														ExpressionStatement: &javascript.Expression{
															Expressions: []javascript.AssignmentExpression{
																{
																	ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
																		IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
																	}),
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ArrowFunction: &javascript.ArrowFunction{
								BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
								AssignmentExpression: &javascript.AssignmentExpression{
									ConditionalExpression: &javascript.ConditionalExpression{
										LogicalORExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
										}).LogicalORExpression,
										True: &javascript.AssignmentExpression{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
											}),
										},
										False: &javascript.AssignmentExpression{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 31
			[]Option{ArrowFn},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
													MemberExpression: &javascript.MemberExpression{
														PrimaryExpression: &javascript.PrimaryExpression{
															IdentifierReference: makeToken(javascript.TokenIdentifier, "m"),
														},
													},
													Arguments: &javascript.Arguments{},
												}),
											},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									Type: javascript.StatementReturn,
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "2"),
												}),
											},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									VariableStatement: &javascript.VariableStatement{
										VariableDeclarationList: []javascript.VariableDeclaration{
											{
												BindingIdentifier: makeToken(javascript.TokenIdentifier, "m"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
													MemberExpression: &javascript.MemberExpression{
														PrimaryExpression: &javascript.PrimaryExpression{
															IdentifierReference: makeToken(javascript.TokenIdentifier, "m"),
														},
													},
													Arguments: &javascript.Arguments{},
												}),
											},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									Type: javascript.StatementReturn,
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "2"),
												}),
											},
										},
									},
								},
							},
							{
								Statement: &javascript.Statement{
									VariableStatement: &javascript.VariableStatement{
										VariableDeclarationList: []javascript.VariableDeclaration{
											{
												BindingIdentifier: makeToken(javascript.TokenIdentifier, "m"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 32
			[]Option{RemoveDebugger},
			&javascript.Statement{
				Type: javascript.StatementDebugger,
			},
			&javascript.Statement{},
		},
		{ // 33
			[]Option{BlocksToStatement},
			&javascript.Statement{
				BlockStatement: &javascript.Block{},
			},
			&javascript.Statement{
				BlockStatement: &javascript.Block{},
			},
		},
		{ // 34
			[]Option{BlocksToStatement},
			&javascript.Statement{
				BlockStatement: &javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								Type: javascript.StatementContinue,
							},
						},
					},
				},
			},
			&javascript.Statement{
				Type: javascript.StatementContinue,
			},
		},
		{ // 35
			[]Option{BlocksToStatement},
			&javascript.Statement{
				BlockStatement: &javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								BlockStatement: &javascript.Block{
									StatementList: []javascript.StatementListItem{
										{
											Statement: &javascript.Statement{
												Type: javascript.StatementContinue,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				Type: javascript.StatementContinue,
			},
		},
		{ // 36
			[]Option{BlocksToStatement},
			&javascript.Statement{
				BlockStatement: &javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								BlockStatement: &javascript.Block{
									StatementList: []javascript.StatementListItem{
										{
											Statement: &javascript.Statement{
												Type: javascript.StatementContinue,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				Type: javascript.StatementContinue,
			},
		},
		{ // 37
			[]Option{BlocksToStatement},
			&javascript.Statement{
				BlockStatement: &javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "hello"),
											}),
										},
									},
								},
							},
						},
						{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "world"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "hello"),
							}),
						},
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "world"),
							}),
						},
					},
				},
			},
		},
		{ // 38
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"a\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "a"),
			},
		},
		{ // 39
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"ab\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "ab"),
			},
		},
		{ // 40
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"&\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"&\""),
			},
		},
		{ // 41
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"Infinity\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "Infinity"),
			},
		},
		{ // 42
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"123\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenNumericLiteral, "123"),
			},
		},
		{ // 43
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"true\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "true"),
			},
		},
		{ // 44
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"false\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "false"),
			},
		},
		{ // 45
			[]Option{Keys},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"null\""),
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "null"),
			},
		},
		{ // 46
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"a\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "a"),
			},
		},
		{ // 47
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"ab\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "ab"),
			},
		},
		{ // 48
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"&\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenStringLiteral, "\"&\""),
			},
		},
		{ // 49
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"Infinity\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "Infinity"),
			},
		},
		{ // 50
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"123\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenNumericLiteral, "123"),
			},
		},
		{ // 51
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"true\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "true"),
			},
		},
		{ // 52
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"false\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "false"),
			},
		},
		{ // 53
			[]Option{Keys},
			&javascript.PropertyName{
				ComputedPropertyName: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenStringLiteral, "\"null\""),
					}),
				},
			},
			&javascript.PropertyName{
				LiteralPropertyName: makeToken(javascript.TokenIdentifier, "null"),
			},
		},
		{ // 54
			[]Option{RemoveExpressionNames},
			&javascript.PrimaryExpression{
				FunctionExpression: &javascript.FunctionDeclaration{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
				},
			},
			&javascript.PrimaryExpression{
				FunctionExpression: &javascript.FunctionDeclaration{},
			},
		},
		{ // 55
			[]Option{RemoveExpressionNames},
			&javascript.PrimaryExpression{
				ClassExpression: &javascript.ClassDeclaration{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
				},
			},
			&javascript.PrimaryExpression{
				ClassExpression: &javascript.ClassDeclaration{},
			},
		},
		{ // 56
			[]Option{FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{}),
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{},
					FunctionBody:     &javascript.Block{},
				},
			},
		},
		{ // 57
			[]Option{FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FunctionBody: javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
												}),
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{},
					FunctionBody: &javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
												}),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 58
			[]Option{FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FunctionBody: javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenKeyword, "this"),
												}),
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FunctionBody: javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenKeyword, "this"),
												}),
											},
										},
									},
								},
							},
						},
					},
				}),
			},
		},
		{ // 59
			[]Option{FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FunctionBody: javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "arguments"),
												}),
											},
										},
									},
								},
							},
						},
					},
				}),
			},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FunctionBody: javascript.Block{
						StatementList: []javascript.StatementListItem{
							{
								Statement: &javascript.Statement{
									ExpressionStatement: &javascript.Expression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "arguments"),
												}),
											},
										},
									},
								},
							},
						},
					},
				}),
			},
		},
		{ // 60
			[]Option{FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FormalParameters: javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
					},
				}),
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					FormalParameters: &javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
					},
					FunctionBody: &javascript.Block{},
				},
			},
		},
		{ // 61
			[]Option{ArrowFn},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
			&javascript.ArrowFunction{
				BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
				FunctionBody:      &javascript.Block{},
			},
		},
		{ // 62
			[]Option{ArrowFn},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
						},
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "b"),
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
						},
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "b"),
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
		},
		{ // 63
			[]Option{ArrowFn},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
				},
				FunctionBody: &javascript.Block{},
			},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
				},
				FunctionBody: &javascript.Block{},
			},
		},
		{ // 64
			[]Option{ArrowFn},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							Initializer: &javascript.AssignmentExpression{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									Literal: makeToken(javascript.TokenNumericLiteral, "1"),
								}),
							},
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							Initializer: &javascript.AssignmentExpression{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									Literal: makeToken(javascript.TokenNumericLiteral, "1"),
								}),
							},
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
		},
		{ // 65
			[]Option{ArrowFn},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							ArrayBindingPattern: &javascript.ArrayBindingPattern{
								BindingElementList: []javascript.BindingElement{
									{
										SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
									},
								},
							},
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
			&javascript.ArrowFunction{
				FormalParameters: &javascript.FormalParameters{
					FormalParameterList: []javascript.BindingElement{
						{
							ArrayBindingPattern: &javascript.ArrayBindingPattern{
								BindingElementList: []javascript.BindingElement{
									{
										SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
									},
								},
							},
						},
					},
				},
				FunctionBody: &javascript.Block{},
			},
		},
		{ // 66
			[]Option{ArrowFn, FunctionExpressionToArrowFunc},
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.FunctionDeclaration{
					FormalParameters: javascript.FormalParameters{
						FormalParameterList: []javascript.BindingElement{
							{
								SingleNameBinding: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
					},
				}),
			},
			&javascript.AssignmentExpression{
				ArrowFunction: &javascript.ArrowFunction{
					BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					FunctionBody:      &javascript.Block{},
				},
			},
		},
		{ // 67
			[]Option{UnwrapParens},
			&javascript.Expression{},
			&javascript.Expression{},
		},
		{ // 68
			[]Option{UnwrapParens},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "1"),
									}),
								},
							},
						}),
					},
				},
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
				},
			},
		},
		{ // 69
			[]Option{UnwrapParens},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "1"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "2"),
									}),
								},
							},
						}),
					},
				},
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "2"),
						}),
					},
				},
			},
		},
		{ // 70
			[]Option{UnwrapParens},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "1"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "2"),
									}),
								},
							},
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "3"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "4"),
									}),
								},
							},
						}),
					},
				},
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "2"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "3"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "4"),
						}),
					},
				},
			},
		},
		{ // 71
			[]Option{UnwrapParens},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "1"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "2"),
												}),
											},
										},
									}),
								},
							},
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
										Expressions: []javascript.AssignmentExpression{
											{
												ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
													Literal: makeToken(javascript.TokenNumericLiteral, "3"),
												}),
											},
										},
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "4"),
									}),
								},
							},
						}),
					},
				},
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "2"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "3"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "4"),
						}),
					},
				},
			},
		},
		{ // 72
			[]Option{UnwrapParens},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										Literal: makeToken(javascript.TokenNumericLiteral, "2"),
									}),
								},
							},
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "3"),
						}),
					},
				},
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "1"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "2"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							Literal: makeToken(javascript.TokenNumericLiteral, "3"),
						}),
					},
				},
			},
		},
		{ // 73
			[]Option{UnwrapParens},
			&javascript.Argument{},
			&javascript.Argument{},
		},
		{ // 74
			[]Option{UnwrapParens},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenNumericLiteral, "1"),
					}),
				},
			},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenNumericLiteral, "1"),
					}),
				},
			},
		},
		{ // 75
			[]Option{UnwrapParens},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									Literal: makeToken(javascript.TokenNumericLiteral, "1"),
								}),
							},
						},
					}),
				},
			},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenNumericLiteral, "1"),
					}),
				},
			},
		},
		{ // 76
			[]Option{UnwrapParens},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												Literal: makeToken(javascript.TokenNumericLiteral, "1"),
											}),
										},
									},
								}),
							},
						},
					}),
				},
			},
			&javascript.Argument{
				AssignmentExpression: javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenNumericLiteral, "1"),
					}),
				},
			},
		},
		{ // 77
			[]Option{UnwrapParens},
			&javascript.AssignmentExpression{
				Yield: true,
				AssignmentExpression: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									Literal: makeToken(javascript.TokenNumericLiteral, "1"),
								}),
							},
						},
					}),
				},
			},
			&javascript.AssignmentExpression{
				Yield: true,
				AssignmentExpression: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						Literal: makeToken(javascript.TokenNumericLiteral, "1"),
					}),
				},
			},
		},
		{ // 78
			[]Option{UnwrapParens},
			&javascript.MemberExpression{},
			&javascript.MemberExpression{},
		},
		{ // 79
			[]Option{UnwrapParens},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
			},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
			},
		},
		{ // 80
			[]Option{UnwrapParens},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "c"),
			},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
									}),
								},
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "c"),
			},
		},
		{ // 81
			[]Option{UnwrapParens},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										ArrayLiteral: &javascript.ArrayLiteral{},
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "length"),
			},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ArrayLiteral: &javascript.ArrayLiteral{},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "length"),
			},
		},
		{ // 82
			[]Option{UnwrapParens},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										TemplateLiteral: &javascript.TemplateLiteral{},
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "length"),
			},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						TemplateLiteral: &javascript.TemplateLiteral{},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "length"),
			},
		},
		{ // 83
			[]Option{UnwrapParens},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.MemberExpression{
										MemberExpression: &javascript.MemberExpression{
											PrimaryExpression: &javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											},
										},
										IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "c"),
			},
			&javascript.MemberExpression{
				MemberExpression: &javascript.MemberExpression{
					MemberExpression: &javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
					IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "c"),
			},
		},
		{ // 84
			[]Option{UnwrapParens},
			&javascript.CallExpression{},
			&javascript.CallExpression{},
		},
		{ // 85
			[]Option{UnwrapParens},
			&javascript.CallExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
										MemberExpression: &javascript.MemberExpression{
											PrimaryExpression: &javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											},
										},
										Arguments: &javascript.Arguments{},
									}),
								},
							},
						},
					},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "a"),
			},
			&javascript.CallExpression{
				CallExpression: &javascript.CallExpression{
					MemberExpression: &javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
					Arguments: &javascript.Arguments{},
				},
				IdentifierName: makeToken(javascript.TokenIdentifier, "a"),
			},
		},
		{ // 86
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{},
			&javascript.LeftHandSideExpression{},
		},
		{ // 87
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							ParenthesizedExpression: &javascript.ParenthesizedExpression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
											MemberExpression: &javascript.MemberExpression{
												PrimaryExpression: &javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
												},
											},
											Arguments: &javascript.Arguments{},
										}),
									},
								},
							},
						},
					},
				},
			},
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{
					MemberExpression: &javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
					Arguments: &javascript.Arguments{},
				},
			},
		},
		{ // 88
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					News: 1,
					MemberExpression: javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							ParenthesizedExpression: &javascript.ParenthesizedExpression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
											MemberExpression: &javascript.MemberExpression{
												PrimaryExpression: &javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
												},
											},
											Arguments: &javascript.Arguments{},
										}),
									},
								},
							},
						},
					},
				},
			},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					News: 1,
					MemberExpression: javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							ParenthesizedExpression: &javascript.ParenthesizedExpression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
											MemberExpression: &javascript.MemberExpression{
												PrimaryExpression: &javascript.PrimaryExpression{
													IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
												},
											},
											Arguments: &javascript.Arguments{},
										}),
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 89
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
					},
				},
			},
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{
					CallExpression: &javascript.CallExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
						Arguments: &javascript.Arguments{},
					},
					IdentifierName: makeToken(javascript.TokenIdentifier, "b"),
				},
			},
		},
		{ // 90
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						Expression: &javascript.Expression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
			},
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{
					CallExpression: &javascript.CallExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
						Arguments: &javascript.Arguments{},
					},
					Expression: &javascript.Expression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
								}),
							},
						},
					},
				},
			},
		},
		{ // 91
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						TemplateLiteral: &javascript.TemplateLiteral{
							NoSubstitutionTemplate: makeToken(javascript.TokenNoSubstitutionTemplate, "`b`"),
						},
					},
				},
			},
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{
					CallExpression: &javascript.CallExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
						Arguments: &javascript.Arguments{},
					},
					TemplateLiteral: &javascript.TemplateLiteral{
						NoSubstitutionTemplate: makeToken(javascript.TokenNoSubstitutionTemplate, "`b`"),
					},
				},
			},
		},
		{ // 92
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						PrivateIdentifier: makeToken(javascript.TokenPrivateIdentifier, "#b"),
					},
				},
			},
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{
					CallExpression: &javascript.CallExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
							},
						},
						Arguments: &javascript.Arguments{},
					},
					PrivateIdentifier: makeToken(javascript.TokenPrivateIdentifier, "#b"),
				},
			},
		},
		{ // 93
			[]Option{UnwrapParens},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						Arguments: &javascript.Arguments{},
					},
				},
			},
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{
					MemberExpression: javascript.MemberExpression{
						MemberExpression: &javascript.MemberExpression{
							PrimaryExpression: &javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
												MemberExpression: &javascript.MemberExpression{
													PrimaryExpression: &javascript.PrimaryExpression{
														IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
													},
												},
												Arguments: &javascript.Arguments{},
											}),
										},
									},
								},
							},
						},
						Arguments: &javascript.Arguments{},
					},
				},
			},
		},
		{ // 94
			[]Option{},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{},
					},
				},
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{},
			},
		},
		{ // 95
			[]Option{},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{},
					},
					{
						Statement: &javascript.Statement{},
					},
				},
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{},
			},
		},
		{ // 96
			[]Option{},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{},
					},
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
										}),
									},
								},
							},
						},
					},
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
										}),
									},
								},
							},
						},
					},
					{
						Statement: &javascript.Statement{},
					},
					{
						Statement: &javascript.Statement{},
					},
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
										}),
									},
								},
							},
						},
					},
				},
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
										}),
									},
								},
							},
						},
					},
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
										}),
									},
								},
							},
						},
					},
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
										}),
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 97
			[]Option{},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
				},
			},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{},
			},
		},
		{ // 98
			[]Option{},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
				},
			},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{},
			},
		},
		{ // 99
			[]Option{},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
											}),
										},
									},
								},
							},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.Module{
				ModuleListItems: []javascript.ModuleItem{
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
											}),
										},
									},
								},
							},
						},
					},
					{
						StatementListItem: &javascript.StatementListItem{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "c"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 100
			[]Option{RemoveLastEmptyReturn},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								Type: javascript.StatementReturn,
							},
						},
					},
				},
			},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{},
				},
			},
		},
		{ // 101
			[]Option{RemoveLastEmptyReturn},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								Type: javascript.StatementReturn,
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								Type: javascript.StatementReturn,
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 102
			[]Option{RemoveLastEmptyReturn},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
						{
							Statement: &javascript.Statement{
								Type: javascript.StatementReturn,
							},
						},
					},
				},
			},
			&javascript.FunctionDeclaration{
				FunctionBody: javascript.Block{
					StatementList: []javascript.StatementListItem{
						{
							Statement: &javascript.Statement{
								ExpressionStatement: &javascript.Expression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
											}),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		w := walker{New(test.Options...)}
		w.Handle(test.Input)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, test.Output, test.Input)
		}
	}
}
