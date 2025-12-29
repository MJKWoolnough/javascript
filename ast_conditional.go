package javascript

import "vimagination.zapto.org/parser"

// ConditionalExpression as defined in TC39
// https://tc39.es/ecma262/#prod-ConditionalExpression
//
// One, and only one, of LogicalORExpression or CoalesceExpression must be non-nil.
//
// If True is non-nil, False must be non-nil also.
type ConditionalExpression struct {
	LogicalORExpression *LogicalORExpression
	CoalesceExpression  *CoalesceExpression
	True                *AssignmentExpression
	False               *AssignmentExpression
	Tokens              Tokens
}

func (ce *ConditionalExpression) parse(j *jsParser, in, yield, await bool) error {
	g := j.NewGoal()

	ce.LogicalORExpression = new(LogicalORExpression)
	if err := ce.LogicalORExpression.parse(&g, in, yield, await); err != nil {
		return j.Error("ConditionalExpression", err)
	}

	j.Score(g)

	g = j.NewGoal()

	g.AcceptRunWhitespace()

	if ce.LogicalORExpression.LogicalORExpression == nil && ce.LogicalORExpression.LogicalANDExpression.LogicalANDExpression == nil && g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "??"}) {
		ce.CoalesceExpression = new(CoalesceExpression)
		if err := ce.CoalesceExpression.parse(j, in, yield, await, ce.LogicalORExpression.LogicalANDExpression.BitwiseORExpression); err != nil {
			return j.Error("ConditionalExpression", err)
		}

		ce.LogicalORExpression = nil
		g = j.NewGoal()

		g.AcceptRunWhitespace()
	}

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "?"}) {
		j.Score(g)
		j.AcceptRunWhitespaceNoComment()

		g = j.NewGoal()

		ce.True = new(AssignmentExpression)
		if err := ce.True.parse(&g, true, yield, await); err != nil {
			return j.Error("ConditionalExpression", err)
		}

		j.Score(g)
		j.AcceptRunWhitespace()

		if !j.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ":"}) {
			return j.Error("ConditionalExpression", ErrMissingColon)
		}

		j.AcceptRunWhitespaceNoComment()

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

func (ce *ConditionalExpression) hasFirstComment() bool {
	if ce.LogicalORExpression != nil {
		return ce.LogicalORExpression.hasFirstComment()
	}

	if ce.CoalesceExpression != nil {
		return ce.CoalesceExpression.hasFirstComment()
	}

	return false
}

func (ce *ConditionalExpression) hasLastComment() bool {
	if ce.False != nil {
		return ce.False.hasLastComment()
	}

	if ce.LogicalORExpression != nil {
		return ce.LogicalORExpression.hasLastComment()
	}

	if ce.CoalesceExpression != nil {
		return ce.CoalesceExpression.hasLastComment()
	}

	return false
}

// CoalesceExpression as defined in TC39
// https://tc39.es/ecma262/#prod-CoalesceExpression
type CoalesceExpression struct {
	CoalesceExpressionHead *CoalesceExpression
	BitwiseORExpression    BitwiseORExpression
	Tokens                 Tokens
}

func (ce *CoalesceExpression) parse(j *jsParser, in, yield, await bool, be BitwiseORExpression) error {
	ce.BitwiseORExpression = be

	for {
		ce.Tokens = j.ToTokens()
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "??"}) {
			break
		}

		g.AcceptRunWhitespaceNoComment()

		nce := new(CoalesceExpression)
		*nce = *ce
		h := g.NewGoal()

		*ce = CoalesceExpression{CoalesceExpressionHead: nce}
		if err := ce.BitwiseORExpression.parse(&h, in, yield, await); err != nil {
			return g.Error("CoalesceExpression", err)
		}

		g.Score(h)
		j.Score(g)
	}

	return nil
}

func (ce *CoalesceExpression) hasFirstComment() bool {
	if ce.CoalesceExpressionHead != nil {
		return ce.CoalesceExpressionHead.hasFirstComment()
	}

	return ce.BitwiseORExpression.hasFirstComment()
}

func (ce *CoalesceExpression) hasLastComment() bool {
	return ce.BitwiseORExpression.hasLastComment()
}

// LogicalORExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-LogicalORExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "||"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nlo := new(LogicalORExpression)
		*nlo = *lo
		*lo = LogicalORExpression{LogicalORExpression: nlo}

		j.Score(g)
	}
}

