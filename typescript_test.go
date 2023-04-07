package javascript

import (
	"testing"

	"vimagination.zapto.org/parser"
)

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
		{`class A<B, C> extends D<E> implements F, G<H> {}`, func(t *test, tk Tokens) { // 15
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[2],
									ClassHeritage: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[12],
													Tokens:              tk[12:13],
												},
												Tokens: tk[12:13],
											},
											Tokens: tk[12:13],
										},
										Tokens: tk[12:13],
									},
									Tokens: tk[:29],
								},
								Tokens: tk[:29],
							},
							Tokens: tk[:29],
						},
						Tokens: tk[:29],
					},
				},
				Tokens: tk[:29],
			}
		}},
		{`class A<B, C> extends D<E> implements F, G<H> {}`, func(t *test, tk Tokens) { // 16
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningBrace,
							Parsing: "ClassDeclaration",
							Token:   tk[3],
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
		{`class A extends D<E> implements F, G<H> {}`, func(t *test, tk Tokens) { // 17
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningBrace,
							Parsing: "ClassDeclaration",
							Token:   tk[7],
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
		{`class A extends D implements F, G<H> {}`, func(t *test, tk Tokens) { // 18
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningBrace,
							Parsing: "ClassDeclaration",
							Token:   tk[8],
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
		{`class A<B, C> implements D, E<F> extends G<H> {}`, func(t *test, tk Tokens) { // 19
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[2],
									ClassHeritage: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[22],
													Tokens:              tk[22:23],
												},
												Tokens: tk[22:23],
											},
											Tokens: tk[22:23],
										},
										Tokens: tk[22:23],
									},
									Tokens: tk[:29],
								},
								Tokens: tk[:29],
							},
							Tokens: tk[:29],
						},
						Tokens: tk[:29],
					},
				},
				Tokens: tk[:29],
			}
		}},
		{`class A implements D, E<F> extends G<H> {}`, func(t *test, tk Tokens) { // 20
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrMissingOpeningBrace,
							Parsing: "ClassDeclaration",
							Token:   tk[4],
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
b(): C {}
get d(): E {}
set f(g): H {}
i <J> () {}
}`, func(t *test, tk Tokens) { // 21
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
											MethodDefinition: &MethodDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[6],
														Tokens:              tk[6:7],
													},
													Tokens: tk[6:7],
												},
												Params: FormalParameters{
													Tokens: tk[7:9],
												},
												FunctionBody: Block{
													Tokens: tk[13:15],
												},
												Tokens: tk[6:15],
											},
											Tokens: tk[6:15],
										},
										{
											MethodDefinition: &MethodDefinition{
												Type: MethodGetter,
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[18],
														Tokens:              tk[18:19],
													},
													Tokens: tk[18:19],
												},
												Params: FormalParameters{
													Tokens: tk[19:21],
												},
												FunctionBody: Block{
													Tokens: tk[25:27],
												},
												Tokens: tk[16:27],
											},
											Tokens: tk[16:27],
										},
										{
											MethodDefinition: &MethodDefinition{
												Type: MethodSetter,
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[30],
														Tokens:              tk[30:31],
													},
													Tokens: tk[30:31],
												},
												Params: FormalParameters{
													FormalParameterList: []BindingElement{
														{
															SingleNameBinding: &tk[32],
															Tokens:            tk[32:33],
														},
													},
													Tokens: tk[31:34],
												},
												FunctionBody: Block{
													Tokens: tk[38:40],
												},
												Tokens: tk[28:40],
											},
											Tokens: tk[28:40],
										},
										{
											MethodDefinition: &MethodDefinition{
												ClassElementName: ClassElementName{
													PropertyName: &PropertyName{
														LiteralPropertyName: &tk[41],
														Tokens:              tk[41:42],
													},
													Tokens: tk[41:42],
												},
												Params: FormalParameters{
													Tokens: tk[47:49],
												},
												FunctionBody: Block{
													Tokens: tk[50:52],
												},
												Tokens: tk[41:52],
											},
											Tokens: tk[41:52],
										},
									},
									Tokens: tk[:54],
								},
								Tokens: tk[:54],
							},
							Tokens: tk[:54],
						},
						Tokens: tk[:54],
					},
				},
				Tokens: tk[:54],
			}
		}},
		{`class A {
b(): C {}
get d(): E {}
set f(g): H {}
i <J> () {}
}`, func(t *test, tk Tokens) { // 22
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err:     ErrMissingOpeningBrace,
										Parsing: "Block",
										Token:   tk[9],
									},
									Parsing: "MethodDefinition",
									Token:   tk[9],
								},
								Parsing: "ClassElement",
								Token:   tk[6],
							},
							Parsing: "ClassDeclaration",
							Token:   tk[6],
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
		{`function A<B>(c: D, [e, f]?: [number, string], {g}: {g: boolean}, ...i: J): K {}`, func(t *test, tk Tokens) { // 23
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								FunctionDeclaration: &FunctionDeclaration{
									BindingIdentifier: &tk[2],
									FormalParameters: FormalParameters{
										FormalParameterList: []BindingElement{
											{
												SingleNameBinding: &tk[7],
												Tokens:            tk[7:11],
											},
											{
												ArrayBindingPattern: &ArrayBindingPattern{
													BindingElementList: []BindingElement{
														{
															SingleNameBinding: &tk[14],
															Tokens:            tk[14:15],
														},
														{
															SingleNameBinding: &tk[17],
															Tokens:            tk[17:18],
														},
													},
													Tokens: tk[13:19],
												},
												Tokens: tk[13:28],
											},
											{
												ObjectBindingPattern: &ObjectBindingPattern{
													BindingPropertyList: []BindingProperty{
														{
															PropertyName: PropertyName{
																LiteralPropertyName: &tk[31],
																Tokens:              tk[31:32],
															},
															BindingElement: BindingElement{
																SingleNameBinding: &tk[31],
																Tokens:            tk[31:32],
															},
															Tokens: tk[31:32],
														},
													},
													Tokens: tk[30:33],
												},
												Tokens: tk[30:41],
											},
										},
										BindingIdentifier: &tk[44],
										Tokens:            tk[6:49],
									},
									FunctionBody: Block{
										Tokens: tk[53:55],
									},
									Tokens: tk[:55],
								},
								Tokens: tk[:55],
							},
							Tokens: tk[:55],
						},
						Tokens: tk[:55],
					},
				},
				Tokens: tk[:55],
			}
		}},
		{`function A<B>(c: D, [e, f]?: [number, string], {g}: {g: boolean}, ...i: J): K {}`, func(t *test, tk Tokens) { // 24
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err:     ErrMissingOpeningParenthesis,
								Parsing: "FormalParameters",
								Token:   tk[3],
							},
							Parsing: "FunctionDeclaration",
							Token:   tk[3],
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
		{`1 as number`, func(t *test, tk Tokens) { // 25
			t.Typescript = true
			one := makeConditionLiteral(tk, 0)
			one.Tokens = tk[:5]
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: &one,
											Tokens:                tk[:5],
										},
									},
									Tokens: tk[:5],
								},
								Tokens: tk[:5],
							},
							Tokens: tk[:5],
						},
						Tokens: tk[:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{`1 as number`, func(t *test, tk Tokens) { // 26
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err:     ErrMissingSemiColon,
						Parsing: "Statement",
						Token:   tk[1],
					},
					Parsing: "StatementListItem",
					Token:   tk[0],
				},
				Parsing: "ModuleItem",
				Token:   tk[0],
			}
		}},
		{`function a(this){}`, func(t *test, tk Tokens) { // 27
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								FunctionDeclaration: &FunctionDeclaration{
									BindingIdentifier: &tk[2],
									FormalParameters: FormalParameters{
										Tokens: tk[3:6],
									},
									FunctionBody: Block{
										Tokens: tk[6:8],
									},
									Tokens: tk[:8],
								},
								Tokens: tk[:8],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{`function a(this){}`, func(t *test, tk Tokens) { // 28
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrNoIdentifier,
									Parsing: "BindingElement",
									Token:   tk[4],
								},
								Parsing: "FormalParameters",
								Token:   tk[4],
							},
							Parsing: "FunctionDeclaration",
							Token:   tk[3],
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
		{`function a(this: T){}`, func(t *test, tk Tokens) { // 29
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								FunctionDeclaration: &FunctionDeclaration{
									BindingIdentifier: &tk[2],
									FormalParameters: FormalParameters{
										Tokens: tk[3:9],
									},
									FunctionBody: Block{
										Tokens: tk[9:11],
									},
									Tokens: tk[:11],
								},
								Tokens: tk[:11],
							},
							Tokens: tk[:11],
						},
						Tokens: tk[:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`function a(this: T){}`, func(t *test, tk Tokens) { // 30
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrNoIdentifier,
									Parsing: "BindingElement",
									Token:   tk[4],
								},
								Parsing: "FormalParameters",
								Token:   tk[4],
							},
							Parsing: "FunctionDeclaration",
							Token:   tk[3],
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
		{`function a(this, b){}`, func(t *test, tk Tokens) { // 31
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								FunctionDeclaration: &FunctionDeclaration{
									BindingIdentifier: &tk[2],
									FormalParameters: FormalParameters{
										FormalParameterList: []BindingElement{
											{
												SingleNameBinding: &tk[7],
												Tokens:            tk[7:8],
											},
										},
										Tokens: tk[3:9],
									},
									FunctionBody: Block{
										Tokens: tk[9:11],
									},
									Tokens: tk[:11],
								},
								Tokens: tk[:11],
							},
							Tokens: tk[:11],
						},
						Tokens: tk[:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`function a(this, b){}`, func(t *test, tk Tokens) { // 32
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err:     ErrNoIdentifier,
									Parsing: "BindingElement",
									Token:   tk[4],
								},
								Parsing: "FormalParameters",
								Token:   tk[4],
							},
							Parsing: "FunctionDeclaration",
							Token:   tk[3],
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
		{`<A>() => {}`, func(t *test, tk Tokens) { // 33
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												FormalParameters: &FormalParameters{
													Tokens: tk[3:5],
												},
												FunctionBody: &Block{
													Tokens: tk[8:10],
												},
												Tokens: tk[:10],
											},
											Tokens: tk[:10],
										},
									},
									Tokens: tk[:10],
								},
								Tokens: tk[:10],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`(): A => {}`, func(t *test, tk Tokens) { // 35
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												FormalParameters: &FormalParameters{
													Tokens: tk[0:2],
												},
												FunctionBody: &Block{
													Tokens: tk[8:10],
												},
												Tokens: tk[:10],
											},
											Tokens: tk[:10],
										},
									},
									Tokens: tk[:10],
								},
								Tokens: tk[:10],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`(a: T) => {}`, func(t *test, tk Tokens) { // 35
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												FormalParameters: &FormalParameters{
													FormalParameterList: []BindingElement{
														{
															SingleNameBinding: &tk[1],
															Tokens:            tk[1:2],
														},
													},
													Tokens: tk[0:6],
												},
												FunctionBody: &Block{
													Tokens: tk[9:11],
												},
												Tokens: tk[:11],
											},
											Tokens: tk[:11],
										},
									},
									Tokens: tk[:11],
								},
								Tokens: tk[:11],
							},
							Tokens: tk[:11],
						},
						Tokens: tk[:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`<A>(a: A, b: C): D => {}`, func(t *test, tk Tokens) { // 36
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												FormalParameters: &FormalParameters{
													FormalParameterList: []BindingElement{
														{
															SingleNameBinding: &tk[4],
															Tokens:            tk[4:5],
														},
														{
															SingleNameBinding: &tk[10],
															Tokens:            tk[10:11],
														},
													},
													Tokens: tk[3:15],
												},
												FunctionBody: &Block{
													Tokens: tk[21:23],
												},
												Tokens: tk[:23],
											},
											Tokens: tk[:23],
										},
									},
									Tokens: tk[:23],
								},
								Tokens: tk[:23],
							},
							Tokens: tk[:23],
						},
						Tokens: tk[:23],
					},
				},
				Tokens: tk[:23],
			}
		}},
		{`async <A>() => {}`, func(t *test, tk Tokens) { // 37
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												Async: true,
												FormalParameters: &FormalParameters{
													Tokens: tk[5:7],
												},
												FunctionBody: &Block{
													Tokens: tk[10:12],
												},
												Tokens: tk[:12],
											},
											Tokens: tk[:12],
										},
									},
									Tokens: tk[:12],
								},
								Tokens: tk[:12],
							},
							Tokens: tk[:12],
						},
						Tokens: tk[:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{`async (): A => {}`, func(t *test, tk Tokens) { // 38
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												Async: true,
												FormalParameters: &FormalParameters{
													Tokens: tk[2:4],
												},
												FunctionBody: &Block{
													Tokens: tk[10:12],
												},
												Tokens: tk[:12],
											},
											Tokens: tk[:12],
										},
									},
									Tokens: tk[:12],
								},
								Tokens: tk[:12],
							},
							Tokens: tk[:12],
						},
						Tokens: tk[:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{`async (a: T) => {}`, func(t *test, tk Tokens) { // 39
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												Async: true,
												FormalParameters: &FormalParameters{
													FormalParameterList: []BindingElement{
														{
															SingleNameBinding: &tk[3],
															Tokens:            tk[3:7],
														},
													},
													Tokens: tk[2:8],
												},
												FunctionBody: &Block{
													Tokens: tk[11:13],
												},
												Tokens: tk[:13],
											},
											Tokens: tk[:13],
										},
									},
									Tokens: tk[:13],
								},
								Tokens: tk[:13],
							},
							Tokens: tk[:13],
						},
						Tokens: tk[:13],
					},
				},
				Tokens: tk[:13],
			}
		}},
		{`async <A>(a: A, b: C): D => {}`, func(t *test, tk Tokens) { // 40
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ArrowFunction: &ArrowFunction{
												Async: true,
												FormalParameters: &FormalParameters{
													FormalParameterList: []BindingElement{
														{
															SingleNameBinding: &tk[6],
															Tokens:            tk[6:10],
														},
														{
															SingleNameBinding: &tk[12],
															Tokens:            tk[12:16],
														},
													},
													Tokens: tk[5:17],
												},
												FunctionBody: &Block{
													Tokens: tk[23:25],
												},
												Tokens: tk[:25],
											},
											Tokens: tk[:25],
										},
									},
									Tokens: tk[:25],
								},
								Tokens: tk[:25],
							},
							Tokens: tk[:25],
						},
						Tokens: tk[:25],
					},
				},
				Tokens: tk[:25],
			}
		}},
		{"let a: B = c as D, [e] = f as const", func(t *test, tk Tokens) { // 41
			t.Typescript = true
			c := WrapConditional(&PrimaryExpression{
				IdentifierReference: &tk[9],
				Tokens:              tk[9:10],
			})
			c.Tokens = tk[9:14]
			f := WrapConditional(&PrimaryExpression{
				IdentifierReference: &tk[22],
				Tokens:              tk[22:23],
			})
			f.Tokens = tk[22:27]
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ConditionalExpression: c,
												Tokens:                tk[9:14],
											},
											Tokens: tk[2:14],
										},
										{
											ArrayBindingPattern: &ArrayBindingPattern{
												BindingElementList: []BindingElement{
													{
														SingleNameBinding: &tk[17],
														Tokens:            tk[17:18],
													},
												},
												Tokens: tk[16:19],
											},
											Initializer: &AssignmentExpression{
												ConditionalExpression: f,
												Tokens:                tk[22:27],
											},
											Tokens: tk[16:27],
										},
									},
									Tokens: tk[:27],
								},
								Tokens: tk[:27],
							},
							Tokens: tk[:27],
						},
						Tokens: tk[:27],
					},
				},
				Tokens: tk[:27],
			}
		}},
		{"let a: B = c as D, [e] = f as const", func(t *test, tk Tokens) { // 42
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err:     ErrInvalidLexicalDeclaration,
							Parsing: "LexicalDeclaration",
							Token:   tk[3],
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
		{`type A = {
	data: any;
	other: number;
}`, func(t *test, tk Tokens) { // 43
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:23],
							},
							Tokens: tk[:23],
						},
						Tokens: tk[:23],
					},
				},
				Tokens: tk[:23],
			}
		}},
		{`const a = B<C, D>()`, func(t *test, tk Tokens) { // 44
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ConditionalExpression: WrapConditional(&CallExpression{
													MemberExpression: &MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: &tk[6],
															Tokens:              tk[6:7],
														},
														Tokens: tk[6:7],
													},
													Arguments: &Arguments{
														Tokens: tk[13:15],
													},
													Tokens: tk[6:15],
												}),
												Tokens: tk[6:15],
											},
											Tokens: tk[2:15],
										},
									},
									Tokens: tk[:15],
								},
								Tokens: tk[:15],
							},
							Tokens: tk[:15],
						},
						Tokens: tk[:15],
					},
				},
				Tokens: tk[:15],
			}
		}},
		{`const a = new B<C, D>()`, func(t *test, tk Tokens) { // 45
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ConditionalExpression: WrapConditional(&NewExpression{
													MemberExpression: MemberExpression{
														MemberExpression: &MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																IdentifierReference: &tk[8],
																Tokens:              tk[8:9],
															},
															Tokens: tk[8:9],
														},
														Arguments: &Arguments{
															Tokens: tk[15:17],
														},
														Tokens: tk[6:17],
													},
													Tokens: tk[6:17],
												}),
												Tokens: tk[6:17],
											},
											Tokens: tk[2:17],
										},
									},
									Tokens: tk[:17],
								},
								Tokens: tk[:17],
							},
							Tokens: tk[:17],
						},
						Tokens: tk[:17],
					},
				},
				Tokens: tk[:17],
			}
		}},
		{`const a = new B<C<D, E>>()`, func(t *test, tk Tokens) { // 46
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ConditionalExpression: WrapConditional(&NewExpression{
													MemberExpression: MemberExpression{
														MemberExpression: &MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																IdentifierReference: &tk[8],
																Tokens:              tk[8:9],
															},
															Tokens: tk[8:9],
														},
														Arguments: &Arguments{
															Tokens: tk[18:20],
														},
														Tokens: tk[6:20],
													},
													Tokens: tk[6:20],
												}),
												Tokens: tk[6:20],
											},
											Tokens: tk[2:20],
										},
									},
									Tokens: tk[:20],
								},
								Tokens: tk[:20],
							},
							Tokens: tk[:20],
						},
						Tokens: tk[:20],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{`const a = (b): b is C => {}`, func(t *test, tk Tokens) { // 47
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ArrowFunction: &ArrowFunction{
													FormalParameters: &FormalParameters{
														FormalParameterList: []BindingElement{
															{
																SingleNameBinding: &tk[7],
																Tokens:            tk[7:8],
															},
														},
														Tokens: tk[6:9],
													},
													FunctionBody: &Block{
														Tokens: tk[19:21],
													},
													Tokens: tk[6:21],
												},
												Tokens: tk[6:21],
											},
											Tokens: tk[2:21],
										},
									},
									Tokens: tk[:21],
								},
								Tokens: tk[:21],
							},
							Tokens: tk[:21],
						},
						Tokens: tk[:21],
					},
				},
				Tokens: tk[:21],
			}
		}},
		{`a!`, func(t *test, tk Tokens) { // 48
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: &tk[0],
													Tokens:              tk[:1],
												},
												Tokens: tk[:2],
											}),
											Tokens: tk[:2],
										},
									},
									Tokens: tk[:2],
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
		{`a!.b`, func(t *test, tk Tokens) { // 49
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&MemberExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[0],
														Tokens:              tk[0:1],
													},
													Tokens: tk[0:2],
												},
												IdentifierName: &tk[3],
												Tokens:         tk[:4],
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
					},
				},
				Tokens: tk[:4],
			}
		}},
		{`a!.b!`, func(t *test, tk Tokens) { // 50
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&MemberExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[0],
														Tokens:              tk[0:1],
													},
													Tokens: tk[0:2],
												},
												IdentifierName: &tk[3],
												Tokens:         tk[:5],
											}),
											Tokens: tk[:5],
										},
									},
									Tokens: tk[:5],
								},
								Tokens: tk[:5],
							},
							Tokens: tk[:5],
						},
						Tokens: tk[:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{`a[0]!`, func(t *test, tk Tokens) { // 51
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&MemberExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[0],
														Tokens:              tk[:1],
													},
													Tokens: tk[:1],
												},
												Expression: &Expression{
													Expressions: []AssignmentExpression{
														{
															ConditionalExpression: WrapConditional(&PrimaryExpression{
																Literal: &tk[2],
																Tokens:  tk[2:3],
															}),
															Tokens: tk[2:3],
														},
													},
													Tokens: tk[2:3],
												},
												Tokens: tk[:5],
											}),
											Tokens: tk[:5],
										},
									},
									Tokens: tk[:5],
								},
								Tokens: tk[:5],
							},
							Tokens: tk[:5],
						},
						Tokens: tk[:5],
					},
				},
				Tokens: tk[:5],
			}
		}},
		{`a[0]!.b`, func(t *test, tk Tokens) { // 52
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&MemberExpression{
												MemberExpression: &MemberExpression{
													MemberExpression: &MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: &tk[0],
															Tokens:              tk[:1],
														},
														Tokens: tk[:1],
													},
													Expression: &Expression{
														Expressions: []AssignmentExpression{
															{
																ConditionalExpression: WrapConditional(&PrimaryExpression{
																	Literal: &tk[2],
																	Tokens:  tk[2:3],
																}),
																Tokens: tk[2:3],
															},
														},
														Tokens: tk[2:3],
													},
													Tokens: tk[:5],
												},
												IdentifierName: &tk[6],
												Tokens:         tk[:7],
											}),
											Tokens: tk[:7],
										},
									},
									Tokens: tk[:7],
								},
								Tokens: tk[:7],
							},
							Tokens: tk[:7],
						},
						Tokens: tk[:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{`a!()`, func(t *test, tk Tokens) { // 53
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&CallExpression{
												MemberExpression: &MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: &tk[0],
														Tokens:              tk[:1],
													},
													Tokens: tk[:2],
												},
												Arguments: &Arguments{
													Tokens: tk[2:4],
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
					},
				},
				Tokens: tk[:4],
			}
		}},
		{`a!().b`, func(t *test, tk Tokens) { // 54
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&CallExpression{
												CallExpression: &CallExpression{
													MemberExpression: &MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: &tk[0],
															Tokens:              tk[:1],
														},
														Tokens: tk[:2],
													},
													Arguments: &Arguments{
														Tokens: tk[2:4],
													},
													Tokens: tk[:4],
												},
												IdentifierName: &tk[5],
												Tokens:         tk[:6],
											}),
											Tokens: tk[:6],
										},
									},
									Tokens: tk[:6],
								},
								Tokens: tk[:6],
							},
							Tokens: tk[:6],
						},
						Tokens: tk[:6],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{`a!.b()`, func(t *test, tk Tokens) { // 55
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&CallExpression{
												MemberExpression: &MemberExpression{
													MemberExpression: &MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: &tk[0],
															Tokens:              tk[:1],
														},
														Tokens: tk[:2],
													},
													IdentifierName: &tk[3],
													Tokens:         tk[:4],
												},
												Arguments: &Arguments{
													Tokens: tk[4:6],
												},
												Tokens: tk[:6],
											}),
											Tokens: tk[:6],
										},
									},
									Tokens: tk[:6],
								},
								Tokens: tk[:6],
							},
							Tokens: tk[:6],
						},
						Tokens: tk[:6],
					},
				},
				Tokens: tk[:6],
			}
		}},
		{`a[0]!()`, func(t *test, tk Tokens) { // 56
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&CallExpression{
												MemberExpression: &MemberExpression{
													MemberExpression: &MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: &tk[0],
															Tokens:              tk[:1],
														},
														Tokens: tk[:1],
													},
													Expression: &Expression{
														Expressions: []AssignmentExpression{
															{
																ConditionalExpression: WrapConditional(&PrimaryExpression{
																	Literal: &tk[2],
																	Tokens:  tk[2:3],
																}),
																Tokens: tk[2:3],
															},
														},
														Tokens: tk[2:3],
													},
													Tokens: tk[:5],
												},
												Arguments: &Arguments{
													Tokens: tk[5:7],
												},
												Tokens: tk[:7],
											}),
											Tokens: tk[:7],
										},
									},
									Tokens: tk[:7],
								},
								Tokens: tk[:7],
							},
							Tokens: tk[:7],
						},
						Tokens: tk[:7],
					},
				},
				Tokens: tk[:7],
			}
		}},
		{`a[0]!().b`, func(t *test, tk Tokens) { // 57
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								ExpressionStatement: &Expression{
									Expressions: []AssignmentExpression{
										{
											ConditionalExpression: WrapConditional(&CallExpression{
												CallExpression: &CallExpression{
													MemberExpression: &MemberExpression{
														MemberExpression: &MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																IdentifierReference: &tk[0],
																Tokens:              tk[:1],
															},
															Tokens: tk[:1],
														},
														Expression: &Expression{
															Expressions: []AssignmentExpression{
																{
																	ConditionalExpression: WrapConditional(&PrimaryExpression{
																		Literal: &tk[2],
																		Tokens:  tk[2:3],
																	}),
																	Tokens: tk[2:3],
																},
															},
															Tokens: tk[2:3],
														},
														Tokens: tk[:5],
													},
													Arguments: &Arguments{
														Tokens: tk[5:7],
													},
													Tokens: tk[:7],
												},
												IdentifierName: &tk[8],
												Tokens:         tk[:9],
											}),
											Tokens: tk[:9],
										},
									},
									Tokens: tk[:9],
								},
								Tokens: tk[:9],
							},
							Tokens: tk[:9],
						},
						Tokens: tk[:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{`type A = {[B]: any}`, func(t *test, tk Tokens) { // 58
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:14],
							},
							Tokens: tk[:14],
						},
						Tokens: tk[:14],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`class A {#B: string}`, func(t *test, tk Tokens) { // 59
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
													PrivateIdentifier: &tk[5],
													Tokens:            tk[5:6],
												},
												Tokens: tk[5:9],
											},
											Tokens: tk[5:9],
										},
									},
									Tokens: tk[:10],
								},
								Tokens: tk[:10],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`abstract class A {}`, func(t *test, tk Tokens) { // 60
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[4],
									Tokens:            tk[:8],
								},
								Tokens: tk[:8],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{`export default abstract class A {}`, func(t *test, tk Tokens) { // 61
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							DefaultClass: &ClassDeclaration{
								BindingIdentifier: &tk[8],
								Tokens:            tk[4:12],
							},
							Tokens: tk[:12],
						},
						Tokens: tk[:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{`const a = abstract class {}`, func(t *test, tk Tokens) { // 62
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ConditionalExpression: WrapConditional(&PrimaryExpression{
													ClassExpression: &ClassDeclaration{
														Tokens: tk[6:12],
													},
													Tokens: tk[6:12],
												}),
												Tokens: tk[6:12],
											},
											Tokens: tk[2:12],
										},
									},
									Tokens: tk[:12],
								},
								Tokens: tk[:12],
							},
							Tokens: tk[:12],
						},
						Tokens: tk[:12],
					},
				},
				Tokens: tk[:12],
			}
		}},
		{`abstract class A {
abstract a(): string;
abstract b;
abstract c: number;
public abstract d;
}`, func(t *test, tk Tokens) { // 63
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								ClassDeclaration: &ClassDeclaration{
									BindingIdentifier: &tk[4],
									Tokens:            tk[:39],
								},
								Tokens: tk[:39],
							},
							Tokens: tk[:39],
						},
						Tokens: tk[:39],
					},
				},
				Tokens: tk[:39],
			}
		}},
		{`export type A = B`, func(t *test, tk Tokens) { // 64
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:9],
							},
							Tokens: tk[:9],
						},
						Tokens: tk[:9],
					},
				},
				Tokens: tk[:9],
			}
		}},
		{`class A {b!: string}`, func(t *test, tk Tokens) { // 65
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
														LiteralPropertyName: &tk[5],
														Tokens:              tk[5:6],
													},
													Tokens: tk[5:6],
												},
												Tokens: tk[5:10],
											},
											Tokens: tk[5:10],
										},
									},
									Tokens: tk[:11],
								},
								Tokens: tk[:11],
							},
							Tokens: tk[:11],
						},
						Tokens: tk[:11],
					},
				},
				Tokens: tk[:11],
			}
		}},
		{`export let a: number;`, func(t *test, tk Tokens) { // 66
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						ExportDeclaration: &ExportDeclaration{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[4],
											Tokens:            tk[4:8],
										},
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
		{`const a: (b: any) => b is C`, func(t *test, tk Tokens) { // 67
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Tokens:            tk[2:19],
										},
									},
									Tokens: tk[:19],
								},
								Tokens: tk[:19],
							},
							Tokens: tk[:19],
						},
						Tokens: tk[:19],
					},
				},
				Tokens: tk[:19],
			}
		}},
		{`const a = (b?) => true`, func(t *test, tk Tokens) { // 68
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ArrowFunction: &ArrowFunction{
													FormalParameters: &FormalParameters{
														FormalParameterList: []BindingElement{
															{
																SingleNameBinding: &tk[7],
																Tokens:            tk[7:8],
															},
														},
														Tokens: tk[6:10],
													},
													AssignmentExpression: &AssignmentExpression{
														ConditionalExpression: WrapConditional(&PrimaryExpression{
															Literal: &tk[13],
															Tokens:  tk[13:14],
														}),
														Tokens: tk[13:14],
													},
													Tokens: tk[6:14],
												},
												Tokens: tk[6:14],
											},
											Tokens: tk[2:14],
										},
									},
									Tokens: tk[:14],
								},
								Tokens: tk[:14],
							},
							Tokens: tk[:14],
						},
						Tokens: tk[:14],
					},
				},
				Tokens: tk[:14],
			}
		}},
		{`const a = (b: any): b is [string, ...number[]] => false`, func(t *test, tk Tokens) { // 69
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ArrowFunction: &ArrowFunction{
													FormalParameters: &FormalParameters{
														FormalParameterList: []BindingElement{
															{
																SingleNameBinding: &tk[7],
																Tokens:            tk[7:8],
															},
														},
														Tokens: tk[6:12],
													},
													AssignmentExpression: &AssignmentExpression{
														ConditionalExpression: WrapConditional(&PrimaryExpression{
															Literal: &tk[30],
															Tokens:  tk[30:31],
														}),
														Tokens: tk[30:31],
													},
													Tokens: tk[6:31],
												},
												Tokens: tk[6:31],
											},
											Tokens: tk[2:31],
										},
									},
									Tokens: tk[:31],
								},
								Tokens: tk[:31],
							},
							Tokens: tk[:31],
						},
						Tokens: tk[:31],
					},
				},
				Tokens: tk[:31],
			}
		}},
		{`const a = (b: any = 1) => false`, func(t *test, tk Tokens) { // 70
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Declaration: &Declaration{
								LexicalDeclaration: &LexicalDeclaration{
									LetOrConst: Const,
									BindingList: []LexicalBinding{
										{
											BindingIdentifier: &tk[2],
											Initializer: &AssignmentExpression{
												ArrowFunction: &ArrowFunction{
													FormalParameters: &FormalParameters{
														FormalParameterList: []BindingElement{
															{
																SingleNameBinding: &tk[7],
																Initializer: &AssignmentExpression{
																	ConditionalExpression: WrapConditional(&PrimaryExpression{
																		Literal: &tk[14],
																		Tokens:  tk[14:15],
																	}),
																	Tokens: tk[14:15],
																},
																Tokens: tk[7:15],
															},
														},
														Tokens: tk[6:16],
													},
													AssignmentExpression: &AssignmentExpression{
														ConditionalExpression: WrapConditional(&PrimaryExpression{
															Literal: &tk[19],
															Tokens:  tk[19:20],
														}),
														Tokens: tk[19:20],
													},
													Tokens: tk[6:20],
												},
												Tokens: tk[6:20],
											},
											Tokens: tk[2:20],
										},
									},
									Tokens: tk[:20],
								},
								Tokens: tk[:20],
							},
							Tokens: tk[:20],
						},
						Tokens: tk[:20],
					},
				},
				Tokens: tk[:20],
			}
		}},
		{`switch (a) {
case "b":
	c = 1;
}`, func(t *test, tk Tokens) { // 71
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								SwitchStatement: &SwitchStatement{
									Expression: Expression{
										Expressions: []AssignmentExpression{
											{
												ConditionalExpression: WrapConditional(&PrimaryExpression{
													IdentifierReference: &tk[3],
													Tokens:              tk[3:4],
												}),
												Tokens: tk[3:4],
											},
										},
										Tokens: tk[3:4],
									},
									CaseClauses: []CaseClause{
										{
											Expression: Expression{
												Expressions: []AssignmentExpression{
													{
														ConditionalExpression: WrapConditional(&PrimaryExpression{
															Literal: &tk[10],
															Tokens:  tk[10:11],
														}),
														Tokens: tk[10:11],
													},
												},
												Tokens: tk[10:11],
											},
											StatementList: []StatementListItem{
												{
													Statement: &Statement{
														ExpressionStatement: &Expression{
															Expressions: []AssignmentExpression{
																{
																	LeftHandSideExpression: &LeftHandSideExpression{
																		NewExpression: &NewExpression{
																			MemberExpression: MemberExpression{
																				PrimaryExpression: &PrimaryExpression{
																					IdentifierReference: &tk[14],
																					Tokens:              tk[14:15],
																				},
																				Tokens: tk[14:15],
																			},
																			Tokens: tk[14:15],
																		},
																		Tokens: tk[14:15],
																	},
																	AssignmentOperator: AssignmentAssign,
																	AssignmentExpression: &AssignmentExpression{
																		ConditionalExpression: WrapConditional(&PrimaryExpression{
																			Literal: &tk[18],
																			Tokens:  tk[18:19],
																		}),
																		Tokens: tk[18:19],
																	},
																	Tokens: tk[14:19],
																},
															},
															Tokens: tk[14:19],
														},
														Tokens: tk[14:20],
													},
													Tokens: tk[14:20],
												},
											},
											Tokens: tk[8:20],
										},
									},
									Tokens: tk[:22],
								},
								Tokens: tk[:22],
							},
							Tokens: tk[:22],
						},
						Tokens: tk[:22],
					},
				},
				Tokens: tk[:22],
			}
		}},
		{`type A = {[b]?: C[];}`, func(t *test, tk Tokens) { // 72
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:18],
							},
							Tokens: tk[:18],
						},
						Tokens: tk[:18],
					},
				},
				Tokens: tk[:18],
			}
		}},
		{`export type {A, B}`, func(t *test, tk Tokens) { // 73
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:10],
							},
							Tokens: tk[:10],
						},
						Tokens: tk[:10],
					},
				},
				Tokens: tk[:10],
			}
		}},
		{`export interface A {}`, func(t *test, tk Tokens) { // 74
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:8],
							},
							Tokens: tk[:8],
						},
						Tokens: tk[:8],
					},
				},
				Tokens: tk[:8],
			}
		}},
		{`type A = (this: B, c: D) => void;`, func(t *test, tk Tokens) { // 75
			t.Typescript = true
			t.Output = Module{
				ModuleListItems: []ModuleItem{
					{
						StatementListItem: &StatementListItem{
							Statement: &Statement{
								Tokens: tk[:23],
							},
							Tokens: tk[:23],
						},
						Tokens: tk[:23],
					},
				},
				Tokens: tk[:23],
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

func TestTypescriptTypes(t *testing.T) {
	for n, test := range [...]struct {
		Fn    func(*jsParser) bool
		Input string
	}{
		{ // 1
			(*jsParser).ReadLiteralType,
			"1",
		},
		{ // 2
			(*jsParser).ReadLiteralType,
			"2",
		},
		{ // 3
			(*jsParser).ReadLiteralType,
			"-1",
		},
		{ // 4
			(*jsParser).ReadLiteralType,
			"null",
		},
		{ // 5
			(*jsParser).ReadLiteralType,
			"true",
		},
		{ // 6
			(*jsParser).ReadLiteralType,
			"false",
		},
		{ // 7
			(*jsParser).ReadLiteralType,
			"\"\"",
		},
		{ // 8
			(*jsParser).ReadLiteralType,
			"\"string\"",
		},
		{ // 9
			(*jsParser).ReadLiteralType,
			"``",
		},
		{ // 10
			(*jsParser).ReadLiteralType,
			"`template`",
		},
		{ // 11
			(*jsParser).ReadPredefinedType,
			"void",
		},
		{ // 12
			(*jsParser).ReadPredefinedType,
			"any",
		},
		{ // 13
			(*jsParser).ReadPredefinedType,
			"number",
		},
		{ // 14
			(*jsParser).ReadPredefinedType,
			"boolean",
		},
		{ // 15
			(*jsParser).ReadPredefinedType,
			"string",
		},
		{ // 16
			(*jsParser).ReadPredefinedType,
			"symbol",
		},
		{ // 17
			(*jsParser).ReadPredefinedType,
			"unknown",
		},
		{ // 18
			(*jsParser).ReadPredefinedType,
			"bigint",
		},
		{ // 19
			(*jsParser).ReadPredefinedType,
			"undefined",
		},
		{ // 20
			(*jsParser).ReadPredefinedType,
			"never",
		},
		{ // 21
			(*jsParser).ReadPredefinedType,
			"object",
		},
		{ // 22
			(*jsParser).ReadThisType,
			"this",
		},
		{ // 23
			(*jsParser).ReadTupleType,
			"[any]",
		},
		{ // 24
			(*jsParser).ReadTupleType,
			"[ any ]",
		},
		{ // 25
			(*jsParser).ReadTupleType,
			"[number,bigint]",
		},
		{ // 26
			(*jsParser).ReadTupleType,
			"[ number , bigint ]",
		},
		{ // 27
			(*jsParser).ReadTupleType,
			"[string,1,...symbol]",
		},
		{ // 28
			(*jsParser).ReadTupleType,
			"[ string , 1 , ... symbol ]",
		},
		{ // 29
			(*jsParser).ReadTupleType,
			"[...boolean]",
		},
		{ // 30
			(*jsParser).ReadTupleType,
			"[ ... boolean ]",
		},
		{ // 31
			(*jsParser).ReadTemplateType,
			"`A${number}B`",
		},
		{ // 32
			(*jsParser).ReadTemplateType,
			"`A${ number }B`",
		},
		{ // 33
			(*jsParser).ReadTemplateType,
			"`A${string}B${boolean}`",
		},
		{ // 34
			(*jsParser).ReadTemplateType,
			"`A${ string }B${ boolean }`",
		},
		{ // 35
			(*jsParser).ReadParenthesizedType,
			"(number)",
		},
		{ // 36
			(*jsParser).ReadParenthesizedType,
			"( number )",
		},
		{ // 37
			(*jsParser).ReadObjectType,
			"{}",
		},
		{ // 38
			(*jsParser).ReadObjectType,
			"{ }",
		},
		{ // 39
			(*jsParser).ReadObjectType,
			"{a:number}",
		},
		{ // 40
			(*jsParser).ReadObjectType,
			"{ a : number }",
		},
		{ // 41
			(*jsParser).ReadObjectType,
			"{a: number;}",
		},
		{ // 42
			(*jsParser).ReadObjectType,
			"{a: number,}",
		},
		{ // 43
			(*jsParser).ReadObjectType,
			"{a: number; }",
		},
		{ // 44
			(*jsParser).ReadObjectType,
			"{a: number, }",
		},
		{ // 45
			(*jsParser).ReadObjectType,
			"{a: number;b: string}",
		},
		{ // 46
			(*jsParser).ReadObjectType,
			"{a: number; b: string}",
		},
		{ // 47
			(*jsParser).ReadObjectType,
			"{a: number,b: string}",
		},
		{ // 48
			(*jsParser).ReadObjectType,
			"{a: number, b: string}",
		},
		{ // 49
			(*jsParser).ReadObjectType,
			"{(a: number)}",
		},
		{ // 50
			(*jsParser).ReadObjectType,
			"{ < B > (a: B)}",
		},
		{ // 51
			(*jsParser).ReadObjectType,
			"{(a: number) : string}",
		},
		{ // 52
			(*jsParser).ReadObjectType,
			"{ < B > () : B}",
		},
		{ // 53
			(*jsParser).ReadObjectType,
			"{new (a: number)}",
		},
		{ // 54
			(*jsParser).ReadObjectType,
			"{ new < B > (a: B)}",
		},
		{ // 55
			(*jsParser).ReadObjectType,
			"{new(a: number) : string}",
		},
		{ // 56
			(*jsParser).ReadObjectType,
			"{ new < B > () : B}",
		},
		{ // 57
			(*jsParser).ReadObjectType,
			"{get(a: number)}",
		},
		{ // 57
			(*jsParser).ReadObjectType,
			"{ get < B > (a: B)}",
		},
		{ // 59
			(*jsParser).ReadObjectType,
			"{get(a: number) : string}",
		},
		{ // 60
			(*jsParser).ReadObjectType,
			"{ get < B > () : B}",
		},
		{ // 61
			(*jsParser).ReadObjectType,
			"{set(a: number)}",
		},
		{ // 62
			(*jsParser).ReadObjectType,
			"{ set < B > (a: B)}",
		},
		{ // 63
			(*jsParser).ReadObjectType,
			"{set(a: number) : string}",
		},
		{ // 64
			(*jsParser).ReadObjectType,
			"{ set < B > () : B}",
		},
		{ // 65
			(*jsParser).ReadObjectType,
			"{ [ A ] }",
		},
		{ // 66
			(*jsParser).ReadObjectType,
			"{[A: boolean]: bigint}",
		},
		{ // 67
			(*jsParser).ReadObjectType,
			"{ [ A : number ] ? : string ; }",
		},
		{ // 68
			(*jsParser).ReadObjectType,
			"{ [ const A : number ] ? : string ; }",
		},
		{ // 69
			(*jsParser).ReadObjectType,
			"{ [ static A : number ] ? : string ; }",
		},
		{ // 70
			(*jsParser).ReadObjectType,
			"{[const static A: number]: string}",
		},
		{ // 71
			(*jsParser).ReadObjectType,
			"{[static const A: number]: string}",
		},
		{ // 72
			(*jsParser).ReadObjectType,
			"{a()}",
		},
		{ // 73
			(*jsParser).ReadObjectType,
			"{a?(): string}",
		},
		{ // 74
			(*jsParser).ReadObjectType,
			"{a<B>(c: D): E}",
		},
		{ // 75
			(*jsParser).ReadObjectType,
			"{ a ? < B > ( c : D ) : E ; }",
		},
		{ // 76
			(*jsParser).ReadObjectType,
			"{'a'()}",
		},
		{ // 77
			(*jsParser).ReadObjectType,
			"{0()}",
		},
		{ // 78
			(*jsParser).ReadObjectType,
			"{0}",
		},
		{ // 79
			(*jsParser).ReadObjectType,
			"{''}",
		},
		{ // 80
			(*jsParser).ReadObjectType,
			"{0?: number}",
		},
		{ // 81
			(*jsParser).ReadObjectType,
			"{'': string}",
		},
		{ // 82
			(*jsParser).ReadObjectType,
			"{a}",
		},
		{ // 83
			(*jsParser).ReadObjectType,
			"{ a : boolean ; }",
		},
		{ // 84
			(*jsParser).ReadMappedType,
			"{readonly [A in B]}",
		},
		{ // 85
			(*jsParser).ReadMappedType,
			"{+readonly [A in B]}",
		},
		{ // 86
			(*jsParser).ReadMappedType,
			"{-readonly [A in B]}",
		},
		{ // 87
			(*jsParser).ReadMappedType,
			"{[A in B]?}",
		},
		{ // 88
			(*jsParser).ReadMappedType,
			"{[A in B]-?}",
		},
		{ // 89
			(*jsParser).ReadMappedType,
			"{[A in B]?: string}",
		},
		{ // 90
			(*jsParser).ReadMappedType,
			"{[A in B]-?: number}",
		},
		{ // 91
			(*jsParser).ReadMappedType,
			"{[A in B as C]}",
		},
		{ // 92
			(*jsParser).ReadTupleType,
			"[]",
		},
		{ // 93
			(*jsParser).ReadTupleType,
			"[ ]",
		},
		{ // 94
			(*jsParser).ReadTupleType,
			"[number]",
		},
		{ // 95
			(*jsParser).ReadTupleType,
			"[ number ]",
		},
		{ // 96
			(*jsParser).ReadTupleType,
			"[number, string]",
		},
		{ // 97
			(*jsParser).ReadTupleType,
			"[ number, string ]",
		},
		{ // 98
			(*jsParser).ReadTupleType,
			"[...number]",
		},
		{ // 99
			(*jsParser).ReadTupleType,
			"[number, ...string]",
		},
		{ // 100
			(*jsParser).ReadThisType,
			"this",
		},
		{ // 101
			(*jsParser).ReadTypeQuery,
			"typeof A",
		},
		{ // 102
			(*jsParser).ReadTypeQuery,
			"typeof A<B>",
		},
		{ // 103
			(*jsParser).ReadTypeQuery,
			"typeof A.B",
		},
		{ // 104
			(*jsParser).ReadTypeQuery,
			"typeof A.const",
		},
		{ // 105
			(*jsParser).ReadTypeQuery,
			"typeof A.#B",
		},
		{ // 106
			(*jsParser).ReadTypeQuery,
			"typeof A . #B . void",
		},
		{ // 107
			(*jsParser).ReadTypeQuery,
			"typeof A.#B.void<C, D>",
		},
		{ // 108
			(*jsParser).ReadTypeReference,
			"A",
		},
		{ // 109
			(*jsParser).ReadTypeReference,
			"A.B",
		},
		{ // 110
			(*jsParser).ReadTypeReference,
			"A<B>",
		},
		{ // 110
			(*jsParser).ReadTypeReference,
			"A.B<C>",
		},
		{ // 111
			(*jsParser).ReadTypeReference,
			"A . B < C >",
		},
	} {
		tk := parser.NewStringTokeniser(test.Input)
		j, err := newJSParser(&tk)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			j[:cap(j)][cap(j)-1].Data = marker
			g := j
			if !test.Fn(&j) {
				t.Errorf("test %d: failed on specific type fn ", n+1)
			} else if !g.ReadType() {
				t.Errorf("test %d: failed on generic type fn", n+1)
			} else if len(j) != len(g) {
				t.Errorf("test %d: inconsistant number of tokens read. %d != %d", n+1, len(j), len(g))
			}
		}
	}
}
