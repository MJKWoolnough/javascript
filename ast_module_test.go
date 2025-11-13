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
											IdentifierName:  &tk[3],
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
											IdentifierName:  &tk[10],
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
											IdentifierName:  &tk[6],
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
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[3],
										Tokens:          tk[3:4],
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
										IdentifierName:  &tk[17],
										EIdentifierName: &tk[17],
										Tokens:          tk[17:18],
									},
									{
										IdentifierName:  &tk[20],
										EIdentifierName: &tk[20],
										Tokens:          tk[20:21],
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
										IdentifierName:  &tk[3],
										EIdentifierName: &tk[3],
										Tokens:          tk[3:4],
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
										IdentifierName:  &tk[17],
										EIdentifierName: &tk[17],
										Tokens:          tk[17:18],
									},
									{
										IdentifierName:  &tk[20],
										EIdentifierName: &tk[20],
										Tokens:          tk[20:21],
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
	}, func(t *test) (Type, error) {
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
		{`await a`, func(t *test, tk Tokens) { // 4
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(UnaryExpression{
												UnaryOperators: []UnaryOperatorComments{
													{
														UnaryOperator: UnaryAwait,
													},
												},
												UpdateExpression: UpdateExpression{
													LeftHandSideExpression: &LeftHandSideExpression{
														NewExpression: &NewExpression{
															MemberExpression: MemberExpression{
																PrimaryExpression: &PrimaryExpression{
																	IdentifierReference: &tk[2],
																	Tokens:              tk[2:3],
																},
																Tokens: tk[2:3],
															},
															Tokens: tk[2:3],
														},
														Tokens: tk[2:3],
													},
													Tokens: tk[2:3],
												},
												Tokens: tk[:3],
											}),
											Tokens: tk[:3],
										},
									},
									Tokens: tk[:3],
								},
								Tokens: tk[:3],
							},
							Tokens: tk[:3],
						},
						Tokens: tk[:3],
					},
				},
				Tokens: tk[:3],
			}
		}},
		{"// A\na\n// B", func(t *test, tk Tokens) { // 5
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&PrimaryExpression{
												IdentifierReference: &tk[2],
												Tokens:              tk[2:3],
											}),
											Tokens: tk[2:3],
										},
									},
									Tokens: tk[2:3],
								},
								Tokens: tk[2:3],
							},
							Tokens: tk[2:3],
						},
						Tokens: tk[2:3],
					},
				},
				Comments: [2]Comments{{tk[0]}, {tk[4]}},
				Tokens:   tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
		var m Module

		err := m.parse(&t.Tokens)

		return m, err
	})
}

