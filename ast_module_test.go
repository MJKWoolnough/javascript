package javascript

import "testing"

func TestModule(t *testing.T) {
	doTests(t, []sourceFn{
		{`import 'a';`, func(t *test, tk Tokens) { // 1
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							FromClause: FromClause{
								ModuleSpecifier: &tk[2],
								Tokens:          tk[2:3],
							},
							Tokens: tk[:4],
						},
						Tokens: tk[:4],
					},
				},
				Tokens: tk[:4],
			}
		}},
		{`import a from 'b';`, func(t *test, tk Tokens) { // 2
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &ImportedBinding{
									Identifier: &tk[2],
								},
								Tokens: tk[2:3],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[6],
								Tokens:          tk[4:7],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{`import * as a from 'b';`, func(t *test, tk Tokens) { // 3
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NameSpaceImport: &ImportedBinding{
									Identifier: &tk[6],
								},
								Tokens: tk[2:7],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[10],
								Tokens:          tk[8:11],
							},
							Tokens: tk[:12],
						},
						Tokens: tk[:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{`import {a} from 'b';`, func(t *test, tk Tokens) { // 4
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											ImportedBinding: ImportedBinding{
												Identifier: &tk[3],
											},
											Tokens: tk[3:4],
										},
									},
									Tokens: tk[2:5],
								},
								Tokens: tk[2:5],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[8],
								Tokens:          tk[6:9],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`import {a as b} from 'c';`, func(t *test, tk Tokens) { // 5
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											IdentifierName: &tk[3],
											ImportedBinding: ImportedBinding{
												Identifier: &tk[7],
											},
											Tokens: tk[3:8],
										},
									},
									Tokens: tk[2:9],
								},
								Tokens: tk[2:9],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[12],
								Tokens:          tk[10:13],
							},
							Tokens: tk[:14],
						},
						Tokens: tk[:14],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`import {a as b, c} from 'd';`, func(t *test, tk Tokens) { // 6
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											IdentifierName: &tk[3],
											ImportedBinding: ImportedBinding{
												Identifier: &tk[7],
											},
											Tokens: tk[3:8],
										},
										{
											ImportedBinding: ImportedBinding{
												Identifier: &tk[10],
											},
											Tokens: tk[10:11],
										},
									},
									Tokens: tk[2:12],
								},
								Tokens: tk[2:12],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[15],
								Tokens:          tk[13:16],
							},
							Tokens: tk[:17],
						},
						Tokens: tk[:17],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{`import a, * as b from 'c';`, func(t *test, tk Tokens) { // 7
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &ImportedBinding{
									Identifier: &tk[2],
								},
								NameSpaceImport: &ImportedBinding{
									Identifier: &tk[9],
								},
								Tokens: tk[2:10],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[13],
								Tokens:          tk[11:14],
							},
							Tokens: tk[:15],
						},
						Tokens: tk[:15],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{`import a, {b} from 'c';`, func(t *test, tk Tokens) { // 8
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &ImportedBinding{
									Identifier: &tk[2],
								},
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											ImportedBinding: ImportedBinding{
												Identifier: &tk[6],
											},
											Tokens: tk[6:7],
										},
									},
									Tokens: tk[5:8],
								},
								Tokens: tk[2:8],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[11],
								Tokens:          tk[9:12],
							},
							Tokens: tk[:13],
						},
						Tokens: tk[:13],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{`export {};`, func(t *test, tk Tokens) { // 9
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								Tokens: tk[2:4],
							},
							Tokens: tk[:5],
						},
						Tokens: tk[:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{`export {a};`, func(t *test, tk Tokens) { // 10
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName: &tk[3],
										Tokens:         tk[3:4],
									},
								},
								Tokens: tk[2:5],
							},
							Tokens: tk[:6],
						},
						Tokens: tk[:6],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{`export {a as b};`, func(t *test, tk Tokens) { // 11
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[7],
										Tokens:          tk[3:8],
									},
								},
								Tokens: tk[2:9],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`export {a as b, c as d, e, f};`, func(t *test, tk Tokens) { // 12
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[7],
										Tokens:          tk[3:8],
									},
									{
										IdentifierName:  &tk[10],
										EIdentifierName: &tk[14],
										Tokens:          tk[10:15],
									},
									{
										IdentifierName: &tk[17],
										Tokens:         tk[17:18],
									},
									{
										IdentifierName: &tk[20],
										Tokens:         tk[20:21],
									},
								},
								Tokens: tk[2:22],
							},
							Tokens: tk[:23],
						},
						Tokens: tk[:23],
					},
				},
				Tokens: tk[:23],
			}
		}},
		{`export * from 'a';`, func(t *test, tk Tokens) { // 13
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							FromClause: &FromClause{
								ModuleSpecifier: &tk[6],
								Tokens:          tk[4:7],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{`export {} from 'a';`, func(t *test, tk Tokens) { // 14
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								Tokens: tk[2:4],
							},
							FromClause: &FromClause{
								ModuleSpecifier: &tk[7],
								Tokens:          tk[5:8],
							},
							Tokens: tk[:9],
						},
						Tokens: tk[:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{`export {a} from 'b';`, func(t *test, tk Tokens) { // 15
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName: &tk[3],
										Tokens:         tk[3:4],
									},
								},
								Tokens: tk[2:5],
							},
							FromClause: &FromClause{
								ModuleSpecifier: &tk[8],
								Tokens:          tk[6:9],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`export {a as b} from 'c';`, func(t *test, tk Tokens) { // 16
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[7],
										Tokens:          tk[3:8],
									},
								},
								Tokens: tk[2:9],
							},
							FromClause: &FromClause{
								ModuleSpecifier: &tk[12],
								Tokens:          tk[10:13],
							},
							Tokens: tk[:14],
						},
						Tokens: tk[:14],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`export {a as b, c as d, e, f} from 'g';`, func(t *test, tk Tokens) { // 17
			t.Output = Module{
				ModuleListItems: []ModuleListItem{
					{
						ExportDeclaration: &ExportDeclaration{
							ExportClause: &ExportClause{
								ExportList: []ExportSpecifier{
									{
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[7],
										Tokens:          tk[3:8],
									},
									{
										IdentifierName:  &tk[10],
										EIdentifierName: &tk[14],
										Tokens:          tk[10:15],
									},
									{
										IdentifierName: &tk[17],
										Tokens:         tk[17:18],
									},
									{
										IdentifierName: &tk[20],
										Tokens:         tk[20:21],
									},
								},
								Tokens: tk[2:22],
							},
							FromClause: &FromClause{
								ModuleSpecifier: &tk[25],
								Tokens:          tk[23:26],
							},
							Tokens: tk[:27],
						},
						Tokens: tk[:27],
					},
				},
				Tokens: tk[:27],
			}
		}},
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseModule()
	})
}