func (lo *LogicalORExpression) hasFirstComment() bool {
	if lo.LogicalORExpression != nil {
		return lo.LogicalORExpression.hasFirstComment()
	}

	return lo.LogicalANDExpression.hasFirstComment()
}

func (lo *LogicalORExpression) hasLastComment() bool {
	return lo.LogicalANDExpression.hasLastComment()
}

// LogicalANDExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-LogicalANDExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&&"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nla := new(LogicalANDExpression)
		*nla = *la
		*la = LogicalANDExpression{LogicalANDExpression: nla}

		j.Score(g)
	}
}

func (la *LogicalANDExpression) hasFirstComment() bool {
	if la.LogicalANDExpression != nil {
		return la.hasFirstComment()
	}

	return la.BitwiseORExpression.hasFirstComment()
}

func (la *LogicalANDExpression) hasLastComment() bool {
	return la.BitwiseORExpression.hasLastComment()
}

// BitwiseORExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-BitwiseORExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "|"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nbo := new(BitwiseORExpression)
		*nbo = *bo
		*bo = BitwiseORExpression{BitwiseORExpression: nbo}

		j.Score(g)
	}
}

func (bo *BitwiseORExpression) hasFirstComment() bool {
	if bo.BitwiseORExpression != nil {
		return bo.hasFirstComment()
	}

	return bo.BitwiseXORExpression.hasFirstComment()
}

func (bo *BitwiseORExpression) hasLastComment() bool {
	return bo.BitwiseXORExpression.hasLastComment()
}

// BitwiseXORExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-BitwiseXORExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "^"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nbx := new(BitwiseXORExpression)
		*nbx = *bx
		*bx = BitwiseXORExpression{BitwiseXORExpression: nbx}

		j.Score(g)
	}
}

func (bx *BitwiseXORExpression) hasFirstComment() bool {
	if bx.BitwiseXORExpression != nil {
		return bx.BitwiseXORExpression.hasFirstComment()
	}

	return bx.BitwiseANDExpression.hasFirstComment()
}

func (bx *BitwiseXORExpression) hasLastComment() bool {
	return bx.BitwiseANDExpression.hasLastComment()
}

// BitwiseANDExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-BitwiseANDExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "&"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nba := new(BitwiseANDExpression)
		*nba = *ba
		*ba = BitwiseANDExpression{BitwiseANDExpression: nba}

		j.Score(g)
	}
}

func (ba *BitwiseANDExpression) hasFirstComment() bool {
	if ba.BitwiseANDExpression != nil {
		return ba.BitwiseANDExpression.hasFirstComment()
	}

	return ba.EqualityExpression.hasFirstComment()
}

func (ba *BitwiseANDExpression) hasLastComment() bool {
	return ba.EqualityExpression.hasLastComment()
}

// EqualityOperator determines the type of EqualityExpression
type EqualityOperator int

// Valid EqualityOperator's
const (
	EqualityNone EqualityOperator = iota
	EqualityEqual
	EqualityNotEqual
	EqualityStrictEqual
	EqualityStrictNotEqual
)

// EqualityExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-EqualityExpression
//
// If EqualityOperator is not EqualityNone, then EqualityExpression must be
// non-nil, and vice-versa.
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
		case parser.Token{Type: TokenPunctuator, Data: "=="}:
			eo = EqualityEqual
		case parser.Token{Type: TokenPunctuator, Data: "!="}:
			eo = EqualityNotEqual
		case parser.Token{Type: TokenPunctuator, Data: "==="}:
			eo = EqualityStrictEqual
		case parser.Token{Type: TokenPunctuator, Data: "!=="}:
			eo = EqualityStrictNotEqual
		default:
			return nil
		}

		g.Skip()
		g.AcceptRunWhitespaceNoComment()

		nee := new(EqualityExpression)
		*nee = *ee
		*ee = EqualityExpression{EqualityExpression: nee, EqualityOperator: eo}

		j.Score(g)
	}
}

func (ee *EqualityExpression) hasFirstComment() bool {
	if ee.EqualityExpression != nil {
		return ee.EqualityExpression.hasFirstComment()
	}

	return ee.RelationalExpression.hasFirstComment()
}

func (ee *EqualityExpression) hasLastComment() bool {
	return ee.RelationalExpression.hasLastComment()
}

// RelationshipOperator determines the relationship type for RelationalExpression
type RelationshipOperator int

// Valid RelationshipOperator's
const (
	RelationshipNone RelationshipOperator = iota
	RelationshipLessThan
	RelationshipGreaterThan
	RelationshipLessThanEqual
	RelationshipGreaterThanEqual
	RelationshipInstanceOf
	RelationshipIn
)