func TestModuleItem(t *testing.T) {
	doTests(t, []sourceFn{
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
		{"import.meta", func(t *test, tk Tokens) { // 8
			t.Output = ModuleItem{
				StatementListItem: &StatementListItem{
					Statement: &Statement{
						ExpressionStatement: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: WrapConditional(MemberExpression{
										ImportMeta: true,
										Tokens:     tk[:3],
									}),
									Tokens: tk[:3],
								},
							},
							Tokens: tk[:3],
						},
						Tokens: tk[:3],
					},
					Tokens: tk[:3],
				},
				Tokens: tk[:3],
			}
		}},
		{"import(a)", func(t *test, tk Tokens) { // 9
			t.Output = ModuleItem{
				StatementListItem: &StatementListItem{
					Statement: &Statement{
						ExpressionStatement: &Expression{
							Expressions: []AssignmentExpression{
								{
									ConditionalExpression: WrapConditional(&CallExpression{
										ImportCall: &AssignmentExpression{
											ConditionalExpression: WrapConditional(&PrimaryExpression{
												IdentifierReference: &tk[2],
												Tokens:              tk[2:3],
											}),
											Tokens: tk[2:3],
										},
										Tokens: tk[:4],
									}),
									Tokens: tk[:4],
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
	}, func(t *test) (Type, error) {
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
		{"import a from 'b' with {c:'d'}", func(t *test, tk Tokens) { // 10
			t.Output = ImportDeclaration{
				ImportClause: &ImportClause{
					ImportedDefaultBinding: &tk[2],
					Tokens:                 tk[2:3],
				},
				FromClause: FromClause{
					ModuleSpecifier: &tk[6],
					Tokens:          tk[4:7],
				},
				WithClause: &WithClause{
					WithEntries: []WithEntry{
						{
							AttributeKey: &tk[11],
							Value:        &tk[13],
							Tokens:       tk[11:14],
						},
					},
					Tokens: tk[8:15],
				},
				Tokens: tk[:15],
			}
		}},
		{"// A\nimport \"a\"; // B\n\n// C", func(t *test, tk Tokens) { // 11
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[4],
					Tokens:          tk[4:5],
				},
				Comments: [4]Comments{{tk[0]}, nil, nil, {tk[7]}},
				Tokens:   tk[:8],
			}
		}},
		{"// A\nimport /* B */\"a\" /* C */; // B\n\n// C\n", func(t *test, tk Tokens) { // 12
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[5],
					Tokens:          tk[5:6],
				},
				Comments: [4]Comments{{tk[0]}, {tk[4]}, {tk[7]}, {tk[10]}},
				Tokens:   tk[:11],
			}
		}},
		{"// A\nimport /* B */ a /* C */ from // D\n'b' /* E */ with /* F */ {c:'d'} /* G */; // G\n", func(t *test, tk Tokens) { // 13
			t.Output = ImportDeclaration{
				ImportClause: &ImportClause{
					ImportedDefaultBinding: &tk[6],
					Comments:               [6]Comments{{tk[4]}, nil, nil, nil, nil, {tk[8]}},
					Tokens:                 tk[4:9],
				},
				FromClause: FromClause{
					ModuleSpecifier: &tk[14],
					Comments:        Comments{tk[12]},
					Tokens:          tk[10:15],
				},
				WithClause: &WithClause{
					WithEntries: []WithEntry{
						{
							AttributeKey: &tk[23],
							Value:        &tk[25],
							Tokens:       tk[23:26],
						},
					},
					Comments: [4]Comments{{tk[16]}, {tk[20]}},
					Tokens:   tk[16:27],
				},
				Comments: [4]Comments{{tk[0]}, nil, {tk[28]}, {tk[31]}},
				Tokens:   tk[:32],
			}
		}},
		{"// A\nimport /* B */ \"\" /* C */; // D\n", func(t *test, tk Tokens) { // 14
			t.Output = ImportDeclaration{
				FromClause: FromClause{
					ModuleSpecifier: &tk[6],
					Tokens:          tk[6:7],
				},
				Comments: [4]Comments{{tk[0]}, {tk[4]}, {tk[8]}, {tk[11]}},
				Tokens:   tk[:12],
			}
		}},
	}, func(t *test) (Type, error) {
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
		{"// A\na // B\n", func(t *test, tk Tokens) { // 16
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[2],
				Comments:               [6]Comments{{tk[0]}, nil, nil, nil, nil, {tk[4]}},
				Tokens:                 tk[:5],
			}
		}},
		{"// A\n{} // B\n", func(t *test, tk Tokens) { // 17
			t.Output = ImportClause{
				NamedImports: &NamedImports{
					Tokens: tk[2:4],
				},
				Comments: [6]Comments{{tk[0]}, nil, nil, nil, nil, {tk[5]}},
				Tokens:   tk[:6],
			}
		}},
		{"// A\na /* B */,/* C */ {} /* D */", func(t *test, tk Tokens) { // 18
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[2],
				NamedImports: &NamedImports{
					Tokens: tk[8:10],
				},
				Comments: [6]Comments{{tk[0]}, {tk[4]}, {tk[6]}, nil, nil, {tk[11]}},
				Tokens:   tk[:12],
			}
		}},
		{"/* A */ * // B\nas // C\nb // D\n", func(t *test, tk Tokens) { // 19
			t.Output = ImportClause{
				NameSpaceImport: &tk[10],
				Comments:        [6]Comments{{tk[0]}, nil, nil, {tk[4]}, {tk[8]}, {tk[12]}},
				Tokens:          tk[:13],
			}
		}},
		{"// A\na // B\n, // C\n* // D\nas // E\n b // F\n", func(t *test, tk Tokens) { // 20
			t.Output = ImportClause{
				ImportedDefaultBinding: &tk[2],
				NameSpaceImport:        &tk[19],
				Comments:               [6]Comments{{tk[0]}, {tk[4]}, {tk[8]}, {tk[12]}, {tk[16]}, {tk[21]}},
				Tokens:                 tk[:22],
			}
		}},
	}, func(t *test) (Type, error) {
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
		{"from /* A */ \"\"", func(t *test, tk Tokens) { // 4
			t.Output = FromClause{
				ModuleSpecifier: &tk[4],
				Comments:        Comments{tk[2]},
				Tokens:          tk[:5],
			}
		}},
	}, func(t *test) (Type, error) {
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
						IdentifierName:  &tk[2],
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
						IdentifierName:  &tk[2],
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
						IdentifierName:  &tk[2],
						ImportedBinding: &tk[2],
						Tokens:          tk[2:3],
					},
					{
						IdentifierName:  &tk[6],
						ImportedBinding: &tk[6],
						Tokens:          tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"{\na\n,\na\n}", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[6],
			}
		}},
		{"{\na\n,\nb as a\n}", func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[6],
			}
		}},
		{"{\na as b\n,\nb\n}", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[10],
			}
		}},
		{"{\na as c\n,\nb as c\n}", func(t *test, tk Tokens) { // 11
			t.Err = Error{
				Err:     ErrInvalidNamedImport,
				Parsing: "NamedImports",
				Token:   tk[10],
			}
		}},
		{"{ // A\na\n // B\n}", func(t *test, tk Tokens) { // 12
			t.Output = NamedImports{
				ImportList: []ImportSpecifier{
					{
						IdentifierName:  &tk[4],
						ImportedBinding: &tk[4],
						Tokens:          tk[4:5],
					},
				},
				Comments: [2]Comments{{tk[2]}, {tk[7]}},
				Tokens:   tk[:10],
			}
		}},
		{"{ // A\n\n/* B */ a /* C */\n\n// D\n, // E\nb // F\n\n/* G */}", func(t *test, tk Tokens) { // 13
			t.Output = NamedImports{
				ImportList: []ImportSpecifier{
					{
						IdentifierName:  &tk[6],
						ImportedBinding: &tk[6],
						Comments:        [4]Comments{{tk[4]}, nil, nil, {tk[8], tk[10]}},
						Tokens:          tk[4:11],
					},
					{
						IdentifierName:  &tk[16],
						ImportedBinding: &tk[16],
						Comments:        [4]Comments{{tk[14]}, nil, nil, {tk[18]}},
						Tokens:          tk[14:19],
					},
				},
				Comments: [2]Comments{{tk[2]}, {tk[20]}},
				Tokens:   tk[:22],
			}
		}},
	}, func(t *test) (Type, error) {
		var ni NamedImports

		err := ni.parse(&t.Tokens)

		return ni, err
	})
}

