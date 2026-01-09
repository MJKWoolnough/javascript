package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
)

func TestBlockAsModule(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Block
		Callback func(*javascript.Module) bool
		Output   *javascript.Block
	}{
		{ // 1
			&javascript.Block{},
			func(m *javascript.Module) bool { return false },
			&javascript.Block{},
		},
		{ // 2
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementDebugger,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementDebugger,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
				},
			},
			func(m *javascript.Module) bool {
				for i := 0; i < len(m.ModuleListItems); i++ {
					if m.ModuleListItems[i].StatementListItem.Statement.Type == javascript.StatementDebugger {
						m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
						i--
					}
				}

				return false
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementContinue,
						},
					},
				},
			},
		},
	} {
		blockAsModule(test.Input, test.Callback)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting output %v, got %v", n+1, test.Output, test.Input)
		}
	}
}

func TestExpressionsAsModule(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Expression
		Callback func(*javascript.Module) bool
		Output   *javascript.Expression
	}{
		{ // 1
			&javascript.Expression{},
			func(m *javascript.Module) bool { return false },
			&javascript.Expression{},
		},
		{ // 2
			&javascript.Expression{
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
							IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
						}),
					},
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "d"),
						}),
					},
				},
			},
			func(m *javascript.Module) bool {
				for i := 0; i < len(m.ModuleListItems); i++ {
					if javascript.UnwrapConditional(m.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference.Data == "b" {
						m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
					}
				}

				return false
			},
			&javascript.Expression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
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
		},
	} {
		expressionsAsModule(test.Input, test.Callback)
		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting output %v, got %v", n+1, test.Output, test.Input)
		}
	}
}

func TestIsReturnStatement(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Statement
		IsReturn bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.Statement{},
			false,
		},
		{ // 3
			&javascript.Statement{
				Type: javascript.StatementReturn,
			},
			true,
		},
		{ // 4
			&javascript.Statement{
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
			true,
		},
	} {
		if isReturnStatement(test.Input) != test.IsReturn {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsReturn, !test.IsReturn)
		}
	}
}

func TestIsNonEmptyReturnStatement(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Statement
		IsReturn bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.Statement{},
			false,
		},
		{ // 3
			&javascript.Statement{
				Type: javascript.StatementReturn,
			},
			false,
		},
		{ // 4
			&javascript.Statement{
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
			true,
		},
	} {
		if isNonEmptyReturnStatement(test.Input) != test.IsReturn {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsReturn, !test.IsReturn)
		}
	}
}

func TestIsStatementExpression(t *testing.T) {
	for n, test := range [...]struct {
		Input                 *javascript.Statement
		IsStatementExpression bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.Statement{},
			false,
		},
		{ // 3
			&javascript.Statement{
				Type: javascript.StatementDebugger,
			},
			false,
		},
		{ // 4
			&javascript.Statement{
				Type: javascript.StatementDebugger,
			},
			false,
		},
		{ // 5
			&javascript.Statement{
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
			false,
		},
		{ // 6
			&javascript.Statement{
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
			true,
		},
	} {
		if isStatementExpression(test.Input) != test.IsStatementExpression {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsStatementExpression, !test.IsStatementExpression)
		}
	}
}

func TestIsSLIExpression(t *testing.T) {
	for n, test := range [...]struct {
		Input           *javascript.StatementListItem
		IsSLIExpression bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.StatementListItem{},
			false,
		},
		{ // 3
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					Type: javascript.StatementDebugger,
				},
			},
			false,
		},
		{ // 4
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					Type: javascript.StatementDebugger,
				},
			},
			false,
		},
		{ // 5
			&javascript.StatementListItem{
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
			false,
		},
		{ // 6
			&javascript.StatementListItem{
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
			true,
		},
		{ // 7
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{
					FunctionDeclaration: &javascript.FunctionDeclaration{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					},
				},
			},
			false,
		},
	} {
		if isSLIExpression(test.Input) != test.IsSLIExpression {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsSLIExpression, !test.IsSLIExpression)
		}
	}
}

