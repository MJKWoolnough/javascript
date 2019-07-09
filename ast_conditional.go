package javascript

import "vimagination.zapto.org/parser"

type ConditionalExpression struct {
	LogicalORExpression LogicalORExpression
	True                *AssignmentExpression
	False               *AssignmentExpression
	Tokens              Tokens
}

func (j *jsParser) parseConditionalExpression(in, yield, await bool) (ConditionalExpression, error) {
	var (
		ce  ConditionalExpression
		err error
	)
	g := j.NewGoal()
	ce.LogicalORExpression, err = g.parseLogicalORExpression(in, yield, await)
	if err != nil {
		return ce, j.Error("ConditionalExpression", err)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "?"}) {
		j.Score(g)
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		t, err := g.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return ce, j.Error("ConditionalExpression", err)
		}
		j.Score(g)
		ce.True = &t
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
			return ce, j.Error("ConditionalExpression", ErrMissingColon)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		f, err := g.parseAssignmentExpression(true, yield, await)
		if err != nil {
			return ce, j.Error("ConditionalExpression", err)
		}
		j.Score(g)
		ce.False = &f
	}
	ce.Tokens = j.ToTokens()
	return ce, nil
}

type LogicalORExpression struct {
	LogicalORExpression  *LogicalORExpression
	LogicalANDExpression LogicalANDExpression
	Tokens               Tokens
}

