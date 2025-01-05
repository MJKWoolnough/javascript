package javascript

import (
	"errors"
	"reflect"
	"testing"

	"vimagination.zapto.org/parser"
)

func TestUnquote(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{ // 1
			"\"\"",
			"",
			nil,
		},
		{ // 2
			"''",
			"",
			nil,
		},
		{ // 3
			"\"a\"",
			"a",
			nil,
		},
		{ // 4
			"'a'",
			"a",
			nil,
		},
		{ // 5
			"\"\\'\\\"\\\\\\b\\f\\n\\r\\t\\v\"",
			"'\"\\\b\f\n\r\t\v",
			nil,
		},
		{ // 6
			"\"\x41\"",
			"A",
			nil,
		},
		{ // 7
			"",
			"",
			ErrInvalidQuoted,
		},
		{ // 8
			"\"\\x41\"",
			"A",
			nil,
		},
		{ // 9
			"\"\\x4G\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 10
			"\"\\xG1\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 11
			"\"\\u0041\"",
			"A",
			nil,
		},
		{ // 12
			"\"\\u004G\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 13
			"\"\\u00G1\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 14
			"\"\\u0G41\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 15
			"\"\\uG041\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 16
			"\"\\c\"",
			"c",
			nil,
		},
		{ // 17
			"\"\n\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 18
			"\"\\0\"",
			"\000",
			nil,
		},
		{ // 19
			"\"\\01\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 20
			"\"\\u{41}\"",
			"A",
			nil,
		},
		{ // 21
			"\"\\u{}\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 22
			"\"\\u{41G}\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 23
			"\"\\u{41\"",
			"",
			ErrInvalidQuoted,
		},
		{ // 24
			"'\\'\\\"\\\\\\b\\f\\n\\r\\t\\v'",
			"'\"\\\b\f\n\r\t\v",
			nil,
		},
		{ // 25
			"'\x41'",
			"A",
			nil,
		},
		{ // 26
			"'\\x41'",
			"A",
			nil,
		},
		{ // 27
			"'\\x4G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 28
			"'\\xG1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 29
			"'\\u0041'",
			"A",
			nil,
		},
		{ // 30
			"'\\u004G'",
			"",
			ErrInvalidQuoted,
		},
		{ // 31
			"'\\u00G1'",
			"",
			ErrInvalidQuoted,
		},
		{ // 32
			"'\\u0G41'",
			"",
			ErrInvalidQuoted,
		},
		{ // 33
			"'\\uG041'",
			"",
			ErrInvalidQuoted,
		},
		{ // 34
			"'\\c'",
			"c",
			nil,
		},
		{ // 35
			"'\n'",
			"",
			ErrInvalidQuoted,
		},
		{ // 36
			"'\\0'",
			"\000",
			nil,
		},
		{ // 37
			"'\\01'",
			"",
			ErrInvalidQuoted,
		},
		{ // 38
			"'\\u{41}'",
			"A",
			nil,
		},
		{ // 39
			"'\\u{}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 40
			"'\\u{41G}'",
			"",
			ErrInvalidQuoted,
		},
		{ // 41
			"'\\u{41'",
			"",
			ErrInvalidQuoted,
		},
	} {
		if o, err := Unquote(test.Input); !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.Err, err)
		} else if o != test.Output {
			t.Errorf("test %d: from %s, expecting output %q, got %q", n+1, test.Input, test.Output, o)
		}
	}
}

func TestUnquoteTemplate(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		Err           error
	}{
		{ // 1
			"``",
			"",
			nil,
		},
		{ // 2
			"}`",
			"",
			nil,
		},
		{ // 3
			"`${",
			"",
			nil,
		},
		{ // 4
			"}${",
			"",
			nil,
		},
		{ // 5
			"`",
			"",
			ErrInvalidQuoted,
		},
		{ // 6
			"}",
			"",
			ErrInvalidQuoted,
		},
		{ // 7
			"${",
			"",
			ErrInvalidQuoted,
		},
		{ // 8
			"`a`",
			"a",
			nil,
		},
		{ // 9
			"`\\'\\\"\\\\\\b\\f\\n\\r\\t\\v`",
			"'\"\\\b\f\n\r\t\v",
			nil,
		},
		{ // 10
			"`\x41`",
			"A",
			nil,
		},
		{ // 11
			"`\n`",
			"\n",
			nil,
		},
		{ // 12
			"`\\x4G`",
			"",
			ErrInvalidQuoted,
		},
		{ // 13
			"`\\u0041`",
			"A",
			nil,
		},
		{ // 14
			"`\\u00G1`",
			"",
			ErrInvalidQuoted,
		},
		{ // 15
			"`\\c`",
			"c",
			nil,
		},
		{ // 16
			"`\\0`",
			"\000",
			nil,
		},
		{ // 17
			"`\\u{41}`",
			"A",
			nil,
		},
		{ // 18
			"`\\u{}`",
			"",
			ErrInvalidQuoted,
		},
		{ // 19
			"`\\${`",
			"${",
			nil,
		},
	} {
		o, err := UnquoteTemplate(test.Input)
		if !errors.Is(err, test.Err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.Err, err)
		} else if o != test.Output {
			t.Errorf("test %d: from %s, expecting output %q, got %q", n+1, test.Input, test.Output, o)
		}
	}
}