func TestImportSpecifier(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidImportSpecifier,
				Parsing: "ImportSpecifier",
				Token:   tk[0],
			}
		}},
		{`,`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrInvalidImportSpecifier,
				Parsing: "ImportSpecifier",
				Token:   tk[0],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 3
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[0],
				ImportedBinding: &tk[0],
				Tokens:          tk[:1],
			}
		}},
		{"for", func(t *test, tk Tokens) { // 4
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[0],
				ImportedBinding: &tk[0],
				Tokens:          tk[:1],
			}
		}},
		{"for\nas", func(t *test, tk Tokens) { // 5
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[0],
				ImportedBinding: &tk[0],
				Tokens:          tk[:1],
			}
		}},
		{"a\nas", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ImportSpecifier",
				Token:   tk[3],
			}
		}},
		{"a\nas\nb", func(t *test, tk Tokens) { // 7
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[0],
				ImportedBinding: &tk[4],
				Tokens:          tk[:5],
			}
		}},
		{"/* A */ a // B", func(t *test, tk Tokens) { // 8
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[2],
				ImportedBinding: &tk[2],
				Comments:        [4]Comments{{tk[0]}, nil, nil, {tk[4]}},
				Tokens:          tk[:5],
			}
		}},
		{"// A\na /* B */ as // C\nb // D\n\n// E\n,", func(t *test, tk Tokens) { // 9
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[2],
				ImportedBinding: &tk[10],
				Comments:        [4]Comments{{tk[0]}, {tk[4]}, {tk[8]}, {tk[12], tk[14]}},
				Tokens:          tk[:15],
			}
		}},
		{"// A\na /* B */ as // C\nb // D\n\n// E", func(t *test, tk Tokens) { // 10
			t.Output = ImportSpecifier{
				IdentifierName:  &tk[2],
				ImportedBinding: &tk[10],
				Comments:        [4]Comments{{tk[0]}, {tk[4]}, {tk[8]}, {tk[12]}},
				Tokens:          tk[:13],
			}
		}},
	}, func(t *test) (Type, error) {
		var is ImportSpecifier

		err := is.parse(&t.Tokens)

		return is, err
	})
}

