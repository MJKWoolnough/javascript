package javascript

import "testing"

func makeConditionLiteral(tk Tokens, pos int) ConditionalExpression {
	return wrapConditional(UpdateExpression{
		LeftHandSideExpression: &LeftHandSideExpression{
			NewExpression: &NewExpression{
				MemberExpression: MemberExpression{
					PrimaryExpression: &PrimaryExpression{
						Literal: &tk[pos],
						Tokens:  tk[pos : pos+1],
					},
					Tokens: tk[pos : pos+1],
				},
				Tokens: tk[pos : pos+1],
			},
			Tokens: tk[pos : pos+1],
		},
		Tokens: tk[pos : pos+1],
	})
}

func wrapConditional(p interface{}) ConditionalExpression {
	var c ConditionalExpression
	switch p := p.(type) {
	case LogicalORExpression:
		c.LogicalORExpression = p
		goto logicalORExpression
	case LogicalANDExpression:
		c.LogicalORExpression.LogicalANDExpression = p
		goto logicalANDExpression
	case BitwiseORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression = p
		goto bitwiseORExpression
	case BitwiseXORExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression = p
		goto bitwiseXORExpression
	case BitwiseANDExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression = p
		goto bitwiseANDExpression
	case EqualityExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression = p
		goto equalityExpression
	case RelationalExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression = p
		goto relationalExpression
	case ShiftExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression = p
		goto shiftExpression
	case AdditiveExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression = p
		goto additiveExpression
	case MultiplicativeExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression = p
		goto multiplicativeExpression
	case ExponentiationExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression = p
		goto exponentiationExpression
	case UnaryExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression = p
		goto unaryExpression
	case UpdateExpression:
		c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression = p
	default:
		panic("invalid conditional type")
	}
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.UpdateExpression.Tokens
unaryExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.UnaryExpression.Tokens
exponentiationExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.ExponentiationExpression.Tokens
multiplicativeExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.MultiplicativeExpression.Tokens
additiveExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.AdditiveExpression.Tokens
shiftExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.ShiftExpression.Tokens
relationalExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.RelationalExpression.Tokens
equalityExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.EqualityExpression.Tokens
bitwiseANDExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.BitwiseANDExpression.Tokens
bitwiseXORExpression:
	c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.BitwiseXORExpression.Tokens
bitwiseORExpression:
	c.LogicalORExpression.LogicalANDExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.BitwiseORExpression.Tokens
logicalANDExpression:
	c.LogicalORExpression.Tokens = c.LogicalORExpression.LogicalANDExpression.Tokens
logicalORExpression:
	c.Tokens = c.LogicalORExpression.Tokens
	return c
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
				LogicalORExpression:  &litA.LogicalORExpression,
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
					LogicalORExpression:  &litA.LogicalORExpression,
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
				LogicalORExpression: &litA.LogicalORExpression,
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
		/*
			{`1 ** 2`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 0)
				litB := makeConditionLiteral(tk, 4)

			}},
			{`1 ** 2 ** -3`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 0)
				litB := makeConditionLiteral(tk, 4)
				litC := makeConditionLiteral(tk, 8)

			}},
			{`delete 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`void 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`typeof 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`+ 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`- 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 3)

			}},
			{`~ 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`! 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`await 1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`await!~1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 3)

			}},
			{`1++`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 0)

			}},
			{`1--`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 0)

			}},
			{`++1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 1)

			}},
			{`++!1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 2)

			}},
			{`--1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 1)

			}},
			{`--!1`, func(t *test, tk Tokens) {
				litA := makeConditionLiteral(tk, 1)

			}},*/
	}, func(t *test) (interface{}, error) {
		return t.Tokens.parseConditionalExpression(t.In, t.Yield, t.Await)
	})
}