func (j *jsParser) parseLogicalORExpression(in, yield, await bool) (LogicalORExpression, error) {
	var (
		lo  LogicalORExpression
		err error
	)
	for {
		g := j.NewGoal()
		lo.LogicalANDExpression, err = g.parseLogicalANDExpression(in, yield, await)
		if err != nil {
			return lo, j.Error("LogicalORExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "||"}) {
			break
		}
		g.AcceptRunWhitespace()
		lo = LogicalORExpression{
			LogicalORExpression: &LogicalORExpression{
				LogicalORExpression:  lo.LogicalORExpression,
				LogicalANDExpression: lo.LogicalANDExpression,
				Tokens:               j.ToTokens(),
			},
		}
		j.Score(g)
	}
	lo.Tokens = j.ToTokens()
	return lo, nil
}

type LogicalANDExpression struct {
	LogicalANDExpression *LogicalANDExpression
	BitwiseORExpression  BitwiseORExpression
	Tokens               Tokens
}

func (j *jsParser) parseLogicalANDExpression(in, yield, await bool) (LogicalANDExpression, error) {
	var (
		la  LogicalANDExpression
		err error
	)
	for {
		g := j.NewGoal()
		la.BitwiseORExpression, err = g.parseBitwiseORExpression(in, yield, await)
		if err != nil {
			return la, j.Error("LogicalANDExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "&&"}) {
			break
		}
		g.AcceptRunWhitespace()
		la = LogicalANDExpression{
			LogicalANDExpression: &LogicalANDExpression{
				LogicalANDExpression: la.LogicalANDExpression,
				BitwiseORExpression:  la.BitwiseORExpression,
				Tokens:               j.ToTokens(),
			},
		}
		j.Score(g)
	}
	la.Tokens = j.ToTokens()
	return la, nil
}

type BitwiseORExpression struct {
	BitwiseORExpression  *BitwiseORExpression
	BitwiseXORExpression BitwiseXORExpression
	Tokens               Tokens
}

func (j *jsParser) parseBitwiseORExpression(in, yield, await bool) (BitwiseORExpression, error) {
	var (
		bo  BitwiseORExpression
		err error
	)
	for {
		g := j.NewGoal()
		bo.BitwiseXORExpression, err = g.parseBitwiseXORExpression(in, yield, await)
		if err != nil {
			return bo, j.Error("BitwiseORExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "|"}) {
			break
		}
		g.AcceptRunWhitespace()
		bo = BitwiseORExpression{
			BitwiseORExpression: &BitwiseORExpression{
				BitwiseORExpression:  bo.BitwiseORExpression,
				BitwiseXORExpression: bo.BitwiseXORExpression,
				Tokens:               j.ToTokens(),
			},
		}
		j.Score(g)
	}
	bo.Tokens = j.ToTokens()
	return bo, nil
}

type BitwiseXORExpression struct {
	BitwiseXORExpression *BitwiseXORExpression
	BitwiseANDExpression BitwiseANDExpression
	Tokens               Tokens
}

func (j *jsParser) parseBitwiseXORExpression(in, yield, await bool) (BitwiseXORExpression, error) {
	var (
		bx  BitwiseXORExpression
		err error
	)
	for {
		g := j.NewGoal()
		bx.BitwiseANDExpression, err = g.parseBitwiseANDExpression(in, yield, await)
		if err != nil {
			return bx, j.Error("BitwiseXORExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "^"}) {
			break
		}
		g.AcceptRunWhitespace()
		bx = BitwiseXORExpression{
			BitwiseXORExpression: &BitwiseXORExpression{
				BitwiseXORExpression: bx.BitwiseXORExpression,
				BitwiseANDExpression: bx.BitwiseANDExpression,
				Tokens:               j.ToTokens(),
			},
		}
		j.Score(g)
	}
	bx.Tokens = j.ToTokens()
	return bx, nil
}

type BitwiseANDExpression struct {
	BitwiseANDExpression *BitwiseANDExpression
	EqualityExpression   EqualityExpression
	Tokens               Tokens
}

func (j *jsParser) parseBitwiseANDExpression(in, yield, await bool) (BitwiseANDExpression, error) {
	var (
		ba  BitwiseANDExpression
		err error
	)
	for {
		g := j.NewGoal()
		ba.EqualityExpression, err = g.parseEqualityExpression(in, yield, await)
		if err != nil {
			return ba, j.Error("BitwiseANDExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "&"}) {
			break
		}
		g.AcceptRunWhitespace()
		ba = BitwiseANDExpression{
			BitwiseANDExpression: &BitwiseANDExpression{
				EqualityExpression:   ba.EqualityExpression,
				BitwiseANDExpression: ba.BitwiseANDExpression,
				Tokens:               j.ToTokens(),
			},
		}
		j.Score(g)
	}
	ba.Tokens = j.ToTokens()
	return ba, nil
}

type EqualityOperator int

const (
	EqualityNone EqualityOperator = iota
	EqualityEqual
	EqualityNotEqual
	EqualityStrictEqual
	EqualityStrictNotEqual
)

type EqualityExpression struct {
	EqualityExpression   *EqualityExpression
	EqualityOperator     EqualityOperator
	RelationalExpression RelationalExpression
	Tokens               Tokens
}

func (j *jsParser) parseEqualityExpression(in, yield, await bool) (EqualityExpression, error) {
	var (
		ee  EqualityExpression
		err error
		eo  EqualityOperator
	)
Loop:
	for {
		g := j.NewGoal()
		ee.RelationalExpression, err = g.parseRelationalExpression(in, yield, await)
		if err != nil {
			return ee, j.Error("EqualityExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "=="}:
			eo = EqualityEqual
		case parser.Token{TokenPunctuator, "!="}:
			eo = EqualityNotEqual
		case parser.Token{TokenPunctuator, "==="}:
			eo = EqualityStrictEqual
		case parser.Token{TokenPunctuator, "!=="}:
			eo = EqualityStrictNotEqual
		default:
			break Loop
		}
		g.Except()
		g.AcceptRunWhitespace()
		ee = EqualityExpression{
			EqualityExpression: &EqualityExpression{
				EqualityExpression:   ee.EqualityExpression,
				EqualityOperator:     ee.EqualityOperator,
				RelationalExpression: ee.RelationalExpression,
				Tokens:               j.ToTokens(),
			},
		}
		ee.EqualityOperator = eo
		j.Score(g)
	}
	ee.Tokens = j.ToTokens()
	return ee, nil
}

type RelationshipOperator int

const (
	RelationshipNone RelationshipOperator = iota
	RelationshipLessThan
	RelationshipGreaterThan
	RelationshipLessThanEqual
	RelationshipGreaterThanEqual
	RelationshipInstanceOf
	RelationshipIn
)

type RelationalExpression struct {
	RelationalExpression *RelationalExpression
	RelationshipOperator RelationshipOperator
	ShiftExpression      ShiftExpression
	Tokens               Tokens
}

func (j *jsParser) parseRelationalExpression(in, yield, await bool) (RelationalExpression, error) {
	var (
		re  RelationalExpression
		err error
	)
Loop:
	for {
		g := j.NewGoal()
		re.ShiftExpression, err = g.parseShiftExpression(yield, await)
		if err != nil {
			return re, j.Error("RelationalExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var ro RelationshipOperator
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "<"}:
			ro = RelationshipLessThan
		case parser.Token{TokenPunctuator, ">"}:
			ro = RelationshipGreaterThan
		case parser.Token{TokenPunctuator, "<="}:
			ro = RelationshipLessThanEqual
		case parser.Token{TokenPunctuator, ">="}:
			ro = RelationshipGreaterThanEqual
		case parser.Token{TokenKeyword, "instanceof"}:
			ro = RelationshipInstanceOf
		case parser.Token{TokenKeyword, "in"}:
			if !in {
				break Loop
			}
			ro = RelationshipIn
		default:
			break Loop
		}
		g.Except()
		g.AcceptRunWhitespace()
		re = RelationalExpression{
			RelationalExpression: &RelationalExpression{
				RelationalExpression: re.RelationalExpression,
				RelationshipOperator: re.RelationshipOperator,
				ShiftExpression:      re.ShiftExpression,
				Tokens:               j.ToTokens(),
			},
		}
		re.RelationshipOperator = ro
		j.Score(g)
	}
	re.Tokens = j.ToTokens()
	return re, nil
}

type ShiftOperator int

const (
	ShiftNone ShiftOperator = iota
	ShiftLeft
	ShiftRight
	ShiftUnsignedRight
)

type ShiftExpression struct {
	ShiftExpression    *ShiftExpression
	ShiftOperator      ShiftOperator
	AdditiveExpression AdditiveExpression
	Tokens             Tokens
}

func (j *jsParser) parseShiftExpression(yield, await bool) (ShiftExpression, error) {
	var (
		se  ShiftExpression
		err error
	)
Loop:
	for {
		g := j.NewGoal()
		se.AdditiveExpression, err = g.parseAdditiveExpression(yield, await)
		if err != nil {
			return se, j.Error("ShiftExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var so ShiftOperator
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "<<"}:
			so = ShiftLeft
		case parser.Token{TokenPunctuator, ">>"}:
			so = ShiftRight
		case parser.Token{TokenPunctuator, ">>>"}:
			so = ShiftUnsignedRight
		default:
			break Loop
		}
		g.Except()
		g.AcceptRunWhitespace()
		se = ShiftExpression{
			ShiftExpression: &ShiftExpression{
				ShiftExpression:    se.ShiftExpression,
				ShiftOperator:      se.ShiftOperator,
				AdditiveExpression: se.AdditiveExpression,
				Tokens:             j.ToTokens(),
			},
		}
		se.ShiftOperator = so
		j.Score(g)
	}
	se.Tokens = j.ToTokens()
	return se, nil
}

type AdditiveOperator int

const (
	AdditiveNone AdditiveOperator = iota
	AdditiveAdd
	AdditiveMinus
)

type AdditiveExpression struct {
	AdditiveExpression       *AdditiveExpression
	AdditiveOperator         AdditiveOperator
	MultiplicativeExpression MultiplicativeExpression
	Tokens                   Tokens
}

func (j *jsParser) parseAdditiveExpression(yield, await bool) (AdditiveExpression, error) {
	var (
		ae  AdditiveExpression
		err error
	)
Loop:
	for {
		g := j.NewGoal()
		ae.MultiplicativeExpression, err = g.parseMultiplicativeExpression(yield, await)
		if err != nil {
			return ae, j.Error("AdditiveExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var ao AdditiveOperator
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "+"}:
			ao = AdditiveAdd
		case parser.Token{TokenPunctuator, "-"}:
			ao = AdditiveMinus
		default:
			break Loop
		}
		g.Except()
		g.AcceptRunWhitespace()
		ae = AdditiveExpression{
			AdditiveExpression: &AdditiveExpression{
				AdditiveExpression:       ae.AdditiveExpression,
				AdditiveOperator:         ae.AdditiveOperator,
				MultiplicativeExpression: ae.MultiplicativeExpression,
				Tokens:                   j.ToTokens(),
			},
		}
		ae.AdditiveOperator = ao
		j.Score(g)
	}
	ae.Tokens = j.ToTokens()
	return ae, nil
}

type MultiplicativeOperator int

const (
	MultiplicativeNone MultiplicativeOperator = iota
	MultiplicativeMultiply
	MultiplicativeDivide
	MultiplicativeRemainder
)

type MultiplicativeExpression struct {
	MultiplicativeExpression *MultiplicativeExpression
	MultiplicativeOperator   MultiplicativeOperator
	ExponentiationExpression ExponentiationExpression
	Tokens                   Tokens
}

func (j *jsParser) parseMultiplicativeExpression(yield, await bool) (MultiplicativeExpression, error) {
	var (
		me  MultiplicativeExpression
		err error
	)
Loop:
	for {
		g := j.NewGoal()
		me.ExponentiationExpression, err = g.parseExponentiationExpression(yield, await)
		if err != nil {
			return me, j.Error("MultiplicativeExpression", err)
		}
		j.Score(g)
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var mo MultiplicativeOperator
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "*"}:
			mo = MultiplicativeMultiply
		case parser.Token{TokenDivPunctuator, "/"}:
			mo = MultiplicativeDivide
		case parser.Token{TokenPunctuator, "%"}:
			mo = MultiplicativeRemainder
		default:
			break Loop
		}
		g.Except()
		g.AcceptRunWhitespace()
		me = MultiplicativeExpression{
			MultiplicativeExpression: &MultiplicativeExpression{
				MultiplicativeExpression: me.MultiplicativeExpression,
				MultiplicativeOperator:   me.MultiplicativeOperator,
				ExponentiationExpression: me.ExponentiationExpression,
				Tokens:                   j.ToTokens(),
			},
		}
		me.MultiplicativeOperator = mo
		j.Score(g)
	}
	me.Tokens = j.ToTokens()
	return me, nil
}

type ExponentiationExpression struct {
	ExponentiationExpression *ExponentiationExpression
	UnaryExpression          UnaryExpression
	Tokens                   Tokens
}

func (j *jsParser) parseExponentiationExpression(yield, await bool) (ExponentiationExpression, error) {
	var (
		ee  ExponentiationExpression
		err error
	)
Loop:
	for {
		g := j.NewGoal()
		ee.UnaryExpression, err = g.parseUnaryExpression(yield, await)
		if err != nil {
			return ee, j.Error("ExponentiationExpression", err)
		}
		j.Score(g)
		if len(ee.UnaryExpression.UnaryOperators) > 0 {
			break
		}
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "**"}) {
			break Loop
		}
		g.AcceptRunWhitespace()
		ee = ExponentiationExpression{
			ExponentiationExpression: &ExponentiationExpression{
				ExponentiationExpression: ee.ExponentiationExpression,
				UnaryExpression:          ee.UnaryExpression,
				Tokens:                   j.ToTokens(),
			},
		}
		j.Score(g)
	}
	ee.Tokens = j.ToTokens()
	return ee, nil
}

type UnaryOperator byte

const (
	UnaryNone UnaryOperator = iota
	UnaryDelete
	UnaryVoid
	UnaryTypeOf
	UnaryAdd
	UnaryMinus
	UnaryBitwiseNot
	UnaryLogicalNot
	UnaryAwait
)

type UnaryExpression struct {
	UnaryOperators   []UnaryOperator
	UpdateExpression UpdateExpression
	Tokens           Tokens
}

func (j *jsParser) parseUnaryExpression(yield, await bool) (UnaryExpression, error) {
	var (
		ue  UnaryExpression
		err error
	)
Loop:
	for {
		switch j.Peek() {
		case parser.Token{TokenKeyword, "delete"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryDelete)
		case parser.Token{TokenKeyword, "void"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryVoid)
		case parser.Token{TokenKeyword, "typeof"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryTypeOf)
		case parser.Token{TokenPunctuator, "+"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryAdd)
		case parser.Token{TokenPunctuator, "-"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryMinus)
		case parser.Token{TokenPunctuator, "~"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryBitwiseNot)
		case parser.Token{TokenPunctuator, "!"}:
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryLogicalNot)
		case parser.Token{TokenKeyword, "await"}:
			if !await {
				break Loop
			}
			ue.UnaryOperators = append(ue.UnaryOperators, UnaryAwait)
		default:
			break Loop
		}
		j.Except()
		j.AcceptRunWhitespace()
	}
	g := j.NewGoal()
	ue.UpdateExpression, err = g.parseUpdateExpression(yield, await)
	if err != nil {
		return ue, j.Error("UnaryExpression", err)
	}
	j.Score(g)
	ue.Tokens = j.ToTokens()
	return ue, nil
}

type UpdateOperator int

const (
	UpdateNone UpdateOperator = iota
	UpdatePostIncrement
	UpdatePostDecrement
	UpdatePreIncrement
	UpdatePreDecrement
)

type UpdateExpression struct {
	LeftHandSideExpression *LeftHandSideExpression
	UpdateOperator         UpdateOperator
	UnaryExpression        *UnaryExpression
	Tokens                 Tokens
}

func (j *jsParser) parseUpdateExpression(yield, await bool) (UpdateExpression, error) {
	var ue UpdateExpression
	if j.AcceptToken(parser.Token{TokenPunctuator, "++"}) || j.AcceptToken(parser.Token{TokenPunctuator, "--"}) {
		if j.GetLastToken().Data == "++" {
			ue.UpdateOperator = UpdatePreIncrement
		} else {
			ue.UpdateOperator = UpdatePreDecrement
		}
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		une, err := g.parseUnaryExpression(yield, await)
		if err != nil {
			return ue, j.Error("UpdateExpression", err)
		}
		j.Score(g)
		ue.UnaryExpression = &une
	} else {
		g := j.NewGoal()
		lhs, err := g.parseLeftHandSideExpression(yield, await)
		if err != nil {
			return ue, j.Error("UpdateExpression", err)
		}
		j.Score(g)
		ue.LeftHandSideExpression = &lhs
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if g.AcceptToken(parser.Token{TokenPunctuator, "++"}) {
			j.Score(g)
			ue.UpdateOperator = UpdatePostIncrement
		} else if g.AcceptToken(parser.Token{TokenPunctuator, "--"}) {
			j.Score(g)
			ue.UpdateOperator = UpdatePostDecrement
		}
	}
	ue.Tokens = j.ToTokens()
	return ue, nil
}