func TestWithClause(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingOpeningBrace,
				Parsing: "WithClause",
				Token:   tk[0],
			}
		}},
		{`{}`, func(t *test, tk Tokens) { // 2
			t.Output = WithClause{
				Tokens: tk[:2],
			}
		}},
		{`{a}`, func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingColon,
					Parsing: "WithEntry",
					Token:   tk[2],
				},
				Parsing: "WithClause",
				Token:   tk[1],
			}
		}},
		{`{a:"b"}`, func(t *test, tk Tokens) { // 4
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[1],
						Value:        &tk[3],
						Tokens:       tk[1:4],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{`{a:"b"c}`, func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingComma,
				Parsing: "WithClause",
				Token:   tk[4],
			}
		}},
		{`{a:"b",}`, func(t *test, tk Tokens) { // 6
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[1],
						Value:        &tk[3],
						Tokens:       tk[1:4],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{`{a:"b" ,}`, func(t *test, tk Tokens) { // 7
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[1],
						Value:        &tk[3],
						Tokens:       tk[1:4],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{`{a:"b", "c": "d"}`, func(t *test, tk Tokens) { // 8
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[1],
						Value:        &tk[3],
						Tokens:       tk[1:4],
					},
					{
						AttributeKey: &tk[6],
						Value:        &tk[9],
						Tokens:       tk[6:10],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{"/* A */{// B\na:\"b\"\n// C\n}", func(t *test, tk Tokens) { // 9
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[4],
						Value:        &tk[6],
						Tokens:       tk[4:7],
					},
				},
				Comments: [4]Comments{{tk[0]}, nil, {tk[2]}, {tk[8]}},
				Tokens:   tk[:11],
			}
		}},
		{"{/* A */\n\n// B\na:\"b\"\n// C\n, \"c\": \"d\"\n// E\n}", func(t *test, tk Tokens) { // 10
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[5],
						Value:        &tk[7],
						Comments:     [4]Comments{{tk[3]}, nil, nil, {tk[9]}},
						Tokens:       tk[3:10],
					},
					{
						AttributeKey: &tk[13],
						Value:        &tk[16],
						Tokens:       tk[13:17],
					},
				},
				Comments: [4]Comments{nil, nil, {tk[1]}, {tk[18]}},
				Tokens:   tk[:21],
			}
		}},
		{"/* A */ with // B\n{// C\na:\"b\"\n// D\n}", func(t *test, tk Tokens) { // 11
			t.Output = WithClause{
				WithEntries: []WithEntry{
					{
						AttributeKey: &tk[9],
						Value:        &tk[11],
						Tokens:       tk[9:12],
					},
				},
				Comments: [4]Comments{{tk[0]}, {tk[4]}, {tk[7]}, {tk[13]}},
				Tokens:   tk[:16],
			}
		}},
	}, func(t *test) (Type, error) {
		var wc WithClause

		err := wc.parse(&t.Tokens)

		return wc, err
	})
}

