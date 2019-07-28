package javascript

import (
	"testing"
)

func TestStatementOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`{}`, func(t *test, tk Tokens) { // 1
			t.Output = StatementListItem{
				Statement: &Statement{
					BlockStatement: &Block{
						Tokens: tk[:2],
					},
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`var a;`, func(t *test, tk Tokens) { // 2
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								BindingIdentifier: &tk[2],
								Tokens:            tk[2:3],
							},
						},
						Tokens: tk[:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`var a = 1;`, func(t *test, tk Tokens) { // 3
			litA := makeConditionLiteral(tk, 6)
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								BindingIdentifier: &tk[2],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[6:7],
								},
								Tokens: tk[2:7],
							},
						},
						Tokens: tk[:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`var a = 1, b;`, func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 6)
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								BindingIdentifier: &tk[2],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[6:7],
								},
								Tokens: tk[2:7],
							},
							{
								BindingIdentifier: &tk[9],
								Tokens:            tk[9:10],
							},
						},
						Tokens: tk[:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{`var a = 1, b = 2;`, func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 6)
			litB := makeConditionLiteral(tk, 13)
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								BindingIdentifier: &tk[2],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[6:7],
								},
								Tokens: tk[2:7],
							},
							{
								BindingIdentifier: &tk[9],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litB,
									Tokens:                tk[13:14],
								},
								Tokens: tk[9:14],
							},
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`var a, b = 1;`, func(t *test, tk Tokens) { // 6
			litB := makeConditionLiteral(tk, 9)
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								BindingIdentifier: &tk[2],
								Tokens:            tk[2:3],
							},
							{
								BindingIdentifier: &tk[5],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litB,
									Tokens:                tk[9:10],
								},
								Tokens: tk[5:10],
							},
						},
						Tokens: tk[:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{`var [a, b] = [1, 2];`, func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 12)
			litB := makeConditionLiteral(tk, 15)
			arr := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								ArrayLiteral: &ArrayLiteral{
									ElementList: []AssignmentExpression{
										{
											ConditionalExpression: &litA,
											Tokens:                tk[12:13],
										},
										{
											ConditionalExpression: &litB,
											Tokens:                tk[15:16],
										},
									},
									Tokens: tk[11:17],
								},
								Tokens: tk[11:17],
							},
							Tokens: tk[11:17],
						},
						Tokens: tk[11:17],
					},
					Tokens: tk[11:17],
				},
				Tokens: tk[11:17],
			})
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								ArrayBindingPattern: &ArrayBindingPattern{
									BindingElementList: []BindingElement{
										{
											SingleNameBinding: &tk[3],
											Tokens:            tk[3:4],
										},
										{
											SingleNameBinding: &tk[6],
											Tokens:            tk[6:7],
										},
									},
									Tokens: tk[2:8],
								},
								Initializer: &AssignmentExpression{
									ConditionalExpression: &arr,
									Tokens:                tk[11:17],
								},
								Tokens: tk[2:17],
							},
						},
						Tokens: tk[:18],
					},
					Tokens: tk[:18],
				},
				Tokens: tk[:18],
			}
		}},
		{`var {a, b} = {c, d};`, func(t *test, tk Tokens) { // 8
			obj := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								ObjectLiteral: &ObjectLiteral{
									PropertyDefinitionList: []PropertyDefinition{
										{
											IdentifierReference: &tk[12],
											Tokens:              tk[12:13],
										},
										{
											IdentifierReference: &tk[15],
											Tokens:              tk[15:16],
										},
									},
									Tokens: tk[11:17],
								},
								Tokens: tk[11:17],
							},
							Tokens: tk[11:17],
						},
						Tokens: tk[11:17],
					},
					Tokens: tk[11:17],
				},
				Tokens: tk[11:17],
			})
			t.Output = StatementListItem{
				Statement: &Statement{
					VariableStatement: &VariableStatement{
						VariableDeclarationList: []VariableDeclaration{
							{
								ObjectBindingPattern: &ObjectBindingPattern{
									BindingPropertyList: []BindingProperty{
										{
											SingleNameBinding: &tk[3],
											Tokens:            tk[3:4],
										},
										{
											SingleNameBinding: &tk[6],
											Tokens:            tk[6:7],
										},
									},
									Tokens: tk[2:8],
								},
								Initializer: &AssignmentExpression{
									ConditionalExpression: &obj,
									Tokens:                tk[11:17],
								},
								Tokens: tk[2:17],
							},
						},
						Tokens: tk[:18],
					},
					Tokens: tk[:18],
				},
				Tokens: tk[:18],
			}

		}},
		{`;`, func(t *test, tk Tokens) { // 9
			t.Output = StatementListItem{
				Statement: &Statement{
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
		{`if (1) {}`, func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 3)
			t.Output = StatementListItem{
				Statement: &Statement{
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
				},
				Tokens: tk[:8],
			}
		}},
		{`if (1) {} else {}`, func(t *test, tk Tokens) { // 11
			litA := makeConditionLiteral(tk, 3)
			t.Output = StatementListItem{
				Statement: &Statement{
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
				},
				Tokens: tk[:13],
			}
		}},
		{`do {} while(1);`, func(t *test, tk Tokens) { // 12
			litA := makeConditionLiteral(tk, 7)
			t.Output = StatementListItem{
				Statement: &Statement{
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
				},
				Tokens: tk[:10],
			}
		}},
		{`while(1){}`, func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 2)
			t.Output = StatementListItem{
				Statement: &Statement{
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
				},
				Tokens: tk[:6],
			}
		}},
		{`for(;;) {}`, func(t *test, tk Tokens) { // 14
			t.Output = StatementListItem{
				Statement: &Statement{
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
				},
				Tokens: tk[:8],
			}
		}},
		{`for(i = a; b < c; d++) {}`, func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 6)
			litB := makeConditionLiteral(tk, 9)
			litC := makeConditionLiteral(tk, 13)
			litD := makeConditionLiteral(tk, 16)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForNormalExpression,
						InitExpression: &Expression{
							Expressions: []AssignmentExpression{
								{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[2],
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
				},
				Tokens: tk[:22],
			}
		}},
		{`for(var i = a; b > c; d--) {}`, func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 18)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForNormalVar,
						InitVar: []VariableDeclaration{
							{
								BindingIdentifier: &tk[4],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[8:9],
								},
								Tokens: tk[4:9],
							},
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
				},
				Tokens: tk[:24],
			}
		}},
		{`for(let i = a; b <= c; ++d) {}`, func(t *test, tk Tokens) { // 17
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 19)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForNormalLexicalDeclaration,
						InitLexical: &LexicalDeclaration{
							LetOrConst: Let,
							BindingList: []LexicalBinding{
								{
									BindingIdentifier: &tk[4],
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
				},
				Tokens: tk[:24],
			}
		}},
		{`for(const i = a; b >= c; ++d) {}`, func(t *test, tk Tokens) { // 18
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 11)
			litC := makeConditionLiteral(tk, 15)
			litD := makeConditionLiteral(tk, 19)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForNormalLexicalDeclaration,
						InitLexical: &LexicalDeclaration{
							LetOrConst: Const,
							BindingList: []LexicalBinding{
								{
									BindingIdentifier: &tk[4],
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
				},
				Tokens: tk[:24],
			}
		}},
		{`for(a in b) {}`, func(t *test, tk Tokens) { // 19
			litB := makeConditionLiteral(tk, 6)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForInLeftHandSide,
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[2],
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
				},
				Tokens: tk[:11],
			}
		}},
		{`for(var a in b) {}`, func(t *test, tk Tokens) { // 20
			litB := makeConditionLiteral(tk, 8)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type:                 ForInVar,
						ForBindingIdentifier: &tk[4],
						In: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[8:9],
								},
							},
							Tokens: tk[8:9],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[11:13],
						},
						Tokens: tk[:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`for(let {a} in b) {}`, func(t *test, tk Tokens) { // 21
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForInLet,
						ForBindingPatternObject: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									SingleNameBinding: &tk[5],
									Tokens:            tk[5:6],
								},
							},
							Tokens: tk[4:7],
						},
						In: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[13:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`for(const [a] in b) {}`, func(t *test, tk Tokens) { // 22
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForInConst,
						ForBindingPatternArray: &ArrayBindingPattern{
							BindingElementList: []BindingElement{
								{
									SingleNameBinding: &tk[5],
									Tokens:            tk[5:6],
								},
							},
							Tokens: tk[4:7],
						},
						In: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[10:11],
								},
							},
							Tokens: tk[10:11],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[13:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`for(a of b) {}`, func(t *test, tk Tokens) { // 23
			litB := makeConditionLiteral(tk, 6)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForOfLeftHandSide,
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[2],
										Tokens:              tk[2:3],
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[6:7],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[9:11],
							},
							Tokens: tk[9:11],
						},
						Tokens: tk[:11],
					},
					Tokens: tk[:11],
				},
				Tokens: tk[:11],
			}
		}},
		{`for(var a of b) {}`, func(t *test, tk Tokens) { // 24
			litB := makeConditionLiteral(tk, 8)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type:                 ForOfVar,
						ForBindingIdentifier: &tk[4],
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[8:9],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[11:13],
						},
						Tokens: tk[:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`for(let {a} of b) {}`, func(t *test, tk Tokens) { // 25
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForOfLet,
						ForBindingPatternObject: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									SingleNameBinding: &tk[5],
									Tokens:            tk[5:6],
								},
							},
							Tokens: tk[4:7],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[13:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`for(const [a] of b) {}`, func(t *test, tk Tokens) { // 26
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForOfConst,
						ForBindingPatternArray: &ArrayBindingPattern{
							BindingElementList: []BindingElement{
								{
									SingleNameBinding: &tk[5],
									Tokens:            tk[5:6],
								},
							},
							Tokens: tk[4:7],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[13:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`switch(true){}`, func(t *test, tk Tokens) { // 27
			litA := makeConditionLiteral(tk, 2)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`switch(a){case 0:case 1:}`, func(t *test, tk Tokens) { // 28
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 11)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						CaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litB,
											Tokens:                tk[7:8],
										},
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[5:9],
							},
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litC,
											Tokens:                tk[11:12],
										},
									},
									Tokens: tk[11:12],
								},
								Tokens: tk[9:13],
							},
						},
						Tokens: tk[:14],
					},
					Tokens: tk[:14],
				},
				Tokens: tk[:14],
			}
		}},
		{`switch(a){default:case 0:case 1:}`, func(t *test, tk Tokens) { // 29
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 9)
			litC := makeConditionLiteral(tk, 13)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						DefaultClause: []StatementListItem{},
						PostDefaultCaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litB,
											Tokens:                tk[9:10],
										},
									},
									Tokens: tk[9:10],
								},
								Tokens: tk[7:11],
							},
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litC,
											Tokens:                tk[13:14],
										},
									},
									Tokens: tk[13:14],
								},
								Tokens: tk[11:15],
							},
						},
						Tokens: tk[:16],
					},
					Tokens: tk[:16],
				},
				Tokens: tk[:16],
			}
		}},
		{`switch(a){case 0:default:case 1:}`, func(t *test, tk Tokens) { // 30
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 13)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						CaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litB,
											Tokens:                tk[7:8],
										},
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[5:9],
							},
						},
						DefaultClause: []StatementListItem{},
						PostDefaultCaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litC,
											Tokens:                tk[13:14],
										},
									},
									Tokens: tk[13:14],
								},
								Tokens: tk[11:15],
							},
						},
						Tokens: tk[:16],
					},
					Tokens: tk[:16],
				},
				Tokens: tk[:16],
			}
		}},
		{`switch(a){case 0:case 1:default:}`, func(t *test, tk Tokens) { // 31
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 11)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						CaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litB,
											Tokens:                tk[7:8],
										},
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[5:9],
							},
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litC,
											Tokens:                tk[11:12],
										},
									},
									Tokens: tk[11:12],
								},
								Tokens: tk[9:13],
							},
						},
						DefaultClause: []StatementListItem{},
						Tokens:        tk[:16],
					},
					Tokens: tk[:16],
				},
				Tokens: tk[:16],
			}
		}},
		{`switch(a){case b:case c:d;default:e;f;case g:h;}`, func(t *test, tk Tokens) { // 32
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 11)
			litD := makeConditionLiteral(tk, 13)
			litE := makeConditionLiteral(tk, 17)
			litF := makeConditionLiteral(tk, 19)
			litG := makeConditionLiteral(tk, 23)
			litH := makeConditionLiteral(tk, 25)
			t.Output = StatementListItem{
				Statement: &Statement{
					SwitchStatement: &SwitchStatement{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litA,
									Tokens:                tk[2:3],
								},
							},
							Tokens: tk[2:3],
						},
						CaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litB,
											Tokens:                tk[7:8],
										},
									},
									Tokens: tk[7:8],
								},
								Tokens: tk[5:9],
							},
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litC,
											Tokens:                tk[11:12],
										},
									},
									Tokens: tk[11:12],
								},
								StatementList: []StatementListItem{
									{
										Statement: &Statement{
											ExpressionStatement: &Expression{
												Expressions: []AssignmentExpression{
													{
														ConditionalExpression: &litD,
														Tokens:                tk[13:14],
													},
												},
												Tokens: tk[13:14],
											},
											Tokens: tk[13:15],
										},
										Tokens: tk[13:15],
									},
								},
								Tokens: tk[9:15],
							},
						},
						DefaultClause: []StatementListItem{
							{
								Statement: &Statement{
									ExpressionStatement: &Expression{
										Expressions: []AssignmentExpression{
											{
												ConditionalExpression: &litE,
												Tokens:                tk[17:18],
											},
										},
										Tokens: tk[17:18],
									},
									Tokens: tk[17:19],
								},
								Tokens: tk[17:19],
							},
							{
								Statement: &Statement{
									ExpressionStatement: &Expression{
										Expressions: []AssignmentExpression{
											{
												ConditionalExpression: &litF,
												Tokens:                tk[19:20],
											},
										},
										Tokens: tk[19:20],
									},
									Tokens: tk[19:21],
								},
								Tokens: tk[19:21],
							},
						},
						PostDefaultCaseClauses: []CaseClause{
							{
								Expression: Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litG,
											Tokens:                tk[23:24],
										},
									},
									Tokens: tk[23:24],
								},
								StatementList: []StatementListItem{
									{
										Statement: &Statement{
											ExpressionStatement: &Expression{
												Expressions: []AssignmentExpression{
													{
														ConditionalExpression: &litH,
														Tokens:                tk[25:26],
													},
												},
												Tokens: tk[25:26],
											},
											Tokens: tk[25:27],
										},
										Tokens: tk[25:27],
									},
								},
								Tokens: tk[21:27],
							},
						},
						Tokens: tk[:28],
					},
					Tokens: tk[:28],
				},
				Tokens: tk[:28],
			}
		}},
		{`continue;`, func(t *test, tk Tokens) { // 33
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementContinue,
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`continue ;`, func(t *test, tk Tokens) { // 34
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementContinue,
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`continue Name;`, func(t *test, tk Tokens) { // 35
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:            StatementContinue,
					LabelIdentifier: &tk[2],
					Tokens:          tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`break;`, func(t *test, tk Tokens) { // 36
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementBreak,
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`break ;`, func(t *test, tk Tokens) { // 37
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementBreak,
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{`break Name;`, func(t *test, tk Tokens) { // 38
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:            StatementBreak,
					LabelIdentifier: &tk[2],
					Tokens:          tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`return;`, func(t *test, tk Tokens) { // 39
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidStatement,
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{`return;`, func(t *test, tk Tokens) { // 40
			t.Ret = true
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementReturn,
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`return 1;`, func(t *test, tk Tokens) { // 41
			t.Ret = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = StatementListItem{
				Statement: &Statement{
					Type: StatementReturn,
					ExpressionStatement: &Expression{
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
		{`with(a){}`, func(t *test, tk Tokens) { // 42
			litA := makeConditionLiteral(tk, 2)
			t.Output = StatementListItem{
				Statement: &Statement{
					WithStatement: &WithStatement{
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
				},
				Tokens: tk[:6],
			}
		}},
		{`throw a;`, func(t *test, tk Tokens) { // 43
			litA := makeConditionLiteral(tk, 2)
			t.Output = StatementListItem{
				Statement: &Statement{
					Type: StatementThrow,
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[2:3],
							},
						},
						Tokens: tk[2:3],
					},
					Tokens: tk[0:4],
				},
				Tokens: tk[0:4],
			}
		}},
		{`try{a;}catch(e){b;}`, func(t *test, tk Tokens) { // 44
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						CatchParameterBindingIdentifier: &tk[7],
						CatchBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[10:11],
												},
											},
											Tokens: tk[10:11],
										},
										Tokens: tk[10:12],
									},
									Tokens: tk[10:12],
								},
							},
							Tokens: tk[9:13],
						},
						Tokens: tk[:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`try{a;}catch({e}){b;}`, func(t *test, tk Tokens) { // 45
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 12)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						CatchParameterObjectBindingPattern: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									SingleNameBinding: &tk[8],
									Tokens:            tk[8:9],
								},
							},
							Tokens: tk[7:10],
						},
						CatchBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[12:13],
												},
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:14],
									},
									Tokens: tk[12:14],
								},
							},
							Tokens: tk[11:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`try{a;}finally{b;}`, func(t *test, tk Tokens) { // 46
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						FinallyBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[7:8],
												},
											},
											Tokens: tk[7:8],
										},
										Tokens: tk[7:9],
									},
									Tokens: tk[7:9],
								},
							},
							Tokens: tk[6:10],
						},
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`try{a;}catch([e]){b;}finally{c;}`, func(t *test, tk Tokens) { // 47
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 12)
			litC := makeConditionLiteral(tk, 17)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						CatchParameterArrayBindingPattern: &ArrayBindingPattern{
							BindingElementList: []BindingElement{
								{
									SingleNameBinding: &tk[8],
									Tokens:            tk[8:9],
								},
							},
							Tokens: tk[7:10],
						},
						CatchBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[12:13],
												},
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:14],
									},
									Tokens: tk[12:14],
								},
							},
							Tokens: tk[11:15],
						},
						FinallyBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litC,
													Tokens:                tk[17:18],
												},
											},
											Tokens: tk[17:18],
										},
										Tokens: tk[17:19],
									},
									Tokens: tk[17:19],
								},
							},
							Tokens: tk[16:20],
						},
						Tokens: tk[:20],
					},
					Tokens: tk[:20],
				},
				Tokens: tk[:20],
			}
		}},
		{`debugger;`, func(t *test, tk Tokens) { // 48
			t.Output = StatementListItem{
				Statement: &Statement{
					Type:   StatementDebugger,
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{`label: debugger;`, func(t *test, tk Tokens) { // 49
			t.Output = StatementListItem{
				Statement: &Statement{
					LabelIdentifier: &tk[0],
					LabelledItemStatement: &Statement{
						Type:   StatementDebugger,
						Tokens: tk[3:5],
					},
					Tokens: tk[0:5],
				},
				Tokens: tk[0:5],
			}
		}},
		{`label: function fn(){}`, func(t *test, tk Tokens) { // 50
			t.Output = StatementListItem{
				Statement: &Statement{
					LabelIdentifier: &tk[0],
					LabelledItemFunction: &FunctionDeclaration{
						BindingIdentifier: &tk[5],
						FormalParameters: FormalParameters{
							Tokens: tk[6:8],
						},
						FunctionBody: Block{
							Tokens: tk[8:10],
						},
						Tokens: tk[3:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`a = b;`, func(t *test, tk Tokens) { // 51
			litB := makeConditionLiteral(tk, 4)
			t.Output = StatementListItem{
				Statement: &Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
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
									ConditionalExpression: &litB,
									Tokens:                tk[4:5],
								},
								Tokens: tk[:5],
							},
						},
						Tokens: tk[:5],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{`class a{}`, func(t *test, tk Tokens) { // 52
			t.Output = StatementListItem{
				Declaration: &Declaration{
					ClassDeclaration: &ClassDeclaration{
						BindingIdentifier: &tk[2],
						Tokens:            tk[:5],
					},
					Tokens: tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{`function a(){}`, func(t *test, tk Tokens) { // 53
			t.Output = StatementListItem{
				Declaration: &Declaration{
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
				},
				Tokens: tk[:7],
			}
		}},
		{`let a;`, func(t *test, tk Tokens) { // 54
			t.Output = StatementListItem{
				Declaration: &Declaration{
					LexicalDeclaration: &LexicalDeclaration{
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &tk[2],
								Tokens:            tk[2:3],
							},
						},
						Tokens: tk[:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{`const a = 1;`, func(t *test, tk Tokens) { // 55
			litA := makeConditionLiteral(tk, 6)
			t.Output = StatementListItem{
				Declaration: &Declaration{
					LexicalDeclaration: &LexicalDeclaration{
						LetOrConst: Const,
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &tk[2],
								Initializer: &AssignmentExpression{
									ConditionalExpression: &litA,
									Tokens:                tk[6:7],
								},
								Tokens: tk[2:7],
							},
						},
						Tokens: tk[:8],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{`async function fn(){}`, func(t *test, tk Tokens) { // 56
			t.Output = StatementListItem{
				Declaration: &Declaration{
					FunctionDeclaration: &FunctionDeclaration{
						Type:              FunctionAsync,
						BindingIdentifier: &tk[4],
						FormalParameters: FormalParameters{
							Tokens: tk[5:7],
						},
						FunctionBody: Block{
							Tokens: tk[7:9],
						},
						Tokens: tk[:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{`async () => {};`, func(t *test, tk Tokens) { // 57
			t.Output = StatementListItem{
				Statement: &Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ArrowFunction: &ArrowFunction{
									Async: true,
									FormalParameters: &FormalParameters{
										Tokens: tk[2:4],
									},
									FunctionBody: &Block{
										Tokens: tk[7:9],
									},
									Tokens: tk[:9],
								},
								Tokens: tk[:9],
							},
						},
						Tokens: tk[:9],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`for await(a of b) {}`, func(t *test, tk Tokens) { // 58
			litB := makeConditionLiteral(tk, 8)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForAwaitOfLeftHandSide,
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[4],
										Tokens:              tk[4:5],
									},
									Tokens: tk[4:5],
								},
								Tokens: tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[8:9],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[11:13],
							},
							Tokens: tk[11:13],
						},
						Tokens: tk[:13],
					},
					Tokens: tk[:13],
				},
				Tokens: tk[:13],
			}
		}},
		{`for await(var a of b) {}`, func(t *test, tk Tokens) { // 59
			litB := makeConditionLiteral(tk, 10)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type:                 ForAwaitOfVar,
						ForBindingIdentifier: &tk[6],
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[13:15],
							},
							Tokens: tk[13:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
		{`for await(let {a} of b) {}`, func(t *test, tk Tokens) { // 60
			litB := makeConditionLiteral(tk, 12)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForAwaitOfLet,
						ForBindingPatternObject: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									SingleNameBinding: &tk[7],
									Tokens:            tk[7:8],
								},
							},
							Tokens: tk[6:9],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[12:13],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[15:17],
						},
						Tokens: tk[:17],
					},
					Tokens: tk[:17],
				},
				Tokens: tk[:17],
			}
		}},
		{`for await(const [a] of b) {}`, func(t *test, tk Tokens) { // 61
			litB := makeConditionLiteral(tk, 12)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForAwaitOfConst,
						ForBindingPatternArray: &ArrayBindingPattern{
							BindingElementList: []BindingElement{
								{
									SingleNameBinding: &tk[7],
									Tokens:            tk[7:8],
								},
							},
							Tokens: tk[6:9],
						},
						Of: &AssignmentExpression{
							ConditionalExpression: &litB,
							Tokens:                tk[12:13],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[15:17],
						},
						Tokens: tk[:17],
					},
					Tokens: tk[:17],
				},
				Tokens: tk[:17],
			}
		}},
		{`try{a;}catch{b;}`, func(t *test, tk Tokens) { // 62
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						CatchBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[7:8],
												},
											},
											Tokens: tk[7:8],
										},
										Tokens: tk[7:9],
									},
									Tokens: tk[7:9],
								},
							},
							Tokens: tk[6:10],
						},
						Tokens: tk[:10],
					},
					Tokens: tk[:10],
				},
				Tokens: tk[:10],
			}
		}},
		{`try{a;}catch{b;}finally{c;}`, func(t *test, tk Tokens) { // 63
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 7)
			litC := makeConditionLiteral(tk, 12)
			t.Output = StatementListItem{
				Statement: &Statement{
					TryStatement: &TryStatement{
						TryBlock: Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litA,
													Tokens:                tk[2:3],
												},
											},
											Tokens: tk[2:3],
										},
										Tokens: tk[2:4],
									},
									Tokens: tk[2:4],
								},
							},
							Tokens: tk[1:5],
						},
						CatchBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litB,
													Tokens:                tk[7:8],
												},
											},
											Tokens: tk[7:8],
										},
										Tokens: tk[7:9],
									},
									Tokens: tk[7:9],
								},
							},
							Tokens: tk[6:10],
						},
						FinallyBlock: &Block{
							StatementList: []StatementListItem{
								{
									Statement: &Statement{
										ExpressionStatement: &Expression{
											Expressions: []AssignmentExpression{
												{
													ConditionalExpression: &litC,
													Tokens:                tk[12:13],
												},
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:14],
									},
									Tokens: tk[12:14],
								},
							},
							Tokens: tk[11:15],
						},
						Tokens: tk[:15],
					},
					Tokens: tk[:15],
				},
				Tokens: tk[:15],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var sl StatementListItem
		err := sl.parse(&t.Tokens, t.Yield, t.Await, t.Ret)
		return sl, err
	})
}

