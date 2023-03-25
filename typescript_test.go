package javascript

import "testing"

func TestTypescriptModule(t *testing.T) {
	doTests(t, []sourceFn{
		{`import def from './a';import type typeDef from './b';import type {typ1, typ2} from './c';import {a} from './d';`, func(t *test, tk Tokens) { // 1
			t.Typescript = true
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
								Tokens: tk[8:18],
							},
							Tokens: tk[8:18],
						},
						Tokens: tk[8:18],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[18:33],
							},
							Tokens: tk[18:33],
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
		{`import def from './a';import type typeDef from './b';import type {typ1, typ2} from './c';import {a} from './d';`, func(t *test, tk Tokens) { // 2
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingFrom,
						Parsing: "FromClause",
						Token:   tk[12],
					},
					Parsing: "ImportDeclaration",
					Token:   tk[12],
				},
				Parsing: "ModuleItem",
				Token:   tk[8],
			}
		}},
		{`class A {
private;
public
protected;
a
private A
protected B;
public C
private static D;
protected static E
public static F
}`, func(t *test, tk Tokens) { // 3
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[2],
									ClassBody: []ClassElement{
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[6],
														Tokens:              tk[6:7],
													},
													Tokens: tk[6:7],
												},
												Tokens: tk[6:7],
											},
											Tokens: tk[6:8],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[9],
														Tokens:              tk[9:10],
													},
													Tokens: tk[9:10],
												},
												Tokens: tk[9:10],
											},
											Tokens: tk[9:10],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[11],
														Tokens:              tk[11:12],
													},
													Tokens: tk[11:12],
												},
												Tokens: tk[11:12],
											},
											Tokens: tk[11:13],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[14],
														Tokens:              tk[14:15],
													},
													Tokens: tk[14:15],
												},
												Tokens: tk[14:15],
											},
											Tokens: tk[14:15],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[18],
														Tokens:              tk[18:19],
													},
													Tokens: tk[18:19],
												},
												Tokens: tk[18:19],
											},
											Tokens: tk[16:19],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[22],
														Tokens:              tk[22:23],
													},
													Tokens: tk[22:23],
												},
												Tokens: tk[22:23],
											},
											Tokens: tk[20:24],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[27],
														Tokens:              tk[27:28],
													},
													Tokens: tk[27:28],
												},
												Tokens: tk[27:28],
											},
											Tokens: tk[25:28],
										},
										{
											Static: true,
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[33],
														Tokens:              tk[33:34],
													},
													Tokens: tk[33:34],
												},
												Tokens: tk[33:34],
											},
											Tokens: tk[29:35],
										},
										{
											Static: true,
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[40],
														Tokens:              tk[40:41],
													},
													Tokens: tk[40:41],
												},
												Tokens: tk[40:41],
											},
											Tokens: tk[36:41],
										},
										{
											Static: true,
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[46],
														Tokens:              tk[46:47],
													},
													Tokens: tk[46:47],
												},
												Tokens: tk[46:47],
											},
											Tokens: tk[42:47],
										},
									},
									Tokens: tk[:49],
								},
								Tokens: tk[:49],
							},
							Tokens: tk[:49],
						},
						Tokens: tk[:49],
					},
				},
				Tokens: tk[:49],
			}
		}},
		{`class A {
private;
public
protected;
a
private A
protected B;
public C
private static D;
protected static E
public static F
}`, func(t *test, tk Tokens) { // 4
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingSemiColon,
								Parsing: "ClassElement",
								Token:   tk[18],
							},
							Parsing: "ClassDeclaration",
							Token:   tk[16],
						},
						Parsing: "Declaration",
						Token:   tk[0],
					},
					Parsing: "StatementListItem",
					Token:   tk[0],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{`class A {
readonly;
a
readonly A
private readonly B
protected readonly C;
public readonly D
static readonly E;
}`, func(t *test, tk Tokens) { // 5
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[2],
									ClassBody: []ClassElement{
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[6],
														Tokens:              tk[6:7],
													},
													Tokens: tk[6:7],
												},
												Tokens: tk[6:7],
											},
											Tokens: tk[6:8],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[9],
														Tokens:              tk[9:10],
													},
													Tokens: tk[9:10],
												},
												Tokens: tk[9:10],
											},
											Tokens: tk[9:10],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[13],
														Tokens:              tk[13:14],
													},
													Tokens: tk[13:14],
												},
												Tokens: tk[13:14],
											},
											Tokens: tk[11:14],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[19],
														Tokens:              tk[19:20],
													},
													Tokens: tk[19:20],
												},
												Tokens: tk[19:20],
											},
											Tokens: tk[15:20],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[25],
														Tokens:              tk[25:26],
													},
													Tokens: tk[25:26],
												},
												Tokens: tk[25:26],
											},
											Tokens: tk[21:27],
										},
										{
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[32],
														Tokens:              tk[32:33],
													},
													Tokens: tk[32:33],
												},
												Tokens: tk[32:33],
											},
											Tokens: tk[28:33],
										},
										{
											Static: true,
											FieldDefinition: &FieldDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[38],
														Tokens:              tk[38:39],
													},
													Tokens: tk[38:39],
												},
												Tokens: tk[38:39],
											},
											Tokens: tk[34:40],
										},
									},
									Tokens: tk[:42],
								},
								Tokens: tk[:42],
							},
							Tokens: tk[:42],
						},
						Tokens: tk[:42],
					},
				},
				Tokens: tk[:42],
			}
		}},
		{`class A {
readonly;
a
readonly A
private readonly B
protected readonly C;
public readonly D
static readonly E;
}`, func(t *test, tk Tokens) { // 6
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingSemiColon,
								Parsing: "ClassElement",
								Token:   tk[13],
							},
							Parsing: "ClassDeclaration",
							Token:   tk[11],
						},
						Parsing: "Declaration",
						Token:   tk[0],
					},
					Parsing: "StatementListItem",
					Token:   tk[0],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{`import {type} from 'a';`, func(t *test, tk Tokens) { // 7
			t.Typescript = true
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
		{`import {type A} from 'a';`, func(t *test, tk Tokens) { // 8
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									Tokens: tk[2:7],
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
		{`import {type A} from 'a';`, func(t *test, tk Tokens) { // 9
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidNamedImport,
							Parsing: "NamedImports",
							Token:   tk[5],
						},
						Parsing: "ImportClause",
						Token:   tk[2],
					},
					Parsing: "ImportDeclaration",
					Token:   tk[2],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{`import {type as} from 'a';`, func(t *test, tk Tokens) { // 10
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									Tokens: tk[2:7],
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
		{`import {type as as} from 'a';`, func(t *test, tk Tokens) { // 11
			t.Typescript = true
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
		{`import {type as as as} from 'a';`, func(t *test, tk Tokens) { // 12
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ImportDeclaration: &ImportDeclaration{
							ImportClause: &ImportClause{
								NamedImports: &NamedImports{
									Tokens: tk[2:11],
								},
								Tokens: tk[2:11],
							},
							FromClause: FromClause{
								ModuleSpecifier: &tk[14],
								Tokens:          tk[12:15],
							},
							Tokens: tk[:16],
						},
						Tokens: tk[:16],
					},
				},
				Tokens: tk[:16],
			}
		}},
		{`{}
type B = number;
{}`, func(t *test, tk Tokens) { // 13
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								BlockStatement: &Block{
									Tokens: tk[:2],
								},
								Tokens: tk[:2],
							},
							Tokens: tk[:2],
						},
						Tokens: tk[:2],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[3:11],
							},
							Tokens: tk[3:11],
						},
						Tokens: tk[3:11],
					},
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								BlockStatement: &Block{
									Tokens: tk[12:14],
								},
								Tokens: tk[12:14],
							},
							Tokens: tk[12:14],
						},
						Tokens: tk[12:14],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`{}
type B = number;
{}`, func(t *test, tk Tokens) { // 14
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingSemiColon,
						Parsing: "Statement",
						Token:   tk[4],
					},
					Parsing: "StatementListItem",
					Token:   tk[3],
				},
				Parsing: "ModuleItem",
				Token:   tk[3],
			}
		}},
	}, func(t *test) (Type, error) {
		if t.Typescript {
			t.Tokens[:cap(t.Tokens)][cap(t.Tokens)-1].Data = marker
		}
		var m Module
		err := m.parse(&t.Tokens)
		return m, err
	})
}
