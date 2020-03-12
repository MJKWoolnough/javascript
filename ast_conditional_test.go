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

func wrapConditional(i interface{}) ConditionalExpression {
	return *WrapConditional(i)
}

func TestConditional(t *testing.T) {
	doTests(t, []sourceFn{
		{`true`, func(t *test, tk Tokens) {
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`true || false`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(LogicalORExpression{
				LogicalORExpression:  litA.LogicalORExpression,
				LogicalANDExpression: litB.LogicalORExpression.LogicalANDExpression,
				Tokens:               tk[:5],
			})
		}},
		{`true && false`, func(t *test, tk Tokens) {
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
		{`1 || 2 || 3`, func(t *test, tk Tokens) {
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
		{`1 && 2 && 3`, func(t *test, tk Tokens) {
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
		{`1 && 2 || 3`, func(t *test, tk Tokens) {
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
		{`1 || 2 && 3`, func(t *test, tk Tokens) {
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
		{`1 | 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseORExpression{
				BitwiseORExpression:  &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression,
				BitwiseXORExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 | 2 | 3`, func(t *test, tk Tokens) {
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
		{`1 ^ 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseXORExpression{
				BitwiseXORExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression,
				BitwiseANDExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 ^ 2 ^ 3`, func(t *test, tk Tokens) {
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
		{`1 & 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(BitwiseANDExpression{
				BitwiseANDExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression,
				EqualityExpression:   litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 & 2 & 3`, func(t *test, tk Tokens) {
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
		{`1 == 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 == 2 != true`, func(t *test, tk Tokens) {
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
		{`1 != 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityNotEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 != 2 == true`, func(t *test, tk Tokens) {
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
		{`1 === 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityStrictEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 === 2 !== true`, func(t *test, tk Tokens) {
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
		{`1 !== 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(EqualityExpression{
				EqualityExpression:   &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression,
				EqualityOperator:     EqualityStrictNotEqual,
				RelationalExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				Tokens:               tk[:5],
			})

		}},
		{`1 !== 2 === true`, func(t *test, tk Tokens) {
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
		{`1 < 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipLessThan,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 < 2 === true`, func(t *test, tk Tokens) {
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
		{`true === 1 < 2`, func(t *test, tk Tokens) {
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
		{`1 > 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipGreaterThan,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})

		}},
		{`1 <= 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipLessThanEqual,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})

		}},
		{`1 >= 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipGreaterThanEqual,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 instanceof 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(RelationalExpression{
				RelationalExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression,
				RelationshipOperator: RelationshipInstanceOf,
				ShiftExpression:      litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				Tokens:               tk[:5],
			})
		}},
		{`1 in 2`, func(t *test, tk Tokens) {
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`1 in 2`, func(t *test, tk Tokens) {
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
		{`1 << 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftLeft,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:5],
			})
		}},
		{`1 << 2 << 3`, func(t *test, tk Tokens) {
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
		{`1 >> 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftRight,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:5],
			})
		}},
		{`1 >> 2 >> 3`, func(t *test, tk Tokens) {

			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression: &ShiftExpression{
					ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					ShiftOperator:      ShiftRight,
					AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					Tokens:             tk[:5],
				},
				ShiftOperator:      ShiftRight,
				AdditiveExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:9],
			})
		}},
		{`1 >>> 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
				ShiftOperator:      ShiftUnsignedRight,
				AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:5],
			})
		}},
		{`1 >>> 2 >>> 3`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			litC := makeConditionLiteral(tk, 8)
			t.Output = wrapConditional(ShiftExpression{
				ShiftExpression: &ShiftExpression{
					ShiftExpression:    &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression,
					ShiftOperator:      ShiftUnsignedRight,
					AdditiveExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
					Tokens:             tk[:5],
				},
				ShiftOperator:      ShiftUnsignedRight,
				AdditiveExpression: litC.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				Tokens:             tk[:9],
			})
		}},
		{`1 + 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveAdd,
				MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 + 2 + 3`, func(t *test, tk Tokens) {
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
		{`1 - 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(AdditiveExpression{
				AdditiveExpression:       &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression,
				AdditiveOperator:         AdditiveMinus,
				MultiplicativeExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 + 2 - 3`, func(t *test, tk Tokens) {
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
		{`1 - 2 - 3`, func(t *test, tk Tokens) {
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
		{`1 - 2 + 3`, func(t *test, tk Tokens) {
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
		{`1 * 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeMultiply,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 * 2 * 3`, func(t *test, tk Tokens) {
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
		{`1 / 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeDivide,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 / 2 / 3`, func(t *test, tk Tokens) {
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
		{`1 % 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(MultiplicativeExpression{
				MultiplicativeExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression,
				MultiplicativeOperator:   MultiplicativeRemainder,
				ExponentiationExpression: litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 % 2 % 3`, func(t *test, tk Tokens) {
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
		{`1 ** 2`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			litB := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(ExponentiationExpression{
				ExponentiationExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression,
				UnaryExpression:          litB.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:                   tk[:5],
			})
		}},
		{`1 ** 2 ** 3`, func(t *test, tk Tokens) {
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
		{`delete 1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryDelete},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`void 1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryVoid},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`typeof 1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryTypeOf},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`+1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAdd},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`-1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryMinus},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`~1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryBitwiseNot},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`!1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryLogicalNot},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:2],
			})
		}},
		{`await 1`, func(t *test, tk Tokens) {
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`await 1`, func(t *test, tk Tokens) {
			t.Await = true
			litA := makeConditionLiteral(tk, 2)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAwait},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:3],
			})
		}},
		{`await!~-1`, func(t *test, tk Tokens) {
			t.Output = makeConditionLiteral(tk, 0)
		}},
		{`await!~-1`, func(t *test, tk Tokens) {
			t.Await = true
			litA := makeConditionLiteral(tk, 4)
			t.Output = wrapConditional(UnaryExpression{
				UnaryOperators:   []UnaryOperator{UnaryAwait, UnaryLogicalNot, UnaryBitwiseNot, UnaryMinus},
				UpdateExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression,
				Tokens:           tk[:5],
			})
		}},
		{`1++`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			t.Output = wrapConditional(UpdateExpression{
				LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
				UpdateOperator:         UpdatePostIncrement,
				Tokens:                 tk[:2],
			})
		}},
		{`1--`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 0)
			t.Output = wrapConditional(UpdateExpression{
				LeftHandSideExpression: litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.LeftHandSideExpression,
				UpdateOperator:         UpdatePostDecrement,
				Tokens:                 tk[:2],
			})
		}},
		{`++1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator:  UpdatePreIncrement,
				UnaryExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[:2],
			})
		}},
		{`--1`, func(t *test, tk Tokens) {
			litA := makeConditionLiteral(tk, 1)
			t.Output = wrapConditional(UpdateExpression{
				UpdateOperator:  UpdatePreDecrement,
				UnaryExpression: &litA.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression,
				Tokens:          tk[:2],
			})
		}},
		{`++!1`, func(t *test, tk Tokens) {
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
		{`--!1`, func(t *test, tk Tokens) {
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
		{`true ? 1 : 2`, func(t *test, tk Tokens) {
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
		{`true ? 1 ? 2 : 3 : 4 ? 5 : 6`, func(t *test, tk Tokens) {
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
		{`true ? :`, func(t *test, tk Tokens) {
			t.Err = Error{
				Err:     assignmentError(tk[4]),
				Parsing: "ConditionalExpression",
				Token:   tk[4],
			}
		}},
		{`true ? a`, func(t *test, tk Tokens) {
			t.Err = Error{
				Err:     ErrMissingColon,
				Parsing: "ConditionalExpression",
				Token:   tk[5],
			}
		}},
		{`true ? a : :`, func(t *test, tk Tokens) {
			t.Err = Error{
				Err:     assignmentError(tk[8]),
				Parsing: "ConditionalExpression",
				Token:   tk[8],
			}
		}},
		{`++?`, func(t *test, tk Tokens) {
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
	}, func(t *test) (interface{}, error) {
		var ce ConditionalExpression
		err := ce.parse(&t.Tokens, t.In, t.Yield, t.Await)
		return ce, err
	})
}
