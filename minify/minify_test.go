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
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "false"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!1"),
			},
		},
		{ // 2
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "true"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenIdentifier, "!0"),
			},
		},
		{ // 3
			[]Option{Literals()},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "undefined"),
			},
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "void 0"),
			},
		},
		{ // 4
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "100"),
			},
		},
		{ // 5
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e3"),
			},
		},
		{ // 6
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123450000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "12345e4"),
			},
		},
		{ // 7
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.01"),
			},
		},
		{ // 8
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e-3"),
			},
		},
		{ // 9
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0.00123400"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1234e-6"),
			},
		},
		{ // 10
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xff"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "255"),
			},
		},
		{ // 11
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "999999999999"),
			},
		},
		{ // 12
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000000"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1e12"),
			},
		},
		{ // 13
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "1000000000001"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0xe8d4a51001"),
			},
		},
		{ // 14
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0o7"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "7"),
			},
		},
		{ // 15
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "0b10"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "2"),
			},
		},
		{ // 16
			[]Option{Numbers()},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123_456"),
			},
			&javascript.PrimaryExpression{
				Literal: makeToken(javascript.TokenNumericLiteral, "123456"),
			},
		},
		{ // 17
			[]Option{ArrowFn()},
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
			[]Option{ArrowFn()},
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
			[]Option{ArrowFn()},
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
			[]Option{ArrowFn()},
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
	} {
		w := walker{New(test.Options...)}
		w.Handle(test.Input)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting \n%+v\n...got...\n%+v", n+1, test.Output, test.Input)
		}
	}
}
