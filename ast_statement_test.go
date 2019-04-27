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
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseStatement(t.Yield, t.Await, t.Ret)
	})
}
