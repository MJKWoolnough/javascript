package javascript

import "testing"

func TestClassDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{`class myClass{}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Tokens: tk[3:3],
				},
				Tokens: tk[:5],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseClassDeclaration(t.Yield, t.Await, t.Def)
	})
}
