package javascript

import "testing"

func TestModuleOld(t *testing.T) {
	doTests(t, []sourceFn{
		{`import 'a';`, func(t *test, tk Tokens) { // 1
			t.Output = Module{
				ModuleListItems: []ModuleItem{
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
				},
				Tokens: tk[:8],
			}
		}},
		{`import * as a from 'b';`, func(t *test, tk Tokens) { // 3
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NameSpaceImport: &tk[6],
								Tokens:          tk[2:7],
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
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											ImportedBinding: &tk[3],
											Tokens:          tk[3:4],
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
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											IdentifierName:  &tk[3],
											ImportedBinding: &tk[7],
											Tokens:          tk[3:8],
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
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											IdentifierName:  &tk[3],
											ImportedBinding: &tk[7],
											Tokens:          tk[3:8],
										},
										{
											ImportedBinding: &tk[10],
											Tokens:          tk[10:11],
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
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &tk[2],
								NameSpaceImport:        &tk[9],
								Tokens:                 tk[2:10],
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
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								ImportedDefaultBinding: &tk[2],
								NamedImports: &NamedImports{
									ImportList: []ImportSpecifier{
										{
											ImportedBinding: &tk[6],
											Tokens:          tk[6:7],
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
				ModuleListItems: []ModuleItem{
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
		{`export var a = 1;`, func(t *test, tk Tokens) { // 18
			litA := makeConditionLiteral(tk, 8)
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							VariableStatement: &VariableStatement{
								VariableDeclarationList: []VariableDeclaration{
									{
										BindingIdentifier: &tk[4],
										Initializer: &AssignmentExpression{
											ConditionalExpression: &litA,
											Tokens:                tk[8:9],
										},
										Tokens: tk[4:9],
									},
								},
								Tokens: tk[2:10],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`export function a(){}`, func(t *test, tk Tokens) { // 19
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							Declaration: &Declaration{
								FunctionDeclaration: &FunctionDeclaration{
									BindingIdentifier: &tk[4],
									FormalParameters: FormalParameters{
										Tokens: tk[5:7],
									},
									FunctionBody: Block{
										Tokens: tk[7:9],
									},
									Tokens: tk[2:9],
								},
								Tokens: tk[2:9],
							},
							Tokens: tk[:9],
						},
						Tokens: tk[:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{`export default function(){}`, func(t *test, tk Tokens) { // 20
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							DefaultFunction: &FunctionDeclaration{
								FormalParameters: FormalParameters{
									Tokens: tk[5:7],
								},
								FunctionBody: Block{
									Tokens: tk[7:9],
								},
								Tokens: tk[4:9],
							},
							Tokens: tk[:9],
						},
						Tokens: tk[:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{`export default class{}`, func(t *test, tk Tokens) { // 21
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							DefaultClass: &ClassDeclaration{
								Tokens: tk[4:7],
							},
							Tokens: tk[:7],
						},
						Tokens: tk[:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{`export default 1;`, func(t *test, tk Tokens) { // 22
			litA := makeConditionLiteral(tk, 4)
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							DefaultAssignmentExpression: &AssignmentExpression{
								ConditionalExpression: &litA,
								Tokens:                tk[4:5],
							},
							Tokens: tk[:6],
						},
						Tokens: tk[:6],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{`1;`, func(t *test, tk Tokens) { // 23
			litA := makeConditionLiteral(tk, 0)
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &litA,
											Tokens:                tk[:1],
										},
									},
									Tokens: tk[:1],
								},
								Tokens: tk[:2],
							},
							Tokens: tk[:2],
						},
						Tokens: tk[:2],
					},
				},
				Tokens: tk[:2],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var m Module
		err := m.parse(&t.Tokens)
		return m, err
	})
}

func TestModule(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Output = Module{
				Tokens: tk[:0],
			}
		}},
		{`import`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidImport,
						Parsing: "ImportClause",
						Token:   tk[1],
					},
					Parsing: "ImportDeclaration",
					Token:   tk[1],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{`;;`, func(t *test, tk Tokens) { // 3
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:1],
							},
							Tokens: tk[:1],
						},
						Tokens: tk[:1],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[1:2],
							},
							Tokens: tk[1:2],
						},
						Tokens: tk[1:2],
					},
				},
				Tokens: tk[:2],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var m Module
		err := m.parse(&t.Tokens)
		return m, err
	})
}

func TestModuleItem(t *testing.T) {
	doTests(t, []sourceFn{ // 1
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     assignmentError(tk[0]),
							Parsing: "Expression",
							Token:   tk[0],
						},
						Parsing: "Statement",
						Token:   tk[0],
					},
					Parsing: "StatementListItem",
					Token:   tk[0],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{"import", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidImport,
						Parsing: "ImportClause",
						Token:   tk[1],
					},
					Parsing: "ImportDeclaration",
					Token:   tk[1],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{"import\n'a';", func(t *test, tk Tokens) { // 3
			t.Output = ModuleItem{
				ImportDeclaration: &ImportDeclaration{
					FromClause: FromClause{
						ModuleSpecifier: &tk[2],
						Tokens:          tk[2:3],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"export", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidDeclaration,
						Parsing: "Declaration",
						Token:   tk[1],
					},
					Parsing: "ExportDeclaration",
					Token:   tk[1],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{"export\n*\nfrom\n'a';", func(t *test, tk Tokens) { // 5
			t.Output = ModuleItem{
				ExportDeclaration: &ExportDeclaration{
					FromClause: &FromClause{
						ModuleSpecifier: &tk[6],
						Tokens:          tk[4:7],
					},
					Tokens: tk[:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"var", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrNoIdentifier,
								Parsing: "LexicalBinding",
								Token:   tk[1],
							},
							Parsing: "VariableStatement",
							Token:   tk[1],
						},
						Parsing: "Statement",
						Token:   tk[0],
					},
					Parsing: "StatementListItem",
					Token:   tk[0],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{"var\na;", func(t *test, tk Tokens) { // 7
			t.Output = ModuleItem{
				StatementListItem: &StatementListItem{
					Statement: &Statement{
						VariableStatement: &VariableStatement{
							VariableDeclarationList: []VariableDeclaration{
								{
									BindingIdentifier: &tk[2],
									Tokens:            tk[2:3],
								},
							},
							Tokens: tk[:4],
						},
						Tokens: tk[:4],
					},
					Tokens: tk[:4],
				},
				Tokens: tk[:4],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var mi ModuleItem
		err := mi.parse(&t.Tokens)
		return mi, err
	})
}

func TestImportDeclaration(t *testing.T) {
	doTests(t, []sourceFn{ // 1
		{``, func(t *test, tk Tokens) {
			t.Err = Error{
				Err:     ErrInvalidImport,
				Parsing: "ImportDeclaration",
				Token:   tk[0],
			}
		}},
		{"import", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidImport,
					Parsing: "ImportClause",
					Token:   tk[1],
				},
				Parsing: "ImportDeclaration",
				Token:   tk[1],
			}
		}},
		{"import\n\"\";", func(t *test, tk Tokens) { // 3
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[2],
					Tokens:          tk[2:3],
				},
				Tokens: tk[:4],
			}
		}},
		{"import\n*\n", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidNameSpaceImport,
					Parsing: "ImportClause",
					Token:   tk[4],
				},
				Parsing: "ImportDeclaration",
				Token:   tk[2],
			}
		}},
		{"import\n*\nas\na\n", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingFrom,
					Parsing: "FromClause",
					Token:   tk[8],
				},
				Parsing: "ImportDeclaration",
				Token:   tk[8],
			}
		}},
		{"import\n*\nas\na\nfrom\n\"\";", func(t *test, tk Tokens) { // 6
			t.Output = ImportDeclaration{
				ImportClause: &ImportClause{
					NameSpaceImport: &tk[6],
					Tokens:          tk[2:7],
				},
				FromClause: FromClause{
					ModuleSpecifier: &tk[10],
					Tokens:          tk[8:11],
				},
				Tokens: tk[:12],
			}
		}},
		{"import\n\"\"", func(t *test, tk Tokens) { // 7
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[2],
					Tokens:          tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"import\n\"\"\na", func(t *test, tk Tokens) { // 8
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[2],
					Tokens:          tk[2:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"import\n\"\" a", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "ImportDeclaration",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var id ImportDeclaration
		err := id.parse(&t.Tokens)
		return id, err
	})
}

func TestImportClause(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidImport,
				Parsing: "ImportClause",
				Token:   tk[0],
			}
		}},
		{`for`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ImportClause",
				Token:   tk[0],
			}
		}},
		{"a\n", func(t *test, tk Tokens) { // 3
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[0],
				Tokens:                 tk[:1],
			}
		}},
		{"a\n,", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrInvalidImport,
				Parsing: "ImportClause",
				Token:   tk[3],
			}
		}},
		{"a\n,\n*", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrInvalidNameSpaceImport,
				Parsing: "ImportClause",
				Token:   tk[5],
			}

		}},
		{"a\n,\n*\nas", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ImportClause",
				Token:   tk[7],
			}

		}},
		{"a\n,\n*\nas\nb", func(t *test, tk Tokens) { // 7
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[0],
				NameSpaceImport:        &tk[8],
				Tokens:                 tk[:9],
			}
		}},
		{",", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrInvalidImport,
				Parsing: "ImportClause",
				Token:   tk[0],
			}
		}},
		{"*", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrInvalidNameSpaceImport,
				Parsing: "ImportClause",
				Token:   tk[1],
			}

		}},
		{"*\nas", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ImportClause",
				Token:   tk[3],
			}

		}},
		{"*\nas\nb", func(t *test, tk Tokens) { // 11
			t.Output = ImportClause{
				NameSpaceImport: &tk[4],
				Tokens:          tk[:5],
			}
		}},
		{"a\n,\n{+}", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidImportSpecifier,
						Parsing: "ImportSpecifier",
						Token:   tk[5],
					},
					Parsing: "NamedImports",
					Token:   tk[5],
				},
				Parsing: "ImportClause",
				Token:   tk[4],
			}
		}},
		{"a\n,\n{}", func(t *test, tk Tokens) { // 13
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[0],
				NamedImports: &NamedImports{
					Tokens: tk[4:6],
				},
				Tokens: tk[:6],
			}
		}},
		{"{+}", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrInvalidImportSpecifier,
						Parsing: "ImportSpecifier",
						Token:   tk[1],
					},
					Parsing: "NamedImports",
					Token:   tk[1],
				},
				Parsing: "ImportClause",
				Token:   tk[0],
			}
		}},
		{"{}", func(t *test, tk Tokens) { // 15
			t.Output = ImportClause{
				NamedImports: &NamedImports{
					Tokens: tk[:2],
				},
				Tokens: tk[:2],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var ic ImportClause
		err := ic.parse(&t.Tokens)
		return ic, err
	})
}

