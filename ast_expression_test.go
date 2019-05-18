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
			t.Output = NewExpression{
				MemberExpression: MemberExpression{
					SuperProperty: true,
					Expression: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
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
																										IdentifierReference: &IdentifierReference{Identifier: &tk[2]},
																										Tokens:              tk[2:3],
																									},
																									Tokens: tk[2:3],
																								},
																								Tokens: tk[2:3],
																							},
																							Tokens: tk[2:3],
																						},
																						Tokens: tk[2:3],
																					},
																					Tokens: tk[2:3],
																				},
																				Tokens: tk[2:3],
																			},
																			Tokens: tk[2:3],
																		},
																		Tokens: tk[2:3],
																	},
																	Tokens: tk[2:3],
																},
																Tokens: tk[2:3],
															},
															Tokens: tk[2:3],
														},
														Tokens: tk[2:3],
													},
													Tokens: tk[2:3],
												},
												Tokens: tk[2:3],
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:3],
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
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
