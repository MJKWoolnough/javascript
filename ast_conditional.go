package javascript

import "vimagination.zapto.org/parser"

type ConditionalExpression struct {
	LogicalORExpression LogicalORExpression
	True                *AssignmentExpression
	False               *AssignmentExpression
	Tokens              Tokens
}

func (ce *ConditionalExpression) parse(j *jsParser, in, yield, await bool) error {
	g := j.NewGoal()
	if err := ce.LogicalORExpression.parse(&g, in, yield, await); err != nil {
		return j.Error("ConditionalExpression", err)
	}
	j.Score(g)
	g = j.NewGoal()
	g.AcceptRunWhitespace()
	if g.AcceptToken(parser.Token{TokenPunctuator, "?"}) {
		j.Score(g)
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		ce.True = new(AssignmentExpression)
		if err := ce.True.parse(&g, true, yield, await); err != nil {
			return j.Error("ConditionalExpression", err)
		}
		j.Score(g)
		j.AcceptRunWhitespace()
		if !j.AcceptToken(parser.Token{TokenPunctuator, ":"}) {
			return j.Error("ConditionalExpression", ErrMissingColon)
		}
		j.AcceptRunWhitespace()
		g = j.NewGoal()
		ce.False = new(AssignmentExpression)
		if err := ce.False.parse(&g, true, yield, await); err != nil {
			return j.Error("ConditionalExpression", err)
		}
		j.Score(g)
	}
	ce.Tokens = j.ToTokens()
	return nil
}

type LogicalORExpression struct {
	LogicalORExpression  *LogicalORExpression
	LogicalANDExpression LogicalANDExpression
	Tokens               Tokens
}

func (lo *LogicalORExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := lo.LogicalANDExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("LogicalORExpression", err)
		}
		j.Score(g)
		lo.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "||"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nlo := new(LogicalORExpression)
		*nlo = *lo
		*lo = LogicalORExpression{
			LogicalORExpression: nlo,
		}
		j.Score(g)
	}
}

type LogicalANDExpression struct {
	LogicalANDExpression *LogicalANDExpression
	BitwiseORExpression  BitwiseORExpression
	Tokens               Tokens
}

func (la *LogicalANDExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := la.BitwiseORExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("LogicalANDExpression", err)
		}
		j.Score(g)
		la.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "&&"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nla := new(LogicalANDExpression)
		*nla = *la
		*la = LogicalANDExpression{
			LogicalANDExpression: nla,
		}
		j.Score(g)
	}
}

type BitwiseORExpression struct {
	BitwiseORExpression  *BitwiseORExpression
	BitwiseXORExpression BitwiseXORExpression
	Tokens               Tokens
}

func (bo *BitwiseORExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := bo.BitwiseXORExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("BitwiseORExpression", err)
		}
		j.Score(g)
		bo.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "|"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nbo := new(BitwiseORExpression)
		*nbo = *bo
		*bo = BitwiseORExpression{
			BitwiseORExpression: nbo,
		}
		j.Score(g)
	}
}

type BitwiseXORExpression struct {
	BitwiseXORExpression *BitwiseXORExpression
	BitwiseANDExpression BitwiseANDExpression
	Tokens               Tokens
}

func (bx *BitwiseXORExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := bx.BitwiseANDExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("BitwiseXORExpression", err)
		}
		j.Score(g)
		bx.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "^"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nbx := new(BitwiseXORExpression)
		*nbx = *bx
		*bx = BitwiseXORExpression{
			BitwiseXORExpression: nbx,
		}
		j.Score(g)
	}
}

type BitwiseANDExpression struct {
	BitwiseANDExpression *BitwiseANDExpression
	EqualityExpression   EqualityExpression
	Tokens               Tokens
}

func (ba *BitwiseANDExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := ba.EqualityExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("BitwiseANDExpression", err)
		}
		j.Score(g)
		ba.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "&"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nba := new(BitwiseANDExpression)
		*nba = *ba
		*ba = BitwiseANDExpression{
			BitwiseANDExpression: nba,
		}
		j.Score(g)
	}
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

func (ee *EqualityExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := ee.RelationalExpression.parse(&g, in, yield, await); err != nil {
			return j.Error("EqualityExpression", err)
		}
		j.Score(g)
		ee.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var eo EqualityOperator
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
			return nil
		}
		g.Except()
		g.AcceptRunWhitespace()
		nee := new(EqualityExpression)
		*nee = *ee
		*ee = EqualityExpression{
			EqualityExpression: nee,
			EqualityOperator:   eo,
		}
		j.Score(g)
	}
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

