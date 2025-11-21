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
									ElementList: []ArrayElement{
										{
											AssignmentExpression: AssignmentExpression{
												ConditionalExpression: &litA,
												Tokens:                tk[12:13],
											},
											Tokens: tk[12:13],
										},
										{
											AssignmentExpression: AssignmentExpression{
												ConditionalExpression: &litB,
												Tokens:                tk[15:16],
											},
											Tokens: tk[15:16],
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
			litC := makeConditionLiteral(tk, 12)
			litD := makeConditionLiteral(tk, 15)
			obj := wrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								ObjectLiteral: &ObjectLiteral{
									PropertyDefinitionList: []PropertyDefinition{
										{
											PropertyName: &PropertyName{
												LiteralPropertyName: &tk[12],
												Tokens:              tk[12:13],
											},
											AssignmentExpression: &AssignmentExpression{
												ConditionalExpression: &litC,
												Tokens:                tk[12:13],
											},
											Tokens: tk[12:13],
										},
										{
											PropertyName: &PropertyName{
												LiteralPropertyName: &tk[15],
												Tokens:              tk[15:16],
											},
											AssignmentExpression: &AssignmentExpression{
												ConditionalExpression: &litD,
												Tokens:                tk[15:16],
											},
											Tokens: tk[15:16],
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
											PropertyName: PropertyName{
												LiteralPropertyName: &tk[3],
												Tokens:              tk[3:4],
											},
											BindingElement: BindingElement{
												SingleNameBinding: &tk[3],
												Tokens:            tk[3:4],
											},
											Tokens: tk[3:4],
										},
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
			litC := makeConditionLiteral(tk, 16)
			litD := makeConditionLiteral(tk, 20)
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
											Tokens:               tk[11:17],
										}).LogicalORExpression,
										Tokens: tk[11:17],
									},
									Tokens: tk[11:17],
								},
							},
							Tokens: tk[11:17],
						},
						Afterthought: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &ConditionalExpression{
										LogicalORExpression: wrapConditional(UpdateExpression{
											UpdateOperator:  UpdatePreIncrement,
											UnaryExpression: &litD.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
											Tokens:          tk[19:21],
										}).LogicalORExpression,
										Tokens: tk[19:21],
									},
									Tokens: tk[19:21],
								},
							},
							Tokens: tk[19:21],
						},
						Statement: Statement{
							BlockStatement: &Block{
								Tokens: tk[23:25],
							},
							Tokens: tk[23:25],
						},
						Tokens: tk[:25],
					},
					Tokens: tk[:25],
				},
				Tokens: tk[:25],
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
									PropertyName: PropertyName{
										LiteralPropertyName: &tk[8],
										Tokens:              tk[8:9],
									},
									BindingElement: BindingElement{
										SingleNameBinding: &tk[8],
										Tokens:            tk[8:9],
									},
									Tokens: tk[8:9],
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
			t.Await = true
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
			t.Await = true
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
			t.Await = true
			litB := makeConditionLiteral(tk, 12)
			t.Output = StatementListItem{
				Statement: &Statement{
					IterationStatementFor: &IterationStatementFor{
						Type: ForAwaitOfLet,
						ForBindingPatternObject: &ObjectBindingPattern{
							BindingPropertyList: []BindingProperty{
								{
									PropertyName: PropertyName{
										LiteralPropertyName: &tk[7],
										Tokens:              tk[7:8],
									},
									BindingElement: BindingElement{
										SingleNameBinding: &tk[7],
										Tokens:            tk[7:8],
									},
									Tokens: tk[7:8],
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
			t.Await = true
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
	}, func(t *test) (Type, error) {
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
		{"{\n,\n}", func(t *test, tk Tokens) { // 3
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
		{"{ // A\n\n// B\na // C\n\n// D\n}", func(t *test, tk Tokens) { // 7
			t.Output = Block{
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: WrapConditional(&MemberExpression{
											PrimaryExpression: &PrimaryExpression{
												IdentifierReference: &tk[6],
												Tokens:              tk[6:7],
											},
											Comments: [5]Comments{nil, nil, nil, nil, {&tk[8]}},
											Tokens:   tk[6:9],
										}),
										Tokens: tk[6:9],
									},
								},
								Tokens: tk[6:9],
							},
							Tokens: tk[6:9],
						},
						Comments: [2]Comments{{&tk[4]}},
						Tokens:   tk[4:9],
					},
				},
				Comments: [2]Comments{{&tk[2]}, {&tk[10]}},
				Tokens:   tk[:13],
			}
		}},
	}, func(t *test) (Type, error) {
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
						Err: assignmentCustomError(tk[0], Error{
							Err:     ErrInvalidFunction,
							Parsing: "FunctionDeclaration",
							Token:   tk[1],
						}),
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
						Err: assignmentCustomError(tk[0], Error{
							Err:     ErrInvalidFunction,
							Parsing: "FunctionDeclaration",
							Token:   tk[1],
						}),
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
		{"function\na(){}", func(t *test, tk Tokens) { // 15
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
		{"// A\n// B\na // C\n\n// D", func(t *test, tk Tokens) { // 17
			t.Output = StatementListItem{
				Statement: &Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[4],
										Tokens:              tk[4:5],
									},
									Comments: [5]Comments{nil, nil, nil, nil, {&tk[6]}},
									Tokens:   tk[4:7],
								}),
								Tokens: tk[4:7],
							},
						},
						Tokens: tk[4:7],
					},
					Tokens: tk[4:7],
				},
				Comments: [2]Comments{{&tk[0], &tk[2]}},
				Tokens:   tk[:7],
			}
		}},
		{"// A\nclass a{} // B\n", func(t *test, tk Tokens) { // 18
			t.Output = StatementListItem{
				Declaration: &Declaration{
					ClassDeclaration: &ClassDeclaration{
						BindingIdentifier: &tk[4],
						Tokens:            tk[2:7],
					},
					Tokens: tk[2:7],
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[8]}},
				Tokens:   tk[:9],
			}
		}},
		{"// A\nfunction a(){} // B\n", func(t *test, tk Tokens) { // 19
			t.Output = StatementListItem{
				Declaration: &Declaration{
					FunctionDeclaration: &FunctionDeclaration{
						BindingIdentifier: &tk[4],
						FormalParameters: FormalParameters{
							Tokens: tk[5:7],
						},
						FunctionBody: Block{
							Tokens: tk[7:9],
						},
						Tokens: tk[2:9],
					},
					Tokens: tk[2:9],
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[10]}},
				Tokens:   tk[:11],
			}
		}},
		{"// A\nconst a = 1; // B\n", func(t *test, tk Tokens) { // 20
			t.Output = StatementListItem{
				Declaration: &Declaration{
					LexicalDeclaration: &LexicalDeclaration{
						LetOrConst: Const,
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &tk[4],
								Initializer: &AssignmentExpression{
									ConditionalExpression: WrapConditional(&PrimaryExpression{
										Literal: &tk[8],
										Tokens:  tk[8:9],
									}),
									Tokens: tk[8:9],
								},
								Tokens: tk[4:9],
							},
						},
						Tokens: tk[2:10],
					},
					Tokens: tk[2:10],
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[11]}},
				Tokens:   tk[:12],
			}
		}},
		{"// A\nlet a = 1; // B\n", func(t *test, tk Tokens) { // 21
			t.Output = StatementListItem{
				Declaration: &Declaration{
					LexicalDeclaration: &LexicalDeclaration{
						LetOrConst: Let,
						BindingList: []LexicalBinding{
							{
								BindingIdentifier: &tk[4],
								Initializer: &AssignmentExpression{
									ConditionalExpression: WrapConditional(&PrimaryExpression{
										Literal: &tk[8],
										Tokens:  tk[8:9],
									}),
									Tokens: tk[8:9],
								},
								Tokens: tk[4:9],
							},
						},
						Tokens: tk[2:10],
					},
					Tokens: tk[2:10],
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[11]}},
				Tokens:   tk[:12],
			}
		}},
	}, func(t *test) (Type, error) {
		var si StatementListItem

		err := si.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return si, err
	})
}

func TestStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[0]),
					Parsing: "Expression",
					Token:   tk[0],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"{,}", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     assignmentError(tk[1]),
								Parsing: "Expression",
								Token:   tk[1],
							},
							Parsing: "Statement",
							Token:   tk[1],
						},
						Parsing: "StatementListItem",
						Token:   tk[1],
					},
					Parsing: "Block",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"{}", func(t *test, tk Tokens) { // 3
			t.Output = Statement{
				BlockStatement: &Block{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
		{"var", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[1],
					},
					Parsing: "VariableStatement",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"var\na", func(t *test, tk Tokens) { // 5
			t.Output = Statement{
				VariableStatement: &VariableStatement{
					VariableDeclarationList: []VariableDeclaration{
						{
							BindingIdentifier: &tk[2],
							Tokens:            tk[2:3],
						},
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{";", func(t *test, tk Tokens) { // 6
			t.Output = Statement{
				Tokens: tk[:1],
			}
		}},
		{"if", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "IfStatement",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"if(a){}", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
				IfStatement: &IfStatement{
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
		{"do", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     assignmentError(tk[1]),
							Parsing: "Expression",
							Token:   tk[1],
						},
						Parsing: "Statement",
						Token:   tk[1],
					},
					Parsing: "IterationStatementDo",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"do {} while(1)", func(t *test, tk Tokens) { // 10
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
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"while", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "IterationStatementWhile",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"while (a) {}", func(t *test, tk Tokens) { // 12
			litA := makeConditionLiteral(tk, 3)
			t.Output = Statement{
				IterationStatementWhile: &IterationStatementWhile{
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
		{"for", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "IterationStatementFor",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"for (;;) {}", func(t *test, tk Tokens) { // 14
			t.Output = Statement{
				IterationStatementFor: &IterationStatementFor{
					Statement: Statement{
						BlockStatement: &Block{
							Tokens: tk[7:9],
						},
						Tokens: tk[7:9],
					},
					Tokens: tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"switch", func(t *test, tk Tokens) { // 15
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "SwitchStatement",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"switch (a) {}", func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 3)
			t.Output = Statement{
				SwitchStatement: &SwitchStatement{
					Expression: Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[3:4],
							},
						},
						Tokens: tk[3:4],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"continue", func(t *test, tk Tokens) { // 17
			t.Output = Statement{
				Type:   StatementContinue,
				Tokens: tk[:1],
			}
		}},
		{"continue;", func(t *test, tk Tokens) { // 18
			t.Output = Statement{
				Type:   StatementContinue,
				Tokens: tk[:2],
			}
		}},
		{"continue\nswitch", func(t *test, tk Tokens) { // 19
			t.Output = Statement{
				Type:   StatementContinue,
				Tokens: tk[:1],
			}
		}},
		{"continue switch", func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "Statement",
				Token:   tk[2],
			}
		}},
		{"continue a", func(t *test, tk Tokens) { // 21
			t.Output = Statement{
				Type:            StatementContinue,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:3],
			}
		}},
		{"continue a\n", func(t *test, tk Tokens) { // 22
			t.Output = Statement{
				Type:            StatementContinue,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:3],
			}
		}},
		{"continue a;", func(t *test, tk Tokens) { // 23
			t.Output = Statement{
				Type:            StatementContinue,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:4],
			}
		}},
		{"continue a b", func(t *test, tk Tokens) { // 24
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[3],
			}
		}},
		{"break", func(t *test, tk Tokens) { // 25
			t.Output = Statement{
				Type:   StatementBreak,
				Tokens: tk[:1],
			}
		}},
		{"break;", func(t *test, tk Tokens) { // 26
			t.Output = Statement{
				Type:   StatementBreak,
				Tokens: tk[:2],
			}
		}},
		{"break\nswitch", func(t *test, tk Tokens) { // 27
			t.Output = Statement{
				Type:   StatementBreak,
				Tokens: tk[:1],
			}
		}},
		{"break switch", func(t *test, tk Tokens) { // 28
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "Statement",
				Token:   tk[2],
			}
		}},
		{"break a", func(t *test, tk Tokens) { // 29
			t.Output = Statement{
				Type:            StatementBreak,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:3],
			}
		}},
		{"break a\n", func(t *test, tk Tokens) { // 30
			t.Output = Statement{
				Type:            StatementBreak,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:3],
			}
		}},
		{"break a;", func(t *test, tk Tokens) { // 31
			t.Output = Statement{
				Type:            StatementBreak,
				LabelIdentifier: &tk[2],
				Tokens:          tk[:4],
			}
		}},
		{"break a b", func(t *test, tk Tokens) { // 32
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[3],
			}
		}},
		{"return", func(t *test, tk Tokens) { // 33
			t.Err = Error{
				Err:     ErrInvalidStatement,
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"return", func(t *test, tk Tokens) { // 34
			t.Ret = true
			t.Output = Statement{
				Type:   StatementReturn,
				Tokens: tk[:1],
			}
		}},
		{"return;", func(t *test, tk Tokens) { // 35
			t.Ret = true
			t.Output = Statement{
				Type:   StatementReturn,
				Tokens: tk[:2],
			}
		}},
		{"return\na", func(t *test, tk Tokens) { // 36
			t.Ret = true
			t.Output = Statement{
				Type:   StatementReturn,
				Tokens: tk[:1],
			}
		}},
		{"return ,", func(t *test, tk Tokens) { // 37
			t.Ret = true
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[2]),
					Parsing: "Expression",
					Token:   tk[2],
				},
				Parsing: "Statement",
				Token:   tk[2],
			}
		}},
		{"return a", func(t *test, tk Tokens) { // 38
			t.Ret = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
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
				Tokens: tk[:3],
			}
		}},
		{"return a;", func(t *test, tk Tokens) { // 39
			t.Ret = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
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
			}
		}},
		{"return a b", func(t *test, tk Tokens) { // 40
			t.Ret = true
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[3],
			}
		}},
		{"with", func(t *test, tk Tokens) { // 41
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningParenthesis,
					Parsing: "WithStatement",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"with (a) {}", func(t *test, tk Tokens) { // 42
			litA := makeConditionLiteral(tk, 3)
			t.Output = Statement{
				WithStatement: &WithStatement{
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
		{"throw", func(t *test, tk Tokens) { // 43
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[1]),
					Parsing: "Expression",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[1],
			}
		}},
		{"throw\na", func(t *test, tk Tokens) { // 44
			t.Err = Error{
				Err:     ErrUnexpectedLineTerminator,
				Parsing: "Statement",
				Token:   tk[1],
			}
		}},
		{"throw a", func(t *test, tk Tokens) { // 45
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
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
				Tokens: tk[:3],
			}
		}},
		{"throw a;", func(t *test, tk Tokens) { // 46
			litA := makeConditionLiteral(tk, 2)
			t.Output = Statement{
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
				Tokens: tk[:4],
			}
		}},
		{"throw a b", func(t *test, tk Tokens) { // 47
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[3],
			}
		}},
		{"try", func(t *test, tk Tokens) { // 48
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningBrace,
						Parsing: "Block",
						Token:   tk[1],
					},
					Parsing: "TryStatement",
					Token:   tk[1],
				},
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"try{}finally{}", func(t *test, tk Tokens) { // 49
			t.Output = Statement{
				TryStatement: &TryStatement{
					TryBlock: Block{
						Tokens: tk[1:3],
					},
					FinallyBlock: &Block{
						Tokens: tk[4:6],
					},
					Tokens: tk[:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"debugger", func(t *test, tk Tokens) { // 50
			t.Output = Statement{
				Type:   StatementDebugger,
				Tokens: tk[:1],
			}
		}},
		{"debugger;", func(t *test, tk Tokens) { // 51
			t.Output = Statement{
				Type:   StatementDebugger,
				Tokens: tk[:2],
			}
		}},
		{"debugger a", func(t *test, tk Tokens) { // 52
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[1],
			}
		}},
		{"debugger\na", func(t *test, tk Tokens) { // 53
			t.Output = Statement{
				Type:   StatementDebugger,
				Tokens: tk[:1],
			}
		}},
		{"a\n:\nfunction", func(t *test, tk Tokens) { // 54
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "FunctionDeclaration",
					Token:   tk[5],
				},
				Parsing: "Statement",
				Token:   tk[4],
			}
		}},
		{"a\n:\nfunction\nb(){}", func(t *test, tk Tokens) { // 55
			t.Output = Statement{
				LabelIdentifier: &tk[0],
				LabelledItemFunction: &FunctionDeclaration{
					BindingIdentifier: &tk[6],
					FormalParameters: FormalParameters{
						Tokens: tk[7:9],
					},
					FunctionBody: Block{
						Tokens: tk[9:11],
					},
					Tokens: tk[4:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"a\n:", func(t *test, tk Tokens) { // 56
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[3]),
						Parsing: "Expression",
						Token:   tk[3],
					},
					Parsing: "Statement",
					Token:   tk[3],
				},
				Parsing: "Statement",
				Token:   tk[3],
			}
		}},
		{"function", func(t *test, tk Tokens) { // 57
			t.Err = Error{
				Err:     ErrInvalidStatement,
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"class", func(t *test, tk Tokens) { // 58
			t.Err = Error{
				Err:     ErrInvalidStatement,
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"async function", func(t *test, tk Tokens) { // 59
			t.Err = Error{
				Err:     ErrInvalidStatement,
				Parsing: "Statement",
				Token:   tk[0],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 60
			litA := makeConditionLiteral(tk, 0)
			t.Output = Statement{
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
			}
		}},
		{"a;", func(t *test, tk Tokens) { // 61
			litA := makeConditionLiteral(tk, 0)
			t.Output = Statement{
				ExpressionStatement: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[:1],
						},
					},
					Tokens: tk[:1],
				},
				Tokens: tk[:2],
			}
		}},
		{"a b", func(t *test, tk Tokens) { // 62
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "Statement",
				Token:   tk[1],
			}
		}},
		{"a\nb", func(t *test, tk Tokens) { // 63
			litA := makeConditionLiteral(tk, 0)
			t.Output = Statement{
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
			}
		}},
		{"continue /* A */ a // B\n", func(t *test, tk Tokens) { // 64
			t.Output = Statement{
				Type:            StatementContinue,
				LabelIdentifier: &tk[4],
				Comments:        [2]Comments{{&tk[2]}, {&tk[6]}},
				Tokens:          tk[:7],
			}
		}},
		{"continue /* A */ a // B\n;", func(t *test, tk Tokens) { // 65
			t.Output = Statement{
				Type:            StatementContinue,
				LabelIdentifier: &tk[4],
				Comments:        [2]Comments{{&tk[2]}, {&tk[6]}},
				Tokens:          tk[:9],
			}
		}},
		{"continue /* A */;", func(t *test, tk Tokens) { // 66
			t.Output = Statement{
				Type:     StatementContinue,
				Comments: [2]Comments{{&tk[2]}},
				Tokens:   tk[:4],
			}
		}},
		{"break /* A */ a // B\n", func(t *test, tk Tokens) { // 67
			t.Output = Statement{
				Type:            StatementBreak,
				LabelIdentifier: &tk[4],
				Comments:        [2]Comments{{&tk[2]}, {&tk[6]}},
				Tokens:          tk[:7],
			}
		}},
		{"break /* A */ a // B\n;", func(t *test, tk Tokens) { // 68
			t.Output = Statement{
				Type:            StatementBreak,
				LabelIdentifier: &tk[4],
				Comments:        [2]Comments{{&tk[2]}, {&tk[6]}},
				Tokens:          tk[:9],
			}
		}},
		{"break /* A */;", func(t *test, tk Tokens) { // 69
			t.Output = Statement{
				Type:     StatementBreak,
				Comments: [2]Comments{{&tk[2]}},
				Tokens:   tk[:4],
			}
		}},
		{"return /* A */ a // B\n", func(t *test, tk Tokens) { // 70
			t.Ret = true
			t.Output = Statement{
				Type: StatementReturn,
				ExpressionStatement: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[4],
									Tokens:              tk[4:5],
								},
								Comments: [5]Comments{{&tk[2]}, nil, nil, nil, {&tk[6]}},
								Tokens:   tk[2:7],
							}),
							Tokens: tk[2:7],
						},
					},
					Tokens: tk[2:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"return /* A */ a // B\n;", func(t *test, tk Tokens) { // 71
			t.Ret = true
			t.Output = Statement{
				Type: StatementReturn,
				ExpressionStatement: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[4],
									Tokens:              tk[4:5],
								},
								Comments: [5]Comments{{&tk[2]}, nil, nil, nil, {&tk[6]}},
								Tokens:   tk[2:7],
							}),
							Tokens: tk[2:7],
						},
					},
					Tokens: tk[2:7],
				},
				Tokens: tk[:9],
			}
		}},
		{"return /* A */;", func(t *test, tk Tokens) { // 72
			t.Ret = true
			t.Output = Statement{
				Type:     StatementReturn,
				Comments: [2]Comments{{&tk[2]}},
				Tokens:   tk[:4],
			}
		}},
		{"debugger // A\n", func(t *test, tk Tokens) { // 73
			t.Output = Statement{
				Type:     StatementDebugger,
				Comments: [2]Comments{{&tk[2]}},
				Tokens:   tk[:3],
			}
		}},
		{"debugger /* A */;", func(t *test, tk Tokens) { // 74
			t.Output = Statement{
				Type:     StatementDebugger,
				Comments: [2]Comments{{&tk[2]}},
				Tokens:   tk[:4],
			}
		}},
		{"a /* A */: // B\ndebugger;", func(t *test, tk Tokens) { // 75
			t.Output = Statement{
				LabelIdentifier: &tk[0],
				LabelledItemStatement: &Statement{
					Type:   StatementDebugger,
					Tokens: tk[7:9],
				},
				Comments: [2]Comments{{&tk[2]}, {&tk[5]}},
				Tokens:   tk[:9],
			}
		}},
	}, func(t *test) (Type, error) {
		var s Statement

		err := s.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return s, err
	})
}

func TestIfStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidIfStatement,
				Parsing: "IfStatement",
				Token:   tk[0],
			}
		}},
		{`if`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IfStatement",
				Token:   tk[1],
			}
		}},
		{"if\n(\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "IfStatement",
				Token:   tk[4],
			}
		}},
		{"if\n(\na\nb\n)", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "IfStatement",
				Token:   tk[6],
			}
		}},
		{"if\n(\na\n)\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[8]),
						Parsing: "Expression",
						Token:   tk[8],
					},
					Parsing: "Statement",
					Token:   tk[8],
				},
				Parsing: "IfStatement",
				Token:   tk[8],
			}
		}},
		{"if\n(\na\n)\nb", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 8)
			t.Output = IfStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"if\n(\na\n)\nb\nelse", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[11]),
						Parsing: "Expression",
						Token:   tk[11],
					},
					Parsing: "Statement",
					Token:   tk[11],
				},
				Parsing: "IfStatement",
				Token:   tk[11],
			}
		}},
		{"if\n(\na\n)\nb\nelse\nc", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 8)
			litC := makeConditionLiteral(tk, 12)
			t.Output = IfStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					Tokens: tk[8:9],
				},
				ElseStatement: &Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litC,
								Tokens:                tk[12:13],
							},
						},
						Tokens: tk[12:13],
					},
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"if (a) b: function c(){}", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "IfStatement",
				Token:   tk[6],
			}
		}},
		{"if (a){b}else c: function d(){}", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "IfStatement",
				Token:   tk[10],
			}
		}},
		{"if // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G", func(t *test, tk Tokens) { // 11
			t.Output = IfStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								},
								Comments: [5]Comments{{&tk[8]}, nil, nil, nil, {&tk[12]}},
								Tokens:   tk[8:13],
							}),
							Tokens: tk[8:13],
						},
					},
					Tokens: tk[8:13],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[20],
										Tokens:              tk[20:21],
									},
									Comments: [5]Comments{nil, nil, nil, nil, {&tk[22]}},
									Tokens:   tk[20:23],
								}),
								Tokens: tk[20:23],
							},
						},
						Tokens: tk[20:23],
					},
					Tokens: tk[20:23],
				},
				Comments: [6]Comments{{&tk[2]}, {&tk[6]}, {&tk[14]}, {&tk[18]}},
				Tokens:   tk[:23],
			}
		}},
		{"if (a){} // A\nelse // B\n{}", func(t *test, tk Tokens) { // 12
			t.Output = IfStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&PrimaryExpression{
								IdentifierReference: &tk[3],
								Tokens:              tk[3:4],
							}),
							Tokens: tk[3:4],
						},
					},
					Tokens: tk[3:4],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[5:7],
					},
					Tokens: tk[5:7],
				},
				ElseStatement: &Statement{
					BlockStatement: &Block{
						Tokens: tk[14:16],
					},
					Tokens: tk[14:16],
				},
				Comments: [6]Comments{nil, nil, nil, nil, {&tk[8]}, {&tk[12]}},
				Tokens:   tk[:16],
			}
		}},
	}, func(t *test) (Type, error) {
		var is IfStatement
		err := is.parse(&t.Tokens, t.Yield, t.Await, t.Ret)
		return is, err
	})
}

func TestIterationStatementDo(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidIterationStatementDo,
				Parsing: "IterationStatementDo",
				Token:   tk[0],
			}
		}},
		{`do`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[1]),
						Parsing: "Expression",
						Token:   tk[1],
					},
					Parsing: "Statement",
					Token:   tk[1],
				},
				Parsing: "IterationStatementDo",
				Token:   tk[1],
			}
		}},
		{"do\na", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrInvalidIterationStatementDo,
				Parsing: "IterationStatementDo",
				Token:   tk[3],
			}
		}},
		{"do\na\nwhile", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IterationStatementDo",
				Token:   tk[5],
			}
		}},
		{"do\na\nwhile\n(\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[8]),
					Parsing: "Expression",
					Token:   tk[8],
				},
				Parsing: "IterationStatementDo",
				Token:   tk[8],
			}
		}},
		{"do\na\nwhile\n(\nb\nc\n)", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "IterationStatementDo",
				Token:   tk[10],
			}
		}},
		{"do\na\nwhile\n(\nb\n)", func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementDo{
				Statement: Statement{
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
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[:11],
			}
		}},
		{"do\na\nwhile\n(\nb\n)\n;", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementDo{
				Statement: Statement{
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
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[:13],
			}
		}},
		{"do\na\nwhile\n(\nb\n) c", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementDo",
				Token:   tk[11],
			}
		}},
		{"do\na\nwhile\n(\nb\n)\nc", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementDo{
				Statement: Statement{
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
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[:11],
			}
		}},
		{"do\na:function b(){}\nwhile\n(\ntrue\n)", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "IterationStatementDo",
				Token:   tk[2],
			}
		}},
	}, func(t *test) (Type, error) {
		var is IterationStatementDo

		err := is.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return is, err
	})
}

func TestIterationStatementWhile(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidIterationStatementWhile,
				Parsing: "IterationStatementWhile",
				Token:   tk[0],
			}
		}},
		{`while`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IterationStatementWhile",
				Token:   tk[1],
			}
		}},
		{"while\n(\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "IterationStatementWhile",
				Token:   tk[4],
			}
		}},
		{"while\n(\na\nb\n)", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "IterationStatementWhile",
				Token:   tk[6],
			}
		}},
		{"while\n(\na\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[7]),
						Parsing: "Expression",
						Token:   tk[7],
					},
					Parsing: "Statement",
					Token:   tk[7],
				},
				Parsing: "IterationStatementWhile",
				Token:   tk[7],
			}
		}},
		{"while\n(\na\n)\nb", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementWhile{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[8:9],
							},
						},
						Tokens: tk[8:9],
					},
					Tokens: tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"while(yield a){}", func(t *test, tk Tokens) { // 7
			t.Yield = true
			litA := makeConditionLiteral(tk, 4)
			t.Output = IterationStatementWhile{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							Yield: true,
							AssignmentExpression: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[4:5],
							},
							Tokens: tk[2:5],
						},
					},
					Tokens: tk[2:5],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[6:8],
					},
					Tokens: tk[6:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"while (true) a:function b(){}", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "IterationStatementWhile",
				Token:   tk[6],
			}
		}},
		{"while // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n", func(t *test, tk Tokens) { // 9
			t.Output = IterationStatementWhile{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								},
								Comments: [5]Comments{{&tk[8]}, nil, nil, nil, {&tk[12]}},
								Tokens:   tk[8:13],
							}),
							Tokens: tk[8:13],
						},
					},
					Tokens: tk[8:13],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[20],
										Tokens:              tk[20:21],
									},
									Comments: [5]Comments{nil, nil, nil, nil, {&tk[22]}},
									Tokens:   tk[20:23],
								}),
								Tokens: tk[20:23],
							},
						},
						Tokens: tk[20:23],
					},
					Tokens: tk[20:23],
				},
				Comments: [4]Comments{{&tk[2]}, {&tk[6]}, {&tk[14]}, {&tk[18]}},
				Tokens:   tk[:23],
			}
		}},
	}, func(t *test) (Type, error) {
		var is IterationStatementWhile

		err := is.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return is, err
	})
}

