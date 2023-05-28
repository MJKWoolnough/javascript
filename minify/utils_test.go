package minify

import (
	"reflect"
	"testing"

	"vimagination.zapto.org/javascript"
)

func TestBlockAsModule(t *testing.T) {
	for n, test := range [...]struct {
		Input    *javascript.Block
		Callback func(*javascript.Module)
		Output   *javascript.Block
	}{
		{ // 1
			&javascript.Block{},
			func(m *javascript.Module) {
			},
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
			func(m *javascript.Module) {
				for i := 0; i < len(m.ModuleListItems); i++ {
					if m.ModuleListItems[i].StatementListItem.Statement.Type == javascript.StatementDebugger {
						m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
						i--
					}
				}
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
		Callback func(*javascript.Module)
		Output   *javascript.Expression
	}{
		{ // 1
			&javascript.Expression{},
			func(m *javascript.Module) {
			},
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
			func(m *javascript.Module) {
				for i := 0; i < len(m.ModuleListItems); i++ {
					if javascript.UnwrapConditional(m.ModuleListItems[i].StatementListItem.Statement.ExpressionStatement.Expressions[0].ConditionalExpression).(*javascript.PrimaryExpression).IdentifierReference.Data == "b" {
						m.ModuleListItems = append(m.ModuleListItems[:i], m.ModuleListItems[i+1:]...)
					}
				}
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
		{ // 3
			&javascript.Statement{
				Type: javascript.StatementDebugger,
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
			false,
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
		{ // 5
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
		{ // 5
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