func TestQuoteTemplate(t *testing.T) {
	for n, test := range [...]struct {
		Input, Output string
		templateType  TemplateType
	}{
		{ // 1
			"",
			"``",
			TemplateNoSubstitution,
		},
		{ // 2
			"a",
			"`a`",
			TemplateNoSubstitution,
		},
		{ // 3
			"abc",
			"`abc`",
			TemplateNoSubstitution,
		},
		{ // 4
			"a\nb	c",
			"`a\nb	c`",
			TemplateNoSubstitution,
		},
		{ // 5
			"\\n",
			"`\\\\n`",
			TemplateNoSubstitution,
		},
		{ // 6
			"a$b",
			"`a$b`",
			TemplateNoSubstitution,
		},
		{ // 7
			"a${b",
			"`a\\${b`",
			TemplateNoSubstitution,
		},
		{ // 8
			"`",
			"`\\``",
			TemplateNoSubstitution,
		},
		{ // 9
			"",
			"`${",
			TemplateHead,
		},
		{ // 10
			"a",
			"`a${",
			TemplateHead,
		},
		{ // 11
			"abc",
			"`abc${",
			TemplateHead,
		},
		{ // 12
			"a\nb	c",
			"`a\nb	c${",
			TemplateHead,
		},
		{ // 13
			"\\n",
			"`\\\\n${",
			TemplateHead,
		},
		{ // 14
			"a$b",
			"`a$b${",
			TemplateHead,
		},
		{ // 15
			"a${b",
			"`a\\${b${",
			TemplateHead,
		},
		{ // 16
			"`",
			"`\\`${",
			TemplateHead,
		},
		{ // 17
			"",
			"}${",
			TemplateMiddle,
		},
		{ // 18
			"a",
			"}a${",
			TemplateMiddle,
		},
		{ // 19
			"abc",
			"}abc${",
			TemplateMiddle,
		},
		{ // 20
			"a\nb	c",
			"}a\nb	c${",
			TemplateMiddle,
		},
		{ // 21
			"\\n",
			"}\\\\n${",
			TemplateMiddle,
		},
		{ // 22
			"a$b",
			"}a$b${",
			TemplateMiddle,
		},
		{ // 23
			"a${b",
			"}a\\${b${",
			TemplateMiddle,
		},
		{ // 24
			"`",
			"}\\`${",
			TemplateMiddle,
		},
		{ // 25
			"",
			"}`",
			TemplateTail,
		},
		{ // 26
			"a",
			"}a`",
			TemplateTail,
		},
		{ // 27
			"abc",
			"}abc`",
			TemplateTail,
		},
		{ // 28
			"a\nb	c",
			"}a\nb	c`",
			TemplateTail,
		},
		{ // 29
			"\\n",
			"}\\\\n`",
			TemplateTail,
		},
		{ // 30
			"a$b",
			"}a$b`",
			TemplateTail,
		},
		{ // 31
			"a${b",
			"}a\\${b`",
			TemplateTail,
		},
		{ // 32
			"`",
			"}\\``",
			TemplateTail,
		},
	} {
		if out := QuoteTemplate(test.Input, test.templateType); out != test.Output {
			t.Errorf("test %d: expecting output %s, got %s", n+1, test.Output, out)
		}
	}
}

