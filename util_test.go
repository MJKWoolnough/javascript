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