func TestIsEmptyStatement(t *testing.T) {
	for n, test := range [...]struct {
		Input            *javascript.Statement
		IsEmptyStatement bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.Statement{},
			true,
		},
		{ // 3
			&javascript.Statement{
				Type: javascript.StatementContinue,
			},
			false,
		},
		{ // 4
			&javascript.Statement{
				BlockStatement: &javascript.Block{},
			},
			false,
		},
	} {
		if isEmptyStatement(test.Input) != test.IsEmptyStatement {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsEmptyStatement, !test.IsEmptyStatement)
		}
	}
}

func TestIsHoistable(t *testing.T) {
	for n, test := range [...]struct {
		Input       *javascript.StatementListItem
		IsHoistable bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.StatementListItem{},
			false,
		},
		{ // 3
			&javascript.StatementListItem{
				Statement: &javascript.Statement{},
			},
			false,
		},
		{ // 4
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					VariableStatement: &javascript.VariableStatement{},
				},
			},
			true,
		},
		{ // 5
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					LabelIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					LabelledItemFunction: &javascript.FunctionDeclaration{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "b"),
					},
				},
			},
			true,
		},
		{ // 6
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{},
			},
			false,
		},
		{ // 7
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{
					LexicalDeclaration: &javascript.LexicalDeclaration{},
				},
			},
			false,
		},
		{ // 8
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{
					FunctionDeclaration: &javascript.FunctionDeclaration{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					},
				},
			},
			true,
		},
		{ // 9
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{
					ClassDeclaration: &javascript.ClassDeclaration{
						BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
					},
				},
			},
			true,
		},
	} {
		if isHoistable(test.Input) != test.IsHoistable {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsHoistable, !test.IsHoistable)
		}
	}
}

func TestStatementsListItemsAsExpressionsAndReturn(t *testing.T) {
	for n, test := range [...]struct {
		StatementListItems []javascript.StatementListItem
		Expressions        []javascript.AssignmentExpression
		Return             bool
	}{
		{ // 1
			nil,
			nil,
			false,
		},
		{ // 2
			[]javascript.StatementListItem{},
			nil,
			false,
		},
		{ // 3
			[]javascript.StatementListItem{
				{
					Declaration: &javascript.Declaration{
						FunctionDeclaration: &javascript.FunctionDeclaration{
							BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
				},
			},
			nil,
			false,
		},
		{ // 4
			[]javascript.StatementListItem{
				{
					Statement: &javascript.Statement{},
				},
			},
			nil,
			false,
		},
		{ // 5
			[]javascript.StatementListItem{
				{
					Statement: &javascript.Statement{
						Type: javascript.StatementReturn,
					},
				},
			},
			nil,
			false,
		},
		{ // 6
			[]javascript.StatementListItem{
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
			[]javascript.AssignmentExpression{
				{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					}),
				},
			},
			true,
		},
		{ // 7
			[]javascript.StatementListItem{
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
			[]javascript.AssignmentExpression{
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
			true,
		},
		{ // 8
			[]javascript.StatementListItem{
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
			[]javascript.AssignmentExpression{
				{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					}),
				},
			},
			false,
		},
		{ // 9
			[]javascript.StatementListItem{
				{
					Statement: &javascript.Statement{
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
							},
						},
					},
				},
			},
			[]javascript.AssignmentExpression{
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
			false,
		},
		{ // 10
			[]javascript.StatementListItem{
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
			},
			[]javascript.AssignmentExpression{
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
			false,
		},
		{ // 11
			[]javascript.StatementListItem{
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
			[]javascript.AssignmentExpression{
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
			true,
		},
		{ // 12
			[]javascript.StatementListItem{
				{
					Declaration: &javascript.Declaration{
						FunctionDeclaration: &javascript.FunctionDeclaration{
							BindingIdentifier: makeToken(javascript.TokenIdentifier, "a"),
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
										IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
									}),
								},
							},
						},
					},
				},
			},
			nil,
			false,
		},
		{ // 13
			[]javascript.StatementListItem{
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
				{
					Declaration: &javascript.Declaration{
						FunctionDeclaration: &javascript.FunctionDeclaration{
							BindingIdentifier: makeToken(javascript.TokenIdentifier, "b"),
						},
					},
				},
			},
			nil,
			true,
		},
		{ // 14
			[]javascript.StatementListItem{
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
			},
			[]javascript.AssignmentExpression{
				{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					}),
				},
			},
			true,
		},
	} {
		aes, ret := statementsListItemsAsExpressionsAndReturn(test.StatementListItems)
		if ret != test.Return {
			t.Errorf("test %d: expecting Return value of %v, got %v", n+1, test.Return, ret)
		} else if !reflect.DeepEqual(aes, test.Expressions) {
			t.Errorf("test %d: expecting AEs of %v, got %v", n+1, test.Expressions, aes)
		}
	}
}