func TestFromClause(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingFrom,
				Parsing: "FromClause",
				Token:   tk[0],
			}
		}},
		{"from\n", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingModuleSpecifier,
				Parsing: "FromClause",
				Token:   tk[2],
			}
		}},
		{"from\n\"\"", func(t *test, tk Tokens) { // 3
			t.Output = FromClause{
				ModuleSpecifier: &tk[2],
				Tokens:          tk[:3],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var fc FromClause
		err := fc.parse(&t.Tokens)
		return fc, err
	})
}

func TestNamedImports(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[0],
			}
		}},
		{"{\n}", func(t *test, tk Tokens) { // 2
			t.Output = NamedImports{
				Tokens: tk[:3],
			}
		}},
		{"{\n,}", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     ErrInvalidImportSpecifier,
					Parsing: "ImportSpecifier",
					Token:   tk[2],
				},
				Parsing: "NamedImports",
				Token:   tk[2],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 4
			t.Output = NamedImports{
				ImportList: []ImportSpecifier{
					{
						ImportedBinding: &tk[2],
						Tokens:          tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"{\na\n,\n}", func(t *test, tk Tokens) { // 5
			t.Output = NamedImports{
				ImportList: []ImportSpecifier{
					{
						ImportedBinding: &tk[2],
						Tokens:          tk[2:3],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"{\na\nb}", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[4],
			}
		}},
		{"{\na\n,\nb\n}", func(t *test, tk Tokens) { // 7
			t.Output = NamedImports{
				ImportList: []ImportSpecifier{
					{
						ImportedBinding: &tk[2],
						Tokens:          tk[2:3],
					},
					{
						ImportedBinding: &tk[6],
						Tokens:          tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
	}, func(t *test) (interface{}, error) {
		var ni NamedImports
		err := ni.parse(&t.Tokens)
		return ni, err
	})
}
