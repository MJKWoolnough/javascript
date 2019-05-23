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
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseStatement(t.Yield, t.Await, t.Ret)
	})
}