func TestAEIsCE(t *testing.T) {
	for n, test := range [...]struct {
		AssignmentExpression    *javascript.AssignmentExpression
		IsConditionalExpression bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.AssignmentExpression{},
			false,
		},
		{ // 3
			&javascript.AssignmentExpression{
				AssignmentOperator: javascript.AssignmentAdd,
				AssignmentExpression: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					}),
				},
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
				}),
			},
			false,
		},
		{ // 4
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
				}),
			},
			true,
		},
	} {
		if aeIsCE(test.AssignmentExpression) != test.IsConditionalExpression {
			t.Errorf("test %d: expecting aeIsCE toreturn %v, got %v", n+1, test.IsConditionalExpression, !test.IsConditionalExpression)
		}
	}
}

func TestAEAsParen(t *testing.T) {
	for n, test := range [...]struct {
		AssignmentExpression    *javascript.AssignmentExpression
		ParenthesizedExpression *javascript.ParenthesizedExpression
	}{
		{ // 1
			nil,
			nil,
		},
		{ // 2
			&javascript.AssignmentExpression{},
			nil,
		},
		{ // 3
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
				}),
			},
			nil,
		},
		{ // 4
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					ParenthesizedExpression: &javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
								}),
							},
						},
					},
				}),
			},
			&javascript.ParenthesizedExpression{
				Expressions: []javascript.AssignmentExpression{
					{
						ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
						}),
					},
				},
			},
		},
		{ // 5
			&javascript.AssignmentExpression{
				AssignmentOperator: javascript.AssignmentAdd,
				AssignmentExpression: &javascript.AssignmentExpression{
					ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
						ParenthesizedExpression: &javascript.ParenthesizedExpression{
							Expressions: []javascript.AssignmentExpression{
								{
									ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
										IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
									}),
								},
							},
						},
					}),
				},
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					ParenthesizedExpression: &javascript.ParenthesizedExpression{
						Expressions: []javascript.AssignmentExpression{
							{
								ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
									IdentifierReference: makeToken(javascript.TokenIdentifier, "b"),
								}),
							},
						},
					},
				}),
			},
			nil,
		},
		{ // 6
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
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
				}),
			},
			&javascript.ParenthesizedExpression{
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
	} {
		if pe := aeAsParen(test.AssignmentExpression); !reflect.DeepEqual(pe, test.ParenthesizedExpression) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.ParenthesizedExpression, pe)
		}
	}
}