// RelationalExpression as defined in ECMA-262
// https://tc39.es/ecma262/#prod-RelationalExpression
//
// If PrivateIdentifier is non-nil, then RelationshipOperator should be
// RelationshipIn.
//
// If PrivateIdentifier is nil and RelationshipOperator does not equal
// RelationshipNone, then RelationalExpression should be non-nil
type RelationalExpression struct {
	PrivateIdentifier    *Token
	RelationalExpression *RelationalExpression
	RelationshipOperator RelationshipOperator
	ShiftExpression      ShiftExpression
	Comments             [2]Comments
	Tokens               Tokens
}

func (re *RelationalExpression) parse(j *jsParser, in, yield, await bool) error {
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if in && g.Accept(TokenPrivateIdentifier) {
		g.AcceptRunWhitespace()

		if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
			re.Comments[0] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()

			re.PrivateIdentifier = j.Next()
			re.Comments[1] = j.AcceptRunWhitespaceComments()

			j.AcceptRunWhitespace()
			j.Skip()
			j.AcceptRunWhitespaceNoComment()

			g = j.NewGoal()

			re.RelationshipOperator = RelationshipIn
			if err := re.ShiftExpression.parse(&g, yield, await); err != nil {
				return j.Error("RelationalExpression", err)
			}

			j.Score(g)

			re.Tokens = j.ToTokens()

			return nil
		}
	}

	g = j.NewGoal()

	for {
		if err := re.ShiftExpression.parse(&g, yield, await); err != nil {
			return j.Error("RelationalExpression", err)
		}

		j.Score(g)

		re.Tokens = j.ToTokens()
		g = j.NewGoal()

		g.AcceptRunWhitespace()

		var ro RelationshipOperator

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<"}) {
			ro = RelationshipLessThan
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "="}) {
				ro = RelationshipGreaterThanEqual
			} else if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: ">"}) {
				return nil
			} else {
				ro = RelationshipGreaterThan
			}
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<="}) {
			ro = RelationshipLessThanEqual
		} else if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "instanceof"}) {
			ro = RelationshipInstanceOf
		} else if g.AcceptToken(parser.Token{Type: TokenKeyword, Data: "in"}) {
			if !in {
				return nil
			}

			ro = RelationshipIn
		} else {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nre := new(RelationalExpression)
		*nre = *re
		*re = RelationalExpression{RelationalExpression: nre, RelationshipOperator: ro}

		j.Score(g)

		g = j.NewGoal()
	}
}

func (re *RelationalExpression) hasFirstComment() bool {
	if re.RelationalExpression != nil {
		return re.RelationalExpression.hasFirstComment()
	}

	return re.ShiftExpression.hasFirstComment()
}

func (re *RelationalExpression) hasLastComment() bool {
	return re.ShiftExpression.hasLastComment()
}

// ShiftOperator determines the shift tyoe for ShiftExpression
type ShiftOperator int

// Valid ShiftOperator's
const (
	ShiftNone ShiftOperator = iota
	ShiftLeft
	ShiftRight
	ShiftUnsignedRight
)

// ShiftExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ShiftExpression
//
// If ShiftOperator is not ShiftNone then ShiftExpression must be non-nil, and
// vice-versa.
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

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "<<"}) {
			so = ShiftLeft
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) && g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
			if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: ">"}) {
				so = ShiftUnsignedRight
			} else {
				so = ShiftRight
			}

			if g.Peek() == (parser.Token{Type: TokenPunctuator, Data: "="}) {
				return nil
			}
		} else {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nse := new(ShiftExpression)
		*nse = *se
		*se = ShiftExpression{ShiftExpression: nse, ShiftOperator: so}

		j.Score(g)
	}
}

func (se *ShiftExpression) hasFirstComment() bool {
	if se.ShiftExpression != nil {
		return se.ShiftExpression.hasFirstComment()
	}

	return se.AdditiveExpression.hasFirstComment()
}

func (se *ShiftExpression) hasLastComment() bool {
	return se.AdditiveExpression.hasLastComment()
}

// AdditiveOperator determines the additive type for AdditiveExpression
type AdditiveOperator int

// Valid AdditiveOperator's
const (
	AdditiveNone AdditiveOperator = iota
	AdditiveAdd
	AdditiveMinus
)

// AdditiveExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-AdditiveExpression
//
// If AdditiveOperator is not AdditiveNone then AdditiveExpression must be
// non-nil, and vice-versa.
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
		case parser.Token{Type: TokenPunctuator, Data: "+"}:
			ao = AdditiveAdd
		case parser.Token{Type: TokenPunctuator, Data: "-"}:
			ao = AdditiveMinus
		default:
			return nil
		}

		g.Skip()
		g.AcceptRunWhitespaceNoComment()

		nae := new(AdditiveExpression)
		*nae = *ae
		*ae = AdditiveExpression{AdditiveExpression: nae, AdditiveOperator: ao}

		j.Score(g)
	}
}

func (ae *AdditiveExpression) hasFirstComment() bool {
	if ae.AdditiveExpression != nil {
		return ae.AdditiveExpression.hasFirstComment()
	}

	return ae.MultiplicativeExpression.hasFirstComment()
}

func (ae *AdditiveExpression) hasLastComment() bool {
	return ae.MultiplicativeExpression.hasLastComment()
}

// MultiplicativeOperator determines the multiplication type for MultiplicativeExpression
type MultiplicativeOperator int

// Valid MultiplicativeOperator's
const (
	MultiplicativeNone MultiplicativeOperator = iota
	MultiplicativeMultiply
	MultiplicativeDivide
	MultiplicativeRemainder
)

// MultiplicativeExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-MultiplicativeExpression
//
// If MultiplicativeOperator is not MultiplicativeNone then
// MultiplicativeExpression must be non-nil, and vice-versa.
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
		case parser.Token{Type: TokenPunctuator, Data: "*"}:
			mo = MultiplicativeMultiply
		case parser.Token{Type: TokenDivPunctuator, Data: "/"}:
			mo = MultiplicativeDivide
		case parser.Token{Type: TokenPunctuator, Data: "%"}:
			mo = MultiplicativeRemainder
		default:
			return nil
		}

		g.Skip()
		g.AcceptRunWhitespaceNoComment()

		nmw := new(MultiplicativeExpression)
		*nmw = *me
		*me = MultiplicativeExpression{MultiplicativeExpression: nmw, MultiplicativeOperator: mo}

		j.Score(g)
	}
}

func (me *MultiplicativeExpression) hasFirstComment() bool {
	if me.MultiplicativeExpression != nil {
		return me.MultiplicativeExpression.hasFirstComment()
	}

	return me.ExponentiationExpression.hasFirstComment()
}

func (me *MultiplicativeExpression) hasLastComment() bool {
	return me.ExponentiationExpression.hasLastComment()
}

// ExponentiationExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-ExponentiationExpression
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

		if !g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "**"}) {
			return nil
		}

		g.AcceptRunWhitespaceNoComment()

		nee := new(ExponentiationExpression)
		*nee = *ee
		*ee = ExponentiationExpression{ExponentiationExpression: nee}

		j.Score(g)
	}
}

func (ee *ExponentiationExpression) hasFirstComment() bool {
	if ee.ExponentiationExpression != nil {
		return ee.ExponentiationExpression.hasFirstComment()
	}

	return ee.UnaryExpression.hasFirstComment()
}

func (ee *ExponentiationExpression) hasLastComment() bool {
	return ee.UnaryExpression.hasLastComment()
}

// UnaryOperator determines a unary operator within UnaryExpression
type UnaryOperator byte

// Valid UnaryOperator's
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

// UnaryOperatorComments is a UnaryOperator and any following comments.
type UnaryOperatorComments struct {
	UnaryOperator
	Comments
}

// UnaryExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-UnaryExpression
type UnaryExpression struct {
	UnaryOperators   []UnaryOperatorComments
	UpdateExpression UpdateExpression
	Tokens           Tokens
}

func (ue *UnaryExpression) parse(j *jsParser, yield, await bool) error {
Loop:
	for {
		g := j.NewGoal()

		g.AcceptRunWhitespace()

		var op UnaryOperator

		switch g.Peek() {
		case parser.Token{Type: TokenKeyword, Data: "delete"}:
			op = UnaryDelete
		case parser.Token{Type: TokenKeyword, Data: "void"}:
			op = UnaryVoid
		case parser.Token{Type: TokenKeyword, Data: "typeof"}:
			op = UnaryTypeOf
		case parser.Token{Type: TokenPunctuator, Data: "+"}:
			op = UnaryAdd
		case parser.Token{Type: TokenPunctuator, Data: "-"}:
			op = UnaryMinus
		case parser.Token{Type: TokenPunctuator, Data: "~"}:
			op = UnaryBitwiseNot
		case parser.Token{Type: TokenPunctuator, Data: "!"}:
			op = UnaryLogicalNot
		case parser.Token{Type: TokenKeyword, Data: "await"}:
			if !await {
				break Loop
			}

			op = UnaryAwait
		default:
			break Loop
		}

		ue.UnaryOperators = append(ue.UnaryOperators, UnaryOperatorComments{
			UnaryOperator: op,
			Comments:      j.AcceptRunWhitespaceComments(),
		})

		j.AcceptRunWhitespace()
		j.Skip()
		j.AcceptRunWhitespaceNoComment()
	}

	g := j.NewGoal()

	if err := ue.UpdateExpression.parse(&g, yield, await); err != nil {
		return j.Error("UnaryExpression", err)
	}

	j.Score(g)

	ue.Tokens = j.ToTokens()

	return nil
}

