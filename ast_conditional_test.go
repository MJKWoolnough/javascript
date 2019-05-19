package javascript

import "testing"

func makeConditionLiteral(tk Tokens, pos int) ConditionalExpression {
	return ConditionalExpression{
		LogicalORExpression: LogicalORExpression{
			LogicalANDExpression: LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			Literal: &tk[pos],
																			Tokens:  tk[pos : pos+1],
																		},
																		Tokens: tk[pos : pos+1],
																	},
																	Tokens: tk[pos : pos+1],
																},
																Tokens: tk[pos : pos+1],
															},
															Tokens: tk[pos : pos+1],
														},
														Tokens: tk[pos : pos+1],
													},
													Tokens: tk[pos : pos+1],
												},
												Tokens: tk[pos : pos+1],
											},
											Tokens: tk[pos : pos+1],
										},
										Tokens: tk[pos : pos+1],
									},
									Tokens: tk[pos : pos+1],
								},
								Tokens: tk[pos : pos+1],
							},
							Tokens: tk[pos : pos+1],
						},
						Tokens: tk[pos : pos+1],
					},
					Tokens: tk[pos : pos+1],
				},
				Tokens: tk[pos : pos+1],
			},
			Tokens: tk[pos : pos+1],
		},
		Tokens: tk[pos : pos+1],
	}
}

func TestConditional(t *testing.T) {
	doTests(t, []sourceFn{
		{`true`, func(t *test, tk Tokens) {
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`true || false`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = ConditionalExpression{
				LogicalORExpression: LogicalORExpression{
					LogicalORExpression:  &litA.LogicalORExpression,
					LogicalANDExpression: litB.LogicalORExpression.LogicalANDExpression,
					Tokens:               tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseConditionalExpression(t.In, t.Yield, t.Await)
	})
}
