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
		{
			(*jsParser).ReadLiteralType,
			"1",
		},
		{
			(*jsParser).ReadLiteralType,
			"2",
		},
		{
			(*jsParser).ReadLiteralType,
			"-1",
		},
		{
			(*jsParser).ReadLiteralType,
			"null",
		},
		{
			(*jsParser).ReadLiteralType,
			"true",
		},
		{
			(*jsParser).ReadLiteralType,
			"false",
		},
		{
			(*jsParser).ReadLiteralType,
			"\"\"",
		},
		{
			(*jsParser).ReadLiteralType,
			"\"string\"",
		},
		{
			(*jsParser).ReadLiteralType,
			"``",
		},
		{
			(*jsParser).ReadLiteralType,
			"`template`",
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
			}
		}
	}
}
