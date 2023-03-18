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
	}, func(t *test) (Type, error) {
		if t.Typescript {
			t.Tokens[:cap(t.Tokens)][cap(t.Tokens)-1].Data = marker
		}
		var m Module
		err := m.parse(&t.Tokens)
		return m, err
	})
}