func (ue *UnaryExpression) hasFirstComment() bool {
	return ue.UpdateExpression.hasFirstComment()
}

func (ue *UnaryExpression) hasLastComment() bool {
	return ue.UpdateExpression.hasLastComment()
}

// UpdateOperator determines the type of update operation for UpdateExpression
type UpdateOperator int

// Valid UpdateOperator's
const (
	UpdateNone UpdateOperator = iota
	UpdatePostIncrement
	UpdatePostDecrement
	UpdatePreIncrement
	UpdatePreDecrement
)

// UpdateExpression as defined in ECMA-262
// https://262.ecma-international.org/11.0/#prod-UpdateExpression
//
// If UpdateOperator is UpdatePreIncrement or UpdatePreDecrement
// UnaryExpression must be non-nil, and vice-versa. In all other cases,
// LeftHandSideExpression must be non-nil.
type UpdateExpression struct {
	LeftHandSideExpression *LeftHandSideExpression
	UpdateOperator         UpdateOperator
	UnaryExpression        *UnaryExpression
	Comments               Comments
	Tokens                 Tokens
}

func (ue *UpdateExpression) parse(j *jsParser, yield, await bool) error {
	g := j.NewGoal()

	g.AcceptRunWhitespace()

	if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "++"}) || g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "--"}) {
		ue.Comments = j.AcceptRunWhitespaceComments()

		j.AcceptRunWhitespace()

		if j.Next().Data == "++" {
			ue.UpdateOperator = UpdatePreIncrement
		} else {
			ue.UpdateOperator = UpdatePreDecrement
		}

		j.AcceptRunWhitespaceNoComment()

		g := j.NewGoal()

		ue.UnaryExpression = new(UnaryExpression)
		if err := ue.UnaryExpression.parse(&g, yield, await); err != nil {
			return j.Error("UpdateExpression", err)
		}

		j.Score(g)
	} else {
		g = j.NewGoal()

		ue.LeftHandSideExpression = new(LeftHandSideExpression)
		if err := ue.LeftHandSideExpression.parse(&g, yield, await); err != nil {
			return j.Error("UpdateExpression", err)
		}

		j.Score(g)

		g = j.NewGoal()

		g.AcceptRunWhitespaceNoNewLineNoComment()

		if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "++"}) {
			if !ue.LeftHandSideExpression.IsSimple() {
				return j.Error("UpdateExpression", ErrNotSimple)
			}

			j.Score(g)

			ue.UpdateOperator = UpdatePostIncrement
		} else if g.AcceptToken(parser.Token{Type: TokenPunctuator, Data: "--"}) {
			if !ue.LeftHandSideExpression.IsSimple() {
				return j.Error("UpdateExpression", ErrNotSimple)
			}

			j.Score(g)

			ue.UpdateOperator = UpdatePostDecrement
		}

		ue.Comments = j.AcceptRunWhitespaceCommentsInList()
	}

	ue.Tokens = j.ToTokens()

	return nil
}

func (ue *UpdateExpression) hasFirstComment() bool {
	if ue.LeftHandSideExpression != nil {
		return ue.LeftHandSideExpression.hasFirstComment()
	}

	if ue.UpdateOperator == UpdatePreIncrement || ue.UpdateOperator == UpdatePreDecrement {
		return len(ue.Comments) > 0
	}

	return ue.UnaryExpression.hasFirstComment()
}

func (ue *UpdateExpression) hasLastComment() bool {
	if ue.UpdateOperator == UpdatePostIncrement || ue.UpdateOperator == UpdatePostDecrement {
		return len(ue.Comments) > 0
	}

	return ue.LeftHandSideExpression != nil && ue.LeftHandSideExpression.hasLastComment()
}