func TestMEIsSinglePE(t *testing.T) {
	for n, test := range [...]struct {
		MemberExpression *javascript.MemberExpression
		IsSinglePE       bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.MemberExpression{},
			false,
		},
		{ // 3
			&javascript.MemberExpression{
				PrimaryExpression: &javascript.PrimaryExpression{
					IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
				},
			},
			false,
		},
		{ // 4
			&javascript.MemberExpression{
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
			true,
		},
		{ // 5
			&javascript.MemberExpression{
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
			false,
		},
		{ // 6
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
			false,
		},
	} {
		if isPE := meIsSinglePe(test.MemberExpression); isPE != test.IsSinglePE {
			t.Errorf("test %d: expecting meIsSinglePe to return %v, got %v", n+1, test.IsSinglePE, !test.IsSinglePE)
		}
	}
}

func TestMEAsCE(t *testing.T) {
	for n, test := range [...]struct {
		MemberExpression *javascript.MemberExpression
		CallExpression   *javascript.CallExpression
	}{
		{ // 1
			nil,
			nil,
		},
		{ // 2
			&javascript.MemberExpression{},
			nil,
		},
		{ // 3
			&javascript.MemberExpression{
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
			&javascript.CallExpression{
				MemberExpression: &javascript.MemberExpression{
					PrimaryExpression: &javascript.PrimaryExpression{
						IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
					},
				},
				Arguments: &javascript.Arguments{},
			},
		},
		{ // 4
			&javascript.MemberExpression{
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
			nil,
		},
		{ // 5
			&javascript.MemberExpression{
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
			&javascript.CallExpression{
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
	} {
		if ce := meAsCE(test.MemberExpression); !reflect.DeepEqual(ce, test.CallExpression) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.MemberExpression, ce)
		}
	}
}

func TestIsStatementListItemExpression(t *testing.T) {
	for n, test := range [...]struct {
		StatementListItem             *javascript.StatementListItem
		IsStatementListItemExpression bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.StatementListItem{},
			false,
		},
		{ // 3
			&javascript.StatementListItem{
				Declaration: &javascript.Declaration{},
			},
			false,
		},
		{ // 4
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					Type: javascript.StatementDebugger,
				},
			},
			false,
		},
		{ // 5
			&javascript.StatementListItem{
				Statement: &javascript.Statement{
					Type: javascript.StatementDebugger,
				},
			},
			false,
		},
		{ // 6
			&javascript.StatementListItem{
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
			false,
		},
		{ // 7
			&javascript.StatementListItem{
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
			true,
		},
	} {
		if isStatementListItemExpression(test.StatementListItem) != test.IsStatementListItemExpression {
			t.Errorf("test %d: expecting return %v, got %v", n+1, test.IsStatementListItemExpression, !test.IsStatementListItemExpression)
		}
	}
}

func TestLeftMostLHS(t *testing.T) {
	for n, test := range [...]struct {
		Input  javascript.ConditionalWrappable
		Output *javascript.LeftHandSideExpression
	}{
		{ // 1
			nil,
			nil,
		},
		{ // 2
			javascript.WrapConditional(&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{},
			}),
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{},
			},
		},
		{ // 3
			javascript.WrapConditional(&javascript.AdditiveExpression{
				MultiplicativeExpression: javascript.MultiplicativeExpression{
					ExponentiationExpression: javascript.ExponentiationExpression{
						UnaryExpression: javascript.UnaryExpression{
							UpdateExpression: javascript.UpdateExpression{
								LeftHandSideExpression: &javascript.LeftHandSideExpression{
									CallExpression: &javascript.CallExpression{},
								},
							},
						},
					},
				},
			}),
			&javascript.LeftHandSideExpression{
				CallExpression: &javascript.CallExpression{},
			},
		},
		{ // 4
			javascript.WrapConditional(&javascript.ExponentiationExpression{
				ExponentiationExpression: &javascript.ExponentiationExpression{
					UnaryExpression: javascript.UnaryExpression{
						UpdateExpression: javascript.UpdateExpression{
							LeftHandSideExpression: &javascript.LeftHandSideExpression{
								NewExpression: &javascript.NewExpression{},
							},
						},
					},
				},
				UnaryExpression: javascript.UnaryExpression{
					UpdateExpression: javascript.UpdateExpression{
						LeftHandSideExpression: &javascript.LeftHandSideExpression{
							CallExpression: &javascript.CallExpression{},
						},
					},
				},
			}),
			&javascript.LeftHandSideExpression{
				NewExpression: &javascript.NewExpression{},
			},
		},
		{ // 5
			javascript.WrapConditional(&javascript.ExponentiationExpression{
				ExponentiationExpression: &javascript.ExponentiationExpression{
					ExponentiationExpression: &javascript.ExponentiationExpression{
						UnaryExpression: javascript.UnaryExpression{
							UpdateExpression: javascript.UpdateExpression{
								LeftHandSideExpression: &javascript.LeftHandSideExpression{
									OptionalExpression: &javascript.OptionalExpression{},
								},
							},
						},
					},
					UnaryExpression: javascript.UnaryExpression{
						UpdateExpression: javascript.UpdateExpression{
							LeftHandSideExpression: &javascript.LeftHandSideExpression{
								NewExpression: &javascript.NewExpression{},
							},
						},
					},
				},
				UnaryExpression: javascript.UnaryExpression{
					UpdateExpression: javascript.UpdateExpression{
						LeftHandSideExpression: &javascript.LeftHandSideExpression{
							CallExpression: &javascript.CallExpression{},
						},
					},
				},
			}),
			&javascript.LeftHandSideExpression{
				OptionalExpression: &javascript.OptionalExpression{},
			},
		},
	} {
		if out := leftMostLHS(test.Input); !reflect.DeepEqual(out, test.Output) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Output, out)
		}
	}
}