func TestBlock(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "Block",
				Token:   tk[0],
			}
		}},
		{"{\n}", func(t *test, tk Tokens) { // 2
			t.Output = Block{
				Tokens: tk[:3],
			}
		}},
		{"{\n,\n}", func(t *test, tk Tokens) { //3
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     assignmentError(tk[2]),
							Parsing: "Expression",
							Token:   tk[2],
						},
						Parsing: "Statement",
						Token:   tk[2],
					},
					Parsing: "StatementListItem",
					Token:   tk[2],
				},
				Parsing: "Block",
				Token:   tk[2],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 2)
			t.Output = Block{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litA,
										Tokens:                tk[2:3],
									},
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"{\na\nfunction}", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidStatement,
						Parsing: "Statement",
						Token:   tk[4],
					},
					Parsing: "StatementListItem",
					Token:   tk[4],
				},
				Parsing: "Block",
				Token:   tk[4],
			}
		}},
		{"{\na\nb\n}", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 4)
			t.Output = Block{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litA,
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
										ConditionalExpression: &litB,
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
			}
		}},
	}, func(t *test) (interface{}, error) {
		var b Block
		err := b.parse(&t.Tokens, t.Yield, t.Await, t.Ret)
		return b, err
	})
}

func TestStatementListItem(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[0]),
						Parsing: "Expression",
						Token:   tk[0],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"let", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
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
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"const", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
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
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"class", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidStatement,
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"class\na", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningBrace,
						Parsing: "ClassDeclaration",
						Token:   tk[3],
					},
					Parsing: "Declaration",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"async", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "ArrowFunction",
								Token:   tk[1],
							},
							Parsing: "AssignmentExpression",
							Token:   tk[0],
						},
						Parsing: "Expression",
						Token:   tk[0],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"async function\na", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningParenthesis,
							Parsing: "FormalParameters",
							Token:   tk[5],
						},
						Parsing: "FunctionDeclaration",
						Token:   tk[5],
					},
					Parsing: "Declaration",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"async function\n*\na", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningParenthesis,
							Parsing: "FormalParameters",
							Token:   tk[7],
						},
						Parsing: "FunctionDeclaration",
						Token:   tk[7],
					},
					Parsing: "Declaration",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"async function\n*", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidStatement,
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"async\nfunction\na", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "ArrowFunction",
								Token:   tk[1],
							},
							Parsing: "AssignmentExpression",
							Token:   tk[0],
						},
						Parsing: "Expression",
						Token:   tk[0],
					},
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"function", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidStatement,
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"function\n*", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidStatement,
					Parsing: "Statement",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"function\n*\na", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningParenthesis,
							Parsing: "FormalParameters",
							Token:   tk[5],
						},
						Parsing: "FunctionDeclaration",
						Token:   tk[5],
					},
					Parsing: "Declaration",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"function\na", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningParenthesis,
							Parsing: "FormalParameters",
							Token:   tk[3],
						},
						Parsing: "FunctionDeclaration",
						Token:   tk[3],
					},
					Parsing: "Declaration",
					Token:   tk[0],
				},
				Parsing: "StatementListItem",
				Token:   tk[0],
			}
		}},
		{"function\na(){}", func(t *test, tk Tokens) { //15
			t.Output = StatementListItem{
				Declaration: &Declaration{
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
				},
				Tokens: tk[:7],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 0)
			t.Output = StatementListItem{
				Statement: &Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[:1],
							},
						},
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:1],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var si StatementListItem
		err := si.parse(&t.Tokens, t.Yield, t.Await, t.Ret)
		return si, err
	})
}
