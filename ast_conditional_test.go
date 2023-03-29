package javascript

import "testing"

func makeConditionLiteral(tk Tokens, pos int) ConditionalExpression {
	p := PrimaryExpression{
		Tokens: tk[pos : pos+1],
	}
	if tk[pos].Type == TokenIdentifier || tk[pos].Type == TokenKeyword {
		p.IdentifierReference = &tk[pos]
	} else {
		p.Literal = &tk[pos]
	}
	return *WrapConditional(&p)
}

func wrapConditional(i ConditionalWrappable) ConditionalExpression {
	return *WrapConditional(i)
}

func TestConditional(t *testing.T) {
	doTests(t, []sourceFn{
		{`true`, func(t *test, tk Tokens) { // 1
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`true || false`, func(t *test, tk Tokens) { // 2
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalORExpression:  litA.LogicalORExpression,
				LogicalANDExpression: litB.LogicalORExpression.LogicalANDExpression,
				Tokens:               tk[:5],
			})
		}},
		{`true && false`, func(t *test, tk Tokens) { // 3
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalANDExpression: LogicalANDExpression{
					LogicalANDExpression: &litA.LogicalORExpression.LogicalANDExpression,
					BitwiseORExpression:  litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					Tokens:               tk[:5],
				},
				Tokens: tk[:5],
			})
		}},
		{`1 || 2 || 3`, func(t *test, tk Tokens) { // 4
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalORExpression: &LogicalORExpression{
					LogicalORExpression:  litA.LogicalORExpression,
					LogicalANDExpression: litB.LogicalORExpression.LogicalANDExpression,
					Tokens:               tk[:5],
				},
				LogicalANDExpression: litC.LogicalORExpression.LogicalANDExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 && 2 && 3`, func(t *test, tk Tokens) { // 5
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(LogicalANDExpression{
				LogicalANDExpression: &LogicalANDExpression{
					LogicalANDExpression: &litA.LogicalORExpression.LogicalANDExpression,
					BitwiseORExpression:  litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					Tokens:               tk[:5],
				},
				BitwiseORExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
				Tokens:              tk[:9],
			})
		}},
		{`1 && 2 || 3`, func(t *test, tk Tokens) { // 6
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalORExpression: &LogicalORExpression{
					LogicalANDExpression: LogicalANDExpression{
						LogicalANDExpression: &litA.LogicalORExpression.LogicalANDExpression,
						BitwiseORExpression:  litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
						Tokens:               tk[:5],
					},
					Tokens: tk[:5],
				},
				LogicalANDExpression: litC.LogicalORExpression.LogicalANDExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 || 2 && 3`, func(t *test, tk Tokens) { // 7
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalORExpression: litA.LogicalORExpression,
				LogicalANDExpression: LogicalANDExpression{
					LogicalANDExpression: &litB.LogicalORExpression.LogicalANDExpression,
					BitwiseORExpression:  litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					Tokens:               tk[4:9],
				},
				Tokens: tk[:9],
			})
		}},
		{`1 | 2`, func(t *test, tk Tokens) { // 8
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseORExpression{
				BitwiseORExpression:  &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
				BitwiseXORExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 | 2 | 3`, func(t *test, tk Tokens) { // 9
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(BitwiseORExpression{
				BitwiseORExpression: &BitwiseORExpression{
					BitwiseORExpression:  &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					BitwiseXORExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
					Tokens:               tk[:5],
				},
				BitwiseXORExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 ^ 2`, func(t *test, tk Tokens) { // 10
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseXORExpression{
				BitwiseXORExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
				BitwiseANDExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 ^ 2 ^ 3`, func(t *test, tk Tokens) { // 11
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(BitwiseXORExpression{
				BitwiseXORExpression: &BitwiseXORExpression{
					BitwiseXORExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
					BitwiseANDExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
					Tokens:               tk[:5],
				},
				BitwiseANDExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 & 2`, func(t *test, tk Tokens) { // 12
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseANDExpression{
				BitwiseANDExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
				EqualityExpression:   litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 & 2 & 3`, func(t *test, tk Tokens) { // 13
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(BitwiseANDExpression{
				BitwiseANDExpression: &BitwiseANDExpression{
					BitwiseANDExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
					EqualityExpression:   litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
					Tokens:               tk[:5],
				},
				EqualityExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				Tokens:             tk[:9],
			})
		}},
		{`1 == 2`, func(t *test, tk Tokens) { // 14
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 == 2 != true`, func(t *test, tk Tokens) { // 15
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &EqualityExpression{
					EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
					EqualityOperator:     EqualityEqual,
					RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
					Tokens:               tk[:5],
				},
				EqualityOperator:     EqualityNotEqual,
				RelationalExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 != 2`, func(t *test, tk Tokens) { // 16
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityNotEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 != 2 == true`, func(t *test, tk Tokens) { // 17
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &EqualityExpression{
					EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
					EqualityOperator:     EqualityNotEqual,
					RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
					Tokens:               tk[:5],
				},
				EqualityOperator:     EqualityEqual,
				RelationalExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 === 2`, func(t *test, tk Tokens) { // 18
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityStrictEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 === 2 !== true`, func(t *test, tk Tokens) { // 19
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &EqualityExpression{
					EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
					EqualityOperator:     EqualityStrictEqual,
					RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
					Tokens:               tk[:5],
				},
				EqualityOperator:     EqualityStrictNotEqual,
				RelationalExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 !== 2`, func(t *test, tk Tokens) { // 20
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityStrictNotEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 !== 2 === true`, func(t *test, tk Tokens) { // 21
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &EqualityExpression{
					EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
					EqualityOperator:     EqualityStrictNotEqual,
					RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
					Tokens:               tk[:5],
				},
				EqualityOperator:     EqualityStrictEqual,
				RelationalExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:9],
			})
		}},
		{`1 < 2`, func(t *test, tk Tokens) { // 22
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipLessThan,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 < 2 === true`, func(t *test, tk Tokens) { // 23
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &EqualityExpression{
					RelationalExpression: RelationalExpression{
						RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
						RelationshipOperator: RelationshipLessThan,
						ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
						Tokens:               tk[:5],
					},
					Tokens: tk[:5],
				},
				EqualityOperator:     EqualityStrictEqual,
				RelationalExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:9],
			})
		}},
		{`true === 1 < 2`, func(t *test, tk Tokens) { // 24
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:   EqualityStrictEqual,
				RelationalExpression: RelationalExpression{
					RelationalExpression: &litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
					RelationshipOperator: RelationshipLessThan,
					ShiftExpression:      litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					Tokens:               tk[4:9],
				},
				Tokens: tk[:9],
			})
		}},
		{`1 > 2`, func(t *test, tk Tokens) { // 25
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipGreaterThan,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 <= 2`, func(t *test, tk Tokens) { // 26
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipLessThanEqual,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 >= 2`, func(t *test, tk Tokens) { // 27
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 5)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipGreaterThanEqual,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:6],
			})
		}},
		{`1 instanceof 2`, func(t *test, tk Tokens) { // 28
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipInstanceOf,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 in 2`, func(t *test, tk Tokens) { // 29
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`1 in 2`, func(t *test, tk Tokens) { // 30
			t.In = true
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipIn,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 << 2`, func(t *test, tk Tokens) { // 31
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftLeft,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:5],
			})
		}},
		{`1 << 2 << 3`, func(t *test, tk Tokens) { // 32
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression: &ShiftExpression{
					ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					ShiftOperator:      ShiftLeft,
					AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					Tokens:             tk[:5],
				},
				ShiftOperator:      ShiftLeft,
				AdditiveExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:9],
			})
		}},
		{`1 >> 2`, func(t *test, tk Tokens) { // 33
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 5)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftRight,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:6],
			})
		}},
		{`1 >> 2 >> 3`, func(t *test, tk Tokens) { // 34
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 5)
			litC := makeConditionLiteral(tk, 10)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression: &ShiftExpression{
					ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					ShiftOperator:      ShiftRight,
					AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					Tokens:             tk[:6],
				},
				ShiftOperator:      ShiftRight,
				AdditiveExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:11],
			})
		}},
		{`1 >>> 2`, func(t *test, tk Tokens) { // 35
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 6)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftUnsignedRight,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:7],
			})
		}},
		{`1 >>> 2 >>> 3`, func(t *test, tk Tokens) { // 36
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 6)
			litC := makeConditionLiteral(tk, 12)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression: &ShiftExpression{
					ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					ShiftOperator:      ShiftUnsignedRight,
					AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					Tokens:             tk[:7],
				},
				ShiftOperator:      ShiftUnsignedRight,
				AdditiveExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:13],
			})
		}},
		{`1 + 2`, func(t *test, tk Tokens) { // 37
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 + 2 + 3`, func(t *test, tk Tokens) { // 38
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression: &AdditiveExpression{
					AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					AdditiveOperator:         AdditiveAdd,
					MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					Tokens:                   tk[:5],
				},
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 - 2`, func(t *test, tk Tokens) { // 39
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveMinus,
				MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 + 2 - 3`, func(t *test, tk Tokens) { // 40
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression: &AdditiveExpression{
					AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					AdditiveOperator:         AdditiveAdd,
					MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					Tokens:                   tk[:5],
				},
				AdditiveOperator:         AdditiveMinus,
				MultiplicativeExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 - 2 - 3`, func(t *test, tk Tokens) { // 41
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)

			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression: &AdditiveExpression{
					AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					AdditiveOperator:         AdditiveMinus,
					MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					Tokens:                   tk[:5],
				},
				AdditiveOperator:         AdditiveMinus,
				MultiplicativeExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 - 2 + 3`, func(t *test, tk Tokens) { // 42
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)

			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression: &AdditiveExpression{
					AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					AdditiveOperator:         AdditiveMinus,
					MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					Tokens:                   tk[:5],
				},
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 * 2`, func(t *test, tk Tokens) { // 43
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeMultiply,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 * 2 * 3`, func(t *test, tk Tokens) { // 44
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &MultiplicativeExpression{
					MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					MultiplicativeOperator:   MultiplicativeMultiply,
					ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
					Tokens:                   tk[:5],
				},
				MultiplicativeOperator:   MultiplicativeMultiply,
				ExponentiationExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 / 2`, func(t *test, tk Tokens) { // 45
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeDivide,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 / 2 / 3`, func(t *test, tk Tokens) { // 46
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &MultiplicativeExpression{
					MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					MultiplicativeOperator:   MultiplicativeDivide,
					ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
					Tokens:                   tk[:5],
				},
				MultiplicativeOperator:   MultiplicativeDivide,
				ExponentiationExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 % 2`, func(t *test, tk Tokens) { // 47
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeRemainder,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 % 2 % 3`, func(t *test, tk Tokens) { // 48
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &MultiplicativeExpression{
					MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
					MultiplicativeOperator:   MultiplicativeRemainder,
					ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
					Tokens:                   tk[:5],
				},
				MultiplicativeOperator:   MultiplicativeRemainder,
				ExponentiationExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:9],
			})
		}},
		{`1 ** 2`, func(t *test, tk Tokens) { // 49
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ExponentiationExpression{
				ExponentiationExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				UnaryExpression:          litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 ** 2 ** 3`, func(t *test, tk Tokens) { // 50
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(ExponentiationExpression{
				ExponentiationExpression: &ExponentiationExpression{
					ExponentiationExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
					UnaryExpression:          litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
					Tokens:                   tk[:5],
				},
				UnaryExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[:9],
			})
		}},
		{`delete 1`, func(t *test, tk Tokens) { // 51
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryDelete},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`void 1`, func(t *test, tk Tokens) { // 52
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryVoid},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`typeof 1`, func(t *test, tk Tokens) { // 53
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryTypeOf},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`+1`, func(t *test, tk Tokens) { // 54
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAdd},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`-1`, func(t *test, tk Tokens) { // 55
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryMinus},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`~1`, func(t *test, tk Tokens) { // 56
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryBitwiseNot},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`!1`, func(t *test, tk Tokens) { // 57
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryLogicalNot},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`await 1`, func(t *test, tk Tokens) { // 58
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`await 1`, func(t *test, tk Tokens) { // 59
			t.Await = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAwait},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`await!~-1`, func(t *test, tk Tokens) { // 60
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`await!~-1`, func(t *test, tk Tokens) { // 61
			t.Await = true
			litA := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAwait, UnaryLogicalNot, UnaryBitwiseNot, UnaryMinus},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:5],
			})
		}},
		{`a++`, func(t *test, tk Tokens) { // 62
			litA := makeConditionLiteral(tk, 0)
			t.Output = wrapConditional(UpdateExpression{
				LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
				UpdateOperator:         UpdatePostIncrement,
				Tokens:                 tk[:2],
			})
		}},
		{`a--`, func(t *test, tk Tokens) { // 63
			litA := makeConditionLiteral(tk, 0)
			t.Output = wrapConditional(UpdateExpression{
				LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
				UpdateOperator:         UpdatePostDecrement,
				Tokens:                 tk[:2],
			})
		}},
		{`++a`, func(t *test, tk Tokens) { // 64
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator:  UpdatePreIncrement,
				UnaryExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[:2],
			})
		}},
		{`--a`, func(t *test, tk Tokens) { // 65
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator:  UpdatePreDecrement,
				UnaryExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[:2],
			})
		}},
		{`++!a`, func(t *test, tk Tokens) { // 66
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator: UpdatePreIncrement,
				UnaryExpression: &UnaryExpression{
					UnaryOperators:   []UnaryOperator{UnaryLogicalNot},
					UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
					Tokens:           tk[1:3],
				},
				Tokens: tk[:3],
			})
		}},
		{`--!a`, func(t *test, tk Tokens) { // 67
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator: UpdatePreDecrement,
				UnaryExpression: &UnaryExpression{
					UnaryOperators:   []UnaryOperator{UnaryLogicalNot},
					UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
					Tokens:           tk[1:3],
				},
				Tokens: tk[:3],
			})
		}},
		{`true ? 1 : 2`, func(t *test, tk Tokens) { // 68
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = ConditionalExpression{
				LogicalORExpression: litA.LogicalORExpression,
				True: &AssignmentExpression{
					ConditionalExpression: &litB,
					Tokens:                tk[4:5],
				},
				False: &AssignmentExpression{
					ConditionalExpression: &litC,
					Tokens:                tk[8:9],
				},
				Tokens: tk[:9],
			}
		}},
		{`true ? 1 ? 2 : 3 : 4 ? 5 : 6`, func(t *test, tk Tokens) { // 69
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			litD := makeConditionLiteral(tk, 12)
			litE := makeConditionLiteral(tk, 16)
			litF := makeConditionLiteral(tk, 20)
			litG := makeConditionLiteral(tk, 24)
			t.Output = ConditionalExpression{
				LogicalORExpression: litA.LogicalORExpression,
				True: &AssignmentExpression{
					ConditionalExpression: &ConditionalExpression{
						LogicalORExpression: litB.LogicalORExpression,
						True: &AssignmentExpression{
							ConditionalExpression: &litC,
							Tokens:                tk[8:9],
						},
						False: &AssignmentExpression{
							ConditionalExpression: &litD,
							Tokens:                tk[12:13],
						},
						Tokens: tk[4:13],
					},
					Tokens: tk[4:13],
				},
				False: &AssignmentExpression{
					ConditionalExpression: &ConditionalExpression{
						LogicalORExpression: litE.LogicalORExpression,
						True: &AssignmentExpression{
							ConditionalExpression: &litF,
							Tokens:                tk[20:21],
						},
						False: &AssignmentExpression{
							ConditionalExpression: &litG,
							Tokens:                tk[24:25],
						},
						Tokens: tk[16:25],
					},
					Tokens: tk[16:25],
				},
				Tokens: tk[:25],
			}
		}},
		{`true ? :`, func(t *test, tk Tokens) { // 70
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "ConditionalExpression",
				Token:   tk[4],
			}
		}},
		{`true ? a`, func(t *test, tk Tokens) { // 71
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "ConditionalExpression",
				Token:   tk[5],
			}
		}},
		{`true ? a : :`, func(t *test, tk Tokens) { // 72
			t.Err = Error{
				Err:     assignmentError(tk[8]),
				Parsing: "ConditionalExpression",
				Token:   tk[8],
			}
		}},
		{`++?`, func(t *test, tk Tokens) { /// 73
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err: Error{
																					Err: Error{
																						Err: Error{
																							Err:     ErrNoIdentifier,
																							Parsing: "PrimaryExpression",
																							Token:   tk[1],
																						},
																						Parsing: "MemberExpression",
																						Token:   tk[1],
																					},
																					Parsing: "NewExpression",
																					Token:   tk[1],
																				},
																				Parsing: "LeftHandSideExpression",
																				Token:   tk[1],
																			},
																			Parsing: "UpdateExpression",
																			Token:   tk[1],
																		},
																		Parsing: "UnaryExpression",
																		Token:   tk[1],
																	},
																	Parsing: "UpdateExpression",
																	Token:   tk[1],
																},
																Parsing: "UnaryExpression",
																Token:   tk[0],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[0],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[0],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[0],
												},
												Parsing: "ShiftExpression",
												Token:   tk[0],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"a\n??\nb", func(t *test, tk Tokens) { // 74
			t.Output = ConditionalExpression{
				CoalesceExpression: &CoalesceExpression{
					CoalesceExpressionHead: &CoalesceExpression{
						BitwiseORExpression: makeConditionLiteral(tk, 0).LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
						Tokens:              tk[:1],
					},
					BitwiseORExpression: makeConditionLiteral(tk, 4).LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					Tokens:              tk[:5],
				},
				Tokens: tk[:5],
			}
		}},
		{"a\n??\nb\n??\nc", func(t *test, tk Tokens) { // 75
			t.Output = ConditionalExpression{
				CoalesceExpression: &CoalesceExpression{
					CoalesceExpressionHead: &CoalesceExpression{
						CoalesceExpressionHead: &CoalesceExpression{
							BitwiseORExpression: makeConditionLiteral(tk, 0).LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
							Tokens:              tk[:1],
						},
						BitwiseORExpression: makeConditionLiteral(tk, 4).LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
						Tokens:              tk[:5],
					},
					BitwiseORExpression: makeConditionLiteral(tk, 8).LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
					Tokens:              tk[:9],
				},
				Tokens: tk[:9],
			}
		}},
		{"a\n??\n!", func(t *test, tk Tokens) { // 76
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err:     ErrNoIdentifier,
																				Parsing: "PrimaryExpression",
																				Token:   tk[5],
																			},
																			Parsing: "MemberExpression",
																			Token:   tk[5],
																		},
																		Parsing: "NewExpression",
																		Token:   tk[5],
																	},
																	Parsing: "LeftHandSideExpression",
																	Token:   tk[5],
																},
																Parsing: "UpdateExpression",
																Token:   tk[5],
															},
															Parsing: "UnaryExpression",
															Token:   tk[5],
														},
														Parsing: "ExponentiationExpression",
														Token:   tk[4],
													},
													Parsing: "MultiplicativeExpression",
													Token:   tk[4],
												},
												Parsing: "AdditiveExpression",
												Token:   tk[4],
											},
											Parsing: "ShiftExpression",
											Token:   tk[4],
										},
										Parsing: "RelationalExpression",
										Token:   tk[4],
									},
									Parsing: "EqualityExpression",
									Token:   tk[4],
								},
								Parsing: "BitwiseANDExpression",
								Token:   tk[4],
							},
							Parsing: "BitwiseXORExpression",
							Token:   tk[4],
						},
						Parsing: "BitwiseORExpression",
						Token:   tk[4],
					},
					Parsing: "CoalesceExpression",
					Token:   tk[4],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[1],
			}
		}},
		{"1++", func(t *test, tk Tokens) { // 77
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err:     ErrNotSimple,
																	Parsing: "UpdateExpression",
																	Token:   tk[1],
																},
																Parsing: "UnaryExpression",
																Token:   tk[0],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[0],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[0],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[0],
												},
												Parsing: "ShiftExpression",
												Token:   tk[0],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"1--", func(t *test, tk Tokens) { // 78
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err:     ErrNotSimple,
																	Parsing: "UpdateExpression",
																	Token:   tk[1],
																},
																Parsing: "UnaryExpression",
																Token:   tk[0],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[0],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[0],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[0],
												},
												Parsing: "ShiftExpression",
												Token:   tk[0],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"#a in b", func(t *test, tk Tokens) { // 79
			t.In = true
			t.Output = *WrapConditional(RelationalExpression{
				PrivateIdentifier:    &tk[0],
				RelationshipOperator: RelationshipIn,
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
														IdentifierReference: &tk[4],
														Tokens:              tk[4:5],
													},
													Tokens: tk[4:5],
												},
												Tokens: tk[4:5],
											},
											Tokens: tk[4:5],
										},
										Tokens: tk[4:5],
									},
									Tokens: tk[4:5],
								},
								Tokens: tk[4:5],
							},
							Tokens: tk[4:5],
						},
						Tokens: tk[4:5],
					},
					Tokens: tk[4:5],
				},
				Tokens: tk[:5],
			})
		}},
		{"#a in b", func(t *test, tk Tokens) { // 80
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err: Error{
																					Err:     ErrNoIdentifier,
																					Parsing: "PrimaryExpression",
																					Token:   tk[0],
																				},
																				Parsing: "MemberExpression",
																				Token:   tk[0],
																			},
																			Parsing: "NewExpression",
																			Token:   tk[0],
																		},
																		Parsing: "LeftHandSideExpression",
																		Token:   tk[0],
																	},
																	Parsing: "UpdateExpression",
																	Token:   tk[0],
																},
																Parsing: "UnaryExpression",
																Token:   tk[0],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[0],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[0],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[0],
												},
												Parsing: "ShiftExpression",
												Token:   tk[0],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"#a of b", func(t *test, tk Tokens) { // 81
			t.In = true
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err: Error{
																					Err:     ErrNoIdentifier,
																					Parsing: "PrimaryExpression",
																					Token:   tk[0],
																				},
																				Parsing: "MemberExpression",
																				Token:   tk[0],
																			},
																			Parsing: "NewExpression",
																			Token:   tk[0],
																		},
																		Parsing: "LeftHandSideExpression",
																		Token:   tk[0],
																	},
																	Parsing: "UpdateExpression",
																	Token:   tk[0],
																},
																Parsing: "UnaryExpression",
																Token:   tk[0],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[0],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[0],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[0],
												},
												Parsing: "ShiftExpression",
												Token:   tk[0],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"#a in #b", func(t *test, tk Tokens) { // 82
			t.In = true
			t.Err = Error{
				Err: Error{
					Err: Error{
						Err: Error{
							Err: Error{
								Err: Error{
									Err: Error{
										Err: Error{
											Err: Error{
												Err: Error{
													Err: Error{
														Err: Error{
															Err: Error{
																Err: Error{
																	Err: Error{
																		Err: Error{
																			Err: Error{
																				Err: Error{
																					Err:     ErrNoIdentifier,
																					Parsing: "PrimaryExpression",
																					Token:   tk[4],
																				},
																				Parsing: "MemberExpression",
																				Token:   tk[4],
																			},
																			Parsing: "NewExpression",
																			Token:   tk[4],
																		},
																		Parsing: "LeftHandSideExpression",
																		Token:   tk[4],
																	},
																	Parsing: "UpdateExpression",
																	Token:   tk[4],
																},
																Parsing: "UnaryExpression",
																Token:   tk[4],
															},
															Parsing: "ExponentiationExpression",
															Token:   tk[4],
														},
														Parsing: "MultiplicativeExpression",
														Token:   tk[4],
													},
													Parsing: "AdditiveExpression",
													Token:   tk[4],
												},
												Parsing: "ShiftExpression",
												Token:   tk[4],
											},
											Parsing: "RelationalExpression",
											Token:   tk[0],
										},
										Parsing: "EqualityExpression",
										Token:   tk[0],
									},
									Parsing: "BitwiseANDExpression",
									Token:   tk[0],
								},
								Parsing: "BitwiseXORExpression",
								Token:   tk[0],
							},
							Parsing: "BitwiseORExpression",
							Token:   tk[0],
						},
						Parsing: "LogicalANDExpression",
						Token:   tk[0],
					},
					Parsing: "LogicalORExpression",
					Token:   tk[0],
				},
				Parsing: "ConditionalExpression",
				Token:   tk[0],
			}
		}},
		{"this.#a", func(t *test, tk Tokens) { // 83
			t.Output = *WrapConditional(MemberExpression{
				MemberExpression: &MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						This:   &tk[0],
						Tokens: tk[:1],
					},
					Tokens: tk[:1],
				},
				PrivateIdentifier: &tk[2],
				Tokens:            tk[:3],
			})
		}},
		{"this.#a++", func(t *test, tk Tokens) { // 84
			t.Output = *WrapConditional(UpdateExpression{
				LeftHandSideExpression: &LeftHandSideExpression{
					NewExpression: &NewExpression{
						MemberExpression: MemberExpression{
							MemberExpression: &MemberExpression{
								PrimaryExpression: &PrimaryExpression{
									This:   &tk[0],
									Tokens: tk[:1],
								},
								Tokens: tk[:1],
							},
							PrivateIdentifier: &tk[2],
							Tokens:            tk[:3],
						},
						Tokens: tk[:3],
					},
					Tokens: tk[:3],
				},
				UpdateOperator: UpdatePostIncrement,
				Tokens:         tk[:4],
			})
		}},
	}, func(t *test) (Type, error) {
		var ce ConditionalExpression
		err := ce.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return ce, err
	})
}