func TestFixWrapping(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output *javascript.Statement
	}{
		{ // 1
			nil,
			nil,
		},
		{ // 2
			&javascript.Statement{},
			&javascript.Statement{},
		},
		{ // 3
			&javascript.Statement{
				BlockStatement: &javascript.Block{},
			},
			&javascript.Statement{
				BlockStatement: &javascript.Block{},
			},
		},
		{ // 4
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{},
			},
		},
		{ // 5
			&javascript.Statement{
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
			&javascript.Statement{
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
		{ // 6
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								ObjectLiteral: &javascript.ObjectLiteral{},
							}),
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												ObjectLiteral: &javascript.ObjectLiteral{},
											}),
										},
									},
								},
							}),
						},
					},
				},
			},
		},
		{ // 7
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								FunctionExpression: &javascript.FunctionDeclaration{},
							}),
						},
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								ObjectLiteral: &javascript.ObjectLiteral{},
							}),
						},
					},
				},
			},
			&javascript.Statement{
				ExpressionStatement: &javascript.Expression{
					Expressions: []javascript.AssignmentExpression{
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								ParenthesizedExpression: &javascript.ParenthesizedExpression{
									Expressions: []javascript.AssignmentExpression{
										{
											ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
												FunctionExpression: &javascript.FunctionDeclaration{},
											}),
										},
									},
								},
							}),
						},
						{
							ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
								ObjectLiteral: &javascript.ObjectLiteral{},
							}),
						},
					},
				},
			},
		},
	} {
		fixWrapping(test.Input)

		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Output, test.Input)
		}
	}
}

func TestScoreCE(t *testing.T) {
	for n, test := range [...]struct {
		Input  javascript.ConditionalWrappable
		Output int
	}{
		{ // 1
			nil,
			-1,
		},
		{ // 2
			&javascript.LogicalORExpression{},
			1,
		},
		{ // 3
			&javascript.EqualityExpression{},
			6,
		},
		{ // 4
			&javascript.UpdateExpression{},
			13,
		},
	} {
		if out := scoreCE(test.Input); out != test.Output {
			t.Errorf("test %d: expecting %d, got %d", n+1, test.Output, out)
		}
	}
}