func (re *RelationalExpression) parse(j *jsParser, in, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := re.ShiftExpression.parse(&g, yield, await); err != nil {
			return j.Error("RelationalExpression", err)
		}
		j.Score(g)
		re.Tokens = j.ToTokens()
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
				return nil
			}
			ro = RelationshipIn
		default:
			return nil
		}
		g.Except()
		g.AcceptRunWhitespace()
		nre := new(RelationalExpression)
		*nre = *re
		*re = RelationalExpression{
			RelationalExpression: nre,
			RelationshipOperator: ro,
		}
		j.Score(g)
	}
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

func (se *ShiftExpression) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := se.AdditiveExpression.parse(&g, yield, await); err != nil {
			return j.Error("ShiftExpression", err)
		}
		j.Score(g)
		se.Tokens = j.ToTokens()
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
			return nil
		}
		g.Except()
		g.AcceptRunWhitespace()
		nse := new(ShiftExpression)
		*nse = *se
		*se = ShiftExpression{
			ShiftExpression: nse,
			ShiftOperator:   so,
		}
		j.Score(g)
	}
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

func (ae *AdditiveExpression) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := ae.MultiplicativeExpression.parse(&g, yield, await); err != nil {
			return j.Error("AdditiveExpression", err)
		}
		j.Score(g)
		ae.Tokens = j.ToTokens()
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		var ao AdditiveOperator
		switch g.Peek() {
		case parser.Token{TokenPunctuator, "+"}:
			ao = AdditiveAdd
		case parser.Token{TokenPunctuator, "-"}:
			ao = AdditiveMinus
		default:
			return nil
		}
		g.Except()
		g.AcceptRunWhitespace()
		nae := new(AdditiveExpression)
		*nae = *ae
		*ae = AdditiveExpression{
			AdditiveExpression: nae,
			AdditiveOperator:   ao,
		}
		j.Score(g)
	}
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

func (me *MultiplicativeExpression) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := me.ExponentiationExpression.parse(&g, yield, await); err != nil {
			return j.Error("MultiplicativeExpression", err)
		}
		j.Score(g)
		me.Tokens = j.ToTokens()
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
			return nil
		}
		g.Except()
		g.AcceptRunWhitespace()
		nmw := new(MultiplicativeExpression)
		*nmw = *me
		*me = MultiplicativeExpression{
			MultiplicativeExpression: nmw,
			MultiplicativeOperator:   mo,
		}
		j.Score(g)
	}
}

type ExponentiationExpression struct {
	ExponentiationExpression *ExponentiationExpression
	UnaryExpression          UnaryExpression
	Tokens                   Tokens
}

func (ee *ExponentiationExpression) parse(j *jsParser, yield, await bool) error {
	for {
		g := j.NewGoal()
		if err := ee.UnaryExpression.parse(&g, yield, await); err != nil {
			return j.Error("ExponentiationExpression", err)
		}
		j.Score(g)
		ee.Tokens = j.ToTokens()
		if len(ee.UnaryExpression.UnaryOperators) > 0 {
			return nil
		}
		g = j.NewGoal()
		g.AcceptRunWhitespace()
		if !g.AcceptToken(parser.Token{TokenPunctuator, "**"}) {
			return nil
		}
		g.AcceptRunWhitespace()
		nee := new(ExponentiationExpression)
		*nee = *ee
		*ee = ExponentiationExpression{
			ExponentiationExpression: nee,
		}
		j.Score(g)
	}
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

func (ue *UnaryExpression) parse(j *jsParser, yield, await bool) error {
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
	if err := ue.UpdateExpression.parse(&g, yield, await); err != nil {
		return j.Error("UnaryExpression", err)
	}
	j.Score(g)
	ue.Tokens = j.ToTokens()
	return nil
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

func (ue *UpdateExpression) parse(j *jsParser, yield, await bool) error {
	if j.AcceptToken(parser.Token{TokenPunctuator, "++"}) || j.AcceptToken(parser.Token{TokenPunctuator, "--"}) {
		if j.GetLastToken().Data == "++" {
			ue.UpdateOperator = UpdatePreIncrement
		} else {
			ue.UpdateOperator = UpdatePreDecrement
		}
		j.AcceptRunWhitespace()
		g := j.NewGoal()
		ue.UnaryExpression = new(UnaryExpression)
		if err := ue.UnaryExpression.parse(&g, yield, await); err != nil {
			return j.Error("UpdateExpression", err)
		}
		j.Score(g)
	} else {
		g := j.NewGoal()
		ue.LeftHandSideExpression = new(LeftHandSideExpression)
		if err := ue.LeftHandSideExpression.parse(&g, yield, await); err != nil {
			return j.Error("UpdateExpression", err)
		}
		j.Score(g)
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
	return nil
}