func TestWithEntry(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrMissingAttributeKey,
				Parsing: "WithEntry",
				Token:   tk[0],
			}
		}},
		{`a`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "WithEntry",
				Token:   tk[1],
			}
		}},
		{`a `, func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "WithEntry",
				Token:   tk[2],
			}
		}},
		{`a:`, func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err:     ErrMissingString,
				Parsing: "WithEntry",
				Token:   tk[2],
			}
		}},
		{`a: `, func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrMissingString,
				Parsing: "WithEntry",
				Token:   tk[3],
			}
		}},
		{`a:b`, func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrMissingString,
				Parsing: "WithEntry",
				Token:   tk[2],
			}
		}},
		{`a:"b"`, func(t *test, tk Tokens) { // 7
			t.Output = WithEntry{
				AttributeKey: &tk[0],
				Value:        &tk[2],
				Tokens:       tk[:3],
			}
		}},
		{`"a":"b"`, func(t *test, tk Tokens) { // 8
			t.Output = WithEntry{
				AttributeKey: &tk[0],
				Value:        &tk[2],
				Tokens:       tk[:3],
			}
		}},
		{"// A\na /* B */:/* C */\"b\" // D\n\n// E\n", func(t *test, tk Tokens) { // 9
			t.Output = WithEntry{
				AttributeKey: &tk[2],
				Value:        &tk[7],
				Comments:     [4]Comments{{tk[0]}, {tk[4]}, {tk[6]}, {tk[9]}},
				Tokens:       tk[:10],
			}
		}},
		{"// A\na /* B */:/* C */\"b\" // D\n\n// E\n,", func(t *test, tk Tokens) { // 10
			t.Output = WithEntry{
				AttributeKey: &tk[2],
				Value:        &tk[7],
				Comments:     [4]Comments{{tk[0]}, {tk[4]}, {tk[6]}, {tk[9], tk[11]}},
				Tokens:       tk[:12],
			}
		}},
	}, func(t *test) (Type, error) {
		var we WithEntry

		err := we.parse(&t.Tokens)

		return we, err
	})
}