func TestIsConditionalWrappingAConditional(t *testing.T) {
	for n, test := range [...]struct {
		Input, Below javascript.ConditionalWrappable
		Output       *javascript.ConditionalExpression
	}{
		{ // 1
			nil,
			nil,
			nil,
		},
		{ // 2
			&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
			},
			nil,
			nil,
		},
		{ // 3
			&javascript.PrimaryExpression{
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
			nil,
			javascript.WrapConditional(&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
			}),
		},
		{ // 4
			&javascript.PrimaryExpression{
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
			nil,
			nil,
		},
		{ // 5
			&javascript.MultiplicativeExpression{
				ExponentiationExpression: javascript.ExponentiationExpression{
					UnaryExpression: javascript.UnaryExpression{
						UpdateExpression: javascript.UpdateExpression{
							LeftHandSideExpression: &javascript.LeftHandSideExpression{
								NewExpression: &javascript.NewExpression{
									MemberExpression: javascript.MemberExpression{
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
								},
							},
						},
					},
				},
			},
			nil,
			javascript.WrapConditional(&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
			}),
		},
		{ // 6
			&javascript.MultiplicativeExpression{
				ExponentiationExpression: javascript.ExponentiationExpression{
					UnaryExpression: javascript.UnaryExpression{
						UpdateExpression: javascript.UpdateExpression{
							LeftHandSideExpression: &javascript.LeftHandSideExpression{
								NewExpression: &javascript.NewExpression{
									MemberExpression: javascript.MemberExpression{
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
								},
							},
						},
					},
				},
			},
			&javascript.AdditiveExpression{},
			javascript.WrapConditional(&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
			}),
		},
		{ // 7
			&javascript.MultiplicativeExpression{
				ExponentiationExpression: javascript.ExponentiationExpression{
					UnaryExpression: javascript.UnaryExpression{
						UpdateExpression: javascript.UpdateExpression{
							LeftHandSideExpression: &javascript.LeftHandSideExpression{
								NewExpression: &javascript.NewExpression{
									MemberExpression: javascript.MemberExpression{
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
								},
							},
						},
					},
				},
			},
			&javascript.ExponentiationExpression{},
			javascript.WrapConditional(&javascript.PrimaryExpression{
				IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
			}),
		},
	} {
		if out := isConditionalWrappingAConditional(test.Input, test.Below); !reflect.DeepEqual(out, test.Output) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Output, out)
		}
	}
}

func TestRemoveLastReturnStatement(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output *javascript.Block
	}{
		{ // 1
			nil,
			nil,
		},
		{ // 2
			&javascript.Block{},
			&javascript.Block{},
		},
		{ // 3
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenLineTerminator, "a"),
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
											IdentifierReference: makeToken(javascript.TokenLineTerminator, "a"),
										}),
									},
								},
							},
						},
					},
				},
			},
		},
		{ // 4
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							Type: javascript.StatementReturn,
						},
					},
				},
			},
			&javascript.Block{
				StatementList: []javascript.StatementListItem{},
			},
		},
		{ // 5
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenLineTerminator, "a"),
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
			&javascript.Block{
				StatementList: []javascript.StatementListItem{
					{
						Statement: &javascript.Statement{
							ExpressionStatement: &javascript.Expression{
								Expressions: []javascript.AssignmentExpression{
									{
										ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
											IdentifierReference: makeToken(javascript.TokenLineTerminator, "a"),
										}),
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		removeLastReturnStatement(test.Input)

		if !reflect.DeepEqual(test.Input, test.Output) {
			t.Errorf("test %d: expecting %b, got %v", n+1, test.Output, test.Input)
		}
	}
}

func TestIsSimpleAE(t *testing.T) {
	for n, test := range [...]struct {
		Input  *javascript.AssignmentExpression
		Output bool
	}{
		{ // 1
			nil,
			false,
		},
		{ // 2
			&javascript.AssignmentExpression{},
			false,
		},
		{ // 3
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
				}),
			},
			true,
		},
		{ // 4
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.PrimaryExpression{
					Literal: makeToken(javascript.TokenNumericLiteral, "1"),
				}),
			},
			true,
		},
		{ // 5
			&javascript.AssignmentExpression{
				ConditionalExpression: javascript.WrapConditional(&javascript.CallExpression{
					MemberExpression: &javascript.MemberExpression{
						PrimaryExpression: &javascript.PrimaryExpression{
							IdentifierReference: makeToken(javascript.TokenIdentifier, "a"),
						},
					},
					Arguments: &javascript.Arguments{},
				}),
			},
			false,
		},
	} {
		if out := isSimpleAE(test.Input); out != test.Output {
			t.Errorf("test %d: expecting output %v, got %v", n+1, test.Output, out)
		}
	}
}
