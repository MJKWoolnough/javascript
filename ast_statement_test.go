package javascript

import "testing"

func TestStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{`;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Tokens: tk[:1],
			}
		}},
		{`continue;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type:   StatementContinue,
				Tokens: tk[:2],
			}
		}},
		{`continue ;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type:   StatementContinue,
				Tokens: tk[:3],
			}
		}},
		{`continue Name;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type: StatementContinue,
				ContinueStatement: &LabelIdentifier{
					Identifier: &tk[2],
				},
				Tokens: tk[:4],
			}
		}},
		{`break;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type:   StatementBreak,
				Tokens: tk[:2],
			}
		}},
		{`break ;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type:   StatementBreak,
				Tokens: tk[:3],
			}
		}},
		{`break Name;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				Type: StatementBreak,
				BreakStatement: &LabelIdentifier{
					Identifier: &tk[2],
				},
				Tokens: tk[:4],
			}
		}},
		{`debugger;`, func(t *test, tk Tokens) {
			t.Output = Statement{
				DebuggerStatement: &tk[0],
				Tokens:            tk[:2],
			}
		}},
		{`if (1) {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 3)
			t.Output = Statement{
				IfStatement: &IfStatement{
					Expression: Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[6:8],
						},
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`if (1) {} else {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 3)
			t.Output = Statement{
				IfStatement: &IfStatement{
					Expression: Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[6:8],
						},
						Tokens: tk[6:8],
					},
					ElseStatement: &Statement{
						BlockStatement: &Block{
							Tokens: tk[11:13],
						},
						Tokens: tk[11:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`do {} while(1);`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 7)
			t.Output = Statement{
				IterationStatementDo: &IterationStatementDo{
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[2:4],
						},
						Tokens: tk[2:4],
					},
					Expression: Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[7:8],
							},
						},
						Tokens: tk[7:8],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`while(1){}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
				IterationStatementWhile: &IterationStatementWhile{
					Expression: Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[4:6],
						},
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`for(;;) {}`, func(t *test, tk Tokens) {
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[6:8],
						},
						Tokens: tk[6:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`for(i = a; b < c; d++) {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 6)
			litB := makeConditionLiteral(tk, 9)
			litC := makeConditionLiteral(tk, 13)
			litD := makeConditionLiteral(tk, 16)
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Type: ForNormalExpression,
					InitExpression: &Expression{
						Expressions: []AssignmentExpression{
							{
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
								AssignmentOperator: AssignmentAssign,
								AssignmentExpression: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[6:7],
								},
								Tokens: tk[2:7],
							},
						},
						Tokens: tk[2:7],
					},
					Conditional: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(RelationalExpression{
										RelationalExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
										RelationshipOperator: RelationshipLessThan,
										ShiftExpression:      litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
										Tokens:               tk[9:14],
									}).LogicalORExpression,
									Tokens: tk[9:14],
								},
								Tokens: tk[9:14],
							},
						},
						Tokens: tk[9:14],
					},
					Afterthought: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(UpdateExpression{
										LeftHandSideExpression: litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
										UpdateOperator:         UpdatePostIncrement,
										Tokens:                 tk[16:18],
									}).LogicalORExpression,
									Tokens: tk[16:18],
								},
								Tokens: tk[16:18],
							},
						},
						Tokens: tk[16:18],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[20:22],
						},
						Tokens: tk[20:22],
					},
					Tokens: tk[:22],
				},
				Tokens: tk[:22],
			}
		}},
		{`for(var i = a; b > c; d--) {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 18)
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Type: ForNormalVar,
					InitVar: &VariableDeclaration{
						BindingIdentifier: &BindingIdentifier{Identifier: &tk[4]},
						Initializer: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[8:9],
						},
						Tokens: tk[4:9],
					},
					Conditional: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(RelationalExpression{
										RelationalExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
										RelationshipOperator: RelationshipGreaterThan,
										ShiftExpression:      litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
										Tokens:               tk[11:16],
									}).LogicalORExpression,
									Tokens: tk[11:16],
								},
								Tokens: tk[11:16],
							},
						},
						Tokens: tk[11:16],
					},
					Afterthought: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(UpdateExpression{
										LeftHandSideExpression: litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
										UpdateOperator:         UpdatePostDecrement,
										Tokens:                 tk[18:20],
									}).LogicalORExpression,
									Tokens: tk[18:20],
								},
								Tokens: tk[18:20],
							},
						},
						Tokens: tk[18:20],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[22:24],
						},
						Tokens: tk[22:24],
					},
					Tokens: tk[:24],
				},
				Tokens: tk[:24],
			}
		}},
		{`for(let i = a; b <= c; ++d) {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 19)
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Type: ForNormalLexicalDeclaration,
					InitLexical: &LexicalDeclaration{
						LetOrConst: Let,
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &BindingIdentifier{Identifier: &tk[4]},
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[8:9],
								},
								Tokens: tk[4:9],
							},
						},
						Tokens: tk[2:10],
					},
					Conditional: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(RelationalExpression{
										RelationalExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
										RelationshipOperator: RelationshipLessThanEqual,
										ShiftExpression:      litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
										Tokens:               tk[11:16],
									}).LogicalORExpression,
									Tokens: tk[11:16],
								},
								Tokens: tk[11:16],
							},
						},
						Tokens: tk[11:16],
					},
					Afterthought: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(UpdateExpression{
										UpdateOperator:  UpdatePreIncrement,
										UnaryExpression: &litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
										Tokens:          tk[18:20],
									}).LogicalORExpression,
									Tokens: tk[18:20],
								},
								Tokens: tk[18:20],
							},
						},
						Tokens: tk[18:20],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[22:24],
						},
						Tokens: tk[22:24],
					},
					Tokens: tk[:24],
				},
				Tokens: tk[:24],
			}
		}},
		{`for(const i = a; b >= c; ++d) {}`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 19)
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Type: ForNormalLexicalDeclaration,
					InitLexical: &LexicalDeclaration{
						LetOrConst: Const,
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &BindingIdentifier{Identifier: &tk[4]},
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[8:9],
								},
								Tokens: tk[4:9],
							},
						},
						Tokens: tk[2:10],
					},
					Conditional: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(RelationalExpression{
										RelationalExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
										RelationshipOperator: RelationshipGreaterThanEqual,
										ShiftExpression:      litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
										Tokens:               tk[11:16],
									}).LogicalORExpression,
									Tokens: tk[11:16],
								},
								Tokens: tk[11:16],
							},
						},
						Tokens: tk[11:16],
					},
					Afterthought: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &ConditionalExpression{
									LogicalORExpression: wrapConditional(UpdateExpression{
										UpdateOperator:  UpdatePreIncrement,
										UnaryExpression: &litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
										Tokens:          tk[18:20],
									}).LogicalORExpression,
									Tokens: tk[18:20],
								},
								Tokens: tk[18:20],
							},
						},
						Tokens: tk[18:20],
					},
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[22:24],
						},
						Tokens: tk[22:24],
					},
					Tokens: tk[:24],
				},
				Tokens: tk[:24],
			}
		}},
		{`for(a in b) {}`, func(t *test, tk Tokens) {
			litB := makeConditionLiteral(tk, 6)
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Type: ForInLeftHandSide,
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
					In: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[6:7],
							},
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[:11],
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[9:11],
						},
						Tokens: tk[9:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseStatement(t.Yield, t.Await, t.Ret)
	})
}