func TestExportDeclaration(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidExportDeclaration,
				Parsing: "ExportDeclaration",
				Token:   tk[0],
			}
		}},
		{"export\ndefault\nfunction", func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingOpeningParenthesis,
						Parsing: "FormalParameters",
						Token:   tk[5],
					},
					Parsing: "FunctionDeclaration",
					Token:   tk[5],
				},
				Parsing: "ExportDeclaration",
				Token:   tk[4],
			}
		}},
		{"export\ndefault\nfunction(){}", func(t *test, tk Tokens) { // 3
			t.Output = ExportDeclaration{
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
			}
		}},
		{"export\ndefault\nclass", func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingOpeningBrace,
					Parsing: "ClassDeclaration",
					Token:   tk[5],
				},
				Parsing: "ExportDeclaration",
				Token:   tk[4],
			}
		}},
		{"export\ndefault\nclass{}", func(t *test, tk Tokens) { // 5
			t.Output = ExportDeclaration{
				DefaultClass: &ClassDeclaration{
					Tokens: tk[4:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"export\ndefault\n,", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "ExportDeclaration",
				Token:   tk[4],
			}
		}},
		{"export\ndefault\n1", func(t *test, tk Tokens) { // 7
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = ExportDeclaration{
				DefaultAssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"export\ndefault\n1 2", func(t *test, tk Tokens) { // 8
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "ExportDeclaration",
				Token:   tk[5],
			}
		}},
		{"export\ndefault\n1; 2", func(t *test, tk Tokens) { // 9
			lit1 := makeConditionLiteral(tk, 4)
			t.Output = ExportDeclaration{
				DefaultAssignmentExpression: &AssignmentExpression{
					ConditionalExpression: &lit1,
					Tokens:                tk[4:5],
				},
				Tokens: tk[:6],
			}
		}},
		{"export\n*", func(t *test, tk Tokens) { // 10
			t.Err = Error{
				Err: Error{
					Err:     ErrMissingFrom,
					Parsing: "FromClause",
					Token:   tk[3],
				},
				Parsing: "ExportDeclaration",
				Token:   tk[3],
			}
		}},
		{"export\n*\nfrom\n''", func(t *test, tk Tokens) { // 11
			t.Output = ExportDeclaration{
				FromClause: &FromClause{
					ModuleSpecifier: &tk[6],
					Tokens:          tk[4:7],
				},
				Tokens: tk[:7],
			}
		}},
		{"export\n*\nfrom\n'' b", func(t *test, tk Tokens) { // 12
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "ExportDeclaration",
				Token:   tk[7],
			}
		}},
		{"export\n*\nfrom\n'';b", func(t *test, tk Tokens) { // 13
			t.Output = ExportDeclaration{
				FromClause: &FromClause{
					ModuleSpecifier: &tk[6],
					Tokens:          tk[4:7],
				},
				Tokens: tk[:8],
			}
		}},
		{"export\nvar", func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "LexicalBinding",
						Token:   tk[3],
					},
					Parsing: "VariableStatement",
					Token:   tk[3],
				},
				Parsing: "ExportDeclaration",
				Token:   tk[2],
			}
		}},
		{"export\nvar\na", func(t *test, tk Tokens) { // 15
			t.Output = ExportDeclaration{
				VariableStatement: &VariableStatement{
					VariableDeclarationList: []VariableDeclaration{
						{
							BindingIdentifier: &tk[4],
							Tokens:            tk[4:5],
						},
					},
					Tokens: tk[2:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"export\n{,}", func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrNoIdentifier,
						Parsing: "ExportSpecifier",
						Token:   tk[3],
					},
					Parsing: "ExportClause",
					Token:   tk[3],
				},
				Parsing: "ExportDeclaration",
				Token:   tk[2],
			}
		}},
		{"export\n{}", func(t *test, tk Tokens) { // 17
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				Tokens: tk[:4],
			}
		}},
		{"export\n{};", func(t *test, tk Tokens) { // 18
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				Tokens: tk[:5],
			}
		}},
		{"export\n{} b", func(t *test, tk Tokens) { // 19
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "ExportDeclaration",
				Token:   tk[4],
			}
		}},
		{"export\n{}\nfrom\n''", func(t *test, tk Tokens) { // 20
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				FromClause: &FromClause{
					ModuleSpecifier: &tk[7],
					Tokens:          tk[5:8],
				},
				Tokens: tk[:8],
			}
		}},
		{"export\n{}\nfrom\n'' a", func(t *test, tk Tokens) { // 21
			t.Err = Error{
				Err:     ErrMissingSemiColon,
				Parsing: "ExportDeclaration",
				Token:   tk[8],
			}
		}},
		{"export\n{}\nfrom\n'';a", func(t *test, tk Tokens) { // 22
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				FromClause: &FromClause{
					ModuleSpecifier: &tk[7],
					Tokens:          tk[5:8],
				},
				Tokens: tk[:9],
			}
		}},
		{"export\n*\nas\na\nfrom\n'';", func(t *test, tk Tokens) { // 23
			t.Output = ExportDeclaration{
				ExportFromClause: &tk[6],
				FromClause: &FromClause{
					ModuleSpecifier: &tk[10],
					Tokens:          tk[8:11],
				},
				Tokens: tk[:12],
			}
		}},
		{"export * as ;", func(t *test, tk Tokens) { // 24
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ExportDeclaration",
				Token:   tk[6],
			}
		}},

		{"// A\nexport /* B */ default /* C */ function(){}", func(t *test, tk Tokens) { // 25
			t.Output = ExportDeclaration{
				DefaultFunction: &FunctionDeclaration{
					FormalParameters: FormalParameters{
						Tokens: tk[11:13],
					},
					FunctionBody: Block{
						Tokens: tk[13:15],
					},
					Tokens: tk[10:15],
				},
				Comments: [7]Comments{{tk[0]}, {tk[4]}, {tk[8]}},
				Tokens:   tk[:15],
			}
		}},
		{"// A\nexport /* B */ default /* C */ 1 /* D */;", func(t *test, tk Tokens) { // 26
			t.Output = ExportDeclaration{
				DefaultAssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							Literal: &tk[10],
							Tokens:  tk[10:11],
						},
						Comments: [5]Comments{nil, nil, nil, nil, {tk[12]}},
						Tokens:   tk[10:13],
					}),
					Tokens: tk[10:13],
				},
				Comments: [7]Comments{{tk[0]}, {tk[4]}, {tk[8]}},
				Tokens:   tk[:14],
			}
		}},
		{"export /* A */ * /* B */ from /* C */ '' /* D */;", func(t *test, tk Tokens) { // 27
			t.Output = ExportDeclaration{
				FromClause: &FromClause{
					ModuleSpecifier: &tk[12],
					Comments:        Comments{tk[10]},
					Tokens:          tk[8:13],
				},
				Comments: [7]Comments{nil, {tk[2]}, {tk[6]}, nil, nil, {tk[14]}},
				Tokens:   tk[:16],
			}
		}},
		{"/* A */export/* B */*/* C */as/* D */a/* E */from/* F */''/* G */;", func(t *test, tk Tokens) { // 28
			t.Output = ExportDeclaration{
				ExportFromClause: &tk[7],
				FromClause: &FromClause{
					ModuleSpecifier: &tk[11],
					Comments:        Comments{tk[10]},
					Tokens:          tk[9:12],
				},
				Comments: [7]Comments{{tk[0]}, {tk[2]}, {tk[4]}, {tk[6]}, {tk[8]}, {tk[12]}},
				Tokens:   tk[:14],
			}
		}},
		{"export/* A */{}/* B */; // C\n\n// D", func(t *test, tk Tokens) { // 29
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				Comments: [7]Comments{nil, {tk[1]}, nil, nil, nil, {tk[4]}, {tk[7]}},
				Tokens:   tk[:8],
			}
		}},
		{"export/* A */{}/* B */from/* C */''/* D */", func(t *test, tk Tokens) { // 30
			t.Output = ExportDeclaration{
				ExportClause: &ExportClause{
					Tokens: tk[2:4],
				},
				FromClause: &FromClause{
					ModuleSpecifier: &tk[7],
					Comments:        Comments{tk[6]},
					Tokens:          tk[5:8],
				},
				Comments: [7]Comments{nil, {tk[1]}, nil, nil, {tk[4]}, {tk[8]}},
				Tokens:   tk[:9],
			}
		}},
	}, func(t *test) (Type, error) {
		var ed ExportDeclaration

		err := ed.parse(&t.Tokens)

		return ed, err
	})
}