func TestIterationStatementFor(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidIterationStatementFor,
				Parsing: "IterationStatementFor",
				Token:   tk[0],
			}
		}},
		{`for`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IterationStatementFor",
				Token:   tk[1],
			}
		}},
		{"for\nawait", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IterationStatementFor",
				Token:   tk[2],
			}
		}},
		{"for\nawait", func(t *test, tk Tokens) { // 4
			t.Await = true
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "IterationStatementFor",
				Token:   tk[3],
			}
		}},
		{"for\n(\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\nawait\n(\n)", func(t *test, tk Tokens) { // 6
			t.Await = true
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "PrimaryExpression",
								Token:   tk[6],
							},
							Parsing: "MemberExpression",
							Token:   tk[6],
						},
						Parsing: "NewExpression",
						Token:   tk[6],
					},
					Parsing: "LeftHandSideExpression",
					Token:   tk[6],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\n;\n)", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[6]),
					Parsing: "Expression",
					Token:   tk[6],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\nawait\n(\n;\n)", func(t *test, tk Tokens) { // 8
			t.Await = true
			t.Err = Error{
				Err:     ErrInvalidForAwaitLoop,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\n;\n;\n)", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[9]),
						Parsing: "Expression",
						Token:   tk[9],
					},
					Parsing: "Statement",
					Token:   tk[9],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[9],
			}
		}},
		{"for\n(\n;\n;\n)\na", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litA,
								Tokens:                tk[10:11],
							},
						},
						Tokens: tk[10:11],
					},
					Tokens: tk[10:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"for\n(\n;\n,\n)", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[6]),
					Parsing: "Expression",
					Token:   tk[6],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\n;\n;\n,)", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[8]),
					Parsing: "Expression",
					Token:   tk[8],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[8],
			}
		}},
		{"for\n(\n;\na\n;\n)\nb", func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 6)
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[12:13],
							},
						},
						Tokens: tk[12:13],
					},
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"for\n(\n;\n;\na\n)\nb", func(t *test, tk Tokens) { // 14
			litA := makeConditionLiteral(tk, 8)
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[8:9],
						},
					},
					Tokens: tk[8:9],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litB,
								Tokens:                tk[12:13],
							},
						},
						Tokens: tk[12:13],
					},
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"for\n(\n;\na\n;\nb\n)\nc", func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 6)
			litB := makeConditionLiteral(tk, 10)
			litC := makeConditionLiteral(tk, 14)
			t.Output = IterationStatementFor{
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[6:7],
						},
					},
					Tokens: tk[6:7],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: &litC,
								Tokens:                tk[14:15],
							},
						},
						Tokens: tk[14:15],
					},
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nvar)", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "LexicalBinding",
					Token:   tk[5],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[5],
			}
		}},
		{"for\n(\nvar\na\nb)", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[8],
			}
		}},
		{"for\n(\nvar\na,\nb)", func(t *test, tk Tokens) { // 18
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\n(\nvar\na\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 19
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[14:16],
					},
					Tokens: tk[14:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nvar\na\n,\nb\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 20
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
					{
						BindingIdentifier: &tk[10],
						Tokens:            tk[10:11],
					},
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nvar\na\n;\nb\n;\n)\n{}", func(t *test, tk Tokens) { // 21
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Conditional: &Expression{
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
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nvar\na\n;\nb)", func(t *test, tk Tokens) { // 22
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nvar\na\n,\nb\n;\nc\n;\n)\n{}", func(t *test, tk Tokens) { // 23
			litC := makeConditionLiteral(tk, 14)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
					{
						BindingIdentifier: &tk[10],
						Tokens:            tk[10:11],
					},
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[14:15],
						},
					},
					Tokens: tk[14:15],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[20:22],
					},
					Tokens: tk[20:22],
				},
				Tokens: tk[:22],
			}
		}},
		{"for\n(\nvar\na\n;\n;\nb\nc)", func(t *test, tk Tokens) { // 24
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "IterationStatementFor",
				Token:   tk[14],
			}
		}},
		{"for\n(\nvar\na\n;\n;\nb\n)\n{}", func(t *test, tk Tokens) { // 25
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[12:13],
						},
					},
					Tokens: tk[12:13],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nvar\na\n,\nb\n;\n;\nc\n)\n{}", func(t *test, tk Tokens) { // 26
			litC := makeConditionLiteral(tk, 16)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
					{
						BindingIdentifier: &tk[10],
						Tokens:            tk[10:11],
					},
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[16:17],
						},
					},
					Tokens: tk[16:17],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[20:22],
					},
					Tokens: tk[20:22],
				},
				Tokens: tk[:22],
			}
		}},
		{"for\n(\nvar\na\n;\nb\n;\nc\n)\n{}", func(t *test, tk Tokens) { // 27
			litB := makeConditionLiteral(tk, 10)
			litC := makeConditionLiteral(tk, 14)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[14:15],
						},
					},
					Tokens: tk[14:15],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nvar\na\n,\nb\n;\nc\n;\nd\n)\n{}", func(t *test, tk Tokens) { // 28
			litC := makeConditionLiteral(tk, 14)
			litD := makeConditionLiteral(tk, 18)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
					{
						BindingIdentifier: &tk[10],
						Tokens:            tk[10:11],
					},
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[14:15],
						},
					},
					Tokens: tk[14:15],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litD,
							Tokens:                tk[18:19],
						},
					},
					Tokens: tk[18:19],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[22:24],
					},
					Tokens: tk[22:24],
				},
				Tokens: tk[:24],
			}
		}},
		{"for\n(\nlet)", func(t *test, tk Tokens) { // 29
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[5],
					},
					Parsing: "LexicalDeclaration",
					Token:   tk[5],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\n(\nlet\na\nb)", func(t *test, tk Tokens) { // 30
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[8],
			}
		}},
		{"for\n(\nlet\na,\nb)", func(t *test, tk Tokens) { // 31
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidLexicalDeclaration,
					Parsing: "LexicalDeclaration",
					Token:   tk[10],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\n(\nlet\na\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 32
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[14:16],
					},
					Tokens: tk[14:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nlet\na\n,\nb\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 33
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
						{
							BindingIdentifier: &tk[10],
							Tokens:            tk[10:11],
						},
					},
					Tokens: tk[4:13],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nlet\na\n;\nb\n;\n)\n{}", func(t *test, tk Tokens) { // 34
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Conditional: &Expression{
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
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nlet\na\n;\n;\nb\n)\n{}", func(t *test, tk Tokens) { // 35
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[12:13],
						},
					},
					Tokens: tk[12:13],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nlet\na\n;\nb\n;\nc\n)\n{}", func(t *test, tk Tokens) { // 36
			litB := makeConditionLiteral(tk, 10)
			litC := makeConditionLiteral(tk, 14)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[14:15],
						},
					},
					Tokens: tk[14:15],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nconst)", func(t *test, tk Tokens) { // 37
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[5],
					},
					Parsing: "LexicalDeclaration",
					Token:   tk[5],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\n(\nconst\na\nb)", func(t *test, tk Tokens) { // 38
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[8],
			}
		}},
		{"for\n(\nconst\na,\nb)", func(t *test, tk Tokens) { // 39
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidLexicalDeclaration,
					Parsing: "LexicalDeclaration",
					Token:   tk[10],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\n(\nconst\na\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 40
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[14:16],
					},
					Tokens: tk[14:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nconst\na\n,\nb\n;\n;\n)\n{}", func(t *test, tk Tokens) { // 41
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
						{
							BindingIdentifier: &tk[10],
							Tokens:            tk[10:11],
						},
					},
					Tokens: tk[4:13],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nconst\na\n;\nb\n;\n)\n{}", func(t *test, tk Tokens) { // 42
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Conditional: &Expression{
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
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nconst\na\n;\n;\nb\n)\n{}", func(t *test, tk Tokens) { // 43
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[12:13],
						},
					},
					Tokens: tk[12:13],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[16:18],
					},
					Tokens: tk[16:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"for\n(\nconst\na\n;\nb\n;\nc\n)\n{}", func(t *test, tk Tokens) { // 44
			litB := makeConditionLiteral(tk, 10)
			litC := makeConditionLiteral(tk, 14)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[6],
							Tokens:            tk[6:7],
						},
					},
					Tokens: tk[4:9],
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[10:11],
						},
					},
					Tokens: tk[10:11],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litC,
							Tokens:                tk[14:15],
						},
					},
					Tokens: tk[14:15],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[18:20],
					},
					Tokens: tk[18:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nvar\n{}\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 45
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						ObjectBindingPattern: &ObjectBindingPattern{
							Tokens: tk[6:8],
						},
						Initializer: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[11:12],
						},
						Tokens: tk[6:12],
					},
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nvar\n[]\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 46
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						ArrayBindingPattern: &ArrayBindingPattern{
							Tokens: tk[6:8],
						},
						Initializer: &AssignmentExpression{
							ConditionalExpression: &litA,
							Tokens:                tk[11:12],
						},
						Tokens: tk[6:12],
					},
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nlet\n{}\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 47
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								Tokens: tk[6:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[11:12],
							},
							Tokens: tk[6:12],
						},
					},
					Tokens: tk[4:14],
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nlet\n[]\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 48
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								Tokens: tk[6:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[11:12],
							},
							Tokens: tk[6:12],
						},
					},
					Tokens: tk[4:14],
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nconst\n{}\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 49
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							ObjectBindingPattern: &ObjectBindingPattern{
								Tokens: tk[6:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[11:12],
							},
							Tokens: tk[6:12],
						},
					},
					Tokens: tk[4:14],
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nconst\n[]\n=\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 50
			litA := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Const,
					BindingList: []LexicalBinding{
						{
							ArrayBindingPattern: &ArrayBindingPattern{
								Tokens: tk[6:8],
							},
							Initializer: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[11:12],
							},
							Tokens: tk[6:12],
						},
					},
					Tokens: tk[4:14],
				},
				Statement: Statement{
					Tokens: tk[19:20],
				},
				Tokens: tk[:20],
			}
		}},
		{"for\n(\nvar\n{,}\nin)", func(t *test, tk Tokens) { // 51
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n[+]\nin)", func(t *test, tk Tokens) { // 52
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n,\nin)", func(t *test, tk Tokens) { // 53
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n{,}\nin)", func(t *test, tk Tokens) { // 54
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n[+]\nin)", func(t *test, tk Tokens) { // 55
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n,\nin)", func(t *test, tk Tokens) { // 56
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n{,}\nin)", func(t *test, tk Tokens) { // 57
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n[+]\nin)", func(t *test, tk Tokens) { // 58
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n,\nin)", func(t *test, tk Tokens) { // 59
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n{}\nin\n)", func(t *test, tk Tokens) { // 60
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[11]),
					Parsing: "Expression",
					Token:   tk[11],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nlet\n{}\nin\n)", func(t *test, tk Tokens) { // 61
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[11]),
					Parsing: "Expression",
					Token:   tk[11],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nconst\n{}\nin\n)", func(t *test, tk Tokens) { // 62
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[11]),
					Parsing: "Expression",
					Token:   tk[11],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nvar\na\nin\nb\n)\n;", func(t *test, tk Tokens) { // 63
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForInVar,
				ForBindingIdentifier: &tk[6],
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
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nlet\na\nin\nb\n)\n;", func(t *test, tk Tokens) { // 64
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForInLet,
				ForBindingIdentifier: &tk[6],
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
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nconst\na\nin\nb\n)\n;", func(t *test, tk Tokens) { // 65
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForInConst,
				ForBindingIdentifier: &tk[6],
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
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nvar\n{}\nin\nb\n)\n;", func(t *test, tk Tokens) { // 66
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInVar,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nlet\n{}\nin\nb\n)\n;", func(t *test, tk Tokens) { // 67
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInLet,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nconst\n{}\nin\nb\n)\n;", func(t *test, tk Tokens) { // 68
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInConst,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nvar\n[]\nin\nb\n)\n;", func(t *test, tk Tokens) { // 69
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInVar,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nlet\n[]\nin\nb\n)\n;", func(t *test, tk Tokens) { // 70
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInLet,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nconst\n[]\nin\nb\n)\n;", func(t *test, tk Tokens) { // 71
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForInConst,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litB,
							Tokens:                tk[11:12],
						},
					},
					Tokens: tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\nawait\n(\nvar\na\nin\nb\n)\n;", func(t *test, tk Tokens) { // 72
			t.Await = true
			t.Err = Error{
				Err:     ErrInvalidForAwaitLoop,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n{,}\nof)", func(t *test, tk Tokens) { // 73
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n[+]\nof)", func(t *test, tk Tokens) { // 74
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n,\nof)", func(t *test, tk Tokens) { // 75
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n{,}\nof)", func(t *test, tk Tokens) { // 76
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n[+]\nof)", func(t *test, tk Tokens) { // 77
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nlet\n,\nof)", func(t *test, tk Tokens) { // 78
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n{,}\nof)", func(t *test, tk Tokens) { // 79
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[7],
						},
						Parsing: "BindingProperty",
						Token:   tk[7],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n[+]\nof)", func(t *test, tk Tokens) { // 80
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[7],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[7],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nconst\n,\nof)", func(t *test, tk Tokens) { // 81
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\nvar\n{}\nof\n)", func(t *test, tk Tokens) { // 82
			t.Err = Error{
				Err:     assignmentError(tk[11]),
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nlet\n{}\nof\n)", func(t *test, tk Tokens) { // 83
			t.Err = Error{
				Err:     assignmentError(tk[11]),
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nconst\n{}\nof\n)", func(t *test, tk Tokens) { // 84
			t.Err = Error{
				Err:     assignmentError(tk[11]),
				Parsing: "IterationStatementFor",
				Token:   tk[11],
			}
		}},
		{"for\n(\nvar\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 85
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForOfVar,
				ForBindingIdentifier: &tk[6],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[10:11],
				},
				Statement: Statement{
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nlet\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 86
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForOfLet,
				ForBindingIdentifier: &tk[6],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[10:11],
				},
				Statement: Statement{
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nconst\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 87
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type:                 ForOfConst,
				ForBindingIdentifier: &tk[6],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[10:11],
				},
				Statement: Statement{
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\nvar\n{}\nof\nb\n)\n;", func(t *test, tk Tokens) { // 88
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfVar,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nlet\n{}\nof\nb\n)\n;", func(t *test, tk Tokens) { // 89
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfLet,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nconst\n{}\nof\nb\n)\n;", func(t *test, tk Tokens) { // 90
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfConst,
				ForBindingPatternObject: &ObjectBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nvar\n[]\nof\nb\n)\n;", func(t *test, tk Tokens) { // 91
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfVar,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nlet\n[]\nof\nb\n)\n;", func(t *test, tk Tokens) { // 92
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfLet,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\n(\nconst\n[]\nof\nb\n)\n;", func(t *test, tk Tokens) { // 93
			litB := makeConditionLiteral(tk, 11)
			t.Output = IterationStatementFor{
				Type: ForOfConst,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[6:8],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[11:12],
				},
				Statement: Statement{
					Tokens: tk[15:16],
				},
				Tokens: tk[:16],
			}
		}},
		{"for\nawait\n(\nvar\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 94
			t.Await = true
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type:                 ForAwaitOfVar,
				ForBindingIdentifier: &tk[8],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[12:13],
				},
				Statement: Statement{
					Tokens: tk[16:17],
				},
				Tokens: tk[:17],
			}
		}},
		{"for\nawait\n(\nlet\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 95
			t.Await = true
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type:                 ForAwaitOfLet,
				ForBindingIdentifier: &tk[8],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[12:13],
				},
				Statement: Statement{
					Tokens: tk[16:17],
				},
				Tokens: tk[:17],
			}
		}},
		{"for\nawait\n(\nconst\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 96
			t.Await = true
			litB := makeConditionLiteral(tk, 12)
			t.Output = IterationStatementFor{
				Type:                 ForAwaitOfConst,
				ForBindingIdentifier: &tk[8],
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[12:13],
				},
				Statement: Statement{
					Tokens: tk[16:17],
				},
				Tokens: tk[:17],
			}
		}},
		{"for\nawait\n(\n)", func(t *test, tk Tokens) { // 97
			t.Await = true
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "PrimaryExpression",
								Token:   tk[6],
							},
							Parsing: "MemberExpression",
							Token:   tk[6],
						},
						Parsing: "NewExpression",
						Token:   tk[6],
					},
					Parsing: "LeftHandSideExpression",
					Token:   tk[6],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\nawait\n(\na\n)", func(t *test, tk Tokens) { // 98
			t.Await = true
			t.Err = Error{
				Err:     ErrInvalidForAwaitLoop,
				Parsing: "IterationStatementFor",
				Token:   tk[8],
			}
		}},
		{"for\nawait\n(\na\nof\n)", func(t *test, tk Tokens) { // 99
			t.Await = true
			t.Err = Error{
				Err:     assignmentError(tk[10]),
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\nawait\n(\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 100
			t.Await = true
			litB := makeConditionLiteral(tk, 10)
			t.Output = IterationStatementFor{
				Type: ForAwaitOfLeftHandSide,
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &tk[6],
								Tokens:              tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[6:7],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[10:11],
				},
				Statement: Statement{
					Tokens: tk[14:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"for\n(\n)\n;", func(t *test, tk Tokens) { // 101
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "IterationStatementFor",
				Token:   tk[4],
			}
		}},
		{"for\n(\na\n)", func(t *test, tk Tokens) { // 102
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[6],
			}
		}},
		{"for\n(\na\n;\n;\n)\n;", func(t *test, tk Tokens) { // 103
			litA := makeConditionLiteral(tk, 4)
			t.Output = IterationStatementFor{
				Type: ForNormalExpression,
				InitExpression: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Statement: Statement{
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"for\n(\na\nin\nb\n)\n;", func(t *test, tk Tokens) { // 104
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementFor{
				Type: ForInLeftHandSide,
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
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"for\n(\na\nof\nb\n)\n;", func(t *test, tk Tokens) { // 105
			litB := makeConditionLiteral(tk, 8)
			t.Output = IterationStatementFor{
				Type: ForOfLeftHandSide,
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
					Tokens: tk[12:13],
				},
				Tokens: tk[:13],
			}
		}},
		{"for\n(\na\n,\nb\nin)", func(t *test, tk Tokens) { // 106
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\n(\na\n,\nb\nof)", func(t *test, tk Tokens) { // 107
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\n(\na\n=\n1\nin)", func(t *test, tk Tokens) { // 108
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\n(\na\n=\n1\nof)", func(t *test, tk Tokens) { // 109
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for\n(\n!a\nin)", func(t *test, tk Tokens) { // 110
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[7],
			}
		}},
		{"for\n(\n!a\nof)", func(t *test, tk Tokens) { // 111
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "IterationStatementFor",
				Token:   tk[7],
			}
		}},
		{"for (c in d) a: function b(){}", func(t *test, tk Tokens) { // 112
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "IterationStatementFor",
				Token:   tk[10],
			}
		}},
		{"for // A\n( // B\n\n// C\n; // D\n; // E\n\n// F\n) // G\na", func(t *test, tk Tokens) { // 113
			t.Await = true
			t.Output = IterationStatementFor{
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									IdentifierReference: &tk[24],
									Tokens:              tk[24:25],
								}),
								Tokens: tk[24:25],
							},
						},
						Tokens: tk[24:25],
					},
					Tokens: tk[24:25],
				},
				Comments: [8]Comments{{&tk[2]}, nil, {&tk[6]}, {&tk[8]}, {&tk[12]}, {&tk[16]}, {&tk[18]}, {&tk[22]}},
				Tokens:   tk[:25],
			}
		}},
		{"for ( // A\n\n// B\nvar // C\na // D\n, // E\nb // F\n; // G\nc // H\n; // I\n\n// J\n) // K\n{}", func(t *test, tk Tokens) { // 114
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						BindingIdentifier: &tk[12],
						Comments:          [2]Comments{{&tk[10]}, {&tk[14]}},
						Tokens:            tk[10:15],
					},
					{
						BindingIdentifier: &tk[20],
						Comments:          [2]Comments{{&tk[18]}, {&tk[22]}},
						Tokens:            tk[18:23],
					},
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[28],
									Tokens:              tk[28:29],
								},
								Comments: [5]Comments{{&tk[26]}, nil, nil, nil, {&tk[30]}},
								Tokens:   tk[26:31],
							}),
							Tokens: tk[26:31],
						},
					},
					Tokens: tk[26:31],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[42:44],
					},
					Tokens: tk[42:44],
				},
				Comments: [8]Comments{nil, nil, {&tk[4]}, {&tk[6]}, nil, {&tk[34]}, {&tk[36]}, {&tk[40]}},
				Tokens:   tk[:44],
			}
		}},
		{"for ( // A\n\n// B\n; // C\na // D\n; // E\nb // F\n\n// G\n) // H\nc", func(t *test, tk Tokens) { // 115
			t.Output = IterationStatementFor{
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[12],
									Tokens:              tk[12:13],
								},
								Comments: [5]Comments{{&tk[10]}, nil, nil, nil, {&tk[14]}},
								Tokens:   tk[10:15],
							}),
							Tokens: tk[10:15],
						},
					},
					Tokens: tk[10:15],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[20],
									Tokens:              tk[20:21],
								},
								Comments: [5]Comments{{&tk[18]}, nil, nil, nil, {&tk[22]}},
								Tokens:   tk[18:23],
							}),
							Tokens: tk[18:23],
						},
					},
					Tokens: tk[18:23],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&PrimaryExpression{
									IdentifierReference: &tk[30],
									Tokens:              tk[30:31],
								}),
								Tokens: tk[30:31],
							},
						},
						Tokens: tk[30:31],
					},
					Tokens: tk[30:31],
				},
				Comments: [8]Comments{nil, nil, {&tk[4]}, {&tk[6]}, nil, nil, {&tk[24]}, {&tk[28]}},
				Tokens:   tk[:31],
			}
		}},
		{"for ( // A\n\n// B\nlet // C\na // D\n; b; c) {}", func(t *test, tk Tokens) { // 116
			t.Output = IterationStatementFor{
				Type: ForNormalLexicalDeclaration,
				InitLexical: &LexicalDeclaration{
					LetOrConst: Let,
					BindingList: []LexicalBinding{
						{
							BindingIdentifier: &tk[12],
							Comments:          [2]Comments{{&tk[10]}, {&tk[14]}},
							Tokens:            tk[10:15],
						},
					},
					Tokens: tk[8:17],
				},
				Conditional: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&PrimaryExpression{
								IdentifierReference: &tk[18],
								Tokens:              tk[18:19],
							}),
							Tokens: tk[18:19],
						},
					},
					Tokens: tk[18:19],
				},
				Afterthought: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&PrimaryExpression{
								IdentifierReference: &tk[21],
								Tokens:              tk[21:22],
							}),
							Tokens: tk[21:22],
						},
					},
					Tokens: tk[21:22],
				},
				Statement: Statement{
					BlockStatement: &Block{
						Tokens: tk[24:26],
					},
					Tokens: tk[24:26],
				},
				Comments: [8]Comments{nil, nil, {&tk[4]}, {&tk[6]}},
				Tokens:   tk[:26],
			}
		}},
		{"for ( // A\n\n// B\nvar // C\n{}// D\n= // E\na // F\n;;);", func(t *test, tk Tokens) { // 117
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						ObjectBindingPattern: &ObjectBindingPattern{
							Tokens: tk[12:14],
						},
						Initializer: &AssignmentExpression{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[20],
									Tokens:              tk[20:21],
								},
								Comments: [5]Comments{{&tk[18]}, nil, nil, nil, {&tk[22]}},
								Tokens:   tk[18:23],
							}),
							Tokens: tk[18:23],
						},
						Comments: [2]Comments{{&tk[10]}, {&tk[14]}},
						Tokens:   tk[10:23],
					},
				},
				Statement: Statement{
					Tokens: tk[27:28],
				},
				Comments: [8]Comments{nil, nil, {&tk[4]}, {&tk[6]}},
				Tokens:   tk[:28],
			}
		}},
		{"for ( // A\n\n// B\nvar // C\n[]// D\n= a // E\n; // F\n; // G\n) // H\n;", func(t *test, tk Tokens) { // 118
			t.Output = IterationStatementFor{
				Type: ForNormalVar,
				InitVar: []VariableDeclaration{
					{
						ArrayBindingPattern: &ArrayBindingPattern{
							Tokens: tk[12:14],
						},
						Initializer: &AssignmentExpression{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[18],
									Tokens:              tk[18:19],
								},
								Comments: [5]Comments{nil, nil, nil, nil, {&tk[20]}},
								Tokens:   tk[18:21],
							}),
							Tokens: tk[18:21],
						},
						Comments: [2]Comments{{&tk[10]}, {&tk[14]}},
						Tokens:   tk[10:21],
					},
				},
				Statement: Statement{
					Tokens: tk[34:35],
				},
				Comments: [8]Comments{nil, nil, {&tk[4]}, {&tk[6]}, {&tk[24]}, {&tk[28]}, nil, {&tk[32]}},
				Tokens:   tk[:35],
			}
		}},
		{"for (\n// A\nvar // B\na // C\nin // D\nb // E\n\n// F\n);", func(t *test, tk Tokens) { // 119
			t.Output = IterationStatementFor{
				Type:                 ForInVar,
				ForBindingIdentifier: &tk[10],
				In: &Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[18],
									Tokens:              tk[18:19],
								},
								Comments: [5]Comments{{&tk[16]}, nil, nil, nil, {&tk[20]}},
								Tokens:   tk[16:21],
							}),
							Tokens: tk[16:21],
						},
					},
					Tokens: tk[16:21],
				},
				Statement: Statement{
					Tokens: tk[25:26],
				},
				Comments: [8]Comments{nil, nil, nil, {&tk[4]}, {&tk[8]}, {&tk[12]}, {&tk[22]}},
				Tokens:   tk[:26],
			}
		}},
		{"for // A\nawait // B\n( // C\n\n// D\nconst // E\n[]// F\nof // G\nb // H\n\n// I\n) // J\n;", func(t *test, tk Tokens) { // 120
			t.Await = true
			t.Output = IterationStatementFor{
				Type: ForAwaitOfConst,
				ForBindingPatternArray: &ArrayBindingPattern{
					Tokens: tk[18:20],
				},
				Of: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							IdentifierReference: &tk[26],
							Tokens:              tk[26:27],
						},
						Comments: [5]Comments{{&tk[24]}, nil, nil, nil, {&tk[28]}},
						Tokens:   tk[24:29],
					}),
					Tokens: tk[24:29],
				},
				Statement: Statement{
					Tokens: tk[36:37],
				},
				Comments: [8]Comments{{&tk[2]}, {&tk[6]}, {&tk[10]}, {&tk[12]}, {&tk[16]}, {&tk[20]}, {&tk[30]}, {&tk[34]}},
				Tokens:   tk[:37],
			}
		}},
	}, func(t *test) (Type, error) {
		var is IterationStatementFor

		err := is.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return is, err
	})
}

func TestSwitchStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidSwitchStatement,
				Parsing: "SwitchStatement",
				Token:   tk[0],
			}
		}},
		{`switch`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "SwitchStatement",
				Token:   tk[1],
			}
		}},
		{"switch\n(\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "SwitchStatement",
				Token:   tk[4],
			}
		}},
		{"switch\n(\na\nb\n)", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "SwitchStatement",
				Token:   tk[6],
			}
		}},
		{"switch\n(\na\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "SwitchStatement",
				Token:   tk[7],
			}
		}},
		{"switch\n(\na\n)\n{\n}", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 4)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:11],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n}", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "SwitchStatement",
				Token:   tk[12],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\n}", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 4)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				DefaultClause: []StatementListItem{},
				Tokens:        tk[:15],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\nlet\n}", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "LexicalBinding",
								Token:   tk[16],
							},
							Parsing: "LexicalDeclaration",
							Token:   tk[16],
						},
						Parsing: "Declaration",
						Token:   tk[14],
					},
					Parsing: "StatementListItem",
					Token:   tk[14],
				},
				Parsing: "SwitchStatement",
				Token:   tk[14],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\nb\n}", func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 14)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				DefaultClause: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litB,
										Tokens:                tk[14:15],
									},
								},
								Tokens: tk[14:15],
							},
							Tokens: tk[14:15],
						},
						Tokens: tk[14:15],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\nb;\nc\n}", func(t *test, tk Tokens) { // 11
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 14)
			litC := makeConditionLiteral(tk, 17)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				DefaultClause: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litB,
										Tokens:                tk[14:15],
									},
								},
								Tokens: tk[14:15],
							},
							Tokens: tk[14:16],
						},
						Tokens: tk[14:16],
					},
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
							Tokens: tk[17:18],
						},
						Tokens: tk[17:18],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\ncase}", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[15]),
						Parsing: "Expression",
						Token:   tk[15],
					},
					Parsing: "CaseClause",
					Token:   tk[15],
				},
				Parsing: "SwitchStatement",
				Token:   tk[14],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\ncase b:\n}", func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 16)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				DefaultClause: []StatementListItem{},
				PostDefaultCaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[16:17],
								},
							},
							Tokens: tk[16:17],
						},
						Tokens: tk[14:18],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{"switch\n(\na\n)\n{\ndefault\n:\ncase b:\ncase c:}", func(t *test, tk Tokens) { // 14
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 16)
			litC := makeConditionLiteral(tk, 21)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				DefaultClause: []StatementListItem{},
				PostDefaultCaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[16:17],
								},
							},
							Tokens: tk[16:17],
						},
						Tokens: tk[14:18],
					},
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litC,
									Tokens:                tk[21:22],
								},
							},
							Tokens: tk[21:22],
						},
						Tokens: tk[19:23],
					},
				},
				Tokens: tk[:24],
			}
		}},
		{"switch\n(\na\n)\n{\ncase b:\n}", func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 12)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				CaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[12:13],
								},
							},
							Tokens: tk[12:13],
						},
						Tokens: tk[10:14],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{"switch\n(\na\n)\n{\ncase b:\ncase c:}", func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 12)
			litC := makeConditionLiteral(tk, 17)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				CaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[12:13],
								},
							},
							Tokens: tk[12:13],
						},
						Tokens: tk[10:14],
					},
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litC,
									Tokens:                tk[17:18],
								},
							},
							Tokens: tk[17:18],
						},
						Tokens: tk[15:19],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{"switch\n(\na\n)\n{\ncase b:\ncase c:\ndefault\n:\nd\ne\ncase f:\ncase g:}", func(t *test, tk Tokens) { // 17
			litA := makeConditionLiteral(tk, 4)
			litB := makeConditionLiteral(tk, 12)
			litC := makeConditionLiteral(tk, 17)
			litD := makeConditionLiteral(tk, 24)
			litE := makeConditionLiteral(tk, 26)
			litF := makeConditionLiteral(tk, 30)
			litG := makeConditionLiteral(tk, 35)
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				CaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litB,
									Tokens:                tk[12:13],
								},
							},
							Tokens: tk[12:13],
						},
						Tokens: tk[10:14],
					},
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litC,
									Tokens:                tk[17:18],
								},
							},
							Tokens: tk[17:18],
						},
						Tokens: tk[15:19],
					},
				},
				DefaultClause: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litD,
										Tokens:                tk[24:25],
									},
								},
								Tokens: tk[24:25],
							},
							Tokens: tk[24:25],
						},
						Tokens: tk[24:25],
					},
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litE,
										Tokens:                tk[26:27],
									},
								},
								Tokens: tk[26:27],
							},
							Tokens: tk[26:27],
						},
						Tokens: tk[26:27],
					},
				},
				PostDefaultCaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litF,
									Tokens:                tk[30:31],
								},
							},
							Tokens: tk[30:31],
						},
						Tokens: tk[28:32],
					},
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: &litG,
									Tokens:                tk[35:36],
								},
							},
							Tokens: tk[35:36],
						},
						Tokens: tk[33:37],
					},
				},
				Tokens: tk[:38],
			}
		}},
		{"switch\n(\na\n)\n{\ncase b:\ncase c:\ndefault\n:\nd;\ne;\ncase f:\ncase g:default:}", func(t *test, tk Tokens) { // 18
			t.Err = Error{
				Err:     ErrDuplicateDefaultClause,
				Parsing: "SwitchStatement",
				Token:   tk[39],
			}
		}},
		{"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\n}", func(t *test, tk Tokens) { // 19
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								},
								Comments: [5]Comments{{&tk[8]}, nil, nil, nil, {&tk[12]}},
								Tokens:   tk[8:13],
							}),
							Tokens: tk[8:13],
						},
					},
					Tokens: tk[8:13],
				},
				Comments: [9]Comments{{&tk[2]}, {&tk[6]}, {&tk[14]}, {&tk[18]}, {&tk[22]}, nil, nil, nil, {&tk[24]}},
				Tokens:   tk[:27],
			}
		}},
		{"switch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{ // G\n\n// H\ncase /* I */ b /* J */ : // K\n\n// J\ncase c:// L\n\n// M\ndefault /* N */ : // O\nd;// P\n\n// Q\n}", func(t *test, tk Tokens) { // 20
			t.Output = SwitchStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								},
								Comments: [5]Comments{{&tk[8]}, nil, nil, nil, {&tk[12]}},
								Tokens:   tk[8:13],
							}),
							Tokens: tk[8:13],
						},
					},
					Tokens: tk[8:13],
				},
				CaseClauses: []CaseClause{
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: WrapConditional(&MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: &tk[30],
											Tokens:              tk[30:31],
										},
										Comments: [5]Comments{{&tk[28]}, nil, nil, nil, {&tk[32]}},
										Tokens:   tk[28:33],
									}),
									Tokens: tk[28:33],
								},
							},
							Tokens: tk[28:33],
						},
						Comments: [2]Comments{{&tk[24]}, {&tk[36]}},
						Tokens:   tk[24:37],
					},
					{
						Expression: Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: WrapConditional(&PrimaryExpression{
										IdentifierReference: &tk[42],
										Tokens:              tk[42:43],
									}),
									Tokens: tk[42:43],
								},
							},
							Tokens: tk[42:43],
						},
						Comments: [2]Comments{{&tk[38]}, {&tk[44]}},
						Tokens:   tk[38:45],
					},
				},
				DefaultClause: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: WrapConditional(&PrimaryExpression{
											IdentifierReference: &tk[56],
											Tokens:              tk[56:57],
										}),
										Tokens: tk[56:57],
									},
								},
								Tokens: tk[56:57],
							},
							Tokens: tk[56:58],
						},
						Comments: [2]Comments{nil, {&tk[58]}},
						Tokens:   tk[56:59],
					},
				},
				Comments: [9]Comments{{&tk[2]}, {&tk[6]}, {&tk[14]}, {&tk[18]}, {&tk[22]}, {&tk[46]}, {&tk[50]}, {&tk[54]}, {&tk[60]}},
				Tokens:   tk[:63],
			}
		}},
	}, func(t *test) (Type, error) {
		var ss SwitchStatement

		err := ss.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return ss, err
	})
}

