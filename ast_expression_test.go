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
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseNewExpression(t.Yield, t.Await)
	})
}
