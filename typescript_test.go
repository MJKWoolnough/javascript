package javascript

import "testing"

func TestTypescriptModule(t *testing.T) {
	doTests(t, []sourceFn{
		{`import def from './a';import type typeDef from './b';import type {typ1, typ2} from './c';import {a} from './d';`, func(t *test, tk Tokens) { // 1
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &tk[2],
								Tokens:                 tk[2:3],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[6],
								Tokens:          tk[4:7],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[17:18],
							},
							Tokens: tk[17:18],
						},
						Tokens: tk[8:18],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[32:33],
							},
							Tokens: tk[32:33],
						},
						Tokens: tk[18:33],
					},
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											IdentifierName:  &tk[36],
											ImportedBinding: &tk[36],
											Tokens:          tk[36:37],
										},
									},
									Tokens: tk[35:38],
								},
								Tokens: tk[35:38],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[41],
								Tokens:          tk[39:42],
							},
							Tokens: tk[33:43],
						},
						Tokens: tk[33:43],
					},
				},
				Tokens: tk[:43],
			}
		}},
	}, func(t *test) (Type, error) {
		t.Tokens[:cap(t.Tokens)][cap(t.Tokens)-1].Data = marker
		var m Module
		err := m.parse(&t.Tokens)
		return m, err
	})
}
