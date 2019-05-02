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
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseNewExpression(t.Yield, t.Await)
	})
}