func TestExportClause(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrInvalidExportClause,
				Parsing: "ExportClause",
				Token:   tk[0],
			}
		}},
		{"{\n}", func(t *test, tk Tokens) { // 2
			t.Output = ExportClause{
				Tokens: tk[:3],
			}
		}},
		{"{\n,}", func(t *test, tk Tokens) { // 3
			t.Err = Error{
				Err: Error{
					Err:     ErrNoIdentifier,
					Parsing: "ExportSpecifier",
					Token:   tk[2],
				},
				Parsing: "ExportClause",
				Token:   tk[2],
			}
		}},
		{"{\na\n}", func(t *test, tk Tokens) { // 4
			t.Output = ExportClause{
				ExportList: []ExportSpecifier{
					{
						IdentifierName:  &tk[2],
						EIdentifierName: &tk[2],
						Tokens:          tk[2:3],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{"{\na\nb}", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrInvalidExportClause,
				Parsing: "ExportClause",
				Token:   tk[4],
			}
		}},
		{"{\na\n,\n}", func(t *test, tk Tokens) { // 6
			t.Output = ExportClause{
				ExportList: []ExportSpecifier{
					{
						IdentifierName:  &tk[2],
						EIdentifierName: &tk[2],
						Tokens:          tk[2:3],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{"{\na\n,\nb\n}", func(t *test, tk Tokens) { // 7
			t.Output = ExportClause{
				ExportList: []ExportSpecifier{
					{
						IdentifierName:  &tk[2],
						EIdentifierName: &tk[2],
						Tokens:          tk[2:3],
					},
					{
						IdentifierName:  &tk[6],
						EIdentifierName: &tk[6],
						Tokens:          tk[6:7],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{"{ // A\n// B\n\n// C\na // D\n// E\n\n// F\n}", func(t *test, tk Tokens) { // 8
			t.Output = ExportClause{
				ExportList: []ExportSpecifier{
					{
						IdentifierName:  &tk[8],
						EIdentifierName: &tk[8],
						Comments:        [4]Comments{{tk[6]}, nil, nil, {tk[10], tk[12]}},
						Tokens:          tk[6:13],
					},
				},
				Comments: [2]Comments{{tk[2], tk[4]}, {tk[14]}},
				Tokens:   tk[:17],
			}
		}},
	}, func(t *test) (Type, error) {
		var ec ExportClause

		err := ec.parse(&t.Tokens)

		return ec, err
	})
}

func TestExportSpecifier(t *testing.T) {
	doTests(t, []sourceFn{
		{``, func(t *test, tk Tokens) { // 1
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ExportSpecifier",
				Token:   tk[0],
			}
		}},
		{`+`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ExportSpecifier",
				Token:   tk[0],
			}
		}},
		{"a", func(t *test, tk Tokens) { // 3
			t.Output = ExportSpecifier{
				IdentifierName:  &tk[0],
				EIdentifierName: &tk[0],
				Tokens:          tk[:1],
			}
		}},
		{"for", func(t *test, tk Tokens) { // 4
			t.Output = ExportSpecifier{
				IdentifierName:  &tk[0],
				EIdentifierName: &tk[0],
				Tokens:          tk[:1],
			}
		}},
		{"a\nas", func(t *test, tk Tokens) { // 5
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ExportSpecifier",
				Token:   tk[3],
			}
		}},
		{"a\nas\n,", func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err:     ErrNoIdentifier,
				Parsing: "ExportSpecifier",
				Token:   tk[4],
			}
		}},
		{"a\nas\nb", func(t *test, tk Tokens) { // 7
			t.Output = ExportSpecifier{
				IdentifierName:  &tk[0],
				EIdentifierName: &tk[4],
				Tokens:          tk[:5],
			}
		}},
		{"// A\na // B\n\n// C\n", func(t *test, tk Tokens) { // 8
			t.Output = ExportSpecifier{
				IdentifierName:  &tk[2],
				EIdentifierName: &tk[2],
				Comments:        [4]Comments{{tk[0]}, nil, nil, {tk[4]}},
				Tokens:          tk[:5],
			}
		}},
		{"// A\na /* B */ as /* C */ b // D\n\n// E\n,", func(t *test, tk Tokens) { // 9
			t.Output = ExportSpecifier{
				IdentifierName:  &tk[2],
				EIdentifierName: &tk[10],
				Comments:        [4]Comments{{tk[0]}, {tk[4]}, {tk[8]}, {tk[12], tk[14]}},
				Tokens:          tk[:15],
			}
		}},
	}, func(t *test) (Type, error) {
		var es ExportSpecifier

		err := es.parse(&t.Tokens)

		return es, err
	})
}
