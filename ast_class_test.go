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
		{`class myClass extends OtherClass{}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				Extends: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: &IdentifierReference{
									Identifier: &tk[6],
								},
								Tokens: tk[6:7],
							},
							Tokens: tk[6:7],
						},
						Tokens: tk[6:7],
					},
					Tokens: tk[6:7],
				},
				ClassBody: ClassBody{
					Tokens: tk[7:7],
				},
				Tokens: tk[:9],
			}
		}},
		{`class myClass {constructor(){}}`, func(t *test, tk Tokens) {
			t.Output = ClassDeclaration{
				BindingIdentifier: &BindingIdentifier{Identifier: &tk[2]},
				ClassBody: ClassBody{
					Methods: []MethodDefinition{
						{
							PropertyName: PropertyName{
								LiteralPropertyName: &tk[5],
								Tokens:              tk[5:6],
							},
							Params: FormalParameters{
								Tokens: tk[7:7],
							},
							FunctionBody: Block{
								Tokens: tk[8:10],
							},
							Tokens: tk[5:10],
						},
					},
					Tokens: tk[5:10],
				},
				Tokens: tk[:11],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseClassDeclaration(t.Yield, t.Await, t.Def)
	})
}