func TestWrapConditional(t *testing.T) {
	tks := Tokens{
		{
			Token: parser.Token{
				Type: TokenNoSubstitutionTemplate,
				Data: "`abc`",
			},
		},
	}
	template := &TemplateLiteral{
		NoSubstitutionTemplate: &tks[0],
		Tokens:                 tks,
	}
	expectedOutput := ConditionalExpression{
		LogicalORExpression: &LogicalORExpression{
			LogicalANDExpression: LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			TemplateLiteral: template,
																			Tokens:          tks,
																		},
																		Tokens: tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		Tokens: tks,
	}

	for n, test := range [...]ConditionalWrappable{
		template,  // 1
		*template, // 2
		&PrimaryExpression{ // 3
			TemplateLiteral: template,
			Tokens:          tks,
		},
		PrimaryExpression{ // 4
			TemplateLiteral: template,
			Tokens:          tks,
		},
		&MemberExpression{ // 5
			PrimaryExpression: &PrimaryExpression{
				TemplateLiteral: template,
				Tokens:          tks,
			},
			Tokens: tks,
		},
		MemberExpression{ // 6
			PrimaryExpression: &PrimaryExpression{
				TemplateLiteral: template,
				Tokens:          tks,
			},
			Tokens: tks,
		},
		&NewExpression{ // 7
			MemberExpression: MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					TemplateLiteral: template,
					Tokens:          tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		NewExpression{ // 8
			MemberExpression: MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					TemplateLiteral: template,
					Tokens:          tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&LeftHandSideExpression{ // 9
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						TemplateLiteral: template,
						Tokens:          tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		LeftHandSideExpression{ // 10
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						TemplateLiteral: template,
						Tokens:          tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&UpdateExpression{ // 11
			LeftHandSideExpression: &LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							TemplateLiteral: template,
							Tokens:          tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		UpdateExpression{ // 12
			LeftHandSideExpression: &LeftHandSideExpression{
				NewExpression: &NewExpression{
					MemberExpression: MemberExpression{
						PrimaryExpression: &PrimaryExpression{
							TemplateLiteral: template,
							Tokens:          tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&UnaryExpression{ // 13
			UpdateExpression: UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								TemplateLiteral: template,
								Tokens:          tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		UnaryExpression{ // 14
			UpdateExpression: UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								TemplateLiteral: template,
								Tokens:          tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&ExponentiationExpression{ // 15
			UnaryExpression: UnaryExpression{
				UpdateExpression: UpdateExpression{
					LeftHandSideExpression: &LeftHandSideExpression{
						NewExpression: &NewExpression{
							MemberExpression: MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									TemplateLiteral: template,
									Tokens:          tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		ExponentiationExpression{ // 16
			UnaryExpression: UnaryExpression{
				UpdateExpression: UpdateExpression{
					LeftHandSideExpression: &LeftHandSideExpression{
						NewExpression: &NewExpression{
							MemberExpression: MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									TemplateLiteral: template,
									Tokens:          tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&MultiplicativeExpression{ // 17
			ExponentiationExpression: ExponentiationExpression{
				UnaryExpression: UnaryExpression{
					UpdateExpression: UpdateExpression{
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										TemplateLiteral: template,
										Tokens:          tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		MultiplicativeExpression{ // 18
			ExponentiationExpression: ExponentiationExpression{
				UnaryExpression: UnaryExpression{
					UpdateExpression: UpdateExpression{
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										TemplateLiteral: template,
										Tokens:          tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&AdditiveExpression{ // 19
			MultiplicativeExpression: MultiplicativeExpression{
				ExponentiationExpression: ExponentiationExpression{
					UnaryExpression: UnaryExpression{
						UpdateExpression: UpdateExpression{
							LeftHandSideExpression: &LeftHandSideExpression{
								NewExpression: &NewExpression{
									MemberExpression: MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											TemplateLiteral: template,
											Tokens:          tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		AdditiveExpression{ // 20
			MultiplicativeExpression: MultiplicativeExpression{
				ExponentiationExpression: ExponentiationExpression{
					UnaryExpression: UnaryExpression{
						UpdateExpression: UpdateExpression{
							LeftHandSideExpression: &LeftHandSideExpression{
								NewExpression: &NewExpression{
									MemberExpression: MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											TemplateLiteral: template,
											Tokens:          tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&ShiftExpression{ // 21
			AdditiveExpression: AdditiveExpression{
				MultiplicativeExpression: MultiplicativeExpression{
					ExponentiationExpression: ExponentiationExpression{
						UnaryExpression: UnaryExpression{
							UpdateExpression: UpdateExpression{
								LeftHandSideExpression: &LeftHandSideExpression{
									NewExpression: &NewExpression{
										MemberExpression: MemberExpression{
											PrimaryExpression: &PrimaryExpression{
												TemplateLiteral: template,
												Tokens:          tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&ShiftExpression{ // 22
			AdditiveExpression: AdditiveExpression{
				MultiplicativeExpression: MultiplicativeExpression{
					ExponentiationExpression: ExponentiationExpression{
						UnaryExpression: UnaryExpression{
							UpdateExpression: UpdateExpression{
								LeftHandSideExpression: &LeftHandSideExpression{
									NewExpression: &NewExpression{
										MemberExpression: MemberExpression{
											PrimaryExpression: &PrimaryExpression{
												TemplateLiteral: template,
												Tokens:          tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&RelationalExpression{ // 23
			ShiftExpression: ShiftExpression{
				AdditiveExpression: AdditiveExpression{
					MultiplicativeExpression: MultiplicativeExpression{
						ExponentiationExpression: ExponentiationExpression{
							UnaryExpression: UnaryExpression{
								UpdateExpression: UpdateExpression{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													TemplateLiteral: template,
													Tokens:          tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		RelationalExpression{ // 24
			ShiftExpression: ShiftExpression{
				AdditiveExpression: AdditiveExpression{
					MultiplicativeExpression: MultiplicativeExpression{
						ExponentiationExpression: ExponentiationExpression{
							UnaryExpression: UnaryExpression{
								UpdateExpression: UpdateExpression{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													TemplateLiteral: template,
													Tokens:          tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&EqualityExpression{ // 25
			RelationalExpression: RelationalExpression{
				ShiftExpression: ShiftExpression{
					AdditiveExpression: AdditiveExpression{
						MultiplicativeExpression: MultiplicativeExpression{
							ExponentiationExpression: ExponentiationExpression{
								UnaryExpression: UnaryExpression{
									UpdateExpression: UpdateExpression{
										LeftHandSideExpression: &LeftHandSideExpression{
											NewExpression: &NewExpression{
												MemberExpression: MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														TemplateLiteral: template,
														Tokens:          tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		EqualityExpression{ // 26
			RelationalExpression: RelationalExpression{
				ShiftExpression: ShiftExpression{
					AdditiveExpression: AdditiveExpression{
						MultiplicativeExpression: MultiplicativeExpression{
							ExponentiationExpression: ExponentiationExpression{
								UnaryExpression: UnaryExpression{
									UpdateExpression: UpdateExpression{
										LeftHandSideExpression: &LeftHandSideExpression{
											NewExpression: &NewExpression{
												MemberExpression: MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														TemplateLiteral: template,
														Tokens:          tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&BitwiseANDExpression{ // 27
			EqualityExpression: EqualityExpression{
				RelationalExpression: RelationalExpression{
					ShiftExpression: ShiftExpression{
						AdditiveExpression: AdditiveExpression{
							MultiplicativeExpression: MultiplicativeExpression{
								ExponentiationExpression: ExponentiationExpression{
									UnaryExpression: UnaryExpression{
										UpdateExpression: UpdateExpression{
											LeftHandSideExpression: &LeftHandSideExpression{
												NewExpression: &NewExpression{
													MemberExpression: MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															TemplateLiteral: template,
															Tokens:          tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		BitwiseANDExpression{ // 28
			EqualityExpression: EqualityExpression{
				RelationalExpression: RelationalExpression{
					ShiftExpression: ShiftExpression{
						AdditiveExpression: AdditiveExpression{
							MultiplicativeExpression: MultiplicativeExpression{
								ExponentiationExpression: ExponentiationExpression{
									UnaryExpression: UnaryExpression{
										UpdateExpression: UpdateExpression{
											LeftHandSideExpression: &LeftHandSideExpression{
												NewExpression: &NewExpression{
													MemberExpression: MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															TemplateLiteral: template,
															Tokens:          tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&BitwiseXORExpression{ // 29
			BitwiseANDExpression: BitwiseANDExpression{
				EqualityExpression: EqualityExpression{
					RelationalExpression: RelationalExpression{
						ShiftExpression: ShiftExpression{
							AdditiveExpression: AdditiveExpression{
								MultiplicativeExpression: MultiplicativeExpression{
									ExponentiationExpression: ExponentiationExpression{
										UnaryExpression: UnaryExpression{
											UpdateExpression: UpdateExpression{
												LeftHandSideExpression: &LeftHandSideExpression{
													NewExpression: &NewExpression{
														MemberExpression: MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																TemplateLiteral: template,
																Tokens:          tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		BitwiseXORExpression{ // 30
			BitwiseANDExpression: BitwiseANDExpression{
				EqualityExpression: EqualityExpression{
					RelationalExpression: RelationalExpression{
						ShiftExpression: ShiftExpression{
							AdditiveExpression: AdditiveExpression{
								MultiplicativeExpression: MultiplicativeExpression{
									ExponentiationExpression: ExponentiationExpression{
										UnaryExpression: UnaryExpression{
											UpdateExpression: UpdateExpression{
												LeftHandSideExpression: &LeftHandSideExpression{
													NewExpression: &NewExpression{
														MemberExpression: MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																TemplateLiteral: template,
																Tokens:          tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&BitwiseORExpression{ // 31
			BitwiseXORExpression: BitwiseXORExpression{
				BitwiseANDExpression: BitwiseANDExpression{
					EqualityExpression: EqualityExpression{
						RelationalExpression: RelationalExpression{
							ShiftExpression: ShiftExpression{
								AdditiveExpression: AdditiveExpression{
									MultiplicativeExpression: MultiplicativeExpression{
										ExponentiationExpression: ExponentiationExpression{
											UnaryExpression: UnaryExpression{
												UpdateExpression: UpdateExpression{
													LeftHandSideExpression: &LeftHandSideExpression{
														NewExpression: &NewExpression{
															MemberExpression: MemberExpression{
																PrimaryExpression: &PrimaryExpression{
																	TemplateLiteral: template,
																	Tokens:          tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		BitwiseORExpression{ // 32
			BitwiseXORExpression: BitwiseXORExpression{
				BitwiseANDExpression: BitwiseANDExpression{
					EqualityExpression: EqualityExpression{
						RelationalExpression: RelationalExpression{
							ShiftExpression: ShiftExpression{
								AdditiveExpression: AdditiveExpression{
									MultiplicativeExpression: MultiplicativeExpression{
										ExponentiationExpression: ExponentiationExpression{
											UnaryExpression: UnaryExpression{
												UpdateExpression: UpdateExpression{
													LeftHandSideExpression: &LeftHandSideExpression{
														NewExpression: &NewExpression{
															MemberExpression: MemberExpression{
																PrimaryExpression: &PrimaryExpression{
																	TemplateLiteral: template,
																	Tokens:          tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&LogicalANDExpression{ // 33
			BitwiseORExpression: BitwiseORExpression{
				BitwiseXORExpression: BitwiseXORExpression{
					BitwiseANDExpression: BitwiseANDExpression{
						EqualityExpression: EqualityExpression{
							RelationalExpression: RelationalExpression{
								ShiftExpression: ShiftExpression{
									AdditiveExpression: AdditiveExpression{
										MultiplicativeExpression: MultiplicativeExpression{
											ExponentiationExpression: ExponentiationExpression{
												UnaryExpression: UnaryExpression{
													UpdateExpression: UpdateExpression{
														LeftHandSideExpression: &LeftHandSideExpression{
															NewExpression: &NewExpression{
																MemberExpression: MemberExpression{
																	PrimaryExpression: &PrimaryExpression{
																		TemplateLiteral: template,
																		Tokens:          tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&LogicalANDExpression{ // 34
			BitwiseORExpression: BitwiseORExpression{
				BitwiseXORExpression: BitwiseXORExpression{
					BitwiseANDExpression: BitwiseANDExpression{
						EqualityExpression: EqualityExpression{
							RelationalExpression: RelationalExpression{
								ShiftExpression: ShiftExpression{
									AdditiveExpression: AdditiveExpression{
										MultiplicativeExpression: MultiplicativeExpression{
											ExponentiationExpression: ExponentiationExpression{
												UnaryExpression: UnaryExpression{
													UpdateExpression: UpdateExpression{
														LeftHandSideExpression: &LeftHandSideExpression{
															NewExpression: &NewExpression{
																MemberExpression: MemberExpression{
																	PrimaryExpression: &PrimaryExpression{
																		TemplateLiteral: template,
																		Tokens:          tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&LogicalORExpression{ // 35
			LogicalANDExpression: LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			TemplateLiteral: template,
																			Tokens:          tks,
																		},
																		Tokens: tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		LogicalORExpression{ // 36
			LogicalANDExpression: LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			TemplateLiteral: template,
																			Tokens:          tks,
																		},
																		Tokens: tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		&ConditionalExpression{ // 37
			LogicalORExpression: &LogicalORExpression{
				LogicalANDExpression: LogicalANDExpression{
					BitwiseORExpression: BitwiseORExpression{
						BitwiseXORExpression: BitwiseXORExpression{
							BitwiseANDExpression: BitwiseANDExpression{
								EqualityExpression: EqualityExpression{
									RelationalExpression: RelationalExpression{
										ShiftExpression: ShiftExpression{
											AdditiveExpression: AdditiveExpression{
												MultiplicativeExpression: MultiplicativeExpression{
													ExponentiationExpression: ExponentiationExpression{
														UnaryExpression: UnaryExpression{
															UpdateExpression: UpdateExpression{
																LeftHandSideExpression: &LeftHandSideExpression{
																	NewExpression: &NewExpression{
																		MemberExpression: MemberExpression{
																			PrimaryExpression: &PrimaryExpression{
																				TemplateLiteral: template,
																				Tokens:          tks,
																			},
																			Tokens: tks,
																		},
																		Tokens: tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
		ConditionalExpression{ // 38
			LogicalORExpression: &LogicalORExpression{
				LogicalANDExpression: LogicalANDExpression{
					BitwiseORExpression: BitwiseORExpression{
						BitwiseXORExpression: BitwiseXORExpression{
							BitwiseANDExpression: BitwiseANDExpression{
								EqualityExpression: EqualityExpression{
									RelationalExpression: RelationalExpression{
										ShiftExpression: ShiftExpression{
											AdditiveExpression: AdditiveExpression{
												MultiplicativeExpression: MultiplicativeExpression{
													ExponentiationExpression: ExponentiationExpression{
														UnaryExpression: UnaryExpression{
															UpdateExpression: UpdateExpression{
																LeftHandSideExpression: &LeftHandSideExpression{
																	NewExpression: &NewExpression{
																		MemberExpression: MemberExpression{
																			PrimaryExpression: &PrimaryExpression{
																				TemplateLiteral: template,
																				Tokens:          tks,
																			},
																			Tokens: tks,
																		},
																		Tokens: tks,
																	},
																	Tokens: tks,
																},
																Tokens: tks,
															},
															Tokens: tks,
														},
														Tokens: tks,
													},
													Tokens: tks,
												},
												Tokens: tks,
											},
											Tokens: tks,
										},
										Tokens: tks,
									},
									Tokens: tks,
								},
								Tokens: tks,
							},
							Tokens: tks,
						},
						Tokens: tks,
					},
					Tokens: tks,
				},
				Tokens: tks,
			},
			Tokens: tks,
		},
	} {
		if output := WrapConditional(test); !reflect.DeepEqual(output, &expectedOutput) {
			t.Errorf("test %d: expecting\n%v\n...got...\n%v", n+1, expectedOutput, output)
		}
	}
}

func TestWrapConditionalExtra(t *testing.T) {
	arrayTokens := Tokens{
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "[",
			},
		},
		{
			Token: parser.Token{
				Type: TokenIdentifier,
				Data: "a",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "]",
			},
		},
	}
	arrayLiteral := &ArrayLiteral{
		ElementList: []ArrayElement{
			{
				AssignmentExpression: AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: &arrayTokens[1],
						Tokens:              arrayTokens[1:2],
					}),
					Tokens: arrayTokens[1:2],
				},
				Tokens: arrayTokens[1:2],
			},
		},
		Tokens: arrayTokens,
	}
	objectTokens := Tokens{
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "{",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "...",
			},
		},
		{
			Token: parser.Token{
				Type: TokenIdentifier,
				Data: "a",
			},
		},
		{
			Token: parser.Token{
				Type: TokenRightBracePunctuator,
				Data: "}",
			},
		},
	}
	objectLiteral := &ObjectLiteral{
		PropertyDefinitionList: []PropertyDefinition{
			{
				AssignmentExpression: &AssignmentExpression{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: &arrayTokens[2],
						Tokens:              arrayTokens[2:3],
					}),
					Tokens: arrayTokens[2:3],
				},
				Tokens: arrayTokens[2:3],
			},
		},
		Tokens: objectTokens,
	}
	functionTokens := Tokens{
		{
			Token: parser.Token{
				Type: TokenKeyword,
				Data: "function",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "(",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: ")",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "{",
			},
		},
		{
			Token: parser.Token{
				Type: TokenRightBracePunctuator,
				Data: "}",
			},
		},
	}
	functionDeclaration := &FunctionDeclaration{
		FormalParameters: FormalParameters{
			Tokens: functionTokens[1:3],
		},
		FunctionBody: Block{
			Tokens: functionTokens[3:5],
		},
		Tokens: functionTokens,
	}
	classTokens := Tokens{
		{
			Token: parser.Token{
				Type: TokenKeyword,
				Data: "class",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "{",
			},
		},
		{
			Token: parser.Token{
				Type: TokenRightBracePunctuator,
				Data: "}",
			},
		},
	}
	classDeclaration := &ClassDeclaration{
		Tokens: classTokens,
	}
	parenthesizedTokens := Tokens{
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: "(",
			},
		},
		{
			Token: parser.Token{
				Type: TokenIdentifier,
				Data: "a",
			},
		},
		{
			Token: parser.Token{
				Type: TokenPunctuator,
				Data: ")",
			},
		},
	}
	parenthesizedExpression := &ParenthesizedExpression{
		Expressions: []AssignmentExpression{
			{
				ConditionalExpression: WrapConditional(&PrimaryExpression{
					IdentifierReference: &arrayTokens[1],
					Tokens:              parenthesizedTokens[1:2],
				}),
				Tokens: parenthesizedTokens[1:2],
			},
		},
		Tokens: parenthesizedTokens,
	}

	for n, test := range [...]struct {
		ConditionalWrappable ConditionalWrappable
		PrimaryExpression    *PrimaryExpression
	}{
		{ // 1
			arrayLiteral,
			&PrimaryExpression{
				ArrayLiteral: arrayLiteral,
				Tokens:       arrayTokens,
			},
		},
		{ // 2
			*arrayLiteral,
			&PrimaryExpression{
				ArrayLiteral: arrayLiteral,
				Tokens:       arrayTokens,
			},
		},
		{ // 3
			objectLiteral,
			&PrimaryExpression{
				ObjectLiteral: objectLiteral,
				Tokens:        objectTokens,
			},
		},
		{ // 4
			*objectLiteral,
			&PrimaryExpression{
				ObjectLiteral: objectLiteral,
				Tokens:        objectTokens,
			},
		},
		{ // 5
			functionDeclaration,
			&PrimaryExpression{
				FunctionExpression: functionDeclaration,
				Tokens:             functionTokens,
			},
		},
		{ // 6
			*functionDeclaration,
			&PrimaryExpression{
				FunctionExpression: functionDeclaration,
				Tokens:             functionTokens,
			},
		},
		{ // 7
			classDeclaration,
			&PrimaryExpression{
				ClassExpression: classDeclaration,
				Tokens:          classTokens,
			},
		},
		{ // 8
			*classDeclaration,
			&PrimaryExpression{
				ClassExpression: classDeclaration,
				Tokens:          classTokens,
			},
		},
		{ // 9
			parenthesizedExpression,
			&PrimaryExpression{
				ParenthesizedExpression: parenthesizedExpression,
				Tokens:                  parenthesizedTokens,
			},
		},
		{ // 10
			*parenthesizedExpression,
			&PrimaryExpression{
				ParenthesizedExpression: parenthesizedExpression,
				Tokens:                  parenthesizedTokens,
			},
		},
	} {
		if output, expectedOutput := WrapConditional(test.ConditionalWrappable), WrapConditional(test.PrimaryExpression); !reflect.DeepEqual(output, expectedOutput) {
			t.Errorf("test %d: expecting\n%v\n...got...\n%v", n+1, expectedOutput, output)
		}
	}
}

func TestUnwrapConditional(t *testing.T) {
	tks := Tokens{
		{
			Token: parser.Token{
				Type: TokenIdentifier,
				Data: "a",
			},
		},
		{
			Token: parser.Token{
				Type: TokenIdentifier,
				Data: "b",
			},
		},
	}
	identA := &tks[0]
	identB := &tks[1]

	for n, test := range [...]ConditionalWrappable{
		&ParenthesizedExpression{ // 1
			Expressions: []AssignmentExpression{
				{
					ConditionalExpression: WrapConditional(&PrimaryExpression{
						IdentifierReference: identA,
						Tokens:              tks[:1],
					}),
					Tokens: tks[:1],
				},
			},
			Tokens: tks[:1],
		},
		&TemplateLiteral{ // 2
			NoSubstitutionTemplate: identA,
			Tokens:                 tks[:1],
		},
		&ClassDeclaration{ // 3
			BindingIdentifier: identA,
			Tokens:            tks[:1],
		},
		&FunctionDeclaration{ // 4
			BindingIdentifier: identA,
			Tokens:            tks[:1],
		},
		&ObjectLiteral{ // 5
			PropertyDefinitionList: []PropertyDefinition{
				{
					AssignmentExpression: &AssignmentExpression{
						ConditionalExpression: WrapConditional(&PrimaryExpression{
							IdentifierReference: identA,
							Tokens:              tks[:1],
						}),
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
			},
			Tokens: tks[:1],
		},
		&ArrayLiteral{ // 6
			ElementList: []ArrayElement{
				{
					AssignmentExpression: AssignmentExpression{
						ConditionalExpression: WrapConditional(&PrimaryExpression{
							IdentifierReference: identA,
							Tokens:              tks[:1],
						}),
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
			},
			Tokens: tks[:1],
		},
		&PrimaryExpression{ // 7
			IdentifierReference: identA,
			Tokens:              tks[:1],
		},
		&MemberExpression{ // 8
			MemberExpression: &MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					IdentifierReference: identA,
					Tokens:              tks[:1],
				},
			},
			IdentifierName: identB,
			Tokens:         tks,
		},
		&NewExpression{ // 9
			News: 1,
			MemberExpression: MemberExpression{
				PrimaryExpression: &PrimaryExpression{
					IdentifierReference: identA,
					Tokens:              tks[:1],
				},
				Tokens: tks[:1],
			},
			Tokens: tks[:1],
		},
		&CallExpression{ // 10
			CallExpression: &CallExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: identA,
						Tokens:              tks[:1],
					},
				},
				Tokens: tks[:1],
			},
			Arguments: &Arguments{
				ArgumentList: []Argument{
					{
						AssignmentExpression: AssignmentExpression{
							ConditionalExpression: WrapConditional(&PrimaryExpression{
								IdentifierReference: identB,
								Tokens:              tks[1:2],
							}),
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&OptionalExpression{ // 11
			OptionalExpression: &OptionalExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						IdentifierReference: identA,
						Tokens:              tks[:1],
					},
				},
				Tokens: tks[:1],
			},
			OptionalChain: OptionalChain{
				IdentifierName: identB,
				Tokens:         tks[1:2],
			},
			Tokens: tks[:2],
		},
		&UpdateExpression{ // 12
			UpdateOperator: UpdatePostIncrement,
			UnaryExpression: &UnaryExpression{
				UpdateExpression: UpdateExpression{
					LeftHandSideExpression: &LeftHandSideExpression{
						NewExpression: &NewExpression{
							MemberExpression: MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: identA,
									Tokens:              tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			Tokens: tks[:1],
		},
		&UnaryExpression{ // 13
			UnaryOperators: []UnaryOperator{UnaryVoid},
			UpdateExpression: UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							PrimaryExpression: &PrimaryExpression{
								IdentifierReference: identA,
								Tokens:              tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			Tokens: tks[:1],
		},
		&ExponentiationExpression{ // 14
			ExponentiationExpression: &ExponentiationExpression{
				UnaryExpression: UnaryExpression{
					UpdateExpression: UpdateExpression{
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: identA,
										Tokens:              tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			UnaryExpression: UnaryExpression{
				UpdateExpression: UpdateExpression{
					LeftHandSideExpression: &LeftHandSideExpression{
						NewExpression: &NewExpression{
							MemberExpression: MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									IdentifierReference: identB,
									Tokens:              tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&MultiplicativeExpression{ // 15
			MultiplicativeExpression: &MultiplicativeExpression{
				ExponentiationExpression: ExponentiationExpression{
					UnaryExpression: UnaryExpression{
						UpdateExpression: UpdateExpression{
							LeftHandSideExpression: &LeftHandSideExpression{
								NewExpression: &NewExpression{
									MemberExpression: MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: identA,
											Tokens:              tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			ExponentiationExpression: ExponentiationExpression{
				UnaryExpression: UnaryExpression{
					UpdateExpression: UpdateExpression{
						LeftHandSideExpression: &LeftHandSideExpression{
							NewExpression: &NewExpression{
								MemberExpression: MemberExpression{
									PrimaryExpression: &PrimaryExpression{
										IdentifierReference: identB,
										Tokens:              tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&AdditiveExpression{ // 16
			AdditiveExpression: &AdditiveExpression{
				MultiplicativeExpression: MultiplicativeExpression{
					ExponentiationExpression: ExponentiationExpression{
						UnaryExpression: UnaryExpression{
							UpdateExpression: UpdateExpression{
								LeftHandSideExpression: &LeftHandSideExpression{
									NewExpression: &NewExpression{
										MemberExpression: MemberExpression{
											PrimaryExpression: &PrimaryExpression{
												IdentifierReference: identA,
												Tokens:              tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			MultiplicativeExpression: MultiplicativeExpression{
				ExponentiationExpression: ExponentiationExpression{
					UnaryExpression: UnaryExpression{
						UpdateExpression: UpdateExpression{
							LeftHandSideExpression: &LeftHandSideExpression{
								NewExpression: &NewExpression{
									MemberExpression: MemberExpression{
										PrimaryExpression: &PrimaryExpression{
											IdentifierReference: identB,
											Tokens:              tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&ShiftExpression{ // 17
			ShiftExpression: &ShiftExpression{
				AdditiveExpression: AdditiveExpression{
					MultiplicativeExpression: MultiplicativeExpression{
						ExponentiationExpression: ExponentiationExpression{
							UnaryExpression: UnaryExpression{
								UpdateExpression: UpdateExpression{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: identA,
													Tokens:              tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			AdditiveExpression: AdditiveExpression{
				MultiplicativeExpression: MultiplicativeExpression{
					ExponentiationExpression: ExponentiationExpression{
						UnaryExpression: UnaryExpression{
							UpdateExpression: UpdateExpression{
								LeftHandSideExpression: &LeftHandSideExpression{
									NewExpression: &NewExpression{
										MemberExpression: MemberExpression{
											PrimaryExpression: &PrimaryExpression{
												IdentifierReference: identB,
												Tokens:              tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&RelationalExpression{ // 18
			RelationalExpression: &RelationalExpression{
				ShiftExpression: ShiftExpression{
					AdditiveExpression: AdditiveExpression{
						MultiplicativeExpression: MultiplicativeExpression{
							ExponentiationExpression: ExponentiationExpression{
								UnaryExpression: UnaryExpression{
									UpdateExpression: UpdateExpression{
										LeftHandSideExpression: &LeftHandSideExpression{
											NewExpression: &NewExpression{
												MemberExpression: MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: identA,
														Tokens:              tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			ShiftExpression: ShiftExpression{
				AdditiveExpression: AdditiveExpression{
					MultiplicativeExpression: MultiplicativeExpression{
						ExponentiationExpression: ExponentiationExpression{
							UnaryExpression: UnaryExpression{
								UpdateExpression: UpdateExpression{
									LeftHandSideExpression: &LeftHandSideExpression{
										NewExpression: &NewExpression{
											MemberExpression: MemberExpression{
												PrimaryExpression: &PrimaryExpression{
													IdentifierReference: identB,
													Tokens:              tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&EqualityExpression{ // 19
			EqualityExpression: &EqualityExpression{
				RelationalExpression: RelationalExpression{
					ShiftExpression: ShiftExpression{
						AdditiveExpression: AdditiveExpression{
							MultiplicativeExpression: MultiplicativeExpression{
								ExponentiationExpression: ExponentiationExpression{
									UnaryExpression: UnaryExpression{
										UpdateExpression: UpdateExpression{
											LeftHandSideExpression: &LeftHandSideExpression{
												NewExpression: &NewExpression{
													MemberExpression: MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: identA,
															Tokens:              tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			RelationalExpression: RelationalExpression{
				ShiftExpression: ShiftExpression{
					AdditiveExpression: AdditiveExpression{
						MultiplicativeExpression: MultiplicativeExpression{
							ExponentiationExpression: ExponentiationExpression{
								UnaryExpression: UnaryExpression{
									UpdateExpression: UpdateExpression{
										LeftHandSideExpression: &LeftHandSideExpression{
											NewExpression: &NewExpression{
												MemberExpression: MemberExpression{
													PrimaryExpression: &PrimaryExpression{
														IdentifierReference: identB,
														Tokens:              tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&BitwiseANDExpression{ // 20
			BitwiseANDExpression: &BitwiseANDExpression{
				EqualityExpression: EqualityExpression{
					RelationalExpression: RelationalExpression{
						ShiftExpression: ShiftExpression{
							AdditiveExpression: AdditiveExpression{
								MultiplicativeExpression: MultiplicativeExpression{
									ExponentiationExpression: ExponentiationExpression{
										UnaryExpression: UnaryExpression{
											UpdateExpression: UpdateExpression{
												LeftHandSideExpression: &LeftHandSideExpression{
													NewExpression: &NewExpression{
														MemberExpression: MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																IdentifierReference: identA,
																Tokens:              tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			EqualityExpression: EqualityExpression{
				RelationalExpression: RelationalExpression{
					ShiftExpression: ShiftExpression{
						AdditiveExpression: AdditiveExpression{
							MultiplicativeExpression: MultiplicativeExpression{
								ExponentiationExpression: ExponentiationExpression{
									UnaryExpression: UnaryExpression{
										UpdateExpression: UpdateExpression{
											LeftHandSideExpression: &LeftHandSideExpression{
												NewExpression: &NewExpression{
													MemberExpression: MemberExpression{
														PrimaryExpression: &PrimaryExpression{
															IdentifierReference: identB,
															Tokens:              tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&BitwiseXORExpression{ // 21
			BitwiseXORExpression: &BitwiseXORExpression{
				BitwiseANDExpression: BitwiseANDExpression{
					EqualityExpression: EqualityExpression{
						RelationalExpression: RelationalExpression{
							ShiftExpression: ShiftExpression{
								AdditiveExpression: AdditiveExpression{
									MultiplicativeExpression: MultiplicativeExpression{
										ExponentiationExpression: ExponentiationExpression{
											UnaryExpression: UnaryExpression{
												UpdateExpression: UpdateExpression{
													LeftHandSideExpression: &LeftHandSideExpression{
														NewExpression: &NewExpression{
															MemberExpression: MemberExpression{
																PrimaryExpression: &PrimaryExpression{
																	IdentifierReference: identA,
																	Tokens:              tks[:1],
																},
																Tokens: tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			BitwiseANDExpression: BitwiseANDExpression{
				EqualityExpression: EqualityExpression{
					RelationalExpression: RelationalExpression{
						ShiftExpression: ShiftExpression{
							AdditiveExpression: AdditiveExpression{
								MultiplicativeExpression: MultiplicativeExpression{
									ExponentiationExpression: ExponentiationExpression{
										UnaryExpression: UnaryExpression{
											UpdateExpression: UpdateExpression{
												LeftHandSideExpression: &LeftHandSideExpression{
													NewExpression: &NewExpression{
														MemberExpression: MemberExpression{
															PrimaryExpression: &PrimaryExpression{
																IdentifierReference: identB,
																Tokens:              tks[1:2],
															},
															Tokens: tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&BitwiseORExpression{ // 22
			BitwiseORExpression: &BitwiseORExpression{
				BitwiseXORExpression: BitwiseXORExpression{
					BitwiseANDExpression: BitwiseANDExpression{
						EqualityExpression: EqualityExpression{
							RelationalExpression: RelationalExpression{
								ShiftExpression: ShiftExpression{
									AdditiveExpression: AdditiveExpression{
										MultiplicativeExpression: MultiplicativeExpression{
											ExponentiationExpression: ExponentiationExpression{
												UnaryExpression: UnaryExpression{
													UpdateExpression: UpdateExpression{
														LeftHandSideExpression: &LeftHandSideExpression{
															NewExpression: &NewExpression{
																MemberExpression: MemberExpression{
																	PrimaryExpression: &PrimaryExpression{
																		IdentifierReference: identA,
																		Tokens:              tks[:1],
																	},
																	Tokens: tks[:1],
																},
																Tokens: tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			BitwiseXORExpression: BitwiseXORExpression{
				BitwiseANDExpression: BitwiseANDExpression{
					EqualityExpression: EqualityExpression{
						RelationalExpression: RelationalExpression{
							ShiftExpression: ShiftExpression{
								AdditiveExpression: AdditiveExpression{
									MultiplicativeExpression: MultiplicativeExpression{
										ExponentiationExpression: ExponentiationExpression{
											UnaryExpression: UnaryExpression{
												UpdateExpression: UpdateExpression{
													LeftHandSideExpression: &LeftHandSideExpression{
														NewExpression: &NewExpression{
															MemberExpression: MemberExpression{
																PrimaryExpression: &PrimaryExpression{
																	IdentifierReference: identB,
																	Tokens:              tks[1:2],
																},
																Tokens: tks[1:2],
															},
															Tokens: tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&LogicalANDExpression{ // 23
			LogicalANDExpression: &LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			IdentifierReference: identA,
																			Tokens:              tks[:1],
																		},
																		Tokens: tks[:1],
																	},
																	Tokens: tks[:1],
																},
																Tokens: tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			BitwiseORExpression: BitwiseORExpression{
				BitwiseXORExpression: BitwiseXORExpression{
					BitwiseANDExpression: BitwiseANDExpression{
						EqualityExpression: EqualityExpression{
							RelationalExpression: RelationalExpression{
								ShiftExpression: ShiftExpression{
									AdditiveExpression: AdditiveExpression{
										MultiplicativeExpression: MultiplicativeExpression{
											ExponentiationExpression: ExponentiationExpression{
												UnaryExpression: UnaryExpression{
													UpdateExpression: UpdateExpression{
														LeftHandSideExpression: &LeftHandSideExpression{
															NewExpression: &NewExpression{
																MemberExpression: MemberExpression{
																	PrimaryExpression: &PrimaryExpression{
																		IdentifierReference: identB,
																		Tokens:              tks[1:2],
																	},
																	Tokens: tks[1:2],
																},
																Tokens: tks[1:2],
															},
															Tokens: tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&LogicalORExpression{ // 24
			LogicalORExpression: &LogicalORExpression{
				LogicalANDExpression: LogicalANDExpression{
					BitwiseORExpression: BitwiseORExpression{
						BitwiseXORExpression: BitwiseXORExpression{
							BitwiseANDExpression: BitwiseANDExpression{
								EqualityExpression: EqualityExpression{
									RelationalExpression: RelationalExpression{
										ShiftExpression: ShiftExpression{
											AdditiveExpression: AdditiveExpression{
												MultiplicativeExpression: MultiplicativeExpression{
													ExponentiationExpression: ExponentiationExpression{
														UnaryExpression: UnaryExpression{
															UpdateExpression: UpdateExpression{
																LeftHandSideExpression: &LeftHandSideExpression{
																	NewExpression: &NewExpression{
																		MemberExpression: MemberExpression{
																			PrimaryExpression: &PrimaryExpression{
																				IdentifierReference: identA,
																				Tokens:              tks[:1],
																			},
																			Tokens: tks[:1],
																		},
																		Tokens: tks[:1],
																	},
																	Tokens: tks[:1],
																},
																Tokens: tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				Tokens: tks[:1],
			},
			LogicalANDExpression: LogicalANDExpression{
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			IdentifierReference: identB,
																			Tokens:              tks[1:2],
																		},
																		Tokens: tks[1:2],
																	},
																	Tokens: tks[1:2],
																},
																Tokens: tks[1:2],
															},
															Tokens: tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
		&ConditionalExpression{ // 25
			CoalesceExpression: &CoalesceExpression{
				CoalesceExpressionHead: &CoalesceExpression{
					BitwiseORExpression: BitwiseORExpression{
						BitwiseXORExpression: BitwiseXORExpression{
							BitwiseANDExpression: BitwiseANDExpression{
								EqualityExpression: EqualityExpression{
									RelationalExpression: RelationalExpression{
										ShiftExpression: ShiftExpression{
											AdditiveExpression: AdditiveExpression{
												MultiplicativeExpression: MultiplicativeExpression{
													ExponentiationExpression: ExponentiationExpression{
														UnaryExpression: UnaryExpression{
															UpdateExpression: UpdateExpression{
																LeftHandSideExpression: &LeftHandSideExpression{
																	NewExpression: &NewExpression{
																		MemberExpression: MemberExpression{
																			PrimaryExpression: &PrimaryExpression{
																				IdentifierReference: identA,
																				Tokens:              tks[:1],
																			},
																			Tokens: tks[:1],
																		},
																		Tokens: tks[:1],
																	},
																	Tokens: tks[:1],
																},
																Tokens: tks[:1],
															},
															Tokens: tks[:1],
														},
														Tokens: tks[:1],
													},
													Tokens: tks[:1],
												},
												Tokens: tks[:1],
											},
											Tokens: tks[:1],
										},
										Tokens: tks[:1],
									},
									Tokens: tks[:1],
								},
								Tokens: tks[:1],
							},
							Tokens: tks[:1],
						},
						Tokens: tks[:1],
					},
					Tokens: tks[:1],
				},
				BitwiseORExpression: BitwiseORExpression{
					BitwiseXORExpression: BitwiseXORExpression{
						BitwiseANDExpression: BitwiseANDExpression{
							EqualityExpression: EqualityExpression{
								RelationalExpression: RelationalExpression{
									ShiftExpression: ShiftExpression{
										AdditiveExpression: AdditiveExpression{
											MultiplicativeExpression: MultiplicativeExpression{
												ExponentiationExpression: ExponentiationExpression{
													UnaryExpression: UnaryExpression{
														UpdateExpression: UpdateExpression{
															LeftHandSideExpression: &LeftHandSideExpression{
																NewExpression: &NewExpression{
																	MemberExpression: MemberExpression{
																		PrimaryExpression: &PrimaryExpression{
																			IdentifierReference: identB,
																			Tokens:              tks[1:2],
																		},
																		Tokens: tks[1:2],
																	},
																	Tokens: tks[1:2],
																},
																Tokens: tks[1:2],
															},
															Tokens: tks[1:2],
														},
														Tokens: tks[1:2],
													},
													Tokens: tks[1:2],
												},
												Tokens: tks[1:2],
											},
											Tokens: tks[1:2],
										},
										Tokens: tks[1:2],
									},
									Tokens: tks[1:2],
								},
								Tokens: tks[1:2],
							},
							Tokens: tks[1:2],
						},
						Tokens: tks[1:2],
					},
					Tokens: tks[1:2],
				},
				Tokens: tks[1:2],
			},
			Tokens: tks[:2],
		},
	} {
		if output := UnwrapConditional(WrapConditional(test)); !reflect.DeepEqual(output, test) {
			t.Errorf("test %d: expecting\n%v\n...got...\n%v", n+1, test, output)
		}
	}
}