func TestCaseClause(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingCaseClause,
				Parsing: "CaseClause",
				Token:   tk[0],
			}
		}},
		{`case`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[1]),
					Parsing: "Expression",
					Token:   tk[1],
				},
				Parsing: "CaseClause",
				Token:   tk[1],
			}
		}},
		{"case\na", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "CaseClause",
				Token:   tk[3],
			}
		}},
		{"case\na\n:", func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 2)
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"case\na\n:case", func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 2)
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"case\na\n:default", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 2)
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				Tokens: tk[:5],
			}
		}},
		{"case\na\n:\nlet", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "LexicalBinding",
								Token:   tk[7],
							},
							Parsing: "LexicalDeclaration",
							Token:   tk[7],
						},
						Parsing: "Declaration",
						Token:   tk[6],
					},
					Parsing: "StatementListItem",
					Token:   tk[6],
				},
				Parsing: "CaseClause",
				Token:   tk[6],
			}
		}},
		{"case\na\n:\nb\nc", func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 2)
			litB := makeConditionLiteral(tk, 6)
			litC := makeConditionLiteral(tk, 8)
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[2:3],
						},
					},
					Tokens: tk[2:3],
				},
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litB,
										Tokens:                tk[6:7],
									},
								},
								Tokens: tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[6:7],
					},
					{
						Statement: &Statement{
							ExpressionStatement: &Expression{
								Expressions: []AssignmentExpression{
									{
										ConditionalExpression: &litC,
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
		{"// A\ncase /* B */ a /* C */: // D\n\n// E\ncase", func(t *test, tk Tokens) { // 9
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[6],
									Tokens:              tk[6:7],
								},
								Comments: [5]Comments{{&tk[4]}, nil, nil, nil, {&tk[8]}},
								Tokens:   tk[4:9],
							}),
							Tokens: tk[4:9],
						},
					},
					Tokens: tk[4:9],
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[11]}},
				Tokens:   tk[:12],
			}
		}},
		{"// A\ncase /* B */ a /* C */: // D\n\n// E\n{} // F\ncase", func(t *test, tk Tokens) { // 10
			t.Output = CaseClause{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[6],
									Tokens:              tk[6:7],
								},
								Comments: [5]Comments{{&tk[4]}, nil, nil, nil, {&tk[8]}},
								Tokens:   tk[4:9],
							}),
							Tokens: tk[4:9],
						},
					},
					Tokens: tk[4:9],
				},
				StatementList: []StatementListItem{
					{
						Statement: &Statement{
							BlockStatement: &Block{
								Tokens: tk[15:17],
							},
							Tokens: tk[15:17],
						},
						Comments: [2]Comments{{&tk[13]}, {&tk[18]}},
						Tokens:   tk[13:19],
					},
				},
				Comments: [2]Comments{{&tk[0]}, {&tk[11]}},
				Tokens:   tk[:19],
			}
		}},
	}, func(t *test) (Type, error) {
		var cc CaseClause

		err := cc.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return cc, err
	})
}

func TestWithStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidWithStatement,
				Parsing: "WithStatement",
				Token:   tk[0],
			}
		}},
		{`with`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingOpeningParenthesis,
				Parsing: "WithStatement",
				Token:   tk[1],
			}
		}},
		{"with\n(\n)", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     assignmentError(tk[4]),
					Parsing: "Expression",
					Token:   tk[4],
				},
				Parsing: "WithStatement",
				Token:   tk[4],
			}
		}},
		{"with\n(\na\nb\n)", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "WithStatement",
				Token:   tk[6],
			}
		}},
		{"with\n(\na\n)", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     assignmentError(tk[7]),
						Parsing: "Expression",
						Token:   tk[7],
					},
					Parsing: "Statement",
					Token:   tk[7],
				},
				Parsing: "WithStatement",
				Token:   tk[7],
			}
		}},
		{"with\n(\na\n)\n;", func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 4)
			t.Output = WithStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: &litA,
							Tokens:                tk[4:5],
						},
					},
					Tokens: tk[4:5],
				},
				Statement: Statement{
					Tokens: tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"with (b) a: function b(){}", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrLabelledFunction,
				Parsing: "WithStatement",
				Token:   tk[6],
			}
		}},
		{"with // A\n( // B\n\n// C\na // D\n\n// E\n) // F\nb // G\n", func(t *test, tk Tokens) { // 8
			t.Output = WithStatement{
				Expression: Expression{
					Expressions: []AssignmentExpression{
						{
							ConditionalExpression: WrapConditional(&MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: &tk[10],
									Tokens:              tk[10:11],
								},
								Comments: [5]Comments{{&tk[8]}, nil, nil, nil, {&tk[12]}},
								Tokens:   tk[8:13],
							}),
							Tokens: tk[8:13],
						},
					},
					Tokens: tk[8:13],
				},
				Statement: Statement{
					ExpressionStatement: &Expression{
						Expressions: []AssignmentExpression{
							{
								ConditionalExpression: WrapConditional(&MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: &tk[20],
										Tokens:              tk[20:21],
									},
									Comments: [5]Comments{nil, nil, nil, nil, {&tk[22]}},
									Tokens:   tk[20:23],
								}),
								Tokens: tk[20:23],
							},
						},
						Tokens: tk[20:23],
					},
					Tokens: tk[20:23],
				},
				Comments: [4]Comments{{&tk[2]}, {&tk[6]}, {&tk[14]}, {&tk[18]}},
				Tokens:   tk[:23],
			}
		}},
	}, func(t *test) (Type, error) {
		var ws WithStatement

		err := ws.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return ws, err
	})
}

func TestTryStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidTryStatement,
				Parsing: "TryStatement",
				Token:   tk[0],
			}
		}},
		{`try`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[1],
				},
				Parsing: "TryStatement",
				Token:   tk[1],
			}
		}},
		{"try\n{\n}", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingCatchFinally,
				Parsing: "TryStatement",
				Token:   tk[5],
			}
		}},
		{"try\n{\n}\ncatch", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[7],
				},
				Parsing: "TryStatement",
				Token:   tk[7],
			}
		}},
		{"try\n{\n}\ncatch\n{,}", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     assignmentError(tk[9]),
								Parsing: "Expression",
								Token:   tk[9],
							},
							Parsing: "Statement",
							Token:   tk[9],
						},
						Parsing: "StatementListItem",
						Token:   tk[9],
					},
					Parsing: "Block",
					Token:   tk[9],
				},
				Parsing: "TryStatement",
				Token:   tk[8],
			}
		}},
		{"try\n{\n}\ncatch\n{\n}", func(t *test, tk Tokens) { // 6
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				CatchBlock: &Block{
					Tokens: tk[8:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"try\n{\n}\ncatch\n(\n)", func(t *test, tk Tokens) { // 7
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "TryStatement",
				Token:   tk[10],
			}
		}},
		{"try\n{\n}\ncatch\n(\na\n)\n{\n}", func(t *test, tk Tokens) { // 8
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				CatchParameterBindingIdentifier: &tk[10],
				CatchBlock: &Block{
					Tokens: tk[14:17],
				},
				Tokens: tk[:17],
			}
		}},
		{"try\n{\n}\ncatch\n(\n{,}\n)", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidPropertyName,
							Parsing: "PropertyName",
							Token:   tk[11],
						},
						Parsing: "BindingProperty",
						Token:   tk[11],
					},
					Parsing: "ObjectBindingPattern",
					Token:   tk[11],
				},
				Parsing: "TryStatement",
				Token:   tk[10],
			}
		}},
		{"try\n{\n}\ncatch\n(\n{}\n)\n{\n}", func(t *test, tk Tokens) { // 10
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				CatchParameterObjectBindingPattern: &ObjectBindingPattern{
					Tokens: tk[10:12],
				},
				CatchBlock: &Block{
					Tokens: tk[15:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"try\n{\n}\ncatch\n(\n[!]\n)", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "BindingElement",
						Token:   tk[11],
					},
					Parsing: "ArrayBindingPattern",
					Token:   tk[11],
				},
				Parsing: "TryStatement",
				Token:   tk[10],
			}
		}},
		{"try\n{\n}\ncatch\n(\n[]\n)\n{\n}", func(t *test, tk Tokens) { // 12
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				CatchParameterArrayBindingPattern: &ArrayBindingPattern{
					Tokens: tk[10:12],
				},
				CatchBlock: &Block{
					Tokens: tk[15:18],
				},
				Tokens: tk[:18],
			}
		}},
		{"try\n{\n}\nfinally", func(t *test, tk Tokens) { // 13
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "Block",
					Token:   tk[7],
				},
				Parsing: "TryStatement",
				Token:   tk[7],
			}
		}},
		{"try\n{\n}\nfinally\n{\n}", func(t *test, tk Tokens) { // 14
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				FinallyBlock: &Block{
					Tokens: tk[8:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"try\n{\n}\nfinally\n{\n}\ncatch\n{\n}", func(t *test, tk Tokens) { // 15
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				FinallyBlock: &Block{
					Tokens: tk[8:11],
				},
				Tokens: tk[:11],
			}
		}},
		{"try\n{\n}\ncatch\n{\n}\nfinally\n{\n}", func(t *test, tk Tokens) { // 16
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:5],
				},
				CatchBlock: &Block{
					Tokens: tk[8:11],
				},
				FinallyBlock: &Block{
					Tokens: tk[14:17],
				},
				Tokens: tk[:17],
			}
		}},
		{"try\n{\n}\ncatch\n(\na\nb\n)", func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err:     ErrMissingClosingParenthesis,
				Parsing: "TryStatement",
				Token:   tk[12],
			}
		}},
		{"try // A\n{} // B\ncatch // C\n{}", func(t *test, tk Tokens) { // 18
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[4:6],
				},
				CatchBlock: &Block{
					Tokens: tk[13:15],
				},
				Comments: [10]Comments{{&tk[2]}, {&tk[7]}, {&tk[11]}},
				Tokens:   tk[:15],
			}
		}},
		{"try {}catch // A\n( // B\n\n// C\na // D\n\n// E\n) // F\n{} // G", func(t *test, tk Tokens) { // 19
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[2:4],
				},
				CatchParameterBindingIdentifier: &tk[14],
				CatchBlock: &Block{
					Tokens: tk[24:26],
				},
				Comments: [10]Comments{nil, nil, {&tk[6]}, {&tk[10]}, {&tk[12]}, {&tk[16]}, {&tk[18]}, {&tk[22]}},
				Tokens:   tk[:26],
			}
		}},
		{"try{} // A\nfinally // B\n{} // C", func(t *test, tk Tokens) { // 20
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[1:3],
				},
				FinallyBlock: &Block{
					Tokens: tk[10:12],
				},
				Comments: [10]Comments{nil, nil, nil, nil, nil, nil, nil, nil, {&tk[4]}, {&tk[8]}},
				Tokens:   tk[:12],
			}
		}},
		{"try // A\n{}// B\ncatch /* C */ ( // D\n\n// E\na // F\n\n// G\n) // H\n{}// I\nfinally /* J */ {} // K", func(t *test, tk Tokens) { // 21
			t.Output = TryStatement{
				TryBlock: Block{
					Tokens: tk[4:6],
				},
				CatchParameterBindingIdentifier: &tk[18],
				CatchBlock: &Block{
					Tokens: tk[28:30],
				},
				FinallyBlock: &Block{
					Tokens: tk[36:38],
				},
				Comments: [10]Comments{{&tk[2]}, {&tk[6]}, {&tk[10]}, {&tk[14]}, {&tk[16]}, {&tk[20]}, {&tk[22]}, {&tk[26]}, {&tk[30]}, {&tk[34]}},
				Tokens:   tk[:38],
			}
		}},
	}, func(t *test) (Type, error) {
		var ts TryStatement

		err := ts.parse(&t.Tokens, t.Yield, t.Await, t.Ret)

		return ts, err
	})
}

func TestVariableStatement(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidVariableStatement,
				Parsing: "VariableStatement",
				Token:   tk[0],
			}
		}},
		{`var`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "LexicalBinding",
					Token:   tk[1],
				},
				Parsing: "VariableStatement",
				Token:   tk[1],
			}
		}},
		{"var\na", func(t *test, tk Tokens) { // 3
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"var\na\n,", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "LexicalBinding",
					Token:   tk[5],
				},
				Parsing: "VariableStatement",
				Token:   tk[5],
			}
		}},
		{"var\na\n", func(t *test, tk Tokens) { // 5
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"var\na\n;", func(t *test, tk Tokens) { // 6
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"var\na\nb", func(t *test, tk Tokens) { // 7
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"var\na\n,\nb", func(t *test, tk Tokens) { // 8
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"var\na\n,\nb\n", func(t *test, tk Tokens) { // 9
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"var\na\n,\nb\n;", func(t *test, tk Tokens) { // 10
			t.Output = VariableStatement{
				VariableDeclarationList: []VariableDeclaration{
					{
						BindingIdentifier: &tk[2],
						Tokens:            tk[2:3],
					},
					{
						BindingIdentifier: &tk[6],
						Tokens:            tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"var\na b\n;", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "VariableStatement",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) {
		var vs VariableStatement

		err := vs.parse(&t.Tokens, t.Yield, t.Await)

		return vs, err
	})
}
